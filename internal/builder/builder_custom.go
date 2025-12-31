package builder

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"sort"

	"html/template"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/otaleghani/kiln/internal/log"
	"github.com/yuin/goldmark"
	"gopkg.in/yaml.v3"
)

// walk the input directory and loads the different kinds of files into s.Files
func (s *CustomSite) walk() error {
	s.Files = Files{
		Env:       FilePaths{},
		Markdown:  []FilePaths{},
		Base:      []FilePaths{},
		Canvas:    []FilePaths{},
		Config:    []FilePaths{},
		Layout:    make(map[string]FilePaths),
		Component: []FilePaths{},
		Static:    []FilePaths{},
	}

	err := filepath.Walk(InputDir, func(path string, info fs.FileInfo, err error) error {
		l := log.Default.WithFile(path)
		l.Debug("Processing file...")

		// Handle permission errors etc.
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(InputDir, path)
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			l.Debug("Skipping file", log.FieldReason, "File is a directory")
			return nil
		}

		// Skip dotfiles
		if strings.HasPrefix(relPath, ".") && path != InputDir {
			l.Debug("Skipping file", log.FieldReason, "File is a hidden file (dotfile)")
			return nil
		}

		processFile := FilePaths{path, relPath}
		fileExt := filepath.Ext(path)
		switch fileExt {
		case ".md":
			s.Files.Markdown = append(s.Files.Markdown, processFile)
		case ".base":
			s.Files.Base = append(s.Files.Base, processFile)
		case ".canvas":
			s.Files.Canvas = append(s.Files.Canvas, processFile)
		case ".json":
			if filepath.Base(path) == "config.json" {
				s.Files.Config = append(s.Files.Config, processFile)
				break
			}
			if path == filepath.Join(InputDir, "env.json") {
				s.Files.Env = processFile
				break
			}
			l.Debug("Found unknown JSON file, added to static files")
			s.Files.Static = append(s.Files.Static, processFile)
		case ".html":
			if filepath.Base(path) == "layout.html" {
				s.Files.Layout[getConfigDirectory(relPath)] = processFile
				break
			}
			if strings.HasPrefix(filepath.Base(path), "_") {
				s.Files.Component = append(s.Files.Component, processFile)
				break
			}
			l.Debug("Found unknown HTML file, skipped")
			// staticPages = append(staticPages, processFile)
		default:
			s.Files.Static = append(s.Files.Static, processFile)
		}

		return nil
	})

	return err
}

// TODO: Handle base files
// TODO: Handle canvas files

// parseComponentFiles takes every found HTML component and loads it, creating a base template
func (s *CustomSite) parseComponentFiles() (err error) {
	log.Info("Loading components...")
	for _, file := range s.Files.Component {
		log.Default.WithFile(file.RelPath).Debug("Processing component...")

		s.Template, err = s.Template.ParseFiles(file.Path)
		if err != nil {
			return err
		}
	}
	return nil
}

// parseLayouts creates for each configuration the specific template to execute
func (s *CustomSite) parseLayouts() error {
	log.Info("Parsing layouts...")
	for _, config := range s.Configs {
		l := log.Default.WithFile(config.RelPath)
		l.Debug("Loading 'layout.html'")
		if layout, exists := s.Files.Layout[getConfigDirectory(config.RelPath)]; exists {

			// Clones the base layout (that contains all components)
			configTemplate, err := s.Template.Clone()
			if err != nil {
				return err
			}

			// Parses layout.html for the specific config
			configTemplate, err = configTemplate.ParseFiles(layout.Path)
			if err != nil {
				return err
			}

			config.LayoutPath = layout.Path
			config.LayoutRelPath = layout.RelPath
			config.Template = configTemplate
		} else {
			return fmt.Errorf("No 'layout.html' file found for the collection %s", config.Name)
		}
	}
	return nil
}

// parseNotes loops through the pages of the CustomSite and parses the raw data discovered
// by loadNoteFile.
func (s *CustomSite) parseNotes() error {
	log.Info("Parsing notes...")
	siblings := make(map[string][]*CustomPage)

	for _, page := range s.Pages {
		l := log.Default.WithFile(page.RelPath)

		// Validate fields
		fields := make(map[string]*FieldContent)
		for key, value := range page.RawFrontmatter {
			// TODO: Fix required
			content, err := s.validateFrontmatterField(key, value, page)
			if err != nil {
				l.Error("Coudn't validate field", log.FieldName, key, log.FieldError, err)
				return err
			}
			fields[key] = &content
		}

		// Validate required fields
		if config, ok := s.Configs[getConfigDirectory(page.ID)]; ok && config != nil {
			for fieldName, field := range config.Fields {
				_, exists := fields[fieldName]
				if !exists && field.Required {
					l.Error("Field is required", log.FieldName, fieldName)
					return ErrorRequiredField
				}
			}
		}

		ext := filepath.Ext(page.Path)
		nameWithoutExt := strings.TrimSuffix(page.RelPath, ext)
		outputPath, webPath := getPageOutputPath(page.RelPath, nameWithoutExt, ext)

		page.OutputPath = outputPath
		page.WebPath = webPath
		page.Fields = fields
		page.IsIndex = strings.TrimSuffix(filepath.Base(page.Path), ".md") == "index"

		// Generates the HTML of the body of the note
		var buf bytes.Buffer
		s.Resolver.CurrentSource = page.WebPath
		if err := s.MarkdownParser.Convert(page.RawContent, &buf); err != nil {
			l.Error("Coudn't parse markdown content", log.FieldError, err)
			return err
		}

		customLayoutPath := strings.TrimSuffix(page.Path, ".md") + ".html"
		if _, err := os.Stat(customLayoutPath); err == nil {
			l.Debug("Found custom layout", log.FieldFile, customLayoutPath)

			customTemplate, err := s.Template.Clone()
			if err != nil {
				l.Error("Error cloning base layout", log.FieldError, err)
				return err
			}

			// Parses layout.html for the specific config
			customTemplate, err = customTemplate.ParseFiles(customLayoutPath)
			if err != nil {
				l.Error(
					"Error parsing custom layout",
					log.FieldFile,
					customLayoutPath,
					log.FieldError,
					err,
				)
				return err
			}
			page.Template = customTemplate
		}

		finalHTML := buf.String()
		finalHTML = transformCallouts(finalHTML)
		finalHTML = transformMermaid(finalHTML)
		finalHTML = transformHighlights(finalHTML)
		page.Content = template.HTML(finalHTML)

		// We append the page for the siblings handling
		if page.Collection != "" && !page.IsIndex {
			siblings[page.Collection] = append(siblings[page.Collection], page)
		}
	}

	// Handling siblings
	for _, page := range s.Pages {
		if siblings, collectionExists := siblings[page.Collection]; collectionExists {
			page.Siblings = siblings
		}
	}

	return nil
}

// parseConfigs loops through the configs of the CustomSite and parses the raw data
// discovered by loadConfigFile.
func (s *CustomSite) parseConfigs() error {
	log.Info("Parsing configurations...")

	for _, config := range s.Configs {
		l := log.Default.WithFile(config.RelPath)
		if filepath.Join(InputDir, "config.json") == config.Path {
			l.Debug("Found config file in root folder, skipping it")
			return nil
		}

		// Parse the configuration and check the fields
		configMap := make(map[string]FieldConfig)
		for key, value := range config.RawFields {
			// Normalize configuration field
			if key == "collection_name" {
				l.Debug("Skipped", log.FieldName, key)
				continue
			}
			field, err := s.normalizeConfigField(value)
			if err != nil {
				l.Error("Couldn't normalized field", log.FieldName, key, log.FieldError, err)
				continue
			}
			configMap[key] = field
			l.Debug("Added field", log.FieldName, key, log.FieldType, field.Type)
		}

		config.Fields = configMap

		l.Info("Configuration validated and loaded")
	}
	return nil
}

// parseStaticFiles creates a rappresentation of the given static file and copies the file over
// the output directory
func (s *CustomSite) parseStaticFiles() error {
	log.Info("Parsing static files...")
	for _, file := range s.Files.Static {
		// Index the asset
		cleanName := filepath.Base(file.Path)
		finalOutPath, webPath := getAssetOutputPath(file.RelPath)

		// Index the file
		asset := &Asset{
			ID:           cleanName,
			Path:         file.Path,
			RelPath:      file.RelPath,
			RelPermalink: webPath,
			OutputPath:   finalOutPath,
		}

		if err := os.MkdirAll(filepath.Dir(finalOutPath), 0755); err != nil {
			return err
		}

		err := copyFile(asset.Path, finalOutPath)
		if err != nil {
			return err
		}

		s.Assets[cleanName] = asset

		log.Default.WithFile(asset.RelPath).Info("Static file parsed correctly")
	}
	return nil
}

// loadNoteFile creates the initial rappresentation of the file (discovery phase)
func (s *CustomSite) loadNoteFiles() error {
	log.Info("Loading notes...")
	for _, file := range s.Files.Markdown {
		cleanName := strings.TrimSuffix(filepath.Base(file.Path), ".md")
		config := s.Configs[getConfigDirectory(file.RelPath)]

		data, err := os.ReadFile(file.Path)
		if err != nil {
			return err
		}
		rawFrontmatter, rawContent := parseFrontmatter(data)

		page := &CustomPage{
			ID:             file.RelPath,
			Title:          cleanName,
			Path:           file.Path,
			RelPath:        file.RelPath,
			RawFrontmatter: rawFrontmatter,
			RawContent:     rawContent,
		}
		if config != nil {
			page.Collection = config.Name
		}
		s.Pages[file.RelPath] = page
		s.PagesLookup[cleanName] = page

		log.Default.WithFile(page.RelPath).Info("Markdown note parsed correctly")
	}
	return nil
}

// parseEnvFile saves the data in the env.json file if found
func (s *CustomSite) parseEnvFile() error {
	log.Info("Loading environment...")
	if s.Files.Env.Path == "" {
		log.Info("No 'env.json' file found")
		return nil
	}

	l := log.Default.WithFile(s.Files.Env.RelPath)
	// Read file
	rawData, err := os.ReadFile(s.Files.Env.Path)
	if err != nil {
		return err
	}

	// Unmarshal the data
	var rawEnvMap map[string]any
	if err := json.Unmarshal(rawData, &rawEnvMap); err != nil {
		return err
	}

	// Validate the data
	envMap := make(map[string]string)
	for key, value := range rawEnvMap {
		strVal, ok := value.(string)
		if !ok {
			l.Warn("Couldn't parse field", "name", key)
			continue
		}
		envMap[key] = strVal
		l.Debug("Added field", "name", key, "value", strVal)
	}
	s.Env = envMap

	l.Info("Environment file parsed correctly")
	return nil
}

// loadConfigFiles loads the configuration found in the config.json file
func (s *CustomSite) loadConfigFiles() error {
	// Configuration files are handled before everything because they are needed (e.g. notes)
	log.Info("Loading configurations...")
	for _, file := range s.Files.Config {
		if filepath.Join(InputDir, "config.json") == file.Path {
			log.Debug("Configuration file is in root folder, skipping it")
			return nil
		}

		rawData, err := os.ReadFile(file.Path)
		if err != nil {
			return err
		}

		var rawConfigMap map[string]any
		if err := json.Unmarshal(rawData, &rawConfigMap); err != nil {
			return err
		}

		collectionName, ok := rawConfigMap["collection_name"].(string)
		if !ok {
			return errors.New("Couln't parse 'collection_name' field")
		}

		// Check if there is another collection with that name
		if _, exists := s.ConfigsLookup[collectionName]; exists {
			return fmt.Errorf("Found two collections with the name %s", collectionName)
		}

		config := &Config{
			ID:        getConfigDirectory(file.RelPath),
			Path:      file.Path,
			RelPath:   file.RelPath,
			Name:      collectionName,
			RawFields: rawConfigMap,
		}
		s.Configs[config.ID] = config
		s.ConfigsLookup[collectionName] = config

		log.Default.WithFile(config.RelPath).Info("Configuration file loaded correctly")
	}
	return nil
}

func (s *CustomSite) render() error {
	for _, page := range s.Pages {
		l := log.Default.WithFile(page.RelPath)

		var tmpl *template.Template
		var tmplPath string
		if page.Template != nil {
			l.Debug("Using custom template")
			tmpl = page.Template
			tmplPath = strings.TrimSuffix(filepath.Base(page.Path), ".md") + ".html"
		} else {
			config := s.ConfigsLookup[page.Collection]
			tmpl = config.Template
			l.Debug("Using collection template")
			tmplPath = filepath.Base(config.LayoutPath)
		}

		if err := os.MkdirAll(filepath.Dir(page.OutputPath), 0755); err != nil {
			l.Error("Error creating dirs", log.FieldPath, page.OutputPath, log.FieldError, err)
			return err
		}

		f, err := os.Create(page.OutputPath)
		if err != nil {
			l.Error("Error creating file", log.FieldPath, page.OutputPath, log.FieldError, err)
			return err
		}
		defer f.Close()

		data := &CustomPageData{
			Page: page,
			Site: s,
		}

		if err := tmpl.ExecuteTemplate(f, tmplPath, data); err != nil {
			l.Error("Error executing template", log.FieldError, err)
			return err
		}
	}
	return nil
}

// BuildCustom executes the user-first generation logic (Obsidian-SSG)
// It takes sourceDir (vault root) and outputDir as arguments.
func buildCustom() error {
	// Initialize Goldmark using the existing parser in markdown.go
	fileIndex, sourceMap, _ := initBuild()
	u, err := url.Parse(BaseURL)
	if err != nil {
		log.Warn("Couldn't parse base URL", "url", BaseURL)
	}
	basePath := u.Path
	markdownParser, resolver := newMarkdownParser(fileIndex, sourceMap, basePath,
		func(path string) ([]byte, error) {
			// Point this to your content folder
			return os.ReadFile(filepath.Join(basePath, path))
		})

	site := &CustomSite{
		Pages:          make(map[string]*CustomPage),
		PagesLookup:    make(map[string]*CustomPage),
		Configs:        make(map[string]*Config),
		ConfigsLookup:  make(map[string]*Config),
		Tags:           make(map[string][]*CustomPage),
		Assets:         make(map[string]*Asset),
		Env:            make(map[string]string),
		MarkdownParser: markdownParser,
		Resolver:       resolver,
		Template:       template.New("base"),
	}
	site.Template.Funcs(site.getFuncMap())

	err = site.walk()
	if err != nil {
		log.Fatal("Error in walk", log.FieldError, err)
	}

	err = site.parseEnvFile()
	if err != nil {
		log.Fatal("Error handling the 'env.json'", log.FieldError, err)
	}

	err = site.loadConfigFiles()
	if err != nil {
		log.Fatal("Error loading a 'config.json'", log.FieldError, err)
	}

	err = site.parseConfigs()
	if err != nil {
		log.Fatal("Error parsing a 'config.json'", log.FieldError, err)
	}

	// Load components files before layouts
	err = site.parseComponentFiles()
	if err != nil {
		log.Fatal("Error parsing components", log.FieldError, err)
	}

	err = site.parseLayouts()
	if err != nil {
		log.Fatal("Error loading layouts", log.FieldError, err)
	}

	err = site.loadNoteFiles()
	if err != nil {
		log.Fatal("Error loading notes", log.FieldError, err)
	}

	err = site.parseStaticFiles()
	if err != nil {
		log.Fatal("Error handling static file", log.FieldError, err)
	}

	err = site.parseNotes()
	if err != nil {
		log.Fatal("Error parsing note", log.FieldError, err)
	}

	err = site.render()
	if err != nil {
		log.Fatal("Error rendering pages", log.FieldError, err)
	}

	return nil
}

// resolveFieldValue extracts the underlying Go value from a FieldContent wrapper
func resolveFieldValue(p *CustomPage, key string) any {
	v, ok := p.Fields[key]
	if !ok {
		// Fallback: Check standard struct fields if not in custom Fields map
		// This allows users to filter by "Title" or "Date" even if they aren't in Fields map
		switch key {
		case "Title":
			return p.Title
		case "Path":
			return p.WebPath
		case "Content":
			return p.Content
		case "Collection":
			return p.Collection
		case "Siblings":
			return p.Siblings
		default:
			return nil
		}
	}

	switch v.Config.Type {
	case TypeImage:
		return v.Image
	case TypeBoolean:
		return v.Boolean
	case TypeInteger:
		return v.Integer
	case TypeFloat:
		return v.Float
	case TypeDate:
		return v.Date
	case TypeDateTime:
		return v.DateTime
	case TypeEnum:
		return v.Enum
	case TypeString:
		return v.String
	case TypeTag:
		return v.Tag
	case TypeTags:
		return v.Tags
	case TypeReference:
		return v.Reference
	case TypeReferences:
		return v.References
	default:
		return v.Config.Data
	}
}

// getFuncMap creates the functions for the template
func (s *CustomSite) getFuncMap() template.FuncMap {
	f := template.FuncMap{}

	f["env"] = tmplFuncEnv
	f["asset"] = tmplFuncAsset
	f["tag"] = tmplFuncTag
	f["page"] = tmplFuncPage
	f["get"] = tmplFuncGet
	f["where"] = tmplFuncWhere
	f["where_not"] = tmplFuncWhereNot
	f["limit"] = tmplFuncLimit
	f["offset"] = tmplFuncOffset
	f["sort"] = tmplFuncSort

	return f
}

func tmplFuncEnv(key string, s *CustomSite) any {
	if value, exists := s.Env[key]; exists {
		return value
	}
	return nil
}

func tmplFuncAsset(key string, s *CustomSite) any {
	if asset, exists := s.Assets[key]; exists {
		return asset
	}
	return nil
}
func tmplFuncTag(key string, s *CustomSite) any {
	if pages, exists := s.Tags[key]; exists {
		return pages
	}
	return nil
}
func tmplFuncPage(key string, s *CustomSite) any {
	if page, exists := s.PagesLookup[key]; exists {
		return page
	}
	return nil
}

// get is the accessor for a CustomPage instance
//
// Usage: {{ .Page | get "title" }}
func tmplFuncGet(key string, p *CustomPage) any {
	if p == nil {
		return nil
	}
	return resolveFieldValue(p, key)
}

// where filters the given []CustomPage
//
// Usage: {{ .Siblings | where "type" "post" }}
func tmplFuncWhere(key string, value any, list []*CustomPage) []*CustomPage {
	var result []*CustomPage
	for _, p := range list {
		if resolveFieldValue(p, key) == value {
			result = append(result, p)
		}
	}
	return result
}

// Usage: {{ range .Siblings | where_not "draft" true }}
func tmplFuncWhereNot(key string, value any, list []*CustomPage) []*CustomPage {
	var result []*CustomPage
	for _, p := range list {
		val := resolveFieldValue(p, key)
		if val != value {
			result = append(result, p)
		}
	}
	return result
}

// limit slices based on the given integer
//
// Usage: {{ .Siblings | limit 3 }}
func tmplFuncLimit(n int, list []*CustomPage) []*CustomPage {
	if n > len(list) {
		return list
	}
	return list[:n]
}

// Usage: {{ range .Siblings | offset 2 }}
func tmplFuncOffset(n int, list []*CustomPage) []*CustomPage {
	if n >= len(list) {
		return []*CustomPage{}
	}
	return list[n:]
}

// isLess returns true if value 'a' is strictly less than value 'b'
//
// This function is used as the engine for tmplFuncSort
func isLess(a, b any) bool {
	if a == nil {
		return b != nil // if a is nil and b is not, a < b is true (a comes first)
	}
	if b == nil {
		return false // if b is nil, a is not less than b
	}

	// Type Switch to handle your specific FieldContent types
	switch va := a.(type) {
	case int:
		if vb, ok := b.(int); ok {
			return va < vb
		}
		// Handle int vs float comparison mixed cases
		if vb, ok := b.(float64); ok {
			return float64(va) < vb
		}

	case float64:
		if vb, ok := b.(float64); ok {
			return va < vb
		}
		if vb, ok := b.(int); ok {
			return va < float64(vb)
		}

	case string:
		if vb, ok := b.(string); ok {
			return va < vb
		}

	case time.Time:
		if vb, ok := b.(time.Time); ok {
			return va.Before(vb)
		}

	case bool:
		if vb, ok := b.(bool); ok {
			// standard: false (0) < true (1)
			return !va && vb
		}
	}

	// Fallback: If types don't match (e.g. string vs int), compare as strings
	return fmt.Sprintf("%v", a) < fmt.Sprintf("%v", b)
}

func tmplFuncSort(key string, direction string, list []*CustomPage) []*CustomPage {
	if len(list) == 0 {
		return list
	}

	sortedList := make([]*CustomPage, len(list))
	copy(sortedList, list)

	// Determine Sort Direction
	desc := strings.ToLower(direction) == "desc"

	sort.SliceStable(sortedList, func(i, j int) bool {
		valI := resolveFieldValue(sortedList[i], key)
		valJ := resolveFieldValue(sortedList[j], key)

		if desc {
			return isLess(valJ, valI)
		}
		return isLess(valI, valJ)
	})

	return sortedList
}

// parseFrontmatter parses the raw markdown data and returns the frontmatter and the content of the note
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

// Global Regex for WikiLinks [[Link]] or [[Link|Label]]
var wikiLinkRegex = regexp.MustCompile(`\[\[(.*?)(?:\|(.*?))?\]\]`)

func extractWikiLink(s string) string {
	matches := wikiLinkRegex.FindStringSubmatch(s)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

var (
	ErrorInvalidField         error = errors.New("Invalid type")
	ErrorNoTypeNameInField    error = errors.New("No type name in field")
	ErrorEnumTypeWithNoValues error = errors.New("Enum type doesn't have values")
	ErrorCustomTypeWithNoData error = errors.New("Custom type has no data")
	ErrorInvalidReference     error = errors.New("Invalid reference")
)

// normalizeConfigField takes in the raw field from the frontmatter and generates a FieldConfig instance
func (s *CustomSite) normalizeConfigField(
	input any,
) (FieldConfig, error) {
	config := FieldConfig{
		Required: false,
	}

	switch v := input.(type) {
	// Shorthand definition
	case string:
		typeName := FieldType(v)
		if !isValidType(typeName) {
			return FieldConfig{}, ErrorInvalidField
		}
		config.Type = typeName
		return config, nil

	// Longhand definition / complex types
	case map[string]any:
		// Extract field type name
		typeVal, ok := v["type"].(string)
		if !ok {
			return FieldConfig{}, ErrorNoTypeNameInField
		}

		typeName := FieldType(typeVal)
		if !isValidType(typeName) {
			return FieldConfig{}, ErrorInvalidField
		}

		config.Type = typeName

		// Extracts the required field
		if req, ok := v["required"].(bool); ok {
			config.Required = req
		}

		// Get optional, type specific fields, and ignore everything else
		switch typeName {
		case TypeEnum:
			rawValues, ok := v["values"].([]any)
			if !ok {
				return FieldConfig{}, ErrorEnumTypeWithNoValues
			}

			for _, val := range rawValues {
				str, ok := val.(string)
				if !ok {
					log.Debug("Couldn't parse enum data for field %s", typeName)
					continue
				}
				config.AllowedValues = append(config.AllowedValues, str)
			}
		case TypeReference, TypeReferences:
			// Parse the reference collection name
			collectionName, ok := v["reference"].(string)
			if !ok {
				return FieldConfig{}, ErrorInvalidReference
			}

			// Check if the collections actually exist
			_, ok = s.ConfigsLookup[collectionName]
			if !ok {
				return FieldConfig{}, ErrorInvalidReference
			}
			config.Reference = collectionName

		case TypeCustom:
			rawData, ok := v["data"]
			if !ok {
				return FieldConfig{}, ErrorCustomTypeWithNoData
			}

			config.Data = rawData
		}
	}

	return config, nil
}

func isValidType(name FieldType) bool {
	switch name {
	case TypeString,
		TypeDate,
		TypeDateTime,
		TypeBoolean,
		TypeInteger,
		TypeFloat,
		TypeImage,
		TypeTag,
		TypeTags,
		TypeReference,
		TypeReferences,
		TypeEnum,
		TypeCustom:
		return true
	default:
		return false
	}
}

var (
	ErrorNoConfig        = errors.New("No configuration found.")
	ErrorNoConfigField   = errors.New("Field does not exist.")
	ErrorRequiredField   = errors.New("Field is required.")
	ErrorParsing         = errors.New("Failed parsing.")
	ErrorWrongTimeLayout = errors.New(
		"Wrong time layout. Use Obsidian default date layout '2000-01-02' or 2000-01-01T12:12:00",
	)
	ErrorAssetNotFound            = errors.New("Asset not found.")
	ErrorReferenceNotExistant     = errors.New("Referenced page does not exist.")
	ErrorUnknownValueEnum         = errors.New("Found a unallowed value in enum.")
	ErrorReferenceWrongCollection = errors.New("Link points to wrong reference collection")
)

// validateFrontmatter takes in a raw frontmatter field and a configuration and returns a validated frontmatter
// TODO: Divide it into different smaller functions
func (s *CustomSite) validateFrontmatterField(
	key string,
	value any,
	currentPage *CustomPage,
) (FieldContent, error) {
	l := log.Default.WithFile(currentPage.RelPath)

	config := s.Configs[getConfigDirectory(currentPage.ID)]
	if config == nil {
		return FieldContent{}, ErrorNoConfig
	}

	// Check if the key exists in the given configuration
	fieldConfig, ok := config.Fields[key]
	if !ok {
		return FieldContent{}, ErrorNoConfigField
	}

	l.Debug("Validating field", log.FieldName, key)

	content := FieldContent{Raw: value, Config: &fieldConfig}

	// Based on the type, try to parse it
	switch fieldConfig.Type {
	case TypeString:
		strVal, ok := value.(string)
		if !ok {
			return FieldContent{}, ErrorParsing
		}
		content.String = strVal
		return content, nil

	case TypeBoolean:
		bolVal, ok := value.(bool)
		if !ok {
			return FieldContent{}, ErrorParsing
		}
		content.Boolean = bolVal
		return content, nil

	case TypeDate, TypeDateTime:
		switch v := value.(type) {
		case string:
			// TODO: Support custom date formats
			dateVal, err := time.Parse("2006-01-02 15:04:05 -0700 MST", v)
			if err != nil {
				return FieldContent{}, err
			}
			content.Date = dateVal
			content.DateTime = dateVal
			return content, nil
		case time.Time:
			content.Date = v
			content.DateTime = v
			return content, nil
		}
		return FieldContent{}, ErrorParsing

	case TypeInteger:
		strVal, ok := value.(string)
		if !ok {
			return FieldContent{}, ErrorParsing
		}
		intVal, err := strconv.Atoi(strVal)
		if err != nil {
			return FieldContent{}, err
		}
		content.Integer = intVal
		return content, nil

	case TypeFloat:
		strVal, ok := value.(string)
		if !ok {
			return FieldContent{}, ErrorParsing
		}
		floatVal, err := strconv.ParseFloat(strVal, 64)
		if err != nil {
			return FieldContent{}, err
		}
		content.Float = floatVal
		return content, nil

	case TypeImage:
		strVal, ok := value.(string)
		if !ok {
			return FieldContent{}, ErrorParsing
		}
		imgVal, ok := s.Assets[extractWikiLink(strVal)]
		if !ok {
			return FieldContent{}, ErrorAssetNotFound
		}
		content.Image = imgVal
		return content, nil

	case TypeTag:
		strVal, ok := value.(string)
		if !ok {
			return FieldContent{}, ErrorParsing
		}
		// Add it to the site-wide tag map
		s.Tags[strVal] = append(s.Tags[strVal], currentPage)
		content.Tag = strVal
		return content, nil

	case TypeTags:
		sliceVal, ok := value.([]string)
		if !ok {
			return FieldContent{}, ErrorParsing
		}

		for _, val := range sliceVal {
			// Add it to the site-wide tag map
			s.Tags[val] = append(s.Tags[val], currentPage)
			content.Tags = append(content.Tags, val)
		}
		return content, nil

	case TypeReference:
		strVal, ok := value.(string)
		if !ok {
			return FieldContent{}, ErrorParsing
		}
		pageVal, ok := s.PagesLookup[extractWikiLink(strVal)]
		if !ok {
			return FieldContent{}, ErrorReferenceNotExistant
		}
		// Are you referencing the correct collection?
		if fieldConfig.Reference != pageVal.Collection {
			return FieldContent{}, ErrorReferenceWrongCollection
		}
		content.Reference = pageVal
		return content, nil

	case TypeReferences:
		sliceVal, err := extractStringSlice(value)
		if err != nil {
			return FieldContent{}, ErrorParsing
		}

		for _, val := range sliceVal {
			pageVal, ok := s.PagesLookup[extractWikiLink(val)]
			if !ok {
				return FieldContent{}, ErrorReferenceNotExistant
			}
			// Are you referencing the correct collection?
			if fieldConfig.Reference != pageVal.Collection {
				return FieldContent{}, ErrorReferenceWrongCollection
			}
			content.References = append(content.References, pageVal)
		}
		return content, nil

	case TypeEnum:
		strVal, ok := value.(string)
		if !ok {
			return FieldContent{}, ErrorParsing
		}
		// Does the given value appear in the values?
		if exists := slices.Contains(fieldConfig.AllowedValues, strVal); !exists {
			return FieldContent{}, ErrorUnknownValueEnum
		}
		content.Enum = strVal
		return content, nil

	case TypeCustom:
		content.Raw = fieldConfig.Data
		return content, nil

	default:
		log.Fatal("Found config with unknown type", "type", fieldConfig.Type)
	}

	return FieldContent{}, nil
}

// extractStringSlice parses the list from a frontmatter
func extractStringSlice(input any) ([]string, error) {
	if input == nil {
		return []string{}, nil
	}

	var result []string

	switch v := input.(type) {
	case string:
		// Single string, we append
		result = append(result, v)
	case []string:
		result = v // I don't know if this can actually happen
	case []any:
		for _, item := range v {
			// RECURSION: This handles the nested [[[Link]]] issue
			// If item is itself a list (due to extra brackets), flatten it
			subSlice, err := extractStringSlice(item)
			if err != nil {
				return nil, err
			}
			result = append(result, subSlice...)
		}
	default:
		return nil, fmt.Errorf("unexpected type %T", input)
	}

	return result, nil
}

// getConfigDirectory get's the configuration directory given a relative path
func getConfigDirectory(relPath string) string {
	dir := filepath.Dir(relPath)
	if dir == "." {
		dir = ""
	}
	return dir
}

// CustomPage represents a single markdown file or index in the custom generation mode
type CustomPage struct {
	ID             string                   // Unique ID (relative path) // TODO: Delete this, use RelPermalink
	Title          string                   // Derived from filename or frontmatter
	Path           string                   // Original file path
	RelPath        string                   // Original file relative path
	WebPath        string                   // Output URL (e.g., /posts/my-post.html)
	Content        template.HTML            // Rendered HTML content
	RawFrontmatter map[string]any           // Raw YAML from the file
	RawContent     []byte                   // Raw content of the note
	Fields         map[string]*FieldContent // Validated collection fields
	IsIndex        bool                     // Is this an index.md?
	Siblings       []*CustomPage            // All the other pages in the current collection
	OutputPath     string                   // Final output path
	Collection     string                   // The collection that the page is a part of
	Template       *template.Template       // If the page has a custom template override it will show here
}

// CustomSite holds the global state for custom generation
type CustomSite struct {
	Pages            map[string]*CustomPage
	PagesLookup      map[string]*CustomPage   // Map of "Clean Filename" -> Page for Wikilink resolution
	Configs          map[string]*Config       // Map of directory path -> config data
	ConfigsLookup    map[string]*Config       // Map of collection_name -> Config for references
	CollectionsNames map[string]struct{}      // TODO: Delete this. We already have the ConfigsLookup for this. Map of the different collections names found in the configs to check
	Env              map[string]string        // Map of the environment variables
	Assets           map[string]*Asset        // Map of "Clean Filename" -> Asset for Asset resolution
	Tags             map[string][]*CustomPage // Map of tag -> array of pages
	MarkdownParser   goldmark.Markdown        // Custom markdown parser
	Resolver         *IndexResolver           // Link resolver
	Template         *template.Template       // Base template with all components loaded
	Files            Files                    // All the paths to files to process
}

type Asset struct {
	ID           string // Unique ID (relative path)
	Path         string // Path of the original path
	RelPath      string // Relative path of the original file
	RelPermalink string // Output URL
	OutputPath   string // Output path of the file
}

// FieldType is the type of the configuration field. Consts are defined for the supported types
type FieldType string

type Config struct {
	ID            string                 // The directory where the config.json file is present
	Name          string                 // The name of the collection
	Path          string                 // The full path of the configuration
	RelPath       string                 // The relative path of the configuration
	Fields        map[string]FieldConfig // The parsed and normalized fields of the configuration
	RawFields     map[string]any         // The raw fields of the configuration
	LayoutPath    string                 // The layout path
	LayoutRelPath string                 // The layout relative path
	Template      *template.Template     // The specific template to execute for every page of the collection
}

// FieldConfig describes the field as displayed in a config.json
type FieldConfig struct {
	Type          FieldType // The type name
	Required      bool      // Defaults to false
	AllowedValues []string  // The allowed values (used in enums)
	Data          any       // The raw JSON of the data field (used in custom)
	Reference     string    // The collection that is references (used in reference and references)
}

// FieldContent holds the content parsed from the frontmatter. You can then access the data
// by using one of the fields, based on the type of your field.
type FieldContent struct {
	Config     *FieldConfig
	Raw        any
	String     string
	Date       time.Time
	DateTime   time.Time
	Boolean    bool
	Integer    int
	Float      float64
	Image      *Asset
	Tag        string
	Tags       []string
	Reference  *CustomPage
	References []*CustomPage
	Enum       string
}

// FilePaths rappresents a file that needs to be processed
type FilePaths struct {
	Path    string
	RelPath string
}

// Files rappresents all the different kinds of files to process
type Files struct {
	Env       FilePaths            // Expected only one 'env.json' file
	Markdown  []FilePaths          // All found '.md' files
	Base      []FilePaths          // All found '.base' files
	Canvas    []FilePaths          // All found '.canvas' files
	Config    []FilePaths          // All found 'config.json' files
	Layout    map[string]FilePaths // Layouts  are related to the collection name
	Component []FilePaths          // All found '_*.html' files
	Static    []FilePaths          // Other files are treated as static
}

// CustomPageData is the struct passed to the templates
type CustomPageData struct {
	Page *CustomPage
	Site *CustomSite
}

const (
	TypeString     FieldType = "string"
	TypeDate       FieldType = "date"
	TypeDateTime   FieldType = "dateTime"
	TypeBoolean    FieldType = "boolean"
	TypeInteger    FieldType = "integer"
	TypeFloat      FieldType = "float"
	TypeImage      FieldType = "image"
	TypeTag        FieldType = "tag"
	TypeTags       FieldType = "tags"
	TypeReference  FieldType = "reference"
	TypeReferences FieldType = "references"
	TypeEnum       FieldType = "enum"
	TypeCustom     FieldType = "custom"
)
