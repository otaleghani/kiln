// @feature:linter Tests for broken link detection across wikilinks and markdown links.
package linter

import (
	"bytes"
	"log/slog"
	"os"
	"path/filepath"
	"testing"
)

// setupVault creates a temp directory with the given files and returns the path.
// Each key is a relative path, each value is the file content.
func setupVault(t *testing.T, files map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	for rel, content := range files {
		full := filepath.Join(dir, rel)
		if err := os.MkdirAll(filepath.Dir(full), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

// runBrokenLinks collects log output from BrokenLinks and returns it as a string.
func runBrokenLinks(t *testing.T, dir string) string {
	t.Helper()
	notes := CollectNotes(dir)
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelWarn}))
	BrokenLinks(dir, notes, logger)
	return buf.String()
}

func TestBrokenWikilink(t *testing.T) {
	dir := setupVault(t, map[string]string{
		"note.md": `[[nonexistent]]`,
	})
	out := runBrokenLinks(t, dir)
	if out == "" {
		t.Error("expected warning for broken wikilink [[nonexistent]], got none")
	}
}

func TestValidWikilink(t *testing.T) {
	dir := setupVault(t, map[string]string{
		"note.md":   `[[target]]`,
		"target.md": `hello`,
	})
	out := runBrokenLinks(t, dir)
	if out != "" {
		t.Errorf("expected no warnings for valid wikilink, got: %s", out)
	}
}

func TestBrokenMarkdownLink(t *testing.T) {
	dir := setupVault(t, map[string]string{
		"note.md": `[text](./nonexistent.md)`,
	})
	out := runBrokenLinks(t, dir)
	if out == "" {
		t.Error("expected warning for broken markdown link [text](./nonexistent.md), got none")
	}
}

func TestValidMarkdownLink(t *testing.T) {
	dir := setupVault(t, map[string]string{
		"note.md":     `[text](./existing.md)`,
		"existing.md": `hello`,
	})
	out := runBrokenLinks(t, dir)
	if out != "" {
		t.Errorf("expected no warnings for valid markdown link, got: %s", out)
	}
}

func TestExternalLinksSkipped(t *testing.T) {
	dir := setupVault(t, map[string]string{
		"note.md": `[a](https://example.com) [b](http://example.com) [c](mailto:user@test.com)`,
	})
	out := runBrokenLinks(t, dir)
	if out != "" {
		t.Errorf("expected no warnings for external links, got: %s", out)
	}
}

func TestAnchorOnlyLinkSkipped(t *testing.T) {
	dir := setupVault(t, map[string]string{
		"note.md": `[heading](#some-heading)`,
	})
	out := runBrokenLinks(t, dir)
	if out != "" {
		t.Errorf("expected no warnings for anchor-only link, got: %s", out)
	}
}

func TestMarkdownLinkInSubdirectory(t *testing.T) {
	dir := setupVault(t, map[string]string{
		"sub/note.md":   `[text](./sibling.md)`,
		"sub/sibling.md": `hello`,
	})
	out := runBrokenLinks(t, dir)
	if out != "" {
		t.Errorf("expected no warnings for valid relative link in subdirectory, got: %s", out)
	}
}

func TestMarkdownLinkToParentDirectory(t *testing.T) {
	dir := setupVault(t, map[string]string{
		"sub/note.md": `[text](../root.md)`,
		"root.md":     `hello`,
	})
	out := runBrokenLinks(t, dir)
	if out != "" {
		t.Errorf("expected no warnings for valid parent directory link, got: %s", out)
	}
}

func TestMarkdownLinkWithAnchorValid(t *testing.T) {
	dir := setupVault(t, map[string]string{
		"note.md":   `[text](./target.md#section)`,
		"target.md": `hello`,
	})
	out := runBrokenLinks(t, dir)
	if out != "" {
		t.Errorf("expected no warnings for valid link with anchor, got: %s", out)
	}
}

func TestMarkdownLinkWithAnchorBroken(t *testing.T) {
	dir := setupVault(t, map[string]string{
		"note.md": `[text](./missing.md#section)`,
	})
	out := runBrokenLinks(t, dir)
	if out == "" {
		t.Error("expected warning for broken markdown link with anchor, got none")
	}
}

func TestMixedWikiAndMarkdownLinks(t *testing.T) {
	dir := setupVault(t, map[string]string{
		"note.md":     `[[valid]] [text](./nonexistent.md)`,
		"valid.md":    `hello`,
	})
	out := runBrokenLinks(t, dir)
	if out == "" {
		t.Error("expected warning for broken markdown link in mixed content, got none")
	}
	// The valid wikilink should not be flagged, so only the markdown link warning.
	// We can't be too specific about format, but "nonexistent.md" should appear.
	if !bytes.Contains([]byte(out), []byte("nonexistent.md")) {
		t.Errorf("expected warning to mention nonexistent.md, got: %s", out)
	}
}
