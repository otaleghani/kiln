package linter

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func CollectNotes(inputDir string) map[string]bool {
	validFiles := make(map[string]bool)
	filepath.WalkDir(inputDir, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() && !strings.HasPrefix(d.Name(), ".") {
			rel, _ := filepath.Rel(inputDir, path)
			validFiles[rel] = true
			// Also handle "Note Name" matching "Note Name.md"
			if filepath.Ext(rel) == ".md" {
				noExt := strings.TrimSuffix(rel, ".md")
				validFiles[noExt] = true
			}
		}
		return nil
	})
	return validFiles
}

func BrokenLinks(inputDir string, notes map[string]bool) {
	// Scan for links
	linkRegex := regexp.MustCompile(`\[\[(.*?)\]\]`)
	issuesFound := 0

	filepath.WalkDir(inputDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || filepath.Ext(path) != ".md" {
			return nil
		}

		content, _ := os.ReadFile(path)
		matches := linkRegex.FindAllStringSubmatch(string(content), -1)

		for _, match := range matches {
			rawLink := match[1] // "Note Name" or "Note|Alias"

			// Remove Alias
			if strings.Contains(rawLink, "|") {
				rawLink = strings.Split(rawLink, "|")[0]
			}
			// Remove Anchor
			if strings.Contains(rawLink, "#") {
				rawLink = strings.Split(rawLink, "#")[0]
			}

			if rawLink == "" {
				continue
			} // Was just an anchor [[#Top]]

			// Check existence
			// Try exact match, or adding .md, or adding extensions for assets
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
				fmt.Printf("Broken link in [%s]: [[%s]]\n", d.Name(), rawLink)
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
