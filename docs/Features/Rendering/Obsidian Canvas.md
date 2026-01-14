---
title: Canvas
description: Kiln natively renders Obsidian .canvas files as interactive diagrams, allowing you to visualize complex relationships directly on your static site.
---
# Obsidian Canvas

Kiln provides out-of-the-box support for **Obsidian Canvas**. It automatically detects `.canvas` files in your vault and renders them as fully interactive diagrams on your website.

This allows you to preserve your mind maps, flowcharts, and organizational diagrams without needing to export them as static images. Users can pan and zoom around the canvas just as they would in Obsidian.

## Supported Elements

Kiln currently supports the core building blocks of a canvas:

* **Text Cards:** Markdown-supported text blocks.
* **Files:** Embedded notes and images from your vault.
* **Edges:** Connections and arrows between nodes.
* **Groups:** Visual groupings for organizing nodes.

## Limitations & Quirks

The Canvas implementation is currently in **Beta**. While most standard diagrams render correctly, there are known limitations regarding advanced features:

### External Embeds (iFrames)
If your canvas includes a "Link Node" (e.g., a card displaying `https://google.com`), Kiln will generate the correct `<iframe>` tag. However, many modern websites (including Google, GitHub, and Twitter) set `X-Frame-Options: DENY` headers, which strictly prevent browsers from loading them inside an iframe for security reasons. These nodes may appear empty or show a "refused to connect" error.

### Custom Colors
Currently, Kiln doesn't support custom colors, but only the default ones.

*We are actively working on expanding support for these styling features in future updates.*

## Quirks
Checkout the [[Quirks#Same filename, multiple extensions]] for more information about files overrides. 