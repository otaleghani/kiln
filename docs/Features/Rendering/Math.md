---
title: Math & LaTeX
description: Kiln supports mathematical notation using MathJax, allowing you to write complex equations with standard LaTeX syntax while maintaining optimal page performance.
---
# Math & LaTeX

Kiln includes full support for mathematical notation, allowing you to render complex equations directly in your notes using **LaTeX** syntax.

It is designed to be fully compatible with Obsidian's math support, so your existing equations will render perfectly without modification.

## Usage

Kiln uses the standard `$` delimiter syntax to distinguish math from regular text.

### Block Equations
Use double dollar signs `$$` to create a centered block equation. This is best for complex formulas that need their own line.

````markdown
$$
\int_0^\infty e^{-x^2} dx = \frac{\sqrt{\pi}}{2}
$$
````

**Result:**
$$
\int_0^\infty e^{-x^2} dx = \frac{\sqrt{\pi}}{2}
$$

### Inline Equations

Use a single dollar sign `$` to embed math directly into a sentence. 

``` markdown
The mass-energy equivalence is described by the famous equation $E=mc^2$.
```

**Result:** The mass-energy equivalence is described by the famous equation $E=mc^2$.

## Performance

Rendering high-quality math on the web requires powerful libraries, which can often slow down websites. Kiln solves this with a **Lazy Loading** architecture.

- **Smart Detection:** Kiln scans your page content. If a page does not contain any math symbols (`$` or `$$`), the math engine is not loaded.
- **On-Demand Loading:** The **[MathJax](https://www.mathjax.org/)** library is only downloaded and executed when a user visits a page that actually requires it.

This ensures that your blog posts and simple text notes remain lightning fast, while your technical documentation retains full rendering capabilities.