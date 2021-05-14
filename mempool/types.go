// Copyright (c) 2018-2019 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package mempool

import (
	"context"
	"time"

	"github.com/decred/dcrd/chaincfg"
	"github.com/decred/dcrd/rpcclient"
	"github.com/planetdecred/dcrextdata/datasync"
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

type Dto struct {
	Time                 string  `json:"time"`
	FirstSeenTime        string  `json:"first_seen_time"`
	NumberOfTransactions int     `json:"number_of_transactions"`
	Voters               int     `json:"voters"`
	Tickets              int     `json:"tickets"`
	Revocations          int     `json:"revocations"`
	Size                 int32   `json:"size"`
	TotalFee             float64 `json:"total_fee"`
	Total                float64 `json:"total"`
}

type Block struct {
	BlockReceiveTime  time.Time
	BlockInternalTime time.Time
	BlockHeight       uint32
	BlockHash         string
}

type BlockDto struct {
	BlockReceiveTime  string    `json:"block_receive_time"`
	BlockInternalTime string    `json:"block_internal_time"`
	Delay             string    `json:"delay"`
	BlockHeight       uint32    `json:"block_height"`
	BlockHash         string    `json:"block_hash"`
	Votes             []VoteDto `json:"votes"`
}

type Vote struct {
	Hash              string
	ReceiveTime       time.Time
	TargetedBlockTime time.Time
	BlockReceiveTime  time.Time
	VotingOn          int64
	BlockHash         string
	ValidatorId       int
	Validity          string
}

type VoteDto struct {
	Hash                  string `json:"hash"`
	ReceiveTime           string `json:"receive_time"`
	TargetedBlockTimeDiff string `json:"block_time_diff"`
	BlockReceiveTimeDiff  string `json:"block_receive_time_diff"`
	VotingOn              int64  `json:"voting_on"`
	BlockHash             string `json:"block_hash"`
	ShortBlockHash        string `json:"short_block_hash"`
	ValidatorId           int    `json:"validator_id"`
	Validity              string `json:"validity"`
}

type PropagationChartData struct {
	BlockHeight    int64     `json:"block_height"`
	TimeDifference float64   `json:"time_difference"`
	BlockTime      time.Time `json:"block_time"`
}

type BlockReceiveTime struct {
	BlockHeight int64 `json:"block_height"`
	ReceiveTime time.Time
}

type DataStore interface {
	MempoolTableName() string
	BlockTableName() string
	VoteTableName() string
	StoreMempool(context.Context, Mempool) error
	LastMempoolTime() (entryTime time.Time, err error)
	FetchMempoolForSync(ctx context.Context, date time.Time, offtset int, limit int) ([]Mempool, int64, error)
	SaveBlock(context.Context, Block) error
	UpdateBlockBinData(context.Context) error
	FetchBlockForSync(ctx context.Context, blockHeight int64, offtset int, limit int) ([]Block, int64, error)
	SaveVote(ctx context.Context, vote Vote) error
	UpdateVoteTimeDeviationData(context.Context) error
	FetchVoteForSync(ctx context.Context, date time.Time, offtset int, limit int) ([]Vote, int64, error)

	datasync.Store
}

type Collector struct {
	collectionInterval float64
	dcrClient          *rpcclient.Client
	dataStore          DataStore
	activeChain        *chaincfg.Params
	syncIsDone         bool
	bestBlockHeight    uint32
}
