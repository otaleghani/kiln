---
title: "Folder Navigation & Auto-Generated Index Pages"
description: Kiln mirrors your Obsidian vault folder structure as clean URLs, auto-generates index pages, and adds folders to your knowledge graph.
---
# Folders

Kiln follows a "Parity First" philosophy with your directory structure. The folder hierarchy you create in your Obsidian vault is exactly what gets built for the web, complete with auto-generated index pages, clean URLs, and [graph visualization](../User Interface/Global Graph.md).

![[folders.png]]

## Structure Mirroring

Kiln preserves your file organization 1:1. If you have a folder path in your vault like `Recipes/Desserts/Cakes`, Kiln generates the corresponding URL structure (e.g., `yoursite.com/recipes/desserts/cakes`).

Folder names are **slugified** for web-friendly URLs: spaces become dashes and letters are lowercased. For example, a folder called `My Cool Project` becomes `my-cool-project` in the URL.

## Auto-Generated Index Pages

When a user navigates to a folder, Kiln automatically generates an index page for that directory. This page displays:

- **Subfolders** listed first, each linking to their own index page.
- **Files** listed below, showing the note title and last modified date.

Empty folders — those with no files and no subfolders — are skipped entirely and do not generate a page.

### Replacing a Folder Page with Custom Content

You can override the auto-generated folder page by creating a file with the same name as the folder. For example, placing a `Recipes.md` note next to a `Recipes/` folder means the `/recipes` URL displays your note instead of the auto-generated listing.

The override follows a strict priority hierarchy:

```
Folder → Canvas → Base → Note
```

A note always wins. See [Kiln quirks](../Quirks.md) for full details on how overrides interact with links.

You can also place an `index.md` file inside a folder to replace its auto-generated page. An `index.md` at the vault root becomes your site homepage.

## Hiding Folders from Output

Folders (and files) that start with a `.` are automatically excluded from the generated site. To hide specific folders without renaming them to dotfiles, add the `_hidden_` prefix to the folder name.

For example, renaming your template folder to `_hidden_Templates` prevents Kiln from processing it. See [hidden files and folders](../Rendering/Hidden files folders.md) for more details.

## Folder Sorting in the Sidebar

The [Explorer](./Explorer.md) sidebar lists folders before files at every level, with alphabetical sorting (A–Z, case-insensitive) within each group. To control the display order, prefix folder names with numbers:

```text
📂 01-Guide
  📄 01-Getting-Started.md
  📄 02-Installation.md
📂 02-Advanced
📄 About.md
```

## Graph Integration

Folders are not just containers; they are part of your knowledge network. Each folder appears as a node in both the [[Global Graph]] and [[Local Graph]], with links connecting it to every file and subfolder it contains. This lets you visualize how clusters of knowledge are structurally grouped.

## Sitemap and SEO

Every folder page is automatically included in your site's [sitemap.xml](../SEO/Sitemap xml.md) with its last-modified date, ensuring search engines can discover and index your full site structure.
