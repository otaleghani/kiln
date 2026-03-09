// @feature:markdown Integration tests for the full rendering pipeline with markdown links and wikilinks.
package markdown

import (
	"strings"
	"testing"

	"github.com/otaleghani/kiln/internal/obsidian"
)

func newTestMarkdown(files ...*obsidian.File) *ObsidianMarkdown {
	index := make(map[string][]*obsidian.File)
	for _, f := range files {
		index[f.Name] = append(index[f.Name], f)
	}
	loader := func(path string) ([]byte, error) {
		return nil, nil
	}
	return New(index, loader)
}

func TestRenderNote_Wikilink(t *testing.T) {
	note := &obsidian.File{
		Name:    "Note",
		RelPath: "Note.md",
		Path:    "/vault/Note.md",
		Ext:     ".md",
		WebPath: "/note",
	}
	md := newTestMarkdown(note)
	md.Resolver.CurrentSource = "/current"

	html, err := md.RenderNote([]byte("[[Note]]"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(html, `<a href="/note"`) {
		t.Errorf("expected href to /note, got: %s", html)
	}
	if !strings.Contains(html, `class="internal-link"`) {
		t.Errorf("expected internal-link class, got: %s", html)
	}
	if !strings.Contains(html, ">Note</a>") {
		t.Errorf("expected link text 'Note', got: %s", html)
	}
}

func TestRenderNote_MarkdownLinkRelative(t *testing.T) {
	md := newTestMarkdown()
	md.Resolver.CurrentSource = "/current-dir/page"

	html, err := md.RenderNote([]byte("[text](./note.md)"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(html, `<a href="/current-dir/note"`) {
		t.Errorf("expected resolved relative path, got: %s", html)
	}
	if !strings.Contains(html, ">text</a>") {
		t.Errorf("expected link text 'text', got: %s", html)
	}
}

func TestRenderNote_MarkdownLinkExternal(t *testing.T) {
	md := newTestMarkdown()
	md.Resolver.CurrentSource = "/page"

	html, err := md.RenderNote([]byte("[text](https://example.com)"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(html, `<a href="https://example.com"`) {
		t.Errorf("expected external URL untouched, got: %s", html)
	}
	if !strings.Contains(html, ">text</a>") {
		t.Errorf("expected link text 'text', got: %s", html)
	}
}

func TestRenderNote_MarkdownImage(t *testing.T) {
	md := newTestMarkdown()
	md.Resolver.CurrentSource = "/docs/page"

	html, err := md.RenderNote([]byte("![alt](./image.png)"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(html, `src="/docs/image.png"`) {
		t.Errorf("expected resolved image path, got: %s", html)
	}
	if !strings.Contains(html, `alt="alt"`) {
		t.Errorf("expected alt text, got: %s", html)
	}
}

func TestRenderNote_MixedWikilinkAndMarkdownLink(t *testing.T) {
	note := &obsidian.File{
		Name:    "Target",
		RelPath: "Target.md",
		Path:    "/vault/Target.md",
		Ext:     ".md",
		WebPath: "/target",
	}
	md := newTestMarkdown(note)
	md.Resolver.CurrentSource = "/folder/source"

	input := "See [[Target]] and also [other](./other.md) for details."
	html, err := md.RenderNote([]byte(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(html, `<a href="/target" class="internal-link">Target</a>`) {
		t.Errorf("expected wikilink rendered, got: %s", html)
	}
	if !strings.Contains(html, `<a href="/folder/other">other</a>`) {
		t.Errorf("expected markdown link rendered, got: %s", html)
	}
}

func TestRenderNote_GraphLinksCollected(t *testing.T) {
	note := &obsidian.File{
		Name:    "Target",
		RelPath: "Target.md",
		Path:    "/vault/Target.md",
		Ext:     ".md",
		WebPath: "/target",
	}
	md := newTestMarkdown(note)
	md.Resolver.CurrentSource = "/source"

	input := "Link to [[Target]] and [md link](./sibling.md)."
	_, err := md.RenderNote([]byte(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	links := md.Resolver.Links
	if len(links) != 2 {
		t.Fatalf("expected 2 graph links, got %d: %+v", len(links), links)
	}

	targets := map[string]bool{}
	for _, link := range links {
		if link.Source != "/source" {
			t.Errorf("expected source '/source', got %q", link.Source)
		}
		targets[link.Target] = true
	}

	if !targets["/target"] {
		t.Errorf("expected graph link to '/target', got links: %+v", links)
	}
	if !targets["/sibling"] {
		t.Errorf("expected graph link to '/sibling', got links: %+v", links)
	}
}
