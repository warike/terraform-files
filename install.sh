#!/bin/bash

# --- Configuration ---
PROJECT_NAME="wtf"
GITHUB_ORG_REPO="warike/terraform-files" # IMPORTANT: Must be "username/repo" or "org/repo"
RELEASE_TAG="v1.0.0" # <<-- SET YOUR DESIRED RELEASE TAG HERE (e.g., "v1.0.0", "v2.3.4")
INSTALL_DIR="/usr/local/bin" # Common installation directory for binaries

### --- ASCII Art ---
ascii_art='
░██       ░██                     ░██░██
░██       ░██                        ░██
░██  ░██  ░██  ░██████   ░██░████ ░██░██    ░██ ░███████
░██ ░████ ░██       ░██  ░███     ░██░██   ░██ ░██    ░██
░██░██ ░██░██  ░███████  ░██      ░██░███████  ░█████████
░████   ░████ ░██   ░██  ░██      ░██░██   ░██ ░██
░███     ░███  ░█████░██ ░██      ░██░██    ░██ ░███████
'

echo -e "\n$ascii_art\n"
echo "---"
echo " ✅ Starting installation for **$PROJECT_NAME**. 🧉 🚀"
echo "This installer supports **Apple Silicon (ARM64) Macs** and **Linux (AMD64)**."
echo "---"

# --- Validate OS and Architecture ---
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

TARGET_OS=""
TARGET_ARCH=""

case "$OS" in
    darwin)
        TARGET_OS="darwin"
        if [ "$ARCH" == "arm64" ]; then
            TARGET_ARCH="arm64"
        elif [ "$ARCH" == "x86_64" ]; then
            echo "❌ Error: Intel Macs (x86_64) are not supported by this installer."
            exit 1
        else
            echo "❌ Error: Unsupported architecture for macOS: $ARCH."
            exit 1
        fi
        echo "✅ Detected macOS (Apple Silicon - ARM64)."
        ;;
    linux)
        TARGET_OS="linux"
        if [ "$ARCH" == "x86_64" ]; then
            TARGET_ARCH="amd64"
            echo "✅ Detected Linux (AMD64)."
        elif [ "$ARCH" == "arm64" ] || [ "$ARCH" == "aarch64" ]; then
            echo "❌ Error: ARM-based Linux systems are not supported by this installer."
            exit 1
        else
            echo "❌ Error: Unsupported architecture for Linux: $ARCH."
            exit 1
        fi
        ;;
    *)
        echo "❌ Error: This script is only for macOS (Apple Silicon) or Linux (AMD64)."
        exit 1
        ;;
esac

# --- Construct Download URL ---
if [ -z "$TARGET_OS" ] || [ -z "$TARGET_ARCH" ]; then
    echo "❌ Critical Error: OS or Architecture not detected correctly. Exiting."
    exit 1
fi

VERSION_NUMBER_ONLY=$(echo "$RELEASE_TAG" | sed 's/^v//')

DOWNLOAD_FILE="${PROJECT_NAME}_${VERSION_NUMBER_ONLY}_${TARGET_OS}_${TARGET_ARCH}.zip"
DOWNLOAD_URL="https://github.com/${GITHUB_ORG_REPO}/releases/download/${RELEASE_TAG}/${DOWNLOAD_FILE}"
BINARY_NAME="${PROJECT_NAME}_${TARGET_OS}_${TARGET_ARCH}" # Name of binary *inside* the zip

echo "---"
echo "Downloading ${PROJECT_NAME} version ${RELEASE_TAG} for ${TARGET_OS}_${TARGET_ARCH}..."
echo "URL: $DOWNLOAD_URL"
echo "---"

# --- Download, Unzip, and Install ---
TEMP_DIR=$(mktemp -d) # Create temp directory
ZIP_PATH="$TEMP_DIR/$DOWNLOAD_FILE"
EXTRACTED_BINARY_PATH="$TEMP_DIR/$BINARY_NAME"

echo "Downloading release asset..."
if ! curl -L -f -o "$ZIP_PATH" "$DOWNLOAD_URL"; then
    echo "❌ Error: Failed to download release asset. Verify URL: $DOWNLOAD_URL"
    rm -rf "$TEMP_DIR"; exit 1
fi
echo "✅ Download successful!"

echo "Extracting binary..."
if ! unzip -o "$ZIP_PATH" -d "$TEMP_DIR"; then
    echo "❌ Error: Failed to unzip file. Is 'unzip' installed?"
    rm -rf "$TEMP_DIR"; exit 1
fi
echo "✅ Extraction successful!"

if [ ! -f "$EXTRACTED_BINARY_PATH" ]; then
    echo "❌ Error: Extracted binary '$BINARY_NAME' not found. Check zip contents."
    rm -rf "$TEMP_DIR"; exit 1
fi

echo "Moving $PROJECT_NAME to $INSTALL_DIR and setting permissions..."
sudo mkdir -p "$INSTALL_DIR" # Ensure target directory exists
sudo mv "$EXTRACTED_BINARY_PATH" "$INSTALL_DIR/$PROJECT_NAME" || { echo "❌ Error: Failed to move binary. Do you have sudo permissions?"; rm -rf "$TEMP_DIR"; exit 1; }
sudo chmod +x "$INSTALL_DIR/$PROJECT_NAME" || { echo "❌ Error: Failed to set executable permissions."; rm -rf "$TEMP_DIR"; exit 1; }

# --- Cleanup ---
echo "Cleaning up temporary files..."
rm -rf "$TEMP_DIR"
echo "✅ Cleanup complete."

echo "---"
echo "🎉 **$PROJECT_NAME** has been successfully installed!"
echo "You can now run it by typing: **$PROJECT_NAME**"
echo "---"