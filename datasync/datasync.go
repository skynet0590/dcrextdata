package datasync

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/raedahgroup/dcrextdata/app"
)

var coordinator *SyncCoordinator

func NewCoordinator(isEnabled bool) *SyncCoordinator {
	coordinator = &SyncCoordinator{
		instances: []instance{}, syncers: map[string]Syncer{}, isEnabled: isEnabled, syncersKeys: map[int]string{},
	}
	return coordinator
}

func (s *SyncCoordinator) AddSyncer(tableName string, syncer Syncer) {
	s.syncers[tableName] = syncer
	s.syncersKeys[len(s.syncersKeys)] = tableName
}

func (s *SyncCoordinator) AddSource(url string, store Store, database string) {
	s.instances = append(s.instances, instance{
		store: store,
		url:   url,
		database: database,
	})
}

func (s *SyncCoordinator) Syncer(tableName string) (Syncer, bool) {
	syncer, found := s.syncers[tableName]
	return syncer, found
}

func (s *SyncCoordinator) StartSyncing(ctx context.Context) {
	log.Info("Starting all registered sync collectors")

	runSyncers := func() {
		for {
			if app.MarkBusyIfFree() {
				break
			}
		}

		defer app.ReleaseForNewModule()

		for _, source := range s.instances {
			for i := 0; i <= len(s.syncersKeys); i++ {
				tableName := s.syncersKeys[i]
				syncer, found := s.syncers[tableName]
				if !found {
					return
				}

				log.Infof("Starting sync operation for %s on %s", tableName, source.url)

				err := s.sync(ctx, source, tableName, syncer)
				if err != nil {
					log.Error(err)
				}

				if err != nil {
					log.Error(err)
				}
			}
		}
	}

	runSyncers()

	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Infof("Stopping sync coordinators")
			return
		case <-ticker.C:
			runSyncers()
		}
	}
}

func (s *SyncCoordinator) sync(ctx context.Context, source instance, tableName string, syncer Syncer) error {
	startTime := time.Now()
	skip := 0
	take := 1000
	lastEntry, err := syncer.LastEntry(ctx, source.store)
	for {
		if err != nil {
			return fmt.Errorf("error in fetching sync history, %s", err.Error())

		}
		url := fmt.Sprintf("%s/api/sync/%s?last=%s&skip=%d&take=%d", source.url, tableName, lastEntry, skip, take)
		result, err := syncer.Collect(ctx, url)
		if err != nil {
			// todo: check if this is a sync disable error before stopping
			return fmt.Errorf("error in fetching data for %s, %s", url, err.Error())
		}

		if !result.Success {
			return fmt.Errorf("sync error, %s", result.Message)
		}

		if result.Records == nil {
			if result.TotalCount == 0 {
				return nil
			}
			duration := time.Now().Sub(startTime).Seconds()
			log.Infof("Synced %d %s records from %s into %s in %.4f seconds", result.TotalCount, tableName,
				source.url, source.database, math.Abs(duration))
			return nil
		}

		syncer.Append(ctx, source.store, result.Records)

		skip += take
	}
}

func Retrieve(ctx context.Context, tableName string, last string, skip, take int) (*Result, error) {
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

	return syncer.Retrieve(ctx, last, skip, take)
}

func DecodeSyncObj(obj interface{}, receiver interface{}) error {
	b, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, receiver)
	return err
}
