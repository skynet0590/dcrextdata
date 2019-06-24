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
)

func NewCollector(config *rpcclient.ConnConfig, dataStore DataStore) *Collector {
	return &Collector{
		dcrdClientConfig: config,
		dataStore:        dataStore,
	}
}

func (c Collector) StartMonitoring(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	client, err := rpcclient.New(c.dcrdClientConfig, nil)
	if err != nil {
		log.Error(err)
		return
	}

	defer client.Shutdown()

	collectMempool := func() {
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
			FirstSeenTime:        time.Now(),
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
