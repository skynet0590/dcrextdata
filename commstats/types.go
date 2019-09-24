// Copyright (c) 2018-2019 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package commstats

import (
	"context"
	"github.com/raedahgroup/dcrextdata/app/config"
	"net/http"
	"time"
)

type CommStat struct {
	Date               time.Time             `json:"date"`
	RedditStats        map[string]RedditStat `json:"reddit_stats"`
	TwitterFollowers   int                   `json:"twitter_followers"`
	YoutubeSubscribers int                   `json:"youtube_subscribers"`
	GithubStars        int                   `json:"github_stars"`
	GithubFolks        int                   `json:"github_folks"`
}

type RedditResponse struct {
	Kind string `json:"kind"`
	Data RedditStat `json:"data"`
}

type RedditStat struct {
	Subscribers    int `json:"subscribers"`
	AccountsActive int `json:"active_user_count"`
}

type DataStore interface {
	StoreCommStat(context.Context, CommStat) error
	LastCommStatEntry() (time time.Time)
}

type Collector struct {
	client    http.Client
	period    time.Duration
	dataStore DataStore
	options   *config.CommunityStatOptions
}
