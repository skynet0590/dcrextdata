// Code generated by SQLBoiler 3.7.1 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
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

	"github.com/friendsofgo/errors"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"github.com/volatiletech/sqlboiler/queries/qmhelper"
	"github.com/volatiletech/sqlboiler/strmangle"
)

// MempoolBin is an object representing the database table.
type MempoolBin struct {
	Time                 int64        `boil:"time" json:"time" toml:"time" yaml:"time"`
	Bin                  string       `boil:"bin" json:"bin" toml:"bin" yaml:"bin"`
	NumberOfTransactions null.Int     `boil:"number_of_transactions" json:"number_of_transactions,omitempty" toml:"number_of_transactions" yaml:"number_of_transactions,omitempty"`
	Size                 null.Int     `boil:"size" json:"size,omitempty" toml:"size" yaml:"size,omitempty"`
	TotalFee             null.Float64 `boil:"total_fee" json:"total_fee,omitempty" toml:"total_fee" yaml:"total_fee,omitempty"`

	R *mempoolBinR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L mempoolBinL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var MempoolBinColumns = struct {
	Time                 string
	Bin                  string
	NumberOfTransactions string
	Size                 string
	TotalFee             string
}{
	Time:                 "time",
	Bin:                  "bin",
	NumberOfTransactions: "number_of_transactions",
	Size:                 "size",
	TotalFee:             "total_fee",
}

// Generated where

var MempoolBinWhere = struct {
	Time                 whereHelperint64
	Bin                  whereHelperstring
	NumberOfTransactions whereHelpernull_Int
	Size                 whereHelpernull_Int
	TotalFee             whereHelpernull_Float64
}{
	Time:                 whereHelperint64{field: "\"mempool_bin\".\"time\""},
	Bin:                  whereHelperstring{field: "\"mempool_bin\".\"bin\""},
	NumberOfTransactions: whereHelpernull_Int{field: "\"mempool_bin\".\"number_of_transactions\""},
	Size:                 whereHelpernull_Int{field: "\"mempool_bin\".\"size\""},
	TotalFee:             whereHelpernull_Float64{field: "\"mempool_bin\".\"total_fee\""},
}

// MempoolBinRels is where relationship names are stored.
var MempoolBinRels = struct {
}{}

// mempoolBinR is where relationships are stored.
type mempoolBinR struct {
}

// NewStruct creates a new relationship struct
func (*mempoolBinR) NewStruct() *mempoolBinR {
	return &mempoolBinR{}
}

// mempoolBinL is where Load methods for each relationship are stored.
type mempoolBinL struct{}

var (
	mempoolBinAllColumns            = []string{"time", "bin", "number_of_transactions", "size", "total_fee"}
	mempoolBinColumnsWithoutDefault = []string{"time", "bin", "number_of_transactions", "size", "total_fee"}
	mempoolBinColumnsWithDefault    = []string{}
	mempoolBinPrimaryKeyColumns     = []string{"time", "bin"}
)

type (
	// MempoolBinSlice is an alias for a slice of pointers to MempoolBin.
	// This should generally be used opposed to []MempoolBin.
	MempoolBinSlice []*MempoolBin

	mempoolBinQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	mempoolBinType                 = reflect.TypeOf(&MempoolBin{})
	mempoolBinMapping              = queries.MakeStructMapping(mempoolBinType)
	mempoolBinPrimaryKeyMapping, _ = queries.BindMapping(mempoolBinType, mempoolBinMapping, mempoolBinPrimaryKeyColumns)
	mempoolBinInsertCacheMut       sync.RWMutex
	mempoolBinInsertCache          = make(map[string]insertCache)
	mempoolBinUpdateCacheMut       sync.RWMutex
	mempoolBinUpdateCache          = make(map[string]updateCache)
	mempoolBinUpsertCacheMut       sync.RWMutex
	mempoolBinUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

// One returns a single mempoolBin record from the query.
func (q mempoolBinQuery) One(ctx context.Context, exec boil.ContextExecutor) (*MempoolBin, error) {
	o := &MempoolBin{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for mempool_bin")
	}

	return o, nil
}

// All returns all MempoolBin records from the query.
func (q mempoolBinQuery) All(ctx context.Context, exec boil.ContextExecutor) (MempoolBinSlice, error) {
	var o []*MempoolBin

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to MempoolBin slice")
	}

	return o, nil
}

// Count returns the count of all MempoolBin records in the query.
func (q mempoolBinQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count mempool_bin rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q mempoolBinQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if mempool_bin exists")
	}

	return count > 0, nil
}

// MempoolBins retrieves all the records using an executor.
func MempoolBins(mods ...qm.QueryMod) mempoolBinQuery {
	mods = append(mods, qm.From("\"mempool_bin\""))
	return mempoolBinQuery{NewQuery(mods...)}
}

// FindMempoolBin retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindMempoolBin(ctx context.Context, exec boil.ContextExecutor, time int64, bin string, selectCols ...string) (*MempoolBin, error) {
	mempoolBinObj := &MempoolBin{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"mempool_bin\" where \"time\"=$1 AND \"bin\"=$2", sel,
	)

	q := queries.Raw(query, time, bin)

	err := q.Bind(ctx, exec, mempoolBinObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from mempool_bin")
	}

	return mempoolBinObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *MempoolBin) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no mempool_bin provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(mempoolBinColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	mempoolBinInsertCacheMut.RLock()
	cache, cached := mempoolBinInsertCache[key]
	mempoolBinInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			mempoolBinAllColumns,
			mempoolBinColumnsWithDefault,
			mempoolBinColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(mempoolBinType, mempoolBinMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(mempoolBinType, mempoolBinMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"mempool_bin\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"mempool_bin\" %sDEFAULT VALUES%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			queryReturning = fmt.Sprintf(" RETURNING \"%s\"", strings.Join(returnColumns, "\",\""))
		}

		cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}

	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}

	if err != nil {
		return errors.Wrap(err, "models: unable to insert into mempool_bin")
	}

	if !cached {
		mempoolBinInsertCacheMut.Lock()
		mempoolBinInsertCache[key] = cache
		mempoolBinInsertCacheMut.Unlock()
	}

	return nil
}

// Update uses an executor to update the MempoolBin.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *MempoolBin) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	key := makeCacheKey(columns, nil)
	mempoolBinUpdateCacheMut.RLock()
	cache, cached := mempoolBinUpdateCache[key]
	mempoolBinUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			mempoolBinAllColumns,
			mempoolBinPrimaryKeyColumns,
		)

		if len(wl) == 0 {
			return 0, errors.New("models: unable to update mempool_bin, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"mempool_bin\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, mempoolBinPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(mempoolBinType, mempoolBinMapping, append(wl, mempoolBinPrimaryKeyColumns...))
		if err != nil {
			return 0, err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, values)
	}
	var result sql.Result
	result, err = exec.ExecContext(ctx, cache.query, values...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update mempool_bin row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for mempool_bin")
	}

	if !cached {
		mempoolBinUpdateCacheMut.Lock()
		mempoolBinUpdateCache[key] = cache
		mempoolBinUpdateCacheMut.Unlock()
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values.
func (q mempoolBinQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for mempool_bin")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for mempool_bin")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o MempoolBinSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), mempoolBinPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"mempool_bin\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, mempoolBinPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in mempoolBin slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all mempoolBin")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *MempoolBin) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no mempool_bin provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(mempoolBinColumnsWithDefault, o)

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

	mempoolBinUpsertCacheMut.RLock()
	cache, cached := mempoolBinUpsertCache[key]
	mempoolBinUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			mempoolBinAllColumns,
			mempoolBinColumnsWithDefault,
			mempoolBinColumnsWithoutDefault,
			nzDefaults,
		)
		update := updateColumns.UpdateColumnSet(
			mempoolBinAllColumns,
			mempoolBinPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert mempool_bin, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(mempoolBinPrimaryKeyColumns))
			copy(conflict, mempoolBinPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"mempool_bin\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(mempoolBinType, mempoolBinMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(mempoolBinType, mempoolBinMapping, ret)
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

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
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
		return errors.Wrap(err, "models: unable to upsert mempool_bin")
	}

	if !cached {
		mempoolBinUpsertCacheMut.Lock()
		mempoolBinUpsertCache[key] = cache
		mempoolBinUpsertCacheMut.Unlock()
	}

	return nil
}

// Delete deletes a single MempoolBin record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *MempoolBin) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no MempoolBin provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), mempoolBinPrimaryKeyMapping)
	sql := "DELETE FROM \"mempool_bin\" WHERE \"time\"=$1 AND \"bin\"=$2"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from mempool_bin")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for mempool_bin")
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q mempoolBinQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no mempoolBinQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from mempool_bin")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for mempool_bin")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o MempoolBinSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), mempoolBinPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"mempool_bin\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, mempoolBinPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from mempoolBin slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for mempool_bin")
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *MempoolBin) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindMempoolBin(ctx, exec, o.Time, o.Bin)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *MempoolBinSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := MempoolBinSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), mempoolBinPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"mempool_bin\".* FROM \"mempool_bin\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, mempoolBinPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in MempoolBinSlice")
	}

	*o = slice

	return nil
}

// MempoolBinExists checks if the MempoolBin row exists.
func MempoolBinExists(ctx context.Context, exec boil.ContextExecutor, time int64, bin string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"mempool_bin\" where \"time\"=$1 AND \"bin\"=$2 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, time, bin)
	}
	row := exec.QueryRowContext(ctx, sql, time, bin)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if mempool_bin exists")
	}

	return exists, nil
}
