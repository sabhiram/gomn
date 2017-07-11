package main

////////////////////////////////////////////////////////////////////////////////

import (
	"flag"
	"log"
	"strings"

	"github.com/sabhiram/gomn/coin"
	"github.com/sabhiram/gomn/version"

	// Include any coins that we want to manage mns for using gomn
	// we can think of these as "plugins".
	_ "github.com/sabhiram/gomn/coin/pivx"
)

////////////////////////////////////////////////////////////////////////////////

var (
	GoMnVersion = version.VersionString
	GoMnHelp    = `
GoMn Usage:
===========

    $ gomn [OPTIONS] COMMAND [COMMAND OPTIONS]

Running this tool with no arguments is effectively the same as invoking
'gomn help'.

Most 'COMMAND's do not have options, but if they do, they are specified after
the command itself.

OPTIONS:
========

    --coin    Specify the coin [Ex: "pivx", "dash", ...]
    --data    Specify the data path for the coin, if empty use coin default
    --bins    Specify the binary path for the coin, if empty use coin default

Not all 'COMMAND's require the above options to be set, however most that query
or setup a node for a given coin will require them.

COMMANDS:
=========

    help         Print this help menu, same as running 'gomn' with no command
    version      Print the application's version information
    list         List coins that gomn is aware of

    info         Get info on a given coin, effectively run the 'getinfo' RPC
    download     Get a copy of the coin's wallet to the bin path
    bootstrap    Fetch the bootstrap bundle (if available) to the data path
    configure    Configure the 'coin'.conf file for mn duty
    ...          Other commands :)
`

	CLI = struct {
		coin     string   // Name of the coin we are operating on
		binPath  string   // Path to coin's binary directory (empty => coin default)
		dataPath string   // Path to coin's data directory (empty => coin default)
		args     []string // Rest of the command line, args[0] is the command
	}{}
)

////////////////////////////////////////////////////////////////////////////////

func fatalOnError(err error) {
	if err != nil {
		log.Fatalf("Fatal error encountered: %s\nAborting...\n", err.Error())
	}
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	cmd := "help"
	opts := []string{}
	if len(CLI.args) > 0 {
		cmd = strings.ToLower(CLI.args[0])
		opts = CLI.args[1:]
	}

	switch cmd {
	case "help":
		log.Printf("%s\n", GoMnHelp)
	case "version":
		log.Printf("%s\n", GoMnVersion)
	case "list":
		cs := coin.RegisteredCoins()
		if len(cs) > 0 {
			log.Printf("Registered coins:\n")
			for k, v := range cs {
				log.Printf("%d. %s\n", k+1, v)
			}
			log.Printf("\n")
		} else {
			log.Printf("No coins registered!\n\n")
		}
	default:
		if err := coin.Command(CLI.coin, CLI.binPath, CLI.dataPath, cmd, opts); err != nil {
			fatalOnError(err)
		}
	}
}

////////////////////////////////////////////////////////////////////////////////

func init() {
	// Setup logger for this application.
	log.SetFlags(0)
	log.SetPrefix("")

	// CLI Argument parsing.
	flag.StringVar(&CLI.coin, "coin", "", "currency you are setting up a mn for")
	flag.StringVar(&CLI.binPath, "bins", "", "path where the coin's binaries should reside (optional)")
	flag.StringVar(&CLI.dataPath, "data", "", "path where the blockchain data should reside (optional)")
	flag.Parse()

	// Normalize and fix-up arguments.
	CLI.args = flag.Args()
	CLI.coin = strings.ToLower(CLI.coin)
}

////////////////////////////////////////////////////////////////////////////////
