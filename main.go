package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"lunarfetch/src/utils"
)

// Version information
const (
	Version     = "1.0.0"
	VersionDate = "2025-03-01"
)

// Color codes for terminal output
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorCyan   = "\033[36m"
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
		Commands:    []string{"lspci"},
		ArchPackage: "pciutils",
		DebPackage:  "pciutils",
	},
	{
		Name:        "theme detection",
		Commands:    []string{"gsettings", "dconf"},
		ArchPackage: "gsettings-desktop-schemas dconf",
		DebPackage:  "gsettings-desktop-schemas dconf-cli",
	},
	{
		Name:        "image processing",
		Commands:    []string{"chafa", "sixel-convert"},
		ArchPackage: "chafa",
		DebPackage:  "chafa",
	},
}

// SimpleConfig is a simplified configuration structure for installation
type SimpleConfig struct {
	Logo Logo `json:"logo"`
	Info Info `json:"info"`
}

// Logo represents the logo configuration for installation
type Logo struct {
	Enabled bool   `json:"enabled"`
	Type    string `json:"type"`
	Path    string `json:"path"`
	Color   string `json:"color"`
}

// Info represents the information display configuration for installation
type Info struct {
	Enabled bool     `json:"enabled"`
	Items   []string `json:"items"`
}

func main() {
	var configPath string

	// Parse command line arguments
	if len(os.Args) > 1 {
		// Check for help flag first
		if os.Args[1] == "--help" || os.Args[1] == "-h" {
			printUsage()
			return
		}

		// Check for version flag
		if os.Args[1] == "--version" || os.Args[1] == "-v" {
			printVersion()
			return
		}

		// Check for debug flag
		if os.Args[1] == "--debug" || os.Args[1] == "-d" {
			// Enable debug mode and shift arguments
			os.Setenv("LUNARFETCH_DEBUG", "1")
			if len(os.Args) > 2 {
				// Shift arguments
				os.Args = append([]string{os.Args[0]}, os.Args[2:]...)
			} else {
				// No more arguments, just run normally
				os.Args = []string{os.Args[0]}
			}
		}

		// Check for config flag
		for i := 1; i < len(os.Args); i++ {
			if os.Args[i] == "--config" || os.Args[i] == "-c" {
				if i+1 < len(os.Args) {
					configPath = os.Args[i+1]
					// Remove these arguments
					newArgs := append([]string{os.Args[0]}, os.Args[1:i]...)
					if i+2 < len(os.Args) {
						newArgs = append(newArgs, os.Args[i+2:]...)
					}
					os.Args = newArgs
					break
				}
			}
		}

		if len(os.Args) > 1 {
			handleCommands(os.Args[1:])
			return
		}
	}

	// Load configuration
	configLoader := utils.NewConfigLoader()
	var config utils.Config
	var err error

	if configPath != "" {
		config, err = configLoader.LoadConfig(configPath)
		if err != nil {
			fmt.Println("Error loading config from", configPath, ":", err)
			config = utils.DefaultConfig()
		}
	} else {
		config, err = configLoader.LoadConfig()
		if err != nil {
			fmt.Println("Error loading config:", err)
			config = utils.DefaultConfig()
		}
	}

	// Initialize display manager
	displayManager := utils.NewDisplayManager(config)
	displayManager.InitializeComponents()

	// Variables to store image and logo outputs
	var imageOutput string
	var logoOutput string
	var sysInfoOutput string

	// Get system information
	sysInfoOutput = displayManager.Display()

	// Load logo if enabled
	if config.Logo.EnableLogo {
		logoLoader := utils.NewLogoLoader(config.Logo.LogoPath)
		var err error
		logoOutput, err = logoLoader.GetRandomLogo()
		if err != nil && os.Getenv("LUNARFETCH_DEBUG") == "1" {
			fmt.Printf("Error loading logo: %v\n", err)
		}
	}

	// Load image if enabled
	if config.Image.Enabled || config.Image.EnableImage {
		imageLoader := utils.NewImageLoader(config)
		var err error

		// Debug output
		if os.Getenv("LUNARFETCH_DEBUG") == "1" {
			fmt.Printf("Image settings:\n")
			fmt.Printf("  Enabled: %v\n", config.Image.Enabled)
			fmt.Printf("  Path: %s\n", config.Image.ImagePath)
			fmt.Printf("  Width: %d\n", config.Image.Width)
			fmt.Printf("  Height: %d\n", config.Image.Height)
			fmt.Printf("  Position: %s\n", config.Image.Position)
		}

		imageOutput, err = imageLoader.RenderImage()
		if err != nil && os.Getenv("LUNARFETCH_DEBUG") == "1" {
			fmt.Printf("Error rendering image: %v\n", err)
		}
	}

	// Display the information based on configuration
	if config.Image.Position == "above" && (config.Image.Enabled || config.Image.EnableImage) && imageOutput != "" {
		// Display image above system info
		fmt.Print(imageOutput)
		fmt.Print(sysInfoOutput)
	} else if config.Logo.Position == "above" && config.Logo.EnableLogo && logoOutput != "" {
		// Display logo above system info
		fmt.Print(logoOutput)
		fmt.Println()
		fmt.Print(sysInfoOutput)
	} else {
		// Display side by side (default behavior)
		// Determine display order based on configuration
		if config.Display.ShowImageFirst && (config.Image.Enabled || config.Image.EnableImage) && imageOutput != "" {
			fmt.Print(imageOutput)
		} else if config.Display.ShowLogoFirst && config.Logo.EnableLogo && logoOutput != "" {
			fmt.Print(logoOutput)
			fmt.Println()
		}

		// Display system information
		fmt.Print(sysInfoOutput)

		// Display the other element if not shown first
		if !config.Display.ShowImageFirst && (config.Image.Enabled || config.Image.EnableImage) && imageOutput != "" {
			fmt.Print(imageOutput)
		} else if !config.Display.ShowLogoFirst && config.Logo.EnableLogo && logoOutput != "" {
			fmt.Print(logoOutput)
			fmt.Println()
		}
	}
}

func handleCommands(args []string) {
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
	case "help", "-h", "--help":
		printUsage()
	case "version", "-v", "--version":
		printVersion()
	case "setup-image":
		setupImage()
	default:
		fmt.Printf("%sUnknown command: %s%s\n", ColorRed, args[0], ColorReset)
		printUsage()
	}
}

// printUsage displays usage information
func printUsage() {
	fmt.Printf("%sLunarFetch%s - A customizable system information display tool\n\n", ColorCyan, ColorReset)

	// Main usage
	fmt.Printf("%sUSAGE:%s\n", ColorYellow, ColorReset)
	fmt.Printf("  lunarfetch [flags] [command]\n\n")

	// Flags section
	fmt.Printf("%sFLAGS:%s\n", ColorYellow, ColorReset)
	fmt.Printf("  %s-c, --config%s <path>    Specify a custom configuration file path\n", ColorGreen, ColorReset)
	fmt.Printf("  %s-d, --debug%s            Enable debug mode for verbose output\n", ColorGreen, ColorReset)
	fmt.Printf("  %s-h, --help%s             Display this help message\n", ColorGreen, ColorReset)
	fmt.Printf("  %s-v, --version%s          Display version information\n\n", ColorGreen, ColorReset)

	// Commands section
	fmt.Printf("%sCOMMANDS:%s\n", ColorYellow, ColorReset)
	fmt.Printf("  %sinstall%s              Install LunarFetch to your system\n", ColorGreen, ColorReset)
	fmt.Printf("                       - Creates configuration directories\n")
	fmt.Printf("                       - Copies default configuration and assets\n")
	fmt.Printf("                       - Installs binary to /usr/local/bin\n\n")

	fmt.Printf("  %suninstall%s            Remove LunarFetch binary from your system\n", ColorGreen, ColorReset)
	fmt.Printf("                       - Removes only the binary file\n")
	fmt.Printf("                       - Keeps configuration files intact\n\n")

	fmt.Printf("  %spurge%s                Completely remove LunarFetch from your system\n", ColorGreen, ColorReset)
	fmt.Printf("                       - Removes the binary file\n")
	fmt.Printf("                       - Deletes all configuration files and directories\n\n")

	fmt.Printf("  %scheck-deps%s           Check for required system dependencies\n", ColorGreen, ColorReset)
	fmt.Printf("                       - Verifies all required commands are available\n")
	fmt.Printf("                       - Reports missing dependencies\n\n")

	fmt.Printf("  %sinstall-deps%s         Install required system dependencies\n", ColorGreen, ColorReset)
	fmt.Printf("                       - Automatically installs missing dependencies\n")
	fmt.Printf("                       - Supports Arch Linux and Debian/Ubuntu\n\n")

	fmt.Printf("  %sbuild%s                Build the binary without installing\n", ColorGreen, ColorReset)
	fmt.Printf("                       - Creates executable in current directory\n")
	fmt.Printf("                       - Useful for development and testing\n\n")

	fmt.Printf("  %ssetup-image%s          Configure image display support\n", ColorGreen, ColorReset)
	fmt.Printf("                       - Sets up image configuration\n")
	fmt.Printf("                       - Creates necessary directories\n")
	fmt.Printf("                       - Updates configuration file\n\n")

	fmt.Printf("  %shelp%s                 Display this help message\n\n", ColorGreen, ColorReset)

	fmt.Printf("  %sversion%s              Display version information\n\n", ColorGreen, ColorReset)

	// Examples section
	fmt.Printf("%sEXAMPLES:%s\n", ColorYellow, ColorReset)
	fmt.Printf("  lunarfetch                          # Display system information with default config\n")
	fmt.Printf("  lunarfetch -c ~/.config/lunarfetch/custom.json  # Use custom config file\n")
	fmt.Printf("  lunarfetch --debug                  # Run with debug output\n")
	fmt.Printf("  lunarfetch install                  # Install LunarFetch to your system\n")
	fmt.Printf("  lunarfetch setup-image              # Configure image display\n\n")

	// Configuration section
	fmt.Printf("%sCONFIGURATION:%s\n", ColorYellow, ColorReset)
	fmt.Printf("  Default config location: ~/.config/lunarfetch/config.json\n")
	fmt.Printf("  Logo directory: ~/.config/lunarfetch/logos/\n")
	fmt.Printf("  Images directory: ~/.config/lunarfetch/images/\n\n")

	fmt.Printf("For more information and updates, visit: %shttps://github.com/Lunaris-Project/lunarfetch%s\n", ColorBlue, ColorReset)
}

// printVersion displays version information
func printVersion() {
	fmt.Printf("%sLunarFetch%s - A customizable system information display tool\n", ColorCyan, ColorReset)
	fmt.Printf("Version: %s\n", Version)
	fmt.Printf("Version Date: %s\n", VersionDate)
	fmt.Printf("For more information, visit: %shttps://github.com/Lunaris-Project/lunarfetch%s\n", ColorBlue, ColorReset)
}

// install installs LunarFetch to the system
func install() {
	// Find the source directory
	sourceDir, err := findSourceDirectory()
	if err != nil {
		fmt.Printf("%sError: %s%s\n", ColorRed, err.Error(), ColorReset)
		fmt.Printf("%sPlease run this command from the LunarFetch source directory or ensure Go modules are installed.%s\n", ColorYellow, ColorReset)
		return
	}

	fmt.Printf("%sFound LunarFetch source directory: %s%s\n", ColorGreen, sourceDir, ColorReset)

	// Save the current directory to return to it later
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("%sError: Could not get current directory: %s%s\n", ColorRed, err.Error(), ColorReset)
		return
	}

	// Change to the source directory for building
	err = os.Chdir(sourceDir)
	if err != nil {
		fmt.Printf("%sError: Could not change to source directory: %s%s\n", ColorRed, err.Error(), ColorReset)
		return
	}

	// Ensure we change back to the original directory when done
	defer func() {
		os.Chdir(currentDir)
	}()

	// Check for missing dependencies
	missingDeps := checkDependencies()
	if len(missingDeps) > 0 {
		fmt.Printf("%sMissing dependencies: %v%s\n", ColorYellow, missingDeps, ColorReset)
		fmt.Printf("%sWould you like to install them? (y/n): %s", ColorYellow, ColorReset)
		var answer string
		fmt.Scanln(&answer)
		if strings.ToLower(answer) == "y" {
			installDependencies()
		} else {
			fmt.Printf("%sSkipping dependency installation. Some features may not work correctly.%s\n", ColorYellow, ColorReset)
		}
	}

	// Create config directory if it doesn't exist
	configDir := filepath.Join(os.Getenv("HOME"), ".config", "lunarfetch")
	err = os.MkdirAll(configDir, 0755)
	if err != nil {
		fmt.Printf("%sError: Could not create config directory: %s%s\n", ColorRed, err.Error(), ColorReset)
		return
	}

	// Create images directory if it doesn't exist
	imagesDir := filepath.Join(configDir, "images")
	err = os.MkdirAll(imagesDir, 0755)
	if err != nil {
		fmt.Printf("%sError: Could not create images directory: %s%s\n", ColorRed, err.Error(), ColorReset)
		return
	}

	// Copy sample config if it doesn't exist
	configPath := filepath.Join(configDir, "config.json")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		sampleConfigPath := filepath.Join(sourceDir, "config.json")
		if _, err := os.Stat(sampleConfigPath); err == nil {
			err = copyFile(sampleConfigPath, configPath)
			if err != nil {
				fmt.Printf("%sError: Could not copy sample config: %s%s\n", ColorRed, err.Error(), ColorReset)
			} else {
				fmt.Printf("%sCreated sample config at %s%s\n", ColorGreen, configPath, ColorReset)
			}
		} else {
			// Create a default config if sample doesn't exist
			defaultConfig := SimpleConfig{
				Logo: Logo{
					Enabled: true,
					Type:    "ascii",
					Path:    "",
					Color:   "blue",
				},
				Info: Info{
					Enabled: true,
					Items: []string{
						"host", "user", "os", "kernel", "uptime", "terminal", "shell",
						"disk", "memory", "packages", "battery", "gpu", "cpu", "resolution",
						"wm", "theme", "icons", "desktop",
					},
				},
			}
			configJSON, _ := json.MarshalIndent(defaultConfig, "", "  ")
			err = os.WriteFile(configPath, configJSON, 0644)
			if err != nil {
				fmt.Printf("%sError: Could not create default config: %s%s\n", ColorRed, err.Error(), ColorReset)
			} else {
				fmt.Printf("%sCreated default config at %s%s\n", ColorGreen, configPath, ColorReset)
			}
		}
	}

	// Copy sample side config if it doesn't exist
	sideConfigPath := filepath.Join(configDir, "side.json")
	if _, err := os.Stat(sideConfigPath); os.IsNotExist(err) {
		sampleSideConfigPath := filepath.Join(sourceDir, "side.json")
		if _, err := os.Stat(sampleSideConfigPath); err == nil {
			err = copyFile(sampleSideConfigPath, sideConfigPath)
			if err != nil {
				fmt.Printf("%sError: Could not copy sample side config: %s%s\n", ColorRed, err.Error(), ColorReset)
			} else {
				fmt.Printf("%sCreated sample side config at %s%s\n", ColorGreen, sideConfigPath, ColorReset)
			}
		}
	}

	// Copy sample test-side config if it doesn't exist
	testSideConfigPath := filepath.Join(configDir, "test-side.json")
	if _, err := os.Stat(testSideConfigPath); os.IsNotExist(err) {
		sampleTestSideConfigPath := filepath.Join(sourceDir, "test-side.json")
		if _, err := os.Stat(sampleTestSideConfigPath); err == nil {
			err = copyFile(sampleTestSideConfigPath, testSideConfigPath)
			if err != nil {
				fmt.Printf("%sError: Could not copy sample test-side config: %s%s\n", ColorRed, err.Error(), ColorReset)
			} else {
				fmt.Printf("%sCreated sample test-side config at %s%s\n", ColorGreen, testSideConfigPath, ColorReset)
			}
		}
	}

	// Copy sample image if it doesn't exist
	sampleImagePath := filepath.Join(sourceDir, "sample.png")
	if _, err := os.Stat(sampleImagePath); err == nil {
		destImagePath := filepath.Join(imagesDir, "sample.png")
		if _, err := os.Stat(destImagePath); os.IsNotExist(err) {
			err = copyFile(sampleImagePath, destImagePath)
			if err != nil {
				fmt.Printf("%sError: Could not copy sample image: %s%s\n", ColorRed, err.Error(), ColorReset)
			} else {
				fmt.Printf("%sCreated sample image at %s%s\n", ColorGreen, destImagePath, ColorReset)
			}
		}
	}

	// Build and install the binary
	fmt.Printf("%sBuilding LunarFetch...%s\n", ColorYellow, ColorReset)
	if buildBinary() {
		fmt.Printf("%sInstalling LunarFetch...%s\n", ColorYellow, ColorReset)
		installBinary()
	}

	// Clean up temporary directory if we created one
	if strings.Contains(sourceDir, "lunarfetch-install") {
		fmt.Printf("%sCleaning up temporary directory...%s\n", ColorYellow, ColorReset)
		os.RemoveAll(sourceDir)
	}
}

// uninstall removes LunarFetch from the system
func uninstall(purge bool) {
	fmt.Printf("%sUninstalling LunarFetch...%s\n", ColorCyan, ColorReset)

	// Remove binary
	cmd := exec.Command("sudo", "rm", "-f", "/usr/local/bin/lunarfetch")
	err := cmd.Run()
	if err != nil {
		fmt.Printf("%sError removing binary: %v%s\n", ColorRed, err, ColorReset)
	} else {
		fmt.Printf("%sRemoved LunarFetch binary%s\n", ColorGreen, ColorReset)
	}

	// Remove configuration if purge is true
	if purge {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("%sError getting home directory: %v%s\n", ColorRed, err, ColorReset)
			return
		}

		configDir := filepath.Join(homeDir, ".config", "lunarfetch")
		err = os.RemoveAll(configDir)
		if err != nil {
			fmt.Printf("%sError removing configuration directory: %v%s\n", ColorRed, err, ColorReset)
		} else {
			fmt.Printf("%sRemoved configuration directory%s\n", ColorGreen, ColorReset)
		}
	}

	fmt.Printf("%sLunarFetch has been uninstalled%s\n", ColorGreen, ColorReset)
}

// checkDependencies checks for required dependencies
func checkDependencies() []Dependency {
	fmt.Printf("%sChecking for required dependencies...%s\n", ColorCyan, ColorReset)

	var missingDeps []Dependency
	for _, dep := range dependencies {
		if !dependencyExists(dep) {
			missingDeps = append(missingDeps, dep)
			fmt.Printf("  %s✗ %s%s\n", ColorRed, dep.Name, ColorReset)
		} else {
			fmt.Printf("  %s✓ %s%s\n", ColorGreen, dep.Name, ColorReset)
		}
	}

	if len(missingDeps) > 0 {
		fmt.Printf("\n%sSome dependencies are missing.%s\n", ColorYellow, ColorReset)
		fmt.Printf("Run '%slunarfetch install-deps%s' to install them.\n", ColorGreen, ColorReset)
	} else {
		fmt.Printf("\n%sAll dependencies are installed.%s\n", ColorGreen, ColorReset)
	}

	return missingDeps
}

// installDependencies installs required dependencies
func installDependencies() {
	fmt.Printf("%sInstalling dependencies...%s\n", ColorCyan, ColorReset)

	// Try to detect distribution
	var distro string

	// First try lsb_release if available
	out, err := exec.Command("lsb_release", "-si").Output()
	if err == nil {
		distro = strings.TrimSpace(string(out))
	} else {
		// If lsb_release fails, try checking for distribution-specific files
		if _, err := os.Stat("/etc/arch-release"); err == nil {
			distro = "arch"
		} else if _, err := os.Stat("/etc/debian_version"); err == nil {
			distro = "debian"
		} else if _, err := os.Stat("/etc/os-release"); err == nil {
			// Read /etc/os-release for more information
			osReleaseContent, err := os.ReadFile("/etc/os-release")
			if err == nil {
				osRelease := string(osReleaseContent)
				if strings.Contains(strings.ToLower(osRelease), "arch") {
					distro = "arch"
				} else if strings.Contains(strings.ToLower(osRelease), "debian") {
					distro = "debian"
				} else if strings.Contains(strings.ToLower(osRelease), "ubuntu") {
					distro = "ubuntu"
				}
			}
		}
	}

	// If we couldn't detect the distribution, prompt the user
	if distro == "" {
		fmt.Printf("%sCould not automatically detect your distribution.%s\n", ColorYellow, ColorReset)
		fmt.Printf("Please select your distribution:\n")
		fmt.Printf("1. Arch Linux\n")
		fmt.Printf("2. Debian/Ubuntu\n")
		fmt.Printf("3. Other\n")
		fmt.Printf("Enter your choice (1-3): ")

		var choice string
		fmt.Scanln(&choice)

		switch choice {
		case "1":
			distro = "arch"
		case "2":
			distro = "debian"
		default:
			fmt.Printf("%sUnsupported distribution. Please install dependencies manually.%s\n", ColorRed, ColorReset)
			return
		}
	}

	// Install dependencies based on distribution
	if strings.Contains(strings.ToLower(distro), "arch") {
		installArchDependencies()
	} else if strings.Contains(strings.ToLower(distro), "debian") || strings.Contains(strings.ToLower(distro), "ubuntu") {
		installDebianDependencies()
	} else {
		fmt.Printf("%sUnsupported distribution: %s. Please install dependencies manually.%s\n", ColorRed, distro, ColorReset)
	}
}

// installArchDependencies installs dependencies for Arch Linux
func installArchDependencies() {
	fmt.Printf("%sInstalling dependencies for Arch Linux...%s\n", ColorCyan, ColorReset)

	var packages []string
	for _, dep := range dependencies {
		if !dependencyExists(dep) && dep.ArchPackage != "" {
			packages = append(packages, strings.Fields(dep.ArchPackage)...)
		}
	}

	if len(packages) == 0 {
		fmt.Printf("%sNo packages to install.%s\n", ColorGreen, ColorReset)
		return
	}

	// Install packages
	args := append([]string{"pacman", "-S", "--noconfirm"}, packages...)
	cmd := exec.Command("sudo", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	if err != nil {
		fmt.Printf("%sError installing packages: %v%s\n", ColorRed, err, ColorReset)
		return
	}

	fmt.Printf("%sSuccessfully installed dependencies.%s\n", ColorGreen, ColorReset)
}

// installDebianDependencies installs dependencies for Debian/Ubuntu
func installDebianDependencies() {
	fmt.Printf("%sInstalling dependencies for Debian/Ubuntu...%s\n", ColorCyan, ColorReset)

	var packages []string
	for _, dep := range dependencies {
		if !dependencyExists(dep) && dep.DebPackage != "" {
			packages = append(packages, strings.Fields(dep.DebPackage)...)
		}
	}

	if len(packages) == 0 {
		fmt.Printf("%sNo packages to install.%s\n", ColorGreen, ColorReset)
		return
	}

	// Update package list
	cmd := exec.Command("sudo", "apt", "update")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Run()

	// Install packages
	args := append([]string{"apt", "install", "-y"}, packages...)
	cmd = exec.Command("sudo", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	if err != nil {
		fmt.Printf("%sError installing packages: %v%s\n", ColorRed, err, ColorReset)
		return
	}

	fmt.Printf("%sSuccessfully installed dependencies.%s\n", ColorGreen, ColorReset)
}

// buildBinary builds the LunarFetch binary
func buildBinary() bool {
	fmt.Printf("%sBuilding LunarFetch...%s\n", ColorCyan, ColorReset)

	// Check if go.mod exists
	if _, err := os.Stat("go.mod"); os.IsNotExist(err) {
		fmt.Printf("%sInitializing Go module...%s\n", ColorYellow, ColorReset)
		cmd := exec.Command("go", "mod", "init", "lunarfetch")
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("%sError initializing Go module: %v%s\n", ColorRed, err, ColorReset)
			fmt.Printf("%sOutput: %s%s\n", ColorRed, string(output), ColorReset)
			fmt.Printf("%sTrying to build without modules...%s\n", ColorYellow, ColorReset)
			return false
		} else {
			// Tidy the module to get dependencies
			fmt.Printf("%sTidying Go module...%s\n", ColorYellow, ColorReset)
			cmd = exec.Command("go", "mod", "tidy")
			output, err = cmd.CombinedOutput()
			if err != nil {
				fmt.Printf("%sWarning: Error tidying Go module: %v%s\n", ColorYellow, err, ColorReset)
				fmt.Printf("%sOutput: %s%s\n", ColorYellow, string(output), ColorReset)
				fmt.Printf("%sContinuing with build...%s\n", ColorYellow, ColorReset)
			}
		}
	}

	// Build the binary
	cmd := exec.Command("go", "build", "-o", "lunarfetch", ".")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("%sError building binary: %v%s\n", ColorRed, err, ColorReset)
		fmt.Printf("%sOutput: %s%s\n", ColorRed, string(output), ColorReset)

		// Try building with absolute module path
		fmt.Printf("%sTrying alternative build method...%s\n", ColorYellow, ColorReset)
		cmd = exec.Command("go", "build", "-o", "lunarfetch")
		output, err = cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("%sError building binary: %v%s\n", ColorRed, err, ColorReset)
			fmt.Printf("%sOutput: %s%s\n", ColorRed, string(output), ColorReset)
			return false
		}
	}

	// Check if the binary was created
	if _, err := os.Stat("lunarfetch"); os.IsNotExist(err) {
		fmt.Printf("%sError: Binary was not created%s\n", ColorRed, ColorReset)
		return false
	}

	fmt.Printf("%sSuccessfully built LunarFetch binary.%s\n", ColorGreen, ColorReset)
	return true
}

// installBinary installs the LunarFetch binary to the system
func installBinary() bool {
	fmt.Printf("%sInstalling LunarFetch binary to /usr/local/bin/...%s\n", ColorCyan, ColorReset)

	// Check if the binary exists
	if _, err := os.Stat("lunarfetch"); os.IsNotExist(err) {
		fmt.Printf("%sError: Binary not found. Please build it first.%s\n", ColorRed, ColorReset)
		return false
	}

	// Get the absolute path to the binary
	absPath, err := filepath.Abs("lunarfetch")
	if err != nil {
		fmt.Printf("%sError getting absolute path to binary: %v%s\n", ColorRed, err, ColorReset)
		return false
	}

	// Install the binary
	cmd := exec.Command("sudo", "cp", absPath, "/usr/local/bin/")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("%sError installing binary: %v%s\n", ColorRed, err, ColorReset)
		fmt.Printf("%sOutput: %s%s\n", ColorRed, string(output), ColorReset)

		// Try alternative installation method
		fmt.Printf("%sTrying alternative installation method...%s\n", ColorYellow, ColorReset)
		cmd = exec.Command("sudo", "install", "-m", "755", absPath, "/usr/local/bin/")
		output, err = cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("%sError installing binary: %v%s\n", ColorRed, err, ColorReset)
			fmt.Printf("%sOutput: %s%s\n", ColorRed, string(output), ColorReset)
			return false
		}
	}

	// Set permissions
	cmd = exec.Command("sudo", "chmod", "+x", "/usr/local/bin/lunarfetch")
	output, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("%sWarning: Could not set executable permissions: %v%s\n", ColorYellow, err, ColorReset)
		fmt.Printf("%sOutput: %s%s\n", ColorYellow, string(output), ColorReset)
	}

	fmt.Printf("%sLunarFetch has been installed successfully!%s\n", ColorGreen, ColorReset)
	fmt.Printf("%sYou can now run it from anywhere with the command: %slunarfetch%s\n", ColorGreen, ColorCyan, ColorReset)
	return true
}

// Helper functions

func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func dependencyExists(dep Dependency) bool {
	// Special case for lsb_release - we don't want to fail the dependency check
	// if lsb_release is not installed, as we're trying to install it
	if dep.Name == "LSB release" {
		// If we're checking for lsb_release and we're on Arch (detected by /etc/arch-release),
		// return false to allow installation to proceed
		if _, err := os.Stat("/etc/arch-release"); err == nil {
			for _, cmd := range dep.Commands {
				if cmd == "lsb_release" {
					// Check if it's actually installed
					if _, err := exec.LookPath(cmd); err == nil {
						return true
					}
					// Not installed, but we'll handle this specially
					return false
				}
			}
		}
	}

	// Normal dependency check
	for _, cmd := range dep.Commands {
		if commandExists(cmd) {
			return true
		}
	}
	return false
}

// findSourceDirectory attempts to find the LunarFetch source directory
func findSourceDirectory() (string, error) {
	// First, check if we're already in the source directory
	if _, err := os.Stat("main.go"); err == nil {
		// We're in the source directory
		absPath, err := filepath.Abs(".")
		if err != nil {
			return "", err
		}
		return absPath, nil
	}

	// Check if we're running from the binary and can find the source
	exePath, err := os.Executable()
	if err == nil {
		exeDir := filepath.Dir(exePath)

		// Check if the source is in the same directory as the binary
		if _, err := os.Stat(filepath.Join(exeDir, "main.go")); err == nil {
			return exeDir, nil
		}

		// Check if the source is in the parent directory of the binary
		parentDir := filepath.Dir(exeDir)
		if _, err := os.Stat(filepath.Join(parentDir, "main.go")); err == nil {
			return parentDir, nil
		}

		// Check if the source is in a sibling directory
		siblingDir := filepath.Join(parentDir, "LunarFetch")
		if _, err := os.Stat(filepath.Join(siblingDir, "main.go")); err == nil {
			return siblingDir, nil
		}
	}

	// Check if the source is in GOPATH
	gopath := os.Getenv("GOPATH")
	if gopath != "" {
		srcPath := filepath.Join(gopath, "src", "lunarfetch")
		if _, err := os.Stat(filepath.Join(srcPath, "main.go")); err == nil {
			return srcPath, nil
		}
	}

	// Try to find the source using go list
	cmd := exec.Command("go", "list", "-f", "{{.Dir}}", "lunarfetch")
	output, err := cmd.Output()
	if err == nil {
		srcPath := strings.TrimSpace(string(output))
		if _, err := os.Stat(filepath.Join(srcPath, "main.go")); err == nil {
			return srcPath, nil
		}
	}

	// If we can't find it, check common installation directories
	commonPaths := []string{
		"/usr/local/src/lunarfetch",
		"/usr/src/lunarfetch",
		"/opt/lunarfetch",
		"/home/nixev/Projects/LunarFetch", // Add the specific path for this user
	}

	for _, path := range commonPaths {
		if _, err := os.Stat(filepath.Join(path, "main.go")); err == nil {
			return path, nil
		}
	}

	// If we still can't find it, create a temporary directory with the necessary files
	tempDir, err := os.MkdirTemp("", "lunarfetch-install")
	if err == nil {
		fmt.Printf("%sCreating temporary installation directory: %s%s\n", ColorYellow, tempDir, ColorReset)

		// Create a minimal go.mod file
		modContent := `module lunarfetch

go 1.16
`
		err = os.WriteFile(filepath.Join(tempDir, "go.mod"), []byte(modContent), 0644)
		if err == nil {
			// Create a minimal main.go file that just imports the binary
			mainContent := `package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	// Just run the binary
	cmd := exec.Command("/usr/local/bin/lunarfetch")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Run()
}
`
			err = os.WriteFile(filepath.Join(tempDir, "main.go"), []byte(mainContent), 0644)
			if err == nil {
				return tempDir, nil
			}
		}
	}

	return "", fmt.Errorf("could not find LunarFetch source directory")
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

// setupImage sets up the image configuration
func setupImage() {
	// Find the source directory
	sourceDir, err := findSourceDirectory()
	if err != nil {
		fmt.Printf("%sError: %s%s\n", ColorRed, err.Error(), ColorReset)
		fmt.Printf("%sPlease run this command from the LunarFetch source directory or ensure Go modules are installed.%s\n", ColorYellow, ColorReset)
		return
	}

	fmt.Printf("%sFound LunarFetch source directory: %s%s\n", ColorGreen, sourceDir, ColorReset)

	// Save the current directory to return to it later
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("%sError: Could not get current directory: %s%s\n", ColorRed, err.Error(), ColorReset)
		return
	}

	// Change to the source directory for building
	err = os.Chdir(sourceDir)
	if err != nil {
		fmt.Printf("%sError: Could not change to source directory: %s%s\n", ColorRed, err.Error(), ColorReset)
		return
	}

	// Ensure we change back to the original directory when done
	defer func() {
		os.Chdir(currentDir)
	}()

	// Create config directory if it doesn't exist
	configDir := filepath.Join(os.Getenv("HOME"), ".config", "lunarfetch")
	err = os.MkdirAll(configDir, 0755)
	if err != nil {
		fmt.Printf("%sError: Could not create config directory: %s%s\n", ColorRed, err.Error(), ColorReset)
		return
	}

	// Create images directory if it doesn't exist
	imagesDir := filepath.Join(configDir, "images")
	err = os.MkdirAll(imagesDir, 0755)
	if err != nil {
		fmt.Printf("%sError: Could not create images directory: %s%s\n", ColorRed, err.Error(), ColorReset)
		return
	}

	// Load existing config or create a new one
	configPath := filepath.Join(configDir, "config.json")
	var config SimpleConfig

	if _, err := os.Stat(configPath); err == nil {
		// Config exists, load it
		configData, err := os.ReadFile(configPath)
		if err != nil {
			fmt.Printf("%sError reading config file: %s%s\n", ColorRed, err.Error(), ColorReset)
			return
		}

		err = json.Unmarshal(configData, &config)
		if err != nil {
			fmt.Printf("%sError parsing config file: %s%s\n", ColorRed, err.Error(), ColorReset)
			return
		}
	} else {
		// Config doesn't exist, create a default one
		config = SimpleConfig{
			Logo: Logo{
				Enabled: true,
				Type:    "ascii",
				Path:    "",
				Color:   "blue",
			},
			Info: Info{
				Enabled: true,
				Items: []string{
					"host", "user", "os", "kernel", "uptime", "terminal", "shell",
					"disk", "memory", "packages", "battery", "gpu", "cpu", "resolution",
					"wm", "theme", "icons", "desktop",
				},
			},
		}
	}

	// Copy sample image if it exists
	sampleImagePath := filepath.Join(sourceDir, "sample.png")
	if _, err := os.Stat(sampleImagePath); err == nil {
		destImagePath := filepath.Join(imagesDir, "sample.png")
		if _, err := os.Stat(destImagePath); os.IsNotExist(err) {
			err = copyFile(sampleImagePath, destImagePath)
			if err != nil {
				fmt.Printf("%sError: Could not copy sample image: %s%s\n", ColorRed, err.Error(), ColorReset)
			} else {
				fmt.Printf("%sCreated sample image at %s%s\n", ColorGreen, destImagePath, ColorReset)
			}
		}
	}

	fmt.Printf("%sImage setup complete!%s\n", ColorGreen, ColorReset)
	fmt.Printf("%sYou can now use images with LunarFetch by setting the image path in your config.%s\n", ColorGreen, ColorReset)
}
