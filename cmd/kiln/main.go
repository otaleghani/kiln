package main

import (
	"log"
	"os"

	"github.com/otaleghani/kiln/internal/cli"
)

// init runs before main() and sets up global application defaults.
// Here, it configures the standard logger to prefix all output with "kiln: ",
// making it easier to distinguish Kiln's logs from other system output.
func init() {
	log.SetPrefix("kiln: ")
}

// main is the entry point of the application.
// It initializes the Cobra command hierarchy via the cli package and executes the root command.
func main() {
	// 1. Initialize the Command Line Interface
	rootCmd := cli.Init()

	// 2. Execute the requested command
	// If Execute returns an error (e.g., unknown command, missing flag),
	// we log it and exit with a non-zero status code to indicate failure.
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
