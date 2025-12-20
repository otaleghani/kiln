---
title: Light & Dark Mode
description: Kiln features a built-in theme toggle that supports Light and Dark modes, respects system preferences, and persists your choice across sessions.
---
# Light & Dark Mode

Kiln comes equipped with a fully responsive theming engine that supports both **Light** and **Dark** visual modes.

A toggle switch is located at the top of the left sidebar (represented by a Sun/Moon icon), allowing users to instantly switch between themes to suit their environment or reading preference.

## Smart Features

### System Sync
By default, Kiln respects the user's operating system settings. If a visitor has their computer set to "Dark Mode," your site will automatically load in Dark Mode to match their expectation.

### Persistence
Once a user manually toggles the theme, Kiln remembers their choice. Using local storage, the site will maintain their preferred theme across page reloads and future visits.

### Deep Integration
The theme switch isn't just a cosmetic overlay; it integrates deeply with all Kiln features. When you toggle the mode, the following elements automatically adapt in real-time:

* **Syntax Highlighting:** Code blocks switch to optimized high-contrast colors.
* **Mermaid Graphs:** Diagrams re-render with appropriate stroke and fill colors.
* **Canvas:** The infinite canvas background and nodes adjust to reduce eye strain.

## Prevention of "Flash"
Kiln includes a lightweight script in the `<head>` of every page that runs before the content loads. This prevents the "Flash of Incorrect Theme" (FOIT)—where a dark mode user is briefly blinded by a white page before the styles load.