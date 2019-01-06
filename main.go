package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

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

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPass, cfg.DBName)

	db, err := NewPgDb(psqlInfo)
	defer db.Close()

	err = db.Ping()
	if err != nil {
		db.Close()
		log.Fatal("Error connecting to Postgresl: ", err)
	}

	if cfg.DropTables {
		log.Info("Dropping tables")
		err = db.DropExchangeDataTable()
		if err != nil {
			db.Close()
			log.Fatal("Could not drop tables: ", err)
		} else {
			log.Info("Tables dropped")
			return err
		}

	}

	log.Info("Attemping to retrieve exchange data")
	data := make([]exchangeDataTick, 0)
	if exists := db.ExchangeDataTableExits(); exists {
		t, err := db.LastExchangeEntryTime()
		if err != nil {
			if strings.Contains(err.Error(), "no rows") {
				t = 0
			} else {
				log.Error("Could not retrieve last entry time: ", err)
				return err
			}
		}
		log.Info("Retireving exchange data from ", time.Unix(t, 0).String())
		if d, err := collectExchangeData(t); err == nil {
			data = d
		} else {
			log.Error("Could not retrieve exchange data: ", err)
			return err
		}
	} else {
		log.Info("Creating new exchange data table")
		if err := db.CreateExchangeDataTable(); err != nil {
			log.Error("Error creating exchange data table: ", err)
			return err
		}
		log.Info("Retrieving exchange data")
		if d, err := collectExchangeData(0); err == nil {
			data = d
		} else {
			log.Error("Could not retrieve exchange data: ", err)
			return err
		}
	}

	log.Debug("Collected exchange entry count: ", len(data))

	err = db.AddExchangeData(data)
	if err != nil {
		log.Error("Error adding exchange entries: ", err)
		return err
	}
	log.Info("Collected entries stored")

	quit := make(chan struct{})

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		<-c
		signal.Stop(c)

		log.Info("CTRL+C hit.  Closing goroutines.")
		close(quit)
	}()

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		last, err := db.LastExchangeEntryTime()

		if err != nil {
			log.Error("Could not retrieve last entry time ", err)
			wg.Done()
			return
		}

		// Sleep till 30 seconds before next collection time
		time.Sleep(time.Duration(last+1730-time.Now().Unix()) * time.Second)

		ticker := time.NewTicker(time.Second * time.Duration(1800)) // Set a timer for every 30 minutes

		defer func() {
			ticker.Stop()
			wg.Done()
		}()

		log.Info("Starting collector")
		for {
			select {
			case t := <-ticker.C:
				log.Info("Collecting recent exchange data")
				data, err := collectExchangeData(last)
				last = t.Unix()
				if err != nil {
					log.Error("Could not retrieve exchange data: ", err)
					return
				}
				err = db.AddExchangeData(data)
				if err != nil {
					log.Error("Error adding exchange data entries: ", err)
					return
				}
				log.Info("Added recent exchange data")
			case <-quit:
				log.Info("Closing collector")
				return
			}
		}
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
