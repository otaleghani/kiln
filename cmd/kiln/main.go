// Obsidian vault to static site generator with themes, wikilinks, and graph visualization. @project
// Go 1.25, Cobra CLI, Goldmark, Chroma, Minify, YAML, embed.FS. @stack
// Single binary distribution. No database. Obsidian-compatible markdown only. @constraint
// Functional options pattern. Package-per-domain. Embedded assets via go:embed. @convention
package main

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/otaleghani/kiln/internal/cli"
)

// main is the entry point of the application.
// It initializes the Cobra command hierarchy via the cli package and executes the root command.
func main() {
	// Initialize the Command Line Interface
	rootCmd := cli.Init()

	// Execute the requested command
	// If Execute returns an error (e.g., unknown command, missing flag),
	// we log it and exit with a non-zero status code to indicate failure.
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
