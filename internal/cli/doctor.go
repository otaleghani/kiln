package cli

import (
	"log"

	"github.com/otaleghani/kiln/internal/builder"
	"github.com/otaleghani/kiln/internal/linter"
	"github.com/spf13/cobra"
)

var cmdDoctor = &cobra.Command{
	Use:   "doctor",
	Short: "Checks for broken links",
	Run:   runDoctor,
}

func init() {
	cmdDoctor.Flags().
		StringVarP(&inputDir, "input", "i", "", "Name of the input directory (defaults to ./vault)")
}

func runDoctor(cmd *cobra.Command, args []string) {
	log.Println("Diagnosing vault...")

	if inputDir != "" {
		builder.InputDir = inputDir
	}

	notes := linter.CollectNotes(builder.OutputDir)
	linter.BrokenLinks(builder.OutputDir, notes)
}
