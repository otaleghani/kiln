package builder

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"
)

func Stats() {
	var noteCount, wordCount, maxWords int
	var longestNote string

	filepath.WalkDir(InputDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || filepath.Ext(path) != ".md" {
			return nil
		}

		content, _ := os.ReadFile(path)
		words := len(strings.Fields(string(content)))

		noteCount++
		wordCount += words
		if words > maxWords {
			maxWords = words
			longestNote = d.Name()
		}
		return nil
	})

	// Use text/tabwriter for a nice table
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "METRIC\tVALUE")
	fmt.Fprintln(w, "------\t-----")
	fmt.Fprintf(w, "Total Notes\t%d\n", noteCount)
	fmt.Fprintf(w, "Total Words\t%d\n", wordCount)
	fmt.Fprintf(w, "Longest Note\t%s (%d words)\n", longestNote, maxWords)
	w.Flush()
}
