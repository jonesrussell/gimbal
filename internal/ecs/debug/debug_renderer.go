package debug

import (
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	v2text "github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/yohamta/donburi"

	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/config"
)

// DebugLevel represents different levels of debug information
type DebugLevel int

const (
	DebugOff      DebugLevel = iota
	DebugBasic               // FPS and entity count only
	DebugDetailed            // Adds 50px grid and entity/collision debug
)

// DebugRenderer handles all debug visualization and metrics
type DebugRenderer struct {
	enabled    bool
	level      DebugLevel
	font       v2text.Face
	config     *config.GameConfig
	mouseX     int
	mouseY     int
	hoverRange float64 // Distance for mouse hover entity info
}

// NewDebugRenderer creates a new debug renderer
func NewDebugRenderer(gameConfig *config.GameConfig) *DebugRenderer {
	return &DebugRenderer{
		enabled:    false,
		level:      DebugBasic,
		config:     gameConfig,
		hoverRange: 50.0,
	}
}

// Toggle switches debug mode on/off
func (dr *DebugRenderer) Toggle() {
	dr.enabled = !dr.enabled
	if dr.enabled {
		dr.level = DebugBasic
	} else {
		dr.level = DebugOff
	}
	log.Printf("[DEBUG] Debug mode toggled enabled=%v level=%v", dr.enabled, dr.level)
}

// CycleLevel cycles through debug levels (Basic -> Detailed -> Off)
func (dr *DebugRenderer) CycleLevel() {
	if !dr.enabled {
		dr.enabled = true
		dr.level = DebugBasic
	} else {
		switch dr.level {
		case DebugBasic:
			dr.level = DebugDetailed
		case DebugDetailed:
			dr.enabled = false
			dr.level = DebugOff
		}
	}
	log.Printf("[DEBUG] Debug level cycled enabled=%v level=%v", dr.enabled, dr.level)
}

// IsEnabled returns whether debug mode is active
func (dr *DebugRenderer) IsEnabled() bool {
	return dr.enabled
}

// SetFont sets the font for debug text rendering
func (dr *DebugRenderer) SetFont(font v2text.Face) {
	dr.font = font
}

// UpdateMousePosition updates the mouse position for hover detection
func (dr *DebugRenderer) UpdateMousePosition() {
	dr.mouseX, dr.mouseY = ebiten.CursorPosition()
}

// shouldShowEntityInfo returns true if the mouse is near the given position
func (dr *DebugRenderer) shouldShowEntityInfo(pos *common.Point) bool {
	if pos == nil {
		return false
	}
	distance := math.Sqrt(math.Pow(float64(dr.mouseX)-pos.X, 2) + math.Pow(float64(dr.mouseY)-pos.Y, 2))
	return distance <= dr.hoverRange
}

// drawTextWithBackground draws text with a semi-transparent background
func (dr *DebugRenderer) drawTextWithBackground(screen *ebiten.Image, text string, x, y float64) {
	if dr.font == nil {
		return
	}

	// Measure text bounds
	width, height := v2text.Measure(text, dr.font, 0)

	// Draw semi-transparent black rectangle behind text
	padding := float32(4.0)
	vector.DrawFilledRect(screen,
		float32(x)-padding,
		float32(y-height)-padding,
		float32(width)+padding*2,
		float32(height)+padding*2,
		color.RGBA{0, 0, 0, 100}, false)

	// Draw text on top
	op := &v2text.DrawOptions{}
	op.GeoM.Translate(x, y)
	v2text.Draw(screen, text, dr.font, op)
}

// Render draws debug information for the current level: Basic = performance metrics; Detailed adds grid, entity debug, and collision debug.
func (dr *DebugRenderer) Render(screen *ebiten.Image, world donburi.World) {
	if !dr.enabled {
		return
	}

	// Update mouse position for hover detection
	dr.UpdateMousePosition()

	// Draw debug grid only in Detailed mode (Basic = FPS/entity count only)
	if dr.level == DebugDetailed {
		dr.drawGrid(screen)
	}

	// Draw performance metrics
	dr.drawPerformanceMetrics(screen, world)

	// Draw entity debug info (only in detailed mode)
	if dr.level == DebugDetailed {
		dr.drawEntityDebug(screen, world)
		dr.drawCollisionDebug(screen, world)
	}
}
