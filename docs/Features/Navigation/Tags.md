---
title: "Tags — Organize and Browse Notes by Topic"
description: "Use Obsidian tags in Kiln to auto-generate topic pages, create cross-folder collections, and visualize note connections in the graph."
---

# Tags

Tags let you group related notes across your entire vault regardless of folder structure. Kiln automatically detects every tag in your Obsidian vault and generates a dedicated page for each one, listing all notes that share that tag and connecting them in your site's [graph visualization](../User Interface/Global Graph.md).

## How Kiln Detects Tags

Kiln picks up tags from two sources in each markdown file:

- **Inline hashtags** written anywhere in the note body, such as `#philosophy` or `#project-alpha`.
- **YAML frontmatter** tags defined in the `tags` field at the top of a file:

```yaml
---
tags:
  - philosophy
  - greek
---
```

Both formats are combined. A note with `#philosophy` in the body and `philosophy` in the frontmatter still counts as one tag. Tag names support letters, numbers, hyphens, and underscores (for example `#my-tag` or `#project_v2`).

## Auto-Generated Tag Pages

For every unique tag found in your vault, Kiln creates a page at `/tags/<tagname>`. These tag pages act as dynamic hubs that list every note containing that tag, sorted with the note name and last-modified date.

- **Cross-folder discovery:** A tag page collects notes from any directory. Clicking `#urgent` shows tasks from both your `Work/` and `Personal/` folders in one place.
- **Clickable inline tags:** When Kiln renders a note, every `#tag` in the body becomes a link pointing to its tag page, so readers can jump straight to related content.

Tag page URLs respect the global URL structure setting. With flat URLs enabled, the path is `/tags/tagname/`; otherwise it is `/tags/tagname`.

## Graph Integration

Tags are full participants in the [local graph](../User Interface/Local Graph.md) and [global graph](../User Interface/Global Graph.md). Each tag appears as its own node, with edges linking it to every note that uses it. This makes hidden connections visible — you might discover that `#productivity` bridges your `Journal` folder to your `Reading List`, revealing patterns that folders alone would miss.

## Practical Example

Suppose your vault contains three notes:

| Note | Location | Tags |
|------|----------|------|
| Stoicism.md | `Philosophy/` | `#philosophy` `#ancient` |
| Marcus Aurelius.md | `People/` | `#philosophy` `#stoicism` |
| Daily Reflection.md | `Journal/` | `#stoicism` `#journaling` |

After running the [Generate Command](../../Commands/generate.md), Kiln produces tag pages for `#philosophy`, `#ancient`, `#stoicism`, and `#journaling`. The `#philosophy` page lists both *Stoicism* and *Marcus Aurelius* even though they live in different folders. In the graph, the `#stoicism` node connects *Marcus Aurelius* to *Daily Reflection*, surfacing a relationship that the folder tree cannot show.

## Tips for Effective Tagging

- **Keep tag names consistent.** `#project-alpha` and `#projectAlpha` are treated as separate tags.
- **Combine tags with [folders](./Folders.md)** for both rigid structure and flexible associations.
- **Use the [Explorer](./Explorer.md)** sidebar to navigate your folder tree, then rely on tags for cross-cutting topics.
- **Check the [global graph](../User Interface/Global Graph.md)** after publishing to spot unexpected connections between tagged notes.
