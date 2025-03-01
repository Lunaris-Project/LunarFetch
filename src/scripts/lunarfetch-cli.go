// This file is a separate executable for installation purposes
//go:build ignore

package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// ANSI color codes
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
)

// Dependency represents a system dependency
type Dependency struct {
	Name        string
	Commands    []string
	ArchPackage string
	DebPackage  string
}

// Dependencies required by LunarFetch
var dependencies = []Dependency{
	{
		Name:        "coreutils",
		Commands:    []string{"df", "free"},
		ArchPackage: "coreutils",
		DebPackage:  "coreutils",
	},
	{
		Name:        "process utilities",
		Commands:    []string{"uptime"},
		ArchPackage: "procps-ng",
		DebPackage:  "procps",
	},
	{
		Name:        "LSB release",
		Commands:    []string{"lsb_release"},
		ArchPackage: "lsb-release",
		DebPackage:  "lsb-release",
	},
	{
		Name:        "display utilities",
		Commands:    []string{"xrandr", "xdpyinfo", "swaymsg", "wlr-randr"},
		ArchPackage: "xorg-xrandr xorg-xdpyinfo",
		DebPackage:  "x11-xserver-utils",
	},
	{
		Name:        "hardware detection",
		Commands:    []string{"lscpu", "lspci"},
		ArchPackage: "pciutils",
		DebPackage:  "pciutils",
	},
	{
		Name:        "theme detection",
		Commands:    []string{"gsettings"},
		ArchPackage: "gsettings-desktop-schemas",
		DebPackage:  "gnome-settings-daemon",
	},
	{
		Name:        "image processing",
		Commands:    []string{"convert"},
		ArchPackage: "imagemagick",
		DebPackage:  "imagemagick",
	},
}

// installMain is the entry point for the installation script
// Renamed from main to avoid conflict with main.go
func installMain() {
	args := os.Args[1:]
	if len(args) == 0 {
		printUsage()
		return
	}

	switch args[0] {
	case "install":
		install()
	case "uninstall":
		uninstall(false)
	case "purge":
		uninstall(true)
	case "check-deps":
		checkDependencies()
	case "install-deps":
		installDependencies()
	case "build":
		buildBinary()
	case "help":
		printUsage()
	case "setup-image":
		setupImageConfig()
	default:
		fmt.Printf("%sUnknown command: %s%s\n", ColorRed, args[0], ColorReset)
		printUsage()
	}
}

// main function calls installMain to maintain compatibility
func main() {
	installMain()
}

func printUsage() {
	fmt.Printf("%sLunarFetch Installer%s\n\n", ColorCyan, ColorReset)
	fmt.Printf("Usage: %sgo run install.go [command]%s\n\n", ColorGreen, ColorReset)
	fmt.Printf("Commands:\n")
	fmt.Printf("  %sinstall%s       Install LunarFetch to your system\n", ColorGreen, ColorReset)
	fmt.Printf("  %suninstall%s     Remove LunarFetch binary\n", ColorGreen, ColorReset)
	fmt.Printf("  %spurge%s         Remove LunarFetch binary and all configuration files\n", ColorGreen, ColorReset)
	fmt.Printf("  %scheck-deps%s    Check for required dependencies\n", ColorGreen, ColorReset)
	fmt.Printf("  %sinstall-deps%s  Install required dependencies\n", ColorGreen, ColorReset)
	fmt.Printf("  %sbuild%s         Build the binary without installing\n", ColorGreen, ColorReset)
	fmt.Printf("  %ssetup-image%s   Configure image display settings\n", ColorGreen, ColorReset)
	fmt.Printf("  %shelp%s          Show this help message\n", ColorGreen, ColorReset)
}

func install() {
	fmt.Printf("%sInstalling LunarFetch...%s\n", ColorCyan, ColorReset)

	// Check dependencies
	missingDeps := checkDependencies()
	if len(missingDeps) > 0 {
		fmt.Printf("\n%sWould you like to install missing dependencies? (y/n):%s ", ColorYellow, ColorReset)
		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(response)
		if strings.ToLower(response) == "y" {
			installDependencies()
		}
	}

	// Build the binary
	buildBinary()

	// Create config directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("%sError getting home directory: %v%s\n", ColorRed, err, ColorReset)
		return
	}

	configDir := filepath.Join(homeDir, ".config", "lunarfetch")
	logosDir := filepath.Join(configDir, "logos")
	imagesDir := filepath.Join(configDir, "images")

	err = os.MkdirAll(logosDir, 0755)
	if err != nil {
		fmt.Printf("%sError creating logos directory: %v%s\n", ColorRed, err, ColorReset)
		return
	}

	err = os.MkdirAll(imagesDir, 0755)
	if err != nil {
		fmt.Printf("%sError creating images directory: %v%s\n", ColorRed, err, ColorReset)
		return
	}

	// Copy config file if it doesn't exist
	configFile := filepath.Join(configDir, "config.json")
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		copyFile("src/assets/config.json", configFile)
		fmt.Printf("%sCreated default configuration file%s\n", ColorGreen, ColorReset)
	}

	// Copy logo file if it doesn't exist
	logoFile := filepath.Join(logosDir, "moon.txt")
	if _, err := os.Stat(logoFile); os.IsNotExist(err) {
		copyFile("src/assets/logo.txt", logoFile)
		fmt.Printf("%sAdded sample logo%s\n", ColorGreen, ColorReset)
	}

	// Copy sample image if it exists
	sampleImage := "src/assets/sample.png"
	if _, err := os.Stat(sampleImage); !os.IsNotExist(err) {
		destImage := filepath.Join(imagesDir, "sample.png")
		copyFile(sampleImage, destImage)
		fmt.Printf("%sAdded sample image%s\n", ColorGreen, ColorReset)
	}

	// Install the binary
	installBinary()
}

// setupImageConfig configures image display settings
func setupImageConfig() {
	fmt.Printf("%sConfiguring image display settings...%s\n", ColorCyan, ColorReset)

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("%sError getting home directory: %v%s\n", ColorRed, err, ColorReset)
		return
	}

	configDir := filepath.Join(homeDir, ".config", "lunarfetch")
	imagesDir := filepath.Join(configDir, "images")

	// Create images directory if it doesn't exist
	err = os.MkdirAll(imagesDir, 0755)
	if err != nil {
		fmt.Printf("%sError creating images directory: %v%s\n", ColorRed, err, ColorReset)
		return
	}

	reader := bufio.NewReader(os.Stdin)

	// Ask for image path
	fmt.Printf("\n%sEnter path to image or directory containing images:%s\n", ColorYellow, ColorReset)
	fmt.Printf("(Default: %s): ", imagesDir)
	imagePath, _ := reader.ReadString('\n')
	imagePath = strings.TrimSpace(imagePath)
	if imagePath == "" {
		imagePath = imagesDir
	}

	// Ask for image width
	fmt.Printf("\n%sEnter image width (in characters):%s\n", ColorYellow, ColorReset)
	fmt.Printf("(Default: 80): ")
	widthStr, _ := reader.ReadString('\n')
	widthStr = strings.TrimSpace(widthStr)
	width := "80"
	if widthStr != "" {
		width = widthStr
	}

	// Ask for image height
	fmt.Printf("\n%sEnter image height (in characters):%s\n", ColorYellow, ColorReset)
	fmt.Printf("(Default: 24): ")
	heightStr, _ := reader.ReadString('\n')
	heightStr = strings.TrimSpace(heightStr)
	height := "24"
	if heightStr != "" {
		height = heightStr
	}

	// Ask for render mode
	fmt.Printf("\n%sSelect render mode:%s\n", ColorYellow, ColorReset)
	fmt.Printf("1. Simple (fewer ASCII characters)\n")
	fmt.Printf("2. Detailed (more ASCII characters)\n")
	fmt.Printf("3. Default (medium detail)\n")
	fmt.Printf("(Default: 3): ")
	renderModeStr, _ := reader.ReadString('\n')
	renderModeStr = strings.TrimSpace(renderModeStr)
	renderMode := "default"
	switch renderModeStr {
	case "1":
		renderMode = "simple"
	case "2":
		renderMode = "detailed"
	}

	// Ask for dither mode
	fmt.Printf("\n%sSelect dither mode:%s\n", ColorYellow, ColorReset)
	fmt.Printf("1. None\n")
	fmt.Printf("2. Floyd-Steinberg\n")
	fmt.Printf("(Default: 1): ")
	ditherModeStr, _ := reader.ReadString('\n')
	ditherModeStr = strings.TrimSpace(ditherModeStr)
	ditherMode := "none"
	if ditherModeStr == "2" {
		ditherMode = "floyd-steinberg"
	}

	// Ask for terminal output
	fmt.Printf("\n%sUse ANSI color output? (y/n):%s\n", ColorYellow, ColorReset)
	fmt.Printf("(Default: n): ")
	terminalOutputStr, _ := reader.ReadString('\n')
	terminalOutputStr = strings.TrimSpace(terminalOutputStr)
	terminalOutput := "false"
	if strings.ToLower(terminalOutputStr) == "y" {
		terminalOutput = "true"
	}

	// Update config file
	configFile := filepath.Join(configDir, "config.json")
	updateConfigWithImageSettings(configFile, imagePath, width, height, renderMode, ditherMode, terminalOutput)

	fmt.Printf("\n%sImage display settings configured successfully!%s\n", ColorGreen, ColorReset)
	fmt.Printf("Run '%slunarfetch%s' to see the results.\n", ColorGreen, ColorReset)
}

// updateConfigWithImageSettings updates the config file with image settings
func updateConfigWithImageSettings(configFile, imagePath, width, height, renderMode, ditherMode, terminalOutput string) {
	// Read the current config file
	content, err := os.ReadFile(configFile)
	if err != nil {
		fmt.Printf("%sError reading config file: %v%s\n", ColorRed, err, ColorReset)
		return
	}

	// Convert the config file to a string
	configStr := string(content)

	// Check if the image section already exists
	if !strings.Contains(configStr, `"image":`) {
		// Add the image section before the modules section
		imageSection := fmt.Sprintf(`  "image": {
    "enableImage": true,
    "imagePath": "%s",
    "width": %s,
    "height": %s,
    "renderMode": "%s",
    "ditherMode": "%s",
    "terminalOutput": %s
  },`, imagePath, width, height, renderMode, ditherMode, terminalOutput)

		// Find the position to insert the image section
		modulesPos := strings.Index(configStr, `  "modules":`)
		if modulesPos != -1 {
			configStr = configStr[:modulesPos] + imageSection + "\n" + configStr[modulesPos:]
		}
	} else {
		// Update the existing image section
		configStr = updateJsonValue(configStr, "image.enableImage", "true")
		configStr = updateJsonValue(configStr, "image.imagePath", fmt.Sprintf(`"%s"`, imagePath))
		configStr = updateJsonValue(configStr, "image.width", width)
		configStr = updateJsonValue(configStr, "image.height", height)
		configStr = updateJsonValue(configStr, "image.renderMode", fmt.Sprintf(`"%s"`, renderMode))
		configStr = updateJsonValue(configStr, "image.ditherMode", fmt.Sprintf(`"%s"`, ditherMode))
		configStr = updateJsonValue(configStr, "image.terminalOutput", terminalOutput)
	}

	// Disable logo if image is enabled
	configStr = updateJsonValue(configStr, "logo.enableLogo", "false")

	// Write the updated config back to the file
	err = os.WriteFile(configFile, []byte(configStr), 0644)
	if err != nil {
		fmt.Printf("%sError writing config file: %v%s\n", ColorRed, err, ColorReset)
		return
	}
}

// updateJsonValue updates a JSON value in a string
func updateJsonValue(jsonStr, path, newValue string) string {
	parts := strings.Split(path, ".")
	if len(parts) != 2 {
		return jsonStr
	}

	section := parts[0]
	key := parts[1]

	// Find the section
	sectionStart := strings.Index(jsonStr, fmt.Sprintf(`"%s": {`, section))
	if sectionStart == -1 {
		return jsonStr
	}

	// Find the key within the section
	sectionEnd := findClosingBrace(jsonStr, sectionStart)
	if sectionEnd == -1 {
		return jsonStr
	}

	sectionContent := jsonStr[sectionStart:sectionEnd]
	keyStart := strings.Index(sectionContent, fmt.Sprintf(`"%s":`, key))
	if keyStart == -1 {
		return jsonStr
	}

	// Find the end of the value
	valueStart := sectionStart + keyStart + len(fmt.Sprintf(`"%s":`, key))
	valueEnd := findValueEnd(jsonStr, valueStart)
	if valueEnd == -1 {
		return jsonStr
	}

	// Replace the value
	return jsonStr[:valueStart] + " " + newValue + jsonStr[valueEnd:]
}

// findClosingBrace finds the closing brace for a JSON object
func findClosingBrace(jsonStr string, start int) int {
	braceCount := 0
	inQuotes := false
	escaped := false

	for i := start; i < len(jsonStr); i++ {
		c := jsonStr[i]

		if escaped {
			escaped = false
			continue
		}

		if c == '\\' {
			escaped = true
			continue
		}

		if c == '"' {
			inQuotes = !inQuotes
			continue
		}

		if !inQuotes {
			if c == '{' {
				braceCount++
			} else if c == '}' {
				braceCount--
				if braceCount == 0 {
					return i + 1
				}
			}
		}
	}

	return -1
}

// findValueEnd finds the end of a JSON value
func findValueEnd(jsonStr string, start int) int {
	inQuotes := false
	escaped := false

	for i := start; i < len(jsonStr); i++ {
		c := jsonStr[i]

		if escaped {
			escaped = false
			continue
		}

		if c == '\\' {
			escaped = true
			continue
		}

		if c == '"' {
			inQuotes = !inQuotes
			continue
		}

		if !inQuotes && (c == ',' || c == '}') {
			return i
		}
	}

	return -1
}

func uninstall(purge bool) {
	fmt.Printf("%sUninstalling LunarFetch...%s\n", ColorCyan, ColorReset)

	// Remove binary
	err := exec.Command("sudo", "rm", "-f", "/usr/local/bin/lunarfetch").Run()
	if err != nil {
		fmt.Printf("%sError removing binary: %v%s\n", ColorRed, err, ColorReset)
	} else {
		fmt.Printf("%sRemoved LunarFetch binary%s\n", ColorGreen, ColorReset)
	}

	// Remove config files if purge is true
	if purge {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("%sError getting home directory: %v%s\n", ColorRed, err, ColorReset)
			return
		}

		configDir := filepath.Join(homeDir, ".config", "lunarfetch")
		err = os.RemoveAll(configDir)
		if err != nil {
			fmt.Printf("%sError removing config directory: %v%s\n", ColorRed, err, ColorReset)
		} else {
			fmt.Printf("%sRemoved LunarFetch configuration files%s\n", ColorGreen, ColorReset)
		}
	}
}

func checkDependencies() []Dependency {
	fmt.Printf("%sChecking for required dependencies...%s\n", ColorCyan, ColorReset)

	var missingDeps []Dependency

	for _, dep := range dependencies {
		found := false
		for _, cmd := range dep.Commands {
			if commandExists(cmd) {
				found = true
				break
			}
		}

		if !found {
			missingDeps = append(missingDeps, dep)
			fmt.Printf("  %s✗ %s%s\n", ColorRed, dep.Name, ColorReset)
		} else {
			fmt.Printf("  %s✓ %s%s\n", ColorGreen, dep.Name, ColorReset)
		}
	}

	if len(missingDeps) > 0 {
		fmt.Printf("\n%sWarning: Some dependencies are missing.%s\n", ColorYellow, ColorReset)
		fmt.Printf("Run '%sgo run install.go install-deps%s' to install them.\n", ColorGreen, ColorReset)
	} else {
		fmt.Printf("\n%sAll dependencies are installed.%s\n", ColorGreen, ColorReset)
	}

	return missingDeps
}

func installDependencies() {
	fmt.Printf("%sInstalling dependencies...%s\n", ColorCyan, ColorReset)

	// Detect distribution
	var distro string
	if commandExists("pacman") {
		distro = "arch"
	} else if commandExists("apt") {
		distro = "debian"
	} else {
		fmt.Printf("%sUnsupported distribution. Please install dependencies manually.%s\n", ColorRed, ColorReset)
		return
	}

	// Install dependencies based on distribution
	if distro == "arch" {
		installArchDependencies()
	} else if distro == "debian" {
		installDebianDependencies()
	}
}

func installArchDependencies() {
	fmt.Printf("%sInstalling dependencies for Arch Linux...%s\n", ColorCyan, ColorReset)

	var packages []string
	for _, dep := range dependencies {
		if !dependencyExists(dep) && dep.ArchPackage != "" {
			packages = append(packages, strings.Split(dep.ArchPackage, " ")...)
		}
	}

	if len(packages) == 0 {
		fmt.Printf("%sNo packages to install.%s\n", ColorGreen, ColorReset)
		return
	}

	args := append([]string{"pacman", "-S", "--needed"}, packages...)
	cmd := exec.Command("sudo", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	if err != nil {
		fmt.Printf("%sError installing packages: %v%s\n", ColorRed, err, ColorReset)
	} else {
		fmt.Printf("%sSuccessfully installed dependencies.%s\n", ColorGreen, ColorReset)
	}
}

func installDebianDependencies() {
	fmt.Printf("%sInstalling dependencies for Debian/Ubuntu...%s\n", ColorCyan, ColorReset)

	var packages []string
	for _, dep := range dependencies {
		if !dependencyExists(dep) && dep.DebPackage != "" {
			packages = append(packages, strings.Split(dep.DebPackage, " ")...)
		}
	}

	if len(packages) == 0 {
		fmt.Printf("%sNo packages to install.%s\n", ColorGreen, ColorReset)
		return
	}

	// Update package lists
	updateCmd := exec.Command("sudo", "apt", "update")
	updateCmd.Stdout = os.Stdout
	updateCmd.Stderr = os.Stderr
	updateCmd.Run()

	// Install packages
	args := append([]string{"apt", "install", "-y"}, packages...)
	cmd := exec.Command("sudo", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	if err != nil {
		fmt.Printf("%sError installing packages: %v%s\n", ColorRed, err, ColorReset)
	} else {
		fmt.Printf("%sSuccessfully installed dependencies.%s\n", ColorGreen, ColorReset)
	}
}

func buildBinary() {
	fmt.Printf("%sBuilding LunarFetch...%s\n", ColorCyan, ColorReset)

	cmd := exec.Command("go", "build", "-o", "lunarfetch", "main.go")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("%sError building binary: %v%s\n", ColorRed, err, ColorReset)
		fmt.Println(string(output))
		return
	}

	fmt.Printf("%sSuccessfully built LunarFetch binary.%s\n", ColorGreen, ColorReset)
}

func installBinary() {
	fmt.Printf("%sInstalling LunarFetch binary...%s\n", ColorCyan, ColorReset)

	cmd := exec.Command("sudo", "mv", "lunarfetch", "/usr/local/bin/")
	err := cmd.Run()
	if err != nil {
		fmt.Printf("%sError installing binary: %v%s\n", ColorRed, err, ColorReset)
		return
	}

	fmt.Printf("%sLunarFetch has been installed to /usr/local/bin/%s\n", ColorGreen, ColorReset)
	fmt.Printf("Run '%slunarfetch%s' to try it out!\n", ColorGreen, ColorReset)
}

// Helper functions

func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func dependencyExists(dep Dependency) bool {
	for _, cmd := range dep.Commands {
		if commandExists(cmd) {
			return true
		}
	}
	return false
}

func copyFile(src, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	err = os.WriteFile(dst, input, 0644)
	if err != nil {
		return err
	}

	return nil
}
