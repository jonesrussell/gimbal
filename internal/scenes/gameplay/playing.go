package gameplay

import (
	"fmt"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	v2text "github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/jonesrussell/gimbal/internal/ecs/core"
	"github.com/jonesrussell/gimbal/internal/ecs/managers"
	resources "github.com/jonesrussell/gimbal/internal/ecs/managers/resource"
	"github.com/jonesrussell/gimbal/internal/scenes"
)

type PlayingScene struct {
	manager      *scenes.SceneManager
	screenShake  float64 // Screen shake intensity (0 = no shake)
	font         v2text.Face
	scoreManager *managers.ScoreManager
	resourceMgr  *resources.ResourceManager

	// Level title display
	levelTitleStartTime time.Time
	showLevelTitle      bool
	currentLevelNumber  int
	levelTitleDuration  float64 // Duration to show title in seconds
}

func NewPlayingScene(
	manager *scenes.SceneManager,
	font v2text.Face,
	scoreManager *managers.ScoreManager,
	resourceMgr *resources.ResourceManager,
) *PlayingScene {
	scene := &PlayingScene{
		manager:            manager,
		font:               font,
		scoreManager:       scoreManager,
		resourceMgr:        resourceMgr,
		levelTitleDuration: 3.0, // Show title for 3 seconds
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

	// Update level title display
	if s.showLevelTitle {
		elapsed := time.Since(s.levelTitleStartTime).Seconds()
		if elapsed >= s.levelTitleDuration {
			s.showLevelTitle = false
		}
	}

	return nil
}

func (s *PlayingScene) Draw(screen *ebiten.Image) {
	// Clear screen
	screen.Fill(color.Black)

	s.manager.GetLogger().Debug("PlayingScene.Draw called", "screen_size", screen.Bounds())

	// Apply screen shake if active
	if s.screenShake > 0 {
		// Get image from pool instead of creating new one
		shakenImage := s.manager.GetImagePool().GetImage(screen.Bounds().Dx(), screen.Bounds().Dy())
		defer s.manager.GetImagePool().ReturnImage(shakenImage)

		s.drawGameContent(shakenImage)
		op := &ebiten.DrawImageOptions{}
		shakeOffset := s.screenShake * 5
		op.GeoM.Translate(shakeOffset, shakeOffset)
		screen.DrawImage(shakenImage, op)
	} else {
		s.drawGameContent(screen)
	}

	// Draw level title overlay if showing
	if s.showLevelTitle {
		s.drawLevelTitle(screen)
	}
}

func (s *PlayingScene) drawGameContent(screen *ebiten.Image) {
	s.manager.GetLogger().Debug("drawGameContent called", "screen_size", screen.Bounds())

	// Use optimized render system if available
	if renderOptimizer := s.manager.GetRenderOptimizer(); renderOptimizer != nil {
		renderOptimizer.OptimizedRenderSystem(s.manager.GetWorld(), screen)
	} else {
		// Fallback to original render system
		renderWrapper := core.NewRenderSystemWrapper(screen)
		if err := renderWrapper.Update(s.manager.GetWorld()); err != nil {
			s.manager.GetLogger().Error("Render system failed", "error", err)
		}
	}

	if s.manager.GetConfig().Debug {
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
	).Each(s.manager.GetWorld(), func(entry *donburi.Entry) {
		players = append(players, entry.Entity())
	})

	if len(players) > 0 {
		playerEntry := s.manager.GetWorld().Entry(players[0])
		if playerEntry.Valid() {
			pos := core.Position.Get(playerEntry)
			orb := core.Orbital.Get(playerEntry)

			// Log debug info
			s.manager.GetLogger().Debug("Debug Info",
				"player_pos", fmt.Sprintf("(%.1f, %.1f)", pos.X, pos.Y),
				"player_angle", fmt.Sprintf("%.1f°", orb.OrbitalAngle),
				"entity_count", s.manager.GetWorld().Len(),
			)
		}
	}
}

func (s *PlayingScene) Enter() {
	s.manager.GetLogger().Debug("Entering playing scene")

	// Show level title when entering playing scene
	if levelManager := s.manager.GetLevelManager(); levelManager != nil {
		if levelConfig := levelManager.GetCurrentLevelConfig(); levelConfig != nil {
			s.ShowLevelTitle(levelConfig.LevelNumber)
		}
	}
}

func (s *PlayingScene) Exit() {
	s.manager.GetLogger().Debug("Exiting playing scene")
	s.showLevelTitle = false
}

// ShowLevelTitle displays the level title overlay
func (s *PlayingScene) ShowLevelTitle(levelNumber int) {
	s.currentLevelNumber = levelNumber
	s.levelTitleStartTime = time.Now()
	s.showLevelTitle = true
	s.manager.GetLogger().Debug("Level title shown", "level", levelNumber)
}

// drawLevelTitle draws the level title overlay
func (s *PlayingScene) drawLevelTitle(screen *ebiten.Image) {
	if s.font == nil {
		return
	}

	elapsed := time.Since(s.levelTitleStartTime).Seconds()
	alpha := s.calculateTitleAlpha(elapsed)
	if alpha <= 0 {
		return
	}

	s.drawTitleOverlay(screen, alpha)
	titleText, descText := s.getTitleText()
	s.drawTitleText(screen, titleText, descText, alpha)
}

// calculateTitleAlpha calculates fade alpha (fade in for first 0.5s, fade out for last 0.5s)
func (s *PlayingScene) calculateTitleAlpha(elapsed float64) float64 {
	alpha := 1.0
	if elapsed < 0.5 {
		alpha = elapsed / 0.5 // Fade in
	} else if elapsed > s.levelTitleDuration-0.5 {
		alpha = (s.levelTitleDuration - elapsed) / 0.5 // Fade out
	}
	return alpha
}

// drawTitleOverlay draws the semi-transparent background overlay
func (s *PlayingScene) drawTitleOverlay(screen *ebiten.Image, alpha float64) {
	bgColor := color.RGBA{0, 0, 0, uint8(200 * alpha)}
	overlay := ebiten.NewImage(screen.Bounds().Dx(), screen.Bounds().Dy())
	overlay.Fill(bgColor)

	op := &ebiten.DrawImageOptions{}
	op.ColorScale.SetA(float32(alpha))
	screen.DrawImage(overlay, op)
}

// getTitleText gets title and description text from level config
func (s *PlayingScene) getTitleText() (titleText, descText string) {
	levelManager := s.manager.GetLevelManager()
	if levelManager == nil {
		titleText = fmt.Sprintf("LEVEL %d", s.currentLevelNumber)
		return titleText, descText
	}

	levelConfig := levelManager.GetCurrentLevelConfig()
	if levelConfig == nil {
		titleText = fmt.Sprintf("LEVEL %d", s.currentLevelNumber)
		return titleText, descText
	}

	titleText = levelConfig.Metadata.Name
	if titleText == "" {
		titleText = fmt.Sprintf("LEVEL %d", levelConfig.LevelNumber)
	} else {
		titleText = fmt.Sprintf("LEVEL %d: %s", levelConfig.LevelNumber, titleText)
	}
	descText = levelConfig.Metadata.Description
	return titleText, descText
}

// drawTitleText draws the title and description text
func (s *PlayingScene) drawTitleText(screen *ebiten.Image, titleText, descText string, alpha float64) {
	titleWidth, titleHeight := v2text.Measure(titleText, s.font, 0)
	screenWidth := float64(s.manager.GetConfig().ScreenSize.Width)
	screenHeight := float64(s.manager.GetConfig().ScreenSize.Height)

	titleX := (screenWidth - float64(titleWidth)) / 2
	titleY := screenHeight/2 - float64(titleHeight) - 30

	// Draw title
	titleOp := &v2text.DrawOptions{}
	titleOp.GeoM.Translate(titleX, titleY)
	titleOp.ColorScale.SetA(float32(alpha))
	v2text.Draw(screen, titleText, s.font, titleOp)

	// Draw description if available
	if descText != "" {
		descWidth, _ := v2text.Measure(descText, s.font, 0)
		descX := (screenWidth - float64(descWidth)) / 2
		descY := titleY + float64(titleHeight) + 20

		descOp := &v2text.DrawOptions{}
		descOp.GeoM.Translate(descX, descY)
		descOp.ColorScale.SetA(float32(alpha * 0.8)) // Slightly more transparent
		v2text.Draw(screen, descText, s.font, descOp)
	}
}

func (s *PlayingScene) GetType() scenes.SceneType {
	return scenes.ScenePlaying
}
