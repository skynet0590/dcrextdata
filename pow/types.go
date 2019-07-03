package pow

import "time"

type PowDataSource struct {
	Source string
}

type PowData struct {
	Time              int64
	NetworkHashrate   int64
	PoolHashrate      float64
	Workers           int64
	NetworkDifficulty float64
	CoinPrice         float64
	BtcPrice          float64
	Source            string
}

type PowDataDto struct {
	Time              time.Time
	NetworkHashrate   int64
	PoolHashrate      float64
	Workers           int64
	NetworkDifficulty float64
	CoinPrice         float64
	BtcPrice          float64
	Source            string
}

type luxorPowData struct {
	Time              string  `json:"time"`
	NetworkHashrate   int64   `json:"network_hashrate"`
	PoolHashrate      float64 `json:"pool_hashrate"`
	Workers           int64   `json:"workers"`
	NetworkDifficulty float64 `json:"network_difficulty"`
	CoinPrice         string  `json:"coin_price"`
	BtcPrice          string  `json:"btc_price"`
}

type luxorAPIResponse struct {
	GlobalStats []luxorPowData `json:"globalStats"`
}

type f2poolPowData map[string]float64

type f2poolAPIResponse struct {
	Hashrate f2poolPowData `json:"hashrate_history"`
}

type coinmineAPIResponse struct {
	NetworkHashrate int64   `json:"network_hashrate"`
	PoolHashrate    float64 `json:"hashrate"`
	Workers         int64   `json:"workers"`
}

type btcData struct {
	NetworkHashrate string              `json:"network_hashrate"`
	PoolHashrate    string              `json:"pool_hashrate"`
	Rates           btcExchangeRateData `json:"exchange_rate"`
}

type btcExchangeRateData struct {
	CoinPrice float64 `json:"DCR2USD"`
}
type btcAPIResponse struct {
	BtcData btcData `json:"data"`
}
