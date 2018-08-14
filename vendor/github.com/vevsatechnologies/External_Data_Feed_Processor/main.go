package main

// go:generate sqlboiler postgres

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"runtime/pprof"
	"time"

	_ "github.com/lib/pq"
	"github.com/vattle/sqlboiler/boil"
	// "github.com/vevsatechnologies/External_Data_Feed_Processor/tree/master/models"
	log15 "gopkg.in/inconshreveable/log15.v2"
)

type Config struct {
	historicData []string
	chartData    string
	posData      string
	powData      []string
}

// Open handle to database like normal
var log = log15.New()

// mainCore does all the work. Deferred functions do not run after os.Exit(),
// so main wraps this function, which returns a code.
func mainCore() error {
	// Parse the configuration file, and setup logger.
	cfg, err := loadConfig()
	if err != nil {
		fmt.Printf("Failed to load dcrdata config: %s\n", err.Error())
		return err
	}

	defer func() {
		if logRotator != nil {
			logRotator.Close()
		}
	}()

	if cfg.CPUProfile != "" {
		var f *os.File
		f, err = os.Create(cfg.CPUProfile)
		if err != nil {
			return err
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
}

func main() {

	if err := mainCore(); err != nil {
		if logRotator != nil {
			log.Error(err)
		}
		os.Exit(1)
	}
	os.Exit(0)

	db, err := sql.Open("postgres", "dbname=data_feed_processor user=postgres host=localhost password=alisha")
	if err != nil {
		panic(err.Error())
		return
	}

	boil.SetDB(db)

	// fetchHistoricData(1, "") //parameters : exchangeID,currency pair, start time, end time
	getPOSdata()
	// for {
	// 	getHistoricData(1, "BTC-DCR", "1514764800", "1514851200")   //parameters : exchangeID,currency pair, start time, end time                                    //parameters :  Currency pair
	// 	getChartData(1,"BTC_DCR", "1405699200", "9999999999") //parameters: exchange id,Currency Pair, start time , end time

	// }

	// getPOWData(2, "") //parameters: pool id

}

// Exchange id = 0 for Poloneix
// Exchange id =1 for Bittrex

func fetchHistoricData(exchangeID int, date string) {

	if exchangeID == 0 { //Poloniex exchange
		user := Poloniex{

			client: &http.Client{},
		}
		user.fetchPoloniexData(date)

	}

	if exchangeID == 1 { //Bittrex Exchange

		user := Bittrex{
			client: &http.Client{},
		}
		user.fetchBittrexData(date)
	}

}

func getPOSdata() {

	user := POS{
		client: &http.Client{},
	}

	user.getPOS()
}

func getPOWData(PoolID int, api_key string) {

	user := POW{
		client: &http.Client{},
	}
	if PoolID == 1 {
		url := url1
		user.getPOW(PoolID, url, api_key)
	}
	if PoolID == 2 {
		url := url2
		user.getPOW(PoolID, url, api_key)
	}

	if PoolID == 3 {
		url := url3
		user.getPOW(PoolID, url, api_key)
	}

	if PoolID == 4 {
		url := url4
		user.getPOW(PoolID, url, api_key)
	}

	if PoolID == 5 {
		url := url5
		user.getPOW(PoolID, url, api_key)
	}

	if PoolID == 6 {
		url := url6
		user.getPOW(PoolID, url, api_key)
	}
	if PoolID == 7 {
		url := url7
		user.getPOW(PoolID, url, api_key)
	}

	if PoolID == 8 {
		url := url8
		user.getPOW(PoolID, url, api_key)
	}

}

// Exchange id = 0 for Poloneix
// Exchange id =1 for Bittrex

func getHistoricData(exchangeID int, currencyPair string, startTime string, endTime string) {

	if exchangeID == 0 { //Poloniex exchange
		user := Poloniex{

			client: &http.Client{},
		}
		user.getPoloniexData(currencyPair, startTime, endTime)

	}

	if exchangeID == 1 { //Bittrex Exchange

		user := Bittrex{
			client: &http.Client{},
		}
		user.getBittrexData(currencyPair)
	}

	//Time delay of 24 hours

	time.Sleep(86400 * time.Second)
}

//Get chart data from exchanges
// if exchange id =0 , get chart data from poloniex exchange
// exchange id =1  get chart data from bittrex exchange

func getChartData(exchangeID int, currencyPair string, startTime string, endTime string) {

	if exchangeID == 0 {
		user := Poloniex{

			client: &http.Client{},
		}
		user.getChartData(currencyPair, startTime, endTime)

	}
	if exchangeID == 1 {
		user := Bittrex{
			client: &http.Client{},
		}
		user.getTicks(currencyPair)

	}

}
