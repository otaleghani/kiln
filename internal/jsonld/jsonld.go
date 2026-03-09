// @feature:jsonld JSON-LD structured data generation for Article and BreadcrumbList schemas.
package jsonld

import (
	"encoding/json"
	"time"
)

// ArticleParams holds the inputs needed to build an Article JSON-LD object.
type ArticleParams struct {
	Title        string
	Description  string
	Author       string
	BaseURL      string
	WebPath      string
	SiteName     string
	DateCreated  time.Time
	DateModified time.Time
	ImageURL     string
}

// BreadcrumbItem represents one breadcrumb for JSON-LD generation.
type BreadcrumbItem struct {
	Label string
	URL   string
}

// BuildArticleJSON returns a JSON-LD string for schema.org/Article.
// Returns empty string if title is empty.
func BuildArticleJSON(p ArticleParams) string {
	if p.Title == "" {
		return ""
	}

	schema := map[string]any{
		"@context":      "https://schema.org",
		"@type":         "Article",
		"headline":      p.Title,
		"datePublished": p.DateCreated.Format(time.RFC3339),
		"dateModified":  p.DateModified.Format(time.RFC3339),
		"publisher": map[string]any{
			"@type": "Organization",
			"name":  p.SiteName,
		},
		"mainEntityOfPage": map[string]any{
			"@type": "WebPage",
			"@id":   p.BaseURL + p.WebPath,
		},
	}

	if p.Description != "" {
		schema["description"] = p.Description
	}
	if p.Author != "" {
		schema["author"] = map[string]any{
			"@type": "Person",
			"name":  p.Author,
		}
	}
	if p.ImageURL != "" {
		schema["image"] = p.ImageURL
	}

	data, err := json.Marshal(schema)
	if err != nil {
		return ""
	}
	return string(data)
}

// BuildBreadcrumbJSON returns a JSON-LD string for schema.org/BreadcrumbList.
// Returns empty string if items is empty.
func BuildBreadcrumbJSON(baseURL string, items []BreadcrumbItem) string {
	if len(items) == 0 {
		return ""
	}

	elements := make([]map[string]any, 0, len(items))
	for i, item := range items {
		entry := map[string]any{
			"@type":    "ListItem",
			"position": i + 1,
			"name":     item.Label,
		}
		if item.URL != "#" {
			entry["item"] = baseURL + item.URL
		}
		elements = append(elements, entry)
	}

	schema := map[string]any{
		"@context":        "https://schema.org",
		"@type":           "BreadcrumbList",
		"itemListElement": elements,
	}

	data, err := json.Marshal(schema)
	if err != nil {
		return ""
	}
	return string(data)
}
