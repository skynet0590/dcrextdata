// Copyright (c) 2018-2019 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package reddit

import (
	"context"
	"net/http"
	"time"
)

type RedditInfo struct {
	Date           time.Time `json:"date"`
	Subscribers    int `json:"subscribers"`
	AccountsActive int `json:"accounts_active"`
}

type Response struct {
	Kind string `json:"kind"`
	Data struct{
		Subscribers    int `json:"subscribers"`
		AccountsActive int `json:"active_user_count"`
	} `json:"data"`
}

type DataStore interface {
	StoreRedditData(context.Context, RedditInfo) error
	LastRedditEntryTime() (time time.Time)
}

type Collector struct {
	client    http.Client
	period    time.Duration
	request   *http.Request
	dataStore DataStore
}
