---
title: Environment Variables (`env.json`)
description: Learn how to manage global site data with env.json. Define site-wide constants like API keys, social links, and titles that are accessible across all your templates and components.
---
# Environment Variables (`env.json`)

Place an `env.json` file in the root of your project to define site-wide constants. This is useful for data that doesn't belong in a specific note but is needed globally (like site title, social links, or API keys).

**Example:**
```json
{
	"site_name": "My Digital Garden",
	"twitter_handle": "@otaleghani",
	"production_url": "https://kiln.so"
}
```

## Accessing Environment Variables
You can access these in your templates via the global environment accessor.

```html
<div>Env: {{ .Site | env "site_name" }}</div>
```
