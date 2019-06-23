// Copyright (c) 2018-2019 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package mempool

import (
	"context"
	"github.com/decred/dcrd/wire"
	"sync"
	"time"

	"github.com/decred/dcrd/dcrjson"
	"github.com/decred/dcrd/rpcclient"
)

func NewCollector(config *rpcclient.ConnConfig, dataStore DataStore) *Collector {
	return &Collector{
		dcrdClientConfig: config,
		dataStore:        dataStore,
	}
}

func (c *Collector) StartMonitoring(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	ntfnHandlers := rpcclient.NotificationHandlers{
		OnTxAcceptedVerbose: func(txDetails *dcrjson.TxRawResult) {
			if c.currentMempool == nil {
				c.currentMempool = &Mempool{
					FirstSeenTime: time.Now(),
				}
			}
			c.currentMempool.NumberOfTransactions++
		},
		OnBlockConnected: func(blockHeaderSerialized []byte, transactions [][]byte) {
			blockHeader := new(wire.BlockHeader)
			err := blockHeader.FromBytes(blockHeaderSerialized)
			if err != nil {
				log.Error("Failed to deserialize blockHeader in new block notification: %v", err)
				return
			}

			if c.currentMempool == nil {
				return
			}

			mempool := *c.currentMempool
			c.currentMempool = nil

			mempool.Size = blockHeader.Size
			mempool.BlockReceiveTime = time.Now()
			mempool.BlockInternalTime = blockHeader.Timestamp
			mempool.BlockHeight = blockHeader.Height
			mempool.BlockHash = blockHeader.BlockHash().String()

			err = c.dataStore.StoreMempool(ctx, mempool)
			if err != nil {
				log.Error(err)
				return
			}
		},
	}

	client, err := rpcclient.New(c.dcrdClientConfig, &ntfnHandlers)
	if err != nil {
		log.Error(err)
		return
	}

	if err := client.NotifyNewTransactions(true); err != nil {
		log.Error(err)
	}

	if err := client.NotifyBlocks(); err != nil {
		log.Error(err)
	}

	defer client.Shutdown()

	for {
		select {
		case <-ctx.Done():
			return
		}
	}
}
