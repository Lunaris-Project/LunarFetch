package main

import (
	"fmt"
	"os"

	"lunarfetch/src/scripts"
	"lunarfetch/src/utils"
)

// Version information
var (
	Version     = scripts.Version
	VersionDate = scripts.VersionDate
)

// Color codes for terminal output
const (
	ColorReset  = scripts.ColorReset
	ColorRed    = scripts.ColorRed
	ColorGreen  = scripts.ColorGreen
	ColorYellow = scripts.ColorYellow
	ColorBlue   = scripts.ColorBlue
	ColorCyan   = scripts.ColorCyan
)

// Dependency represents a system dependency
type Dependency = scripts.Dependency

// SimpleConfig is a simplified configuration structure for installation
type SimpleConfig = scripts.SimpleConfig

// Logo represents the logo configuration for installation
type Logo = scripts.Logo

// Info represents the information display configuration for installation
type Info = scripts.Info

func main() {
	// Parse command line arguments and get config path if specified
	configPath, shouldExit := parseCommandLineArgs()
	if shouldExit {
		return
	}

	// Load configuration
	config := loadConfiguration(configPath)

	// Run the main application
	runLunarFetch(config)
}

// parseCommandLineArgs processes command line arguments and returns the config path if specified
// Returns a boolean indicating if the program should exit after argument processing
func parseCommandLineArgs() (string, bool) {
	var configPath string

	// No arguments provided, just run normally
	if len(os.Args) <= 1 {
		return "", false
	}

	// Check for help flag first
	if os.Args[1] == "--help" || os.Args[1] == "-h" {
		scripts.PrintUsage()
		return "", true
	}

	// Check for version flag
	if os.Args[1] == "--version" || os.Args[1] == "-v" {
		scripts.PrintVersion()
		return "", true
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

	// Check if there are any command arguments left
	if len(os.Args) > 1 {
		scripts.HandleCommands(os.Args[1:])
		return "", true
	}

	return configPath, false
}

// loadConfiguration loads the configuration from the specified path or default location
func loadConfiguration(configPath string) utils.Config {
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

	return config
}

// runLunarFetch runs the main application with the provided configuration
func runLunarFetch(config utils.Config) {
	// Initialize display manager
	displayManager := utils.NewDisplayManager(config)
	displayManager.InitializeComponents()

	// Get system information
	sysInfoOutput := displayManager.Display()

	// Load logo and image if enabled
	logoOutput := loadLogo(config)
	imageOutput := loadImage(config)

	// Display the information based on configuration
	displayOutput(config, sysInfoOutput, logoOutput, imageOutput)
}

// loadLogo loads the logo if enabled in the configuration
func loadLogo(config utils.Config) string {
	if !config.Logo.EnableLogo {
		return ""
	}

	logoLoader := utils.NewLogoLoader(config.Logo.LogoPath)
	logoOutput, err := logoLoader.GetRandomLogo()
	if err != nil && os.Getenv("LUNARFETCH_DEBUG") == "1" {
		fmt.Printf("Error loading logo: %v\n", err)
	}
	return logoOutput
}

// loadImage loads and renders the image if enabled in the configuration
func loadImage(config utils.Config) string {
	if !config.Image.Enabled && !config.Image.EnableImage {
		return ""
	}

	imageLoader := utils.NewImageLoader(config)

	// Debug output
	if os.Getenv("LUNARFETCH_DEBUG") == "1" {
		printImageDebugInfo(config)
	}

	// Use GetRandomImage if random is enabled, otherwise use RenderImage
	var imageOutput string
	var err error

	if config.Image.Random {
		imageOutput, err = imageLoader.GetRandomImage()
	} else {
		imageOutput, err = imageLoader.RenderImage()
	}

	if err != nil && os.Getenv("LUNARFETCH_DEBUG") == "1" {
		fmt.Printf("Error rendering image: %v\n", err)
	}

	return imageOutput
}

// printImageDebugInfo prints debug information about image configuration
func printImageDebugInfo(config utils.Config) {
	fmt.Printf("Image settings:\n")
	fmt.Printf("  Enabled: %v\n", config.Image.EnableImage)
	fmt.Printf("  Path: %s\n", config.Image.ImagePath)
	fmt.Printf("  Width: %d\n", config.Image.Width)
	fmt.Printf("  Height: %d\n", config.Image.Height)
	fmt.Printf("  Position: %s\n", config.Image.Position)
	fmt.Printf("  Random: %v\n", config.Image.Random)
}

// displayOutput displays the system information, logo, and image based on configuration
func displayOutput(config utils.Config, sysInfoOutput, logoOutput, imageOutput string) {
	// Check if image should be displayed above
	if config.Image.Position == "above" && (config.Image.Enabled || config.Image.EnableImage) && imageOutput != "" {
		// Display image above system info
		fmt.Print(imageOutput)
		// Add a newline to ensure separation between image and system info
		fmt.Println()
		fmt.Print(sysInfoOutput)
		return
	}

	// Check if logo should be displayed above
	if config.Logo.Position == "above" && config.Logo.EnableLogo && logoOutput != "" {
		// Display logo above system info
		fmt.Print(logoOutput)
		fmt.Println()
		fmt.Print(sysInfoOutput)
		return
	}

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
