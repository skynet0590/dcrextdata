package netsnapshot

import (
	"context"
	"fmt"
	"math"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/decred/dcrd/chaincfg/v2"
	"github.com/raedahgroup/dcrextdata/app"
	"github.com/raedahgroup/dcrextdata/app/config"
	"github.com/raedahgroup/dcrextdata/app/helpers"
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

	// update all reachable nodes
	loadLiveNodes := func() {
		nodes, err := t.dataStore.GetAvailableNodes(ctx)
		if err != nil {
			log.Errorf("Error in taking network snapshot, %s", err.Error())
		}
		amgr.setLiveNodes(nodes)
		log.Infof("%d nodes loaded to live reload", len(nodes))
	}

	// enqueue previous known ips
	loadLiveNodes()

	go runSeeder(t.cfg, netParams)

	var mtx sync.Mutex
	var bestBlockHeight int64

	var count int
	var timestamp = time.Now().UTC().Unix()
	snapshot := SnapShot{
		Timestamp: timestamp,
		Height:    bestBlockHeight,
	}

	lastSnapshot, err := t.dataStore.LastSnapshot(ctx)
	if err == nil {
		minutesPassed := math.Abs(time.Since(time.Unix(lastSnapshot.Timestamp, 0)).Minutes())
		if minutesPassed < float64(t.cfg.SnapshotInterval)/2 {
			snapshot = *lastSnapshot
		}
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
			})

			if err != nil {
				t.dataStore.DeleteSnapshot(ctx, timestamp)
				log.Errorf("Error in saving network snapshot, %s", err.Error())
			}

			mtx.Lock()
			count = 0
			log.Infof("Took a new network snapshot, recorded %d discoverable nodes.", count)
			timestamp = time.Now().UTC().Unix()
			mtx.Unlock()
			// update all reachable nodes
			loadLiveNodes()


		case node := <-amgr.goodPeer:
			if node.IP.String() == "127.0.0.1" { // do not add the local IP
				break
			}

			networkPeer := NetworkPeer{
				Timestamp:       timestamp,
				Address:         node.IP.String(),
				LastSeen:        node.LastSeen.UTC().Unix(),
				LastSuccess: 	 node.LastSuccess.UTC().Unix(),
				ConnectionTime:  node.ConnectionTime,
				ProtocolVersion: node.ProtocolVersion,
				UserAgent:       node.UserAgent,
				StartingHeight:  node.StartingHeight,
				CurrentHeight:   node.CurrentHeight,
				Services: 		 node.Services.String(),
				Latency: 		 int(node.Latency),
			}

			if exists, _ := t.dataStore.NodeExists(ctx, networkPeer.Address); exists {
				err = t.dataStore.UpdateNode(ctx, networkPeer)
				if err != nil {
					log.Errorf("Error in saving node info, %s.", err.Error())
				}
			} else {
				geoLoc, err := t.geolocation(ctx, node.IP)
				if err == nil {
					networkPeer.IPInfo = *geoLoc
					// networkPeer.Country = geoLoc.CountryName
					if geoLoc.Type == "ipv4" {
						networkPeer.IPVersion = 4
					} else if geoLoc.Type == "ipv6" {
						networkPeer.IPVersion = 6
					}
				}

				err = t.dataStore.SaveNode(ctx, networkPeer)
				if err != nil {
					log.Errorf("Error in saving node info, %s.", err.Error())
				}
			}

			// if this node is reachable within the current timestamp, save heartbeat
			if networkPeer.LastSuccess >= timestamp {
				err = t.dataStore.SaveHeartbeat(ctx, Heartbeat{
					Timestamp: timestamp,
					Address:   node.IP.String(),
					LastSeen:  node.LastSeen.UTC().Unix(),
					Latency:   int(node.Latency),
				})
				if err != nil {
					log.Errorf("Error in saving node info, %s.", err.Error())
				} else {
					mtx.Lock()
					count++
					if node.CurrentHeight > bestBlockHeight {
						bestBlockHeight = node.CurrentHeight
					}

					snapshot := SnapShot{
						Timestamp: timestamp,
						Height:    bestBlockHeight,
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
			}
		case <-ctx.Done():
			log.Info("Shutting down network seeder")
			amgr.quit <- struct{}{}
			return
		}
	}
}

func (t taker) geolocation(ctx context.Context, ip net.IP) (*IPInfo, error) {
	url := fmt.Sprintf("http://api.ipstack.com/%s?access_key=fcd33d8814206ce1f0a255a2204ad71e&format=1", ip.String())
	var geo IPInfo
	err := helpers.GetResponse(ctx, &http.Client{Timeout: 3 * time.Second}, url, &geo)
	return &geo, err
}
