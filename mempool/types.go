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
	FirstSeenTime        time.Time
	NumberOfTransactions int
	Voters               uint32
	Size                 int32
	Fee 				 float64
	Total 				 float64
	BlockReceiveTime     time.Time
	BlockInternalTime    time.Time
	BlockHeight          uint32
	BlockHash            string
	Time                 time.Time
}

type DataStore interface {
	StoreMempool(context.Context, Mempool) error
	LastMempoolBlockHeight() (height int64, err error)
}

type Collector struct {
	dcrdClientConfig *rpcclient.ConnConfig
	dataStore        DataStore
}
