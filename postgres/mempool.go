package postgres

import (
	"context"
	"strings"

	"github.com/raedahgroup/dcrextdata/mempool"
	"github.com/raedahgroup/dcrextdata/postgres/models"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
)

func (pg PgDb) StoreMempool(ctx context.Context, mempoolDto mempool.Mempool) error {
	mempoolModel := mempoolDtoToModel(mempoolDto)
	err := mempoolModel.Insert(ctx, pg.db, boil.Infer())
	if err != nil {
		if !strings.Contains(err.Error(), "unique constraint") { // Ignore duplicate entries
			return err
		}
		return err
	}
	log.Infof("Added mempool entry, Block Height: %6d, Tx Count: %2d, Size: %2d, Timestamp: %s",
		mempoolDto.BlockHeight, mempoolDto.NumberOfTransactions, mempoolDto.Size, mempoolDto.BlockInternalTime.Format(dateTemplate))
	return nil
}

func mempoolDtoToModel(mempoolDto mempool.Mempool) models.Mempool {
	return models.Mempool{
		FirstSeenTime:        null.Int64From(mempoolDto.FirstSeenTime.Unix()),
		BlockReceiveTime:     null.Int64From(mempoolDto.BlockReceiveTime.Unix()),
		BlockInternalTime:    null.Int64From(mempoolDto.BlockInternalTime.Unix()),
		BlockHeight:          int(mempoolDto.BlockHeight),
		Size:                 null.IntFrom(int(mempoolDto.Size)),
		NumberOfTransactions: null.IntFrom(mempoolDto.NumberOfTransactions),
		BlockHash:            mempoolDto.BlockHash,
	}
}

func (pg *PgDb) LastMempoolBlockHeight() (height int64, err error) {
	rows := pg.db.QueryRow(lastMempoolBlockHeight)
	err = rows.Scan(&height)
	return
}
