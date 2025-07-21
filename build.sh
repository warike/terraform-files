#!/bin/bash

# --- Project Configuration ---
APP_NAME="wtf" # Your application name
# --- End Configuration ---

# --- 1. Get the Version ---
# Attempts to get the latest Git tag (e.g., "v1.0.0", "v2.3.4").
# If no tags are found (e.g., a new repository without tags), it defaults to "1.0.0".
VERSION=$(git describe --tags --abbrev=0 2>/dev/null || echo "v1.0.0")
echo "Preparing build for version: ${VERSION}"

# --- 2. Prepare the Build Directory ---
# Creates the 'build' folder if it doesn't exist. This is where binaries will be stored.
mkdir -p build
echo "Build directory 'build/' created or already exists."

# Define linker flags to inject the version variable.
# This sets the 'Version' variable in the 'main' Go package to the value of our ${VERSION} shell variable.
LDFLAGS="-X main.Version=${VERSION}"
echo "Linker flags configured: ${LDFLAGS}"

# --- 3. Compile the Binaries ---
echo "Starting binary compilation for multiple platforms..."

# Compile for macOS (Apple Silicon - ARM64)
echo "  -> Compiling for macOS (Apple Silicon - ARM64)..."
GOOS=darwin GOARCH=arm64 go build -ldflags="${LDFLAGS}" -o "build/${APP_NAME}_darwin_arm64" .

# Compile for macOS (Intel - x86_64)
echo "  -> Compiling for macOS (Intel - x86_64)..."
GOOS=darwin GOARCH=amd64 go build -ldflags="${LDFLAGS}" -o "build/${APP_NAME}_darwin_amd64" .

# You can add more platforms if you need (e.g., Linux, Windows):
# echo "  -> Compiling for Linux (AMD64)..."
# GOOS=linux GOARCH=amd64 go build -ldflags="${LDFLAGS}" -o "build/${APP_NAME}_linux_amd64" .
#
# echo "  -> Compiling for Windows (AMD64)..."
# GOOS=windows GOARCH=amd64 go build -ldflags="${LDFLAGS}" -o "build/${APP_NAME}_windows_amd64.exe" .

echo "Binary compilation completed."
echo "Binaries created in 'build/':"
ls -lh build/

# --- 4. Package Binaries for Release ---
# Change into the build directory to compress files directly there.
echo "Packaging binaries into .zip files for release..."
cd build || { echo "Error: Could not enter 'build' directory."; exit 1; }

# Clean up old .zip files to prevent conflicts
rm -f *.zip

# Create a .zip file for each binary, including the version in the name.
zip -r "${APP_NAME}_${VERSION}_darwin_arm64.zip" "${APP_NAME}_darwin_arm64"
zip -r "${APP_NAME}_${VERSION}_darwin_amd64.zip" "${APP_NAME}_darwin_amd64"
# If you compiled for more platforms, add their zips here:
# zip -r "${APP_NAME}_${VERSION}_linux_amd64.zip" "${APP_NAME}_linux_amd64"
# zip -r "${APP_NAME}_${VERSION}_windows_amd64.zip" "${APP_NAME}_windows_amd64.exe"

echo "Release archives created:"
ls -lh *.zip
