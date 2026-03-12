#!/bin/sh
set -e

REPO="hardhackerlabs/podwise-cli"
FILE_NAME="podwise-skills.tar.gz"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color

echo "==> Installing Podwise skills from $REPO..."

# Fetch latest version if not specified
if [ -z "$VERSION" ]; then
    echo "==> Fetching latest release version..."
    VERSION=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    if [ -z "$VERSION" ]; then
        echo "${RED}Error: Could not determine the latest version from GitHub.${NC}"
        exit 1
    fi
fi

echo "==> Version: $VERSION"

DOWNLOAD_URL="https://github.com/$REPO/releases/download/$VERSION/$FILE_NAME"

# Create a temporary directory
TMP_DIR=$(mktemp -d)
trap 'rm -rf "$TMP_DIR"' EXIT

echo "==> Downloading $FILE_NAME..."
curl -sL -o "$TMP_DIR/$FILE_NAME" "$DOWNLOAD_URL"

echo "==> Extracting to current directory..."
tar xzf "$TMP_DIR/$FILE_NAME" -C .

echo "${GREEN}==> Skills installed successfully!${NC}"
