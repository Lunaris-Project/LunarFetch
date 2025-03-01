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

**System Information:**
- `show_os` - Operating system
- `show_host` - Hostname
- `show_kernel` - Kernel version
- `show_uptime` - System uptime
- `show_packages` - Package count
- `show_user` - Username
- `show_shell` - Current shell

**Hardware Information:**
- `show_cpu` - CPU information
- `show_gpu` - GPU information
- `show_memory` - Memory usage
- `show_disk` - Disk usage
- `show_battery` - Battery status
- `show_resolution` - Screen resolution

**Desktop Environment:**
- `show_de` - Desktop environment
- `show_wm_theme` - Window manager theme
- `show_theme` - GTK theme
- `show_icons` - Icon theme
- `show_terminal` - Terminal name

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
