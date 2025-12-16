package cli

import (
	"github.com/otaleghani/kiln/internal/builder"
	"github.com/spf13/cobra"
)

var cmdClean = &cobra.Command{
	Use:   "clean",
	Short: "Removes the public output directory",
	Run:   runClean,
}

func init() {
	cmdClean.Flags().
		StringVarP(&outputDir, "output", "o", "", "Name of the output directory (defaults to ./public)")
}

func runClean(cmd *cobra.Command, args []string) {
	if outputDir != "" {
		builder.OutputDir = outputDir
	}

	builder.CleanOutDir()
}
