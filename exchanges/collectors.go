// Copyright (c) 2018-2019 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package exchanges

import (
	"net/http"
	"sync"
	"time"

	"github.com/raedahgroup/dcrextdata/helpers"
)

type Collector struct {
	exchanges []Exchange
	period    int64
}

func NewCollector(exchangeLasts map[string]int64, period int64) (*Collector, error) {
	exchanges := make([]Exchange, 0, len(exchangeLasts))

	for exchange, last := range exchangeLasts {
		if contructor, ok := ExchangeConstructors[exchange]; ok {

			ex, err := contructor(&http.Client{Timeout: 300 * time.Second}, last, period) // Consider if sharing a single client is better
			if err != nil {
				return nil, err
			}
			lastStr := helpers.UnixTimeToString(ex.LastUpdateTime())
			if last == 0 {
				lastStr = "never"
			}
			log.Infof("Starting exchange collector for %s, last collect time: %s", exchange, lastStr)
			exchanges = append(exchanges, ex)
		}
	}

	return &Collector{
		exchanges: exchanges,
		period:    period,
	}, nil
}

func (ec *Collector) HistoricSync(data chan []DataTick) []error {
	now := time.Now().Unix()
	wg := new(sync.WaitGroup)
	errs := make([]error, 0)
	errMtx := new(sync.Mutex)
	for _, ex := range ec.exchanges {
		l := ex.LastUpdateTime()
		if now-l <= ec.period {
			continue
		}
		wg.Add(1)
		go func(ex Exchange, errMtx *sync.Mutex, wg *sync.WaitGroup) {
			err := ex.Historic(data)
			if err != nil {
				errMtx.Lock()
				errs = append(errs, err)
				errMtx.Unlock()
			} else {
				log.Infof("Completed historic sync for %s", ex.Name())
			}
			wg.Done()
		}(ex, errMtx, wg)
	}

	wg.Wait()
	return errs
}

func (ec *Collector) Collect(data chan []DataTick, wg *sync.WaitGroup, quit chan struct{}) {
	ticker := time.NewTicker(time.Duration(ec.period) * time.Second)
	for {
		select {
		case <-ticker.C:
			log.Trace("Triggering exchange collectors")
			for _, ex := range ec.exchanges {
				go func(exh Exchange) {
					err := exh.Collect(data)
					if err != nil {
						log.Error(err)
					}
				}(ex)
			}
		case <-quit:
			log.Infof("Stopping collector")
			wg.Done()
			return
		}

	}
}
