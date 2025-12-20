---
title: Global Graph
description: Visualize the connections between your notes with Kiln's interactive Global Graph, a force-directed network view of your entire vault.
---
# Global Graph

Kiln automatically generates a **Global Graph** view for your website. This is an interactive visualization of your entire knowledge base, modeled after the popular graph view in Obsidian.

This feature helps you and your visitors understand the structure of your notes, identifying clusters of related topics and how different ideas connect.

## Accessing the Graph

The graph is generated as a standalone page located at `/graph`.

If your site is hosted at `https://example.com`, your graph will be available at:
```text
https://example.com/graph
```

### Adding a Link

The right sidebar has a link to the graph page. If you wish to add another link you can manually add a link to it in your homepage or any other note:
```markdown
Check out the [Map of my Brain](/graph)!
```

## How it Works

The Global Graph uses a force-directed algorithm to organize your notes visually:

- **Nodes (Circles):** Each node represents a single page or note in your vault.
- **Edges (Lines):** A line is drawn between two nodes whenever one note links to another using a [[Wikilinks|Wikilink]].
- **Interactivity:** The graph is fully interactive. Users can zoom in, pan around, and hover over nodes to see the note titles. Clicking a node will navigate directly to that page.

## Use Cases

- **Discovery:** Users can find connections between topics they might not have realized were related.
- **Navigation:** It serves as a visual alternative to the sidebar for browsing content.
- **Overview:** It provides a "bird's-eye view" of your work, highlighting which areas of your documentation are the most dense and interconnected.
- **Looks cool:** Okay? Okay.