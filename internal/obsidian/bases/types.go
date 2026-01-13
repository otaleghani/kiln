package bases

import "github.com/otaleghani/kiln/internal/obsidian"

type FilterGroup struct {
	Operator   string      `yaml:"operator"` // "and", "or" (simplified for this example)
	Conditions []Condition `yaml:"conditions"`
}

// Individual filter condition (e.g., status == "done")
type Condition struct {
	Field    string `yaml:"field"`    // e.g., "file.folder" or "status"
	Operator string `yaml:"operator"` // "==", "!=", "contains", ">"
	Value    any    `yaml:"value"`
}

type ViewConfig struct {
	Type    string              `yaml:"type"` // "table", "list", "cards"
	Name    string              `yaml:"name"`
	Order   []string            `yaml:"order"`
	Filters map[string][]string `yaml:"filters"` // View-specific filters
	Sort    []SortConfig        `yaml:"sort"`
	GroupBy GroupConfig         `yaml:"groupBy"`
}

type GroupConfig struct {
	Property  string `yaml:"property"`
	Direction string `yaml:"direction"`
}

// SortConfig handles the complex object structure of the sort field
type SortConfig struct {
	Property  string `yaml:"property"`
	Direction string `yaml:"direction"`
}

type ColumnDef struct {
	Field string `yaml:"field"` // The property key to show
	Title string `yaml:"title"` // Optional override
}

type PropConfig struct {
	DisplayName string `yaml:"displayName"`
}

// FileGroup represents one section of grouped notes
type FileGroup struct {
	Key   string           // The value we grouped by (e.g., "book", "2023-01")
	Notes []*obsidian.File // The files in this group
}
