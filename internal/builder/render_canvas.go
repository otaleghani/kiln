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

// CanvasNoteRenderer handles the conversion of Obsidian .canvas files into HTML pages.
// It parses the JSON structure, hydrates referenced Markdown/Images, and renders the layout.
type CanvasNoteRenderer struct {
	path           string             // Absolute path to source file
	relPath        string             // Path relative to input directory
	siteName       string             // Global site name
	nameWithoutExt string             // Filename without extension (slug)
	template       *template.Template // Main HTML layout
	renderer       goldmark.Markdown  // Shared Markdown parser
	minifier       *minify.M          // HTML minifier
	rootNode       *Node              // Navigation tree root
	baseURL        string             // Site Base URL
	theme          Theme              // Visual theme settings
}

// CanvasNode represents a single element (card, file, group) within the Canvas JSON.
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

	// Fields injected during build time for the frontend:
	HtmlContent string `json:"htmlContent,omitempty"` // Rendered HTML for markdown files
	IsImage     bool   `json:"isImage,omitempty"`     // Flag for image nodes
	Src         string `json:"src,omitempty"`         // Web path for image source
	URL         string `json:"url,omitempty"`         // Link URL
}

// CanvasData represents the top-level structure of an Obsidian Canvas file.
type CanvasData struct {
	Nodes []CanvasNode     `json:"nodes"`
	Edges []map[string]any `json:"edges"`
}

// render processes the canvas file and writes the final HTML to disk.
// It returns the web-accessible path of the generated page.
func (c *CanvasNoteRenderer) render() string {
	source, err := os.ReadFile(c.path)
	if err != nil {
		log.Printf("Failed to create canvas note. Error: %s", err.Error())
	}

	// Parse the raw Canvas JSON
	var canvasData CanvasData
	if err := json.Unmarshal(source, &canvasData); err == nil {

		// Hydrate nodes: Read linked files and inject content directly into the JSON
		for i, node := range canvasData.Nodes {
			if node.Type == "file" && node.File != "" {

				// Resolve full path to the linked file
				linkedFilePath := filepath.Join(InputDir, node.File)
				linkedExt := strings.ToLower(filepath.Ext(linkedFilePath))

				// Verify file exists before processing
				if _, err := os.Stat(linkedFilePath); err == nil {

					// Case 1: Handle Images
					if isImageExt(linkedExt) {
						canvasData.Nodes[i].IsImage = true

						// Generate web path (must match the logic used by the static asset copier)
						relImgPath, _ := filepath.Rel(InputDir, linkedFilePath)
						parts := strings.Split(relImgPath, string(os.PathSeparator))
						for k, p := range parts {
							parts[k] = slugify(p)
						}

						// TODO: Check if c.baseURL path needs to be prepended here if site is not at root
						canvasData.Nodes[i].Src = "/" + strings.Join(parts, "/")

						// Case 2: Handle Markdown Notes
					} else if linkedExt == ".md" {
						noteContent, err := os.ReadFile(linkedFilePath)
						if err == nil {
							// Render Markdown to HTML using the shared renderer
							var buf bytes.Buffer
							if err := c.renderer.Convert(noteContent, &buf); err == nil {
								htmlContent := buf.String()

								// Apply post-processing hooks
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

		// Re-marshal the hydrated data to inject into the template
		if hydratedJson, err := json.Marshal(canvasData); err == nil {
			source = hydratedJson
		}
	}

	// Prepare output paths and navigation
	outPath, webPath := getPageOutputPath(c.relPath, c.nameWithoutExt, ".canvas")
	breadcrumbs := getBreadcrumbs(c.relPath, c.nameWithoutExt)
	setTreeActive(c.rootNode.Children, webPath)

	// Open output file
	f, _ := os.Create(outPath)
	defer f.Close()

	// Wrap writer with minifier
	mw := c.minifier.Writer("text/html", f)
	defer mw.Close()

	data := PageData{
		Title:         c.nameWithoutExt,
		BaseURL:       c.baseURL,
		SiteName:      c.siteName,
		CanvasContent: template.JS(string(source)), // Inject JSON as raw JS
		Breadcrumbs:   breadcrumbs,
		IsCanvas:      true,
		IsGraph:       false,
		Sidebar:       c.rootNode.Children,
		Font:          c.theme.Font,
	}

	c.template.Execute(mw, data)
	return webPath
}
