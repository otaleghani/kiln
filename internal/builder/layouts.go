package builder

import (
	"html/template"
	"strings"
	textTemplate "text/template"

	"github.com/otaleghani/kiln/assets"
	"github.com/otaleghani/kiln/internal/log"
)

// resolveLayout looks up a Layout by name.
func resolveLayout(name string) *Layout {
	log.Info("Resolving layout...", "name", name)
	layout, ok := layouts[strings.ToLower(name)]
	if !ok {
		log.Warn("Layout not found, using default layout", log.FieldName, "default")
		name = "default"
		layout = layouts[name]
	}

	return layout
}

// layouts is a key-value pairs of all available layouts
var layouts = map[string]*Layout{
	"default": {
		Name:     "default",
		HtmlPath: "default_layout.html",
		CssPath:  "default_style.css",
		JsPath:   "default_app.js",
	},
}

// loadLayout loads the given layout files into memory
func (l *Layout) loadLayout() error {
	log.Debug("Loading layout", "name", l.Name)

	// Load and parse the base HTML layout
	layoutContent, err := assets.TemplateFS.ReadFile(l.HtmlPath)
	if err != nil {
		return err
	}
	tmplLayout, err := template.New("layout").Parse(string(layoutContent))
	if err != nil {
		return err
	}
	l.HtmlTemplate = tmplLayout

	// Load and parse the CSS template
	cssContent, err := assets.TemplateFS.ReadFile(l.CssPath)
	if err != nil {
		return err
	}
	tmplCSS, err := textTemplate.New("css").Parse(string(cssContent))
	if err != nil {
		return err
	}
	l.CssTemplate = tmplCSS

	// Load and parse the JS template
	jsContent, err := assets.TemplateFS.ReadFile(l.JsPath)
	if err != nil {
		return err
	}
	tmplJS, err := textTemplate.New("js").Parse(string(jsContent))
	if err != nil {
		return err
	}
	l.JsTemplate = tmplJS

	// Load and parse the graph JS template
	jsGraphTemplate, err := assets.TemplateFS.ReadFile("graph.js")
	if err != nil {
		return err
	}
	tmplGraphJS, err := textTemplate.New("js").Parse(string(jsGraphTemplate))
	if err != nil {
		return err
	}
	l.JsGraphTemplate = tmplGraphJS

	// Load and parse the graph JS template
	jsCanvasTemplate, err := assets.TemplateFS.ReadFile("canvas.js")
	if err != nil {
		return err
	}
	tmplCanvasJS, err := textTemplate.New("js").Parse(string(jsCanvasTemplate))
	if err != nil {
		return err
	}
	l.JsCanvasTemplate = tmplCanvasJS

	return nil
}

// Layout contains the paths for the layout
//
// Every Layout is made by 3 different files, layout.html, style.css, and app.js.
// If you are creating a new layout create the needed files by prepending the name of the
// layout to the name of the file. If you have a layout called "default" it should have the
// following files: default_layout.html, default_style.css and default_app.js
type Layout struct {
	Name             string
	HtmlPath         string                 // Path of the HTML file
	CssPath          string                 // Path of the CSS file
	JsPath           string                 // Path of the JS file
	HtmlTemplate     *template.Template     // The template
	CssTemplate      *textTemplate.Template // Used to add the theme variables
	JsTemplate       *textTemplate.Template // If you need to change some data
	JsGraphTemplate  *textTemplate.Template // Usually you'll need to update the graph base url
	JsCanvasTemplate *textTemplate.Template // If you need to change some data
}
