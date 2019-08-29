package datasync

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"
)

var coordinator *SyncCoordinator

func NewCoordinator(store HistoryStore, sources []string, isEnabled bool) *SyncCoordinator {
	coordinator = &SyncCoordinator{
		sources: sources, historyStore: store, syncers: map[string]Syncer{}, isEnabled: isEnabled,
	}
	return coordinator
}

func (s *SyncCoordinator) AddSyncer(tableName string, syncer Syncer) {
	s.syncers[tableName] = syncer
}

func (s *SyncCoordinator) Syncer(tableName string) (Syncer, bool) {
	syncer, found := s.syncers[tableName]
	return syncer, found
}

func (s *SyncCoordinator) StartSyncing(ctx context.Context) {
	for _, source := range s.sources {
		for tableName, syncer := range s.syncers {
			err := s.sync(ctx, source, tableName, syncer)
			if err != nil {
				log.Error(err)
			}
		}
	}
}

func (s *SyncCoordinator) sync(ctx context.Context, source string, tableName string, syncer Syncer) error {
	syncHistory, err := s.historyStore.FetchSyncHistory(ctx, tableName, source)
	if err != nil {
		return fmt.Errorf("error in fetching sync history, %s", err.Error())

	}
	startTime := time.Now()
	skip := 0
	take := 100
	for {
		url := fmt.Sprint("%s?date=%s&skip=%d&take=%d", source, syncHistory.Date.Format(time.RFC3339Nano), 0, 10)
		result, err := syncer.Collect(ctx, url)
		if err != nil {
			// todo: check if this is a sync disable error before stopping
			return err
		}

		if !result.Success {
			return fmt.Errorf("sync error, %s", result.Message)
		}

		syncer.Append(ctx, result.Record)

		skip += take
		if result.TotalCount <= int64(skip) {
			duration := time.Now().Sub(startTime).Seconds()
			log.Infof("Synced %d %s records from %s in %d seconds", result.TotalCount, tableName, source, math.Abs(duration))
			return nil
		}
	}
}

func Retrieve(ctx context.Context, tableName string, date time.Time, skip, take int) (*Result, error) {
	if coordinator == nil {
		return nil, errors.New("syncer not initialized")
	}

	if !coordinator.isEnabled {
		return nil, ErrSyncDisabled
	}

	syncer, found := coordinator.syncers[tableName]
	if !found {
		return nil, errors.New("syncer not found for " + tableName)
	}

	return syncer.Retrieve(ctx, date, skip, take)
}
