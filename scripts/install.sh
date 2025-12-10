#!/bin/sh
set -e

# ==============================================================================
# BTXZ Installer (Linux/macOS)
# ==============================================================================

REPO="BlackTechX011/BTXZ"
VERSION_URL="https://raw.githubusercontent.com/$REPO/main/version.json"

# ANSI Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Print Banner
echo ""
printf "${CYAN}"
echo "██████╗ ████████╗██╗  ██╗███████╗"
echo "██╔══██╗╚══██╔══╝╚██╗██╔╝╚══███╔╝"
echo "██████╔╝   ██║    ╚███╔╝   ███╔╝ "
echo "██╔══██╗   ██║    ██╔██╗  ███╔╝  "
echo "██████╔╝   ██║   ██╔╝ ██╗███████╗"
echo "╚═════╝    ╚═╝   ╚═╝  ╚═╝╚══════╝"
printf "${NC}\n"
echo ":: Professional Secure Archiver Installer ::"
echo ""

# 1. Detect Architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
    x86_64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    *) 
        printf "${RED}Error: Unsupported architecture $ARCH${NC}\n"
        exit 1 
        ;;
esac

echo "Detected System: ${OS}/${ARCH}"

# 2. Fetch Latest Version Info
printf "${BLUE}Fetching latest release info...${NC}\n"

# Use curl or wget
if command -v curl >/dev/null 2>&1; then
    JSON=$(curl -s $VERSION_URL)
else
    JSON=$(wget -qO- $VERSION_URL)
fi

# Simple JSON parsing using grep/sed to avoid jq dependency
VERSION=$(echo "$JSON" | grep '"version":' | sed -E 's/.*"version": "([^"]+)".*/\1/')
PLATFORM_KEY="${OS}-${ARCH}"

# Extract URL and Checksum for this platform
# We match the block for the platform key and extract the fields
BLOCK=$(echo "$JSON" | grep -A 4 "\"$PLATFORM_KEY\"")
DOWNLOAD_URL=$(echo "$BLOCK" | grep '"url":' | sed -E 's/.*"url": "([^"]+)".*/\1/')
EXPECTED_SHA256=$(echo "$BLOCK" | grep '"sha256":' | sed -E 's/.*"sha256": "([^"]+)".*/\1/')

if [ -z "$DOWNLOAD_URL" ]; then
    printf "${RED}Error: Could not find download URL for $PLATFORM_KEY${NC}\n"
    exit 1
fi

echo "Target Version: ${GREEN}v${VERSION}${NC}"
echo "Download URL:   $DOWNLOAD_URL"

# 3. Download Binary
TMP_DIR=$(mktemp -d)
TMP_FILE="$TMP_DIR/btxz"

printf "${BLUE}Downloading binary...${NC}\n"
if command -v curl >/dev/null 2>&1; then
    curl -L "$DOWNLOAD_URL" -o "$TMP_FILE"
else
    wget -O "$TMP_FILE" "$DOWNLOAD_URL"
fi

# 4. Verify Checksum
if [ -n "$EXPECTED_SHA256" ]; then
    printf "${BLUE}Verifying checksum...${NC} "
    
    if command -v sha256sum >/dev/null 2>&1; then
        CALCULATED_SHA256=$(sha256sum "$TMP_FILE" | awk '{print $1}')
    elif command -v shasum >/dev/null 2>&1; then
        CALCULATED_SHA256=$(shasum -a 256 "$TMP_FILE" | awk '{print $1}')
    else
        printf "${YELLOW}Skipped (SHA256 tools not found)${NC}\n"
    fi

    if [ -n "$CALCULATED_SHA256" ]; then
        if [ "$CALCULATED_SHA256" = "$EXPECTED_SHA256" ]; then
            printf "${GREEN}OK${NC}\n"
        else
            printf "${RED}FAILED${NC}\n"
            printf "${RED}Expected: $EXPECTED_SHA256${NC}\n"
            printf "${RED}Got:      $CALCULATED_SHA256${NC}\n"
            exit 1
        fi
    fi
else
    printf "${YELLOW}Skipped (No checksum in manifest)${NC}\n"
fi

# 5. Install
INSTALL_DIR="/usr/local/bin"
if [ "$OS" = "android" ]; then
   # Termux support
   INSTALL_DIR="$PREFIX/bin"
fi

printf "${BLUE}Installing to $INSTALL_DIR...${NC}\n"

chmod +x "$TMP_FILE"

if [ -w "$INSTALL_DIR" ]; then
    mv "$TMP_FILE" "$INSTALL_DIR/btxz"
else
    echo "Sudo permissions required to write to $INSTALL_DIR"
    sudo mv "$TMP_FILE" "$INSTALL_DIR/btxz"
fi

# Cleanup
rm -rf "$TMP_DIR"

# 6. Verify Installation
echo ""
echo "==========================================="
if command -v btxz >/dev/null 2>&1; then
    INSTALLED_VERSION=$(btxz --version)
    printf "${GREEN}Success! BTXZ installed.${NC}\n"
    echo "$INSTALLED_VERSION"
else
    printf "${RED}Installation failed. 'btxz' not found in PATH.${NC}\n"
    exit 1
fi
echo "==========================================="
echo "Run 'btxz help' to get started."
