---
title: Development Environment
description: Internal guide for setting up the Kiln development environment. Instructions for using Nix, Air for hot-reloading, and Tailwind CSS for layout development.
---
# Development Environment

This section outlines the steps to set up the development environment for Kiln. It is intended for contributors or those looking to modify the core binary and default layouts.

## Prerequisites & Setup
Kiln uses **Nix** to manage development dependencies and **npm** for code formatting and styling tools.
```bash
# 1. Enter the Nix shell to install Go, Air, and other system deps
nix develop

# 2. Install Node packages (primarily for Prettier plugins)
npm install
```

## Hot Reloading (Air)
We use **Air** for live reloading during Go development. This will watch your `.go` files and automatically rebuild/restart the binary when changes are detected.

```bash
# Start the auto-reloading dev environment
air
```

## Styling & Layouts (Tailwind)

If you are modifying the built-in layouts (e.g., `default` or `simple`), you need to run the Tailwind CLI to regenerate the CSS.

> [!note] 
> Pay attention to the output flag (`-o`). You must point it to the CSS file corresponding to the specific layout you are working on.

```bash
# Watch mode for Tailwind
# Replace 'test_style.css' with the actual CSS file for your layout
tailwindcss -i ./assets/input.css -o ./assets/test_style.css --watch
```


## Configuration Files
Key configuration files found in the repository root:

- **`.prettierrc`**: Configuration for Prettier (used for formatting HTML/JS/CSS).
- **`.air.toml`**: Configuration for Air (build commands, watch exclusions, and binary locations).