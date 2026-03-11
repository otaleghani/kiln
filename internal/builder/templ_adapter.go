// @feature:builder Adapter functions to convert builder types to template types.
package builder

import (
	"fmt"
	"sort"
	"strings"

	"github.com/otaleghani/kiln/internal/i18n"
	"github.com/otaleghani/kiln/internal/obsidian"
	"github.com/otaleghani/kiln/internal/templates"
)

// toTemplPageData maps a DefaultSitePageData to the templ-compatible PageData.
func toTemplPageData(p *DefaultSitePageData) *templates.PageData {
	data := &templates.PageData{
		Content:     string(p.Content),
		TOC:         string(p.TOC),
		CanvasJSON:  string(p.CanvasContent),
		Breadcrumbs: p.Breadcrumbs,
		File:        p.File,
		Folder:      p.Folder,
		Tag:         p.Tag,
		IsGraph:     p.IsGraph,
		IsCanvas:    p.IsCanvas,
		IsBase:      p.IsBase,
		IsNote:      p.IsNote,
		IsFolder:    p.IsFolder,
		IsTag:       p.IsTag,
		Is404:       p.Is404,
		Frontmatter: p.Frontmatter,
		Site: &templates.SiteData{
			BaseURL:           p.Site.BaseURL,
			SiteName:          p.Site.SiteName,
			Theme:             toTemplTheme(p.Site.Theme),
			NavbarRoot:        p.Site.NavbarRoot,
			DisableLocalGraph: p.Site.DisableLocalGraph,
			DisableTOC:        p.Site.DisableTOC,
			DisableBacklinks:  p.Site.DisableBacklinks,
			FlatURLs:          p.Site.FlatURLs,
			Lang:              p.Site.Lang,
			Labels:            i18n.Resolve(p.Site.Lang),
		},
	}

	if p.IsNote && p.File != nil {
		wc := templates.WordCount(p.File.Content)
		tags := make([]string, 0, len(p.File.Tags))
		for t := range p.File.Tags {
			tags = append(tags, t)
		}
		sort.Strings(tags)
		data.Meta = &templates.NoteMeta{
			WordCount:   wc,
			ReadingTime: templates.ReadingTimeFromWords(wc),
			Created:     p.File.Created,
			Modified:    p.File.Modified,
			Tags:        tags,
		}

		if !p.Site.DisableBacklinks && len(p.File.Backlinks) > 0 {
			for _, bl := range p.File.Backlinks {
				name := strings.TrimPrefix(bl, "[[")
				name = strings.TrimSuffix(name, "]]")
				if files, ok := p.Site.Obsidian.Vault.FileIndex[name]; ok && len(files) > 0 {
					data.Backlinks = append(data.Backlinks, templates.Backlink{
						Name:    name,
						WebPath: files[0].WebPath,
					})
				}
			}
		}
	}

	if p.IsBase && p.Base.File != nil && len(p.Base.File.Views) > 0 {
		data.Base = templates.BaseViewData{
			Groups:  p.Base.Groups,
			Notes:   p.Base.Notes,
			Columns: p.Base.Columns,
			ViewType: p.Base.File.Views[0].Type,
			DisplayNameFn: func(field string) string {
				return GetDisplayName(p.Base.File, field)
			},
			ValueFn: func(note *obsidian.File, field string) string {
				return fmt.Sprintf("%v", GetValue(p.Site, note, field))
			},
		}
	}

	return data
}

// toTemplTheme maps a builder Theme to the templ-compatible ThemeData.
func toTemplTheme(t *Theme) *templates.ThemeData {
	return &templates.ThemeData{
		Light:       toTemplThemeColors(t.Light),
		Dark:        toTemplThemeColors(t.Dark),
		FontFamily:  string(t.Font.Family),
		FontFaceCSS: string(t.Font.FontFaceReplaced),
	}
}

func toTemplThemeColors(c *ThemeColors) *templates.ThemeColors {
	return &templates.ThemeColors{
		Bg:            c.Bg,
		Text:          c.Text,
		SidebarBg:     c.SidebarBg,
		SidebarBorder: c.SidebarBorder,
		Accent:        c.Accent,
		Hover:         c.Hover,
		Comment:       c.Comment,
		Red:           c.Red,
		Orange:        c.Orange,
		Yellow:        c.Yellow,
		Green:         c.Green,
		Blue:          c.Blue,
		Purple:        c.Purple,
		Cyan:          c.Cyan,
	}
}
