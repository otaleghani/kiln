---
title: Fonts & Typography
description: Kiln includes a selection of high-quality, self-hosted Google Fonts that are baked directly into your site for maximum performance and privacy.
---

# Fonts & Typography

Typography is the foundation of a readable site. Kiln allows you to choose from a curated selection of typefaces to match your site's aesthetic.

Crucially, Kiln uses a **Zero-Dependency** approach for fonts.

* **Self-Hosted:** When you choose a font, Kiln embeds the font files directly into your website's output.
* **Privacy Focused:** Your visitors never make requests to third-party servers (like Google Fonts), ensuring better privacy (GDPR compliance) and speed.
* **Offline Ready:** Your site looks perfect even if the user is offline.

## Available Fonts

You can select one of the following typeface options during the build process:

| ID | Name | Style | Best For |
| :--- | :--- | :--- | :--- |
| `system` | **System UI** | Sans-serif | Performance purists. Uses the native font of the user's OS (San Francisco on Mac, Segoe UI on Windows). |
| `inter` | **Inter** | Sans-serif | Clean, modern, highly legible. The standard for modern site. |
| `lato` | **Lato** | Sans-serif | Friendly and stable. Good for blogs and softer aesthetics. |
| `merriweather` | **Merriweather** | Serif | Excellent for long-form reading and dense text. |

## Configuration

To apply a font, you simply pass its **ID** (from the table above) to the `--font` flag when running the `generate` command.

**Example: Using Inter**
```bash
./kiln generate --font "inter"
```

**Example: Using System Fonts (Default)** If you do not specify a flag, Kiln defaults to `system` for the fastest possible load times.
```bash
./kiln generate
```