package datasync

import (
	"context"
	"errors"
	"time"
)

var (
	ErrSyncDisabled = errors.New("data sharing is disabled on this instance")
)

type SyncCoordinator struct {
	syncers      map[string]Syncer
	historyStore HistoryStore
	sources      []string
	isEnabled    bool
}

type Syncer struct {
	Collect  func(ctx context.Context, url string) (*Result, error)
	Retrieve func(ctx context.Context, date time.Time, skip, take int) (*Result, error)
	Append   func(ctx context.Context, data interface{})
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
	Success    bool        `json:"success"`
	Message    string      `json:"message,omitempty"`
	Records    interface{} `json:"records,omitempty"`
	TotalCount int64       `json:"total_count,omitempty"`
}
