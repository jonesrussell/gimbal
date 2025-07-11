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
	return &PlayingScene{
		manager:      manager,
		font:         font,
		scoreManager: scoreManager,
		resourceMgr:  resourceMgr,
	}
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

	// Draw lives display
	s.drawLivesDisplay(screen)

	// Draw score display
	s.drawScore(screen)

	// Draw debug info if enabled
	if s.manager.config.Debug {
		s.drawDebugInfo(screen)
	}
}

func (s *PlayingScene) drawScore(screen *ebiten.Image) {
	score := s.scoreManager.GetScore()
	scoreText := fmt.Sprintf("Score: %d", score)

	op := &v2text.DrawOptions{}
	// Position: top-right, 150px from right, 30px from top
	w := screen.Bounds().Dx()
	op.GeoM.Translate(float64(w-150), 30)
	op.ColorScale.SetR(1)
	op.ColorScale.SetG(1)
	op.ColorScale.SetB(1)
	op.ColorScale.SetA(1)
	v2text.Draw(screen, scoreText, s.font, op)
}

// TriggerScreenShake triggers a screen shake effect
func (s *PlayingScene) TriggerScreenShake() {
	s.screenShake = 1.0 // Set shake intensity
}

// drawLivesDisplay renders the player's remaining lives in the top-left corner
func (s *PlayingScene) drawLivesDisplay(screen *ebiten.Image) {
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

		// Get heart sprite from resource manager
		heartSprite, exists := s.resourceMgr.GetSprite("heart")
		if !exists {
			// Fallback to text if sprite not found
			s.drawLivesText(screen, current, maximum)
			return
		}

		// Draw heart sprites
		s.drawLivesSprites(screen, heartSprite, current, maximum)
	}
}

// drawLivesText draws lives as text (fallback method)
func (s *PlayingScene) drawLivesText(screen *ebiten.Image, current, maximum int) {
	// Create lives text with heart emojis
	livesText := "Lives: "
	for i := 0; i < maximum; i++ {
		if i < current {
			livesText += "❤️"
		} else {
			livesText += "🖤" // Empty heart
		}
	}

	// Draw the text in the top-left corner
	op := &v2text.DrawOptions{}
	op.GeoM.Translate(20, 30)
	op.ColorScale.SetR(1)
	op.ColorScale.SetG(1)
	op.ColorScale.SetB(1)
	op.ColorScale.SetA(1)
	v2text.Draw(screen, livesText, s.font, op)
}

// drawLivesSprites draws lives as heart sprites
func (s *PlayingScene) drawLivesSprites(screen, heartSprite *ebiten.Image, current, maximum int) {
	// Draw "Lives:" text first
	label := "Lives: "
	op := &v2text.DrawOptions{}
	op.GeoM.Translate(20, 30)
	op.ColorScale.SetR(1)
	op.ColorScale.SetG(1)
	op.ColorScale.SetB(1)
	op.ColorScale.SetA(1)
	v2text.Draw(screen, label, s.font, op)

	// Measure the actual text width using text/v2
	labelWidth, _ := v2text.Measure(label, s.font, 0)

	heartSize := 96 // Desired heart sprite size in pixels (increased from 32)
	spacing := 16   // Space between hearts (increased from 8)

	// Get original sprite size
	origW, origH := heartSprite.Bounds().Dx(), heartSprite.Bounds().Dy()
	scaleX := float64(heartSize) / float64(origW)
	scaleY := float64(heartSize) / float64(origH)

	// Debug logging
	s.manager.logger.Debug("Heart positioning",
		"label_width", labelWidth,
		"heart_size", heartSize,
		"orig_sprite_size", fmt.Sprintf("%dx%d", origW, origH),
		"scale", fmt.Sprintf("%.3fx%.3f", scaleX, scaleY),
	)

	// Draw heart sprites with proper spacing
	heartX := 20 + labelWidth
	for i := 0; i < maximum; i++ {
		heartOp := &ebiten.DrawImageOptions{}
		x := heartX + float64(i*(heartSize+spacing))
		y := float64(30 - heartSize/2) // Center vertically with text
		heartOp.GeoM.Translate(x, y)
		heartOp.GeoM.Scale(scaleX, scaleY)

		// Debug logging for each heart
		s.manager.logger.Debug("Drawing heart",
			"heart_index", i,
			"current_lives", current,
			"max_lives", maximum,
			"position", fmt.Sprintf("(%.1f, %.1f)", x, y),
			"is_full", i < current,
		)

		if i < current {
			// Full heart
			heartOp.ColorScale.SetR(1)
			heartOp.ColorScale.SetG(0)
			heartOp.ColorScale.SetB(0)
			heartOp.ColorScale.SetA(1)
		} else {
			// Empty heart (grayed out)
			heartOp.ColorScale.SetR(0.5)
			heartOp.ColorScale.SetG(0.5)
			heartOp.ColorScale.SetB(0.5)
			heartOp.ColorScale.SetA(0.5)
		}

		screen.DrawImage(heartSprite, heartOp)
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
