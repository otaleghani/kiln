package obsidian

import (
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"time"

	"bytes"
	"regexp"

	"github.com/djherbis/times"
	"gopkg.in/yaml.v3" // Requires: go get gopkg.in/yaml.v3
)

// Regex patterns for Obsidian syntax
var (
	// linkRegex     = regexp.MustCompile(`\[\[([^\]]+)\]\]`)
	wikilinkRegex = regexp.MustCompile(`(!?)\[\[([^\]]+)\]\]`)
	// tagRegex  = regexp.MustCompile(`#(\w+)`)
	// 1. (?m) enables multi-line mode so ^ matches start of line
	// 2. (?:^|\s) is a non-capturing group matching Start-of-Line OR Whitespace
	// 3. (#[a-zA-Z0-9_\-]+) is the capturing group for the actual tag
	tagRegex = regexp.MustCompile(`(?m)(?:^|\s)(#[a-zA-Z0-9_\-]+)`)
)

func (o *Obsidian) NewFolder(path string) (*Folder, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	// Calculate basic paths
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	relPath, err := filepath.Rel(o.InputDir, path)
	if err != nil {
		o.log.Warn("Couldn't get relative path", "path", path)
		relPath = filepath.Base(absPath) // Fallback
	}

	slugPath := o.GetSlugPath(relPath)
	webPath, err := o.GetPageWebPath(slugPath, "")
	if err != nil {
		return nil, err
	}

	outPath, err := o.GetPageOutputPath(slugPath, "")
	if err != nil {
		return nil, err
	}

	parts := strings.Split(relPath, "/")

	f := &Folder{
		Name:     parts[len(parts)-1],
		Path:     path,
		Files:    []*File{},
		WebPath:  webPath,
		OutPath:  outPath,
		RelPath:  relPath,
		Created:  info.ModTime(), // Fallback: OS creation time is not standard in Go
		Modified: info.ModTime(),
	}

	return f, nil
}

func (o *Obsidian) NewFile(path string) (*File, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	// Calculate basic paths
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	relPath, err := filepath.Rel(o.InputDir, path)
	if err != nil {
		o.log.Warn("Couldn't get relative path", "path", path)
		relPath = filepath.Base(absPath) // Fallback
	}

	ext := filepath.Ext(relPath)
	fullName := filepath.Base(relPath)

	// Handle double extensions properly if needed, otherwise standard logic:
	name := fullName
	if ext == ".md" {
		name = strings.TrimSuffix(fullName, ext)
	}
	folder := filepath.Dir(relPath)

	slugPath := o.GetSlugPath(relPath)
	webPath, err := o.GetPageWebPath(slugPath, ext)
	if err != nil {
		return nil, err
	}

	outPath, err := o.GetPageOutputPath(slugPath, ext)
	if err != nil {
		return nil, err
	}

	t, err := times.Stat(path)
	if err != nil {
		return nil, err
	}

	birthTime := t.ModTime()
	modTime := t.ModTime()
	if t.HasBirthTime() {
		birthTime = t.BirthTime()
	}

	f := &File{
		Path:     absPath,
		RelPath:  relPath,
		Ext:      ext,
		Name:     name,
		OutPath:  outPath,
		WebPath:  webPath,
		Created:  birthTime, // Fallback: OS creation time is not standard in Go
		Modified: modTime,
		Folder:   folder,
		FullName: fullName,
		Size:     info.Size(),
		// Initialize slices to ensure they are empty JSON arrays [] instead of null
		Links:     []string{},
		Backlinks: []string{},
		Tags:      make(map[string]struct{}),
		Embeds:    []string{},
	}

	// 3. Conditional Processing for Markdown
	if ext == ".md" {
		if err := f.processMarkdown(); err != nil {
			return nil, fmt.Errorf("failed to process markdown for %s: %w", fullName, err)
		}
	}

	return f, nil
}

func (f *File) processMarkdown() error {
	rawContent, err := os.ReadFile(f.Path)
	if err != nil {
		return err
	}

	// --- A. Parse Frontmatter ---
	var bodyContent []byte
	if bytes.HasPrefix(rawContent, []byte("---\n")) ||
		bytes.HasPrefix(rawContent, []byte("---\r\n")) {
		parts := bytes.SplitN(rawContent, []byte("---"), 3)
		if len(parts) >= 3 {
			if err := yaml.Unmarshal(parts[1], &f.Frontmatter); err != nil {
				fmt.Printf("Warning: Invalid frontmatter in %s\n", f.Name)
			}
			bodyContent = bytes.TrimSpace(parts[2])
		} else {
			bodyContent = rawContent
		}
	} else {
		bodyContent = rawContent
	}
	f.Content = bodyContent

	bodyString := string(bodyContent)

	// --- B. Extract Links and Embeds ---
	// We use the unified regex to capture both and sort them.
	wikiMatches := wikilinkRegex.FindAllStringSubmatch(bodyString, -1)

	for _, m := range wikiMatches {
		// m[0] = Full match (e.g. "![[Image.png]]" or "[[Link]]")
		// m[1] = The bang "!" (if it's an embed)
		// m[2] = The inner content (e.g. "Image.png")

		if m[1] == "!" {
			// It is an Embed -> Add to f.Embeds
			// Here we strip of the "!" because if not it will be rendered as an embed
			f.Embeds = append(f.Embeds, strings.TrimPrefix(m[0], "!"))

			// In bases embeds are considered both Embeds and Links
			f.Links = append(f.Links, strings.TrimPrefix(m[0], "!"))
		} else {
			// It is a Link -> Add to f.Links
			f.Links = append(f.Links, m[0])
		}
	}

	// --- C. Extract Tags ---
	// Assuming tagRegex is: (?m)(^|\s|>)(#[a-zA-Z0-9_\-]+)
	tagMatches := tagRegex.FindAllStringSubmatch(bodyString, -1)

	// Add YAML tags
	if fmTags, ok := f.Frontmatter["tags"]; ok {
		if tList, ok := fmTags.([]any); ok {
			for _, t := range tList {
				f.Tags[fmt.Sprint(t)] = struct{}{}
			}
		}
	}

	// Add inline tags
	for _, m := range tagMatches {
		if len(m) > 1 {
			f.Tags[m[1]] = struct{}{}
		}
	}

	return nil
}

// GenerateBacklinks populates the Backlinks field for every file in the slice.
func GenerateBacklinks(files []*File) {
	// 1. Create a Lookup Map (Name -> *File)
	fileMap := make(map[string]*File)
	for _, file := range files {
		key := strings.ToLower(file.Name)
		if file.Ext != ".md" {
			fileMap[key+file.Ext] = file
		} else {
			fileMap[key] = file
		}
	}

	// 2. Process Outgoing Links AND Embeds
	for _, sourceFile := range files {

		// Create a temporary list of lists to iterate over both Links and Embeds
		// using the same logic.
		referenceGroups := [][]string{sourceFile.Links, sourceFile.Embeds}

		for _, group := range referenceGroups {
			for _, rawLink := range group {
				// --- A. CLEAN THE LINK ---
				// rawLink is "[[Link]]" or "![[Embed]]"

				// 1. Remove optional Embed "!" prefix
				cleanLink := strings.TrimPrefix(rawLink, "!")

				// 2. Remove surrounding brackets [[ and ]]
				cleanLink = strings.TrimPrefix(cleanLink, "[[")
				cleanLink = strings.TrimSuffix(cleanLink, "]]")

				// 3. Remove Alias (if exists) -> "My Page#Heading"
				if idx := strings.Index(cleanLink, "|"); idx != -1 {
					cleanLink = cleanLink[:idx]
				}

				// 4. Remove Anchor/Header (if exists) -> "My Page"
				if idx := strings.Index(cleanLink, "#"); idx != -1 {
					cleanLink = cleanLink[:idx]
				}

				// 5. Normalize for Lookup
				// Handle paths like "Folder/Note" -> just "Note"
				if strings.Contains(cleanLink, "/") {
					cleanLink = filepath.Base(cleanLink)
				}

				lookupKey := strings.ToLower(cleanLink)

				// --- B. FIND TARGET AND ADD BACKLINK ---
				if targetFile, found := fileMap[lookupKey]; found {
					// Prevent self-linking
					if targetFile.Path == sourceFile.Path {
						continue
					}

					// Format the backlink: [[Name]]
					newBacklink := fmt.Sprintf("[[%s]]", sourceFile.Name)

					// Check for duplicates before adding
					if !slices.Contains(targetFile.Backlinks, newBacklink) {
						targetFile.Backlinks = append(targetFile.Backlinks, newBacklink)
					}
				}
			}
		}
	}
}

func WithLogger(l *slog.Logger) Option {
	return func(o *Obsidian) {
		o.log = l
	}
}

func WithInputDir(s string) Option {
	return func(o *Obsidian) {
		o.InputDir = s
	}
}

func WithOutputDir(s string) Option {
	return func(o *Obsidian) {
		o.OutputDir = s
	}
}

func WithBaseURL(s string) Option {
	return func(o *Obsidian) {
		o.BaseURL = s
	}
}

func WithFlatURLs(b bool) Option {
	return func(o *Obsidian) {
		o.FlatURLs = b
	}
}

func New(opts ...Option) *Obsidian {
	// Default to the standard no-op or default logger
	o := &Obsidian{
		log: slog.Default(),
	}

	for _, opt := range opts {
		opt(o)
	}
	return o
}

// Scans the obsidian vault
func (o *Obsidian) Scan() error {
	o.Vault = &Vault{
		FileIndex:  make(map[string][]*File),
		Tags:       make(map[string]*Tag),
		Folders:    make(map[string]*Folder),
		SourceMap:  make(map[string]string),
		GraphNodes: []GraphNode{},
		Files:      []*File{},
		Sitemap: &Sitemap{
			Path: filepath.Join(o.BaseURL, "/sitemap.xml"),
		},
	}

	filepath.WalkDir(o.InputDir, func(path string, info fs.DirEntry, err error) error {
		l := o.log.With("file", path)

		// Handle permission errors and other related problems
		if err != nil {
			return nil
		}

		// Create relative path
		relPath, err := filepath.Rel(o.InputDir, path)
		if err != nil {
			return err
		}

		// Skip root file
		if path == o.InputDir {
			l.Debug("Skipping file", "reason", "Root folder")
			return nil
		}

		// Skip dotfiles
		if strings.HasPrefix(relPath, ".") && path != o.InputDir {
			l.Debug("Skipping file", "reason", "File is a hidden file or directory")
			return nil
		}

		// Skip directories
		if info.IsDir() {
			folder, err := o.NewFolder(path)
			if err != nil {
				l.Error("Couldn't create new folder", "error", err)
				return nil
			}
			o.log.Debug("Processed folder", "folder", folder.LogValue())
			o.Vault.Folders[folder.RelPath] = folder

			o.Vault.GraphNodes = append(o.Vault.GraphNodes, GraphNode{
				ID:    folder.WebPath,
				Label: folder.Name,
				URL:   folder.WebPath,
				Val:   1,
				Type:  "folder",
			})

			return nil
		}

		file, err := o.NewFile(path)
		if err != nil {
			l.Error("Couldn't create new file", "error", err)
			return nil
		}

		// l.Debug("Processed file", "file", file)

		// Register the file in the global index (filename -> public URL)
		// This is used later for resolving [[WikiLinks]]
		o.Vault.Files = append(o.Vault.Files, file)
		if _, exists := o.Vault.FileIndex[file.Name]; !exists {
			o.Vault.FileIndex[file.Name] = []*File{}
		}
		o.Vault.FileIndex[file.Name] = append(o.Vault.FileIndex[file.Name], file)

		// Used to resolve the real path of the original file (public URL -> original vault file path)
		// This is used later for resolving text embeds ![[Note#heading]]
		o.Vault.SourceMap[file.WebPath] = relPath

		o.Vault.GraphNodes = append(o.Vault.GraphNodes, GraphNode{
			ID:    file.WebPath,
			Label: file.Name,
			URL:   file.WebPath,
			Val:   1,
			Type:  file.Ext,
		})

		o.log.Debug("Processed file", "file", file.LogValue())

		GenerateBacklinks(o.Vault.Files)

		return nil
	})

	// Adds files to folders
	for _, file := range o.Vault.Files {
		switch file.Ext {
		case ".md", ".canvas", ".base":
			if folder, exists := o.Vault.Folders[file.Folder]; exists {
				o.log.Debug("Added file to folder", "file", file.RelPath, "folder", folder.RelPath)
				folder.Files = append(folder.Files, file)
			}
		}
	}

	// Adds subfolders to folders
	for path, currentFolder := range o.Vault.Folders {
		// e.g. if path is "content/blog", parentPath is "content"
		parentPath := filepath.Dir(path)
		// Look for the parent
		if parentFolder, ok := o.Vault.Folders[parentPath]; ok {
			// Safety check: Avoid adding root to itself (filepath.Dir("/") returns "/")
			if parentPath == path {
				continue
			}
			parentFolder.Folders = append(parentFolder.Folders, currentFolder)
		}
	}

	// Sort folders (maps are random order)
	for _, folder := range o.Vault.Folders {
		sort.Slice(folder.Folders, func(i, j int) bool {
			return folder.Folders[i].Name < folder.Folders[j].Name
		})
	}

	// Generates tag map
	for _, file := range o.Vault.Files {
		if len(file.Tags) != 0 {
			o.log.Debug("Adding tags from file", "path", file.RelPath)
			for tagString := range file.Tags {
				if _, exists := o.Vault.Tags[tagString]; !exists {
					webPath, err := o.getTagWebPath(tagString)
					if err != nil {
						o.log.Warn(
							"Couldn't create web path for tag",
							"name",
							tagString,
							"error",
							err,
						)
						continue
					}

					outPath, err := o.getTagOutputPath(tagString)
					if err != nil {
						o.log.Warn(
							"Couldn't create output path for tag",
							"name",
							tagString,
							"error",
							err,
						)
						continue
					}

					newTag := &Tag{
						Name:    tagString,
						WebPath: webPath,
						OutPath: outPath,
						Files:   []*File{},
					}

					o.Vault.Tags[tagString] = newTag
				}

				o.log.Debug("Added tag", "tag", tagString, "file", file.RelPath)
				o.Vault.Tags[tagString].Files = append(o.Vault.Tags[tagString].Files, file)
			}
		}
	}

	// Add tags to graph nodes
	for _, tag := range o.Vault.Tags {
		o.Vault.GraphNodes = append(o.Vault.GraphNodes, GraphNode{
			ID:    tag.WebPath,
			Label: tag.Name,
			URL:   tag.WebPath,
			Val:   1,
			Type:  "tag",
		})
	}

	return nil
}

// loadFavicon loads the favicon.ico file if it exists
func (o *Obsidian) LoadFavicon() error {
	faviconSrc := filepath.Join(o.InputDir, "favicon.ico")
	if _, err := os.Stat(faviconSrc); err != nil {
		return err
	}
	err := CopyFile(faviconSrc, filepath.Join(o.OutputDir, "favicon.ico"))
	if err != nil {
		return err
	}
	o.log.Debug("'favicon.ico' file loaded correctly")
	return nil
}

// loadRedirects loads the _redirects file if it exists. Used for cloudflare pages
// deployment for handling redirects.
//
// For more information check out this link:
// https://developers.cloudflare.com/pages/configuration/redirects/
func (o *Obsidian) LoadRedirects() error {
	redirectsSrc := filepath.Join(o.InputDir, "_redirects")
	if _, err := os.Stat(redirectsSrc); err != nil {
		return err
	}
	err := os.RemoveAll(filepath.Join(o.OutputDir, "_redirects"))
	if err != nil {
		return err
	}

	err = CopyFile(redirectsSrc, filepath.Join(o.OutputDir, "_redirects"))
	if err != nil {
		return err
	}
	o.log.Debug("'_redirects' file loaded correctly")
	return nil
}

// loadCname loads the CNAME file if it exists
func (o *Obsidian) LoadCname() error {
	faviconSrc := filepath.Join(o.InputDir, "CNAME")
	if _, err := os.Stat(faviconSrc); err != nil {
		return err
	}
	err := CopyFile(faviconSrc, filepath.Join(o.OutputDir, "CNAME"))
	if err != nil {
		return err
	}
	o.log.Debug("'CNAME' file loaded correctly")
	return nil
}

// CopyFile is a simple wrapper to copy files from src to dst.
func CopyFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}

// GetFolderLinks returns all of the links between folders
//
// This is done because folder notes are linked with every of it's children.
func (o *Obsidian) GetFolderLinks() []GraphLink {
	links := []GraphLink{}

	for _, folder := range o.Vault.Folders {
		for _, file := range folder.Files {
			links = append(links, GraphLink{Source: folder.WebPath, Target: file.WebPath})
		}

		for _, subFolder := range folder.Folders {
			links = append(links, GraphLink{Source: folder.WebPath, Target: subFolder.WebPath})
		}
	}

	// Prune empty links

	o.log.Debug("Folder links", "amount", len(links))
	return links
}

// GetTagLinks returns all the links between tags
func (o *Obsidian) GetTagLinks() []GraphLink {
	links := []GraphLink{}

	for _, tag := range o.Vault.Tags {
		for _, file := range tag.Files {
			links = append(links, GraphLink{Source: tag.WebPath, Target: file.WebPath})
		}
	}

	return links
}

// File represents a file that needs to be processed
type File struct {
	Path        string              // Complete path of the file
	RelPath     string              // Relative path from input directory
	Ext         string              // Extension of the file
	Name        string              // Name of the file (no extension)
	OutPath     string              // Final output path of the file (e.g. /public/folder/page.html)
	WebPath     string              // Final web path of the page (e.g. /folder/page)
	Created     time.Time           // When the file was created
	Modified    time.Time           // When the file was last modified
	Folder      string              // The folder of note
	FullName    string              // Filename with extension
	Size        int64               // Size of the file
	Frontmatter map[string]any      // Frontmatter, only for notes
	Content     []byte              // Content, only for notes
	Links       []string            // Outgoing links
	Backlinks   []string            // Backlinks to the file
	Tags        map[string]struct{} // Tags
	Embeds      []string            // Embed files
	Breadcrumbs []Breadcrumb
}

// LogValue is used to log out the file
func (f File) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("name", f.Name),
		slog.String("relPath", f.RelPath),
		slog.String("dir", f.Folder),
		slog.String("outPath", f.OutPath),
		slog.String("webPath", f.WebPath),
	)
}

// Folder represents a folder that needs to be processed
type Folder struct {
	Path     string // Complete path of the folder
	Name     string
	RelPath  string    // Relative path from input directory
	OutPath  string    // Final output path of the folder (e.g. /public/folder/index.html)
	WebPath  string    // Final web path of the folder (e.g. /folder)
	Files    []*File   // List of files
	Folders  []*Folder // List of folders
	Created  time.Time // When the folder was created
	Modified time.Time // When the folder was last modified
}

// LogValue is used to log out the folder
func (f Folder) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("relPath", f.RelPath),
		slog.String("path", f.Path),
		slog.String("outPath", f.OutPath),
		slog.String("webPath", f.WebPath),
	)
}

// Vault represents the vault scan
type Vault struct {
	FileIndex  map[string][]*File // Used to resolve wikilinks. "Link" => []Candidates
	SourceMap  map[string]string  // Used to resolve the real disk path. "path.html" => "./vault/path.md"
	GraphNodes []GraphNode        // Lists of all pages for graph
	Files      []*File            // List of all the files found in the vault
	Sitemap    *Sitemap           // Sitemap entity
	Folders    map[string]*Folder // Map of all the folder -> Name of folder -> Folder
	Tags       map[string]*Tag    //
}

// Tag rappresents a tag instance
type Tag struct {
	Name    string  // #something
	Files   []*File // List of files that have that tag
	WebPath string  // Website path
	OutPath string  // Outpath
}

// GraphNode represents a single node in the interactive graph view.
type GraphNode struct {
	ID    string `json:"id"`
	Label string `json:"label"`
	URL   string `json:"url"`
	Val   int    `json:"val"`
	Type  string `json:"type"`
}

// GraphLink represents a directed edge in the note graph.
//
// Used to generate the part of the json file
type GraphLink struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

// Obsidian represents the configs regarding the generation
type Obsidian struct {
	log       *slog.Logger // Custom log instance
	InputDir  string       // Input directory of the vault
	OutputDir string       // Output directory of the parsed vault
	BaseURL   string       // BaseURL of the parsed vault (e.g. https://something.com/folder)
	FlatURLs  bool         // True if flat urls are active (e.g. /folder/note/index.html)
	Vault     *Vault       // Vault scan
}

// Option allows users to configure the Worker
type Option func(*Obsidian)
