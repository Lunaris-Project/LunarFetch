package components

import (
	"fmt"
)

// InfoProvider is an interface for components that provide system information
type InfoProvider interface {
	GetInfo() string
	GetName() string
}

// SystemInfo provides basic system information
type SystemInfo struct {
	Name string
}

// GetName returns the name of this info provider
func (s *SystemInfo) GetName() string {
	return s.Name
}

// FormatBytes formats bytes to a human-readable string
func FormatBytes(bytes uint64) string {
	units := []string{"B", "KB", "MB", "GB", "TB"}
	var unit string
	var value float64 = float64(bytes)

	for _, u := range units {
		if value < 1024 {
			unit = u
			break
		}
		value /= 1024
	}
	return fmt.Sprintf("%.2f %s", value, unit)
}
