package scenes

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	v2text "github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/jonesrussell/gimbal/internal/ecs/core"
	"github.com/jonesrussell/gimbal/internal/ecs/managers"
	resources "github.com/jonesrussell/gimbal/internal/ecs/managers/resource"
)

type PlayingScene struct {
	manager      *SceneManager
	screenShake  float64 // Screen shake intensity (0 = no shake)
	font         v2text.Face
	scoreManager *managers.ScoreManager
	resourceMgr  *resources.ResourceManager
}

func NewPlayingScene(
	manager *SceneManager,
	font v2text.Face,
	scoreManager *managers.ScoreManager,
	resourceMgr *resources.ResourceManager,
) *PlayingScene {
	scene := &PlayingScene{
		manager:      manager,
		font:         font,
		scoreManager: scoreManager,
		resourceMgr:  resourceMgr,
	}

	// UI is now handled by the main game's EbitenUI system

	return scene
}

func (s *PlayingScene) Update() error {
	// Update screen shake
	if s.screenShake > 0 {
		s.screenShake -= 0.1 // Reduce shake over time
		if s.screenShake < 0 {
			s.screenShake = 0
		}
	}
	return nil
}

func (s *PlayingScene) Draw(screen *ebiten.Image) {
	// Clear screen
	screen.Fill(color.Black)

	// Apply screen shake if active
	if s.screenShake > 0 {
		// Create a temporary image for the shaken content
		shakenImage := ebiten.NewImage(screen.Bounds().Dx(), screen.Bounds().Dy())

		// Draw everything to the shaken image
		s.drawGameContent(shakenImage)

		// Apply shake offset when drawing to screen
		op := &ebiten.DrawImageOptions{}
		shakeOffset := s.screenShake * 5 // Scale shake intensity
		op.GeoM.Translate(shakeOffset, shakeOffset)
		screen.DrawImage(shakenImage, op)
	} else {
		// Draw normally without shake
		s.drawGameContent(screen)
	}
}

// drawGameContent draws the main game content (separated for screen shake)
func (s *PlayingScene) drawGameContent(screen *ebiten.Image) {
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

// UI elements are now handled by the main game's EbitenUI system

// TriggerScreenShake triggers a screen shake effect
func (s *PlayingScene) TriggerScreenShake() {
	s.screenShake = 1.0 // Set shake intensity
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
