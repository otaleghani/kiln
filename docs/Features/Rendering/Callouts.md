---
title: "Callouts: Info Boxes, Warnings, and Collapsible Blocks"
description: Use Obsidian-style callouts in Kiln to add info boxes, warnings, tips, and collapsible sections to your static site pages.
---

# Callouts — Info Boxes, Warnings, and Collapsible Blocks

Kiln renders **Obsidian-style callouts** so you can highlight important information, surface warnings, and organize long-form content with collapsible sections. Callouts carry over from your vault to the published site with matching colors, icons, and fold behavior — no extra configuration needed.

Callouts are part of the [Obsidian Markdown](./Obsidian Markdown.md) syntax that Kiln supports alongside features like [Wikilinks & Embeds](../Navigation/Wikilinks.md) and [Math & LaTeX](./Math.md).

## How to Create a Callout

Use the standard blockquote syntax `>` followed by a callout type in square brackets `[!type]`.

```markdown
> [!info]
> This is an info callout.
> It can contain **Markdown**, links, and other elements.
```

### Add a Custom Title

By default, the title matches the callout type (e.g., "Info"). Add text after the brackets to set a custom title.

```markdown
> [!warning] Heads Up!
> This callout has a custom title.
```

## Collapsible Callouts

Make any callout foldable by adding `+` or `-` right after the type keyword.

- `+` — the callout starts **expanded**.
- `-` — the callout starts **collapsed**.

Collapsible callouts are useful for FAQs, spoiler sections, or long reference material that readers can reveal on demand.

```markdown
> [!faq]- Can I hide this content?
> Yes! This content is hidden until the reader clicks to expand it.

> [!faq]+ Can I leave the content open?
> Yes! Use `+` instead of `-` to default to expanded.
```

> [!faq]- Can I hide this content?
> Yes! This content is hidden until the user clicks the arrow to expand it.
> This is perfect for FAQs, spoilers, or long reference material.

> [!faq]+ Can I leave the content open?
> Yes! Just use the `+` instead of the `-`

## Supported Callout Types

Kiln supports the full set of standard Obsidian callout types and their aliases. Each type renders with a distinct color and icon so readers can quickly identify the purpose of the block.

| **Type**   | **Aliases**            | **Color** |
| ---------- | ---------------------- | --------- |
| `note`     |                        | Blue      |
| `abstract` | `summary`, `tldr`      | Cyan      |
| `info`     |                        | Blue      |
| `todo`     |                        | Blue      |
| `tip`      | `hint`, `important`    | Cyan      |
| `success`  | `check`, `done`        | Green     |
| `question` | `help`, `faq`          | Orange    |
| `warning`  | `caution`, `attention` | Orange    |
| `failure`  | `fail`, `missing`      | Red       |
| `danger`   | `error`                | Red       |
| `bug`      |                        | Red       |
| `example`  |                        | Purple    |
| `quote`    | `cite`                 | Gray      |

### Callout Examples

> [!note]

> [!abstract]
> You could also use `summary` or `tldr`

> [!info]

> [!todo]

> [!tip]
> You could also use `hint` or `important`

> [!success]
> You could also use `check` or `done`

> [!question]
> You could also use `help` or `faq`

> [!warning]
> You could also use `caution` or `attention`

> [!failure]
> You could also use `fail` or `missing`

> [!danger]
> You could also use `error`

> [!bug]

> [!example]

> [!quote]
> You could also use `cite`

### Callout Icons

Kiln automatically pairs each callout type with its corresponding icon from the [Lucide](https://lucide.dev/) library, keeping visual consistency with Obsidian.

## Styling Callouts with Themes

Callout colors and appearance adapt to the active [theme](../User Interface/Themes.md) and respect [Light & Dark Mode](../User Interface/Light-Dark Mode.md) settings, so they look correct in both modes without extra work.
