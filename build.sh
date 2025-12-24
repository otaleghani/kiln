#!/bin/bash

INPUT_DIR="./docs"
DEPLOYMENT_URL="https://kiln.talesign.com"

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
# Adjust the flags as needed for your site
echo "Building site..."
./kiln generate \
  --input "$INPUT_DIR" \
  --output ./public \
  --flat-urls=true \
  --name "Kiln" \
  --url "$DEPLOYMENT_URL"

# Note: CF_PAGES_URL is a magic variable provided automatically by Cloudflare

echo "Kiln build complete successfully"
