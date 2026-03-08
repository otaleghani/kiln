// Cobra init command that scaffolds a new vault directory. @feature:cli
package cli

import (
	"github.com/otaleghani/kiln/internal/builder"
	"github.com/spf13/cobra"
)

// cmdInit represents the command to scaffold a new project structure.
// It checks for the existence of a vault directory and creates a welcome note if one isn't found.
var cmdInit = &cobra.Command{
	Use:   "init",
	Short: "Initializes a new Kiln project",
	Run:   runInit,
}

func init() {
	// Register flags for the init command.
	// Allows the user to specify a custom directory name for their new vault.
	cmdInit.Flags().
		StringVarP(&inputDir, FlagInputDir, FlagInputDirShort, DefaultInputDir, "Name of the input directory (defaults to ./vault)")
	cmdInit.Flags().
		StringVarP(&logger, FlagLog, FlagLogShort, DefaultLog, "Logging level. Choose between 'debug' or 'info'. Defaults to 'info'.")
}

// runInit executes the initialization logic.
func runInit(cmd *cobra.Command, args []string) {
	cfg := loadConfig(cmd)
	applyStringFlag(cmd, FlagInputDir, &inputDir, cfg, DefaultInputDir)
	applyStringFlag(cmd, FlagLog, &logger, cfg, DefaultLog)

	builder.InputDir = inputDir

	log := getLogger()
	builder.Init(log)
}
