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
		StringVarP(&outputDir, "output", "o", "", "Name of the output directory (defaults to ./public)")
}

func runDoctor(cmd *cobra.Command, args []string) {
	log.Println("Diagnosing vault...")

	if outputDir != "" {
		builder.OutputDir = outputDir
	}

	notes := linter.CollectNotes(builder.OutputDir)
	linter.BrokenLinks(builder.OutputDir, notes)
}
