package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
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

type POSID struct {
	Delta    POSData `json:"Delta"` //0
	Echo     POSData `json:"Echo"`
	Mike     POSData `json:"Mike"`
	Ray      POSData `json:"Ray"`
	Ethan    POSData `json:"Ethan"`
	Bravo    POSData `json:"Bravo"`
	Hotel    POSData `json:"Hotel"`
	Lima     POSData `json:"Lima"`
	November POSData `json:"November"`
	Golf     POSData `json:"Golf"`
	Oscar    POSData `json:"Oscar"`
	India    POSData `json:"India"`
	Papa     POSData `json:"Papa"`
	Kilo     POSData `json:"Kilo"`
	Juliett  POSData `json:"Juliett"`
}

//POSID Each POS Id maps to the POS data
//

type Data map[string]POSData

func (p *POS) getPOS() {

	url := url
	request, err := http.NewRequest("GET", url, nil)

	res, _ := p.client.Do(request)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}

	var data POSID

	json.Unmarshal(body, &data)

	fmt.Printf("Results: %v\n", data)
	fmt.Println(reflect.TypeOf(data.Delta))

	// ids := [] struct{Golf, Hotel, Juliett, Lima, November, Ray,India, Papa, Echo, Mike, Oscar, Ethan, Bravo, Delta,Kilo}
	//Loop over the entire list to insert data into the table

	// for i := 0; i < 15; i++ {
	sqlStatement := `INSERT INTO POSdataTable(POSid,APIEnabled,APIVersionsSupported,Network,Url,Launched,LastUpdated,Immature,Live,Voted,Missed,PoolFees,ProportionLive,ProportionMissed,UserCount,UserCountActive,timestamp) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)`
	_, err = db.Exec(sqlStatement, "0", data.Delta.APIEnabled, data.Delta.APIVersionsSupported, data.Delta.Network, data.Delta.URL, data.Delta.Launched, data.Delta.LastUpdated, data.Delta.Immature, data.Delta.Live, data.Delta.Voted, data.Delta.Missed, data.Delta.PoolFees, data.Delta.ProportionLive, data.Delta.ProportionMissed, data.Delta.UserCount, data.Delta.UserCountActive, "NOW()")
	sqlStatement = `INSERT INTO POSdataTable(POSid,APIEnabled,APIVersionsSupported,Network,Url,Launched,LastUpdated,Immature,Live,Voted,Missed,PoolFees,ProportionLive,ProportionMissed,UserCount,UserCountActive,timestamp) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)`
	_, err = db.Exec(sqlStatement, "1", data.Echo.APIEnabled, data.Echo.APIVersionsSupported, data.Echo.Network, data.Echo.URL, data.Echo.Launched, data.Echo.LastUpdated, data.Echo.Immature, data.Echo.Live, data.Echo.Voted, data.Echo.Missed, data.Echo.PoolFees, data.Echo.ProportionLive, data.Echo.ProportionMissed, data.Echo.UserCount, data.Echo.UserCountActive)
	sqlStatement = `INSERT INTO POSdataTable(POSid,APIEnabled,APIVersionsSupported,Network,Url,Launched,LastUpdated,Immature,Live,Voted,Missed,PoolFees,ProportionLive,ProportionMissed,UserCount,UserCountActive,timestamp) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)`
	_, err = db.Exec(sqlStatement, "2", data.Mike.APIEnabled, data.Mike.APIVersionsSupported, data.Mike.Network, data.Mike.URL, data.Mike.Launched, data.Mike.LastUpdated, data.Mike.Immature, data.Mike.Live, data.Mike.Voted, data.Mike.Missed, data.Mike.PoolFees, data.Mike.ProportionLive, data.Mike.ProportionMissed, data.Mike.UserCount, data.Mike.UserCountActive)
	sqlStatement = `INSERT INTO POSdataTable(POSid,APIEnabled,APIVersionsSupported,Network,Url,Launched,LastUpdated,Immature,Live,Voted,Missed,PoolFees,ProportionLive,ProportionMissed,UserCount,UserCountActive,timestamp) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)`
	_, err = db.Exec(sqlStatement, "3", data.Ray.APIEnabled, data.Ray.APIVersionsSupported, data.Ray.Network, data.Ray.URL, data.Ray.Launched, data.Ray.LastUpdated, data.Ray.Immature, data.Ray.Live, data.Ray.Voted, data.Ray.Missed, data.Ray.PoolFees, data.Ray.ProportionLive, data.Ray.ProportionMissed, data.Ray.UserCount, data.Ray.UserCountActive)
	sqlStatement = `INSERT INTO POSdataTable(POSid,APIEnabled,APIVersionsSupported,Network,Url,Launched,LastUpdated,Immature,Live,Voted,Missed,PoolFees,ProportionLive,ProportionMissed,UserCount,UserCountActive,timestamp) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)`
	_, err = db.Exec(sqlStatement, "4", data.Ethan.APIEnabled, data.Ethan.APIVersionsSupported, data.Ethan.Network, data.Ethan.URL, data.Ethan.Launched, data.Ethan.LastUpdated, data.Ethan.Immature, data.Ethan.Live, data.Ethan.Voted, data.Ethan.Missed, data.Ethan.PoolFees, data.Ethan.ProportionLive, data.Ethan.ProportionMissed, data.Ethan.UserCount, data.Ethan.UserCountActive)
	sqlStatement = `INSERT INTO POSdataTable(POSid,APIEnabled,APIVersionsSupported,Network,Url,Launched,LastUpdated,Immature,Live,Voted,Missed,PoolFees,ProportionLive,ProportionMissed,UserCount,UserCountActive,timestamp) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)`
	_, err = db.Exec(sqlStatement, "5", data.Bravo.APIEnabled, data.Bravo.APIVersionsSupported, data.Bravo.Network, data.Bravo.URL, data.Bravo.Launched, data.Bravo.LastUpdated, data.Bravo.Immature, data.Bravo.Live, data.Bravo.Voted, data.Bravo.Missed, data.Bravo.PoolFees, data.Bravo.ProportionLive, data.Bravo.ProportionMissed, data.Bravo.UserCount, data.Bravo.UserCountActive)
	sqlStatement = `INSERT INTO POSdataTable(POSid,APIEnabled,APIVersionsSupported,Network,Url,Launched,LastUpdated,Immature,Live,Voted,Missed,PoolFees,ProportionLive,ProportionMissed,UserCount,UserCountActive,timestamp) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)`
	_, err = db.Exec(sqlStatement, "6", data.Hotel.APIEnabled, data.Hotel.APIVersionsSupported, data.Hotel.Network, data.Hotel.URL, data.Hotel.Launched, data.Hotel.LastUpdated, data.Hotel.Immature, data.Hotel.Live, data.Hotel.Voted, data.Hotel.Missed, data.Hotel.PoolFees, data.Hotel.ProportionLive, data.Hotel.ProportionMissed, data.Hotel.UserCount, data.Hotel.UserCountActive)
	sqlStatement = `INSERT INTO POSdataTable(POSid,APIEnabled,APIVersionsSupported,Network,Url,Launched,LastUpdated,Immature,Live,Voted,Missed,PoolFees,ProportionLive,ProportionMissed,UserCount,UserCountActive,timestamp) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)`
	_, err = db.Exec(sqlStatement, "7", data.Lima.APIEnabled, data.Lima.APIVersionsSupported, data.Lima.Network, data.Lima.URL, data.Lima.Launched, data.Lima.LastUpdated, data.Lima.Immature, data.Lima.Live, data.Lima.Voted, data.Lima.Missed, data.Lima.PoolFees, data.Lima.ProportionLive, data.Lima.ProportionMissed, data.Lima.UserCount, data.Lima.UserCountActive)
	sqlStatement = `INSERT INTO POSdataTable(POSid,APIEnabled,APIVersionsSupported,Network,Url,Launched,LastUpdated,Immature,Live,Voted,Missed,PoolFees,ProportionLive,ProportionMissed,UserCount,UserCountActive,timestamp) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)`
	_, err = db.Exec(sqlStatement, "8", data.November.APIEnabled, data.November.APIVersionsSupported, data.November.Network, data.November.URL, data.November.Launched, data.November.LastUpdated, data.November.Immature, data.November.Live, data.November.Voted, data.November.Missed, data.November.PoolFees, data.November.ProportionLive, data.November.ProportionMissed, data.November.UserCount, data.November.UserCountActive)
	sqlStatement = `INSERT INTO POSdataTable(POSid,APIEnabled,APIVersionsSupported,Network,Url,Launched,LastUpdated,Immature,Live,Voted,Missed,PoolFees,ProportionLive,ProportionMissed,UserCount,UserCountActive,timestamp) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)`
	_, err = db.Exec(sqlStatement, "9", data.Golf.APIEnabled, data.Golf.APIVersionsSupported, data.Golf.Network, data.Golf.URL, data.Golf.Launched, data.Golf.LastUpdated, data.Golf.Immature, data.Golf.Live, data.Golf.Voted, data.Golf.Missed, data.Golf.PoolFees, data.Golf.ProportionLive, data.Golf.ProportionMissed, data.Golf.UserCount, data.Golf.UserCountActive)
	sqlStatement = `INSERT INTO POSdataTable(POSid,APIEnabled,APIVersionsSupported,Network,Url,Launched,LastUpdated,Immature,Live,Voted,Missed,PoolFees,ProportionLive,ProportionMissed,UserCount,UserCountActive,timestamp) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)`
	_, err = db.Exec(sqlStatement, "10", data.Oscar.APIEnabled, data.Oscar.APIVersionsSupported, data.Oscar.Network, data.Oscar.URL, data.Oscar.Launched, data.Oscar.LastUpdated, data.Oscar.Immature, data.Oscar.Live, data.Oscar.Voted, data.Oscar.Missed, data.Oscar.PoolFees, data.Oscar.ProportionLive, data.Oscar.ProportionMissed, data.Oscar.UserCount, data.Oscar.UserCountActive)
	sqlStatement = `INSERT INTO POSdataTable(POSid,APIEnabled,APIVersionsSupported,Network,Url,Launched,LastUpdated,Immature,Live,Voted,Missed,PoolFees,ProportionLive,ProportionMissed,UserCount,UserCountActive,timestamp) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)`
	_, err = db.Exec(sqlStatement, "11", data.India.APIEnabled, data.India.APIVersionsSupported, data.India.Network, data.India.URL, data.India.Launched, data.India.LastUpdated, data.India.Immature, data.India.Live, data.India.Voted, data.India.Missed, data.India.PoolFees, data.India.ProportionLive, data.India.ProportionMissed, data.India.UserCount, data.India.UserCountActive)
	sqlStatement = `INSERT INTO POSdataTable(POSid,APIEnabled,APIVersionsSupported,Network,Url,Launched,LastUpdated,Immature,Live,Voted,Missed,PoolFees,ProportionLive,ProportionMissed,UserCount,UserCountActive,timestamp) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)`
	_, err = db.Exec(sqlStatement, "12", data.Papa.APIEnabled, data.Papa.APIVersionsSupported, data.Papa.Network, data.Papa.URL, data.Papa.Launched, data.Papa.LastUpdated, data.Papa.Immature, data.Papa.Live, data.Papa.Voted, data.Papa.Missed, data.Papa.PoolFees, data.Papa.ProportionLive, data.Papa.ProportionMissed, data.Papa.UserCount, data.Papa.UserCountActive)
	sqlStatement = `INSERT INTO POSdataTable(POSid,APIEnabled,APIVersionsSupported,Network,Url,Launched,LastUpdated,Immature,Live,Voted,Missed,PoolFees,ProportionLive,ProportionMissed,UserCount,UserCountActive,timestamp) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)`
	_, err = db.Exec(sqlStatement, "13", data.Kilo.APIEnabled, data.Kilo.APIVersionsSupported, data.Kilo.Network, data.Kilo.URL, data.Kilo.Launched, data.Kilo.LastUpdated, data.Kilo.Immature, data.Kilo.Live, data.Kilo.Voted, data.Kilo.Missed, data.Kilo.PoolFees, data.Kilo.ProportionLive, data.Kilo.ProportionMissed, data.Kilo.UserCount, data.Kilo.UserCountActive)
	sqlStatement = `INSERT INTO POSdataTable(POSid,APIEnabled,APIVersionsSupported,Network,Url,Launched,LastUpdated,Immature,Live,Voted,Missed,PoolFees,ProportionLive,ProportionMissed,UserCount,UserCountActive,timestamp) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)`
	_, err = db.Exec(sqlStatement, "14", data.Juliett.APIEnabled, data.Juliett.APIVersionsSupported, data.Juliett.Network, data.Juliett.URL, data.Juliett.Launched, data.Juliett.LastUpdated, data.Juliett.Immature, data.Juliett.Live, data.Juliett.Voted, data.Juliett.Missed, data.Juliett.PoolFees, data.Juliett.ProportionLive, data.Juliett.ProportionMissed, data.Juliett.UserCount, data.Juliett.UserCountActive)

	// _, err = db.Exec(sqlStatement, "0", data.ids[i].APIEnabled, data.ids[i].APIVersionsSupported, data.ids[i].Network, data.ids[i].URL, data.ids[i].Launched, data.ids[i].LastUpdated, data.ids[i].Immature, data.ids[i].Live, data.ids[i].Voted, data.ids[i].Missed, data.ids[i].PoolFees, data.ids[i].ProportionLive, data.ids[i].ProportionMissed, data.ids[i].UserCount, data.ids[i].UserCountActive, "NOW()")
	// }
	// fmt.Pr`i`ntln("Data fetched")

}
