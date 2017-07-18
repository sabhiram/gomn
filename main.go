package main

////////////////////////////////////////////////////////////////////////////////

import (
	"flag"
	"log"
	"strings"

	"github.com/sabhiram/gomn/coin"
	"github.com/sabhiram/gomn/monitor"
	"github.com/sabhiram/gomn/types"

	// Include any coins that we want to manage mns for using gomn
	// we can think of these as "plugins".
	_ "github.com/sabhiram/gomn/coin/pivx"
)

////////////////////////////////////////////////////////////////////////////////

var (
	cli      = &types.CLI{}
	GoMnHelp = `
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
    --wallet  Specify where the wallet binaries will be fetched
    --bins    Specify the subpath within the wallet where the bins exist

Not all 'COMMAND's require the above options to be set, however most that query
or setup a node for a given coin will require them.

COMMANDS:
=========

  gomn status:
  ------------
    help         Print this help menu, same as running 'gomn' with no command
    version      Print the application's version information
    list         List coins that gomn is aware of

  coin specific commands:
  -----------------------
    info         Get info on a given coin, effectively run the 'getinfo' RPC.

    download     Get a copy of the coin's wallet to the bin path. To override
                 the coin specified defaults, use '--url' to specify a source
                 url to fetch the wallet from and '--type' to specify the type
                 of compression (if any).  If you have a shasum to verify the
                 download against, specify that with '--shasum'.

    bootstrap    Fetch the bootstrap bundle (if available) to the data path. To
                 override the coin specified defaults, use '--url' to specify a
                 source url to fetch the bootstrap from, and use '--type' to
                 specify the type of compression (if any).

    configure    Configure the 'coin'.conf file for mn duty.  You must specify

    monitor      Once all other things are setup, this will monitor your MN.
                 If '--callbackurl' is specified, updates are sent to the URL
                 as the node's state changes.  If '--start' is specified, this
                 will kick off the node's specified daemon.  If '--start' is not
                 specified and the server is not running, this will abort.

`
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
	if len(cli.Args) > 0 {
		cmd = strings.ToLower(cli.Args[0])
		opts = cli.Args[1:]
	}

	switch cmd {
	case "help":
		log.Printf("%s\n", GoMnHelp)
	case "version":
		log.Printf("%s\n", Version)
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
	case "monitor":
		m, err := monitor.New(cli, opts)
		fatalOnError(err)
		m.Start() // Run indefinitely
	default:
		fatalOnError(coin.Command(cli, cmd, opts))
	}
}

////////////////////////////////////////////////////////////////////////////////

func init() {
	// Setup logger for this application.
	log.SetFlags(0)
	log.SetPrefix("")

	// CLI Argument parsing.
	flag.StringVar(&cli.Coin, "coin", "", "currency you are setting up a mn for")
	flag.StringVar(&cli.Wallet, "wallet", "", "base path to where wallet binaries will be extracted (optional)")
	flag.StringVar(&cli.BinPath, "bins", "", "path where the coin's binaries should reside (optional)")
	flag.StringVar(&cli.DataPath, "data", "", "path where the blockchain data should reside (optional)")
	flag.Parse()

	// Normalize and fix-up arguments.
	cli.Args = flag.Args()
	cli.Coin = strings.ToLower(cli.Coin)
}

////////////////////////////////////////////////////////////////////////////////
