package cli

import (
	"github.com/otaleghani/kiln/internal/builder"
	"github.com/spf13/cobra"
)

var themeName string
var fontName string
var baseUrl string
var siteName string
var inputDir string
var outputDir string

var cmdGenerate = &cobra.Command{
	Use:   "generate",
	Short: "Builds the static site from your vault",
	Run:   runGenerate,
}

func init() {
	cmdGenerate.Flags().
		StringVarP(&themeName, "theme", "t", "default", "Color theme (default, dracula, catppuccin, nord)")
	cmdGenerate.Flags().
		StringVarP(&fontName, "font", "f", "inter", "Font family (inter, merriweather, lato, system)")
	cmdGenerate.Flags().
		StringVarP(&baseUrl, "url", "u", "", "Base URL for sitemap generation (e.g. https://example.com)")
	cmdGenerate.Flags().
		StringVarP(&siteName, "name", "n", "My Notes", "Name of the website (e.g. 'My Obsidian Vault')")
	cmdGenerate.Flags().
		StringVarP(&inputDir, "input", "i", "", "Name of the input directory (defaults to ./vault)")
	cmdGenerate.Flags().
		StringVarP(&outputDir, "output", "o", "", "Name of the output directory (defaults to ./public)")
}

func runGenerate(cmd *cobra.Command, args []string) {
	if outputDir != "" {
		builder.OutputDir = outputDir
	}
	if inputDir != "" {
		builder.InputDir = inputDir
	}

	builder.Build(themeName, fontName, baseUrl, siteName)
}
