package cli

import (
	"github.com/otaleghani/kiln/internal/builder"
	"github.com/spf13/cobra"
)

// cmdStats represents the command to analyze and report vault statistics.
// It provides metrics such as total note count, link density, or word counts.
var cmdStats = &cobra.Command{
	Use:   "stats",
	Short: "Displays statistics about your vault",
	Run:   runStats,
}

func init() {
	// Register flags for the stats command.
	// Allows running statistics on a custom vault location.
	cmdStats.Flags().
		StringVarP(&inputDir, FlagInputDir, FlagInputDirShort, DefaultInputDir, "Name of the input directory (defaults to ./vault)")
}

// runStats executes the statistics calculation logic.
func runStats(cmd *cobra.Command, args []string) {
	// If a custom input directory is provided via flags, update the builder configuration.
	if inputDir != "" {
		builder.InputDir = inputDir
	}

	// Trigger the stats analysis and print the results to the console.
	builder.Stats()
}
