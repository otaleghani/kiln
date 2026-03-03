---
title: Canonical Tags — Prevent Duplicate Content in Search Engines
description: Add canonical tags to your Kiln site to prevent duplicate content penalties. Learn how the --url flag generates rel=canonical links automatically.
---
# Canonical Tags

A canonical tag is an HTML `<link rel="canonical">` element that tells search engines which version of a URL is the "official" one. Kiln automatically adds a canonical tag to every page on your site when you provide a base URL during the build.

## Why canonical tags matter

Search engines treat each distinct URL as a separate page. Without canonical tags, your content may appear duplicated under multiple URLs:

- `http://example.com/page` vs `https://example.com/page`
- `www.example.com/page` vs `example.com/page`
- `example.com/page` vs `example.com/page/` (trailing slash)

When Google detects these duplicates, it splits ranking signals between them, which can lower your position in search results. A canonical tag consolidates all variations into one authoritative URL.

## Configuration

To enable canonical tags, pass the `--url` flag to the [Generate Command](../../Commands/generate.md) with your site's final public URL:

```bash
kiln generate --url "https://www.example.com"
```

The URL must be the root of your live site, without a trailing slash. Kiln combines this base URL with each page's path to produce the full canonical URL.

> [!warning] No URL, no canonical tag
> If you omit the `--url` flag, Kiln skips canonical tag generation entirely. The same flag is also required for [Sitemap.xml](./Sitemap xml.md) and [Robots.txt](./Robots txt.md) generation.

## How Kiln generates canonical URLs

Kiln produces a different canonical URL depending on the page type and your URL style setting.

### Notes and markdown pages

For a note located at `vault/guides/setup.md` with the base URL `https://example.com`:

```html
<link rel="canonical" href="https://example.com/guides/setup" />
```

### Folder pages

Folder index pages always receive a trailing slash in their canonical URL:

```html
<link rel="canonical" href="https://example.com/guides/" />
```

### Tag pages

Tag pages follow the same logic as notes. A tag page for `#javascript` produces:

```html
<link rel="canonical" href="https://example.com/tags/javascript" />
```

### Flat URLs mode

If you build your site with the `--flat-urls` flag, Kiln appends a trailing slash to note and tag canonical URLs to match the flat file output format:

```bash
kiln generate --url "https://example.com" --flat-urls
```

```html
<!-- With --flat-urls enabled -->
<link rel="canonical" href="https://example.com/guides/setup/" />
```

## Base path support

If your site lives under a subpath (for example, on GitHub Pages at `https://user.github.io/my-notes/`), include that path in the `--url` flag:

```bash
kiln generate --url "https://user.github.io/my-notes"
```

Kiln extracts the path prefix from the URL and prepends it to every page path, so canonical URLs correctly reflect your deployment structure. See [Deploy on GitHub Pages](../../Deployment/GitHub Pages.md) for a full deployment walkthrough.

## Complete production build example

A typical production command that enables all SEO features — canonical tags, [Meta Tags & SEO](./Meta Tags.md), sitemap, and robots.txt — looks like this:

```bash
kiln generate \
  --url "https://notes.example.com" \
  --name "My Digital Garden" \
  --theme "nord" \
  --font "inter"
```

## Verifying your canonical tags

After building your site, open any generated HTML file and look for the `<link rel="canonical">` tag inside the `<head>` section. Confirm that:

1. The URL uses `https` (not `http`), matching your live site.
2. The domain matches exactly — including or excluding `www` consistently.
3. The path corresponds to the page content.

You can also use browser developer tools or an SEO audit tool to inspect the canonical tag on your deployed site.
