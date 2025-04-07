package utils

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const (
	DefaultConfigDir  = "~/.config/lunarfetch"
	DefaultConfigFile = "config.json"
	DefaultLogoPath   = "~/.config/lunarfetch/logos"
	DefaultImagePath  = "~/.config/lunarfetch/images"
)

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

	Icons struct {
		Host       string `json:"host"`
		User       string `json:"user"`
		OS         string `json:"os"`
		Kernel     string `json:"kernel"`
		Uptime     string `json:"uptime"`
		Terminal   string `json:"terminal"`
		Shell      string `json:"shell"`
		Disk       string `json:"disk"`
		Memory     string `json:"memory"`
		Packages   string `json:"packages"`
		Battery    string `json:"battery"`
		GPU        string `json:"gpu"`
		CPU        string `json:"cpu"`
		Resolution string `json:"resolution"`
		DE         string `json:"de"`
		WMTheme    string `json:"wm_theme"`
		Theme      string `json:"theme"`
		Icons      string `json:"icons"`
	} `json:"icons"`

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

type ConfigLoader struct{}

func NewConfigLoader() *ConfigLoader {
	return &ConfigLoader{}
}

func (c *ConfigLoader) LoadConfig(paths ...string) (Config, error) {
	var configPath string

	if len(paths) > 0 && paths[0] != "" {
		configPath = paths[0]
	} else {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return DefaultConfig(), err
		}

		configDir := filepath.Join(homeDir, ".config", "lunarfetch")
		configPath = filepath.Join(configDir, "config.json")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return DefaultConfig(), err
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return DefaultConfig(), err
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return DefaultConfig(), err
	}

	config = applyConfigDefaults(config)

	return config, nil
}

func applyConfigDefaults(config Config) Config {
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

	return config
}

func DefaultConfig() Config {
	var config Config

	config.Decorations.TopLeft = "╭"
	config.Decorations.TopRight = "╮"
	config.Decorations.BottomLeft = "╰"
	config.Decorations.BottomRight = "╯"
	config.Decorations.TopEdge = "─"
	config.Decorations.BottomEdge = "─"
	config.Decorations.LeftEdge = "│"
	config.Decorations.RightEdge = "│"
	config.Decorations.Separator = ": "

	config.Logo.EnableLogo = true
	config.Logo.Type = "ascii"
	config.Logo.Location = "center"
	config.Logo.Position = "side"
	homeDir, _ := os.UserHomeDir()
	config.Logo.LogoPath = filepath.Join(homeDir, ".config", "lunarfetch", "logos")

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

	config.Display.ShowLogoFirst = true
	config.Display.ShowImageFirst = false

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
