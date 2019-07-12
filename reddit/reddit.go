// Copyright (c) 2018-2019 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package reddit

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/raedahgroup/dcrextdata/app"
	"github.com/raedahgroup/dcrextdata/app/helpers"
)

const (
	requestURL = "https://www.reddit.com/r/decred/about.json"
	retryLimit = 3
)

func NewRedditCollector(period int64, store DataStore) (*Collector, error) {
	request, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}

	return &Collector{
		client:    http.Client{Timeout: 10 * time.Second},
		period:    time.Duration(period),
		request:   request,
		dataStore: store,
	}, nil
}

func (c *Collector) fetch(ctx context.Context, response interface{}) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	// log.Tracef("GET %v", requestURL)
	resp, err := c.client.Do(c.request.WithContext(ctx))
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

func (c *Collector) Run(ctx context.Context) {
	if ctx.Err() != nil {
		return
	}

	lastCollectionDate := c.dataStore.LastRedditEntryTime()
	secondsPassed := time.Since(lastCollectionDate)
	period := c.period * time.Second

	log.Info("Starting VSP collection cycle.")

	if secondsPassed < period {
		timeLeft := period - secondsPassed
		//Fetching VSPs every 5m, collected 1m7.99s ago, will fetch in 3m52.01s
		log.Infof("Fetching reddit data every %dm, collected %s ago, will fetch in %s.", c.period/60, helpers.DurationToString(secondsPassed),
			helpers.DurationToString(timeLeft))

		time.Sleep(timeLeft)
	}

	// continually check the state of the app until its free to run this module
	for {
		if app.MarkBusyIfFree() {
			break
		}
	}

	log.Info("Fetching RedditInfo data...")

	err := c.collectAndStore(ctx)
	app.ReleaseForNewModule()
	if err != nil {
		log.Errorf("Could not start collection: %v", err)
		return
	}

	ticker := time.NewTicker(c.period * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Infof("Shutting down reddit data collector")
			return
		case <-ticker.C:
			// continually check the state of the app until its free to run this module
			for {
				if app.MarkBusyIfFree() {
					break
				}
			}
			log.Info("Fetching RedditInfo data...")
			err := c.collectAndStore(ctx)
			app.ReleaseForNewModule()
			if err != nil {
				return
			}
		}
	}
}

func (c *Collector) collectAndStore(ctx context.Context) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	resp := new(Response)
	err := c.fetch(ctx, resp)
	for retry := 0; err != nil; retry++ {
		if retry == retryLimit {
			return err
		}
		log.Warn(err)
		err = c.fetch(ctx, resp)
	}

	redditData := RedditInfo{
		Date: time.Now(),
		Subscribers: resp.Data.Subscribers,
		AccountsActive: resp.Data.AccountsActive,
	}
	// log.Infof("Collected data for %d vsps", len(*resp))

	return c.dataStore.StoreRedditData(ctx, redditData)
}
