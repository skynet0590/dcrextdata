package postgres

import (
	"context"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"strings"
	"time"

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
	}
	log.Infof("Added mempool entry, Timestamp: %s, Tx Count: %2d, Size: %6d, Total Fee: %f",
		mempoolDto.Time.Format(dateTemplate), mempoolDto.NumberOfTransactions, mempoolDto.Size, mempoolDto.TotalFee)
	return nil
}

func mempoolDtoToModel(mempoolDto mempool.Mempool) models.Mempool {
	return models.Mempool{
		Time:                 mempoolDto.Time.Unix(),
		FirstSeenTime:        null.Int64From(mempoolDto.FirstSeenTime.Unix()),
		Size:                 null.IntFrom(int(mempoolDto.Size)),
		NumberOfTransactions: null.IntFrom(mempoolDto.NumberOfTransactions),
		Revocations:          null.IntFrom(mempoolDto.Revocations),
		Tickets:              null.IntFrom(mempoolDto.Tickets),
		Voters:               null.IntFrom(mempoolDto.Voters),
		Total:                null.Float64From(mempoolDto.Total),
		TotalFee:             null.Float64From(mempoolDto.TotalFee),
	}
}

func (pg *PgDb) LastMempoolBlockHeight() (height int64, err error) {
	rows := pg.db.QueryRow(lastMempoolBlockHeight)
	err = rows.Scan(&height)
	return
}

func (pg *PgDb) MempoolCount(ctx context.Context) (int64, error) {
	return models.Mempools().Count(ctx, pg.db)
}

func (pg *PgDb) Mempools(ctx context.Context, offtset int, limit int) ([]mempool.Mempool, error) {
	mempoolSlice, err := models.Mempools(qm.Offset(offtset), qm.Limit(limit)).All(ctx, pg.db)
	if err != nil {
		return nil, err
	}
	var result []mempool.Mempool
	for _, m := range mempoolSlice {
		result = append(result, mempool.Mempool{
			TotalFee:             m.TotalFee.Float64,
			FirstSeenTime:        int64ToTime(m.FirstSeenTime.Int64),
			Total:                m.Total.Float64,
			Voters:               m.Voters.Int,
			Tickets:              m.Tickets.Int,
			Revocations:          m.Revocations.Int,
			Time:                 time.Unix(m.Time, 0),
			Size:                 int32(m.Size.Int),
			NumberOfTransactions: m.NumberOfTransactions.Int,
		})
	}
	return result, nil
}

func (pg *PgDb) SaveBlock(ctx context.Context, block mempool.Block) error  {
	blockModel := blockDtoToModel(block)
	err := blockModel.Insert(ctx, pg.db, boil.Infer())
	if err != nil {
		if !strings.Contains(err.Error(), "unique constraint") { // Ignore duplicate entries
			return err
		}
	}
	log.Infof("New block received, Timestamp: %s, Height: %d, Hash: ...%s",
		block.BlockInternalTime.Format(dateTemplate), block.BlockHeight, block.BlockHash[len(block.BlockHash) - 23:])
	return nil
}

func blockDtoToModel(block mempool.Block) models.Block {
	return models.Block{
		Height: int(block.BlockHeight),
		Hash: null.StringFrom(block.BlockHash),
		InternalTimestamp: null.Int64From(block.BlockInternalTime.Unix()),
		ReceiveTime: null.Int64From(block.BlockInternalTime.Unix()),
	}
}

func (pg *PgDb) BlockCount(ctx context.Context) (int64, error) {
	return models.Blocks().Count(ctx, pg.db)
}

func (pg *PgDb) Blocks(ctx context.Context, offset int, limit int) ([]mempool.Block, error) {
	blockSlice, err := models.Blocks(qm.Offset(offset), qm.Limit(limit)).All(ctx, pg.db)
	if err != nil {
		return nil, err
	}

	var blocks []mempool.Block

	for _, block := range blockSlice {
		blocks = append(blocks, mempool.Block{
			BlockHash:block.Hash.String,
			BlockHeight:uint32(block.Height),
			BlockInternalTime:int64ToTime(block.InternalTimestamp.Int64),
			BlockReceiveTime:int64ToTime(block.ReceiveTime.Int64),
		})
	}

	return blocks, nil
}
