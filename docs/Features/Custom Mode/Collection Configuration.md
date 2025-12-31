---
title: Collection Configuration (`config.json`)
description: Master the config.json file. Learn how to define strict data schemas for your content collections, from simple string fields to complex relationships, enums, and custom data types.
---
# Collection Configuration (`config.json`)

The `config.json` file is the heart of a collection. It tells Kiln how to read the frontmatter of your Markdown files and validates the data.

## Schema Definition
You can define fields using **Shorthand** (simple key-value) or **Longhand** (detailed objects). Some types require longhand definition because they need additional information.

## Required Fields
To be a valid collection you need to define the following fields.

- `collection_name`: A unique identifier for this collection (e.g., "blog", "projects"). This is used to avoid duplicate collections.

**Example `config.json`:**
```json
{
	"collection_name": "library",
	"title": "string",
	"publish_date": {
		"type": "date",
		"required": true
	},
	"rating": "integer",
	"cover": "image",
	"tags": "tags",
	"status": {
		"type": "enum",
		"values": ["reading", "completed", "wishlist"]
	},
	"author": {
		"type": "reference",
	    "reference": "authors"
	}
}
```

### Supported Types
Here' a table with all supported types.

| **Type**     | **Description**                       | **Definition Style**                     |
| ------------ | ------------------------------------- | ---------------------------------------- |
| `string`     | Basic text.                           | Shorthand                                |
| `date`       | YYYY-MM-DD.                           | Shorthand                                |
| `dateTime`   | ISO 8601 timestamp.                   | Shorthand                                |
| `boolean`    | True/False.                           | Shorthand                                |
| `integer`    | Whole numbers.                        | Shorthand                                |
| `float`      | Decimal numbers.                      | Shorthand                                |
| `image`      | Path to an image in the vault.        | Shorthand                                |
| `tag`        | Single Obsidian tag.                  | Shorthand                                |
| `tags`       | List of Obsidian tags.                | Shorthand                                |
| `enum`       | One value from a strict list.         | **Longhand** (requires `values` array)   |
| `reference`  | Link to a note in another collection. | **Longhand** (requires `reference` name) |
| `references` | List of links to notes.               | **Longhand** (requires `reference` name) |
| `custom`     | Arbitrary JSON object.                | **Longhand** (requires `data` object)    |

**Note on References:** The `reference` type allows you to create relational data. If you have an `authors` collection, you can link a book note directly to an author note, allowing your templates to pull data across collections.

## Example definitions
Here's a complete list of snippets for the collections definition and the resulting data that you can use in the [[Templating System|templates]].

### String
Used for basic text fields like titles, summaries, or names.

**Config:**
```json
{
	"simple_text": "string",
	"required_text": {
		"type": "string",
		"required": true
	}
}
```

**Template:**
```html
<h1>{{ .Page | get "simple_text" }}</h1>
<p>{{ .Page | get "required_text" }}</p>
```


### Boolean
True or false values. Useful for conditional logic like "featured" posts or "draft" states.

**Config:**
```json
{
	"is_featured": "boolean",
	"show_sidebar": {
		"type": "boolean",
		"required": true
	}
}
```

**Template:**
```html
{{ if .Page | get "is_featured" }}
  <span class="badge">Featured</span>
{{ end }}

<span>Sidebar Active: {{ .Page | get "show_sidebar" }}</span>
```

### Integer
Whole numbers. Used for counts, sorting orders, or simple ratings.

**Config:**
```json
{
    "sort_order": "integer",
    "rating": {
	    "type": "integer",
		"required": true
    }
}
```

**Template:**
```html
<div class="stars" data-rating="{{ .Page | get "rating" }}">
	Rating: {{ .Page | get "rating" }}/5
</div>
```

### Float
Decimal numbers. Used for prices, weights, or precise measurements.

**Config:**
```json
{
	"price": "float",
	"weight_kg": {
		"type": "float",
		"required": true
	}
}
```

**Template:**
```html
<span>Price: ${{ .Page | get "price" }}</span>
<span>Weight: {{ .Page | get "weight_kg" }} kg</span>
```

### Date
Dates without time (YYYY-MM-DD). Useful for publication dates or events. Use the Obsidian field "Date" to be sure about the formatting.

**Config:**
```json
{
	"publish_date": "date",
	"event_start": {
		"type": "date",
		"required": true
	}
}
```

**Template:**
```html
<time>{{ .Page | get "publish_date" }}</time>
```

### Date and Time
Dates with time (ISO 8601). Used for specific timestamps. Use the Obsidian field "Date" to be sure about the formatting.

**Config:**
```json
{
	"created_at": "dateTime",
	"last_modified": {
	    "type": "dateTime",
		"required": true
	}
}
```

**Template:**
```html
<span class="timestamp">{{ .Page | get "created_at" }}</span>
```

### Image
A path to an image file within your vault. Returns an `Asset` object that you can use to find the asset in your vault.

**Config:**
```json
{
	"cover_image": "image",
	"avatar": {
		"type": "image",
		"required": true
	}
}
```

**Template:**
```html
{{ with $img := .Page | get "cover_image" }}
<img src="{{ $img.RelPermalink }}" alt="Cover image" />
{{ end }}
```

### Tag
A single Obsidian tag string. Whenever you define a tag the note representation get's added to the site's tag field. This allows you to take every note tagged with a certain tag.

**Config:**
```json
{
	"main_category": "tag",
	"status": {
		"type": "tag",
		"required": true
	}
}
```

**Template:**
```html
<span class="category">
	{{ .Page | get "main_category" }} <!-- This prints the tag like #something -->
</span>

<!-- Knowing the tag, you can retrieve all the pages with that tag like this -->
{{ range $page := .Site | tag (.Page | get "tag_field") }}
<div>Tag: {{ $page | get "Title"}}</div>
{{ end }}
```

### Tags
A list (array) of Obsidian tags.

**Config:**
```json
{
	"categories": "tags",
	"meta_tags": {
		"type": "tags",
		"required": true
	}
}
```

**Template:**
```html
<ul>
	{{ range $tag := (.Page | get "categories") }}
	<li>$tag</li>
	{{ end }}
</ul>
```

### Enum
Enforces that the value must be one of a specific set of strings. Great for strict status control.

**Config:**
```json
{
	"status": {
		"type": "enum",
	    "values": ["draft", "published", "archived"],
	    "required": true
	},
	"size": {
		"type": "enum",
		"values": ["S", "M", "L", "XL"]
	}
}
```

**Template:**
```html
<div class='status-badge {{ .Page | get "status" }}'>
	{{ .Page | get "status" }}
</div>
```

### Reference
Links the current note to a **single** note in another collection. This creates a relationship (e.g., A Book has one Author).

**Config:**
```json
{
	"author": {
	    "type": "reference",
	    "reference": "authors", 
	    "required": true
	}
}
```

**Template:**
```html
{{ $author := .Page | get "author" }}
<div class="author-box">
  Written by: <a href="{{ $author.Permalink }}">{{ $author | get "name" }}</a>
</div>
```

### References
Links the current note to **multiple** notes in another collection. (e.g., A Movie has multiple Actors).

**Config:**
```json
{
	"cast": {
	    "type": "references",
	    "reference": "actors",
	    "required": true
	}
}
```

**Template:**
```html
<h3>Cast</h3>
<ul>
	{{ range (.Page | get "cast") }}
	<li>
		<a href="{{ .Permalink }}">{{ . | get "name" }}</a>
	</li>
	{{ end }}
</ul>
```

### Custom
Allows you to store arbitrary nested JSON data. This is useful for complex structures that don't map to a simple field. To have this field available you still need to add it to the frontmatter.

**Config:**
```json
{
	"metadata": {
		"type": "custom",
		"data": {
			"seo_score": 90,
			"reviewed_by": "admin"
		}
	}
}
```

**Template:**
```html
{{ $meta := .Page | get "metadata" }}

<div>SEO Score: {{ $meta.seo_score }}</div>
<div>Reviewer: {{ $meta.reviewed_by }}</div>
```