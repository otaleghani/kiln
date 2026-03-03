---
title: "Bases — Obsidian Database Views as HTML Tables, Cards & Lists"
description: "Render Obsidian Bases as interactive HTML tables, cards, and lists with filtering, sorting, and grouping on your static site."
---

> [!attention] **Beta**
> The Bases feature is currently in active development. While functional, you may encounter edge cases or visual inconsistencies as we refine the rendering engine.

# Bases — Database Views on Your Static Site

[**Bases**](https://help.obsidian.md/bases) bring the power of Obsidian's database views to your static site. Kiln automatically detects `.base` files in your vault and renders the **first configured view** as a high-fidelity HTML reproduction — complete with filtering, sorting, and grouping. Your site visitors can browse your data just as you see it in Obsidian, without any manual configuration.

If you're setting up a new site, see the [Installation](../../Installation.md) guide. For rendering other Obsidian-specific content, check out [Obsidian Markdown](./Obsidian Markdown.md) and [Callouts](./Callouts.md).

## Supported View Types

Kiln renders three view layouts from your Base definitions:

### Table View

A spreadsheet-like view for dense data. Columns map to your note properties, making it easy to compare values across many files at once.

![[bases_view_table.png]]

### Cards View

A Kanban-style or gallery layout that emphasizes visual content and summaries. Each card represents one note with its key properties displayed.

![[bases_view_cards.png]]

### List View

A clean vertical list for simple collections where a full table isn't needed.

![[bases_view_lists.png]]

## Filtering, Sorting, and Grouping

Bases let you control which notes appear and how they're organized. Kiln supports the following data manipulation functions:

- **Global Filters** — Applied to all views in a Base. Notes that don't match are excluded from every view.
- **View-specific Filters** — Applied only to one particular view, leaving other views unaffected.
- **Group by** — Cluster items based on shared property values (e.g., group by folder, tag, or a custom frontmatter field).
- **Sorting** — Order results by a specific field. Sorting is partially supported at the moment.

### Supported Filter Operators

Kiln's filter engine supports a wide range of comparison and containment operators:

| Operator | Example |
|---|---|
| `is` / `==` | `status is "done"` |
| `is not` / `!=` | `status is not "draft"` |
| `>`, `>=`, `<`, `<=` | `file.size > 1000` |
| `contains` | `file.name contains "project"` |
| `does not contain` | `tags does not contain "archive"` |
| `contains any of` | `tags contains any of ["book", "article"]` |
| `contains all of` | `tags contains all of ["review", "done"]` |
| `starts with` / `ends with` | `file.name starts with "2024"` |
| `is empty` / `is not empty` | `status is not empty` |
| `on` / `not on` | `file.ctime on "2024-01-15"` |

Conditions can be combined with `and`, `or`, and `not` logic.

### Available File Fields

These built-in fields can be used in filters, sorting, and grouping:

| Field | Description |
|---|---|
| `file.name` | Note filename |
| `file.folder` | Parent folder path |
| `file.path` | Full file path |
| `file.ext` | File extension |
| `file.size` | File size in bytes |
| `file.ctime` / `file.mtime` | Creation / modification date |
| `file.tags` | All tags in the note |
| `file.links` | Outgoing wikilinks |
| `file.embeds` | Embedded files |

Any [frontmatter](./Obsidian Markdown.md) property (like `status`, `rating`, or `category`) can also be used as a field.

### Filter Methods

You can call methods directly on fields using dot notation:

- `file.hasTag("book")` — Check if a note has a specific tag
- `file.hasLink("My Note.md")` — Check if a note links to another note
- `file.inFolder("projects")` — Check if a note lives in a folder
- `file.hasProperty("date")` — Check if a frontmatter property exists
- `file.name.contains("draft")` — Substring check on a string field
- `file.name.startsWith("2024")` / `file.name.endsWith(".md")`
- `file.tags.isEmpty()` — Check if a collection field is empty

## Known Quirks

- The `file, links to` filter uses simple string containment rather than resolved paths. If multiple notes share the same name, all of them will match. This is not an issue if your note names are unique.
- Link-based filters use case-insensitive matching. Regular [Obsidian Markdown](./Obsidian Markdown.md) notes support this, but Bases in Obsidian itself are case-sensitive.
- Some filters on `file.links`, `file.embeds`, and `file.backlinks` (like `is exactly` and `is not exactly`) are not yet implemented.
- See [[Quirks#Same filename, multiple extensions]] for details on how Kiln handles files that share a name but differ in extension.

For other rendering features, explore [Mermaid Graphs](./Mermaid Graphs.md), [Math & LaTeX](./Math.md), and [Syntax Highlighting](./Syntax Highlighting.md).
