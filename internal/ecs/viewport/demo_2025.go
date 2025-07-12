package viewport

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// Demo2025 showcases all 2025 responsive design features
type Demo2025 struct {
	viewport           *AdvancedViewportManager
	fluidGrid          *FluidGrid
	responsiveHUD      *GameHUD
	responsiveRenderer *ResponsiveRenderer
	accessibility      *AccessibilityConfig

	// Demo state
	demoTime    float64
	currentDemo int
	demoCount   int
	font        text.Face
}

// NewDemo2025 creates a new 2025 responsive design demo
func NewDemo2025(font text.Face) *Demo2025 {
	return &Demo2025{
		viewport:           NewAdvancedViewportManager(),
		fluidGrid:          NewFluidGrid(),
		responsiveHUD:      NewGameHUD(),
		responsiveRenderer: NewResponsiveRenderer(),
		accessibility:      NewAccessibilityConfig(),
		font:               font,
		demoCount:          4,
	}
}

// Update updates the demo
func (d *Demo2025) Update(deltaTime float64) {
	d.demoTime += deltaTime

	// Cycle through demos every 5 seconds
	if int(d.demoTime/5.0) != d.currentDemo {
		d.currentDemo = int(d.demoTime/5.0) % d.demoCount
	}

	// Update HUD animations
	d.responsiveHUD.Update(deltaTime)
}

// Draw renders the demo
func (d *Demo2025) Draw(screen *ebiten.Image) {
	// Clear screen with dark background
	screen.Fill(color.RGBA{10, 10, 20, 255})

	// Update viewport with current screen dimensions
	width, height := screen.Bounds().Dx(), screen.Bounds().Dy()
	d.viewport.UpdateAdvanced(width, height)

	// Update fluid grid
	d.fluidGrid.UpdateContainer(float64(width), float64(height))

	// Draw demo content based on current demo
	switch d.currentDemo {
	case 0:
		d.drawDeviceClassificationDemo(screen)
	case 1:
		d.drawFluidGridDemo(screen)
	case 2:
		d.drawSciFiHUDDemo(screen)
	case 3:
		d.drawAccessibilityDemo(screen)
	}

	// Draw demo info
	d.drawDemoInfo(screen)
}

// drawDeviceClassificationDemo demonstrates device classification
func (d *Demo2025) drawDeviceClassificationDemo(screen *ebiten.Image) {
	width, height := d.viewport.GetCurrentDimensions()
	deviceClass := d.viewport.GetDeviceClass()
	orientation := d.viewport.GetOrientation()
	scale := d.viewport.GetIntrinsicScale()

	// Draw device info
	info := fmt.Sprintf("Device: %s", deviceClass)
	op := &text.DrawOptions{}
	op.GeoM.Translate(50, 100)
	text.Draw(screen, info, d.font, op)

	info = fmt.Sprintf("Orientation: %s", orientation)
	op = &text.DrawOptions{}
	op.GeoM.Translate(50, 130)
	text.Draw(screen, info, d.font, op)

	info = fmt.Sprintf("Scale: %.2f", scale)
	op = &text.DrawOptions{}
	op.GeoM.Translate(50, 160)
	text.Draw(screen, info, d.font, op)

	info = fmt.Sprintf("Resolution: %dx%d", width, height)
	op = &text.DrawOptions{}
	op.GeoM.Translate(50, 190)
	text.Draw(screen, info, d.font, op)

	// Draw device-specific visual indicator
	d.drawDeviceIndicator(screen, deviceClass)
}

// drawFluidGridDemo demonstrates fluid grid system
func (d *Demo2025) drawFluidGridDemo(screen *ebiten.Image) {
	breakpoint := d.fluidGrid.GetBreakpoint()
	columns := d.fluidGrid.GetColumnsForBreakpoint(breakpoint)
	gutterSize := d.fluidGrid.GetGutterSize()

	// Draw grid info
	info := fmt.Sprintf("Breakpoint: %s", breakpoint)
	op := &text.DrawOptions{}
	op.GeoM.Translate(50, 100)
	text.Draw(screen, info, d.font, op)

	info = fmt.Sprintf("Columns: %d", columns)
	op = &text.DrawOptions{}
	op.GeoM.Translate(50, 130)
	text.Draw(screen, info, d.font, op)

	info = fmt.Sprintf("Gutter: %.1fpx", gutterSize)
	op = &text.DrawOptions{}
	op.GeoM.Translate(50, 160)
	text.Draw(screen, info, d.font, op)

	// Draw visual grid
	d.drawVisualGrid(screen, columns, gutterSize)
}

// drawSciFiHUDDemo demonstrates sci-fi HUD effects
func (d *Demo2025) drawSciFiHUDDemo(screen *ebiten.Image) {
	// Create demo HUD elements
	elements := []UIElement{
		{
			Position:    Position2D{X: 50, Y: 100},
			Size:        Position2D{X: 300, Y: 30},
			Layer:       1,
			ElementType: "health_bar",
			Content:     0.75,
		},
		{
			Position:    Position2D{X: 50, Y: 150},
			Size:        Position2D{X: 200, Y: 25},
			Layer:       2,
			ElementType: "ammo_counter",
			Content:     8,
		},
		{
			Position:    Position2D{X: 50, Y: 200},
			Size:        Position2D{X: 250, Y: 40},
			Layer:       1,
			ElementType: "score_display",
			Content:     "Score: 12,450",
		},
		{
			Position:    Position2D{X: 50, Y: 260},
			Size:        Position2D{X: 30, Y: 30},
			Layer:       3,
			ElementType: "status_indicator",
			Content:     "ready",
		},
	}

	// Render all elements with HUD effects
	for _, element := range elements {
		d.responsiveHUD.RenderWithDepth(screen, element, d.viewport)
	}
}

// drawAccessibilityDemo demonstrates accessibility features
func (d *Demo2025) drawAccessibilityDemo(screen *ebiten.Image) {
	// Calculate safe areas
	top, right, bottom, left := d.accessibility.CalculateSafeArea(d.viewport)

	// Draw safe area indicators
	safeAreaColor := color.RGBA{255, 255, 0, 100} // Semi-transparent yellow

	// Top safe area
	ebitenutil.DrawRect(screen, 0, 0, float64(screen.Bounds().Dx()), top, safeAreaColor)

	// Bottom safe area
	ebitenutil.DrawRect(screen, 0, float64(screen.Bounds().Dy())-bottom,
		float64(screen.Bounds().Dx()), bottom, safeAreaColor)

	// Left safe area
	ebitenutil.DrawRect(screen, 0, 0, left, float64(screen.Bounds().Dy()), safeAreaColor)

	// Right safe area
	ebitenutil.DrawRect(screen, float64(screen.Bounds().Dx())-right, 0,
		right, float64(screen.Bounds().Dy()), safeAreaColor)

	// Draw accessibility info
	info := fmt.Sprintf("Safe Areas - T:%.0f R:%.0f B:%.0f L:%.0f", top, right, bottom, left)
	op := &text.DrawOptions{}
	op.GeoM.Translate(50, 100)
	text.Draw(screen, info, d.font, op)

	// Draw touch target examples
	d.drawTouchTargetExamples(screen)
}

// drawDemoInfo draws demo navigation info
func (d *Demo2025) drawDemoInfo(screen *ebiten.Image) {
	width := screen.Bounds().Dx()

	// Draw demo title
	title := "2025 Responsive Design Demo"
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(width)/2-100, 30)
	text.Draw(screen, title, d.font, op)

	// Draw demo progress
	progress := fmt.Sprintf("Demo %d/%d", d.currentDemo+1, d.demoCount)
	op = &text.DrawOptions{}
	op.GeoM.Translate(float64(width)/2-30, 60)
	text.Draw(screen, progress, d.font, op)

	// Draw performance metrics
	metrics := d.responsiveRenderer.GetPerformanceMetrics()
	metricsText := fmt.Sprintf("Cache: %.1f%% (%d hits, %d misses)",
		metrics["cache_ratio"].(float64)*100,
		metrics["cache_hits"].(int),
		metrics["cache_misses"].(int))

	op = &text.DrawOptions{}
	op.GeoM.Translate(50, float64(screen.Bounds().Dy())-50)
	text.Draw(screen, metricsText, d.font, op)
}

// drawDeviceIndicator draws a visual indicator for the current device class
func (d *Demo2025) drawDeviceIndicator(screen *ebiten.Image, deviceClass string) {
	width, height := screen.Bounds().Dx(), screen.Bounds().Dy()

	var indicatorColor color.RGBA
	var size float64

	switch deviceClass {
	case string(DeviceClassMobile):
		indicatorColor = color.RGBA{255, 100, 100, 255} // Red
		size = 50
	case string(DeviceClassTablet):
		indicatorColor = color.RGBA{100, 255, 100, 255} // Green
		size = 80
	case string(DeviceClassUltrawide):
		indicatorColor = color.RGBA{100, 100, 255, 255} // Blue
		size = 120
	default:
		indicatorColor = color.RGBA{255, 255, 100, 255} // Yellow
		size = 100
	}

	// Draw device indicator in center
	x := float64(width)/2 - size/2
	y := float64(height)/2 - size/2
	ebitenutil.DrawCircle(screen, x+size/2, y+size/2, size/2, indicatorColor)
}

// drawVisualGrid draws a visual representation of the fluid grid
func (d *Demo2025) drawVisualGrid(screen *ebiten.Image, columns int, gutterSize float64) {
	width, height := screen.Bounds().Dx(), screen.Bounds().Dy()
	startY := float64(height) * 0.4

	// Calculate column width
	availableWidth := float64(width) - 100 // 50px margin on each side
	gutterSpace := gutterSize * float64(columns-1)
	columnWidth := (availableWidth - gutterSpace) / float64(columns)

	// Draw grid lines
	gridColor := color.RGBA{100, 100, 255, 100}

	for i := 0; i <= columns; i++ {
		x := 50 + float64(i)*(columnWidth+gutterSize)
		ebitenutil.DrawLine(screen, x, startY, x, startY+200, gridColor)
	}

	// Draw column labels
	for i := 0; i < columns; i++ {
		x := 50 + float64(i)*(columnWidth+gutterSize) + columnWidth/2
		label := fmt.Sprintf("%d", i+1)
		op := &text.DrawOptions{}
		op.GeoM.Translate(x-10, startY+220)
		text.Draw(screen, label, d.font, op)
	}
}

// drawTouchTargetExamples draws examples of accessible touch targets
func (d *Demo2025) drawTouchTargetExamples(screen *ebiten.Image) {
	// Draw small touch target (inaccessible)
	smallTarget := 30.0
	ebitenutil.DrawRect(screen, 50, 150, smallTarget, smallTarget, color.RGBA{255, 100, 100, 255})
	op := &text.DrawOptions{}
	op.GeoM.Translate(50, 150+smallTarget+20)
	text.Draw(screen, "Too Small", d.font, op)

	// Draw proper touch target (accessible)
	properTarget := d.accessibility.EnsureTouchTarget(30.0, d.viewport)
	ebitenutil.DrawRect(screen, 200, 150, properTarget, properTarget, color.RGBA{100, 255, 100, 255})
	op = &text.DrawOptions{}
	op.GeoM.Translate(200, 150+properTarget+20)
	text.Draw(screen, "Accessible", d.font, op)

	// Draw size comparison
	comparison := fmt.Sprintf("Small: %.0fpx → Accessible: %.0fpx", smallTarget, properTarget)
	op = &text.DrawOptions{}
	op.GeoM.Translate(50, 250)
	text.Draw(screen, comparison, d.font, op)
}

// GetViewport returns the viewport manager for external access
func (d *Demo2025) GetViewport() *AdvancedViewportManager {
	return d.viewport
}

// GetFluidGrid returns the fluid grid for external access
func (d *Demo2025) GetFluidGrid() *FluidGrid {
	return d.fluidGrid
}

// GetResponsiveHUD returns the responsive HUD for external access
func (d *Demo2025) GetResponsiveHUD() *GameHUD {
	return d.responsiveHUD
}
