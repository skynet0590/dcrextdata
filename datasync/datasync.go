package datasync

import (
	"context"
)

func NewSyncCoordination(sources []string, store Store, syncers map[string]Syncer) *SyncCoordinator {
	return &SyncCoordinator{
		sources: sources,
		store: store, 
		syncers: syncers, 
	}
}

func (s *SyncCoordinator) Run(ctx context.Context) {
	for name, syncher := range s.syncers {
		
	}
}