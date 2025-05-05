#!/usr/bin/env bash
set -euo pipefail

ARCH=$(dpkg --print-architecture)
REPO="Zweih/qp"
LATEST_RELEASE_TAG=$(curl -s https://api.github.com/repos/$REPO/releases/latest | grep -oP '"tag_name": "\K(.*)(?=")')
DEB_NAME="qp-${LATEST_RELEASE_TAG}-${ARCH}.deb"
URL="https://github.com/$REPO/releases/download/${LATEST_RELEASE_TAG}/${DEB_NAME}"

echo "Downloading $DEB_NAME..."
curl -L -o "$DEB_NAME" "$URL"

echo "Installing $DEB_NAME..."
sudo dpkg -i "$DEB_NAME" || sudo apt-get install -f -y

echo "Cleaning up..."
rm "$DEB_NAME"

echo "qp installed successfully!"
