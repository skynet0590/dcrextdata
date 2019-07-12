// Copyright (c) 2018-2019 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package reddit

import (
	"context"
	"net/http"
	"time"
)

type Reddit struct {
	Date           time.Time `json:"date"`
	Subscribers    int64 `json:"subscribers"`
	AccountsActive int64 `json:"accounts_active"`
}

type Response struct {
	Kind string `json:"kind"`
	Data struct{
		Subscribers    int64 `json:"subscribers"`
		AccountsActive int64 `json:"accounts_active"`
	} `json:"data"`
}

type DataStore interface {
	StoreRedditData(context.Context, Reddit) error
	LastRedditEntryTime() (time time.Time)
}

type Collector struct {
	client    http.Client
	period    time.Duration
	request   *http.Request
	dataStore DataStore
}
