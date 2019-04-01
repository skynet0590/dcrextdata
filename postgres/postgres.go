// Copyright (c) 2018-2019 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package postgres

//go:generate sqlboiler --wipe psql --add-global-variants --no-hooks --no-context --no-auto-timestamps

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/raedahgroup/dcrextdata/exchanges"
	"github.com/volatiletech/sqlboiler/boil"
)

const (
	// Exchange Table
	LastExchangeEntryTime = `SELECT time FROM exchange_data WHERE exchange=$1 ORDER BY time DESC LIMIT 1`

	InsertExchangeDataTick = `INSERT INTO exchange_data (
		high, low, open, close, time, exchange)
	VALUES ($1, $2, $3, $4, $5, $6)`

	CreateExchangeDataTable = `CREATE TABLE IF NOT EXISTS exchange_data (
		high FLOAT8,
		low FLOAT8,
		open FLOAT8,
		close FLOAT8,
		time INT,
		exchange TEXT, 
		CONSTRAINT tick PRIMARY KEY (time, exchange));`
)

type PgDb struct {
	db *sql.DB
}

// Core methods

//
func NewPgDb(host, port, user, pass, dbname string) (*PgDb, error) {
	db, err := Connect(host, port, user, pass, dbname)
	if err != nil {
		return nil, err
	}
	boil.SetDB(db)
	return &PgDb{db}, nil
}

func (pg *PgDb) Close() error {
	log.Trace("Closing postgresql connection")
	return pg.db.Close()
}

func (pg *PgDb) dropTable(name string) error {
	log.Tracef("Dropping table %s", name)
	_, err := pg.db.Exec(fmt.Sprintf(`DROP TABLE IF EXISTS %s;`, name))
	return err
}

func (pg *PgDb) DropAllTables() error {
	if err := pg.dropTable("vsp_data"); err != nil {
		return err
	}
	if err := pg.dropTable("vsp"); err != nil {
		return err
	}
	return pg.dropTable("exchange_data")
}

// Exchange methods

//
func (pg *PgDb) AddExchangeData(data []exchanges.DataTick) error {
	added := 0
	for _, v := range data {
		_, err := pg.db.Exec(InsertExchangeDataTick, v.High, v.Low, v.Open, v.Close, v.Time, v.Exchange)
		if err != nil {
			if !strings.Contains(err.Error(), "unique constraint") { // Ignore duplicate entries
				return err
			}
		}
		added++
	}
	if len(data) == 1 {
		log.Infof("Added %d entry from %s (%s)", added, data[0].Exchange,
			UnixTimeToString(data[0].Time))
	} else {
		last := data[len(data)-1]
		log.Infof("Added %d entries from %s (%s to %s)", added, last.Exchange,
			UnixTimeToString(data[0].Time), UnixTimeToString(last.Time))
	}

	return nil
}

func (pg *PgDb) LastExchangeEntryTime(exchange string) (time int64) {
	rows := pg.db.QueryRow(LastExchangeEntryTime, exchange)
	_ = rows.Scan(&time)
	return
}

func (pg *PgDb) CreateExchangeDataTable() error {
	log.Trace("Creating exchange data table")
	_, err := pg.db.Exec(CreateExchangeDataTable)
	return err
}

func (pg *PgDb) tableExists(name string) (bool, error) {
	rows, err := pg.db.Query(`SELECT relname FROM pg_class WHERE relname = $1`, name)
	if err == nil {
		defer func() {
			if e := rows.Close(); e != nil {
				log.Error("Close of Query failed: ", e)
			}
		}()
		return rows.Next(), nil
	}
	return false, err
}

func (pg *PgDb) ExchangeDataTableExits() bool {
	exists, _ := pg.tableExists("exchange_data")
	return exists
}
