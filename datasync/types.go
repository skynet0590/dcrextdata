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
	syncers     map[string]Syncer
	syncersKeys map[int]string
	instances   []instance
	current     instance
	isEnabled   bool
}

type instance struct {
	db  Store
	url string
}

type Syncer struct {
	Collect  func(ctx context.Context, url string) (*Result, error)
	Retrieve func(ctx context.Context, last string, skip, take int) (*Result, error)
	Append   func(ctx context.Context, db Store, data interface{})
}

type Store interface {
	TableNames() []string
	LastEntry(ctx context.Context, tableName string) (string, error)
	SaveExchangeFromSync(ctx context.Context, exchange interface{}) error
	SaveExchangeTickFromSync(ctx context.Context, tick interface{}) error
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
