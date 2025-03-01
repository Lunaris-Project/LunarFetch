# LunarFetch

<div align="center">
  <img src="assets/image.png" alt="LunarFetch Logo" width="1080">
</div>

A highly customizable system information tool written in Go.

---

## English Documentation

### Project Structure

```
lunarfetch/
├── assets/
│   ├── config.json     # Default configuration file
│   └── logo.txt        # Sample ASCII art logo
├── scripts/
│   ├── install.sh      # Installation script
│   └── Makefile        # Make targets for common operations
├── main.go             # Main source code
├── go.mod              # Go module definition
├── go.sum              # Go module checksums
└── README.md           # This documentation
```

### Dependencies

LunarFetch requires the following Linux packages to function properly:

#### Core Dependencies

- `coreutils` - For basic system commands
- `procps` or `procps-ng` - For process information (uptime, memory)
- `util-linux` - For system utilities

#### Feature-specific Dependencies

- `lsb-release` - For OS detection
- `pacman` - For package counting on Arch-based systems
- `xorg-xrandr` or `xrandr` - For display resolution detection
- `xorg-xdpyinfo` or `xdpyinfo` - Alternative for resolution detection
- `sway` - For Sway window manager support
- `wlroots` - For Wayland compositor support
- `gsettings-desktop-schemas` - For theme detection
- `lm_sensors` - For hardware monitoring
- `pciutils` - For hardware detection (GPU)

#### Installing Dependencies

You can install the required dependencies using the provided Makefile targets:

##### For Arch Linux:

```bash
make -f scripts/Makefile install-deps-arch
```

##### For Debian/Ubuntu:

```bash
make -f scripts/Makefile install-deps-debian
```

##### Manual Installation on Arch Linux:

```bash
sudo pacman -S coreutils procps-ng util-linux lsb-release xorg-xrandr xorg-xdpyinfo pciutils gsettings-desktop-schemas
```

##### Manual Installation on Debian/Ubuntu:

```bash
sudo apt install coreutils procps lsb-release x11-xserver-utils mesa-utils pciutils gnome-settings-daemon
```

### Installation

#### Using the Installation Script

```bash
# Clone the repository
git clone https://github.com/Lunaris-Project/lunarfetch.git
cd lunarfetch

# Run the installation script
./scripts/install.sh
```

#### Using Make with Dependency Check

```bash
# Clone the repository
git clone https://github.com/Lunaris-Project/lunarfetch.git
cd lunarfetch

# Check dependencies first
make -f scripts/Makefile check-deps

# Install using make
make -f scripts/Makefile install
```

#### Manual Installation

```bash
# Clone the repository
git clone https://github.com/Lunaris-Project/lunarfetch.git
cd lunarfetch

# Build the binary
go build -o lunarfetch

# Create config directory and copy files
mkdir -p ~/.config/lunarfetch/logos
cp assets/config.json ~/.config/lunarfetch/config.json
cp assets/logo.txt ~/.config/lunarfetch/logos/moon.txt

# Install to your PATH
sudo mv lunarfetch /usr/local/bin/
```

#### Using Package Manager

Coming soon for various distributions.

### Usage

Simply run the command in your terminal:

```bash
lunarfetch
```

### Features

- **Customizable UI**: Change borders, separators, and layout
- **Random ASCII Art**: Display random ASCII art from a directory
- **Modular Information**: Show/hide specific system information
- **Fast Performance**: Written in Go for speed and efficiency

LunarFetch can display:

- OS information
- Hostname
- Kernel version
- Uptime
- Package count
- Shell
- Resolution
- Desktop environment
- Window manager theme
- GTK theme
- Icon theme
- Terminal
- CPU information
- GPU information
- Memory usage
- Disk usage
- Battery status

### Configuration

LunarFetch looks for its configuration file at `~/.config/lunarfetch/config.json`.

Create the directory if it doesn't exist:

```bash
mkdir -p ~/.config/lunarfetch
```

Create a basic configuration file:

```bash
cat > ~/.config/lunarfetch/config.json << 'EOF'
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
    "separator": "─"
  },
  "logo": {
    "enableLogo": true,
    "type": "file",
    "content": "",
    "location": "left",
    "logoPath": "~/.config/lunarfetch/logos"
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
EOF
```

### Customization

#### ASCII Art Logos

Create a directory for your ASCII art logos:

```bash
mkdir -p ~/.config/lunarfetch/logos
```

Add ASCII art files with `.txt` extension to this directory. LunarFetch will randomly select one each time it runs.

#### Configuration Options

##### Decorations

Change the box characters used to draw the information box:

- `topLeft`, `topRight`, `bottomLeft`, `bottomRight`: Corner characters
- `topEdge`, `bottomEdge`, `leftEdge`, `rightEdge`: Edge characters
- `separator`: Character used for the separator line

##### Logo Options

- `enableLogo`: Set to `true` to display a logo, `false` to hide it
- `type`: Currently supports `file` for reading from files
- `location`: Position of the logo (currently supports `left`)
- `logoPath`: Directory containing ASCII art logo files

##### Modules

Enable or disable specific information modules:

- `show_user`: Username
- `show_cpu`: CPU information
- `show_gpu`: GPU information
- `show_uptime`: System uptime
- `show_shell`: Current shell
- `show_memory`: Memory usage
- `show_packages`: Installed package count
- `show_os`: Operating system
- `show_host`: Hostname
- `show_kernel`: Kernel version
- `show_battery`: Battery status
- `show_disk`: Disk usage
- `show_resolution`: Screen resolution
- `show_de`: Desktop environment
- `show_wm_theme`: Window manager theme
- `show_theme`: GTK theme
- `show_icons`: Icon theme
- `show_terminal`: Terminal name

### Examples

#### Minimal Configuration

```json
{
  "decorations": {
    "topLeft": "┌",
    "topRight": "┐",
    "bottomLeft": "└",
    "bottomRight": "┘",
    "topEdge": "─",
    "bottomEdge": "─",
    "leftEdge": "│",
    "rightEdge": "│",
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
    "show_memory": true,
    "show_cpu": false,
    "show_gpu": false,
    "show_shell": false,
    "show_host": false,
    "show_battery": false,
    "show_disk": false,
    "show_resolution": false,
    "show_de": false,
    "show_wm_theme": false,
    "show_theme": false,
    "show_icons": false,
    "show_terminal": false
  }
}
```

#### Rounded Box Style

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
    "separator": "─"
  }
}
```

#### Double Line Style

```json
{
  "decorations": {
    "topLeft": "╔",
    "topRight": "╗",
    "bottomLeft": "╚",
    "bottomRight": "╝",
    "topEdge": "═",
    "bottomEdge": "═",
    "leftEdge": "║",
    "rightEdge": "║",
    "separator": "═"
  }
}
```

### Uninstallation

To uninstall LunarFetch:

```bash
# Remove just the binary
make -f scripts/Makefile uninstall

# Or remove everything including configuration
make -f scripts/Makefile purge
```
