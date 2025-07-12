package game

import (
	"github.com/yohamta/donburi"

	"github.com/jonesrussell/gimbal/internal/ecs/events"
	"github.com/jonesrussell/gimbal/internal/scenes"
)

// setupEventSubscriptions sets up event handlers
func (g *ECSGame) setupEventSubscriptions() {
	// Subscribe to player movement events
	g.eventSystem.SubscribeToPlayerMoved(func(w donburi.World, event events.PlayerMovedEvent) {
		g.logger.Debug("Player moved",
			"position", event.Position,
			"angle", event.Angle)
	})

	// Subscribe to game state events
	g.eventSystem.SubscribeToGameState(func(w donburi.World, event events.GameStateEvent) {
		g.logger.Debug("Game state changed", "is_paused", event.IsPaused)
	})

	// Subscribe to score changes
	g.eventSystem.SubscribeToScoreChanged(func(w donburi.World, event events.ScoreChangedEvent) {
		g.logger.Debug("Score changed",
			"old_score", event.OldScore,
			"new_score", event.NewScore,
			"delta", event.Delta)
	})

	// Subscribe to game over events
	g.eventSystem.SubscribeToGameOver(func(w donburi.World, event events.GameOverEvent) {
		g.logger.Debug("Game over triggered", "reason", event.Reason)
		g.sceneManager.SwitchScene(scenes.SceneGameOver)
	})

	// Subscribe to player damage events for screen shake
	g.eventSystem.SubscribeToPlayerDamaged(func(w donburi.World, event events.PlayerDamagedEvent) {
		g.logger.Debug("Player damaged", "damage", event.Damage, "remaining_lives", event.RemainingLives)
		// Trigger screen shake if we're in the playing scene
		if g.sceneManager.GetCurrentScene().GetType() == scenes.ScenePlaying {
			// TODO: Re-enable after PlayingScene is moved to gameplay package
			// if playingScene, ok := g.sceneManager.GetCurrentScene().(*scenes.PlayingScene); ok {
			// 	playingScene.TriggerScreenShake()
			// }
		}
	})

	// Subscribe to enemy destroyed events for scoring
	g.eventSystem.SubscribeToEnemyDestroyed(func(w donburi.World, event events.EnemyDestroyedEvent) {
		g.scoreManager.AddScore(event.Points)
		g.logger.Debug(
			"Score added from enemy destruction",
			"points", event.Points,
			"total_score", g.scoreManager.GetScore(),
		)
	})
}
