// Copyright (c) 2018-2019 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package mempool

import (
	"context"
	"sync"

	"github.com/decred/dcrd/rpcclient"
)

func NewCollector(config *rpcclient.ConnConfig, dataStore DataStore) (*Collector, error) {
	client, err := rpcclient.New(config, nil)
	if err != nil {
		return nil, err
	}

	return &Collector{
		dcrdClient: client,
		dataStore:dataStore,
	}, nil
}

func (c Collector) Collect(ctx context.Context, wg sync.WaitGroup) error {
	defer wg.Done()
	for {

	}
}