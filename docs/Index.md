---
title: "Kiln — Obsidian Static Site Generator"
description: "Turn your Obsidian vault into a fast, SEO-ready static website. Supports Canvas, Wikilinks, Graphs, Math, themes, and deployment anywhere."
---
# Kiln: Obsidian Vault to Static Website Generator

> **Zero Config. Zero Compromise. Blazing Fast.**

**Kiln** is an open-source static site generator that turns your [Obsidian](https://obsidian.md) vault into a fast, interactive website. Built with a **"Parity First"** philosophy: if it works in Obsidian, it works in your browser — [[Wikilinks|wikilinks]], [[Obsidian Canvas|canvas files]], [[Math|LaTeX math]], [[Callouts|callouts]], and all.

No config files, no rigid folder structures, and no broken links. Write your notes in Obsidian, generate a site with a single command, and deploy anywhere.

[Get Started](#get-started-in-seconds) · [View Demo](https://kiln.talesign.com) · [[Roadmap]]

---
## Why Choose Kiln to Publish Your Obsidian Vault?

Most static site generators force you to rewrite your content. Kiln renders your vault as-is — every feature Obsidian supports, Kiln publishes faithfully.

### Obsidian Feature Support — Canvas, Wikilinks, Math & More

Publish your notes without worrying about compatibility. Kiln natively renders the Obsidian features you rely on:

* **[[Obsidian Canvas|Interactive Canvas]]**: Render `.canvas` files as zoomable, pannable diagrams.
* **[[Mermaid Graphs]]**: Flowcharts, sequence diagrams, and Gantt charts work out of the box.
* **[[Math|Math & LaTeX]]**: Complex equations rendered via MathJax.
* **[[Callouts]]**: Styled info boxes, warnings, and collapsible blocks.
* **[[Wikilinks]]**: Internal links and embeds with full [[Obsidian Markdown]] support, including [[Syntax Highlighting|syntax-highlighted code blocks]].

### Instant Page Loads with Client-Side Navigation

Kiln sites feel like apps, not static documents.

* **Powered by HTMX**: [[Client Side Navigation]] loads pages instantly without full-page refreshes.
* **Single Binary**: Written in Go with zero runtime dependencies — no Node.js, no Ruby gems. Just download and run.

### Graphs, Explorer & Search

Your knowledge is a network, not a list. Kiln provides multiple ways to navigate your published notes:

* **[[Global Graph]]**: Visualize your entire vault's connections at a glance.
* **[[Local Graph]]**: Context-aware network views on every note.
* **[[Explorer|File Explorer]]**: A sidebar that mirrors your vault's [[Folders|folder structure]] exactly.
* **[[Search]]**: Filter and find notes across your entire site.
* **[[Tags]]**: Browse notes by topic using Obsidian's tag system.

### Themes, Fonts & Customization

Every vault is different — make your site reflect that.

- **[[Themes]]**: Choose from a large collection of built-in themes.
- **[[Fonts]]**: Pick from dozens of bundled fonts and typography options.
- **[[Light-Dark Mode]]**: Every theme ships with both light and dark variants.
- **[[Layouts]]**: Control the structure of your pages with flexible layout options.

---

## Get Started in Seconds

Kiln ships as a single binary. Go from vault to website in under a minute.

### Install Kiln

```bash
# Install with Go (recommended)
go install github.com/otaleghani/kiln/cmd/kiln@latest
```

Pre-compiled binaries are available for macOS, Linux, and Windows — see the full [[Installation]] guide for download links and checksum verification.

### Generate and Preview Your Site

Point Kiln at your Obsidian vault and run two commands:

```bash
# Generate the static site
kiln generate --input ./my-vault --output ./public

# Preview it locally
kiln serve ./public
```

Open `http://localhost:8080` and your site is live. The [[Generate Command]] accepts flags for themes, fonts, base URL, and more — run `kiln generate --help` to see all options. Use the [[Serve Command]] to preview changes locally before deploying.

Want to scaffold a new project from scratch? The [[Init Command]] creates a vault pre-configured for Kiln.

---

## Deploy Your Obsidian Site Anywhere

Kiln outputs standard HTML, CSS, and JS — host your site on any platform. Follow one of the step-by-step deployment guides:

- [[Cloudflare Pages|Deploy on Cloudflare Pages]]
- [[GitHub Pages|Deploy on GitHub Pages]]
- [[Netlify|Deploy on Netlify]]
- [[Vercel|Deploy on Vercel]]
- [[Web Servers|Deploy on Nginx, Apache, or Caddy]]

---

## Built-In SEO for Obsidian Sites

Kiln handles search engine optimization automatically so your published notes get discovered:

- Automatic **[[Meta Tags]]** and Open Graph tags for rich previews when shared on social media.
- Auto-generated **[[Sitemap xml|Sitemap.xml]]** and **[[Robots txt|Robots.txt]]** so search engines can crawl your site.
- **[[Canonical tag|Canonical URLs]]** to prevent duplicate content issues across your pages.

---

## Advanced: Custom Mode

Ready to go beyond a default knowledge base?

> [!new] NEW: Custom Mode
> Use Obsidian as a headless CMS. Pass the `--mode "custom"` flag to take full control of your output with custom HTML templates. Learn more in [[What is Custom Mode]] or follow the [[Quick Start Guide]] to build a blog from scratch.

### Data-Driven Pages with Bases

Organize your knowledge with **[[Bases]]**. Group, filter, and sort your notes like a database — perfect for project trackers, book lists, or research logs.

---

## Maintain Your Site

Keep your published vault in shape with Kiln's built-in utilities:

- **[[Doctor Command|Doctor]]**: Scan your vault for broken links and common issues before publishing.
- **[[Clean Command|Clean]]**: Remove stale build output and start fresh.
- **[[Stats Command|Stats]]**: View word counts and note metrics across your vault.

---

## Community and Contributing

Kiln is open source under the MIT license.

- **Found a bug?** [Open an issue on GitHub](https://github.com/otaleghani/kiln/issues).
- **Want to help?** Check the [[Roadmap]] for planned features and open tasks.
- **Curious about the stack?** See the [[Credits]] page.
- **View the [[Changelog]]** for recent updates and release notes.
