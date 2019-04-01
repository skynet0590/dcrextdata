// Copyright (c) 2018-2019 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package postgres

import (
	"database/sql"
	"time"

	"github.com/volatiletech/sqlboiler/queries/qm"

	"github.com/raedahgroup/dcrextdata/postgres/models"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/types"

	"github.com/raedahgroup/dcrextdata/vsp"
)

type insertableG interface {
	InsertG(boil.Columns) error
}

type upsertableG interface {
	UpsertG(updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error
}

// StoreVSPs attemps to store the vsp responses by calling storeVspResponseG and returning
// a slice of errors
func (pg *PgDb) StoreVSPs(data vsp.Response) []error {
	errs := make([]error, 0, len(data))
	for name, tick := range data {
		err := storeVspResponseG(name, tick)
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

func storeVspResponseG(name string, resp *vsp.ResposeData) error {
	txr, err := boil.Begin()
	if err != nil {
		return err
	}
	pool, err := models.VSPS(models.VSPWhere.Name.EQ(name)).OneG()
	if err == sql.ErrNoRows {
		pool = responseToVSP(name, resp)
		err := tryInsertG(txr, pool)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	vspTick := responseToVSPTick(pool.ID, resp)
	t := int64ToTime(resp.LastUpdated)

	tw := models.VSPTickWhere
	vTick, err := models.VSPTicks(qm.Expr(
		tw.VSPID.EQ(pool.ID),
		tw.Immature.EQ(vspTick.Immature),
		tw.Live.EQ(vspTick.Live),
		tw.Voted.EQ(vspTick.Voted),
		tw.Missed.EQ(vspTick.Missed),
		tw.PoolFees.EQ(vspTick.PoolFees),
		tw.ProportionLive.EQ(vspTick.ProportionLive),
		tw.ProportionMissed.EQ(vspTick.ProportionMissed),
		tw.UserCount.EQ(vspTick.UserCount),
		tw.UsersActive.EQ(vspTick.UsersActive))).OneG()

	if err == sql.ErrNoRows {
		err = tryInsertG(txr, vspTick)
		if err != nil {
			return err
		}
		vTick = vspTick
	} else if err != nil {
		return err
	}

	tickTimeExists, err := models.VSPTickTimeExistsG(vTick.ID, t)
	if err != nil {
		return err
	}
	if tickTimeExists {
		return vsp.PoolTickTimeExistsError{
			PoolName: name,
			TickTime: t,
		}
	}

	tickTime := &models.VSPTickTime{
		VSPTickID:  vTick.ID,
		UpdateTime: t,
	}

	err = tryInsertG(txr, tickTime)
	if err != nil {
		return err
	}

	pool.LastUpdate = t

	err = pool.UpsertG(true, nil, boil.Infer(), boil.Infer())
	if err != nil {
		return err
	}

	err = txr.Commit()
	if err != nil {
		return err
	}

	log.Infof("Added complete pool data for %s at %s", name, t.UTC().String())
	return nil
}

func tryInsertG(txr boil.Transactor, data insertableG) error {
	err := data.InsertG(boil.Infer())
	if err != nil {
		errT := txr.Rollback()
		if errT != nil {
			return errT
		}
		return err
	}
	return nil
}

func tryUpsertG(txr boil.Transactor, data upsertableG) error {
	err := data.UpsertG(true, nil, boil.Infer(), boil.Infer())
	if err != nil {
		errT := txr.Rollback()
		if errT != nil {
			return errT
		}
		return err
	}
	return nil
}

func responseToVSP(name string, resp *vsp.ResposeData) *models.VSP {
	return &models.VSP{
		Name:                 name,
		APIEnabled:           resp.APIEnabled,
		APIVersionsSupported: types.Int64Array(resp.APIVersionsSupported),
		Network:              resp.Network,
		URL:                  resp.URL,
		Launched:             int64ToTime(resp.Launched),
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
	}
}

func int64ToTime(t int64) time.Time {
	return time.Unix(t, 0)
}
