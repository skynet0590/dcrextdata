// Copyright (c) 2018-2019 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/raedahgroup/dcrextdata/cache"
	"github.com/raedahgroup/dcrextdata/datasync"
	"github.com/raedahgroup/dcrextdata/postgres/models"
	"github.com/raedahgroup/dcrextdata/vsp"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"github.com/volatiletech/sqlboiler/types"
)

var (
	vspTickExistsErr = fmt.Errorf("VSPTick exists")
)

func (pg *PgDb) VspTableName() string {
	return models.TableNames.VSP
}

func (pg *PgDb) VspTickTableName() string {
	return models.TableNames.VSPTick
}

// StoreVSPs attempts to store the vsp responses by calling storeVspResponseG and returning
// a slice of errors
func (pg *PgDb) StoreVSPs(ctx context.Context, data vsp.Response) (int, []error) {
	if ctx.Err() != nil {
		return 0, []error{ctx.Err()}
	}
	errs := make([]error, 0, len(data))
	completed := 0
	for name, tick := range data {
		err := pg.storeVspResponse(ctx, name, tick)
		if err == nil {
			completed++
		} else if err != vspTickExistsErr {
			log.Trace(err)
			errs = append(errs, err)
		}
		if ctx.Err() != nil {
			return 0, append(errs, ctx.Err())
		}
	}
	if completed == 0 {
		log.Info("Unable to store any vsp entry")
	}
	return completed, errs
}

func (pg *PgDb) storeVspResponse(ctx context.Context, name string, resp *vsp.ResposeData) error {
	txr, err := pg.db.Begin()
	if err != nil {
		return err
	}

	pool, err := models.VSPS(models.VSPWhere.Name.EQ(null.StringFrom(name))).One(ctx, pg.db)
	if err == sql.ErrNoRows {
		pool = responseToVSP(name, resp)
		err := pg.tryInsert(ctx, txr, pool)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	vspTick := responseToVSPTick(pool.ID, resp)

	err = vspTick.Insert(ctx, pg.db, boil.Infer())
	if err != nil {
		errR := txr.Rollback()
		if errR != nil {
			return err
		}
		if strings.Contains(err.Error(), "unique constraint") {
			return vspTickExistsErr
		}
		return err
	}

	err = txr.Commit()
	if err != nil {
		return txr.Rollback()
	}
	return nil
}

func responseToVSP(name string, resp *vsp.ResposeData) *models.VSP {
	return &models.VSP{
		Name:                 null.StringFrom(name),
		APIEnabled:           null.BoolFrom(resp.APIEnabled),
		APIVersionsSupported: types.Int64Array(resp.APIVersionsSupported),
		Network:              null.StringFrom(resp.Network),
		URL:                  null.StringFrom(resp.URL),
		Launched:             null.TimeFrom(time.Unix(resp.Launched, 0)),
	}
}

func responseToVSPTick(poolID int, resp *vsp.ResposeData) *models.VSPTick {
	return &models.VSPTick{
		VSPID:            poolID,
		Immature:         resp.Immature,
		Live:             resp.Live,
		Voted:            resp.Voted,
		Missed:           resp.Missed,
		PoolFees:         resp.PoolFees,
		ProportionLive:   resp.ProportionLive,
		ProportionMissed: resp.ProportionMissed,
		UserCount:        resp.UserCount,
		UsersActive:      resp.UserCountActive,
		Time:             time.Unix(resp.LastUpdated, 0).UTC(),
	}
}

func (pg *PgDb) FetchVSPs(ctx context.Context) ([]vsp.VSPDto, error) {
	vspData, err := models.VSPS(qm.OrderBy(models.VSPColumns.URL), qm.OrderBy(models.VSPColumns.Name)).All(ctx, pg.db)
	if err != nil {
		return nil, err
	}

	var result []vsp.VSPDto
	for _, item := range vspData {
		parsedURL, err := url.Parse(item.URL.String)
		if err != nil {
			return nil, err
		}
		result = append(result, vsp.VSPDto{
			Name:                 item.Name.String,
			APIEnabled:           item.APIEnabled.Bool,
			APIVersionsSupported: item.APIVersionsSupported,
			Network:              item.Network.String,
			URL:                  item.URL.String,
			Host:                 parsedURL.Host,
			Launched:             item.Launched.Time,
		})
	}

	return result, nil
}

func (pg *PgDb) AddVspSourceFromSync(ctx context.Context, vspData interface{}) error {
	vspDto := vspData.(vsp.VSPDto)
	count, _ := models.VSPS(models.VSPWhere.Name.EQ(null.StringFrom(vspDto.Name))).Count(ctx, pg.db)
	if count > 0 {
		return nil
	}
	vspModel := models.VSP{
		ID:                   vspDto.ID,
		Name:                 null.StringFrom(vspDto.Name),
		APIEnabled:           null.BoolFrom(vspDto.APIEnabled),
		APIVersionsSupported: vspDto.APIVersionsSupported,
		Network:              null.StringFrom(vspDto.Network),
		URL:                  null.StringFrom(vspDto.URL),
		Launched:             null.TimeFrom(vspDto.Launched),
	}
	err := vspModel.Insert(ctx, pg.db, boil.Infer())
	return err
}

func (pg *PgDb) FetchVspSourcesForSync(ctx context.Context, lastID int64, skip, take int) ([]vsp.VSPDto, int64, error) {
	vspData, err := models.VSPS(
		models.VSPWhere.ID.GT(int(lastID)),
		qm.Offset(skip), qm.Limit(take)).All(ctx, pg.db)
	if err != nil {
		return nil, 0, err
	}

	var result []vsp.VSPDto
	for _, item := range vspData {
		parsedURL, err := url.Parse(item.URL.String)
		if err != nil {
			return nil, 0, err
		}
		result = append(result, vsp.VSPDto{
			ID:                   item.ID,
			Name:                 item.Name.String,
			APIEnabled:           item.APIEnabled.Bool,
			APIVersionsSupported: item.APIVersionsSupported,
			Network:              item.Network.String,
			URL:                  item.URL.String,
			Host:                 parsedURL.Host,
			Launched:             item.Launched.Time,
		})
	}

	totalCount, err := models.VSPS(models.VSPWhere.ID.GT(int(lastID))).Count(ctx, pg.db)

	return result, totalCount, err
}

// VSPTicks
func (pg *PgDb) FetchVspTicksForSync(ctx context.Context, lastID int64, skip, take int) ([]datasync.VSPTickSyncDto, int64, error) {
	vspIdQuery := models.VSPTickWhere.ID.GT(int(lastID))

	vspTickSlice, err := models.VSPTicks(
		vspIdQuery,
		qm.OrderBy(models.VSPTickColumns.ID),
		qm.Limit(take), qm.Offset(skip)).All(ctx, pg.db)
	if err != nil {
		return nil, 0, err
	}

	vspTickCount, err := models.VSPTicks(vspIdQuery).Count(ctx, pg.db)

	var vspTicks []datasync.VSPTickSyncDto
	for _, tick := range vspTickSlice {
		vspTicks = append(vspTicks, pg.vspTickModelToSyncDto(tick))
	}

	return vspTicks, vspTickCount, nil
}

func (pg *PgDb) AddVspTicksFromSync(ctx context.Context, tick datasync.VSPTickSyncDto) error {
	if _, err := models.VSPTicks(models.VSPTickWhere.VSPID.EQ(tick.VSPID),
		models.VSPTickWhere.Time.EQ(tick.Time)).One(ctx, pg.db); err == nil {
		return nil // record exists
	}
	tickModel := models.VSPTick{
		ID:               tick.ID,
		VSPID:            tick.VSPID,
		Immature:         tick.Immature,
		Live:             tick.Live,
		Voted:            tick.Voted,
		Missed:           tick.Missed,
		PoolFees:         tick.PoolFees,
		ProportionLive:   tick.ProportionLive,
		ProportionMissed: tick.ProportionMissed,
		UserCount:        tick.UserCount,
		UsersActive:      tick.UsersActive,
		Time:             tick.Time,
	}

	return tickModel.Insert(ctx, pg.db, boil.Infer())
}

func (pg *PgDb) FiltredVSPTicks(ctx context.Context, vspName string, offset, limit int) ([]vsp.VSPTickDto, int64, error) {
	vspInfo, err := models.VSPS(models.VSPWhere.Name.EQ(null.StringFrom(vspName))).One(ctx, pg.db)
	if err != nil {
		return nil, 0, err
	}

	vspIdQuery := models.VSPTickWhere.VSPID.EQ(vspInfo.ID)
	vspTickSlice, err := models.VSPTicks(qm.Load(models.VSPTickRels.VSP), vspIdQuery, qm.Limit(limit), qm.Offset(offset), qm.OrderBy(fmt.Sprintf("%s DESC", models.VSPTickColumns.Time))).All(ctx, pg.db)
	if err != nil {
		return nil, 0, err
	}

	vspTickCount, err := models.VSPTicks(qm.Load(models.VSPTickRels.VSP), vspIdQuery).Count(ctx, pg.db)

	vspTicks := []vsp.VSPTickDto{}
	for _, tick := range vspTickSlice {
		vspTicks = append(vspTicks, pg.vspTickModelToDto(tick))
	}

	return vspTicks, vspTickCount, nil
}

// VSPTicks
// todo impliment sorting for VSP ticks as it is currently been sorted by time
func (pg *PgDb) AllVSPTicks(ctx context.Context, offset, limit int) ([]vsp.VSPTickDto, int64, error) {
	vspTickSlice, err := models.VSPTicks(qm.Load(models.VSPTickRels.VSP), qm.Limit(limit), qm.Offset(offset), qm.OrderBy(fmt.Sprintf("%s DESC", models.VSPTickColumns.Time))).All(ctx, pg.db)
	if err != nil {
		return nil, 0, err
	}

	vspTickCount, err := models.VSPTicks().Count(ctx, pg.db)
	if err != nil {
		return nil, 0, err
	}

	vspTicks := []vsp.VSPTickDto{}
	for _, tick := range vspTickSlice {
		vspTicks = append(vspTicks, pg.vspTickModelToDto(tick))
	}

	return vspTicks, vspTickCount, nil
}

func (pg *PgDb) vspTickModelToDto(tick *models.VSPTick) vsp.VSPTickDto {
	return vsp.VSPTickDto{
		ID:               tick.ID,
		VSP:              tick.R.VSP.Name.String,
		Time:             tick.Time.Format(dateTemplate),
		Immature:         tick.Immature,
		Live:             tick.Live,
		Missed:           tick.Missed,
		PoolFees:         tick.PoolFees,
		ProportionLive:   RoundValue(tick.ProportionLive),
		ProportionMissed: RoundValue(tick.ProportionMissed),
		UserCount:        tick.UserCount,
		UsersActive:      tick.UsersActive,
		Voted:            tick.Voted,
	}
}

func (pg *PgDb) vspTickModelToSyncDto(tick *models.VSPTick) datasync.VSPTickSyncDto {
	return datasync.VSPTickSyncDto{
		ID:               tick.ID,
		VSPID:            tick.VSPID,
		Time:             tick.Time,
		Immature:         tick.Immature,
		Live:             tick.Live,
		Missed:           tick.Missed,
		PoolFees:         tick.PoolFees,
		ProportionLive:   tick.ProportionLive,
		ProportionMissed: tick.ProportionMissed,
		UserCount:        tick.UserCount,
		UsersActive:      tick.UsersActive,
		Voted:            tick.Voted,
	}
}

func (pg *PgDb) LastVspTickEntryTime() (time time.Time) {
	rows := pg.db.QueryRow(lastVspTickEntryTime)
	_ = rows.Scan(&time)
	return
}

func (pg *PgDb) VspTickCount(ctx context.Context) (int64, error) {
	return models.VSPTicks().Count(ctx, pg.db)
}

func (pg *PgDb) FetchChartData(ctx context.Context, attribute, vspName string) (records []vsp.ChartData, err error) {
	attribute = strings.ToLower(attribute)
	vspInfo, err := models.VSPS(models.VSPWhere.Name.EQ(null.StringFrom(vspName))).One(ctx, pg.db)
	if err != nil {
		return nil, err
	}
	query := fmt.Sprintf("SELECT time as date, %s as record FROM vsp_tick where %s = %d ORDER BY time",
		attribute, models.VSPTickColumns.VSPID, vspInfo.ID)
	rows, err := pg.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var rec vsp.ChartData
		err = rows.Scan(&rec.Date, &rec.Record)
		if err != nil {
			return nil, err
		}
		if attribute == models.VSPTickColumns.ProportionLive || attribute == models.VSPTickColumns.ProportionMissed {
			value, err := strconv.ParseFloat(rec.Record, 64)
			if err != nil {
				return nil, err
			}
			rec.Record = RoundValue(value)
		}
		records = append(records, rec)
	}
	return
}

func (pg *PgDb) fetchChartData(ctx context.Context, vspName string, start time.Time) (records models.VSPTickSlice, err error) {
	vspInfo, err := models.VSPS(models.VSPWhere.Name.EQ(null.StringFrom(vspName))).One(ctx, pg.db)
	if err != nil {
		return nil, err
	}
	return models.VSPTicks(models.VSPTickWhere.VSPID.EQ(vspInfo.ID), models.VSPTickWhere.Time.GT(start)).All(ctx, pg.db)
}

func (pg *PgDb) GetVspTickDistinctDates(ctx context.Context, vsps []string) ([]time.Time, error) {
	var vspIds []string
	for _, vspName := range vsps {
		id, err := pg.vspIdByName(ctx, vspName)
		if err != nil {
			return nil, err
		}
		vspIds = append(vspIds, strconv.Itoa(id))
	}

	query := fmt.Sprintf("SELECT DISTINCT time FROM vsp_tick WHERE vsp_id IN ('%s') ORDER BY time", strings.Join(vspIds, "', '"))
	rows, err := pg.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	var dates []time.Time

	for rows.Next() {
		var date time.Time
		err = rows.Scan(&date)
		if err != nil {
			return nil, err
		}
		dates = append(dates, date)
	}
	return dates, nil
}

func (pg *PgDb) allVspTickDates(ctx context.Context, start time.Time) ([]time.Time, error) {
	vspDates, err := models.VSPTicks(
		qm.Select(models.VSPTickColumns.Time),
		models.VSPTickWhere.Time.GT(start),
		qm.OrderBy(models.VSPTickColumns.Time),
		).All(ctx, pg.db)

	if err != nil {
		return nil, err
	}

	var dates []time.Time
	var unique = map[time.Time]bool{}

	for _, data := range vspDates {
		if _, found := unique[data.Time]; !found {
			dates = append(dates, data.Time)
			unique[data.Time] = true
		}
	}
	return dates, nil
}

func (pg *PgDb) vspIdByName(ctx context.Context, name string) (id int, err error) {
	vspModel, err := models.VSPS(models.VSPWhere.Name.EQ(null.StringFrom(name))).One(ctx, pg.db)
	if err != nil {
		return 0, err
	}
	return vspModel.ID, nil
}

type vspSet struct {
	time     []uint64
	immature         map[string][]*null.Uint64
	live             map[string][]*null.Uint64
	voted            map[string][]*null.Uint64
	missed           map[string][]*null.Uint64
	poolFees         map[string][]*null.Float64
	proportionLive   map[string][]*null.Float64
	proportionMissed map[string][]*null.Float64
	userCount        map[string][]*null.Uint64
	usersActive      map[string][]*null.Uint64
}

func (pg *PgDb) fetchVspChart(ctx context.Context, charts *cache.ChartData) (interface{}, func(), error) {
	cancelFun := func() {}
	var vspDataSet = vspSet{
		time:             []uint64{},
		immature:         make(map[string][]*null.Uint64),
		live:             make(map[string][]*null.Uint64),
		voted:            make(map[string][]*null.Uint64),
		missed:           make(map[string][]*null.Uint64),
		poolFees:         make(map[string][]*null.Float64),
		proportionLive:   make(map[string][]*null.Float64),
		proportionMissed: make(map[string][]*null.Float64),
		userCount:        make(map[string][]*null.Uint64),
		usersActive:      make(map[string][]*null.Uint64),
	}

	allVspData, err := pg.FetchVSPs(ctx)
	if err != nil {
		return nil, cancelFun, err
	}

	var vsps []string
	for _, vspSource := range allVspData {
		vsps = append(vsps, vspSource.Name)
	}

	dates, err := pg.allVspTickDates(ctx, time.Unix(int64(charts.VspTime()), 0))
	if err != nil && err != sql.ErrNoRows{
		return nil, cancelFun, err
	}

	for _, date := range dates {
		vspDataSet.time = append(vspDataSet.time, uint64(date.Unix()))
	}

	for _, vspSource := range allVspData {
		points, err := pg.fetchChartData(ctx, vspSource.Name, int64ToTime(int64(charts.VspTime())))
		if err != nil {
			return nil, cancelFun, fmt.Errorf("error in fetching records for %s: %s", vspSource.Name, err.Error())
		}

		var pointsMap= map[time.Time]*models.VSPTick{}
		for _, record := range points {
			pointsMap[record.Time] = record
		}

		var hasFoundOne bool
		for _, date := range dates {
			if record, found := pointsMap[date]; found {
				vspDataSet.immature[vspSource.Name] = append(vspDataSet.immature[vspSource.Name], &null.Uint64{Valid: true, Uint64: uint64(record.Immature)})
				vspDataSet.live[vspSource.Name] = append(vspDataSet.live[vspSource.Name], &null.Uint64{Valid: true, Uint64: uint64(record.Live)})
				vspDataSet.voted[vspSource.Name] = append(vspDataSet.voted[vspSource.Name], &null.Uint64{Valid: true, Uint64: uint64(record.Voted)})
				vspDataSet.missed[vspSource.Name] = append(vspDataSet.missed[vspSource.Name], &null.Uint64{Valid: true, Uint64: uint64(record.Missed)})
				vspDataSet.poolFees[vspSource.Name] = append(vspDataSet.poolFees[vspSource.Name], &null.Float64{Valid: true, Float64: record.PoolFees})
				vspDataSet.proportionLive[vspSource.Name] = append(vspDataSet.proportionLive[vspSource.Name], &null.Float64{Valid: true, Float64: record.ProportionLive})
				vspDataSet.proportionMissed[vspSource.Name] = append(vspDataSet.proportionMissed[vspSource.Name], &null.Float64{Valid: true, Float64: record.ProportionMissed})
				vspDataSet.userCount[vspSource.Name] = append(vspDataSet.userCount[vspSource.Name], &null.Uint64{Valid: true, Uint64: uint64(record.UserCount)})
				vspDataSet.usersActive[vspSource.Name] = append(vspDataSet.usersActive[vspSource.Name], &null.Uint64{Valid: true, Uint64: uint64(record.UsersActive)})
				hasFoundOne = true
			} else {
				if hasFoundOne {
					vspDataSet.immature[vspSource.Name] = append(vspDataSet.immature[vspSource.Name], &null.Uint64{Valid: false})
					vspDataSet.live[vspSource.Name] = append(vspDataSet.live[vspSource.Name], &null.Uint64{Valid: false})
					vspDataSet.voted[vspSource.Name] = append(vspDataSet.voted[vspSource.Name], &null.Uint64{Valid: false})
					vspDataSet.missed[vspSource.Name] = append(vspDataSet.missed[vspSource.Name], &null.Uint64{Valid: false})
					vspDataSet.poolFees[vspSource.Name] = append(vspDataSet.poolFees[vspSource.Name], &null.Float64{Valid: false})
					vspDataSet.proportionLive[vspSource.Name] = append(vspDataSet.proportionLive[vspSource.Name], &null.Float64{Valid: false})
					vspDataSet.proportionMissed[vspSource.Name] = append(vspDataSet.proportionMissed[vspSource.Name], &null.Float64{Valid: false})
					vspDataSet.userCount[vspSource.Name] = append(vspDataSet.userCount[vspSource.Name], &null.Uint64{Valid: false})
					vspDataSet.usersActive[vspSource.Name] = append(vspDataSet.usersActive[vspSource.Name], &null.Uint64{Valid: false})
				} else {
					vspDataSet.immature[vspSource.Name] = append(vspDataSet.immature[vspSource.Name], nil)
					vspDataSet.live[vspSource.Name] = append(vspDataSet.live[vspSource.Name], nil)
					vspDataSet.voted[vspSource.Name] = append(vspDataSet.voted[vspSource.Name], nil)
					vspDataSet.missed[vspSource.Name] = append(vspDataSet.missed[vspSource.Name], nil)
					vspDataSet.poolFees[vspSource.Name] = append(vspDataSet.poolFees[vspSource.Name], nil)
					vspDataSet.proportionLive[vspSource.Name] = append(vspDataSet.proportionLive[vspSource.Name], nil)
					vspDataSet.proportionMissed[vspSource.Name] = append(vspDataSet.proportionMissed[vspSource.Name], nil)
					vspDataSet.userCount[vspSource.Name] = append(vspDataSet.userCount[vspSource.Name], nil)
					vspDataSet.usersActive[vspSource.Name] = append(vspDataSet.usersActive[vspSource.Name], nil)
				}
			}
		}
	}

	return vspDataSet, cancelFun, nil
}

func appendVspChart(charts *cache.ChartData, data interface{}) error {
	vspDataSet := data.(vspSet)

	charts.Vsp.Time = append(charts.Vsp.Time, vspDataSet.time...)

	for vspSource, record := range vspDataSet.immature {
		if charts.Vsp.Immature == nil {
			charts.Vsp.Immature = map[string]cache.ChartNullUints{}
		}

		charts.Vsp.Immature[vspSource] = append(charts.Vsp.Immature[vspSource], record...)
	}

	for vspSource, record := range vspDataSet.live {
		if charts.Vsp.Live == nil {
			charts.Vsp.Live = map[string]cache.ChartNullUints{}
		}

		charts.Vsp.Live[vspSource] = append(charts.Vsp.Live[vspSource], record...)
	}

	for vspSource, record := range vspDataSet.voted {
		if charts.Vsp.Voted == nil {
			charts.Vsp.Voted = map[string]cache.ChartNullUints{}
		}

		charts.Vsp.Voted[vspSource] = append(charts.Vsp.Voted[vspSource], record...)
	}

	for vspSource, record := range vspDataSet.missed {
		if charts.Vsp.Missed == nil {
			charts.Vsp.Missed = map[string]cache.ChartNullUints{}
		}

		charts.Vsp.Missed[vspSource] = append(charts.Vsp.Missed[vspSource], record...)
	}

	for vspSource, record := range vspDataSet.poolFees {
		if charts.Vsp.PoolFees == nil {
			charts.Vsp.PoolFees = map[string]cache.ChartNullFloats{}
		}

		charts.Vsp.PoolFees[vspSource] = append(charts.Vsp.PoolFees[vspSource], record...)
	}

	return nil
}
