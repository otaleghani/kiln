package cli

import (
	"github.com/otaleghani/kiln/internal/builder"
	"github.com/otaleghani/kiln/internal/server"
	"github.com/spf13/cobra"
)

// cmdServe represents the command to start a local development server.
// It allows users to preview their generated static site before deploying it.
var cmdServe = &cobra.Command{
	Use:   "serve",
	Short: "Serves the generated site locally",
	Run:   runServe,
}

// port stores the port number specified by the user (default: 8080).
var port string

func init() {
	// Register flags for the serve command.
	// Users can customize the listening port and the directory being served.
	cmdServe.Flags().StringVarP(&port, FlagPort, FlagPortShort, DefaultPort, "Port to serve on")
	cmdServe.Flags().
		StringVarP(&outputDir, FlagOutputDir, FlagOutputDirShort, DefaultOutputDir, "Name of the output directory to serve(defaults to ./public)")
}

// runServe executes the server logic.
func runServe(cmd *cobra.Command, args []string) {
	// If the user specified a custom output directory via flags, update the builder config.
	builder.OutputDir = outputDir

	// Construct the local base URL (e.g., http://localhost:8080).
	// This helps ensure absolute links or assets resolve correctly during local preview.
	localBaseURL := "http://localhost:" + port

	// Start the static file server.
	server.Serve(port, builder.OutputDir, localBaseURL)
}
