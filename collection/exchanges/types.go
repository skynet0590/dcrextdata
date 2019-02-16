package exchanges

import "time"

type DataTick struct {
	High     float64
	Low      float64
	Open     float64
	Close    float64
	Volume   float64
	Time     int64
	Exchange string
}

type poloniexAPIResponse []poloniexDataTick

type poloniexDataTick struct {
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Open   float64 `json:"open"`
	Close  float64 `json:"close"`
	Volume float64 `json:"volume"`
	Time   int64   `json:"date"`
}

func (resp poloniexAPIResponse) DataTicks() []DataTick {
	edata := make([]DataTick, 0, len(resp))

	for _, v := range resp {
		edata = append(edata, DataTick{
			High:     v.High,
			Low:      v.Low,
			Open:     v.Open,
			Close:    v.Close,
			Time:     v.Time,
			Exchange: "poloniex",
		})
	}

	return edata
}

type bittrexDataTick struct {
	High   float64 `json:"H"`
	Low    float64 `json:"L"`
	Open   float64 `json:"O"`
	Close  float64 `json:"C"`
	Volume float64 `json:"BV"`
	Time   string  `json:"T"`
}

type bittrexAPIResponse struct {
	Result []bittrexDataTick `json:"result"`
}

func (resp bittrexAPIResponse) DataTicks(start int64) []DataTick {
	edata := make([]DataTick, 0, len(resp.Result))
	for _, v := range resp.Result {
		t, _ := time.Parse("2006-01-02T15:04:05", v.Time)

		// Skip all entries before the required start time
		if t.Unix() < start {
			continue
		}

		edata = append(edata, DataTick{
			High:     v.High,
			Low:      v.Low,
			Open:     v.Open,
			Close:    v.Close,
			Time:     t.Unix(),
			Exchange: "bittrex",
		})
	}

	return edata
}
