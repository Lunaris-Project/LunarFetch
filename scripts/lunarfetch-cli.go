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

	err = os.MkdirAll(logosDir, 0755)
	if err != nil {
		fmt.Printf("%sError creating config directory: %v%s\n", ColorRed, err, ColorReset)
		return
	}

	// Copy config file if it doesn't exist
	configFile := filepath.Join(configDir, "config.json")
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		copyFile("assets/config.json", configFile)
		fmt.Printf("%sCreated default configuration file%s\n", ColorGreen, ColorReset)
	}

	// Copy logo file if it doesn't exist
	logoFile := filepath.Join(logosDir, "moon.txt")
	if _, err := os.Stat(logoFile); os.IsNotExist(err) {
		copyFile("assets/logo.txt", logoFile)
		fmt.Printf("%sAdded sample logo%s\n", ColorGreen, ColorReset)
	}

	// Install the binary
	installBinary()
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
