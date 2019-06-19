// Code generated by SQLBoiler (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"github.com/volatiletech/sqlboiler/queries/qmhelper"
	"github.com/volatiletech/sqlboiler/strmangle"
)

// ExchangeTick is an object representing the database table.
type ExchangeTick struct {
	ID           int       `boil:"id" json:"id" toml:"id" yaml:"id"`
	ExchangeID   int       `boil:"exchange_id" json:"exchange_id" toml:"exchange_id" yaml:"exchange_id"`
	Interval     int       `boil:"interval" json:"interval" toml:"interval" yaml:"interval"`
	High         float64   `boil:"high" json:"high" toml:"high" yaml:"high"`
	Low          float64   `boil:"low" json:"low" toml:"low" yaml:"low"`
	Open         float64   `boil:"open" json:"open" toml:"open" yaml:"open"`
	Close        float64   `boil:"close" json:"close" toml:"close" yaml:"close"`
	Volume       float64   `boil:"volume" json:"volume" toml:"volume" yaml:"volume"`
	CurrencyPair string    `boil:"currency_pair" json:"currency_pair" toml:"currency_pair" yaml:"currency_pair"`
	Time         time.Time `boil:"time" json:"time" toml:"time" yaml:"time"`

	R *exchangeTickR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L exchangeTickL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var ExchangeTickColumns = struct {
	ID           string
	ExchangeID   string
	Interval     string
	High         string
	Low          string
	Open         string
	Close        string
	Volume       string
	CurrencyPair string
	Time         string
}{
	ID:           "id",
	ExchangeID:   "exchange_id",
	Interval:     "interval",
	High:         "high",
	Low:          "low",
	Open:         "open",
	Close:        "close",
	Volume:       "volume",
	CurrencyPair: "currency_pair",
	Time:         "time",
}

// Generated where

type whereHelperfloat64 struct{ field string }

func (w whereHelperfloat64) EQ(x float64) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.EQ, x) }
func (w whereHelperfloat64) NEQ(x float64) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.NEQ, x)
}
func (w whereHelperfloat64) LT(x float64) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.LT, x) }
func (w whereHelperfloat64) LTE(x float64) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LTE, x)
}
func (w whereHelperfloat64) GT(x float64) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.GT, x) }
func (w whereHelperfloat64) GTE(x float64) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GTE, x)
}

type whereHelpertime_Time struct{ field string }

func (w whereHelpertime_Time) EQ(x time.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.EQ, x)
}
func (w whereHelpertime_Time) NEQ(x time.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.NEQ, x)
}
func (w whereHelpertime_Time) LT(x time.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LT, x)
}
func (w whereHelpertime_Time) LTE(x time.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LTE, x)
}
func (w whereHelpertime_Time) GT(x time.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GT, x)
}
func (w whereHelpertime_Time) GTE(x time.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GTE, x)
}

var ExchangeTickWhere = struct {
	ID           whereHelperint
	ExchangeID   whereHelperint
	Interval     whereHelperint
	High         whereHelperfloat64
	Low          whereHelperfloat64
	Open         whereHelperfloat64
	Close        whereHelperfloat64
	Volume       whereHelperfloat64
	CurrencyPair whereHelperstring
	Time         whereHelpertime_Time
}{
	ID:           whereHelperint{field: "\"exchange_tick\".\"id\""},
	ExchangeID:   whereHelperint{field: "\"exchange_tick\".\"exchange_id\""},
	Interval:     whereHelperint{field: "\"exchange_tick\".\"interval\""},
	High:         whereHelperfloat64{field: "\"exchange_tick\".\"high\""},
	Low:          whereHelperfloat64{field: "\"exchange_tick\".\"low\""},
	Open:         whereHelperfloat64{field: "\"exchange_tick\".\"open\""},
	Close:        whereHelperfloat64{field: "\"exchange_tick\".\"close\""},
	Volume:       whereHelperfloat64{field: "\"exchange_tick\".\"volume\""},
	CurrencyPair: whereHelperstring{field: "\"exchange_tick\".\"currency_pair\""},
	Time:         whereHelpertime_Time{field: "\"exchange_tick\".\"time\""},
}

// ExchangeTickRels is where relationship names are stored.
var ExchangeTickRels = struct {
	Exchange string
}{
	Exchange: "Exchange",
}

// exchangeTickR is where relationships are stored.
type exchangeTickR struct {
	Exchange *Exchange
}

// NewStruct creates a new relationship struct
func (*exchangeTickR) NewStruct() *exchangeTickR {
	return &exchangeTickR{}
}

// exchangeTickL is where Load methods for each relationship are stored.
type exchangeTickL struct{}

var (
	exchangeTickAllColumns            = []string{"id", "exchange_id", "interval", "high", "low", "open", "close", "volume", "currency_pair", "time"}
	exchangeTickColumnsWithoutDefault = []string{"exchange_id", "interval", "high", "low", "open", "close", "volume", "currency_pair", "time"}
	exchangeTickColumnsWithDefault    = []string{"id"}
	exchangeTickPrimaryKeyColumns     = []string{"id"}
)

type (
	// ExchangeTickSlice is an alias for a slice of pointers to ExchangeTick.
	// This should generally be used opposed to []ExchangeTick.
	ExchangeTickSlice []*ExchangeTick

	exchangeTickQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	exchangeTickType                 = reflect.TypeOf(&ExchangeTick{})
	exchangeTickMapping              = queries.MakeStructMapping(exchangeTickType)
	exchangeTickPrimaryKeyMapping, _ = queries.BindMapping(exchangeTickType, exchangeTickMapping, exchangeTickPrimaryKeyColumns)
	exchangeTickInsertCacheMut       sync.RWMutex
	exchangeTickInsertCache          = make(map[string]insertCache)
	exchangeTickUpdateCacheMut       sync.RWMutex
	exchangeTickUpdateCache          = make(map[string]updateCache)
	exchangeTickUpsertCacheMut       sync.RWMutex
	exchangeTickUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

// One returns a single exchangeTick record from the query.
func (q exchangeTickQuery) One(ctx context.Context, exec boil.ContextExecutor) (*ExchangeTick, error) {
	o := &ExchangeTick{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for exchange_tick")
	}

	return o, nil
}

// All returns all ExchangeTick records from the query.
func (q exchangeTickQuery) All(ctx context.Context, exec boil.ContextExecutor) (ExchangeTickSlice, error) {
	var o []*ExchangeTick

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to ExchangeTick slice")
	}

	return o, nil
}

// Count returns the count of all ExchangeTick records in the query.
func (q exchangeTickQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count exchange_tick rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q exchangeTickQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if exchange_tick exists")
	}

	return count > 0, nil
}

// Exchange pointed to by the foreign key.
func (o *ExchangeTick) Exchange(mods ...qm.QueryMod) exchangeQuery {
	queryMods := []qm.QueryMod{
		qm.Where("id=?", o.ExchangeID),
	}

	queryMods = append(queryMods, mods...)

	query := Exchanges(queryMods...)
	queries.SetFrom(query.Query, "\"exchange\"")

	return query
}

// LoadExchange allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (exchangeTickL) LoadExchange(ctx context.Context, e boil.ContextExecutor, singular bool, maybeExchangeTick interface{}, mods queries.Applicator) error {
	var slice []*ExchangeTick
	var object *ExchangeTick

	if singular {
		object = maybeExchangeTick.(*ExchangeTick)
	} else {
		slice = *maybeExchangeTick.(*[]*ExchangeTick)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &exchangeTickR{}
		}
		args = append(args, object.ExchangeID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &exchangeTickR{}
			}

			for _, a := range args {
				if a == obj.ExchangeID {
					continue Outer
				}
			}

			args = append(args, obj.ExchangeID)

		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(qm.From(`exchange`), qm.WhereIn(`id in ?`, args...))
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load Exchange")
	}

	var resultSlice []*Exchange
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice Exchange")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for exchange")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for exchange")
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		foreign := resultSlice[0]
		object.R.Exchange = foreign
		if foreign.R == nil {
			foreign.R = &exchangeR{}
		}
		foreign.R.ExchangeTicks = append(foreign.R.ExchangeTicks, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.ExchangeID == foreign.ID {
				local.R.Exchange = foreign
				if foreign.R == nil {
					foreign.R = &exchangeR{}
				}
				foreign.R.ExchangeTicks = append(foreign.R.ExchangeTicks, local)
				break
			}
		}
	}

	return nil
}

// SetExchange of the exchangeTick to the related item.
// Sets o.R.Exchange to related.
// Adds o to related.R.ExchangeTicks.
func (o *ExchangeTick) SetExchange(ctx context.Context, exec boil.ContextExecutor, insert bool, related *Exchange) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"exchange_tick\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"exchange_id"}),
		strmangle.WhereClause("\"", "\"", 2, exchangeTickPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.ID}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, updateQuery)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.ExchangeID = related.ID
	if o.R == nil {
		o.R = &exchangeTickR{
			Exchange: related,
		}
	} else {
		o.R.Exchange = related
	}

	if related.R == nil {
		related.R = &exchangeR{
			ExchangeTicks: ExchangeTickSlice{o},
		}
	} else {
		related.R.ExchangeTicks = append(related.R.ExchangeTicks, o)
	}

	return nil
}

// ExchangeTicks retrieves all the records using an executor.
func ExchangeTicks(mods ...qm.QueryMod) exchangeTickQuery {
	mods = append(mods, qm.From("\"exchange_tick\""))
	return exchangeTickQuery{NewQuery(mods...)}
}

// FindExchangeTick retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindExchangeTick(ctx context.Context, exec boil.ContextExecutor, iD int, selectCols ...string) (*ExchangeTick, error) {
	exchangeTickObj := &ExchangeTick{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"exchange_tick\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, exchangeTickObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from exchange_tick")
	}

	return exchangeTickObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *ExchangeTick) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no exchange_tick provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(exchangeTickColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	exchangeTickInsertCacheMut.RLock()
	cache, cached := exchangeTickInsertCache[key]
	exchangeTickInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			exchangeTickAllColumns,
			exchangeTickColumnsWithDefault,
			exchangeTickColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(exchangeTickType, exchangeTickMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(exchangeTickType, exchangeTickMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"exchange_tick\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"exchange_tick\" %sDEFAULT VALUES%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			queryReturning = fmt.Sprintf(" RETURNING \"%s\"", strings.Join(returnColumns, "\",\""))
		}

		cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, vals)
	}

	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}

	if err != nil {
		return errors.Wrap(err, "models: unable to insert into exchange_tick")
	}

	if !cached {
		exchangeTickInsertCacheMut.Lock()
		exchangeTickInsertCache[key] = cache
		exchangeTickInsertCacheMut.Unlock()
	}

	return nil
}

// Update uses an executor to update the ExchangeTick.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *ExchangeTick) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	key := makeCacheKey(columns, nil)
	exchangeTickUpdateCacheMut.RLock()
	cache, cached := exchangeTickUpdateCache[key]
	exchangeTickUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			exchangeTickAllColumns,
			exchangeTickPrimaryKeyColumns,
		)

		if len(wl) == 0 {
			return 0, errors.New("models: unable to update exchange_tick, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"exchange_tick\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, exchangeTickPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(exchangeTickType, exchangeTickMapping, append(wl, exchangeTickPrimaryKeyColumns...))
		if err != nil {
			return 0, err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	var result sql.Result
	result, err = exec.ExecContext(ctx, cache.query, values...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update exchange_tick row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for exchange_tick")
	}

	if !cached {
		exchangeTickUpdateCacheMut.Lock()
		exchangeTickUpdateCache[key] = cache
		exchangeTickUpdateCacheMut.Unlock()
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values.
func (q exchangeTickQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for exchange_tick")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for exchange_tick")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o ExchangeTickSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	ln := int64(len(o))
	if ln == 0 {
		return 0, nil
	}

	if len(cols) == 0 {
		return 0, errors.New("models: update all requires at least one column argument")
	}

	colNames := make([]string, len(cols))
	args := make([]interface{}, len(cols))

	i := 0
	for name, value := range cols {
		colNames[i] = name
		args[i] = value
		i++
	}

	// Append all of the primary key values for each column
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), exchangeTickPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"exchange_tick\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, exchangeTickPrimaryKeyColumns, len(o)))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in exchangeTick slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all exchangeTick")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *ExchangeTick) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no exchange_tick provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(exchangeTickColumnsWithDefault, o)

	// Build cache key in-line uglily - mysql vs psql problems
	buf := strmangle.GetBuffer()
	if updateOnConflict {
		buf.WriteByte('t')
	} else {
		buf.WriteByte('f')
	}
	buf.WriteByte('.')
	for _, c := range conflictColumns {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(updateColumns.Kind))
	for _, c := range updateColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(insertColumns.Kind))
	for _, c := range insertColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	exchangeTickUpsertCacheMut.RLock()
	cache, cached := exchangeTickUpsertCache[key]
	exchangeTickUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			exchangeTickAllColumns,
			exchangeTickColumnsWithDefault,
			exchangeTickColumnsWithoutDefault,
			nzDefaults,
		)
		update := updateColumns.UpdateColumnSet(
			exchangeTickAllColumns,
			exchangeTickPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert exchange_tick, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(exchangeTickPrimaryKeyColumns))
			copy(conflict, exchangeTickPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"exchange_tick\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(exchangeTickType, exchangeTickMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(exchangeTickType, exchangeTickMapping, ret)
			if err != nil {
				return err
			}
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)
	var returns []interface{}
	if len(cache.retMapping) != 0 {
		returns = queries.PtrsFromMapping(value, cache.retMapping)
	}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, vals)
	}

	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(returns...)
		if err == sql.ErrNoRows {
			err = nil // Postgres doesn't return anything when there's no update
		}
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}
	if err != nil {
		return errors.Wrap(err, "models: unable to upsert exchange_tick")
	}

	if !cached {
		exchangeTickUpsertCacheMut.Lock()
		exchangeTickUpsertCache[key] = cache
		exchangeTickUpsertCacheMut.Unlock()
	}

	return nil
}

// Delete deletes a single ExchangeTick record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *ExchangeTick) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no ExchangeTick provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), exchangeTickPrimaryKeyMapping)
	sql := "DELETE FROM \"exchange_tick\" WHERE \"id\"=$1"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from exchange_tick")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for exchange_tick")
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q exchangeTickQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no exchangeTickQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from exchange_tick")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for exchange_tick")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o ExchangeTickSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), exchangeTickPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"exchange_tick\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, exchangeTickPrimaryKeyColumns, len(o))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}

	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from exchangeTick slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for exchange_tick")
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *ExchangeTick) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindExchangeTick(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *ExchangeTickSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := ExchangeTickSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), exchangeTickPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"exchange_tick\".* FROM \"exchange_tick\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, exchangeTickPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in ExchangeTickSlice")
	}

	*o = slice

	return nil
}

// ExchangeTickExists checks if the ExchangeTick row exists.
func ExchangeTickExists(ctx context.Context, exec boil.ContextExecutor, iD int) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"exchange_tick\" where \"id\"=$1 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, iD)
	}

	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if exchange_tick exists")
	}

	return exists, nil
}
