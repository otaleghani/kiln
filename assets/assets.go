package assets

import "embed"

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
var TemplateFS embed.FS
