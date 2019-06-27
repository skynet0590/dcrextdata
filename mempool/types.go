// Copyright (c) 2018-2019 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package mempool

import (
	"context"
	"time"

	"github.com/decred/dcrd/chaincfg"
	"github.com/decred/dcrd/rpcclient"
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
}

type Block struct {
	BlockReceiveTime  time.Time
	BlockInternalTime time.Time
	BlockHeight       uint32
	BlockHash         string
}

type Vote struct {
	Hash        string
	ReceiveTime time.Time
	VotingOn    int64
	ValidatorId int
}

type DataStore interface {
	StoreMempool(context.Context, Mempool) error
	SaveBlock(context.Context, Block) error
	SaveVote(ctx context.Context, vote Vote) error
}

type Collector struct {
	dcrdClientConfig *rpcclient.ConnConfig
	dataStore        DataStore
	activeChain  *chaincfg.Params
}
