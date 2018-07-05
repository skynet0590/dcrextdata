package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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
	O  float64 `json:"O"`
	H  float64 `json:"H"`
	L  float64 `json:"L"`
	C  float64 `json:"C"`
	V  float64 `json:"V"`
	T  string  `json:"T"`
	BV float64 `json:"BV"`
}

//ResultArray Export the values to ResultArray struct
type ResultArray struct {
	ID        int64   `json:"Id"`
	Timestamp string  `json:"TimeStamp"`
	Quantity  float64 `json:"Quantity"`
	Price     float64 `json:"Price"`
	Total     float64 `json:"Total"`
	Filltype  string  `json:"FillType"`
	Ordertype string  `json:"OrderType"`
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

		sqlStatement := `INSERT INTO historicData(exchangeID,globaltradeid,tradeid,timestamp,quantity,price,total,fill_type,order_type) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9)`
		_, err = db.Exec(sqlStatement, "1", data.Result[i].ID, "nil", data.Result[i].Timestamp, data.Result[i].Quantity, data.Result[i].Price, data.Result[i].Total, data.Result[i].Filltype, data.Result[i].Ordertype)

		fmt.Println()
	}
	return

}

func (b *Bittrex) fetchBittrexData(date string) {

	//Fetch Data from historicData Table

	sqlStatement := `SELECT * FROM historicData where timestamp=$1`
	err = db.QueryRow(sqlStatement).Scan(&date)
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

		sqlStatement := `INSERT INTO chartData(exhcangeID,date,high,low,open,close,volume,quoteVolume,baseVolume,weightedAverage) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`
		_, err = db.Exec(sqlStatement, "1", data.Result[i].T, data.Result[i].H, data.Result[i].L, data.Result[i].O, data.Result[i].C, data.Result[i].V, "nil", data.Result[i].BV, "nil")

		fmt.Println()
	}
	return
}
