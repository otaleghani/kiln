package linter

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// CollectNotes scans the input directory and creates an index of all valid files.
// It returns a map acting as a set of file paths relative to the root.
// For Markdown files, it stores both the full filename ("note.md") and the base name ("note")
// to allow linking without extensions.
func CollectNotes(inputDir string) map[string]bool {
	validFiles := make(map[string]bool)

	filepath.WalkDir(inputDir, func(path string, d fs.DirEntry, err error) error {
		// Skip directories and hidden files (starting with .)
		if !d.IsDir() && !strings.HasPrefix(d.Name(), ".") {
			rel, _ := filepath.Rel(inputDir, path)

			// Register the exact file path (e.g., "images/logo.png" or "notes/idea.md")
			validFiles[rel] = true

			// For Markdown files, also register the path without the extension.
			// This supports the common Obsidian style of linking [[Note Name]] instead of [[Note Name.md]].
			if filepath.Ext(rel) == ".md" {
				noExt := strings.TrimSuffix(rel, ".md")
				validFiles[noExt] = true
			}
		}
		return nil
	})
	return validFiles
}

// BrokenLinks iterates through all Markdown files in the directory and validates their WikiLinks.
// It reports any links that point to files not present in the 'notes' index.
func BrokenLinks(inputDir string, notes map[string]bool) {
	// Regex to find standard WikiLink syntax: [[Target]]
	linkRegex := regexp.MustCompile(`\[\[(.*?)\]\]`)
	issuesFound := 0

	filepath.WalkDir(inputDir, func(path string, d fs.DirEntry, err error) error {
		// specific checks: Only scan .md files for links
		if err != nil || d.IsDir() || filepath.Ext(path) != ".md" {
			return nil
		}

		content, _ := os.ReadFile(path)
		matches := linkRegex.FindAllStringSubmatch(string(content), -1)

		for _, match := range matches {
			rawLink := match[1] // The content inside brackets: "Note Name", "Note|Alias", or "Note#Header"

			// Handle Aliasing: [[Note Name|Custom Text]] -> we only care about "Note Name"
			if strings.Contains(rawLink, "|") {
				rawLink = strings.Split(rawLink, "|")[0]
			}

			// Handle Anchors: [[Note Name#Section]] -> we only care about "Note Name"
			if strings.Contains(rawLink, "#") {
				rawLink = strings.Split(rawLink, "#")[0]
			}

			// If the link was just an anchor to the current page (e.g., [[#Top]]), skip validation.
			if rawLink == "" {
				continue
			}

			// Resolution Strategy:
			// Check if the link exists in the index using various common suffixes.
			// 1. Exact match (e.g., "Folder/Note")
			// 2. Markdown extension (e.g., "Folder/Note" -> "Folder/Note.md")
			// 3. Image extensions (e.g., "image" -> "image.png")
			exists := false
			candidates := []string{
				rawLink,
				rawLink + ".md",
				rawLink + ".png",
				rawLink + ".jpg",
			}

			for _, c := range candidates {
				if notes[c] {
					exists = true
					break
				}
			}

			if !exists {
				log.Printf("Broken link in [%s]: [[%s]]\n", d.Name(), rawLink)
				issuesFound++
			}
		}
		return nil
	})

	if issuesFound == 0 {
		log.Println("No broken links found")
	} else {
		log.Printf("Found %d broken links.", issuesFound)
	}
}
