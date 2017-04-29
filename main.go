package main

import (
	// "flags"
	"log"

	pivx "github.com/sabhiram/mnwatch/rpc/pivx"
)

func main() {
	log.Printf("MAIN\n")

	prpc, err := pivx.New()
	if err != nil {
		log.Fatalf("Fatal error: %s\n", err.Error())
	}

	println("GET MN status")
	_, err = prpc.MasternodeStatus()
	if err != nil {
		log.Printf("Error: %s\n", err.Error())
	}
	//     println("GET INFO")
	//     prpc.GetInfo()
}

func init() {
	log.SetFlags(0)
	log.SetPrefix("")
}
