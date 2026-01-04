package cli

import (
	"log"

	"github.com/otaleghani/kiln/internal/builder"
	"github.com/otaleghani/kiln/internal/linter"
	"github.com/spf13/cobra"
)

// cmdDoctor represents the diagnostic command.
// It scans the vault to identify issues such as broken wiki-links or missing references.
var cmdDoctor = &cobra.Command{
	Use:   "doctor",
	Short: "Checks for broken links",
	Run:   runDoctor,
}

func init() {
	// Register flags for the doctor command.
	// Allows running diagnostics on a custom vault location.
	cmdDoctor.Flags().
		StringVarP(&inputDir, FlagInputDir, FlagInputDirShort, DefaultInputDir, "Name of the input directory (defaults to ./vault)")
	cmdDoctor.Flags().
		StringVarP(&logger, FlagLog, FlagLogShort, DefaultLog, "Logging level. Choose between 'debug' or 'info'. Defaults to 'info'.")
}

// runDoctor executes the linting logic.
func runDoctor(cmd *cobra.Command, args []string) {
	log.Println("Diagnosing vault...")

	// Override the default input directory if the flag is set.
	builder.InputDir = inputDir

	setLogger()

	// Collect all valid note paths to build a reference index.
	notes := linter.CollectNotes(builder.InputDir)

	// Scan content for links that point to non-existent notes.
	linter.BrokenLinks(builder.InputDir, notes)
}
