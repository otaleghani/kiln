// @feature:search Tests for search index building, markdown stripping, and JSON output.
package search

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/otaleghani/kiln/internal/obsidian"
)

func TestStripMarkdown(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "headings",
			input: "# Title\n## Subtitle\n### Deep",
			want:  "Title Subtitle Deep",
		},
		{
			name:  "bold and italic",
			input: "This is **bold** and *italic* and __also bold__ and _also italic_",
			want:  "This is bold and italic and also bold and also italic",
		},
		{
			name:  "links",
			input: "Click [here](https://example.com) for info",
			want:  "Click here for info",
		},
		{
			name:  "wikilinks with alias",
			input: "See [[target page|alias text]] for details",
			want:  "See alias text for details",
		},
		{
			name:  "wikilinks without alias",
			input: "See [[target page]] for details",
			want:  "See target page for details",
		},
		{
			name:  "code blocks",
			input: "Before\n```go\nfunc main() {}\n```\nAfter",
			want:  "Before After",
		},
		{
			name:  "inline code",
			input: "Use `fmt.Println` to print",
			want:  "Use fmt.Println to print",
		},
		{
			name:  "images",
			input: "Look at ![alt text](image.png) here",
			want:  "Look at here",
		},
		{
			name:  "HTML tags",
			input: "Some <b>bold</b> and <em>emphasis</em> text",
			want:  "Some bold and emphasis text",
		},
		{
			name:  "collapse whitespace",
			input: "Too   many    spaces\n\n\nnewlines",
			want:  "Too many spaces newlines",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stripMarkdown([]byte(tt.input))
			if got != tt.want {
				t.Errorf("stripMarkdown(%q)\ngot:  %q\nwant: %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestBuildIndex(t *testing.T) {
	files := []*obsidian.File{
		{
			Name:    "My Note",
			WebPath: "/notes/my-note",
			Folder:  "notes",
			Content: []byte("# My Note\nThis is **important** content."),
			Tags:    map[string]struct{}{"go": {}, "search": {}},
		},
		{
			Name:    "Another",
			WebPath: "/docs/another",
			Folder:  "docs",
			Content: []byte("Plain text here."),
			Tags:    map[string]struct{}{},
		},
	}

	entries := BuildIndex(files)

	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}

	first := entries[0]
	if first.Title != "My Note" {
		t.Errorf("Title = %q, want %q", first.Title, "My Note")
	}
	if first.URL != "/notes/my-note" {
		t.Errorf("URL = %q, want %q", first.URL, "/notes/my-note")
	}
	if first.Folder != "notes" {
		t.Errorf("Folder = %q, want %q", first.Folder, "notes")
	}
	// Content should be stripped markdown
	if first.Content != "My Note This is important content." {
		t.Errorf("Content = %q, want %q", first.Content, "My Note This is important content.")
	}
	// Tags should contain both keys (order may vary)
	if len(first.Tags) != 2 {
		t.Errorf("expected 2 tags, got %d: %v", len(first.Tags), first.Tags)
	}
	tagSet := map[string]bool{}
	for _, tag := range first.Tags {
		tagSet[tag] = true
	}
	if !tagSet["go"] || !tagSet["search"] {
		t.Errorf("expected tags [go, search], got %v", first.Tags)
	}

	second := entries[1]
	if second.Title != "Another" {
		t.Errorf("Title = %q, want %q", second.Title, "Another")
	}
	if len(second.Tags) != 0 {
		t.Errorf("expected 0 tags, got %d", len(second.Tags))
	}
}

func TestBuildIndex_EmptyContent(t *testing.T) {
	files := []*obsidian.File{
		{
			Name:    "Empty",
			WebPath: "/empty",
			Folder:  "",
			Content: []byte{},
			Tags:    map[string]struct{}{},
		},
	}

	entries := BuildIndex(files)

	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Content != "" {
		t.Errorf("Content = %q, want empty string", entries[0].Content)
	}
	if entries[0].Title != "Empty" {
		t.Errorf("Title = %q, want %q", entries[0].Title, "Empty")
	}
}

func TestWriteIndex(t *testing.T) {
	dir := t.TempDir()

	entries := []SearchEntry{
		{
			Title:   "Test Note",
			URL:     "/test-note",
			Content: "Some content here",
			Tags:    []string{"tag1", "tag2"},
			Folder:  "notes",
		},
		{
			Title:   "Another Note",
			URL:     "/another",
			Content: "More content",
		},
	}

	if err := WriteIndex(entries, dir); err != nil {
		t.Fatalf("WriteIndex failed: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(dir, "search-index.json"))
	if err != nil {
		t.Fatalf("failed to read search-index.json: %v", err)
	}

	var got []SearchEntry
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}

	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(got))
	}

	if got[0].Title != "Test Note" {
		t.Errorf("Title = %q, want %q", got[0].Title, "Test Note")
	}
	if got[0].URL != "/test-note" {
		t.Errorf("URL = %q, want %q", got[0].URL, "/test-note")
	}
	if got[0].Content != "Some content here" {
		t.Errorf("Content = %q, want %q", got[0].Content, "Some content here")
	}
	if len(got[0].Tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(got[0].Tags))
	}
	if got[0].Folder != "notes" {
		t.Errorf("Folder = %q, want %q", got[0].Folder, "notes")
	}

	// Second entry should have empty tags and folder (omitempty)
	if got[1].Tags != nil {
		t.Errorf("expected nil tags for second entry, got %v", got[1].Tags)
	}
	if got[1].Folder != "" {
		t.Errorf("expected empty folder for second entry, got %q", got[1].Folder)
	}
}
