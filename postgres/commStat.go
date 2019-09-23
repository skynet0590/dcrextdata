package postgres

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/raedahgroup/dcrextdata/commstats"
	"github.com/raedahgroup/dcrextdata/postgres/models"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
)

func (pg *PgDb) StoreCommStat(ctx context.Context, stat commstats.CommStat) error {
	commStat := models.CommStat{
		Date:                 stat.Date,
		RedditSubscribers:    stat.RedditSubscribers,
		RedditAccountsActive: stat.RedditAccountsActive,
		TwitterFollowers:     stat.TwitterFollowers,
		YoutubeSubscribers:   stat.YoutubeSubscribers,
		GithubStars:          stat.GithubStars,
		GithubFolks:          stat.GithubFolks,
	}
	err := commStat.Insert(ctx, pg.db, boil.Infer())
	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") { // Ignore duplicate entries
			return nil
		}
		return err
	}
	log.Infof("Added a new Community stat entry received at %s, \n"+
		"\t\t  Reddit Subscribers  %d\n"+
		"\t\t  Reddit Active Users %d\n"+
		"\t\t  Twitter Followers   %d\n"+
		"\t\t  Youtube Subscribers   %d",
		stat.Date.Format(dateMiliTemplate), stat.RedditSubscribers, stat.RedditAccountsActive, stat.TwitterFollowers,
		stat.YoutubeSubscribers)
	return nil
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
	redditInfoSlice, err := models.CommStats(qm.OrderBy(fmt.Sprintf("%s DESC", models.CommStatColumns.Date)),
		qm.Offset(offtset), qm.Limit(limit)).All(ctx, pg.db)
	if err != nil {
		return nil, err
	}

	var result []commstats.CommStat
	for _, record := range redditInfoSlice {
		result = append(result, commstats.CommStat{
			Date:                 record.Date,
			RedditSubscribers:    record.RedditSubscribers,
			RedditAccountsActive: record.RedditAccountsActive,
			TwitterFollowers:     record.TwitterFollowers,
			YoutubeSubscribers:   record.YoutubeSubscribers,
			GithubStars:          record.GithubStars,
			GithubFolks:          record.GithubFolks,
		})
	}
	return result, nil
}
