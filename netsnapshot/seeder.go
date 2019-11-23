// Copyright (c) 2018 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package netsnapshot

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/decred/dcrd/chaincfg/v2"
	"github.com/decred/dcrd/dcrutil/v2"
	"github.com/decred/dcrd/peer/v2"
	"github.com/decred/dcrd/wire"
	"github.com/raedahgroup/dcrextdata/app/config"
)

const (
	// defaultAddressTimeout defines the duration to wait
	// for new addresses.
	defaultAddressTimeout = time.Minute * 10

	// defaultNodeTimeout defines the timeout time waiting for
	// a response from a node.
	defaultNodeTimeout = time.Second * 3
)

var (
	defaultHomeDir    = dcrutil.AppDataDir("dcrextdata", false)

	amgr *Manager
	wg   sync.WaitGroup
	seederIsReadyMtx sync.Mutex
	seederIsReady bool
)

func setSeederIsReady (val bool) {
	seederIsReadyMtx.Lock()
	defer seederIsReadyMtx.Unlock()
	seederIsReady = val
}

func creep(netParams *chaincfg.Params) {
	defer wg.Done()

	onaddr := make(chan struct{})
	verack := make(chan struct{})
	peerConfig := peer.Config{
		UserAgentName:    "dcrpeersniffer",
		UserAgentVersion: "0.0.1",
		Net:              netParams.Net,
		DisableRelayTx:   true,

		Listeners: peer.MessageListeners{
			OnAddr: func(p *peer.Peer, msg *wire.MsgAddr) {
				n := make([]net.IP, 0, len(msg.AddrList))
				for _, addr := range msg.AddrList {
					n = append(n, addr.IP)
				}
				added := amgr.AddAddresses(n)
				log.Infof("Peer %v sent %v addresses, %d new",
					p.Addr(), len(msg.AddrList), added)
				onaddr <- struct{}{}
			},
			OnVerAck: func(p *peer.Peer, msg *wire.MsgVerAck) {
				setSeederIsReady(false)
				log.Infof("Adding peer %v with services %v",
					p.NA().IP.String(), p.Services())

				verack <- struct{}{}
			},
		},
	}

	var wg sync.WaitGroup
	for {
		ips := amgr.Addresses()
		if len(ips) == 0 {
			seederIsReady = true
			log.Infof("No stale addresses -- sleeping for %v",
				defaultAddressTimeout)
			time.Sleep(defaultAddressTimeout)
			continue
		}

		wg.Add(len(ips))

		for _, ip := range ips {
			go func(ip net.IP) {
				defer wg.Done()

				host := net.JoinHostPort(ip.String(),
					netParams.DefaultPort)
				p, err := peer.NewOutboundPeer(&peerConfig, host)
				if err != nil {
					log.Infof("NewOutboundPeer on %v: %v",
						host, err)
					return
				}
				amgr.Attempt(ip)
				conn, err := net.DialTimeout("tcp", p.Addr(),
					defaultNodeTimeout)
				if err != nil {
					return
				}
				p.AssociateConnection(conn)

				// Wait for the verack message or timeout in case of
				// failure.
				select {
				case <-verack:
					// Mark this peer as a good node.
					amgr.Good(p)

					// Ask peer for some addresses.
					p.QueueMessage(wire.NewMsgGetAddr(), nil)

				case <-time.After(defaultNodeTimeout):
					log.Infof("verack timeout on peer %v",
						p.Addr())
					p.Disconnect()
					return
				}

				select {
				case <-onaddr:
				case <-time.After(defaultNodeTimeout):
					log.Infof("getaddr timeout on peer %v",
						p.Addr())
					p.Disconnect()
					return
				}
				p.Disconnect()
			}(ip)
		}
		wg.Wait()
	}
}

func runSeeder(cfg config.NetworkSnapshotOptions) {
	var netParams = chaincfg.MainNetParams()
	if cfg.TestNet {
		netParams = chaincfg.TestNet3Params()
	}

	var err error
	amgr, err = NewManager(filepath.Join(defaultHomeDir,
		netParams.Name))
	if err != nil {
		fmt.Fprintf(os.Stderr, "NewManager: %v\n", err)
		os.Exit(1)
	}

	amgr.AddAddresses([]net.IP{net.ParseIP(cfg.Seeder)})

	wg.Add(1)
	go creep(netParams)

	// todo remove dns related feature
	/*dnsServer := NewDNSServer(cfg.SeederHost, cfg.Nameserver, cfg.Listen)
	go dnsServer.Start()*/

	wg.Wait()
}

func nodes() map[string]*Node {
	amgr.mtx.RLock()
	defer amgr.mtx.RUnlock()
	return amgr.nodes
}
