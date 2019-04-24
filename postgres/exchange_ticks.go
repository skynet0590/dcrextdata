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
	if err != nil {
		if err == sql.ErrNoRows {
			newXch := models.Exchange{
				Name:                 exchange.Name,
				URL:                  exchange.WebsiteURL,
				TickShortInterval:    int(exchange.ShortInterval.Seconds()),
				TickLongInterval:     int(exchange.LongInterval.Seconds()),
				TickHistoricInterval: int(exchange.HistoricInterval.Seconds()),
			}
			err = newXch.Insert(ctx, pg.db, boil.Infer())
		}
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

	xch, err := models.Exchanges(models.ExchangeWhere.Name.EQ(name)).One(ctx, pg.db)
	if err != nil {
		return zeroTime, err
	}

	var lastTime time.Time
	lastTick, err := models.ExchangeTicks(models.ExchangeTickWhere.ExchangeID.EQ(xch.ID),
		models.ExchangeTickWhere.Interval.EQ(intervalString),
		models.ExchangeTickWhere.CurrencyPair.EQ(pair),
		qm.OrderBy(models.ExchangeTickColumns.Time)).One(ctx, pg.db)

	if err == sql.ErrNoRows {
		lastTime = ticks[0].Time.Add(-interval)
	} else if err != nil {
		return lastTime, err
	} else {
		lastTime = lastTick.Time
	}

	firstTime := ticks[0].Time
	added := 0
	for _, tick := range ticks {
		// if tick.Time.Unix() <= lastTime.Unix() {
		// 	continue
		// }
		xcTick := tickToExchangeTick(xch.ID, pair, intervalString, tick)
		err = xcTick.Insert(ctx, pg.db, boil.Infer())
		if err != nil && !strings.Contains(err.Error(), "unique constraint") {
			return lastTime, err
		}
		lastTime = xcTick.Time
		added++
	}

	if added == 0 {
		log.Infof("No new ticks for %s(%s)", name, intervalString)
	} else if added == 1 {
		log.Infof("Stored an exchange tick for %s(%s) at %v", name, intervalString, firstTime)
	} else {
		log.Infof("Stored %d exchange ticks for %s(%s) from %v to %v", added, name, intervalString,
			firstTime, lastTime)
	}
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
		Time:         tick.Time.UTC(),
		CurrencyPair: pair,
		Interval:     interval,
	}
}
