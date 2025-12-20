---
title: Robots.txt
description: Learn how Kiln automatically generates a robots.txt file to guide search engine crawlers and why the base URL configuration is required.
---
# Robots.txt

Kiln automatically generates a standard `robots.txt` file at the root of your site during the build process. This file serves as the "gatekeeper" for your website, giving instructions to web crawlers (like Googlebot) about which pages they are allowed to access.

## Requirement

To enable the generation of this file, you **must** specify your site's public URL during the build command.

```bash
./kiln generate --url "[https://your-domain.com](https://your-domain.com)"
```

> [!warning] URL Flag Missing 
> If you do not provide the `--url` flag, Kiln will skip generating the `robots.txt` file entirely.

## Default Configuration

Kiln generates a standard, permissive configuration designed for public documentation sites and blogs. It allows all search engines to crawl all content.

**Output:**
```txt
User-agent: *
Allow: /
Sitemap: https://your-domain.com/sitemap.xml
```

_Note: The `Sitemap` directive is automatically appended to help crawlers discover your content map instantly._