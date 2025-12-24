---
title: Deploy on Cloudflare Pages
description: A step-by-step guide to automatically building and hosting your Kiln site using Cloudflare Pages.
---
# Deploy on Cloudflare Pages

[Cloudflare Pages](https://pages.cloudflare.com/) is an excellent hosting choice for Kiln sites due to its speed and simple "git-push-to-deploy" workflow.

Instead of committing the Kiln binary to your repository, we use a simple build script. This script runs automatically on Cloudflare's servers to download the latest version of Kiln, generate your site, and publish it.

## Prerequisites

1. A GitHub repository containing your Obsidian vault.
2. A Cloudflare account.

## Step 1: Add the Build Script

Create a new file in the root of your repository named `build.sh`.

Paste the following script into it. This script automatically fetches the latest version of Kiln and builds your site.
```bash
#!/bin/bash

# --- CONFIGURATION START ---
SITE_NAME="Kiln"
INPUT_DIR="./docs"
DEPLOYMENT_URL="https://kiln.talesign.com"
# --- CONFIGURATION END ---

# Exit immediately if any command fails
set -e

echo "Kiln build script"

# Find the latest version of Kiln
LATEST_TAG=$(curl -s https://api.github.com/repos/otaleghani/kiln/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST_TAG" ]; then
  echo "Error: Could not determine latest Kiln version."
  exit 1
fi

echo "Detected latest version: $LATEST_TAG"

# Download the binary
URL="https://github.com/otaleghani/kiln/releases/download/${LATEST_TAG}/kiln_linux_amd64"

echo "Downloading binary from $URL..."
curl -L -o ./kiln "$URL"
chmod +x ./kiln

# Run the build
echo "Building site..."
./kiln generate \
  --input "$INPUT_DIR" \
  --output ./public \
  --flat-urls=true \
  --name "$SITE_NAME" \
  --url "$DEPLOYMENT_URL"

echo "Kiln build complete successfully"
```

## Step 2: Critical Configuration

There are two variables at the top of the `build.sh` file that you **must** customize for your specific project:

### The Input Directory (`INPUT_DIR`)

Change `INPUT_DIR` to point to the folder containing your markdown notes.

- **Example:** If your notes are in the root of the repo, use `"."`.
- **Example:** If your notes are in a folder named `content`, use `"./content"`.

### The Deployment URL (`DEPLOYMENT_URL`)

Change `DEPLOYMENT_URL` to the actual address where your site will be hosted.

- **Why?** Kiln uses this to generate the [[Sitemap xml|sitemap.xml]], [[Robots txt|robots.txt]], and canonical meta tags.
- **Tip:** Cloudflare provides a `*.pages.dev` subdomain (e.g., `https://my-site.pages.dev`), or you can use your own custom domain.

### The Website Name (`SITE_NAME`)

Change `SITE_NAME` to change the navbar name and site meta tags.

## Step 3: Configure Cloudflare Pages

1. Push the `build.sh` file to your GitHub repository.
2. Log in to the **Cloudflare Dashboard** and go to **Compute (Workers & Pages)**.
3. Click **Create application** > **Pages** > **Connect to Git**.
4. Select your repository.
5. In the **Build settings** section, configure the following:

|**Setting**|**Value**|
|---|---|
|**Framework preset**|`None`|
|**Build command**|`bash ./build.sh`|
|**Build output directory**|`public`|

6. Click **Save and Deploy**.

Cloudflare will now clone your repository, run your script (which downloads Kiln), and deploy the resulting `./public` folder to the global edge network.