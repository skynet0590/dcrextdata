package datasync

import (
	"context"
	"time"
)

type SyncCoordinator struct {
	sources []string{}
	store Store
	syncers map[string]Syncer
}

type Syncer interface {
	// Sync fetches infomation from the given source and stores it for its table
	FetchSyncData(ctx context.Context, url string) (SyncResult, error)

	// Store save data gotten from the sync operation
	StoreSynceData(ctx context.Context, record interface{}) (error)
}

type Store interface {
	TableNames() []string
	SaveSyncHistory(ctx context.Context, history SyncHistory) error
	FetchSyncHistory(ctx context.Context, tableName string) (SyncHistory, error)

	StoreSyncResult(ctx context.Context, result SyncResult) error
	FetchSyncResult(ctx context.Context, request SyncRequest) (SyncResult, error)
}

type SyncHistory struct {
	Source string
	Table  string
	Date   time.Time
}

type SyncRequest struct {
	Table        string
	Date         time.Time
	MaxSkipCount int
	MaxTakeCount int
}

type SyncResult struct {
	Success bool
	Record interface{}
	TotalCount int
}