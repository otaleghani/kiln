package builder

import (
	"html/template"
	"log"
	"os"
	"path/filepath"

	"github.com/tdewolff/minify/v2"
)

// GraphRenderer generates the dedicated full-screen interactive graph page.
type GraphRenderer struct {
	rootNode *Node              // Navigation tree root
	minifier *minify.M          // HTML minifier
	siteName string             // Global site name
	template *template.Template // Main HTML layout
	baseURL  string             // Site Base URL
	theme    Theme              // Visual theme settings
}

// render creates the /graph/index.html page.
// It sets up the specific container ID that the frontend JS looks for to initialize the D3 graph.
func (g *GraphRenderer) render() {
	graphOutPath := filepath.Join(OutputDir, "graph", "index.html")
	os.MkdirAll(filepath.Dir(graphOutPath), 0755)

	// Highlight the "Graph" item in the sidebar navigation
	setTreeActive(g.rootNode.Children, "/graph")

	fGraph, err := os.Create(graphOutPath)
	if err != nil {
		log.Printf("Failed to create graph view. Error: %s", err.Error())
		return
	}
	defer fGraph.Close()

	// Wrap writer with minifier
	mwGraph := g.minifier.Writer("text/html", fGraph)
	defer mwGraph.Close()

	// The frontend script targets this ID to mount the global graph
	graphHTML := `<div id="global-graph-container" style=""></div>`

	dataGraph := PageData{
		Title:       "Graph View",
		BaseURL:     g.baseURL,
		SiteName:    g.siteName,
		Content:     template.HTML(graphHTML),
		Breadcrumbs: []string{"Home", "Graph"},
		IsCanvas:    false,
		IsGraph:     true, // Flags the template to load graph-specific scripts
		Sidebar:     g.rootNode.Children,
		Font:        g.theme.Font,
	}

	g.template.Execute(mwGraph, dataGraph)
}
