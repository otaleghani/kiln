package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var cmdVersion = &cobra.Command{
	Use:   "version",
	Short: "Displays kiln version",
	Run:   runVersion,
}

var version = "dev"

// runStats executes the statistics calculation logic.
func runVersion(cmd *cobra.Command, args []string) {
	fmt.Println("kiln", version)
}
