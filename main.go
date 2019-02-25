package main

import (
	"os"
	"os/signal"
	"strings"
	"sync"
)

const dcrlaunchtime int64 = 1454889600

func init() {
	// log.SetFormatter(&log.TextFormatter{
	// 	FullTimestamp:          true,
	// 	DisableLevelTruncation: true,
	// 	TimestampFormat:        "2006-01-02 15:04:05",
	// })
	// log.SetOutput(os.Stdout)
	// log.SetLevel(log.DebugLevel)

	//initLogRotator("logs/main.log")
}

func main() {
	cfg, err := loadConfig()
	if err != nil {
		log.Error("Unable to load config: ", err)
		return
	}

	defer func() {
		if logRotator != nil {
			logRotator.Close()
		}
	}()

	db, err := NewPgDb(cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPass, cfg.DBName)
	defer db.Close()

	if cfg.Reset {
		log.Info("Dropping tables")
		err = db.DropAllTables()
		if err != nil {
			db.Close()
			log.Error("Could not drop tables: ", err)
			return
		} else {
			log.Info("Tables dropped")
			// return err
		}
	}

	if exists := db.ExchangeDataTableExits(); !exists {
		log.Info("Creating new exchange data table")
		if err := db.CreateExchangeDataTable(); err != nil {
			log.Error("Error creating exchange data table: ", err)
			return
		}
	}

	//retrievers := make([]exchanges.Retriever, 0, 2)

	//exchangeCollector := exchanges.Collector{Retrievers: []exchanges.Retriever{poloniex, bittrex}}

	resultChan := make(chan []DataTick)

	quit := make(chan struct{})
	wg := new(sync.WaitGroup)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		<-c
		signal.Stop(c)
		log.Info("CTRL+C hit. Closing goroutines.")
		close(quit)
	}()

	wg.Add(1)
	log.Info("Starting exchange storage goroutine")
	go storeExchangeData(db, resultChan, quit, wg)

	// Exchange collection enabled
	if len(cfg.Exchanges) > 0 {
		exchanges := make(map[string]int64)
		for _, ex := range strings.Split(cfg.Exchanges, ",") {
			exchanges[ex] = db.LastExchangeEntryTime(ex)
		}

		//log.Debugf("exchangeMap: %v", exchanges)
		collector, err := NewExchangeCollector(exchanges, cfg.CollectionInterval)

		if err != nil {
			log.Error(err)
			return
		}

		if collector.HistoricSyncRequired() {
			log.Info("Starting historic sync")
			if collector.HistoricSync(resultChan) {
				excLog.Error("Historic sync failed")
				close(quit)
				return
			}
			excLog.Info("Completed historic sync")
		}

		go collector.Collect(resultChan, quit)
	}

	//last := db.LastExchangeEntryTime()
	// Sleep till 30 seconds before next collection time
	//time.Sleep(time.Duration(last+1730-time.Now().Unix()) * time.Second)

	wg.Wait()
	log.Info("Goodbye")
}

func storeExchangeData(db *PgDb, resultChan chan []DataTick, quit chan struct{}, wg *sync.WaitGroup) {
	for {
		select {
		case dataTick := <-resultChan:
			err := db.AddExchangeData(dataTick)
			if err != nil {
				log.Errorf("Could not store exchange entry: %v", err)
			}
		case <-quit:
			log.Debug("Retrieved quit signal")
			wg.Done()
			return
		}
	}
}
