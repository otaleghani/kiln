package builder

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"gopkg.in/yaml.v3"
	"html/template"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"
)

// walk walks the input directory and handles different kinds of files
func (s *CustomSite) walk() error {
	err := filepath.Walk(InputDir, func(path string, info fs.FileInfo, err error) error {
		relPath, err := filepath.Rel(InputDir, path)
		if err != nil {
			return err
		}

		log.SetPrefix(relPath)
		log.Debug("Processing start...")
		if err != nil {
			return err
		}

		fileExt := filepath.Ext(path)
		switch fileExt {
		case ".md":
			// Index the page
		case ".json":
			if path == filepath.Join(InputDir, "env.json") {
				// handleEnvironment
			}
			if filepath.Base(path) == "config.json" {
				// handleCollection
			}
		case ".html":
			return nil
		default:
			// Copy over
		}

		// Handle env.json
		if path == filepath.Join(InputDir, "env.json") {
			log.SetPrefix(path)
			log.Debug("Found environment variables file")
			return nil
		}

		// Handle other configurations
		if filepath.Base(path) == "config.json" {
		}

		return nil
	})
	return err
}

// handleEnvironmentFile saves the data in the env.json file if found
func (s *CustomSite) handleEnvironmentFile() error {
	return nil
}

// handleConfigFile loads the configuration found in the config.json file
func (s *CustomSite) handleConfigFile(path string, relPath string) error {
	// log.Debug("Found configuration file")
	// if filepath.Join(InputDir, "config.json") == path {
	// 	log.Debug("Configuration file is in root folder. Skipping it.")
	// 	return nil
	// }
	//
	// data, err := os.ReadFile(path)
	// if err != nil {
	// 	log.Error("Coudn't read file", "error", err)
	// 	return nil
	// }
	//
	// // var configMap map[string]any
	// // if err := json.Unmarshal(data, &configMap); err != nil {
	// // 	log.Error("Couldn't unmarshal configuration", "error", err)
	// // 	return nil
	// // }
	//
	// if filepath.Join(InputDir, "config.json") == path {
	// 	log.Debug("Found config file in root folder. Skipping it")
	// 	return nil
	// }
	//
	// log.Debug("Handling configuration")
	// dir := filepath.Dir(relPath)
	// if dir == "." {
	// 	dir = ""
	// }
	//
	// // data, err := os.ReadFile(path)
	// // if err != nil {
	// // 	log.Error("Error reading the configuration", "error", err)
	// // 	return nil
	// // }
	//
	// var rawConfigMap map[string]any
	// if err := json.Unmarshal(data, &rawConfigMap); err != nil {
	// 	log.Error("Coudn't unmarshal the configuration", "error", err)
	// 	return nil
	// }
	//
	// collectionName := ""
	//
	// // Parse configuration, check the field and add it to the site.Configs
	// configMap := make(map[string]FieldConfig)
	// for key, value := range rawConfigMap {
	// 	// Skip collection_name
	// 	if key == "collection_name" {
	// 		strVal, ok := value.(string)
	// 		if !ok {
	// 			log.Error("Missing collection_name field")
	// 		}
	// 		collectionName = strVal
	// 		continue
	// 	}
	// 	// Normalize configuration field
	// 	field, err := normalizeConfigField(value, collections)
	// 	if err != nil {
	// 		log.Error("Couldn't normalized field", "name", key, "error", err)
	// 		continue
	// 	}
	// 	configMap[key] = field
	// }

	// config := &Config{Fields: configMap, Name: collectionName}
	// site.Configs[dir] = config
	// site.ConfigsLookup[collectionName] = config
	// log.Info("Configuration validated and loaded")

	// // Check if there is a field called "collection_name"
	// value, exists := configMap["collection_name"]
	// if !exists {
	// 	log.Warn("Found config.json without 'collection_name'")
	// }
	//
	// // Parse the "collection_name" field
	// collectionName, ok := value.(string)
	// if !ok {
	// 	log.Error("Couldn't parse 'collection_name' from config")
	// }
	//
	// // Check if there is no other collection with that name
	// _, exists = collections[collectionName]
	// if exists {
	// 	log.Error("Duplicate collection", "name", collectionName)
	// }
	//
	// // Add the collection name to the collections map
	// collections[collectionName] = struct{}{}
	log.Info("Successffully loaded config file")
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
	mdParser, resolver := newMarkdownParser(fileIndex, sourceMap, basePath,
		func(path string) ([]byte, error) {
			// Point this to your content folder
			return os.ReadFile(filepath.Join(basePath, path))
		})

	site := &CustomSite{
		Pages:         make(map[string]*CustomPage),
		PagesLookup:   make(map[string]*CustomPage),
		Configs:       make(map[string]*Config),
		ConfigsLookup: make(map[string]*Config),
		Env:           make(map[string]string),
		Assets:        make(map[string]*Asset),
	}
	collections := make(map[string]struct{})

	// TODO: Create the Config struct here instead of the next step
	log.Print(titleStyle.Render("PHASE 0: Find all collection files"))
	err = filepath.Walk(InputDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Handle env.json
		if path == filepath.Join(InputDir, "env.json") {
			log.SetPrefix(path)
			log.Debug("Found environment variables file")
			return nil
		}

		// Handle other configurations
		if filepath.Base(path) == "config.json" {
			if filepath.Join(InputDir, "config.json") == path {
				log.Debug("Found config file in root folder. Skipping it")
				return nil
			}

			log.SetPrefix(path)
			log.Debug("Found config file")

			data, err := os.ReadFile(path)
			if err != nil {
				log.Error("Coudn't read config", "error", err)
				return nil
			}

			var configMap map[string]any
			if err := json.Unmarshal(data, &configMap); err != nil {
				log.Error("Couldn't parse config", "error", err)
				return nil
			}

			// Check if there is a field called "collection_name"
			value, exists := configMap["collection_name"]
			if !exists {
				log.Warn("Found config.json without 'collection_name'")
			}

			// Parse the "collection_name" field
			collectionName, ok := value.(string)
			if !ok {
				log.Error("Couldn't parse 'collection_name' from config")
			}

			// Check if there is no other collection with that name
			_, exists = collections[collectionName]
			if exists {
				log.Error("Duplicate collection", "name", collectionName)
			}

			// Add the collection name to the collections map
			collections[collectionName] = struct{}{}
			log.Info("Successffully loaded config file")
		}

		return nil
	})
	log.SetPrefix("")
	if err != nil {
		log.Fatal("Failed to find all collections configuration")
	}

	// TODO: Instead of rescanning the vault, use the Config struct to take the path to the file
	log.Print(titleStyle.Render("PHASE 1: Scanning configs"))
	err = filepath.Walk(InputDir, func(path string, info fs.FileInfo, err error) error {
		relPath, err := filepath.Rel(InputDir, path)
		if err != nil {
			return err
		}

		// Handle env.json
		if path == filepath.Join(InputDir, "env.json") {
			log.SetPrefix(path)
			log.Debug("Found env.json file")

			// Parse the file
			data, err := os.ReadFile(path)
			if err != nil {
				log.Error("Error reading the environment file", "error", err)
				return nil
			}

			// Unmarshal the file
			var rawEnvMap map[string]any
			if err := json.Unmarshal(data, &rawEnvMap); err != nil {
				log.Error("Error unmarshaling the environment file", "error", err)
				return nil
			}

			envMap := make(map[string]string)
			for key, value := range rawEnvMap {
				strVal, ok := value.(string)
				if !ok {
					log.Warn("Couldn't parse field", "name", key)
					continue
				}
				envMap[key] = strVal
				log.Debug("Added field", "name", key, "value", strVal)
			}
			site.Env = envMap

			log.Info("Added env.json file to site configuration")
			return nil
		}

		// Handle Configs
		if filepath.Base(path) == "config.json" {
			log.SetPrefix(path)
			if filepath.Join(InputDir, "config.json") == path {
				log.Debug("Found config file in root folder. Skipping it")
				return nil
			}

			log.Debug("Handling configuration")
			dir := filepath.Dir(relPath)
			if dir == "." {
				dir = ""
			}

			data, err := os.ReadFile(path)
			if err != nil {
				log.Error("Error reading the configuration", "error", err)
				return nil
			}

			var rawConfigMap map[string]any
			if err := json.Unmarshal(data, &rawConfigMap); err != nil {
				log.Error("Error unmarshaling the configuration", "error", err)
				return nil
			}

			collectionName := ""

			// Parse configuration, check the field and add it to the site.Configs
			configMap := make(map[string]FieldConfig)
			for key, value := range rawConfigMap {
				// Skip collection_name
				if key == "collection_name" {
					strVal, ok := value.(string)
					if !ok {
						log.Error("Missing collection_name field")
					}
					collectionName = strVal
					continue
				}
				// Normalize configuration field
				field, err := normalizeConfigField(value, collections)
				if err != nil {
					log.Error("Couldn't normalized field", "name", key, "error", err)
					continue
				}
				configMap[key] = field
			}

			config := &Config{Fields: configMap, Name: collectionName}
			site.Configs[dir] = config
			site.ConfigsLookup[collectionName] = config
			log.Info("Configuration validated and loaded")
		}

		return nil
	})
	log.SetPrefix("")

	if err != nil {
		log.Fatal("PHASE 1 failed")
	}

	// TODO: Maybe unite this with the first scan?
	log.Print(titleStyle.Render("PHASE 2: Copying and indexing static assets"))
	err = filepath.Walk(InputDir, func(path string, info fs.FileInfo, err error) error {
		log.SetPrefix(path)
		log.Debug("Processing file")
		if err != nil {
			log.Error("Encountered error", "error", err)
			return err
		}

		relPath, err := filepath.Rel(InputDir, path)
		if err != nil {
			log.Error("Couldn't create relative path")
			return err
		}

		// Handle Assets (Copy CSS, JS, Images)
		// Skip directories, markdown, html templates, and json configs
		if !info.IsDir() &&
			!strings.HasSuffix(path, ".md") &&
			!strings.HasSuffix(path, ".html") &&
			!strings.HasSuffix(path, ".json") {

			// destPath := filepath.Join(OutputDir, relPath)

			// Index the asset
			cleanName := filepath.Base(path)
			finalOutPath, webPath := getAssetOutputPath(relPath)

			// Index the file
			file := &Asset{
				ID:           cleanName,
				Path:         path,
				RelPermalink: webPath,
				OutputPath:   finalOutPath,
			}

			if err := os.MkdirAll(filepath.Dir(finalOutPath), 0755); err != nil {
				log.Error(
					"Coudn't create the directory",
					"destination",
					finalOutPath,
					"error",
					err,
				)
				return err
			}

			err := copyFile(path, finalOutPath)
			if err != nil {
				log.Error(
					"Coudn't copy file",
					"destination",
					finalOutPath,
					"error",
					err,
				)
			}

			site.Assets[cleanName] = file
			log.Info("Copied static file to destination", "destination", finalOutPath)
			return nil
		}

		log.Debug("File skipped")
		return nil
	})
	log.SetPrefix("")
	if err != nil {
		log.Fatal("PHASE 2: Failed")
	}

	// Assign collection to pages!
	log.Print(titleStyle.Render("PHASE 3: Discover pages"))
	err = filepath.Walk(InputDir, func(path string, d fs.FileInfo, err error) error {
		log.SetPrefix(path)
		if err != nil || d.IsDir() || filepath.Ext(path) != ".md" {
			log.Debug("File skipped")
			return nil
		}
		relPath, err := filepath.Rel(InputDir, path)
		cleanName := strings.TrimSuffix(filepath.Base(path), ".md")
		config := site.Configs[getConfigDirectory(relPath)]
		page := &CustomPage{
			ID:   relPath,
			Path: path,
		}
		if config != nil {
			page.Collection = config.Name
		}
		site.Pages[relPath] = page
		site.PagesLookup[cleanName] = page
		log.Info("Added file")
		return nil
	})
	log.SetPrefix("")
	if err != nil {
		log.Fatal("PHASE 3: Failed")
	}

	log.Print(titleStyle.Render("PHASE 4: Indexing pages"))
	err = filepath.Walk(InputDir, func(path string, d fs.FileInfo, err error) error {
		log.SetPrefix(path)
		if err != nil || d.IsDir() || filepath.Ext(path) != ".md" {
			log.Debug("File skipped")
			return nil
		}

		log.Debug("Indexing page")
		data, err := os.ReadFile(path)
		if err != nil {
			log.Error("Couldn't read file")
			return err
		}

		relPath, err := filepath.Rel(InputDir, path)
		if err != nil {
			log.Debug("Couldn't create relative path")
			return err
		}
		dir := filepath.Dir(relPath)
		if dir == "." {
			dir = ""
		}
		fm, rawContent := parseFrontmatter(data)

		page := &CustomPage{
			ID:          relPath,
			Path:        path,
			Frontmatter: fm,
		}

		validFm := make(map[string]FieldContent)
		// Validate frontmatter
		for key, value := range fm {
			content, err := validateFrontmatterField(
				key,
				value,
				site,
				page,
			)
			if err != nil {
				log.Warn("Coudn't validate field", "name", key, "error", err)
			}
			validFm[key] = content
		}

		//relPath, err := filepath.Rel(InputDir, path)
		cleanName := strings.TrimSuffix(filepath.Base(path), ".md")
		ext := filepath.Ext(path)
		nameWithoutExt := strings.TrimSuffix(relPath, ext)
		_, webPath := getPageOutputPath(relPath, nameWithoutExt, ext)

		page.Title = cleanName
		page.RelPermalink = webPath
		page.Fields = validFm
		page.IsIndex = cleanName == "index"

		// If the field title is present in the frontmatter, use that as the page Title
		if val, ok := fm["title"]; ok {
			page.Title = fmt.Sprintf("%v", val)
		}

		// Generates the HTML of the body of the note
		log.Debug("Generating HTML")
		var buf bytes.Buffer
		// Set the current file context for the link resolver
		resolver.CurrentSource = nameWithoutExt
		if err := mdParser.Convert(rawContent, &buf); err != nil {
			return err
		}
		finalHTML := buf.String()
		finalHTML = transformCallouts(finalHTML)
		finalHTML = transformMermaid(finalHTML)
		finalHTML = transformHighlights(finalHTML)
		page.Content = template.HTML(finalHTML)

		// Adds the page to the site
		site.Pages[relPath] = page
		site.PagesLookup[cleanName] = page

		log.Info("Added")

		return nil
	})
	log.SetPrefix("")
	if err != nil {
		log.Fatal("PHASE 4 failed")
	}

	log.Print(titleStyle.Render("PHASE 5: Merging & relation resolution"))
	// Link siblings pages
	for path, page := range site.Pages {
		log.SetPrefix(page.ID)
		log.Debug("Processing page")
		// if page.IsIndex {
		count := 0
		pageDir := filepath.Dir(path)
		for otherPath, otherPage := range site.Pages {
			otherDir := filepath.Dir(otherPath)
			if pageDir == otherDir && !otherPage.IsIndex {
				count += 1
				page.Siblings = append(page.Siblings, otherPage)
			}
		}
		log.Debug("Siblings processed and added", "count", count)
		// }
	}

	log.SetPrefix("")
	log.Info("PHASE 6: Rendering")
	for _, page := range site.Pages {
		log.SetPrefix(page.ID)
		dir := filepath.Dir(page.ID)

		tmplName := "layout.html"
		if page.IsIndex {
			log.Debug("Using index.html layout")
			tmplName = "index.html"
		}

		// Finds a template, if it doesn't exist defaults to an empty page to just render the content
		tmplPath := findTemplate(InputDir, dir, tmplName)
		tmplContent := "{{ .Page.Content }}"
		if tmplPath != "" {
			b, err := os.ReadFile(tmplPath)
			if err != nil {
				log.Error("Couldn't read template", "file", tmplPath)
			}
			tmplContent = string(b)
			log.Debug("Loaded template", "file", tmplPath)
		} else {
			log.Warn("No template found for this page. Defaulting to empty layout.")
		}

		// Create function map
		funcMap := template.FuncMap{
			"upper": strings.ToUpper,
			"param": func(p *CustomPage, key string) any {
				if v, ok := p.Fields[key]; ok {
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
			"relation": func(v string) *CustomPage {
				return site.PagesLookup[v]
			},
			// TODO: Add new features, like "limit", some data retrieval logic etc.
		}

		// Parses the template
		tmpl, err := template.New(tmplName).Funcs(funcMap).Parse(tmplContent)
		if err != nil {
			log.Error("Error parsing template", "file", tmplPath, "error", err)
			continue
		}

		// Creates the necessary folders and saves the final file
		ext := filepath.Ext(page.Path)
		nameWithoutExt := strings.TrimSuffix(page.ID, ext)
		finalOutPath, _ := getPageOutputPath(page.ID, nameWithoutExt, ext)

		if err := os.MkdirAll(filepath.Dir(finalOutPath), 0755); err != nil {
			log.Error("Error creating directories", "file", tmplPath, "error", err)
			return err
		}

		f, err := os.Create(finalOutPath)
		if err != nil {
			log.Error("Error creating file", "path", finalOutPath, "error", err)
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
			log.Error("Error executing template", "error", err)
		}
		f.Close()
	}

	log.SetPrefix("")
	log.Info("Rendering finished")
	return nil
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

// TODO: Change this, every collection should have it's own template
// findTemplate finds the first template to use
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

// Global Regex for WikiLinks [[Link]] or [[Link|Label]]
// var wikiLinkRegex = regexp.MustCompile(`\[\[(.*?)(?:\|(.*?))?\]\]`)
// func extractWikiLink(s string) string {
// 	matches := wikiLinkRegex.FindStringSubmatch(s)
// 	if len(matches) > 1 {
// 		return matches[1]
// 	}
// 	return ""
// }

var (
	ErrorInvalidField         error = errors.New("Invalid type")
	ErrorNoTypeNameInField    error = errors.New("No type name in field")
	ErrorEnumTypeWithNoValues error = errors.New("Enum type doesn't have values")
	ErrorCustomTypeWithNoData error = errors.New("Custom type has no data")
	ErrorInvalidReference     error = errors.New("Invalid reference")
)

// normalizeConfigField takes in the raw field from the frontmatter and generates a FieldConfig instance
// TODO: Logging of found fields
func normalizeConfigField(
	input any,
	collections map[string]struct{},
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
		case Enum:
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
		case Reference, References:
			// Parse the reference collection name
			collectionName, ok := v["reference"].(string)
			if !ok {
				return FieldConfig{}, ErrorInvalidReference
			}

			// Check if the collections actually exist
			_, ok = collections[collectionName]
			if !ok {
				return FieldConfig{}, ErrorInvalidReference
			}
			config.Reference = collectionName

		case Custom:
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
	case String,
		Date,
		Boolean,
		Integer,
		Float,
		Image,
		Tag,
		Tags,
		Reference,
		References,
		Enum,
		Custom:
		return true
	default:
		return false
	}
}

var titleStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#FAFAFA")).
	Background(lipgloss.Color("#7D56F4")).
	PaddingLeft(0).
	Width(60)

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
func validateFrontmatterField(
	key string,
	value any,
	site *CustomSite,
	currentPage *CustomPage,
) (FieldContent, error) {
	config := site.Configs[getConfigDirectory(currentPage.ID)]
	if config == nil {
		return FieldContent{}, ErrorNoConfig
	}

	// Check if the key exists in the given configuration
	fieldConfig, ok := config.Fields[key]
	if !ok {
		return FieldContent{}, ErrorNoConfigField
	}

	if fieldConfig.Required && value == "" {
		return FieldContent{}, ErrorRequiredField
	}

	content := FieldContent{Raw: value}

	// Based on the type, try to parse it
	switch fieldConfig.Type {
	case String:
		strVal, ok := value.(string)
		if !ok {
			return FieldContent{}, ErrorParsing
		}
		content.String = strVal
		return content, nil

	case Boolean:
		bolVal, ok := value.(bool)
		if !ok {
			return FieldContent{}, ErrorParsing
		}
		content.Boolean = bolVal
		return content, nil

	case Date:
		strVal, ok := value.(string)
		if !ok {
			return FieldContent{}, ErrorParsing
		}
		// TODO: Support custom date formats
		dateVal, err := time.Parse("2000-01-01", strVal)
		if err != nil {
			return FieldContent{}, err
		}
		content.Date = dateVal
		content.DateTime = dateVal
		return content, nil

	case DateTime:
		strVal, ok := value.(string)
		if !ok {
			return FieldContent{}, ErrorParsing
		}
		// TODO: Support custom date-time formats
		dateVal, err := time.Parse("2000-01-01T12:12:00", strVal)
		if err != nil {
			return FieldContent{}, err
		}
		content.Date = dateVal
		content.DateTime = dateVal
		return content, nil

	case Integer:
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

	case Float:
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

	case Image:
		strVal, ok := value.(string)
		if !ok {
			return FieldContent{}, ErrorParsing
		}
		imgVal, ok := site.Assets[strVal]
		if !ok {
			return FieldContent{}, ErrorAssetNotFound
		}
		content.Image = imgVal
		return content, nil

	case Tag:
		strVal, ok := value.(string)
		if !ok {
			return FieldContent{}, ErrorParsing
		}
		// Add it to the site-wide tag map
		site.Tags[strVal] = append(site.Tags[strVal], currentPage)
		content.Tag = strVal
		return content, nil

	case Tags:
		sliceVal, ok := value.([]string)
		if !ok {
			return FieldContent{}, ErrorParsing
		}

		for _, val := range sliceVal {
			// Add it to the site-wide tag map
			site.Tags[val] = append(site.Tags[val], currentPage)
			content.Tags = append(content.Tags, val)
		}
		return content, nil

	case Reference:
		strVal, ok := value.(string)
		if !ok {
			return FieldContent{}, ErrorParsing
		}
		// Does the requested link exists?
		pageVal, ok := site.PagesLookup[strVal]
		if !ok {
			return FieldContent{}, ErrorReferenceNotExistant
		}
		// Are you referencing the correct collection?
		if fieldConfig.Reference != pageVal.Collection {
			return FieldContent{}, ErrorReferenceWrongCollection
		}
		content.Reference = pageVal
		return content, nil

	case References:
		sliceVal, ok := value.([]string)
		if !ok {
			return FieldContent{}, ErrorParsing
		}
		for _, val := range sliceVal {
			pageVal, ok := site.PagesLookup[val]
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

	case Enum:
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

	case Custom:
		return content, nil

	default:
		log.Fatal("Found config with unknown type", "type", fieldConfig.Type)
	}

	return FieldContent{}, nil
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
	ID           string                  // Unique ID (relative path)
	Title        string                  // Derived from filename or frontmatter
	Path         string                  // Original file path
	RelPermalink string                  // Output URL (e.g., /posts/my-post.html)
	Content      template.HTML           // Rendered HTML content
	Frontmatter  map[string]any          // Raw YAML from the file
	Fields       map[string]FieldContent // Validated collection fields
	IsIndex      bool                    // Is this an index.md?
	Siblings     []*CustomPage           // All the other pages in the current collection
	OutputPath   string                  // Final output path
	Collection   string                  // The collection that the page is a part of
}

// CustomSite holds the global state for custom generation
type CustomSite struct {
	Pages            map[string]*CustomPage
	PagesLookup      map[string]*CustomPage   // Map of "Clean Filename" -> Page for Wikilink resolution
	Configs          map[string]*Config       // Map of directory path -> config data
	ConfigsLookup    map[string]*Config       // Map of collection_name -> Config for references
	CollectionsNames map[string]struct{}      // Map of the different collections names found in the configs
	Env              map[string]string        // Map of the environment variables
	Assets           map[string]*Asset        // Map of "Clean Filename" -> Asset for Asset resolution
	Tags             map[string][]*CustomPage // Map of tag -> array of pages
}

type Asset struct {
	ID           string // Unique ID (relative path)
	Alt          string // Alternative text
	Path         string // Original file path
	RelPermalink string // Output URL
	OutputPath   string // Output path of the file
}

// FieldType is the type of the configuration field. Consts are defined for the supported types
type FieldType string

const (
	String     FieldType = "string"
	Date       FieldType = "date"
	DateTime   FieldType = "dateTime"
	Boolean    FieldType = "boolean"
	Integer    FieldType = "integer"
	Float      FieldType = "float"
	Image      FieldType = "image"
	Tag        FieldType = "tag"
	Tags       FieldType = "tags"
	Reference  FieldType = "reference"
	References FieldType = "references"
	Enum       FieldType = "enum"
	Custom     FieldType = "custom"
)

type Config struct {
	Name      string
	Fields    map[string]FieldConfig
	RawFields map[string]any
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
