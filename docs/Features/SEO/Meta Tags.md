---
title: Meta Tags & SEO
description: Learn how to optimize your site's SEO by customizing page titles and meta descriptions using standard Markdown frontmatter.
---

# Meta Tags & SEO

Kiln is built to be **SEO-native**. It automatically generates the necessary `<meta>` tags to ensure your content looks great in search engine results and on social media.

This is controlled entirely through the **Frontmatter** (the YAML block at the very top of your Markdown files).

## Configuration

To customize how a page appears to search engines, simply define the `title` and `description` fields in your note.

```yaml
---
title: My Custom Page Title
description: A short summary of the page content (ideal length is 150-160 characters).
---
```

### The Title Field

The `title` field serves two purposes:

1. **Browser Tab:** It sets the text displayed in the user's browser tab.
2. **Search Result:** It becomes the clickable headline in Google or Bing search results.
3. **Open Graph:** It is used as the title for social media cards (Twitter/X, LinkedIn, Slack).

> **Fallback:** If you do not provide a `title`, Kiln will automatically use the filename of the note (e.g., `My Note.md` becomes "My Note").

### The Description Field

The `description` field populates the `<meta name="description">` tag.

1. **Search Snippet:** Search engines use this text to describe your page below the link.
2. **Social Preview:** When you share a link on Discord or Twitter, this text appears under the title in the preview card.

## Automatic HTML Generation

When Kiln builds your site, it transforms your frontmatter into standard HTML tags automatically.

``` yaml
---
title: Guide to Kiln
description: The ultimate guide to using the Kiln static site generator.
---
```

**Output (HTML):**

``` html
<head>
  <title>Guide to Kiln</title>
  <meta name="description" content="The ultimate guide to using the Kiln static site generator.">
  
  <meta property="og:title" content="Guide to Kiln">
  <meta property="og:description" content="The ultimate guide to using the Kiln static site generator.">
</head>
```
