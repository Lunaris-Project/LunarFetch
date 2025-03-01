package utils

import (
	"fmt"
	"strings"
)

// BoxConfig holds the configuration for box drawing
type BoxConfig struct {
	TopLeft     string
	TopRight    string
	BottomLeft  string
	BottomRight string
	TopEdge     string
	BottomEdge  string
	LeftEdge    string
	RightEdge   string
	Separator   string
}

// BoxDrawer handles drawing boxes around content
type BoxDrawer struct {
	Config BoxConfig
}

// NewBoxDrawer creates a new BoxDrawer with the specified configuration
func NewBoxDrawer(config BoxConfig) *BoxDrawer {
	return &BoxDrawer{
		Config: config,
	}
}

// Draw draws a box around the provided content
func (b *BoxDrawer) Draw(content string) string {
	lines := strings.Split(content, "\n")
	maxLen := 0
	for _, line := range lines {
		if len(line) > maxLen {
			maxLen = len(line)
		}
	}

	var box strings.Builder

	box.WriteString(b.Config.TopLeft)
	box.WriteString(strings.Repeat(b.Config.TopEdge, maxLen+2))
	box.WriteString(b.Config.TopRight + "\n")

	for _, line := range lines {
		box.WriteString(b.Config.LeftEdge + " ")
		box.WriteString(fmt.Sprintf("%-*s", maxLen, line))
		box.WriteString(" " + b.Config.RightEdge + "\n")
	}

	box.WriteString(b.Config.BottomLeft)
	box.WriteString(strings.Repeat(b.Config.BottomEdge, maxLen+2))
	box.WriteString(b.Config.BottomRight + "\n")

	return box.String()
}
