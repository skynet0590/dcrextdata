package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"
)

const tickInterval int64 = 3600

func mainCore() error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPass, cfg.DBName)

	db, err := NewPgDb(psqlInfo)
	defer db.Close()

	if err != nil {
		log.Printf("Error: %v", err)
		return err
	}

	if cfg.DropTables {
		log.Print("Dropping tables")
		err = db.DropExchangeDataTable()
		if err != nil {
			log.Printf("Could not drop tables: %v", err)
		} else {
			log.Print("Tables dropped")
		}
		return err
	}

	data := make([]exchangeDataTick, 0)
	if exists := db.ExchangeDataTableExits(); exists {
		t, err := db.LastExchangeEntryTime()
		if err != nil {
			if strings.Contains(err.Error(), "no rows") {
				t = 0
			} else {
				log.Printf("Could not retrieve last entry time: %v", err)
				return err
			}
		}
		log.Printf("Retireving exchange data from %s", time.Unix(t, 0).String())
		if d, err := collectExchangeData(t); err == nil {
			data = d
		} else {
			log.Print("Could not retrieve exchange data")
			return err
		}
	} else {
		log.Printf("Creating new exchange data table")
		if err := db.CreateExchangeDataTable(); err != nil {
			log.Printf("Error: %v", err)
			return err
		}
		log.Print("Retrieving exchange data")
		if d, err := collectExchangeData(0); err == nil {
			data = d
		} else {
			log.Print("Could not retrieve exchange data")
			return err
		}
	}

	log.Print("Attempting to store entries...")
	err = db.AddExchangeData(data)
	if err != nil {
		log.Printf("Error: %v", err)
		return err
	}
	log.Print("All entries successfully stored")

	quit := make(chan struct{})

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		<-c
		signal.Stop(c)

		log.Print("CTRL+C hit.  Closing goroutines.")
		close(quit)
	}()

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		ticker := time.NewTicker(time.Second * time.Duration(1860))

		defer func() {
			ticker.Stop()
			wg.Done()
		}()

		last := time.Now().Unix()
		log.Printf("Starting collector")
		for {
			select {
			case t := <-ticker.C:
				log.Print("Collecting recent exchange data")
				data, err := collectExchangeData(last)
				last = t.Unix()
				if err != nil {
					log.Print("Could not retrieve exchange data")
					return
				}
				err = db.AddExchangeData(data)
				if err != nil {
					log.Printf("Error: %v", err)
					return
				}
				log.Print("Added recent exchange data")
			case <-quit:
				log.Printf("Closing collector")
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
