package viewport

import (
	"math"
)

// AdvancedViewportManager implements 2025 responsive design techniques
type AdvancedViewportManager struct {
	// Base design resolution (1920x1080)
	baseWidth  int
	baseHeight int

	// Current viewport dimensions
	currentWidth  int
	currentHeight int

	// Device characteristics
	pixelRatio  float64 // Device pixel ratio support
	orientation string  // "portrait" | "landscape"
	deviceClass string  // "mobile" | "tablet" | "desktop" | "ultrawide"

	// 2025: Content-aware intrinsic scaling
	intrinsicScale float64

	// Performance optimization
	lastWidth       int
	lastHeight      int
	lastOrientation string
	needsRelayout   bool
}

// DeviceClass represents different device categories
type DeviceClass string

const (
	DeviceClassMobile    DeviceClass = "mobile"
	DeviceClassTablet    DeviceClass = "tablet"
	DeviceClassDesktop   DeviceClass = "desktop"
	DeviceClassUltrawide DeviceClass = "ultrawide"
)

// Orientation represents screen orientation
type Orientation string

const (
	OrientationPortrait  Orientation = "portrait"
	OrientationLandscape Orientation = "landscape"
)

// NewAdvancedViewportManager creates a new viewport manager with 2025 techniques
func NewAdvancedViewportManager() *AdvancedViewportManager {
	return &AdvancedViewportManager{
		baseWidth:      1920,
		baseHeight:     1080,
		pixelRatio:     1.0,
		orientation:    string(OrientationLandscape),
		deviceClass:    string(DeviceClassDesktop),
		intrinsicScale: 1.0,
		needsRelayout:  true,
	}
}

// UpdateAdvanced updates the viewport with 2025 responsive techniques
func (vm *AdvancedViewportManager) UpdateAdvanced(width, height int) {
	vm.currentWidth = width
	vm.currentHeight = height

	// Detect orientation
	if width > height {
		vm.orientation = string(OrientationLandscape)
	} else {
		vm.orientation = string(OrientationPortrait)
	}

	// 2025: Advanced device classification
	vm.classifyDevice(width, height)

	// 2025: Calculate intrinsic scaling (content-aware)
	vm.calculateIntrinsicScale()

	// Performance optimization: only relayout on significant changes
	vm.needsRelayout = vm.shouldRelayout()

	if vm.needsRelayout {
		vm.lastWidth = width
		vm.lastHeight = height
		vm.lastOrientation = vm.orientation
	}
}

// classifyDevice implements 2025 device classification
func (vm *AdvancedViewportManager) classifyDevice(width, height int) {
	// 2025: More sophisticated device classification
	diagonal := math.Sqrt(float64(width*width + height*height))

	// Account for pixel density and physical size
	effectiveDiagonal := diagonal / vm.pixelRatio

	if effectiveDiagonal < 7.0 {
		// Small mobile devices
		vm.deviceClass = string(DeviceClassMobile)
	} else if effectiveDiagonal < 12.0 {
		// Tablets and large phones
		vm.deviceClass = string(DeviceClassTablet)
	} else if width > 2560 {
		// Ultrawide displays
		vm.deviceClass = string(DeviceClassUltrawide)
	} else {
		// Standard desktop displays
		vm.deviceClass = string(DeviceClassDesktop)
	}
}

// calculateIntrinsicScale implements 2025 content-aware scaling
func (vm *AdvancedViewportManager) calculateIntrinsicScale() {
	aspectRatio := float64(vm.currentWidth) / float64(vm.currentHeight)
	baseAspectRatio := float64(vm.baseWidth) / float64(vm.baseHeight)

	// 2025: Content-aware scaling based on available space
	if aspectRatio > baseAspectRatio {
		// Wide screen: scale by height, center horizontally
		vm.intrinsicScale = float64(vm.currentHeight) / float64(vm.baseHeight)
	} else {
		// Tall screen: scale by width, center vertically
		vm.intrinsicScale = float64(vm.currentWidth) / float64(vm.baseWidth)
	}

	// Apply device-specific scaling adjustments
	vm.applyDeviceSpecificScaling()
}

// applyDeviceSpecificScaling applies 2025 device-specific optimizations
func (vm *AdvancedViewportManager) applyDeviceSpecificScaling() {
	switch vm.deviceClass {
	case string(DeviceClassMobile):
		// Mobile: Slightly larger scaling for touch targets
		vm.intrinsicScale *= 1.1
	case string(DeviceClassTablet):
		// Tablet: Balanced scaling
		vm.intrinsicScale *= 1.05
	case string(DeviceClassUltrawide):
		// Ultrawide: Slightly smaller to avoid UI stretching
		vm.intrinsicScale *= 0.95
	}
}

// shouldRelayout implements 2025 performance optimization
func (vm *AdvancedViewportManager) shouldRelayout() bool {
	if vm.lastWidth == 0 || vm.lastHeight == 0 {
		return true
	}

	// Threshold-based re-layout to prevent jitter
	widthChange := math.Abs(float64(vm.currentWidth - vm.lastWidth))
	heightChange := math.Abs(float64(vm.currentHeight - vm.lastHeight))

	threshold := 50.0 // 50px threshold
	return widthChange > threshold || heightChange > threshold ||
		vm.orientation != vm.lastOrientation
}

// GetLogicalScreenSize returns the logical screen size for the current device class
func (vm *AdvancedViewportManager) GetLogicalScreenSize() (width, height int) {
	switch vm.deviceClass {
	case string(DeviceClassMobile):
		if vm.orientation == string(OrientationPortrait) {
			return 1080, 1920 // Mobile portrait
		}
		return 1920, 1080 // Mobile landscape
	case string(DeviceClassTablet):
		return 1440, 1080 // Tablet optimized
	case string(DeviceClassUltrawide):
		return 2560, 1080 // Ultrawide support
	default:
		return 1920, 1080 // Standard desktop
	}
}

// GetDeviceClass returns the current device class
func (vm *AdvancedViewportManager) GetDeviceClass() string {
	return vm.deviceClass
}

// GetOrientation returns the current orientation
func (vm *AdvancedViewportManager) GetOrientation() string {
	return vm.orientation
}

// GetIntrinsicScale returns the current intrinsic scale
func (vm *AdvancedViewportManager) GetIntrinsicScale() float64 {
	return vm.intrinsicScale
}

// NeedsRelayout returns whether a relayout is needed
func (vm *AdvancedViewportManager) NeedsRelayout() bool {
	return vm.needsRelayout
}

// GetBaseDimensions returns the base design dimensions
func (vm *AdvancedViewportManager) GetBaseDimensions() (width, height int) {
	return vm.baseWidth, vm.baseHeight
}

// GetCurrentDimensions returns the current viewport dimensions
func (vm *AdvancedViewportManager) GetCurrentDimensions() (width, height int) {
	return vm.currentWidth, vm.currentHeight
}
