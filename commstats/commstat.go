// Copyright (c) 2018-2019 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package commstats

import (
	"context"
	"net/http"
	"time"

	"github.com/raedahgroup/dcrextdata/app"
	"github.com/raedahgroup/dcrextdata/app/config"
)

const (
	dateTemplate     = "2006-01-02 15:04"
	dateMiliTemplate = "2006-01-02 15:04:05.99"
	retryLimit       = 3
)

func NewCommStatCollector(store DataStore, options *config.CommunityStatOptions) (*Collector, error) {

	if len(options.Subreddit) == 0 {
		options.Subreddit = append(options.Subreddit, "decred")
	}

	subreddits = options.Subreddit
	twitterHandles = options.TwitterHandles
	repositories = options.GithubRepositories

	return &Collector{
		client:    http.Client{Timeout: 10 * time.Second},
		dataStore: store,
		options:   options,
	}, nil
}

func (c *Collector) Run(ctx context.Context) {
	if ctx.Err() != nil {
		return
	}

	/*lastCollectionDate := c.dataStore.LastCommStatEntry()
	secondsPassed := time.Since(lastCollectionDate)
	period := c.period * time.Second

	log.Info("Starting community stats collection cycle.")

	if secondsPassed < period {
		timeLeft := period - secondsPassed
		log.Infof("Fetching community stats every %dm, collected %s ago, will fetch in %s.", c.period/60,
			helpers.DurationToString(secondsPassed),
			helpers.DurationToString(timeLeft))

		time.Sleep(timeLeft)
	}*/

	// continually check the state of the app until its free to run this module
	for {
		if app.MarkBusyIfFree() {
			break
		}
	}

	log.Info("Fetching community stats...")

	go c.startTwitterCollector(ctx)

	go c.startYoutubeCollector(ctx)

	// github
	go c.startGithubCollector(ctx)

	go c.startRedditCollector(ctx)
}
