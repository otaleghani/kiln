package builder

import (
	"encoding/json"
	"html/template"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	textTemplate "text/template"
	"time"

	"github.com/otaleghani/kiln/assets"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/html"
)

var OutputDir string = "./public"
var InputDir string = "./vault"

type GraphNode struct {
	ID    string `json:"id"`
	Label string `json:"label"`
	URL   string `json:"url"`
	Val   int    `json:"val"`
}

type Builder struct {
	ThemeName string
	FontName  string
	BaseURL   string
	SiteName  string
}

func Build(themeName, fontName, baseURL, siteName string) {
	start := time.Now()
	theme := resolveTheme(themeName, fontName)

	fileIndex, graphNodes := initBuild()
	rootNode := getRootNode(InputDir, baseURL)

	// Parses layouts
	layoutContent, err := assets.TemplateFS.ReadFile("layout.html")
	if err != nil {
		log.Fatal("Could not read layout.html: ", err)
	}
	tmplLayout, err := template.New("layout").Parse(string(layoutContent))
	if err != nil {
		log.Fatal("Layout parsing failed: ", err)
	}

	cssContent, err := assets.TemplateFS.ReadFile("style.css")
	if err != nil {
		log.Fatal("Could not read style.css: ", err)
	}
	tmplCSS, err := textTemplate.New("css").Parse(string(cssContent))
	if err != nil {
		log.Fatal("CSS parsing failed: ", err)
	}

	fileCount := 0
	var sitemapEntries []SitemapEntry

	markdownRenderer, resolver := newMarkdownParser(fileIndex, baseURL)

	minifier := minify.New()
	minifier.AddFunc("text/html", html.Minify)

	err = filepath.WalkDir(InputDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip hidden files and folders (e.g. .obsidian, .trash, .git)
		if strings.HasPrefix(d.Name(), ".") && path != InputDir {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if d.IsDir() {
			return nil
		}

		relPath, _ := filepath.Rel(InputDir, path)
		ext := filepath.Ext(path)
		nameWithoutExt := strings.TrimSuffix(d.Name(), ext)

		// Set context for resolver
		resolver.CurrentSource = nameWithoutExt

		switch ext {
		case ".md":
			fileCount++
			m := MarkdownNoteRenderer{
				path:           path,
				nameWithoutExt: nameWithoutExt,
				siteName:       siteName,
				relPath:        relPath,
				renderer:       markdownRenderer,
				rootNode:       rootNode,
				baseURL:        baseURL,
				template:       tmplLayout,
				minifier:       minifier,
			}
			webPath := m.render()
			if baseURL != "" {
				addToSitemap(d, baseURL, webPath, &sitemapEntries)
			}
		case ".canvas":
			fileCount++
			c := CanvasNoteRenderer{
				path:           path,
				relPath:        relPath,
				siteName:       siteName,
				nameWithoutExt: nameWithoutExt,
				template:       tmplLayout,
				renderer:       markdownRenderer,
				baseURL:        baseURL,
				minifier:       minifier,
				rootNode:       rootNode,
			}
			webPath := c.render()
			if baseURL != "" {
				addToSitemap(d, baseURL, webPath, &sitemapEntries)
			}
		default:
			outPath := filepath.Join(OutputDir, getSlugPath(relPath))
			os.MkdirAll(filepath.Dir(outPath), 0755)
			copyFile(path, outPath)
		}

		return nil
	})

	if err != nil {
		log.Fatal("Walk failed: ", err)
	}

	// Generate CSS based on the selected theme and font
	cssOut, _ := os.Create(filepath.Join(OutputDir, "style.css"))
	defer cssOut.Close()
	tmplCSS.Execute(cssOut, theme)

	// Static files - App javascript
	appJsContent, _ := assets.TemplateFS.ReadFile("app.js")
	os.WriteFile(filepath.Join(OutputDir, "app.js"), appJsContent, 0644)

	// Static files - Canvas javascript
	canvasJsContent, _ := assets.TemplateFS.ReadFile("canvas.js")
	os.WriteFile(filepath.Join(OutputDir, "canvas.js"), canvasJsContent, 0644)

	// Static files - Graph javascript
	graphJsContent, _ := assets.TemplateFS.ReadFile("graph.js")
	os.WriteFile(filepath.Join(OutputDir, "graph.js"), graphJsContent, 0644)

	// Static files - Graph JSON
	graphJSON := map[string]any{
		"nodes": graphNodes,
		"links": resolver.Links,
	}
	jsonBytes, _ := json.Marshal(graphJSON)
	os.WriteFile(filepath.Join(OutputDir, "graph.json"), jsonBytes, 0644)

	// Generates graph page
	graphPage := GraphRenderer{
		rootNode: rootNode,
		minifier: minifier,
		baseURL:  baseURL,
		siteName: siteName,
		template: tmplLayout,
	}
	graphPage.render()

	if baseURL != "" {
		generateSitemap(sitemapEntries)
		generateRobots(baseURL)
	}

	log.Printf("Build complete! Processed %d files in %v", fileCount, time.Since(start))
}
