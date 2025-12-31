---
title: Templating System
description: Deep dive into Kiln's pipe-driven templating engine. Learn how to build custom HTML layouts, create reusable components, and use powerful functions like where, sort, and limit to control your content.
---
# Templating System
Kiln uses a powerful pipe-driven templating engine tailored for content. There are a couple of types of files and data structures to keep in mind.

## Data structures
When you're creating templates you'll need to access data from your Obsidian vault. To do that you'll need to use one of the available data structures:

- `.Page`: Contains data of the current page, like data from the frontmatter parsed from the [[Collection Configuration|configuration files]], sibling pages that are part of the same collection, the original name the note, the rendered content of the note and the webpath. 
- `.Site`: Contains data of the hole site, like data that was in the [[Environment Variables|environment variables file]], parsed assets, specific pages and tags.

We'll explore how to get that data shortly.

## Important files

### `layout.html`
This file, found in note collections, defines the standard look for every note for that collection. It wraps your Markdown content 

```html
<!DOCTYPE html>
<html>
	<head>
		<title>{{ .Page | get "Title" }}</title>
	</head>
	<body>
		<article>
			{{ .Page | get "Content" }}
		</article>
		{{ template "footer" . }}
	</body>
</html>
```

### `_components.html`
Any HTML file starting with `_` (e.g., `_card.html`) is treated as a component file. Inside these files, you must explicitly define your template using standard Go syntax.

**File:** `_card.html`
```html
{{ define "card" }}
<div class="post-card">
	<h3><a href="{{ .Permalink }}">{{ .Page | get "title" }}</a></h3>
	<p>{{ .Page | get "summary" }}</p>
</div>
{{ end }}
```

**Usage in Layout:**
```html
<div>
	{{ range .Page | get "Siblings" }}
		{{ template "card" . }}
	{{ end }}
</div>
```

### `same-name.html`
Sometimes a specific page needs a unique layout (e.g., a Landing Page inside a "Pages" collection, an "About" page or even the "Home"). If you have a file named `contact.md` and you create `contact.html` in the same folder, Kiln will use `contact.html` instead of `layout.html` for that specific note. 

This note will have all the available fields from the `.Page` and `.Site` fields like any other. 

## Template Functions
We provide a set of Go template functions to query and manipulate your content.

### `get`
Safely retrieves frontmatter fields defined in your config and other fields available to the `.Page` or `.Site`.

```html
{{ .Page | get "Title" }} <!-- Original title of the note -->
{{ .Page | get "Siblings" }} <!-- Siblings pages (pages that are part of the same collection -->
{{ .Page | get "Path" }} <!-- Webpath of the page -->
{{ .Page | get "Content" }} <!-- Rendered content of the note -->
{{ .Page | get "custom_field" }}
```

#### Filtering: `where` & `where_not`
Filter lists of pages based on frontmatter values.

```html
{{ range $page := .Page | get "Siblings" | where "published" true }}
	<a href="{{ $page | get .Path }}">{{ $page | get "Title" }}</a>
{{ end }}

{{ range .Siblings | where_not "status" "draft" }}
	...
{{ end }}
```

#### Slicing: `limit` & `offset`
Control pagination or list sizes.

```html
{{ range .Page | get "Siblings" | limit 3 }}
  ...
{{ end }}

{{ range .Page | get "Siblings" | offset 1 | limit 5 }}
  ...
{{ end }}
```

#### Sorting: `sort`
Sort a list of pages by a specific field key.

```html
{{ range .Page | get "Siblings" | sort "date" "desc" }}
  ...
{{ end }}

{{ range .Page | get "Siblings" | sort "title" "asc" }}
  ...
{{ end }}
```