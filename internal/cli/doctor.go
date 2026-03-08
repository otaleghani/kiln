// Cobra doctor command that scans the vault for broken links. @feature:cli
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
	cfg := loadConfig(cmd)
	applyStringFlag(cmd, FlagInputDir, &inputDir, cfg, DefaultInputDir)
	applyStringFlag(cmd, FlagLog, &logger, cfg, DefaultLog)

	log.Println("Diagnosing vault...")

	builder.InputDir = inputDir

	slogger := getLogger()

	notes := linter.CollectNotes(builder.InputDir)
	linter.BrokenLinks(builder.InputDir, notes, slogger)
}
