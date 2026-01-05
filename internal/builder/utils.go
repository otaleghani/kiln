package builder

import (
	"io"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/otaleghani/kiln/internal/log"
	obsidianmarkdown "github.com/otaleghani/kiln/pkg/obsidian-markdown"
)

// Init checks if the input directory (vault) exists.
// If not, it creates the directory and a default "Home.md" welcome note.
func Init() {
	_, err := os.Stat(InputDir)

	if err == nil {
		log.Error("Vault directory already exists")
		return
	}

	if !os.IsNotExist(err) {
		log.Error("Couldn't read information about directory", log.FieldError, err)
		return
	}

	err = os.Mkdir(InputDir, 0755)
	if err != nil {
		log.Error("Couldn't create folder", log.FieldError, err)
		return
	}

	log.Info("Created vault directory")

	// Create a welcome note to get the user started
	welcomeText := "# Welcome to Kiln\n\nThis is your new vault. Run `kiln generate` to build it!"
	err = os.WriteFile(filepath.Join(InputDir, "Home.md"), []byte(welcomeText), 0644)
	if err != nil {
		log.Error("Couldn't create welcome note", log.FieldError, err)
		return
	}

	log.Info("Initialization complete")
}

// CleanOutputDir removes the entire output directory to ensure a clean build.
// This prevents stale files from persisting in the generated site.
func CleanOutputDir() {
	err := os.RemoveAll(OutputDir)
	if err != nil {
		log.Error("Couldn't remove output directory", log.FieldError, err)
	} else {
		log.Info("Cleaned directory", log.FieldPath, OutputDir)
	}
}

func scanVault() (*VaultScan, error) {
	log.Info("Scanning vault...")
	scan := &VaultScan{
		FileIndex:  make(map[string][]*File),
		GraphNodes: []GraphNode{},
		SourceMap:  make(map[string]string),
		Sitemap:    &Sitemap{},
	}
	scan.Sitemap.Path = filepath.Join(BaseURL, "/sitemap.xml")

	filepath.WalkDir(InputDir, func(path string, info fs.DirEntry, err error) error {
		l := log.Default.WithFile(path)
		l.Debug("Processing file...")

		// Handle permission errors and other related problems
		if err != nil {
			return nil
		}

		// Create relative path
		relPath, err := filepath.Rel(InputDir, path)
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			l.Debug("Skipping file", log.FieldReason, "File is a directory")
			return nil
		}

		// Skip dotfiles
		if strings.HasPrefix(relPath, ".") && path != InputDir {
			l.Debug("Skipping file", log.FieldReason, "File is a hidden file (dotfile)")
			return nil
		}

		slugPath := getSlugPath(relPath)
		name, ext := obsidianmarkdown.SplitExt(relPath)
		// ext := filepath.Ext(path)
		// name := strings.TrimSuffix(info.Name(), ext) // Clean name of file
		webPath := getPageWebPath(slugPath, ext)
		outputPath := getPageOutputPath(slugPath, ext)

		switch ext {
		case ".md", ".canvas", ".base":
			// Add these to graph nodes
			scan.GraphNodes = append(scan.GraphNodes, GraphNode{
				ID:    webPath,
				Label: name,
				URL:   webPath,
				Val:   1,
			})
			info, err := info.Info()
			modTime := time.Now()
			if err == nil {
				modTime = info.ModTime()
			}
			scan.Sitemap.addEntry(modTime, BaseURL, webPath)

		default:
			// For everything else
		}

		// Adds found file to the scan
		file := &File{
			Path:    path,
			RelPath: relPath,
			Ext:     ext,
			Name:    name,
			OutPath: outputPath,
			WebPath: webPath,
		}
		scan.Files = append(scan.Files, file)

		// Register the file in the global index (filename -> public URL)
		// This is used later for resolving [[WikiLinks]]
		if _, exists := scan.FileIndex[name]; !exists {
			scan.FileIndex[name] = []*File{}
		}
		scan.FileIndex[name] = append(scan.FileIndex[name], file)

		// Used to resolve the real path of the original file (public URL -> original vault file path)
		// This is used later for resolving text embeds ![[Note#heading]]
		scan.SourceMap[webPath] = relPath

		return nil
	})

	return scan, nil
}

// loadFavicon loads the favicon.ico file if it exists
func loadFavicon() error {
	faviconSrc := filepath.Join(InputDir, "favicon.ico")
	if _, err := os.Stat(faviconSrc); err != nil {
		return err
	}
	err := copyFile(faviconSrc, filepath.Join(OutputDir, "favicon.ico"))
	if err != nil {
		return err
	}
	log.Debug("'favicon.ico' file loaded correctly")
	return nil
}

// loadCname loads the CNAME file if it exists
func loadCname() error {
	faviconSrc := filepath.Join(InputDir, "CNAME")
	if _, err := os.Stat(faviconSrc); err != nil {
		return err
	}
	err := copyFile(faviconSrc, filepath.Join(OutputDir, "CNAME"))
	if err != nil {
		return err
	}
	log.Debug("'CNAME' file loaded correctly")
	return nil
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
//
// E.g.: /Directory/Cool File.md -> /directory/cool-file.md
func getSlugPath(relPath string) string {
	parts := strings.Split(relPath, string(os.PathSeparator))
	for i, p := range parts {
		parts[i] = slugify(p)
	}
	slugPath := filepath.Join(parts...)
	return slugPath
}

// getPageWebPath returns the relative webpath of the given slugifies path
//
// E.g.: /directory/cool-file.md -> /directory/cool-file
func getPageWebPath(slugPath, ext string) string {
	var webPath string
	pathPrefix := "/"

	u, err := url.Parse(BaseURL)
	if err != nil {
		log.Fatal("Couldn't parse given base URL", "url", BaseURL, "error", err)
	}
	if u.Path != "" {
		pathPrefix = u.Path
	}

	switch ext {
	case ".md", ".canvas", ".base":
		// Handle the root homepage
		if strings.EqualFold(slugPath, "index.md") {
			return pathPrefix
		}

		// Handle other indexes
		if strings.HasSuffix(slugPath, "index.md") {
			webPath = filepath.Join(pathPrefix, strings.TrimSuffix(slugPath, "index.md"))
			break
		}

		// Remove extension
		slugPath = strings.TrimSuffix(slugPath, ext)
		webPath = filepath.Join(pathPrefix, slugPath)

	default:
		// If static file, just append together pathPrefix and slugPath
		webPath = filepath.Join(pathPrefix, slugPath)
	}

	// Replace all os specific path separator. Thanks windows
	return strings.ReplaceAll(webPath, string(os.PathSeparator), "/")
}

// getOutputPath returns the relative output path of the given slugifies path
//
// E.g. Flat urls: /directory/cool-file.md -> /directory/cool-file.html
// E.g. With directories: /directory/cool-file.md -> /directory/cool-file/index.html
func getPageOutputPath(slugPath, ext string) string {
	var outputPath string

	// Check which extension we are dealing with
	switch ext {
	case ".md", ".canvas", ".base":
		// Handles markdown, canvases and bases to become .html files
		fileName := strings.TrimSuffix(slugPath, ext)

		// Handle root index page
		if strings.EqualFold(fileName, "index") {
			outputPath = filepath.Join(OutputDir, "index.html")
			break
		}

		// Handle non flat urls by using the slug path as a folder
		if FlatUrls {
			outputPath = filepath.Join(OutputDir, fileName, "/index.html")
			break
		}

		outputPath = filepath.Join(OutputDir, fileName+".html")

	default:
		// Handles static assets by just adding the slugPath as is
		outputPath = filepath.Join(OutputDir, slugPath)
	}

	// Ensure the parent directory exists before returning
	err := os.MkdirAll(filepath.Dir(outputPath), 0755)
	if err != nil {
		log.Fatal("Couldn't create parent directory for file", "path", slugPath, "error", err)
	}

	// If flat url, just return the slugPath.html
	return outputPath
}

// Breadcrumbs constructs a list of navigation steps for the current page.
// It walks up the directory tree from the current file location.
func (f *File) Breadcrumbs() []string {
	var breadcrumbs []string
	breadcrumbs = append(breadcrumbs, "Home")

	// Add intermediate directories
	dir := filepath.Dir(f.RelPath)
	if dir != "." && dir != "" {
		breadcrumbs = append(breadcrumbs, strings.Split(dir, string(os.PathSeparator))...)
	}

	// Add current page (unless it is the home page itself)
	if !strings.EqualFold(f.Name, "index") {
		breadcrumbs = append(breadcrumbs, f.Name)
	}

	return breadcrumbs
}

// File rappresents a file that needs to be processed
type File struct {
	Path    string // Complete path of the file
	RelPath string // Relative path from input directory
	Ext     string // Extension of the file
	Name    string // Name of the file (no extension)
	OutPath string // Final output path of the file (e.g. /public/folder/page.html)
	WebPath string // Final web path of the page (e.g. /folder/page)
}

type VaultScan struct {
	FileIndex  map[string][]*File // Used to resolve wikilinks
	SourceMap  map[string]string  // Used to resolve the real disk path
	GraphNodes []GraphNode        // Lists of all pages for graph
	Files      []*File            // List of all the files found in the vault
	Sitemap    *Sitemap           // Sitemap entity
}
