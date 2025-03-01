package components

import (
	"fmt"
	"strings"

	"lunarfetch/src/common"
)

// PackagesInfo provides package count information
type PackagesInfo struct {
	SystemInfo
}

// GetInfo returns the number of installed packages
func (p *PackagesInfo) GetInfo() string {
	out, err := common.GlobalCommandExecutor.Execute("pacman", "-Qq")
	if err != nil {
		return "Unknown"
	}
	packages := strings.Split(out, "\n")
	return fmt.Sprintf("%d", len(packages))
}
