---
title: Serve Command
description: The serve command launches a lightweight local web server, allowing you to preview your generated site in the browser before deploying it.
---

# Serve Command

The `serve` command spins up a lightweight, local HTTP server to preview your website.

Because Kiln generates a static site, you cannot simply open the HTML files directly in your browser (features like Client Side Navigation and AJAX requests will fail due to security protocols). You must view them through a web server. This command provides that environment instantly.

## Usage

This command is typically the final step in your local workflow:

```bash
# Build the site
./kiln generate

# Preview the site
./kiln serve
```

Once running, you can view your site by opening `http://localhost:8080` in your web browser.

## Flags

| Flag       | Short version | Default value | Description                                                                                         |
| ---------- | ------------- | ------------- | --------------------------------------------------------------------------------------------------- |
| `--port`   | `-p`          | `"8080"`      | The port number to listen on. Useful if port 8080 is already in use by another application.         |
| `--output` | `-o`          | `"./public"`  | The directory to serve. This should match the output folder you specified during the generate step. |
| `--log`    | `-l`          | `info`        | Sets the log level. You can choose between `info` or `debug`.                                       |

## Production Note

While this command is perfect for local development and previewing, it is **not recommended** for high-traffic production environments.

For deploying to the live web, it is recommended to upload the contents of your `public` folder to a dedicated static host (like [[GitHub Pages]], [[Netlify]], or [[Deployment/Vercel]]) or serve them using a production-grade server like [[Web Servers#Nginx|Nginx]], [[Web Servers#Caddy|Caddy]] or [[Web Servers#Apache|Apache]].