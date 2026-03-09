// @feature:jsonld Tests for JSON-LD structured data generation.
package jsonld

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestBuildArticleJSON_FullParams(t *testing.T) {
	created := time.Date(2024, 6, 15, 10, 30, 0, 0, time.UTC)
	modified := time.Date(2024, 7, 20, 14, 0, 0, 0, time.UTC)

	p := ArticleParams{
		Title:        "My Test Article",
		Description:  "A description of the article",
		Author:       "John Doe",
		BaseURL:      "https://example.com",
		WebPath:      "/blog/my-test-article",
		SiteName:     "Example Site",
		DateCreated:  created,
		DateModified: modified,
		ImageURL:     "https://example.com/og/my-test-article.png",
	}

	got := BuildArticleJSON(p)
	if got == "" {
		t.Fatal("expected non-empty JSON-LD, got empty string")
	}

	var m map[string]any
	if err := json.Unmarshal([]byte(got), &m); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if m["@context"] != "https://schema.org" {
		t.Errorf("@context = %v, want https://schema.org", m["@context"])
	}
	if m["@type"] != "Article" {
		t.Errorf("@type = %v, want Article", m["@type"])
	}
	if m["headline"] != "My Test Article" {
		t.Errorf("headline = %v, want My Test Article", m["headline"])
	}
	if m["description"] != "A description of the article" {
		t.Errorf("description = %v, want A description of the article", m["description"])
	}
	if m["datePublished"] != created.Format(time.RFC3339) {
		t.Errorf("datePublished = %v, want %v", m["datePublished"], created.Format(time.RFC3339))
	}
	if m["dateModified"] != modified.Format(time.RFC3339) {
		t.Errorf("dateModified = %v, want %v", m["dateModified"], modified.Format(time.RFC3339))
	}
	if m["image"] != "https://example.com/og/my-test-article.png" {
		t.Errorf("image = %v, want https://example.com/og/my-test-article.png", m["image"])
	}

	author, ok := m["author"].(map[string]any)
	if !ok {
		t.Fatal("author is not an object")
	}
	if author["@type"] != "Person" {
		t.Errorf("author.@type = %v, want Person", author["@type"])
	}
	if author["name"] != "John Doe" {
		t.Errorf("author.name = %v, want John Doe", author["name"])
	}

	publisher, ok := m["publisher"].(map[string]any)
	if !ok {
		t.Fatal("publisher is not an object")
	}
	if publisher["@type"] != "Organization" {
		t.Errorf("publisher.@type = %v, want Organization", publisher["@type"])
	}
	if publisher["name"] != "Example Site" {
		t.Errorf("publisher.name = %v, want Example Site", publisher["name"])
	}

	mainEntity, ok := m["mainEntityOfPage"].(map[string]any)
	if !ok {
		t.Fatal("mainEntityOfPage is not an object")
	}
	if mainEntity["@type"] != "WebPage" {
		t.Errorf("mainEntityOfPage.@type = %v, want WebPage", mainEntity["@type"])
	}
	if mainEntity["@id"] != "https://example.com/blog/my-test-article" {
		t.Errorf("mainEntityOfPage.@id = %v, want https://example.com/blog/my-test-article", mainEntity["@id"])
	}
}

func TestBuildArticleJSON_EmptyTitle(t *testing.T) {
	p := ArticleParams{
		Title:   "",
		Author:  "John Doe",
		BaseURL: "https://example.com",
		WebPath: "/page",
	}

	got := BuildArticleJSON(p)
	if got != "" {
		t.Errorf("expected empty string for empty title, got %q", got)
	}
}

func TestBuildArticleJSON_EmptyAuthor(t *testing.T) {
	p := ArticleParams{
		Title:        "Some Title",
		Author:       "",
		BaseURL:      "https://example.com",
		WebPath:      "/page",
		SiteName:     "My Site",
		DateCreated:  time.Now(),
		DateModified: time.Now(),
	}

	got := BuildArticleJSON(p)
	if got == "" {
		t.Fatal("expected non-empty JSON-LD")
	}

	if strings.Contains(got, `"author"`) {
		t.Error("expected author field to be omitted when Author is empty")
	}
}

func TestBuildArticleJSON_EmptyDescription(t *testing.T) {
	p := ArticleParams{
		Title:        "Some Title",
		Description:  "",
		Author:       "Jane",
		BaseURL:      "https://example.com",
		WebPath:      "/page",
		SiteName:     "My Site",
		DateCreated:  time.Now(),
		DateModified: time.Now(),
	}

	got := BuildArticleJSON(p)
	if got == "" {
		t.Fatal("expected non-empty JSON-LD")
	}

	if strings.Contains(got, `"description"`) {
		t.Error("expected description field to be omitted when Description is empty")
	}
}

func TestBuildArticleJSON_EmptyImageURL(t *testing.T) {
	p := ArticleParams{
		Title:        "Some Title",
		Author:       "Jane",
		BaseURL:      "https://example.com",
		WebPath:      "/page",
		SiteName:     "My Site",
		ImageURL:     "",
		DateCreated:  time.Now(),
		DateModified: time.Now(),
	}

	got := BuildArticleJSON(p)
	if got == "" {
		t.Fatal("expected non-empty JSON-LD")
	}

	if strings.Contains(got, `"image"`) {
		t.Error("expected image field to be omitted when ImageURL is empty")
	}
}

func TestBuildBreadcrumbJSON_ThreeItems(t *testing.T) {
	items := []BreadcrumbItem{
		{Label: "Home", URL: "/"},
		{Label: "Blog", URL: "/blog"},
		{Label: "My Post", URL: "/blog/my-post"},
	}

	got := BuildBreadcrumbJSON("https://example.com", items)
	if got == "" {
		t.Fatal("expected non-empty JSON-LD, got empty string")
	}

	var m map[string]any
	if err := json.Unmarshal([]byte(got), &m); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if m["@context"] != "https://schema.org" {
		t.Errorf("@context = %v, want https://schema.org", m["@context"])
	}
	if m["@type"] != "BreadcrumbList" {
		t.Errorf("@type = %v, want BreadcrumbList", m["@type"])
	}

	elements, ok := m["itemListElement"].([]any)
	if !ok {
		t.Fatal("itemListElement is not an array")
	}
	if len(elements) != 3 {
		t.Fatalf("itemListElement length = %d, want 3", len(elements))
	}

	for i, el := range elements {
		item, ok := el.(map[string]any)
		if !ok {
			t.Fatalf("itemListElement[%d] is not an object", i)
		}
		if item["@type"] != "ListItem" {
			t.Errorf("itemListElement[%d].@type = %v, want ListItem", i, item["@type"])
		}
		wantPos := float64(i + 1)
		if item["position"] != wantPos {
			t.Errorf("itemListElement[%d].position = %v, want %v", i, item["position"], wantPos)
		}
		if item["name"] != items[i].Label {
			t.Errorf("itemListElement[%d].name = %v, want %v", i, item["name"], items[i].Label)
		}
		wantURL := "https://example.com" + items[i].URL
		if item["item"] != wantURL {
			t.Errorf("itemListElement[%d].item = %v, want %v", i, item["item"], wantURL)
		}
	}
}

func TestBuildBreadcrumbJSON_EmptyItems(t *testing.T) {
	got := BuildBreadcrumbJSON("https://example.com", nil)
	if got != "" {
		t.Errorf("expected empty string for empty items, got %q", got)
	}

	got = BuildBreadcrumbJSON("https://example.com", []BreadcrumbItem{})
	if got != "" {
		t.Errorf("expected empty string for empty slice, got %q", got)
	}
}

func TestBuildBreadcrumbJSON_HashURLOmitsItem(t *testing.T) {
	items := []BreadcrumbItem{
		{Label: "Home", URL: "/"},
		{Label: "Tags", URL: "/tags"},
		{Label: "Current Tag", URL: "#"},
	}

	got := BuildBreadcrumbJSON("https://example.com", items)
	if got == "" {
		t.Fatal("expected non-empty JSON-LD")
	}

	var m map[string]any
	if err := json.Unmarshal([]byte(got), &m); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	elements := m["itemListElement"].([]any)
	lastItem := elements[2].(map[string]any)

	if _, exists := lastItem["item"]; exists {
		t.Error("expected item field to be omitted for URL '#'")
	}
	if lastItem["name"] != "Current Tag" {
		t.Errorf("name = %v, want Current Tag", lastItem["name"])
	}
	if lastItem["position"] != float64(3) {
		t.Errorf("position = %v, want 3", lastItem["position"])
	}
}
