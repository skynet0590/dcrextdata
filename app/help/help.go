package help

import (
	"fmt"
	"io"

	"github.com/jessevdk/go-flags"
	"github.com/raedahgroup/dcrextdata/app"
	"github.com/raedahgroup/dcrextdata/app/config"
)

type GeneralHelpData struct {
	config.CommandLineOptions `group:"Command-Line options:"`
	config.ConfigFileOptions  `group:"Configuration file options:"`
}

func PrintGeneralHelp(output io.Writer, parser *flags.Parser) {
	tabWriter := TabWriter(output)

	// print version
	fmt.Fprintf(tabWriter, "%s v%s\n", app.AppName, app.Version())
	fmt.Fprintln(tabWriter)

	// print general app options
	printOptionGroups(tabWriter, parser.Groups())
}

func HelpParser() *flags.Parser {
	helpData := GeneralHelpData{}
	return flags.NewParser(&helpData, flags.HelpFlag)
}
