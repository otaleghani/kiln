---
title: Kiln quirks
description: A series of different quirks to keep in mind when using Kiln.
---
# Quirks
Here's a series of known quirks to keep in mind when using Kiln. 

## Same filename, multiple extensions

Obsidian allows you to have multiple items with the exact same name living in the same directory. For example, you might have a **folder** named `Example`, a **canvas** named `Example.canvas`, and a **note** named `Example.md` all sitting side-by-side.

Instead of generating separate, confusing pages for each version, Kiln treats this as a battle for dominance. It consolidates them into a single URL based on a strict override hierarchy.

Right now, the pecking order looks like this (from lowest priority to highest):
```
Folder -> Canvas -> Base -> Note
```

This means a **Note** will override everything else. If there is no note, the [[Bases|Base]] takes over. If there is no base, the [[Obsidian Canvas|Canvas]] wins. The [[Folders|Folder]] is the bottom of the food chain; it only displays if nothing else overrides it.

If you do not wish to have this override happen, just use a different name for your note, you marsupial.

### Why do this?
This hierarchy allows you to replace a standard Folder page with a page of your choice:

- **Canvases:** Use a Canvas to replace a boring file list with a fun, visual map of the folder's contents.
- **Bases:** Use a Base to create a structured database view of the notes inside.
- **Notes:** If you want great SEO, nothing beats a good ol' note with high-quality content and smart internal linkage.

## About links

In the case of non `--flat-urls` version of the generated sites, the `/index.html` file of the folder is still present. It's not linked by any other page. True overrides are only present in `--flat-urls` mode, where the `/index.html` of the folder get's overrode by your other files.

So this begs the question: **are folder pages still accessible?** Yes, but only if you manually search for them in the search bar of your browser, by adding the `/` to the end of the page.