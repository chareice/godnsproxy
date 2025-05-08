#!/bin/sh

# Configuration
GITHUB_REPO="chareice/godnsproxy"
INSTALL_DIR="/opt/godnsproxy"
SERVICE_NAME="godnsproxy"
CONFIG_DIR="/etc/config"
ARCH="$(uname -m)"

# Detect architecture
case "$ARCH" in
    "x86_64")
        ARCH_NAME="amd64"
        ;;
    "aarch64")
        ARCH_NAME="arm64"
        ;;
    "armv7l")
        ARCH_NAME="armv7"
        ;;
    "mips"*)
        ARCH_NAME="mips"
        ;;
    *)
        echo "Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

# Create install directory
mkdir -p "$INSTALL_DIR"

# Get latest version
echo "Checking for latest version..."
LATEST_VERSION=$(curl -kfsSL "https://api.github.com/repos/$GITHUB_REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST_VERSION" ]; then
    echo "Failed to get latest version info"
    exit 1
fi

echo "Latest version: $LATEST_VERSION"

# Check current version
CURRENT_VERSION=""
if [ -f "$INSTALL_DIR/godnsproxy" ]; then
    CURRENT_VERSION=$("$INSTALL_DIR/godnsproxy" --version 2>/dev/null || echo "")
fi

if [ "$CURRENT_VERSION" = "$LATEST_VERSION" ]; then
    echo "Already on latest version: $LATEST_VERSION"
    exit 0
fi

# Download latest version
DOWNLOAD_URL="https://github.com/$GITHUB_REPO/releases/download/$LATEST_VERSION/godnsproxy_${LATEST_VERSION#v}_linux_${ARCH_NAME}.tar.gz"
echo "Downloading: $DOWNLOAD_URL"

TMP_DIR=$(mktemp -d)
echo "Using temp directory: $TMP_DIR"

# Download with error checking
curl -kfSL "$DOWNLOAD_URL" -o "$TMP_DIR/godnsproxy.tar.gz"
CURL_EXIT_CODE=$?

if [ $CURL_EXIT_CODE -ne 0 ]; then
    echo "Download failed, curl exit code: $CURL_EXIT_CODE"
    echo "Please check:"
    echo "1. Network connection"
    echo "2. Version number: $LATEST_VERSION"
    echo "3. Architecture: $ARCH_NAME"
    echo "4. Full download URL: $DOWNLOAD_URL"
    rm -rf "$TMP_DIR"
    exit 1
fi

# Verify download size
FILE_SIZE=$(ls -l "$TMP_DIR/godnsproxy.tar.gz" | awk '{print $5}')
echo "Download size: $FILE_SIZE bytes"

if [ $FILE_SIZE -lt 1000 ]; then
    echo "Download too small, may be invalid"
    rm -rf "$TMP_DIR"
    exit 1
fi

# Extract and install
cd "$TMP_DIR"
echo "Extracting files..."
tar xzf godnsproxy.tar.gz
TAR_EXIT_CODE=$?

if [ $TAR_EXIT_CODE -ne 0 ]; then
    echo "Extraction failed, tar exit code: $TAR_EXIT_CODE"
    rm -rf "$TMP_DIR"
    exit 1
fi

# Verify extracted files
if [ ! -f "godnsproxy" ]; then
    echo "No executable found after extraction"
    rm -rf "$TMP_DIR"
    exit 1
fi

chmod +x godnsproxy

# Stop existing service if running
if [ -f "/etc/init.d/$SERVICE_NAME" ] && [ -x "/etc/init.d/$SERVICE_NAME" ]; then
    echo "Stopping existing service..."
    /etc/init.d/$SERVICE_NAME stop
fi

# Install binary
echo "Installing binary..."
mv godnsproxy "$INSTALL_DIR/"

# Install service file
echo "Installing service..."
curl -kfsSL "https://raw.githubusercontent.com/$GITHUB_REPO/main/scripts/openwrt-init.d" \
    -o "/etc/init.d/$SERVICE_NAME"
chmod +x "/etc/init.d/$SERVICE_NAME"

# Create default domains file if not exists
[ -f "$INSTALL_DIR/domains.txt" ] || touch "$INSTALL_DIR/domains.txt"

# Create config file if not exists
if [ ! -f "$CONFIG_DIR/$SERVICE_NAME" ]; then
    echo "Creating default config..."
    cat > "$CONFIG_DIR/$SERVICE_NAME" << 'EOF'
config godnsproxy 'main'
    option enabled '1'
    option port '53'
    option china_server '223.5.5.5'
    option trust_server 'https://1.1.1.1/dns-query' 
    option domains_file '/opt/godnsproxy/domains.txt'
    option log_level 'info'
EOF
fi

# Cleanup
rm -rf "$TMP_DIR"

# Enable and start service
echo "Enabling and starting service..."
/etc/init.d/$SERVICE_NAME enable
/etc/init.d/$SERVICE_NAME start

echo "Installation complete!"
echo ""
echo "Usage:"
echo "  Edit domains: vi $INSTALL_DIR/domains.txt"
echo "  Edit config: vi $CONFIG_DIR/$SERVICE_NAME"
echo "  Restart service: /etc/init.d/$SERVICE_NAME restart"
echo "  Stop service: /etc/init.d/$SERVICE_NAME stop"
echo "  Upgrade: Re-run this install script"
