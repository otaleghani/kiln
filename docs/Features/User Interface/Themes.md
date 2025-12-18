---
title: Themes
description: Custom description
---
# Themes

Kiln comes baked in with some different themes. More to come at a later day. Right now it has:

- Default (Obsidian colors)
- Cattpuccin
- Nord
- Dracula

> [!note]
> Custom themes will be available soon enough.

## How to use a theme
Kiln bakes into the CSS of your website the chosen theme, so you'll need to add the `--theme` flag when you use the `generate` command, like this:
``` bash
kiln generate --theme "catppuccin"
```