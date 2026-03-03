---
title: Wikilinks & Embeds — Link Notes, Embed Content, and Resize Images
description: Use Obsidian-style wikilinks to connect notes, embed pages and sections, link to headers and block IDs, and resize images in your Kiln site.
---

# Wikilinks & Embeds

Kiln fully supports **wikilinks**, the double-bracket linking syntax used by Obsidian. Wikilinks let you connect notes, embed content from other pages, and display images — all with a concise, readable notation that replaces verbose Markdown link paths.

Every wikilink you create is automatically resolved to the correct URL and tracked for the [Local Graph](../User Interface/Local Graph.md) and [Global Graph](../User Interface/Global Graph.md), so your site's interactive graph visualization reflects the connections between your notes.

## Linking Notes

Wrap a filename in double square brackets to link to another note. Kiln resolves the path automatically — you do not need to include the `.md` extension or the full folder path.

| Syntax               | Description                                                                        | Result             |
| :------------------- | :--------------------------------------------------------------------------------- | :----------------- |
| `[[Index]]`          | **Standard Link** — Links to a file named `Index.md`.                              | [[Index]]          |
| `[[Index\|Home]]`    | **Aliased Link** — Links to `Index.md` but displays "Home" as the clickable text.  | [[Index\|Home]]    |
| `[[Index#Features]]` | **Header Link** — Links directly to the "Features" heading inside `Index.md`.      | [[Index#Features]] |

When multiple files share the same name, Kiln picks the best match: root-level files take priority, then the file with the shortest path. You can disambiguate by including a partial path, such as `[[Deployment/Vercel]]`.

## Embedding Notes and Sections

Add an exclamation mark before a wikilink to embed content from another page directly into the current one. This renders the target page (or a specific section) inline, with a header linking back to the original.

* `![[My Note]]` — Embeds the entire content of `My Note.md`.
* `![[My Note#Introduction]]` — Embeds only the "Introduction" heading and its content.
* `![[My Note#^ref123]]` — Embeds the specific block marked with `^ref123`.

Block IDs work by appending `^identifier` to any paragraph or list item in your source note. For example, writing `This is important. ^key-point` creates a referenceable block you can embed elsewhere with `![[My Note#^key-point]]`.

Embedded content appears in a styled container with the section title and a link to open the full original page.

## Embedding Images

Embed images into your pages by adding `!` before an image wikilink. This tells Kiln to render the image inline rather than creating a text link.

### Syntax

`![[photo.png]]`

Supported formats: `.png`, `.jpg`, `.jpeg`, `.gif`, `.svg`, and `.webp`.

### Resize Options

Control the display size of embedded images using the pipe `|` syntax:

* `![[image.png|300]]` — Sets the image width to **300 pixels**.
* `![[image.png|100x100]]` — Sets the image to exactly **100x100 pixels**.

The text after the pipe also serves as the image's alt text when no resize value is detected, which helps with accessibility and [SEO](../SEO/Meta Tags.md).

## How Path Resolution Works

Kiln resolves wikilink targets using a case-insensitive lookup against all files in your vault. The resolution follows these rules:

1. **Exact name match** — A file whose name matches the link target exactly is preferred.
2. **Root-level priority** — If multiple files share the same name, a file at the vault root wins.
3. **Shortest path** — Among non-root candidates, the file with the shortest path is selected.
4. **Partial path matching** — Including folder segments (e.g., `[[folder/note]]`) narrows the match to files containing that path.

This mirrors Obsidian's own resolution behavior, so links that work in Obsidian will work on your generated site. If a link target cannot be found, Kiln renders it as plain text. You can use the [Doctor Command](../../Commands/doctor.md) to scan your vault for broken wikilinks before publishing.

## Graph Integration

Every resolved wikilink creates a directed edge in Kiln's graph data. These connections appear in both the [Local Graph](../User Interface/Local Graph.md) (showing a page's direct neighbors) and the [Global Graph](../User Interface/Global Graph.md) (showing all vault connections). Nodes with more wikilinks pointing to them appear larger, making heavily referenced pages easy to spot. The graph also includes connections from [Folders](./Folders.md) and [Tags](./Tags.md).
