---
title: Sitemap.xml
description: Kiln automatically builds an XML sitemap to help search engines discover and index your content efficiently. Learn how to enable and configure it.
---

# Sitemap.xml

A **Sitemap** is a blueprint of your website that helps search engines find, crawl, and index all of your website's content. Kiln automatically generates a standards-compliant `sitemap.xml` file every time you build your site.

## Why it matters

Without a sitemap, search engines like Google rely on finding links on one page to discover another. If you have a new note that isn't linked from anywhere yet (an "orphan page"), Google might never find it.

The `sitemap.xml` solves this by providing a complete list of every single public page on your site, ensuring 100% of your content is discoverable.

## Configuration

Because a sitemap requires **absolute URLs** (e.g., `https://example.com/page` instead of just `/page`), Kiln needs to know your domain name to generate it.

You must provide the `--url` flag when running the generate command:

```bash
./kiln generate --url "[https://kiln.talesign.com](https://kiln.talesign.com)"
```

> [!warning] Missing URL Flag 
> If you do not provide the `--url` flag, Kiln will skip generating the sitemap entirely to prevent creating an invalid file.

## Inclusion Logic

Kiln automatically scans your vault to build the sitemap.

- **Included:** All Markdown files (`.md`) and Canvas files (`.canvas`) are added automatically.
- **Excluded:** Any file or folder starting with a dot `.` (hidden files) is automatically excluded from the sitemap to protect your private data or drafts.