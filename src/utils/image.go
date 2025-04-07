package utils

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
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
	_ "golang.org/x/image/webp"
)

const (
	FormatPNG  = ".png"
	FormatJPG  = ".jpg"
	FormatJPEG = ".jpeg"
	FormatWEBP = ".webp"
)

const (
	ProtocolAuto    = "auto"
	ProtocolSixel   = "sixel"
	ProtocolKitty   = "kitty"
	ProtocolITerm2  = "iterm2"
	ProtocolChafa   = "chafa"
	ProtocolUberzug = "uberzug"
)

const (
	DisplayModeAuto  = "auto"
	DisplayModeBlock = "block"
	DisplayModeASCII = "ascii"
)

const (
	RenderModeDetailed = "detailed"
	RenderModeSimple   = "simple"
	RenderModeBlock    = "block"
	RenderModeASCII    = "ascii"
)

const (
	DitherModeNone           = "none"
	DitherModeFloydSteinberg = "floyd-steinberg"
)

type ImageLoader struct {
	Config ImageConfig
}

type ImageConfig struct {
	ImagePath      string
	Width          int
	Height         int
	RenderMode     string
	DitherMode     string
	TerminalOutput bool
	DisplayMode    string
	Protocol       string
	Scale          int
	Offset         int
	Background     string
}

func NewImageLoader(config Config) *ImageLoader {

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

func (i *ImageLoader) LoadImage(imagePath string) (image.Image, error) {

	fullPath, err := expandPath(imagePath)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(fullPath)
	if err != nil {
		return nil, fmt.Errorf("error opening image: %v", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("error decoding image: %v", err)
	}

	return img, nil
}

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

func (i *ImageLoader) ResizeImage(img image.Image) image.Image {
	width := i.Config.Width
	height := i.Config.Height

	if width <= 0 {
		width = 80
	}
	if height <= 0 {
		height = 24
	}

	if i.Config.Scale > 0 {
		width *= i.Config.Scale
		height *= i.Config.Scale
	}

	resized := imaging.Fit(img, width, height, imaging.Lanczos)
	return resized
}

func (i *ImageLoader) CalculateOptimalDimensions(img image.Image) (int, int) {

	imgWidth := img.Bounds().Dx()
	imgHeight := img.Bounds().Dy()

	termWidth, termHeight := getTerminalSize()

	if termHeight <= 1 {
		termWidth = 80
		termHeight = 24
	}

	aspectRatio := float64(imgWidth) / float64(imgHeight) * 0.5

	maxWidth := int(float64(termWidth) * 0.6)
	maxHeight := int(float64(termHeight) * 0.6)

	if i.Config.Width > 0 {
		maxWidth = i.Config.Width
	}
	if i.Config.Height > 0 {
		maxHeight = i.Config.Height
	}

	if i.Config.Scale > 0 {
		maxWidth *= i.Config.Scale
		maxHeight *= i.Config.Scale
	}

	displayWidth := maxWidth
	displayHeight := int(float64(displayWidth) / aspectRatio)

	if displayHeight > maxHeight {
		displayHeight = maxHeight
		displayWidth = int(float64(displayHeight) * aspectRatio)
	}

	if displayWidth < 20 {
		displayWidth = 20
	}
	if displayHeight < 10 {
		displayHeight = 10
	}

	if displayWidth > termWidth {
		displayWidth = termWidth - 5
	}
	if displayHeight > termHeight {
		displayHeight = termHeight - 5
	}

	return displayWidth, displayHeight
}

func getTerminalSize() (int, int) {

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

	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {

		return 80, 24
	}

	parts := strings.Split(strings.TrimSpace(string(out)), " ")
	if len(parts) != 2 {
		return 80, 24
	}

	height, _ := strconv.Atoi(parts[0])
	width, _ := strconv.Atoi(parts[1])

	return width, height
}

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

func (i *ImageLoader) DisplayWithKitty(img image.Image) (string, error) {

	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	if err != nil {
		return "", fmt.Errorf("error encoding image to PNG: %v", err)
	}

	encoded := base64Encode(buf.Bytes())

	cmd := fmt.Sprintf("\033_Ga=T,f=100,s=%d,v=%d;%s\033\\\n",
		img.Bounds().Dx(), img.Bounds().Dy(), encoded)

	return cmd, nil
}

func (i *ImageLoader) DisplayWithITerm2(img image.Image) (string, error) {

	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	if err != nil {
		return "", fmt.Errorf("error encoding image to PNG: %v", err)
	}

	encoded := base64Encode(buf.Bytes())

	cmd := fmt.Sprintf("\033]1337;File=inline=1;width=auto;height=auto;preserveAspectRatio=1:%s\a", encoded)

	return cmd, nil
}

func (i *ImageLoader) DisplayWithUberzug(img image.Image, path string) (string, error) {

	_, err := exec.LookPath("ueberzug")
	if err != nil {
		return "", fmt.Errorf("uberzug is not installed: %v", err)
	}

	tmpFile, err := os.CreateTemp("", "lunarfetch-*.png")
	if err != nil {
		return "", fmt.Errorf("error creating temporary file: %v", err)
	}
	defer tmpFile.Close()

	err = png.Encode(tmpFile, img)
	if err != nil {
		return "", fmt.Errorf("error saving image: %v", err)
	}

	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("error getting terminal size: %v", err)
	}

	var rows, cols int
	fmt.Sscanf(string(out), "%d %d", &rows, &cols)

	x := i.Config.Offset
	y := i.Config.Offset
	if x <= 0 {
		x = 1
	}
	if y <= 0 {
		y = 1
	}

	uberzugCmd := fmt.Sprintf("ueberzug layer --parser json <<EOF\n"+
		"{\"action\": \"add\", \"identifier\": \"lunarfetch\", \"x\": %d, \"y\": %d, \"path\": \"%s\"}\n"+
		"sleep 5\n"+
		"{\"action\": \"remove\", \"identifier\": \"lunarfetch\"}\n"+
		"EOF", x, y, tmpFile.Name())

	go func() {
		cmd := exec.Command("bash", "-c", uberzugCmd)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	}()

	return "", nil
}

func (i *ImageLoader) DisplayWithChafa(img image.Image) (string, error) {

	_, err := exec.LookPath("chafa")
	if err != nil {
		return "", fmt.Errorf("chafa is not installed")
	}

	tmpFile, err := ioutil.TempFile("", "lunarfetch-*.png")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmpFile.Name())

	err = png.Encode(tmpFile, img)
	if err != nil {
		return "", err
	}
	tmpFile.Close()

	var chafaArgs []string

	chafaArgs = append(chafaArgs, "--size", fmt.Sprintf("%dx%d", i.Config.Width, i.Config.Height))

	chafaArgs = append(chafaArgs, tmpFile.Name())

	cmd := exec.Command("chafa", chafaArgs...)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	result := string(output)

	if !strings.HasSuffix(result, "\n") {
		result += "\n"
	}

	return result, nil
}

func (i *ImageLoader) ApplyDithering(img image.Image) image.Image {
	if i.Config.DitherMode == "none" {
		return img
	}

	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	dithered := imaging.New(width, height, color.NRGBA{0, 0, 0, 0})

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			dithered.Set(x, y, img.At(x, y))
		}
	}

	if i.Config.DitherMode == "floyd-steinberg" {

		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {

				oldColor := dithered.At(x, y)
				r, g, b, a := oldColor.RGBA()

				oldR := uint8(r >> 8)
				oldG := uint8(g >> 8)
				oldB := uint8(b >> 8)
				oldA := uint8(a >> 8)

				newR := uint8((int(oldR) / 32) * 32)
				newG := uint8((int(oldG) / 32) * 32)
				newB := uint8((int(oldB) / 32) * 32)

				newColor := color.RGBA{newR, newG, newB, oldA}
				dithered.Set(x, y, newColor)

				errR := int(oldR) - int(newR)
				errG := int(oldG) - int(newG)
				errB := int(oldB) - int(newB)

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

func clamp(value int) int {
	if value < 0 {
		return 0
	}
	if value > 255 {
		return 255
	}
	return value
}

func (i *ImageLoader) RenderImage() (string, error) {

	img, err := i.LoadImage(i.Config.ImagePath)
	if err != nil {
		return "", err
	}

	optWidth, optHeight := i.CalculateOptimalDimensions(img)

	originalWidth := i.Config.Width
	originalHeight := i.Config.Height
	i.Config.Width = optWidth
	i.Config.Height = optHeight

	img = i.ResizeImage(img)

	if i.Config.DitherMode != "none" {
		img = i.ApplyDithering(img)
	}

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

		output, err = i.AutoDetectProtocol(img)
	}

	i.Config.Width = originalWidth
	i.Config.Height = originalHeight

	return output, err
}

func (i *ImageLoader) AutoDetectProtocol(img image.Image) (string, error) {

	term := os.Getenv("TERM")

	if os.Getenv("KITTY_WINDOW_ID") != "" {
		return i.DisplayWithKitty(img)
	}

	if os.Getenv("ITERM_SESSION_ID") != "" {
		return i.DisplayWithITerm2(img)
	}

	if strings.Contains(term, "sixel") || strings.Contains(term, "mlterm") {
		return i.DisplayWithSixel(img)
	}

	_, err := exec.LookPath("chafa")
	if err == nil {
		return i.DisplayWithChafa(img)
	}

	return "", fmt.Errorf("no suitable image display protocol detected")
}

func (i *ImageLoader) GetRandomImage() (string, error) {
	expandedPath, err := expandPath(i.Config.ImagePath)
	if err != nil {
		return "", fmt.Errorf("error expanding path: %v", err)
	}

	fileInfo, err := os.Stat(expandedPath)
	if err != nil {
		return "", fmt.Errorf("error accessing image path: %v", err)
	}

	if !fileInfo.IsDir() {
		return i.RenderImage()
	}

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

	originalPath := i.Config.ImagePath

	randomIndex := rand.Intn(len(imageFiles))
	i.Config.ImagePath = imageFiles[randomIndex]

	result, err := i.RenderImage()

	i.Config.ImagePath = originalPath

	return result, err
}

func base64Encode(data []byte) string {
	const base64Chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	var result strings.Builder

	for i := 0; i < len(data); i += 3 {

		b1 := data[i]
		b2 := byte(0)
		b3 := byte(0)

		if i+1 < len(data) {
			b2 = data[i+1]
		}
		if i+2 < len(data) {
			b3 = data[i+2]
		}

		c1 := b1 >> 2
		c2 := ((b1 & 0x3) << 4) | (b2 >> 4)
		c3 := ((b2 & 0xF) << 2) | (b3 >> 6)
		c4 := b3 & 0x3F

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
