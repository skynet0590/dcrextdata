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
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"github.com/volatiletech/sqlboiler/queries/qmhelper"
	"github.com/volatiletech/sqlboiler/strmangle"
)

// Propagation is an object representing the database table.
type Propagation struct {
	Height    int64   `boil:"height" json:"height" toml:"height" yaml:"height"`
	Time      int64   `boil:"time" json:"time" toml:"time" yaml:"time"`
	Bin       string  `boil:"bin" json:"bin" toml:"bin" yaml:"bin"`
	Source    string  `boil:"source" json:"source" toml:"source" yaml:"source"`
	Deviation float64 `boil:"deviation" json:"deviation" toml:"deviation" yaml:"deviation"`

	R *propagationR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L propagationL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var PropagationColumns = struct {
	Height    string
	Time      string
	Bin       string
	Source    string
	Deviation string
}{
	Height:    "height",
	Time:      "time",
	Bin:       "bin",
	Source:    "source",
	Deviation: "deviation",
}

// Generated where

var PropagationWhere = struct {
	Height    whereHelperint64
	Time      whereHelperint64
	Bin       whereHelperstring
	Source    whereHelperstring
	Deviation whereHelperfloat64
}{
	Height:    whereHelperint64{field: "\"propagation\".\"height\""},
	Time:      whereHelperint64{field: "\"propagation\".\"time\""},
	Bin:       whereHelperstring{field: "\"propagation\".\"bin\""},
	Source:    whereHelperstring{field: "\"propagation\".\"source\""},
	Deviation: whereHelperfloat64{field: "\"propagation\".\"deviation\""},
}

// PropagationRels is where relationship names are stored.
var PropagationRels = struct {
}{}

// propagationR is where relationships are stored.
type propagationR struct {
}

// NewStruct creates a new relationship struct
func (*propagationR) NewStruct() *propagationR {
	return &propagationR{}
}

// propagationL is where Load methods for each relationship are stored.
type propagationL struct{}

var (
	propagationAllColumns            = []string{"height", "time", "bin", "source", "deviation"}
	propagationColumnsWithoutDefault = []string{"height", "time", "bin", "source", "deviation"}
	propagationColumnsWithDefault    = []string{}
	propagationPrimaryKeyColumns     = []string{"height", "source", "bin"}
)

type (
	// PropagationSlice is an alias for a slice of pointers to Propagation.
	// This should generally be used opposed to []Propagation.
	PropagationSlice []*Propagation

	propagationQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	propagationType                 = reflect.TypeOf(&Propagation{})
	propagationMapping              = queries.MakeStructMapping(propagationType)
	propagationPrimaryKeyMapping, _ = queries.BindMapping(propagationType, propagationMapping, propagationPrimaryKeyColumns)
	propagationInsertCacheMut       sync.RWMutex
	propagationInsertCache          = make(map[string]insertCache)
	propagationUpdateCacheMut       sync.RWMutex
	propagationUpdateCache          = make(map[string]updateCache)
	propagationUpsertCacheMut       sync.RWMutex
	propagationUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

// One returns a single propagation record from the query.
func (q propagationQuery) One(ctx context.Context, exec boil.ContextExecutor) (*Propagation, error) {
	o := &Propagation{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for propagation")
	}

	return o, nil
}

// All returns all Propagation records from the query.
func (q propagationQuery) All(ctx context.Context, exec boil.ContextExecutor) (PropagationSlice, error) {
	var o []*Propagation

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to Propagation slice")
	}

	return o, nil
}

// Count returns the count of all Propagation records in the query.
func (q propagationQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count propagation rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q propagationQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if propagation exists")
	}

	return count > 0, nil
}

// Propagations retrieves all the records using an executor.
func Propagations(mods ...qm.QueryMod) propagationQuery {
	mods = append(mods, qm.From("\"propagation\""))
	return propagationQuery{NewQuery(mods...)}
}

// FindPropagation retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindPropagation(ctx context.Context, exec boil.ContextExecutor, height int64, source string, bin string, selectCols ...string) (*Propagation, error) {
	propagationObj := &Propagation{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"propagation\" where \"height\"=$1 AND \"source\"=$2 AND \"bin\"=$3", sel,
	)

	q := queries.Raw(query, height, source, bin)

	err := q.Bind(ctx, exec, propagationObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from propagation")
	}

	return propagationObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *Propagation) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no propagation provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(propagationColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	propagationInsertCacheMut.RLock()
	cache, cached := propagationInsertCache[key]
	propagationInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			propagationAllColumns,
			propagationColumnsWithDefault,
			propagationColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(propagationType, propagationMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(propagationType, propagationMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"propagation\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"propagation\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "models: unable to insert into propagation")
	}

	if !cached {
		propagationInsertCacheMut.Lock()
		propagationInsertCache[key] = cache
		propagationInsertCacheMut.Unlock()
	}

	return nil
}

// Update uses an executor to update the Propagation.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *Propagation) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	key := makeCacheKey(columns, nil)
	propagationUpdateCacheMut.RLock()
	cache, cached := propagationUpdateCache[key]
	propagationUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			propagationAllColumns,
			propagationPrimaryKeyColumns,
		)

		if len(wl) == 0 {
			return 0, errors.New("models: unable to update propagation, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"propagation\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, propagationPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(propagationType, propagationMapping, append(wl, propagationPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "models: unable to update propagation row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for propagation")
	}

	if !cached {
		propagationUpdateCacheMut.Lock()
		propagationUpdateCache[key] = cache
		propagationUpdateCacheMut.Unlock()
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values.
func (q propagationQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for propagation")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for propagation")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o PropagationSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), propagationPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"propagation\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, propagationPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in propagation slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all propagation")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *Propagation) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no propagation provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(propagationColumnsWithDefault, o)

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

	propagationUpsertCacheMut.RLock()
	cache, cached := propagationUpsertCache[key]
	propagationUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			propagationAllColumns,
			propagationColumnsWithDefault,
			propagationColumnsWithoutDefault,
			nzDefaults,
		)
		update := updateColumns.UpdateColumnSet(
			propagationAllColumns,
			propagationPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert propagation, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(propagationPrimaryKeyColumns))
			copy(conflict, propagationPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"propagation\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(propagationType, propagationMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(propagationType, propagationMapping, ret)
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
		return errors.Wrap(err, "models: unable to upsert propagation")
	}

	if !cached {
		propagationUpsertCacheMut.Lock()
		propagationUpsertCache[key] = cache
		propagationUpsertCacheMut.Unlock()
	}

	return nil
}

// Delete deletes a single Propagation record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Propagation) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no Propagation provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), propagationPrimaryKeyMapping)
	sql := "DELETE FROM \"propagation\" WHERE \"height\"=$1 AND \"source\"=$2 AND \"bin\"=$3"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from propagation")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for propagation")
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q propagationQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no propagationQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from propagation")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for propagation")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o PropagationSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), propagationPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"propagation\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, propagationPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from propagation slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for propagation")
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *Propagation) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindPropagation(ctx, exec, o.Height, o.Source, o.Bin)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *PropagationSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := PropagationSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), propagationPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"propagation\".* FROM \"propagation\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, propagationPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in PropagationSlice")
	}

	*o = slice

	return nil
}

// PropagationExists checks if the Propagation row exists.
func PropagationExists(ctx context.Context, exec boil.ContextExecutor, height int64, source string, bin string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"propagation\" where \"height\"=$1 AND \"source\"=$2 AND \"bin\"=$3 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, height, source, bin)
	}
	row := exec.QueryRowContext(ctx, sql, height, source, bin)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if propagation exists")
	}

	return exists, nil
}
