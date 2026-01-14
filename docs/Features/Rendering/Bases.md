---
title: Bases
description: Explore Kiln's experimental Bases feature. Transform your Obsidian data views into interactive HTML tables, cards, and lists with filtering, sorting, and grouping capabilities.
---

> [!attention] **Beta** 
> The Bases feature is currently in active development. While functional, you may encounter edge cases or visual inconsistencies as we refine the rendering engine.

# Bases

[**Bases**](https://help.obsidian.md/bases) bring the power of databases to your static site. Kiln automatically detects and reproduces the "Base" views you have configured in Obsidian, transforming static lists into interactive, data-rich dashboards.

When Kiln encounters a Base definition, it mimics the **first view** configured in Obsidian and generates a high-fidelity HTML reproduction. This allows your site visitors to interact with your data just as you do inside your vault.

## Supported views
Kiln currently supports the following visualization layouts:

### Table View
A classic spreadsheet-like view perfect for dense data and comparing properties.
![[bases_view_table.png]]

### Cards View
A Kanban-style or gallery view that emphasizes visual content and summaries.
![[bases_view_cards.png]]

### List View
A clean, vertical list for simple collections.
![[bases_view_lists.png]]

## Supported functions
Bases are a complex obsidian feature that allow you to manipulate the view by filtering, grouping and sorting your notes and files. Kiln right now supports the following functions:

- **Global Filters**: Filters applied to all views of a base.
- **View-specific Filters**: Filters applied only to one specific view.
- **Group by**: Cluster items based on shared values.
- **Sorting**: Sort data based on a specific field. This is partially supported at the moment.

## Quirks

- The filter `file, links to` uses a simple contain filter and not relative paths like Obsidian. This means that if you have multiple notes with the same name they will all be picked up. This poses no problems if you have unique names.
- Links filters use case-insensitive equality. This is something that regular notes support but not bases.
- `file.links`, `file.embeds`, `file.backlinks` have some filters that I'm still trying to figure out, like `is exactly`, `is not exactly`
- Checkout the [[Quirks#Same filename, multiple extensions]] for more information about files overrides. 