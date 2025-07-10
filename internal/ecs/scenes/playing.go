package ecs

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/jonesrussell/gimbal/internal/ecs/core"
)

type PlayingScene struct {
	manager *SceneManager
}

func NewPlayingScene(manager *SceneManager) *PlayingScene {
	return &PlayingScene{manager: manager}
}

func (s *PlayingScene) Update() error {
	// This will be handled by the main game loop
	return nil
}

func (s *PlayingScene) Draw(screen *ebiten.Image) {
	// Clear screen
	screen.Fill(color.Black)

	// Run render system through wrapper
	renderWrapper := core.NewRenderSystemWrapper(screen)
	if err := renderWrapper.Update(s.manager.world); err != nil {
		s.manager.logger.Error("Render system failed", "error", err)
	}

	// Draw debug info if enabled
	if s.manager.config.Debug {
		s.drawDebugInfo(screen)
	}
}

// drawDebugInfo renders debug information
func (s *PlayingScene) drawDebugInfo(screen *ebiten.Image) {
	// Get player info for debug display
	players := make([]donburi.Entity, 0)
	query.NewQuery(
		filter.And(
			filter.Contains(core.PlayerTag),
			filter.Contains(core.Position),
			filter.Contains(core.Orbital),
		),
	).Each(s.manager.world, func(entry *donburi.Entry) {
		players = append(players, entry.Entity())
	})

	if len(players) > 0 {
		playerEntry := s.manager.world.Entry(players[0])
		if playerEntry.Valid() {
			pos := core.Position.Get(playerEntry)
			orb := core.Orbital.Get(playerEntry)

			// Log debug info
			s.manager.logger.Debug("Debug Info",
				"player_pos", fmt.Sprintf("(%.1f, %.1f)", pos.X, pos.Y),
				"player_angle", fmt.Sprintf("%.1f°", orb.OrbitalAngle),
				"entity_count", s.manager.world.Len(),
			)
		}
	}
}

func (s *PlayingScene) Enter() {
	s.manager.logger.Debug("Entering playing scene")
}

func (s *PlayingScene) Exit() {
	s.manager.logger.Debug("Exiting playing scene")
}

func (s *PlayingScene) GetType() SceneType {
	return ScenePlaying
}
