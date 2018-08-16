package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spf13/viper"
	"github.com/vattle/sqlboiler/boil"
	"github.com/vevsatechnologies/External_Data_Feed_Processor/models"
	null "gopkg.in/nullbio/null.v6"
)

// Structure containing Poloniex client data

type Poloniex struct {
	client *http.Client
}

//Structure containing Poloniex Historic Data

type poloniexData struct {
	Result []struct {
		GlobalTradeID null.String `json:"globalTradeID"`
		TradeID       null.String `json:"tradeID"`
		Date          null.Time   `json:"date"`
		Types         null.String `json:"type"`
		Rate          null.String `json:"rate"`
		Amount        null.String `json:"amount"`
		Total         null.String `json:"total"`
	}
}

// Structure containing Poloniex Chart Data

type chartData struct {
	Result []struct {
		Date            null.Time   `json:"date"`
		High            null.String `json:"high"`
		Low             null.String `json:"low"`
		Open            null.String `json:"open"`
		Close           null.String `json:"close"`
		Volume          null.String `json:"volume"`
		QuoteVolume     null.String `json:"quoteVolume"`
		WeightedAverage null.String `json:"weightedAverage"`
	}
}

// Get Poloniex Historic Data
// parameters : Currency pair, Start time , End time

func (p *Poloniex) getPoloniexData(currencyPair string, start string, end string) string {

	dbInfo := fmt.Sprintf("host=%s port=%s user=%s "+"password=%s dbname=%s sslmode=disable", viper.Get("Database.pghost"), viper.Get("Database.pgport"), viper.Get("Database.pguser"), viper.Get("Database.pgpass"), viper.Get("Database.pgdbname"))

	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		panic(err.Error())

	}
	boil.SetDB(db)

	//Get Url of Poloniex API

	url := viper.Get("ExchangeData.0").(string)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err.Error())
	}

	//Append user provided parameters in the URL

	q := req.URL.Query()
	q.Add("command", "returnTradeHistory")
	q.Add("currencyPair", currencyPair)
	q.Add("start", start)
	q.Add("end", end)
	req.URL.RawQuery = q.Encode()

	//Get Historic Data from the API

	request, err := http.NewRequest("GET", req.URL.String(), nil)

	res, _ := p.client.Do(request)

	//Get response of the request as []byte

	fmt.Println(res.StatusCode)
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}

	//Map the data to poloniexData struct and unmarshal the contents

	var data poloniexData
	json.Unmarshal(body, &data)

	fmt.Printf("Results: %v\n", data)

	for i := range data.Result {
		var p1 models.HistoricDatum

		// p1.Exchangeid = 0
		p1.Globaltradeid = data.Result[i].GlobalTradeID
		p1.Tradeid = data.Result[i].TradeID
		// p1.timest = data.Result[i].Date
		p1.Quantity = data.Result[i].Amount
		p1.Price = data.Result[i].Rate
		p1.Total = data.Result[i].Total
		p1.OrderType = data.Result[i].Types
		err := p1.Insert(db)
		panic(err.Error())
	}

	return "Saved poloneix historic data!"
}

//Returns data from Poloniex exchange
//Parameters : currency pair , start time , end time

// func (p *Poloniex) fetchPoloniexData(date string) {

// 	err := models.HistoricDatum(qm.Where("timestamp=$1", date)).All()

// }

//Returns Poloniex Chart Data
//Parameters : Currency pair, Start time , End time

func (p *Poloniex) getChartData(currencyPair string, start string, end string) {

	//Get the base URL

	url := viper.Get("ExchangeData.0").(string)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err.Error())
	}

	//Append the user defined parameters to the url

	q := req.URL.Query()
	q.Add("command", "returnChartData")
	q.Add("currencyPair", currencyPair)
	q.Add("start", start)
	q.Add("end", end)
	req.URL.RawQuery = q.Encode()

	request, err := http.NewRequest("GET", req.URL.String(), nil)

	//Get the data from API and convert the data to byte array

	res, _ := p.client.Do(request)

	fmt.Println(res.StatusCode)
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}

	//Store the data to charData struct

	var data chartData
	json.Unmarshal(body, &data)
	fmt.Printf("Results: %v\n", data)

	//Loop over the entire data and store it in the table
	for i := range data.Result {

		var p2 models.ChartDatum

		// p2.exchangeID = "0"
		//p2.date = data.Result[i].Date

		p2.High = data.Result[i].High
		p2.Low = data.Result[i].Low
		p2.Opening = data.Result[i].Open
		p2.Closing = data.Result[i].Close
		p2.Volume = data.Result[i].Volume
		p2.Quotevolume = data.Result[i].QuoteVolume
		// p2.Basevolume = data.Result[i].BaseVolume
		p2.Weightedaverage = data.Result[i].WeightedAverage

		err := p2.Insert(db)
		panic(err.Error())

	}

}
