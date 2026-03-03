---
title: "Doctor Command — Find Broken Links in Your Vault"
description: "Run kiln doctor to scan your Obsidian vault for broken wikilinks before deploying. Catches missing references, renamed notes, and dead links."
---

# Doctor Command

The `doctor` command scans your Obsidian vault for broken [wikilinks](../Features/Navigation/Wikilinks.md) and missing references. Run it before deploying your site to catch dead links that would result in 404 errors for your readers.

It is best practice to run `doctor` after renaming or deleting notes, or before any production build with the [generate command](./generate.md).

## What It Checks

The doctor performs a full scan of every `.md` file in your vault and validates each wikilink against an index of all known files. It catches:

- **Broken wikilinks** — links like `[[Note Name]]` where the target file no longer exists.
- **Renamed notes** — if you renamed a file but forgot to update links pointing to it, the doctor flags each one.
- **Missing image references** — wikilinks to `.png`, `.jpg`, or `.jpeg` files that cannot be found.
- **Missing canvas files** — references to `.canvas` files that have been moved or deleted.

The doctor also handles common wikilink variations correctly. Aliased links like `[[Note Name|Custom Text]]` and anchor links like `[[Note Name#Section]]` are both resolved to the target note before validation.

## Usage

```bash
./kiln doctor
```

When all links are valid, you will see:

```
kiln: 2025/12/20 11:11:16 Diagnosing vault...
kiln: 2025/12/20 11:11:16 No broken links found
```

When broken links are found, the doctor reports each one with the file path and the unresolved link target:

```
kiln: WRN Found broken link path=notes/index.md link="Old Note Name"
kiln: ERR Found broken links number=1
```

Use the `debug` log level to get more detailed output during the scan:

```bash
./kiln doctor --log debug
```

## Flags

| Flag      | Short | Default   | Description                                                   |
| --------- | ----- | --------- | ------------------------------------------------------------- |
| `--input` | `-i`  | `./vault` | Path to the directory containing your vault.                  |
| `--log`   | `-l`  | `info`    | Sets the log level. Choose between `info` or `debug`.         |

## Recommended Workflow

Run `doctor` as part of your build process to prevent broken links from reaching production. A typical workflow looks like this:

```bash
# Check for broken links first
./kiln doctor --input ./vault

# If no issues, build the site
./kiln generate --name "My Notes" --url "https://notes.example.com"
```

You can also run it after a large reorganization of your vault to verify that all internal references still resolve. If you need to start fresh after fixing issues, the [clean command](./clean.md) removes old build output before regenerating.
