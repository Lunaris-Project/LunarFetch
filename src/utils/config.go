package utils

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Config represents the application configuration
type Config struct {
	Decorations struct {
		TopLeft     string `json:"topLeft"`
		TopRight    string `json:"topRight"`
		BottomLeft  string `json:"bottomLeft"`
		BottomRight string `json:"bottomRight"`
		TopEdge     string `json:"topEdge"`
		BottomEdge  string `json:"bottomEdge"`
		LeftEdge    string `json:"leftEdge"`
		RightEdge   string `json:"rightEdge"`
		Separator   string `json:"separator"`
	} `json:"decorations"`
	Logo struct {
		EnableLogo bool   `json:"enableLogo"`
		Type       string `json:"type"`
		Content    string `json:"content"`
		Location   string `json:"location"`
		LogoPath   string `json:"logoPath"`
		Position   string `json:"position"`
	} `json:"logo"`
	Image struct {
		EnableImage    bool   `json:"enableImage"`
		Enabled        bool   `json:"enabled"`
		Random         bool   `json:"random"`
		ImagePath      string `json:"imagePath"`
		Width          int    `json:"width"`
		Height         int    `json:"height"`
		RenderMode     string `json:"renderMode"`
		DitherMode     string `json:"ditherMode"`
		TerminalOutput bool   `json:"terminalOutput"`
		DisplayMode    string `json:"displayMode"`
		Protocol       string `json:"protocol"`
		Scale          int    `json:"scale"`
		Offset         int    `json:"offset"`
		Background     string `json:"background"`
		Position       string `json:"position"`
	} `json:"image"`
	Display struct {
		ShowLogoFirst  bool `json:"showLogoFirst"`
		ShowImageFirst bool `json:"showImageFirst"`
	} `json:"display"`
	Modules struct {
		ShowUser       bool `json:"show_user"`
		ShowCPU        bool `json:"show_cpu"`
		ShowGPU        bool `json:"show_gpu"`
		ShowUptime     bool `json:"show_uptime"`
		ShowShell      bool `json:"show_shell"`
		ShowMemory     bool `json:"show_memory"`
		ShowPackages   bool `json:"show_packages"`
		ShowOS         bool `json:"show_os"`
		ShowHost       bool `json:"show_host"`
		ShowKernel     bool `json:"show_kernel"`
		ShowBattery    bool `json:"show_battery"`
		ShowDisk       bool `json:"show_disk"`
		ShowResolution bool `json:"show_resolution"`
		ShowDE         bool `json:"show_de"`
		ShowWMTheme    bool `json:"show_wm_theme"`
		ShowTheme      bool `json:"show_theme"`
		ShowIcons      bool `json:"show_icons"`
		ShowTerminal   bool `json:"show_terminal"`
	} `json:"modules"`
}

// ConfigLoader handles loading configuration from files
type ConfigLoader struct{}

// NewConfigLoader creates a new ConfigLoader
func NewConfigLoader() *ConfigLoader {
	return &ConfigLoader{}
}

// LoadConfig loads configuration from the default location or specified file
func (c *ConfigLoader) LoadConfig(filename ...string) (Config, error) {
	var configPath string

	if len(filename) > 0 && filename[0] != "" {
		configPath = filename[0]
	} else {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return Config{}, err
		}
		configPath = filepath.Join(homeDir, ".config", "lunarfetch", "config.json")
	}

	// Create default config
	config := DefaultConfig()

	// Try to open and read the config file
	file, err := os.Open(configPath)
	if err != nil {
		// If file doesn't exist, create the directory structure and return default config
		if os.IsNotExist(err) {
			// Create config directory if it doesn't exist
			configDir := filepath.Dir(configPath)
			if err := os.MkdirAll(configDir, 0755); err != nil {
				return config, err
			}

			// Save default config
			if err := c.SaveConfig(config, configPath); err != nil {
				return config, err
			}

			return config, nil
		}
		return config, err
	}
	defer file.Close()

	// Decode the config file
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return config, err
	}

	return config, nil
}

// SaveConfig saves the configuration to the specified file
func (c *ConfigLoader) SaveConfig(config Config, filename string) error {
	// Create directory if it doesn't exist
	configDir := filepath.Dir(filename)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	// Marshal config to JSON
	configJSON, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	// Write to file
	return os.WriteFile(filename, configJSON, 0644)
}

// GetBoxConfig extracts box drawing configuration from the main config
func (c *ConfigLoader) GetBoxConfig(config Config) BoxConfig {
	return BoxConfig{
		TopLeft:     config.Decorations.TopLeft,
		TopRight:    config.Decorations.TopRight,
		BottomLeft:  config.Decorations.BottomLeft,
		BottomRight: config.Decorations.BottomRight,
		TopEdge:     config.Decorations.TopEdge,
		BottomEdge:  config.Decorations.BottomEdge,
		LeftEdge:    config.Decorations.LeftEdge,
		RightEdge:   config.Decorations.RightEdge,
		Separator:   config.Decorations.Separator,
	}
}

// DefaultConfig returns a default configuration
func DefaultConfig() Config {
	var config Config

	// Default decorations
	config.Decorations.TopLeft = "╭"
	config.Decorations.TopRight = "╮"
	config.Decorations.BottomLeft = "╰"
	config.Decorations.BottomRight = "╯"
	config.Decorations.TopEdge = "─"
	config.Decorations.BottomEdge = "─"
	config.Decorations.LeftEdge = "│"
	config.Decorations.RightEdge = "│"
	config.Decorations.Separator = ": "

	// Default logo settings
	config.Logo.EnableLogo = true
	config.Logo.Type = "ascii"
	config.Logo.LogoPath = "~/.config/lunarfetch/logos"
	config.Logo.Position = "side"

	// Default image settings
	config.Image.Enabled = false
	config.Image.EnableImage = false
	config.Image.Random = false
	config.Image.ImagePath = "~/.config/lunarfetch/images"
	config.Image.Width = 40
	config.Image.Height = 20
	config.Image.RenderMode = "detailed"
	config.Image.DitherMode = "none"
	config.Image.TerminalOutput = false
	config.Image.Protocol = "auto"
	config.Image.Scale = 1
	config.Image.Offset = 2
	config.Image.Position = "side"

	// Default display settings
	config.Display.ShowLogoFirst = true
	config.Display.ShowImageFirst = false

	// Default modules
	config.Modules.ShowUser = true
	config.Modules.ShowOS = true
	config.Modules.ShowHost = true
	config.Modules.ShowKernel = true
	config.Modules.ShowUptime = true
	config.Modules.ShowPackages = true
	config.Modules.ShowShell = true
	config.Modules.ShowResolution = true
	config.Modules.ShowDE = true
	config.Modules.ShowWMTheme = true
	config.Modules.ShowTheme = true
	config.Modules.ShowIcons = true
	config.Modules.ShowTerminal = true
	config.Modules.ShowCPU = true
	config.Modules.ShowGPU = true
	config.Modules.ShowMemory = true
	config.Modules.ShowBattery = true
	config.Modules.ShowDisk = true

	return config
}
