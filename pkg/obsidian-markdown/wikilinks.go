package obsidianmarkdown

import (
	"strings"

	"errors"
	"path/filepath"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"

	"github.com/yuin/goldmark"
	"go.abhg.dev/goldmark/wikilink"
)

// RegisterFuncs registers custom renderers for standard links and images.
// This allows us to intercept [text](url) and ![alt](url) to fix paths.
func (r *IndexResolver) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindLink, r.renderLink)
	reg.Register(ast.KindImage, r.renderImage)
	reg.Register(wikilink.Kind, r.renderWikilink)
}

// ResolveWikilink tries to find the best File instance for the given target.
//
// If it finds a valid candidate and that candidate is a rendered page it adds a link to the GraphLinks.
func (r *IndexResolver) ResolveWikilink(n *wikilink.Node) ([]byte, error) {
	// Reconstruct input
	rawInput := n.Target
	if len(n.Fragment) > 0 {
		rawInput = append(rawInput, '#')
		rawInput = append(rawInput, n.Fragment...)
	}

	// Find the real file path
	file, anchor, err := r.FindFile(rawInput)

	// Default anchor return if not found
	if err != nil {
		return []byte(anchor), nil
	}

	// Only record the link if we are dealing with a content file (optional safety)
	if file.Ext == ".md" || file.Ext == ".canvas" || file.Ext == ".base" {
		r.recordGraphLink(file.WebPath)
	}

	return []byte(file.WebPath + anchor), nil
}

var ErrorCandidateNotFound = errors.New("Candidate not found")

// FindFile returns the best File match for the given target and the anchor
func (r *IndexResolver) FindFile(target []byte) (*File, string, error) {
	dest := string(target)
	dest = strings.TrimSpace(dest)

	// Separate anchor
	anchor := ""
	if idx := strings.Index(dest, "#"); idx != -1 {
		anchor = dest[idx:]
		dest = dest[:idx]
	}

	name, ext := SplitExt(dest)
	// Remove extension
	if ext == ".md" || ext == ".canvas" {
		dest = strings.TrimSuffix(dest, ext)
	}

	// Lookup in Index (case insensitive)
	candidates, ok := r.Index[name]
	// Fallback: Try looking up the full path lowercased (rare, but good for safety)
	if !ok {
		lowerDest := strings.ToLower(dest)
		if alt, ok := r.Index[lowerDest]; ok {
			candidates = alt
		} else {
			return &File{}, anchor, ErrorCandidateNotFound // Not found
		}
	}

	// Select best match
	var bestMatch *File

	// If user provided a path (e.g. "Deployment/Cloudflare Pages")
	if strings.Contains(dest, "/") || strings.Contains(dest, "\\") {
		for _, file := range candidates {
			if strings.Contains(file.RelPath, dest) {
				bestMatch = file
				break
			}
		}

		// Fallback: If contains failed, return first match
		if bestMatch.Name == "" && len(candidates) > 0 {
			bestMatch = candidates[0]
		}

	} else {
		// No path provided (e.g. "Cloudflare Pages"), use standard priority (first root match, then shortest path)

		// Root Match
		for _, file := range candidates {
			if !strings.Contains(file.RelPath, "/") && !strings.Contains(file.RelPath, "\\") {
				bestMatch = file
				break
			}
		}

		// Shortest Path
		if bestMatch == nil {
			shortest := candidates[0]
			for _, file := range candidates {
				if len(file.Path) < len(shortest.Path) {
					shortest = file
				}
			}
			bestMatch = shortest
		}
	}

	return bestMatch, anchor, nil
}

// recordLink adds a directed edge to the graph if the source and target differ.
func (r *IndexResolver) recordGraphLink(target string) {
	if r.CurrentSource != "" && target != "" {
		if !strings.EqualFold(r.CurrentSource, target) {
			r.Links = append(r.Links, GraphLink{
				Source: r.CurrentSource,
				Target: target,
			})
		}
	}
}

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
			"<a href=\"" + destUrl + "\" class=\"markdown-embed-link\" title=\"Open Original\">",
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

// renderLink writes the HTML for standard markdown links.
func (r *IndexResolver) renderLink(
	w util.BufWriter,
	source []byte,
	node ast.Node,
	entering bool,
) (ast.WalkStatus, error) {
	n := node.(*ast.Link)
	if entering {
		newDest := string(n.Destination)

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

// renderImage writes the HTML for standard images
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
	newDest := string(n.Destination)

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

// renderWikilink renders a wikilink and handles text embedding
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

	// Get the webpath (e.g. "/books/harry-potter")
	// This is correct for the browser <a href="..."> or <img src="...">
	destBytes, err := r.ResolveWikilink(n)
	if err != nil {
		return ast.WalkStop, err
	}
	webPath := string(destBytes)

	// If link but NOT embed
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

	// If embed

	// Resolve the real disk path using the SourceMap
	// We strip the hash (#header) from the webPath for the lookup
	cleanWebPath := webPath
	var fragment string
	if idx := strings.Index(webPath, "#"); idx != -1 {
		fragment = webPath[idx:]
		cleanWebPath = webPath[:idx]
	}

	// Default to webPath if not found (fallback)
	realFilePath := cleanWebPath
	if src, ok := r.SourceMap[cleanWebPath]; ok {
		realFilePath = src
	} else {
		return ast.WalkSkipChildren, nil
		// log.Warn("Couldn't find webpath in SourceMap", log.FieldPath, cleanWebPath)
	}

	// Handle images
	// For images, we usually want to render the tag with the web path (for the browser),
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

	// Handle Markdown Notes
	// Read the note to create the embed
	content, err := r.ReadFile(realFilePath)
	if err != nil {
		// log.Warn(
		// 	"Failed to read file",
		// 	log.FieldFile,
		// 	realFilePath,
		// 	"url",
		// 	webPath,
		// 	log.FieldError,
		// 	err,
		// )
		w.WriteString(
			"<a href=\"" + webPath + "\" class=\"broken-embed\">" + string(n.Target) + "</a>",
		)
		return ast.WalkSkipChildren, nil
	}

	if err := r.renderSelection(w, content, fragment, webPath, string(n.Target)); err != nil {
		w.WriteString("")
	}

	return ast.WalkSkipChildren, nil
}

// isImageFile returns true if the given extension is an image extension
func isImageFile(ext string) bool {
	return ext == ".png" || ext == ".jpg" || ext == ".jpeg" ||
		ext == ".gif" || ext == ".svg" || ext == ".webp"
}

// SplitExt returns the base name and the extension of the given path
func SplitExt(path string) (name, ext string) {
	ext = filepath.Ext(path)
	name = strings.TrimSuffix(strings.ToLower(filepath.Base(path)), ext)
	return
}

// GraphLink represents a directed edge in the note graph.
//
// Used to generate the part of the json file
type GraphLink struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

// File rappresents a file that needs to be processed
type File struct {
	Path    string // Complete path of the file
	RelPath string // Relative path from input directory
	Ext     string // Extension of the file
	Name    string // Name of the file (no extension)
	OutPath string // Final output path of the file (e.g. /public/folder/page.html)
	WebPath string // Final web path of the page (e.g. /folder/page)
}

// IndexResolver handles link resolution and graph tracking.
// It intercepts both Wikilinks ([[...]]) and standard Markdown links to ensure
// they point to the correct URL, respecting the site's BasePath.
type IndexResolver struct {
	Index         map[string][]*File                // Lowercase filename -> Candidate files
	SourceMap     map[string]string                 // Webpath -> Real file (used for text embedding)
	Links         []GraphLink                       // All of the graph links
	CurrentSource string                            // The current source
	BasePath      string                            // The basepath // TODO: Delete this, you
	Engine        goldmark.Markdown                 // Needed for recursion
	ReadFile      func(path string) ([]byte, error) // Needed for loading files
}
