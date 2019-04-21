// Copyright (c) 2018-2019 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"os"
	"runtime"
	"sync"

	"github.com/raedahgroup/dcrextdata/exchanges"
	"github.com/raedahgroup/dcrextdata/postgres"
	"github.com/raedahgroup/dcrextdata/version"
	"github.com/raedahgroup/dcrextdata/vsp"
)

// const dcrlaunchtime int64 = 1454889600

func _main(ctx context.Context) error {
	cfg, err := loadConfig()
	if err != nil {
		// fmt.Printf("Unable to load config: %v\n", err)
		return err
	}

	defer func() {
		if logRotator != nil {
			logRotator.Close()
		}
	}()

	// Display app version.
	log.Infof("%s version %v (Go version %s)", version.AppName,
		version.Version(), runtime.Version())

	db, err := postgres.NewPgDb(cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPass, cfg.DBName)
	defer func(db *postgres.PgDb) {
		err := db.Close()
		if err != nil {
			log.Error("Could not close database connection: %v", err)
		}
	}(db)

	if err != nil {
		return err
	}

	wg := new(sync.WaitGroup)

	if cfg.VSPEnabled {
		log.Info("Starting VSP data collection")
		vspCollector, err := vsp.NewVspCollector(cfg.VSPInterval, db)
		if err == nil {
			wg.Add(1)
			go vspCollector.Run(ctx, wg)
		} else {
			log.Error(err)
		}
	}

	if cfg.ExchangesEnabled {
		log.Infof("Starting exchange ticker data collection for %v", cfg.Exchanges)
		ticksHub, err := exchanges.NewTickHub(ctx, cfg.Exchanges, db)
		if err == nil {
			wg.Add(1)
			go ticksHub.Run(ctx, wg)
			// ticksHub.Collect(ctx)
		} else {
			log.Error(err)
		}
	}

	wg.Wait()
	log.Info("Goodbye")
	return nil
}

func main() {
	// Create a context that is cancelled when a shutdown request is received
	// via requestShutdown.
	ctx := withShutdownCancel(context.Background())
	// Listen for both interrupt signals and shutdown requests.
	go shutdownListener()

	if err := _main(ctx); err != nil {
		if logRotator != nil {
			log.Error(err)
		}
		os.Exit(1)
	}
	os.Exit(0)
}
