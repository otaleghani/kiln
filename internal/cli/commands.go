package cli

import (
	"github.com/spf13/cobra"
)

const (
	OutputDir  = "./public"
	InputVault = "./vault"
)

func Init() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "kiln",
		Short: "A lightweight Obsidian static site generator",
	}

	rootCmd.AddCommand(cmdGenerate)
	rootCmd.AddCommand(cmdServe)
	rootCmd.AddCommand(cmdInit)
	rootCmd.AddCommand(cmdClean)
	rootCmd.AddCommand(cmdDoctor)
	rootCmd.AddCommand(cmdStats)

	return rootCmd
}
