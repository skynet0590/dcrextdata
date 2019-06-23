// Copyright (c) 2018-2019 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package mempool

import (
	"context"
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

func (c Collector) StartMonitoring(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	client, err := rpcclient.New(c.dcrdClientConfig, nil)
	if err != nil {
		log.Error(err)
		return
	}

	defer client.Shutdown()

	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			mempoolTransactionMap, err := client.GetRawMempoolVerbose(dcrjson.GRMAll)
			if err != nil {
				log.Error(err)
				break
			}

			mempoolDto := Mempool{
				NumberOfTransactions:len(mempoolTransactionMap),
				Time : time.Now(),
				FirstSeenTime:time.Now(),
			}

			for _, tx := range mempoolTransactionMap {
				mempoolDto.Fee += tx.Fee
				mempoolDto.Size += tx.Size
				if mempoolDto.FirstSeenTime.Unix() < tx.Time {
					mempoolDto.FirstSeenTime = time.Unix(tx.Time, 0)
				}
			}
			break
		case <-ctx.Done():
			return
		}
	}
}
