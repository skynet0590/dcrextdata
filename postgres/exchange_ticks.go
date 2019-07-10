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
				Name: exchange.Name,
				URL:  exchange.WebsiteURL,
			}
			err = newXch.Insert(ctx, pg.db, boil.Infer())
		}
		return zeroTime, zeroTime, zeroTime, err
	}
	var shortTime, longTime, historicTime time.Time
	toMin := func(t time.Duration) int {
		return int(t.Minutes())
	}
	timeDesc := qm.OrderBy("time desc")
	lastShort, err := models.ExchangeTicks(qm.Expr(models.ExchangeTickWhere.ExchangeID.EQ(xch.ID), models.ExchangeTickWhere.Interval.EQ(toMin(exchange.ShortInterval)), timeDesc)).One(ctx, pg.db)
	if err == nil {
		shortTime = lastShort.Time
	}
	lastLong, err := models.ExchangeTicks(qm.Expr(models.ExchangeTickWhere.ExchangeID.EQ(xch.ID), models.ExchangeTickWhere.Interval.EQ(toMin(exchange.LongInterval)), timeDesc)).One(ctx, pg.db)
	if err == nil {
		longTime = lastLong.Time
	}
	lastHistoric, err := models.ExchangeTicks(qm.Expr(models.ExchangeTickWhere.ExchangeID.EQ(xch.ID), models.ExchangeTickWhere.Interval.EQ(toMin(exchange.HistoricInterval)), timeDesc)).One(ctx, pg.db)
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
func (pg *PgDb) StoreExchangeTicks(ctx context.Context, name string, interval int, pair string, ticks []ticks.Tick) (time.Time, error) {
	if len(ticks) == 0 {
		return zeroTime, fmt.Errorf("No ticks recieved for %s", name)
	}

	exchange, err := models.Exchanges(models.ExchangeWhere.Name.EQ(name)).One(ctx, pg.db)
	if err != nil {
		return zeroTime, err
	}

	var lastTime time.Time
	lastTick, err := models.ExchangeTicks(models.ExchangeTickWhere.ExchangeID.EQ(exchange.ID),
		models.ExchangeTickWhere.Interval.EQ(interval),
		models.ExchangeTickWhere.CurrencyPair.EQ(pair),
		qm.OrderBy(models.ExchangeTickColumns.Time)).One(ctx, pg.db)

	if err == sql.ErrNoRows {
		lastTime = ticks[0].Time.Add(-time.Duration(interval))
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
		xcTick := tickToExchangeTick(exchange.ID, pair, interval, tick)
		err = xcTick.Insert(ctx, pg.db, boil.Infer())
		if err != nil && !strings.Contains(err.Error(), "unique constraint") {
			return lastTime, err
		}
		lastTime = xcTick.Time
		added++
	}

	if added == 0 {
		log.Infof("No new ticks on %s(%dm) for", name, pair, interval)
	} else if added == 1 {
		log.Infof("%-9s %7s, received %6dm ticks, storing      1 entries %s", name, pair,
			interval, firstTime.Format(dateTemplate))

		/*log.Infof("%10s %7s, received      1  tick %14s %s", name, pair,
		fmt.Sprintf("(%dm)", interval), firstTime.Format(dateTemplate))*/
	} else {
		log.Infof("%-9s %7s, received %6dm ticks, storing %6v entries %s to %s", name, pair,
			interval, added, firstTime.Format(dateTemplate), lastTime.Format(dateTemplate))

		/*log.Infof("%10s %7s, received %6v ticks %14s %s to %s",
		name, pair, added, fmt.Sprintf("(%dm each)", interval), firstTime.Format(dateTemplate), lastTime.Format(dateTemplate))*/
	}
	return lastTime, nil
}

// AllExchange fetches a slice of all exchange from the db
func (pg *PgDb) AllExchange(ctx context.Context) (models.ExchangeSlice, error) {
	exchangeSlice, err := models.Exchanges().All(ctx, pg.db)
	return exchangeSlice, err
}

// FetchExchangeTicks fetches a slice exchange ticks of the supplied exchange name
func (pg *PgDb) FetchExchangeTicks(ctx context.Context, currencyPair, name string, offset int, limit int) ([]ticks.TickDto, int64, error) {
	exchange, err := models.Exchanges(models.ExchangeWhere.Name.EQ(name)).One(ctx, pg.db)
	if err != nil {
		return nil, 0, err
	}

	var idQuery qm.QueryMod
	if currencyPair != "" {
		idQuery = qm.Where("(exchange_id=? and currency_pair=?)", exchange.ID, currencyPair)
	} else {
		idQuery = models.ExchangeTickWhere.ExchangeID.EQ(exchange.ID)
	}

	exchangeTickSlice, err := models.ExchangeTicks(qm.Load("Exchange"), idQuery, qm.Limit(limit), qm.Offset(offset), qm.OrderBy("time DESC")).All(ctx, pg.db)
	if err != nil {
		return nil, 0, err
	}

	exchangeTickSliceCount, err := models.ExchangeTicks(qm.Load("Exchange"), idQuery, qm.Limit(limit), qm.Offset(offset)).Count(ctx, pg.db)
	if err != nil {
		return nil, 0, err
	}

	tickDtos := []ticks.TickDto{}
	for _, tick := range exchangeTickSlice {
		tickDtos = append(tickDtos, ticks.TickDto{
			ExchangeID:   tick.ExchangeID,
			Interval:     tick.Interval,
			CurrencyPair: tick.CurrencyPair,
			Time:         tick.Time.UTC(),
			Close:        tick.Close,
			ExchangeName: tick.R.Exchange.Name,
			High:         tick.High,
			Low:          tick.Low,
			Open:         tick.Open,
			Volume:       tick.Volume,
		})
	}

	return tickDtos, exchangeTickSliceCount, err
}

// FetchExchangeTicks fetches a slice exchange ticks of the supplied exchange name
// todo impliment sorting for Exchange ticks as it is currently been sorted by time
func (pg *PgDb) AllExchangeTicks(ctx context.Context, currencyPair string, offset int, limit int) ([]ticks.TickDto, int64, error) {
	var exchangeTickSlice models.ExchangeTickSlice
	var exchangeTickSliceCount int64
	var err error

	if currencyPair != "" {
		idQuery := models.ExchangeTickWhere.CurrencyPair.EQ(currencyPair)
		exchangeTickSlice, err = models.ExchangeTicks(qm.Load("Exchange"), idQuery, qm.Limit(limit), qm.Offset(offset), qm.OrderBy("time DESC")).All(ctx, pg.db)
		if err != nil {
			return nil, 0, err
		}

		exchangeTickSliceCount, err = models.ExchangeTicks(qm.Load("Exchange"), idQuery).Count(ctx, pg.db)
		if err != nil {
			return nil, 0, err
		}
	} else {
		exchangeTickSlice, err = models.ExchangeTicks(qm.Load("Exchange"), qm.Limit(limit), qm.Offset(offset), qm.OrderBy("time DESC")).All(ctx, pg.db)
		if err != nil {
			return nil, 0, err
		}

		exchangeTickSliceCount, err = models.ExchangeTicks(qm.Load("Exchange")).Count(ctx, pg.db)
		if err != nil {
			return nil, 0, err
		}
	}

	tickDtos := []ticks.TickDto{}
	for _, tick := range exchangeTickSlice {
		tickDtos = append(tickDtos, ticks.TickDto{
			ExchangeID:   tick.ExchangeID,
			Interval:     tick.Interval,
			CurrencyPair: tick.CurrencyPair,
			Time:         tick.Time.UTC(),
			Close:        tick.Close,
			ExchangeName: tick.R.Exchange.Name,
			High:         tick.High,
			Low:          tick.Low,
			Open:         tick.Open,
			Volume:       tick.Volume,
		})
	}

	return tickDtos, exchangeTickSliceCount, err
}

func (pg *PgDb) AllExchangeTicksCurrencyPair(ctx context.Context) ([]ticks.TickDtoCP, error) {
	exchangeTickCPSlice, err := models.ExchangeTicks(qm.Select("currency_pair"), qm.GroupBy("currency_pair")).All(ctx, pg.db)

	if err != nil {
		return nil, err
	}

	TickDtoCP := []ticks.TickDtoCP{}
	for _, cp := range exchangeTickCPSlice {
		TickDtoCP = append(TickDtoCP, ticks.TickDtoCP{
			CurrencyPair: cp.CurrencyPair,
		})
	}

	return TickDtoCP, err
}

func tickToExchangeTick(exchangeID int, pair string, interval int, tick ticks.Tick) *models.ExchangeTick {
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
