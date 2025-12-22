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
func CollectNotes(inputDir string) map[string]bool {
	validFiles := make(map[string]bool)

	filepath.WalkDir(inputDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		// Skip directories and hidden files
		if strings.HasPrefix(d.Name(), ".") && path != inputDir {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if !d.IsDir() {
			rel, _ := filepath.Rel(inputDir, path)

			// 1. Index full relative path (e.g., "folder/note.md")
			validFiles[rel] = true

			// 2. Index relative path without extension (e.g., "folder/note")
			noExtRel := strings.TrimSuffix(rel, filepath.Ext(rel))
			validFiles[noExtRel] = true

			// --- THE FIX ---
			// 3. Index the filename itself (e.g., "note.md")
			// This allows [[note.md]] to match, even if it's inside a folder.
			validFiles[d.Name()] = true

			// 4. Index the filename without extension (e.g., "note")
			// This allows [[note]] to match, even if it's inside a folder.
			if filepath.Ext(d.Name()) == ".md" || filepath.Ext(d.Name()) == ".canvas" {
				noExtName := strings.TrimSuffix(d.Name(), filepath.Ext(d.Name()))
				validFiles[noExtName] = true
			}
		}
		return nil
	})
	return validFiles
}

// BrokenLinks iterates through all Markdown files in the directory and validates their WikiLinks.
func BrokenLinks(inputDir string, notes map[string]bool) {
	// Regex to find standard WikiLink syntax: [[Target]]
	linkRegex := regexp.MustCompile(`\[\[(.*?)\]\]`)
	issuesFound := 0

	filepath.WalkDir(inputDir, func(path string, d fs.DirEntry, err error) error {
		// Skip hidden folders
		if strings.HasPrefix(d.Name(), ".") && d.IsDir() {
			return filepath.SkipDir
		}

		// Only scan .md files for links
		if err != nil || d.IsDir() || filepath.Ext(path) != ".md" {
			return nil
		}

		content, _ := os.ReadFile(path)
		matches := linkRegex.FindAllStringSubmatch(string(content), -1)

		for _, match := range matches {
			rawLink := match[1]

			// Handle Aliasing: [[Note Name|Custom Text]] -> "Note Name"
			if strings.Contains(rawLink, "|") {
				rawLink = strings.Split(rawLink, "|")[0]
			}

			// Handle Anchors: [[Note Name#Section]] -> "Note Name"
			if strings.Contains(rawLink, "#") {
				rawLink = strings.Split(rawLink, "#")[0]
			}

			if rawLink == "" {
				continue
			}

			// Check candidates against the now robust 'notes' map
			exists := false
			candidates := []string{
				rawLink,             // "Note"
				rawLink + ".md",     // "Note.md"
				rawLink + ".png",    // "Note.png"
				rawLink + ".jpg",    // "Note.jpg"
				rawLink + ".jpeg",   // "Note.jpeg"
				rawLink + ".canvas", // "Note.canvas"
			}

			for _, c := range candidates {
				if notes[c] {
					exists = true
					break
				}
			}

			if !exists {
				// Use Rel path for cleaner logging
				relPath, _ := filepath.Rel(inputDir, path)
				log.Printf("Broken link in [%s]: [[%s]]\n", relPath, rawLink)
				issuesFound++
			}
		}
		return nil
	})

	if issuesFound == 0 {
		log.Println("No broken links found")
	} else {
		// Optional: Fail the build if broken links are critical
		// os.Exit(1)
		log.Printf("Found %d broken links.", issuesFound)
	}
}
