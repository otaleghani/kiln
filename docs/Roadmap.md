---
title: Roadmap
description: Explore the future of Kiln. A list of upcoming features, planned fixes, and long-term goals for the project.
---
# Roadmap

Kiln is an evolving project. Below is the current plan for upcoming features and improvements.

> [!note] Contributing
> Have a suggestion? Feel free to open an issue or contribute to the discussion on our GitHub repository.

## Immediate Fixes
*Priority items to stabilize the current beta.*

- [ ] **Better logs**: Create better logs management using `slog`
- [ ] **Support bases**: Support for Obsidian bases.

## Core Features
*Major functionality upgrades.*

- [ ] **Configuration File:** Support for `kiln.yaml` to save preferences (flags) persistently.
- [ ] **Performance:** Automatic image optimization (WebP conversion).
- [ ] **Customization:** Full support for custom CSS themes and user-defined fonts.
- [ ] **Flexible Paths:** Configurable input/output directories in the config file.
- [ ] **Localization:** i18n support for translating the UI.

## UI & UX Improvements
*Enhancing the reading experience.*

- [ ] **Page Transitions:** Smooth View Transitions API support for navigation.
- [ ] **Link Previews:** Obsidian-style hover previews when mousing over Wikilinks.
- [ ] **Search:** Fuzzy finder implementation for full-text search.
- [ ] **Navigation:** Breadcrumb trails and "Back to Top" buttons.
- [ ] **Animations:** Subtle CSS animations for interactive elements.
- [ ] **Canvas:** Custom color support for Canvas nodes.

## Content Management
*Tools to help you manage your digital garden.*

- [ ] **Metadata:** Display "Last Updated" and "Time to Read" on pages.
- [ ] **Tags:** Dedicated tag pages and tag-based search.
- [ ] **Backlinks:** A dedicated section in the sidebar showing incoming links.
- [ ] **RSS:** Automatic RSS feed generation for blog posts.
- [ ] **Drafts:** Better handling for draft pages (exclude from build).
- [ ] **Folder Pages:** Auto-generated index pages for folders.
- [ ] **Privacy:** Ability to hide specific pages or folders via frontmatter.

## Advanced Configuration
- [ ] **Permalinks:** Custom URL slugs via frontmatter.
- [ ] **Social Cards:** Auto-generated Open Graph images for social sharing.
- [ ] **Custom 404:** Support for a user-defined `404.md`.
- [ ] **Layouts:** Support for alternative page layouts (e.g., full-width, landing page).
- [ ] **Custom layouts**: Create custom layouts for specific folders.