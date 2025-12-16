package cli

import (
	"github.com/otaleghani/kiln/internal/builder"
	"github.com/otaleghani/kiln/internal/server"
	"github.com/spf13/cobra"
)

var cmdServe = &cobra.Command{
	Use:   "serve",
	Short: "Serves the generated site locally",
	Run:   runServe,
}

var port string

func init() {
	cmdServe.Flags().StringVarP(&port, "port", "p", "8080", "Port to serve on")
	cmdServe.Flags().
		StringVarP(&outputDir, "output", "o", "", "Name of the output directory to serve(defaults to ./public)")
}

func runServe(cmd *cobra.Command, args []string) {
	if outputDir != "" {
		builder.OutputDir = outputDir
	}

	server.Serve(port, builder.OutputDir)
}
