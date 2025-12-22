---
title: Deploy on Vercel
description: How to deploy your Kiln site to Vercel using a custom build script.
---
# Deploy on Vercel

Vercel is optimized for frontend frameworks, but it works perfectly with static generators like Kiln. Similar to [[Netlify]], we will use a small script to handle the build process.

## Step 1: Add a Build Script
Add a file named `build.sh` to the root of your repository.

```bash
#!/bin/bash

# 1. Download Kiln (Replace URL with your actual Linux binary release URL)
curl -L -o kiln [https://github.com/otaleghani/kiln/releases/latest/download/kiln_linux_amd64](https://github.com/otaleghani/kiln/releases/latest/download/kiln_linux_amd64)

# 2. Make it executable
chmod +x kiln

# 3. Generate the site
# You can hardcode your URL here, or set it as an Environment Variable in Vercel settings
./kiln generate --url "https://your-project.vercel.app"
```

## Step 2: Configure Vercel

1. Log in to Vercel and click **"Add New..."** -> **"Project"**.
2. Import your GitHub repository.
3. In the **Configure Project** screen, expand the **"Build and Output Settings"** section.
4. Toggle **Override** for the following settings:

|**Setting**|**Value**|
|---|---|
|**Build Command**|`bash build.sh`|
|**Output Directory**|`public`|

5. Click **Deploy**.