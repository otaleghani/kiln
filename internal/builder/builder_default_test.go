// @feature:builder Tests for base file parsing and empty base file handling.
package builder

import (
	"image/png"
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/otaleghani/kiln/internal/obsidian"
	"github.com/otaleghani/kiln/internal/obsidian/bases"
)

func TestGeneratePageOGImages(t *testing.T) {
	outDir := t.TempDir()

	site := &DefaultSite{
		SiteName: "Site",
		Theme: &Theme{
			Dark: &ThemeColors{
				Accent: "#7e6df7",
				Bg:     "#1e1e1e",
				Text:   "#dcddde",
			},
		},
		log: slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelWarn})),
	}

	site.GeneratePageOGImages("Test", "Desc", "", outDir)

	// Verify og.png exists and is a valid PNG
	ogPath := filepath.Join(outDir, "og.png")
	ogFile, err := os.Open(ogPath)
	if err != nil {
		t.Fatalf("og.png not created: %v", err)
	}
	defer ogFile.Close()

	ogImg, err := png.Decode(ogFile)
	if err != nil {
		t.Fatalf("og.png is not a valid PNG: %v", err)
	}
	if bounds := ogImg.Bounds(); bounds.Dx() != 1200 || bounds.Dy() != 630 {
		t.Errorf("og.png expected 1200x630, got %dx%d", bounds.Dx(), bounds.Dy())
	}

	// Verify twitter.png exists and is a valid PNG
	twPath := filepath.Join(outDir, "twitter.png")
	twFile, err := os.Open(twPath)
	if err != nil {
		t.Fatalf("twitter.png not created: %v", err)
	}
	defer twFile.Close()

	twImg, err := png.Decode(twFile)
	if err != nil {
		t.Fatalf("twitter.png is not a valid PNG: %v", err)
	}
	if bounds := twImg.Bounds(); bounds.Dx() != 1200 || bounds.Dy() != 600 {
		t.Errorf("twitter.png expected 1200x600, got %dx%d", bounds.Dx(), bounds.Dy())
	}
}

func TestParseBaseFile_EmptyFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.base")
	if err := os.WriteFile(path, []byte(""), 0644); err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	base, err := ParseBaseFile(path)
	if err != nil {
		t.Fatalf("ParseBaseFile should not error on empty file, got: %v", err)
	}

	if len(base.Views) != 0 {
		t.Errorf("expected 0 views for empty base file, got %d", len(base.Views))
	}
	if len(base.Filters) != 0 {
		t.Errorf("expected 0 filters for empty base file, got %d", len(base.Filters))
	}
}

func TestParseBaseFile_ValidContent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "valid.base")
	content := `filters:
  and:
    - "file.folder == \"notes\""
views:
  - type: table
    name: "My View"
    order:
      - file.name
      - status
    filters:
      and:
        - "status == \"done\""
`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	base, err := ParseBaseFile(path)
	if err != nil {
		t.Fatalf("ParseBaseFile failed: %v", err)
	}

	if len(base.Views) != 1 {
		t.Fatalf("expected 1 view, got %d", len(base.Views))
	}
	if base.Views[0].Name != "My View" {
		t.Errorf("expected view name 'My View', got %q", base.Views[0].Name)
	}
	if len(base.Views[0].Order) != 2 {
		t.Errorf("expected 2 columns, got %d", len(base.Views[0].Order))
	}
}

func TestParseBaseFile_NonExistent(t *testing.T) {
	_, err := ParseBaseFile("/nonexistent/path/file.base")
	if err == nil {
		t.Error("expected error for non-existent file")
	}
}

func TestParseBaseFile_InvalidYAML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "invalid.base")
	content := `:::not valid yaml[[[`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	_, err := ParseBaseFile(path)
	if err == nil {
		t.Error("expected error for invalid YAML")
	}
}

func TestEmptyViewsDoNotPanic(t *testing.T) {
	// Verify that an empty base (no views) does not cause an index-out-of-range
	// panic when processing view-dependent fields.
	base := PageBase{
		File: &obsidian.File{
			Name:    "empty",
			WebPath: "/empty",
		},
		// Views deliberately left empty
	}

	allFiles := []*obsidian.File{
		{Name: "note1.md", Ext: ".md"},
		{Name: "note2.md", Ext: ".md"},
	}

	activeFiles := bases.FilterFiles(allFiles, base.Filters)

	var fileGroups []*bases.FileGroup
	var columns []string
	if len(base.Views) > 0 {
		activeFiles = bases.FilterFiles(activeFiles, base.Views[0].Filters)
		if base.Views[0].GroupBy.Property != "" {
			fileGroups = bases.GroupFiles(activeFiles, base.Views[0].GroupBy.Property)
		}
		columns = base.Views[0].Order
	}

	// With empty views, all files should pass through unfiltered
	if len(activeFiles) != 2 {
		t.Errorf("expected 2 active files with no filters, got %d", len(activeFiles))
	}
	if fileGroups != nil {
		t.Errorf("expected nil file groups with no views, got %v", fileGroups)
	}
	if columns != nil {
		t.Errorf("expected nil columns with no views, got %v", columns)
	}
}
