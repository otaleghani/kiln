---
title: "Search & Quick Find — Filter Notes in the Sidebar"
description: "Use Kiln's real-time Quick Find search to instantly filter notes and folders in the sidebar by name. Client-side, zero latency, with smart folder expansion."
---
# Search & Quick Find

Kiln includes a built-in **Quick Find** search bar that lets you instantly filter notes and folders in the [Explorer](./Explorer.md) sidebar. As you type, the file tree updates in real time to show only matching items — no server round-trips, no waiting.

Quick Find is available in the **default** and **legacy** layouts. The [simple layout](../User Interface/Layouts.md) does not include a search bar.

## How to Use Quick Find

The search input sits at the top of the left sidebar, directly above your file tree. Start typing a word or phrase and the Explorer immediately hides every item that does not match.

1. **Click the search field** (or tab into it) at the top of the sidebar.
2. **Type your query** — results filter with every keystroke.
3. **Click a result** to navigate to that page.
4. **Clear the field** to restore the full file tree.

Because navigation in Kiln is powered by [client-side page swaps](./Client Side Navigation.md), the search field retains focus even after you navigate to a new page, so you can continue refining your query without clicking back into the input.

## How the Search Works

Quick Find runs entirely in the **browser**. It performs a case-insensitive substring match against every item visible in the sidebar — notes, folders, canvas files, and base files. There is no server request and no external search index, which means **zero latency**: results appear the instant you press a key.

### Smart Folder Expansion

When a matching file is buried inside collapsed directories, Kiln automatically opens every parent folder along the path so the result is visible. Non-matching items are hidden to keep the list focused.

- **Filtering:** Files and folders whose names do not contain the search term are hidden.
- **Auto-Expand:** Parent directories of every match are opened automatically, no matter how deeply nested the file is.
- **Auto-Collapse on Clear:** Clearing the search field restores the sidebar to its previous state.

### What Quick Find Searches

Quick Find currently matches against **file and folder names** (the titles displayed in the sidebar). It does not yet search note content, tags, or frontmatter fields. Full-text content search is planned for a future release — see the [Roadmap](../../Roadmap.md) for upcoming features.

## Practical Examples

| You type | What appears |
|----------|-------------|
| `recipe` | Every note and folder whose name contains "recipe" — e.g., `Pasta Recipes`, `Recipe Index`, `recipes/` |
| `2024` | Any file or folder with "2024" in its name, useful for date-based vault organization |
| `setup` | Notes like `Dev Setup.md`, `Server Setup Guide.md`, and the `setup/` folder |

Matching is **case-insensitive**, so typing `setup`, `Setup`, or `SETUP` all return the same results.

## Tips for Navigating Large Vaults

- **Use specific terms.** In a vault with hundreds of notes, a short query like `a` matches almost everything. Start with two or three characters for faster filtering.
- **Combine with [tags](./Tags.md).** If Quick Find returns too many results, open a [tag page](./Tags.md) to browse notes by topic instead of by name.
- **Prefix file names with numbers** to control sort order in the [Explorer](./Explorer.md). Searching for the prefix (e.g., `01-`) lets you jump straight to a section.
- **Leverage [folders](./Folders.md)** for broad categories and Quick Find for pinpointing a specific note within them.

## Limitations

- **Name-only search** — file content, tags, and frontmatter are not searched.
- **Substring matching** — there is no fuzzy matching or typo tolerance. The query must appear exactly as a substring of the item name.
- **No keyboard shortcut** — the search field must be clicked or tabbed into; there is no global hotkey to focus it.
