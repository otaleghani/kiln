// Package assets handles the embedding of static resources into the compiled binary.
// This allows the Kiln CLI to be distributed as a single, standalone executable
// without requiring users to manage external template or css files.
package assets

import "embed"

// TemplateFS is a virtual filesystem containing all the embedded static assets.
// The //go:embed directives below tell the Go compiler to include these specific files
// directly into the binary at build time.
//
// Included Assets:
// - JavaScript: Logic for the interactive graph, canvas rendering, and app behavior.
// - HTML/CSS: The base layout templates and styling.
// - Fonts: WOFF2 font files for the supported typography themes (Inter, Lato, Merriweather).
//
//go:embed canvas.js
//go:embed graph.js
//go:embed layout.html
//go:embed style.css
//go:embed app.js
//go:embed Inter-Regular.woff2
//go:embed Inter-Bold.woff2
//go:embed Lato-Regular.woff2
//go:embed Lato-Bold.woff2
//go:embed Merriweather-Regular.woff2
//go:embed Merriweather-Bold.woff2
//go:embed default_layout.html
//go:embed default_style.css
//go:embed default_app.js
var TemplateFS embed.FS
