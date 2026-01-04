---
title: Deploy on Netlify
description: How to deploy your Kiln site to Netlify using a custom build script.
---
# Deploy on Netlify

Netlify is a popular platform for Continuous Deployment. You connect your GitHub repository, and Netlify builds your site every time you push a commit.

Since Kiln is a standalone binary (and not a standard Node.js package), we need to tell Netlify how to download and run it.

## Step 1: Add a Build Script
Add a file named `build.sh` to the root of your repository. This script downloads Kiln and runs the generator.

```bash
#!/bin/bash

# 1. Download Kiln (Replace URL with your actual Linux binary release URL)
curl -L -o kiln [https://github.com/otaleghani/kiln/releases/latest/download/kiln_linux_amd64](https://github.com/otaleghani/kiln/releases/latest/download/kiln_linux_amd64)

# 2. Make it executable
chmod +x kiln

# 3. Generate the site
# Remember to set the --url to your Netlify subdomain (e.g., [https://my-site.netlify.app](https://my-site.netlify.app))
./kiln generate --url $URL
```

_Note: Netlify automatically provides the `$URL` environment variable during build, so the script above grabs the correct URL automatically!_

## Step 2: Configure Netlify

1. Log in to Netlify and click **"Add new site"** -> **"Import an existing project"**.
2. Select your GitHub repository.
3. In the **Build settings** screen, configure the following:

|**Setting**|**Value**|
|---|---|
|**Build Command**|`bash build.sh`|
|**Publish Directory**|`public`|

4. Click **Deploy site**.
