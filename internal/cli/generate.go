// Cobra generate command that triggers static site building. @feature:cli
package cli

import (
	"log/slog"

	"github.com/otaleghani/kiln/internal/builder"
	"github.com/otaleghani/kiln/internal/config"
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
		StringVarP(&themeName, FlagTheme, FlagThemeShort, DefaultThemeName, "Color theme (default, dracula, catppuccin, nord)")
	cmdGenerate.Flags().
		StringVarP(&fontName, FlagFont, FlagFontShort, DefaultFontName, "Font family (inter, merriweather, lato, system)")
	cmdGenerate.Flags().
		StringVarP(&baseURL, FlagUrl, FlagUrlShort, DefaultBaseURL, "Base URL for sitemap generation (e.g. https://example.com)")
	cmdGenerate.Flags().
		StringVarP(&siteName, FlagSiteName, FlagSiteNameShort, DefaultSiteName, "Name of the website (e.g. 'My Obsidian Vault')")
	cmdGenerate.Flags().
		StringVarP(&inputDir, FlagInputDir, FlagInputDirShort, DefaultInputDir, "Name of the input directory (defaults to ./vault)")
	cmdGenerate.Flags().
		StringVarP(&outputDir, FlagOutputDir, FlagOutputDirShort, DefaultOutputDir, "Name of the output directory (defaults to ./public)")
	cmdGenerate.Flags().
		BoolVar(&flatUrls, FlagFlatURLS, DefaultFlatURLS, "Generate flat HTML files (note.html) instead of pretty directories (note/index.html)")
	cmdGenerate.Flags().
		StringVarP(&mode, FlagMode, FlagModeShort, DefaultMode, "The mode to use for the generation. Available modes 'default' and 'custom' (defaults to 'default')")
	cmdGenerate.Flags().
		StringVarP(&logger, FlagLog, FlagLogShort, DefaultLog, "Logging level. Choose between 'debug' or 'info'. Defaults to 'info'.")
	cmdGenerate.Flags().
		StringVarP(&layout, FlagLayout, FlagLayoutShort, DefaultLayout, "Layout to use. Choose between 'default' and others.")
	cmdGenerate.Flags().
		BoolVar(&disableTOC, FlagDisableTOC, DefaultDisableTOC, "Disables the Table of contents on the right sidebar. If the local graph is disabled too, hides the right sidebar.")
	cmdGenerate.Flags().
		BoolVar(&disableLocalGraph, FlagDisableLocalGraph, DefaultDisableLocalGraph, "Disables the Local graph. If the table of contents is disabled too, hides the right sidebar.")
}

// runGenerate executes the build logic.
func runGenerate(cmd *cobra.Command, args []string) {
	// Load config file if present; values act as defaults that CLI flags override.
	applyConfig(cmd)

	// Apply overrides
	builder.OutputDir = outputDir
	builder.InputDir = inputDir
	builder.FlatUrls = flatUrls
	builder.ThemeName = themeName
	builder.FontName = fontName
	builder.BaseURL = baseURL
	builder.SiteName = siteName
	builder.Mode = mode
	builder.LayoutName = layout
	builder.DisableTOC = disableTOC
	builder.DisableLocalGraph = disableLocalGraph

	log := getLogger()

	// Trigger the build
	builder.Build(log)
}

// applyConfig discovers and loads kiln.yaml, applying its values as defaults
// for any flags the user did not explicitly set on the command line.
func applyConfig(cmd *cobra.Command) {
	path, err := config.FindFile(".")
	if err != nil {
		slog.Warn("Error searching for config file", "error", err)
		return
	}
	if path == "" {
		return
	}

	cfg, err := config.Load(path)
	if err != nil {
		slog.Warn("Failed to parse config file", "file", path, "error", err)
		return
	}
	if cfg == nil {
		return
	}

	slog.Info("Loaded configuration", "file", path)

	if !cmd.Flags().Changed(FlagTheme) {
		themeName = cfg.ValueOr("theme", DefaultThemeName)
	}
	if !cmd.Flags().Changed(FlagFont) {
		fontName = cfg.ValueOr("font", DefaultFontName)
	}
	if !cmd.Flags().Changed(FlagUrl) {
		baseURL = cfg.ValueOr("url", DefaultBaseURL)
	}
	if !cmd.Flags().Changed(FlagSiteName) {
		siteName = cfg.ValueOr("name", DefaultSiteName)
	}
	if !cmd.Flags().Changed(FlagInputDir) {
		inputDir = cfg.ValueOr("input", DefaultInputDir)
	}
	if !cmd.Flags().Changed(FlagOutputDir) {
		outputDir = cfg.ValueOr("output", DefaultOutputDir)
	}
	if !cmd.Flags().Changed(FlagMode) {
		mode = cfg.ValueOr("mode", DefaultMode)
	}
	if !cmd.Flags().Changed(FlagLayout) {
		layout = cfg.ValueOr("layout", DefaultLayout)
	}
	if !cmd.Flags().Changed(FlagLog) {
		logger = cfg.ValueOr("log", DefaultLog)
	}
	if !cmd.Flags().Changed(FlagFlatURLS) {
		flatUrls = cfg.BoolOr("flat-urls", DefaultFlatURLS)
	}
	if !cmd.Flags().Changed(FlagDisableTOC) {
		disableTOC = cfg.BoolOr("disable-toc", DefaultDisableTOC)
	}
	if !cmd.Flags().Changed(FlagDisableLocalGraph) {
		disableLocalGraph = cfg.BoolOr("disable-local-graph", DefaultDisableLocalGraph)
	}
}
