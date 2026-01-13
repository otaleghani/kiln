package obsidian

import (
	"encoding/xml"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// generateRobots creates a robots.txt file in the output directory.
// It points crawlers to the Sitemap location.
func (o *Obsidian) GenerateRobots() error {
	robotsFile, err := os.Create(filepath.Join(o.OutputDir, "robots.txt"))
	if err != nil {
		return err
	}
	defer robotsFile.Close()
	robotsFile.WriteString("User-agent: *\n")
	robotsFile.WriteString("Allow: /\n")
	robotsFile.WriteString("Sitemap: " + o.Vault.Sitemap.Path + "\n")
	return nil
}

func (o *Obsidian) GenerateSitemap() error {
	sitemapFile, err := os.Create(filepath.Join(o.OutputDir, "sitemap.xml"))
	if err != nil {
		return err
	}
	defer sitemapFile.Close()
	sitemapFile.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	sitemapFile.WriteString(
		`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">` + "\n",
	)
	for _, entry := range o.Vault.Sitemap.Entries {
		output, _ := xml.MarshalIndent(entry, "  ", "  ")
		sitemapFile.Write(output)
		sitemapFile.WriteString("\n")
	}
	sitemapFile.WriteString(`</urlset>`)

	return nil
}

// addEntry appends a new entry to the sitemap slice.
// It retrieves the file's modification time to populate 'lastmod'.
func (o *Obsidian) AddSitemapEntry(modTime time.Time, baseURL, webPath string) {
	fullURL := strings.TrimRight(baseURL, "/") + webPath
	o.Vault.Sitemap.Entries = append(o.Vault.Sitemap.Entries, SitemapEntry{
		Loc:     fullURL,
		LastMod: modTime.Format("2006-01-02"),
	})
	log.Println("Added to sitemap", "url", webPath)
}

// SitemapEntry represents a single URL entry in the sitemap.xml.
type SitemapEntry struct {
	XMLName xml.Name `xml:"url"`
	Loc     string   `xml:"loc"`     // The absolute URL
	LastMod string   `xml:"lastmod"` // The last modification date
}

// Sitemap holds all of the entries to generate the sitemap
type Sitemap struct {
	Entries []SitemapEntry // Entries of the sitemap
	Path    string         // Path of the sitemap
}
