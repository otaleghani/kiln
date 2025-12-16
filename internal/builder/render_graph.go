package builder

import (
	"html/template"
	"log"
	"os"
	"path/filepath"

	"github.com/tdewolff/minify/v2"
)

type GraphRenderer struct {
	rootNode *Node              // The root node of the vault
	minifier *minify.M          // Minifier
	siteName string             // The name of the website
	template *template.Template // Golang html template
	baseURL  string             // The base URL for links
}

func (g *GraphRenderer) render() {
	graphOutPath := filepath.Join(OutputDir, "graph", "index.html")
	os.MkdirAll(filepath.Dir(graphOutPath), 0755)

	// Set active for graph
	setTreeActive(g.rootNode.Children, "/graph")

	fGraph, err := os.Create(graphOutPath)
	if err != nil {
		log.Printf("Failed to create graph view. Error: %s", err.Error())
	}
	defer fGraph.Close()
	mwGraph := g.minifier.Writer("text/html", fGraph)
	defer mwGraph.Close()

	// Use a special div ID for the global graph
	graphHTML := `<div id="global-graph-container" style=""></div>`

	dataGraph := PageData{
		Title:       "Graph View",
		BaseURL:     g.baseURL,
		SiteName:    g.siteName,
		Content:     template.HTML(graphHTML),
		Breadcrumbs: []string{"Home", "Graph"},
		IsCanvas:    false,
		IsGraph:     true,
		Sidebar:     g.rootNode.Children,
	}
	g.template.Execute(mwGraph, dataGraph)
}
