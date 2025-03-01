package components

import (
	"os"
	"path/filepath"
	"strings"

	"lunarfetch/src/common"
)

// ShellInfo provides shell information
type ShellInfo struct {
	SystemInfo
}

// GetInfo returns the shell name
func (s *ShellInfo) GetInfo() string {
	shell := os.Getenv("SHELL")
	if shell == "" {
		out, err := common.GlobalCommandExecutor.Execute("getent", "passwd", os.Getenv("USER"))
		if err == nil {
			fields := strings.Split(out, ":")
			if len(fields) >= 7 {
				shell = fields[6]
			}
		}
	}
	return filepath.Base(shell)
}
