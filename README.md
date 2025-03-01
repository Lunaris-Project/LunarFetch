# LunarFetch

<div align="center">
  <img src="src/assets/image.png" alt="LunarFetch Logo" width="1080">
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
- Image support (PNG, JPG, JPEG, WebP)
- Modular information display
- Fast performance
- Built-in installation commands

## Project Structure

```
lunarfetch/
├── src/
│   ├── assets/        # Assets like logos and default config
│   ├── components/    # Core components
│   └── scripts/       # CLI functionality scripts
│   └── utils/         # Utility functions
├── tests/
│   ├── image/         # Image display tests
│   └── config/        # Configuration tests
├── main.go            # Main source code
├── go.mod             # Go module definition
├── go.sum             # Go module checksums
└── README.md          # This documentation
```

## Dependencies

Core: `coreutils`, `procps`/`procps-ng`
Optional: `lsb-release`, `xorg-xrandr`, `pciutils`, `gsettings-desktop-schemas`, `imagemagick`

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
lunarfetch setup-image  # Configure image display settings
```

## Configuration

Configuration file: `~/.config/lunarfetch/config.json`

### Basic Setup

```bash
mkdir -p ~/.config/lunarfetch/logos ~/.config/lunarfetch/images
cp src/assets/config.json ~/.config/lunarfetch/config.json
cp src/assets/logo.txt ~/.config/lunarfetch/logos/moon.txt
# If you have images, copy them to the images directory
cp your-image.png ~/.config/lunarfetch/images/
```

### Configuration Options

#### Decorations
Change box characters: corners, edges, and separators

#### Logo
Enable/disable ASCII art logos and set logo directory

#### Image
Enable/disable image display and configure image settings:
- `enableImage` - Enable or disable image display
- `imagePath` - Path to image or directory containing images
- `width` - Width of the image in characters
- `height` - Height of the image in characters
- `renderMode` - Rendering mode: "simple", "detailed", or "default"
- `ditherMode` - Dithering mode: "none" or "floyd-steinberg"
- `terminalOutput` - Use ANSI color output instead of ASCII art

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

### With Image Support
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
  "image": {
    "enableImage": true,
    "imagePath": "~/.config/lunarfetch/images",
    "width": 80,
    "height": 24,
    "renderMode": "detailed",
    "ditherMode": "floyd-steinberg",
    "terminalOutput": false
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

## Image Display

LunarFetch supports displaying actual images in the terminal using various protocols:

- **Sixel**: For terminals that support Sixel graphics (xterm with sixel support, mlterm)
- **Kitty**: For Kitty terminal
- **iTerm2**: For iTerm2 terminal on macOS
- **Chafa**: Uses the Chafa tool to display images
- **Überzug**: Uses Überzug to display images (Linux only)
- **terminal-image**: Uses the terminal-image tool

To set up image support:

1. Run `./lunarfetch setup-image`
2. Place your images in `~/.config/lunarfetch/images/`
3. Edit the config file to customize image settings:

```json
"image": {
  "enabled": true,
  "random": false,
  "imagePath": "~/.config/lunarfetch/images/myimage.png",
  "width": 40,
  "height": 20,
  "protocol": "auto",
  "scale": 1,
  "offset": 2,
  "position": "side"
}
```

#### Image Settings

- `enabled`: Enable image display
- `random`: Randomly select an image from the directory specified in `imagePath`
- `imagePath`: Path to an image file or directory
- `width`: Width of the image in characters
- `height`: Height of the image in characters
- `protocol`: Image display protocol ("auto", "sixel", "kitty", "iterm2", "chafa", "uberzug", "terminal-image")
- `scale`: Scale factor for the image
- `offset`: Offset from the edge of the terminal
- `position`: Position of the image relative to system info ("side" or "above")

#### Display Options

LunarFetch now supports displaying images or logos either above or beside the system information:

```json
"logo": {
  "enableLogo": true,
  "type": "ascii",
  "logoPath": "~/.config/lunarfetch/logos",
  "position": "side"
},
"image": {
  "enabled": true,
  "random": false,
  "imagePath": "~/.config/lunarfetch/images/myimage.png",
  "width": 40,
  "height": 20,
  "protocol": "auto",
  "scale": 1,
  "offset": 2,
  "position": "side"
},
"display": {
  "showLogoFirst": true,
  "showImageFirst": false
}
```

#### Display Configuration Options

- **Logo Position**:
  - `logo.position`: Set to "side" (default) to display the logo beside system info, or "above" to display it above
  
- **Image Position**:
  - `image.position`: Set to "side" (default) to display the image beside system info, or "above" to display it above
  
- **Display Order** (when using "side" position):
  - `display.showLogoFirst`: When true, displays the logo before system info (default: true)
  - `display.showImageFirst`: When true, displays the image before system info (default: false)

When both logo and image are enabled:
- If both have position "above", the logo takes precedence
- If both have position "side", the display order is determined by `showLogoFirst` and `showImageFirst`

### Example with Image Above System Info

```json
{
  "decorations": {
    "topLeft": "╭", "topRight": "╮",
    "bottomLeft": "╰", "bottomRight": "╯",
    "topEdge": "─", "bottomEdge": "─",
    "leftEdge": "│", "rightEdge": "│",
    "separator": "─"
  },
  "logo": {
    "enableLogo": false
  },
  "image": {
    "enabled": true,
    "random": false,
    "imagePath": "~/.config/lunarfetch/images/myimage.png",
    "width": 40,
    "height": 20,
    "protocol": "auto",
    "scale": 1,
    "offset": 2,
    "position": "above"
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

## Modules

LunarFetch can display various system information modules:

- OS information
- Host information
- Kernel version
- CPU information
- GPU information
- Memory usage
- Disk usage
- Shell
- Terminal
- Package count
- Battery status
- Desktop environment
- Theme information
- Icon theme
- Resolution

## Testing

LunarFetch includes a comprehensive test suite to verify its functionality:

### Image Display Tests

Test the image display functionality:

```bash
go run tests/image/test-image.go ~/.config/lunarfetch/images/test.png
```

### Configuration Tests

Test different configuration options:

```bash
# Generate test configurations
go run tests/config/main.go

# Test image above configuration
go run main.go -c ~/.config/lunarfetch/test-above.json

# Test image beside configuration
go run main.go -c ~/.config/lunarfetch/test-side.json
```

For more details, see the [Tests README](tests/README.md).

## License

MIT
