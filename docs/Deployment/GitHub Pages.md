---
title: Deploy on GitHub Pages
description: A step-by-step guide to automatically building and hosting your Kiln site using GitHub Actions.
---
# Deploy on GitHub Pages

The most robust way to deploy Kiln is using **GitHub Actions**. This method allows you to keep your repository clean (just your notes) and have GitHub's servers automatically download Kiln, build your site, and publish it whenever you push changes.

## Prerequisites

1.  A GitHub repository containing your Obsidian vault.
2.  **No binary needed:** You do *not* need to commit the `kiln` file to your repository. The workflow below will download it for you.

## Step 1: Configure the Workflow

Create a new file in your repository at: `.github/workflows/deploy.yml`.

Paste the following configuration. Note the new **"Setup Kiln"** step.

```yaml
name: Deploy to GitHub Pages

on:
  push:
    branches: ["main"]
  workflow_dispatch:

permissions:
  contents: read
  pages: write
  id-token: write

concurrency:
  group: "pages"
  cancel-in-progress: false

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v4

      - name: Setup Pages
        uses: actions/configure-pages@v5

      - name: Setup Kiln
        run: |
          # Download the latest Linux binary
          curl -L -o kiln https://github.com/otaleghani/kiln/releases/latest/download/kiln_linux_amd64
          
          # Make it executable
          chmod +x kiln

      - name: Build Site
        # Replace the URL below with your actual GitHub Pages URL
        run: ./kiln generate --url "https://username.github.io/repo"

      - name: Upload artifact
        uses: actions/upload-pages-artifact@v3
        with:
          path: ./public

  deploy:
    environment:
      name: github-pages
      url: ${{ steps.deployment.outputs.page_url }}
    runs-on: ubuntu-latest
    needs: build
    steps:
      - name: Deploy to GitHub Pages
        id: deployment
        uses: actions/deploy-pages@v4
```

## Step 2: Critical Configuration

There are two lines in the YAML above that you **must** customize for your specific project:

### 1. The Download URL

In the `Setup Kiln` step, you must point `curl` to the actual location of your Kiln binary release.

- **Target:** GitHub Actions runners use Linux. Ensure you link to the **Linux (AMD64)** version of the binary.
- **Format:** `https://github.com/otaleghani/kiln/releases/latest/download/kiln_linux_amd64`
- 

### 2. The Site URL

In the `Build Site` step, you must update the `--url` flag to match your specific GitHub Pages address.

- **Format:** `https://<username>.github.io/<repository-name>`
- **Why?** If this is incorrect, your CSS and links will break because GitHub hosts project sites in a sub-folder.

## Step 3: Enable Actions

1. Push the `.github/workflows/deploy.yml` file to your repository.
2. Go to your repository **Settings** > **Pages**.
3. Under **Build and deployment**, switch the source to **GitHub Actions**.

The next time you push a note to your repository, GitHub will automatically spin up a server, download Kiln, build your site, and deploy it.