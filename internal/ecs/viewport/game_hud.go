package viewport

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// GameHUD implements 2025 gaming-inspired HUD design
type GameHUD struct {
	glowEffects     bool // Neon glow effects
	hologramStyle   bool // Translucent panels
	cinematicBars   bool // Letterbox effect
	microAnimations bool // Subtle hover effects
	depthLayers     int  // 3D-like layering

	// 2025: Immersive UI elements
	neonColors    map[string]color.RGBA
	glowIntensity float64
	panelOpacity  float64
	animationTime float64
}

// UIElement represents a UI element with depth and effects
type UIElement struct {
	Position    Position2D
	Size        Position2D
	Layer       int
	ElementType string
	Content     interface{}
}

// NewGameHUD creates a new gaming-inspired HUD
func NewGameHUD() *GameHUD {
	return &GameHUD{
		glowEffects:     true,
		hologramStyle:   true,
		cinematicBars:   false,
		microAnimations: true,
		depthLayers:     5,
		glowIntensity:   0.3,
		panelOpacity:    0.8,
		animationTime:   0.0,
		neonColors: map[string]color.RGBA{
			"primary":   {0, 255, 255, 255}, // Cyan
			"secondary": {255, 0, 255, 255}, // Magenta
			"accent":    {255, 255, 0, 255}, // Yellow
			"warning":   {255, 100, 0, 255}, // Orange
			"success":   {0, 255, 100, 255}, // Green
		},
	}
}

// RenderWithDepth renders a UI element with 2025 depth effects
func (hud *GameHUD) RenderWithDepth(screen *ebiten.Image, element UIElement, viewport *AdvancedViewportManager) {
	// Apply depth-based scaling and positioning
	depthScale := 1.0 - (float64(element.Layer) * 0.02) // Subtle depth scaling

	// Scale position and size by viewport
	scaledPos := ScalePosition(element.Position, viewport)
	scaledSize := Position2D{
		X: element.Size.X * viewport.GetIntrinsicScale() * depthScale,
		Y: element.Size.Y * viewport.GetIntrinsicScale() * depthScale,
	}

	if hud.glowEffects {
		// Add subtle glow for sci-fi aesthetic
		hud.applyGlowEffect(screen, scaledPos, scaledSize, element.Layer)
	}

	if hud.hologramStyle {
		// Semi-transparent panels with edge highlights
		hud.applyHologramEffect(screen, scaledPos, scaledSize, depthScale)
	}

	// Render the actual element content
	hud.renderElementContent(screen, element, scaledPos, scaledSize)
}

// applyGlowEffect applies 2025 neon glow effects
func (hud *GameHUD) applyGlowEffect(screen *ebiten.Image, pos, size Position2D, layer int) {
	glowColor := hud.neonColors["primary"]
	glowColor.A = uint8(255 * hud.glowIntensity * (1.0 - float64(layer)*0.1))

	// Outer glow
	glowSize := 8.0 * (1.0 - float64(layer)*0.1)
	ebitenutil.DrawRect(screen,
		pos.X-glowSize, pos.Y-glowSize,
		size.X+glowSize*2, size.Y+glowSize*2,
		glowColor)

	// Inner glow
	innerGlowColor := glowColor
	innerGlowColor.A = uint8(255 * hud.glowIntensity * 0.5)
	ebitenutil.DrawRect(screen,
		pos.X+2, pos.Y+2,
		size.X-4, size.Y-4,
		innerGlowColor)
}

// applyHologramEffect applies 2025 holographic panel effects
func (hud *GameHUD) applyHologramEffect(screen *ebiten.Image, pos, size Position2D, depthScale float64) {
	// Semi-transparent background
	panelColor := color.RGBA{0, 100, 200, uint8(255 * hud.panelOpacity * depthScale)}
	ebitenutil.DrawRect(screen, pos.X, pos.Y, size.X, size.Y, panelColor)

	// Edge highlights for holographic effect
	edgeColor := hud.neonColors["primary"]
	edgeColor.A = uint8(255 * hud.panelOpacity * depthScale)

	// Top and left edges
	ebitenutil.DrawLine(screen, pos.X, pos.Y, pos.X+size.X, pos.Y, edgeColor)
	ebitenutil.DrawLine(screen, pos.X, pos.Y, pos.X, pos.Y+size.Y, edgeColor)

	// Bottom and right edges (darker)
	darkEdgeColor := edgeColor
	darkEdgeColor.A = uint8(255 * hud.panelOpacity * 0.3 * depthScale)
	ebitenutil.DrawLine(screen, pos.X+size.X, pos.Y, pos.X+size.X, pos.Y+size.Y, darkEdgeColor)
	ebitenutil.DrawLine(screen, pos.X, pos.Y+size.Y, pos.X+size.X, pos.Y+size.Y, darkEdgeColor)
}

// renderElementContent renders the actual content of a UI element
func (hud *GameHUD) renderElementContent(screen *ebiten.Image, element UIElement, pos, size Position2D) {
	switch element.ElementType {
	case "health_bar":
		hud.renderHealthBar(screen, pos, size, element.Content)
	case "score_display":
		hud.renderScoreDisplay(screen, pos, size, element.Content)
	case "ammo_counter":
		hud.renderAmmoCounter(screen, pos, size, element.Content)
	case "status_indicator":
		hud.renderStatusIndicator(screen, pos, size, element.Content)
	}
}

// renderHealthBar renders a sci-fi health bar
func (hud *GameHUD) renderHealthBar(screen *ebiten.Image, pos, size Position2D, content interface{}) {
	// Background
	bgColor := color.RGBA{20, 20, 40, 200}
	ebitenutil.DrawRect(screen, pos.X, pos.Y, size.X, size.Y, bgColor)

	// Health fill (assuming content is a float between 0 and 1)
	if healthPercent, ok := content.(float64); ok {
		fillWidth := size.X * healthPercent
		fillColor := hud.neonColors["success"]
		if healthPercent < 0.3 {
			fillColor = hud.neonColors["warning"]
		}
		if healthPercent < 0.1 {
			fillColor = hud.neonColors["accent"]
		}

		ebitenutil.DrawRect(screen, pos.X, pos.Y, fillWidth, size.Y, fillColor)

		// Pulsing effect for low health
		if healthPercent < 0.3 {
			hud.renderPulsingEffect(screen, pos, size, fillColor)
		}
	}
}

// renderScoreDisplay renders a futuristic score display
func (hud *GameHUD) renderScoreDisplay(screen *ebiten.Image, pos, size Position2D, content interface{}) {
	// Background panel
	panelColor := color.RGBA{0, 50, 100, 180}
	ebitenutil.DrawRect(screen, pos.X, pos.Y, size.X, size.Y, panelColor)

	// Score text would be rendered here using text/v2
	// For now, just show a placeholder
	textColor := hud.neonColors["primary"]
	ebitenutil.DrawRect(screen, pos.X+5, pos.Y+5, size.X-10, size.Y-10, textColor)
}

// renderAmmoCounter renders a sci-fi ammo counter
func (hud *GameHUD) renderAmmoCounter(screen *ebiten.Image, pos, size Position2D, content interface{}) {
	// Ammo indicator dots
	if ammoCount, ok := content.(int); ok {
		dotSize := 6.0
		dotSpacing := 8.0
		startX := pos.X + 5
		startY := pos.Y + size.Y/2 - dotSize/2

		for i := 0; i < ammoCount && i < 10; i++ {
			dotColor := hud.neonColors["secondary"]
			if i < 3 {
				dotColor = hud.neonColors["warning"]
			}

			dotX := startX + float64(i)*dotSpacing
			ebitenutil.DrawCircle(screen, dotX, startY, dotSize, dotColor)
		}
	}
}

// renderStatusIndicator renders a status indicator with animations
func (hud *GameHUD) renderStatusIndicator(screen *ebiten.Image, pos, size Position2D, content interface{}) {
	// Status indicator with pulsing effect
	indicatorColor := hud.neonColors["success"]
	ebitenutil.DrawCircle(screen, pos.X+size.X/2, pos.Y+size.Y/2, size.X/2, indicatorColor)

	// Pulsing ring effect
	hud.renderPulsingEffect(screen, pos, size, indicatorColor)
}

// renderPulsingEffect renders a pulsing animation effect
func (hud *GameHUD) renderPulsingEffect(screen *ebiten.Image, pos, size Position2D, baseColor color.RGBA) {
	if !hud.microAnimations {
		return
	}

	// Simple pulsing based on time
	pulseIntensity := 0.5 + 0.5*math.Sin(hud.animationTime*4.0)
	pulseColor := baseColor
	pulseColor.A = uint8(255 * pulseIntensity * 0.3)

	// Draw pulsing ring
	ringSize := size.X * (1.0 + pulseIntensity*0.2)
	ringPos := Position2D{
		X: pos.X + (size.X-ringSize)/2,
		Y: pos.Y + (size.Y-ringSize)/2,
	}

	ebitenutil.DrawCircle(screen, ringPos.X+ringSize/2, ringPos.Y+ringSize/2, ringSize/2, pulseColor)
}

// Update updates the HUD animations
func (hud *GameHUD) Update(deltaTime float64) {
	hud.animationTime += deltaTime
}

// SetGlowEffects enables or disables glow effects
func (hud *GameHUD) SetGlowEffects(enabled bool) {
	hud.glowEffects = enabled
}

// SetHologramStyle enables or disables holographic panels
func (hud *GameHUD) SetHologramStyle(enabled bool) {
	hud.hologramStyle = enabled
}

// SetMicroAnimations enables or disables micro animations
func (hud *GameHUD) SetMicroAnimations(enabled bool) {
	hud.microAnimations = enabled
}

// GetNeonColor returns a neon color by name
func (hud *GameHUD) GetNeonColor(name string) color.RGBA {
	if color, exists := hud.neonColors[name]; exists {
		return color
	}
	return hud.neonColors["primary"] // Default
}
