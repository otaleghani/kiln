// @feature:watch Dependency graph for incremental rebuild invalidation.
package watch

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/otaleghani/kiln/internal/obsidian"
)

// DepGraph tracks dependencies between vault files. Forward maps each source
// RelPath to the set of target names it links to. Reverse maps each target
// name to the set of source RelPaths that reference it.
type DepGraph struct {
	Forward map[string]map[string]struct{}
	Reverse map[string]map[string]struct{}
}

func NewDepGraph() *DepGraph {
	return &DepGraph{
		Forward: make(map[string]map[string]struct{}),
		Reverse: make(map[string]map[string]struct{}),
	}
}

func (g *DepGraph) AddEdge(sourceRelPath, targetName string) {
	if _, ok := g.Forward[sourceRelPath]; !ok {
		g.Forward[sourceRelPath] = make(map[string]struct{})
	}
	g.Forward[sourceRelPath][targetName] = struct{}{}

	if _, ok := g.Reverse[targetName]; !ok {
		g.Reverse[targetName] = make(map[string]struct{})
	}
	g.Reverse[targetName][sourceRelPath] = struct{}{}
}

// RemoveSource removes all edges originating from sourceRelPath.
func (g *DepGraph) RemoveSource(sourceRelPath string) {
	targets, ok := g.Forward[sourceRelPath]
	if !ok {
		return
	}
	for target := range targets {
		if srcs, exists := g.Reverse[target]; exists {
			delete(srcs, sourceRelPath)
			if len(srcs) == 0 {
				delete(g.Reverse, target)
			}
		}
	}
	delete(g.Forward, sourceRelPath)
}

// Dependents returns the list of source RelPaths that link to the given name.
func (g *DepGraph) Dependents(name string) []string {
	srcs, ok := g.Reverse[name]
	if !ok {
		return nil
	}
	result := make([]string, 0, len(srcs))
	for s := range srcs {
		result = append(result, s)
	}
	return result
}

// UpdateFiles refreshes the graph for the given files, removing stale edges
// and adding new ones based on current links.
func (g *DepGraph) UpdateFiles(files []*obsidian.File) {
	for _, f := range files {
		if f.Ext != ".md" {
			continue
		}
		g.RemoveSource(f.RelPath)
		for _, link := range f.Links {
			target := parseTarget(link)
			if target == "" {
				continue
			}
			g.AddEdge(f.RelPath, target)
		}
	}
}

var mdLinkRe = regexp.MustCompile(`\[([^\]]*)\]\(([^)]+)\)`)

// BuildFromFiles populates the graph from a slice of obsidian files.
func (g *DepGraph) BuildFromFiles(files []*obsidian.File) {
	for _, f := range files {
		if f.Ext != ".md" {
			continue
		}
		for _, link := range f.Links {
			target := parseTarget(link)
			if target == "" {
				continue
			}
			g.AddEdge(f.RelPath, target)
		}
	}
}

// parseTarget extracts a normalised target name from a raw link string.
func parseTarget(link string) string {
	if strings.HasPrefix(link, "[[") {
		return parseWikilink(link)
	}
	if m := mdLinkRe.FindStringSubmatch(link); m != nil {
		return parseMdLink(m[2])
	}
	return ""
}

func parseWikilink(link string) string {
	inner := strings.TrimSuffix(strings.TrimPrefix(link, "[["), "]]")
	if i := strings.Index(inner, "|"); i >= 0 {
		inner = inner[:i]
	}
	if i := strings.Index(inner, "#"); i >= 0 {
		inner = inner[:i]
	}
	inner = strings.TrimSpace(inner)
	if strings.Contains(inner, "/") {
		inner = filepath.Base(inner)
	}
	return strings.ToLower(inner)
}

func parseMdLink(path string) string {
	for _, prefix := range []string{"https://", "http://", "mailto:"} {
		if strings.HasPrefix(path, prefix) {
			return ""
		}
	}
	if i := strings.Index(path, "#"); i >= 0 {
		path = path[:i]
	}
	path = filepath.Base(path)
	path = strings.TrimSuffix(path, ".md")
	return strings.ToLower(path)
}
