package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const (
	url = "https://api.decred.org/?c=gsd"
)

type POS struct {
	client *http.Client
}

type POSData struct {
	APIEnabled           bool    `json:"APIEnabled"`
	APIVersionsSupported []int   `json:"APIVersionsSupported"`
	Network              string  `json:"Network"`
	URL                  string  `json:"URL"`
	Launched             int64   `json:"Launched"`
	LastUpdated          int64   `json:"LastUpdated"`
	Immature             int64   `json:"Immature"`
	Live                 int64   `json:"Live"`
	Voted                int64   `json:"Voted"`
	Missed               int64   `json:"Missed"`
	PoolFees             float64 `json:"PoolFees"`
	ProportionLive       float64 `json:"ProportionLive"`
	ProportionMissed     float64 `json:"ProportionMissed"`
	UserCount            int64   `json:"UserCount"`
	UserCountActive      int64   `json:"UserCountActive"`
}

type Data map[string]POSData

func (p *POS) getPOS() {

	url := url
	request, err := http.NewRequest("GET", url, nil)

	res, _ := p.client.Do(request)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}

	var data Data

	json.Unmarshal(body, &data)

	//Loop over the entire list to insert data into the table

	for key, value := range data {

		sqlStatement := `INSERT INTO POSdataTable(POSid,APIEnabled,APIVersionsSupported,Network,Url,Launched,LastUpdated,Immature,Live,Voted,Missed,PoolFees,ProportionLive,ProportionMissed,UserCount,UserCountActive,timestamp) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)`
		_, err = db.Exec(sqlStatement, key, value.APIEnabled, value.APIVersionsSupported, value.Network, value.URL, value.Launched, value.LastUpdated, value.Immature, value.Live, value.Voted, value.Missed, value.PoolFees, value.ProportionLive, value.ProportionMissed, value.UserCount, value.UserCountActive, "NOW()")

	}

}
