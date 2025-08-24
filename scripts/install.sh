#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
BOLD='\033[1m'
NC='\033[0m' # No Color

# Configuration
REPO="philipptpunkt/ledger-live-starter"
INSTALL_DIR="$HOME/.ledger-live"
BINARY_NAME="ledger-live"

echo -e "${BOLD}${BLUE}Ledger Live Starter Installer${NC}"
echo -e "Installing to: ${YELLOW}$INSTALL_DIR${NC}"
echo

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case $ARCH in
    x86_64) ARCH="amd64" ;;
    arm64|aarch64) ARCH="arm64" ;;
    *) 
        echo -e "${RED}Error: Unsupported architecture: $ARCH${NC}"
        exit 1 
        ;;
esac

case $OS in
    linux) OS="linux" ;;
    darwin) OS="darwin" ;;
    *)
        echo -e "${RED}Error: Unsupported OS: $OS${NC}"
        exit 1
        ;;
esac

echo -e "Detected: ${GREEN}$OS-$ARCH${NC}"

# Create installation directory
echo -e "${BLUE}Creating installation directory...${NC}"
mkdir -p "$INSTALL_DIR"

# Get latest release info
echo -e "${BLUE}Fetching latest release...${NC}"
RELEASE_URL="https://api.github.com/repos/$REPO/releases/latest"
DOWNLOAD_URL=$(curl -s "$RELEASE_URL" | grep -o "https://.*ledger-live-$OS-$ARCH[^\"]*")

if [ -z "$DOWNLOAD_URL" ]; then
    echo -e "${RED}Error: Could not find binary for $OS-$ARCH${NC}"
    echo -e "${YELLOW}Available releases: https://github.com/$REPO/releases${NC}"
    exit 1
fi

echo -e "Download URL: ${GREEN}$DOWNLOAD_URL${NC}"

# Download binary
echo -e "${BLUE}Downloading ledger-live-starter...${NC}"
TEMP_FILE=$(mktemp)
curl -L -o "$TEMP_FILE" "$DOWNLOAD_URL"

# Install binary
echo -e "${BLUE}Installing binary...${NC}"
mv "$TEMP_FILE" "$INSTALL_DIR/$BINARY_NAME"
chmod +x "$INSTALL_DIR/$BINARY_NAME"

# Add to PATH
echo -e "${BLUE}Adding to PATH...${NC}"
SHELL_CONFIG=""
case "$SHELL" in
    */zsh) SHELL_CONFIG="$HOME/.zshrc" ;;
    */bash) SHELL_CONFIG="$HOME/.bashrc" ;;
    *) SHELL_CONFIG="$HOME/.profile" ;;
esac

# Check if already in PATH
if ! echo "$PATH" | grep -q "$INSTALL_DIR"; then
    echo "export PATH=\"$INSTALL_DIR:\$PATH\"" >> "$SHELL_CONFIG"
    export PATH="$INSTALL_DIR:$PATH"
    echo -e "${GREEN}Added $INSTALL_DIR to PATH in $SHELL_CONFIG${NC}"
else
    echo -e "${YELLOW}$INSTALL_DIR already in PATH${NC}"
fi

# Verify installation
echo -e "${BLUE}Verifying installation...${NC}"
if "$INSTALL_DIR/$BINARY_NAME" version >/dev/null 2>&1; then
    echo -e "${GREEN}✓ Installation successful!${NC}"
else
    echo -e "${YELLOW}⚠ Installation completed, but verification failed${NC}"
    echo -e "${YELLOW}  You may need to restart your terminal${NC}"
fi

echo
echo -e "${BOLD}${GREEN}Installation Complete!${NC}"
echo
echo -e "${BOLD}Usage:${NC}"
echo -e "  ${GREEN}ledger-live start${NC}    - Start with interactive menu"
echo -e "  ${GREEN}ledger-live setup${NC}    - Run initial setup"
echo -e "  ${GREEN}ledger-live --help${NC}   - Show help"
echo
echo -e "${YELLOW}Note: You may need to restart your terminal or run:${NC}"
echo -e "  ${BLUE}source $SHELL_CONFIG${NC}"
echo
echo -e "${BLUE}Config location: ${YELLOW}$INSTALL_DIR/config.json${NC}"
