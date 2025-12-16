package builder

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"

	"github.com/tdewolff/minify/v2"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

// PageData is passed to the HTML template
type PageData struct {
	Title         string
	Content       template.HTML
	SiteName      string
	CanvasContent template.JS
	IsCanvas      bool
	IsGraph       bool
	Breadcrumbs   []string
	Sidebar       []*Node
	Frontmatter   map[string]any
	TOC           template.HTML
	BaseURL       string
}

type MarkdownNoteRenderer struct {
	path           string             // Path of the file to render
	nameWithoutExt string             // Name of the file without extention
	siteName       string             // The name of the website
	relPath        string             // Relative path of the file
	baseURL        string             // The base URL for links
	renderer       goldmark.Markdown  // Markdown renderer
	rootNode       *Node              // The root node of the vault
	template       *template.Template // Golang html template
	minifier       *minify.M          // Minifier
}

func (m *MarkdownNoteRenderer) render() string {
	breadcrumbs := getBreadcrumbs(m.relPath, m.nameWithoutExt)

	source, err := os.ReadFile(m.path)
	if err != nil {
		log.Printf("Failed to create markdown note at path %s. Error %s", m.path, err.Error())
		return ""
	}

	outPath, webPath := getOutputPaths(m.relPath, m.nameWithoutExt, ".md")

	// Separate Parse and Render to allow AST inspection
	ctx := parser.NewContext()
	doc := m.renderer.Parser().Parse(text.NewReader(source), parser.WithContext(ctx))

	// Extract TOC from AST
	tocHTML := extractTOC(doc, source)

	var buf bytes.Buffer
	if err := m.renderer.Renderer().Render(&buf, source, doc); err != nil {
		log.Printf("Failed to create markdown note at path %s. Error %s", m.path, err.Error())
		return ""
	}

	finalHTML := buf.String()
	finalHTML = transformCallouts(finalHTML)
	finalHTML = transformMermaid(finalHTML)
	finalHTML = transformHighlights(finalHTML)

	setTreeActive(m.rootNode.Children, webPath)

	f, err := os.Create(outPath)
	if err != nil {
		log.Printf("Failed to create markdown note at path %s. Error %s", m.path, err.Error())
		return ""
	}
	defer f.Close()
	mw := m.minifier.Writer("text/html", f) // minifies HTML
	defer mw.Close()

	data := PageData{
		BaseURL:     m.baseURL,
		Title:       m.nameWithoutExt,
		SiteName:    m.siteName,
		Content:     template.HTML(finalHTML),
		Breadcrumbs: breadcrumbs,
		IsCanvas:    false,
		IsGraph:     false,
		TOC:         tocHTML,
		Sidebar:     m.rootNode.Children,
		Frontmatter: meta.Get(ctx),
	}
	m.template.Execute(mw, data)

	return webPath
}

// extractTOC generates a flat HTML list with classes indicating nesting level
func extractTOC(doc ast.Node, source []byte) template.HTML {
	var listBuf bytes.Buffer
	hasHeadings := false

	ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering || n.Kind() != ast.KindHeading {
			return ast.WalkContinue, nil
		}

		h := n.(*ast.Heading)
		// Extract ID
		var id string
		if val, ok := h.Attribute([]byte("id")); ok {
			if idBytes, ok := val.([]byte); ok {
				id = string(idBytes)
			}
		}

		// Extract Text
		text := string(h.Text(source))

		// Render List Item
		// using class 'toc-level-N' so CSS can handle indentation (e.g. margin-left: 1rem * N)
		listBuf.WriteString(fmt.Sprintf(
			"<li class=\"toc-level-%d\"><a href=\"#%s\">%s</a></li>",
			h.Level, id, text,
		))
		hasHeadings = true

		return ast.WalkContinue, nil
	})

	if !hasHeadings {
		return template.HTML("")
	}

	return template.HTML("<nav class=\"toc\"><ul>" + listBuf.String() + "</ul></nav>")
}
