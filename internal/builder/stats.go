package builder

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"
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

	// Print the results in a clean, aligned table
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "METRIC\tVALUE")
	fmt.Fprintln(w, "------\t-----")
	fmt.Fprintf(w, "Total Notes\t%d\n", noteCount)
	fmt.Fprintf(w, "Total Words\t%d\n", wordCount)
	fmt.Fprintf(w, "Longest Note\t%s (%d words)\n", longestNote, maxWords)
	w.Flush()
}
