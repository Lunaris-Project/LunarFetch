# LunarFetch

<div align="center">
  <img src="assets/image.png" alt="LunarFetch Logo" width="1080">
</div>

A customizable system information tool written in Go.

## Quick Start

```bash
# Install
git clone https://github.com/Lunaris-Project/lunarfetch.git
cd lunarfetch
go run main.go install

# Run
lunarfetch
```

## Features

- Customizable UI with different box styles
- Random ASCII art logos
- Modular information display
- Fast performance
- Built-in installation commands

## Dependencies

Core: `coreutils`, `procps`/`procps-ng`
Optional: `lsb-release`, `xorg-xrandr`, `pciutils`, `gsettings-desktop-schemas`

Check and install dependencies:
```bash
lunarfetch check-deps
lunarfetch install-deps
```

## Commands

```bash
lunarfetch              # Display system info
lunarfetch help         # Show help
lunarfetch install      # Install to system
lunarfetch uninstall    # Remove binary
lunarfetch purge        # Remove binary and config
lunarfetch check-deps   # Check dependencies
lunarfetch install-deps # Install dependencies
lunarfetch build        # Build without installing
```

## Configuration

Configuration file: `~/.config/lunarfetch/config.json`

### Basic Setup

```bash
mkdir -p ~/.config/lunarfetch/logos
cp assets/config.json ~/.config/lunarfetch/config.json
cp assets/logo.txt ~/.config/lunarfetch/logos/moon.txt
```

### Configuration Options

#### Decorations
Change box characters: corners, edges, and separators

#### Logo
Enable/disable logos and set logo directory

#### Modules
Toggle display of system information:
- System: OS, hostname, kernel, uptime, packages
- Hardware: CPU, GPU, memory, disk, battery
- Desktop: DE, resolution, themes, icons, terminal

## Example Configurations

### Minimal
```json
{
  "decorations": {
    "topLeft": "┌", "topRight": "┐",
    "bottomLeft": "└", "bottomRight": "┘",
    "topEdge": "─", "bottomEdge": "─",
    "leftEdge": "│", "rightEdge": "│",
    "separator": "─"
  },
  "logo": {
    "enableLogo": false
  },
  "modules": {
    "show_user": true,
    "show_os": true,
    "show_kernel": true,
    "show_uptime": true,
    "show_packages": true,
    "show_memory": true
  }
}
```

### Rounded Box
```json
{
  "decorations": {
    "topLeft": "╭", "topRight": "╮",
    "bottomLeft": "╰", "bottomRight": "╯",
    "topEdge": "─", "bottomEdge": "─",
    "leftEdge": "│", "rightEdge": "│",
    "separator": "─"
  }
}
```

### Double Line
```json
{
  "decorations": {
    "topLeft": "╔", "topRight": "╗",
    "bottomLeft": "╚", "bottomRight": "╝",
    "topEdge": "═", "bottomEdge": "═",
    "leftEdge": "║", "rightEdge": "║",
    "separator": "═"
  }
}
```
