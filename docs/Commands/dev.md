---
title: "Dev Command — Build, Watch, and Serve in One Step"
description: "Use kiln dev to build your Obsidian vault, watch for file changes, and serve the site locally — all in a single command for a seamless development workflow."
---

# Dev Command

The `dev` command combines the [Generate Command](./generate.md) and [Serve Command](./serve.md) into a single workflow. It performs an initial full build of your vault, then watches for file changes and automatically rebuilds while serving the site on a local HTTP server. This gives you a live development loop where edits to your Obsidian notes are reflected in the browser without running separate commands.

The `dev` command accepts all the same flags as `generate` plus `--port` from `serve`, so you can customize themes, fonts, layouts, and other settings exactly as you would with a standalone build.

## Usage

```bash
kiln dev [flags]
```

A minimal dev session with default settings:

```bash
kiln dev
```

This reads from `./vault`, writes to `./public`, applies the default theme, font, and layout, and starts a local server on port `8080`. Open `http://localhost:8080` in your browser to preview the site.

## Flags

| Flag                    | Short | Default   | Description                                                                                                                              |
| ----------------------- | ----- | --------- | ---------------------------------------------------------------------------------------------------------------------------------------- |
| `--theme`               | `-t`  | `default` | Sets the color scheme. See [Themes & Visuals](../Features/User Interface/Themes.md) for available options.                               |
| `--font`                | `-f`  | `inter`   | Sets the font family. See [Fonts & Typography](../Features/User Interface/Fonts.md) for available options.                               |
| `--url`                 | `-u`  | `""`      | The public URL of your site (e.g., `https://example.com`). Required for [Sitemap.xml](../Features/SEO/Sitemap xml.md) and [Robots.txt](../Features/SEO/Robots txt.md) generation. |
| `--name`                | `-n`  | `My Notes` | The site name displayed in the browser tab and [Meta Tags & SEO](../Features/SEO/Meta Tags.md).                                        |
| `--input`               | `-i`  | `./vault` | Path to the source directory containing your Markdown notes.                                                                             |
| `--output`              | `-o`  | `./public`| Path where the generated HTML files are saved.                                                                                           |
| `--mode`                | `-m`  | `default` | Build mode. Use `default` for standard vault rendering or `custom` for [Custom Mode](../Features/Custom Mode/What is Custom Mode.md) with collection configs and templates. |
| `--layout`              | `-L`  | `default` | Page layout to use. See [Layouts](../Features/User Interface/Layouts.md) for available options.                                          |
| `--flat-urls`           |       | `false`   | Generate flat files (`note.html`) instead of directories (`note/index.html`).                                                            |
| `--disable-toc`         |       | `false`   | Hides the [Table of Contents](../Features/User Interface/Table of Contents.md) from the right sidebar.                                   |
| `--disable-local-graph` |       | `false`   | Hides the [Local Graph](../Features/User Interface/Local Graph.md) from the right sidebar.                                               |
| `--disable-backlinks`   |       | `false`   | Hides the Backlinks panel from the right sidebar.                                                                                        |
| `--lang`                | `-g`  | `en`      | Language code for the site (e.g., `en`, `it`, `fr`).                                                                                     |
| `--accent-color`        | `-a`  | `""`      | Accent color from the theme palette (`red`, `orange`, `yellow`, `green`, `blue`, `purple`, `cyan`). Defaults to the theme's built-in accent. |
| `--log`                 | `-l`  | `info`    | Log verbosity. Choose `info` or `debug`.                                                                                                 |
| `--host`                | `-H`  | `127.0.0.1` | Host address to bind the local development server to.                                                                                    |
| `--port`                | `-p`  | `8080`    | Port number for the local development server.                                                                                            |

## How It Works

When you run `kiln dev`, the following steps happen in order:

1. **Initial full build** — the entire vault is built exactly as the [Generate Command](./generate.md) would, producing a complete static site in the output directory.
2. **Modification-time baseline** — a snapshot of every file's last-modified timestamp is recorded. This baseline is used to detect which files change on subsequent edits.
3. **Dependency graph** — [wikilinks](../Features/Navigation/Wikilinks.md) between notes are parsed to build a graph of dependencies. When a note changes, the graph determines which other pages need to be rebuilt (for example, pages that link to or embed the changed note).
4. **Filesystem watcher** — a watcher (powered by [fsnotify](https://github.com/fsnotify/fsnotify)) is started on the input directory. File system events are debounced to avoid redundant rebuilds during rapid edits.
5. **Local HTTP server** — a development server starts on the configured port, serving the output directory with the same [clean URL support](./serve.md) as the standalone `serve` command.
6. **Incremental rebuild** — when the watcher detects a file change, it compares current modification times against the baseline, computes a changeset of affected files using the dependency graph, and triggers a rebuild.

Press `Ctrl+C` to cleanly shut down both the watcher and the server. The command intercepts `SIGINT` and `SIGTERM` signals for a graceful exit.

## Examples

### Basic Development Session

Start a dev server with the default settings and open `http://localhost:8080`:

```bash
kiln dev
```

### Themed Development with a Custom Port

Preview your vault with the Nord theme on port 3000:

```bash
kiln dev \
  --theme "nord" \
  --font "merriweather" \
  --name "My Digital Garden" \
  --port 3000
```

### Custom Input and Output Directories

Work with a non-default vault location and a separate build folder:

```bash
kiln dev \
  --input ./notes \
  --output ./dist \
  --url "https://example.com/docs"
```

## Related Commands

- [Generate Command](./generate.md) — run a one-off full build without watching or serving
- [Serve Command](./serve.md) — serve a previously built site without rebuilding
- [Clean Command](./clean.md) — remove build output before rebuilding
