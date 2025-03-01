package components

import (
	"strings"

	"lunarfetch/src/common"
)

// ThemeInfo provides theme information
type ThemeInfo struct {
	SystemInfo
}

// GetInfo returns the current theme
func (t *ThemeInfo) GetInfo() string {
	out, err := common.GlobalCommandExecutor.Execute("gsettings", "get", "org.gnome.desktop.interface", "gtk-theme")
	if err == nil {
		theme := strings.TrimSpace(out)
		theme = strings.Trim(theme, "'")
		return theme
	}

	out, err = common.GlobalCommandExecutor.Execute("dconf", "read", "/org/gnome/desktop/interface/gtk-theme")
	if err == nil {
		theme := strings.TrimSpace(out)
		theme = strings.Trim(theme, "'")
		return theme
	}

	out, err = common.GlobalCommandExecutor.Execute("grep", "gtk-theme-name", "~/.gtkrc-2.0")
	if err == nil {
		parts := strings.Split(out, "=")
		if len(parts) > 1 {
			theme := strings.TrimSpace(parts[1])
			theme = strings.Trim(theme, "\"")
			return theme
		}
	}

	return "Unknown"
}

// WMThemeInfo provides window manager theme information
type WMThemeInfo struct {
	SystemInfo
}

// GetInfo returns the window manager theme
func (w *WMThemeInfo) GetInfo() string {
	out, err := common.GlobalCommandExecutor.Execute("gsettings", "get", "org.gnome.desktop.wm.preferences", "theme")
	if err != nil {
		return "Unknown"
	}
	theme := strings.TrimSpace(out)
	theme = strings.Trim(theme, "'")
	return theme
}

// IconsInfo provides icon theme information
type IconsInfo struct {
	SystemInfo
}

// GetInfo returns the icon theme
func (i *IconsInfo) GetInfo() string {
	out, err := common.GlobalCommandExecutor.Execute("gsettings", "get", "org.gnome.desktop.interface", "icon-theme")
	if err != nil {
		return "Unknown"
	}
	theme := strings.TrimSpace(out)
	theme = strings.Trim(theme, "'")
	return theme
}
