// Copyright (c) 2018-2019 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package ticks

import (
	"context"
	"strconv"
	"time"
)

type Collector interface {
	GetShort(context.Context) error
	GetLong(context.Context) error
	GetHistoric(context.Context) error
}

type Store interface {
	RegisterExchange(ctx context.Context, exchange ExchangeData) (lastShort, lastLong, lastHistoric time.Time, err error)
	StoreExchangeTicks(ctx context.Context, exchange string, interval time.Duration, intervalString, pair string, data []Tick) (time.Time, error)
}

type tickable interface {
	toTicks() []Tick
}

// Tick represents an exchange data tick
type Tick struct {
	High   float64
	Low    float64
	Open   float64
	Close  float64
	Volume float64
	Time   time.Time
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

func (resp poloniexAPIResponse) toTicks() []Tick {
	res := []poloniexDataTick(resp)
	dataTicks := make([]Tick, 0, len(res))
	for _, v := range res {
		dataTicks = append(dataTicks, Tick{
			High:   v.High,
			Low:    v.Low,
			Open:   v.Open,
			Close:  v.Close,
			Volume: v.Volume,
			Time:   time.Unix(v.Time, 0),
		})
	}
	return dataTicks
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

func (resp bittrexAPIResponse) toTicks() []Tick {
	bTicks := resp.Result
	dataTicks := make([]Tick, 0, len(bTicks))
	for _, v := range bTicks {
		t, err := time.Parse("2006-01-02T15:04:05", v.Time)
		if err != nil {
			continue
		}
		dataTicks = append(dataTicks, Tick{
			High:   v.High,
			Low:    v.Low,
			Open:   v.Open,
			Close:  v.Close,
			Volume: v.Volume,
			Time:   t,
		})
	}
	return dataTicks
}

type bleutradeDataTick struct {
	High   float64 `json:"High"`
	Low    float64 `json:"Low"`
	Open   float64 `json:"Open"`
	Close  float64 `json:"Close"`
	Volume float64 `json:"Volume"`
	Time   string  `json:"TimeStamp"`
}

type bleutradeAPIResponse struct {
	Result []bleutradeDataTick `json:"result"`
}

func (resp bleutradeAPIResponse) toTicks() []Tick {
	res := resp.Result
	dataTicks := make([]Tick, 0, len(res))
	for _, v := range res {
		t, err := time.Parse("2006-01-02 15:04:05", v.Time)
		if err != nil {
			continue
		}
		dataTicks = append(dataTicks, Tick{
			High:   v.High,
			Low:    v.Low,
			Open:   v.Open,
			Close:  v.Close,
			Volume: v.Volume,
			Time:   t,
		})
	}
	return dataTicks
}

type binanceAPIResponse []binanceDataTick
type binanceDataTick []interface{}

func (resp binanceAPIResponse) toTicks() []Tick {
	res := []binanceDataTick(resp)
	dataTicks := make([]Tick, 0, len(res))
	for _, j := range res {
		high, err := strconv.ParseFloat(j[2].(string), 64)
		if err != nil {
			continue
		}
		low, err := strconv.ParseFloat(j[3].(string), 64)
		if err != nil {
			continue
		}
		open, err := strconv.ParseFloat(j[1].(string), 64)
		if err != nil {
			continue
		}
		close, err := strconv.ParseFloat(j[4].(string), 64)
		if err != nil {
			continue
		}
		volume, err := strconv.ParseFloat(j[5].(string), 64)
		if err != nil {
			continue
		}
		// Converting unix time from milliseconds to seconds
		t := time.Unix(int64(j[0].(float64)/1000), 0)
		dataTicks = append(dataTicks, Tick{
			High:   high,
			Low:    low,
			Open:   open,
			Close:  close,
			Volume: volume,
			Time:   t,
		})
	}
	return dataTicks
}
