// @feature:cli Cobra dev command: build, watch, and serve for local development.
package cli

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/otaleghani/kiln/internal/builder"
	"github.com/otaleghani/kiln/internal/obsidian"
	"github.com/otaleghani/kiln/internal/server"
	"github.com/otaleghani/kiln/internal/watch"
	"github.com/spf13/cobra"
)

var cmdDev = &cobra.Command{
	Use:   "dev",
	Short: "Build, watch for changes, and serve locally",
	Run:   runDev,
}

func init() {
	cmdDev.Flags().
		StringVarP(&themeName, FlagTheme, FlagThemeShort, DefaultThemeName, "Color theme (default, dracula, catppuccin, nord)")
	cmdDev.Flags().
		StringVarP(&fontName, FlagFont, FlagFontShort, DefaultFontName, "Font family (inter, merriweather, lato, system)")
	cmdDev.Flags().
		StringVarP(&baseURL, FlagUrl, FlagUrlShort, DefaultBaseURL, "Base URL for sitemap generation (e.g. https://example.com)")
	cmdDev.Flags().
		StringVarP(&siteName, FlagSiteName, FlagSiteNameShort, DefaultSiteName, "Name of the website (e.g. 'My Obsidian Vault')")
	cmdDev.Flags().
		StringVarP(&inputDir, FlagInputDir, FlagInputDirShort, DefaultInputDir, "Name of the input directory (defaults to ./vault)")
	cmdDev.Flags().
		StringVarP(&outputDir, FlagOutputDir, FlagOutputDirShort, DefaultOutputDir, "Name of the output directory (defaults to ./public)")
	cmdDev.Flags().
		BoolVar(&flatUrls, FlagFlatURLS, DefaultFlatURLS, "Generate flat HTML files (note.html) instead of pretty directories (note/index.html)")
	cmdDev.Flags().
		StringVarP(&mode, FlagMode, FlagModeShort, DefaultMode, "The mode to use for the generation. Available modes 'default' and 'custom' (defaults to 'default')")
	cmdDev.Flags().
		StringVarP(&logger, FlagLog, FlagLogShort, DefaultLog, "Logging level. Choose between 'debug' or 'info'. Defaults to 'info'.")
	cmdDev.Flags().
		StringVarP(&layout, FlagLayout, FlagLayoutShort, DefaultLayout, "Layout to use. Choose between 'default' and others.")
	cmdDev.Flags().
		BoolVar(&disableTOC, FlagDisableTOC, DefaultDisableTOC, "Disables the Table of contents on the right sidebar.")
	cmdDev.Flags().
		BoolVar(&disableLocalGraph, FlagDisableLocalGraph, DefaultDisableLocalGraph, "Disables the Local graph.")
	cmdDev.Flags().
		BoolVar(&disableBacklinks, FlagDisableBacklinks, DefaultDisableBacklinks, "Disables the Backlinks panel on the right sidebar.")
	cmdDev.Flags().
		StringVarP(&lang, FlagLang, FlagLangShort, DefaultLang, "Language code for the site (e.g. en, it, fr)")
	cmdDev.Flags().
		StringVarP(&accentColor, FlagAccentColor, FlagAccentColorShort, DefaultAccentColor, "Accent color from theme palette (red, orange, yellow, green, blue, purple, cyan)")
	cmdDev.Flags().
		StringVarP(&port, FlagPort, FlagPortShort, DefaultPort, "Port to serve on")
}

func runDev(cmd *cobra.Command, args []string) {
	cfg := loadConfig(cmd)
	applyStringFlag(cmd, FlagTheme, &themeName, cfg, DefaultThemeName)
	applyStringFlag(cmd, FlagFont, &fontName, cfg, DefaultFontName)
	applyStringFlag(cmd, FlagUrl, &baseURL, cfg, DefaultBaseURL)
	applyStringFlag(cmd, FlagSiteName, &siteName, cfg, DefaultSiteName)
	applyStringFlag(cmd, FlagInputDir, &inputDir, cfg, DefaultInputDir)
	applyStringFlag(cmd, FlagOutputDir, &outputDir, cfg, DefaultOutputDir)
	applyStringFlag(cmd, FlagMode, &mode, cfg, DefaultMode)
	applyStringFlag(cmd, FlagLayout, &layout, cfg, DefaultLayout)
	applyStringFlag(cmd, FlagLog, &logger, cfg, DefaultLog)
	applyBoolFlag(cmd, FlagFlatURLS, &flatUrls, cfg, DefaultFlatURLS)
	applyBoolFlag(cmd, FlagDisableTOC, &disableTOC, cfg, DefaultDisableTOC)
	applyBoolFlag(cmd, FlagDisableLocalGraph, &disableLocalGraph, cfg, DefaultDisableLocalGraph)
	applyBoolFlag(cmd, FlagDisableBacklinks, &disableBacklinks, cfg, DefaultDisableBacklinks)
	applyStringFlag(cmd, FlagLang, &lang, cfg, DefaultLang)
	applyStringFlag(cmd, FlagAccentColor, &accentColor, cfg, DefaultAccentColor)
	applyStringFlag(cmd, FlagPort, &port, cfg, DefaultPort)

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
	builder.DisableBacklinks = disableBacklinks
	builder.Lang = lang
	builder.AccentColorName = accentColor

	log := getLogger()

	// Initial full build
	log.Info("Running initial build")
	builder.Build(log)

	// Populate mtime baseline
	mtimeStore := watch.NewMtimeStore()
	if _, _, err := mtimeStore.Update(inputDir); err != nil {
		log.Error("failed to populate mtime baseline", "err", err)
		return
	}

	// Build dependency graph from vault scan
	graph := watch.NewDepGraph()
	vault := obsidian.New(
		obsidian.WithInputDir(inputDir),
		obsidian.WithOutputDir(outputDir),
		obsidian.WithBaseURL(baseURL),
		obsidian.WithFlatURLs(flatUrls),
		obsidian.WithLogger(log),
	)
	if err := vault.Scan(); err != nil {
		log.Error("failed to scan vault for dependency graph", "err", err)
		return
	}
	graph.BuildFromFiles(vault.Vault.Files)

	// Set up watcher with rebuild callback
	watcher := &watch.Watcher{
		InputDir: inputDir,
		Log:      log,
		OnRebuild: func() error {
			changed, removed, err := mtimeStore.Update(inputDir)
			if err != nil {
				log.Error("mtime update failed", "err", err)
				return err
			}

			cs := watch.ComputeChangeSet(changed, removed, graph)
			log.Info("changeset",
				"rebuild", len(cs.Rebuild),
				"remove", len(cs.Remove),
			)

			builder.IncrementalBuild(log, cs.Rebuild, cs.Remove)

			// Refresh dependency graph for changed files
			vault := obsidian.New(
				obsidian.WithInputDir(inputDir),
				obsidian.WithOutputDir(outputDir),
				obsidian.WithBaseURL(baseURL),
				obsidian.WithFlatURLs(flatUrls),
				obsidian.WithLogger(log),
			)
			if err := vault.Scan(); err != nil {
				log.Error("failed to rescan vault", "err", err)
				return nil
			}
			graph.UpdateFiles(vault.Vault.Files)

			return nil
		},
	}

	// Clean shutdown on interrupt
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Start watcher in background
	go func() {
		if err := watcher.Watch(ctx); err != nil {
			log.Error("watcher error", "err", err)
		}
	}()

	// Serve on main goroutine
	localBaseURL := "http://localhost:" + port
	server.Serve(ctx, port, builder.OutputDir, localBaseURL, log)
}
