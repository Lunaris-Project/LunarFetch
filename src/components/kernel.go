package components

import (
	"strings"

	"lunarfetch/src/common"
)

// KernelInfo provides kernel version information
type KernelInfo struct {
	SystemInfo
}

// GetInfo returns the kernel version
func (k *KernelInfo) GetInfo() string {
	out, err := common.GlobalCommandExecutor.Execute("uname", "-r")
	if err != nil {
		return "Unknown"
	}
	return strings.TrimSpace(out)
}
