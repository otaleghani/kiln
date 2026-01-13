---
title: Folders
description: Learn how Kiln mirrors your Obsidian folder structure. Folders automatically generate index pages, listing contents and sub-contents, and act as structural nodes in your knowledge graph.
---
# Folders

Kiln follows a "Parity First" philosophy with your directory structure. The hierarchy you create in your Obsidian vault is exactly what gets built for the web.

![[folders.png]]

## Structure Mirroring
Kiln preserves your file organization 1:1. If you have a folder path in your vault like `Recipes/Desserts/Cakes`, Kiln generates the corresponding URL structure (e.g., `yoursite.com/recipes/desserts/cakes`). This ensures your "physical" organization remains intact online.

## Content Listing
When a user navigates to a folder—either via the File Explorer or a link—Kiln generates an index page for that directory.

- **Direct Contents:** It lists all notes found directly in that folder.
- **Recursive Contents:** It also displays subfolders for easy navigation.

## Graph Integration
Folders are not just containers; they are part of your knowledge network. Folders appear as nodes in both the **[[Global Graph]]** and **[[Local Graph]]**. This allows you to visualize how clusters of knowledge are structurally grouped and how different directories relate to one another.