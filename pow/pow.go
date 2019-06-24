// Copyright (c) 2018-2019 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package pow

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/raedahgroup/dcrextdata/helpers"
)

const (
	Luxor    = "luxor"
	LuxorUrl = "http://mining.luxor.tech/API/DCR/stats"

	F2pool    = "f2pool"
	F2poolUrl = "https://api.f2pool.com/decred/"

	Coinmine    = "coinmine"
	CoinmineUrl = "https://www2.coinmine.pl/dcr/index.php?page=api&action=public"

	Btc    = "btc"
	BtcUrl = "https://pool.api.btc.com/v1/pool/status"
)

var PowConstructors = map[string]func(*http.Client, int64) (Pow, error){
	Luxor:    NewLuxor,
	F2pool:   NewF2pool,
	Coinmine: NewCoinmine,
	Btc:      NewBtc,
}

type Pow interface {
	Collect(ctx context.Context) ([]PowData, error)
	LastUpdateTime() int64
	Name() string
}

type CommonInfo struct {
	client     *http.Client
	lastUpdate int64
	baseUrl    string
}

func (in *CommonInfo) LastUpdateTime() int64 {
	return in.lastUpdate
}

type LuxorPow struct {
	CommonInfo
}

func NewLuxor(client *http.Client, lastUpdate int64) (Pow, error) {
	if client == nil {
		return nil, nilClientError
	}
	return &LuxorPow{
		CommonInfo: CommonInfo{
			client:     client,
			lastUpdate: lastUpdate,
			baseUrl:    LuxorUrl,
		},
	}, nil
}

func (in *LuxorPow) Collect(ctx context.Context) ([]PowData, error) {
	res := new(luxorAPIResponse)
	err := helpers.GetResponse(ctx, in.client, in.baseUrl, res)

	if err != nil {
		return nil, err
	}

	result := in.fetch(res, in.lastUpdate)
	in.lastUpdate = result[len(result)-1].Time

	return result, nil
}

func (LuxorPow) fetch(res *luxorAPIResponse, start int64) []PowData {
	data := make([]PowData, 0, len(res.GlobalStats))
	for _, j := range res.GlobalStats {
		t, _ := time.Parse(time.RFC3339, j.Time)

		if t.Unix() < start {
			continue
		}

		coinPrice, err := strconv.ParseFloat(j.CoinPrice, 64)
		if err != nil {
			continue
		}
		btcPrice, err := strconv.ParseFloat(j.BtcPrice, 64)
		if err != nil {
			continue
		}

		data = append(data, PowData{
			Time:              t.Unix(),
			NetworkHashrate:   j.NetworkHashrate,
			PoolHashrate:      j.PoolHashrate,
			Workers:           j.Workers,
			NetworkDifficulty: j.NetworkDifficulty,
			CoinPrice:         coinPrice,
			BtcPrice:          btcPrice,
			Source:            "luxor",
		})
	}
	return data
}

func (*LuxorPow) Name() string { return Luxor }

type F2poolPow struct {
	CommonInfo
}

func NewF2pool(client *http.Client, lastUpdate int64) (Pow, error) {
	if client == nil {
		return nil, nilClientError
	}
	return &F2poolPow{
		CommonInfo: CommonInfo{
			client:     client,
			lastUpdate: lastUpdate,
			baseUrl:    F2poolUrl,
		},
	}, nil
}

func (in *F2poolPow) Collect(ctx context.Context) ([]PowData, error) {
	res := new(f2poolAPIResponse)
	err := helpers.GetResponse(ctx, in.client, in.baseUrl, res)

	if err != nil {
		return nil, err
	}

	result := in.fetch(res, in.lastUpdate)
	in.lastUpdate = result[len(result)-1].Time

	return result, nil
}

func (F2poolPow) fetch(res *f2poolAPIResponse, start int64) []PowData {
	data := make([]PowData, 0, len(res.Hashrate))
	for k, v := range res.Hashrate {
		t, _ := time.Parse(time.RFC3339, k)

		if t.Unix() < start {
			continue
		}

		data = append(data, PowData{
			Time:              t.Unix(),
			NetworkHashrate:   0,
			PoolHashrate:      v,
			Workers:           0,
			NetworkDifficulty: 0,
			CoinPrice:         0,
			BtcPrice:          0,
			Source:            "f2pool",
		})
	}
	return data
}

func (*F2poolPow) Name() string { return F2pool }

type CoinminePow struct {
	CommonInfo
}

func NewCoinmine(client *http.Client, lastUpdate int64) (Pow, error) {
	if client == nil {
		return nil, nilClientError
	}
	return &CoinminePow{
		CommonInfo: CommonInfo{
			client:     client,
			lastUpdate: lastUpdate,
			baseUrl:    CoinmineUrl,
		},
	}, nil
}

func (in *CoinminePow) Collect(ctx context.Context) ([]PowData, error) {
	res := new(coinmineAPIResponse)
	err := helpers.GetResponse(ctx, in.client, in.baseUrl, res)

	if err != nil {
		return nil, err
	}

	result := in.fetch(res, in.lastUpdate)
	in.lastUpdate = result[len(result)-1].Time

	return result, nil
}

func (CoinminePow) fetch(res *coinmineAPIResponse, start int64) []PowData {
	data := make([]PowData, 0, 1)
	t := time.Now().Unix()

	data = append(data, PowData{
		Time:              t,
		NetworkHashrate:   res.NetworkHashrate,
		PoolHashrate:      res.PoolHashrate,
		Workers:           res.Workers,
		NetworkDifficulty: 0,
		CoinPrice:         0,
		BtcPrice:          0,
		Source:            "coinmine",
	})
	return data
}

func (*CoinminePow) Name() string { return Coinmine }

type BtcPow struct {
	CommonInfo
}

func NewBtc(client *http.Client, lastUpdate int64) (Pow, error) {
	if client == nil {
		return nil, nilClientError
	}
	return &BtcPow{
		CommonInfo: CommonInfo{
			client:     client,
			lastUpdate: lastUpdate,
			baseUrl:    BtcUrl,
		},
	}, nil
}

func (in *BtcPow) Collect(ctx context.Context) ([]PowData, error) {
	res := new(btcAPIResponse)
	err := helpers.GetResponse(ctx, in.client, in.baseUrl, res)

	if err != nil {
		return nil, err
	}

	result := in.fetch(res, in.lastUpdate)
	if len(result) > 0 {
		in.lastUpdate = result[len(result)-1].Time
	}

	return result, nil
}

func (BtcPow) fetch(res *btcAPIResponse, start int64) []PowData {
	data := make([]PowData, 0, 1)
	t := time.Now().Unix()

	n, err := strconv.ParseFloat(res.BtcData.NetworkHashrate, 64)
	if err != nil {
		return nil
	}
	p, err := strconv.ParseFloat(res.BtcData.PoolHashrate, 64)
	if err != nil {
		return nil
	}

	networkHashrate := int64(1000000000000000 * n)
	poolHashrate := 1000000000000000 * p

	data = append(data, PowData{
		Time:              t,
		NetworkHashrate:   networkHashrate,
		PoolHashrate:      poolHashrate,
		Workers:           0,
		NetworkDifficulty: 0,
		CoinPrice:         0,
		BtcPrice:          res.BtcData.Rates.CoinPrice,
		Source:            "btc",
	})
	return data
}

func (*BtcPow) Name() string { return Btc }
