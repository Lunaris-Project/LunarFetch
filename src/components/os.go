package components

import (
	"strings"

	"lunarfetch/src/common"
)

// OSInfo provides operating system information
type OSInfo struct {
	SystemInfo
}

// GetInfo returns the operating system information
func (o *OSInfo) GetInfo() string {
	out, err := common.GlobalCommandExecutor.Execute("lsb_release", "-si")
	if err != nil {
		out, err = common.GlobalCommandExecutor.Execute("grep", "^NAME=", "/etc/os-release")
		if err != nil {
			return "Unknown"
		}
		return strings.Trim(strings.TrimPrefix(strings.TrimSpace(out), "NAME=\""), "\"")
	}
	return strings.TrimSpace(out)
}
