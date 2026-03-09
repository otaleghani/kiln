// @feature:rss RSS feed XML generation integrated with vault scanning.
package obsidian

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/otaleghani/kiln/internal/rss"
)

// GenerateRSS builds an RSS 2.0 feed from entries collected during vault scanning
// and writes it to feed.xml in the output directory.
func (o *Obsidian) GenerateRSS() error {
	o.log.Debug("Generating RSS feed...")

	sort.Slice(o.Vault.RSS, func(i, j int) bool {
		return o.Vault.RSS[i].PubDate.After(o.Vault.RSS[j].PubDate)
	})

	entries := o.Vault.RSS
	if len(entries) > 50 {
		entries = entries[:50]
	}

	items := make([]rss.ItemParams, 0, len(entries))
	baseURL := strings.TrimRight(o.BaseURL, "/")
	for _, entry := range entries {
		link := baseURL + entry.WebPath
		items = append(items, rss.ItemParams{
			Title:       entry.Title,
			Link:        link,
			Description: entry.Description,
			PubDate:     entry.PubDate,
			GUID:        link,
		})
	}

	xmlStr, err := rss.BuildFeedXML(rss.FeedParams{
		Title: o.BaseURL,
		Link:  o.BaseURL,
	}, items)
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(o.OutputDir, "feed.xml"), []byte(xmlStr), 0644)
}
