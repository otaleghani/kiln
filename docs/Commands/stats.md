---
title: "Stats Command — Vault Word Count and Note Metrics"
description: "Use the kiln stats command to count notes, total words, and find the longest note in your Obsidian vault before building your site."
---

# Stats Command

The `stats` command gives you a quick overview of your Obsidian vault's size and content. It scans every Markdown file in your input directory and reports the total number of notes, the combined word count, and which note is the longest. This is useful for tracking your writing progress over time, estimating site build size, or getting a sense of scale before running [Generate Command](./generate.md).

## What It Reports

Kiln calculates three metrics from your vault:

- **Total Notes** — the number of `.md` files found in the input directory, including all subdirectories.
- **Total Words** — the combined word count across every note, estimated by splitting on whitespace. Frontmatter content is included in the count.
- **Longest Note** — the file name and word count of the single largest note in your vault.

These metrics cover all Markdown files recursively. Non-Markdown files (images, PDFs, attachments) are ignored.

## Usage

Run the command from the same directory where your vault is located:

```bash
kiln stats
```

**Example output:**

```
METRIC        VALUE
------        -----
Total Notes   38
Total Words   6533
Longest Note  Client Side Navigation.md (412 words)
```

To analyze a vault in a different location, use the `--input` flag:

```bash
kiln stats --input ./my-notes
```

## Flags

| Flag      | Short | Default   | Description                                                              |
| --------- | ----- | --------- | ------------------------------------------------------------------------ |
| `--input` | `-i`  | `"vault"` | The path to the directory containing your Markdown notes to be analyzed. |
| `--log`   | `-l`  | `info`    | Sets the log level. Choose between `info` or `debug`.                    |

## Common Use Cases

### Tracking writing progress

Run `kiln stats` periodically to see how your vault grows. Comparing the total word count and note count over time helps you stay motivated and measure output.

### Auditing before deployment

Before generating your site with the [Generate Command](./generate.md), run `stats` to confirm the expected number of notes will be processed. If the count seems off, you may have notes in the wrong directory or files that aren't saved as `.md`.

### Checking vault health alongside Doctor

Pair `stats` with the [Doctor Command](./doctor.md) for a full pre-build checkup. While `doctor` catches broken wikilinks and structural issues, `stats` confirms that all your content is being picked up by the scanner.

### Debugging input directory issues

If `stats` reports zero notes, double-check that the `--input` path points to the correct folder. The default is `./vault`, so if your Obsidian vault lives elsewhere, pass the right path:

```bash
kiln stats --input ./content
```

## How Word Count Works

Kiln counts words by reading each `.md` file and splitting its full text content on whitespace. This means frontmatter fields (title, tags, description) contribute to the word count. The count is an approximation — code blocks, YAML, and metadata are not excluded. For most vaults this gives a reliable ballpark figure.

## Related Commands

After reviewing your vault stats, you can proceed with building and previewing your site:

- [Generate Command](./generate.md) — build your vault into a static site
- [Serve Command](./serve.md) — preview the generated site locally
- [Doctor Command](./doctor.md) — scan for broken wikilinks before deploying
- [Clean Command](./clean.md) — remove previous build output before regenerating
