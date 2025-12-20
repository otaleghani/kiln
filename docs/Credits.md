---
title: Credits & Acknowledgements
description: A list of the open-source software, libraries, and resources that make Kiln possible.
---

# Credits

Kiln is built on the shoulders of giants. It relies on a robust ecosystem of open-source tools to parse your notes, highlight your code, and render your diagrams.

## Core Technology

* **[Go](https://go.dev/):** Kiln is written in Go, chosen for its speed, concurrency model, and ability to compile into a single, dependency-free binary.

## Backend Libraries

These libraries run during the build process to transform your Markdown into HTML.

* **[Goldmark](https://github.com/yuin/goldmark):** A standards-compliant, extensible Markdown parser. It handles the heavy lifting of converting your notes into web pages.
* **[Chroma](https://github.com/alecthomas/chroma):** A general-purpose syntax highlighter. It powers the coloring of code blocks without relying on external CSS or Javascript.

## Frontend Libraries

These libraries are loaded by the user's browser to provide interactivity and visual rendering.

* **[Mermaid.js](https://mermaid.js.org/):** A Javascript-based diagramming and charting tool that renders Markdown definitions into dynamic visualizations.
* **[MathJax](https://www.mathjax.org/):** An open-source JavaScript display engine for LaTeX, enabling high-quality mathematical typesetting.

## Typography

Kiln embeds these open-source fonts directly into the build to ensure privacy and performance.

* **[Inter](https://rsms.me/inter/):** Designed by Rasmus Andersson. A typeface carefully crafted for computer screens.
* **[Merriweather](https://fonts.google.com/specimen/Merriweather):** Designed by Sorkin Type. A serif font designed to be highly readable on screens.
* **[Lato](https://fonts.google.com/specimen/Lato):** Designed by Łukasz Dziedzic. A humanist sans-serif font.


## Visual Themes

Kiln includes ports of these beautiful community-created color palettes.

* **[Catppuccin](https://catppuccin.com/):** A community-driven pastel theme that aims to be the middle ground between high and low contrast.
* **[Nord](https://www.nordtheme.com/):** An arctic, north-bluish color palette. Designed for a clean and clutter-free workflow.
* **[Dracula](https://draculatheme.com/):** A dark theme for many editors, shells, and more. Famous for its vibrant colors and high contrast.