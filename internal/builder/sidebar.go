package builder

import (
	"log"
	"net/url"
	"os"
	"path" // Use 'path' for URLs, 'filepath' for OS files
	"path/filepath"
	"sort"
	"strings"
)

// Node represents a file or folder in the sidebar tree
type Node struct {
	Name     string
	Path     string
	IsFolder bool
	Active   bool
	Children []*Node
}

func getRootNode(dir string, fullURL string) *Node {
	// Parse the Full URL to extract ONLY the path
	// Input: "https://otaleghani.github.io/kiln" -> basePaths: "/kiln"
	// Input: "http://localhost:8080"             -> basePaths: "/"
	u, err := url.Parse(fullURL)
	basePath := "/"
	if err == nil {
		basePath = u.Path
	}

	// Sanitize the Path (ensure it starts/ends correctly for joining)
	// e.g. "/kiln" or "/"
	if !strings.HasPrefix(basePath, "/") {
		basePath = "/" + basePath
	}
	// We trim the suffix here so path.Join doesn't get confused later,
	// unless it's just root "/"
	basePath = strings.TrimSuffix(basePath, "/")
	if basePath == "" {
		basePath = "/"
	}

	rootNode := &Node{
		Name:     "Home",
		IsFolder: true,
		Path:     basePath, // Set the root path to "/kiln", not the full URL
	}

	buildTree(dir, rootNode)
	rootNode.Children = pruneTree(rootNode.Children)
	sortTree(rootNode.Children)
	log.Println("File tree constructed, pruned, and sorted")

	return rootNode
}

// buildTree recursively walks the directory to create the sidebar structure
func buildTree(dir string, parent *Node) {
	entries, _ := os.ReadDir(dir)
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		node := &Node{
			Name:     entry.Name(),
			IsFolder: entry.IsDir(),
		}

		nameForSlug := node.Name
		if !node.IsFolder {
			ext := filepath.Ext(node.Name)
			if ext != ".md" && ext != ".canvas" {
				continue
			}
			nameForSlug = strings.TrimSuffix(node.Name, ext)
			node.Name = nameForSlug
		}

		// URL GENERATION
		// Handle Home/Index override
		if !node.IsFolder &&
			(strings.EqualFold(nameForSlug, "Home") || strings.EqualFold(nameForSlug, "index")) {
			// Inherit parent path (e.g. "/kiln/" or "/kiln/docs")
			node.Path = parent.Path
		} else {
			// Standard Join
			currentSlug := slugify(nameForSlug)
			// path.Join handles the slashes intelligently
			node.Path = path.Join(parent.Path, currentSlug)
		}

		parent.Children = append(parent.Children, node)

		if entry.IsDir() {
			fullPath := filepath.Join(dir, entry.Name())
			buildTree(fullPath, node)
		}
	}
}

// pruneTree removes empty folders or folders that do not contain .md files or .canvas files
func pruneTree(nodes []*Node) []*Node {
	var kept []*Node
	for _, n := range nodes {
		if n.IsFolder {
			n.Children = pruneTree(n.Children)
			if len(n.Children) > 0 {
				kept = append(kept, n)
			}
		} else {
			// Use path.Ext for URL paths
			ext := path.Ext(n.Path)
			// If path is root or just base url, ext might be empty, which is fine for folders/root
			if !n.IsFolder && ext == "" {
				// If it's a file but has no extension in the Path (slugified), keep it
				kept = append(kept, n)
			} else {
				// Original logic
				// lowerExt := strings.ToLower(filepath.Ext(n.Name)) // Check extension of original Name, not Path
				// Or if you only store extension-less names now, you might need to check IsFolder
				// The original logic checked the Path extension, but our URLs don't have .md anymore!
				// Correct Logic: Just keep all files that made it this far (since buildTree filters .md/.canvas)
				kept = append(kept, n)
			}
		}
	}
	return kept
}

// sortTree recursively sorts nodes
func sortTree(nodes []*Node) {
	sort.Slice(nodes, func(i, j int) bool {
		if nodes[i].IsFolder && !nodes[j].IsFolder {
			return true
		}
		if !nodes[i].IsFolder && nodes[j].IsFolder {
			return false
		}
		return strings.ToLower(nodes[i].Name) < strings.ToLower(nodes[j].Name)
	})

	for _, n := range nodes {
		if n.IsFolder && len(n.Children) > 0 {
			sortTree(n.Children)
		}
	}
}

func setTreeActive(nodes []*Node, currentPath string) {
	for _, n := range nodes {
		// Ensure strict slash matching or clean paths
		n.Active = (n.Path == currentPath)
		if n.IsFolder {
			setTreeActive(n.Children, currentPath)
		}
	}
}
