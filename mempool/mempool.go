// Copyright (c) 2018-2019 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package mempool

import (
	"context"
	"database/sql"
	"sync"

	"github.com/decred/dcrd/chaincfg/chainhash"
	"github.com/decred/dcrd/dcrjson"
	"github.com/decred/dcrd/rpcclient"
)

const (
	genesisHashStr = "298e5cc3d985bfe7f81dc135f360abe089edd4396b86d2de66b0cef42b21d980"
)

func NewCollector(config *rpcclient.ConnConfig, dataStore DataStore) (*Collector) {
	return &Collector{
		dcrdClientConfig: config,
		dataStore:dataStore,
	}
}

func (c Collector) Collect(ctx context.Context, wg sync.WaitGroup) error {
	defer wg.Done()
	log.Info("Starting historic mempool data collector")
	client, err := rpcclient.New(c.dcrdClientConfig, nil)
	if err != nil {
		return err
	}

	defer client.Shutdown()

	if ctx.Err() != nil {
		return ctx.Err()
	}

	for {
		select {
		default:
			var nextBlockHash *chainhash.Hash

			lastMempool, err := c.dataStore.LastMempool(ctx)
			if err == nil {
				lastBlockHash, err := client.GetBlockHash(lastMempool.LastBlockHeight)
				if err != nil {
					return err
				}

				lastBlockInfo, err := client.GetBlockVerbose(lastBlockHash, false)
				if err != nil {
					return err
				}

				if lastBlockInfo.NextHash == "" {
					// we are done fetching historic data, start checking at intervals to fetch current data

					return nil
				}

				nextBlockHash, err = chainhash.NewHashFromStr(lastBlockInfo.NextHash)
				if err != nil {
					return err
				}
			} else if err == sql.ErrNoRows {
				nextBlockHash, err = chainhash.NewHashFromStr(genesisHashStr)
				if err != nil {
					return err
				}
			} else {
				return err
			}

			nextBlock, err := client.GetBlockVerbose(nextBlockHash, true)
			if err != nil {
				return err
			}

			total := sumOutsTxRawResult(nextBlock.RawTx) + sumOutsTxRawResult(nextBlock.RawSTx)

			mempool := Mempool{
				LastBlockHeight:nextBlock.Height,
				BlockReceiveTime:nextBlock.Time, // todo this is the time of the next block
				RegularTransactionCount:len(nextBlock.Tx),
				RevocationCount:nextBlock.Revocations,
				Size:nextBlock.Size,
				VoteCount:nextBlock.Voters,
				TicketsCount:nextBlock.FreshStake,
				TotalSent: total,
				FirstSeenTime:0,// todo this should be gotten from the first tx
			}

			err = c.dataStore.StoreMempool(ctx, mempool)
			if err != nil {
				return err
			}
			break
			case <- ctx.Done():
				return ctx.Err()
		}
	}
}

func sumOutsTxRawResult(txs []dcrjson.TxRawResult) (sum float64) {
	for _, tx := range txs {
		for _, vout := range tx.Vout {
			sum += vout.Value
		}
	}
	return
}