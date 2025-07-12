package viewport

import (
	"math"
)

// FluidGrid implements 2025 advanced fluid grid with container-query concepts
type FluidGrid struct {
	columns    int     // 12-column system
	gutterBase float64 // Base gutter size
	maxWidth   float64 // Maximum container width
	minWidth   float64 // Minimum container width

	// 2025: Intrinsic sizing support
	contentAdaptive bool
	flowDirection   string // "row" | "column" | "adaptive"

	// Container query concepts
	containerWidth  float64
	containerHeight float64
	breakpoints     map[string]float64
}

// GridBreakpoint represents responsive breakpoints
type GridBreakpoint string

const (
	BreakpointMobile    GridBreakpoint = "mobile"    // < 768px
	BreakpointTablet    GridBreakpoint = "tablet"    // 768-1024px
	BreakpointDesktop   GridBreakpoint = "desktop"   // 1024-1440px
	BreakpointUltrawide GridBreakpoint = "ultrawide" // > 1440px
)

// NewFluidGrid creates a new fluid grid system
func NewFluidGrid() *FluidGrid {
	return &FluidGrid{
		columns:         12,
		gutterBase:      16.0,
		maxWidth:        1920.0,
		minWidth:        320.0,
		contentAdaptive: true,
		flowDirection:   "adaptive",
		breakpoints: map[string]float64{
			string(BreakpointMobile):    768.0,
			string(BreakpointTablet):    1024.0,
			string(BreakpointDesktop):   1440.0,
			string(BreakpointUltrawide): 2560.0,
		},
	}
}

// UpdateContainer updates the container dimensions for container-query behavior
func (fg *FluidGrid) UpdateContainer(width, height float64) {
	fg.containerWidth = math.Min(width, fg.maxWidth)
	fg.containerWidth = math.Max(fg.containerWidth, fg.minWidth)
	fg.containerHeight = height
}

// GetColumnWidth calculates column width with container-query concepts
func (fg *FluidGrid) GetColumnWidth(viewport *AdvancedViewportManager, span int) float64 {
	// Container-query-like behavior for game UI
	gutterSpace := fg.gutterBase * float64(fg.columns-1)
	availableSpace := fg.containerWidth - gutterSpace
	columnWidth := availableSpace / float64(fg.columns)

	// Apply device-specific adjustments
	adjustedWidth := fg.applyDeviceAdjustments(columnWidth, viewport)

	return adjustedWidth * float64(span)
}

// applyDeviceAdjustments applies 2025 device-specific grid adjustments
func (fg *FluidGrid) applyDeviceAdjustments(baseWidth float64, viewport *AdvancedViewportManager) float64 {
	switch viewport.GetDeviceClass() {
	case string(DeviceClassMobile):
		// Mobile: Slightly larger columns for touch targets
		return baseWidth * 1.1
	case string(DeviceClassTablet):
		// Tablet: Balanced sizing
		return baseWidth * 1.05
	case string(DeviceClassUltrawide):
		// Ultrawide: Slightly smaller to prevent stretching
		return baseWidth * 0.95
	default:
		return baseWidth
	}
}

// GetGutterSize returns the current gutter size based on container width
func (fg *FluidGrid) GetGutterSize() float64 {
	// 2025: Responsive gutters
	scale := fg.containerWidth / fg.maxWidth
	return fg.gutterBase * math.Max(0.5, math.Min(1.5, scale))
}

// GetBreakpoint returns the current breakpoint based on container width
func (fg *FluidGrid) GetBreakpoint() GridBreakpoint {
	if fg.containerWidth < fg.breakpoints[string(BreakpointMobile)] {
		return BreakpointMobile
	} else if fg.containerWidth < fg.breakpoints[string(BreakpointTablet)] {
		return BreakpointTablet
	} else if fg.containerWidth < fg.breakpoints[string(BreakpointDesktop)] {
		return BreakpointDesktop
	} else {
		return BreakpointUltrawide
	}
}

// GetColumnsForBreakpoint returns the number of columns for a specific breakpoint
func (fg *FluidGrid) GetColumnsForBreakpoint(breakpoint GridBreakpoint) int {
	switch breakpoint {
	case BreakpointMobile:
		return 4 // Mobile: 4 columns
	case BreakpointTablet:
		return 8 // Tablet: 8 columns
	case BreakpointDesktop:
		return 12 // Desktop: 12 columns
	case BreakpointUltrawide:
		return 16 // Ultrawide: 16 columns
	default:
		return 12
	}
}

// CalculatePosition calculates the position for a grid item
func (fg *FluidGrid) CalculatePosition(col, row int, span int) (x, y float64) {
	columnWidth := fg.GetColumnWidth(nil, 1) // Use nil viewport for base calculation
	gutterSize := fg.GetGutterSize()

	x = float64(col) * (columnWidth + gutterSize)
	y = float64(row) * (columnWidth + gutterSize)

	return x, y
}

// GetResponsiveSpan returns the appropriate span for different breakpoints
func (fg *FluidGrid) GetResponsiveSpan(mobile, tablet, desktop, ultrawide int) int {
	breakpoint := fg.GetBreakpoint()

	switch breakpoint {
	case BreakpointMobile:
		return mobile
	case BreakpointTablet:
		return tablet
	case BreakpointDesktop:
		return desktop
	case BreakpointUltrawide:
		return ultrawide
	default:
		return desktop
	}
}

// IsContainerQuery returns true if container-query behavior should be used
func (fg *FluidGrid) IsContainerQuery() bool {
	return fg.contentAdaptive && fg.containerWidth > 0
}

// GetContainerAspectRatio returns the container's aspect ratio
func (fg *FluidGrid) GetContainerAspectRatio() float64 {
	if fg.containerHeight == 0 {
		return 16.0 / 9.0 // Default aspect ratio
	}
	return fg.containerWidth / fg.containerHeight
}

// SetContentAdaptive enables or disables content-adaptive behavior
func (fg *FluidGrid) SetContentAdaptive(adaptive bool) {
	fg.contentAdaptive = adaptive
}

// SetFlowDirection sets the flow direction for the grid
func (fg *FluidGrid) SetFlowDirection(direction string) {
	fg.flowDirection = direction
}
