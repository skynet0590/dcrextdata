// Copyright (c) 2018-2019 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package pow

import (
	"context"
	"net/http"
	"time"

	"github.com/raedahgroup/dcrextdata/app"
	"github.com/raedahgroup/dcrextdata/app/helpers"
	"github.com/raedahgroup/dcrextdata/datasync"
)

var (
	availablePows = []string{
		Coinmine,
		Luxor,
		F2pool,
		Uupool,
	}
)

type PowDataStore interface {
	PowTableName() string
	AddPowData(context.Context, []PowData) error
	AddPowDataFromSync(ctx context.Context, data PowData) error
	LastPowEntryTime(source string) (time int64)
	FetchPowDataForSync(ctx context.Context, date time.Time, skip, take int) ([]PowData, int64, error)
}

type Collector struct {
	pows   []Pow
	period int64
	store  PowDataStore
}

func NewCollector(disabledPows []string, period int64, store PowDataStore) (*Collector, error) {
	pows := make([]Pow, 0, len(availablePows)-len(disabledPows))
	disabledMap := make(map[string]struct{})
	for _, pow := range disabledPows {
		disabledMap[pow] = struct{}{}
	}

	for _, pow := range availablePows {
		if _, disabled := disabledMap[pow]; disabled {
			continue
		}

		if contructor, ok := PowConstructors[pow]; ok {
			lastEntryTime := store.LastPowEntryTime(pow)
			in, err := contructor(&http.Client{Timeout: 300 * time.Second}, lastEntryTime) // Consider if sharing a single client is better
			if err != nil {
				return nil, err
			}
			pows = append(pows, in)
		}
	}

	return &Collector{
		pows:   pows,
		period: period,
		store:  store,
	}, nil
}

func (pc *Collector) Run(ctx context.Context) {
	log.Info("Triggering PoW collectors.")

	lastCollectionDateUnix := pc.store.LastPowEntryTime("")
	lastCollectionDate := time.Unix(lastCollectionDateUnix, 0)
	secondsPassed := time.Since(lastCollectionDate)
	period := time.Duration(pc.period) * time.Second

	if lastCollectionDateUnix > 0 && secondsPassed < period {
		timeLeft := period - secondsPassed
		log.Infof("Fetching PoW data every %dm, collected %s ago, will fetch in %s.", pc.period/60,
			helpers.DurationToString(secondsPassed), helpers.DurationToString(timeLeft))

		time.Sleep(timeLeft)
	}
	// continually check the state of the app until its free to run this module
	for {
		if app.MarkBusyIfFree() {
			break
		}
	}
	pc.Collect(ctx)
	app.ReleaseForNewModule()
	go pc.CollectAsync(ctx)
}

func (pc *Collector) CollectAsync(ctx context.Context) {
	if ctx.Err() != nil {
		return
	}

	ticker := time.NewTicker(time.Duration(pc.period) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Infof("Stopping PoW collectors")
			return
		case <-ticker.C:
			// continually check the state of the app until its free to run this module
			for {
				if app.MarkBusyIfFree() {
					break
				}
			}
			completeCollectionCycle := pc.store.LastPowEntryTime("")
			collectionCycleDate := time.Unix(completeCollectionCycle, 0)
			timeInterval := time.Since(collectionCycleDate)
			log.Info("The next collection cycle begins in", timeInterval)

			log.Info("Starting a new PoW collection cycle")
			pc.Collect(ctx)
			app.ReleaseForNewModule()
		}
	}
}

func (pc *Collector) Collect(ctx context.Context) {
	log.Info("Fetching PoW data.")
	for _, powInfo := range pc.pows {
		select {
		case <-ctx.Done():
			return
		default:
			data, err := powInfo.Collect(ctx)
			if err != nil {
				log.Error(err)
			}
			err = pc.store.AddPowData(ctx, data)
			if err != nil {
				log.Error(err)
			}
		}
	}
}

func (pc *Collector) RegisterSyncer(syncCoordinator *datasync.SyncCoordinator) {
	syncCoordinator.AddSyncer(pc.store.PowTableName(), datasync.Syncer{
		Collect: func(ctx context.Context, url string) (result *datasync.Result, err error) {
			result = new(datasync.Result)
			result.Records = []PowData{}
			err = helpers.GetResponse(ctx, &http.Client{}, url, result)
			return
		},
		Retrieve: func(ctx context.Context, date time.Time, skip, take int) (result *datasync.Result, err error) {
			result = new(datasync.Result)
			powDatum, totalCount, err := pc.store.FetchPowDataForSync(ctx, date, skip, take)
			if err != nil {
				result.Message = err.Error()
				return
			}
			result.Records = powDatum
			result.TotalCount = totalCount
			result.Success = true
			return
		},
		Append: func(ctx context.Context, data interface{}) {
			mappedData := data.([]interface{})
			var powDataSlice []PowData
			for _, item := range mappedData {
				var powData PowData
				err := datasync.DecodeSyncObj(item, &powData)
				if err != nil {
					log.Errorf("Error in decoding the received PoW data, %s", err.Error())
					return
				}
				powDataSlice = append(powDataSlice, powData)
			}

			for _, powData := range powDataSlice {
				err := pc.store.AddPowDataFromSync(ctx, powData)
				if err != nil {
					log.Errorf("Error while appending PoW synced data, %s", err.Error())
				}
			}
		},
	})
}
