package game

import (
	"context"
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	v2text "github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

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
