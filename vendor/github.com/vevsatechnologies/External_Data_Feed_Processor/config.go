package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/decred/dcrd/dcrutil"
)

const (
	defaultConfigFilename = "dcrextdata.conf"
)

var (
	defaultHomeDir = dcrutil.AppDataDir("dcrdata", false)

	defaultConfigFile = filepath.Join(defaultHomeDir, defaultConfigFilename)

	defaultDBFileName = "data_feed_processor.sql"

	defaultPGHost   = "127.0.0.1:5432"
	defaultPGUser   = "postgres"
	defaultPGPass   = ""
	defaultPGDBName = "data_feed_processor"
)

type config struct {
	// General application behavior
	HomeDir    string `short:"A" long:"appdata" description:"Path to application home directory"`
	DBFileName string `long:"dbfile" description:"SQLite DB file name (default is dcrdata.sqlt.db)."`
	ConfigFile string `short:"C" long:"configfile" description:"Path to configuration file"`
	PGDBName   string `long:"pgdbname" description:"PostgreSQL DB name."`
	PGUser     string `long:"pguser" description:"PostgreSQL DB user."`
	PGPass     string `long:"pgpass" description:"PostgreSQL DB password."`
	PGHost     string `long:"pghost" description:"PostgreSQL server host:port or UNIX socket (e.g. /run/postgresql)."`
}

var (
	defaultConfig = config{
		HomeDir:    defaultHomeDir,
		ConfigFile: defaultConfigFile,
		DBFileName: defaultDBFileName,
		PGDBName:   defaultPGDBName,
		PGUser:     defaultPGUser,
		PGPass:     defaultPGPass,
		PGHost:     defaultPGHost,
	}
)

// loadConfig initializes and parses the config using a config file and command
// line options.
func loadConfig() (*config, error) {
	loadConfigError := func(err error) (*config, error) {
		return nil, err
	}

	// Default config.
	cfg := defaultConfig

	// Pre-parse the command line options to see if an alternative config
	// file or the version flag was specified.
	preCfg := cfg
	preParser := flags.NewParser(&preCfg, flags.HelpFlag|flags.PassDoubleDash)
	_, err := preParser.Parse()
	if err != nil {
		e, ok := err.(*flags.Error)
		if !ok || e.Type != flags.ErrHelp {
			preParser.WriteHelp(os.Stderr)
		}
		if ok && e.Type == flags.ErrHelp {
			preParser.WriteHelp(os.Stdout)
			os.Exit(0)
		}
		return loadConfigError(err)
	}

	// Show the version and exit if the version flag was specified.
	appName := filepath.Base(os.Args[0])
	appName = strings.TrimSuffix(appName, filepath.Ext(appName))
	if preCfg.ShowVersion {
		fmt.Printf("%s version %s (Go version %s)\n", appName,
			version.Ver.String(), runtime.Version())
		os.Exit(0)
	}

	// Load additional config from file.
	var configFileError error
	// Config file name for logging.
	configFile := "NONE (defaults)"
	parser := flags.NewParser(&cfg, flags.Default)
	if _, err := os.Stat(preCfg.ConfigFile); os.IsNotExist(err) {
		// Non-default config file must exist
		if defaultConfig.ConfigFile != preCfg.ConfigFile {
			fmt.Fprintln(os.Stderr, err)
			return loadConfigError(err)
		}
		// Warn about missing default config file, but continue
		fmt.Printf("Config file (%s) does not exist. Using defaults.\n",
			preCfg.ConfigFile)
	} else {
		// The config file exists, so attempt to parse it.
		err = flags.NewIniParser(parser).ParseFile(preCfg.ConfigFile)
		if err != nil {
			if _, ok := err.(*os.PathError); !ok {
				fmt.Fprintln(os.Stderr, err)
				parser.WriteHelp(os.Stderr)
				return loadConfigError(err)
			}
			configFileError = err
		}
		configFile = preCfg.ConfigFile
	}

	// Parse command line options again to ensure they take precedence.
	_, err = parser.Parse()
	if err != nil {
		if e, ok := err.(*flags.Error); !ok || e.Type != flags.ErrHelp {
			parser.WriteHelp(os.Stderr)
		}
		return loadConfigError(err)
	}

	// Create the home directory if it doesn't already exist.
	funcName := "loadConfig"
	err = os.MkdirAll(cfg.HomeDir, 0700)
	if err != nil {
		// Show a nicer error message if it's because a symlink is linked to a
		// directory that does not exist (probably because it's not mounted).
		if e, ok := err.(*os.PathError); ok && os.IsExist(err) {
			if link, lerr := os.Readlink(e.Path); lerr == nil {
				str := "is symlink %s -> %s mounted?"
				err = fmt.Errorf(str, e.Path, link)
			}
		}

		str := "%s: failed to create home directory: %v"
		err := fmt.Errorf(str, funcName, err)
		fmt.Fprintln(os.Stderr, err)
		return nil, err
	}

	// Warn about missing config file after the final command line parse
	// succeeds.  This prevents the warning on help messages and invalid
	// options.
	if configFileError != nil {
		fmt.Printf("%v\n", configFileError)
		return loadConfigError(configFileError)
	}

	return &cfg, nil
}
