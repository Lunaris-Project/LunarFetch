package components

import (
	"fmt"
	"strconv"
	"strings"

	"lunarfetch/src/common"
)

// UptimeInfo provides system uptime information
type UptimeInfo struct {
	SystemInfo
}

// GetInfo returns the system uptime
func (u *UptimeInfo) GetInfo() string {
	out, err := common.GlobalCommandExecutor.Execute("cat", "/proc/uptime")
	if err != nil {
		return "Unknown"
	}

	fields := strings.Fields(out)
	if len(fields) < 1 {
		return "Unknown"
	}

	uptime, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return "Unknown"
	}

	// Convert to seconds
	seconds := int(uptime)
	days := seconds / 86400
	seconds %= 86400
	hours := seconds / 3600
	seconds %= 3600
	minutes := seconds / 60

	if days > 0 {
		return fmt.Sprintf("%d days, %d hours, %d minutes", days, hours, minutes)
	} else if hours > 0 {
		return fmt.Sprintf("%d hours, %d minutes", hours, minutes)
	} else {
		return fmt.Sprintf("%d minutes", minutes)
	}
}
