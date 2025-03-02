package components

import (
	"fmt"
	"os/exec"
	"strings"

	"lunarfetch/src/common"
)

// PackagesInfo provides package count information
type PackagesInfo struct {
	SystemInfo
}

// GetInfo returns the number of installed packages
func (p *PackagesInfo) GetInfo() string {
	// Check for pacman (Arch-based distributions)
	if _, err := exec.LookPath("pacman"); err == nil {
		out, err := common.GlobalCommandExecutor.Execute("pacman", "-Qq")
		if err == nil {
			packages := strings.Split(out, "\n")
			return fmt.Sprintf("%d", len(packages)-1) // Subtract 1 to account for the empty line at the end
		}
	}

	// Check for dpkg (Debian-based distributions)
	if _, err := exec.LookPath("dpkg"); err == nil {
		out, err := common.GlobalCommandExecutor.Execute("dpkg-query", "-f", "${binary:Package}\n", "-W")
		if err == nil {
			packages := strings.Split(out, "\n")
			return fmt.Sprintf("%d", len(packages)-1) // Subtract 1 to account for the empty line at the end
		}
	}

	// Check for rpm (Red Hat-based distributions)
	if _, err := exec.LookPath("rpm"); err == nil {
		out, err := common.GlobalCommandExecutor.Execute("rpm", "-qa")
		if err == nil {
			packages := strings.Split(out, "\n")
			return fmt.Sprintf("%d", len(packages)-1) // Subtract 1 to account for the empty line at the end
		}
	}

	// Check for flatpak
	if _, err := exec.LookPath("flatpak"); err == nil {
		out, err := common.GlobalCommandExecutor.Execute("flatpak", "list")
		if err == nil {
			packages := strings.Split(out, "\n")
			if len(packages) > 1 { // If there's at least one package (plus the empty line)
				return fmt.Sprintf("%d", len(packages)-1)
			}
		}
	}

	return "Unknown"
}
