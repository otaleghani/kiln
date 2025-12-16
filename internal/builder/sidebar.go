package builder

import (
	"log"
	"os"
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

func getRootNode(dir string) *Node {
	rootNode := &Node{
		Name:     "Home",
		IsFolder: true,
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
		// Skip hidden files/folders
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		node := &Node{
			Name:     entry.Name(),
			IsFolder: entry.IsDir(),
		}

		// 1. Calculate the SLUG for this node (Folder OR File)
		// Clean the name (remove extension if it's a file)
		nameForSlug := node.Name
		if !node.IsFolder {
			ext := filepath.Ext(node.Name)
			// ALLOW both .md and .canvas
			if ext != ".md" && ext != ".canvas" {
				continue
			}
			nameForSlug = strings.TrimSuffix(node.Name, ext)
			node.Name = nameForSlug // Update display name
		}

		// 2. Build the URL Path relative to the PARENT
		// This ensures folders get paths, and files preserve their hierarchy.
		currentSlug := slugify(nameForSlug)

		// Special handling to avoid double slashes if parent is root "/"
		parentPath := parent.Path
		if parentPath == "/" {
			parentPath = ""
		}

		// Assign path to BOTH folders and files
		node.Path = parentPath + "/" + currentSlug

		// Special case: "Home" or "index" overrides the path to root
		if !node.IsFolder &&
			(strings.EqualFold(nameForSlug, "Home") || strings.EqualFold(nameForSlug, "index")) {
			node.Path = "/"
		}

		parent.Children = append(parent.Children, node)

		if entry.IsDir() {
			// Calculate full physical path for the next recursion step
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
			ext := filepath.Ext(n.Path)
			lowerExt := strings.ToLower(ext)
			if ext == "" || lowerExt == ".md" || lowerExt == ".canvas" {
				kept = append(kept, n)
			}
		}
	}
	return kept
}

// sortTree recursively sorts nodes: Folders first, then alphabetically (A-Z)
func sortTree(nodes []*Node) {
	sort.Slice(nodes, func(i, j int) bool {
		// Folders always come before files
		if nodes[i].IsFolder && !nodes[j].IsFolder {
			return true
		}
		if !nodes[i].IsFolder && nodes[j].IsFolder {
			return false
		}

		// Alphabetical sort (case-insensitive)
		return strings.ToLower(nodes[i].Name) < strings.ToLower(nodes[j].Name)
	})

	// Sort children of folder nodes
	for _, n := range nodes {
		if n.IsFolder && len(n.Children) > 0 {
			sortTree(n.Children)
		}
	}
}

// setTreeActive sets the node as active if it has the currentPath
func setTreeActive(nodes []*Node, currentPath string) {
	for _, n := range nodes {
		n.Active = (n.Path == currentPath)
		if n.IsFolder {
			setTreeActive(n.Children, currentPath)
		}
	}
}
