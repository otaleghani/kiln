package builder

import (
	"encoding/xml"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Init checks if the input directory (vault) exists.
// If not, it creates the directory and a default "Home.md" welcome note.
func Init() {
	if _, err := os.Stat(InputDir); os.IsNotExist(err) {
		os.Mkdir(InputDir, 0755)
		log.Println("Created vault directory.")

		// Create a welcome note to get the user started
		welcomeText := "# Welcome to Kiln\n\nThis is your new vault. Run `kiln generate` to build it!"
		os.WriteFile(filepath.Join(InputDir, "Home.md"), []byte(welcomeText), 0644)
	} else {
		log.Println("Vault directory already exists.")
	}
	log.Println("Initialization complete.")
}

// CleanOutDir removes the entire output directory to ensure a clean build.
// This prevents stale files from persisting in the generated site.
func CleanOutDir() {
	err := os.RemoveAll(OutputDir)
	if err != nil {
		log.Printf("Error cleaning output: %v", err)
	} else {
		log.Println("Cleaned ./public directory")
	}
}

// initBuild prepares the environment for a new build.
// It cleans the output directory, copies global assets (favicon, CNAME),
// and traverses the input directory to build a file index and a graph of nodes.
func initBuild() (map[string]string, []GraphNode) {
	// 1. Reset output state
	CleanOutDir()
	os.MkdirAll(OutputDir, 0755)

	// 2. Copy Global Assets
	// Copy over favicon.ico if it exists in the root
	faviconSrc := filepath.Join(InputDir, "favicon.ico")
	if _, err := os.Stat(faviconSrc); err == nil {
		if err := copyFile(faviconSrc, filepath.Join(OutputDir, "favicon.ico")); err != nil {
			log.Printf("Warning: Found favicon.ico but failed to copy it: %v", err)
		} else {
			log.Println("Found and copied favicon.ico")
		}
	}

	// Copy over CNAME file (for custom domains on GitHub Pages/Netlify) if it exists
	cnameSrc := filepath.Join(InputDir, "CNAME")
	if _, err := os.Stat(cnameSrc); err == nil {
		if err := copyFile(cnameSrc, filepath.Join(OutputDir, "CNAME")); err != nil {
			log.Printf("Warning: Found CNAME but failed to copy it: %v", err)
		} else {
			log.Println("Found and copied CNAME")
		}
	}

	// 3. Index Content
	// Creates file index (name -> url) and file nodes (for the graph view)
	fileIndex := make(map[string]string)
	graphNodes := []GraphNode{}

	// Map to track unique nodes and avoid duplicates in the graph
	nodeSet := make(map[string]bool)

	filepath.WalkDir(InputDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		// Skip hidden files and folders (e.g. .obsidian, .trash, .git)
		// We check if the name starts with "." but ensure we aren't skipping the root directory itself.
		if strings.HasPrefix(d.Name(), ".") && path != InputDir {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// We only process files in this pass, not directories
		if d.IsDir() {
			return nil
		}

		ext := filepath.Ext(path)
		relPath, _ := filepath.Rel(InputDir, path)

		// Normalize paths for URLs (slugify every component of the path)
		parts := strings.Split(relPath, string(os.PathSeparator))
		for i, p := range parts {
			parts[i] = slugify(p)
		}

		// Join parts back together. Note: slugPath currently includes the extension (e.g., "my-note-md")
		slugPath := filepath.Join(parts...)
		webPath := "/" + strings.ReplaceAll(slugPath, string(os.PathSeparator), "/")

		key := d.Name()

		// Handle Content Files (Markdown and Canvas)
		if ext == ".md" || ext == ".canvas" {
			key = strings.TrimSuffix(key, ext)

			// Remove the extension from the web path to create clean URLs
			// e.g., "notes/my-note-md" -> "notes/my-note"
			slugFolder := strings.TrimSuffix(slugPath, slugify(ext))
			webPath = "/" + strings.ReplaceAll(slugFolder, string(os.PathSeparator), "/")

			// Special handling for index/Home notes to map them to the root URL "/"
			if strings.EqualFold(key, "Home") || strings.EqualFold(key, "index") {
				webPath = "/"
			}

			// Add to Graph Nodes list
			if !nodeSet[key] {
				graphNodes = append(graphNodes, GraphNode{
					ID:    key,
					Label: key,
					URL:   webPath,
					Val:   1, // Default weight for visualization
				})
				nodeSet[key] = true
			}
		}

		// Register the file in the global index (filename -> public URL)
		// This is used later for resolving [[WikiLinks]].
		if _, exists := fileIndex[key]; !exists {
			fileIndex[key] = webPath
		}
		return nil
	})

	return fileIndex, graphNodes
}

// copyFile is a simple wrapper to copy files from src to dst.
func copyFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}

// slugify converts a string to a URL-friendly format.
// It lowercases the string and replaces spaces with dashes.
// Example: "My Note" -> "my-note"
func slugify(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	return s
}

// isImageExt checks if the given file extension corresponds to a supported image format.
func isImageExt(ext string) bool {
	switch strings.ToLower(ext) {
	case ".png", ".jpg", ".jpeg", ".gif", ".webp", ".svg":
		return true
	default:
		return false
	}
}

// getSlugPath converts a relative file path into a slugified path string.
// It iterates through every directory in the path and slugifies it.
func getSlugPath(relPath string) string {
	parts := strings.Split(relPath, string(os.PathSeparator))
	for i, p := range parts {
		parts[i] = slugify(p)
	}
	slugPath := filepath.Join(parts...)
	return slugPath
}

// getOutputPaths determines the filesystem output location and the public web URL for a file.
// It handles "Pretty URLs" by creating index.html files inside named directories.
// Returns:
// - outPath: where to write the file on disk (e.g., ./public/notes/my-note/index.html)
// - webPath: the URL path (e.g., /notes/my-note/)
func getOutputPaths(relPath, nameWithoutExt, ext string) (outPath string, webPath string) {
	slugPath := getSlugPath(relPath)

	// Case 1: Root files (Home or index) -> ./public/index.html
	if strings.EqualFold(nameWithoutExt, "Home") ||
		strings.EqualFold(nameWithoutExt, "index") {
		outPath = filepath.Join(OutputDir, "index.html")
		webPath = "/"
	} else {
		// Regular files -> ./public/folder/note-name/index.html
		slugFolder := strings.TrimSuffix(slugPath, slugify(ext))
		outPath = filepath.Join(OutputDir, slugFolder, "index.html")
		webPath = "/" + strings.ReplaceAll(slugFolder, string(os.PathSeparator), "/")

		// Regular files -> ./public/folder/note-name.html
		// slugName := strings.TrimSuffix(slugPath, slugify(ext))
		// outPath = filepath.Join(OutputDir, slugName+".html")
		// webPath = "/" + strings.ReplaceAll(slugName, string(os.PathSeparator), "/")
	}

	// Ensure the parent directory exists before returning
	os.MkdirAll(filepath.Dir(outPath), 0755)
	return
}

// getBreadcrumbs constructs a list of navigation steps for the current page.
// It walks up the directory tree from the current file location.
func getBreadcrumbs(relPath, nameWithoutExt string) []string {
	var breadcrumbs []string
	breadcrumbs = append(breadcrumbs, "Home")

	// Add intermediate directories
	dir := filepath.Dir(relPath)
	if dir != "." && dir != "" {
		breadcrumbs = append(breadcrumbs, strings.Split(dir, string(os.PathSeparator))...)
	}

	// Add current page (unless it is the home page itself)
	if !strings.EqualFold(nameWithoutExt, "index") &&
		!strings.EqualFold(nameWithoutExt, "Home") {
		breadcrumbs = append(breadcrumbs, nameWithoutExt)
	}
	return breadcrumbs
}

// SitemapEntry represents a single URL entry in the sitemap.xml.
type SitemapEntry struct {
	XMLName xml.Name `xml:"url"`
	Loc     string   `xml:"loc"`     // The absolute URL
	LastMod string   `xml:"lastmod"` // The last modification date
}

// addToSitemap appends a new entry to the sitemap slice.
// It retrieves the file's modification time to populate 'lastmod'.
func addToSitemap(d fs.DirEntry, baseURL, webPath string, sitemapEntries *[]SitemapEntry) {
	info, err := d.Info()
	modTime := time.Now()
	if err == nil {
		modTime = info.ModTime()
	}

	fullURL := strings.TrimRight(baseURL, "/") + webPath
	*sitemapEntries = append(*sitemapEntries, SitemapEntry{
		Loc:     fullURL,
		LastMod: modTime.Format("2006-01-02"),
	})
}

// generateRobots creates a robots.txt file in the output directory.
// It points crawlers to the Sitemap location.
func generateRobots(baseURL string) {
	robotsFile, _ := os.Create(filepath.Join(OutputDir, "robots.txt"))
	defer robotsFile.Close()
	robotsFile.WriteString("User-agent: *\n")
	robotsFile.WriteString("Allow: /\n")
	robotsFile.WriteString("Sitemap: " + strings.TrimRight(baseURL, "/") + "/sitemap.xml\n")
}

// generateSitemap marshals the list of sitemap entries into XML format
// and writes it to sitemap.xml in the output directory.
func generateSitemap(sitemapEntries []SitemapEntry) {
	sitemapFile, _ := os.Create(filepath.Join(OutputDir, "sitemap.xml"))
	defer sitemapFile.Close()
	sitemapFile.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	sitemapFile.WriteString(
		`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">` + "\n",
	)
	for _, entry := range sitemapEntries {
		output, _ := xml.MarshalIndent(entry, "  ", "  ")
		sitemapFile.Write(output)
		sitemapFile.WriteString("\n")
	}
	sitemapFile.WriteString(`</urlset>`)

	log.Println("Generated sitemap.xml and robots.txt")
}
