package obsidianmarkdown

import (
	"bytes"
	"html/template"
	"os"

	chromaHTML "github.com/alecthomas/chroma/v2/formatters/html"
	mathjax "github.com/litao91/goldmark-mathjax"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
	"go.abhg.dev/goldmark/wikilink"
)

// newMarkdownParser creates a Goldmark instance configured for Obsidian compatibility.
// It enables GFM, MathJax, syntax highlighting, and custom link resolution.
func New(
	fileIndex map[string][]*File,
	loader func(path string) ([]byte, error),
) *ObsidianMarkdown {
	resolver := &IndexResolver{
		Index:    fileIndex,
		Links:    []GraphLink{},
		ReadFile: loader,
	}
	sourceMap := make(map[string]string)
	for _, candidate := range fileIndex {
		for _, file := range candidate {
			sourceMap[file.WebPath] = file.RelPath
		}
	}
	resolver.SourceMap = sourceMap

	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Footnote,
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
			html.WithUnsafe(),
			html.WithHardWraps(),
			// Standard links and images
			renderer.WithNodeRenderers(
				// util.Prioritized(resolver, 1000),
				util.Prioritized(resolver, 0),
			),
		),
	)

	resolver.Engine = md

	return &ObsidianMarkdown{markdown: md, Resolver: resolver}
}

// Render processes the given path
func (o *ObsidianMarkdown) Render(file File) (NoteData, error) {
	// Sets the current source
	o.Resolver.CurrentSource = file.WebPath

	// Read the file source
	source, err := os.ReadFile(file.Path)
	if err != nil {
		return NoteData{}, err
	}

	// We separate Parsing and Rendering to inspect the AST (Abstract Syntax Tree)
	// This allows us to extract metadata (Frontmatter) and generate the TOC before rendering HTML.
	ctx := parser.NewContext()
	doc := o.markdown.Parser().Parse(text.NewReader(source), parser.WithContext(ctx))

	// Generate Table of Contents from the AST
	tocHTML := extractTOC(doc, source)

	// Render the actual HTML body
	var buf bytes.Buffer
	if err := o.markdown.Renderer().Render(&buf, source, doc); err != nil {
		return NoteData{}, err
	}

	// Apply post-processing transformers
	finalHTML := buf.String()
	finalHTML = applyTransforms(finalHTML)

	return NoteData{Content: finalHTML, TOC: tocHTML, Frontmatter: meta.Get(ctx)}, nil
}

// RenderNote processes the given content
func (o *ObsidianMarkdown) RenderNote(content []byte) (string, error) {
	var buf bytes.Buffer
	err := o.markdown.Convert(content, &buf)
	if err != nil {
		return "", err
	}

	htmlContent := buf.String()
	htmlContent = applyTransforms(htmlContent)

	return htmlContent, nil
}

// ObsidianMarkdown is the main entrypoint
type ObsidianMarkdown struct {
	markdown goldmark.Markdown
	Resolver *IndexResolver
}

// NoteData holds the data of the rendered note, like the body HTML and the table of contents
type NoteData struct {
	Content     string         // Content of the note
	TOC         template.HTML  // Table of contents
	Frontmatter map[string]any // Extracted frontmatter
}
