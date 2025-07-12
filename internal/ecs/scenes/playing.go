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
	"github.com/jonesrussell/gimbal/internal/ecs/resources"
	"github.com/jonesrussell/gimbal/internal/ecs/ui"
)

type PlayingScene struct {
	manager      *SceneManager
	screenShake  float64 // Screen shake intensity (0 = no shake)
	font         v2text.Face
	scoreManager *managers.ScoreManager
	resourceMgr  *resources.ResourceManager
	uiRenderer   *ui.UIRenderer
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

	// Initialize UI renderer with default theme
	scene.uiRenderer = ui.NewUIRenderer(nil, ui.DefaultTheme)

	// Set up theme fonts
	ui.DefaultTheme.SetFonts(font, font, font)

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
	// Set the screen for the UI renderer
	s.uiRenderer.SetScreen(screen)

	// Set up heart sprite for lives display
	if heartSprite, exists := s.resourceMgr.GetSprite("heart"); exists {
		s.uiRenderer.SetHeartSprite(heartSprite)
	}

	// Enable debug mode if configured
	s.uiRenderer.SetDebug(s.manager.config.Debug)

	// Run render system through wrapper
	renderWrapper := core.NewRenderSystemWrapper(screen)
	if err := renderWrapper.Update(s.manager.world); err != nil {
		s.manager.logger.Error("Render system failed", "error", err)
	}

	// Draw UI elements using the new renderer
	s.drawUIElements()

	// Draw debug info if enabled
	if s.manager.config.Debug {
		s.drawDebugInfo(screen)
	}
}

// drawUIElements draws all UI elements using the new renderer system
func (s *PlayingScene) drawUIElements() {
	// Get health system from scene manager
	healthSystem := s.manager.GetHealthSystem()
	if healthSystem == nil {
		return
	}

	// Type assert to get current and maximum lives
	if hs, ok := healthSystem.(interface {
		GetPlayerHealth() (int, int)
	}); ok {
		current, maximum := hs.GetPlayerHealth()

		// Debug logging for health values
		s.manager.logger.Debug("Health system values",
			"current_lives", current,
			"maximum_lives", maximum,
		)

		// Draw lives display using UI renderer
		livesDisplay := ui.NewLivesDisplay(current, maximum)
		s.uiRenderer.Draw(livesDisplay, ui.TopLeft(20, 20))
	}

	// Draw score display using UI renderer
	score := s.scoreManager.GetScore()
	scoreDisplay := ui.NewScoreDisplay(score)

	// Get screen width for proper top-right positioning
	screenWidth := s.uiRenderer.GetScreenWidth()
	s.uiRenderer.Draw(scoreDisplay, ui.TopRightRelative(screenWidth, 150, 20))
}

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
