package builder

import (
	"fmt"
	"net/url"
	"path"
	"path/filepath"
	"strings"

	chromaHTML "github.com/alecthomas/chroma/v2/formatters/html"
	mathjax "github.com/litao91/goldmark-mathjax"

	// "github.com/otaleghani/kiln/internal/markdown"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	htmlRenderer "github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
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
	RefMap        map[string][]string // Lowercase Filename -> Real Paths
	SourceMap     map[string]string
	Links         []GraphLink
	CurrentSource string
	BasePath      string

	// TEST: Wikilink text embedding
	Engine   goldmark.Markdown                 // Needed for recursion
	ReadFile func(path string) ([]byte, error) // Needed for loading files
}

func (r *IndexResolver) ResolveWikilink(n *wikilink.Node) ([]byte, error) {
	// 1. Reconstruct Input
	rawInput := n.Target
	if len(n.Fragment) > 0 {
		rawInput = append(rawInput, '#')
		rawInput = append(rawInput, n.Fragment...)
	}

	// 2. Find the REAL file path
	filePath, anchor := r.FindFilePath(rawInput)

	// Default anchor return if not found
	if filePath == "" {
		return []byte(anchor), nil
	}

	// 3. Generate the Web URL (The unique ID)
	ext := filepath.Ext(filePath)
	isContentFile := ext == ".md" || ext == ".canvas"
	pathNoExt := strings.TrimSuffix(filePath, ext)

	// Slugify
	urlPath := slugifyPath(pathNoExt)

	// Keep extensions for assets
	if !isContentFile && ext != "" {
		urlPath += strings.ToLower(ext)
	}

	// Create Final URL
	finalLink := path.Join(r.BasePath, urlPath)
	if isContentFile {
		finalLink = strings.TrimSuffix(finalLink, "/index")
	}

	// --- FIX: Record the Link using the URL ---
	// We strip the anchor so the link points to the Page, not the Header
	graphTarget := finalLink

	// Ensure we don't record self-links if you don't want them (optional)
	// r.recordLink(graphTarget)

	// DEBUG PRINT
	// fmt.Printf("DEBUG LINK: Source='%s' Target='%s'\n", r.CurrentSource, graphTarget)

	// Only record if we are dealing with a content file (optional safety)
	// if isContentFile {
	// 	r.recordLink(graphTarget)
	// }
	// -----------------------------------------

	// TODO: Delete all of the isContentFile logic, because we are only running this in .md and .canvas files
	// anyways AND we are passing to the resolver the final link
	r.recordLink(graphTarget)

	return []byte(finalLink + anchor), nil
}

// Simple helper
func slugifyPath(p string) string {
	// Split by slash, slugify each component, join back
	parts := strings.Split(p, "/")
	for i, part := range parts {
		parts[i] = slugify(part) // Assuming you have the slugify function from your original code
	}
	return strings.Join(parts, "/")
}

// RegisterFuncs registers custom renderers for standard links and images.
// This allows us to intercept [text](url) and ![alt](url) to fix paths.
func (r *IndexResolver) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindLink, r.renderLink)
	reg.Register(ast.KindImage, r.renderImage)
	reg.Register(wikilink.Kind, r.renderWikilink)
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

func (r *IndexResolver) renderImage(
	w util.BufWriter,
	source []byte,
	node ast.Node,
	entering bool,
) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkSkipChildren, nil
	}
	n := node.(*ast.Image)
	newDest := r.resolvePath(string(n.Destination))

	// Use the helper
	r.writeImage(w, newDest, n.Text(source), n.Title)

	return ast.WalkSkipChildren, nil
}

// writeImage is a helper to ensure consistent image rendering for both
// standard markdown images and Obsidian wikilink images.
func (r *IndexResolver) writeImage(w util.BufWriter, src string, alt []byte, title []byte) {
	w.WriteString("<img src=\"")
	w.WriteString(src)
	w.WriteString("\" alt=\"")
	w.Write(alt)
	w.WriteString("\"")
	if len(title) > 0 {
		w.WriteString(" title=\"")
		w.Write(title)
		w.WriteString("\"")
	}
	w.WriteString(">")
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

	if dest == "index" {
		return ""
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
	// Update this type to match the output of initBuild
	fileIndex map[string][]string,
	sourceMap map[string]string,
	basePath string,
	loader func(path string) ([]byte, error),
) (goldmark.Markdown, *IndexResolver) {

	if basePath == "" {
		basePath = "/"
	}
	// --- 1. BUILD THE REFERENCE MAP ---
	// Maps "credits" -> ["content/Credits.md"]
	refMap := make(map[string][]string)

	for filename, paths := range fileIndex {
		// Convert "Credits" -> "credits"
		lowerName := strings.ToLower(filename)

		// If multiple files have same name (diff casing), append them
		refMap[lowerName] = append(refMap[lowerName], paths...)
	}

	resolver := &IndexResolver{
		Index:     fileIndex,
		RefMap:    refMap,
		SourceMap: sourceMap,
		Links:     []GraphLink{},
		BasePath:  basePath,
		ReadFile:  loader,
	}

	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			meta.Meta,
			&wikilink.Extender{Resolver: resolver}, // Hook 1: Wikilinks
			// markdown.Comments,
			// markdown.Mermaid,
			// markdown.Callouts,
			extension.Footnote,
			// Highlights,                             // Handles ==text==
			// Mermaid,                                // Handles Mermaid
			// Callouts,                               // Handles > [!type] blocks
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
				// util.Prioritized(resolver, 1000),
				util.Prioritized(resolver, 0),
			),
		),
	)

	resolver.Engine = md

	return md, resolver
}

// TEST
func (r *IndexResolver) renderWikilink(
	w util.BufWriter,
	source []byte,
	node ast.Node,
	entering bool,
) (ast.WalkStatus, error) {
	n := node.(*wikilink.Node)

	if !entering {
		return ast.WalkContinue, nil
	}

	// 1. Get the Web Path (e.g. "/books/harry-potter")
	// This is correct for the browser <a href="..."> or <img src="...">
	destBytes, err := r.ResolveWikilink(n)
	if err != nil {
		return ast.WalkStop, err
	}
	webPath := string(destBytes)

	// --- IF LINK (NOT EMBED) ---
	// Check for the "!" hack on previous sibling
	isEmbed := n.Embed
	if !isEmbed {
		prev := n.PreviousSibling()
		if prev != nil && prev.Kind() == ast.KindText {
			textBytes := prev.Text(source)
			if len(textBytes) > 0 && textBytes[len(textBytes)-1] == '!' {
				isEmbed = true
				segment := prev.Lines().At(prev.Lines().Len() - 1)
				prev.Lines().
					Set(prev.Lines().Len()-1, text.NewSegment(segment.Start, segment.Stop-1))
			}
		}
	}

	if !isEmbed {
		w.WriteString("<a href=\"")
		w.WriteString(webPath)
		w.WriteString("\" class=\"internal-link\">")
		w.Write(n.Target)
		w.WriteString("</a>")
		return ast.WalkSkipChildren, nil
	}

	// --- IF EMBED (Transclusion) ---

	// 2. Resolve the REAL Disk Path using the SourceMap
	// We strip the hash (#header) from the webPath for the lookup
	cleanWebPath := webPath
	var fragment string
	if idx := strings.Index(webPath, "#"); idx != -1 {
		fragment = webPath[idx:]
		cleanWebPath = webPath[:idx]
	}

	// KEY FIX: Look up the real file path!
	// Default to webPath if not found (fallback)
	realFilePath := cleanWebPath
	if src, ok := r.SourceMap[cleanWebPath]; ok {
		realFilePath = src
	} else {
		// Try adding/removing BasePath if lookup failed
		// (Depends on how you populate SourceMap vs BasePath)
		// altKey := path.Join(r.BasePath, cleanWebPath)
		// if src, ok := r.SourceMap[altKey]; ok {
		// 	realFilePath = src
		// }
		u, _ := url.Parse(BaseURL)
		basePath := u.Path
		if src, ok := r.SourceMap[strings.TrimPrefix(realFilePath, basePath)]; ok {
			realFilePath = src
		} else {
			fmt.Println("Error trimming prefix")
		}
	}
	// Determine the base path for link resolution (e.g., "/kiln" from "https://domain.com/kiln")

	// fmt.Printf("DEBUG: Web: %s | Disk: %s\n", cleanWebPath, realFilePath)

	// 3. Handle Images
	// For images, we usually want to render the tag with the WEB PATH (for the browser),
	// NOT the file path.
	ext := strings.ToLower(filepath.Ext(realFilePath))
	if isImageFile(ext) {
		// Extract Alt Text: In Wikilinks ([[Image.png|Alt Text]]),
		// the text after the pipe is stored as the node's children.
		var altText []byte
		if n.HasChildren() {
			// Iterate over children to reconstruct the label/alt text
			for child := n.FirstChild(); child != nil; child = child.NextSibling() {
				// We assume text nodes here; simpler than full recursion
				if child.Kind() == ast.KindText {
					altText = append(altText, child.Text(source)...)
				}
			}
		} else {
			// If no alt text is provided, some people prefer using the filename,
			// others prefer empty. Obsidian uses the filename often.
			altText = n.Target
		}

		// Note: Wikilinks don't officially support "Title" attributes (tooltip),
		// so we pass nil. If you want title, you'd have to parse it from the alt text manually.
		r.writeImage(w, webPath, altText, nil)
		return ast.WalkSkipChildren, nil
	}

	// 4. Handle Markdown Notes
	// Now we use realFilePath to read from disk
	content, err := r.ReadFile(realFilePath)
	if err != nil {
		fmt.Printf("ERROR: ReadFile failed for '%s' -> '%s' (%v)\n", webPath, realFilePath, err)
		w.WriteString(
			"<a href=\"" + webPath + "\" class=\"broken-embed\">" + string(n.Target) + "</a>",
		)
		return ast.WalkSkipChildren, nil
	}

	// w.WriteString(
	// 	"<a href=\"" + webPath + "\" class=\"\">" + string(n.Target) + "</a>",
	// )
	if err := r.renderSelection(w, content, fragment, webPath, string(n.Target)); err != nil {
		w.WriteString("")
	}
	// if err := r.renderSelection(w, content, fragment); err != nil {
	// 	w.WriteString("")
	// }

	return ast.WalkSkipChildren, nil
}

func isImageFile(ext string) bool {
	return ext == ".png" || ext == ".jpg" || ext == ".jpeg" ||
		ext == ".gif" || ext == ".svg" || ext == ".webp"
}

// Update the function signature
func (r *IndexResolver) renderSelection(
	w util.BufWriter,
	source []byte,
	fragment string,
	destUrl string, // New: The href for the links
	displayName string, // New: The text to display (n.Target)
) error {
	reader := text.NewReader(source)
	doc := r.Engine.Parser().Parse(reader)

	nodes := r.extractNodes(doc, source, fragment)

	if len(nodes) > 0 {
		w.WriteString("<div class=\"markdown-embed\">")

		// --- EMBED HEADER (Title + Button) ---
		w.WriteString("<div class=\"markdown-embed-header\">")

		// 1. Left Side: Title or Link
		if fragment == "" {
			// Case A: Whole Page -> Display Link to Page
			w.WriteString("<a href=\"" + destUrl + "\" class=\"markdown-embed-title\">")
			w.WriteString(displayName)
			w.WriteString("</a>")
		} else {
			// Case B: Fragment -> Display Section Name (or displayName > fragment)
			w.WriteString("<div class=\"markdown-embed-title\">")
			w.WriteString(strings.TrimPrefix(fragment, "#"))
			w.WriteString("</div>")
		}

		// 2. Right Side: Maximize Button (Always visible)
		// Links to the specific anchor or page
		w.WriteString(
			"<a href=\"" + slugify(
				destUrl,
			) + "\" class=\"markdown-embed-link\" title=\"Open Original\">",
		)
		w.WriteString("<i class=\"\" data-lucide=\"maximize-2\"></i>")
		w.WriteString("</a>")

		w.WriteString("</div>") // End Header
		// -------------------------------------

		w.WriteString("<div class=\"markdown-embed-content\">")

		renderer := r.Engine.Renderer()
		for _, n := range nodes {
			if err := renderer.Render(w, source, n); err != nil {
				return err
			}
		}

		w.WriteString("</div></div>")
	}
	return nil
}

// extractNodes finds the specific Heading or BlockID within the AST.
func (r *IndexResolver) extractNodes(root ast.Node, source []byte, fragment string) []ast.Node {
	if fragment == "" {
		return []ast.Node{root}
	}

	fragment = strings.TrimPrefix(fragment, "#")

	// 1. Block ID Lookup (^blockid)
	if strings.HasPrefix(fragment, "^") {
		targetID := fragment[1:]
		var found ast.Node

		ast.Walk(root, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
			if !entering || n.Type() != ast.TypeBlock {
				return ast.WalkContinue, nil
			}
			// Look for ^id in the raw text of the block
			raw := string(n.Text(source))
			if strings.Contains(raw, "^"+targetID) {
				found = n
				return ast.WalkStop, nil
			}
			return ast.WalkContinue, nil
		})

		if found != nil {
			return []ast.Node{found}
		}
		return nil
	}

	// 2. Heading Lookup (Header Text)
	var startNode ast.Node
	ast.Walk(root, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering || n.Kind() != ast.KindHeading {
			return ast.WalkContinue, nil
		}
		h := n.(*ast.Heading)

		// Match by text content (case insensitive)
		// Note: A more robust approach handles slugification matching here
		headerText := string(h.Text(source))
		if strings.EqualFold(headerText, fragment) {
			startNode = n
			return ast.WalkStop, nil
		}
		return ast.WalkContinue, nil
	})

	if startNode != nil {
		results := []ast.Node{startNode}
		startLevel := startNode.(*ast.Heading).Level

		// Collect siblings until next header of same level
		curr := startNode.NextSibling()
		for curr != nil {
			if h, ok := curr.(*ast.Heading); ok {
				if h.Level <= startLevel {
					break
				}
			}
			results = append(results, curr)
			curr = curr.NextSibling()
		}
		return results
	}

	return nil
}

func (r *IndexResolver) FindFilePath(target []byte) (string, string) {
	dest := string(target)
	dest = strings.TrimSpace(dest)

	// 1. Separate Anchor
	anchor := ""
	if idx := strings.Index(dest, "#"); idx != -1 {
		anchor = dest[idx:]
		dest = dest[:idx]
	}

	// 2. Prepare for Lookup
	// User Input: "Deployment/Cloudflare Pages"

	// Remove extension
	ext := filepath.Ext(dest)
	if ext == ".md" || ext == ".canvas" {
		dest = strings.TrimSuffix(dest, ext)
	}

	// 3. Extract Key (Filename) & Lowercase it
	// "Deployment/Cloudflare Pages" -> base: "Cloudflare Pages" -> lower: "cloudflare pages"
	baseName := filepath.Base(dest)
	lowerKey := strings.ToLower(baseName)

	// 4. Lookup in RefMap (Case Insensitive!)
	candidates, ok := r.RefMap[lowerKey]

	// Fallback: Try looking up the full path lowercased (rare, but good for safety)
	if !ok {
		lowerDest := strings.ToLower(dest)
		if alt, ok := r.RefMap[lowerDest]; ok {
			candidates = alt
		} else {
			return "", anchor // Not found
		}
	}

	// 5. Select Best Match (Path Filtering)
	var bestMatch string

	// If user provided a path (e.g. "Deployment/Cloudflare Pages")
	if strings.Contains(dest, "/") || strings.Contains(dest, "\\") {
		// Normalize input path for comparison
		searchSuffix := strings.ToLower(filepath.ToSlash(dest))

		// Ensure suffix expects an extension if candidates have one
		if len(candidates) > 0 && filepath.Ext(candidates[0]) != "" &&
			filepath.Ext(searchSuffix) == "" {
			searchSuffix += ".md"
		}

		for _, pathStr := range candidates {
			// Compare lowercase suffixes
			// Candidate: "content/Deployment/Cloudflare Pages.md"
			// Suffix:    "deployment/cloudflare pages.md"
			normPath := strings.ToLower(filepath.ToSlash(pathStr))

			if strings.HasSuffix(normPath, searchSuffix) {
				bestMatch = pathStr
				break
			}
		}

		// Fallback: If strict suffix failed, return first candidate
		if bestMatch == "" && len(candidates) > 0 {
			bestMatch = candidates[0]
		}

	} else {
		// No path provided (e.g. "Cloudflare Pages"), use standard priority

		// Priority A: Root Match
		for _, pathStr := range candidates {
			if !strings.Contains(pathStr, "/") && !strings.Contains(pathStr, "\\") {
				bestMatch = pathStr
				break
			}
		}
		// Priority B: Shortest Path
		if bestMatch == "" {
			shortest := candidates[0]
			for _, pathStr := range candidates {
				if len(pathStr) < len(shortest) {
					shortest = pathStr
				}
			}
			bestMatch = shortest
		}
	}

	return bestMatch, anchor
}
