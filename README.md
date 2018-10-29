# Dcrextdata
- Collect and store poloniex and bittrex exchange data
- Watch and update exchange data records
- Collect and store mining pool stats

## Requirements
- [Go](golang.org/dl) 1.9+
- [Dep](https://github.com/golang/dep)
- [Postgresql](postgresql.org/download)

## Getting Started
With postgresql installed you can start the client with `sudo -u postgres psql` or you could `su` into the postgres user and run `psql` then execute the following sql statements:
```sql
    CREATE USER {user} WITH PASSWORD '{password}' CREATEDB;
    CREATE DATABASE {database} OWNER {user};
```

Copy the contents of the sample config file `sample-dcrextdata.conf` to the main config `dcrextdata.conf` (or just create a new one) and set the parameters for `dbuser`, `dbpass` and `dbname` to `{user}`, `{password}` and `{database}`.

Run `dep ensure` to get dependencies, `go build` to build and `./dcrextdata` to run. You can also drop the table(s) with `./dcrextdata  -D` .






