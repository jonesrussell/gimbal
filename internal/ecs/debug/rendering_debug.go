package debug

import (
	"fmt"
	"image/color"
	"log"
	"runtime"

	"github.com/hajimehoshi/ebiten/v2"
	v2text "github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/jonesrussell/gimbal/internal/ecs/core"
)

// RenderingDebugger provides detailed rendering diagnostics
type RenderingDebugger struct {
	enabled     bool
	font        v2text.Face
	drawCalls   int
	entityCount int
	spriteCount int
}

// NewRenderingDebugger creates a new rendering debugger
func NewRenderingDebugger(font v2text.Face) *RenderingDebugger {
	return &RenderingDebugger{
		enabled: false,
		font:    font,
	}
}

// Toggle enables/disables rendering debug mode
func (rd *RenderingDebugger) Toggle() {
	rd.enabled = !rd.enabled
	log.Printf("[DEBUG] Rendering debug toggled enabled=%v", rd.enabled)
}

// IsEnabled returns whether rendering debug is active
func (rd *RenderingDebugger) IsEnabled() bool {
	return rd.enabled
}

// StartFrame resets frame statistics
func (rd *RenderingDebugger) StartFrame() {
	rd.drawCalls = 0
	rd.entityCount = 0
	rd.spriteCount = 0
}

// RecordDrawCall records a draw call for statistics
func (rd *RenderingDebugger) RecordDrawCall() {
	rd.drawCalls++
}

// RecordEntity records an entity for statistics
func (rd *RenderingDebugger) RecordEntity() {
	rd.entityCount++
}

// RecordSprite records a sprite for statistics
func (rd *RenderingDebugger) RecordSprite() {
	rd.spriteCount++
}

// RenderDebugInfo renders comprehensive rendering debug information
func (rd *RenderingDebugger) RenderDebugInfo(screen *ebiten.Image, world donburi.World) {
	if !rd.enabled || rd.font == nil {
		return
	}

	// Get memory stats
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// Analyze rendering performance
	rd.analyzeRendering(world)

	// Draw debug information
	rd.drawRenderingStats(screen, &m)
	rd.drawEntityAnalysis(screen, world)
	rd.drawPerformanceWarnings(screen, &m)
}

// analyzeRendering analyzes the current rendering state
func (rd *RenderingDebugger) analyzeRendering(world donburi.World) {
	// Count entities with different components
	entitiesWithSprite := 0
	entitiesWithPosition := 0
	entitiesWithSize := 0
	entitiesWithScale := 0
	entitiesWithRotation := 0

	query.NewQuery(
		filter.And(
			filter.Contains(core.Position),
			filter.Contains(core.Sprite),
		),
	).Each(world, func(entry *donburi.Entry) {
		entitiesWithSprite++
		entitiesWithPosition++

		if entry.HasComponent(core.Size) {
			entitiesWithSize++
		}
		if entry.HasComponent(core.Scale) {
			entitiesWithScale++
		}
		if entry.HasComponent(core.Orbital) || entry.HasComponent(core.Angle) {
			entitiesWithRotation++
		}
	})

	// Log analysis
	log.Printf(
		"[DEBUG] Rendering analysis entities_with_sprite=%d entities_with_size=%d "+
			"entities_with_scale=%d entities_with_rotation=%d draw_calls=%d",
		entitiesWithSprite, entitiesWithSize, entitiesWithScale, entitiesWithRotation, rd.drawCalls,
	)
}

// drawRenderingStats draws rendering statistics
func (rd *RenderingDebugger) drawRenderingStats(screen *ebiten.Image, m *runtime.MemStats) {
	stats := []string{
		fmt.Sprintf("Draw Calls: %d", rd.drawCalls),
		fmt.Sprintf("Entities: %d", rd.entityCount),
		fmt.Sprintf("Sprites: %d", rd.spriteCount),
		fmt.Sprintf("Memory: %dKB", m.Alloc/1024),
		fmt.Sprintf("FPS: %.0f", ebiten.ActualFPS()),
	}

	y := float64(100)
	for _, stat := range stats {
		rd.drawTextWithBackground(screen, stat, 10, y)
		y += 20
	}
}

// drawEntityAnalysis draws entity-specific analysis
func (rd *RenderingDebugger) drawEntityAnalysis(screen *ebiten.Image, world donburi.World) {
	// Count entities by type
	playerCount := 0
	enemyCount := 0
	projectileCount := 0
	starCount := 0

	query.NewQuery(
		filter.And(
			filter.Contains(core.Position),
			filter.Contains(core.Sprite),
		),
	).Each(world, func(entry *donburi.Entry) {
		if entry.HasComponent(core.PlayerTag) {
			playerCount++
		}
		if entry.HasComponent(core.EnemyTag) {
			enemyCount++
		}
		if entry.HasComponent(core.ProjectileTag) {
			projectileCount++
		}
		if entry.HasComponent(core.StarTag) {
			starCount++
		}
	})

	entityStats := []string{
		fmt.Sprintf("Players: %d", playerCount),
		fmt.Sprintf("Enemies: %d", enemyCount),
		fmt.Sprintf("Projectiles: %d", projectileCount),
		fmt.Sprintf("Stars: %d", starCount),
	}

	y := float64(220)
	for _, stat := range entityStats {
		rd.drawTextWithBackground(screen, stat, 10, y)
		y += 20
	}
}

// drawPerformanceWarnings draws performance warnings
func (rd *RenderingDebugger) drawPerformanceWarnings(screen *ebiten.Image, m *runtime.MemStats) {
	warnings := []string{}

	// Check for performance issues
	if rd.drawCalls > 1000 {
		warnings = append(warnings, "WARNING: High draw call count")
	}
	if rd.entityCount > 500 {
		warnings = append(warnings, "WARNING: High entity count")
	}
	if m.Alloc > 100*1024*1024 { // 100MB
		warnings = append(warnings, "WARNING: High memory usage")
	}
	if ebiten.ActualFPS() < 50 {
		warnings = append(warnings, "WARNING: Low FPS detected")
	}

	// Draw warnings in red
	y := float64(350)
	for _, warning := range warnings {
		rd.drawWarningText(screen, warning, 10, y)
		y += 20
	}
}

// drawTextWithBackground draws text with a semi-transparent background
func (rd *RenderingDebugger) drawTextWithBackground(screen *ebiten.Image, str string, x, y float64) {
	if rd.font == nil {
		return
	}

	// Measure text size using v2text
	width, height := v2text.Measure(str, rd.font, 0)

	// Draw semi-transparent black rectangle behind text
	padding := float32(4.0)
	vector.DrawFilledRect(screen,
		float32(x)-padding,
		float32(y)-float32(height)-padding,
		float32(width)+padding*2,
		float32(height)+padding*2,
		color.RGBA{0, 0, 0, 100}, false)

	// Draw text on top using v2text
	op := &v2text.DrawOptions{}
	op.GeoM.Translate(x, y)
	v2text.Draw(screen, str, rd.font, op)
}

// drawWarningText draws warning text in red
func (rd *RenderingDebugger) drawWarningText(screen *ebiten.Image, str string, x, y float64) {
	if rd.font == nil {
		return
	}

	// Measure text size using v2text
	width, height := v2text.Measure(str, rd.font, 0)

	// Draw semi-transparent red rectangle behind text
	padding := float32(4.0)
	vector.DrawFilledRect(screen,
		float32(x)-padding,
		float32(y)-float32(height)-padding,
		float32(width)+padding*2,
		float32(height)+padding*2,
		color.RGBA{255, 0, 0, 100}, false)

	// Draw text on top using v2text
	op := &v2text.DrawOptions{}
	op.GeoM.Translate(x, y)
	v2text.Draw(screen, str, rd.font, op)
}
