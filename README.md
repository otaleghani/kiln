# Kiln

> **Bake your Obsidian vault into a blazing fast static site.**

![Build Status](https://img.shields.io/github/actions/workflow/status/otaleghani/kiln/deploy-docs.yml?style=flat-square) ![Go Version](https://img.shields.io/badge/go-1.25+-00ADD8?style=flat-square&logo=go)
![License](https://img.shields.io/badge/license-MIT-blue?style=flat-square)

**Kiln** is a static site generator designed specifically for [Obsidian](https://obsidian.md/). It takes your Markdown vault—images, canvases, graphs, and math included—and "bakes" it into a highly optimized, interactive HTML website. 

## ✨ Features

**Content & Rendering**
- 📝 **Full Markdown Support:** Renders almost every Obsidian-flavored markdown tag.
- 🎨 **Obsidian Canvas:** Renders `.canvas` files directly into interactive diagrams.
- 📊 **Mermaid Graphs:** Native support for flowcharts, sequence diagrams, and more.
- ➗ **Math Expressions:** Beautiful LaTeX rendering via MathJax.
- 🌗 **Theming:** Built-in Light and Dark modes (respects system preference).

**Navigation & UX**
- 🕸️ **Graph View:** Interactive global graph (`/graph`) and local graph per note.
- ⚡ **HTMX Powered:** Instant, client-side page reloading without full refreshes.
- 🔍 **Search:** Instant note search functionality.
- 📱 **Responsive:** Mobile-friendly design out of the box.

**Technical & SEO**
- 🚀 **SEO Optimized:** Automatic meta tags and social sharing previews.
- 🤖 **Bot Friendly:** Auto-generates `sitemap.xml` and `robots.txt`.
- ⚡ **Zero Config:** Smart defaults for everything.

---

## 📦 Installation

### Option 1: Download Binary (Recommended)
Download the latest version for your operating system from the [Releases Page](https://github.com/otaleghani/kiln/releases).

### Option 2: Install via Go
If you have Go installed, you can build from source:

```bash
go install github.com/otaleghani/kiln@latest
```

---

## 🚀 Quick Start

Build your vault
```bash
kiln generate
```
This creates a output/ folder ready to be uploaded to Netlify, Vercel, or GitHub Pages.

Preview your site locally
```bash
kiln serve
```
Open http://localhost:8080 to see your vault baked into a website.

---

## 🛠 Command Reference
```bash
kiln init	    # Initializes a new Kiln project in the current directory
kiln generate	# Builds the static site from your vault into the output folder
kiln serve	    # Starts a local web server to preview your generated site
kiln doctor	    # Scans your vault for broken internal links 
kiln stats	    # Displays insights and statistics about your vault (longest note, total words etc.)
kiln clear	    # Removes the public output directory to ensure a clean build
```

___

## ⚠️ Quirks & Requirements

To keep Kiln simple, it has a few opinionated behaviors:

- Ignored Folders: Kiln automatically ignores the .obsidian folder and any other file or folder starting with a dot (.).
- Favicon: To have a custom favicon, place a favicon.ico file in the root of your vault. Kiln will automatically detect and use it.
- CNAME: Kiln copies over a CNAME file if it finds it

---

## ☁️ Deployment (GitHub Actions)

You can automate the building of your site using GitHub Actions. Create a file at .github/workflows/build.yml in your vault's repository.

See the Wiki for full deployment examples.

___

## 🗺 Roadmap

- [ ] Fix hover colors
- [ ] Configuration file support (kiln.yaml)
- [ ] Automatic image optimization (WebP conversion)
- [ ] PDF preview support
- [ ] Custom themes support
- [ ] Custom font selection
- [ ] Configurable input/output directories
- [ ] Localization / i18n
- [ ] Frontmatter display toggle
- [ ] Custom colors for Canvas elements
- [ ] View transitions for page navigation
- [ ] CSS animations for buttons and links
- [ ] Folder pages
- [ ] Page metadata (last update, time to read, etc.)
- [ ] Hover on links to preview note
- [ ] Backlinks on right sidebar
- [ ] Breadcrumbs links
- [ ] Tags management
- [ ] Search for tags
- [ ] Fuzzy finder (full-text search)
- [ ] Hide certain pages
- [ ] RSS feed
- [ ] Social media preview cards
- [ ] Draft management
- [ ] Custom permalink

___

## 📄 License

Distributed under the MIT License. See LICENSE for more information.
