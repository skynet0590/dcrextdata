package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spf13/viper"
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

	url := viper.Get("POS").(string)
	request, err := http.NewRequest("GET", url, nil)

	res, _ := p.client.Do(request)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}

	dbInfo := fmt.Sprintf("host=%s port=%s user=%s "+"password=%s dbname=%s sslmode=disable",
		viper.Get("Database.pghost"), 5432, viper.Get("Database.pguser"), viper.Get("Database.pgpass"),
		viper.Get("Database.pgdbname"))

	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		panic(err.Error())
		return
	}
	var data Data

	json.Unmarshal(body, &data)

	//Loop over the entire list to insert data into the table

	for key, value := range data {

		err := db.QueryRow("Insert into pos_data Values $1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17",
			key, value.APIEnabled, value.APIVersionsSupported, value.Network, value.URL, value.Launched,
			value.LastUpdated, value.Immature, value.Live, value.Voted, value.Missed, value.PoolFees,
			value.ProportionLive, value.ProportionMissed, value.UserCount, value.UserCountActive, "NOW()")

		err := p1.Insert(db)
	}

}
