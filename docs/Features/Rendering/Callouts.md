---
title: Callouts
description: Learn how to use Obsidian-style callouts to highlight important information, create warnings, and organize content with collapsible blocks.
---

# Callouts

Kiln supports standard **Obsidian-style callouts**. These allow you to highlight specific parts of your content with distinct colors and icons, making your notes more readable and structured.

## Syntax

To create a callout, use the blockquote syntax `>` followed by the callout type in square brackets `[!type]`.

```markdown
> [!info]
> This is an info callout.
> It can contain **Markdown**, links, and other elements.
```

### Custom Titles

By default, the title of the callout matches its type (e.g., "Info"). You can define a custom title by adding text after the brackets.

```
> [!warning] Heads Up!
> This callout has a custom title.
```

## Collapsible Callouts

You can make callouts foldable (collapsible) by adding a `+` or `-` symbol immediately after the callout type.

- `+`: The callout is **expanded** by default.
- `-`: The callout is **collapsed** by default.

> [!faq]- Can I hide this content?
> Yes! This content is hidden until the user clicks the arrow to expand it. 
> This is perfect for FAQs, spoilers, or long reference material.

> [!faq]+ Can I leave the content open? 
> Yes! Just use the `+` instead of the `-`

## Supported Types

Kiln supports the full set of standard Obsidian callout types and their aliases.

|**Type**|**Aliases**|**Color Profile**|
|---|---|---|
|`note`||🔵 Blue|
|`abstract`|`summary`, `tldr`|🟢 Cyan|
|`info`||🔵 Blue|
|`todo`||🔵 Blue|
|`tip`|`hint`, `important`|🟢 Cyan|
|`success`|`check`, `done`|🟢 Green|
|`question`|`help`, `faq`|🟠 Orange|
|`warning`|`caution`, `attention`|🟠 Orange|
|`failure`|`fail`, `missing`|🔴 Red|
|`danger`|`error`|🔴 Red|
|`bug`||🔴 Red|
|`example`||🟣 Purple|
|`quote`|`cite`|⚪ Gray|

### Icons

Kiln automatically pairs each callout type with its corresponding icon from the [Lucide](https://lucide.dev/) library, ensuring visual consistency with the Obsidian ecosystem.