package builder

import (
	"path"
	"path/filepath"
	"strings"

	chromaHTML "github.com/alecthomas/chroma/v2/formatters/html"
	mathjax "github.com/litao91/goldmark-mathjax"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	htmlRenderer "github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"
	"go.abhg.dev/goldmark/wikilink"
)

// GraphLink represents a directed edge in the note graph.
type GraphLink struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

// IndexResolver handles link resolution and graph tracking.
// It intercepts both Wikilinks ([[...]]) and standard Markdown links to ensure
// they point to the correct URL, respecting the site's BasePath.
type IndexResolver struct {
	Index         map[string][]string
	Links         []GraphLink
	CurrentSource string
	BasePath      string
}

// ResolveWikilink implements wikilink.Resolver.
// It maps Obsidian-style [[links]] to their final HTML paths and records the relationship.
func (r *IndexResolver) ResolveWikilink(n *wikilink.Node) ([]byte, error) {
	dest := string(n.Target)
	dest = strings.TrimSpace(dest)

	// Handle anchors (e.g., [[Page#Section]])
	if len(n.Fragment) > 0 {
		dest += "#" + string(n.Fragment)
	}
	anchor := ""
	if idx := strings.Index(dest, "#"); idx != -1 {
		anchor = dest[idx:]
		dest = dest[:idx]
	}

	cleanDest := dest
	ext := filepath.Ext(cleanDest)
	if ext == ".md" || ext == ".canvas" {
		cleanDest = strings.TrimSuffix(cleanDest, ext)
	}

	// Track this connection for the graph view
	r.recordLink(cleanDest)

	if dest == "" {
		return []byte(anchor), nil
	}

	// Logic to handle multiple matches
	if candidates, ok := r.Index[cleanDest]; ok && len(candidates) > 0 {
		var bestMatch string

		// 1. Priority: Check for Root Match
		// If the file is simply "Page.md", it lives in root.
		// If it is "Folder/Page.md", it does not.
		for _, pathStr := range candidates {
			// Check if path has no directory separators
			if !strings.Contains(pathStr, "/") && !strings.Contains(pathStr, "\\") {
				bestMatch = pathStr
				break
			}
		}

		// 2. Priority: If no root match, find the shortest path (closest to root)
		if bestMatch == "" {
			shortest := candidates[0]
			for _, pathStr := range candidates {
				// Simply comparing string length is a good proxy for folder depth here
				if len(pathStr) < len(shortest) {
					shortest = pathStr
				}
			}
			bestMatch = shortest
		}

		finalLink := path.Join(r.BasePath, bestMatch)
		return []byte(finalLink + anchor), nil
	}

	// Fallback: assume the slug matches the destination
	finalPath := path.Join(r.BasePath, slugify(dest))
	return []byte(finalPath + anchor), nil
}

// RegisterFuncs registers custom renderers for standard links and images.
// This allows us to intercept [text](url) and ![alt](url) to fix paths.
func (r *IndexResolver) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindLink, r.renderLink)
	reg.Register(ast.KindImage, r.renderImage)
}

// renderLink writes the HTML for standard markdown links.
func (r *IndexResolver) renderLink(
	w util.BufWriter,
	source []byte,
	node ast.Node,
	entering bool,
) (ast.WalkStatus, error) {
	n := node.(*ast.Link)
	if entering {
		newDest := r.resolvePath(string(n.Destination))

		w.WriteString("<a href=\"")
		w.WriteString(newDest)
		w.WriteString("\"")
		if n.Title != nil {
			w.WriteString(" title=\"")
			w.Write(n.Title)
			w.WriteString("\"")
		}
		w.WriteString(">")
	} else {
		w.WriteString("</a>")
	}
	return ast.WalkContinue, nil
}

// renderImage writes the HTML for standard markdown images.
func (r *IndexResolver) renderImage(
	w util.BufWriter,
	source []byte,
	node ast.Node,
	entering bool,
) (ast.WalkStatus, error) {
	n := node.(*ast.Image)
	if entering {
		newDest := r.resolvePath(string(n.Destination))

		w.WriteString("<img src=\"")
		w.WriteString(newDest)
		w.WriteString("\" alt=\"")
		w.Write(n.Text(source))
		w.WriteString("\"")
		if n.Title != nil {
			w.WriteString(" title=\"")
			w.Write(n.Title)
			w.WriteString("\"")
		}
		w.WriteString(">")
	}
	return ast.WalkSkipChildren, nil
}

// resolvePath ensures internal links respect the site's BasePath.
// It leaves external links, mailto, and anchors untouched.
func (r *IndexResolver) resolvePath(dest string) string {
	if strings.Contains(dest, "://") || strings.HasPrefix(dest, "mailto:") ||
		strings.HasPrefix(dest, "#") {
		return dest
	}

	// If it's an absolute path to the site root, prepend BasePath.
	// e.g. /style.css -> /kiln/style.css
	if strings.HasPrefix(dest, "/") {
		if r.BasePath != "/" && strings.HasPrefix(dest, r.BasePath+"/") {
			return dest
		}
		return path.Join(r.BasePath, dest)
	}

	return dest
}

// recordLink adds a directed edge to the graph if the source and target differ.
func (r *IndexResolver) recordLink(target string) {
	if r.CurrentSource != "" && target != "" {
		if !strings.EqualFold(r.CurrentSource, target) {
			r.Links = append(r.Links, GraphLink{
				Source: r.CurrentSource,
				Target: target,
			})
		}
	}
}

// newMarkdownParser creates a Goldmark instance configured for Obsidian compatibility.
// It enables GFM, MathJax, syntax highlighting, and custom link resolution.
func newMarkdownParser(
	// CHANGED: Update this type to match the output of initBuild
	fileIndex map[string][]string,
	basePath string,
) (goldmark.Markdown, *IndexResolver) {

	if basePath == "" {
		basePath = "/"
	}

	resolver := &IndexResolver{
		Index:    fileIndex,
		Links:    []GraphLink{},
		BasePath: basePath,
	}

	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			meta.Meta,
			&wikilink.Extender{Resolver: resolver}, // Hook 1: Wikilinks
			highlighting.NewHighlighting(
				highlighting.WithFormatOptions(
					chromaHTML.WithClasses(true),
				),
			),
			mathjax.NewMathJax(
				mathjax.WithInlineDelim("$", "$"),
				mathjax.WithBlockDelim("$$", "$$"),
			),
		),
		goldmark.WithParserOptions(parser.WithAutoHeadingID()),
		goldmark.WithRendererOptions(
			htmlRenderer.WithUnsafe(),
			htmlRenderer.WithHardWraps(),
			// Hook 2: Standard Links & Images
			renderer.WithNodeRenderers(
				util.Prioritized(resolver, 1000),
			),
		),
	)

	return md, resolver
}
