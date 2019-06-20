package postgres

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/raedahgroup/dcrextdata/postgres/models"
	"github.com/raedahgroup/dcrextdata/pow"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
)

func (pg *PgDb) LastPowEntryTime(source string) (time int64) {
	rows := pg.db.QueryRow(LastPowEntryTime, source)
	_ = rows.Scan(&time)
	return
}

//
func (pg *PgDb) AddPowData(ctx context.Context, data []pow.PowData) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	added := 0
	for _, d := range data {
		powModel, err := responseToPowModel(d)
		if err != nil {
			return err
		}

		err = powModel.Insert(ctx, pg.db, boil.Infer())
		if err != nil {
			if !strings.Contains(err.Error(), "unique constraint") { // Ignore duplicate entries
				return err
			}
		}
		added++
	}
	if len(data) == 1 {
		log.Infof("Added %4d PoW   entry from %10s %s", added, data[0].Source, UnixTimeToString(data[0].Time))
	} else if len(data) > 1 {
		last := data[len(data)-1]
		log.Infof("Added %4d PoW entries from %10s %s to %s",
			added, last.Source, UnixTimeToString(data[0].Time), UnixTimeToString(last.Time))
	}

	return nil
}

func responseToPowModel(data pow.PowData) (models.PowDatum, error) {
	return models.PowDatum{
		BTCPrice:          null.StringFrom(fmt.Sprint(data.BtcPrice)),
		CoinPrice:         null.StringFrom(fmt.Sprint(data.CoinPrice)),
		NetworkDifficulty: null.Float64From(data.NetworkDifficulty),
		NetworkHashrate:   null.StringFrom(fmt.Sprint(data.NetworkHashrate)),
		PoolHashrate:      null.StringFrom(fmt.Sprint(data.PoolHashrate)),
		Source:            data.Source,
		Time:              int(data.Time),
		Workers:           null.IntFrom(int(data.Workers)),
	}, nil
}


func (pg *PgDb) FetchPowData(ctx context.Context, offset int, limit int) ([]pow.PowDataDto, error) {
	powDatum, err := models.PowData(qm.Offset(offset), qm.Limit(limit)).All(ctx, pg.db)
	if err != nil {
		return  nil, err
	}

	var result []pow.PowDataDto
	for _, item := range powDatum {
		networkHashRate, err := strconv.ParseInt(item.NetworkHashrate.String, 10, 64)
		if err != nil {
			return nil, err
		}

		poolHashRate, err := strconv.ParseFloat(item.PoolHashrate.String, 10)
		if err != nil {
			return nil, err
		}

		networkHashrate, err := strconv.ParseFloat(item.NetworkHashrate.String, 10)
		if err != nil {
			return nil, err
		}

		coinPrice, err := strconv.ParseFloat(item.CoinPrice.String, 10)
		if err != nil {
			return nil, err
		}

		bTCPrice, err := strconv.ParseFloat(item.BTCPrice.String, 10)
		if err != nil {
			return nil, err
		}

		result = append(result, pow.PowDataDto{
			Time: time.Unix(int64(item.Time), 0).UTC(),
			NetworkHashrate: networkHashRate,
			PoolHashrate:poolHashRate,
			Workers: int64(item.Workers.Int),
			Source: item.Source,
			NetworkDifficulty: networkHashrate,
			CoinPrice: coinPrice,
			BtcPrice: bTCPrice,
		})
	}

	return result, nil
}

func (pg *PgDb) CountPowData(ctx context.Context) (int64, error) {
	return models.PowData().Count(ctx, pg.db)
}

func (pg *PgDb) FetchPowDataBySource(ctx context.Context, source string, offset int, limit int) ([]pow.PowDataDto, error) {
	powDatum, err := models.PowData(models.PowDatumWhere.Source.EQ(source), qm.Offset(offset), qm.Limit(limit)).All(ctx, pg.db)
	if err != nil {
		return  nil, err
	}

	var result []pow.PowDataDto
	for _, item := range powDatum {
		networkHashRate, err := strconv.ParseInt(item.NetworkHashrate.String, 10, 64)
		if err != nil {
			return nil, err
		}

		poolHashRate, err := strconv.ParseFloat(item.PoolHashrate.String, 10)
		if err != nil {
			return nil, err
		}

		networkHashrate, err := strconv.ParseFloat(item.NetworkHashrate.String, 10)
		if err != nil {
			return nil, err
		}

		coinPrice, err := strconv.ParseFloat(item.CoinPrice.String, 10)
		if err != nil {
			return nil, err
		}

		bTCPrice, err := strconv.ParseFloat(item.BTCPrice.String, 10)
		if err != nil {
			return nil, err
		}

		result = append(result, pow.PowDataDto{
			Time: time.Unix(int64(item.Time), 0).UTC(),
			NetworkHashrate: networkHashRate,
			PoolHashrate:poolHashRate,
			Workers: int64(item.Workers.Int),
			Source: item.Source,
			NetworkDifficulty: networkHashrate,
			CoinPrice: coinPrice,
			BtcPrice: bTCPrice,
		})
	}

	return result, nil
}

func (pg *PgDb) CountPowDataBySource(ctx context.Context, source string) (int64, error) {
	return models.PowData(models.PowDatumWhere.Source.EQ(source)).Count(ctx, pg.db)
}

