---
title: Install Kiln — Obsidian Static Site Generator for macOS, Linux, Windows
description: Install Kiln in under a minute. Download a single binary or use go install to turn your Obsidian vault into a static website on macOS, Linux, or Windows.
---
# Installation

**Kiln** is distributed as a single binary with zero dependencies — no Node.js, no Ruby, no Docker required. Choose the method that fits your workflow to get started on macOS, Linux, or Windows.

## Install with Go (Recommended)

If you have Go 1.25 or later installed, this is the fastest way to install Kiln:

```bash
go install github.com/otaleghani/kiln/cmd/kiln@latest
```

This downloads, compiles, and places the `kiln` binary in your `$GOPATH/bin` directory.

## macOS (Apple Silicon / ARM64)

### Download the binary
Download the latest release for Apple Silicon:
```bash
curl -LO https://github.com/otaleghani/kiln/releases/latest/download/kiln_darwin_arm64
```

### Verify the Checksum (Recommended)
Ensure the file was downloaded correctly and has not been tampered with:
```bash
curl -LO https://github.com/otaleghani/kiln/releases/latest/download/checksums.txt

sha256sum -c checksums.txt --ignore-missing
```
_You should see `kiln_darwin_arm64: OK`._

### Install
Make the binary executable and move it to a directory in your `PATH`:
```bash
chmod +x kiln_darwin_arm64
sudo mv kiln_darwin_arm64 /usr/local/bin/kiln
```

### Allow Execution (First Run Only)
Since this binary is not notarized by Apple, you may need to allow it to run. Go to **System Settings** > **Privacy & Security**, scroll down, and click **Allow Anyway** next to the notification about `kiln`. Alternatively, remove the quarantine attribute via terminal:
```bash
xattr -d com.apple.quarantine /usr/local/bin/kiln
```

## Linux (AMD64)

### Download the binary
```bash
curl -LO https://github.com/otaleghani/kiln/releases/latest/download/kiln_linux_amd64
```

### Verify the Checksum (Recommended)
```bash
curl -LO https://github.com/otaleghani/kiln/releases/latest/download/checksums.txt

sha256sum -c checksums.txt --ignore-missing
```
_You should see `kiln_linux_amd64: OK`._

### Install
Make the binary executable and move it to `/usr/local/bin`:
```bash
chmod +x kiln_linux_amd64
sudo mv kiln_linux_amd64 /usr/local/bin/kiln
```

## Windows (AMD64)

### Download the binary
Download `kiln_windows_amd64.exe` from the [Releases Page](https://github.com/otaleghani/kiln/releases/latest) or via PowerShell:
```powershell
Invoke-WebRequest -Uri "https://github.com/otaleghani/kiln/releases/latest/download/kiln_windows_amd64.exe" -OutFile "kiln.exe"
```

### Verify the Checksum (Recommended)
Run the following in PowerShell to verify the hash matches:
```powershell
Invoke-WebRequest -Uri "https://github.com/otaleghani/kiln/releases/latest/download/checksums.txt" -OutFile "checksums.txt"
$expected = Select-String -Path .\checksums.txt -Pattern "kiln_windows_amd64.exe" | ForEach-Object { $_.Line.Split(' ')[0] };
(Get-FileHash .\kiln_windows_amd64.exe -Algorithm SHA256).Hash.ToLower() -eq $expected
```
_This should return `True`._

### Install
Move `kiln.exe` to a folder of your choice (e.g., `C:\Program Files\Kiln\`) and add that folder to your System `PATH` environment variable so you can run `kiln` from any terminal window.

## Verify the Installation

After installing, confirm Kiln is available by checking the version:

```bash
kiln version
```

## Quick Start: Generate Your First Site

Once installed, you can turn your Obsidian vault into a website with two commands:

```bash
kiln generate --input ./my-vault --output ./public
kiln serve --output ./public
```

Open `http://localhost:8080` to preview your site locally. The [Generate Command](./Commands/generate.md) accepts flags for [themes](./Features/User Interface/Themes.md), fonts, site name, and base URL — run `kiln generate --help` to see all options.

When you are ready to publish, check the deployment guides for [Cloudflare Pages](./Deployment/Cloudflare Pages.md), [GitHub Pages](./Deployment/GitHub Pages.md), [Netlify](./Deployment/Netlify.md), or [Vercel](./Deployment/Vercel.md).

## Troubleshooting

If something doesn't look right after generating your site, run the [Doctor Command](./Commands/doctor.md) to scan your vault for broken links and common issues:

```bash
kiln doctor --input ./my-vault
```
