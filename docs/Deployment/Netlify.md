---
title: Deploy Kiln on Netlify — Continuous Deployment for Obsidian Sites
description: Deploy your Obsidian vault as a static site on Netlify. Step-by-step build script setup, environment variables, custom domains, and troubleshooting.
---
# Deploy on Netlify

Netlify is a popular platform for hosting static sites with continuous deployment. You connect your GitHub repository, and Netlify automatically builds and publishes your site every time you push a commit. Since Kiln is a standalone binary (not a standard Node.js package), you create a small build script that downloads and runs Kiln during each deployment.

This guide walks you through connecting your GitHub repository to Netlify so that every push automatically rebuilds and publishes your Obsidian vault as a website.

## Prerequisites

1. A GitHub repository containing your Obsidian vault.
2. A [Netlify account](https://www.netlify.com/) (the free Starter plan works fine for most documentation sites).
3. **No binary needed:** You do *not* need to commit the `kiln` binary to your repository. The build script below downloads it automatically.

## Step 1: Add a Build Script

Add a file named `build.sh` to the root of your repository. This script runs on Netlify's servers during each deployment to download Kiln and generate your static site.

```bash
#!/bin/bash
set -e

# 1. Download Kiln (Linux AMD64 binary for Netlify's build environment)
curl -L -o kiln https://github.com/otaleghani/kiln/releases/latest/download/kiln_linux_amd64

# 2. Make it executable
chmod +x kiln

# 3. Generate the site
./kiln generate \
  --url "$URL" \
  --name "My Digital Garden" \
  --input "." \
  --output "./public"
```

Netlify automatically provides the `$URL` environment variable during builds, so the script picks up the correct deployment URL without any extra configuration.

### Customize the Build Flags

The [Generate Command](../Commands/generate.md) accepts several flags to control the output. The most important ones for deployment are:

| **Flag** | **Example** | **Purpose** |
|---|---|---|
| `--url` | `https://my-site.netlify.app` | Sets the base URL for your [sitemap.xml](../Features/SEO/Sitemap xml.md), [robots.txt](../Features/SEO/Robots txt.md), and [canonical tags](../Features/SEO/Canonical.md) |
| `--name` | `"My Digital Garden"` | Sets the site name in [meta tags](../Features/SEO/Meta Tags.md) and the navigation bar |
| `--input` | `./vault` | Path to the folder containing your Markdown notes (defaults to `./vault`) |
| `--output` | `./public` | Path where generated HTML files are saved (defaults to `./public`) |
| `--theme` | `dracula` | Color theme for your site — see [Themes & Visuals](../Features/User Interface/Themes.md) for all options |
| `--font` | `merriweather` | Typography family — see [Fonts & Typography](../Features/User Interface/Fonts.md) for all options |

## Step 2: Configure Environment Variables

Instead of hardcoding your production URL in the build script, you can use a Netlify environment variable. This keeps your script portable and makes it easy to use different URLs for preview and production deployments.

1. In your Netlify site dashboard, go to **Site configuration** > **Environment variables**.
2. Add a new variable:

| **Name** | **Value** | **Scope** |
|---|---|---|
| `URL` | `https://your-site.netlify.app` | All deploys |

3. Update your `build.sh` to reference it (already done in the script above):
```bash
./kiln generate --url "$URL"
```

If you prefer simplicity, you can hardcode the URL directly in the script instead:
```bash
./kiln generate --url "https://your-site.netlify.app"
```

## Step 3: Configure Netlify

1. Log in to [Netlify](https://www.netlify.com/) and click **"Add new site"** > **"Import an existing project"**.
2. Select your GitHub repository.
3. In the **Build settings** screen, configure the following:

| **Setting** | **Value** |
|---|---|
| **Build Command** | `bash build.sh` |
| **Publish Directory** | `public` |

4. Click **Deploy site**.

Netlify will clone your repository, run `build.sh` (which downloads Kiln and generates the HTML), and publish the contents of the `public` folder to its global CDN.

## Custom Domains

Netlify assigns a `.netlify.app` subdomain to every site by default. To use your own domain:

1. Go to your site **Domain management** > **Add a domain**.
2. Add your custom domain and follow the DNS configuration instructions.
3. Update the `URL` environment variable (or the hardcoded `--url` flag) to match your custom domain so that your [sitemap.xml](../Features/SEO/Sitemap xml.md) and [canonical links](../Features/SEO/Canonical.md) point to the correct address.

## Troubleshooting

### Build Fails with "Permission Denied"

Make sure the `build.sh` file is executable. You can fix this locally and commit:
```bash
chmod +x build.sh
git add build.sh
git commit -m "Make build script executable"
```

### Broken CSS or Links

If your styles or internal links are broken after deployment, the `--url` flag is likely incorrect or missing. Kiln uses this value to generate absolute paths for assets and SEO files. Double-check that it matches your actual deployment URL (including `https://`).

### Site Shows Stale Content

Netlify caches deployments at the CDN level. After updating your vault, push your changes to GitHub — Netlify will automatically trigger a new build. You can also trigger a manual redeploy from the Netlify dashboard under **Deploys** > **Trigger deploy**.

## Other Hosting Options

Kiln generates standard HTML, CSS, and JavaScript that works on any static hosting platform. See the deployment guides for [Cloudflare Pages](./Cloudflare%20Pages.md), [GitHub Pages](./GitHub%20Pages.md), [Vercel](./Vercel.md), or self-hosted [Web Servers](./Web%20Servers.md) like Nginx, Apache, and Caddy.
