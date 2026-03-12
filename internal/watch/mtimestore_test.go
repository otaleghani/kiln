// @feature:watch Tests for mtime store.
package watch

import (
	"os"
	"path/filepath"
	"slices"
	"sort"
	"testing"
	"time"
)

func TestUpdateNewFiles(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "a.md"), "alpha")
	writeFile(t, filepath.Join(dir, "b.md"), "bravo")

	s := NewMtimeStore()
	changed, removed, err := s.Update(dir)
	if err != nil {
		t.Fatal(err)
	}

	sort.Strings(changed)
	want := []string{"a.md", "b.md"}
	if !slices.Equal(changed, want) {
		t.Errorf("changed = %v, want %v", changed, want)
	}
	if len(removed) != 0 {
		t.Errorf("removed = %v, want []", removed)
	}
}

func TestUpdateUnchanged(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "a.md"), "alpha")

	s := NewMtimeStore()
	if _, _, err := s.Update(dir); err != nil {
		t.Fatal(err)
	}

	changed, removed, err := s.Update(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(changed) != 0 {
		t.Errorf("changed = %v, want []", changed)
	}
	if len(removed) != 0 {
		t.Errorf("removed = %v, want []", removed)
	}
}

func TestUpdateModifiedFile(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "a.md"), "alpha")
	writeFile(t, filepath.Join(dir, "b.md"), "bravo")

	s := NewMtimeStore()
	if _, _, err := s.Update(dir); err != nil {
		t.Fatal(err)
	}

	// Bump mtime on a.md by rewriting with a future timestamp.
	path := filepath.Join(dir, "a.md")
	writeFile(t, path, "alpha-modified")
	future := time.Now().Add(2 * time.Second)
	if err := os.Chtimes(path, future, future); err != nil {
		t.Fatal(err)
	}

	changed, removed, err := s.Update(dir)
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(changed, []string{"a.md"}) {
		t.Errorf("changed = %v, want [a.md]", changed)
	}
	if len(removed) != 0 {
		t.Errorf("removed = %v, want []", removed)
	}
}

func TestUpdateRemovedFile(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "a.md"), "alpha")
	writeFile(t, filepath.Join(dir, "b.md"), "bravo")

	s := NewMtimeStore()
	if _, _, err := s.Update(dir); err != nil {
		t.Fatal(err)
	}

	if err := os.Remove(filepath.Join(dir, "b.md")); err != nil {
		t.Fatal(err)
	}

	changed, removed, err := s.Update(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(changed) != 0 {
		t.Errorf("changed = %v, want []", changed)
	}
	if !slices.Equal(removed, []string{"b.md"}) {
		t.Errorf("removed = %v, want [b.md]", removed)
	}
}

func TestSkipsDotfiles(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "visible.md"), "yes")

	hiddenDir := filepath.Join(dir, ".hidden")
	if err := os.Mkdir(hiddenDir, 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(hiddenDir, "secret.md"), "no")
	writeFile(t, filepath.Join(dir, ".dotfile"), "no")

	s := NewMtimeStore()
	changed, _, err := s.Update(dir)
	if err != nil {
		t.Fatal(err)
	}

	if !slices.Equal(changed, []string{"visible.md"}) {
		t.Errorf("changed = %v, want [visible.md]", changed)
	}
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
