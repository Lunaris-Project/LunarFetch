package components

import (
	"strings"

	"lunarfetch/src/common"
)

// CPUInfo provides CPU information
type CPUInfo struct {
	SystemInfo
}

// GetInfo returns the CPU model
func (c *CPUInfo) GetInfo() string {
	out, err := common.GlobalCommandExecutor.Execute("lscpu")
	if err != nil {
		return "Unknown"
	}
	lines := strings.Split(out, "\n")
	for _, line := range lines {
		if strings.Contains(line, "Model name:") {
			fields := strings.Fields(line)
			return strings.Join(fields[2:], " ")
		}
	}
	return "Unknown"
}
