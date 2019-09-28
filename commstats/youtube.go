package commstats

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/raedahgroup/dcrextdata/app"
	"github.com/raedahgroup/dcrextdata/app/helpers"
	"github.com/raedahgroup/dcrextdata/postgres/models"
)

const (
	youtubeChannelId = "UCJ2bYDaPYHpSmJPh_M5dNSg"
)

func (c *Collector) startYoutubeCollector(ctx context.Context) {
	if c.options.YoutubeDataApiKey == "" {
		log.Error("youtubedataapikey is required for the youtube stat collector to work")
		return
	}
	var lastCollectionDate time.Time
	err := c.dataStore.LastEntry(ctx, models.TableNames.Youtube, &lastCollectionDate)
	if err != nil && err != sql.ErrNoRows {
		log.Errorf("Cannot fetch last Youtube entry time, %s", err.Error())
		return
	}

	secondsPassed := time.Since(lastCollectionDate)
	period := time.Duration(c.options.TwitterStatInterval) * time.Minute

	if secondsPassed < period {
		timeLeft := period - secondsPassed
		log.Infof("Fetching Youtube stats every %dm, collected %s ago, will fetch in %s.", period/time.Minute, helpers.DurationToString(secondsPassed),
			helpers.DurationToString(timeLeft))

		time.Sleep(timeLeft)
	}

	registerStarter := func() {
		// continually check the state of the app until its free to run this module
		for {
			if app.MarkBusyIfFree() {
				break
			}
		}
	}

	registerStarter()
	c.collectAndStoreYoutubeStat(ctx)
	app.ReleaseForNewModule()

	ticker := time.NewTicker(time.Duration(c.options.YoutubeStatInterval) * time.Minute)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			registerStarter()
			c.collectAndStoreYoutubeStat(ctx)
			app.ReleaseForNewModule()
		}
	}
}

func (c *Collector) collectAndStoreYoutubeStat(ctx context.Context) {
	log.Info("Starting Github stats collection cycle")
	// youtube
	youtubeSubscribers, err := c.getYoutubeSubscriberCount(ctx)
	for retry := 0; err != nil; retry++ {
		if retry == retryLimit {
			return
		}
		log.Warn(err)
		youtubeSubscribers, err = c.getYoutubeSubscriberCount(ctx)
	}

	youtubeStat := Youtube{
		Date:        time.Now().UTC(),
		Subscribers: youtubeSubscribers,
	}
	err = c.dataStore.StoreYoutubeStat(ctx, youtubeStat)
	if err != nil {
		log.Error("Unable to save Youtube stat, %s", err.Error())
		return
	}

	log.Infof("New Youtube stat collected at %s, Subscribers %d",
		youtubeStat.Date.Format(dateMiliTemplate), youtubeSubscribers)
}

func (c *Collector) getYoutubeSubscriberCount(ctx context.Context) (int, error) {
	if ctx.Err() != nil {
		return 0, ctx.Err()
	}

	youtubeUrl := fmt.Sprintf("https://content.googleapis.com/youtube/v3/channels?key=%s&part=statistics&id=%s",
		c.options.YoutubeDataApiKey, youtubeChannelId)

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
		Items []struct {
			Statistics struct {
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
