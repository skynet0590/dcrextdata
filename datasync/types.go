package datasync

import (
	"context"
	"time"
)

type SyncCoordinator struct {
	store Store
	syncers map[string]Syncer
}

type Syncer interface {

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