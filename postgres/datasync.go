package postgres

import (
	"context"
	"fmt"
	"github.com/raedahgroup/dcrextdata/postgres/models"
)

func (pg *PgDb) TableNames() []string {
	return []string{
		models.TableNames.Vote,
		models.TableNames.Block,
		models.TableNames.Mempool,
		models.TableNames.Exchange,
		models.TableNames.ExchangeTick,
		models.TableNames.VSP,
		models.TableNames.VSPTick,
		models.TableNames.PowData,
	}
}

func (pg *PgDb) LastEntry(ctx context.Context, tableName string, receiver interface{}) error {
	var columnName string
	switch tableName {
	case models.TableNames.Exchange:
		columnName = models.ExchangeColumns.ID
		break
	case models.TableNames.ExchangeTick:
		columnName = models.ExchangeTickColumns.ID
		break
	case models.TableNames.Mempool:
		columnName = models.MempoolColumns.Time
		break
	case models.TableNames.Block:
		columnName = models.BlockColumns.Height
		break
	case models.TableNames.Vote:
		columnName = models.VoteColumns.ReceiveTime
		break
	case models.TableNames.PowData:
		columnName = models.PowDatumColumns.Time
		break
	case models.TableNames.VSP:
		columnName = models.VSPColumns.ID
		break
	case models.TableNames.VSPTick:
		columnName = models.VSPTickColumns.ID
		break
	case models.TableNames.Reddit:
		columnName = models.RedditColumns.Date
		break
	case models.TableNames.Twitter:
		columnName = models.TwitterColumns.Date
		break
	case models.TableNames.Github:
		columnName = models.GithubColumns.Date
		break
	case models.TableNames.Youtube:
		columnName = models.YoutubeColumns.Date
		break
	case models.TableNames.NetworkSnapshot:
		columnName = models.NetworkSnapshotColumns.Timestamp
		break

	}

	rows := pg.db.QueryRow(fmt.Sprintf("SELECT %s FROM %s ORDER BY %s DESC LIMIT 1", columnName, tableName, columnName))
	err := rows.Scan(receiver)
	if err != nil || receiver == nil {
		receiver = 0
	}
	return err

}
