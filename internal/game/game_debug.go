package game

import (
	"context"
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	v2text "github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/jonesrussell/gimbal/internal/ecs/core"
	enemysys "github.com/jonesrussell/gimbal/internal/ecs/systems/enemy"
)

// drawWaveDebugInfo draws wave information at the bottom of the screen
func (g *ECSGame) drawWaveDebugInfo(screen *ebiten.Image) {
	if g.enemySystem == nil {
		return
	}

	waveManager := g.enemySystem.GetWaveManager()
	if waveManager == nil {
		return
	}

	screenHeight := float64(g.config.ScreenSize.Height)
	lineHeight := 20.0
	x := 10.0

	currentWave := waveManager.GetCurrentWave()
	if currentWave == nil {
		g.drawNoWaveDebugInfo(screen, waveManager, x, screenHeight, lineHeight)
		return
	}

	g.drawActiveWaveDebugInfo(screen, currentWave, waveManager, struct {
		x, screenHeight, lineHeight float64
	}{x, screenHeight, lineHeight})
}

// drawNoWaveDebugInfo handles debug info when no active wave
func (g *ECSGame) drawNoWaveDebugInfo(
	screen *ebiten.Image,
	waveManager *enemysys.WaveManager,
	x, screenHeight, lineHeight float64,
) {
	if !waveManager.HasMoreWaves() {
		// All waves complete - show boss info if boss is active or spawning
		if g.enemySystem.IsBossActive() || g.enemySystem.WasBossSpawned() {
			g.drawBossDebugInfo(screen, x, screenHeight, lineHeight)
			return
		}
		// Boss not spawned yet - show spawn timer
		if g.enemySystem.WasBossSpawned() {
			// Boss was spawned but is dead
			g.drawDebugText(screen, "Boss: Defeated", x, screenHeight-lineHeight)
			return
		}
		// Boss spawning soon
		g.drawDebugText(screen, "Boss: Spawning soon...", x, screenHeight-lineHeight)
		return
	}

	// Still have waves - show waiting status
	var statusText string
	if waveManager.IsWaiting() {
		statusText = "Wave: Waiting for next wave..."
	} else {
		statusText = "Wave: Starting..."
	}
	g.drawDebugText(screen, statusText, x, screenHeight-lineHeight)
}

// drawActiveWaveDebugInfo handles debug info for active wave
func (g *ECSGame) drawActiveWaveDebugInfo(
	screen *ebiten.Image,
	currentWave *enemysys.WaveState,
	waveManager *enemysys.WaveManager,
	pos struct{ x, screenHeight, lineHeight float64 },
) {
	x := pos.x
	screenHeight := pos.screenHeight
	lineHeight := pos.lineHeight
	// Format formation type
	formationName := g.formatFormationType(currentWave.Config.FormationType)

	// Format enemy types
	enemyTypesStr := g.formatEnemyTypes(currentWave.Config.EnemyTypes)

	// Calculate progress
	progress := float64(currentWave.EnemiesKilled) / float64(currentWave.Config.EnemyCount) * 100
	if currentWave.Config.EnemyCount == 0 {
		progress = 0
	}

	// Calculate number of lines to determine starting Y position
	// Wave, Formation, Enemies, Spawned, Types, Pattern, Status, Timer
	numLines := 8
	startY := screenHeight - float64(numLines)*lineHeight - 20 // Increased margin to prevent cutoff

	// Draw wave information from bottom up
	y := startY
	g.drawDebugText(screen, fmt.Sprintf("Wave %d/%d", currentWave.WaveIndex+1, waveManager.GetWaveCount()), x, y)
	y += lineHeight
	g.drawDebugText(screen, fmt.Sprintf("Formation: %s", formationName), x, y)
	y += lineHeight
	enemyText := fmt.Sprintf("Enemies: %d/%d (%.0f%%)",
		currentWave.EnemiesKilled, currentWave.Config.EnemyCount, progress)
	g.drawDebugText(screen, enemyText, x, y)
	y += lineHeight
	g.drawDebugText(screen, fmt.Sprintf("Spawned: %d", currentWave.EnemiesSpawned), x, y)
	y += lineHeight
	g.drawDebugText(screen, fmt.Sprintf("Types: %s", enemyTypesStr), x, y)
	y += lineHeight
	patternText := fmt.Sprintf("Pattern: %s",
		g.formatMovementPattern(currentWave.Config.MovementPattern))
	g.drawDebugText(screen, patternText, x, y)
	y += lineHeight
	if currentWave.IsSpawning {
		g.drawDebugText(screen, "Status: Spawning", x, y)
	} else if currentWave.IsComplete {
		g.drawDebugText(screen, "Status: Complete", x, y)
	} else {
		g.drawDebugText(screen, "Status: Active", x, y)
	}
	y += lineHeight
	g.drawDebugText(screen, fmt.Sprintf("Timer: %.1fs", currentWave.WaveTimer), x, y)
}

// drawDebugText draws text with a semi-transparent background
func (g *ECSGame) drawDebugText(screen *ebiten.Image, text string, x, y float64) {
	ctx := g.ctx
	if ctx == nil {
		ctx = context.Background()
	}

	fontFace, err := g.resourceManager.GetDefaultFont(ctx)
	if err != nil {
		return
	}

	// Measure text size
	width, height := v2text.Measure(text, fontFace, 0)

	// Draw semi-transparent black rectangle behind text
	padding := float32(4.0)
	vector.DrawFilledRect(screen,
		float32(x)-padding,
		float32(y)-float32(height)-padding,
		float32(width)+padding*2,
		float32(height)+padding*2,
		color.RGBA{0, 0, 0, 150}, false)

	// Draw text on top
	op := &v2text.DrawOptions{}
	op.GeoM.Translate(x, y)
	v2text.Draw(screen, text, fontFace, op)
}

// formatFormationType formats a formation type as a string
func (g *ECSGame) formatFormationType(ft enemysys.FormationType) string {
	switch ft {
	case enemysys.FormationLine:
		return "Line"
	case enemysys.FormationCircle:
		return "Circle"
	case enemysys.FormationV:
		return "V"
	case enemysys.FormationDiamond:
		return "Diamond"
	case enemysys.FormationDiagonal:
		return "Diagonal"
	case enemysys.FormationSpiral:
		return "Spiral"
	case enemysys.FormationRandom:
		return "Random"
	default:
		return "Unknown"
	}
}

// formatEnemyTypes formats enemy types as a string
func (g *ECSGame) formatEnemyTypes(types []enemysys.EnemyType) string {
	if len(types) == 0 {
		return "None"
	}

	typeCounts := make(map[string]int)
	for _, t := range types {
		typeCounts[t.String()]++
	}

	result := ""
	first := true
	for name, count := range typeCounts {
		if !first {
			result += ", "
		}
		if count > 1 {
			result += fmt.Sprintf("%s x%d", name, count)
		} else {
			result += name
		}
		first = false
	}
	return result
}

// formatMovementPattern formats a movement pattern as a string
func (g *ECSGame) formatMovementPattern(mp enemysys.MovementPattern) string {
	switch mp {
	case enemysys.MovementPatternNormal:
		return "Normal"
	case enemysys.MovementPatternZigzag:
		return "Zigzag"
	case enemysys.MovementPatternAccelerating:
		return "Accelerating"
	case enemysys.MovementPatternPulsing:
		return "Pulsing"
	default:
		return "Unknown"
	}
}

// drawBossDebugInfo draws boss debug information
func (g *ECSGame) drawBossDebugInfo(screen *ebiten.Image, x, screenHeight, lineHeight float64) {
	bossEntry := g.findBossEntity()
	if bossEntry == nil {
		g.drawBossStatus(screen, x, screenHeight, lineHeight)
		return
	}
	g.drawBossDetails(screen, bossEntry, x, screenHeight, lineHeight)
}

// findBossEntity finds the boss entity in the world
func (g *ECSGame) findBossEntity() *donburi.Entry {
	var bossEntry *donburi.Entry
	query.NewQuery(
		filter.And(
			filter.Contains(core.EnemyTag),
			filter.Contains(core.EnemyTypeID),
		),
	).Each(g.world, func(entry *donburi.Entry) {
		typeID := core.EnemyTypeID.Get(entry)
		if enemysys.EnemyType(*typeID) == enemysys.EnemyTypeBoss {
			bossEntry = entry
		}
	})
	return bossEntry
}

// drawBossStatus draws boss spawn/defeat status
func (g *ECSGame) drawBossStatus(screen *ebiten.Image, x, screenHeight, lineHeight float64) {
	if g.enemySystem.WasBossSpawned() {
		g.drawDebugText(screen, "Boss: Defeated", x, screenHeight-lineHeight)
	} else {
		g.drawDebugText(screen, "Boss: Spawning soon...", x, screenHeight-lineHeight)
	}
}

// drawBossDetails draws detailed boss information
func (g *ECSGame) drawBossDetails(screen *ebiten.Image, bossEntry *donburi.Entry, x, screenHeight, lineHeight float64) {
	pos := core.Position.Get(bossEntry)
	health := core.Health.Get(bossEntry)
	orbital := core.Orbital.Get(bossEntry)
	size := core.Size.Get(bossEntry)

	// Calculate number of lines for boss info
	numLines := 6 // Boss, Health, Position, Orbital Angle, Size, Status
	startY := screenHeight - float64(numLines)*lineHeight - 20

	// Draw boss information from bottom up
	y := startY
	g.drawDebugText(screen, "BOSS", x, y)
	y += lineHeight

	if health != nil {
		healthPercent := float64(health.Current) / float64(health.Maximum) * 100
		healthText := fmt.Sprintf("Health: %d/%d (%.0f%%)",
			health.Current, health.Maximum, healthPercent)
		g.drawDebugText(screen, healthText, x, y)
	} else {
		g.drawDebugText(screen, "Health: Unknown", x, y)
	}
	y += lineHeight

	if pos != nil {
		g.drawDebugText(screen, fmt.Sprintf("Position: (%.0f, %.0f)", pos.X, pos.Y), x, y)
	} else {
		g.drawDebugText(screen, "Position: Unknown", x, y)
	}
	y += lineHeight

	if orbital != nil {
		g.drawDebugText(screen, fmt.Sprintf("Orbital Angle: %.1f°", float64(orbital.OrbitalAngle)), x, y)
	} else {
		g.drawDebugText(screen, "Orbital: Unknown", x, y)
	}
	y += lineHeight

	if size != nil {
		g.drawDebugText(screen, fmt.Sprintf("Size: %dx%d", size.Width, size.Height), x, y)
	} else {
		g.drawDebugText(screen, "Size: Unknown", x, y)
	}
	y += lineHeight

	status := "Active"
	if health != nil && health.Current <= 0 {
		status = "Defeated"
	}
	g.drawDebugText(screen, fmt.Sprintf("Status: %s", status), x, y)
}

