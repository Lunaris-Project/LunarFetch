package utils

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg" // Register JPEG format
	"image/png"
	_ "image/png" // Register PNG format
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"io/ioutil"

	"github.com/disintegration/imaging"
	"github.com/mattn/go-sixel"
	_ "golang.org/x/image/webp" // Register WebP format
)

// Supported image formats
const (
	FormatPNG  = ".png"
	FormatJPG  = ".jpg"
	FormatJPEG = ".jpeg"
	FormatWEBP = ".webp"
)

// Supported rendering protocols
const (
	ProtocolAuto    = "auto"
	ProtocolSixel   = "sixel"
	ProtocolKitty   = "kitty"
	ProtocolITerm2  = "iterm2"
	ProtocolChafa   = "chafa"
	ProtocolUberzug = "uberzug"
)

// Supported display modes
const (
	DisplayModeAuto  = "auto"
	DisplayModeBlock = "block"
	DisplayModeASCII = "ascii"
)

// Supported render modes
const (
	RenderModeDetailed = "detailed"
	RenderModeSimple   = "simple"
	RenderModeBlock    = "block"
	RenderModeASCII    = "ascii"
)

// Supported dither modes
const (
	DitherModeNone           = "none"
	DitherModeFloydSteinberg = "floyd-steinberg"
)

// ImageLoader handles loading and processing images
type ImageLoader struct {
	Config ImageConfig
}

// ImageConfig holds configuration for image processing
type ImageConfig struct {
	ImagePath      string // Path to the image file or directory
	Width          int    // Width in terminal characters
	Height         int    // Height in terminal characters
	RenderMode     string // Rendering mode (detailed, simple, block, ascii)
	DitherMode     string // Dithering algorithm (none, floyd-steinberg)
	TerminalOutput bool   // Whether to output directly to terminal
	DisplayMode    string // Display mode (auto, block, ascii)
	Protocol       string // Image protocol (auto, sixel, kitty, iterm2, chafa, uberzug)
	Scale          int    // Image scaling factor
	Offset         int    // Offset from terminal edge
	Background     string // Background color (transparent or color value)
}

// NewImageLoader creates a new ImageLoader with the specified configuration
func NewImageLoader(config Config) *ImageLoader {
	// Initialize random seed for random image selection
	rand.Seed(time.Now().UnixNano())

	return &ImageLoader{
		Config: ImageConfig{
			ImagePath:      config.Image.ImagePath,
			Width:          config.Image.Width,
			Height:         config.Image.Height,
			RenderMode:     config.Image.RenderMode,
			DitherMode:     config.Image.DitherMode,
			TerminalOutput: config.Image.TerminalOutput,
			DisplayMode:    config.Image.DisplayMode,
			Protocol:       config.Image.Protocol,
			Scale:          config.Image.Scale,
			Offset:         config.Image.Offset,
			Background:     config.Image.Background,
		},
	}
}

// LoadImage loads an image from the specified path
func (i *ImageLoader) LoadImage(imagePath string) (image.Image, error) {
	// Expand home directory if path starts with ~
	fullPath, err := expandPath(imagePath)
	if err != nil {
		return nil, err
	}

	// Open the image file
	file, err := os.Open(fullPath)
	if err != nil {
		return nil, fmt.Errorf("error opening image: %v", err)
	}
	defer file.Close()

	// Decode the image
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("error decoding image: %v", err)
	}

	return img, nil
}

// expandPath expands the ~ in a path to the user's home directory
func expandPath(path string) (string, error) {
	if !strings.HasPrefix(path, "~") {
		return path, nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return strings.Replace(path, "~", homeDir, 1), nil
}

// ResizeImage resizes an image to the specified dimensions
func (i *ImageLoader) ResizeImage(img image.Image) image.Image {
	width := i.Config.Width
	height := i.Config.Height

	// If dimensions are not specified, use defaults
	if width <= 0 {
		width = 80
	}
	if height <= 0 {
		height = 24
	}

	// Apply scaling factor if specified
	if i.Config.Scale > 0 {
		width *= i.Config.Scale
		height *= i.Config.Scale
	}

	// Resize the image while preserving aspect ratio
	resized := imaging.Fit(img, width, height, imaging.Lanczos)
	return resized
}

// CalculateOptimalDimensions calculates the optimal dimensions for displaying an image
// based on its actual dimensions and the terminal size
func (i *ImageLoader) CalculateOptimalDimensions(img image.Image) (int, int) {
	// Get image dimensions
	imgWidth := img.Bounds().Dx()
	imgHeight := img.Bounds().Dy()

	// Get terminal dimensions
	termWidth, termHeight := getTerminalSize()

	// If terminal size detection failed, use reasonable defaults
	if termHeight <= 1 {
		termWidth = 80
		termHeight = 24
	}

	// Terminal characters are roughly twice as tall as they are wide
	// Adjust the aspect ratio to account for terminal character dimensions
	aspectRatio := float64(imgWidth) / float64(imgHeight) * 0.5

	// Use a reasonable portion of the terminal for the image
	maxWidth := int(float64(termWidth) * 0.6)
	maxHeight := int(float64(termHeight) * 0.6)

	// If user specified dimensions, respect them
	if i.Config.Width > 0 {
		maxWidth = i.Config.Width
	}
	if i.Config.Height > 0 {
		maxHeight = i.Config.Height
	}

	// Apply scaling factor if specified
	if i.Config.Scale > 0 {
		maxWidth *= i.Config.Scale
		maxHeight *= i.Config.Scale
	}

	// Start with maximum width
	displayWidth := maxWidth
	displayHeight := int(float64(displayWidth) / aspectRatio)

	// If height is too large, scale down
	if displayHeight > maxHeight {
		displayHeight = maxHeight
		displayWidth = int(float64(displayHeight) * aspectRatio)
	}

	// Ensure minimum dimensions
	if displayWidth < 20 {
		displayWidth = 20
	}
	if displayHeight < 10 {
		displayHeight = 10
	}

	// Ensure dimensions are reasonable for terminal display
	if displayWidth > termWidth {
		displayWidth = termWidth - 5
	}
	if displayHeight > termHeight {
		displayHeight = termHeight - 5
	}

	return displayWidth, displayHeight
}

// getTerminalSize returns the width and height of the terminal in characters
func getTerminalSize() (int, int) {
	// Try tput first as it's more reliable
	widthCmd := exec.Command("tput", "cols")
	widthCmd.Stdin = os.Stdin
	widthOut, widthErr := widthCmd.Output()

	heightCmd := exec.Command("tput", "lines")
	heightCmd.Stdin = os.Stdin
	heightOut, heightErr := heightCmd.Output()

	if widthErr == nil && heightErr == nil {
		width, _ := strconv.Atoi(strings.TrimSpace(string(widthOut)))
		height, _ := strconv.Atoi(strings.TrimSpace(string(heightOut)))
		return width, height
	}

	// Fallback to stty
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		// Default fallback values
		return 80, 24
	}

	parts := strings.Split(strings.TrimSpace(string(out)), " ")
	if len(parts) != 2 {
		return 80, 24
	}

	// stty returns rows columns (height width)
	height, _ := strconv.Atoi(parts[0])
	width, _ := strconv.Atoi(parts[1])

	return width, height
}

// DisplayWithSixel displays an image using Sixel graphics
func (i *ImageLoader) DisplayWithSixel(img image.Image) (string, error) {
	var buf bytes.Buffer
	enc := sixel.NewEncoder(&buf)
	enc.Dither = true

	err := enc.Encode(img)
	if err != nil {
		return "", fmt.Errorf("error encoding image to sixel: %v", err)
	}

	return buf.String(), nil
}

// DisplayWithKitty displays an image using Kitty graphics protocol
func (i *ImageLoader) DisplayWithKitty(img image.Image) (string, error) {
	// Convert image to PNG
	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	if err != nil {
		return "", fmt.Errorf("error encoding image to PNG: %v", err)
	}

	// Base64 encode the PNG data
	encoded := base64Encode(buf.Bytes())

	// Format the Kitty graphics protocol command
	// Use 'a=T' for temporary image and add explicit newline after the image
	cmd := fmt.Sprintf("\033_Ga=T,f=100,s=%d,v=%d;%s\033\\\n",
		img.Bounds().Dx(), img.Bounds().Dy(), encoded)

	return cmd, nil
}

// DisplayWithITerm2 displays an image using iTerm2 graphics protocol
func (i *ImageLoader) DisplayWithITerm2(img image.Image) (string, error) {
	// Convert image to PNG
	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	if err != nil {
		return "", fmt.Errorf("error encoding image to PNG: %v", err)
	}

	// Base64 encode the PNG data
	encoded := base64Encode(buf.Bytes())

	// Format the iTerm2 graphics protocol command
	cmd := fmt.Sprintf("\033]1337;File=inline=1;width=auto;height=auto;preserveAspectRatio=1:%s\a", encoded)

	return cmd, nil
}

// DisplayWithUberzug displays an image using Überzug
func (i *ImageLoader) DisplayWithUberzug(img image.Image, path string) (string, error) {
	// Check if Überzug is installed
	_, err := exec.LookPath("ueberzug")
	if err != nil {
		return "", fmt.Errorf("Überzug is not installed: %v", err)
	}

	// Create a temporary file for the image
	tmpFile, err := os.CreateTemp("", "lunarfetch-*.png")
	if err != nil {
		return "", fmt.Errorf("error creating temporary file: %v", err)
	}
	defer tmpFile.Close()

	// Save the image to the temporary file
	err = png.Encode(tmpFile, img)
	if err != nil {
		return "", fmt.Errorf("error saving image: %v", err)
	}

	// Get the terminal size
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("error getting terminal size: %v", err)
	}

	var rows, cols int
	fmt.Sscanf(string(out), "%d %d", &rows, &cols)

	// Calculate the position
	x := i.Config.Offset
	y := i.Config.Offset
	if x <= 0 {
		x = 1
	}
	if y <= 0 {
		y = 1
	}

	// Create the Überzug command
	uberzugCmd := fmt.Sprintf("ueberzug layer --parser json <<EOF\n"+
		"{\"action\": \"add\", \"identifier\": \"lunarfetch\", \"x\": %d, \"y\": %d, \"path\": \"%s\"}\n"+
		"sleep 5\n"+
		"{\"action\": \"remove\", \"identifier\": \"lunarfetch\"}\n"+
		"EOF", x, y, tmpFile.Name())

	// Execute the command in the background
	go func() {
		cmd := exec.Command("bash", "-c", uberzugCmd)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	}()

	return "", nil
}

// DisplayWithChafa displays an image using Chafa
func (i *ImageLoader) DisplayWithChafa(img image.Image) (string, error) {
	// Check if chafa is installed
	_, err := exec.LookPath("chafa")
	if err != nil {
		return "", fmt.Errorf("chafa is not installed")
	}

	// Create a temporary file for the image
	tmpFile, err := ioutil.TempFile("", "lunarfetch-*.png")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmpFile.Name())

	// Save the image to the temporary file
	err = png.Encode(tmpFile, img)
	if err != nil {
		return "", err
	}
	tmpFile.Close()

	// Simply run chafa with the image path - no extra parameters
	cmd := exec.Command("chafa", tmpFile.Name())
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	// Ensure the output ends with a newline
	result := string(output)
	if !strings.HasSuffix(result, "\n") {
		result += "\n"
	}

	return result, nil
}

// ApplyDithering applies dithering to an image
func (i *ImageLoader) ApplyDithering(img image.Image) image.Image {
	if i.Config.DitherMode == "none" {
		return img
	}

	// Implement Floyd-Steinberg dithering
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// Create a new RGBA image to hold the dithered result
	dithered := imaging.New(width, height, color.NRGBA{0, 0, 0, 0})

	// Copy the original image to the new one
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			dithered.Set(x, y, img.At(x, y))
		}
	}

	// Convert to grayscale if needed for simpler dithering
	if i.Config.DitherMode == "floyd-steinberg" {
		// Floyd-Steinberg dithering implementation
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				// Get the original color
				oldColor := dithered.At(x, y)
				r, g, b, a := oldColor.RGBA()

				// Convert to 8-bit color space
				oldR := uint8(r >> 8)
				oldG := uint8(g >> 8)
				oldB := uint8(b >> 8)
				oldA := uint8(a >> 8)

				// Find the closest color in the palette (for 256 colors)
				// This is a simplified version - just quantizing to 8 levels per channel
				newR := uint8((int(oldR) / 32) * 32)
				newG := uint8((int(oldG) / 32) * 32)
				newB := uint8((int(oldB) / 32) * 32)

				// Set the new color
				newColor := color.RGBA{newR, newG, newB, oldA}
				dithered.Set(x, y, newColor)

				// Calculate the error
				errR := int(oldR) - int(newR)
				errG := int(oldG) - int(newG)
				errB := int(oldB) - int(newB)

				// Distribute the error to neighboring pixels
				// Floyd-Steinberg distribution pattern:
				//     X   7/16
				// 3/16 5/16 1/16

				// Right pixel
				if x+1 < width {
					c := dithered.At(x+1, y)
					r, g, b, a := c.RGBA()
					dithered.Set(x+1, y, color.RGBA{
						uint8(clamp(int(r>>8) + errR*7/16)),
						uint8(clamp(int(g>>8) + errG*7/16)),
						uint8(clamp(int(b>>8) + errB*7/16)),
						uint8(a >> 8),
					})
				}

				// Bottom-left pixel
				if x-1 >= 0 && y+1 < height {
					c := dithered.At(x-1, y+1)
					r, g, b, a := c.RGBA()
					dithered.Set(x-1, y+1, color.RGBA{
						uint8(clamp(int(r>>8) + errR*3/16)),
						uint8(clamp(int(g>>8) + errG*3/16)),
						uint8(clamp(int(b>>8) + errB*3/16)),
						uint8(a >> 8),
					})
				}

				// Bottom pixel
				if y+1 < height {
					c := dithered.At(x, y+1)
					r, g, b, a := c.RGBA()
					dithered.Set(x, y+1, color.RGBA{
						uint8(clamp(int(r>>8) + errR*5/16)),
						uint8(clamp(int(g>>8) + errG*5/16)),
						uint8(clamp(int(b>>8) + errB*5/16)),
						uint8(a >> 8),
					})
				}

				// Bottom-right pixel
				if x+1 < width && y+1 < height {
					c := dithered.At(x+1, y+1)
					r, g, b, a := c.RGBA()
					dithered.Set(x+1, y+1, color.RGBA{
						uint8(clamp(int(r>>8) + errR*1/16)),
						uint8(clamp(int(g>>8) + errG*1/16)),
						uint8(clamp(int(b>>8) + errB*1/16)),
						uint8(a >> 8),
					})
				}
			}
		}
	}

	return dithered
}

// clamp ensures a value is between 0 and 255
func clamp(value int) int {
	if value < 0 {
		return 0
	}
	if value > 255 {
		return 255
	}
	return value
}

// RenderImage loads, processes, and renders an image
func (i *ImageLoader) RenderImage() (string, error) {
	// Load the image
	img, err := i.LoadImage(i.Config.ImagePath)
	if err != nil {
		return "", err
	}

	// Calculate optimal dimensions based on the actual image
	optWidth, optHeight := i.CalculateOptimalDimensions(img)

	// Update the config with the optimal dimensions
	originalWidth := i.Config.Width
	originalHeight := i.Config.Height
	i.Config.Width = optWidth
	i.Config.Height = optHeight

	// Resize the image
	img = i.ResizeImage(img)

	// Apply dithering if requested
	if i.Config.DitherMode != "none" {
		img = i.ApplyDithering(img)
	}

	// Render the image based on the protocol
	var output string

	switch i.Config.Protocol {
	case "sixel":
		output, err = i.DisplayWithSixel(img)
	case "kitty":
		output, err = i.DisplayWithKitty(img)
	case "iterm2":
		output, err = i.DisplayWithITerm2(img)
	case "uberzug":
		output, err = i.DisplayWithUberzug(img, i.Config.ImagePath)
	case "chafa":
		output, err = i.DisplayWithChafa(img)
	default:
		// Try to detect the best protocol
		output, err = i.AutoDetectProtocol(img)
	}

	// Restore original dimensions
	i.Config.Width = originalWidth
	i.Config.Height = originalHeight

	return output, err
}

// AutoDetectProtocol tries to detect the best protocol for the current terminal
func (i *ImageLoader) AutoDetectProtocol(img image.Image) (string, error) {
	// Check for TERM environment variable
	term := os.Getenv("TERM")

	// Check for terminal-specific environment variables
	if os.Getenv("KITTY_WINDOW_ID") != "" {
		return i.DisplayWithKitty(img)
	}

	if os.Getenv("ITERM_SESSION_ID") != "" {
		return i.DisplayWithITerm2(img)
	}

	// Check for sixel support
	if strings.Contains(term, "sixel") || strings.Contains(term, "mlterm") {
		return i.DisplayWithSixel(img)
	}

	// Try Chafa as a fallback
	_, err := exec.LookPath("chafa")
	if err == nil {
		return i.DisplayWithChafa(img)
	}

	// If all else fails, return an error
	return "", fmt.Errorf("no suitable image display protocol detected")
}

// GetRandomImage returns a random image from the image directory
func (i *ImageLoader) GetRandomImage() (string, error) {
	expandedPath, err := expandPath(i.Config.ImagePath)
	if err != nil {
		return "", fmt.Errorf("error expanding path: %v", err)
	}

	// Check if the path is a directory
	fileInfo, err := os.Stat(expandedPath)
	if err != nil {
		return "", fmt.Errorf("error accessing image path: %v", err)
	}

	// If it's a single file, just render it
	if !fileInfo.IsDir() {
		return i.RenderImage()
	}

	// If it's a directory, pick a random image
	files, err := os.ReadDir(expandedPath)
	if err != nil {
		return "", fmt.Errorf("error reading directory: %v", err)
	}

	var imageFiles []string
	for _, file := range files {
		ext := strings.ToLower(filepath.Ext(file.Name()))
		if ext == FormatPNG || ext == FormatJPG || ext == FormatJPEG || ext == FormatWEBP {
			imageFiles = append(imageFiles, filepath.Join(expandedPath, file.Name()))
		}
	}

	if len(imageFiles) == 0 {
		return "", fmt.Errorf("no image files found in %s", expandedPath)
	}

	// Save the original path
	originalPath := i.Config.ImagePath

	// Pick a random image
	randomIndex := rand.Intn(len(imageFiles))
	i.Config.ImagePath = imageFiles[randomIndex]

	// Render the image
	result, err := i.RenderImage()

	// Restore the original path
	i.Config.ImagePath = originalPath

	return result, err
}

// base64Encode encodes data to base64
func base64Encode(data []byte) string {
	const base64Chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	var result strings.Builder

	// Process 3 bytes at a time
	for i := 0; i < len(data); i += 3 {
		// Get the next 3 bytes (or fewer if at the end)
		b1 := data[i]
		b2 := byte(0)
		b3 := byte(0)

		if i+1 < len(data) {
			b2 = data[i+1]
		}
		if i+2 < len(data) {
			b3 = data[i+2]
		}

		// Convert to 4 base64 characters
		c1 := b1 >> 2
		c2 := ((b1 & 0x3) << 4) | (b2 >> 4)
		c3 := ((b2 & 0xF) << 2) | (b3 >> 6)
		c4 := b3 & 0x3F

		// Append the characters
		result.WriteByte(base64Chars[c1])
		result.WriteByte(base64Chars[c2])

		if i+1 < len(data) {
			result.WriteByte(base64Chars[c3])
		} else {
			result.WriteByte('=')
		}

		if i+2 < len(data) {
			result.WriteByte(base64Chars[c4])
		} else {
			result.WriteByte('=')
		}
	}

	return result.String()
}
