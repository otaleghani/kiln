// Vault diagnostics that scans for broken wikilinks across markdown files. @feature:linter
package linter

import (
	"io/fs"
	"log/slog"
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

			// Index full relative path (e.g., "folder/note.md")
			validFiles[rel] = true

			// Index relative path without extension (e.g., "folder/note")
			noExtRel := strings.TrimSuffix(rel, filepath.Ext(rel))
			validFiles[noExtRel] = true

			// Index the filename itself (e.g., "note.md")
			// This allows [[note.md]] to match, even if it's inside a folder.
			validFiles[d.Name()] = true

			// Index the filename without extension (e.g., "note")
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

func checkCandidates(rawLink string, notes map[string]bool) bool {
	candidates := []string{
		rawLink,
		rawLink + ".md",
		rawLink + ".png",
		rawLink + ".jpg",
		rawLink + ".jpeg",
		rawLink + ".canvas",
	}
	for _, c := range candidates {
		if notes[c] {
			return true
		}
	}
	return false
}

// BrokenLinks iterates through all Markdown files in the directory and validates their links.
func BrokenLinks(inputDir string, notes map[string]bool, log *slog.Logger) {
	wikiLinkRegex := regexp.MustCompile(`\[\[(.*?)\]\]`)
	mdLinkRegex := regexp.MustCompile(`\[([^\]]*)\]\(([^)]+)\)`)
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
		relPath, _ := filepath.Rel(inputDir, path)

		// Check wikilinks
		for _, match := range wikiLinkRegex.FindAllStringSubmatch(string(content), -1) {
			rawLink := match[1]

			if strings.Contains(rawLink, "|") {
				rawLink = strings.Split(rawLink, "|")[0]
			}
			if strings.Contains(rawLink, "#") {
				rawLink = strings.Split(rawLink, "#")[0]
			}
			if rawLink == "" {
				continue
			}

			if !checkCandidates(rawLink, notes) {
				log.Warn("Found broken link", "path", relPath, "link", rawLink)
				issuesFound++
			}
		}

		// Check markdown links
		for _, match := range mdLinkRegex.FindAllStringSubmatch(string(content), -1) {
			linkPath := match[2]

			if strings.HasPrefix(linkPath, "http://") ||
				strings.HasPrefix(linkPath, "https://") ||
				strings.HasPrefix(linkPath, "mailto:") {
				continue
			}
			if strings.HasPrefix(linkPath, "#") {
				continue
			}

			// Strip anchor
			if idx := strings.Index(linkPath, "#"); idx != -1 {
				linkPath = linkPath[:idx]
			}

			// Resolve relative to the directory of the current file
			resolved := filepath.Clean(filepath.Join(filepath.Dir(relPath), linkPath))

			if !checkCandidates(resolved, notes) {
				log.Warn("Found broken link", "path", relPath, "link", linkPath)
				issuesFound++
			}
		}

		return nil
	})

	if issuesFound == 0 {
		log.Info("No broken links found")
	} else {
		// Optional: Fail the build if broken links are critical
		// os.Exit(1)
		log.Error("Found broken links", "number", issuesFound)
	}
}
