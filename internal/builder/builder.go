package builder

import "log/slog"

// Build orchestrates the static site generation process.
func Build(log *slog.Logger) {
	CleanOutputDir(log)
	switch Mode {
	case "custom":
		log.Info("Building site in Custom mode")
		buildCustom(log)
	default:
		log.Info("Building site in Default mode")
		buildDefault(log)
	}
}

var (
	OutputDir         string // Destination directory
	InputDir          string // Source directory
	FlatUrls          bool   // Defines if the user opted in for flat urls
	ThemeName         string // Theme name
	FontName          string // Font name
	BaseURL           string // Base URL of the application
	SiteName          string // Sitename
	Mode              string // Mode, either default or custom
	LayoutName        string // Layout name
	DisableTOC        bool   // Disables table of contents
	DisableLocalGraph bool   // Disables local graph
)
