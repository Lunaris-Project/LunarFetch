package components

import (
	"fmt"
	"os"
	"strings"

	"lunarfetch/src/common"
)

// BatteryInfo provides battery status information
type BatteryInfo struct {
	SystemInfo
}

// GetInfo returns the battery status
func (b *BatteryInfo) GetInfo() string {
	if _, err := os.Stat("/sys/class/power_supply/BAT0"); err != nil {
		return "No battery"
	}

	capacity, err := os.ReadFile("/sys/class/power_supply/BAT0/capacity")
	if err != nil {
		// Try using command execution as fallback
		out, err := common.GlobalCommandExecutor.Execute("cat", "/sys/class/power_supply/BAT0/capacity")
		if err != nil {
			return "Unknown"
		}
		capacity = []byte(out)
	}

	status, err := os.ReadFile("/sys/class/power_supply/BAT0/status")
	if err != nil {
		// Try using command execution as fallback
		out, err := common.GlobalCommandExecutor.Execute("cat", "/sys/class/power_supply/BAT0/status")
		if err != nil {
			return fmt.Sprintf("%s%%", strings.TrimSpace(string(capacity)))
		}
		status = []byte(out)
	}

	return fmt.Sprintf("%s%% (%s)", strings.TrimSpace(string(capacity)), strings.TrimSpace(string(status)))
}
