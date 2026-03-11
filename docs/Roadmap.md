---
title: "Kiln Roadmap — Upcoming Features and Planned Improvements"
description: "Explore planned features for Kiln: RSS feeds, full-text search, image optimization, View Transitions, and more. See what's next for the Obsidian static site generator."
---
# Roadmap

Kiln is an actively developed Obsidian-to-website generator with frequent releases. This roadmap outlines planned features, upcoming improvements, and long-term goals. For a full list of what has already shipped, see the [Changelog](./Changelog.md).

> [!note] Contributing
> Have a suggestion or feature request? Feel free to open an issue or contribute to the discussion on our [GitHub repository](https://github.com/otaleghani/kiln/issues).

## Recently Completed

These features were on the roadmap and have since been shipped:

- [x] **Custom 404 Pages:** User-defined `404.md` pages are now supported out of the box.
- [x] **Themes & Fonts:** A full collection of [themes](./Features/User Interface/Themes.md) and [fonts](./Features/User Interface/Fonts.md) is available, including Catppuccin, Nord, Dracula, and more.
- [x] **Comments:** Giscus-powered [comment sections](./Features/User Interface/Comments.md) can be added to any page.
- [x] **SEO Essentials:** Automatic [sitemap.xml](./Features/SEO/Sitemap xml.md), [robots.txt](./Features/SEO/Robots txt.md), [canonical tags](./Features/SEO/Canonical.md), and [meta tags](./Features/SEO/Meta Tags.md) generation.

## Core Features
*Major functionality upgrades planned for upcoming releases.*

- [x] **Configuration File:** Support for `kiln.yaml` to save preferences (flags) persistently, so you don't need to pass CLI flags on every [generate](./Commands/generate.md) run.
- [x] **Image Optimization:** Automatic WebP conversion and responsive image sizing to improve page load speed.
- [x] **Localization:** i18n support for translating the UI into multiple languages.
- [x] **Structured Data:** Automatically generate JSON-LD schema markup for articles and breadcrumbs to improve search engine visibility.
- [x] Support for other kinds of links, not only wikilinks 

## UI & UX Improvements
*Enhancing the reading and navigation experience.*

- [x] **Page Transitions:** Smooth View Transitions API support to complement the existing [client-side navigation](./Features/Navigation/Client Side Navigation.md).
- [x] **Link Previews:** Obsidian-style hover previews when mousing over [wikilinks](./Features/Navigation/Wikilinks.md).
- [x] **Full-Text Search:** Fuzzy finder with content search to go beyond the current name-based [search](./Features/Navigation/Search.md).
- [x] **Navigation Enhancements:** "Back to Top" buttons and improved breadcrumb trails.
- [x] **Animations:** Subtle CSS animations for interactive elements.
- [x] **Canvas Colors:** Custom color support for [Canvas](./Features/Rendering/Obsidian Canvas.md) nodes.

## Content Management
*Tools to help you manage and publish your digital garden.*

- [x] **Reading Metadata:** Display "Last Updated" and "Time to Read" indicators on every page.
- [ ] **Backlinks Panel:** A dedicated sidebar section showing all incoming links to the current page.
- [x] **RSS Feeds:** Automatic RSS/Atom feed generation for blog-style vaults, making it easy for readers to subscribe.
- [ ] **Draft Support:** Exclude work-in-progress pages from the build using frontmatter flags. For now, you can use [hidden files and folders](./Features/Rendering/Hidden files folders.md) to keep content out of the generated site.
- [ ] **Privacy Controls:** Hide specific pages or folders from the published site via frontmatter.

## Advanced Configuration
*Power-user features for fine-grained control over the generated site.*

- [ ] **Permalinks:** Custom URL slugs defined in frontmatter for complete control over your site's URL structure.
- [x] **Social Cards:** Auto-generated Open Graph images for rich social media previews when sharing links.

## How to Stay Updated

New features and fixes ship regularly. Check the [Changelog](./Changelog.md) for the latest releases, or watch the [GitHub repository](https://github.com/otaleghani/kiln) for updates. To get started with Kiln today, see the [Installation](./Installation.md) guide.
