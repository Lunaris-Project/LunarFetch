package main

import (
	"fmt"
	"os"
	"strings"

	"lunarfetch/tests/config/testabove"
	"lunarfetch/tests/config/testside"
)

func main() {
	if len(os.Args) > 1 {
		arg := strings.ToLower(os.Args[1])

		if arg == "above" {
			fmt.Println("Creating test configuration with image above system info...")
			testabove.CreateConfig()
		} else if arg == "side" {
			fmt.Println("Creating test configuration with image beside system info...")
			testside.CreateConfig()
		} else {
			printUsage()
		}
	} else {
		// Default to creating both configurations
		fmt.Println("Creating both test configurations...")
		fmt.Println("\n=== Configuration with image above system info ===")
		testabove.CreateConfig()

		fmt.Println("\n=== Configuration with image beside system info ===")
		testside.CreateConfig()
	}
}

func printUsage() {
	fmt.Println("Usage: go run tests/config/main.go [option]")
	fmt.Println("Options:")
	fmt.Println("  above    Create test configuration with image above system info")
	fmt.Println("  side     Create test configuration with image beside system info")
	fmt.Println("  (none)   Create both test configurations")
}
