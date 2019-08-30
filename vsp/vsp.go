// Copyright (c) 2018-2019 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package vsp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/raedahgroup/dcrextdata/app"
	"github.com/raedahgroup/dcrextdata/app/helpers"
	"github.com/raedahgroup/dcrextdata/datasync"
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

	return &Collector{
		client:    http.Client{Timeout: time.Minute},
		period:    time.Duration(period),
		request:   request,
		dataStore: store,
	}, nil
}

func (vsp *Collector) fetch(ctx context.Context, response interface{}) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	// log.Tracef("GET %v", requestURL)
	resp, err := vsp.client.Do(vsp.request.WithContext(ctx))
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

func (vsp *Collector) Run(ctx context.Context) {
	if ctx.Err() != nil {
		return
	}

	lastCollectionDate := vsp.dataStore.LastVspTickEntryTime()
	secondsPassed := time.Since(lastCollectionDate)
	period := vsp.period * time.Second

	log.Info("Starting VSP collection cycle.")

	if secondsPassed < period {
		timeLeft := period - secondsPassed
		log.Infof("Fetching VSPs every %dm, collected %s ago, will fetch in %s.", vsp.period/60, helpers.DurationToString(secondsPassed),
			helpers.DurationToString(timeLeft))

		time.Sleep(timeLeft)
	}

	// continually check the state of the app until its free to run this module
	for {
		if app.MarkBusyIfFree() {
			break
		}
	}

	err := vsp.collectAndStore(ctx)
	app.ReleaseForNewModule()
	if err != nil {
		log.Errorf("Could not start collection: %s", err.Error())
		return
	}

	go func() {
		ticker := time.NewTicker(vsp.period * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Infof("Shutting down VSP collector")
				return
			case <-ticker.C:
				// continually check the state of the app until its free to run this module
				for {
					if app.MarkBusyIfFree() {
						break
					}
				}
				err := vsp.collectAndStore(ctx)
				app.ReleaseForNewModule()
				if err != nil {
					return
				}
			}
		}
	}()
}

func (vsp *Collector) collectAndStore(ctx context.Context) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	log.Info("Fetching VSP from source")

	resp := new(Response)
	err := vsp.fetch(ctx, resp)
	for retry := 0; err != nil; retry++ {
		if retry == retryLimit {
			return err
		}
		log.Warn(err)
		err = vsp.fetch(ctx, resp)
	}

	numberStored, errs := vsp.dataStore.StoreVSPs(ctx, *resp)
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

	log.Infof("Saved ticks for %d VSPs from %s", numberStored, requestURL)
	return nil
}

func (vsp *Collector) RegisterSyncer(syncCoordinator *datasync.SyncCoordinator) {
	vsp.registerVspSyncer(syncCoordinator)
	vsp.registerVspTickSyncer(syncCoordinator)
}

func (vsp *Collector) registerVspSyncer(syncCoordinator *datasync.SyncCoordinator) {
	syncCoordinator.AddSyncer(vsp.dataStore.VspTableName(), datasync.Syncer{
		Collect: func(ctx context.Context, url string) (result *datasync.Result, err error) {
			result = new(datasync.Result)
			result.Records = []VSPDto{}
			err = helpers.GetResponse(ctx, &http.Client{}, url, result)
			return
		},
		Retrieve: func(ctx context.Context, date time.Time, skip, take int) (result *datasync.Result, err error) {
			result = new(datasync.Result)
			vspSources, totalCount, err := vsp.dataStore.FetchVspSourcesForSync(ctx, date, skip, take)
			if err != nil {
				result.Message = err.Error()
				return
			}
			result.Records = vspSources
			result.TotalCount = totalCount
			result.Success = true
			return
		},
		Append: func(ctx context.Context, data interface{}) {
			vspSources := data.([]VSPDto)
			for _, vspSource := range vspSources {
				err := vsp.dataStore.AddVspSourceFromSync(ctx, vspSource)
				if err != nil {
					log.Errorf("Error while appending vsp source synced data, %s", err.Error())
				}
			}
		},
	})
}
func (vsp *Collector) registerVspTickSyncer(syncCoordinator *datasync.SyncCoordinator) {
	syncCoordinator.AddSyncer(vsp.dataStore.VspTickTableName(), datasync.Syncer{
		Collect: func(ctx context.Context, url string) (result *datasync.Result, err error) {
			result = new(datasync.Result)
			result.Records = []VSPTickSyncDto{}
			err = helpers.GetResponse(ctx, &http.Client{}, url, result)
			return
		},
		Retrieve: func(ctx context.Context, date time.Time, skip, take int) (result *datasync.Result, err error) {
			result = new(datasync.Result)
			vspTicks, totalCount, err := vsp.dataStore.FetchVspTicksForSync(ctx, date, skip, take)
			if err != nil {
				result.Message = err.Error()
				return
			}
			result.Records = vspTicks
			result.TotalCount = totalCount
			result.Success = true
			return
		},
		Append: func(ctx context.Context, data interface{}) {
			vspTicks := data.([]VSPTickSyncDto)
			for _, tick := range vspTicks {
				err := vsp.dataStore.AddVspTicksFromSync(ctx, tick)
				if err != nil {
					log.Errorf("Error while appending vsp tick synced data, %s", err.Error())
				}
			}
		},
	})
}
