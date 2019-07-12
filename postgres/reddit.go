package postgres

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/raedahgroup/dcrextdata/postgres/models"
	"github.com/raedahgroup/dcrextdata/reddit"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
)

func (pg *PgDb) StoreRedditData(ctx context.Context, redditInfo reddit.RedditInfo) error {
	redditModel := models.RedditInfo{
		Date:              redditInfo.Date,
		Subscribers:          redditInfo.Subscribers,
		AccountsActive:      redditInfo.AccountsActive,
	}
	err := redditModel.Insert(ctx, pg.db, boil.Infer())
	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") { // Ignore duplicate entries
			return nil
		}
		return err
	}
	log.Infof("Added a new reddit info received at %s, Subscribers %d, Active Users %d",
		redditInfo.Date.Format(dateMiliTemplate), redditInfo.Subscribers, redditInfo.AccountsActive)
	return nil
}

func (pg *PgDb) LastRedditEntryTime() (entryTime time.Time) {
	rows := pg.db.QueryRow(lastRedditEntryTime)
	_ = rows.Scan(&entryTime)
	return
}

func (pg *PgDb) RedditInfoCount(ctx context.Context) (int64, error) {
	return models.RedditInfos().Count(ctx, pg.db)
}

func (pg *PgDb) RedditInfos(ctx context.Context, offtset int, limit int) ([]reddit.RedditInfo, error) {
	redditInfoSlice, err := models.RedditInfos(qm.OrderBy(fmt.Sprintf("%s DESC", models.RedditInfoColumns.Date)),
		qm.Offset(offtset), qm.Limit(limit)).All(ctx, pg.db)
	if err != nil {
		return nil, err
	}

	var result []reddit.RedditInfo
	for _, record := range redditInfoSlice {
		result = append(result,reddit.RedditInfo{
			Date:             record.Date,
			Subscribers:        record.Subscribers,
			AccountsActive:                record.AccountsActive,
		})
	}
	return result, nil
}
