// Copyright (c) 2018-2019 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package commstats

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/raedahgroup/dcrextdata/app"
)

const (
	redditRequestURL  = "https://www.reddit.com/r/decred/about.json"
	twitterRequestURL = "https://cdn.syndication.twimg.com/widgets/followbutton/info.json?screen_names=decredproject"
	youtubeChannelId  = "UCJ2bYDaPYHpSmJPh_M5dNSg"
	youtubeKey		  = "AIzaSyBmUyMNtZUqReP2NTs39UlTjd9aUjXWKq0"
	retryLimit        = 3
)

func NewCommStatCollector(period int64, store DataStore) (*Collector, error) {
	if period < 300 {
		log.Info("The minimum value for community stat collector interval is 300s(5m), setting interval to 300")
		period = 300
	}

	if period > 1800 {
		log.Info("The minimum value for community stat collector interval is 1800s(30m), setting interval to 1800")
		period = 1800
	}

	return &Collector{
		client:    http.Client{Timeout: 10 * time.Second},
		period:    time.Duration(period),
		dataStore: store,
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

	// reddit
	resp := new(RedditResponse)
	err := c.fetchRedditStat(ctx, resp)
	for retry := 0; err != nil; retry++ {
		if retry == retryLimit {
			return err
		}
		log.Warn(err)
		err = c.fetchRedditStat(ctx, resp)
	}

	stat := CommStat{
		Date:                 time.Now().UTC(),
		RedditSubscribers:    resp.Data.Subscribers,
		RedditAccountsActive: resp.Data.AccountsActive,
	}

	// twitter
	stat.TwitterFollowers, err = c.getTwitterFollowers(ctx)
	for retry := 0; err != nil; retry++ {
		if retry == retryLimit {
			return err
		}
		log.Warn(err)
		stat.TwitterFollowers, err = c.getTwitterFollowers(ctx)
	}

	// youtube
	stat.YoutubeSubscribers, err = c.getYoutubeSubscriberCount(ctx)
	for retry := 0; err != nil; retry++ {
		if retry == retryLimit {
			return err
		}
		log.Warn(err)
		stat.YoutubeSubscribers, err = c.getYoutubeSubscriberCount(ctx)
	}

	// github
	stat.GithubStars, stat.GithubFolks, err = c.getGithubData(ctx)
	for retry := 0; err != nil; retry++ {
		if retry == retryLimit {
			return err
		}
		log.Warn(err)
		stat.GithubStars, stat.GithubFolks, err = c.getGithubData(ctx)
	}

	return c.dataStore.StoreCommStat(ctx, stat)
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