package components

import (
	"os"
	"strings"

	"lunarfetch/src/common"
)

// TerminalInfo provides terminal information
type TerminalInfo struct {
	SystemInfo
}

// GetInfo returns the terminal name
func (t *TerminalInfo) GetInfo() string {
	term := os.Getenv("TERM")
	if term == "" {
		out, err := common.GlobalCommandExecutor.Execute("ps", "-p", os.Getenv("$"), "-o", "args=")
		if err == nil && out != "" {
			term = strings.TrimSpace(out)
		} else {
			term = "Unknown"
		}
	}
	return term
}
