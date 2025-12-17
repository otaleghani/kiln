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

func Init() {
	if _, err := os.Stat(InputDir); os.IsNotExist(err) {
		os.Mkdir(InputDir, 0755)
		log.Println("Created vault directory.")

		// Create a welcome note
		welcomeText := "# Welcome to Kiln\n\nThis is your new vault. Run `kiln generate` to build it!"
		os.WriteFile(filepath.Join(InputDir, "Home.md"), []byte(welcomeText), 0644)
	} else {
		log.Println("Vault directory already exists.")
	}
	log.Println("Initialization complete.")

}

func CleanOutDir() {
	err := os.RemoveAll(OutputDir)
	if err != nil {
		log.Printf("Error cleaning output: %v", err)
	} else {
		log.Println("Cleaned ./public directory")
	}
}

// initBuild clears the output directory, copies over specific files and returns the file index and graph nodes
func initBuild() (map[string]string, []GraphNode) {
	// Remove output directory
	CleanOutDir()
	os.MkdirAll(OutputDir, 0755)

	// Copy over favicon.ico if it exists
	faviconSrc := filepath.Join(InputDir, "favicon.ico")
	if _, err := os.Stat(faviconSrc); err == nil {
		if err := copyFile(faviconSrc, filepath.Join(OutputDir, "favicon.ico")); err != nil {
			log.Printf("Warning: Found favicon.ico but failed to copy it: %v", err)
		} else {
			log.Println("Found and copied favicon.ico")
		}
	}

	// Creates file index and file nodes by traversing the input directory
	fileIndex := make(map[string]string)
	graphNodes := []GraphNode{}

	// Map to ensure unique nodes
	nodeSet := make(map[string]bool)

	filepath.WalkDir(InputDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		// Skip hidden files and folders (e.g. .obsidian, .trash, .git)
		if strings.HasPrefix(d.Name(), ".") && path != InputDir {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if d.IsDir() {
			return nil
		}

		ext := filepath.Ext(path)
		relPath, _ := filepath.Rel(InputDir, path)

		// Replicate Step 8 logic for consistent URL generation
		parts := strings.Split(relPath, string(os.PathSeparator))
		for i, p := range parts {
			parts[i] = slugify(p)
		}
		// slugPath here might contain the extension in the last part (e.g. "note-md")
		slugPath := filepath.Join(parts...)
		webPath := "/" + strings.ReplaceAll(slugPath, string(os.PathSeparator), "/")

		key := d.Name()
		if ext == ".md" || ext == ".canvas" {
			key = strings.TrimSuffix(key, ext)

			// Fix: Strip extension from URL for content files
			// We remove the slugified extension from the end of the slugPath
			slugFolder := strings.TrimSuffix(slugPath, slugify(ext))
			webPath = "/" + strings.ReplaceAll(slugFolder, string(os.PathSeparator), "/")

			if strings.EqualFold(key, "Home") || strings.EqualFold(key, "index") {
				webPath = "/"
			}

			// Add to Graph Nodes
			if !nodeSet[key] {
				graphNodes = append(graphNodes, GraphNode{
					ID:    key,
					Label: key,
					URL:   webPath,
					Val:   1, // Default weight
				})
				nodeSet[key] = true
			}
		}

		if _, exists := fileIndex[key]; !exists {
			fileIndex[key] = webPath
		}
		return nil
	})

	return fileIndex, graphNodes
}

// copyFile is a simple wrapper to copy files
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

// slugify converts a string to a URL-friendly format:
// "My Note" -> "my-note"
func slugify(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	return s
}

// isImageExt tells you if the given extension is an image or not
func isImageExt(ext string) bool {
	switch strings.ToLower(ext) {
	case ".png", ".jpg", ".jpeg", ".gif", ".webp", ".svg":
		return true
	default:
		return false
	}
}

// Based on the relative path returns the slugify version of it
func getSlugPath(relPath string) string {
	// Determine output path
	parts := strings.Split(relPath, string(os.PathSeparator))
	for i, p := range parts {
		parts[i] = slugify(p)
	}
	slugPath := filepath.Join(parts...)

	return slugPath
}

// Given the name of the file, the extension and the relative path, returns the output path and the web path
func getOutputPaths(relPath, nameWithoutExt, ext string) (outPath string, webPath string) {
	slugPath := getSlugPath(relPath)

	if strings.EqualFold(nameWithoutExt, "Home") ||
		strings.EqualFold(nameWithoutExt, "index") {
		outPath = filepath.Join(OutputDir, "index.html")
		webPath = "/"
	} else {
		slugFolder := strings.TrimSuffix(slugPath, slugify(ext))
		outPath = filepath.Join(OutputDir, slugFolder, "index.html")
		webPath = "/" + strings.ReplaceAll(slugFolder, string(os.PathSeparator), "/")
	}

	os.MkdirAll(filepath.Dir(outPath), 0755)
	return
}

// Returns the breadcrumbs for the given note
func getBreadcrumbs(relPath, nameWithoutExt string) []string {
	var breadcrumbs []string
	breadcrumbs = append(breadcrumbs, "Home")
	dir := filepath.Dir(relPath)
	if dir != "." && dir != "" {
		breadcrumbs = append(breadcrumbs, strings.Split(dir, string(os.PathSeparator))...)
	}
	if !strings.EqualFold(nameWithoutExt, "index") &&
		!strings.EqualFold(nameWithoutExt, "Home") {
		breadcrumbs = append(breadcrumbs, nameWithoutExt)
	}
	return breadcrumbs
}

type SitemapEntry struct {
	XMLName xml.Name `xml:"url"`
	Loc     string   `xml:"loc"`
	LastMod string   `xml:"lastmod"`
}

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

func generateRobots(baseURL string) {
	robotsFile, _ := os.Create(filepath.Join(OutputDir, "robots.txt"))
	defer robotsFile.Close()
	robotsFile.WriteString("User-agent: *\n")
	robotsFile.WriteString("Allow: /\n")
	robotsFile.WriteString("Sitemap: " + strings.TrimRight(baseURL, "/") + "/sitemap.xml\n")
}

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
