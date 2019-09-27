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

func NewCommStatCollector(period int64, store DataStore, options *config.CommunityStatOptions) (*Collector, error) {
	if period < 300 {
		log.Info("The minimum value for community stat collector interval is 300s(5m), setting interval to 300")
		period = 300
	}

	if period > 1800 {
		log.Info("The minimum value for community stat collector interval is 1800s(30m), setting interval to 1800")
		period = 1800
	}

	if len(options.Subreddit) == 0 {
		options.Subreddit = append(options.Subreddit, "decred")
	}

	subreddits = options.Subreddit
	twitterHandles = options.TwitterHandles
	repositories = options.GithubRepositories

	return &Collector{
		client:    http.Client{Timeout: 10 * time.Second},
		period:    time.Duration(period),
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

	for _, subreddit := range c.options.Subreddit {
		// reddit
		resp := new(RedditResponse)
		resp, err := c.fetchRedditStat(ctx, subreddit)
		for retry := 0; err != nil; retry++ {
			if retry == retryLimit {
				return err
			}
			log.Warn(err)
			resp, err = c.fetchRedditStat(ctx, subreddit)
		}

		err = c.dataStore.StoreRedditStat(ctx, Reddit{
			Date:           time.Now().UTC(),
			Subscribers:    resp.Data.Subscribers,
			AccountsActive: resp.Data.AccountsActive,
			Subreddit:      subreddit,
		})
		if err != nil {
			log.Error("Unable to save reddit stat, %s", err.Error())
			return err
		}
		log.Infof("New Reddit stat collected for %s at %s, Subscribers  %d, Active Users %d", subreddit,
			time.Now().Format(dateMiliTemplate), resp.Data.Subscribers, resp.Data.AccountsActive)
	}

	go c.startTwitterCollector(ctx)

	go c.startYoutubeCollector(ctx)

	// github
	go c.startGithubCollector(ctx)

	return nil
}
