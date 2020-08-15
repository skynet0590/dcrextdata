package postgres

import (
	"fmt"
	"strconv"
	"strings"

	migrate "github.com/rubenv/sql-migrate"
)

const (
	migrateUp   = "up"
	migrateDown = "down"

	lastExchangeTickEntryTime = `SELECT time FROM exchange_tick ORDER BY time DESC LIMIT 1`
	lastExchangeEntryID       = `SELECT id FROM exchange ORDER BY id DESC LIMIT 1`
	lastVspTickEntryTime      = `SELECT time FROM vsp_tick ORDER BY time DESC LIMIT 1`
	lastPowEntryTimeBySource  = `SELECT time FROM pow_data WHERE source=$1 ORDER BY time DESC LIMIT 1`
	lastPowEntryTime          = `SELECT time FROM pow_data ORDER BY time DESC LIMIT 1`
	lastMempoolBlockHeight    = `SELECT last_block_height FROM mempool ORDER BY last_block_height DESC LIMIT 1`
	lastMempoolEntryTime      = `SELECT time FROM mempool ORDER BY time DESC LIMIT 1`
	lastCommStatEntryTime     = `SELECT date FROM reddit ORDER BY date DESC LIMIT 1`
)

var (
	ErrInvalidMigrateConvention = fmt.Errorf("Invalid migrate convention")
)

func (pg *PgDb) DropCacheTables() error {
	// vsp_tick
	if err := pg.dropTable("vsp_tick_bin"); err != nil {
		return err
	}

	// pow_bin
	if err := pg.dropTable("pow_bin"); err != nil {
		return err
	}

	// mempool_bin
	if err := pg.dropTable("mempool_bin"); err != nil {
		return err
	}

	// propagation
	if err := pg.dropTable("propagation"); err != nil {
		return err
	}

	// vote_receive_time_deviation
	if err := pg.dropTable("vote_receive_time_deviation"); err != nil {
		return err
	}

	//network_snapshot_bin
	if err := pg.dropTable("network_snapshot_bin"); err != nil {
		return err
	}

	return nil
}

func (pg *PgDb) dropTable(name string) error {
	log.Tracef("Dropping table %s", name)
	_, err := pg.db.Exec(fmt.Sprintf(`DROP TABLE IF EXISTS %s;`, name))
	return err
}

func (pg *PgDb) dropIndex(name string) error {
	log.Tracef("Dropping table %s", name)
	_, err := pg.db.Exec(fmt.Sprintf(`DROP INDEX IF EXISTS %s;`, name))
	return err
}

func (pg *PgDb) MigrateDatabase(migrateCode string) error {
	migrateInfos := strings.Split(strings.ToLower(migrateCode), ":")
	var migrateAction migrate.MigrationDirection
	switch migrateInfos[0] {
	case migrateUp:
		migrateAction = migrate.Up
	case migrateDown:
		migrateAction = migrate.Down
	default:
		return ErrInvalidMigrateConvention
	}

	migrations := &migrate.FileMigrationSource{
		Dir: "postgres/migrations",
	}
	var n int
	var err error
	if len(migrateInfos) == 1 {
		n, err = migrate.Exec(pg.db, "postgres", migrations, migrateAction)
	} else if len(migrateInfos) == 2 {
		var limit int
		limit, err = strconv.Atoi(migrateInfos[1])
		if err != nil {
			return ErrInvalidMigrateConvention
		}
		n, err = migrate.ExecMax(pg.db, "postgres", migrations, migrateAction, limit)
	}
	if err != nil {
		return err
	}
	fmt.Printf("Applied %d migrations!\n", n)
	return nil
}
