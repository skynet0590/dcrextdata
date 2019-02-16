package exchanges

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/raedahgroup/dcrextdata/collection/internal"
)

const (
	poloniexStartTime   = 1463364000 // Approximate poloniex data collection start time
	poloniexVolumeLimit = 20000      // Approximate poloniex data response limit
)

var bittrexTimeString = map[int64]string{
	300:  "fiveMin",
	1800: "thirtyMin",
}

type Poloniex struct {
	client   *http.Client
	mtx      *sync.Mutex
	lastTime int64
	period   int64
}

func NewPoloniex(client *http.Client, lastTime, period int64) *Poloniex {
	if client == nil {
		return nil
	}

	if lastTime < 1 {
		lastTime = poloniexStartTime
	}
	return &Poloniex{
		client:   client,
		lastTime: lastTime,
		period:   period,
		mtx:      new(sync.Mutex),
	}
}

func (collector *Poloniex) Retrieve() ([]DataTick, error) {
	collector.mtx.Lock()
	defer collector.mtx.Unlock()

	now := time.Now().Unix()

	if now-collector.lastTime < collector.period {
		return nil, nil
	}

	var result []DataTick

	for (now-collector.lastTime)/collector.period >= poloniexVolumeLimit {
		end := collector.lastTime + poloniexVolumeLimit*collector.period

		resp, err := collector.retrieve(collector.lastTime, end, collector.period)

		if err != nil {
			return result, err
		}

		result = append(result, resp...)

		collector.lastTime = result[len(result)-1].Time
	}

	resp, err := collector.retrieve(collector.lastTime, now, collector.period)

	if err != nil {
		return result, err
	}

	result = append(result, resp...)

	collector.lastTime = result[len(result)-1].Time

	return result, nil
}

func (collector *Poloniex) retrieve(start, end, period int64) ([]DataTick, error) {
	resp := new(poloniexAPIResponse)
	requestStr := fmt.Sprintf("https://poloniex.com/public?command=returnChartData&currencyPair=BTC_DCR&start=%d&end=%d&period=%d", collector.lastTime, end, collector.period)

	err := internal.ResponseToType(collector.client, requestStr, resp)

	if err != nil {
		return nil, err
	}
	return resp.DataTicks(), nil
}

type Bittrex struct {
	client   *http.Client
	lastTime int64
	url      string
	mtx      *sync.Mutex
}

func NewBittrex(client *http.Client, lastTime, period int64) *Bittrex {
	if client == nil {
		return nil
	}

	return &Bittrex{
		client:   client,
		lastTime: lastTime,
		url:      fmt.Sprint("https://bittrex.com/Api/v2.0/pub/market/GetTicks?marketName=BTC-DCR&tickInterval=" + bittrexTimeString[period]),
		mtx:      new(sync.Mutex),
	}
}

func (collector *Bittrex) Retrieve() ([]DataTick, error) {
	collector.mtx.Lock()
	defer collector.mtx.Unlock()

	resp := new(bittrexAPIResponse)

	err := internal.ResponseToType(collector.client, collector.url, resp)

	if err != nil {
		return nil, err
	}

	result := resp.DataTicks(collector.lastTime)

	collector.lastTime = result[len(result)-1].Time

	return result, nil
}
