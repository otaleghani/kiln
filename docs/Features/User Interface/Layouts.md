---
title: Layouts
description: Customize the interface of your generated site. Choose between the standard Obsidian-like 'Default' layout or the focused 'Simple' layout using the --layout CLI flag.
---
# Layouts

Kiln allows you to control the "shell" of your website—the navigation, sidebars, and tools that surround your content. While the content stays the same, the user interface can be adapted to fit your aesthetic preferences or the complexity of your vault.

You can select a layout during the build process using the `--layout` flag.

## Usage
To switch layouts, simply pass the name of the layout to the `generate` command.
```bash
# Use the minimalist "simple" layout
kiln generate --input ./vault --output ./dist --layout simple

# If no flag is provided, Kiln defaults to "default"
kiln generate --input ./vault --output ./dist
```

## Available Layouts
Kiln currently ships with these core layouts.

|**Layout Name**|**Description**|**Best For**|
|---|---|---|
|**`default`**|A faithful reproduction of the Obsidian client interface. Includes persistent left and right sidebars for navigation and context.|Documentation, Wikis, and Digital Gardens where context and quick navigation are key.|
|**`simple`**|A minimalist, distraction-free interface. Sidebars are removed in favor of floating action buttons to access tools on demand.|Blogs, Portfolios, or simple reading experiences where the content should take center stage.|
## Need more control?
If just changing the layout is not enough, you can take a look at [[What is Custom Mode|custom mode]], which allows you to use Obsidian as an headless CMS while you design different layouts for different content collections.

## Layout Details

### Default (Obsidian Parity)
The **Default** layout is designed to feel exactly like Obsidian Publish or the Obsidian desktop app. It maximizes information density.

- **Left Sidebar:** Contains the **[[Explorer|File Explorer]]** and Search.
- **Right Sidebar:** Contains the **[[Table of Contents]]** and **[[Local Graph]]**.
- **Behavior:** Sidebars are collapsible but intended to be open to provide context.

### Simple (Minimalist)
The **Simple** layout removes the "app-like" chrome to focus purely on the text.

- **No Fixed Sidebars:** The screen is dedicated to your content.
- **Action Buttons:** Tools like the Graph, Table of Contents, and Search are tucked away behind floating buttons or a simplified header. They appear as overlays or modals when needed, rather than taking up permanent screen real estate.