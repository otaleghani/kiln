// Package cli provides the command-line interface entry points for the Kiln static site generator.
// It orchestrates the various subcommands (generate, serve, etc.) using the Cobra library.
package cli

import (
	"github.com/spf13/cobra"
)

// Default configuration constants for the build process.
const (
	OutputDir  = "./public" // The target directory for the generated static site
	InputVault = "./vault"  // The default source directory containing Obsidian markdown files
)

// Global variables to store the values of command-line flags.
// These are populated by Cobra when the command is executed.
var (
	themeName string // The visual theme (e.g., "dracula")
	fontName  string // The font family to use (e.g., "inter")
	baseURL   string // The base URL for SEO and sitemap generation
	siteName  string // The display name of the generated site
	inputDir  string // Custom path to the source vault
	outputDir string // Custom path for the build output
	flatUrls  bool   // Choose between pretty (e.g. note/index.html) or flat URLs (e.g. note.html)
)

// Init constructs and returns the root command for the application.
// It registers all available subcommands, establishing the CLI hierarchy.
func Init() *cobra.Command {
	// rootCmd represents the base command when called without any subcommands
	rootCmd := &cobra.Command{
		Use:   "kiln",
		Short: "A lightweight Obsidian static site generator",
		Long: `Kiln is a tool that turns your Obsidian vault into a fast, static website.
It supports wikilinks, callouts, mermaid diagrams, and graph visualization.`,
	}

	// Register subcommands to the root
	rootCmd.AddCommand(cmdGenerate) // Builds the static site
	rootCmd.AddCommand(cmdServe)    // Starts a local preview server
	rootCmd.AddCommand(cmdInit)     // Initializes a new vault structure
	rootCmd.AddCommand(cmdClean)    // Removes generated artifacts
	rootCmd.AddCommand(cmdDoctor)   // Checks for common issues
	rootCmd.AddCommand(cmdStats)    // Displays vault statistics

	return rootCmd
}
