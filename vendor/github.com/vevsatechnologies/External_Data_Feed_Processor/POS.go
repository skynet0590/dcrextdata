package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/vevsatechnologies/External_Data_Feed_Processor/models"
)

const (
	url = "https://api.decred.org/?c=gsd"
)

type POS struct {
	client *http.Client
}

type POSData struct {
	APIEnabled           string  `json:"APIEnabled"`
	APIVersionsSupported []int   `json:"APIVersionsSupported"`
	Network              string  `json:"Network"`
	URL                  string  `json:"URL"`
	Launched             string  `json:"Launched"`
	LastUpdated          string  `json:"LastUpdated"`
	Immature             string  `json:"Immature"`
	Live                 string  `json:"Live"`
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

		var p1 models.Posdatatable

		p1.Posid = key
		p1.Apienabled = value.APIEnabled
		p1.Apiversionssupported = value.APIVersionsSupported
		p1.Network = value.Network
		p1.URL = value.URL
		p1.Launched = value.Launched
		p1.Lastupdated = value.LastUpdated
		p1.Immature = value.Immature
		p1.Live = value.Live
		p1.Voted = value.Voted
		p1.Missed = value.Missed
		p1.Poolfees = value.PoolFees
		p1.Proportionlive = value.ProportionLive
		p1.Proportionmissed = value.ProportionMissed
		p1.Usercount = value.UserCount
		p1.Usercountactive = value.UserCountActive
		p1.Timestamp = NOW()

	}

}
