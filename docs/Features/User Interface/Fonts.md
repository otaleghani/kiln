---
title: Fonts & Typography
description: Kiln includes a selection of high-quality, self-hosted Google Fonts that are baked directly into your site for maximum performance and privacy.
tags:
---
#user-interface
# Fonts & Typography

Typography is the foundation of a readable site. Kiln allows you to choose from a curated selection of typefaces to match your site's aesthetic.

Crucially, Kiln uses a **Zero-Dependency** approach for fonts.

* **Self-Hosted:** When you choose a font, Kiln embeds the font files directly into your website's output.
* **Privacy Focused:** Your visitors never make requests to third-party servers (like Google Fonts), ensuring better privacy (GDPR compliance) and speed.
* **Offline Ready:** Your site looks perfect even if the user is offline.

## Available Fonts

You can select one of the following typeface options during the build process:

| **ID**              | **Name**              | **Style**  | **Best For**                                                                                            |
| ------------------- | --------------------- | ---------- | ------------------------------------------------------------------------------------------------------- |
| `system`            | **System UI**         | Sans-serif | Performance purists. Uses the native font of the user's OS (San Francisco on Mac, Segoe UI on Windows). |
| `inter`             | **Inter**             | Sans-serif | Clean, modern, highly legible. The standard for modern sites.                                           |
| `lato`              | **Lato**              | Sans-serif | Friendly and stable. Good for blogs and softer aesthetics.                                              |
| `merriweather`      | **Merriweather**      | Serif      | Excellent for long-form reading and dense text.                                                         |
| `lora`              | **Lora**              | Serif      | Contemporary with calligraphic roots. Elegant and perfect for storytelling or essays.                   |
| `libre-baskerville` | **Libre Baskerville** | Serif      | Traditional and open. Specifically optimized for reading body text on screens.                          |
| `noto-serif`        | **Noto Serif**        | Serif      | A universal, reliable serif. Clean, balanced, and highly compatible.                                    |
| `ibm-plex-sans`     | **IBM Plex Sans**     | Sans-serif | Neutral and technical. Excellent for technical documentation and data-heavy sites.                      |
| `google-sans`       | **Google Sans**       | Sans-serif | Geometric and simple. The distinctive, clean look used across Google products.                          |
| `roboto`            | **Roboto**            | Sans-serif | Mechanical skeleton with friendly curves. A widely used, neutral standard.                              |

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