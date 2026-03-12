// @feature:watch Changeset computation combining mtime changes with dependency graph invalidation.
package watch

import (
	"path/filepath"
	"sort"
	"strings"
)

type ChangeSet struct {
	Rebuild []string // RelPaths to re-render
	Remove  []string // RelPaths to delete from output
}

// ComputeChangeSet expands a set of changed and removed file paths into a full
// ChangeSet by walking the dependency graph. Each changed or removed file's
// dependents are added to the rebuild set. Removed paths are excluded from
// rebuild and placed in Remove instead.
func ComputeChangeSet(changed, removed []string, graph *DepGraph) *ChangeSet {
	rebuildSet := make(map[string]struct{})
	removedSet := make(map[string]struct{})

	for _, relPath := range changed {
		rebuildSet[relPath] = struct{}{}
		for _, dep := range graph.Dependents(nameFromRelPath(relPath)) {
			rebuildSet[dep] = struct{}{}
		}
	}

	for _, relPath := range removed {
		removedSet[relPath] = struct{}{}
		for _, dep := range graph.Dependents(nameFromRelPath(relPath)) {
			rebuildSet[dep] = struct{}{}
		}
		graph.RemoveSource(relPath)
	}

	rebuild := make([]string, 0, len(rebuildSet))
	for p := range rebuildSet {
		if _, isRemoved := removedSet[p]; !isRemoved {
			rebuild = append(rebuild, p)
		}
	}
	sort.Strings(rebuild)

	remove := make([]string, 0, len(removedSet))
	for p := range removedSet {
		remove = append(remove, p)
	}
	sort.Strings(remove)

	return &ChangeSet{
		Rebuild: rebuild,
		Remove:  remove,
	}
}

// nameFromRelPath extracts the normalised name from a relative path:
// base filename, lowercased, with .md extension stripped.
func nameFromRelPath(relPath string) string {
	base := filepath.Base(relPath)
	base = strings.TrimSuffix(base, ".md")
	return strings.ToLower(base)
}
