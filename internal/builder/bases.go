package builder

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/otaleghani/kiln/internal/obsidian"
	"github.com/otaleghani/kiln/internal/obsidian/bases"
)

// func Get

func GetValue(site *DefaultSite, note *obsidian.File, field string) any {
	switch field {
	case "file":
		return note.Name
	// Metadata Mappings
	case "file.name":
		return note.Name
	case "file.fullname":
		return note.FullName
	case "file.folder":
		if note.Folder == "." {
			return "/"
		}
		return note.Folder
	case "file.ext":
		return strings.TrimPrefix(note.Ext, ".")
	case "file.basename":
		return note.Name
	case "file.path":
		return note.Path
	case "file.ctime":
		return note.Created.Format("2006-01-02") // Format dates for display
	case "file.mtime":
		return note.Modified.Format("2006-01-02")
	case "file.size":
		return note.Size
	case "file.backlinks":
		res := strings.Join(note.Backlinks, " ")
		backlinks, err := site.Markdown.RenderNote([]byte(res))
		if err != nil {
			site.log.Warn("Couldn't render backlinks")
			return ""
		}
		return backlinks
	case "file.embeds":
		res := strings.Join(note.Embeds, " ")
		embeds, err := site.Markdown.RenderNote([]byte(res))
		if err != nil {
			site.log.Warn("Couldn't render embeds")
			return ""
		}
		return embeds
	case "file.links":
		res := []byte(strings.Join(note.Links, " "))
		links, err := site.Markdown.RenderNote(res)
		if err != nil {
			site.log.Warn("Couldn't render links", "path", note.RelPath, "error", err)
			return ""
		}
		return links
	}

	// Frontmatter Fallback
	if val, ok := note.Frontmatter[field]; ok {
		return stringify(val)
	}

	// return nil // Return empty string for nil/missing values to avoid crashes
	return ""
}

// stringify converts arbitrary Frontmatter data into a displayable string.
func stringify(val any) string {
	if val == nil {
		return ""
	}

	switch v := val.(type) {
	case string:
		return v
	case int, int8, int16, int32, int64:
		return fmt.Sprintf("%d", v)
	case float32, float64:
		// %v handles floats cleanly (e.g., 1.5 stays 1.5, 1.0 becomes 1)
		return fmt.Sprintf("%v", v)
	case bool:
		return fmt.Sprintf("%t", v)
	case time.Time:
		return v.Format("2006-01-02")
	case []any:
		// Handle YAML lists (e.g. tags: [a, b]) -> "a, b"
		var parts []string
		for _, item := range v {
			parts = append(parts, stringify(item))
		}
		return strings.Join(parts, ", ")
	case []string:
		return strings.Join(v, ", ")
	default:
		// Fallback for complex maps or unknown types
		return fmt.Sprintf("%v", v)
	}
}

func FilterNotes(allFiles []*obsidian.File, baseFilters map[string][]string) []*obsidian.File {
	filteredBaseFiles := bases.FilterFiles(allFiles, baseFilters)
	return filteredBaseFiles
}

// Gets the display name. Searches for both the normal field name or note.field
func GetDisplayName(note *PageBase, field string) string {
	if property, exists := note.Properties[field]; exists {
		return property.DisplayName
	}
	if property, exists := note.Properties["note."+field]; exists {
		return property.DisplayName
	}
	return field
}

// dict allows passing multiple values to a template: {{ template "x" (dict "Key1" Val1 "Key2" Val2) }}
func dict(values ...any) (map[string]any, error) {
	if len(values)%2 != 0 {
		return nil, errors.New("dict expects an even number of arguments")
	}

	d := make(map[string]any, len(values)/2)

	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, fmt.Errorf("dict key at index %d must be a string", i)
		}
		d[key] = values[i+1]
	}

	return d, nil
}

// PageBase represents a '.base' file to be rendered
type PageBase struct {
	File *obsidian.File // Original file
	// Filters    bases.FilterGroup           `yaml:"filters"`
	Filters    map[string][]string         `yaml:"filters"`
	Properties map[string]bases.PropConfig `yaml:"properties"`
	Views      []bases.ViewConfig          `yaml:"views"`
}
