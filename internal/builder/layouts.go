// Layout resolution, template loading, and template function registration. @feature:layouts
package builder

import (
	"log/slog"
	"strings"
	textTemplate "text/template"

	"github.com/a-h/templ"
	"github.com/otaleghani/kiln/assets"
	"github.com/otaleghani/kiln/internal/templates"
)

// resolveLayout looks up a Layout by name.
func resolveLayout(name string, log *slog.Logger) *Layout {
	log.Info("Resolving layout...", "name", name)
	layout, ok := layouts[strings.ToLower(name)]
	if !ok {
		log.Warn("Layout not found, using default layout", "name", "default")
		name = "default"
		layout = layouts[name]
	}

	layout.log = log

	return layout
}

// layouts is a key-value pairs of all available layouts
var layouts = map[string]*Layout{
	"default": {
		Name:    "default",
		CssPath: "default_style.css",
		JsPath:  "default_app.js",
		TemplRender: func(d *templates.PageData) templ.Component {
			return templates.DefaultLayout(d)
		},
	},

	"simple": {
		Name:    "simple",
		CssPath: "simple_style.css",
		JsPath:  "simple_app.js",
		TemplRender: func(d *templates.PageData) templ.Component {
			return templates.SimpleLayout(d)
		},
	},
}

// loadLayout loads the given layout files into memory
func (l *Layout) loadLayout() error {
	l.log.Debug("Loading layout", "name", l.Name)

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

	// Load and parse the canvas JS template
	jsCanvasTemplate, err := assets.TemplateFS.ReadFile("canvas.js")
	if err != nil {
		return err
	}
	tmplCanvasJS, err := textTemplate.New("js").Parse(string(jsCanvasTemplate))
	if err != nil {
		return err
	}
	l.JsCanvasTemplate = tmplCanvasJS

	// Load and parse the search JS template
	jsSearchContent, err := assets.TemplateFS.ReadFile("search.js")
	if err != nil {
		return err
	}
	tmplSearchJS, err := textTemplate.New("js").Parse(string(jsSearchContent))
	if err != nil {
		return err
	}
	l.JsSearchTemplate = tmplSearchJS

	// Load and parse the giscus CSS template
	cssGiscusLightContent, err := assets.TemplateFS.ReadFile("giscus_theme_light.css")
	if err != nil {
		return err
	}
	tmplGiscusLightCSS, err := textTemplate.New("css").Parse(string(cssGiscusLightContent))
	if err != nil {
		return err
	}
	l.CssGiscusLightTemplate = tmplGiscusLightCSS

	// Discus dark theme
	cssGiscusDarkContent, err := assets.TemplateFS.ReadFile("giscus_theme_dark.css")
	if err != nil {
		return err
	}
	tmplGiscusDarkCSS, err := textTemplate.New("css").Parse(string(cssGiscusDarkContent))
	if err != nil {
		return err
	}
	l.CssGiscusDarkTemplate = tmplGiscusDarkCSS

	return nil
}

// Layout contains the paths for the layout
//
// Every Layout is made by 3 different files, style.css, and app.js.
// If you are creating a new layout create the needed files by prepending the name of the
// layout to the name of the file. If you have a layout called "default" it should have the
// following files: default_style.css and default_app.js
type Layout struct {
	Name                   string
	CssPath                string                                        // Path of the CSS file
	JsPath                 string                                        // Path of the JS file
	TemplRender            func(*templates.PageData) templ.Component     // Templ layout renderer
	CssTemplate            *textTemplate.Template                        // Used to add the theme variables
	JsTemplate             *textTemplate.Template                        // If you need to change some data
	JsGraphTemplate        *textTemplate.Template                        // Usually you'll need to update the graph base url
	JsCanvasTemplate       *textTemplate.Template                        // If you need to change some data
	JsSearchTemplate       *textTemplate.Template                        // Full-text search JS
	CssGiscusLightTemplate *textTemplate.Template                        // Giscus template
	CssGiscusDarkTemplate  *textTemplate.Template                        // Giscus template
	log                    *slog.Logger
}
