package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"
)

const tickInterval int64 = 300

func mainCore() error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPass, cfg.DBName)

	client, err := initClient(psqlInfo)
	defer client.close()

	if err != nil {
		log.Printf("Error: %v", err)
		return err
	}

	data := make([]exchangeDataTick, 0)
	if exists, _ := tableExists(client.db, "exchangedata"); exists {
		if d := collectExchangeData(time.Now().Unix()); d != nil {
			data = d
		} else {
			log.Print("Could not retrieve exchange data")
		}
	} else {
		if err := client.createExchangetable(); err != nil {
			log.Printf("Error: %v", err)
			return err
		}
		if d := collectExchangeData(0); d != nil {
			data = d
		} else {
			log.Print("Could not retrieve exchange data")
		}
	}

	err = client.addEntries(data)
	if err != nil {
		log.Printf("Error: %v", err)
		return err
	}
	log.Print("All entries successfully stored")

	quit := make(chan struct{})
	// Only accept a single CTRL+C
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Start waiting for the interrupt signal
	go func() {
		<-c
		signal.Stop(c)
		// Close the channel so multiple goroutines can get the message
		log.Print("CTRL+C hit.  Closing goroutines.")
		close(quit)
	}()

	var wg sync.WaitGroup

	wg.Add(1)

	go collectAtInterval(client, tickInterval, &wg, quit)

	wg.Wait()
	return nil
}
func main() {
	if err := mainCore(); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}

func collectAtInterval(client *pgClient, interval int64, wg *sync.WaitGroup, quit chan struct{}) {
	ticker := time.NewTicker(time.Second * time.Duration(interval))

	defer func() {
		ticker.Stop()
		wg.Done()
	}()
	last := time.Now().Unix()
	log.Printf("Starting collector")
	for {
		select {
		case t := <-ticker.C:
			log.Print("Collecting exchange data")
			data := collectExchangeData(last)
			last = t.Unix()
			if data == nil {
				log.Print("Could not retrieve exchange data")
				return
			}
			err := client.addEntries(data)
			if err != nil {
				log.Printf("Error: %v", err)
				return
			}

		case <-quit:
			log.Printf("Closing collector")
			return
		}
	}
}
