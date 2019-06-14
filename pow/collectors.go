// Copyright (c) 2018-2019 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package pow

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/raedahgroup/dcrextdata/helpers"
)

type PowDataStore interface {
	AddPowData(context.Context, []PowData) error
	CreatePowDataTable() error
}

type Collector struct {
	pows   []Pow
	period int64
	store  PowDataStore
}

func NewCollector(powLasts map[string]int64, period int64, store PowDataStore) (*Collector, error) {
	pows := make([]Pow, 0, len(powLasts))

	for pow, last := range powLasts {
		if contructor, ok := PowConstructors[pow]; ok {

			in, err := contructor(&http.Client{Timeout: 300 * time.Second}, last) // Consider if sharing a single client is better
			if err != nil {
				return nil, err
			}
			lastStr := helpers.UnixTimeToString(in.LastUpdateTime())
			if last == 0 {
				lastStr = "never"
			}
			log.Infof("Starting PoW collector for %s, last collect time: %s", pow, lastStr)
			pows = append(pows, in)
		}
	}

	return &Collector{
		pows:   pows,
		period: period,
		store:  store,
	}, nil
}

func (pc *Collector) Collect(ctx context.Context, wg *sync.WaitGroup) {
	ticker := time.NewTicker(time.Duration(pc.period) * time.Second)
	for {
		select {
		case <-ticker.C:
			log.Trace("Triggering PoW collectors")
			for _, in := range pc.pows {
				go func(info Pow) {
					data, err := info.Collect(ctx)
					if err != nil {
						log.Error(err)
					}
					err = pc.store.AddPowData(ctx, data)
					if err != nil {
						log.Error(err)
					}
				}(in)
			}
		case <-ctx.Done():
			log.Infof("Stopping collector")
			wg.Done()
			return
		}

	}
}
