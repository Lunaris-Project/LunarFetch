package components

import (
	"fmt"
	"strings"

	"lunarfetch/src/common"
)

// MemoryInfo provides memory usage information
type MemoryInfo struct {
	SystemInfo
}

// GetInfo returns the memory usage
func (m *MemoryInfo) GetInfo() string {
	out, err := common.GlobalCommandExecutor.Execute("free", "-m")
	if err != nil {
		return "Unknown"
	}
	lines := strings.Split(out, "\n")
	fields := strings.Fields(lines[1])
	used := fields[2]
	total := fields[1]
	return fmt.Sprintf("%sMiB / %sMiB", used, total)
}
