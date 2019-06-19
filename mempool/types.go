// Copyright (c) 2018-2019 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package mempool

import (
	"context"
	"time"

	"github.com/decred/dcrd/rpcclient"
)

type Mempool struct {
	Date time.Time
	TotalSent float64
	LastBlockHeight int64
	Size float64
	RegularTransactionCount int
	TicketsCount int
	VoteCount int
	RevocationCount int
}

type DataStore interface {
	StoreMempool(context.Context, Mempool) []error
	LastMempool(ctx context.Context) (Mempool, error)
}

type Collector struct {
	dcrdClient  *rpcclient.Client
	dataStore     DataStore
}