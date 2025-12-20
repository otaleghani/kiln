---
title: Local Graph
description: Kiln generates a context-aware Local Graph for every page, visualizing the immediate connections, backlinks, and related topics in the sidebar.
---

# Local Graph

While the [[Global Graph]] provides a view of your entire knowledge base, the **Local Graph** focuses on context.

Located at the top of the **Right Sidebar**, this feature generates a dynamic network visualization specific to the page you are currently reading.

## Why it is useful

In a dense documentation site or digital garden, notes rarely exist in isolation. The Local Graph helps you answer:
* "What other concepts link to this page?" (**Backlinks**)
* "Where does this page lead next?" (**Outgoing Links**)

It allows users to visualize the immediate "neighborhood" of a topic without getting lost in the noise of the entire vault.

## Behavior

The Local Graph updates automatically as you navigate.

1.  **Center Node:** The current page is always the central node, highlighted for clarity.
2.  **Neighbors:** Any note that directly links to *or* is linked from the current page is displayed as a connected node.
3.  **Interactivity:** Just like the global graph, you can hover over nodes to see titles and click them to navigate directly to that page.

*Note: On smaller screens (mobile devices), the right sidebar—and thus the Local Graph—is hidden to preserve reading space.*