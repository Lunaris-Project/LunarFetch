package utils

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// LogoLoader handles loading ASCII art logos
type LogoLoader struct {
	LogoPath string
}

// NewLogoLoader creates a new LogoLoader with the specified path
func NewLogoLoader(logoPath string) *LogoLoader {
	return &LogoLoader{
		LogoPath: logoPath,
	}
}

// GetRandomLogo returns a random logo from the logo directory
func (l *LogoLoader) GetRandomLogo() (string, error) {
	expandedPath, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	logoPath := strings.Replace(l.LogoPath, "~", expandedPath, 1)

	files, err := ioutil.ReadDir(logoPath)
	if err != nil {
		log.Printf("Error reading directory: %v\n", err)
		return "", err
	}

	var logoFiles []string
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".txt" {
			logoFiles = append(logoFiles, filepath.Join(logoPath, file.Name()))
		}
	}

	if len(logoFiles) == 0 {
		return "", fmt.Errorf("no logo files found")
	}

	rand.Seed(time.Now().UnixNano())
	randomFile := logoFiles[rand.Intn(len(logoFiles))]
	content, err := ioutil.ReadFile(randomFile)
	if err != nil {
		log.Printf("Error reading logo file: %v\n", err)
		return "", err
	}

	return string(content), nil
}
