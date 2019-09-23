// Copyright (c) 2018-2019 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package commstats

import (
	"context"
	"net/http"
	"time"
)

type CommStat struct {
	Date                 time.Time `json:"date"`
	RedditSubscribers    int       `json:"reddit_subscribers"`
	RedditAccountsActive int       `json:"reddit_accounts_active"`
	TwitterFollowers     int       `json:"twitter_followers"`
	YoutubeSubscribers   int       `json:"youtube_subscribers"`
	GithubStars          int       `json:"github_stars"`
	GithubFolks          int       `json:"github_folks"`
}

type RedditResponse struct {
	Kind string `json:"kind"`
	Data struct {
		Subscribers    int `json:"subscribers"`
		AccountsActive int `json:"active_user_count"`
	} `json:"data"`
}

type DataStore interface {
	StoreCommStat(context.Context, CommStat) error
	LastCommStatEntry() (time time.Time)
}

type Collector struct {
	client    http.Client
	period    time.Duration
	dataStore DataStore
}
