#!/bin/sh

set -e

echo "Installing DNS Proxy for OpenWRT..."

# Install dependencies
opkg update
opkg install coreutils-nohup curl

# Detect architecture
ARCH=$(uname -m)
case "$ARCH" in
    aarch64) ARCH="arm64";;
    armv7l) ARCH="armv7";;
    x86_64) ARCH="amd64";;
    mips*) ARCH="mips";;
    *) echo "Unsupported architecture: $ARCH"; exit 1;;
esac

# Create install directory
mkdir -p /opt/godnsproxy

# Download and extract latest binary
echo "Downloading and extracting latest version..."
curl -kfsSL "https://github.com/chareice/godnsproxy/releases/latest/download/godnsproxy_$(curl -kfsSL https://api.github.com/repos/chareice/godnsproxy/releases/latest | grep tag_name | cut -d'"' -f4)_linux_$ARCH.tar.gz" \
    | tar -xz -C /opt/godnsproxy
chmod +x /opt/godnsproxy/godnsproxy

# Install init script
echo "Installing service..."
curl -kfsSL "https://raw.githubusercontent.com/chareice/godnsproxy/main/scripts/openwrt-init.d" \
    -o /etc/init.d/godnsproxy
chmod +x /etc/init.d/godnsproxy

# Create default domain file if not exists
[ -f /opt/godnsproxy/domains.txt ] || touch /opt/godnsproxy/domains.txt

# Enable and start service
/etc/init.d/godnsproxy enable
/etc/init.d/godnsproxy start

echo "Installation complete!"
echo "Usage:"
echo "  Edit domains: vi /opt/godnsproxy/domains.txt"
echo "  Restart service: /etc/init.d/godnsproxy restart"
echo "  Upgrade: Re-run this install script"
