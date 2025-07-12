package viewport

import (
	"math"
)

// Position2D represents a 2D coordinate position
type Position2D struct {
	X, Y float64
}

// ResponsiveCoordinate implements 2025 device-class aware positioning
type ResponsiveCoordinate struct {
	Mobile    Position2D // < 768px
	Tablet    Position2D // 768-1024px
	Desktop   Position2D // 1024-1440px
	Ultrawide Position2D // > 1440px
	Adaptive  bool       // Use AI-like content adaptation
}

// AccessibilityConfig implements 2025 universal design principles
type AccessibilityConfig struct {
	minTouchTarget  float64 // 44pt minimum (iOS) / 48dp (Android)
	contrastRatio   float64 // WCAG 2.2 AA compliance
	motionReduction bool    // Respect prefers-reduced-motion
	focusVisible    bool    // Enhanced focus indicators
	textScaling     bool    // Support system text scaling
}

// NewResponsiveCoordinate creates a new responsive coordinate
func NewResponsiveCoordinate(mobile, tablet, desktop, ultrawide Position2D) *ResponsiveCoordinate {
	return &ResponsiveCoordinate{
		Mobile:    mobile,
		Tablet:    tablet,
		Desktop:   desktop,
		Ultrawide: ultrawide,
		Adaptive:  true,
	}
}

// Resolve resolves the coordinate for the current device class
func (rc ResponsiveCoordinate) Resolve(viewport *AdvancedViewportManager) Position2D {
	switch viewport.GetDeviceClass() {
	case string(DeviceClassMobile):
		if rc.Adaptive {
			// 2025: Smart positioning for touch interfaces
			return rc.adjustForTouch(rc.Mobile, viewport)
		}
		return rc.Mobile
	case string(DeviceClassTablet):
		return rc.Tablet
	case string(DeviceClassUltrawide):
		return rc.Ultrawide
	default:
		return rc.Desktop
	}
}

// adjustForTouch implements 2025 touch-optimized positioning
func (rc ResponsiveCoordinate) adjustForTouch(base Position2D, viewport *AdvancedViewportManager) Position2D {
	// Move UI elements away from edges on mobile for thumb reach
	safeX := math.Max(base.X, 0.05) // 5% margin minimum
	safeY := math.Max(base.Y, 0.1)  // 10% margin minimum for status bar

	// Apply device-specific adjustments
	if viewport.GetOrientation() == string(OrientationPortrait) {
		// Portrait: Adjust for one-handed use
		safeX = math.Min(safeX, 0.85) // Keep within thumb reach
	}

	return Position2D{X: safeX, Y: safeY}
}

// NewAccessibilityConfig creates a new accessibility configuration
func NewAccessibilityConfig() *AccessibilityConfig {
	return &AccessibilityConfig{
		minTouchTarget:  48.0, // Android standard
		contrastRatio:   4.5,  // WCAG 2.2 AA
		motionReduction: false,
		focusVisible:    true,
		textScaling:     true,
	}
}

// EnsureTouchTarget ensures minimum touch target size for accessibility
func (ac *AccessibilityConfig) EnsureTouchTarget(size float64, viewport *AdvancedViewportManager) float64 {
	scaledMin := ac.minTouchTarget * viewport.GetIntrinsicScale()

	if viewport.GetDeviceClass() == string(DeviceClassMobile) ||
		viewport.GetDeviceClass() == string(DeviceClassTablet) {
		return math.Max(size, scaledMin)
	}
	return size // Desktop doesn't need enlarged touch targets
}

// CalculateSafeArea calculates the safe area for UI elements
func (ac *AccessibilityConfig) CalculateSafeArea(viewport *AdvancedViewportManager) (top, right, bottom, left float64) {
	width, height := viewport.GetCurrentDimensions()

	switch viewport.GetDeviceClass() {
	case string(DeviceClassMobile):
		// Mobile: Account for status bar, notch, home indicator
		top = 44.0 * viewport.GetIntrinsicScale()    // Status bar
		bottom = 34.0 * viewport.GetIntrinsicScale() // Home indicator
		left = 16.0 * viewport.GetIntrinsicScale()   // Side margin
		right = 16.0 * viewport.GetIntrinsicScale()  // Side margin
	case string(DeviceClassTablet):
		// Tablet: Minimal safe areas
		top = 20.0 * viewport.GetIntrinsicScale()
		bottom = 20.0 * viewport.GetIntrinsicScale()
		left = 24.0 * viewport.GetIntrinsicScale()
		right = 24.0 * viewport.GetIntrinsicScale()
	default:
		// Desktop: No special safe areas needed
		top = 0
		bottom = 0
		left = 0
		right = 0
	}

	// Ensure safe areas don't exceed screen bounds
	top = math.Min(top, float64(height)*0.1)
	bottom = math.Min(bottom, float64(height)*0.1)
	left = math.Min(left, float64(width)*0.05)
	right = math.Min(right, float64(width)*0.05)

	return top, right, bottom, left
}

// Position2D helpers for common UI locations
func TopLeft(x, y float64) Position2D {
	return Position2D{X: x, Y: y}
}

func TopRight(screenWidth, x, y float64) Position2D {
	return Position2D{X: screenWidth - x, Y: y}
}

func BottomLeft(x, screenHeight, y float64) Position2D {
	return Position2D{X: x, Y: screenHeight - y}
}

func BottomRight(screenWidth, screenHeight, x, y float64) Position2D {
	return Position2D{X: screenWidth - x, Y: screenHeight - y}
}

func Center(screenWidth, screenHeight float64) Position2D {
	return Position2D{X: screenWidth / 2, Y: screenHeight / 2}
}

// CreateResponsivePosition creates a responsive position that adapts to device class
func CreateResponsivePosition(mobileX, mobileY, tabletX, tabletY, desktopX, desktopY, ultrawideX, ultrawideY float64) *ResponsiveCoordinate {
	return &ResponsiveCoordinate{
		Mobile:    Position2D{X: mobileX, Y: mobileY},
		Tablet:    Position2D{X: tabletX, Y: tabletY},
		Desktop:   Position2D{X: desktopX, Y: desktopY},
		Ultrawide: Position2D{X: ultrawideX, Y: ultrawideY},
		Adaptive:  true,
	}
}

// ScalePosition scales a position by the viewport's intrinsic scale
func ScalePosition(pos Position2D, viewport *AdvancedViewportManager) Position2D {
	scale := viewport.GetIntrinsicScale()
	return Position2D{
		X: pos.X * scale,
		Y: pos.Y * scale,
	}
}

// OffsetPosition offsets a position by the given amount
func OffsetPosition(pos Position2D, offsetX, offsetY float64) Position2D {
	return Position2D{
		X: pos.X + offsetX,
		Y: pos.Y + offsetY,
	}
}
