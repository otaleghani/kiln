// @feature:watch Tests for changeset computation.
package watch

import (
	"slices"
	"sort"
	"testing"
)

func TestComputeChangeSetSimple(t *testing.T) {
	g := NewDepGraph()

	cs := ComputeChangeSet([]string{"notes/hello.md"}, nil, g)

	if !slices.Equal(cs.Rebuild, []string{"notes/hello.md"}) {
		t.Errorf("Rebuild = %v, want [notes/hello.md]", cs.Rebuild)
	}
	if len(cs.Remove) != 0 {
		t.Errorf("Remove = %v, want []", cs.Remove)
	}
}

func TestComputeChangeSetWithDependents(t *testing.T) {
	g := NewDepGraph()
	g.AddEdge("notes/b.md", "a")

	cs := ComputeChangeSet([]string{"notes/a.md"}, nil, g)

	sort.Strings(cs.Rebuild)
	want := []string{"notes/a.md", "notes/b.md"}
	if !slices.Equal(cs.Rebuild, want) {
		t.Errorf("Rebuild = %v, want %v", cs.Rebuild, want)
	}
	if len(cs.Remove) != 0 {
		t.Errorf("Remove = %v, want []", cs.Remove)
	}
}

func TestComputeChangeSetRemoved(t *testing.T) {
	g := NewDepGraph()
	g.AddEdge("notes/b.md", "a")

	cs := ComputeChangeSet(nil, []string{"notes/a.md"}, g)

	if !slices.Equal(cs.Rebuild, []string{"notes/b.md"}) {
		t.Errorf("Rebuild = %v, want [notes/b.md]", cs.Rebuild)
	}
	if !slices.Equal(cs.Remove, []string{"notes/a.md"}) {
		t.Errorf("Remove = %v, want [notes/a.md]", cs.Remove)
	}
}

func TestComputeChangeSetDedup(t *testing.T) {
	g := NewDepGraph()
	// b depends on both a and c
	g.AddEdge("notes/b.md", "a")
	g.AddEdge("notes/b.md", "c")

	// Both a and c changed — b should appear only once
	cs := ComputeChangeSet([]string{"notes/a.md", "notes/c.md"}, nil, g)

	sort.Strings(cs.Rebuild)
	want := []string{"notes/a.md", "notes/b.md", "notes/c.md"}
	if !slices.Equal(cs.Rebuild, want) {
		t.Errorf("Rebuild = %v, want %v", cs.Rebuild, want)
	}
	if len(cs.Remove) != 0 {
		t.Errorf("Remove = %v, want []", cs.Remove)
	}
}
