---
title: "Generate Command — Build Your Obsidian Vault into a Static Site"
description: "Use kiln generate to convert your Obsidian vault into a fast static website. Customize themes, fonts, layouts, SEO settings, and URL structure with CLI flags."
---

# Generate Command

The `generate` command converts your Obsidian vault into a complete static HTML website. It processes every Markdown note, resolves [wikilinks](../Features/Navigation/Wikilinks.md), renders [Obsidian-flavored Markdown](../Features/Rendering/Obsidian Markdown.md), generates navigation and graph data, and outputs a ready-to-deploy site.

Run this command every time you add or update content in your vault.

## Usage

```bash
kiln generate [flags]
```

A minimal build with default settings:

```bash
kiln generate
```

This reads from `./vault`, writes to `./public`, and applies the default theme, font, and layout.

## Flags

| Flag                    | Short | Default   | Description                                                                                                                              |
| ----------------------- | ----- | --------- | ---------------------------------------------------------------------------------------------------------------------------------------- |
| `--theme`               | `-t`  | `default` | Sets the color scheme. See [Themes & Visuals](../Features/User Interface/Themes.md) for available options.                               |
| `--font`                | `-f`  | `inter`   | Sets the font family. See [Fonts & Typography](../Features/User Interface/Fonts.md) for available options.                               |
| `--url`                 | `-u`  | `""`      | The public URL of your site (e.g., `https://example.com`). Required for [Sitemap.xml](../Features/SEO/Sitemap xml.md) and [Robots.txt](../Features/SEO/Robots txt.md) generation. |
| `--name`                | `-n`  | `My Notes` | The site name displayed in the browser tab and [Meta Tags & SEO](../Features/SEO/Meta Tags.md).                                        |
| `--input`               | `-i`  | `./vault` | Path to the source directory containing your Markdown notes.                                                                             |
| `--output`              | `-o`  | `./public`| Path where the generated HTML files are saved.                                                                                           |
| `--mode`                | `-m`  | `default` | Build mode. Use `default` for standard vault rendering or `custom` for [Custom Mode](../Features/Custom Mode/What is Custom Mode.md) with collection configs and templates. |
| `--layout`              | `-L`  | `default` | Page layout to use. See [Layouts](../Features/User Interface/Layouts.md) for available options.                                          |
| `--flat-urls`           |       | `false`   | Generate flat files (`note.html`) instead of directories (`note/index.html`).                                                            |
| `--disable-toc`         |       | `false`   | Hides the [Table of Contents](../Features/User Interface/Table of Contents.md) from the right sidebar.                                   |
| `--disable-local-graph` |       | `false`   | Hides the [Local Graph](../Features/User Interface/Local Graph.md) from the right sidebar. Disabling TOC, local graph, and backlinks removes the right sidebar entirely. |
| `--disable-backlinks`   |       | `false`   | Hides the [[Backlinks]] panel from the right sidebar.                                                                    |
| `--log`                 | `-l`  | `info`    | Log verbosity. Choose `info` or `debug`.                                                                                                 |

## What Gets Generated

The build produces a complete static site including:

- **HTML pages** for every Markdown note, folder index, tag page, and [Canvas](../Features/Rendering/Obsidian Canvas.md) file in your vault.
- **CSS and JavaScript** — theme styles, font files, sidebar navigation, [search](../Features/Navigation/Search.md), and the interactive [Global Graph](../Features/User Interface/Global Graph.md).
- **SEO files** — [sitemap.xml](../Features/SEO/Sitemap xml.md) and [robots.txt](../Features/SEO/Robots txt.md) when you provide `--url`.
- **Static assets** — images, PDFs, and attachments are copied to the output directory.
- **Special files** — `CNAME`, `favicon.ico`, and `_redirects` are carried over if present in your vault.

The output directory is cleaned automatically before each build, so there are no stale files from previous runs.

## Examples

### Production Build

When deploying to the web, always include `--url` and `--name` so that SEO features, sitemaps, and meta tags work correctly:

```bash
kiln generate \
  --name "My Digital Garden" \
  --url "https://notes.mysite.com" \
  --theme "nord" \
  --font "inter"
```

### Custom Layout and Subdirectory Deployment

Build with a specific layout and a base path for hosting under a subdirectory:

```bash
kiln generate \
  --name "Documentation" \
  --url "https://example.com/docs" \
  --layout "simple"
```

The [Serve Command](./serve.md) respects this base path when previewing locally.

### Minimal Sidebar

Remove the table of contents and local graph to create a cleaner reading experience:

```bash
kiln generate --disable-toc --disable-local-graph
```

### Custom Mode

Use [Custom Mode](../Features/Custom Mode/What is Custom Mode.md) to build a site with collection configs and custom templates instead of the default vault layout:

```bash
kiln generate --mode custom --input ./my-project
```

## Full Workflow

A typical workflow from a fresh project to a live local preview:

```bash
# 1. Scaffold a new vault
kiln init

# 2. Check for broken links
kiln doctor --input ./vault

# 3. Build the site
kiln generate --name "My Notes" --url "https://notes.example.com" --theme dracula

# 4. Preview locally
kiln serve
```

See the [Init Command](./init.md), [Doctor Command](./doctor.md), and [Serve Command](./serve.md) for details on each step. After verifying the output, deploy to [GitHub Pages](../Deployment/GitHub Pages.md), [Netlify](../Deployment/Netlify.md), [Vercel](../Deployment/Vercel.md), [Cloudflare Pages](../Deployment/Cloudflare Pages.md), or any [static web server](../Deployment/Web Servers.md).
