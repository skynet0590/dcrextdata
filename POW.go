package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	url1 = "http://api.f2pool.com/decred/address"
	url2 = "https://www2.coinmine.pl/dcr/index.php?page=api&action=getpoolstatus"
	url3 = "http://mining.luxor.tech/API/DCR/stats/"
	url4 = "https://dcr.suprnova.cc/index.php?page=api&action=getpoolstatus"
	url5 = "https://altpool.net/api/status"
	url6 = "https://altpool.net/api/currencies"
	url7 = "https://pool.mn/dcr/index.php?page=api&action=getpoolstatus"
	url8 = "https://dcred.eu/api/api?q=network"
)

type POW struct {
	client *http.Client
}

type POWdata struct {
	date          string              `json : "date"`
	hashper       string              `json : "hashper" `
	blocksper     string              `json:"blocksper"`
	luck          string              `json:"luck"`
	miners        string              `json:"miners"`
	pphash        string              `json:"pphash"`
	ppshare       string              `json:"ppshare"`
	totalKickback string              `json:"total_kickback"`
	price         float64             `json:"price"`
	hashrate      float64             `json:"hashrate"`
	blocksfound   int64               `json:"blocksfound"`
	totalminers   int64               `json:"totalminers"`
	globalStats   []globalStatsValues `json:"globalStats"`
	dataVal       dataVal             `json:"data"`
	decred        altpool             `json:"decred"`
	dcr           altpoolCurrency     `json:"DCR"`
	success       string              `json:"success"`
	lastUpdate    int64               `json:"lastUpdate"`
	mainnet       mainnet             `json:"mainnet"`
	blockReward   blockReward         `json:"blockReward"`
}

type mainnet struct {
	currentHeight     int64 `json:"currentHeight"`
	networkHashrate   int64 `json:"networkHashrate"`
	networkDifficulty int64 `json:"networkDifficulty"`
}

type blockReward struct {
	total float64 `json:"total"`
	pow   float64 `json:"pow"`
	pos   float64 `json:"pos"`
	dev   float64 `json:"dev"`
}

type globalStatsValues struct {
	time              string  `json:"time"`
	networkHashrate   float64 `json:"network_hashrate"`
	poolHashrate      float64 `json:"pool_hashrate"`
	workers           int64   `json:"workers"`
	networkDifficulty float64 `json:"network_difficulty"`
	coinPrice         float64 `json:"coin_price"`
	btcPrice          float64 `json:"btc_price"`
}

type dataVal struct {
	poolName            string  `json:"pool_name"`
	hashrate            float64 `json:"hashrate"`
	efficiency          float64 `json:'efficiency"`
	progress            float64 `json:"progress"`
	workers             int64   `json:"workers"`
	currentnetworkblock int64   `json:"currentnetworkblock"`
	nextnetworkblock    int64   `json:"nextnetworkblock"`
	lastblock           int64   `json:"lastblock"`
	networkdiff         float64 `json:"networkdiff"`
	esttime             float64 `json:"esttime"`
	estshares           int64   `json:"estshares"`
	timesincelast       int64   `json:"timesincelast"`
	nethashrate         int64   `json:"nethashrate"`
}

type altpool struct {
	name             string  `json:"name"`
	port             int64   `json:"port"`
	coins            int64   `json:"coins"`
	fees             int64   `json:"fees"`
	hashrate         int64   `json:"hashrate"`
	workers          int64   `json:"workers"`
	estimate_current float64 `json:"estimate_current"`
	estimate_last24h float64 `json"estimate_last24h"`
	actual_last24h   float64 `json:"actual_last24h"`
	mbtc_mh_factor   int64   `json:"mbtc_mh_factor"`
	hashrate_last24h float64 `json:"hashrate_last24h"`
	rental_current   float64 `json:"rental_current"`
}

type altpoolCurrency struct {
	algo          string  `json:"algo"`
	port          int64   `json:"port"`
	name          string  `json:"name"`
	height        int64   `json:"height"`
	workers       int64   `json:"workers"`
	shares        int64   `json:"shares"`
	hashrate      int64   `json:"hashrate"`
	estimate      float64 `json:"estimate"`
	blocks24h     int64   `json:"24h_blocks"`
	btc24h        float64 `json:"24h_btc"`
	lastblock     int64   `json:"lastblock"`
	timesincelast int64   `json:"timesincelast"`
}

func (p *POW) getPOW(id int, url string, api_key string) {

	req, err := http.NewRequest("GET", url, nil)

	if len(api_key) != 0 {
		q := req.URL.Query()
		q.Add("api_key", api_key)
		req.URL.RawQuery = q.Encode()
	}

	request, err := http.NewRequest("GET", req.URL.String(), nil)

	res, _ := p.client.Do(request)

	fmt.Println(res.StatusCode)
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}

	var data POWdata
	json.Unmarshal(body, &data)

	fmt.Println(string(body))

	fmt.Printf("Results: %v\n", data)

	//Loop over the entire list to insert data into the table
	for i := 0; i < 15; i++ {
		sqlStatement := `INSERT INTO POWDataTable(POSId,date,hashper,blocksper,luck,miners,pphash,ppshare,totalKickback,price,hashrate,blocksfound,totalminers,success,lastUpdate,time,networkHashrate,poolHashrate,workers,networkDifficulty,coinPrice,btcPrice,poolName,efficiency,progress,workers,currentnetworkblock,nextnetworkblock,lastblock,networkdiff,esttime,estshares,timesincelast,nethashrate,name,port,coins,fees,estimate_current,estimate_last24h,actual_24h,mbtc_mh_factor,hashrate_last24h,rental_current,shares,height,estimate,blocks24h,btc24h,currentHeight,total,pos,pow,dev) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26,$27,$28,$29,$30,$31,$32,$33,$34,$35,$36,$37,$38,$39,$40,$41,$42,$43,$44,$45,$46,$47,$48,$49,$50,$51,$52,$53,$54)`
		_, err = db.Exec(sqlStatement, id, data.date, data.hashper, data.blocksper, data.luck, data.miners, data.pphash, data.ppshare, data.totalKickback, data.hashrate, data.blocksfound, data.totalminers, data.success, data.lastUpdate, data.globalStats[0].time, data.globalStats[0].networkHashrate, data.globalStats[0].poolHashrate, data.globalStats[0].workers, data.globalStats[0].networkDifficulty, data.globalStats[0].coinPrice, data.globalStats[0].btcPrice, data.dataVal.poolName, data.dataVal.efficiency, data.dataVal.progress, data.dataVal.workers, data.dataVal.currentnetworkblock, data.dataVal.nextnetworkblock, data.dataVal.lastblock, data.dataVal.networkdiff, data.dataVal.esttime, data.dataVal.estshares, data.dataVal.timesincelast, data.dataVal.nethashrate, data.decred.name, data.decred.port, data.decred.coins, data.decred.fees, data.decred.estimate_current, data.decred.estimate_last24h, data.decred.actual_last24h, data.decred.mbtc_mh_factor, data.decred.hashrate_last24h, data.decred.rental_current, data.dcr.shares, data.dcr.height, data.dcr.estimate, data.dcr.blocks24h, data.dcr.btc24h, data.mainnet.currentHeight, data.blockReward.total, data.blockReward.pos, data.blockReward.pow, data.blockReward.dev)

	}

}
