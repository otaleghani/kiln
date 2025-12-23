package builder

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

// CustomPage represents a single markdown file or index in the custom generation mode
type CustomPage struct {
	ID           string         // Unique ID (relative path)
	Title        string         // Derived from filename or frontmatter
	Path         string         // Original file path
	RelPermalink string         // Output URL (e.g., /posts/my-post.html)
	Content      template.HTML  // Rendered HTML content
	Frontmatter  map[string]any // Raw YAML from the file
	Params       map[string]any // Merged Config + Frontmatter
	IsIndex      bool           // Is this an index.md?
	Children     []*CustomPage  // If this is an index, list of child pages
}

// CustomSite holds the global state for custom generation
type CustomSite struct {
	Pages      map[string]*CustomPage
	NameLookup map[string]*CustomPage    // Map of "Clean Filename" -> Page for Wikilink resolution
	Configs    map[string]map[string]any // Map of directory path -> config data
}

// Global Regex for WikiLinks [[Link]] or [[Link|Label]]
var wikiLinkRegex = regexp.MustCompile(`\[\[(.*?)(?:\|(.*?))?\]\]`)

// BuildCustom executes the user-first generation logic (Obsidian-SSG)
// It takes sourceDir (vault root) and outputDir as arguments.
func buildCustom() error {
	// Initialize Goldmark using the existing parser in markdown.go
	fileIndex, _ := initBuild()
	u, _ := url.Parse(BaseURL)
	basePath := u.Path
	mdParser, resolver := newMarkdownParser(fileIndex, basePath)

	site := &CustomSite{
		Pages:      make(map[string]*CustomPage),
		NameLookup: make(map[string]*CustomPage),
		Configs:    make(map[string]map[string]any),
	}

	log.Println("PHASE 1: Scanning Configs & Assets")

	err := filepath.Walk(InputDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(InputDir, path)
		if err != nil {
			return err
		}

		// Handle Assets (Copy CSS, JS, Images)
		// Skip directories, markdown, html templates, and json configs
		if !info.IsDir() &&
			!strings.HasSuffix(path, ".md") &&
			!strings.HasSuffix(path, ".html") &&
			!strings.HasSuffix(path, ".json") {

			destPath := filepath.Join(OutputDir, relPath)
			if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
				return err
			}
			copyFile(path, destPath)
			return nil
		}

		// Handle Configs
		if filepath.Base(path) == "config.json" {
			dir := filepath.Dir(relPath)
			if dir == "." {
				dir = ""
			}

			data, err := os.ReadFile(path)
			if err != nil {
				log.Printf("Error reading config %s: %v\n", path, err)
				return nil
			}

			var configMap map[string]any
			if err := json.Unmarshal(data, &configMap); err != nil {
				log.Printf("Error parsing config %s: %v\n", path, err)
				return nil
			}
			site.Configs[dir] = configMap
			log.Printf("Loaded config for dir: '%s'\n", dir)
		}

		return nil
	})

	if err != nil {
		log.Fatalln("PHASE 1: Failed to scan configs and assets")
	}

	fmt.Println("--- PHASE 2: Indexing Pages ---")

	err = filepath.Walk(InputDir, func(path string, d fs.FileInfo, err error) error {
		if err != nil || d.IsDir() || filepath.Ext(path) != ".md" {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		fm, rawContent := parseFrontmatter(data)

		relPath, _ := filepath.Rel(InputDir, path)
		cleanName := strings.TrimSuffix(filepath.Base(path), ".md")
		outPath := strings.Replace(relPath, ".md", ".html", 1)

		permLink := "/" + filepath.ToSlash(outPath)
		if strings.HasSuffix(permLink, "/index.html") {
			permLink = strings.TrimSuffix(permLink, "index.html")
		} else if permLink == "/index.html" {
			permLink = "/"
		}

		page := &CustomPage{
			ID:           relPath,
			Path:         path,
			Title:        cleanName,
			RelPermalink: permLink,
			Frontmatter:  fm,
			IsIndex:      cleanName == "index",
		}

		if val, ok := fm["title"]; ok {
			page.Title = fmt.Sprintf("%v", val)
		}

		var buf bytes.Buffer

		ext := filepath.Ext(path)
		nameWithoutExt := strings.TrimSuffix(d.Name(), ext)
		resolver.CurrentSource = nameWithoutExt
		if err := mdParser.Convert(rawContent, &buf); err != nil {
			return err
		}
		finalHTML := buf.String()
		finalHTML = transformCallouts(finalHTML)
		finalHTML = transformMermaid(finalHTML)
		finalHTML = transformHighlights(finalHTML)
		page.Content = template.HTML(finalHTML)

		site.Pages[relPath] = page
		site.NameLookup[cleanName] = page

		return nil
	})

	if err != nil {
		return err
	}

	fmt.Println("--- PHASE 3: Merging & Relation Resolution ---")

	for path, page := range site.Pages {
		// 3a. Merge Configs
		page.Params = make(map[string]any)

		dir := filepath.Dir(path)
		segments := strings.Split(dir, string(os.PathSeparator))

		currentPath := ""
		for i, seg := range segments {
			if i == 0 && seg == "." {
				continue
			}
			if i > 0 {
				currentPath = filepath.Join(currentPath, seg)
			} else {
				currentPath = seg
			}

			if cfg, ok := site.Configs[currentPath]; ok {
				for k, v := range cfg {
					page.Params[k] = v
				}
			}
		}

		for k, v := range page.Frontmatter {
			page.Params[k] = v
		}

		// 3b. Resolve Linked Notes in Frontmatter
		for key, value := range page.Params {

			// Case 1: Single Relation (key starts with "link_")
			if strings.HasPrefix(key, "link_") {
				if strVal, ok := value.(string); ok {
					targetName := extractWikiLink(strVal)
					if targetName != "" {
						if targetPage, exists := site.NameLookup[targetName]; exists {
							page.Params[key] = targetPage
						}
					}
				}
			}

			// Case 2: List Relation (key starts with "links_")
			if strings.HasPrefix(key, "links_") {
				if listVal, ok := value.([]any); ok {
					var resolvedPages []*CustomPage

					for _, item := range listVal {
						if strItem, ok := item.(string); ok {
							targetName := extractWikiLink(strItem)
							if targetName != "" {
								if targetPage, exists := site.NameLookup[targetName]; exists {
									resolvedPages = append(resolvedPages, targetPage)
								}
							}
						}
					}
					if len(resolvedPages) > 0 {
						page.Params[key] = resolvedPages
					}
				}
			}
		}

		// 3c. Link Children
		if page.IsIndex {
			pageDir := filepath.Dir(path)
			for otherPath, otherPage := range site.Pages {
				otherDir := filepath.Dir(otherPath)
				if pageDir == otherDir && !otherPage.IsIndex {
					page.Children = append(page.Children, otherPage)
				}
			}
		}
	}

	fmt.Println("--- PHASE 4: Rendering ---")

	for _, page := range site.Pages {
		dir := filepath.Dir(page.ID)

		tmplName := "layout.html"
		if page.IsIndex {
			tmplName = "index.html"
		}

		tmplPath := findTemplate(InputDir, dir, tmplName)

		tmplContent := "{{ .Page.Content }}"
		if tmplPath != "" {
			b, _ := os.ReadFile(tmplPath)
			tmplContent = string(b)
		} else {
			fmt.Printf("Warning: No template found for %s\n", page.ID)
		}

		funcMap := template.FuncMap{
			"upper": strings.ToUpper,
			"param": func(p *CustomPage, key string) any {
				if v, ok := p.Params[key]; ok {
					return v
				}
				return nil
			},
			"isPageList": func(v any) bool {
				_, ok := v.([]*CustomPage)
				return ok
			},
			"isPage": func(v any) bool {
				_, ok := v.(*CustomPage)
				return ok
			},
		}

		tmpl, err := template.New(tmplName).Funcs(funcMap).Parse(tmplContent)
		if err != nil {
			fmt.Printf("Error parsing template %s: %v\n", tmplPath, err)
			continue
		}

		ext := filepath.Ext(page.Path)
		nameWithoutExt := strings.TrimSuffix(page.ID, ext)
		finalOutPath, _ := getOutputPaths(page.ID, nameWithoutExt, ext)
		// finalOutPath := filepath.Join(OutputDir, filepath.Dir(page.ID), filepath.Base(page.ID))
		// finalOutPath = strings.Replace(finalOutPath, ".md", ".html", 1)

		if err := os.MkdirAll(filepath.Dir(finalOutPath), 0755); err != nil {
			return err
		}

		f, err := os.Create(finalOutPath)
		if err != nil {
			return err
		}

		data := struct {
			Page *CustomPage
			Site *CustomSite
		}{
			Page: page,
			Site: site,
		}

		if err := tmpl.Execute(f, data); err != nil {
			fmt.Printf("Error executing template for %s: %v\n", page.ID, err)
		}
		f.Close()
	}

	fmt.Println("Done!")
	return nil
}

// Helpers

func parseFrontmatter(data []byte) (map[string]any, []byte) {
	fm := make(map[string]any)
	if !bytes.HasPrefix(data, []byte("---\n")) && !bytes.HasPrefix(data, []byte("---\r\n")) {
		return fm, data
	}
	parts := bytes.SplitN(data, []byte("---"), 3)
	if len(parts) == 3 {
		if err := yaml.Unmarshal(parts[1], &fm); err == nil {
			return fm, bytes.TrimSpace(parts[2])
		}
	}
	return fm, data
}

func findTemplate(root, startDir, filename string) string {
	curr := startDir
	for {
		checkPath := filepath.Join(root, curr, filename)
		// Handle root specially if curr is empty or dot
		if curr == "." || curr == "" {
			checkPath = filepath.Join(root, filename)
		}

		if _, err := os.Stat(checkPath); err == nil {
			return checkPath
		}

		if curr == "." || curr == "" {
			break
		}
		curr = filepath.Dir(curr)
	}
	return ""
}

func extractWikiLink(s string) string {
	matches := wikiLinkRegex.FindStringSubmatch(s)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}
