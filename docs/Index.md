---
title: Bake your Obsidian vault into a website
description: Turn your Obsidian vault into a blazing fast, interactive website. Kiln natively supports Canvas, Graphs, Math, and Wikilinks with zero configuration required.
---
# Kiln: Your Obsidian Vault, Online.

> **Zero Config. Zero Compromise. Blazing Fast.**

**Kiln** turns your local Obsidian vault into a high-performance, interactive website. It is built with a **"Parity First"** philosophy: if it works in Obsidian, it works in your browser.

No complex config files, no rigid folder structures, and no broken links. Just write, bake, and ship.

[Get Started](#installation) · [View Demo](https://kiln.talesign.com) · [Roadmap]([[Roadmap]])

---
## Why Kiln?
Most static site generators require you to fight against your tools. Kiln embraces them.

### Parity First Rendering
Stop worrying if your notes will break when you publish them. Kiln natively supports the features you rely on to think:
* **[[Obsidian Canvas|Interactive Canvas]]**: Render your `.canvas` files as fully zoomable, pan-able diagrams.
* **[[Mermaid Graphs]]**: Flowcharts, sequence diagrams, and Gantt charts work out of the box.
* **[[Math|Math & LaTeX]]**: Complex equations rendered beautifully via MathJax.
* **[[Callouts]] & [[Wikilinks]]**: First-class citizens, not afterthoughts.

### Built for Speed (SPA Feel)
Kiln sites don't feel like static documents; they feel like apps.
* **Powered by HTMX**: We use [[Client Side Navigation]] to load pages instantly without full refreshes. 
* **Single Binary**: Written in Go. No Node_modules, no Ruby gems, no dependencies. Just a simple binary.

### Visual Navigation
Your knowledge is a network, not a list.
* **[[Global Graph]]**: Visualize your entire vault's connections.
* **[[Local Graph]]**: Context-aware network views for every note.
* **[[Explorer|File Explorer]]**: A "What You See Is What You Get" sidebar that mirrors your vault structure exactly.

### Make it yours
Every vault is different, so let's reflect that!
- [[Themes]]: Choose from a big collection of themes.
- [[Fonts]]: We have them! Tons of them!
- [[Light-Dark Mode]]: Every theme comes with a light and a dark mode.

---

## Get Started in Seconds
Kiln is a single binary. You can go from Vault to Website in under 60 seconds.

### Installation
```bash
# Install via Go
go install github.com/otaleghani/kiln/cmd/kiln@latest
```

(See [[Installation]] for pre-compiled binaries for Windows, Mac, and Linux).

### Bake & Serve
Locate your Obsidian vault and run:

```bash
# Generate the site
kiln generate --input ./my-vault --output ./public

# Preview locally
kiln serve ./public
```

Open `http://localhost:8080`. Your digital garden is now live. 

## Deployment
Kiln outputs standard HTML/CSS/JS. You can host it _anywhere_. Check our guides for [[Cloudflare Pages]], [[GitHub Pages]], [[Vercel]], [[Netlify]], or standard [[Web Servers]].

---

## Advanced Features
Ready to go beyond a simple knowledge base?

> [!warning] NEW: Custom Mode 
> Want to use Obsidian as a Headless CMS? You can now use the flag `--mode "custom"` to take full control of the output. Check out [[What is Custom Mode]] or follow the [[Quick Start Guide]] to build completely custom layouts.

### Data as a First-Class Citizen
Organize your knowledge with **[[Bases]]**. Group, filter, and view your notes like a database—perfect for project trackers, book lists, or research logs.

### Zero-Config SEO
Search optimization made easy. Kiln handles it for you:

- Automatic **[[Meta Tags]]** & Open Graph support (_soon_).
- Auto-generated **[[Sitemap xml|Sitemap.xml]]** and **[[Robots txt|Robots.txt]]**.
- Built-in **[[Themes|Light/Dark Mode]]**.

---

## Community & Contributing
Kiln is open source (MIT).

- **Found a bug?** [Open an issue on GitHub](https://github.com/otaleghani/kiln/issues).
- **Want to help?** Check the [[Roadmap]].
- **Love the project?** Give us a star!
- **Interested about the how?**: Checkout the [[Credits]]!