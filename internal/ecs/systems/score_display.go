package systems

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features"
	"github.com/yohamta/donburi/filter"

	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/ecs"
)

// ScoreDisplaySystem displays the current score in the HUD
func ScoreDisplaySystem(world *donburi.World, font *ebiten.Font, scoreManager *ecs.ScoreManager) donburi.System {
	return func() {
		// Get score values
		score := scoreManager.GetScore()
		highScore := scoreManager.GetHighScore()
		multiplier := scoreManager.GetMultiplier()

		// Draw score
		x := 10
		y := 20
		common.DrawText(world, font, "Score:", x, y, common.ColorWhite)
		common.DrawText(world, font, common.FormatNumber(score), x+70, y, common.ColorYellow)

		// Draw high score
		y += 20
		common.DrawText(world, font, "High Score:", x, y, common.ColorWhite)
		common.DrawText(world, font, common.FormatNumber(highScore), x+100, y, common.ColorYellow)

		// Draw multiplier
		y += 20
		common.DrawText(world, font, "Multiplier:", x, y, common.ColorWhite)
		common.DrawText(world, font, fmt.Sprintf("x%d", multiplier), x+90, y, common.ColorYellow)

		// Draw bonus life threshold
		y += 20
		bonusLifeScore := scoreManager.bonusLifeScore
		bonusLifeCount := score / bonusLifeScore
		common.DrawText(world, font, "Bonus Lives:", x, y, common.ColorWhite)
		common.DrawText(world, font, fmt.Sprintf("%d", bonusLifeCount), x+100, y, common.ColorYellow)
	}
}
