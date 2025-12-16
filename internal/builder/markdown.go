package builder

import (
	"path/filepath"
	"strings"

	chromaHTML "github.com/alecthomas/chroma/v2/formatters/html"
	mathjax "github.com/litao91/goldmark-mathjax"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	htmlRenderer "github.com/yuin/goldmark/renderer/html"
	"go.abhg.dev/goldmark/wikilink"
)

type GraphLink struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

type IndexResolver struct {
	Index         map[string]string
	Links         []GraphLink
	CurrentSource string
}

// ResolveWikilink is an helper function to resolve the different kinds of wikilinks in obsidian
func (r *IndexResolver) ResolveWikilink(n *wikilink.Node) ([]byte, error) {
	dest := string(n.Target)
	dest = strings.TrimSpace(dest)
	if len(n.Fragment) > 0 {
		dest += "#" + string(n.Fragment)
	}
	anchor := ""
	if idx := strings.Index(dest, "#"); idx != -1 {
		anchor = dest[idx:]
		dest = dest[:idx]
	}

	// Graph Linking Logic
	cleanDest := dest
	ext := filepath.Ext(cleanDest)
	if ext == ".md" || ext == ".canvas" {
		cleanDest = strings.TrimSuffix(cleanDest, ext)
	}

	// Record the link for the graph
	if r.CurrentSource != "" && cleanDest != "" {
		// Only record if we aren't linking to self
		if !strings.EqualFold(r.CurrentSource, cleanDest) {
			r.Links = append(r.Links, GraphLink{
				Source: r.CurrentSource,
				Target: cleanDest,
			})
		}
	}

	if dest == "" {
		return []byte(anchor), nil
	}

	if link, ok := r.Index[cleanDest]; ok {
		return []byte(link + anchor), nil
	}
	return []byte("/" + slugify(dest) + anchor), nil
}

// newMarkdownParser creates a new markdown and index resolver
func newMarkdownParser(fileIndex map[string]string) (goldmark.Markdown, *IndexResolver) {
	resolver := &IndexResolver{
		Index: fileIndex,
		Links: []GraphLink{}, // Initialize empty links
	}

	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			meta.Meta,
			&wikilink.Extender{Resolver: resolver},
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
		),
	)

	return md, resolver
}
