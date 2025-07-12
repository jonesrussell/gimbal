package debug

import (
	"fmt"
	"image/color"
	"math"
	"runtime"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	v2text "github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/ecs/core"
)

// DebugLevel represents different levels of debug information
type DebugLevel int

const (
	DebugOff      DebugLevel = iota
	DebugBasic               // Just grid + performance
	DebugDetailed            // Everything
)

// DebugRenderer handles all debug visualization and metrics
type DebugRenderer struct {
	enabled    bool
	level      DebugLevel
	font       v2text.Face
	config     *common.GameConfig
	logger     common.Logger
	mouseX     int
	mouseY     int
	hoverRange float64 // Distance for mouse hover entity info
}

// NewDebugRenderer creates a new debug renderer
func NewDebugRenderer(config *common.GameConfig, logger common.Logger) *DebugRenderer {
	return &DebugRenderer{
		enabled:    false,
		level:      DebugBasic,
		config:     config,
		logger:     logger,
		hoverRange: 50.0, // Show entity info within 50 pixels of mouse
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
	dr.logger.Debug("Debug mode toggled", "enabled", dr.enabled, "level", dr.level)
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
	dr.logger.Debug("Debug level cycled", "enabled", dr.enabled, "level", dr.level)
}

// IsEnabled returns whether debug mode is active
func (dr *DebugRenderer) IsEnabled() bool {
	return dr.enabled
}

// SetFont sets the font for debug text rendering
func (dr *DebugRenderer) SetFont(font v2text.Face) {
	dr.font = font
}

// UpdateMousePosition updates the current mouse position for hover detection
func (dr *DebugRenderer) UpdateMousePosition() {
	dr.mouseX, dr.mouseY = ebiten.CursorPosition()
}

// shouldShowEntityInfo determines if entity info should be shown based on mouse proximity
func (dr *DebugRenderer) shouldShowEntityInfo(pos *common.Point) bool {
	if dr.level != DebugDetailed {
		return false
	}

	// Calculate distance from mouse cursor
	dist := math.Sqrt(math.Pow(pos.X-float64(dr.mouseX), 2) + math.Pow(pos.Y-float64(dr.mouseY), 2))
	return dist < dr.hoverRange
}

// drawTextWithBackground draws text with a subtle dark background for better readability
func (dr *DebugRenderer) drawTextWithBackground(screen *ebiten.Image, text string, x, y float64) {
	if dr.font == nil {
		return
	}

	// Measure text bounds
	width, height := v2text.Measure(text, dr.font, 0)

	// Draw semi-transparent black rectangle behind text
	padding := 4.0
	ebitenutil.DrawRect(screen,
		x-padding,
		y-height-padding,
		width+padding*2,
		height+padding*2,
		color.RGBA{0, 0, 0, 100})

	// Draw text on top
	op := &v2text.DrawOptions{}
	op.GeoM.Translate(x, y)
	v2text.Draw(screen, text, dr.font, op)
}

// Render draws all debug information
func (dr *DebugRenderer) Render(screen *ebiten.Image, world donburi.World) {
	if !dr.enabled {
		return
	}

	// Update mouse position for hover detection
	dr.UpdateMousePosition()

	// Draw debug grid
	dr.drawGrid(screen)

	// Draw performance metrics
	dr.drawPerformanceMetrics(screen, world)

	// Draw entity debug info (only in detailed mode)
	if dr.level == DebugDetailed {
		dr.drawEntityDebug(screen, world)
		dr.drawCollisionDebug(screen, world)
	}
}

// drawGrid draws a debug grid overlay
func (dr *DebugRenderer) drawGrid(screen *ebiten.Image) {
	bounds := screen.Bounds()
	gridSize := 50

	// Draw vertical lines - barely visible guide lines
	for x := 0; x < bounds.Dx(); x += gridSize {
		ebitenutil.DrawLine(screen, float64(x), 0, float64(x), float64(bounds.Dy()),
			color.RGBA{255, 255, 255, 20})
	}

	// Draw horizontal lines - barely visible guide lines
	for y := 0; y < bounds.Dy(); y += gridSize {
		ebitenutil.DrawLine(screen, 0, float64(y), float64(bounds.Dx()), float64(y),
			color.RGBA{255, 255, 255, 20})
	}
}

// drawPerformanceMetrics draws condensed performance information
func (dr *DebugRenderer) drawPerformanceMetrics(screen *ebiten.Image, world donburi.World) {
	if dr.font == nil {
		return
	}

	// Get memory stats
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// Get entity count
	entityCount := world.Len()

	// Format condensed debug text
	debugText := fmt.Sprintf("FPS:%.0f TPS:%.0f E:%d M:%dK [F3]",
		ebiten.ActualFPS(),
		ebiten.ActualTPS(),
		entityCount,
		m.Alloc/1024)

	// Draw text with background for better readability
	dr.drawTextWithBackground(screen, debugText, 10, 30)
}

// drawEntityDebug draws entity positions, bounds, and states (only for nearby entities)
func (dr *DebugRenderer) drawEntityDebug(screen *ebiten.Image, world donburi.World) {
	// Draw all entities with position and sprite components
	query.NewQuery(
		filter.And(
			filter.Contains(core.Position),
			filter.Contains(core.Sprite),
		),
	).Each(world, func(entry *donburi.Entry) {
		pos := core.Position.Get(entry)
		sprite := core.Sprite.Get(entry)

		if pos == nil || sprite == nil {
			return
		}

		// Determine entity type and color
		entityColor := dr.getEntityColor(entry)

		// Draw entity center point - tiny colored dot
		ebitenutil.DrawCircle(screen, pos.X, pos.Y, 1, entityColor)

		// Draw bounding box if size component exists
		if entry.HasComponent(core.Size) {
			size := core.Size.Get(entry)
			if size != nil {
				// Calculate bounds
				boundsX := pos.X - float64(size.Width)/2
				boundsY := pos.Y - float64(size.Height)/2

				// Draw bounding box - very thin colored outline
				ebitenutil.DrawRect(screen, boundsX, boundsY, float64(size.Width), float64(size.Height),
					entityColor)

				// Only show entity info text if mouse is nearby
				if dr.shouldShowEntityInfo(pos) {
					entityInfo := fmt.Sprintf("Pos: (%.1f,%.1f)\nSize: %dx%d", pos.X, pos.Y, size.Width, size.Height)
					dr.drawTextWithBackground(screen, entityInfo, pos.X+10, pos.Y-20)
				}
			}
		}

		// Draw sprite bounds (only if mouse is nearby)
		if dr.shouldShowEntityInfo(pos) {
			dr.drawSpriteDebug(screen, *sprite, pos.X, pos.Y)
		}
	})
}

// getEntityColor returns the appropriate color for different entity types
func (dr *DebugRenderer) getEntityColor(entry *donburi.Entry) color.RGBA {
	// Player entities - Green
	if entry.HasComponent(core.PlayerTag) {
		return color.RGBA{0, 255, 0, 80}
	}

	// Enemy entities - Red
	if entry.HasComponent(core.EnemyTag) {
		return color.RGBA{255, 0, 0, 80}
	}

	// Projectile entities - Yellow
	if entry.HasComponent(core.ProjectileTag) {
		return color.RGBA{255, 255, 0, 60}
	}

	// Star entities - Blue
	if entry.HasComponent(core.StarTag) {
		return color.RGBA{0, 150, 255, 40}
	}

	// Default - White
	return color.RGBA{255, 255, 255, 50}
}

// drawSpriteDebug draws sprite boundaries and center points
func (dr *DebugRenderer) drawSpriteDebug(screen *ebiten.Image, sprite *ebiten.Image, x, y float64) {
	bounds := sprite.Bounds()

	// Calculate sprite position (assuming sprite is centered on entity)
	spriteX := x - float64(bounds.Dx())/2
	spriteY := y - float64(bounds.Dy())/2

	// Draw sprite boundary rectangle - very thin white outline
	ebitenutil.DrawRect(screen, spriteX, spriteY, float64(bounds.Dx()), float64(bounds.Dy()),
		color.RGBA{255, 255, 255, 30})

	// Draw sprite center point - tiny white dot
	centerX, centerY := x, y
	ebitenutil.DrawCircle(screen, centerX, centerY, 1, color.RGBA{255, 255, 255, 80})
}

// drawCollisionDebug draws collision boundaries and detection ranges
func (dr *DebugRenderer) drawCollisionDebug(screen *ebiten.Image, world donburi.World) {
	// Draw player collision area
	query.NewQuery(
		filter.And(
			filter.Contains(core.PlayerTag),
			filter.Contains(core.Position),
			filter.Contains(core.Size),
		),
	).Each(world, func(entry *donburi.Entry) {
		pos := core.Position.Get(entry)
		size := core.Size.Get(entry)

		if pos == nil || size == nil {
			return
		}

		// Only show collision debug if mouse is nearby
		if !dr.shouldShowEntityInfo(pos) {
			return
		}

		// Draw player collision box - very thin green outline
		boundsX := pos.X - float64(size.Width)/2
		boundsY := pos.Y - float64(size.Height)/2
		ebitenutil.DrawRect(screen, boundsX, boundsY, float64(size.Width), float64(size.Height),
			color.RGBA{0, 255, 0, 60})
	})

	// Draw enemy collision areas
	query.NewQuery(
		filter.And(
			filter.Contains(core.EnemyTag),
			filter.Contains(core.Position),
			filter.Contains(core.Size),
		),
	).Each(world, func(entry *donburi.Entry) {
		pos := core.Position.Get(entry)
		size := core.Size.Get(entry)

		if pos == nil || size == nil {
			return
		}

		// Only show collision debug if mouse is nearby
		if !dr.shouldShowEntityInfo(pos) {
			return
		}

		// Draw enemy collision box - very thin red outline
		boundsX := pos.X - float64(size.Width)/2
		boundsY := pos.Y - float64(size.Height)/2
		ebitenutil.DrawRect(screen, boundsX, boundsY, float64(size.Width), float64(size.Height),
			color.RGBA{255, 0, 0, 60})
	})
}
