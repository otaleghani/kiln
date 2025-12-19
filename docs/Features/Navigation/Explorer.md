---
title: Explorer
description: Kiln automatically generates a hierarchical file explorer in the sidebar, mirroring your local folder structure for intuitive content navigation.
---

# Explorer

Kiln generates a hierarchical **File Explorer** in the left sidebar, serving as the primary navigation menu for your site.

This explorer is designed to be a direct reflection of your local file structure (your Obsidian Vault or directory). This "What You See Is What You Get" approach ensures that if your content is organized logically on your computer, it will be organized logically on your website.

## Ordering Logic

The Explorer automatically sorts your content to maintain consistency and predictability. The sorting algorithm follows two strict rules:

1.  **Type Priority:** Folders are always listed **before** individual files.
2.  **Alphabetical Order:** Within each group (folders or files), items are sorted alphabetically (A-Z) by their filename.

### Customizing the Order
Because the Explorer relies on filenames for sorting, you can control the order of your pages by using numeric prefixes in your file or folder names.

**Example:**
If you want a specific page to appear at the top, you might name it `01-Introduction.md`.

```text
📂 01-Guide
  📄 01-Getting-Started.md
  📄 02-Installation.md
📂 02-Advanced
📄 About.md
📄 Changelog.md
```

Future versions of Kiln will give you the ability to customize the Explorer further.

## Features

### Collapsible Directories

To keep the navigation clean, folders in the Explorer are collapsible. The state of these folders (open or closed) is managed automatically:

- **Auto-Expand:** When you navigate to a page nested deep within folders, Kiln automatically expands the directory tree to show the active page's location.
- **Search Interaction:** Typing in the [[Search|sidebar search bar]] will automatically expand all folders containing matching results.

### Active State

The Explorer highlights the currently viewed page, providing users with immediate visual context regarding their location within the site hierarchy.