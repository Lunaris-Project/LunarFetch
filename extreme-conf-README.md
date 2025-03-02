# LunarFetch Extreme Configuration Guide

This document explains all the configuration options available in the `extreme-conf.json` file, which demonstrates every possible setting for LunarFetch.

## How to Use

Copy the `extreme-conf.json` file to `~/.config/lunarfetch/config.json` or use it with the `-c` flag:

```bash
./lunarfetch -c /path/to/extreme-conf.json
```

## Configuration Sections

### Decorations

Box drawing characters for the system information display.

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

Alternative styles:
- Regular box: `"topLeft": "┌", "topRight": "┐", "bottomLeft": "└", "bottomRight": "┘"`
- Rounded box (shown above): `"topLeft": "╭", "topRight": "╮", "bottomLeft": "╰", "bottomRight": "╯"`
- Double line: `"topLeft": "╔", "topRight": "╗", "bottomLeft": "╚", "bottomRight": "╝", "topEdge": "═", "bottomEdge": "═", "leftEdge": "║", "rightEdge": "║"`

### Logo

ASCII art logo configuration.

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

Options:
- `enableLogo`: Enable/disable logo display (`true` or `false`)
- `type`: Logo type (`"ascii"` or `"file"` to load from a file)
- `content`: Custom ASCII content (when type is `"ascii"`)
- `location`: Text alignment (`"center"`, `"left"`, or `"right"`)
- `logoPath`: Directory containing logo files
- `position`: Position relative to system info (`"side"` or `"above"`)

### Image

Image display configuration.

```json
"image": {
  "enableImage": true,
  "enabled": true,
  "random": true,
  "imagePath": "~/.config/lunarfetch/images",
  "width": 40,
  "height": 20,
  "renderMode": "detailed",
  "ditherMode": "floyd-steinberg",
  "terminalOutput": false,
  "displayMode": "block",
  "protocol": "auto",
  "scale": 1,
  "offset": 2,
  "background": "transparent",
  "position": "side"
}
```

Options:
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
  - `"sixel"`: For terminals with Sixel support (xterm with sixel, mlterm)
  - `"kitty"`: For Kitty terminal
  - `"iterm2"`: For iTerm2 terminal on macOS
  - `"chafa"`: Uses the Chafa tool
  - `"uberzug"`: Uses Überzug (Linux only)
  - `"terminal-image"`: Uses the terminal-image tool
- `scale`: Image scaling factor (integer)
- `offset`: Offset from terminal edge (integer)
- `background`: Background color (`"transparent"` or a color value)
- `position`: Position relative to system info (`"side"` or `"above"`)

### Display

Controls display order when both logo and image are enabled.

```json
"display": {
  "showLogoFirst": true,
  "showImageFirst": false
}
```

Options:
- `showLogoFirst`: When `true`, logo appears before system info
- `showImageFirst`: When `true`, image appears before system info

Note: If both are `true`, logo takes precedence.

### Icons

Icons displayed next to each system information item.

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

These can be replaced with emoji or other Unicode characters. For example:
```json
"host": "🏠", "user": "👤", "os": "🐧", "kernel": "🧠", "uptime": "⏱️", etc.
```

### Modules

Enable/disable individual system information components.

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

## Example Configurations

### Minimal Configuration

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

### Image-Only Configuration

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

### Emoji Icons Configuration

```json
{
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
}
``` 