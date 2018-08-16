package main

// go:generate sqlboiler postgres

import (
	"database/sql"
	"dcrextdata/models"
	"fmt"
	"net/http"
	"time"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/queries/qm"
	log15 "gopkg.in/inconshreveable/log15.v2"
)

// Open handle to database like normal
var log = log15.New()
var psqlInfo = fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", viper.Get("Database.pghost"), viper.Get("Database.pgport"), viper.Get("Database.pguser"), viper.Get("Database.pgpass"), viper.Get("Database.pgdbname"))
var db, err = sql.Open("postgres", psqlInfo)

func main() {

	viper.SetConfigFile("./config.json")

	// Searches for config file in given paths and read it
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	viper.SetDefault("POW", "http://api.f2pool.com/decred/address")
	viper.SetDefault("ExchangeData", "https://bittrex.com/api/v1.1/public/getmarkethistory")

	if err != nil {
		panic(err.Error())
		return
	}

	boil.SetDB(db)

	// getHistoricData("bittrex", "BTC-DCR", "1514764800", "1514851200") //parameters : exchangeID,currency pair, start time, end time
	getPOSdata()
	// for {
	// getHistoricData(1, "BTC-DCR", "1514764800", "1514851200") //parameters : exchangeID,currency pair, start time, end time                                    //parameters :  Currency pair
	// 	getChartData(1, "BTC_DCR", "1405699200", "9999999999")    //parameters: exchange id,Currency Pair, start time , end time

	// }

	// getPOWData(2, "") //parameters: pool id

}

func fetchHistoricData(date string) {

	Result, err := models.HistoricDatum(qm.Where("Timest=?", date)).One(ctx, db)

	fmt.Print(Result)

}

func getPOSdata() {

	user := pos{
		client: &http.Client{},
	}

	user.getPOS()
}

func getPOWData(PoolID int, apiKey string) {

	user := pow{
		client: &http.Client{},
	}

	user.getPOW(PoolID, viper.GetString("POW"+"["+string(PoolID)+"]"), apiKey)

}

func getHistoricData(exchangeName string, currencyPair string, startTime string, endTime string) {

	if exchangeName == "poloniex" {
		user := Poloniex{

			client: &http.Client{},
		}
		user.getPoloniexData(currencyPair, startTime, endTime)

	}

	if exchangeName == "bittrex" {

		user := Bittrex{
			client: &http.Client{},
		}
		user.getBittrexData(currencyPair)
	}

	//Time delay of 24 hours

	time.Sleep(86400 * time.Second)
}

//Get chart data from exchanges

func getChartData(exchangeName string, currencyPair string, startTime string, endTime string) {

	if exchangeName == "poloniex" {
		user := Poloniex{

			client: &http.Client{},
		}
		user.getChartData(currencyPair, startTime, endTime)

	}
	if exchangeName == "Bittrex" {
		user := Bittrex{
			client: &http.Client{},
		}
		user.getTicks(currencyPair)

	}

}
