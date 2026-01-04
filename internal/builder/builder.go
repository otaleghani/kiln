package builder

import "github.com/otaleghani/kiln/internal/log"

// Build orchestrates the static site generation process.
func Build() {
	CleanOutputDir()
	switch Mode {
	case "custom":
		log.Info("Building site in Custom mode")
		buildCustom()
	default:
		log.Info("Building site in Default mode")
		buildDefault()
	}
}

var (
	OutputDir  string // Destination directory
	InputDir   string // Source directory
	FlatUrls   bool   // Defines if the user opted in for flat urls
	ThemeName  string // Theme name
	FontName   string // Font name
	BaseURL    string // Base URL of the application
	SiteName   string // Sitename
	Mode       string // Mode, either default or custom
	LayoutName string // Layout name
)

// GraphNode represents a single node in the interactive graph view.
type GraphNode struct {
	ID    string `json:"id"`
	Label string `json:"label"`
	URL   string `json:"url"`
	Val   int    `json:"val"`
}
