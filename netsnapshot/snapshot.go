package netsnapshot

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/decred/dcrd/chaincfg/v2"
	"github.com/raedahgroup/dcrextdata/app"
	"github.com/raedahgroup/dcrextdata/app/config"
)

func NewTaker(store DataStore, cfg config.NetworkSnapshotOptions) *taker {
	return &taker{
		dataStore: store,
		cfg:       cfg,
	}
}

func (t taker) Start(ctx context.Context) {
	for {
		if app.MarkBusyIfFree() {
			break
		}
	}
	log.Info("Triggering network snapshot taker.")

	var netParams = chaincfg.MainNetParams()
	if t.cfg.TestNet {
		netParams = chaincfg.TestNet3Params()
	}

	// defaultStaleTimeout = time.Minute * time.Duration(t.cfg.SnapshotInterval)
	// pruneExpireTimeout = defaultStaleTimeout * 2

	var err error
	amgr, err = NewManager(filepath.Join(defaultHomeDir,
		netParams.Name))
	if err != nil {
		fmt.Fprintf(os.Stderr, "NewManager: %v\n", err)
		os.Exit(1)
	}

	go runSeeder(t.cfg, netParams)

	var mtx sync.Mutex
	var bestBlockHeight int64
	var count int
	var timestamp = time.Now().UTC().Unix()

	snapshot := SnapShot{
		Timestamp: timestamp,
		Height:    bestBlockHeight,
		Nodes:     count,
	}
	err = t.dataStore.SaveSnapshot(ctx, snapshot)

	if err != nil {
		// todo delete all the related node info
		t.dataStore.DeleteSnapshot(ctx, timestamp)
		log.Errorf("Error in saving network snapshot, %s", err.Error())
	}

	ticker := time.NewTicker(time.Duration(t.cfg.SnapshotInterval) * time.Minute)
	defer ticker.Stop()

	for {
		// start listening for node heartbeat
		select {
		case <-ticker.C:
			err := t.dataStore.SaveSnapshot(ctx, SnapShot{
				Timestamp: timestamp,
				Height:    bestBlockHeight,
				Nodes:     count,
			})

			if err != nil {
				// todo delete all the related node info
				t.dataStore.DeleteSnapshot(ctx, timestamp)
				log.Errorf("Error in saving network snapshot, %s", err.Error())
			}

			mtx.Lock()
			count = 0
			log.Infof("Took a new network snapshot, recorded %d discoverable nodes.", count)
			timestamp = time.Now().UTC().Unix()
			mtx.Unlock()

		case node := <-amgr.goodPeer:
			err := t.dataStore.SaveNetworkPeer(ctx, NetworkPeer{
				Timestamp:       timestamp,
				Address:         node.IP.String(),
				LastSeen:        node.LastSeen.UTC().Unix(),
				ConnectionTime:  node.ConnectionTime,
				ProtocolVersion: node.ProtocolVersion,
				UserAgent:       node.UserAgent,
				StartingHeight:  node.StartingHeight,
				CurrentHeight:   node.CurrentHeight,
			})
			if err != nil {
				log.Errorf("Error in saving node info, %s.", err.Error())
			} else {
				mtx.Lock()
				count++
				if node.CurrentHeight > bestBlockHeight {
					bestBlockHeight = node.CurrentHeight
				}
				snapshot.Nodes = count
				snapshot.Height = bestBlockHeight

				snapshot := SnapShot{
					Timestamp: timestamp,
					Height:    bestBlockHeight,
					Nodes:     count,
				}

				err = t.dataStore.SaveSnapshot(ctx, snapshot)
				if err != nil {
					// todo delete all the related node info
					t.dataStore.DeleteSnapshot(ctx, timestamp)
					log.Errorf("Error in saving network snapshot, %s", err.Error())
				}

				mtx.Unlock()
				log.Infof("New heartbeat recorded for node: %s, %s, %d", node.IP.String(), node.UserAgent, node.ProtocolVersion)
			}
		case <-ctx.Done():
			log.Info("Shutting down network seeder")
			amgr.quit <- struct{}{}
			return
		}
	}
}
