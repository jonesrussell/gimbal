package game

import (
	"github.com/jonesrussell/gimbal/internal/scenes"
)

// checkLevelCompletion checks if the stage is complete and advances to the next
func (g *ECSGame) checkLevelCompletion() {
	// Check if current stage is complete (boss defeated)
	if !g.gyrussSystem.IsStageComplete() {
		return
	}

	g.handleLevelComplete()
}

// handleLevelComplete handles level completion actions
func (g *ECSGame) handleLevelComplete() {
	currentStage := g.gyrussSystem.GetCurrentStage()
	g.logger.Debug("Stage complete", "stage", currentStage)

	// Check if final stage (stage 6)
	if currentStage >= 6 {
		// Show victory sequence
		g.sceneManager.SwitchScene(scenes.SceneVictory)
		return
	}

	// Update level manager to track progression
	g.levelManager.IncrementLevel()

	// Show between-stage transition
	// The transition scene will load the next stage and show stage intro
	g.sceneManager.SwitchScene(scenes.SceneStageTransition)
}
