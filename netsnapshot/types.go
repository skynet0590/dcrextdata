package netsnapshot

import (
	"context"
	
	"github.com/decred/dcrd/chaincfg"
	"github.com/decred/dcrd/rpcclient"
)

type NetworkPeer struct {
	ID              int    `json:"id"`
	Address         string `json:"address"`
	LastReceiveTime int64  `json:"last_receive_time"`
	LastSendTime    int64  `json:"last_send_time"`
	ConnectionTime  int64  `json:"connection_time"`
	ProtocolVersion uint32 `json:"protocol_version"`
	UserAgent       string `json:"user_agent"`
	StartingHeight  int64  `json:"starting_height"`
	CurrentHeight   int64  `json:"current_height"`
}

type DataStore interface {
	SaveNetworkPeer(ctx context.Context, peer NetworkPeer) error
	NetworkPeers(ctx context.Context, q string, offset int, limit int) ([]NetworkPeer, int64, error)
	TotalPeerCount(ctx context.Context) (int64, error)
	TotalPeerCountByProtocol(ctx context.Context, protocolVersion int) (int64, error)
	PeerCountByUserAgents(ctx context.Context) (counts map[string]int64, err error)
}

type taker struct {
	dcrClient          *rpcclient.Client
	dataStore          DataStore
	activeChain        *chaincfg.Params
}
