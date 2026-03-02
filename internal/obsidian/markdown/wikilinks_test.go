// @feature:wikilinks Tests for wikilink resolution and nil pointer safety.
package markdown

import (
	"errors"
	"testing"

	"github.com/otaleghani/kiln/internal/obsidian"
)

func newTestResolver(index map[string][]*obsidian.File) *IndexResolver {
	return &IndexResolver{
		Index:     index,
		SourceMap: map[string]string{},
		Links:     []obsidian.GraphLink{},
	}
}

func TestFindFile_NotInIndex(t *testing.T) {
	r := newTestResolver(map[string][]*obsidian.File{})

	file, anchor, err := r.FindFile([]byte("nonexistent"))
	if !errors.Is(err, ErrorCandidateNotFound) {
		t.Fatalf("expected ErrorCandidateNotFound, got %v", err)
	}
	if file != nil {
		t.Errorf("expected nil file, got %+v", file)
	}
	if anchor != "" {
		t.Errorf("expected empty anchor, got %q", anchor)
	}
}

func TestFindFile_NotInIndexWithAnchor(t *testing.T) {
	r := newTestResolver(map[string][]*obsidian.File{})

	file, anchor, err := r.FindFile([]byte("nonexistent#heading"))
	if !errors.Is(err, ErrorCandidateNotFound) {
		t.Fatalf("expected ErrorCandidateNotFound, got %v", err)
	}
	if file != nil {
		t.Errorf("expected nil file, got %+v", file)
	}
	if anchor != "#heading" {
		t.Errorf("expected anchor '#heading', got %q", anchor)
	}
}

func TestFindFile_EmptyCandidateSlice(t *testing.T) {
	r := newTestResolver(map[string][]*obsidian.File{
		"empty": {},
	})

	file, _, err := r.FindFile([]byte("empty"))
	if !errors.Is(err, ErrorCandidateNotFound) {
		t.Fatalf("expected ErrorCandidateNotFound, got %v", err)
	}
	if file != nil {
		t.Errorf("expected nil file, got %+v", file)
	}
}

func TestFindFile_SingleCandidate(t *testing.T) {
	expected := &obsidian.File{
		Name:    "note",
		RelPath: "note.md",
		Path:    "/vault/note.md",
		Ext:     ".md",
		WebPath: "/note",
	}
	r := newTestResolver(map[string][]*obsidian.File{
		"note": {expected},
	})

	file, anchor, err := r.FindFile([]byte("note"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if file != expected {
		t.Errorf("expected %+v, got %+v", expected, file)
	}
	if anchor != "" {
		t.Errorf("expected empty anchor, got %q", anchor)
	}
}

func TestFindFile_PathBasedLookupNoMatch(t *testing.T) {
	candidate := &obsidian.File{
		Name:    "note",
		RelPath: "other/note.md",
		Path:    "/vault/other/note.md",
		Ext:     ".md",
		WebPath: "/other/note",
	}
	r := newTestResolver(map[string][]*obsidian.File{
		"note": {candidate},
	})

	// Path-based lookup where the path doesn't match — should fallback to first candidate
	file, _, err := r.FindFile([]byte("folder/note"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if file != candidate {
		t.Errorf("expected fallback to first candidate, got %+v", file)
	}
}

func TestFindFile_PathBasedLookupMatch(t *testing.T) {
	c1 := &obsidian.File{
		Name:    "note",
		RelPath: "folder/note.md",
		Path:    "/vault/folder/note.md",
		Ext:     ".md",
		WebPath: "/folder/note",
	}
	c2 := &obsidian.File{
		Name:    "note",
		RelPath: "other/note.md",
		Path:    "/vault/other/note.md",
		Ext:     ".md",
		WebPath: "/other/note",
	}
	r := newTestResolver(map[string][]*obsidian.File{
		"note": {c1, c2},
	})

	file, _, err := r.FindFile([]byte("folder/note"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if file != c1 {
		t.Errorf("expected c1, got %+v", file)
	}
}

func TestFindFile_ShortestPathFallback(t *testing.T) {
	long := &obsidian.File{
		Name:    "note",
		RelPath: "a/b/c/note.md",
		Path:    "/vault/a/b/c/note.md",
		Ext:     ".md",
		WebPath: "/a/b/c/note",
	}
	short := &obsidian.File{
		Name:    "note",
		RelPath: "a/note.md",
		Path:    "/vault/a/note.md",
		Ext:     ".md",
		WebPath: "/a/note",
	}
	r := newTestResolver(map[string][]*obsidian.File{
		"note": {long, short},
	})

	file, _, err := r.FindFile([]byte("note"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if file != short {
		t.Errorf("expected shortest path candidate, got %+v", file)
	}
}

func TestFindFile_RootFilePreferred(t *testing.T) {
	nested := &obsidian.File{
		Name:    "note",
		RelPath: "folder/note.md",
		Path:    "/vault/folder/note.md",
		Ext:     ".md",
		WebPath: "/folder/note",
	}
	root := &obsidian.File{
		Name:    "note",
		RelPath: "note.md",
		Path:    "/vault/note.md",
		Ext:     ".md",
		WebPath: "/note",
	}
	r := newTestResolver(map[string][]*obsidian.File{
		"note": {nested, root},
	})

	file, _, err := r.FindFile([]byte("note"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if file != root {
		t.Errorf("expected root file, got %+v", file)
	}
}

func TestFindFile_CaseInsensitiveFallback(t *testing.T) {
	expected := &obsidian.File{
		Name:    "MyNote",
		RelPath: "MyNote.md",
		Path:    "/vault/MyNote.md",
		Ext:     ".md",
		WebPath: "/mynote",
	}
	r := newTestResolver(map[string][]*obsidian.File{
		"mynote": {expected},
	})

	file, _, err := r.FindFile([]byte("MyNote"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if file != expected {
		t.Errorf("expected %+v, got %+v", expected, file)
	}
}

func TestFindFile_StripsMdExtension(t *testing.T) {
	expected := &obsidian.File{
		Name:    "note",
		RelPath: "note.md",
		Path:    "/vault/note.md",
		Ext:     ".md",
		WebPath: "/note",
	}
	r := newTestResolver(map[string][]*obsidian.File{
		"note": {expected},
	})

	file, _, err := r.FindFile([]byte("note.md"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if file != expected {
		t.Errorf("expected %+v, got %+v", expected, file)
	}
}
