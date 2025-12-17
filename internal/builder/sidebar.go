package builder

import (
	"log"
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

// Added baseURL parameter
func getRootNode(dir string, baseURL string) *Node {
	// 1. Sanitize BaseURL
	// Ensure it starts with / and clean trailing slashes
	cleanBase := baseURL
	if cleanBase == "" {
		cleanBase = "/"
	} else {
		// Ensure it starts with / if not empty
		if !strings.HasPrefix(cleanBase, "/") {
			cleanBase = "/" + cleanBase
		}
		// Remove trailing slash for consistency (we add it back when joining)
		cleanBase = strings.TrimSuffix(cleanBase, "/")
		if cleanBase == "" {
			cleanBase = "/"
		}
	}

	rootNode := &Node{
		Name:     "Home",
		IsFolder: true,
		Path:     cleanBase, // 2. Set the Root Node's path explicitly
	}

	buildTree(dir, rootNode, cleanBase) // 3. Pass baseURL down
	rootNode.Children = pruneTree(rootNode.Children)
	sortTree(rootNode.Children)
	log.Println("File tree constructed, pruned, and sorted")

	return rootNode
}

// buildTree recursively walks the directory to create the sidebar structure
func buildTree(dir string, parent *Node, baseURL string) {
	entries, _ := os.ReadDir(dir)
	for _, entry := range entries {
		// Skip hidden files/folders
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		node := &Node{
			Name:     entry.Name(),
			IsFolder: entry.IsDir(),
		}

		// Calculate the SLUG for this node
		nameForSlug := node.Name
		if !node.IsFolder {
			ext := filepath.Ext(node.Name)
			if ext != ".md" && ext != ".canvas" {
				continue
			}
			nameForSlug = strings.TrimSuffix(node.Name, ext)
			node.Name = nameForSlug
		}

		// --- URL GENERATION LOGIC ---

		// 1. Check for "Home" or "index" overrides FIRST
		// Instead of hardcoding "/", we use the parent's path.
		// This respects the BaseURL and correctly handles nested index files (e.g. /docs/index)
		if !node.IsFolder &&
			(strings.EqualFold(nameForSlug, "Home") || strings.EqualFold(nameForSlug, "index")) {
			node.Path = parent.Path
		} else {
			// 2. Standard Slug Generation
			currentSlug := slugify(nameForSlug)

			// We use path.Join for URLs (it handles forward slashes correctly on all OSs)
			// It also handles the "/" + "/" case automatically.
			node.Path = path.Join(parent.Path, currentSlug)
		}

		parent.Children = append(parent.Children, node)

		if entry.IsDir() {
			fullPath := filepath.Join(dir, entry.Name())
			buildTree(fullPath, node, baseURL)
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
