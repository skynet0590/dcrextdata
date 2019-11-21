package netsnapshot

import (
	"context"
	"time"

	"github.com/decred/dcrd/rpcclient"
	"github.com/raedahgroup/dcrextdata/app"
	"github.com/raedahgroup/dcrextdata/app/helpers"
)

func NewTaker(dcrClient *rpcclient.Client, store DataStore, period int) *taker {
	return &taker{
		dcrClient: dcrClient,
		dataStore: store,
		period: period,
	}
}

func (t taker) Start(ctx context.Context) {
	for {
		if app.MarkBusyIfFree() {
			break
		}
	}
	log.Info("Triggering network snapshot taker.")

	lastCollectionDateUnix := t.dataStore.LastSnapshotTime(ctx)
	lastCollectionDate := time.Unix(lastCollectionDateUnix, 0)
	timePassed := time.Since(lastCollectionDate)
	period := time.Duration(t.period) * time.Minute

	if lastCollectionDateUnix > 0 && timePassed < period {
		timeLeft := period - timePassed
		log.Infof("Taking network snapshot every %dm, took a snapshot %s ago, will take in %s.", t.period,
			helpers.DurationToString(timePassed), helpers.DurationToString(timeLeft))

		app.ReleaseForNewModule()
		time.Sleep(timeLeft)
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

	ticker := time.NewTicker(time.Duration(t.period) * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Infof("Stopping network snapshot taker.")
			return
		case <-ticker.C:
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
	peerInfo, err := t.dcrClient.GetPeerInfo()
	if err != nil {
		log.Errorf("Error in getting peer info, %s.", err.Error())
		return
	}

	for _, peer := range peerInfo {
		if peer.CurrentHeight > bestBlockHeight {
			bestBlockHeight = peer.CurrentHeight
		}
		err := t.dataStore.SaveNetworkPeer(ctx, NetworkPeer{
			Timestamp:       timestamp,
			Address:         peer.Addr,
			LastReceiveTime: peer.LastRecv,
			LastSendTime:    peer.LastSend,
			ConnectionTime:  peer.ConnTime,
			ProtocolVersion: peer.Version,
			UserAgent:       peer.SubVer,
			StartingHeight:  peer.StartingHeight,
			CurrentHeight:   peer.CurrentHeight,
		})
		if err != nil {
			log.Errorf("Error in saving peer info, %s.", err.Error())
		}
	}

	err = t.dataStore.SaveSnapshot(ctx, SnapShot{
		Timestamp: timestamp,
		Height:    bestBlockHeight,
		Nodes:     len(peerInfo),
	})

	if err != nil {
		// todo delete all the related peer info
		t.dataStore.DeleteSnapshot(ctx, timestamp)
		log.Errorf("Error in saving network snapshot, %s", err.Error())
	}

	log.Infof("Took a new network snapshot, recorded %d discoverable nodes.", len(peerInfo))
}
