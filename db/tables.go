package db

import (
	"fmt"

	"github.com/raedahgroup/dcrextdata/db/internal"
	log "github.com/sirupsen/logrus"
)

func (pg *PgDb) CreateExchangeDataTable() error {
	_, err := pg.db.Exec(internal.CreateExchangeDataTable)
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

func (pg *PgDb) dropTable(name string) error {
	_, err := pg.db.Exec(fmt.Sprintf(`DROP TABLE IF EXISTS %s;`, name))
	return err
}

func (pg *PgDb) DropAllTables() error {
	// TODO: Add the other tables
	return pg.dropTable("exchange_data")
}
