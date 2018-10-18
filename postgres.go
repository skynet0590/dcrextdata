package main

import (
	"database/sql"
	"log"
	"strings"

	_ "github.com/lib/pq"
)

type pgClient struct {
	db *sql.DB
}

var insertexchangedataStmt = "INSERT INTO exchangedata (high, low, open, close, time, exchange) VALUES ($1, $2, $3, $4, $5, $6);"
var createexchangedataStmt = "CREATE TABLE IF NOT EXISTS exchangedata (high FLOAT8, low FLOAT8, open FLOAT8, close FLOAT8, time INT, exchange VARCHAR(25), CONSTRAINT tick PRIMARY KEY (time, exchange));"

func initClient(psqlInfo string) (*pgClient, error) {
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}
	client := &pgClient{
		db: db,
	}

	return client, nil
}

func (c *pgClient) createExchangetable() error {
	_, err := c.db.Exec(createexchangedataStmt)
	return err
}
func tableExists(db *sql.DB, tableName string) (bool, error) {
	rows, err := db.Query(`select relname from pg_class where relname = $1`,
		tableName)
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

func (c *pgClient) close() {
	c.db.Close()
}
func (c *pgClient) addEntries(data []exchangeDataTick) error {
	for _, v := range data {
		_, err := c.db.Exec(insertexchangedataStmt, v.High, v.Low, v.Open, v.Close, v.Time, v.Exchange)
		if err != nil && !strings.Contains(err.Error(), "unique constraint") {
			return err
		}
	}
	return nil
}
