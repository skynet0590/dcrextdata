package main

import (
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/raedahgroup/dcrextdata/collection/exchanges"
	"github.com/raedahgroup/dcrextdata/db"
	log "github.com/sirupsen/logrus"
)

const dcrlaunchtime int64 = 1454889600

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:          true,
		DisableLevelTruncation: true,
		TimestampFormat:        "2006-01-02 15:04:05",
	})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func mainCore() error {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatal("Unable to load config: ", err)
	}

	if cfg.Quiet {
		log.SetLevel(log.ErrorLevel)
	}

	db, err := db.NewPgDb(cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPass, cfg.DBName)
	defer db.Close()

	if cfg.Reset {
		log.Info("Dropping tables")
		err = db.DropAllTables()
		if err != nil {
			db.Close()
			log.Fatal("Could not drop tables: ", err)
		} else {
			log.Info("Tables dropped")
			// return err
		}
	}

	log.Info("Attemping to retrieve exchange data")

	client := &http.Client{Timeout: 300 * time.Second}

	if exists := db.ExchangeDataTableExits(); !exists {
		log.Info("Creating new exchange data table")
		if err := db.CreateExchangeDataTable(); err != nil {
			log.Error("Error creating exchange data table: ", err)
			return err
		}
	}

	quit := make(chan struct{})
	wg := new(sync.WaitGroup)

	//retrievers := make([]exchanges.Retriever, 0, 2)

	poloniex := exchanges.NewPoloniex(client, db.LastExchangeEntryTime("poloniex")-cfg.CollectionInterval, cfg.CollectionInterval)
	bittrex := exchanges.NewBittrex(client, db.LastExchangeEntryTime("bittex")-cfg.CollectionInterval, cfg.CollectionInterval)

	exchangeCollector := exchanges.Collector{Retrievers: []exchanges.Retriever{poloniex, bittrex}}

	log.Info("Starting Collector")
	resultChan, errChan := exchangeCollector.CollectAtInterval(time.Duration(cfg.CollectionInterval)*time.Second, wg, quit)

	wg.Add(1)
	go storeExchangeData(db, resultChan, errChan, quit)

	//last := db.LastExchangeEntryTime()
	// Sleep till 30 seconds before next collection time
	//time.Sleep(time.Duration(last+1730-time.Now().Unix()) * time.Second)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		<-c
		signal.Stop(c)

		log.Info("CTRL+C hit.  Closing goroutines.")
		close(quit)
	}()

	wg.Wait()
	return nil
}

func main() {
	if err := mainCore(); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}

func storeExchangeData(db *db.PgDb, resultChan chan []exchanges.DataTick, errChan chan error, quit chan struct{}) {
	for {
		select {
		case dataTick := <-resultChan:
			added, err := db.AddExchangeData(dataTick)
			if err != nil {
				log.Errorf("Could not store exchange entry: %v", err)
			}
			log.Infof("Added %d entries", added)
		case err := <-errChan:
			log.Error(err)
		case <-quit:
			log.Info("Quitting storage goroutine")
			return
		}
	}
}
