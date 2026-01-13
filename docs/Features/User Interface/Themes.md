---
title: Themes & Visuals
description: Customize the look of your site with Kiln's built-in color palettes, including popular themes like Catppuccin, Nord, and Dracula.
---
# Themes

Kiln allows you to completely transform the visual aesthetic of your site to match your personal taste.

Instead of relying on heavy client-side JavaScript to manage styles, Kiln **bakes your chosen theme directly into the CSS** during the build process. This ensures that your site remains lightweight and fast, regardless of which color palette you choose.

## Available Themes

Kiln currently ships with four professionally curated themes popular in the developer and Obsidian communities.

| ID           | Name            | Aesthetic                                                                                                                       | **Example**                 |
| :----------- | :-------------- | :------------------------------------------------------------------------------------------------------------------------------ | --------------------------- |
| `default`    | **Obsidian**    | The classic Obsidian experience. Dark purples and charcoal grays. Ideal if you want your site to look exactly like your editor. | ![[palette_default.png]]    |
| `catppuccin` | **Catppuccin**  | A soft, pastel theme. Low contrast and easy on the eyes, featuring warm grays and vibrant accents.                              | ![[palette_catppuccin.png]] |
| `nord`       | **Nord**        | An arctic, north-bluish color palette. Cool, professional, and clean.                                                           | ![[palette_nord.png]]       |
| `dracula`    | **Dracula**     | A famous dark theme with high contrast and vibrant pink/green accents.                                                          | ![[palette_dracula.png]]    |
| `tokyonight` | **Tokyo Night** | A clean visual style celebrating the lights of downtown Tokyo. Deep blues, vibrant purples, and high clarity.                   | ![[palette_tokyonight.png]] |
| `rose-pine`  | **Rosé Pine**   | A classy, warm palette with muted pinks and organic hues. Minimalist and comfortably low-contrast.                              | ![[palette_rose-pine.png]]  |
| `gruvbox`    | **Gruvbox**     | A retro groove color scheme with reddish-browns and warm earth tones. Classic, cozy, and distinctive.                           | ![[palette_gruvbox.png]]    |
| `everforest` | **Everforest**  | A nature-inspired theme based on green and warm earthy tones. Softer than Gruvbox; organic and pleasant.                        | ![[palette_everforest.png]] |
| `cyberdream` | **Cyberdream**  | A high-contrast, futuristic theme. Deep blacks and vivid neons for a sharp cyberpunk aesthetic.                                 | ![[palette_cyberdream.png]] |

## Configuration

To apply a theme, use the `--theme` flag followed by the theme's **ID** when running the [[generate]] command.

**Example: Using Catppuccin**
```bash
./kiln generate --theme "catppuccin"
```

**Example: Using Default (Obsidian)** If you do not specify a flag, Kiln defaults to the standard Obsidian look.
```bash
./kiln generate
```