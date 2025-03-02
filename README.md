# LunarFetch

A modern, customizable system information display tool with image rendering capabilities.

<div align="center">
<a href="#features">Features</a> •
<a href="#installation">Installation</a> •
<a href="#usage">Usage</a> •
<a href="#configuration">Configuration</a> •
<a href="#contributing">Contributing</a> •
<a href="#license">License</a>
</div>

## ✨ Features

- **System Information Display**: Shows detailed system information including OS, kernel, CPU, GPU, memory usage, and more
- **Image Rendering**: Supports multiple image rendering protocols (Sixel, Kitty, iTerm2, Chafa)
- **ASCII Art Logos**: Display custom ASCII art logos alongside system information
- **Customizable UI**: Configure colors, layout, and information displayed
- **Cross-Platform**: Works on various Linux distributions
- **Modular Design**: Easily extendable with new information modules

## 📦 Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/Lunaris-Project/lunarfetch.git
cd lunarfetch

# Build the binary
go build -o lunarfetch

# Install (optional)
./lunarfetch install
```

### Dependencies

LunarFetch requires the following dependencies for full functionality:

- `go` (for building)
- `chafa` (for image rendering)

You can install dependencies with:

```bash
./lunarfetch install-deps
```

## 🚀 Usage

```bash
# Basic usage
lunarfetch

# With debug information
lunarfetch --debug

# Using a custom configuration file
lunarfetch -c /path/to/config.json

# Display version information
lunarfetch --version
```

### Command Line Options

```
Usage: lunarfetch [options]

Options:
  -c, --config <file>   Use custom configuration file
  -d, --debug           Enable debug mode
  -v, --version         Display version information
  -h, --help            Show this help message

Commands:
  install               Install LunarFetch to your system
  uninstall             Remove LunarFetch from your system
  install-deps          Install required system dependencies
  build                 Build the binary without installing
  setup-image           Configure image display support
```

## ⚙️ Configuration

LunarFetch can be configured using a JSON configuration file located at `~/.config/lunarfetch/config.json`.

### Configuration Overview

<details>
<summary><b>🔳 Decorations</b> - Box drawing characters for display</summary>

```json
"decorations": {
  "topLeft": "╭",
  "topRight": "╮",
  "bottomLeft": "╰",
  "bottomRight": "╯",
  "topEdge": "─",
  "bottomEdge": "─",
  "leftEdge": "│",
  "rightEdge": "│",
  "separator": ": "
}
```

**Alternative styles:**
- Regular box: `"topLeft": "┌", "topRight": "┐", "bottomLeft": "└", "bottomRight": "┘"`
- Rounded box (shown above): `"topLeft": "╭", "topRight": "╮", "bottomLeft": "╰", "bottomRight": "╯"`
- Double line: `"topLeft": "╔", "topRight": "╗", "bottomLeft": "╚", "bottomRight": "╝", "topEdge": "═", "bottomEdge": "═", "leftEdge": "║", "rightEdge": "║"`
</details>

<details>
<summary><b>🖼️ Logo</b> - ASCII art logo configuration</summary>

```json
"logo": {
  "enableLogo": true,
  "type": "ascii",
  "content": "",
  "location": "center",
  "logoPath": "~/.config/lunarfetch/logos",
  "position": "side"
}
```

**Options:**
- `enableLogo`: Enable/disable logo display (`true` or `false`)
- `type`: Logo type (`"ascii"` or `"file"` to load from a file)
- `content`: Custom ASCII content (when type is `"ascii"`)
- `location`: Text alignment (`"center"`, `"left"`, or `"right"`)
- `logoPath`: Directory containing logo files
- `position`: Position relative to system info (`"side"` or `"above"`)
</details>

<details>
<summary><b>🖼️ Image</b> - Image display configuration</summary>

```json
"image": {
  "enableImage": true,
  "enabled": true,
  "random": true,
  "imagePath": "~/.config/lunarfetch/images",
  "width": 40,
  "height": 20,
  "renderMode": "block",
  "ditherMode": "floyd-steinberg",
  "terminalOutput": true,
  "displayMode": "block",
  "protocol": "chafa",
  "scale": 1,
  "offset": 2,
  "background": "transparent",
  "position": "side"
}
```

**Options:**
- `enableImage`/`enabled`: Enable/disable image display (`true` or `false`)
- `random`: Randomly select an image from the `imagePath` directory (`true` or `false`)
- `imagePath`: Path to image file or directory (for random selection)
- `width`/`height`: Dimensions in terminal characters
- `renderMode`: Image rendering detail level (`"detailed"`, `"simple"`, `"block"`, or `"ascii"`)
- `ditherMode`: Dithering algorithm (`"none"` or `"floyd-steinberg"`)
- `terminalOutput`: Output to terminal directly (`true` or `false`)
- `displayMode`: How to display the image (`"auto"`, `"block"`, or `"ascii"`)
- `protocol`: Image display protocol:
  - `"auto"`: Auto-detect the best protocol
  - `"sixel"`: For terminals with Sixel support
  - `"kitty"`: For Kitty terminal
  - `"iterm2"`: For iTerm2 terminal on macOS
  - `"chafa"`: Uses the Chafa tool
  - `"uberzug"`: Uses Überzug (Linux only)
  - `"terminal-image"`: Uses the terminal-image tool
- `scale`: Image scaling factor (integer)
- `offset`: Offset from terminal edge (integer)
- `background`: Background color (`"transparent"` or a color value)
- `position`: Position relative to system info (`"side"` or `"above"`)
</details>

<details>
<summary><b>📋 Display</b> - Controls display order</summary>

```json
"display": {
  "showLogoFirst": false,
  "showImageFirst": true
}
```

**Options:**
- `showLogoFirst`: When `true`, logo appears before system info
- `showImageFirst`: When `true`, image appears before system info

Note: If both are `true`, logo takes precedence.
</details>

<details>
<summary><b>🔣 Icons</b> - Icons for system information</summary>

```json
"icons": {
  "host": "󰒋",
  "user": "󰀄",
  "os": "󰣇",
  "kernel": "󰣇",
  "uptime": "󰔟",
  "terminal": "󰆍",
  "shell": "󰆍",
  "disk": "󰋊",
  "memory": "󰍛",
  "packages": "󰏗",
  "battery": "󰂄",
  "gpu": "󰢮",
  "cpu": "󰘚",
  "resolution": "󰍹",
  "de": "󰧨",
  "wm_theme": "󰏘",
  "theme": "󰔯",
  "icons": "󰀻"
}
```

You can also use emoji instead of Nerd Font icons:
```json
"icons": {
  "host": "🏠",
  "user": "👤",
  "os": "🐧",
  "kernel": "🧠",
  "uptime": "⏱️",
  "terminal": "💻",
  "shell": "🐚",
  "disk": "💾",
  "memory": "🧮",
  "packages": "📦",
  "battery": "🔋",
  "gpu": "🎮",
  "cpu": "⚙️",
  "resolution": "🖥️",
  "de": "🖼️",
  "wm_theme": "🎨",
  "theme": "🎭",
  "icons": "🔍"
}
```
</details>

<details>
<summary><b>📊 Modules</b> - Enable/disable information components</summary>

```json
"modules": {
  "show_user": true,
  "show_cpu": true,
  "show_gpu": true,
  "show_uptime": true,
  "show_shell": true,
  "show_memory": true,
  "show_packages": true,
  "show_os": true,
  "show_host": true,
  "show_kernel": true,
  "show_battery": true,
  "show_disk": true,
  "show_resolution": true,
  "show_de": true,
  "show_wm_theme": true,
  "show_theme": true,
  "show_icons": true,
  "show_terminal": true
}
```

Set any option to `false` to hide that specific information.
</details>

### Example Configurations

<details>
<summary><b>Minimal Configuration</b></summary>

```json
{
  "decorations": {
    "topLeft": "┌", "topRight": "┐",
    "bottomLeft": "└", "bottomRight": "┘",
    "topEdge": "─", "bottomEdge": "─",
    "leftEdge": "│", "rightEdge": "│",
    "separator": ": "
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
</details>

<details>
<summary><b>Image-Only Configuration</b></summary>

```json
{
  "decorations": {
    "topLeft": "┌", "topRight": "┐",
    "bottomLeft": "└", "bottomRight": "┘",
    "topEdge": "─", "bottomEdge": "─",
    "leftEdge": "│", "rightEdge": "│",
    "separator": ": "
  },
  "logo": {
    "enableLogo": false
  },
  "image": {
    "enableImage": true,
    "random": true,
    "imagePath": "~/.config/lunarfetch/images",
    "width": 40,
    "height": 20,
    "protocol": "chafa",
    "position": "side"
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
</details>

<details>
<summary><b>Full Configuration</b></summary>

```json
{
  "decorations": {
    "topLeft": "╭",
    "topRight": "╮",
    "bottomLeft": "╰",
    "bottomRight": "╯",
    "topEdge": "─",
    "bottomEdge": "─",
    "leftEdge": "│",
    "rightEdge": "│",
    "separator": ": "
  },
  "logo": {
    "enableLogo": true,
    "type": "ascii",
    "content": "",
    "location": "center",
    "logoPath": "~/.config/lunarfetch/logos",
    "position": "side"
  },
  "image": {
    "enableImage": true,
    "enabled": true,
    "random": true,
    "imagePath": "~/.config/lunarfetch/images",
    "width": 40,
    "height": 20,
    "renderMode": "block",
    "ditherMode": "floyd-steinberg",
    "terminalOutput": true,
    "displayMode": "block",
    "protocol": "chafa",
    "scale": 1,
    "offset": 2,
    "background": "transparent",
    "position": "side"
  },
  "display": {
    "showLogoFirst": false,
    "showImageFirst": true
  },
  "icons": {
    "host": "󰒋",
    "user": "󰀄",
    "os": "󰣇",
    "kernel": "󰣇",
    "uptime": "󰔟",
    "terminal": "󰆍",
    "shell": "󰆍",
    "disk": "󰋊",
    "memory": "󰍛",
    "packages": "󰏗",
    "battery": "󰂄",
    "gpu": "󰢮",
    "cpu": "󰘚",
    "resolution": "󰍹",
    "de": "󰧨",
    "wm_theme": "󰏘",
    "theme": "󰔯",
    "icons": "󰀻"
  },
  "modules": {
    "show_user": true,
    "show_cpu": true,
    "show_gpu": true,
    "show_uptime": true,
    "show_shell": true,
    "show_memory": true,
    "show_packages": true,
    "show_os": true,
    "show_host": true,
    "show_kernel": true,
    "show_battery": true,
    "show_disk": true,
    "show_resolution": true,
    "show_de": true,
    "show_wm_theme": true,
    "show_theme": true,
    "show_icons": true,
    "show_terminal": true
  }
}
```
</details>

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.
