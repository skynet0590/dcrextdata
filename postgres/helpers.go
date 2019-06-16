// Copyright (c) 2013-2015 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package postgres

import (
	"context"
	"strings"
	"time"

	"github.com/volatiletech/sqlboiler/boil"
)

const dateTemplate = "2006-01-02 15:04"

type insertable interface {
	Insert(context.Context, boil.ContextExecutor, boil.Columns) error
}

type upsertable interface {
	Upsert(ctx context.Context, db boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error
}

func UnixTimeToString(t int64) string {
	return time.Unix(t, 0).UTC().Format(dateTemplate)
}

func int64ToTime(t int64) time.Time {
	return time.Unix(t, 0)
}

func (pg *PgDb) tryInsert(ctx context.Context, txr boil.Transactor, data insertable) error {
	err := data.Insert(ctx, pg.db, boil.Infer())
	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			return err
		}
		errT := txr.Rollback()
		if errT != nil {
			return errT
		}
		return err
	}
	return nil
}

func (pg *PgDb) tryUpsert(ctx context.Context, txr boil.Transactor, data upsertable) error {
	err := data.Upsert(ctx, pg.db, true, nil, boil.Infer(), boil.Infer())
	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			return err
		}
		errT := txr.Rollback()
		if errT != nil {
			return errT
		}
		return err
	}
	return nil
}
