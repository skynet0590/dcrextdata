package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	cfg, err := loadConfig()
	if err != nil {
		return
	}
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPass, cfg.DBName)

	client, err := initClient(psqlInfo)
	defer client.close()

	if err != nil {
		log.Printf("Error: %v", err)
		return
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
			return
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
		return
	}
	log.Print("All entries successfully stored")
}

func collectExchangeData(start int64) []exchangeDataTick {
	data := make([]exchangeDataTick, 0)

	poloniexdata, err := collectPoloniexData(start)
	if err != nil {
		log.Printf("Error: %v", err)
		return nil
	}
	bittrexdata, err := collectBittrexData(start)
	if err != nil {
		log.Printf("Error: %v", err)
		return nil
	}
	data = append(data, poloniexdata...)
	data = append(data, bittrexdata...)
	return data
}
