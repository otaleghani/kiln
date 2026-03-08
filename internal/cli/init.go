// Cobra init command that scaffolds a new vault directory. @feature:cli
package cli

import (
	"os"

	"github.com/otaleghani/kiln/internal/builder"
	"github.com/otaleghani/kiln/internal/config"
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

	// Scaffold a kiln.yaml in the current directory if one doesn't exist.
	if _, err := os.Stat(config.DefaultFilename); os.IsNotExist(err) {
		content := `# Kiln configuration file
# Uncomment and edit the options below to set defaults for your project.
# CLI flags will override these values.

# theme: default
# font: inter
# url: ""
# name: "My Notes"
# input: ./vault
# output: ./public
# mode: default
# layout: default
# flat-urls: false
# disable-toc: false
# disable-local-graph: false
`
		if err := os.WriteFile(config.DefaultFilename, []byte(content), 0o644); err != nil {
			log.Error("Couldn't create config file", "error", err)
			return
		}
		log.Info("Created kiln.yaml configuration file")
	}
}
