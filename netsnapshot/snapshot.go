package netsnapshot

import (
	"context"
	"time"

	"github.com/raedahgroup/dcrextdata/app"
	"github.com/raedahgroup/dcrextdata/app/config"
	"github.com/raedahgroup/dcrextdata/app/helpers"
)

func NewTaker(store DataStore, cfg config.NetworkSnapshotOptions) *taker {
	return &taker{
		dataStore: store,
		cfg: cfg,
	}
}

func (t taker) Start(ctx context.Context) {
	for {
		if app.MarkBusyIfFree() {
			break
		}
	}
	log.Info("Triggering network snapshot taker.")

	go runSeeder(t.cfg)

	lastCollectionDateUnix := t.dataStore.LastSnapshotTime(ctx)
	lastCollectionDate := time.Unix(lastCollectionDateUnix, 0)
	timePassed := time.Since(lastCollectionDate)
	period := time.Duration(t.cfg.SnapshotInterval) * time.Minute

	if lastCollectionDateUnix > 0 && timePassed < period {
		timeLeft := period - timePassed
		log.Infof("Taking network snapshot every %dm, took a snapshot %s ago, will take in %s.", t.cfg.SnapshotInterval,
			helpers.DurationToString(timePassed), helpers.DurationToString(timeLeft))

		app.ReleaseForNewModule()
		time.Sleep(timeLeft)
	}

	// wait for the first node discovery trip to complete
	for !seederIsReady {
		if ctx.Err() != nil {
			return
		}
		time.Sleep(2 * time.Second)
	}

	if lastCollectionDateUnix > 0 && timePassed < period {
		// continually check the state of the app until its free to run this module
		for {
			if app.MarkBusyIfFree() {
				break
			}
		}
	}

	t.takeSnapshot(ctx)
	app.ReleaseForNewModule()
	go t.takeSnapshotAtIntervals(ctx)
}

func (t taker) takeSnapshotAtIntervals(ctx context.Context) {
	if ctx.Err() != nil {
		return
	}

	ticker := time.NewTicker(time.Duration(t.cfg.SnapshotInterval) * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Infof("Stopping network snapshot taker.")
			return
		case <-ticker.C:
			// wait for the first node discovery trip to complete
			for !seederIsReady {
				if ctx.Err() != nil {
					return
				}
				time.Sleep(2 * time.Second)
			}

			// continually check the state of the app until its free to run this module
			for {
				if app.MarkBusyIfFree() {
					break
				}
			}

			t.takeSnapshot(ctx)
			app.ReleaseForNewModule()
		}
	}
}

func (t taker) takeSnapshot(ctx context.Context) {
	log.Info("Taking a new network snapshot.")

	timestamp := helpers.NowUtc().Unix()
	var bestBlockHeight int64 = 0
	peers := nodes()

	for _, peer := range peers {
		if peer.CurrentHeight > bestBlockHeight {
			bestBlockHeight = peer.CurrentHeight
		}
		err := t.dataStore.SaveNetworkPeer(ctx, NetworkPeer{
			Timestamp:       timestamp,
			Address:         peer.IP.String(),
			LastSeen:        peer.LastSeen.UTC().Unix(),
			ConnectionTime:  peer.ConnectionTime,
			ProtocolVersion: peer.ProtocolVersion,
			UserAgent:       peer.UserAgent,
			StartingHeight:  peer.StartingHeight,
			CurrentHeight:   peer.CurrentHeight,
		})
		if err != nil {
			log.Errorf("Error in saving peer info, %s.", err.Error())
		}
	}

	err := t.dataStore.SaveSnapshot(ctx, SnapShot{
		Timestamp: timestamp,
		Height:    bestBlockHeight,
		Nodes:     len(peers),
	})

	if err != nil {
		// todo delete all the related peer info
		t.dataStore.DeleteSnapshot(ctx, timestamp)
		log.Errorf("Error in saving network snapshot, %s", err.Error())
	}

	log.Infof("Took a new network snapshot, recorded %d discoverable nodes.", len(peers))
}
