package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/vattle/sqlboiler/queries/qm"
	"github.com/vevsatechnologies/External_Data_Feed_Processor/models"
)

const (
	bittrexBaseURL = "https://bittrex.com/api/v1.1/public/getmarkethistory"

	bittrexTicksURL = "https://bittrex.com/Api/v2.0/pub/market/GetTicks"
)

type Bittrex struct {
	client *http.Client
}
type bittrexData struct {
	Success string `json:"success"`
	Message string `json:"message"`

	Result []ResultArray `json:"result"`
}

type ticksData struct {
	Success string `json:"success"`
	Message string `json:"message"`

	Result []tickDataArray `json:"result"`
}

type tickDataArray struct {
	O  string `json:"O"`
	H  string `json:"H"`
	L  string `json:"L"`
	C  string `json:"C"`
	V  string `json:"V"`
	T  string `json:"T"`
	BV string `json:"BV"`
}

//ResultArray Export the values to ResultArray struct
type ResultArray struct {
	ID        int64  `json:"Id"`
	Timestamp string `json:"TimeStamp"`
	Quantity  string `json:"Quantity"`
	Price     string `json:"Price"`
	Total     string `json:"Total"`
	Filltype  string `json:"FillType"`
	Ordertype string `json:"OrderType"`
}

//Function to Return Historic Pricing Data from Bittrex Exchange
//Parameters : Currency Pair

func (b *Bittrex) getBittrexData(currencyPair string) {

	//Get the base url

	url := bittrexBaseURL
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err.Error())
	}
	q := req.URL.Query()

	//Append the user defined parameters to complete the url

	q.Add("market", currencyPair)

	req.URL.RawQuery = q.Encode()

	//Sends the GET request to the API

	request, err := http.NewRequest("GET", req.URL.String(), nil)

	res, _ := b.client.Do(request)

	// To check the status code of response
	fmt.Println(res.StatusCode)

	//Store the response in body variable as a byte array
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}

	//Store the data in bittrexData struct
	var data bittrexData

	json.Unmarshal(body, &data)
	fmt.Printf("Results: %v\n", data.Result)

	//Loop over array of struct and store them in the table

	for i := range data.Result {
		var p1 models.HistoricDatum

		p1.Exchangeid = 1
		p1.Globaltradeid = data.Result[i].ID
		p1.Tradeid = "nil"
		p1.Timestamp = data.Result[i].Timestamp
		p1.Quantity = data.Result[i].Quantity
		p1.Price = data.Result[i].Price
		p1.Total = data.Result[i].Total
		p1.fill_type = data.Result[i].Filltype
		p1.order_type = data.Result[i].Ordertype
		err := p1.Insert(db)
	}
	return

}

func (b *Bittrex) fetchBittrexData(date string) {

	//Fetch Data from historicData Table

	err := models.HistoricDatum(qm.Where("timestamp=?", date)).All()
}

//To get Ticks from Bittrex Exchange every 24 hours
//Parameters : Currency Pair

func (b *Bittrex) getTicks(currencyPair string) {

	url := bittrexTicksURL
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err.Error())
	}
	q := req.URL.Query()

	//Append user defined parameters to the base URL

	q.Add("marketName", currencyPair)
	q.Add("tickInterval", "day")

	req.URL.RawQuery = q.Encode()

	request, err := http.NewRequest("GET", req.URL.String(), nil)

	//Sends the GET request to the API and stores the response

	res, _ := b.client.Do(request)

	// To check the status code of response

	fmt.Println(res.StatusCode)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}

	//Stores the response in ticksData struct

	var data ticksData

	json.Unmarshal(body, &data)
	fmt.Printf("Results: %v\n", data.Result)

	//Loop over array of struct and stores the response in table

	for i := range data.Result {
		var p2 models.ChartDatum

		p1.Exchangeid = 1
		p1.date = data.Result[i].T
		p1.high = data.Result[i].H
		p1.low = data.Result[i].O
		p1.open = data.Result[i].C
		p1.close = data.Result[i].V
		p1.volume = "nil"
		p1.quoteVolume = data.Result[i].BV
		p1.weightedAverage = "nil"
		err := p1.Insert(db)
	}
	return
}
