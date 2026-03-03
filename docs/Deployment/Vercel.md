---
title: Deploy Kiln on Vercel — Static Site Hosting Guide
description: Deploy your Obsidian vault as a static site on Vercel. Step-by-step build script setup, environment variables, custom domains, and troubleshooting.
---
# Deploy on Vercel

Vercel is a popular hosting platform for static sites and frontend frameworks. Although Vercel is optimized for Node.js frameworks like Next.js, it works perfectly with static site generators like Kiln. You create a small build script that downloads the Kiln binary and generates your site during each deployment.

This guide walks you through connecting your GitHub repository to Vercel so that every push automatically rebuilds and publishes your Obsidian vault as a website.

## Prerequisites

1. A GitHub repository containing your Obsidian vault.
2. A [Vercel account](https://vercel.com/) (the free Hobby plan works fine for most documentation sites).
3. **No binary needed:** You do *not* need to commit the `kiln` binary to your repository. The build script below downloads it automatically.

## Step 1: Add a Build Script

Add a file named `build.sh` to the root of your repository. This script runs on Vercel's servers during each deployment to download Kiln and generate your static site.

```bash
#!/bin/bash
set -e

# 1. Download Kiln (Linux AMD64 binary for Vercel's build environment)
curl -L -o kiln https://github.com/otaleghani/kiln/releases/latest/download/kiln_linux_amd64

# 2. Make it executable
chmod +x kiln

# 3. Generate the site
./kiln generate \
  --url "$VERCEL_URL_FULL" \
  --name "My Digital Garden" \
  --input "." \
  --output "./public"
```

### Customize the Build Flags

The [[Generate Command]] accepts several flags to control the output. The most important ones for deployment are:

| **Flag** | **Example** | **Purpose** |
|---|---|---|
| `--url` | `https://my-site.vercel.app` | Sets the base URL for your [[Sitemap xml\|sitemap.xml]], [[Robots txt\|robots.txt]], and [[Canonical tag\|canonical tags]] |
| `--name` | `"My Digital Garden"` | Sets the site name in [[Meta Tags\|meta tags]] and the navigation bar |
| `--input` | `./vault` | Path to the folder containing your Markdown notes (defaults to `./vault`) |
| `--output` | `./public` | Path where generated HTML files are saved (defaults to `./public`) |
| `--theme` | `dracula` | Color theme for your site — see [[Themes]] for all options |
| `--font` | `merriweather` | Typography family — see [[Fonts]] for all options |

## Step 2: Configure Environment Variables

Instead of hardcoding your production URL in the build script, you can use a Vercel environment variable. This keeps your script portable and makes it easy to use different URLs for preview and production deployments.

1. In your Vercel project dashboard, go to **Settings** > **Environment Variables**.
2. Add a new variable:

| **Name** | **Value** | **Environment** |
|---|---|---|
| `VERCEL_URL_FULL` | `https://your-project.vercel.app` | Production |

3. Update your `build.sh` to reference it (already done in the script above):
```bash
./kiln generate --url "$VERCEL_URL_FULL"
```

If you prefer simplicity, you can hardcode the URL directly in the script instead:
```bash
./kiln generate --url "https://your-project.vercel.app"
```

## Step 3: Configure Vercel

1. Log in to [Vercel](https://vercel.com/) and click **"Add New..."** > **"Project"**.
2. Import your GitHub repository.
3. In the **Configure Project** screen, expand the **"Build and Output Settings"** section.
4. Toggle **Override** for the following settings:

| **Setting** | **Value** |
|---|---|
| **Build Command** | `bash build.sh` |
| **Output Directory** | `public` |

5. Click **Deploy**.

Vercel will clone your repository, run `build.sh` (which downloads Kiln and generates the HTML), and publish the contents of the `public` folder to its global edge network.

## Custom Domains

Vercel assigns a `.vercel.app` subdomain to every project by default. To use your own domain:

1. Go to your project **Settings** > **Domains**.
2. Add your custom domain and follow the DNS configuration instructions.
3. Update the `VERCEL_URL_FULL` environment variable (or the hardcoded `--url` flag) to match your custom domain so that your [[Sitemap xml|sitemap.xml]] and [[Canonical tag|canonical links]] point to the correct address.

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

Vercel caches deployments aggressively. After updating your vault, push your changes to GitHub — Vercel will automatically trigger a new build. You can also trigger a manual redeploy from the Vercel dashboard under **Deployments** > **Redeploy**.

## Other Hosting Options

Kiln generates standard HTML, CSS, and JavaScript that works on any static hosting platform. See the deployment guides for [Cloudflare Pages](./Cloudflare%20Pages.md), [GitHub Pages](./GitHub%20Pages.md), [Netlify](./Netlify.md), or self-hosted [Web Servers](./Web%20Servers.md) like Nginx, Apache, and Caddy.
