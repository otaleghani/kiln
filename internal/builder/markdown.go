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

//
// type GraphLink struct {
// 	Source string `json:"source"`
// 	Target string `json:"target"`
// }
//
// type IndexResolver struct {
// 	Index         map[string]string
// 	Links         []GraphLink
// 	CurrentSource string
// }
//
// // ResolveWikilink is an helper function to resolve the different kinds of wikilinks in obsidian
// func (r *IndexResolver) ResolveWikilink(n *wikilink.Node) ([]byte, error) {
// 	dest := string(n.Target)
// 	dest = strings.TrimSpace(dest)
// 	if len(n.Fragment) > 0 {
// 		dest += "#" + string(n.Fragment)
// 	}
// 	anchor := ""
// 	if idx := strings.Index(dest, "#"); idx != -1 {
// 		anchor = dest[idx:]
// 		dest = dest[:idx]
// 	}
//
// 	// Graph Linking Logic
// 	cleanDest := dest
// 	ext := filepath.Ext(cleanDest)
// 	if ext == ".md" || ext == ".canvas" {
// 		cleanDest = strings.TrimSuffix(cleanDest, ext)
// 	}
//
// 	// Record the link for the graph
// 	if r.CurrentSource != "" && cleanDest != "" {
// 		// Only record if we aren't linking to self
// 		if !strings.EqualFold(r.CurrentSource, cleanDest) {
// 			r.Links = append(r.Links, GraphLink{
// 				Source: r.CurrentSource,
// 				Target: cleanDest,
// 			})
// 		}
// 	}
//
// 	if dest == "" {
// 		return []byte(anchor), nil
// 	}
//
// 	if link, ok := r.Index[cleanDest]; ok {
// 		return []byte(link + anchor), nil
// 	}
// 	return []byte("/" + slugify(dest) + anchor), nil
// }
//
// // newMarkdownParser creates a new markdown and index resolver
// func newMarkdownParser(fileIndex map[string]string) (goldmark.Markdown, *IndexResolver) {
// 	resolver := &IndexResolver{
// 		Index: fileIndex,
// 		Links: []GraphLink{}, // Initialize empty links
// 	}
//
// 	md := goldmark.New(
// 		goldmark.WithExtensions(
// 			extension.GFM,
// 			meta.Meta,
// 			&wikilink.Extender{Resolver: resolver},
// 			highlighting.NewHighlighting(
// 				highlighting.WithFormatOptions(
// 					chromaHTML.WithClasses(true),
// 				),
// 			),
// 			mathjax.NewMathJax(
// 				mathjax.WithInlineDelim("$", "$"),
// 				mathjax.WithBlockDelim("$$", "$$"),
// 			),
// 		),
// 		goldmark.WithParserOptions(parser.WithAutoHeadingID()),
// 		goldmark.WithRendererOptions(
// 			htmlRenderer.WithUnsafe(),
// 			htmlRenderer.WithHardWraps(),
// 		),
// 	)
//
// 	return md, resolver
// }

type GraphLink struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

// IndexResolver handles Wikilinks, Standard Links, and Graph tracking
type IndexResolver struct {
	Index         map[string]string
	Links         []GraphLink
	CurrentSource string
	BasePath      string // The path prefix (e.g., "/kiln")
}

// --- 1. WIKILINK RESOLVER (Implements wikilink.Resolver) ---

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

	cleanDest := dest
	ext := filepath.Ext(cleanDest)
	if ext == ".md" || ext == ".canvas" {
		cleanDest = strings.TrimSuffix(cleanDest, ext)
	}

	// Track this link for the Graph
	r.recordLink(cleanDest)

	if dest == "" {
		return []byte(anchor), nil
	}

	// Check the File Index first
	if link, ok := r.Index[cleanDest]; ok {
		// If index has a full path, ensure it respects BasePath if it's relative?
		// Usually index stores the slug. Let's prepend BasePath safely.
		finalLink := path.Join(r.BasePath, link)
		return []byte(finalLink + anchor), nil
	}

	// Fallback: construct slugified path
	// OLD: return []byte("/" + slugify(dest) + anchor), nil
	// NEW: Join with BasePath
	finalPath := path.Join(r.BasePath, slugify(dest))
	return []byte(finalPath + anchor), nil
}

// --- 2. STANDARD NODE RENDERER (Implements renderer.NodeRenderer) ---
// This intercepts [Link](/foo) and ![Image](/foo.png)

func (r *IndexResolver) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindLink, r.renderLink)
	reg.Register(ast.KindImage, r.renderImage)
}

func (r *IndexResolver) renderLink(
	w util.BufWriter,
	source []byte,
	node ast.Node,
	entering bool,
) (ast.WalkStatus, error) {
	n := node.(*ast.Link)
	if entering {
		// Resolve path using BaseURL logic
		newDest := r.resolvePath(string(n.Destination))

		// Optional: Track standard links in graph too?
		// If yes: r.recordLink(cleanPath(string(n.Destination)))

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

// --- 3. HELPER FUNCTIONS ---

// resolvePath prepends BasePath to absolute links (starting with /)
func (r *IndexResolver) resolvePath(dest string) string {
	if strings.Contains(dest, "://") || strings.HasPrefix(dest, "mailto:") ||
		strings.HasPrefix(dest, "#") {
		return dest
	}

	// If it starts with /, it's absolute to the site root. Prepend BasePath.
	if strings.HasPrefix(dest, "/") {
		// Prevent double-prefixing if BasePath is already there
		// e.g. BasePath=/kiln, dest=/kiln/style.css
		if r.BasePath != "/" && strings.HasPrefix(dest, r.BasePath+"/") {
			return dest
		}
		return path.Join(r.BasePath, dest)
	}

	return dest
}

// recordLink centralized logic for graph recording
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

// Updated signature: accepts basePath
func newMarkdownParser(
	fileIndex map[string]string,
	basePath string,
) (goldmark.Markdown, *IndexResolver) {

	// Ensure BasePath is clean for usage (not empty)
	if basePath == "" {
		basePath = "/"
	}

	resolver := &IndexResolver{
		Index:    fileIndex,
		Links:    []GraphLink{},
		BasePath: basePath, // Set the path prefix
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
				util.Prioritized(resolver, 1000), // High priority to override default HTML renderer
			),
		),
	)

	return md, resolver
}
