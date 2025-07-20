# --- Configuration ---
PROJECT_NAME="wtf"
GITHUB_ORG_REPO="warike/terraform-files" 
INSTALL_DIR="/usr/local/bin"


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
# --- Get the LATEST Release Tag and Asset Names ---
# This requires querying the GitHub API to find the latest release details
echo "---"
echo "Fetching latest release information for ${PROJECT_NAME} from ${GITHUB_ORG_REPO}..."

# Query GitHub API for the latest release JSON
LATEST_RELEASE_JSON=$(curl -s "https://api.github.com/repos/${GITHUB_ORG_REPO}/releases/latest")
LATEST_RELEASE_TAG=$(echo "$LATEST_RELEASE_JSON" | grep '"tag_name":' | head -n 1 | sed -e 's/"tag_name": "//' -e 's/",//' -e 's/ //g')

# If we couldn't get the tag, something went wrong
if [ -z "$LATEST_RELEASE_TAG" ]; then
    echo "❌ Error: Could not fetch latest release tag from GitHub API. Check GITHUB_ORG_REPO or API limits."
    exit 1
fi

echo "✅ Latest release tag found: ${LATEST_RELEASE_TAG}"

# Construct the expected binary file name *within* the zip
BINARY_NAME="${PROJECT_NAME}_${LATEST_RELEASE_TAG}_${TARGET_OS}_${TARGET_ARCH}"
DOWNLOAD_FILE="${BINARY_NAME}.zip"

# Construct the download URL for the specific asset within the latest release
DOWNLOAD_URL="https://github.com/${GITHUB_ORG_REPO}/releases/download/${LATEST_RELEASE_TAG}/${DOWNLOAD_FILE}"


echo "---"
echo "Attempting to download ${PROJECT_NAME} version ${LATEST_RELEASE_TAG} for ${TARGET_OS}_${TARGET_ARCH}..."
echo "Download URL: $DOWNLOAD_URL"
echo "---"

# --- Download, Unzip, and Install (rest of your script is the same) ---
TEMP_DIR=$(mktemp -d) # Create temp directory
ZIP_PATH="$TEMP_DIR/$DOWNLOAD_FILE"
EXTRACTED_BINARY_PATH="$TEMP_DIR/$BINARY_NAME" # This path is correct as is

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
sudo mkdir -p "$INSTALL_DIR"
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