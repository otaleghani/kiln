---
title: Client Side Navigation
description: Learn how Kiln uses HTMX hx-boost to provide instant, single-page application (SPA) navigation while maintaining perfect SEO and progressive enhancement.
---
# Client Side Navigation

Kiln transforms your static site into a high-performance **Single Page Application (SPA)** using [HTMX](https://htmx.org/) and its powerful [`hx-boost`](https://htmx.org/attributes/hx-boost/) feature.

This allows your site to navigate as fast as a native app without the complexity of client-side routers or heavy JavaScript frameworks.

## Why use Client Side Navigation?

In a standard website, clicking a link forces the browser to destroy the current page, download the new one, and re-parse every script and stylesheet from scratch. This process often causes a **Flash of Unstyled Content (FOUC)**—a brief moment where the screen goes white or layout jumps before the page settles.

Kiln's boosted navigation solves this:

* **Eliminates FOUC:** Transitions are seamless because the browser never "leaves" the page context.
* **Persistent Layouts:** Your sidebar scroll position and header state remain intact while only the main content updates.
* **Reduced Bandwidth:** Global assets (like fonts, CSS, and heavy scripts) are downloaded only once, saving data and speeding up interactions.

## How `hx-boost` Works

When `hx-boost` is active, Kiln acts as a lightweight browser controller:

1.  **Interception:** When a user clicks an internal link (e.g., `<a href="/about">`), HTMX intercepts the event and prevents the default full-page load.
2.  **AJAX Fetch:** It silently requests the new HTML page in the background via AJAX.
3.  **Smart Swap:** Upon receiving the response, it extracts the `<body>` content and seamlessly swaps it into your current container.
4.  **History Update:** It updates the browser's URL bar and pushes the new state to the History API, ensuring the "Back" and "Forward" buttons work exactly as users expect.

### Progressive Enhancement

Kiln follows the philosophy of **[Progressive Enhancement](https://developer.mozilla.org/en-US/docs/Glossary/Progressive_Enhancement)**. If a user visits your site with JavaScript disabled (or if the script fails to load), `hx-boost` simply steps aside. Your links function as standard HTML anchors, ensuring your site remains fully accessible and robust under any conditions.

## Developer Notes

Because the page does not undergo a full reload during navigation, the standard `DOMContentLoaded` event will **not** fire when users navigate between pages.

If you are adding custom JavaScript, you must listen for HTMX-specific events to know when new content has loaded, like so:

```javascript
// Run on initial load
document.addEventListener('DOMContentLoaded', initMyScript);

// Run after every HTMX navigation (page swap)
document.addEventListener('htmx:afterSwap', initMyScript);

function initMyScript() {
    console.log("Page content initialized!");
    // Your logic here...
}