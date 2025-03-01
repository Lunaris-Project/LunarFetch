package components

import (
	"os"
	"strings"

	"lunarfetch/src/common"
)

// UserInfo provides user information
type UserInfo struct {
	SystemInfo
}

// GetInfo returns the current user
func (u *UserInfo) GetInfo() string {
	user := os.Getenv("USER")
	if user == "" {
		out, err := common.GlobalCommandExecutor.Execute("whoami")
		if err == nil {
			user = strings.TrimSpace(out)
		} else {
			user = "Unknown"
		}
	}
	return user
}
