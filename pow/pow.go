// Copyright (c) 2018-2019 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package pow

import (
	"net/http"
	"time"

	"github.com/raedahgroup/dcrextdata/helpers"
)

const (
	Luxor    = "luxor"
	LuxorUrl = "http://mining.luxor.tech/API/DCR/stats"

	F2pool    = "f2pool"
	F2poolUrl = "https://api.f2pool.com/decred/"
)

var PowConstructors = map[string]func(*http.Client, int64) (Pow, error){
	Luxor: NewLuxor,
}

type Pow interface {
	Collect() ([]PowData, error)
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
		return nil, new(NilClientError)
	}
	return &LuxorPow{
		CommonInfo: CommonInfo{
			client:     client,
			lastUpdate: lastUpdate,
			baseUrl:    LuxorUrl,
		},
	}, nil
}

func (in *LuxorPow) Collect() ([]PowData, error) {
	res := new(luxorAPIResponse)
	err := helpers.GetResponse(in.client, LuxorUrl, res)

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

		data = append(data, PowData{
			Time:              t.Unix(),
			NetworkHashrate:   j.NetworkHashrate,
			PoolHashrate:      j.PoolHashrate,
			Workers:           j.Workers,
			NetworkDifficulty: j.NetworkDifficulty,
			CoinPrice:         j.CoinPrice,
			BtcPrice:          j.BtcPrice,
			Source:            "luxor",
		})
	}
	return data
}

func (*LuxorPow) Name() string { return Luxor }
