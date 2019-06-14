package postgres

import "fmt"

const (
	CreateExchangeDataTable = `CREATE TABLE IF NOT EXISTS exchange_data (
		high FLOAT8,
		low FLOAT8,
		open FLOAT8,
		close FLOAT8,
		time INT,
		exchange TEXT, 
		CONSTRAINT tick PRIMARY KEY (time, exchange));`

	createVSPInfoTable = `CREATE TABLE IF NOT EXISTS vsp (
		id SERIAL PRIMARY KEY,
		name TEXT,
		api_enabled BOOLEAN,
		api_versions_supported INT8[],
		network TEXT,
		url TEXT,
		launched TIMESTAMPTZ
	);`

	createVSPDataTable = `CREATE TABLE IF NOT EXISTS vsp_data (
		id SERIAL PRIMARY KEY,
		vsp_id INT REFERENCES vsp(id),
		last_updated TIMESTAMPTZ,
		immature INT8,
		live INT8,
		voted INT8,
		missed INT8,
		pool_fees FLOAT8,
		proportion_live FLOAT8,
		proportion_missed FLOAT8,
		user_count INT8,
		users_active INT8,
		time TIMESTAMPTZ
	);`

	// PoW table
	CreatePowDataTable = `CREATE TABLE IF NOT EXISTS pow_data (
		time INT,
		network_hashrate INT,
		pool_hashrate INT8,
		workers INT,
		network_difficulty FLOAT8,
		coin_price VARCHAR(25),
		btc_price VARCHAR(25),
		source VARCHAR(25),
		PRIMARY KEY (time, source)
	);`

	InsertPowData = `INSERT INTO pow_data (
		time, network_hashrate, pool_hashrate, workers, network_difficulty, coin_price, btc_price, source)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	LastPowEntryTime = `SELECT time FROM pow_data WHERE source=$1 ORDER BY time DESC LIMIT 1`
)

func (pg *PgDb) CreatePowDataTable() error {
	_, err := pg.db.Exec(CreatePowDataTable)
	return err
}

func (pg *PgDb) PowDataTableExits() bool {
	exists, _ := pg.tableExists("pow_data")
	return exists
}

func (pg *PgDb) CreateExchangeDataTable() error {
	log.Trace("Creating exchange data table")
	_, err := pg.db.Exec(CreateExchangeDataTable)
	return err
}

func (pg *PgDb) ExchangeDataTableExits() bool {
	exists, _ := pg.tableExists("exchange_data")
	return exists
}

func (pg *PgDb) CreateVSPTables() error {
	_, err := pg.db.Exec(createVSPInfoTable)
	if err != nil {
		return err
	}

	_, err = pg.db.Exec(createVSPDataTable)
	if err != nil {
		return err
	}

	return nil
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
	if err := pg.dropTable("pow_data"); err != nil {
		return err
	}
	return pg.dropTable("exchange_data")
}