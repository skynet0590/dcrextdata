// Copyright (c) 2018-2019 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package commstats

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/raedahgroup/dcrextdata/app/config"
	"net/http"
	"strconv"
	"time"

	"github.com/raedahgroup/dcrextdata/app"
)

const (
	dateTemplate = "2006-01-02 15:04"
	dateMiliTemplate = "2006-01-02 15:04:05.99"
	twitterRequestURL = "https://cdn.syndication.twimg.com/widgets/followbutton/info.json?screen_names=decredproject"
	youtubeChannelId  = "UCJ2bYDaPYHpSmJPh_M5dNSg"
	youtubeKey		  = "AIzaSyBmUyMNtZUqReP2NTs39UlTjd9aUjXWKq0"
	retryLimit        = 3
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

	return &Collector{
		client:    http.Client{Timeout: 10 * time.Second},
		period:    time.Duration(period),
		dataStore: store,
		options: options,
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
			Date:           time.Now(),
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

	// twitter
	followers, err := c.getTwitterFollowers(ctx)
	for retry := 0; err != nil; retry++ {
		if retry == retryLimit {
			return err
		}
		log.Warn(err)
		followers, err = c.getTwitterFollowers(ctx)
	}

	var twitterStat = Twitter{Date:time.Now(), Followers:followers}
	err = c.dataStore.StoreTwitterStat(ctx, twitterStat)
	if err != nil {
		log.Error("Unable to save twitter stat, %s", err.Error())
		return err
	}

	log.Infof("New Twitter stat collected at %s, Followers %d",
		twitterStat.Date.Format(dateMiliTemplate), twitterStat.Followers)

	// youtube
	youtubeSubscribers, err := c.getYoutubeSubscriberCount(ctx)
	for retry := 0; err != nil; retry++ {
		if retry == retryLimit {
			return err
		}
		log.Warn(err)
		youtubeSubscribers, err = c.getYoutubeSubscriberCount(ctx)
	}

	youtubeStat := Youtube{
		Date:        time.Now(),
		Subscribers: youtubeSubscribers,
	}
	err = c.dataStore.StoreYoutubeStat(ctx, youtubeStat)
	if err != nil {
		log.Error("Unable to save Youtube stat, %s", err.Error())
		return err
	}

	log.Infof("New Youtube stat collected at %s, Subscribers %d",
		youtubeStat.Date.Format(dateMiliTemplate), youtubeSubscribers)

	// github
	githubStars, githubFolks, err := c.getGithubData(ctx)
	for retry := 0; err != nil; retry++ {
		if retry == retryLimit {
			return err
		}
		log.Warn(err)
		githubStars, githubFolks, err = c.getGithubData(ctx)
	}

	githubStat := Github{
		Date:  time.Now(),
		Stars: githubStars,
		Folks: githubFolks,
	}
	err = c.dataStore.StoreGithubStat(ctx, githubStat)
	if err != nil {
		log.Error("Unable to save Github stat, %s", err.Error())
		return err
	}

	log.Infof("New Github stat collected at %s, Stars %d, Folks %d",
		githubStat.Date.Format(dateMiliTemplate), githubStars, githubFolks)

	return nil
}

func (c *Collector) getTwitterFollowers(ctx context.Context) (int, error) {
	if ctx.Err() != nil {
		return 0, ctx.Err()
	}

	request, err := http.NewRequest(http.MethodGet, twitterRequestURL, nil)
	if err != nil {
		return 0, err
	}

	// reddit returns too many redditRequest http status if user agent is not set
	request.Header.Set("user-agent",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Safari/537.36")

	// log.Tracef("GET %v", redditRequestURL)
	resp, err := c.client.Do(request.WithContext(ctx))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	var response []struct {
		Followers int `json:"followers_count"`
	}
	if resp.StatusCode == http.StatusOK {
		err = json.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			return 0, fmt.Errorf(fmt.Sprintf("Failed to decode json: %v", err))
		}
	} else {
		return 0, fmt.Errorf("unable to fetch twitter followers: %s", resp.Status)
	}

	if len(response) < 1 {
		return 0, errors.New("unable to fetch twitter followers, no response")
	}

	return response[0].Followers, nil
}

func (c *Collector) getYoutubeSubscriberCount(ctx context.Context) (int, error) {
	if ctx.Err() != nil {
		return 0, ctx.Err()
	}

	youtubeUrl := fmt.Sprintf("https://content.googleapis.com/youtube/v3/channels?key=%s&part=statistics&id=%s",
		youtubeKey, youtubeChannelId)

	request, err := http.NewRequest(http.MethodGet, youtubeUrl, nil)
	if err != nil {
		return 0, err
	}

	request.Header.Set("user-agent",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Safari/537.36")

	resp, err := c.client.Do(request.WithContext(ctx))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	var response struct {
		Items []struct{
			Statistics struct{
				SubscriberCount string `json:"subscriberCount"`
			} `json:"statistics"`
		} `json:"items"`
	}
	if resp.StatusCode == http.StatusOK {
		err = json.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			return 0, fmt.Errorf(fmt.Sprintf("Failed to decode json: %v", err))
		}
	} else {
		return 0, fmt.Errorf("unable to fetch youtube subscribers: %s", resp.Status)
	}

	if len(response.Items) < 1 {
		return 0, errors.New("unable to fetch youtube subscribers, no response")
	}

	subscribers, err := strconv.Atoi(response.Items[0].Statistics.SubscriberCount)
	if err != nil {
		return 0, errors.New("unable to fetch youtube subscribers, no response")
	}

	return subscribers, nil
}

func (c *Collector) getGithubData(ctx context.Context) (int, int, error) {
	if ctx.Err() != nil {
		return 0, 0, ctx.Err()
	}

	request, err := http.NewRequest(http.MethodGet, "https://api.github.com/repos/decred/dcrd", nil)
	if err != nil {
		return 0, 0, err
	}

	request.Header.Set("user-agent",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Safari/537.36")

	resp, err := c.client.Do(request.WithContext(ctx))
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()
	var response struct {
		Stars int `json:"stargazers_count"`
		Folks int `json:"forks_count"`
	}
	if resp.StatusCode == http.StatusOK {
		err = json.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			return 0, 0, fmt.Errorf(fmt.Sprintf("Failed to decode json: %v", err))
		}
	} else {
		return 0, 0, fmt.Errorf("unable to fetch youtube subscribers: %s", resp.Status)
	}

	return response.Stars, response.Folks, nil
}