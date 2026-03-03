---
title: "Kiln Quirks — File Overrides, URL Conflicts & Edge Cases"
description: "Known quirks when converting an Obsidian vault with Kiln: how duplicate filenames resolve, folder page overrides, and flat vs pretty URL behavior."
---
# Kiln Quirks

When Kiln converts your Obsidian vault into a static site, a few edge cases arise from the differences between how Obsidian organizes files and how websites serve pages. This page covers the known quirks you should keep in mind.

## Same filename, multiple extensions

Obsidian allows you to have multiple items with the exact same name living in the same directory. For example, you might have a **folder** named `Example`, a **canvas** named `Example.canvas`, and a **note** named `Example.md` all sitting side-by-side.

Instead of generating separate, conflicting pages for each, Kiln consolidates them into a single URL based on a strict override hierarchy (lowest to highest priority):

```
Folder → Canvas → Base → Note
```

A **Note** overrides everything else. If there is no note, the [[Bases|Base]] takes over. If there is no base, the [[Obsidian Canvas|Canvas]] wins. The [[Folders|Folder]] only displays if nothing else overrides it.

To avoid overrides entirely, give each item a unique name.

### Why use file overrides?

This hierarchy lets you replace a standard auto-generated folder index page with richer content:

- **Canvases:** Replace a plain file list with a visual, interactive map of the folder's contents.
- **Bases:** Create a structured database view that filters and sorts the notes inside a folder.
- **Notes:** Write a hand-crafted landing page with keyword-rich content and internal links for better SEO — nothing beats a well-written note for discoverability.

## Folder page access with pretty URLs

When using the default pretty URL mode (without `--flat-urls`), overridden folder pages are still generated — they just aren't linked from anywhere in the navigation. The folder's `/index.html` remains on disk alongside the overriding file.

In `--flat-urls` mode, the override is absolute: the folder's `index.html` is replaced entirely by the overriding file. See the [Generate Command](../Commands/generate.md) for details on the `--flat-urls` flag.

**Can you still reach an overridden folder page?** In pretty URL mode, yes — by manually entering the URL with a trailing `/` in your browser's address bar. In flat URL mode, no — the folder page does not exist.
