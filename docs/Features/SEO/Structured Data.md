---
title: Structured Data (SEO)
description: Learn how Kiln automatically generates JSON-LD structured data to help search engines understand your articles and breadcrumb trails.
---

# Structured Data (SEO)

Kiln helps your site stand out in search engine results by automatically generating **JSON-LD structured data** for your pages. 

Instead of requiring you to manually write complex SEO tags or install plugins, Kiln **bakes this data directly into your HTML** during the build process. Search engines like Google use this hidden information to display rich search results—such as showing your article's author, publish date, and a neat trail of navigation links right on the search page.

## Article Data

Whenever you build your site, Kiln generates an `Article` schema for your pages. It intelligently pulls this information from your site's main settings and the properties defined at the top of your markdown files.

To get the richest search results, simply include these standard properties in your page's frontmatter:

| Property      | Description                                                                         |
| :------------ | :---------------------------------------------------------------------------------- |
| `title`       | **Required.** The main headline. *If left empty, no structured data is generated.* |
| `description` | A short summary of the page or article.                                             |
| `author`      | The name of the person who wrote the content.                                       |
| `image`       | A URL link to the article's main cover image.                                       |

*(Note: Kiln automatically handles the creation and modification dates based on your file data!)*

**Example:**
```yaml
---
title: My First Post
description: A quick look at setting up my new site.
author: Jane Doe
image: [https://mysite.com/assets/cover.png](https://mysite.com/assets/cover.png)
---
```

## Breadcrumb Trails

Kiln also automatically generates BreadcrumbList data. Breadcrumbs help search engines understand the exact folder structure and hierarchy of your website (e.g., Home > Blog > My First Post).

You do not need to configure anything for this to work. When you run the [[generate]] command, Kiln maps out your folder structure and builds the proper navigation trail behind the scenes, ensuring search engines always know exactly where a page lives within your site.
