---
title: Credits & Open-Source Acknowledgements
description: Open-source libraries, fonts, themes, and tools that power Kiln — the Obsidian vault to static site generator.
---

# Credits & Acknowledgements

Kiln is built on the shoulders of giants. Every feature — from [Markdown rendering](./Features/Rendering/Obsidian Markdown.md) to [syntax highlighting](./Features/Rendering/Syntax Highlighting.md) to [graph visualization](./Features/User Interface/Global Graph.md) — relies on a robust ecosystem of open-source tools. This page lists the libraries, fonts, and color palettes that make Kiln possible.

## Core Technology

* **[Go](https://go.dev/):** Kiln is written in Go, chosen for its speed, concurrency model, and ability to compile into a single, dependency-free binary. This means [installing Kiln](./Installation.md) requires no runtime dependencies — just download and run.
* **[Cobra](https://github.com/spf13/cobra):** The CLI framework behind every Kiln command, from [generate](./Commands/generate.md) to [serve](./Commands/serve.md) to [doctor](./Commands/doctor.md).

## Markdown & Build Libraries

These libraries run during the [build process](./Commands/generate.md) to transform your Obsidian Markdown into HTML.

* **[Goldmark](https://github.com/yuin/goldmark):** A standards-compliant, extensible Markdown parser. Goldmark handles the heavy lifting of converting your notes into web pages, with extensions for GFM tables, strikethrough, and task lists.
* **[goldmark-meta](https://github.com/yuin/goldmark-meta):** Parses YAML frontmatter from your notes, powering [Meta Tags & SEO](./Features/SEO/Meta Tags.md) features like page titles and descriptions.
* **[goldmark-wikilink](https://go.abhg.dev/goldmark/wikilink):** Resolves `[[wikilinks]]` and embeds, enabling Kiln's native [Wikilinks & Embeds](./Features/Navigation/Wikilinks.md) support.
* **[goldmark-mathjax](https://github.com/litao91/goldmark-mathjax):** Detects LaTeX delimiters in Markdown so Kiln can render [Math & LaTeX](./Features/Rendering/Math.md) expressions.
* **[Chroma](https://github.com/alecthomas/chroma):** A general-purpose syntax highlighter that powers [Syntax Highlighting](./Features/Rendering/Syntax Highlighting.md) of code blocks without relying on external CSS or JavaScript.
* **[Minify](https://github.com/tdewolff/minify):** Compresses the generated HTML, CSS, and JavaScript output to minimize page size and improve load times.

## Frontend Libraries

These libraries are loaded by the browser to provide interactivity and visual rendering.

* **[HTMX](https://htmx.org/):** Powers [Client Side Navigation](./Features/Navigation/Client Side Navigation.md), loading pages instantly without full browser refreshes — giving Kiln sites an app-like feel.
* **[D3.js](https://d3js.org/):** The force simulation engine behind Kiln's interactive [Global Graph](./Features/User Interface/Global Graph.md) and [Local Graph](./Features/User Interface/Local Graph.md) visualizations.
* **[Mermaid.js](https://mermaid.js.org/):** A JavaScript-based diagramming and charting tool that renders Markdown definitions into [Mermaid Graphs](./Features/Rendering/Mermaid Graphs.md).
* **[MathJax](https://www.mathjax.org/):** An open-source JavaScript display engine for LaTeX, enabling high-quality [mathematical typesetting](./Features/Rendering/Math.md).
* **[Giscus](https://giscus.app/):** A GitHub Discussions-powered commenting system that enables [Comments](./Features/User Interface/Comments.md) on your published notes.

## Typography

Kiln embeds open-source fonts directly into the build output, so your site loads without external network requests. You can switch between them using the [Fonts & Typography](./Features/User Interface/Fonts.md) settings.

* **[Inter](https://rsms.me/inter/):** Designed by Rasmus Andersson. A typeface carefully crafted for computer screens. This is the default font.
* **[Merriweather](https://fonts.google.com/specimen/Merriweather):** Designed by Sorkin Type. A serif font designed to be highly readable on screens.
* **[Lato](https://fonts.google.com/specimen/Lato):** Designed by Łukasz Dziedzic. A humanist sans-serif font.
* **[Lora](https://fonts.google.com/specimen/Lora):** A calligraphy-inspired serif with roots in 18th-century typography.
* **[Libre Baskerville](https://fonts.google.com/specimen/Libre+Baskerville):** A webfont optimized for body text, based on the classic Baskerville typeface.
* **[Noto Serif](https://fonts.google.com/noto/specimen/Noto+Serif):** Part of Google's Noto family, designed for broad language coverage.
* **[IBM Plex Sans](https://fonts.google.com/specimen/IBM+Plex+Sans):** IBM's corporate typeface, open-sourced for everyone.
* **[Roboto](https://fonts.google.com/specimen/Roboto):** Google's default Android font, widely used across the web.

## Visual Themes

Kiln includes ports of community-created color palettes, each with full [Light & Dark Mode](./Features/User Interface/Light-Dark Mode.md) support. Browse all options on the [Themes & Visuals](./Features/User Interface/Themes.md) page.

* **[Catppuccin](https://catppuccin.com/):** A community-driven pastel theme that aims to be the middle ground between high and low contrast.
* **[Nord](https://www.nordtheme.com/):** An arctic, north-bluish color palette designed for a clean and clutter-free workflow.
* **[Dracula](https://draculatheme.com/):** A dark theme famous for its vibrant colors and high contrast.
* **[Tokyo Night](https://github.com/enkia/tokyo-night-vscode-theme):** Inspired by the lights of Tokyo at night, with rich blues and purples.
* **[Rosé Pine](https://rosepinetheme.com/):** A low-contrast palette with muted, natural tones.
* **[Gruvbox](https://github.com/morhetz/gruvbox):** A retro-groove color scheme with warm, earthy tones.
* **[Everforest](https://github.com/sainnhe/everforest):** A green-tinted theme inspired by natural forest colors.
* **[Cyberdream](https://github.com/scottmckendry/cyberdream.nvim):** A high-contrast, futuristic color scheme with vivid neon accents.
