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
//go:embed giscus_theme_light.css
//go:embed giscus_theme_dark.css
//go:embed canvas.js
//go:embed graph.js
//go:embed shared.css
//go:embed Inter-Regular.woff2
//go:embed Inter-Bold.woff2
//go:embed Lato-Regular.woff2
//go:embed Lato-Bold.woff2
//go:embed Merriweather-Regular.woff2
//go:embed Merriweather-Bold.woff2
//go:embed GoogleSans-Regular.woff2
//go:embed GoogleSans-Bold.woff2
//go:embed IBMPlexSans-Regular.woff2
//go:embed IBMPlexSans-Bold.woff2
//go:embed LibreBaskerville-Regular.woff2
//go:embed LibreBaskerville-Bold.woff2
//go:embed Lora-Regular.woff2
//go:embed Lora-Bold.woff2
//go:embed NotoSerif-Regular.woff2
//go:embed NotoSerif-Bold.woff2
//go:embed Roboto-Regular.woff2
//go:embed Roboto-Bold.woff2
//go:embed default_layout.html
//go:embed default_style.css
//go:embed default_app.js
//go:embed simple_layout.html
//go:embed simple_style.css
//go:embed simple_app.js
var TemplateFS embed.FS
