---
title: "Init Command — Scaffold a New Kiln Project"
description: "Use kiln init to create a new Obsidian vault structure ready for static site generation. Get started with Kiln in seconds."
---

# Init Command

The `init` command creates a new Obsidian vault directory pre-configured for Kiln. It scaffolds the folder structure and a starter note so you can begin writing content and generating your site immediately. This is the recommended first step if you are starting a fresh project after [installing Kiln](../Installation.md).

## Usage

```bash
kiln init
```

By default, this creates a `vault` directory in your current working directory. You can specify a different name with the `--input` flag:

```bash
kiln init --input my-notes
```

## What It Creates

Running `init` generates the following structure:

```
vault/
└── Home.md
```

The `Home.md` file contains a welcome note with a heading and a prompt to run your first build. This file becomes the homepage of your generated site.

If a directory with the target name already exists, `init` exits with an error to prevent accidentally overwriting your content.

## Flags

| Flag      | Short | Default   | Description                                              |
| --------- | ----- | --------- | -------------------------------------------------------- |
| `--input` | `-i`  | `./vault` | Name of the directory to create.                         |
| `--log`   | `-l`  | `info`    | Sets the log level. Choose between `info` or `debug`.    |

## Full Workflow: From Init to Preview

After scaffolding your vault, you can build and preview it with two more commands. Here is the complete workflow from an empty directory to a running local site:

```bash
# 1. Scaffold a new vault
kiln init

# 2. Build the site
kiln generate --input ./vault --output ./public

# 3. Preview in your browser
kiln serve --output ./public
```

Open `http://localhost:8080` to see your site. From here you can add Markdown notes, organize them into folders, and rebuild with the [Generate Command](./generate.md) to see your changes.

## Adding Content to Your Vault

The scaffolded vault is intentionally minimal. To build it into a full site, add `.md` files to the vault directory — either by creating them manually or by pointing Obsidian at the folder. Kiln supports standard Obsidian features including [wikilinks](../Features/Navigation/Wikilinks.md), [tags](../Features/Navigation/Tags.md), callouts, and math expressions out of the box.

You can also use an existing Obsidian vault instead of running `init`. Just pass its path to the `generate` command directly:

```bash
kiln generate --input ./my-existing-vault
```

## Troubleshooting

If `init` reports that the vault directory already exists, either choose a different name with `--input` or remove the existing directory first. To check your vault for broken links before building, run the [Doctor Command](./doctor.md):

```bash
kiln doctor --input ./vault
```
