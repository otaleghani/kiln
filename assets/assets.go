package assets

import "embed"

//go:embed canvas.js
//go:embed graph.js
//go:embed layout.html
//go:embed style.css
//go:embed app.js
var TemplateFS embed.FS
