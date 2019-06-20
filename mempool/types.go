// Copyright (c) 2018-2019 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package mempool

import (
	"context"
	"github.com/decred/dcrd/rpcclient"
)

type Mempool struct {
	FirstSeenTime           int64
	BlockReceiveTime        int64
	TotalSent               float64
	LastBlockHeight         int64
	Size                    int32
	RegularTransactionCount int
	TicketsCount            uint8
	VoteCount               uint16
	RevocationCount         uint8
}

type DataStore interface {
	StoreMempool(context.Context, Mempool) error
	LastMempool(ctx context.Context) (Mempool, error)
}

type Collector struct {
	dcrdClientConfig  *rpcclient.ConnConfig
	dataStore     DataStore
}