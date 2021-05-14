// Code generated by SQLBoiler 3.7.1 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"bytes"
	"context"
	"reflect"
	"testing"

	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries"
	"github.com/volatiletech/sqlboiler/randomize"
	"github.com/volatiletech/sqlboiler/strmangle"
)

var (
	// Relationships sometimes use the reflection helper queries.Equal/queries.Assign
	// so force a package dependency in case they don't.
	_ = queries.Equal
)

func testMempoolBins(t *testing.T) {
	t.Parallel()

	query := MempoolBins()

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}

func testMempoolBinsDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &MempoolBin{}
	if err = randomize.Struct(seed, o, mempoolBinDBTypes, true, mempoolBinColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize MempoolBin struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := o.Delete(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := MempoolBins().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testMempoolBinsQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &MempoolBin{}
	if err = randomize.Struct(seed, o, mempoolBinDBTypes, true, mempoolBinColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize MempoolBin struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := MempoolBins().DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := MempoolBins().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testMempoolBinsSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &MempoolBin{}
	if err = randomize.Struct(seed, o, mempoolBinDBTypes, true, mempoolBinColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize MempoolBin struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := MempoolBinSlice{o}

	if rowsAff, err := slice.DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := MempoolBins().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testMempoolBinsExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &MempoolBin{}
	if err = randomize.Struct(seed, o, mempoolBinDBTypes, true, mempoolBinColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize MempoolBin struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	e, err := MempoolBinExists(ctx, tx, o.Time, o.Bin)
	if err != nil {
		t.Errorf("Unable to check if MempoolBin exists: %s", err)
	}
	if !e {
		t.Errorf("Expected MempoolBinExists to return true, but got false.")
	}
}

func testMempoolBinsFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &MempoolBin{}
	if err = randomize.Struct(seed, o, mempoolBinDBTypes, true, mempoolBinColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize MempoolBin struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	mempoolBinFound, err := FindMempoolBin(ctx, tx, o.Time, o.Bin)
	if err != nil {
		t.Error(err)
	}

	if mempoolBinFound == nil {
		t.Error("want a record, got nil")
	}
}

func testMempoolBinsBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &MempoolBin{}
	if err = randomize.Struct(seed, o, mempoolBinDBTypes, true, mempoolBinColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize MempoolBin struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = MempoolBins().Bind(ctx, tx, o); err != nil {
		t.Error(err)
	}
}

func testMempoolBinsOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &MempoolBin{}
	if err = randomize.Struct(seed, o, mempoolBinDBTypes, true, mempoolBinColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize MempoolBin struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if x, err := MempoolBins().One(ctx, tx); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testMempoolBinsAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	mempoolBinOne := &MempoolBin{}
	mempoolBinTwo := &MempoolBin{}
	if err = randomize.Struct(seed, mempoolBinOne, mempoolBinDBTypes, false, mempoolBinColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize MempoolBin struct: %s", err)
	}
	if err = randomize.Struct(seed, mempoolBinTwo, mempoolBinDBTypes, false, mempoolBinColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize MempoolBin struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = mempoolBinOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = mempoolBinTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := MempoolBins().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testMempoolBinsCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	mempoolBinOne := &MempoolBin{}
	mempoolBinTwo := &MempoolBin{}
	if err = randomize.Struct(seed, mempoolBinOne, mempoolBinDBTypes, false, mempoolBinColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize MempoolBin struct: %s", err)
	}
	if err = randomize.Struct(seed, mempoolBinTwo, mempoolBinDBTypes, false, mempoolBinColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize MempoolBin struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = mempoolBinOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = mempoolBinTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := MempoolBins().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func testMempoolBinsInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &MempoolBin{}
	if err = randomize.Struct(seed, o, mempoolBinDBTypes, true, mempoolBinColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize MempoolBin struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := MempoolBins().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testMempoolBinsInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &MempoolBin{}
	if err = randomize.Struct(seed, o, mempoolBinDBTypes, true); err != nil {
		t.Errorf("Unable to randomize MempoolBin struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Whitelist(mempoolBinColumnsWithoutDefault...)); err != nil {
		t.Error(err)
	}

	count, err := MempoolBins().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testMempoolBinsReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &MempoolBin{}
	if err = randomize.Struct(seed, o, mempoolBinDBTypes, true, mempoolBinColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize MempoolBin struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = o.Reload(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testMempoolBinsReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &MempoolBin{}
	if err = randomize.Struct(seed, o, mempoolBinDBTypes, true, mempoolBinColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize MempoolBin struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := MempoolBinSlice{o}

	if err = slice.ReloadAll(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testMempoolBinsSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &MempoolBin{}
	if err = randomize.Struct(seed, o, mempoolBinDBTypes, true, mempoolBinColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize MempoolBin struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := MempoolBins().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	mempoolBinDBTypes = map[string]string{`Time`: `bigint`, `Bin`: `character varying`, `NumberOfTransactions`: `integer`, `Size`: `integer`, `TotalFee`: `double precision`}
	_                 = bytes.MinRead
)

func testMempoolBinsUpdate(t *testing.T) {
	t.Parallel()

	if 0 == len(mempoolBinPrimaryKeyColumns) {
		t.Skip("Skipping table with no primary key columns")
	}
	if len(mempoolBinAllColumns) == len(mempoolBinPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &MempoolBin{}
	if err = randomize.Struct(seed, o, mempoolBinDBTypes, true, mempoolBinColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize MempoolBin struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := MempoolBins().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, mempoolBinDBTypes, true, mempoolBinPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize MempoolBin struct: %s", err)
	}

	if rowsAff, err := o.Update(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only affect one row but affected", rowsAff)
	}
}

func testMempoolBinsSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(mempoolBinAllColumns) == len(mempoolBinPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &MempoolBin{}
	if err = randomize.Struct(seed, o, mempoolBinDBTypes, true, mempoolBinColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize MempoolBin struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := MempoolBins().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, mempoolBinDBTypes, true, mempoolBinPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize MempoolBin struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(mempoolBinAllColumns, mempoolBinPrimaryKeyColumns) {
		fields = mempoolBinAllColumns
	} else {
		fields = strmangle.SetComplement(
			mempoolBinAllColumns,
			mempoolBinPrimaryKeyColumns,
		)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	typ := reflect.TypeOf(o).Elem()
	n := typ.NumField()

	updateMap := M{}
	for _, col := range fields {
		for i := 0; i < n; i++ {
			f := typ.Field(i)
			if f.Tag.Get("boil") == col {
				updateMap[col] = value.Field(i).Interface()
			}
		}
	}

	slice := MempoolBinSlice{o}
	if rowsAff, err := slice.UpdateAll(ctx, tx, updateMap); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("wanted one record updated but got", rowsAff)
	}
}

func testMempoolBinsUpsert(t *testing.T) {
	t.Parallel()

	if len(mempoolBinAllColumns) == len(mempoolBinPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	o := MempoolBin{}
	if err = randomize.Struct(seed, &o, mempoolBinDBTypes, true); err != nil {
		t.Errorf("Unable to randomize MempoolBin struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Upsert(ctx, tx, false, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert MempoolBin: %s", err)
	}

	count, err := MempoolBins().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &o, mempoolBinDBTypes, false, mempoolBinPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize MempoolBin struct: %s", err)
	}

	if err = o.Upsert(ctx, tx, true, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert MempoolBin: %s", err)
	}

	count, err = MempoolBins().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
