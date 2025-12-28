package cli

import (
	"github.com/otaleghani/kiln/internal/builder"
	"github.com/spf13/cobra"
)

// cmdClean represents the command to remove generated build artifacts.
// It is useful for ensuring a fresh build state or removing old files before a new generation.
var cmdClean = &cobra.Command{
	Use:   "clean",
	Short: "Removes the public output directory",
	Run:   runClean,
}

func init() {
	// Register flags specific to the clean command.
	// We allow the user to specify a custom output directory to clean,
	// in case they generated their site into a non-standard location.
	cmdClean.Flags().
		StringVarP(&outputDir, FlagOutputDir, FlagOutputDirShort, DefaultOutputDir, "Name of the output directory (defaults to ./public)")
}

// runClean executes the cleanup logic.
func runClean(cmd *cobra.Command, args []string) {
	// If a custom output directory was provided via flags, update the builder configuration.
	if outputDir != "" {
		builder.OutputDir = outputDir
	}

	// Delegate the actual deletion logic to the builder package.
	builder.CleanOutDir()
}
