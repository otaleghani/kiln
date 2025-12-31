---
title: What is Custom Mode
description: "Discover Kiln's Custom Mode: transform Obsidian into a powerful headless CMS. Learn how to move beyond simple vault mirroring to build structured, data-driven static websites with custom schemas and layouts."
---
> [!warning] 
> Custom mode is still in it's early stages. If you find any bugs please report them in the [Github Issues page](https://github.com/otaleghani/kiln/issues).
# What is Custom Mode

Kiln operates in two distinct modes:

1. **Default Mode (Vault Mirror):** This is the classic behavior. It takes your Obsidian vault and mirrors it 1:1. If it renders in Obsidian, it renders in Kiln. This is perfect for digital gardens, wikis, and documentation.
2. **Custom Mode (Obsidian CMS):** This transforms Obsidian into a **Headless CMS**. Instead of just mirroring your notes, Kiln treats your vault as a database and your folders as **Collections**. You define strict data schemas for your frontmatter, write custom HTML templates, and build fully bespoke websites while keeping the writing experience inside Obsidian.

**Choose Custom Mode if:**
- You want a blog, portfolio, or any content focused site with a specific design that doesn't look like a "knowledge base".
- You need to enforce data structure (e.g., every "Book Review" _must_ have a rating and an author).
- You want full control over the HTML/CSS output.

The core idea of custom mode is to have a 1:1 representation of your websites content on Obsidian. This allows you to use all the cool features of Obsidian, like [[Wikilinks]], and have that reflected on your site.

To access custom mode just add the flag `--mode "custom"` to your [[generate]] command, like so:
```bash
kiln generate --mode "custom"
```

## Project Structure
In Custom Mode, your Obsidian vault becomes the source of truth for content, but the structure relies on a few key files.

Here's an example file structure:
```bash
my-blog/
├── index.md             # Index note of our vault
├── index.html           # Template to use for the index note
├── about.md             # About note in our vault
├── about.html           # Template to use for the about note
├── env.json             # Global site variables
└── posts/               # Our blog collection
    ├── config.json      # The schema for blog posts
    ├── layout.html      # The template for blog posts
    ├── _card.html       # A reusable component
    └── first-post.md    # Content
└── authors/             # Our authors collection
    ├── config.json      # The schema for author posts
    ├── layout.html      # The template for author pages
    └── john-obsidian.md # Content
```

## The Collection Concept
Every folder in your vault is treated as a **Collection**. To activate a folder as a valid collection, you simply add a `config.json` file inside it. This file will be parsed and used to validate the content found in the frontmatter of notes in that collection.

## File Types
- **`config.json`** (required for collections): Defines the schema (fields and types) for the notes in that folder.
- **`layout.html`** (required for collections): The default HTML template used for every note in that folder.
- **`filename.html`**: A specific template override. If you have `about.md`, creating `about.html` will force Kiln to use that specific layout instead of the generic one.
- **`_component.html`**: Any HTML file starting with an underscore is treated as a reusable component (partial), available to all layouts.
- **`env.json`**: A global file in the project root for site-wide variables.

## Build Process Overview
When you run `kiln build` in Custom Mode, the engine performs the following steps:

1. **Scan:** Maps the file system, ignoring hidden files and hidden directories.
2. **Environment:** Loads `env.json`.
3. **Config:** Validates all `config.json` files to establish the schema.
4. **Validation:** Parses layouts, components, and notes. It strictly enforces the types defined in your configs (e.g., if a field is `integer`, a string value will throw an error).
5. **Render:** Generates the static HTML using your templates and data.
6. **Assets:** Copies all static files (images, CSS, JS) to the build directory.

## Learn more
Here's a list of resources related to custom mode to help get started:

- Understand the [[Collection Configuration]] to create content collections.
- Setup the `env.json` file to enable site-wide [[Environment Variables]].
- How to work with [[Styles]] and setup [Tailwind CSS](https://tailwindcss.com/).
- Learn how to leverage the [[Templating System]].
- Get your hands dirty with this fast and easy to follow [[Quick Start Guide]].