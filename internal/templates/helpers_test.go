// @feature:layouts Tests for template helper functions.
package templates

import (
	"bytes"
	"context"
	"strings"
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

func TestHead_OGMetaTags(t *testing.T) {
	data := &PageData{
		File: &obsidian.File{
			Name:    "my-post",
			WebPath: "/blog/my-post",
		},
		Frontmatter: map[string]any{
			"title":       "My Title",
			"description": "My Desc",
		},
		Site: &SiteData{
			BaseURL:  "https://example.com",
			SiteName: "My Site",
		},
	}

	var buf bytes.Buffer
	err := Head(data).Render(context.Background(), &buf)
	if err != nil {
		t.Fatalf("Head() returned error: %v", err)
	}

	html := buf.String()
	expected := []string{
		`<meta property="og:title" content="My Title"`,
		`<meta property="og:description" content="My Desc"`,
		`<meta property="og:image" content="https://example.com/blog/my-post/og.png"`,
		`<meta property="og:url" content="https://example.com/blog/my-post"`,
		`<meta property="og:type" content="article"`,
		`<meta name="twitter:card" content="summary_large_image"`,
		`<meta name="twitter:title" content="My Title"`,
		`<meta name="twitter:description" content="My Desc"`,
		`<meta name="twitter:image" content="https://example.com/blog/my-post/twitter.png"`,
	}

	for _, want := range expected {
		if !strings.Contains(html, want) {
			t.Errorf("expected HTML to contain %q, got:\n%s", want, html)
		}
	}
}

func TestHead_OGMetaTags_Folder(t *testing.T) {
	data := &PageData{
		IsFolder: true,
		Folder: &obsidian.Folder{
			RelPath: "my-folder",
			WebPath: "/my-folder",
		},
		Site: &SiteData{
			BaseURL:  "https://example.com",
			SiteName: "My Site",
		},
	}

	var buf bytes.Buffer
	err := Head(data).Render(context.Background(), &buf)
	if err != nil {
		t.Fatalf("Head() returned error: %v", err)
	}

	html := buf.String()
	expected := []string{
		`<meta property="og:title" content="my-folder"`,
		`<meta property="og:type" content="website"`,
		`<meta property="og:url" content="https://example.com/my-folder"`,
	}

	for _, want := range expected {
		if !strings.Contains(html, want) {
			t.Errorf("expected HTML to contain %q, got:\n%s", want, html)
		}
	}
}

func TestHead_OGMetaTags_Tag(t *testing.T) {
	data := &PageData{
		IsTag: true,
		Tag: &obsidian.Tag{
			Name:    "golang",
			WebPath: "/tags/golang",
		},
		Site: &SiteData{
			BaseURL:  "https://example.com",
			SiteName: "My Site",
		},
	}

	var buf bytes.Buffer
	err := Head(data).Render(context.Background(), &buf)
	if err != nil {
		t.Fatalf("Head() returned error: %v", err)
	}

	html := buf.String()
	expected := []string{
		`<meta property="og:title" content="golang"`,
		`<meta property="og:type" content="website"`,
		`<meta property="og:url" content="https://example.com/tags/golang"`,
	}

	for _, want := range expected {
		if !strings.Contains(html, want) {
			t.Errorf("expected HTML to contain %q, got:\n%s", want, html)
		}
	}
}

func TestHead_OGMetaTags_NoDescription(t *testing.T) {
	data := &PageData{
		File: &obsidian.File{
			Name:    "my-post",
			WebPath: "/blog/my-post",
		},
		Frontmatter: map[string]any{
			"title": "My Title",
		},
		Site: &SiteData{
			BaseURL:  "https://example.com",
			SiteName: "My Site",
		},
	}

	var buf bytes.Buffer
	err := Head(data).Render(context.Background(), &buf)
	if err != nil {
		t.Fatalf("Head() returned error: %v", err)
	}

	html := buf.String()
	want := `<meta property="og:description" content="Notes on my-post"`
	if !strings.Contains(html, want) {
		t.Errorf("expected HTML to contain %q, got:\n%s", want, html)
	}
}
