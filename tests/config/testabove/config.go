package testabove

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"lunarfetch/src/utils"
)

// CreateConfig creates a test configuration with image above system info
func CreateConfig() {
	// Create a test configuration
	config := utils.DefaultConfig()

	// Get home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting home directory:", err)
		os.Exit(1)
	}

	// Set up test configuration
	config.Logo.EnableLogo = false

	config.Image.Enabled = true
	config.Image.EnableImage = true
	config.Image.Position = "above" // Test image above system info
	config.Image.Protocol = "auto"
	config.Image.Width = 40
	config.Image.Height = 20
	config.Image.ImagePath = filepath.Join(homeDir, ".config", "lunarfetch", "images", "test.png")

	// Create a temporary config file
	configDir := filepath.Join(homeDir, ".config", "lunarfetch")
	err = os.MkdirAll(configDir, 0755)
	if err != nil {
		fmt.Println("Error creating config directory:", err)
		os.Exit(1)
	}

	// Save the test configuration
	configFile := filepath.Join(configDir, "test-above.json")
	configJSON, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling config:", err)
		os.Exit(1)
	}

	err = os.WriteFile(configFile, configJSON, 0644)
	if err != nil {
		fmt.Println("Error writing config file:", err)
		os.Exit(1)
	}

	fmt.Println("Test configuration created at:", configFile)
	fmt.Println("\nTest configuration settings:")
	fmt.Println("- Logo enabled:", config.Logo.EnableLogo)
	fmt.Println("- Image enabled:", config.Image.Enabled)
	fmt.Println("- Image position:", config.Image.Position)
	fmt.Println("- Image path:", config.Image.ImagePath)

	fmt.Println("\nTo test this configuration, run:")
	fmt.Printf("go run main.go -c %s\n", configFile)

	// Check if the test image exists
	if _, err := os.Stat(config.Image.ImagePath); os.IsNotExist(err) {
		fmt.Println("\nWarning: Test image not found at:", config.Image.ImagePath)
		fmt.Println("Please make sure the image exists before running the test.")
	} else {
		fmt.Println("\nTest image found at:", config.Image.ImagePath)
	}
}
