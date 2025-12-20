---
title: Obsidian Markdown
description: Kiln is built to support Obsidian Flavored Markdown, a rich superset of standard Markdown that includes wikilinks, callouts, and math support.
---

# Obsidian Markdown

Kiln is designed with a "Parity First" philosophy: **If it renders in Obsidian, it should render in Kiln.**

While Kiln is built on robust, standard-compliant parsers, it specifically targets **Obsidian Flavored Markdown**. This ensures that your personal knowledge base translates seamlessly to the web without you having to rewrite your notes to fit a stricter standard.

## What is a "Flavor"?

Markdown was originally created in 2004 with a very simple goal: to be an easy-to-write format for HTML. However, the original specification did not include features that modern users rely on, like tables, footnotes, or mathematical formulas.

To solve this, different platforms created their own "Flavors"—extensions of the original language to add these missing features.

### CommonMark
This is the "strict" standard. It is highly compatible but very basic. It supports headers, lists, and bold text, but lacks tables or task lists.

### GitHub Flavored Markdown (GFM)
This is the industry standard for developers. It adds **Tables**, **Task Lists** (`- [ ]`), **Strikethrough** (`~~text~~`), and auto-linking URLs. Kiln supports all GFM features by default.

### Obsidian Flavored Markdown
Obsidian extends GFM even further to turn Markdown into a tool for networked thought. This is the flavor Kiln supports.

## Supported Features

Kiln's parser includes support for the following Obsidian-specific syntax extensions:

| Feature | Syntax | Description |
| :--- | :--- | :--- |
| [[Wikilinks]] | `[[Note]]` | Internal linking between files. |
| [[Wikilinks#Embedding Media|Embed]] | `![[Image.png]]` | Displaying images or transcluding notes. |
| [[Callouts]] | `> [!info]` | Colored blockquotes for distinct content. |
| [[Math]] | `$E=mc^2$` | LaTeX rendering via MathJax. |
| **Highlighting** (Not yet supported, soon to be) | `==text==` | Visual highlighting of text. |
| **Comments** (Not yet supported, soon to be) | `%% comment %%` | Text visible in the editor but hidden in the output. |

## Future Compatibility

Kiln uses `goldmark`, a highly extensible Markdown parser for Go. While the current focus is strictly on Obsidian compatibility, the architecture allows for future support of other flavors (such as Hugo-specific shortcodes or Pandoc extensions) should the need arise.