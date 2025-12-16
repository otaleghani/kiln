package cli

import (
	"github.com/otaleghani/kiln/internal/builder"
	"github.com/spf13/cobra"
)

var cmdStats = &cobra.Command{
	Use:   "stats",
	Short: "Displays statistics about your vault",
	Run:   runStats,
}

func init() {
	cmdStats.Flags().
		StringVarP(&inputDir, "input", "i", "", "Name of the input directory (defaults to ./vault)")
}

func runStats(cmd *cobra.Command, args []string) {
	if inputDir != "" {
		builder.InputDir = inputDir
	}

	builder.Stats()
}
