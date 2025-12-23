package builder

import "log"

var (
	// OutputDir is the destination directory for the generated site.
	OutputDir string = "./public"
	// InputDir is the source directory containing the Obsidian vault.
	InputDir string = "./vault"
	// FlatUrls defines if the user opted in for flat urls.
	FlatUrls  bool
	ThemeName string
	FontName  string
	BaseURL   string
	SiteName  string
	Mode      string
)

// GraphNode represents a single node in the interactive graph view.
type GraphNode struct {
	ID    string `json:"id"`
	Label string `json:"label"`
	URL   string `json:"url"`
	Val   int    `json:"val"`
}

// Build orchestrates the static site generation process.
func Build() {
	switch Mode {
	case "custom":
		log.Println("Building site in Custom mode")
		buildCustom()
	default:
		log.Println("Building site in Default mode")
		buildDefault()
	}
}
