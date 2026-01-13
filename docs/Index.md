---
title: Bake your Obsidian vault into a website
description: Turn your Obsidian vault into a blazing fast, interactive website. Kiln natively supports Canvas, Graphs, Math, and Wikilinks with zero configuration required.
---
> [!Info] New: Custom Mode
> You can now use the flag `--mode "custom"` to activate custom mode, a feature that allows you to use Obsidian as an headless CMS. Check it out at [[What is Custom Mode]] or follow the [[Quick Start Guide]] to get started building custom static websites with Kiln. 
# Kiln: an Obsidian static website generator

> **Bake your Obsidian vault into a blazing fast static site.**

**Kiln** is a zero-config static site generator built specifically for [Obsidian](https://obsidian.md/). It takes your local Markdown vault—including your images, canvases, graphs, and math—and "bakes" it into a highly optimized, interactive HTML website.

Unlike other generators that require complex configuration or rigid folder structures, Kiln follows a **"Parity First"** philosophy: if it renders in Obsidian, it should render in Kiln.

## Get Started

Kiln is distributed as a single binary with no dependencies. You can be up and running in seconds.

### Installation
The recommended way to install Kiln is via Go:
```bash
go install github.com/otaleghani/kiln/cmd/kiln@latest
```

(See [[Installation]] for binary downloads).

### Bake your Site
Locate your Obsidian vault folder using [[generate]] and preview it locally with [[serve]]:
```bash
# Generate the static website from your vault
kiln generate --input ./path-to-your-vault --output ./output-directory

# Preview your site locally
kiln serve ./output-directory
```

Open `http://localhost:8080` to see your vault transformed into a website. 

## Features

Kiln is packed with features designed to bridge the gap between your personal knowledge base and a public-facing website.
### Content & Rendering

We support the tools you use to think.

- **[[Obsidian Markdown|Obsidian Parity]]**: Full support for **[[Wikilinks]]**, **[[Callouts]]**, and standard Markdown.
- **[[Obsidian Canvas|Interactive Canvas]]**: Render `.canvas` files as zoomable, pan-able diagrams directly on your site. Check it out [[Canvas.canvas|Canvas]]
- **[[Mermaid Graphs]]**: Native support for flowcharts, sequence diagrams, and Gantt charts.
- **[[Math|Math & LaTeX]]**: Beautiful equation rendering via MathJax.
- **[[Syntax Highlighting]]**: Automatic language detection and coloring for code blocks.
- **[[Bases]]**: Group, filter and view your notes like a database.

### Navigation & UX

Your site behaves like a modern app, not a static document.

- **[[Client Side Navigation|Instant Navigation]]**: Powered by **HTMX**, pages load instantly without full refreshes (SPA feel).
- **[[Explorer|File Explorer]]**: A "What You See Is What You Get" sidebar that mirrors your vault structure.
- **[[Search|Quick Find]]**: Real-time, fuzzy search to filter your file tree instantly.
- **[[Table of Contents|Table of Contents]]**: Auto-generated right sidebar navigation for every note.
- **[[Folders]]**: Navigate your notes with ease.
- **[[Tags]]**: Group together different notes from different folders.

### Visual Knowledge

Visualize how your ideas connect.

- **[[Global Graph]]**: An interactive visualization of your entire vault's connections.
- **[[Local Graph]]**: A context-aware network view specific to the current note.

### Technical & SEO

Built for performance and discoverability.

- **Zero-Config SEO**: Automatic generation of **[[Meta Tags]]**, **[[Sitemap xml|Sitemap xml]]**, and **[[Robots txt|Robots.txt]]**.
- **[[Themes|Theming]]**: Built-in [[Light-Dark Mode]] and customizable color palettes baked directly into CSS.
- **Linter**: A simple to use CLI linter using the [[doctor]] command to check your vault.
- Other commands: Check out the stats of your vault with the [[stats]] command or initialize a new project with the [[init]] command.

## Deployment

The output of Kiln is a static website, which makes it really easy to deploy on your server. You can follow one of our guide for some popular services, like [[Cloudflare Pages]], [[GitHub Pages]], [[Vercel]], [[Netlify]] or just one of the many [[Web Servers]].

## Contributing

Kiln is open source and distributed under the MIT License. If you spot a bug or have a feature request, please check our [[Roadmap]] or [open an issue on GitHub](https://github.com/otaleghani/kiln/releases/latest).

![[Roadmap#Planned immediate enhancements]]

## Credits

Check out the [[Credits]] page for a full list of credits.