---
title: Fonts
description: Custom description
---
# Fonts

Kiln comes with some Google Fonts embedded. More will be embedded soon enough. Right now we have:

- Default (uses system fonts)
- Merriweather
- Inter
- Lato

## How to use a font
Kiln bakes into the CSS of your website the chosen font, so you'll need to add the `--font` flag when you use the `generate` command, like this:
``` bash
kiln generate --font "inter"
```