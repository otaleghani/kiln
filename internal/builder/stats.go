package builder

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/otaleghani/kiln/internal/log"
)

// Stats calculates and prints summary statistics for the vault (count, words, etc.).
func Stats() {
	var noteCount, wordCount, maxWords int
	var longestNote string

	// Walk through the vault to gather metrics
	filepath.WalkDir(InputDir, func(path string, d fs.DirEntry, err error) error {
		// Skip errors, directories, and non-markdown files
		if err != nil || d.IsDir() || filepath.Ext(path) != ".md" {
			return nil
		}

		content, _ := os.ReadFile(path)
		// Simple word count estimation using whitespace splitting
		words := len(strings.Fields(string(content)))

		noteCount++
		wordCount += words

		// Track the longest note
		if words > maxWords {
			maxWords = words
			longestNote = d.Name()
		}
		return nil
	})

	log.Info("Total notes", "notes", noteCount)
	log.Info("Total words", "words", wordCount)
	log.Info("Longest note", "note", longestNote, "words", maxWords)
}
