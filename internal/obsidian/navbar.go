package obsidian

import (
	"net/url"
	"path/filepath"
	"sort"
	"strings"
)

// GenerateNavbar constructs the navigation tree using the existing Vault data.
// It replaces getSidebarRootNode and buildSidebarTree.
func (o *Obsidian) GenerateNavbar() *NavbarNode {
	// Calculate base path from BaseURL
	basePath := "/"
	if u, err := url.Parse(o.BaseURL); err == nil {
		basePath = u.Path
	} else {
		o.log.Warn("Failed to parse baseURL", "error", err)
	}

	// Normalize basePath
	if !strings.HasPrefix(basePath, "/") {
		basePath = "/" + basePath
	}
	basePath = strings.TrimSuffix(basePath, "/")
	if basePath == "" {
		basePath = "/"
	}

	// Initialize Root Node
	rootNode := &NavbarNode{
		Name:     "Home",
		IsFolder: true,
		Path:     basePath,
		Children: []*NavbarNode{},
	}

	// Create a registry of Folder Nodes
	// Key: RelPath (e.g., "folder/subfolder"), Value: *SidebarNode
	folderMap := make(map[string]*NavbarNode)

	// Create nodes for all Folders in the Vault
	for _, folder := range o.Vault.Folders {
		// Clean the path to handle potential "./" or OS-specific oddities
		rel := filepath.Clean(folder.RelPath)

		// Skip if this represents the root dir itself (we use rootNode for that)
		if rel == "." || rel == "/" || folder.RelPath == "" {
			continue
		}

		node := &NavbarNode{
			Name:     filepath.Base(folder.RelPath),
			Path:     folder.WebPath, // Use pre-calculated WebPath
			IsFolder: true,
			Children: []*NavbarNode{},
		}
		folderMap[folder.RelPath] = node
	}

	// Populate Files into their respective parent Nodes
	// We iterate over All Files to ensure we catch everything, including root files.
	for _, file := range o.Vault.Files {
		// Filter allowed extensions
		var isNote, isCanvas, isBase bool
		switch file.Ext {
		case ".md":
			isNote = true
		case ".canvas":
			isCanvas = true
		case ".base":
			isBase = true
		default:
			continue
		}

		if node, exists := folderMap[strings.TrimSuffix(file.RelPath, file.Ext)]; exists {
			switch file.Ext {
			case ".md":
				node.IsNote = true
			case ".canvas":
				node.IsCanvas = true
			case ".base":
				node.IsBase = true
			default:
				continue
			}
			// Do not add again the same if a folder exists
			continue
		}

		node := &NavbarNode{
			Name: strings.TrimSuffix(
				file.Name,
				file.Ext,
			), // TODO: file.Name should be without ext, no?
			Path:     file.WebPath,
			IsNote:   isNote,
			IsCanvas: isCanvas,
			IsBase:   isBase,
			IsFolder: false,
		}

		// Determine the parent folder based on Relative Path
		parentDir := filepath.Dir(file.RelPath)

		// If parent is root ("." or "/"), add to rootNode
		if parentDir == "." || parentDir == "/" || parentDir == "" {
			rootNode.Children = append(rootNode.Children, node)
		} else {
			// Otherwise, find the parent folder node
			if parentNode, ok := folderMap[parentDir]; ok {
				parentNode.Children = append(parentNode.Children, node)
			} else {
				// Edge case: File exists in a folder not mapped in Vault.Folders.
				// Should not happen if Vault is consistent.
				o.log.Debug("Orphaned file found (parent folder missing in map)", "file", file.RelPath)
			}
		}
	}

	// Link folder nodes to their parents
	for relPath, node := range folderMap {
		parentDir := filepath.Dir(relPath)

		// If parent is root, attach to rootNode
		if parentDir == "." || parentDir == "/" || parentDir == "" {
			rootNode.Children = append(rootNode.Children, node)
		} else {
			// Otherwise, attach to the parent folder node
			if parentNode, ok := folderMap[parentDir]; ok {
				parentNode.Children = append(parentNode.Children, node)
			}
		}
	}

	// Prune and Sort
	rootNode.Children = pruneNavbarTree(rootNode.Children)
	sortNavbarTree(rootNode.Children)

	o.log.Info("Sidebar tree constructed from Vault data")

	return rootNode
}

// --- Helper Functions (Preserved logic) ---

// pruneNavbarTree removes folders that end up empty.
func pruneNavbarTree(nodes []*NavbarNode) []*NavbarNode {
	var kept []*NavbarNode
	for _, n := range nodes {
		if n.IsFolder {
			// Recursively prune children first
			n.Children = pruneNavbarTree(n.Children)
			// Only keep the folder if it still has children
			if len(n.Children) > 0 {
				kept = append(kept, n)
			}
		} else {
			kept = append(kept, n)
		}
	}
	return kept
}

// sortNavbarTree sorts nodes in place: Folders top, then files, both alphabetically.
func sortNavbarTree(nodes []*NavbarNode) {
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
			sortNavbarTree(n.Children)
		}
	}
}

// TODO: Delete this, is deprecated in favor of a small client-side script
// SetNavbarNodeActive traverses the tree and marks the node matching currentPath as Active.
func SetNavbarNodeActive(nodes []*NavbarNode, currentPath string) {
	for _, n := range nodes {
		// You might need to normalize paths here depending on strictness
		n.Active = (n.Path == currentPath)

		if n.IsFolder {
			SetNavbarNodeActive(n.Children, currentPath)
		}
	}
}

// NavbarNode represents a single item (file or folder) in the sidebar navigation tree.
type NavbarNode struct {
	Name     string
	Path     string
	IsFolder bool
	IsCanvas bool
	IsBase   bool
	IsNote   bool
	Active   bool
	Children []*NavbarNode
}
