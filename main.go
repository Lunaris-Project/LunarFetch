package main

import (
	"fmt"
	"os"
	"strings"

	"lunarfetch/src/scripts"
	"lunarfetch/src/utils"
)

var (
	Version     = scripts.Version
	VersionDate = scripts.VersionDate
)

const (
	ColorReset  = scripts.ColorReset
	ColorRed    = scripts.ColorRed
	ColorGreen  = scripts.ColorGreen
	ColorYellow = scripts.ColorYellow
	ColorBlue   = scripts.ColorBlue
	ColorCyan   = scripts.ColorCyan
)

type Dependency = scripts.Dependency

type SimpleConfig = scripts.SimpleConfig

type Logo = scripts.Logo

type Info = scripts.Info

func main() {

	configPath, shouldExit := parseCommandLineArgs()
	if shouldExit {
		return
	}

	config := loadConfiguration(configPath)

	runLunarFetch(config)
}

func parseCommandLineArgs() (string, bool) {
	var configPath string

	if len(os.Args) <= 1 {
		return "", false
	}

	if os.Args[1] == "--help" || os.Args[1] == "-h" {
		scripts.PrintUsage()
		return "", true
	}

	if os.Args[1] == "--version" || os.Args[1] == "-v" {
		scripts.PrintVersion()
		return "", true
	}

	if os.Args[1] == "--debug" || os.Args[1] == "-d" {

		os.Setenv("LUNARFETCH_DEBUG", "1")
		if len(os.Args) > 2 {

			os.Args = append([]string{os.Args[0]}, os.Args[2:]...)
		} else {

			os.Args = []string{os.Args[0]}
		}
	}

	for i := 1; i < len(os.Args); i++ {
		if os.Args[i] == "--config" || os.Args[i] == "-c" {
			if i+1 < len(os.Args) {
				configPath = os.Args[i+1]

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
		scripts.HandleCommands(os.Args[1:])
		return "", true
	}

	return configPath, false
}

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

func runLunarFetch(config utils.Config) {

	displayManager := utils.NewDisplayManager(config)
	displayManager.InitializeComponents()

	sysInfoOutput := displayManager.Display()

	logoOutput := loadLogo(config)
	imageOutput := loadImage(config)

	displayOutput(config, sysInfoOutput, logoOutput, imageOutput)
}

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

func loadImage(config utils.Config) string {
	if !config.Image.EnableImage {
		return ""
	}

	imageLoader := utils.NewImageLoader(config)

	if os.Getenv("LUNARFETCH_DEBUG") == "1" {
		printImageDebugInfo(config)
	}

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

func printImageDebugInfo(config utils.Config) {
	fmt.Printf("Image settings:\n")
	fmt.Printf("  Enabled: %v\n", config.Image.EnableImage)
	fmt.Printf("  Path: %s\n", config.Image.ImagePath)
	fmt.Printf("  Width: %d\n", config.Image.Width)
	fmt.Printf("  Height: %d\n", config.Image.Height)
	fmt.Printf("  Position: %s\n", config.Image.Position)
	fmt.Printf("  Random: %v\n", config.Image.Random)
}

func addSimpleMargin(sysInfo string, spaces int) string {
	margin := strings.Repeat(" ", spaces)
	lines := strings.Split(sysInfo, "\n")
	for i, line := range lines {
		lines[i] = margin + line
	}
	return strings.Join(lines, "\n")
}

func displayOutput(config utils.Config, sysInfo, logoOutput, imageOutput string) {
	if os.Getenv("LUNARFETCH_DEBUG") == "1" {
		fmt.Printf("Logo enabled: %v, position: %s\n", config.Logo.EnableLogo, config.Logo.Position)
		fmt.Printf("Image enabled: %v, position: %s\n", config.Image.EnableImage, config.Image.Position)
	}

	sysInfo = strings.TrimSpace(sysInfo)
	logoOutput = strings.TrimSpace(logoOutput)
	imageOutput = strings.TrimSpace(imageOutput)

	var topContent, middleContent, bottomContent string
	middleContent = sysInfo

	if config.Logo.EnableLogo {
		switch config.Logo.Position {
		case "above":
			if topContent == "" {
				topContent = logoOutput
			} else {
				topContent = topContent + "\n" + logoOutput
			}
		case "below":
			if bottomContent == "" {
				bottomContent = logoOutput
			} else {
				bottomContent = bottomContent + "\n" + logoOutput
			}
		}
	}

	if config.Image.EnableImage {
		switch config.Image.Position {
		case "above":
			if topContent == "" {
				topContent = imageOutput
			} else {
				topContent = topContent + "\n" + imageOutput
			}
		case "below":
			if bottomContent == "" {
				bottomContent = imageOutput
			} else {
				bottomContent = bottomContent + "\n" + imageOutput
			}
		}
	}

	var result strings.Builder

	if topContent != "" {
		result.WriteString(topContent + "\n")
	}

	if (config.Logo.EnableLogo && (config.Logo.Position == "left" || config.Logo.Position == "right")) ||
		(config.Image.EnableImage && (config.Image.Position == "left" || config.Image.Position == "right")) {

		if config.Logo.EnableLogo && config.Logo.Position == "right" &&
			config.Image.EnableImage && config.Image.Position == "left" {

			combined := mergeSideBySide(normalizeOutput(imageOutput), sysInfo)

			result.WriteString(mergeSideBySide(combined, logoOutput))
		} else if config.Logo.EnableLogo && config.Logo.Position == "left" &&
			config.Image.EnableImage && config.Image.Position == "right" {

			combined := mergeSideBySide(logoOutput, sysInfo)

			result.WriteString(mergeSideBySide(combined, normalizeOutput(imageOutput)))
		} else if config.Logo.EnableLogo && config.Logo.Position == "left" {

			result.WriteString(mergeSideBySide(logoOutput, sysInfo))
		} else if config.Logo.EnableLogo && config.Logo.Position == "right" {

			result.WriteString(mergeSideBySide(sysInfo, logoOutput))
		} else if config.Image.EnableImage && config.Image.Position == "left" {

			normalizedImage := normalizeOutput(imageOutput)
			result.WriteString(mergeSideBySide(normalizedImage, sysInfo))
		} else if config.Image.EnableImage && config.Image.Position == "right" {

			result.WriteString(mergeSideBySide(sysInfo, normalizeOutput(imageOutput)))
		}
	} else {

		result.WriteString(middleContent)
	}

	if bottomContent != "" {
		result.WriteString("\n" + bottomContent)
	}

	finalOutput := result.String()
	if !strings.HasSuffix(finalOutput, "\n") {
		finalOutput += "\n"
	}

	fmt.Print(finalOutput)
}

func mergeSideBySide(left, right string) string {
	if left == "" {
		return right
	}
	if right == "" {
		return left
	}

	leftLines := strings.Split(strings.TrimRight(left, "\n"), "\n")
	rightLines := strings.Split(strings.TrimRight(right, "\n"), "\n")

	for i := range leftLines {
		leftLines[i] = strings.TrimRight(leftLines[i], " ")
	}
	for i := range rightLines {
		rightLines[i] = strings.TrimRight(rightLines[i], " ")
	}

	for len(leftLines) > 0 && strings.TrimSpace(leftLines[len(leftLines)-1]) == "" {
		leftLines = leftLines[:len(leftLines)-1]
	}
	for len(rightLines) > 0 && strings.TrimSpace(rightLines[len(rightLines)-1]) == "" {
		rightLines = rightLines[:len(rightLines)-1]
	}

	maxLeftWidth := 0
	for _, line := range leftLines {
		width := len([]rune(line))
		if width > maxLeftWidth {
			maxLeftWidth = width
		}
	}

	padding := 2
	totalPadding := maxLeftWidth + padding

	var result strings.Builder
	maxLines := len(leftLines)
	if len(rightLines) > maxLines {
		maxLines = len(rightLines)
	}

	for i := 0; i < maxLines; i++ {
		var leftLine, rightLine string

		if i < len(leftLines) {
			leftLine = leftLines[i]
		}
		if i < len(rightLines) {
			rightLine = rightLines[i]
		}

		if leftLine == "" && rightLine == "" {
			continue
		}

		currentPadding := totalPadding - len([]rune(leftLine))
		if currentPadding < 0 {
			currentPadding = padding
		}

		if leftLine != "" {
			result.WriteString(leftLine)
			if rightLine != "" {
				result.WriteString(strings.Repeat(" ", currentPadding))
			}
		} else if rightLine != "" {

			result.WriteString(strings.Repeat(" ", totalPadding))
		}

		if rightLine != "" {
			result.WriteString(rightLine)
		}

		if i < maxLines-1 {
			result.WriteString("\n")
		}
	}

	return result.String()
}

func normalizeOutput(output string) string {
	lines := strings.Split(strings.TrimRight(output, "\n"), "\n")

	for i := range lines {
		lines[i] = strings.TrimRight(lines[i], " ")
	}

	for len(lines) > 0 && strings.TrimSpace(lines[len(lines)-1]) == "" {
		lines = lines[:len(lines)-1]
	}

	maxLength := 0
	for _, line := range lines {
		lineLength := len([]rune(line))
		if lineLength > maxLength {
			maxLength = lineLength
		}
	}

	for i, line := range lines {
		currentLength := len([]rune(line))
		if currentLength < maxLength {

			lines[i] = line + strings.Repeat(" ", maxLength-currentLength)
		}
	}

	return strings.Join(lines, "\n")
}
