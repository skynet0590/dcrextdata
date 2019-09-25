package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/raedahgroup/dcrextdata/commstats"
	"github.com/raedahgroup/dcrextdata/postgres/models"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
)

func (pg *PgDb) StoreCommStat(ctx context.Context, stat commstats.CommStat) error {
	redditJson, err := json.Marshal(stat.RedditStats)
	if err != nil {
		return fmt.Errorf("error in saving stat, cannot decode reddit stat, %s", err.Error())
	}
	commStat := models.CommStat{
		Date:               stat.Date,
		RedditStat:         string(redditJson),
		TwitterFollowers:   stat.TwitterFollowers,
		YoutubeSubscribers: stat.YoutubeSubscribers,
		GithubStars:        stat.GithubStars,
		GithubFolks:        stat.GithubFolks,
	}

	err = commStat.Insert(ctx, pg.db, boil.Infer())
	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") { // Ignore duplicate entries
			return nil
		}
	}

	return err
}

func (pg *PgDb) LastCommStatEntry() (entryTime time.Time) {
	rows := pg.db.QueryRow(lastCommStatEntryTime)
	_ = rows.Scan(&entryTime)
	return
}

func (pg *PgDb) CommStatCount(ctx context.Context) (int64, error) {
	return models.CommStats().Count(ctx, pg.db)
}

func (pg *PgDb) CommStats(ctx context.Context, offtset int, limit int) ([]commstats.CommStat, error) {
	commStatSlices, err := models.CommStats(qm.OrderBy(fmt.Sprintf("%s DESC", models.CommStatColumns.Date)),
		qm.Offset(offtset), qm.Limit(limit)).All(ctx, pg.db)
	if err != nil {
		return nil, err
	}

	var result []commstats.CommStat
	for _, record := range commStatSlices {
		stat := commstats.CommStat{
			Date:                 record.Date,
			TwitterFollowers:     record.TwitterFollowers,
			YoutubeSubscribers:   record.YoutubeSubscribers,
			GithubStars:          record.GithubStars,
			GithubFolks:          record.GithubFolks,
		}
		if record.RedditStat != "" {
			var redditStat map[string]commstats.Reddit
			err := json.Unmarshal([]byte(record.RedditStat), &redditStat)
			if err != nil {
				return nil, fmt.Errorf("cannot decode reddit data, %s", err.Error())
			}
			stat.RedditStats = redditStat
		}

		result = append(result, stat)
	}
	return result, nil
}
