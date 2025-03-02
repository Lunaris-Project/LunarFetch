package utils

import (
	"fmt"
	"strings"
	"sync"

	"lunarfetch/src/components"
)

// DisplayManager handles the initialization and display of system information
type DisplayManager struct {
	Config        Config
	InfoProviders map[string]components.InfoProvider
	infoCache     map[string]string
	cacheMutex    sync.RWMutex
}

// NewDisplayManager creates a new DisplayManager with the specified configuration
func NewDisplayManager(config Config) *DisplayManager {
	return &DisplayManager{
		Config:        config,
		InfoProviders: make(map[string]components.InfoProvider),
		infoCache:     make(map[string]string),
	}
}

// InitializeComponents initializes all system information components
func (d *DisplayManager) InitializeComponents() {
	// Initialize components with their names
	osInfo := &components.OSInfo{SystemInfo: components.SystemInfo{Name: "OS"}}
	hostInfo := &components.HostInfo{SystemInfo: components.SystemInfo{Name: "Host"}}
	kernelInfo := &components.KernelInfo{SystemInfo: components.SystemInfo{Name: "Kernel"}}
	uptimeInfo := &components.UptimeInfo{SystemInfo: components.SystemInfo{Name: "Uptime"}}
	shellInfo := &components.ShellInfo{SystemInfo: components.SystemInfo{Name: "Shell"}}
	memoryInfo := &components.MemoryInfo{SystemInfo: components.SystemInfo{Name: "Memory"}}
	packagesInfo := &components.PackagesInfo{SystemInfo: components.SystemInfo{Name: "Packages"}}
	userInfo := &components.UserInfo{SystemInfo: components.SystemInfo{Name: "User"}}
	diskInfo := &components.DiskInfo{SystemInfo: components.SystemInfo{Name: "Disk"}}
	batteryInfo := &components.BatteryInfo{SystemInfo: components.SystemInfo{Name: "Battery"}}
	cpuInfo := &components.CPUInfo{SystemInfo: components.SystemInfo{Name: "CPU"}}
	gpuInfo := &components.GPUInfo{SystemInfo: components.SystemInfo{Name: "GPU"}}
	resolutionInfo := &components.ResolutionInfo{SystemInfo: components.SystemInfo{Name: "Resolution"}}
	deInfo := &components.DEInfo{SystemInfo: components.SystemInfo{Name: "Desktop"}}
	terminalInfo := &components.TerminalInfo{SystemInfo: components.SystemInfo{Name: "Terminal"}}
	themeInfo := &components.ThemeInfo{SystemInfo: components.SystemInfo{Name: "Theme"}}
	wmThemeInfo := &components.WMThemeInfo{SystemInfo: components.SystemInfo{Name: "WM Theme"}}
	iconsInfo := &components.IconsInfo{SystemInfo: components.SystemInfo{Name: "Icons"}}

	// Add components to the map
	d.InfoProviders["OS"] = osInfo
	d.InfoProviders["Host"] = hostInfo
	d.InfoProviders["Kernel"] = kernelInfo
	d.InfoProviders["Uptime"] = uptimeInfo
	d.InfoProviders["Shell"] = shellInfo
	d.InfoProviders["Memory"] = memoryInfo
	d.InfoProviders["Packages"] = packagesInfo
	d.InfoProviders["User"] = userInfo
	d.InfoProviders["Disk"] = diskInfo
	d.InfoProviders["Battery"] = batteryInfo
	d.InfoProviders["CPU"] = cpuInfo
	d.InfoProviders["GPU"] = gpuInfo
	d.InfoProviders["Resolution"] = resolutionInfo
	d.InfoProviders["Desktop"] = deInfo
	d.InfoProviders["Terminal"] = terminalInfo
	d.InfoProviders["Theme"] = themeInfo
	d.InfoProviders["WM Theme"] = wmThemeInfo
	d.InfoProviders["Icons"] = iconsInfo
}

// GetInfoParallel gathers all system information in parallel
func (d *DisplayManager) GetInfoParallel() {
	// Create a list of components to gather info for
	var components []string

	if d.Config.Modules.ShowHost {
		components = append(components, "Host")
	}
	if d.Config.Modules.ShowUser {
		components = append(components, "User")
	}
	if d.Config.Modules.ShowOS {
		components = append(components, "OS")
	}
	if d.Config.Modules.ShowKernel {
		components = append(components, "Kernel")
	}
	if d.Config.Modules.ShowUptime {
		components = append(components, "Uptime")
	}
	if d.Config.Modules.ShowTerminal {
		components = append(components, "Terminal")
	}
	if d.Config.Modules.ShowShell {
		components = append(components, "Shell")
	}
	if d.Config.Modules.ShowDisk {
		components = append(components, "Disk")
	}
	if d.Config.Modules.ShowMemory {
		components = append(components, "Memory")
	}
	if d.Config.Modules.ShowPackages {
		components = append(components, "Packages")
	}
	if d.Config.Modules.ShowBattery {
		components = append(components, "Battery")
	}
	if d.Config.Modules.ShowGPU {
		components = append(components, "GPU")
	}
	if d.Config.Modules.ShowCPU {
		components = append(components, "CPU")
	}
	if d.Config.Modules.ShowResolution {
		components = append(components, "Resolution")
	}
	if d.Config.Modules.ShowWMTheme {
		components = append(components, "WM Theme")
	}
	if d.Config.Modules.ShowTheme {
		components = append(components, "Theme")
	}
	if d.Config.Modules.ShowIcons {
		components = append(components, "Icons")
	}
	if d.Config.Modules.ShowDE {
		components = append(components, "Desktop")
	}

	// Create a wait group to wait for all goroutines to finish
	var wg sync.WaitGroup
	wg.Add(len(components))

	// Gather info for each component in parallel
	for _, component := range components {
		go func(comp string) {
			defer wg.Done()
			info := d.InfoProviders[comp].GetInfo()
			d.cacheMutex.Lock()
			d.infoCache[comp] = info
			d.cacheMutex.Unlock()
		}(component)
	}

	// Wait for all goroutines to finish
	wg.Wait()
}

// GenerateContent generates the content to be displayed
func (d *DisplayManager) GenerateContent() string {
	// Gather all info in parallel first
	d.GetInfoParallel()

	var content strings.Builder

	// Use the cached info to build the content
	d.cacheMutex.RLock()
	defer d.cacheMutex.RUnlock()

	if d.Config.Modules.ShowHost {
		content.WriteString(fmt.Sprintf(" %s Host: %s\n", d.Config.Icons.Host, d.infoCache["Host"]))
	}
	if d.Config.Modules.ShowUser {
		content.WriteString(fmt.Sprintf(" %s User: %s\n", d.Config.Icons.User, d.infoCache["User"]))
	}
	if d.Config.Modules.ShowOS {
		content.WriteString(fmt.Sprintf(" %s OS: %s\n", d.Config.Icons.OS, d.infoCache["OS"]))
	}
	if d.Config.Modules.ShowKernel {
		content.WriteString(fmt.Sprintf(" %s Kernel: %s\n", d.Config.Icons.Kernel, d.infoCache["Kernel"]))
	}
	if d.Config.Modules.ShowUptime {
		content.WriteString(fmt.Sprintf(" %s Uptime: %s\n", d.Config.Icons.Uptime, d.infoCache["Uptime"]))
	}
	if d.Config.Modules.ShowTerminal {
		content.WriteString(fmt.Sprintf(" %s Terminal: %s\n", d.Config.Icons.Terminal, d.infoCache["Terminal"]))
	}
	if d.Config.Modules.ShowShell {
		content.WriteString(fmt.Sprintf(" %s Shell: %s\n", d.Config.Icons.Shell, d.infoCache["Shell"]))
	}
	if d.Config.Modules.ShowDisk {
		content.WriteString(fmt.Sprintf(" %s Disk: %s\n", d.Config.Icons.Disk, d.infoCache["Disk"]))
	}
	if d.Config.Modules.ShowMemory {
		content.WriteString(fmt.Sprintf(" %s Memory: %s\n", d.Config.Icons.Memory, d.infoCache["Memory"]))
	}
	if d.Config.Modules.ShowPackages {
		content.WriteString(fmt.Sprintf(" %s Packages: %s\n", d.Config.Icons.Packages, d.infoCache["Packages"]))
	}
	if d.Config.Modules.ShowBattery {
		content.WriteString(fmt.Sprintf(" %s Battery: %s\n", d.Config.Icons.Battery, d.infoCache["Battery"]))
	}
	if d.Config.Modules.ShowGPU {
		content.WriteString(fmt.Sprintf(" %s GPU: %s\n", d.Config.Icons.GPU, d.infoCache["GPU"]))
	}
	if d.Config.Modules.ShowCPU {
		content.WriteString(fmt.Sprintf(" %s CPU: %s\n", d.Config.Icons.CPU, d.infoCache["CPU"]))
	}
	if d.Config.Modules.ShowResolution {
		content.WriteString(fmt.Sprintf(" %s Resolution: %s\n", d.Config.Icons.Resolution, d.infoCache["Resolution"]))
	}
	if d.Config.Modules.ShowWMTheme {
		content.WriteString(fmt.Sprintf(" %s WM Theme: %s\n", d.Config.Icons.WMTheme, d.infoCache["WM Theme"]))
	}
	if d.Config.Modules.ShowTheme {
		content.WriteString(fmt.Sprintf(" %s Theme: %s\n", d.Config.Icons.Theme, d.infoCache["Theme"]))
	}
	if d.Config.Modules.ShowIcons {
		content.WriteString(fmt.Sprintf(" %s Icons: %s\n", d.Config.Icons.Icons, d.infoCache["Icons"]))
	}
	if d.Config.Modules.ShowDE {
		content.WriteString(fmt.Sprintf(" %s Desktop: %s\n", d.Config.Icons.DE, d.infoCache["Desktop"]))
	}

	content.WriteString(strings.Repeat(d.Config.Decorations.Separator, 30) + "\n")

	return strings.TrimSpace(content.String())
}

// Display displays the system information
func (d *DisplayManager) Display() string {
	// Create a box drawer
	boxConfig := BoxConfig{
		TopLeft:     d.Config.Decorations.TopLeft,
		TopRight:    d.Config.Decorations.TopRight,
		BottomLeft:  d.Config.Decorations.BottomLeft,
		BottomRight: d.Config.Decorations.BottomRight,
		TopEdge:     d.Config.Decorations.TopEdge,
		BottomEdge:  d.Config.Decorations.BottomEdge,
		LeftEdge:    d.Config.Decorations.LeftEdge,
		RightEdge:   d.Config.Decorations.RightEdge,
		Separator:   d.Config.Decorations.Separator,
	}
	boxDrawer := NewBoxDrawer(boxConfig)

	// Generate content and draw box
	content := d.GenerateContent()
	return boxDrawer.Draw(content)
}
