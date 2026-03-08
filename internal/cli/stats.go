// Cobra stats command that displays vault statistics. @feature:cli
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
	cmdStats.Flags().
		StringVarP(&logger, FlagLog, FlagLogShort, DefaultLog, "Logging level. Choose between 'debug' or 'info'. Defaults to 'info'.")
}

// runStats executes the statistics calculation logic.
func runStats(cmd *cobra.Command, args []string) {
	cfg := loadConfig(cmd)
	applyStringFlag(cmd, FlagInputDir, &inputDir, cfg, DefaultInputDir)
	applyStringFlag(cmd, FlagLog, &logger, cfg, DefaultLog)

	builder.InputDir = inputDir

	log := getLogger()
	builder.Stats(log)
}
