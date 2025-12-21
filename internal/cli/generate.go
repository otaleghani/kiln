package cli

import (
	"github.com/otaleghani/kiln/internal/builder"
	"github.com/spf13/cobra"
)

// cmdGenerate represents the primary build command.
// It triggers the static site generation process, converting Markdown to HTML.
var cmdGenerate = &cobra.Command{
	Use:   "generate",
	Short: "Builds the static site from your vault",
	Run:   runGenerate,
}

func init() {
	// Register flags to allow users to customize the build without changing code.
	// We support short flags (e.g., -t) and long flags (e.g., --theme).
	cmdGenerate.Flags().
		StringVarP(&themeName, "theme", "t", "default", "Color theme (default, dracula, catppuccin, nord)")
	cmdGenerate.Flags().
		StringVarP(&fontName, "font", "f", "inter", "Font family (inter, merriweather, lato, system)")
	cmdGenerate.Flags().
		StringVarP(&baseURL, "url", "u", "", "Base URL for sitemap generation (e.g. https://example.com)")
	cmdGenerate.Flags().
		StringVarP(&siteName, "name", "n", "My Notes", "Name of the website (e.g. 'My Obsidian Vault')")
	cmdGenerate.Flags().
		StringVarP(&inputDir, "input", "i", "", "Name of the input directory (defaults to ./vault)")
	cmdGenerate.Flags().
		StringVarP(&outputDir, "output", "o", "", "Name of the output directory (defaults to ./public)")
	cmdGenerate.Flags().
		BoolVar(&flatUrls, "flat-urls", false, "Generate flat HTML files (note.html) instead of pretty directories (note/index.html)")
}

// runGenerate executes the build logic.
func runGenerate(cmd *cobra.Command, args []string) {
	// Apply overrides
	// If the user specified custom directories via flags, update the builder configuration.
	if outputDir != "" {
		builder.OutputDir = outputDir
	}
	if inputDir != "" {
		builder.InputDir = inputDir
	}
	builder.FlatUrls = flatUrls

	// Trigger the build
	// Pass the cosmetic and metadata configurations to the builder.
	builder.Build(themeName, fontName, baseURL, siteName)
}
