package builder

import (
	"encoding/json"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/otaleghani/kiln/internal/log"
	obsidian "github.com/otaleghani/kiln/pkg/obsidian-markdown"
	obsidianmarkdown "github.com/otaleghani/kiln/pkg/obsidian-markdown"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/html"
)

func buildDefault() {
	start := time.Now()

	// Resolves the theme
	theme := resolveTheme(ThemeName, FontName)

	// Resolves and loads layout
	layout := resolveLayout(LayoutName)
	err := layout.loadLayout()
	if err != nil {
		log.Fatal("Error loading layout", log.FieldError, err)
	}

	// Scans vault
	vaultScan, err := scanVault()
	if err != nil {
		log.Fatal("Error scanning vault", log.FieldError, err)
	}

	// Convert scanned files into obsidian-markdown.File
	index := make(map[string][]*obsidian.File)
	for name, files := range vaultScan.FileIndex {
		index[name] = []*obsidian.File{}
		for _, file := range files {
			index[name] = append(
				index[name],
				&obsidian.File{
					RelPath: file.RelPath,
					Path:    file.Path,
					WebPath: file.WebPath,
					Name:    file.Name,
					Ext:     file.Ext,
					OutPath: file.OutPath,
				},
			)
		}
	}

	// Creates markdown renderer
	obsidianMd := obsidian.New(index, func(path string) ([]byte, error) {
		return os.ReadFile(filepath.Join(InputDir, path))
	})

	// Get's the sidebar root node
	rootNode := getSidebarRootNode(InputDir, BaseURL)

	site := &DefaultSite{
		Scan:              vaultScan,
		BaseURL:           BaseURL,
		SiteName:          SiteName,
		Theme:             theme,
		Layout:            layout,
		Markdown:          obsidianMd,
		Minifier:          minify.New(),
		SidebarRootNode:   rootNode,
		DisableLocalGraph: DisableLocalGraph,
		DisableTOC:        DisableTOC,
	}
	site.Minifier.AddFunc("text/html", html.Minify)

	log.Info("Rendering pages...")
	for _, page := range site.Scan.Files {
		l := log.Default.WithFile(page.RelPath)
		l.Debug("Rendering...")
		err := site.Render(page)
		if err != nil {
			l.Error("Couldn't render note", log.FieldError, err)
		}
	}

	log.Info("Rendering static files...")
	// Generate CSS based on the given theme/font settings
	cssOut, err := os.Create(filepath.Join(OutputDir, "style.css"))
	if err != nil {
		log.Error("Couldn't create 'style.css'", log.FieldError, err)
	}
	defer cssOut.Close()
	err = site.Layout.CssTemplate.Execute(cssOut, site)
	if err != nil {
		log.Error("Couldn't execute template for 'style.css'", log.FieldError, err)
	}

	// Generate app JS
	jsOut, err := os.Create(filepath.Join(OutputDir, "app.js"))
	if err != nil {
		log.Error("Couldn't create 'app.js'", log.FieldError, err)
	}
	defer jsOut.Close()
	err = site.Layout.JsTemplate.Execute(jsOut, site)
	if err != nil {
		log.Error("Couldn't execute template for 'app.js'", log.FieldError, err)
	}

	// Generate graph JS
	graphJsOut, err := os.Create(filepath.Join(OutputDir, "graph.js"))
	if err != nil {
		log.Error("Couldn't create 'graph.js'", log.FieldError, err)
	}
	defer graphJsOut.Close()
	err = site.Layout.JsGraphTemplate.Execute(graphJsOut, site)
	if err != nil {
		log.Error("Couldn't execute template for 'graph.js'", log.FieldError, err)
	}

	// Generate canvas JS
	canvasJsOut, err := os.Create(filepath.Join(OutputDir, "canvas.js"))
	if err != nil {
		log.Error("Couldn't create 'canvas.js'", log.FieldError, err)
	}
	defer canvasJsOut.Close()
	err = site.Layout.JsCanvasTemplate.Execute(canvasJsOut, site)
	if err != nil {
		log.Error("Couldn't execute template for 'canvas.js'", log.FieldError, err)
	}

	// Extracts fonts
	site.Theme.extractFonts(OutputDir)

	// Generate Graph JSON data
	graphJSON := map[string]any{
		"nodes": site.Scan.GraphNodes,
		"links": site.Markdown.Resolver.Links,
	}
	jsonBytes, err := json.Marshal(graphJSON)
	if err != nil {
		log.Error("Couldn't marshal JSON", log.FieldError, err)
	}
	err = os.WriteFile(filepath.Join(OutputDir, "graph.json"), jsonBytes, 0644)
	if err != nil {
		log.Error("Couldn't create 'graph.json'", log.FieldError, err)
	}

	err = site.RenderGraph()
	if err != nil {
		log.Error("Couldn't render 'graph.html'", log.FieldError, err)
	}

	err = site.Scan.Sitemap.generate()
	if err != nil {
		log.Error("Couldn't render 'sitemap.xml'", log.FieldError, err)
	}

	err = site.Scan.Sitemap.generateRobots()
	if err != nil {
		log.Error("Couldn't render 'robots.txt'", log.FieldError, err)
	}

	err = loadCname()
	if err != nil {
		log.Error("Couldn't transfer 'CNAME' file", log.FieldError, err)
	}

	err = loadFavicon()
	if err != nil {
		log.Error("Couldn't transfer 'favicon.ico' file", log.FieldError, err)
	}

	log.Info(
		"Build complete",
		"seconds",
		time.Since(start).Seconds(),
	)
}

// Render renders the given file based on it's extension.
// It saves the result in the output directory.
func (s *DefaultSite) Render(f *File) error {
	switch f.Ext {
	case ".md":
		return s.RenderNote(f)
	case ".canvas":
		return s.RenderCanvas(f)
	default:
		return nil
	}
}

func (s *DefaultSite) RenderGraph() error {
	graphOutPath := ""
	if FlatUrls {
		err := os.MkdirAll(filepath.Join(OutputDir, "graph"), 0755)
		if err != nil {
			return err
		}
		graphOutPath = filepath.Join(OutputDir, "graph", "index.html")
	} else {
		graphOutPath = filepath.Join(OutputDir, "graph.html")
	}
	setSidebarNodeActive(s.SidebarRootNode.Children, "/graph")

	fGraph, err := os.Create(graphOutPath)
	if err != nil {
		return err
	}
	defer fGraph.Close()

	// Wrap writer with minifier
	minifiedWriter := s.Minifier.Writer("text/html", fGraph)
	defer minifiedWriter.Close()

	// The frontend script targets this ID to mount the global graph
	graphHTML := `<div id="global-graph-container" style=""></div>`

	data := DefaultSitePageData{
		Site:        s,
		Content:     template.HTML(graphHTML),
		IsGraph:     true,
		Frontmatter: make(map[string]any),
		File: &File{
			Name:    "Graph",
			WebPath: BaseURL + "/graph",
		},
		Breadcrumbs: []string{"Home", "Graph"},
	}

	s.Layout.HtmlTemplate.Execute(minifiedWriter, data)
	return nil
}

// RenderNote renders the given markdown file
func (s *DefaultSite) RenderNote(f *File) error {
	setSidebarNodeActive(s.SidebarRootNode.Children, f.WebPath)

	// Creates outfile
	outFile, err := os.Create(f.OutPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// Setups the minifier
	minifierWriter := s.Minifier.Writer("text/html", outFile)
	defer minifierWriter.Close()

	// Renders markdown
	obsidianFile := obsidianmarkdown.File{
		Path:    f.Path,
		RelPath: f.RelPath,
		Ext:     f.Ext,
		Name:    f.Name,
		OutPath: f.OutPath,
		WebPath: f.WebPath,
	}
	obsidianData, err := s.Markdown.Render(obsidianFile)
	if err != nil {
		return err
	}

	// Creates breadcrumbs
	breadcrumbs := f.Breadcrumbs()

	// Executes the template
	pageData := DefaultSitePageData{
		Site:        s,
		Frontmatter: obsidianData.Frontmatter,
		Content:     template.HTML(obsidianData.Content),
		TOC:         obsidianData.TOC,
		File:        f,
		Breadcrumbs: breadcrumbs,
		IsNote:      true,
	}

	// Executes the template
	err = s.Layout.HtmlTemplate.Execute(minifierWriter, pageData)
	if err != nil {
		return err
	}

	return nil
}

func (s *DefaultSite) RenderCanvas(f *File) error {
	setSidebarNodeActive(s.SidebarRootNode.Children, f.WebPath)

	l := log.Default.WithFile(f.RelPath)

	// Read file
	source, err := os.ReadFile(f.Path)
	if err != nil {
		return err
	}

	// Unmarshal JSON data
	var canvasData CanvasData
	err = json.Unmarshal(source, &canvasData)
	if err != nil {
		return err
	}

	// Read linked files and inject content in the JSON
	for i, node := range canvasData.Nodes {
		if node.Type == "file" && node.File != "" {
			// Resolve full path to the linked file
			linkedFilePath := filepath.Join(InputDir, node.File)
			linkedExt := strings.ToLower(filepath.Ext(linkedFilePath))

			// Verify file exists before processing
			_, err := os.Stat(linkedFilePath)
			if err != nil {
				l.Warn("Canvas links to non-existant file", "missing", linkedFilePath)
				continue
			}

			if isImageExt(linkedExt) {
				// Handle images
				canvasData.Nodes[i].IsImage = true

				// Generate web path (must match the logic used by the static asset copier)
				relImgPath, err := filepath.Rel(InputDir, linkedFilePath)
				if err != nil {
					l.Warn(
						"Canvas rendering: Couldn't create relative path for image",
						"image",
						linkedFilePath,
					)
					continue
				}
				slugPath := getSlugPath(relImgPath)
				webPath := getPageWebPath(slugPath, linkedExt)

				canvasData.Nodes[i].Src = webPath
			} else if linkedExt == ".md" {
				// Handle notes
				noteContent, err := os.ReadFile(linkedFilePath)
				if err != nil {
					l.Warn("Canvas rendering: Couldn't read note data", "note", linkedFilePath, log.FieldError, err)
					continue
				}

				// Render Markdown to HTML using the shared renderer
				renderedNote, err := s.Markdown.RenderNote(noteContent)
				if err != nil {
					l.Warn("Canvas rendering: Couldn't render note", "note", linkedFilePath, log.FieldError, err)
					continue
				}
				canvasData.Nodes[i].HtmlContent = renderedNote
			}
		}
	}

	// Creates outfile
	outFile, err := os.Create(f.OutPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// Setups the minifier
	minifierWriter := s.Minifier.Writer("text/html", outFile)
	defer minifierWriter.Close()

	// Creates breadcrumbs
	breadcrumbs := f.Breadcrumbs()

	// Executes the template
	pageData := DefaultSitePageData{
		Site:          s,
		CanvasContent: template.JS(string(source)),
		File:          f,
		Breadcrumbs:   breadcrumbs,
		IsCanvas:      true,
	}

	s.Layout.HtmlTemplate.Execute(minifierWriter, pageData)

	return nil
}

// DefaultSite holds the global state for a default generation
type DefaultSite struct {
	Scan              *VaultScan                 // The result of scanVault
	BaseURL           string                     // Base URL, used for sitemap.xml, robots.txt etc.
	SiteName          string                     // The Name of the site
	Theme             *Theme                     // Selected theme
	Layout            *Layout                    // Selected and loaded layout
	SidebarRootNode   *SidebarNode               // The root node of the sidebar
	Markdown          *obsidian.ObsidianMarkdown // Handles the rendering of obsidian markdown
	Minifier          *minify.M                  // Minifier to minify html pages
	DisableLocalGraph bool                       // If set, disables the local graph
	DisableTOC        bool                       // If set, disables the Table of contents
}

// DefaultSitePage represents a page to be generated
type DefaultSitePageData struct {
	Site          *DefaultSite   // Default site data (to get theme, font, base path etc.)
	Content       template.HTML  // Rendered HTML content of the page
	TOC           template.HTML  // Table of contents
	Breadcrumbs   []string       // Breadcrumbs of the page
	File          *File          // Infomation about the file
	CanvasContent template.JS    // Raw JS content for canvas hydration
	IsGraph       bool           // Is the page a graph page?
	IsCanvas      bool           // Is the page a canvas page?
	IsBase        bool           // Is the page a base page?
	IsNote        bool           // Is the page a note page?
	Frontmatter   map[string]any // Frontmatter data
}

// CanvasData represents the top-level structure of an Obsidian Canvas file.
type CanvasData struct {
	Nodes []CanvasNode     `json:"nodes"`
	Edges []map[string]any `json:"edges"`
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
