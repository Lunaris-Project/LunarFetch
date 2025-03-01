package components

import (
	"encoding/json"
	"fmt"
	"strings"

	"lunarfetch/src/common"
)

// ResolutionInfo provides screen resolution information
type ResolutionInfo struct {
	SystemInfo
}

// GetInfo returns the screen resolution
func (r *ResolutionInfo) GetInfo() string {
	out, err := common.GlobalCommandExecutor.Execute("xrandr")
	if err == nil {
		lines := strings.Split(out, "\n")
		for _, line := range lines {
			if strings.Contains(line, " connected") {
				fields := strings.Fields(line)
				for _, field := range fields {
					if strings.Contains(field, "x") {
						return field
					}
				}
			}
		}
	}

	if out, err := common.GlobalCommandExecutor.Execute("swaymsg", "-t", "get_outputs"); err == nil {
		var outputs []struct {
			CurrentMode struct {
				Width  int `json:"width"`
				Height int `json:"height"`
			} `json:"current_mode"`
		}
		if err := json.Unmarshal([]byte(out), &outputs); err == nil && len(outputs) > 0 {
			return fmt.Sprintf("%dx%d", outputs[0].CurrentMode.Width, outputs[0].CurrentMode.Height)
		}
	}

	if out, err := common.GlobalCommandExecutor.Execute("wlr-randr"); err == nil {
		lines := strings.Split(out, "\n")
		for _, line := range lines {
			if strings.Contains(line, "current") {
				fields := strings.Fields(line)
				for _, field := range fields {
					if strings.Contains(field, "x") {
						return field
					}
				}
			}
		}
	}

	if out, err := common.GlobalCommandExecutor.Execute("xdpyinfo"); err == nil {
		lines := strings.Split(out, "\n")
		for _, line := range lines {
			if strings.Contains(line, "dimensions:") {
				fields := strings.Fields(line)
				return fields[1]
			}
		}
	}

	return "Unknown"
}
