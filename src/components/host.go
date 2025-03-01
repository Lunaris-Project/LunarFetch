package components

import (
	"os"
	"strings"

	"lunarfetch/src/common"
)

// HostInfo provides hostname information
type HostInfo struct {
	SystemInfo
}

// GetInfo returns the hostname
func (h *HostInfo) GetInfo() string {
	hostname, err := os.Hostname()
	if err != nil {
		// Try using command execution as fallback
		out, err := common.GlobalCommandExecutor.Execute("hostname")
		if err != nil {
			return "Unknown"
		}
		return strings.TrimSpace(out)
	}
	return hostname
}
