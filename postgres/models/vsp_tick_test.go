// Code generated by SQLBoiler 3.5.0 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
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

func testVSPTicks(t *testing.T) {
	t.Parallel()

	query := VSPTicks()

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}

func testVSPTicksDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &VSPTick{}
	if err = randomize.Struct(seed, o, vspTickDBTypes, true, vspTickColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize VSPTick struct: %s", err)
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

	count, err := VSPTicks().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testVSPTicksQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &VSPTick{}
	if err = randomize.Struct(seed, o, vspTickDBTypes, true, vspTickColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize VSPTick struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := VSPTicks().DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := VSPTicks().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testVSPTicksSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &VSPTick{}
	if err = randomize.Struct(seed, o, vspTickDBTypes, true, vspTickColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize VSPTick struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := VSPTickSlice{o}

	if rowsAff, err := slice.DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := VSPTicks().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testVSPTicksExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &VSPTick{}
	if err = randomize.Struct(seed, o, vspTickDBTypes, true, vspTickColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize VSPTick struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	e, err := VSPTickExists(ctx, tx, o.ID)
	if err != nil {
		t.Errorf("Unable to check if VSPTick exists: %s", err)
	}
	if !e {
		t.Errorf("Expected VSPTickExists to return true, but got false.")
	}
}

func testVSPTicksFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &VSPTick{}
	if err = randomize.Struct(seed, o, vspTickDBTypes, true, vspTickColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize VSPTick struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	vspTickFound, err := FindVSPTick(ctx, tx, o.ID)
	if err != nil {
		t.Error(err)
	}

	if vspTickFound == nil {
		t.Error("want a record, got nil")
	}
}

func testVSPTicksBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &VSPTick{}
	if err = randomize.Struct(seed, o, vspTickDBTypes, true, vspTickColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize VSPTick struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = VSPTicks().Bind(ctx, tx, o); err != nil {
		t.Error(err)
	}
}

func testVSPTicksOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &VSPTick{}
	if err = randomize.Struct(seed, o, vspTickDBTypes, true, vspTickColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize VSPTick struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if x, err := VSPTicks().One(ctx, tx); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testVSPTicksAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	vspTickOne := &VSPTick{}
	vspTickTwo := &VSPTick{}
	if err = randomize.Struct(seed, vspTickOne, vspTickDBTypes, false, vspTickColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize VSPTick struct: %s", err)
	}
	if err = randomize.Struct(seed, vspTickTwo, vspTickDBTypes, false, vspTickColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize VSPTick struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = vspTickOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = vspTickTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := VSPTicks().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testVSPTicksCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	vspTickOne := &VSPTick{}
	vspTickTwo := &VSPTick{}
	if err = randomize.Struct(seed, vspTickOne, vspTickDBTypes, false, vspTickColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize VSPTick struct: %s", err)
	}
	if err = randomize.Struct(seed, vspTickTwo, vspTickDBTypes, false, vspTickColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize VSPTick struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = vspTickOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = vspTickTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := VSPTicks().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func testVSPTicksInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &VSPTick{}
	if err = randomize.Struct(seed, o, vspTickDBTypes, true, vspTickColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize VSPTick struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := VSPTicks().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testVSPTicksInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &VSPTick{}
	if err = randomize.Struct(seed, o, vspTickDBTypes, true); err != nil {
		t.Errorf("Unable to randomize VSPTick struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Whitelist(vspTickColumnsWithoutDefault...)); err != nil {
		t.Error(err)
	}

	count, err := VSPTicks().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testVSPTickToOneVSPUsingVSP(t *testing.T) {
	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var local VSPTick
	var foreign VSP

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, vspTickDBTypes, false, vspTickColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize VSPTick struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, vspDBTypes, false, vspColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize VSP struct: %s", err)
	}

	if err := foreign.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	local.VSPID = foreign.ID
	if err := local.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := local.VSP().One(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}

	if check.ID != foreign.ID {
		t.Errorf("want: %v, got %v", foreign.ID, check.ID)
	}

	slice := VSPTickSlice{&local}
	if err = local.L.LoadVSP(ctx, tx, false, (*[]*VSPTick)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if local.R.VSP == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.VSP = nil
	if err = local.L.LoadVSP(ctx, tx, true, &local, nil); err != nil {
		t.Fatal(err)
	}
	if local.R.VSP == nil {
		t.Error("struct should have been eager loaded")
	}
}

func testVSPTickToOneSetOpVSPUsingVSP(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a VSPTick
	var b, c VSP

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, vspTickDBTypes, false, strmangle.SetComplement(vspTickPrimaryKeyColumns, vspTickColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, vspDBTypes, false, strmangle.SetComplement(vspPrimaryKeyColumns, vspColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, vspDBTypes, false, strmangle.SetComplement(vspPrimaryKeyColumns, vspColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	for i, x := range []*VSP{&b, &c} {
		err = a.SetVSP(ctx, tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.VSP != x {
			t.Error("relationship struct not set to correct value")
		}

		if x.R.VSPTicks[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if a.VSPID != x.ID {
			t.Error("foreign key was wrong value", a.VSPID)
		}

		zero := reflect.Zero(reflect.TypeOf(a.VSPID))
		reflect.Indirect(reflect.ValueOf(&a.VSPID)).Set(zero)

		if err = a.Reload(ctx, tx); err != nil {
			t.Fatal("failed to reload", err)
		}

		if a.VSPID != x.ID {
			t.Error("foreign key was wrong value", a.VSPID, x.ID)
		}
	}
}

func testVSPTicksReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &VSPTick{}
	if err = randomize.Struct(seed, o, vspTickDBTypes, true, vspTickColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize VSPTick struct: %s", err)
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

func testVSPTicksReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &VSPTick{}
	if err = randomize.Struct(seed, o, vspTickDBTypes, true, vspTickColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize VSPTick struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := VSPTickSlice{o}

	if err = slice.ReloadAll(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testVSPTicksSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &VSPTick{}
	if err = randomize.Struct(seed, o, vspTickDBTypes, true, vspTickColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize VSPTick struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := VSPTicks().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	vspTickDBTypes = map[string]string{`ID`: `integer`, `VSPID`: `integer`, `Immature`: `integer`, `Live`: `integer`, `Voted`: `integer`, `Missed`: `integer`, `PoolFees`: `double precision`, `ProportionLive`: `double precision`, `ProportionMissed`: `double precision`, `UserCount`: `integer`, `UsersActive`: `integer`, `Time`: `timestamp with time zone`}
	_              = bytes.MinRead
)

func testVSPTicksUpdate(t *testing.T) {
	t.Parallel()

	if 0 == len(vspTickPrimaryKeyColumns) {
		t.Skip("Skipping table with no primary key columns")
	}
	if len(vspTickAllColumns) == len(vspTickPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &VSPTick{}
	if err = randomize.Struct(seed, o, vspTickDBTypes, true, vspTickColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize VSPTick struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := VSPTicks().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, vspTickDBTypes, true, vspTickPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize VSPTick struct: %s", err)
	}

	if rowsAff, err := o.Update(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only affect one row but affected", rowsAff)
	}
}

func testVSPTicksSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(vspTickAllColumns) == len(vspTickPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &VSPTick{}
	if err = randomize.Struct(seed, o, vspTickDBTypes, true, vspTickColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize VSPTick struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := VSPTicks().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, vspTickDBTypes, true, vspTickPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize VSPTick struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(vspTickAllColumns, vspTickPrimaryKeyColumns) {
		fields = vspTickAllColumns
	} else {
		fields = strmangle.SetComplement(
			vspTickAllColumns,
			vspTickPrimaryKeyColumns,
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

	slice := VSPTickSlice{o}
	if rowsAff, err := slice.UpdateAll(ctx, tx, updateMap); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("wanted one record updated but got", rowsAff)
	}
}

func testVSPTicksUpsert(t *testing.T) {
	t.Parallel()

	if len(vspTickAllColumns) == len(vspTickPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	o := VSPTick{}
	if err = randomize.Struct(seed, &o, vspTickDBTypes, true); err != nil {
		t.Errorf("Unable to randomize VSPTick struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Upsert(ctx, tx, false, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert VSPTick: %s", err)
	}

	count, err := VSPTicks().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &o, vspTickDBTypes, false, vspTickPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize VSPTick struct: %s", err)
	}

	if err = o.Upsert(ctx, tx, true, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert VSPTick: %s", err)
	}

	count, err = VSPTicks().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
