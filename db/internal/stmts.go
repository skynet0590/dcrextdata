package internal

const (
	// Exchange Table
	LastExchangeEntryTime = `SELECT time FROM exchange_data WHERE exchange=$1 ORDER BY time DESC LIMIT 1`
	//LastExchangeEntryTime   = `SELECT time FROM exchange_data ORDER BY time DESC LIMIT 1`
	InsertExchangeDataTick  = `INSERT INTO exchange_data (high, low, open, close, time, exchange) VALUES ($1, $2, $3, $4, $5, $6)`
	CreateExchangeDataTable = `CREATE TABLE IF NOT EXISTS exchange_data (high FLOAT8, low FLOAT8, open FLOAT8, close FLOAT8, time INT, exchange VARCHAR(25), CONSTRAINT tick PRIMARY KEY (time, exchange))`
)
