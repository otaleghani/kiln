---
title: "Sidebar File Explorer — Browse Your Site Structure"
description: "Kiln generates a hierarchical file explorer in the sidebar that mirrors your vault's folder structure with collapsible directories, active page highlighting, and live search filtering."
---

# Sidebar File Explorer

Kiln generates a hierarchical **File Explorer** in the left sidebar that serves as the primary navigation menu for your site. The explorer mirrors your local file structure — your Obsidian vault or directory — so if your content is organized logically on your computer, it appears the same way on your website.

The explorer is available in the **default** and **legacy** [layouts](../User Interface/Layouts.md). The simple layout does not include a sidebar file tree.

## Sorting and Display Order

The explorer sorts your content with two rules applied at every level of the tree:

1. **Folders first.** Directories always appear above individual files.
2. **Alphabetical order.** Within each group (folders or files), items are sorted A–Z using a case-insensitive comparison.

### Controlling the Order with Numeric Prefixes

Because the explorer relies on filenames for sorting, you can control the order of your pages by adding numeric prefixes to file or folder names.

For example, if you want a specific section to appear at the top of the sidebar:

```text
📂 01-Guide
  📄 01-Getting-Started.md
  📄 02-Installation.md
📂 02-Advanced
📄 About.md
📄 Changelog.md
```

The prefixes determine the position. Without them, items fall into standard alphabetical order. This technique works at any depth of the tree — you can prefix both folders and individual notes.

## Collapsible Directories

Folders in the explorer are collapsible. Each folder renders as a toggleable section with a chevron icon that rotates when the folder opens. This keeps the sidebar clean, especially for vaults with deep nesting.

Kiln manages the open/closed state of folders automatically in two situations:

- **Auto-expand on navigation.** When you navigate to a page nested deep within folders, the explorer expands every parent directory so the active page is always visible.
- **Auto-expand during search.** Typing in the [Quick Find search bar](./Search.md) expands all folders that contain matching results, then collapses them again when you clear the search.

You can also open and close any folder manually by clicking its name or chevron.

## Active Page Highlighting

The explorer highlights the page you are currently viewing by applying your theme's accent color to the corresponding sidebar link. This runs entirely in the browser, so the highlight updates instantly — even during [client-side navigation](./Client Side Navigation.md) where the page does not fully reload.

When the active page sits inside nested folders, the explorer opens each parent folder automatically so the highlighted link is never hidden behind a collapsed directory.

## Sidebar Toggle

On **desktop** screens (1280 px and wider), the left sidebar is visible by default. You can collapse and expand it using the toggle button in the header. Kiln saves your preference in the browser's `localStorage`, so the sidebar stays in the state you chose across page loads and sessions.

On **mobile** screens, the sidebar starts hidden and slides in as an overlay when you tap the toggle. Clicking any link inside the sidebar closes it automatically so the content area is not blocked.

## File Type Indicators

When a folder has an associated note with the same name — for example, a `Recipes.md` file next to a `Recipes/` folder — the explorer displays a small badge on the folder entry indicating the file type (Note, Canvas, or Base). This tells you that clicking the folder link opens your custom content rather than the auto-generated [folder index page](./Folders.md).

## Hidden Files and Folders

Files and folders prefixed with a dot (`.`) or with `_hidden_` are excluded from the explorer entirely. Use this to keep template folders, draft directories, or other internal content out of your published site navigation. See [hidden files and folders](../Rendering/Hidden files folders.md) for the full details.

## Tips for Organizing a Large Vault

- **Use numeric prefixes** on your most important folders so they sort to the top of the sidebar instead of falling into alphabetical order.
- **Combine the explorer with [tags](./Tags.md)** for two complementary navigation paths — folders for structure, tags for cross-cutting topics.
- **Use [Quick Find](./Search.md)** to jump to a specific note without scrolling through the tree. The search field sits directly above the explorer.
- **Keep your folder depth shallow.** Deeply nested structures still work, but two or three levels are easier for readers to scan in the sidebar.
