// @feature:layouts Shared template data types for templ-based rendering.
package templates

import (
	"time"

	"github.com/otaleghani/kiln/internal/i18n"
	"github.com/otaleghani/kiln/internal/obsidian"
	"github.com/otaleghani/kiln/internal/obsidian/bases"
)

// PageData is the top-level data passed to every page template.
type PageData struct {
	Site        *SiteData
	Content     string
	TOC         string
	CanvasJSON  string
	Breadcrumbs []obsidian.Breadcrumb
	File        *obsidian.File
	Folder      *obsidian.Folder
	Tag         *obsidian.Tag
	IsGraph     bool
	IsCanvas    bool
	IsBase      bool
	IsNote      bool
	IsFolder    bool
	IsTag       bool
	Is404       bool
	Frontmatter map[string]any
	Meta        *NoteMeta
	Backlinks   []Backlink
	Base        BaseViewData
}

// NoteMeta holds reading metadata for a note page.
type NoteMeta struct {
	WordCount   int
	ReadingTime int // minutes
	Created     time.Time
	Modified    time.Time
	Tags        []string
}

// Backlink represents a resolved incoming link to the current page.
type Backlink struct {
	Name    string
	WebPath string
}

// SiteData holds global site configuration used across all pages.
type SiteData struct {
	BaseURL           string
	SiteName          string
	Theme             *ThemeData
	NavbarRoot        *obsidian.NavbarNode
	DisableLocalGraph bool
	DisableTOC        bool
	DisableBacklinks  bool
	FlatURLs          bool
	Lang              string
	Labels            *i18n.Labels
}

// ThemeData bundles color schemes and typography for the site.
type ThemeData struct {
	Light       *ThemeColors
	Dark        *ThemeColors
	FontFamily  string
	FontFaceCSS string
}

// ThemeColors defines the color palette for a single mode (light or dark).
type ThemeColors struct {
	Bg            string
	Text          string
	SidebarBg     string
	SidebarBorder string
	Accent        string
	Hover         string
	Comment       string
	Red           string
	Orange        string
	Yellow        string
	Green         string
	Blue          string
	Purple        string
	Cyan          string
}

// BaseViewData holds data for base/database-style views.
type BaseViewData struct {
	Groups        []*bases.FileGroup
	Notes         []*obsidian.File
	Columns       []string
	ViewType      string
	DisplayNameFn func(string) string
	ValueFn       func(*obsidian.File, string) string
}
