// @feature:watch Tests for dependency graph.
package watch

import (
	"slices"
	"sort"
	"testing"

	"github.com/otaleghani/kiln/internal/obsidian"
)

func TestNewDepGraph(t *testing.T) {
	g := NewDepGraph()

	if g.Forward == nil {
		t.Fatal("Forward map should not be nil")
	}
	if g.Reverse == nil {
		t.Fatal("Reverse map should not be nil")
	}
	if got := g.Dependents("anything"); len(got) != 0 {
		t.Errorf("expected no dependents on empty graph, got %v", got)
	}
}

func TestAddEdgeAndDependents(t *testing.T) {
	g := NewDepGraph()
	g.AddEdge("a.md", "b")
	g.AddEdge("c.md", "b")

	got := g.Dependents("b")
	sort.Strings(got)

	want := []string{"a.md", "c.md"}
	if !slices.Equal(got, want) {
		t.Errorf("Dependents(b) = %v, want %v", got, want)
	}
}

func TestAddEdgeDuplicate(t *testing.T) {
	g := NewDepGraph()
	g.AddEdge("a.md", "b")
	g.AddEdge("a.md", "b")

	got := g.Dependents("b")
	if len(got) != 1 {
		t.Errorf("expected 1 dependent after duplicate add, got %v", got)
	}
}

func TestRemoveSource(t *testing.T) {
	g := NewDepGraph()
	g.AddEdge("a.md", "b")
	g.AddEdge("c.md", "b")
	g.AddEdge("a.md", "d")

	g.RemoveSource("a.md")

	gotB := g.Dependents("b")
	if !slices.Equal(gotB, []string{"c.md"}) {
		t.Errorf("Dependents(b) after remove = %v, want [c.md]", gotB)
	}

	gotD := g.Dependents("d")
	if len(gotD) != 0 {
		t.Errorf("Dependents(d) after remove = %v, want []", gotD)
	}

	if _, ok := g.Forward["a.md"]; ok {
		t.Error("Forward[a.md] should be deleted after RemoveSource")
	}
}

func TestRemoveSourceNonexistent(t *testing.T) {
	g := NewDepGraph()
	g.RemoveSource("nonexistent.md") // should not panic
}

func TestBuildFromFiles(t *testing.T) {
	files := []*obsidian.File{
		{
			RelPath: "notes/a.md",
			Ext:     ".md",
			Name:    "a",
			Links:   []string{"[[NoteB]]", "[[NoteB|alias]]", "[[deep#heading]]"},
		},
		{
			RelPath: "notes/b.md",
			Ext:     ".md",
			Name:    "b",
			Links:   []string{"[text](./notec.md)", "[other](https://example.com)"},
		},
		{
			RelPath: "images/photo.png",
			Ext:     ".png",
			Name:    "photo",
			Links:   []string{},
		},
		{
			RelPath: "notes/c.md",
			Ext:     ".md",
			Name:    "c",
			Links:   []string{"[[folder/Sub Page]]", "[ref](../other/page.md#section)"},
		},
	}

	g := NewDepGraph()
	g.BuildFromFiles(files)

	// a.md links to noteb (twice via alias, should deduplicate) and deep
	gotNoteB := g.Dependents("noteb")
	sort.Strings(gotNoteB)
	if !slices.Equal(gotNoteB, []string{"notes/a.md"}) {
		t.Errorf("Dependents(noteb) = %v, want [notes/a.md]", gotNoteB)
	}

	gotDeep := g.Dependents("deep")
	if !slices.Equal(gotDeep, []string{"notes/a.md"}) {
		t.Errorf("Dependents(deep) = %v, want [notes/a.md]", gotDeep)
	}

	// b.md links to notec (markdown link), external skipped
	gotNoteC := g.Dependents("notec")
	if !slices.Equal(gotNoteC, []string{"notes/b.md"}) {
		t.Errorf("Dependents(notec) = %v, want [notes/b.md]", gotNoteC)
	}

	// c.md links to "sub page" (base of folder/Sub Page, lowercased) and "page" (base of ../other/page.md)
	gotSubPage := g.Dependents("sub page")
	if !slices.Equal(gotSubPage, []string{"notes/c.md"}) {
		t.Errorf("Dependents(sub page) = %v, want [notes/c.md]", gotSubPage)
	}

	gotPage := g.Dependents("page")
	if !slices.Equal(gotPage, []string{"notes/c.md"}) {
		t.Errorf("Dependents(page) = %v, want [notes/c.md]", gotPage)
	}

	// photo.png is not .md, should not be processed
	if len(g.Forward) != 3 {
		t.Errorf("expected 3 sources in Forward, got %d", len(g.Forward))
	}
}

func TestUpdateFiles(t *testing.T) {
	g := NewDepGraph()
	g.AddEdge("a.md", "b")

	files := []*obsidian.File{
		{RelPath: "a.md", Ext: ".md", Links: []string{"[[c]]"}},
	}
	g.UpdateFiles(files)

	if deps := g.Dependents("b"); len(deps) != 0 {
		t.Errorf("expected no dependents for b, got %v", deps)
	}
	if deps := g.Dependents("c"); len(deps) != 1 || deps[0] != "a.md" {
		t.Errorf("expected [a.md] for c, got %v", deps)
	}
}

func TestBuildFromFilesMailtoSkipped(t *testing.T) {
	files := []*obsidian.File{
		{
			RelPath: "note.md",
			Ext:     ".md",
			Name:    "note",
			Links:   []string{"[email](mailto:user@example.com)"},
		},
	}

	g := NewDepGraph()
	g.BuildFromFiles(files)

	if len(g.Forward["note.md"]) != 0 {
		t.Errorf("expected no forward edges for mailto link, got %v", g.Forward["note.md"])
	}
}
