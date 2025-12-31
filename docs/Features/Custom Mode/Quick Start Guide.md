---
title: "Quick Start: Building a Blog"
description: Get up and running with Kiln Custom Mode in minutes. A step-by-step tutorial on building a fully functional, type-safe blog using Obsidian as your CMS, from project setup to deployment.
---
# Quick Start: Building a Blog
This guide will walk you through building a simple, data-driven blog using **Kiln Custom Mode**. We will create a `posts` collection, set up the data schema, create a reusable post card component, and build the layout.

## Project Setup
Create a new folder for your project. Inside, we will create a folder structure that looks like this:

```Plaintext
my-blog/
├── env.json             # Global site variables
└── posts/               # Our blog collection
    ├── config.json      # The schema for blog posts
    ├── layout.html      # The template for blog posts
    ├── _card.html       # A reusable component
    └── first-post.md    # Content
```

## Global Variables (`env.json`)
First, let's define some site-wide data. Create `env.json` in the root.

```json
{
	"site_name": "The Kiln Chronicle",
	"author_name": "John Obsidian",
	"twitter_url": "https://twitter.com/kiln_generated"
}
```

## The Collection Config (`config.json`)
Navigate to the `posts` folder. We need to tell Kiln what a "Post" actually looks like. Create `config.json`.

We want every post to have a Title, a Publish Date, a Summary, and a list of Tags.
```json
{
	"collection_name": "posts",
	"title": "string",
	"date": {
		"type": "date",
		"required": true
	},
	"summary": "string",
	"tags": "tags",
	"featured": "boolean"
}
```

## Creating Content
Now that the rules are set, create a Markdown file: `posts/hello-world.md`.

```markdown
---
title: Hello World
date: 2023-10-27
summary: This is my first post generated with Kiln Custom Mode.
tags: [kiln, update]
featured: true
---

# Welcome to the blog

This is standard Markdown content. Because we are in **Custom Mode**, this text will be injected into the `{{ .Page | get "Content" }}` variable in our layout.
```

## The Component (`_card.html`)
We want to display a list of "Related Posts" at the bottom of every article. Instead of writing HTML twice, let's make a component.

Create `posts/_card.html`. Note that we use `{{ define }}` to name the template.
```html
{{ define "post_card" }}
<div class="card">
	<h3>
		<a href="{{ .Permalink }}">{{ .Page | get "Title" }}</a>
	</h3>
	<small>{{ .Page | get "Date" }}</small>
	<p>{{ .Page | get "Summary" }}</p>
</div>
{{ end }}
```

## The Layout (`layout.html`)
Finally, create `posts/layout.html`. This will be the wrapper for `hello-world.md` and any future Markdown files in this folder.

We will use:

1. **Global Env** to get the site name.
2. **`get`** to retrieve frontmatter from the current page.
3. **`where_not`** to filter the sibling list (to avoid linking to the current page).
4. **`sort`** to order posts by date.
5. **`limit`** to show only 3 recent posts.

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>{{ .Page | get "Title" }} | {{ .Site | env "site_name" }}</title>
</head>
<body>

    <header>
        <h1>{{ .Site | env "site_name" }}</h1>
    </header>

    <main>
        <article>
            <h1>{{ .Page | get "Title" }}</h1>
            <p class="meta">By {{ .Site | env "author_name" }} on {{ .Page | get "Date" }}</p>
            
            <div class="content">
                {{ .Page | get "Content" }}
            </div>
        </article>
        
        <hr>

        <section class="recent-posts">
            <h2>More to read</h2>
            <div class="grid">
                {{ range .Page | get "Siblings" | where_not "Title" (.Page | get "Title") | sort "date" "desc" | limit 3 }}
                    {{ template "post_card" . }}
                {{ end }}
            </div>
        </section>
    </main>

</body>
</html>
```

## Build and Serve
Go to your terminal and run:

```bash
# generate the site
kiln generate --input ./my-blog --output ./dist

# preview it
kiln serve ./dist
```

Open your browser to `http://localhost:8080/posts/hello-world.html`.

You have successfully built a structured, type-safe blog using Obsidian as your CMS. To add more posts, simply duplicate your markdown file—Kiln will automatically validate the data and update the lists.

