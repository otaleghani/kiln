package obsidian

import (
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// Slugify converts a string to a URL-friendly format.
// It lowercases the string and replaces spaces with dashes.
// Example: "My Note" -> "my-note"
func Slugify(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	return s
}

// getSlugPath converts a relative file path into a slugified path string.
// It iterates through every directory in the path and slugifies it.
//
// E.g.: /Directory/Cool File.md -> /directory/cool-file.md
func (o *Obsidian) GetSlugPath(relPath string) string {
	parts := strings.Split(relPath, string(os.PathSeparator))
	for i, p := range parts {
		parts[i] = Slugify(p)
	}
	slugPath := filepath.Join(parts...)
	return slugPath
}

// GetPageWebPath returns the relative webpath of the given slugifies path
//
// E.g.: /directory/cool-file.md -> /directory/cool-file
func (o *Obsidian) GetPageWebPath(slugPath, ext string) (string, error) {
	var webPath string
	pathPrefix := "/"

	u, err := url.Parse(o.BaseURL)
	if err != nil {
		return "", err
		// log.Fatal("Couldn't parse given base URL", "url", o.BaseURL, "error", err)
	}
	if u.Path != "" {
		pathPrefix = u.Path
	}

	switch ext {
	case ".md", ".canvas", ".base":
		// Handle the root homepage
		if strings.EqualFold(slugPath, "index.md") {
			return pathPrefix, nil
		}

		// Handle other indexes
		if strings.HasSuffix(slugPath, "index.md") {
			webPath = filepath.Join(pathPrefix, strings.TrimSuffix(slugPath, "index.md"))
			break
		}

		// Remove extension
		slugPath = strings.TrimSuffix(slugPath, ext)
		webPath = filepath.Join(pathPrefix, slugPath)

	case "":
		// Handle folder webpaths
		webPath = filepath.Join(pathPrefix, slugPath)

	default:
		// If static file, just append together pathPrefix and slugPath
		webPath = filepath.Join(pathPrefix, slugPath)
	}

	// Replace all os specific path separator. Thanks windows
	return strings.ReplaceAll(webPath, string(os.PathSeparator), "/"), nil
}

// getOutputPath returns the relative output path of the given slugifies path
//
// E.g. Flat urls: /directory/cool-file.md -> /directory/cool-file.html
// E.g. With directories: /directory/cool-file.md -> /directory/cool-file/index.html
func (o *Obsidian) GetPageOutputPath(slugPath, ext string) (string, error) {
	var outputPath string

	// Check which extension we are dealing with
	switch ext {
	case ".md", ".canvas", ".base":
		// Handles markdown, canvases and bases to become .html files
		fileName := strings.TrimSuffix(slugPath, ext)

		// Handle root index page
		if strings.EqualFold(fileName, "index") {
			outputPath = filepath.Join(o.OutputDir, "index.html")
			break
		}

		// Handle non flat urls by using the slug path as a folder
		if o.FlatURLs {
			outputPath = filepath.Join(o.OutputDir, fileName, "/index.html")
			break
		}

		outputPath = filepath.Join(o.OutputDir, fileName+".html")

	case "":
		// Handle folder webpaths
		outputPath = filepath.Join(o.OutputDir, slugPath, "/index.html")

	default:
		// Handles static assets by just adding the slugPath as is
		outputPath = filepath.Join(o.OutputDir, slugPath)
	}

	// Ensure the parent directory exists before returning
	// TODO: Ensure not to craete a path when a file doesn't end in an extension.
	err := os.MkdirAll(filepath.Dir(outputPath), 0755)
	if err != nil {
		return "", err
		// log.Fatal("Couldn't create parent directory for file", "path", slugPath, "error", err)
	}

	// If flat url, just return the slugPath.html
	return outputPath, nil
}

// getFolderWebPath calculates the web path for a directory.
// Per instructions: /folder -> /folder/index.html
func (o *Obsidian) getFolderWebPath(folderRelPath string) (string, error) {
	pathPrefix := "/"

	// Reuse the BaseURL logic from your existing code
	if o.BaseURL != "" {
		u, err := url.Parse(o.BaseURL)
		if err != nil {
			return "", err
			// log.Printf("Warning: Couldn't parse BaseURL: %v", err)
		} else if u.Path != "" {
			pathPrefix = u.Path
		}
	}

	// 1. Join the prefix and the folder path
	// 2. Append "index.html" as requested
	webPath := filepath.Join(pathPrefix, folderRelPath)

	// Replace OS separators with forward slashes for the web
	return strings.ReplaceAll(webPath, string(os.PathSeparator), "/"), nil
}

// getTagWebPath calculates the web path for a tag
// E.g.: #example -> /tags/example
func (o *Obsidian) getTagWebPath(tagName string) (string, error) {
	var webPath string
	pathPrefix := "/"

	u, err := url.Parse(o.BaseURL)
	if err != nil {
		return "", err
	}
	if u.Path != "" {
		pathPrefix = u.Path
	}

	name := strings.TrimPrefix(tagName, "#")
	webPath = filepath.Join(pathPrefix, "tags", name)

	// Replace all os specific path separator. Thanks windows
	return strings.ReplaceAll(webPath, string(os.PathSeparator), "/"), nil
}

// getTagWebPath calculates the output path for a tag
// E.g. Flat urls: #example -> /tags/example.html
// E.g. With directories: #example -> /tags/example/index.html
func (o *Obsidian) getTagOutputPath(tagName string) (string, error) {
	var outputPath string

	// Handles markdown, canvases and bases to become .html files
	name := strings.TrimPrefix(tagName, "#")

	// Handle non flat urls by using the slug path as a folder
	if o.FlatURLs {
		outputPath = filepath.Join(o.OutputDir, "/tags", name, "/index.html")
	} else {
		outputPath = filepath.Join(o.OutputDir, "/tags", name+".html")
	}

	// Ensure the parent directory exists before returning
	err := os.MkdirAll(filepath.Dir(outputPath), 0755)
	if err != nil {
		return "", err
	}

	return outputPath, nil
}

// Breadcrumbs constructs a list of navigation steps (Breadcrumb structs) for the current page.
func (o *Obsidian) GetBreadcrumbs(f *File) ([]Breadcrumb, error) {
	var crumbs []Breadcrumb

	// Resolve base path for the "Home" link
	homePath := "/"
	if o.BaseURL != "" {
		u, _ := url.Parse(o.BaseURL)
		if u != nil && u.Path != "" {
			homePath = u.Path
		}
	}

	// Add "Home"
	crumbs = append(crumbs, Breadcrumb{
		Label: "Home",
		Url:   homePath,
	})

	// Add intermediate directories
	dir := filepath.Dir(f.RelPath)

	// Check if we are in a subfolder (ignore "." or empty root)
	if dir != "." && dir != "" {
		// Split the directory into parts (e.g., "docs/guides" -> ["docs", "guides"])
		parts := strings.Split(dir, string(os.PathSeparator))

		// Accumulator to build the path for every parent folder
		// e.g., first "docs", then "docs/guides"
		currentPath := ""

		for _, part := range parts {
			if currentPath == "" {
				currentPath = part
			} else {
				currentPath = filepath.Join(currentPath, part)
			}

			url, err := o.getFolderWebPath(o.GetSlugPath(currentPath))
			if err != nil {
				return []Breadcrumb{}, err
			}

			crumbs = append(crumbs, Breadcrumb{
				Label: part,
				Url:   url,
			})
		}
	}

	// Add current page
	// If the current file is "index", it conceptually represents the folder
	// we just added, so we don't add it again.
	if !strings.EqualFold(f.Name, "index") {
		crumbs = append(crumbs, Breadcrumb{
			Label: f.Name,
			Url:   f.WebPath,
		})
	}

	return crumbs, nil
}

// GetFolderBreadcrumbs constructs a list of navigation steps for the current folder.
func (o *Obsidian) GetFolderBreadcrumbs(f *Folder) ([]Breadcrumb, error) {
	var crumbs []Breadcrumb

	// Resolve base path for the "Home" link
	// (This logic remains identical to the File version)
	homePath := "/"
	if o.BaseURL != "" {
		u, _ := url.Parse(o.BaseURL)
		if u != nil && u.Path != "" {
			homePath = u.Path
		}
	}

	// Add "Home"
	crumbs = append(crumbs, Breadcrumb{
		Label: "Home",
		Url:   homePath,
	})

	// Process the Folder Path
	if f.RelPath != "." && f.RelPath != "" {
		// Split the directory into parts (e.g., "docs/guides" -> ["docs", "guides"])
		parts := strings.Split(f.RelPath, string(os.PathSeparator))

		// Accumulator to build the path for every level
		currentPath := ""

		for _, part := range parts {
			if currentPath == "" {
				currentPath = part
			} else {
				currentPath = filepath.Join(currentPath, part)
			}

			// Resolve the web path for this specific folder level
			url, err := o.getFolderWebPath(o.GetSlugPath(currentPath))
			if err != nil {
				return []Breadcrumb{}, err
			}

			crumbs = append(crumbs, Breadcrumb{
				Label: part,
				Url:   url,
			})
		}
	}

	return crumbs, nil
}

type Breadcrumb struct {
	Label string
	Url   string
}
