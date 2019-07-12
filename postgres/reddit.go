package postgres

import (
	"context"
	"github.com/raedahgroup/dcrextdata/reddit"
	"time"
)

func (pg *PgDb) StoreRedditData(context.Context, reddit.RedditInfo) error {
	panic(1)
}

func (pg *PgDb) LastRedditEntryTime() (time time.Time) {
	panic(2)
}
