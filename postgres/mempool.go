package postgres

import (
	"context"
	"strings"
	"time"

	"github.com/raedahgroup/dcrextdata/mempool"
	"github.com/raedahgroup/dcrextdata/postgres/models"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
)

func (pg PgDb) StoreMempool(ctx context.Context, mempoolDto mempool.Mempool) error {
	mempoolModel := mempoolDtoToModel(mempoolDto)
	err := mempoolModel.Insert(ctx, pg.db, boil.Infer())
	if err != nil {
		if !strings.Contains(err.Error(), "unique constraint") { // Ignore duplicate entries
			return err
		}
	}
	//  tx count 76, total size 54205 B, fees 0.00367100
	log.Infof("Added mempool entry at %s, tx count %2d, total size: %6d B, Total Fee: %010.8f",
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

func (pg *PgDb) Mempools(ctx context.Context, offtset int, limit int) ([]mempool.MempoolDto, error) {
	mempoolSlice, err := models.Mempools(qm.OrderBy(models.MempoolColumns.Time), qm.Offset(offtset), qm.Limit(limit)).All(ctx, pg.db)
	if err != nil {
		return nil, err
	}
	var result []mempool.MempoolDto
	for _, m := range mempoolSlice {
		result = append(result, mempool.MempoolDto{
			TotalFee:             m.TotalFee.Float64,
			FirstSeenTime:        int64ToTime(m.FirstSeenTime.Int64).Format(dateMiliTemplate),
			Total:                m.Total.Float64,
			Voters:               m.Voters.Int,
			Tickets:              m.Tickets.Int,
			Revocations:          m.Revocations.Int,
			Time:                 time.Unix(m.Time, 0).Format(dateTemplate),
			Size:                 int32(m.Size.Int),
			NumberOfTransactions: m.NumberOfTransactions.Int,
		})
	}
	return result, nil
}

func (pg *PgDb) SaveBlock(ctx context.Context, block mempool.Block) error {
	blockModel := blockDtoToModel(block)
	err := blockModel.Insert(ctx, pg.db, boil.Infer())
	if err != nil {
		if !strings.Contains(err.Error(), "unique constraint") { // Ignore duplicate entries
			return err
		}
	}
	log.Infof("New block received at %s, Height: %d, Hash: ...%s",
		block.BlockInternalTime.Format(dateMiliTemplate), block.BlockHeight, block.BlockHash[len(block.BlockHash)-23:])
	return nil
}

func blockDtoToModel(block mempool.Block) models.Block {
	return models.Block{
		Height:            int(block.BlockHeight),
		Hash:              null.StringFrom(block.BlockHash),
		InternalTimestamp: null.Int64From(block.BlockInternalTime.Unix()),
		ReceiveTime:       null.Int64From(block.BlockInternalTime.Unix()),
	}
}

func (pg *PgDb) BlockCount(ctx context.Context) (int64, error) {
	return models.Blocks().Count(ctx, pg.db)
}

func (pg *PgDb) Blocks(ctx context.Context, offset int, limit int) ([]mempool.Block, error) {
	blockSlice, err := models.Blocks(qm.OrderBy(models.BlockColumns.ReceiveTime), qm.Offset(offset), qm.Limit(limit)).All(ctx, pg.db)
	if err != nil {
		return nil, err
	}

	var blocks []mempool.Block

	for _, block := range blockSlice {
		blocks = append(blocks, mempool.Block{
			BlockHash:         block.Hash.String,
			BlockHeight:       uint32(block.Height),
			BlockInternalTime: int64ToTime(block.InternalTimestamp.Int64),
			BlockReceiveTime:  int64ToTime(block.ReceiveTime.Int64),
		})
	}

	return blocks, nil
}

func (pg *PgDb) SaveVote(ctx context.Context, vote mempool.Vote) error {
	voteModel := models.Vote{
		Hash:        vote.Hash,
		VotingOn:    null.Int64From(int64(vote.VotingOn)),
		ReceiveTime: null.Int64From(vote.ReceiveTime.Unix()),
		ValidatorID: null.IntFrom(vote.ValidatorId),
	}
	err := voteModel.Insert(ctx, pg.db, boil.Infer())
	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") { // Ignore duplicate entries
			return nil
		}
		return err
	}
	log.Infof("New vote received at %s for %d, Validator Id %d, Hash ...%s",
		vote.ReceiveTime.Format(dateMiliTemplate), vote.VotingOn, vote.ValidatorId, vote.Hash[len(vote.Hash)-23:])
	return nil
}

func (pg *PgDb) Votes(ctx context.Context, offset int, limit int) ([]mempool.Vote, error) {
	voteSlice, err := models.Votes(qm.OrderBy(models.VoteColumns.ReceiveTime), qm.Offset(offset), qm.Limit(limit)).All(ctx, pg.db)
	if err != nil {
		return nil, err
	}

	var votes []mempool.Vote
	for _, vote := range voteSlice {
		votes = append(votes, mempool.Vote{
			Hash:        vote.Hash,
			ReceiveTime: int64ToTime(vote.ReceiveTime.Int64),
			VotingOn:    vote.VotingOn.Int64,
			ValidatorId: vote.ValidatorID.Int,
		})
	}

	return votes, nil
}

func (pg *PgDb) VotesCount(ctx context.Context) (int64, error) {
	return models.Votes().Count(ctx, pg.db)
}
