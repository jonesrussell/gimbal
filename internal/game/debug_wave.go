package game

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	v2text "github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/jonesrussell/gimbal/internal/ecs/systems/stage"
)

// drawWaveDebugInfo draws wave information at the bottom of the screen
func (g *ECSGame) drawWaveDebugInfo(screen *ebiten.Image) {
	if g.stageStateMachine == nil {
		return
	}

	screenHeight := float64(g.config.ScreenSize.Height)
	lineHeight := 20.0
	x := 10.0

	// Get stage info from stage state machine
	stageConfig := g.stageStateMachine.StageConfig()
	if stageConfig == nil {
		g.drawDebugText(screen, "Stage: Not loaded", x, screenHeight-lineHeight)
		return
	}

	// Boss lifecycle from stage state (StageStateMachine is single source of truth)
	st := g.stageStateMachine.State()
	switch st {
	case stage.StageStateBossSpawning:
		g.drawDebugText(screen, "Boss: Spawning soon...", x, screenHeight-lineHeight)
		return
	case stage.StageStateBossActive:
		g.drawBossDebugInfo(screen, x, screenHeight, lineHeight)
		return
	case stage.StageStateBossDefeated, stage.StageStateStageCompleted:
		g.drawDebugText(screen, "Boss: Defeated", x, screenHeight-lineHeight)
		return
	}

	// PreWave (level start or inter-wave)
	if st == stage.StageStatePreWave {
		g.drawDebugText(screen, "Stage: Starting...", x, screenHeight-lineHeight)
		return
	}

	// WaveCompleted (brief transition) or wave info
	waveCount := len(stageConfig.Waves)
	currentWaveIndex := g.stageStateMachine.CurrentWaveIndex()
	if currentWaveIndex >= waveCount && st != stage.StageStateWaveInProgress {
		g.drawDebugText(screen, "Waves: Complete", x, screenHeight-lineHeight)
		return
	}

	// Calculate number of lines
	numLines := 4
	startY := screenHeight - float64(numLines)*lineHeight - 20

	y := startY
	g.drawDebugText(screen, fmt.Sprintf("Stage %d: %s", stageConfig.StageNumber, stageConfig.Metadata.Name), x, y)
	y += lineHeight
	g.drawDebugText(screen, fmt.Sprintf("Wave %d/%d", currentWaveIndex+1, waveCount), x, y)
	y += lineHeight

	// Show current wave description if available
	if currentWaveIndex < len(stageConfig.Waves) {
		wave := stageConfig.Waves[currentWaveIndex]
		g.drawDebugText(screen, fmt.Sprintf("Wave: %s", wave.Description), x, y)
		y += lineHeight
	}

	g.drawDebugText(screen, fmt.Sprintf("Stage: %s", stageConfig.Planet), x, y)
}

// drawDebugText draws text with a semi-transparent background
func (g *ECSGame) drawDebugText(screen *ebiten.Image, text string, x, y float64) {
	// Context is always initialized in NewECSGame, so no nil check needed
	fontFace, err := g.resourceManager.GetDefaultFont(g.ctx)
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
