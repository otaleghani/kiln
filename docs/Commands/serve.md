---
title: "Serve Command — Preview Your Site Locally"
description: "Use kiln serve to launch a local development server and preview your generated Obsidian vault site in the browser with clean URLs and custom 404 pages."
---

# Serve Command

The `serve` command starts a local HTTP server so you can preview your generated site in the browser before deploying it.

Because Kiln produces a static site, opening the HTML files directly from disk will not work — features like [Client Side Navigation](../Features/Navigation/Client Side Navigation.md) and AJAX requests fail due to browser security restrictions. The `serve` command provides a proper web server environment to test your site locally.

## Usage

Run `serve` after building your site with the [Generate Command](./generate.md):

```bash
# Build the site
kiln generate

# Preview the site
kiln serve
```

Once running, open `http://localhost:8080` in your browser to view your site.

## Flags

| Flag       | Short | Default    | Description                                                                                     |
| ---------- | ----- | ---------- | ----------------------------------------------------------------------------------------------- |
| `--port`   | `-p`  | `8080`     | Port number to listen on. Change this if port 8080 is already in use.                           |
| `--output` | `-o`  | `./public` | Directory to serve. Should match the output path used during the [generate](./generate.md) step.|
| `--log`    | `-l`  | `info`     | Sets the log level. Choose between `info` or `debug`.                                           |

## Clean URL Support

The development server automatically handles clean URLs, matching the behavior of production static hosts. When a browser requests a path without a file extension (for example `/my-note`), the server resolves it by looking for:

1. **An HTML file** with the same name — `/my-note` serves `my-note.html`.
2. **A directory** with an `index.html` — `/folder` serves `folder/index.html`.

Trailing slashes are canonicalized with a redirect. A request to `/about/` becomes `/about`, keeping your URLs consistent.

## Custom 404 Pages

If your vault contains a `404.md` note, Kiln generates a `404.html` file during the build step. The serve command detects this file and returns it whenever a visitor hits a missing page, so you can test your custom error page locally before deploying.

## Base Path Handling

When you set a `--url` flag with a path during the build (for example `https://example.com/docs`), the server automatically mounts your site under that same path prefix locally. This lets you verify that all assets and links resolve correctly under a subdirectory, exactly as they will in production. With the example above, your local preview would be at `http://localhost:8080/docs/`.

## Examples

Preview on a custom port:

```bash
kiln serve --port 3000
```

Serve from a different output directory:

```bash
kiln serve --output ./dist
```

Full workflow from build to preview:

```bash
kiln doctor --input ./vault
kiln generate --name "My Notes" --url "https://notes.example.com" --theme nord
kiln serve
```

## Production Deployment

The `serve` command is designed for local development and previewing — it is **not recommended** for production traffic.

To deploy your site to the web, upload the contents of your output directory to a dedicated static host like [GitHub Pages](../Deployment/GitHub Pages.md), [Netlify](../Deployment/Netlify.md), [Vercel](../Deployment/Vercel.md), or [Cloudflare Pages](../Deployment/Cloudflare Pages.md). You can also serve the files with a production-grade web server such as [Nginx, Caddy, or Apache](../Deployment/Web Servers.md).

## Related Commands

- [Generate Command](./generate.md) — build your vault into a static site
- [Dev Command](./dev.md) — build, watch, and serve in a single step
- [Clean Command](./clean.md) — remove build output before rebuilding
- [Doctor Command](./doctor.md) — scan for broken wikilinks before deploying
