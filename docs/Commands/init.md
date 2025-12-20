---
title: Init Command
description: The init command scaffolds a new Kiln project, generating a sample vault structure to help you get started quickly.
---

# Init Command

The `init` command is designed to help you hit the ground running. It scaffolds a new, empty Kiln project that you can then open with Obsidian. This is the perfect starting point if you are testing Kiln for the first time.

> [!warning] Work in progress
> This command is rather useless in the current state. It just creates an empty vault. 
> Next iterations of this will scaffold a starter site.

## Usage

```bash
./kiln init
```

## Flags

| Flag      | Short version | Default value | Description                                  |
| --------- | ------------- | ------------- | -------------------------------------------- |
| `--input` | `-i`          | `"vault"`     | The name of the directory to be created. |

## What it Generates

Running this command will create a new folder (defaulting to `./vault`) containing a standard Obsidian-style structure:

```
📂 vault/
└── 📄 Home.md
```

You can immediately run the `generate` command against this folder to build your first site.

```
# Initialize
./kiln init

# Build
./kiln generate --input "vault"
```