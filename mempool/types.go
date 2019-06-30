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
	Time                 time.Time `json:"time"`
	FirstSeenTime        time.Time `json:"first_seen_time"`
	NumberOfTransactions int       `json:"number_of_transactions"`
	Voters               int       `json:"voters"`
	Tickets              int       `json:"tickets"`
	Revocations          int       `json:"revocations"`
	Size                 int32     `json:"size"`
	TotalFee             float64   `json:"total_fee"`
	Total                float64   `json:"total"`
}

type MempoolDto struct {
	Time                 string `json:"time"`
	FirstSeenTime        string `json:"first_seen_time"`
	NumberOfTransactions int       `json:"number_of_transactions"`
	Voters               int       `json:"voters"`
	Tickets              int       `json:"tickets"`
	Revocations          int       `json:"revocations"`
	Size                 int32     `json:"size"`
	TotalFee             float64   `json:"total_fee"`
	Total                float64   `json:"total"`
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
	activeChain      *chaincfg.Params
}
