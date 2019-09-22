// Copyright (c) 2018-2019 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package commstats

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
	redditRequestURL = "https://www.reddit.com/r/decred/about.json"
	retryLimit       = 3
)

func NewCommStatCollector(period int64, store DataStore) (*Collector, error) {
	if period < 300 {
		log.Info("The minimum value for community stat collector interval is 300s(5m), setting interval to 300")
		period = 300
	}

	if period > 1800{
		log.Info("The minimum value for community stat collector interval is 1800s(30m), setting interval to 1800")
		period = 1800
	}

	return &Collector{
		client:        http.Client{Timeout: 10 * time.Second},
		period:        time.Duration(period),
		dataStore:     store,
	}, nil
}

func (c *Collector) fetchRedditStat(ctx context.Context, response *RedditResponse) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	request, err := http.NewRequest(http.MethodGet, redditRequestURL, nil)
	if err != nil {
		return err
	}

	// reddit returns too many redditRequest http status if user agent is not set
	request.Header.Set("user-agent",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Safari/537.36")

	// log.Tracef("GET %v", redditRequestURL)
	resp, err := c.client.Do(request.WithContext(ctx))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		err = json.NewDecoder(resp.Body).Decode(response)
		if err != nil {
			return fmt.Errorf(fmt.Sprintf("Failed to decode json: %v", err))
		}
	} else {
		log.Infof("Unable to fetchRedditStat data from reddit: %s", resp.Status)
	}

	return nil
}

func (c *Collector) Run(ctx context.Context) {
	if ctx.Err() != nil {
		return
	}

	lastCollectionDate := c.dataStore.LastCommStatEntry()
	lastCollectionDate = lastCollectionDate.Add(time.Hour) // todo: this need justification
	secondsPassed := time.Since(lastCollectionDate)
	period := c.period * time.Second

	log.Info("Starting community stats collection cycle.")

	if secondsPassed < period {
		timeLeft := period - secondsPassed
		log.Infof("Fetching community stats every %dm, collected %s ago, will fetch in %s.", c.period/60,
			helpers.DurationToString(secondsPassed),
			helpers.DurationToString(timeLeft))

		time.Sleep(timeLeft)
	}

	// continually check the state of the app until its free to run this module
	for {
		if app.MarkBusyIfFree() {
			break
		}
	}

	log.Info("Fetching community stats...")

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
			log.Infof("Shutting down community stats collector")
			return
		case <-ticker.C:
			// continually check the state of the app until its free to run this module
			for {
				if app.MarkBusyIfFree() {
					break
				}
			}
			log.Info("Fetching community stats...")
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

	resp := new(RedditResponse)
	err := c.fetchRedditStat(ctx, resp)
	for retry := 0; err != nil; retry++ {
		if retry == retryLimit {
			return err
		}
		log.Warn(err)
		err = c.fetchRedditStat(ctx, resp)
	}

	redditData := CommStat{
		Date:                 time.Now().UTC(),
		RedditSubscribers:    resp.Data.Subscribers,
		RedditAccountsActive: resp.Data.AccountsActive,
	}
	return c.dataStore.StoreCommStat(ctx, redditData)
}
