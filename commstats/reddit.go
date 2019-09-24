package commstats

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	redditRequestURL  = "https://www.reddit.com/r/%s/about.json"
)

func (c *Collector) fetchRedditStat(ctx context.Context, subreddit string) (response *RedditResponse, err error) {
	if ctx.Err() != nil {
		err = ctx.Err()
		return
	}

	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf(redditRequestURL, subreddit), nil)
	if err != nil {
		return
	}

	// reddit returns too many redditRequest http status if user agent is not set
	request.Header.Set("user-agent",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Safari/537.36")

	// log.Tracef("GET %v", redditRequestURL)
	resp, err := c.client.Do(request.WithContext(ctx))
	if err != nil {
		return
	}
	defer resp.Body.Close()
	response = new(RedditResponse)
	if resp.StatusCode == http.StatusOK {
		err = json.NewDecoder(resp.Body).Decode(response)
		if err != nil {
			return nil, fmt.Errorf(fmt.Sprintf("Failed to decode json: %v", err))
		}
	} else {
		log.Infof("Unable to fetchRedditStat data from reddit: %s", resp.Status)
	}

	return
}
