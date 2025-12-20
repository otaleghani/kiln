---
title: Wikilinks & Embeds
description: Learn how to use standard Wikilink syntax to connect notes, create aliases, link to specific headers, and embed images directly into your content.
---

# Wikilinks

Kiln fully supports **Wikilinks**, the standard linking syntax used by tools like Obsidian. This allows you to connect your notes and embed media using a simple, bracket-based notation that is easier to read and write than standard Markdown links.

## Linking Notes

To link to another note in your vault, simply wrap the filename in double square brackets. Kiln will automatically resolve the path to the correct file.

| Syntax               | Description                                                                        | Result             |
| :------------------- | :--------------------------------------------------------------------------------- | :----------------- |
| `[[Index]]`          | **Standard Link**<br>Links to a file named `Index.md`.                             | [[Home]]          |
| `[[Index\|Home]]`    | **Aliased Link**<br>Links to `Index.md` but displays "Home" as the clickable text. | [[Home\|Home]]    |
| `[[Index#Features]]` | **Anchor Link**<br>Links directly to a specific header section inside the note.    | [[Home#Features]] |

> **Note:** You do not need to include the `.md` extension inside the brackets.

## Embedding Media

You can embed images directly into your pages by adding an exclamation mark `!` before the link. This tells Kiln to render the file instead of just linking to it.

### Syntax
`![[Image.png]]`

This syntax supports standard image formats including `.png`, `.jpg`, `.jpeg`, `.gif`, `.svg`, and `.webp`.

### Resize Options
You can control the display size of your embedded images using the pipe `|` syntax:

* `![[image.png|300]]`: Resizes the image to a width of **300 pixels**.
* `![[image.png|100x100]]`: Resizes the image to exactly **100x100 pixels**.