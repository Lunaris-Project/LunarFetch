package components

import (
	"strings"

	"lunarfetch/src/common"
)

// WMInfo provides window manager information
type WMInfo struct {
	SystemInfo
}

// GetInfo returns the window manager
func (w *WMInfo) GetInfo() string {
	// Try to get from environment variable
	if wm := strings.TrimSpace(strings.ToLower(strings.Join([]string{
		strings.TrimSpace(strings.ToLower(common.GetEnv("XDG_CURRENT_DESKTOP", ""))),
		strings.TrimSpace(strings.ToLower(common.GetEnv("DESKTOP_SESSION", ""))),
		strings.TrimSpace(strings.ToLower(common.GetEnv("GDMSESSION", ""))),
		strings.TrimSpace(strings.ToLower(common.GetEnv("XDG_SESSION_DESKTOP", ""))),
	}, " "))); wm != "" && wm != " " {
		return strings.TrimSpace(wm)
	}

	// Try to get from wmctrl
	out, err := common.GlobalCommandExecutor.Execute("wmctrl", "-m")
	if err == nil {
		lines := strings.Split(out, "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "Name:") {
				return strings.TrimSpace(strings.TrimPrefix(line, "Name:"))
			}
		}
	}

	// Try to get from ps
	out, err = common.GlobalCommandExecutor.Execute("ps", "-e")
	if err == nil {
		lines := strings.Split(out, "\n")
		for _, line := range lines {
			for _, wm := range []string{
				"2bwm", "9wm", "awesome", "beryl", "blackbox", "bspwm", "budgie-wm", "cinnamon", "compiz", "deepin-wm",
				"dtwm", "dwm", "enlightenment", "fluxbox", "fvwm", "herbstluftwm", "i3", "icewm", "jwm", "kwin",
				"metacity", "monsterwm", "musca", "mutter", "mwm", "openbox", "pekwm", "qtile", "ratpoison", "sawfish",
				"scrotwm", "spectrwm", "stumpwm", "subtle", "sway", "wmaker", "wmfs", "wmii", "xfwm4", "xmonad",
			} {
				if strings.Contains(line, wm) {
					return wm
				}
			}
		}
	}

	return "Unknown"
}
