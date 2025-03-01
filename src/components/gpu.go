package components

import (
	"strings"

	"lunarfetch/src/common"
)

// GPUInfo provides GPU information
type GPUInfo struct {
	SystemInfo
}

// GetInfo returns the GPU model
func (g *GPUInfo) GetInfo() string {
	out, err := common.GlobalCommandExecutor.Execute("lspci")
	if err != nil {
		return "Unknown"
	}
	lines := strings.Split(out, "\n")
	for _, line := range lines {
		if strings.Contains(line, "VGA") || strings.Contains(line, "3D") {
			fields := strings.Fields(line)
			return strings.Join(fields[4:], " ")
		}
	}
	return "Unknown"
}
