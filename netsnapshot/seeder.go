// Copyright (c) 2018 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package netsnapshot

import (
	"net"
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
	defaultNodeTimeout = time.Second * 10
)

var (
	defaultHomeDir = dcrutil.AppDataDir("dcrextdata", false)

	amgr             *Manager
	wg               sync.WaitGroup
)

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
			log.Infof("No stale addresses -- sleeping for %v",
				defaultAddressTimeout)
			time.Sleep(defaultAddressTimeout)
			continue
		}

		wg.Add(len(ips))

		for _, ip := range ips {
			go func(ip net.IP) {
				defer wg.Done()

				port := netParams.DefaultPort
				if ip.String() == amgr.Seeder && amgr.SeederPort != "" {
					port = amgr.SeederPort
				}
				host := net.JoinHostPort(ip.String(),
					port)
				p, err := peer.NewOutboundPeer(&peerConfig, host)
				if err != nil {
					log.Infof("NewOutboundPeer on %v: %v",
						host, err)
					return
				}
				amgr.Attempt(ip)
				t := time.Now()
				conn, err := net.DialTimeout("tcp", p.Addr(),
					defaultNodeTimeout)
				if err != nil {
					currHeight := p.LastBlock()
					if currHeight == 0 {
						currHeight = p.StartingHeight()
					}
					amgr.goodPeer <- &Node{
						IP:              ip,
						Services:        p.Services(),
						LastAttempt:     time.Now().UTC(),
						LastSuccess:     p.TimeConnected(),
						LastSeen:        p.TimeConnected(),
						Latency:		 -1, // peer is down
						ConnectionTime:  p.TimeConnected().Unix(),
						ProtocolVersion: p.ProtocolVersion(),
						UserAgent:       p.UserAgent(),
						StartingHeight:  p.StartingHeight(),
						CurrentHeight:   currHeight,
					}
					return
				}
				latency := time.Since(t).Milliseconds()
				p.AssociateConnection(conn)

				// Wait for the verack message or timeout in case of
				// failure.
				select {
				case <-verack:
					// Mark this peer as a good node.
					amgr.Good(p)
					amgr.goodPeer <- &Node{
						IP:              ip,
						Services:        p.Services(),
						LastAttempt:     time.Now().UTC(),
						LastSuccess:     time.Now().UTC(),
						LastSeen:        time.Now().UTC(),
						Latency:		 latency,
						ConnectionTime:  p.TimeConnected().Unix(),
						ProtocolVersion: p.ProtocolVersion(),
						UserAgent:       p.UserAgent(),
						StartingHeight:  p.StartingHeight(),
						CurrentHeight:   p.LastBlock(),
					}

					// Ask peer for some addresses.
					p.QueueMessage(wire.NewMsgGetAddr(), nil)

				case <-time.After(defaultNodeTimeout):
					log.Infof("verack timeout on peer %v",
						p.Addr())
						currHeight := p.LastBlock()
						if currHeight == 0 {
							currHeight = p.StartingHeight()
						}
						amgr.goodPeer <- &Node{
							IP:              ip,
							Services:        p.Services(),
							LastAttempt:     time.Now().UTC(),
							LastSuccess:     p.TimeConnected(),
							LastSeen:        p.TimeConnected(),
							Latency:		 latency,
							ConnectionTime:  p.TimeConnected().Unix(),
							ProtocolVersion: p.ProtocolVersion(),
							UserAgent:       p.UserAgent(),
							StartingHeight:  p.StartingHeight(),
							CurrentHeight:   currHeight,
						}
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

func runSeeder(cfg config.NetworkSnapshotOptions, netParams *chaincfg.Params) {
	amgr.AddAddresses([]net.IP{net.ParseIP(cfg.Seeder)})
	amgr.Seeder = cfg.Seeder
	amgr.SeederPort = cfg.SeederPort

	wg.Add(1)
	go creep(netParams)

	wg.Wait()
}
