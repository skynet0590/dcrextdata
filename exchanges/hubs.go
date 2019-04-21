// Copyright (c) 2018-2019 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package exchanges

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/raedahgroup/dcrextdata/exchanges/ticks"
)

const (
	clientTimeout = time.Minute
)

type TickHub struct {
	collectors []ticks.Collector
	client     *http.Client
}

var NoTickCollectorsError = fmt.Errorf("No tick collectors")

func NewTickHub(ctx context.Context, exchanges []string, store ticks.Store) (*TickHub, error) {
	if len(exchanges) == 0 {
		return nil, NoTickCollectorsError
	}

	collectors := make([]ticks.Collector, 0, len(exchanges))
	for _, exchange := range exchanges {
		if contructor, ok := ticks.CollectorConstructors[exchange]; ok {
			collector, err := contructor(ctx, store)
			if err != nil {
				log.Error(err)
				continue
			}
			collectors = append(collectors, collector)
		}
	}

	if len(collectors) == 0 {
		return nil, NoTickCollectorsError
	}

	return &TickHub{
		collectors: collectors,
		client:     &http.Client{Timeout: clientTimeout},
	}, nil
}

func (hub *TickHub) CollectShort(ctx context.Context) {
	wg := new(sync.WaitGroup)
	for _, collector := range hub.collectors {
		wg.Add(1)
		go func(wg *sync.WaitGroup, collector ticks.Collector) {
			err := collector.GetShort(ctx)
			if err != nil {
				log.Error(err)
			}
			wg.Done()
		}(wg, collector)
	}
	wg.Wait()
	log.Info("Completed short collection")
}

func (hub *TickHub) CollectLong(ctx context.Context) {
	wg := new(sync.WaitGroup)
	for _, collector := range hub.collectors {
		wg.Add(1)
		go func(wg *sync.WaitGroup, collector ticks.Collector) {
			err := collector.GetLong(ctx)
			if err != nil {
				log.Error(err)
			}
			wg.Done()
		}(wg, collector)
	}
	wg.Wait()
	log.Info("Completed long collection")
}

func (hub *TickHub) CollectHistoric(ctx context.Context) {
	wg := new(sync.WaitGroup)
	for _, collector := range hub.collectors {
		wg.Add(1)
		go func(wg *sync.WaitGroup, collector ticks.Collector) {
			err := collector.GetHistoric(ctx)
			if err != nil {
				log.Error(err)
			}
			wg.Done()
		}(wg, collector)
	}
	wg.Wait()
	log.Info("Completed historic collection")
}

func (hub *TickHub) CollectAll(ctx context.Context) {
	hub.CollectShort(ctx)
	hub.CollectLong(ctx)
	hub.CollectHistoric(ctx)
}

func (hub *TickHub) Run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	shortTicker := time.NewTicker(5 * time.Minute)
	longTicker := time.NewTicker(time.Hour)
	dayTicker := time.NewTicker(24 * time.Hour)
	defer shortTicker.Stop()
	defer longTicker.Stop()
	defer dayTicker.Stop()

	if ctx.Err() != nil {
		log.Error(ctx.Err())
		return
	}
	hub.CollectAll(ctx)

	for {
		select {
		case <-shortTicker.C:
			hub.CollectShort(ctx)
		case <-longTicker.C:
			hub.CollectLong(ctx)
		case <-dayTicker.C:
			hub.CollectHistoric(ctx)
		case <-ctx.Done():
			return
		}
	}
}
