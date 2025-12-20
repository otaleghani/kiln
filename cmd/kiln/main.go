package main

import (
	"github.com/otaleghani/kiln/internal/cli"
	"log"
	"os"
)

func init() {
	// log.SetFlags(0) // Remove default timestamps (we can add our own or keep it clean)
	log.SetPrefix("kiln: ")
}

func main() {
	rootCmd := cli.Init()

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
