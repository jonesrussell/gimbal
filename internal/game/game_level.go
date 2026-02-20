package game

import (
	"github.com/jonesrussell/gimbal/internal/scenes"
)

// checkLevelCompletion checks if the stage is complete and advances to the next
func (g *ECSGame) checkLevelCompletion() {
	if !g.stageStateMachine.IsStageCompleted() {
		return
	}

	g.handleLevelComplete()
}

// handleLevelComplete handles level completion actions
func (g *ECSGame) handleLevelComplete() {
	currentStage := g.stageStateMachine.StageNumber()
	g.logger.Debug("Stage complete", "stage", currentStage)

	// Check if final stage (stage 6)
	if currentStage >= 6 {
		// Show victory sequence
		g.sceneManager.SwitchScene(scenes.SceneVictory)
		return
	}

	// Update level manager to track progression
	g.levelManager.IncrementLevel()

	// Load next stage via stage state machine (delegates to GyrussSystem)
	if err := g.stageStateMachine.LoadNextStage(); err != nil {
		g.logger.Error("Failed to load next stage", "error", err)
		// Still switch scene so the game doesn't hang; next play will retry same stage
	}

	// Show between-stage transition
	// The transition scene will load the next stage and show stage intro
	g.sceneManager.SwitchScene(scenes.SceneStageTransition)
}
