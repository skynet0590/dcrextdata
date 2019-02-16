package db

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
)

type PgDb struct {
	db *sql.DB
}

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

	return &PgDb{db}, nil
}

func (pg *PgDb) Close() error {
	return pg.db.Close()
}
