DROP TABLE IF EXISTS vsp_tick, vsp, exchange, exchange_tick, vsp_tick_time;

CREATE TABLE IF NOT EXISTS exchange (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    url TEXT NOT NULL,
    start_time TIMESTAMPTZ NOT NULL,
    last_updated TIMESTAMPTZ NOT NULL
);


CREATE TABLE IF NOT EXISTS exchange_tick (
    id SERIAL PRIMARY KEY,
    exchange_id INT REFERENCES exchange(id), 
	high FLOAT NOT NULL,
	low FLOAT NOT NULL,
	open FLOAT NOT NULL,
	close FLOAT NOT NULL,
	time TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS vsp (
	id SERIAL PRIMARY KEY,
	name TEXT NOT NULL,
	api_enabled BOOLEAN NOT NULL,
	api_versions_supported INT[] NOT NULL,
	network TEXT NOT NULL,
	url TEXT NOT NULL,
	launched TIMESTAMPTZ NOT NULL,
    last_update TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS vsp_tick (
	id SERIAL PRIMARY KEY,
	vsp_id INT REFERENCES vsp(id) NOT NULL,
	immature INT NOT NULL,
	live INT NOT NULL,
	voted INT NOT NULL,
	missed INT NOT NULL,
	pool_fees FLOAT NOT NULL,
	proportion_live FLOAT NOT NULL,
	proportion_missed FLOAT NOT NULL,
	user_count INT NOT NULL,
	users_active INT NOT NULL
);

CREATE TABLE IF NOT EXISTS vsp_tick_time (
	vsp_tick_id INT REFERENCES vsp_tick(id) NOT NULL,
	update_time TIMESTAMPTZ NOT NULL,
	CONSTRAINT tick_time PRIMARY KEY (vsp_tick_id, update_time)
);