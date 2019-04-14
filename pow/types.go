package pow

type PowData struct {
	Time              int64
	NetworkHashrate   int64
	PoolHashrate      float64
	Workers           int64
	NetworkDifficulty float64
	CoinPrice         string
	BtcPrice          string
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
