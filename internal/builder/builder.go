package builder

import (
	"encoding/json"
	"html/template"
	"io/fs"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	textTemplate "text/template"
	"time"

	"github.com/otaleghani/kiln/assets"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/html"
)

var (
	// OutputDir is the destination directory for the generated site.
	OutputDir string = "./public"
	// InputDir is the source directory containing the Obsidian vault.
	InputDir string = "./vault"
	// FlatUrls defines if the user opted in for flat urls.
	FlatUrls bool
)

// GraphNode represents a single node in the interactive graph view.
type GraphNode struct {
	ID    string `json:"id"`
	Label string `json:"label"`
	URL   string `json:"url"`
	Val   int    `json:"val"`
}

// Builder holds configuration for the site generation process.
type Builder struct {
	ThemeName string
	FontName  string
	BaseURL   string
	SiteName  string
}

// Build orchestrates the static site generation process.
// It walks the input directory, processes Markdown/Canvas files, and generates assets.
func Build(themeName, fontName, baseURL, siteName string) {
	start := time.Now()

	theme := resolveTheme(themeName, fontName)
	fileIndex, graphNodes := initBuild()
	rootNode := getRootNode(InputDir, baseURL)

	// Load and parse the base HTML layout
	layoutContent, err := assets.TemplateFS.ReadFile("layout.html")
	if err != nil {
		log.Fatal("Could not read layout.html: ", err)
	}
	tmplLayout, err := template.New("layout").Parse(string(layoutContent))
	if err != nil {
		log.Fatal("Layout parsing failed: ", err)
	}

	// Load and parse the CSS template
	cssContent, err := assets.TemplateFS.ReadFile("style.css")
	if err != nil {
		log.Fatal("Could not read style.css: ", err)
	}
	tmplCSS, err := textTemplate.New("css").Parse(string(cssContent))
	if err != nil {
		log.Fatal("CSS parsing failed: ", err)
	}

	// Load and parse the Graph JS template
	graphJsContent, err := assets.TemplateFS.ReadFile("graph.js")
	if err != nil {
		log.Fatal("Could not read graph.js: ", err)
	}
	tmplGraphJs, err := textTemplate.New("js").Parse(string(graphJsContent))
	if err != nil {
		log.Fatal("Graph JS parsing failed: ", err)
	}

	fileCount := 0
	var sitemapEntries []SitemapEntry

	// Determine the base path for link resolution (e.g., "/kiln" from "https://domain.com/kiln")
	u, _ := url.Parse(baseURL)
	basePath := u.Path
	markdownRenderer, resolver := newMarkdownParser(fileIndex, basePath)

	// Configure the HTML minifier
	minifier := minify.New()
	minifier.AddFunc("text/html", html.Minify)

	// Walk the input directory to process files
	err = filepath.WalkDir(InputDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip hidden files and directories (dotfiles)
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

		// Set the current file context for the link resolver
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
				theme:          theme,
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
				theme:          theme,
			}
			webPath := c.render()
			if baseURL != "" {
				addToSitemap(d, baseURL, webPath, &sitemapEntries)
			}
		default:
			// Copy static assets (images, PDFs, etc.) directly to output
			outPath := filepath.Join(OutputDir, getSlugPath(relPath))
			os.MkdirAll(filepath.Dir(outPath), 0755)
			copyFile(path, outPath)
		}

		return nil
	})

	if err != nil {
		log.Fatal("Walk failed: ", err)
	}

	// Generate the final CSS based on theme and font settings
	cssOut, _ := os.Create(filepath.Join(OutputDir, "style.css"))
	defer cssOut.Close()
	tmplCSS.Execute(cssOut, theme)

	// Write static JS files
	appJsContent, _ := assets.TemplateFS.ReadFile("app.js")
	os.WriteFile(filepath.Join(OutputDir, "app.js"), appJsContent, 0644)

	canvasJsContent, _ := assets.TemplateFS.ReadFile("canvas.js")
	os.WriteFile(filepath.Join(OutputDir, "canvas.js"), canvasJsContent, 0644)

	extractFonts(theme.Font, OutputDir)

	// Generate Graph JS with the correct BaseURL
	graphJsOut, _ := os.Create(filepath.Join(OutputDir, "graph.js"))
	defer graphJsOut.Close()

	type GraphJsTemplate struct {
		BaseURL string
	}
	tmplGraphJs.Execute(graphJsOut, GraphJsTemplate{BaseURL: baseURL})

	// Generate Graph JSON data
	graphJSON := map[string]any{
		"nodes": graphNodes,
		"links": resolver.Links,
	}
	jsonBytes, _ := json.Marshal(graphJSON)
	os.WriteFile(filepath.Join(OutputDir, "graph.json"), jsonBytes, 0644)

	// Render the dedicated Graph page
	graphPage := GraphRenderer{
		rootNode: rootNode,
		minifier: minifier,
		baseURL:  baseURL,
		siteName: siteName,
		template: tmplLayout,
		theme:    theme,
	}
	graphPage.render()

	if baseURL != "" {
		generateSitemap(sitemapEntries)
		generateRobots(baseURL)
	}

	log.Printf("Build complete! Processed %d files in %v", fileCount, time.Since(start))
}
