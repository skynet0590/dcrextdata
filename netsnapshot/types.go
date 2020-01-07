package netsnapshot

import (
	"context"

	"github.com/decred/dcrd/rpcclient"
	"github.com/raedahgroup/dcrextdata/app/config"
)

type SnapShot struct {
	Timestamp int64 `json:"timestamp"`
	Height    int64 `json:"height"`
}

type UserAgentInfo struct {
	UserAgent  string  `json:"user_agent"`
	Nodes 	   int64   `json:"nodes"`
	Percentage int64 `json:"percentage"`
}

type CountryInfo struct {
	Country    string  `json:"country"`
	Nodes 	   int64   `json:"nodes"`
	Percentage int64   `json:"percentage"`
}

type NetworkPeer struct {
	Timestamp       int64  `json:"timestamp"`
	Address         string `json:"address"`
	Country         string `json:"country"`
	UserAgent       string `json:"user_agent"`
	StartingHeight  int64  `json:"starting_height"`
	CurrentHeight   int64  `json:"current_height"`
	ConnectionTime  int64  `json:"connection_time"`
	ProtocolVersion uint32 `json:"protocol_version"`
	LastSeen        int64  `json:"last_seen"`
	LastSuccess     int64  `json:"last_success"`
	IsDead          bool   `json:"is_dead"`
	Latency         int    `json:"latency"`
	IPVersion       int    `json:"ip_version"`
	Services        string `json:"services"`
}

type geoIP struct {
	CountryCode string  `json:"country_code"`
	CountryName string  `json:"country_name"`
	Type 		string  `json:"type"`
}

type DataStore interface {
	LastSnapshotTime(ctx context.Context) (timestamp int64)
	DeleteSnapshot(ctx context.Context, timestamp int64)
	SaveSnapshot(ctx context.Context, snapShot SnapShot) error
	SaveNetworkPeer(ctx context.Context, peer NetworkPeer) error
	LastSnapshot(ctx context.Context) (*SnapShot, error)
	GetIPLocation(ctx context.Context, ip string) (string, int, error)
}

type taker struct {
	dcrClient *rpcclient.Client
	dataStore DataStore
	cfg       config.NetworkSnapshotOptions
}
