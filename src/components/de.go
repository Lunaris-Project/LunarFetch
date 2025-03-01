package components

import (
	"lunarfetch/src/common"
)

// DEInfo provides desktop environment information
type DEInfo struct {
	SystemInfo
}

// GetInfo returns the desktop environment
func (d *DEInfo) GetInfo() string {
	de := common.GetEnv("XDG_CURRENT_DESKTOP", "")
	if de != "" {
		return de
	}

	de = common.GetEnv("DESKTOP_SESSION", "")
	if de != "" {
		return de
	}

	return "Unknown"
}
