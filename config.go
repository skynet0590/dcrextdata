package main

import (
	"log"
	"os"

	flags "github.com/jessevdk/go-flags"
)

type config struct {
	DBHost     string `long:"dbhost" description:"Database host"`
	DBPort     int    `long:"dbport" description:"Database port"`
	DBUser     string `long:"dbuser" description:"Database username"`
	DBPass     string `long:"dbpass" description:"Database password"`
	DBName     string `long:"dbname" description:"Database name"`
	DropTables bool   `short:"D" long:"droptables" descripton:"Drop all database tables"`
	Quiet      bool   `short:"q" long:"quiet"`
}

var configfile = "dcrextdata.conf"

func loadConfig() (*config, error) {
	var cfg config
	parser := flags.NewParser(&cfg, flags.Default)
	err := flags.NewIniParser(parser).ParseFile(configfile)
	if err != nil {
		if _, ok := err.(*os.PathError); ok {
			log.Printf("Missing config file %s in current directory", configfile)
		} else {
			return nil, err
		}
	}

	_, err = parser.Parse()
	if err != nil {
		if e, ok := err.(*flags.Error); !ok || e.Type != flags.ErrHelp {
			parser.WriteHelp(os.Stderr)
		}
		return nil, err
	}

	return &cfg, nil
}
