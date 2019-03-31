// Copyright (c) 2018-2019 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
)

const (
	// Exchange Table
	LastExchangeEntryTime   = `SELECT time FROM exchange_data WHERE exchange=$1 ORDER BY time DESC LIMIT 1`
	InsertExchangeDataTick  = `INSERT INTO exchange_data (high, low, open, close, time, exchange) VALUES ($1, $2, $3, $4, $5, $6)`
	CreateExchangeDataTable = `CREATE TABLE IF NOT EXISTS exchange_data (high FLOAT8, low FLOAT8, open FLOAT8, close FLOAT8, time INT, exchange VARCHAR(25), CONSTRAINT tick PRIMARY KEY (time, exchange))`
)

type PgDb struct {
	db *sql.DB
}

// Core methods

//
func NewPgDb(host, port, user, pass, dbname string) (*PgDb, error) {
	var psqlInfo string
	if pass == "" {
		psqlInfo = fmt.Sprintf("host=%s user=%s "+
			"dbname=%s sslmode=disable",
			host, user, dbname)
	} else {
		psqlInfo = fmt.Sprintf("host=%s user=%s "+
			"password=%s dbname=%s sslmode=disable",
			host, user, pass, dbname)
	}
	// Only add port arg fot TCP connection since UNIX domain sockets (specified
	// by a "/" prefix) do not have a port.
	if !strings.HasPrefix(host, "/") {
		psqlInfo += fmt.Sprintf(" port=%s", port)
	}

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	err = db.Ping()

	return &PgDb{db}, err
}

func (pg *PgDb) Close() error {
	pqLog.Trace("Closing postgresql connection")
	return pg.db.Close()
}

func (pg *PgDb) dropTable(name string) error {
	pqLog.Tracef("Dropping table %s", name)
	_, err := pg.db.Exec(fmt.Sprintf(`DROP TABLE IF EXISTS %s;`, name))
	return err
}

func (pg *PgDb) DropAllTables() error {
	// TODO: Add the other tables
	return pg.dropTable("exchange_data")
}

// Exchange methods

//
func (pg *PgDb) AddExchangeData(data []DataTick) error {
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
		pqLog.Infof("Added %d entry from %s (%s)", added, data[0].Exchange, UnixTimeToString(data[0].Time))
	} else {
		last := data[len(data)-1]
		pqLog.Infof("Added %d entries from %s (%s to %s)", added, last.Exchange, UnixTimeToString(data[0].Time), UnixTimeToString(last.Time))
	}

	return nil
}

func (pg *PgDb) LastExchangeEntryTime(exchange string) (time int64) {
	rows := pg.db.QueryRow(LastExchangeEntryTime, exchange)
	_ = rows.Scan(&time)
	return
}

func (pg *PgDb) CreateExchangeDataTable() error {
	pqLog.Trace("Creating exchange data table")
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