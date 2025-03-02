package utils

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Default configuration values
const (
	DefaultConfigDir  = "~/.config/lunarfetch"
	DefaultConfigFile = "config.json"
	DefaultLogoPath   = "~/.config/lunarfetch/logos"
	DefaultImagePath  = "~/.config/lunarfetch/images"
)

// Config represents the application configuration
type Config struct {
	// Box drawing characters for display
	Decorations struct {
		TopLeft     string `json:"topLeft"`     // Top-left corner character
		TopRight    string `json:"topRight"`    // Top-right corner character
		BottomLeft  string `json:"bottomLeft"`  // Bottom-left corner character
		BottomRight string `json:"bottomRight"` // Bottom-right corner character
		TopEdge     string `json:"topEdge"`     // Top edge character
		BottomEdge  string `json:"bottomEdge"`  // Bottom edge character
		LeftEdge    string `json:"leftEdge"`    // Left edge character
		RightEdge   string `json:"rightEdge"`   // Right edge character
		Separator   string `json:"separator"`   // Separator between key and value
	} `json:"decorations"`

	// Logo configuration
	Logo struct {
		EnableLogo bool   `json:"enableLogo"` // Whether to display a logo
		Type       string `json:"type"`       // Logo type (ascii or file)
		Content    string `json:"content"`    // Custom ASCII content
		Location   string `json:"location"`   // Text alignment (center, left, right)
		LogoPath   string `json:"logoPath"`   // Path to logo files
		Position   string `json:"position"`   // Position (side or above)
	} `json:"logo"`

	// Image configuration
	Image struct {
		EnableImage    bool   `json:"enableImage"`    // Whether to display an image
		Enabled        bool   `json:"enabled"`        // Legacy field for backward compatibility
		Random         bool   `json:"random"`         // Whether to select a random image
		ImagePath      string `json:"imagePath"`      // Path to image file or directory
		Width          int    `json:"width"`          // Width in terminal characters
		Height         int    `json:"height"`         // Height in terminal characters
		RenderMode     string `json:"renderMode"`     // Rendering mode
		DitherMode     string `json:"ditherMode"`     // Dithering algorithm
		TerminalOutput bool   `json:"terminalOutput"` // Whether to output directly to terminal
		DisplayMode    string `json:"displayMode"`    // Display mode
		Protocol       string `json:"protocol"`       // Image protocol
		Scale          int    `json:"scale"`          // Image scaling factor
		Offset         int    `json:"offset"`         // Offset from terminal edge
		Background     string `json:"background"`     // Background color
		Position       string `json:"position"`       // Position (side or above)
	} `json:"image"`

	// Display order configuration
	Display struct {
		ShowLogoFirst  bool `json:"showLogoFirst"`  // Whether to show logo before system info
		ShowImageFirst bool `json:"showImageFirst"` // Whether to show image before system info
	} `json:"display"`

	// Icons for system information
	Icons struct {
		Host       string `json:"host"`       // Host icon
		User       string `json:"user"`       // User icon
		OS         string `json:"os"`         // OS icon
		Kernel     string `json:"kernel"`     // Kernel icon
		Uptime     string `json:"uptime"`     // Uptime icon
		Terminal   string `json:"terminal"`   // Terminal icon
		Shell      string `json:"shell"`      // Shell icon
		Disk       string `json:"disk"`       // Disk icon
		Memory     string `json:"memory"`     // Memory icon
		Packages   string `json:"packages"`   // Packages icon
		Battery    string `json:"battery"`    // Battery icon
		GPU        string `json:"gpu"`        // GPU icon
		CPU        string `json:"cpu"`        // CPU icon
		Resolution string `json:"resolution"` // Resolution icon
		DE         string `json:"de"`         // Desktop Environment icon
		WMTheme    string `json:"wm_theme"`   // Window Manager Theme icon
		Theme      string `json:"theme"`      // Theme icon
		Icons      string `json:"icons"`      // Icons icon
	} `json:"icons"`

	// Module visibility configuration
	Modules struct {
		ShowUser       bool `json:"show_user"`       // Whether to show user information
		ShowCPU        bool `json:"show_cpu"`        // Whether to show CPU information
		ShowGPU        bool `json:"show_gpu"`        // Whether to show GPU information
		ShowUptime     bool `json:"show_uptime"`     // Whether to show uptime information
		ShowShell      bool `json:"show_shell"`      // Whether to show shell information
		ShowMemory     bool `json:"show_memory"`     // Whether to show memory information
		ShowPackages   bool `json:"show_packages"`   // Whether to show packages information
		ShowOS         bool `json:"show_os"`         // Whether to show OS information
		ShowHost       bool `json:"show_host"`       // Whether to show host information
		ShowKernel     bool `json:"show_kernel"`     // Whether to show kernel information
		ShowBattery    bool `json:"show_battery"`    // Whether to show battery information
		ShowDisk       bool `json:"show_disk"`       // Whether to show disk information
		ShowResolution bool `json:"show_resolution"` // Whether to show resolution information
		ShowDE         bool `json:"show_de"`         // Whether to show desktop environment information
		ShowWMTheme    bool `json:"show_wm_theme"`   // Whether to show window manager theme information
		ShowTheme      bool `json:"show_theme"`      // Whether to show theme information
		ShowIcons      bool `json:"show_icons"`      // Whether to show icons information
		ShowTerminal   bool `json:"show_terminal"`   // Whether to show terminal information
	} `json:"modules"`
}

// ConfigLoader handles loading configuration from files
type ConfigLoader struct{}

// NewConfigLoader creates a new ConfigLoader
func NewConfigLoader() *ConfigLoader {
	return &ConfigLoader{}
}

// LoadConfig loads the configuration from the specified path or default location
func (c *ConfigLoader) LoadConfig(paths ...string) (Config, error) {
	var configPath string

	// If a path is provided, use it
	if len(paths) > 0 && paths[0] != "" {
		configPath = paths[0]
	} else {
		// Otherwise, use the default path
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return DefaultConfig(), err
		}

		configDir := filepath.Join(homeDir, ".config", "lunarfetch")
		configPath = filepath.Join(configDir, "config.json")
	}

	// Check if the file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// If the file doesn't exist, create a default config
		return DefaultConfig(), err
	}

	// Read the file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return DefaultConfig(), err
	}

	// Parse the JSON
	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return DefaultConfig(), err
	}

	// Apply any necessary fixes or defaults to the loaded config
	config = applyConfigDefaults(config)

	return config, nil
}

// applyConfigDefaults ensures all required fields have values
func applyConfigDefaults(config Config) Config {
	// Apply defaults for decorations if not set
	if config.Decorations.TopLeft == "" {
		config.Decorations.TopLeft = "╭"
	}
	if config.Decorations.TopRight == "" {
		config.Decorations.TopRight = "╮"
	}
	if config.Decorations.BottomLeft == "" {
		config.Decorations.BottomLeft = "╰"
	}
	if config.Decorations.BottomRight == "" {
		config.Decorations.BottomRight = "╯"
	}
	if config.Decorations.TopEdge == "" {
		config.Decorations.TopEdge = "─"
	}
	if config.Decorations.BottomEdge == "" {
		config.Decorations.BottomEdge = "─"
	}
	if config.Decorations.LeftEdge == "" {
		config.Decorations.LeftEdge = "│"
	}
	if config.Decorations.RightEdge == "" {
		config.Decorations.RightEdge = "│"
	}
	if config.Decorations.Separator == "" {
		config.Decorations.Separator = ": "
	}

	// Apply defaults for logo if not set
	if config.Logo.LogoPath == "" {
		homeDir, _ := os.UserHomeDir()
		config.Logo.LogoPath = filepath.Join(homeDir, ".config", "lunarfetch", "logos")
	}
	if config.Logo.Type == "" {
		config.Logo.Type = "ascii"
	}
	if config.Logo.Location == "" {
		config.Logo.Location = "center"
	}
	if config.Logo.Position == "" {
		config.Logo.Position = "side"
	}

	// Apply defaults for image if not set
	if config.Image.ImagePath == "" {
		homeDir, _ := os.UserHomeDir()
		config.Image.ImagePath = filepath.Join(homeDir, ".config", "lunarfetch", "images")
	}
	if config.Image.Width <= 0 {
		config.Image.Width = 40
	}
	if config.Image.Height <= 0 {
		config.Image.Height = 20
	}
	if config.Image.RenderMode == "" {
		config.Image.RenderMode = "detailed"
	}
	if config.Image.DitherMode == "" {
		config.Image.DitherMode = "floyd-steinberg"
	}
	if config.Image.DisplayMode == "" {
		config.Image.DisplayMode = "block"
	}
	if config.Image.Protocol == "" {
		config.Image.Protocol = "auto"
	}
	if config.Image.Scale <= 0 {
		config.Image.Scale = 1
	}
	if config.Image.Background == "" {
		config.Image.Background = "transparent"
	}
	if config.Image.Position == "" {
		config.Image.Position = "side"
	}

	// For backward compatibility

	return config
}

// DefaultConfig returns the default configuration
func DefaultConfig() Config {
	var config Config

	// Set default decorations
	config.Decorations.TopLeft = "╭"
	config.Decorations.TopRight = "╮"
	config.Decorations.BottomLeft = "╰"
	config.Decorations.BottomRight = "╯"
	config.Decorations.TopEdge = "─"
	config.Decorations.BottomEdge = "─"
	config.Decorations.LeftEdge = "│"
	config.Decorations.RightEdge = "│"
	config.Decorations.Separator = ": "

	// Set default logo settings
	config.Logo.EnableLogo = true
	config.Logo.Type = "ascii"
	config.Logo.Location = "center"
	config.Logo.Position = "side"
	homeDir, _ := os.UserHomeDir()
	config.Logo.LogoPath = filepath.Join(homeDir, ".config", "lunarfetch", "logos")

	// Set default image settings
	config.Image.EnableImage = true
	config.Image.Enabled = true
	config.Image.Random = true
	config.Image.ImagePath = filepath.Join(homeDir, ".config", "lunarfetch", "images")
	config.Image.Width = 40
	config.Image.Height = 20
	config.Image.RenderMode = "detailed"
	config.Image.DitherMode = "floyd-steinberg"
	config.Image.TerminalOutput = false
	config.Image.DisplayMode = "block"
	config.Image.Protocol = "auto"
	config.Image.Scale = 1
	config.Image.Offset = 2
	config.Image.Background = "transparent"
	config.Image.Position = "side"

	// Set default display settings
	config.Display.ShowLogoFirst = true
	config.Display.ShowImageFirst = false

	// Set default icons
	config.Icons.Host = "󰒋"
	config.Icons.User = "󰀄"
	config.Icons.OS = "󰣇"
	config.Icons.Kernel = "󰣇"
	config.Icons.Uptime = "󰔟"
	config.Icons.Terminal = "󰆍"
	config.Icons.Shell = "󰆍"
	config.Icons.Disk = "󰋊"
	config.Icons.Memory = "󰍛"
	config.Icons.Packages = "󰏗"
	config.Icons.Battery = "󰂄"
	config.Icons.GPU = "󰢮"
	config.Icons.CPU = "󰘚"
	config.Icons.Resolution = "󰍹"
	config.Icons.DE = "󰧨"
	config.Icons.WMTheme = "󰏘"
	config.Icons.Theme = "󰔯"
	config.Icons.Icons = "󰀻"

	// Set default module visibility
	config.Modules.ShowUser = true
	config.Modules.ShowCPU = true
	config.Modules.ShowGPU = true
	config.Modules.ShowUptime = true
	config.Modules.ShowShell = true
	config.Modules.ShowMemory = true
	config.Modules.ShowPackages = true
	config.Modules.ShowOS = true
	config.Modules.ShowHost = true
	config.Modules.ShowKernel = true
	config.Modules.ShowBattery = true
	config.Modules.ShowDisk = true
	config.Modules.ShowResolution = true
	config.Modules.ShowDE = true
	config.Modules.ShowWMTheme = true
	config.Modules.ShowTheme = true
	config.Modules.ShowIcons = true
	config.Modules.ShowTerminal = true

	return config
}
