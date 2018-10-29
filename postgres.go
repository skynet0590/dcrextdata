package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/lib/pq"
)

type PgDb struct {
	*sql.DB
}

var (
	insertExchangeDataStmt      = `INSERT INTO exchange_data (high, low, open, close, time, exchange) VALUES ($1, $2, $3, $4, $5, $6)`
	createExchangeDataStmt      = `CREATE TABLE IF NOT EXISTS exchange_data (high FLOAT8, low FLOAT8, open FLOAT8, close FLOAT8, time INT, exchange VARCHAR(25), CONSTRAINT tick PRIMARY KEY (time, exchange))`
	getLastExchangeDataTimeStmt = `SELECT time FROM exchange_data ORDER BY time DESC LIMIT 1`
)

func NewPgDb(psqlInfo string) (PgDb, error) {
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return PgDb{nil}, err
	}
	return PgDb{db}, nil
}

func (db *PgDb) CreateExchangeDataTable() error {
	_, err := db.Exec(createExchangeDataStmt)
	return err
}

func (db *PgDb) tableExists(name string) (bool, error) {
	rows, err := db.Query(`SELECT relname FROM pg_class WHERE relname = $1`, name)
	if err == nil {
		defer func() {
			if e := rows.Close(); e != nil {
				log.Printf("Close of Query failed: %v", e)
			}
		}()
		return rows.Next(), nil
	}
	return false, err
}

func (db *PgDb) ExchangeDataTableExits() bool {
	exists, _ := db.tableExists("exchange_data")
	return exists
}

func (db *PgDb) AddExchangeData(data []exchangeDataTick) error {
	for _, v := range data {
		_, err := db.Exec(insertExchangeDataStmt, v.High, v.Low, v.Open, v.Close, v.Time, v.Exchange)
		if err != nil && !strings.Contains(err.Error(), "unique constraint") {
			return err
		}
	}
	return nil
}

func (db *PgDb) LastExchangeEntryTime() (int64, error) {
	var time int64 = -1
	rows := db.QueryRow(getLastExchangeDataTimeStmt)
	err := rows.Scan(&time)

	if err != nil {
		return time, err
	}
	return time, nil
}

func (db *PgDb) DropTable(name string) error {
	_, err := db.Exec(fmt.Sprintf(`DROP TABLE IF EXISTS %s;`, name))
	return err
}

func (db *PgDb) DropExchangeDataTable() error {
	return db.DropTable("exchange_data")
}
