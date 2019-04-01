// Copyright (c) 2018-2019 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package vsp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

const (
	requestURL = "https://api.decred.org/?c=gsd"
	retryLimit = 3
)

func NewVspCollector(period int64, store DataStore) (*Collector, error) {
	request, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}
	return &Collector{
		client:    &http.Client{Timeout: 300 * time.Second},
		period:    time.Duration(period),
		request:   request,
		dataStore: store,
	}, nil
}

func (vsp *Collector) fetch(response interface{}) error {
	log.Tracef("GET %v", requestURL)
	resp, err := vsp.client.Do(vsp.request)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(response)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("Failed to decode json: %v", err))
	}

	return nil
}

func (vsp *Collector) Run(quit chan struct{}, wg *sync.WaitGroup) {
	if err := vsp.CollectAndStore(time.Now()); err != nil {
		log.Errorf("Could not start collection: %v", err)
	}

	ticker := time.NewTicker(vsp.period * time.Second)

	defer func(wg *sync.WaitGroup) {
		log.Info("Stopping collector")
		ticker.Stop()
		wg.Done()
	}(wg)

	for {
		select {
		case t := <-ticker.C:
			err := vsp.CollectAndStore(t)
			if err != nil {
				log.Error(err)
			}
		case <-quit:
			return
		}
	}
}

func (vsp *Collector) CollectAndStore(t time.Time) error {
	resp := new(Response)
	err := vsp.fetch(resp)
	for retry := 0; err != nil; retry++ {
		if retry == retryLimit {
			return err
		}
		log.Warn(err)
		err = vsp.fetch(resp)
	}

	if resp != nil {
		errs := vsp.dataStore.StoreVSPs(*resp)
		for _, err = range errs {
			if err != nil {
				if e, ok := err.(PoolTickTimeExistsError); ok {
					log.Trace(e)
				} else {
					log.Error(err)
					return err
				}
			}
		}
	}
	return nil
}
