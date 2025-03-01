package components

import (
	"fmt"
	"strconv"
	"strings"

	"lunarfetch/src/common"
)

// DiskInfo provides disk usage information
type DiskInfo struct {
	SystemInfo
}

// GetInfo returns the disk usage
func (d *DiskInfo) GetInfo() string {
	out, err := common.GlobalCommandExecutor.Execute("df", "-B1")
	if err != nil {
		return "Unknown"
	}
	lines := strings.Split(out, "\n")
	var totalUsed uint64
	var totalSize uint64
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) > 3 && fields[0] != "Filesystem" {
			used, _ := strconv.ParseUint(fields[2], 10, 64)
			size, _ := strconv.ParseUint(fields[1], 10, 64)
			totalUsed += used
			totalSize += size
		}
	}
	return fmt.Sprintf("%s / %s", FormatBytes(totalUsed), FormatBytes(totalSize))
}
