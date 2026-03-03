---
title: Math & LaTeX Rendering
description: Render LaTeX equations in your Obsidian vault with Kiln. Supports inline and block math using MathJax 4 with lazy loading.
---
# Math & LaTeX

Kiln renders mathematical notation directly from your Obsidian vault using **LaTeX** syntax and [MathJax 4](https://www.mathjax.org/). Inline formulas, block equations, and complex notation all work out of the box with zero configuration — write math in Obsidian, and Kiln publishes it to the web.

This is part of Kiln's support for [Obsidian Markdown](./Obsidian Markdown.md), which targets full parity with how your notes look in the Obsidian editor.

## Supported Syntax

Kiln recognizes four delimiter styles for LaTeX math, matching the syntax Obsidian uses.

### Inline Equations

Use single dollar signs `$...$` or escaped parentheses `\(...\)` to embed math within a sentence.

```markdown
The quadratic formula is $x = \frac{-b \pm \sqrt{b^2 - 4ac}}{2a}$ and applies to any second-degree polynomial.
```

**Result:** The quadratic formula is $x = \frac{-b \pm \sqrt{b^2 - 4ac}}{2a}$ and applies to any second-degree polynomial.

### Block Equations

Use double dollar signs `$$...$$` or escaped brackets `\[...\]` to create centered, standalone equations.

````markdown
$$
\int_0^\infty e^{-x^2} dx = \frac{\sqrt{\pi}}{2}
$$
````

**Result:**
$$
\int_0^\infty e^{-x^2} dx = \frac{\sqrt{\pi}}{2}
$$

Block equations are ideal for complex formulas, proofs, or any expression that benefits from its own line.

## Common LaTeX Examples

Below are copy-pasteable examples for frequently used notation.

### Fractions and Roots

```markdown
$$
\frac{a}{b} \quad \sqrt{x^2 + y^2} \quad \sqrt[3]{27}
$$
```

**Result:**
$$
\frac{a}{b} \quad \sqrt{x^2 + y^2} \quad \sqrt[3]{27}
$$

### Summations and Products

```markdown
$$
\sum_{i=1}^{n} i = \frac{n(n+1)}{2} \qquad \prod_{k=1}^{n} k = n!
$$
```

**Result:**
$$
\sum_{i=1}^{n} i = \frac{n(n+1)}{2} \qquad \prod_{k=1}^{n} k = n!
$$

### Matrices

```markdown
$$
\begin{pmatrix} a & b \\ c & d \end{pmatrix}
\begin{bmatrix} 1 & 0 \\ 0 & 1 \end{bmatrix}
$$
```

**Result:**
$$
\begin{pmatrix} a & b \\ c & d \end{pmatrix}
\begin{bmatrix} 1 & 0 \\ 0 & 1 \end{bmatrix}
$$

### Greek Letters and Symbols

```markdown
Inline symbols work naturally: $\alpha$, $\beta$, $\gamma$, $\Delta$, $\Omega$, $\nabla$, $\partial$, $\infty$.
```

**Result:** Inline symbols work naturally: $\alpha$, $\beta$, $\gamma$, $\Delta$, $\Omega$, $\nabla$, $\partial$, $\infty$.

## Escaping Dollar Signs

Because `$` triggers math rendering, you need to escape literal dollar signs with a backslash when you don't want math mode.

```markdown
The price is \$9.99, not $9.99$ (which would render as math).
```

The `processEscapes` option is enabled by default, so `\$` always produces a plain dollar sign.

## Performance

Rendering high-quality math on the web requires a powerful library, which can slow down pages that don't need it. Kiln solves this with a **lazy loading** strategy, similar to how it handles [Mermaid Graphs](./Mermaid Graphs.md).

- **Smart detection:** Before loading anything, Kiln scans the page content for math delimiters (`$$`, `\(`, `\[`). If none are found, the MathJax library is never downloaded.
- **On-demand loading:** When math is detected, MathJax 4 is fetched asynchronously from a CDN and renders equations as crisp SVGs with a global font cache.
- **Navigation aware:** With [client-side navigation](../Navigation/Client Side Navigation.md) enabled, MathJax re-renders automatically when you navigate to a new page containing equations — no full page reload needed.

This means blog posts, simple text notes, and pages without math have zero overhead from the math engine.

## How It Works

Kiln uses [Goldmark](https://github.com/yuin/goldmark) with the `goldmark-mathjax` extension to parse LaTeX delimiters during the Markdown-to-HTML conversion. The parsed math expressions are preserved in the HTML output, and MathJax 4 picks them up on the client side to render SVG output.

The SVG renderer produces resolution-independent equations that look sharp on any screen, and the global font cache ensures that repeated symbols don't increase page weight.
