// Copyright (c) 2018-2019 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package config

import (
	"fmt"
	"os"
	
	flags "github.com/jessevdk/go-flags"
)

const (
	DefaultConfigFilename = "dcrextdata.conf"
	defaultLogFilename    = "dcrextdata.log"
	defaultLogLevel       = "info"
	Hint                  = `Run dcrextdata < --http > to start http server or dcrextdata < --help > for help.`
)

type Config struct {
	ConfigFileOptions
	CommandLineOptions
}

type ConfigFileOptions struct {
	// General application behaviour
	LogFile    string `short:"L" long:"logfile" description:"File name of the log file"`
	ConfigFile string `short:"C" long:"Configfile" description:"Path to Configuration file"`
	DebugLevel string `short:"d" long:"debuglevel" description:"Logging level {trace, debug, info, warn, error, critical}"`
	Quiet      bool   `short:"q" long:"quiet" description:"Easy way to set debuglevel to error"`

	// Postgresql Configuration
	DBHost string `long:"dbhost" description:"Database host"`
	DBPort string `long:"dbport" description:"Database port"`
	DBUser string `long:"dbuser" description:"Database username"`
	DBPass string `long:"dbpass" description:"Database password"`
	DBName string `long:"dbname" description:"Database name"`

	// Http Server
	HTTPHost string `long:"httphost" description:"HTTP server host address or IP when running godcr in http mode."`
	HTTPPort string `long:"httpport" description:"HTTP server port when running godcr in http mode."`
	// Exchange collector
	DisableExchangeTicks bool     `long:"disablexcticks" decription:"Disables collection of ticker data from exchanges"`
	DisabledExchanges    []string `long:"disableexchange" description:"Disable data collection for this exchange"`

	// PoW collector
	DisablePow   bool     `long:"disablepow" description:"Disables collection of data for pows"`
	DisabledPows []string `long:"disabledpow" description:"Disable data collection for this Pow"`
	PowInterval  int64    `long:"powI" description:"Collection interval for Pow"`

	// VSP
	DisableVSP  bool  `long:"disablevsp" description:"Disables periodic voting service pool status collection"`
	VSPInterval int64 `long:"vspinterval" description:"Collection interval for pool status collection"`

	// Mempool
	DisableMempool  bool   `long:"disablemempool" description:"Disable mempool data collection"`
	DcrdRpcServer   string `long:"dcrdrpcserver" description:"Dcrd rpc server host"`
	DcrdNetworkType string `long:"dcrdnetworktype" description:"Dcrd rpc network type"`
	DcrdRpcUser     string `long:"dcrdrpcuser" description:"Your Dcrd rpc username"`
	DcrdRpcPassword string `long:"dcrdrpcpaswword" description:"Your Dcrd rpc password"`
}

// CommandLineOptions holds the top-level options/flags that are displayed on the command-line menu
type CommandLineOptions struct {
	Reset      bool   `short:"R" long:"reset" description:"Drop all database tables and start over"`
	HttpMode bool `long:"http" description:"Launch http server"`
}

func defaultFileOptions() ConfigFileOptions {
	return ConfigFileOptions{
		LogFile:         defaultLogFilename,
		ConfigFile:      DefaultConfigFilename,
		DebugLevel:      defaultLogLevel,
		VSPInterval:     300,
		PowInterval:     300,
		DcrdNetworkType: "mainnet",
	}
}

func defaultConfig() Config {
	return Config{
		ConfigFileOptions: defaultFileOptions(),
	}
}

func LoadConfig() (*Config, []string, error) {
	cfg := defaultConfig()
	parser := flags.NewParser(&cfg, flags.IgnoreUnknown)
	err := flags.NewIniParser(parser).ParseFile(cfg.ConfigFileOptions.ConfigFile)
	if err != nil {
		if _, ok := err.(*os.PathError); ok {
			fmt.Printf("Missing Config file %s in current directory\n", cfg.ConfigFileOptions.ConfigFile)
		} else {
			return nil, nil, err
		}
	}

	unknownArg, err := parser.Parse()
	if err != nil {
		e, ok := err.(*flags.Error)
		if ok && e.Type == flags.ErrHelp {
			os.Exit(0)
		}
		return nil, nil, err
	}

	return &cfg, unknownArg, nil
}
