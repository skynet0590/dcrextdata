// Copyright (c) 2018-2019 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package mempool

import (
	"context"
	"github.com/decred/dcrd/rpcclient"
	"time"
)

type Mempool struct {
	Time                 time.Time
	FirstSeenTime        time.Time
	NumberOfTransactions int
	Voters               int
	Tickets              int
	Revocations          int
	Size                 int32
	TotalFee             float64
	Total                float64

	BlockReceiveTime     time.Time
	BlockInternalTime    time.Time
	BlockHeight          uint32
	BlockHash            string
}

type DataStore interface {
	StoreMempool(context.Context, Mempool) error
	LastMempoolBlockHeight() (height int64, err error)
}

type Collector struct {
	dcrdClientConfig *rpcclient.ConnConfig
	dataStore        DataStore
}
