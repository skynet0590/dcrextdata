package datasync

import (
	"context"
	"time"
)

type SyncCoordinator struct {
	syncers      map[string]Syncer
	historyStore HistoryStore
	sources      []string
}

type Syncer interface {
	Collect(ctx context.Context, urs string) (*Result, error)
	Retrieve(ctx context.Context, date time.Time, skip, take int) (*Result, error)
	Append(ctx context.Context, result Result) error
} 

type HistoryStore interface {
	TableNames() []string
	SaveSyncHistory(ctx context.Context, history History) error
	FetchSyncHistory(ctx context.Context, tableName string, source string) (History, error)
}

type History struct {
	Source string
	Table  string
	Date   time.Time
}

type Request struct {
	Table        string
	Date         time.Time
	MaxSkipCount int
	MaxTakeCount int
}

type Result struct {
	Success    bool
	Message    string
	Record     interface{}
	TotalCount int
}
