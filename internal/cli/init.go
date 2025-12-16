package cli

import (
	"github.com/otaleghani/kiln/internal/builder"
	"github.com/spf13/cobra"
)

var cmdInit = &cobra.Command{
	Use:   "init",
	Short: "Initializes a new Kiln project",
	Run:   runInit,
}

func init() {
	cmdInit.Flags().
		StringVarP(&inputDir, "input", "i", "", "Name of the input directory (defaults to ./vault)")
}

func runInit(cmd *cobra.Command, args []string) {
	if inputDir != "" {
		builder.InputDir = inputDir
	}

	builder.Init()
}
