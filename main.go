package main

////////////////////////////////////////////////////////////////////////////////

import (
	"log"

	"github.com/sabhiram/gomn/coin"

	// Include any coins that we want to manage mns for using gomn
	// we can think of these as "plugins".
	_ "github.com/sabhiram/gomn/coin/pivx"
)

////////////////////////////////////////////////////////////////////////////////

func fatalOnError(err error) {
	if err != nil {
		log.Fatalf("Fatal error encountered: %s\nAborting...\n", err.Error())
	}
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	log.Printf("MAIN\n")

	fatalOnError(coin.BootstrapCoin("pivx"))
	fatalOnError(coin.BootstrapCoin("tx"))

	// prpc, err := pivx.New()
	// if err != nil {
	// 	log.Fatalf("Fatal error: %s\n", err.Error())
	// }

	// println("GET MN status")
	// _, err = prpc.MasternodeStatus()
	// if err != nil {
	// 	log.Printf("Error: %s\n", err.Error())
	// }
	//     println("GET INFO")
	//     prpc.GetInfo()
}

////////////////////////////////////////////////////////////////////////////////

func init() {
	log.SetFlags(0)
	log.SetPrefix("")
}

////////////////////////////////////////////////////////////////////////////////
