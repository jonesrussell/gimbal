package ui

import (
	"github.com/jonesrussell/gimbal/internal/ui/core"
)

// Layout handles responsive layout calculations
type Layout struct {
	screenWidth  int
	screenHeight int
	deviceClass  string
}

// NewLayout creates a new layout manager
func NewLayout() *Layout {
	return &Layout{}
}

// Update updates the layout based on screen dimensions
func (l *Layout) Update(width, height int) {
	l.screenWidth = width
	l.screenHeight = height
	l.deviceClass = l.calculateDeviceClass(width)
}

// GetDeviceClass returns the current device class
func (l *Layout) GetDeviceClass() string {
	return l.deviceClass
}

// GetDimensions returns the current screen dimensions
func (l *Layout) GetDimensions() (width, height int) {
	return l.screenWidth, l.screenHeight
}

// calculateDeviceClass determines the device class based on width
func (l *Layout) calculateDeviceClass(width int) string {
	switch {
	case width < core.MobileBreakpoint:
		return core.DeviceMobile
	case width < core.TabletBreakpoint:
		return core.DeviceTablet
	case width > core.UltrawideBreakpoint:
		return core.DeviceUltrawide
	default:
		return core.DeviceDesktop
	}
}
