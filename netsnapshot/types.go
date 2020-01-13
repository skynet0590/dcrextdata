package netsnapshot

import (
	"context"
	"net"

	"github.com/decred/dcrd/rpcclient"
	"github.com/raedahgroup/dcrextdata/app/config"
)

type SnapShot struct {
	Timestamp int64 `json:"timestamp"`
	Height    int64 `json:"height"`
}

type NodeCount struct {
	Timestamp int64 `json:"timestamp"`
	Count     int64 `json:"count"`
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
	LastAttempt     int64  `json:"last_attempt"`

	IPInfo
}

type Heartbeat struct {
	Timestamp     int64  `json:"timestamp"`
	Address       string `json:"address"`
	LastSeen      int64  `json:"last_seen"`
	Latency       int    `json:"latency"`
	CurrentHeight int64  `json:"current_height"`
}

type IPInfo struct {
	Type        string `json:"type"`
	CountryCode string `json:"country_code"`
	CountryName string `json:"country_name"`
	RegionCode  string `json:"region_code"`
	RegionName  string `json:"region_name"`
	City        string `json:"city"`
	Zip         string `json:"zip"`
}

type DataStore interface {
	LastSnapshotTime(ctx context.Context) (timestamp int64)
	DeleteSnapshot(ctx context.Context, timestamp int64)
	SaveSnapshot(ctx context.Context, snapShot SnapShot) error
	SaveHeartbeat(ctx context.Context, peer Heartbeat) error
	SaveNode(ctx context.Context, peer NetworkPeer) error
	UpdateNode(ctx context.Context, peer NetworkPeer) error
	GetAvailableNodes(ctx context.Context) ([]net.IP, error)
	LastSnapshot(ctx context.Context) (*SnapShot, error)
	GetIPLocation(ctx context.Context, ip string) (string, int, error)
	NodeExists(ctx context.Context, address string) (bool, error)
}

type taker struct {
	dcrClient *rpcclient.Client
	dataStore DataStore
	cfg       config.NetworkSnapshotOptions
}
