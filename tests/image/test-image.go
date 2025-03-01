package main

import (
	"fmt"
	"os"

	"lunarfetch/src/utils"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run test-image.go <image-path>")
		os.Exit(1)
	}

	imagePath := os.Args[1]

	// Create a test configuration
	config := utils.DefaultConfig()
	config.Image.Enabled = true
	config.Image.ImagePath = imagePath
	config.Image.Protocol = "auto" // Will auto-detect the best protocol

	// Create an image loader
	imageLoader := utils.NewImageLoader(config)

	// Try different protocols in order of preference if auto fails
	protocols := []string{"auto", "sixel", "kitty", "chafa", "terminal-image"}

	var output string
	var err error
	var lastError error

	// Try each protocol until one works
	for _, protocol := range protocols {
		config.Image.Protocol = protocol
		imageLoader = utils.NewImageLoader(config)
		output, err = imageLoader.RenderImage()
		if err == nil {
			fmt.Printf("Successfully rendered with protocol: %s\n", protocol)
			break
		}
		lastError = err
		fmt.Printf("Protocol %s failed: %v\n", protocol, err)
	}

	if err != nil {
		fmt.Printf("All protocols failed. Last error: %v\n", lastError)
		os.Exit(1)
	}

	// Display the image
	fmt.Println(output)
}
