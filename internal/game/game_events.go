package game

import (
	"github.com/yohamta/donburi"

	"github.com/jonesrussell/gimbal/internal/dbg"
	"github.com/jonesrussell/gimbal/internal/ecs/events"
	"github.com/jonesrussell/gimbal/internal/scenes"
	"github.com/jonesrussell/gimbal/internal/ui/presenter"
)

// setupEventSubscriptions sets up event handlers
func (g *ECSGame) setupEventSubscriptions() {
	// Initialize HUD presenter with current game state
	g.initializeHUDPresenter()
	// Subscribe to player movement events (no per-move logging to avoid 60/s spam)
	g.eventSystem.SubscribeToPlayerMoved(func(w donburi.World, event events.PlayerMovedEvent) {})

	// Subscribe to game state events
	g.eventSystem.SubscribeToGameState(func(w donburi.World, event events.GameStateEvent) {
		dbg.Log(dbg.Event, "Game state changed (is_paused=%v)", event.IsPaused)
	})

	// Subscribe to score changes
	g.eventSystem.SubscribeToScoreChanged(func(w donburi.World, event events.ScoreChangedEvent) {
		dbg.Log(dbg.Event, "Score changed (old=%d new=%d delta=%d)", event.OldScore, event.NewScore, event.Delta)
	})

	// Subscribe to game over events
	g.eventSystem.SubscribeToGameOver(func(w donburi.World, event events.GameOverEvent) {
		dbg.Log(dbg.Event, "Game over triggered (reason=%s)", event.Reason)
		g.sceneManager.SwitchScene(scenes.SceneGameOver)
	})

	// Subscribe to player damage events for screen shake
	g.eventSystem.SubscribeToPlayerDamaged(func(w donburi.World, event events.PlayerDamagedEvent) {
		dbg.Log(dbg.Event, "Player damaged (damage=%d remaining_lives=%d)", event.Damage, event.RemainingLives)
		// TODO: Re-enable after PlayingScene is moved to gameplay package
		// Trigger screen shake if we're in the playing scene
		// if g.sceneManager.GetCurrentScene().GetType() == scenes.ScenePlaying {
		// 	if playingScene, ok := g.sceneManager.GetCurrentScene().(*scenes.PlayingScene); ok {
		// 		playingScene.TriggerScreenShake()
		// 	}
		// }
	})

	// Subscribe to enemy destroyed events for scoring
	g.eventSystem.SubscribeToEnemyDestroyed(func(w donburi.World, event events.EnemyDestroyedEvent) {
		g.scoreManager.AddScore(event.Points)
		dbg.Log(dbg.Event, "Score added from enemy (points=%d total=%d)", event.Points, g.scoreManager.GetScore())
	})
}

// initializeHUDPresenter creates and subscribes the HUD presenter to events
func (g *ECSGame) initializeHUDPresenter() {
	// Get initial values from game state
	current, maximum := g.healthSystem.GetPlayerHealth()
	initialHealth := 1.0
	if maximum > 0 {
		initialHealth = float64(current) / float64(maximum)
	}

	g.hudPresenter = presenter.NewHUDPresenter(
		g.scoreManager.GetScore(),
		current,
		g.levelManager.GetLevel(),
		initialHealth,
	)

	// Subscribe to events
	g.hudPresenter.Subscribe(g.eventSystem)
	dbg.Log(dbg.System, "HUD presenter initialized and subscribed to events")
}
