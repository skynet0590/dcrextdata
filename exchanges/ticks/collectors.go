// Copyright (c) 2018-2019 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package ticks

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/raedahgroup/dcrextdata/helpers"
)

const ()

const (
	Bittrex   = "bittrex"
	Poloniex  = "poloniex"
	Bleutrade = "bleutrade"
	Binance   = "binance"

	btcdcrPair = "BTC/DCR"
	usdbtcPair = "USD/BTC"

	fiveMin = time.Minute * 5
	oneHour = time.Hour
	oneDay  = oneHour * 24

	apprxBinanceStart  int64 = 1540353600
	binanceVolumeLimit int64 = 1000

	apprxPoloniexStart  int64 = 1463364000
	poloniexVolumeLimit int64 = 20000

	clientTimeout = time.Minute

	IntervalShort    = "short"
	IntervalLong     = "long"
	IntervalHistoric = "historic"
)

type ExchangeData struct {
	Name             string
	WebsiteURL       string
	apiURL           string
	availableCPairs  map[string]string
	ShortInterval    time.Duration
	LongInterval     time.Duration
	HistoricInterval time.Duration
}

var (
	zeroTime time.Time

	CollectorConstructors = map[string]func(context.Context, Store) (Collector, error){
		Bittrex:         NewBittrexCollector,
		Bittrex + "usd": NewBittrexUSDCollector,
		Poloniex:        NewPoloniexCollector,
		Bleutrade:       NewBleutradeCollector,
		Binance:         NewBinanceCollector,
	}

	bittrexIntervals = map[float64]string{
		300:   "fiveMin",
		1800:  "thirtyMin",
		3600:  "hour",
		86400: "day",
	}

	bleutradeIntervals = map[float64]string{
		3600:  "1h",
		14400: "4h",
		86400: "1d",
	}

	binanceIntervals = map[float64]string{
		300:   "5m",
		3600:  "1h",
		86400: "1d",
	}
)

var (
	poloniexData = ExchangeData{
		Name:       Poloniex,
		WebsiteURL: "https://poloniex.com",
		apiURL:     "https://poloniex.com/public",
		availableCPairs: map[string]string{
			btcdcrPair: "BTC_DCR",
		},
		ShortInterval:    fiveMin,
		LongInterval:     2 * oneHour,
		HistoricInterval: oneDay,
	}

	binanceData = ExchangeData{
		Name:       Binance,
		WebsiteURL: "https://binance.com",
		apiURL:     "https://api.binance.com/api/v1/klines",
		availableCPairs: map[string]string{
			btcdcrPair: "DCRBTC",
		},
		ShortInterval:    fiveMin,
		LongInterval:     oneHour,
		HistoricInterval: oneDay,
	}

	bittrexData = ExchangeData{
		Name:       Bittrex,
		WebsiteURL: "https://bittrex.com",
		apiURL:     "https://bittrex.com/Api/v2.0/pub/market/GetTicks",
		availableCPairs: map[string]string{
			btcdcrPair: "BTC-DCR",
			usdbtcPair: "USD-BTC",
		},
		ShortInterval:    fiveMin,
		LongInterval:     oneHour,
		HistoricInterval: oneDay,
	}

	bleutradeData = ExchangeData{
		Name:       Bleutrade,
		WebsiteURL: "https://bleutrade.com",
		apiURL:     "https://bleutrade.com/api/v3/public/getcandles",
		availableCPairs: map[string]string{
			btcdcrPair: "BTC_DCR",
			// usdbtcPair: "USD-BTC",
		},
		ShortInterval:    oneHour,
		LongInterval:     4 * oneHour,
		HistoricInterval: oneDay,
	}
)

type commonExchange struct {
	*ExchangeData
	currencyPair string
	store        Store
	client       *http.Client
	lastShort    time.Time
	lastLong     time.Time
	lastHistoric time.Time
	respLock     sync.Mutex
	apiResp      tickable
	requestURLer func(time.Time, time.Duration, string) (string, error)
}

func newCommonExchange(exchange ExchangeData, store Store, resp tickable, lastShort, lastLong, lastHistoric time.Time, cpair string, requestUrler func(time.Time, time.Duration, string) (string, error)) *commonExchange {
	return &commonExchange{
		ExchangeData: &exchange,
		client:       new(http.Client),
		store:        store,
		lastShort:    lastShort,
		lastLong:     lastLong,
		lastHistoric: lastHistoric,
		requestURLer: requestUrler,
		apiResp:      resp,
		currencyPair: cpair,
	}
}

func (xc *commonExchange) GetShort(ctx context.Context) error {
	return xc.Get(ctx, &xc.lastShort, xc.ShortInterval, IntervalShort)
}

func (xc *commonExchange) GetLong(ctx context.Context) error {
	return xc.Get(ctx, &xc.lastLong, xc.LongInterval, IntervalLong)
}

func (xc *commonExchange) GetHistoric(ctx context.Context) error {
	return xc.Get(ctx, &xc.lastHistoric, xc.HistoricInterval, IntervalHistoric)
}

func (xc *commonExchange) Get(ctx context.Context, last *time.Time, interval time.Duration, intervalStr string) error {
	xc.respLock.Lock()
	defer xc.respLock.Unlock()
	for time.Now().Add(-interval).Unix() > last.Unix() {
		requestURL, err := xc.requestURLer(*last, interval, xc.availableCPairs[xc.currencyPair])
		if err != nil {
			return err
		}

		err = helpers.GetResponse(xc.client, requestURL, xc.apiResp)
		if err != nil {
			return err
		}

		ticks := xc.apiResp.toTicks()

		newLast, err := xc.store.StoreExchangeTicks(ctx, xc.Name, interval, intervalStr, xc.currencyPair, ticks)

		if newLast != zeroTime {
			*last = newLast
		}

		if err != nil {
			return err
		}
	}
	return nil
}

func NewPoloniexCollector(ctx context.Context, store Store) (Collector, error) {
	lshort, llong, lhistoric, err := store.RegisterExchange(ctx, poloniexData)
	if err != nil {
		return nil, err
	}
	if lhistoric == zeroTime {
		lhistoric = time.Unix(apprxPoloniexStart, 0)
	}

	if llong == zeroTime {
		llong = time.Now().Add((-14) * oneDay)
	}

	if lshort == zeroTime {
		lshort = time.Now().Add((-30) * oneDay)
	}

	requestUrler := func(last time.Time, interval time.Duration, cpair string) (string, error) {
		return helpers.AddParams(poloniexData.apiURL, map[string]interface{}{
			"command":      "returnChartData",
			"currencyPair": cpair,
			"start":        last.Unix(),
			"end":          time.Now().Unix(),
			"period":       int(interval.Seconds()),
		})
	}

	return newCommonExchange(poloniexData, store, new(poloniexAPIResponse), lshort, llong, lhistoric, btcdcrPair, requestUrler), nil
}

// func newBasicCommonExchange(exchange ExchangeData, store Store, resp tickable, pair string, reqUrler func(time.Time, time.Duration, string) (string, error)) *commonExchange {
// 	return newCommonExchange(exchange, store, resp, zeroTime, zeroTime, zeroTime, cpair, reqUrler)
// }

func NewBittrexCollector(ctx context.Context, store Store) (Collector, error) {
	s, l, h, err := store.RegisterExchange(ctx, bittrexData)
	if err != nil {
		return nil, err
	}

	requestUrler := func(last time.Time, interval time.Duration, cpair string) (string, error) {
		return helpers.AddParams(bittrexData.apiURL, map[string]interface{}{
			"marketName":   cpair,
			"tickInterval": bittrexIntervals[interval.Seconds()],
		})

	}
	return newCommonExchange(bittrexData, store, new(bittrexAPIResponse), s, l, h, btcdcrPair, requestUrler), nil
}

func NewBittrexUSDCollector(ctx context.Context, store Store) (Collector, error) {
	s, l, h, err := store.RegisterExchange(ctx, bittrexData)
	if err != nil {
		return nil, err
	}

	requestUrler := func(last time.Time, interval time.Duration, cpair string) (string, error) {
		return helpers.AddParams(bittrexData.apiURL, map[string]interface{}{
			"marketName":   cpair,
			"tickInterval": bittrexIntervals[interval.Seconds()],
		})

	}
	return newCommonExchange(bittrexData, store, new(bittrexAPIResponse), s, l, h, usdbtcPair, requestUrler), nil
}

func NewBleutradeCollector(ctx context.Context, store Store) (Collector, error) {
	s, l, h, err := store.RegisterExchange(ctx, bleutradeData)
	if err != nil {
		return nil, err
	}

	requestUrler := func(last time.Time, interval time.Duration, cpair string) (string, error) {
		return helpers.AddParams(bleutradeData.apiURL, map[string]interface{}{
			"market": cpair,
			"period": bleutradeIntervals[interval.Seconds()],
		})

	}
	return newCommonExchange(bleutradeData, store, new(bleutradeAPIResponse), s, l, h, btcdcrPair, requestUrler), nil
}

func NewBinanceCollector(ctx context.Context, store Store) (Collector, error) {
	lshort, llong, lhistoric, err := store.RegisterExchange(ctx, binanceData)
	if err != nil {
		return nil, err
	}
	if lhistoric == zeroTime {
		lhistoric = time.Unix(apprxBinanceStart, 0)
	}

	if llong == zeroTime {
		llong = time.Now().Add((-14) * oneDay)
	}

	if lshort == zeroTime {
		lshort = time.Now().Add((-30) * oneDay)
	}

	requestUrler := func(last time.Time, interval time.Duration, cpair string) (string, error) {
		return helpers.AddParams(binanceData.apiURL, map[string]interface{}{
			"symbol":    cpair,
			"startTime": last.Unix() * 1000,
			"endTime":   time.Now().Unix() * 1000,
			"interval":  binanceIntervals[interval.Seconds()],
			"limit":     binanceVolumeLimit,
		})

	}
	return newCommonExchange(binanceData, store, new(binanceAPIResponse), lshort, llong, lhistoric, btcdcrPair, requestUrler), nil
}
