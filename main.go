package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
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
	} `json:"logo"`
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

func loadConfig(filename string) (Config, error) {
	var config Config
	file, err := os.Open(filename)
	if err != nil {
		return config, err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	return config, err
}

func drawBox(config Config, content string) string {
	lines := strings.Split(content, "\n")
	maxLen := 0
	for _, line := range lines {
		if len(line) > maxLen {
			maxLen = len(line)
		}
	}

	var box strings.Builder

	box.WriteString(config.Decorations.TopLeft)
	box.WriteString(strings.Repeat(config.Decorations.TopEdge, maxLen+2))
	box.WriteString(config.Decorations.TopRight + "\n")

	for _, line := range lines {
		box.WriteString(config.Decorations.LeftEdge + " ")
		box.WriteString(fmt.Sprintf("%-*s", maxLen, line))
		box.WriteString(" " + config.Decorations.RightEdge + "\n")
	}

	box.WriteString(config.Decorations.BottomLeft)
	box.WriteString(strings.Repeat(config.Decorations.BottomEdge, maxLen+2))
	box.WriteString(config.Decorations.BottomRight + "\n")

	return box.String()
}

func getRandomLogo(logoPath string) (string, error) {
	expandedPath, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	logoPath = strings.Replace(logoPath, "~", expandedPath, 1)

	files, err := ioutil.ReadDir(logoPath)
	if err != nil {
		log.Printf("Error reading directory: %v\n", err)
		return "", err
	}

	var logoFiles []string
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".txt" {
			logoFiles = append(logoFiles, filepath.Join(logoPath, file.Name()))
		}
	}

	if len(logoFiles) == 0 {
		return "", fmt.Errorf("no logo files found")
	}

	rand.Seed(time.Now().UnixNano())
	randomFile := logoFiles[rand.Intn(len(logoFiles))]
	content, err := ioutil.ReadFile(randomFile)
	if err != nil {
		log.Printf("Error reading logo file: %v\n", err)
		return "", err
	}

	return string(content), nil
}

func writeLogo(config Config, writer *strings.Builder) error {
	logo, err := getRandomLogo(config.Logo.LogoPath)
	if err != nil {
		return err
	}
	writer.WriteString(logo)
	return nil
}

func displayInfo(config Config) {

	info := make(map[string]string)

	info["OS"] = getOS()
	info["User"] = getUser()
	info["Host"] = getHostname()
	info["Kernel"] = getKernel()
	info["Uptime"] = getUptime()
	info["Packages"] = getPackages()
	info["Shell"] = getShell()
	info["Resolution"] = getResolution()

	// Combine DE and WM into a single Desktop entry
	if config.Modules.ShowDE {
		info["Desktop"] = getDE()
	}

	var content strings.Builder

	if config.Modules.ShowHost {
		content.WriteString(fmt.Sprintf(" 󰒋 Host: %s\n", getHostname()))
	}
	if config.Modules.ShowUser {
		content.WriteString(fmt.Sprintf(" 󰀄 User: %s\n", getUser()))
	}
	if config.Modules.ShowOS {
		content.WriteString(fmt.Sprintf(" 󰣇 OS: %s\n", getOS()))
	}
	if config.Modules.ShowKernel {
		content.WriteString(fmt.Sprintf(" 󰣇 Kernel: %s\n", getKernel()))
	}
	if config.Modules.ShowUptime {
		content.WriteString(fmt.Sprintf(" 󰔟 Uptime: %s\n", getUptime()))
	}
	if config.Modules.ShowTerminal {
		content.WriteString(fmt.Sprintf(" 󰆍 Terminal: %s\n", getTerminal()))
	}
	if config.Modules.ShowShell {
		content.WriteString(fmt.Sprintf(" 󰆍 Shell: %s\n", getShell()))
	}
	if config.Modules.ShowDisk {
		content.WriteString(fmt.Sprintf(" 󰋊 Disk: %s\n", getDiskUsage()))
	}
	if config.Modules.ShowMemory {
		content.WriteString(fmt.Sprintf(" 󰍛 Memory: %s\n", getMemory()))
	}
	if config.Modules.ShowPackages {
		content.WriteString(fmt.Sprintf(" 󰏗 Packages: %s\n", getPackages()))
	}
	if config.Modules.ShowBattery {
		content.WriteString(fmt.Sprintf(" 󰂄 Battery: %s\n", getBattery()))
	}
	if config.Modules.ShowGPU {
		content.WriteString(fmt.Sprintf(" 󰢮 GPU: %s\n", getGPU()))
	}
	if config.Modules.ShowCPU {
		content.WriteString(fmt.Sprintf(" 󰘚 CPU: %s\n", getCPU()))
	}
	if config.Modules.ShowResolution {
		content.WriteString(fmt.Sprintf(" 󰍹 Resolution: %s\n", getResolution()))
	}
	if config.Modules.ShowWMTheme {
		content.WriteString(fmt.Sprintf(" 󰏘 WM Theme: %s\n", getWMTheme()))
	}
	if config.Modules.ShowTheme {
		content.WriteString(fmt.Sprintf(" 󰔯 Theme: %s\n", getTheme()))
	}
	if config.Modules.ShowIcons {
		content.WriteString(fmt.Sprintf(" 󰀻 Icons: %s\n", getIcons()))
	}
	if config.Modules.ShowDE {
		content.WriteString(fmt.Sprintf(" 󰧨 Desktop: %s\n", getDE()))
	}

	content.WriteString(strings.Repeat(config.Decorations.Separator, 30) + "\n")

	fmt.Print(drawBox(config, strings.TrimSpace(content.String())))
}

func getOS() string {
	out, err := exec.Command("lsb_release", "-si").Output()
	if err != nil {
		out, err = exec.Command("grep", "^NAME=", "/etc/os-release").Output()
		if err != nil {
			return "Unknown"
		}
		return strings.Trim(strings.TrimPrefix(strings.TrimSpace(string(out)), "NAME=\""), "\"")
	}
	return strings.TrimSpace(string(out))
}

func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "Unknown"
	}
	return hostname
}

func getKernel() string {
	out, err := exec.Command("uname", "-r").Output()
	if err != nil {
		return "Unknown"
	}
	return strings.TrimSpace(string(out))
}

func getUptime() string {
	out, err := exec.Command("uptime", "-p").Output()
	if err != nil {
		return "Unknown"
	}
	return strings.TrimSpace(string(out))
}

func getPackages() string {
	out, err := exec.Command("pacman", "-Qq").Output()
	if err != nil {
		return "Unknown"
	}
	packages := strings.Split(string(out), "\n")
	return fmt.Sprintf("%d", len(packages))
}

func getShell() string {
	shell := os.Getenv("SHELL")
	if shell == "" {
		return "Unknown"
	}
	return filepath.Base(shell)
}

func getResolution() string {
	out, err := exec.Command("xrandr").Output()
	if err == nil {
		lines := strings.Split(string(out), "\n")
		for _, line := range lines {
			if strings.Contains(line, " connected") {
				fields := strings.Fields(line)
				for _, field := range fields {
					if strings.Contains(field, "x") {
						return field
					}
				}
			}
		}
	}

	if out, err := exec.Command("swaymsg", "-t", "get_outputs").Output(); err == nil {
		var outputs []struct {
			CurrentMode struct {
				Width  int `json:"width"`
				Height int `json:"height"`
			} `json:"current_mode"`
		}
		if err := json.Unmarshal(out, &outputs); err == nil && len(outputs) > 0 {
			return fmt.Sprintf("%dx%d", outputs[0].CurrentMode.Width, outputs[0].CurrentMode.Height)
		}
	}

	if out, err := exec.Command("wlr-randr").Output(); err == nil {
		lines := strings.Split(string(out), "\n")
		for _, line := range lines {
			if strings.Contains(line, "current") {
				fields := strings.Fields(line)
				for _, field := range fields {
					if strings.Contains(field, "x") {
						return field
					}
				}
			}
		}
	}

	if out, err := exec.Command("xdpyinfo").Output(); err == nil {
		lines := strings.Split(string(out), "\n")
		for _, line := range lines {
			if strings.Contains(line, "dimensions:") {
				fields := strings.Fields(line)
				return fields[1]
			}
		}
	}

	return "Unknown"
}

func getDiskUsage() string {
	out, err := exec.Command("df", "-B1").Output()
	if err != nil {
		return "Unknown"
	}
	lines := strings.Split(string(out), "\n")
	var totalUsed uint64
	var totalSize uint64
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) > 3 && fields[0] != "Filesystem" {
			used, _ := strconv.ParseUint(fields[2], 10, 64)
			size, _ := strconv.ParseUint(fields[1], 10, 64)
			totalUsed += used
			totalSize += size
		}
	}
	return fmt.Sprintf("%s / %s", formatBytes(totalUsed), formatBytes(totalSize))
}

func getBattery() string {
	if _, err := os.Stat("/sys/class/power_supply/BAT0"); err != nil {
		return "No battery"
	}

	capacity, err := os.ReadFile("/sys/class/power_supply/BAT0/capacity")
	if err != nil {
		return "Unknown"
	}

	status, err := os.ReadFile("/sys/class/power_supply/BAT0/status")
	if err != nil {
		return "Unknown"
	}

	return fmt.Sprintf("%s%% (%s)", strings.TrimSpace(string(capacity)), strings.TrimSpace(string(status)))
}

func getDE() string {
	de := os.Getenv("XDG_CURRENT_DESKTOP")
	if de != "" {
		return de
	}
	return "Unknown"
}

func getWMTheme() string {
	out, err := exec.Command("gsettings", "get", "org.gnome.desktop.wm.preferences", "theme").Output()
	if err != nil {
		return "Unknown"
	}
	return strings.TrimSpace(string(out))
}

func getTheme() string {
	out, err := exec.Command("gsettings", "get", "org.gnome.desktop.interface", "gtk-theme").Output()
	if err != nil {
		return "Unknown"
	}
	return strings.TrimSpace(string(out))
}

func getIcons() string {
	out, err := exec.Command("gsettings", "get", "org.gnome.desktop.interface", "icon-theme").Output()
	if err != nil {
		return "Unknown"
	}
	return strings.TrimSpace(string(out))
}

func getTerminal() string {
	term := os.Getenv("TERM")
	if term == "" {
		return "Unknown"
	}
	return term
}

func getCPU() string {
	out, err := exec.Command("lscpu").Output()
	if err != nil {
		return "Unknown"
	}
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if strings.Contains(line, "Model name:") {
			fields := strings.Fields(line)
			return strings.Join(fields[2:], " ")
		}
	}
	return "Unknown"
}

func getGPU() string {
	out, err := exec.Command("lspci").Output()
	if err != nil {
		return "Unknown"
	}
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if strings.Contains(line, "VGA") || strings.Contains(line, "3D") {
			fields := strings.Fields(line)
			return strings.Join(fields[4:], " ")
		}
	}
	return "Unknown"
}

func getMemory() string {
	out, err := exec.Command("free", "-m").Output()
	if err != nil {
		return "Unknown"
	}
	lines := strings.Split(string(out), "\n")
	fields := strings.Fields(lines[1])
	used := fields[2]
	total := fields[1]
	return fmt.Sprintf("%sMiB / %sMiB", used, total)
}

func getUser() string {
	user := os.Getenv("USER")
	if user == "" {
		return "Unknown"
	}
	return user
}

func formatBytes(bytes uint64) string {
	units := []string{"B", "KB", "MB", "GB", "TB"}
	var unit string
	var value float64 = float64(bytes)

	for _, u := range units {
		if value < 1024 {
			unit = u
			break
		}
		value /= 1024
	}
	return fmt.Sprintf("%.2f %s", value, unit)
}

func main() {
	// Check if any command-line arguments were provided
	if len(os.Args) > 1 {
		// If arguments exist, pass them to the CLI handler
		handleCLICommands()
		return
	}

	// No arguments, run the normal fetcher
	configPath := ".config/lunarfetch/config.json"
	configFilePath, err := os.UserHomeDir()
	if err != nil {
		log.Println("Error expanding config path:", err)
		return
	}
	configFilePath = filepath.Join(configFilePath, configPath)

	config, err := loadConfig(configFilePath)
	if err != nil {
		log.Println("Error loading config:", err)
		return
	}

	var content strings.Builder

	err = writeLogo(config, &content)
	if err != nil {
		log.Println(err)
	}

	fmt.Print(content.String() + "\n")
	displayInfo(config)
}

// handleCLICommands executes the CLI functionality from lunarfetch-cli.go
func handleCLICommands() {
	// Get the path to the current executable
	exePath, err := os.Executable()
	if err != nil {
		fmt.Printf("Error getting executable path: %v\n", err)
		os.Exit(1)
	}

	// Get the directory containing the executable
	exeDir := filepath.Dir(exePath)

	// Path to lunarfetch-cli.go relative to the executable
	cliPath := filepath.Join(exeDir, "scripts", "lunarfetch-cli.go")

	// If we're running from the source directory (not installed)
	if _, err := os.Stat(cliPath); os.IsNotExist(err) {
		// Try to find it relative to the current working directory
		cliPath = "scripts/lunarfetch-cli.go"
	}

	// Prepare the command to run the lunarfetch-cli.go script
	args := os.Args[1:] // Skip the program name
	cmd := exec.Command("go", append([]string{"run", cliPath}, args...)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	// Run the command
	err = cmd.Run()
	if err != nil {
		fmt.Printf("Error running CLI commands: %v\n", err)
		os.Exit(1)
	}

	// Exit after running the CLI command
	os.Exit(0)
}
