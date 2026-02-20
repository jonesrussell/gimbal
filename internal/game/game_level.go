package game

import (
	"log"

	"github.com/jonesrussell/gimbal/internal/dbg"
	"github.com/jonesrussell/gimbal/internal/ecs/systems/stage"
	"github.com/jonesrussell/gimbal/internal/scenes"
)

// checkLevelCompletion checks if the stage is complete and advances to the next
func (g *ECSGame) checkLevelCompletion() {
	state := g.stageStateMachine.State()
	completed := g.stageStateMachine.IsStageCompleted()
	// Only log when state is boss-related or completed to avoid per-frame spam
	if completed || state == stage.StageStateBossActive || state == stage.StageStateBossDefeated || state == stage.StageStateStageCompleted {
		log.Printf("[BOSS_DEBUG] checkLevelCompletion: IsStageCompleted=%v state=%v", completed, state)
	}
	if !completed {
		return
	}
	log.Printf("[BOSS_DEBUG] checkLevelCompletion: calling handleLevelComplete")
	g.handleLevelComplete()
}

// handleLevelComplete handles level completion actions
func (g *ECSGame) handleLevelComplete() {
	currentStage := g.stageStateMachine.StageNumber()
	log.Printf("[BOSS_DEBUG] handleLevelComplete: stage=%d", currentStage)
	dbg.Log(dbg.State, "Stage complete (stage=%d)", currentStage)

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
		log.Printf("[ERROR] Failed to load next stage: %v", err)
		// Still switch scene so the game doesn't hang; next play will retry same stage
	}

	// Show between-stage transition
	// The transition scene will load the next stage and show stage intro
	g.sceneManager.SwitchScene(scenes.SceneStageTransition)
}
