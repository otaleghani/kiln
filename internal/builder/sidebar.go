package builder

import (
	"net/url"
	"os"
	"path"          // Used for URL path construction (forward slashes)
	"path/filepath" // Used for OS file system operations
	"sort"
	"strings"

	"github.com/otaleghani/kiln/internal/log"
)

// getSidebarRootNode initializes the navigation tree starting from the input directory.
// It parses the BaseURL to ensure all node paths are prefixed correctly (e.g., "/kiln/foo").
func getSidebarRootNode(dir string, baseURL string) *SidebarNode {
	// Parse the Base URL to extract the path component.
	// Input: "https://example.com/kiln" -> basePath: "/kiln"
	// Input: "http://localhost:8080"    -> basePath: "/"
	u, err := url.Parse(baseURL)
	basePath := "/"
	if err == nil {
		basePath = u.Path
	} else {
		log.Warn("Failed to parse baseURL", log.FieldError, err)
	}

	// Normalize basePath to ensure it is clean for joining.
	if !strings.HasPrefix(basePath, "/") {
		basePath = "/" + basePath
	}
	basePath = strings.TrimSuffix(basePath, "/")
	if basePath == "" {
		basePath = "/"
	}

	rootNode := &SidebarNode{
		Name:     "Home",
		IsFolder: true,
		Path:     basePath,
	}

	// Construct the raw tree from the file system
	buildSidebarTree(dir, rootNode)

	// Remove empty folders
	rootNode.Children = pruneSidebarTree(rootNode.Children)

	// Sort folders first, then alphabetical
	sortSidebarTree(rootNode.Children)

	log.Info("File tree constructed, pruned, and sorted")

	return rootNode
}

// buildSidebarTree recursively walks the directory structure to populate the Node tree.
// It filters for .md and .canvas files and handles URL slug generation.
func buildSidebarTree(dir string, parent *SidebarNode) {
	entries, _ := os.ReadDir(dir)
	for _, entry := range entries {
		// Skip dotfiles (hidden files)
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		node := &SidebarNode{
			Name:     entry.Name(),
			IsFolder: entry.IsDir(),
		}

		nameForSlug := node.Name

		if !node.IsFolder {
			ext := filepath.Ext(node.Name)
			// Only include Markdown and Canvas files in the sidebar
			if ext != ".md" && ext != ".canvas" {
				continue
			}
			// Store display name without extension
			nameForSlug = strings.TrimSuffix(node.Name, ext)
			node.Name = nameForSlug
		}

		// URL Generation Logic
		// If the file is named "Home" or "index", it adopts the parent's path.
		if !node.IsFolder &&
			(strings.EqualFold(nameForSlug, "Home") || strings.EqualFold(nameForSlug, "index")) {
			node.Path = parent.Path
		} else {
			// Construct URL path: parent path + slugified name
			currentSlug := slugify(nameForSlug)
			node.Path = path.Join(parent.Path, currentSlug)
		}

		parent.Children = append(parent.Children, node)

		if entry.IsDir() {
			fullPath := filepath.Join(dir, entry.Name())
			buildSidebarTree(fullPath, node)
		}
	}
}

// pruneSidebarTree removes folders that end up empty (containing no valid .md or .canvas files).
// This is necessary because we might scan a folder that only contains images or unrelated files.
func pruneSidebarTree(nodes []*SidebarNode) []*SidebarNode {
	var kept []*SidebarNode
	for _, n := range nodes {
		if n.IsFolder {
			// Recursively prune children first
			n.Children = pruneSidebarTree(n.Children)
			// Only keep the folder if it still has children
			if len(n.Children) > 0 {
				kept = append(kept, n)
			}
		} else {
			// Since buildTree already filters non-md/canvas files, we keep all leaf nodes.
			kept = append(kept, n)
		}
	}
	return kept
}

// sortSidebarTree sorts nodes in place: Folders top, then files, both alphabetically.
func sortSidebarTree(nodes []*SidebarNode) {
	sort.Slice(nodes, func(i, j int) bool {
		// Prioritize Folders over Files
		if nodes[i].IsFolder && !nodes[j].IsFolder {
			return true
		}
		if !nodes[i].IsFolder && nodes[j].IsFolder {
			return false
		}
		// Alphabetical sort for same types
		return strings.ToLower(nodes[i].Name) < strings.ToLower(nodes[j].Name)
	})

	// Recursively sort children
	for _, n := range nodes {
		if n.IsFolder && len(n.Children) > 0 {
			sortSidebarTree(n.Children)
		}
	}
}

// setSidebarNodeActive traverses the tree and marks the node matching currentPath as Active.
// This is used by the template to highlight the current page in the sidebar.
func setSidebarNodeActive(nodes []*SidebarNode, currentPath string) {
	for _, n := range nodes {
		n.Active = (n.Path == currentPath)
		if n.IsFolder {
			setSidebarNodeActive(n.Children, currentPath)
		}
	}
}

// SidebarNode represents a single item (file or folder) in the sidebar navigation tree.
type SidebarNode struct {
	Name     string
	Path     string
	IsFolder bool
	Active   bool
	Children []*SidebarNode
}
