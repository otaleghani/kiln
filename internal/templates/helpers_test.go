// @feature:layouts Tests for template helper functions.
package templates

import (
	"bytes"
	"context"
	"testing"

	"github.com/otaleghani/kiln/internal/obsidian"
)

func TestToStr(t *testing.T) {
	tests := []struct {
		name string
		val  any
		want string
	}{
		{"nil value", nil, ""},
		{"string value", "hello", "hello"},
		{"empty string", "", ""},
		{"int value", 42, "42"},
		{"bool value", true, "true"},
		{"float value", 3.14, "3.14"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := toStr(tt.val)
			if got != tt.want {
				t.Errorf("toStr(%v) = %q, want %q", tt.val, got, tt.want)
			}
		})
	}
}

func TestHead_NilFrontmatterValues(t *testing.T) {
	tests := []struct {
		name        string
		frontmatter map[string]any
	}{
		{
			name:        "nil title and description",
			frontmatter: map[string]any{"title": nil, "description": nil},
		},
		{
			name:        "nil frontmatter map",
			frontmatter: nil,
		},
		{
			name:        "empty frontmatter",
			frontmatter: map[string]any{},
		},
		{
			name:        "valid title, nil description",
			frontmatter: map[string]any{"title": "My Page", "description": nil},
		},
		{
			name:        "nil title, valid description",
			frontmatter: map[string]any{"title": nil, "description": "A description"},
		},
		{
			name:        "non-string title",
			frontmatter: map[string]any{"title": 123},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := &PageData{
				File:        &obsidian.File{Name: "test-note"},
				Frontmatter: tt.frontmatter,
				Site: &SiteData{
					SiteName: "Test Site",
				},
			}

			var buf bytes.Buffer
			err := Head(data).Render(context.Background(), &buf)
			if err != nil {
				t.Fatalf("Head() returned error: %v", err)
			}
		})
	}
}
