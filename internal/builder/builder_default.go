package builder

import (
	"encoding/json"
	"html/template"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/otaleghani/kiln/assets"
	"github.com/otaleghani/kiln/internal/obsidian"
	"github.com/otaleghani/kiln/internal/obsidian/bases"
	"github.com/otaleghani/kiln/internal/obsidian/markdown"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/html"
	"gopkg.in/yaml.v3"
)

func buildDefault(log *slog.Logger) {
	start := time.Now()

	// Resolves the theme
	theme := ResolveTheme(ThemeName, FontName, log)

	// Resolves and loads layout
	layout := resolveLayout(LayoutName, log)
	err := layout.loadLayout()
	if err != nil {
		log.Error("Error loading layout", "error", err)
		os.Exit(1)
	}

	// Scans vault
	obs := obsidian.New(
		obsidian.WithBaseURL(BaseURL),
		obsidian.WithFlatURLs(FlatUrls),
		obsidian.WithInputDir(InputDir),
		obsidian.WithOutputDir(OutputDir),
		obsidian.WithLogger(log),
	)
	// vaultScan, err := scanVault()
	err = obs.Scan()
	if err != nil {
		log.Error("Error scanning vault", "error", err)
		os.Exit(1)
	}

	// Creates markdown renderer
	obsidianMd := markdown.New(obs.Vault.FileIndex, func(path string) ([]byte, error) {
		return os.ReadFile(filepath.Join(InputDir, path))
	})

	// Get's the sidebar root node
	rootNode := obs.GenerateNavbar()

	site := &DefaultSite{
		// Scan:              vaultScan,
		BaseURL:           BaseURL,
		SiteName:          SiteName,
		Theme:             theme,
		Layout:            layout,
		Markdown:          obsidianMd,
		Minifier:          minify.New(),
		NavbarRoot:        rootNode,
		DisableLocalGraph: DisableLocalGraph,
		DisableTOC:        DisableTOC,
		FlatURLs:          FlatUrls,
		log:               log,
		Obsidian:          obs,
	}
	site.Minifier.AddFunc("text/html", html.Minify)

	// Divide the different files and call "RenderType" directly, instaed of calling Render
	notePages := []*obsidian.File{}
	basePages := []PageBase{}
	canvasPages := []*obsidian.File{}
	staticFiles := []*obsidian.File{}

	for _, file := range site.Obsidian.Vault.Files {
		l := log.With("file", file.Path)
		switch file.Ext {
		case ".md":
			notePages = append(notePages, file)
		case ".base":
			base, err := ParseBaseFile(file.Path)
			if err != nil {
				l.Error("Couldn't parse base file", "error", err)
			}
			base.File = file
			basePages = append(basePages, base)
		case ".canvas":
			canvasPages = append(canvasPages, file)
		default:
			if isAllowedExt(file.Ext) {
				staticFiles = append(staticFiles, file)
			}
		}
	}

	// Adds base URL to fonts
	site.Theme.Font.FontFaceReplaced = template.CSS(strings.Replace(
		string(site.Theme.Font.FontFace),
		"{{.Site.BaseURL}}",
		site.BaseURL,
		4,
	))

	nodes := []obsidian.GraphNode{}

	log.Info("Copying static assets...")
	for _, file := range staticFiles {
		l := log.With("file", file.Path)
		err := obsidian.CopyFile(file.Path, file.OutPath)
		if err != nil {
			l.Error("Couldn't copy file", "error", err)
		}
	}

	log.Info("Rendering canvas pages...")
	for _, file := range canvasPages {
		l := log.With("file", file.Path)
		err := site.RenderCanvas(file)
		if err != nil {
			l.Error("Couldn't copy file", "error", err)
			continue
		}
		nodes = append(nodes, obsidian.GraphNode{
			ID:    file.WebPath,
			Label: file.Name,
			URL:   file.WebPath,
			Val:   1,
			Type:  file.Ext,
		})
	}

	log.Info("Rendering base pages...")
	for _, base := range basePages {
		l := log.With("file", base.File.RelPath)
		err := site.RenderBase(&base, site.Obsidian.Vault.Files)
		if err != nil {
			l.Error("Couldn't render base", "error", err)
			continue
		}
		nodes = append(nodes, obsidian.GraphNode{
			ID:    base.File.WebPath,
			Label: base.File.Name,
			URL:   base.File.WebPath,
			Val:   1,
			Type:  base.File.Ext,
		})
	}

	log.Info("Rendering markdown pages...")
	for _, note := range notePages {
		l := log.With("file", note.RelPath)

		err := site.RenderNote(note)
		if err != nil {
			l.Error("Couldn't render note", "error", err)
			continue
		}
		nodes = append(nodes, obsidian.GraphNode{
			ID:    note.WebPath,
			Label: note.Name,
			URL:   note.WebPath,
			Val:   1,
			Type:  note.Ext,
		})
	}

	log.Info("Rendering folder pages...")
	for _, folder := range site.Obsidian.Vault.Folders {
		l := log.With("folder", folder.RelPath)
		if len(folder.Files) == 0 && len(folder.Folders) == 0 {
			l.Debug("Skipped empty folder", "folder", folder.Name)
			continue
		}

		err := site.RenderFolder(folder)
		if err != nil {
			l.Error("Couldn't render folder", "error", err)
			continue
		}
		nodes = append(nodes, obsidian.GraphNode{
			ID:    folder.WebPath,
			Label: folder.Name,
			URL:   folder.WebPath,
			Val:   1,
			Type:  "folder",
		})
	}

	log.Info("Rendering tag pages...")
	for _, tag := range site.Obsidian.Vault.Tags {
		l := log.With("tag", tag.Name)
		err := site.RenderTag(tag)
		if err != nil {
			l.Error("Couldn't render tag", "error", err)
		}
		nodes = append(nodes, obsidian.GraphNode{
			ID:    tag.WebPath,
			Label: tag.Name,
			URL:   tag.WebPath,
			Val:   1,
			Type:  "tag",
		})
	}

	log.Info("Rendering static files...")
	// Generate CSS based on the given theme/font settings
	cssOut, err := os.Create(filepath.Join(OutputDir, "style.css"))
	if err != nil {
		log.Error("Couldn't create 'style.css'", "error", err)
	}
	defer cssOut.Close()
	err = site.Layout.CssTemplate.Execute(cssOut, site)
	if err != nil {
		log.Error("Couldn't execute template for 'style.css'", "error", err)
	}

	// Copies over shared.css - Contains shared styles between layouts
	cssContent, err := assets.TemplateFS.ReadFile("shared.css")
	if err != nil {
		log.Error("Couldn't read 'shared.css'", "error", err)
	}
	err = os.WriteFile(filepath.Join(OutputDir, "shared.css"), cssContent, 0644)
	if err != nil {
		log.Error("Couldn't write 'shared.css'", "error", err)
	}

	// Generate app JS
	jsOut, err := os.Create(filepath.Join(OutputDir, "app.js"))
	if err != nil {
		log.Error("Couldn't create 'app.js'", "error", err)
	}
	defer jsOut.Close()
	err = site.Layout.JsTemplate.Execute(jsOut, site)
	if err != nil {
		log.Error("Couldn't execute template for 'app.js'", "error", err)
	}

	// Generate graph JS
	graphJsOut, err := os.Create(filepath.Join(OutputDir, "graph.js"))
	if err != nil {
		log.Error("Couldn't create 'graph.js'", "error", err)
	}
	defer graphJsOut.Close()
	err = site.Layout.JsGraphTemplate.Execute(graphJsOut, site)
	if err != nil {
		log.Error("Couldn't execute template for 'graph.js'", "error", err)
	}

	// Generate canvas JS
	canvasJsOut, err := os.Create(filepath.Join(OutputDir, "canvas.js"))
	if err != nil {
		log.Error("Couldn't create 'canvas.js'", "error", err)
	}
	defer canvasJsOut.Close()
	err = site.Layout.JsCanvasTemplate.Execute(canvasJsOut, site)
	if err != nil {
		log.Error("Couldn't execute template for 'canvas.js'", "error", err)
	}

	// Extracts fonts
	site.Theme.extractFonts(OutputDir, log)

	// Generate Graph JSON data
	markdownLinks := site.Markdown.Resolver.Links
	log.Debug("Markdown links", "amount", len(markdownLinks))
	links := append(site.Obsidian.GetFolderLinks(), markdownLinks...)
	links = append(site.Obsidian.GetTagLinks(), links...)
	log.Debug("Total links", "amount", len(links))
	graphJSON := map[string]any{
		"nodes": nodes,
		"links": links,
	}
	jsonBytes, err := json.Marshal(graphJSON)
	if err != nil {
		log.Error("Couldn't marshal JSON", "error", err)
	}
	err = os.WriteFile(filepath.Join(OutputDir, "graph.json"), jsonBytes, 0644)
	if err != nil {
		log.Error("Couldn't create 'graph.json'", "error", err)
	}

	err = site.RenderGraph()
	if err != nil {
		log.Error("Couldn't render 'graph.html'", "error", err)
	}

	err = site.Obsidian.GenerateSitemap()
	if err != nil {
		log.Error("Couldn't render 'sitemap.xml'", "error", err)
	}

	err = site.Obsidian.GenerateRobots()
	if err != nil {
		log.Error("Couldn't render 'robots.txt'", "error", err)
	}

	err = site.Obsidian.LoadCname()
	if err != nil {
		log.Error("Couldn't transfer 'CNAME' file", "error", err)
	}

	err = site.Obsidian.LoadFavicon()
	if err != nil {
		log.Error("Couldn't transfer 'favicon.ico' file", "error", err)
	}

	log.Info(
		"Build complete",
		"seconds",
		time.Since(start).Seconds(),
	)
}

// Render folders
func (s *DefaultSite) RenderFolder(f *obsidian.Folder) error {
	obsidian.SetNavbarNodeActive(s.NavbarRoot.Children, f.WebPath)

	// Creates outfile
	outFile, err := os.Create(f.OutPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// Setups the minifier
	minifierWriter := s.Minifier.Writer("text/html", outFile)
	defer minifierWriter.Close()

	// Get breadcrumbs
	breadcrumbs, err := s.Obsidian.GetFolderBreadcrumbs(f)
	if err != nil {
		return err
	}

	// Executes the template
	pageData := DefaultSitePageData{
		Site:        s,
		Folder:      f,
		IsFolder:    true,
		Breadcrumbs: breadcrumbs,
	}

	// Executes the template
	err = s.Layout.HtmlTemplate.Execute(minifierWriter, pageData)
	if err != nil {
		return err
	}

	return nil
}

// Render tags
func (s *DefaultSite) RenderTag(t *obsidian.Tag) error {
	s.log.Debug("Rendering tag", "name", t.Name, "files", len(t.Files))

	// Creates outfile
	outFile, err := os.Create(t.OutPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// Setups the minifier
	minifierWriter := s.Minifier.Writer("text/html", outFile)
	defer minifierWriter.Close()

	// Get breadcrumbs
	breadcrumbs := []obsidian.Breadcrumb{
		{Label: t.Name, Url: "#"},
	}

	// Executes the template
	pageData := DefaultSitePageData{
		Site:        s,
		Tag:         t,
		IsTag:       true,
		Breadcrumbs: breadcrumbs,
	}

	// Executes the template
	err = s.Layout.HtmlTemplate.Execute(minifierWriter, pageData)
	if err != nil {
		return err
	}

	return nil
}

func (s *DefaultSite) RenderBase(b *PageBase, allFiles []*obsidian.File) error {
	s.log.Info("Found base", "path", b.File.Path)
	obsidian.SetNavbarNodeActive(s.NavbarRoot.Children, b.File.WebPath)

	// Creates outfile
	outFile, err := os.Create(b.File.OutPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	breadcrumbs, err := s.Obsidian.GetBreadcrumbs(b.File)
	if err != nil {
		return err
	}

	// Setups the minifier
	minifierWriter := s.Minifier.Writer("text/html", outFile)
	defer minifierWriter.Close()

	activeFiles := bases.FilterFiles(allFiles, b.Filters)
	activeFiles = bases.FilterFiles(activeFiles, b.Views[0].Filters)
	var fileGroups []*bases.FileGroup
	if b.Views[0].GroupBy.Property != "" {
		fileGroups = bases.GroupFiles(activeFiles, b.Views[0].GroupBy.Property)
	}

	// Executes the template
	pageData := DefaultSitePageData{
		Site:        s,
		File:        b.File,
		IsBase:      true,
		Breadcrumbs: breadcrumbs,
		Base: BaseData{
			Groups:  fileGroups,
			Notes:   activeFiles,
			File:    b,
			Columns: b.Views[0].Order,
		},
	}

	// Executes the template
	err = s.Layout.HtmlTemplate.Execute(minifierWriter, pageData)
	if err != nil {
		return err
	}

	return nil
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
	obsidian.SetNavbarNodeActive(s.NavbarRoot.Children, "/graph")

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
		File: &obsidian.File{
			Name:    "Graph",
			WebPath: BaseURL + "/graph",
		},
		Breadcrumbs: []obsidian.Breadcrumb{
			{Label: "Home", Url: "/"}, {Label: "Graph", Url: "/graph"}},
	}

	s.Layout.HtmlTemplate.Execute(minifiedWriter, data)
	return nil
}

// RenderNote renders the given markdown file
func (s *DefaultSite) RenderNote(f *obsidian.File) error {
	obsidian.SetNavbarNodeActive(s.NavbarRoot.Children, f.WebPath)

	// Creates outfile
	outFile, err := os.Create(f.OutPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// Setups the minifier
	minifierWriter := s.Minifier.Writer("text/html", outFile)
	defer minifierWriter.Close()

	obsidianData, err := s.Markdown.Render(*f)
	if err != nil {
		return err
	}

	// Creates breadcrumbs
	breadcrumbs, err := s.Obsidian.GetBreadcrumbs(f)
	if err != nil {
		return err
	}

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

func (s *DefaultSite) RenderCanvas(f *obsidian.File) error {
	obsidian.SetNavbarNodeActive(s.NavbarRoot.Children, f.WebPath)

	l := s.log.With("path", f.RelPath)

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
				slugPath := s.Obsidian.GetSlugPath(relImgPath)
				webPath, err := s.Obsidian.GetPageWebPath(slugPath, linkedExt)
				if err != nil {
					l.Warn("Canvas rendering: Couldn't get web path", "error", err)
					continue
				}

				canvasData.Nodes[i].Src = webPath
			} else if linkedExt == ".md" {
				// Handle notes
				noteContent, err := os.ReadFile(linkedFilePath)
				if err != nil {
					l.Warn("Canvas rendering: Couldn't read note data", "note", linkedFilePath, "error", err)
					continue
				}

				// Render Markdown to HTML using the shared renderer
				renderedNote, err := s.Markdown.RenderNote(noteContent)
				if err != nil {
					l.Warn("Canvas rendering: Couldn't render note", "note", linkedFilePath, "error", err)
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
	breadcrumbs, err := s.Obsidian.GetBreadcrumbs(f)
	if err != nil {
		return err
	}

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

// ParseBaseFile parses the given base
func ParseBaseFile(filePath string) (PageBase, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return PageBase{}, err
	}

	var base PageBase
	err = yaml.Unmarshal(content, &base)
	return base, err
}

// DefaultSite holds the global state for a default generation
type DefaultSite struct {
	BaseURL           string                     // Base URL, used for sitemap.xml, robots.txt etc.
	SiteName          string                     // The Name of the site
	Theme             *Theme                     // Selected theme
	Layout            *Layout                    // Selected and loaded layout
	NavbarRoot        *obsidian.NavbarNode       // The root node of the sidebar
	Markdown          *markdown.ObsidianMarkdown // Handles the rendering of obsidian markdown
	Minifier          *minify.M                  // Minifier to minify html pages
	DisableLocalGraph bool                       // If set, disables the local graph
	DisableTOC        bool                       // If set, disables the Table of contents
	FlatURLs          bool                       // If set, handles flat urls (used in canonical)
	log               *slog.Logger
	Obsidian          *obsidian.Obsidian
}

// DefaultSitePage represents a page to be generated
type DefaultSitePageData struct {
	Site          *DefaultSite          // Default site data (to get theme, font, base path etc.)
	Content       template.HTML         // Rendered HTML content of the page
	TOC           template.HTML         // Table of contents
	Breadcrumbs   []obsidian.Breadcrumb // Breadcrumbs of the page
	File          *obsidian.File        // Infomation about the file
	Folder        *obsidian.Folder      // Information about the folder
	Tag           *obsidian.Tag         // Information about the folder
	CanvasContent template.JS           // Raw JS content for canvas hydration
	IsGraph       bool                  // Is the page a graph page?
	IsCanvas      bool                  // Is the page a canvas page?
	IsBase        bool                  // Is the page a base page?
	IsNote        bool                  // Is the page a note page?
	IsFolder      bool                  // Is the page a folder page?
	IsTag         bool                  // Is the page a tag page?
	Frontmatter   map[string]any        // Frontmatter data
	Base          BaseData
}

// BaseData is the data of the base
type BaseData struct {
	Groups  []*bases.FileGroup // File group
	Notes   []*obsidian.File   // Used for rendering bases
	File    *PageBase          // Used for rendering bases
	Columns []string
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
