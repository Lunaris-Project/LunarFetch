#!/bin/bash

echo "Installing LunarFetch..."

# Check for required dependencies
echo "Checking for required dependencies..."
MISSING_DEPS=()

# Function to check if a command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Check for core dependencies
if ! command_exists df || ! command_exists free; then
    MISSING_DEPS+=("coreutils")
fi

if ! command_exists uptime; then
    MISSING_DEPS+=("procps/procps-ng")
fi

# Check for feature-specific dependencies
if ! command_exists lsb_release; then
    MISSING_DEPS+=("lsb-release")
fi

if ! command_exists xrandr && ! command_exists xdpyinfo && ! command_exists swaymsg && ! command_exists wlr-randr; then
    MISSING_DEPS+=("xorg-xrandr/xdpyinfo/sway/wlroots")
fi

if ! command_exists lscpu || ! command_exists lspci; then
    MISSING_DEPS+=("pciutils")
fi

if ! command_exists gsettings; then
    MISSING_DEPS+=("gsettings-desktop-schemas")
fi

# Warn about missing dependencies
if [ ${#MISSING_DEPS[@]} -gt 0 ]; then
    echo "Warning: The following dependencies appear to be missing:"
    for dep in "${MISSING_DEPS[@]}"; do
        echo "  - $dep"
    done
    echo "Some features may not work correctly. See README.md for installation instructions."
    echo "Press Enter to continue anyway, or Ctrl+C to abort..."
    read -r
fi

# Build the binary
echo "Building LunarFetch..."
go build -o lunarfetch

# Create config directory if it doesn't exist
mkdir -p ~/.config/lunarfetch/logos

# Copy the config file if it doesn't exist
if [ ! -f ~/.config/lunarfetch/config.json ]; then
    cp assets/config.json ~/.config/lunarfetch/config.json
    echo "Created default configuration file"
fi

# Copy the sample logo if it doesn't exist
if [ ! -f ~/.config/lunarfetch/logos/moon.txt ]; then
    cp assets/logo.txt ~/.config/lunarfetch/logos/moon.txt
    echo "Added sample logo"
fi

# Install the binary
sudo mv lunarfetch /usr/local/bin/
echo "LunarFetch has been installed to /usr/local/bin/"
echo "Run 'lunarfetch' to try it out!" 