// Copyright (c) 2018-2019 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package mempool

import (
	"context"
	"sync"
	"time"

	"github.com/decred/dcrd/chaincfg/chainhash"
	"github.com/decred/dcrd/dcrjson"
	"github.com/decred/dcrd/rpcclient"
	"github.com/decred/dcrd/wire"
	"github.com/decred/dcrdata/txhelpers/v2"
)

func NewCollector(config *rpcclient.ConnConfig, dataStore DataStore) *Collector {
	return &Collector{
		dcrdClientConfig: config,
		dataStore:        dataStore,
	}
}

func (c Collector) StartMonitoring(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	var lastBlockHeight uint32

	ntfnHandlers := rpcclient.NotificationHandlers{
		OnTxAcceptedVerbose: func(txDetails *dcrjson.TxRawResult) {
			msgTx, err := txhelpers.MsgTxFromHex(txDetails.Hex)
			if err != nil {
				log.Errorf("Failed to decode transaction hex: %v", err)
				return
			}

			if txType := txhelpers.DetermineTxTypeString(msgTx); txType != "Vote" {
				return
			}

			vote := Vote{
				ReceiveTime:time.Now(),
				BlockHeight: int64(lastBlockHeight),
				Hash: txDetails.Txid,
			}

			if err = c.dataStore.SaveVote(ctx, vote); err != nil {
				log.Error(err)
			}
		},
		OnBlockConnected: func(blockHeaderSerialized []byte, transactions [][]byte) {
			blockHeader := new(wire.BlockHeader)
			err := blockHeader.FromBytes(blockHeaderSerialized)
			if err != nil {
				log.Error("Failed to deserialize blockHeader in new block notification: %v", err)
				return
			}

			lastBlockHeight = blockHeader.Height
			block := Block{
				BlockInternalTime:blockHeader.Timestamp,
				BlockReceiveTime:time.Now(),
				BlockHash:blockHeader.BlockHash().String(),
				BlockHeight: blockHeader.Height,
			}
			if err = c.dataStore.SaveBlock(ctx, block); err != nil {
				log.Error(err)
			}
		},
	}

	client, err := rpcclient.New(c.dcrdClientConfig, &ntfnHandlers)
	if err != nil {
		log.Error(err)
		return
	}

	defer client.Shutdown()

	if err := client.NotifyNewTransactions(true); err != nil {
		log.Error(err)
	}

	if err := client.NotifyBlocks(); err != nil {
		log.Error(err)
	}

	var mu sync.Mutex

	collectMempool := func() {
		mu.Lock()
		defer mu.Unlock()

		mempoolTransactionMap, err := client.GetRawMempoolVerbose(dcrjson.GRMAll)
		if err != nil {
			log.Error(err)
			return
		}

		if len(mempoolTransactionMap) == 0 {
			return
		}

		mempoolDto := Mempool{
			NumberOfTransactions: len(mempoolTransactionMap),
			Time:                 time.Now(),
			FirstSeenTime:        time.Now(),//todo: use the time of the first tx in the mempool
		}

		for hashString, tx := range mempoolTransactionMap {
			hash, err := chainhash.NewHashFromStr(hashString)
			if err != nil {
				log.Error(err)
				continue
			}
			rawTx, err := client.GetRawTransactionVerbose(hash)
			if err != nil {
				log.Error(err)
				continue
			}

			totalOut := 0.0
			for _, v := range rawTx.Vout {
				totalOut += v.Value
			}

			mempoolDto.Total += totalOut
			mempoolDto.TotalFee += tx.Fee
			mempoolDto.Size += tx.Size
			if mempoolDto.FirstSeenTime.Unix() > tx.Time {
				mempoolDto.FirstSeenTime = time.Unix(tx.Time, 0)
			}

		}

		votes, err := client.GetRawMempool(dcrjson.GRMVotes)
		if err != nil {
			log.Error(err)
			return
		}
		mempoolDto.Voters = len(votes)

		tickets, err := client.GetRawMempool(dcrjson.GRMTickets)
		if err != nil {
			log.Error(err)
			return
		}
		mempoolDto.Tickets = len(tickets)

		revocations, err := client.GetRawMempool(dcrjson.GRMRevocations)
		if err != nil {
			log.Error(err)
			return
		}
		mempoolDto.Revocations = len(revocations)

		err = c.dataStore.StoreMempool(ctx, mempoolDto)
		if err != nil {
			log.Error(err)
		}
	}

	collectMempool()
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			collectMempool()
			break
		case <-ctx.Done():
			return
		}
	}
}
