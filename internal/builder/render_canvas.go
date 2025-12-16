package builder

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/tdewolff/minify/v2"
	"github.com/yuin/goldmark"
)

type CanvasNoteRenderer struct {
	path           string             // Path of the file to render
	relPath        string             // Relative path of the file
	siteName       string             // The name of the website
	nameWithoutExt string             // Name of the file without extention
	template       *template.Template // Golang html template
	renderer       goldmark.Markdown  // Markdown renderer
	minifier       *minify.M          // Minifier
	rootNode       *Node              // The root node of the vault
	baseURL        string             // The base URL for links
}

type CanvasNode struct {
	ID     string `json:"id"`
	Type   string `json:"type"`
	Text   string `json:"text,omitempty"`
	File   string `json:"file,omitempty"`
	Label  string `json:"label,omitempty"`
	X      int    `json:"x"`
	Y      int    `json:"y"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Color  string `json:"color,omitempty"`
	// Injected by render function
	HtmlContent string `json:"htmlContent,omitempty"`
	IsImage     bool   `json:"isImage,omitempty"`
	Src         string `json:"src,omitempty"`
}

type CanvasData struct {
	Nodes []CanvasNode     `json:"nodes"`
	Edges []map[string]any `json:"edges"`
}

func (c *CanvasNoteRenderer) render() string {
	source, err := os.ReadFile(c.path)
	if err != nil {
		log.Printf("Failed to create canvas note. Error: %s", err.Error())
	}

	// Parse the Canvas JSON
	var canvasData CanvasData
	if err := json.Unmarshal(source, &canvasData); err == nil {
		// Iterate over nodes to hydrate them
		for i, node := range canvasData.Nodes {
			if node.Type == "file" && node.File != "" {
				// Resolve full path to the linked file
				linkedFilePath := filepath.Join(InputDir, node.File)
				linkedExt := strings.ToLower(filepath.Ext(linkedFilePath))

				// Check if file exists
				if _, err := os.Stat(linkedFilePath); err == nil {
					// Image handling
					if isImageExt(linkedExt) {
						canvasData.Nodes[i].IsImage = true
						// Generate the public web path for the image
						// Since the walker copies assets to outputDir structure, we replicate that path logic
						relImgPath, _ := filepath.Rel(InputDir, linkedFilePath)
						parts := strings.Split(relImgPath, string(os.PathSeparator))
						for k, p := range parts {
							parts[k] = slugify(p)
						}
						// Reconstruct as web path
						canvasData.Nodes[i].Src = "/" + strings.Join(parts, "/")
					} else if linkedExt == ".md" {
						// Markdown handling
						noteContent, err := os.ReadFile(linkedFilePath)
						if err == nil {
							// Parse Markdown to HTML
							var buf bytes.Buffer
							// Use existing MD parser context
							if err := c.renderer.Convert(noteContent, &buf); err == nil {
								htmlContent := buf.String()
								// Apply Post-processing (Callouts, etc)
								htmlContent = transformCallouts(htmlContent)
								htmlContent = transformMermaid(htmlContent)
								htmlContent = transformHighlights(htmlContent)
								canvasData.Nodes[i].HtmlContent = htmlContent
							}
						}
					}
				}
			}
		}
		// Re-marshal the hydrated data back to JSON string
		if hydratedJson, err := json.Marshal(canvasData); err == nil {
			source = hydratedJson
		}
	}

	outPath, webPath := getOutputPaths(c.relPath, c.nameWithoutExt, ".canvas")
	breadcrumbs := getBreadcrumbs(c.relPath, c.nameWithoutExt)
	setTreeActive(c.rootNode.Children, webPath)
	f, _ := os.Create(outPath)
	defer f.Close()
	mw := c.minifier.Writer("text/html", f)
	defer mw.Close()

	data := PageData{
		Title:         c.nameWithoutExt,
		BaseURL:       c.baseURL,
		SiteName:      c.siteName,
		CanvasContent: template.JS(string(source)),
		Breadcrumbs:   breadcrumbs,
		IsCanvas:      true,
		IsGraph:       false,
		Sidebar:       c.rootNode.Children,
	}
	c.template.Execute(mw, data)
	return webPath
}
