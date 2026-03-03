// @feature:search Search index building for client-side full-text search.
package search

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/otaleghani/kiln/internal/obsidian"
)

type SearchEntry struct {
	Title   string   `json:"title"`
	URL     string   `json:"url"`
	Content string   `json:"content"`
	Tags    []string `json:"tags,omitempty"`
	Folder  string   `json:"folder,omitempty"`
}

var (
	reCodeBlock  = regexp.MustCompile("(?s)```.*?```")
	reImage      = regexp.MustCompile(`!\[[^\]]*\]\([^)]*\)`)
	reLink       = regexp.MustCompile(`\[([^\]]+)\]\([^)]*\)`)
	reWikiAlias  = regexp.MustCompile(`\[\[([^|\]]+)\|([^\]]+)\]\]`)
	reWikiPlain  = regexp.MustCompile(`\[\[([^\]]+)\]\]`)
	reHTML       = regexp.MustCompile(`<[^>]+>`)
	reHeading    = regexp.MustCompile(`(?m)^#{1,6}\s+`)
	reBoldItalic = regexp.MustCompile(`[*_]{1,3}`)
	reInlineCode = regexp.MustCompile("`[^`]+`")
	reWhitespace = regexp.MustCompile(`\s+`)
)

func stripMarkdown(source []byte) string {
	s := string(source)
	s = reCodeBlock.ReplaceAllString(s, "")
	s = reImage.ReplaceAllString(s, "")
	s = reLink.ReplaceAllString(s, "$1")
	s = reWikiAlias.ReplaceAllString(s, "$2")
	s = reWikiPlain.ReplaceAllString(s, "$1")
	s = reHTML.ReplaceAllString(s, "")
	s = reHeading.ReplaceAllString(s, "")
	s = reBoldItalic.ReplaceAllString(s, "")
	s = reInlineCode.ReplaceAllStringFunc(s, func(m string) string {
		return m[1 : len(m)-1]
	})
	s = reWhitespace.ReplaceAllString(s, " ")
	return strings.TrimSpace(s)
}

func BuildIndex(files []*obsidian.File) []SearchEntry {
	entries := make([]SearchEntry, 0, len(files))
	for _, f := range files {
		content := stripMarkdown(f.Content)
		if len(content) > 500 {
			content = content[:500]
		}

		var tags []string
		for tag := range f.Tags {
			tags = append(tags, tag)
		}
		if len(tags) > 0 {
			sort.Strings(tags)
		}

		entries = append(entries, SearchEntry{
			Title:   f.Name,
			URL:     f.WebPath,
			Content: content,
			Tags:    tags,
			Folder:  f.Folder,
		})
	}
	return entries
}

func WriteIndex(entries []SearchEntry, outputDir string) error {
	data, err := json.Marshal(entries)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(outputDir, "search-index.json"), data, 0644)
}
