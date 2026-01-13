---
title: Generate Command
description: The generate command is the core of Kiln. Learn how to build your static site and customize it using flags for themes, fonts, and SEO settings.
---
# Generate Command

The `generate` command is the primary workhorse of Kiln. It reads your Markdown files from the source directory and compiles them into a fully functional, static HTML website ready for deployment.

You will run this command every time you want to update your site with new content.

## Usage

```bash
./kiln generate [flags]
```

## Flags
There are some flags that you can use to customize the output:

| **Flag**                | **Short** | **Default** | **Description**                                                                                                                  |
| ----------------------- | --------- | ----------- | -------------------------------------------------------------------------------------------------------------------------------- |
| `--theme`               | `-t`      | `default`   | Sets the visual color scheme. See [[Themes]] for options.                                                                        |
| `--font`                | `-f`      | `inter`     | Sets the typography family. See [[Fonts]] for options.                                                                           |
| `--url`                 | `-u`      | `""`        | The final public URL of your site (e.g., `https://example.com`). Required for generating the [[Sitemap xml]] and [[Robots txt]]. |
| `--name`                | `-n`      | `My Notes`  | The global name of your site. This appears in the browser tab and [[Meta Tags]].                                                 |
| `--input`               | `-i`      | `vault`     | The path to your source folder containing the Markdown notes.                                                                    |
| `--output`              | `-o`      | `public`    | The path where the generated HTML files will be saved.                                                                           |
| `--flat-urls`           |           | `false`     | Generate flat HTML files (`note.html`) instead of pretty directories (`note/index.html`)                                         |
| `--log`                 | `-l`      | `info`      | Sets the log level. You can choose between `info` or `debug`.                                                                    |
| `--disable-toc`         |           | `false`     | Disables the Table of contents on the right sidebar. If the local graph is disabled too, hides the right sidebar.                |
| `--disable-local-graph` |           | `false`     | Disables the Local graph. If the table of contents is disabled too, hides the right sidebar.                                     |
| `layout`                | `-L`      | `default`   | Layout to use. Choose between 'default' and the others. Find out more about which layout is available at [[Layouts]].            |

## Examples

### Basic Test

For a quick local test, you can run the command without arguments (uses default settings):
```bash
./kiln generate
```

### Production Build

When deploying your site to the web, you should always include the `url` and `name` flags to ensure SEO features work correctly.
```bash
./kiln generate \
  --name "My Digital Garden" \
  --url "https://notes.mysite.com" \
  --theme "nord" \
  --font "inter"
```