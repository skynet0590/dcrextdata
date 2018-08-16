package main

import (
	"dcrextdata/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spf13/viper"

	null "gopkg.in/nullbio/null.v6"
)

type pos struct {
	client *http.Client
}

type posData struct {
	APIEnabled           null.String  `json:"APIEnabled"`
	APIVersionsSupported []int        `json:"APIVersionsSupported"`
	Network              null.String  `json:"Network"`
	URL                  null.String  `json:"URL"`
	Launched             null.String  `json:"Launched"`
	LastUpdated          null.String  `json:"LastUpdated"`
	Immature             null.String  `json:"Immature"`
	Live                 null.String  `json:"Live"`
	Voted                null.Float64 `json:"Voted"`
	Missed               null.Float64 `json:"Missed"`
	PoolFees             null.Float64 `json:"PoolFees"`
	ProportionLive       null.Float64 `json:"ProportionLive"`
	ProportionMissed     null.Float64 `json:"ProportionMissed"`
	UserCount            null.Float64 `json:"UserCount"`
	UserCountActive      null.Float64 `json:"UserCountActive"`
}

type Data map[string]posData

func (p *pos) getPOS() {

	url := viper.Get("POS").(string)
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

		var p1 models.PosDatum

		fmt.Println(key)
		// p1.Posid = key
		p1.Apienabled = value.APIEnabled
		// p1.Apiversionssupported = value.APIVersionsSupported
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
		// p1.Timestamp = NOW()
		err := p1.Insert(db)

		panic(err.Error())
	}

}
