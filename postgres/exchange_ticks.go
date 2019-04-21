// Copyright (c) 2018-2019 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/raedahgroup/dcrextdata/postgres/models"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"

	"github.com/raedahgroup/dcrextdata/exchanges/ticks"
)

const (
	NegativeFiveMin = time.Duration(-5) * time.Minute
	NegativeOneHour = time.Duration(-1) * time.Hour
	NegativeOneDay  = time.Duration(-24) * time.Hour
)

var (
	ErrNonConsecutiveTicks = errors.New("postgres/exchanges: Non consecutive exchange ticks")
	zeroTime               time.Time
)

func (pg *PgDb) RegisterExchange(ctx context.Context, exchange ticks.ExchangeData) (time.Time, time.Time, time.Time, error) {
	xch, err := models.Exchanges(models.ExchangeWhere.Name.EQ(exchange.Name)).One(ctx, pg.db)
	if err == sql.ErrNoRows {
		newXch := models.Exchange{
			Name:                 exchange.Name,
			URL:                  exchange.WebsiteURL,
			TickShortInterval:    int(exchange.ShortInterval.Seconds()),
			TickLongInterval:     int(exchange.LongInterval.Seconds()),
			TickHistoricInterval: int(exchange.HistoricInterval.Seconds()),
		}
		err = newXch.Insert(ctx, pg.db, boil.Infer())
		return zeroTime, zeroTime, zeroTime, err
	} else if err != nil {
		return zeroTime, zeroTime, zeroTime, err
	}
	var shortTime, longTime, historicTime time.Time
	timeAsc := qm.OrderBy("time desc")
	lastShort, err := models.ExchangeTicks(qm.Expr(models.ExchangeTickWhere.ExchangeID.EQ(xch.ID), models.ExchangeTickWhere.Interval.EQ(ticks.IntervalShort), timeAsc)).One(ctx, pg.db)
	if err == nil {
		shortTime = lastShort.Time
	}
	lastLong, err := models.ExchangeTicks(qm.Expr(models.ExchangeTickWhere.ExchangeID.EQ(xch.ID), models.ExchangeTickWhere.Interval.EQ(ticks.IntervalLong), timeAsc)).One(ctx, pg.db)
	if err == nil {
		longTime = lastLong.Time
	}
	lastHistoric, err := models.ExchangeTicks(qm.Expr(models.ExchangeTickWhere.ExchangeID.EQ(xch.ID), models.ExchangeTickWhere.Interval.EQ(ticks.IntervalHistoric), timeAsc)).One(ctx, pg.db)
	if err == nil {
		historicTime = lastHistoric.Time
	}
	if err != nil && err == sql.ErrNoRows {
		err = nil
	}

	// log.Debugf("Exchange %s, %v, %v, %v", exchange.Name, shortTime.UTC(), longTime.UTC(), historicTime.UTC())
	return shortTime, longTime, historicTime, err
}

// StoreExchangeTicks
func (pg *PgDb) StoreExchangeTicks(ctx context.Context, name string, interval time.Duration, intervalString string, pair string, ticks []ticks.Tick) (time.Time, error) {
	if len(ticks) == 0 {
		return zeroTime, fmt.Errorf("No ticks recieved for %s", name)
	}

	var lastTime time.Time
	lastTick, err := models.ExchangeTicks(qm.OrderBy(models.ExchangeTickColumns.Time)).One(ctx, pg.db)

	if err == sql.ErrNoRows {
		lastTime = ticks[0].Time.Add(-1)
	} else if err != nil {
		return lastTime, err
	} else {
		lastTime = lastTick.Time
	}

	xch, err := models.Exchanges(models.ExchangeWhere.Name.EQ(name)).One(ctx, pg.db)

	if err != nil {
		return lastTime, nil
	}

	for _, tick := range ticks {
		// if lastTime != tick.Time.Add(NegativeFiveMin) {
		// 	return lastTime, ErrNonConsecutiveTicks
		// }

		xcTick := tickToExchangeTick(xch.ID, pair, intervalString, tick)

		err = xcTick.Insert(ctx, pg.db, boil.Infer())

		if err != nil && !strings.Contains(err.Error(), "unique constraint") {
			return lastTime, err
		}
		lastTime = xcTick.Time
	}

	log.Infof("Added all exchange ticks from %s to %s", ticks[0].Time.String(),
		ticks[len(ticks)-1].Time.String())
	return lastTime, nil
}

func tickToExchangeTick(exchangeID int, pair string, interval string, tick ticks.Tick) *models.ExchangeTick {
	return &models.ExchangeTick{
		ExchangeID:   exchangeID,
		High:         tick.High,
		Low:          tick.Low,
		Open:         tick.Open,
		Close:        tick.Close,
		Volume:       tick.Volume,
		Time:         tick.Time,
		CurrencyPair: pair,
		Interval:     interval,
	}
}
