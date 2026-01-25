package intro

import (
	"context"
	"image/color"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/ecs/managers"
	resources "github.com/jonesrussell/gimbal/internal/ecs/managers/resource"
	scenes "github.com/jonesrussell/gimbal/internal/scenes"
	"github.com/jonesrussell/gimbal/internal/scenes/effects"
)

type TitleScreenScene struct {
	manager      *scenes.SceneManager
	startTime    time.Time
	font         text.Face
	starfield    *effects.Starfield
	resourceMgr  *resources.ResourceManager
	scoreManager *managers.ScoreManager
	titleLogo    *ebiten.Image
	musicPlaying bool
}

func NewTitleScreenScene(
	manager *scenes.SceneManager,
	font text.Face,
	scoreManager *managers.ScoreManager,
	resourceMgr *resources.ResourceManager,
) *TitleScreenScene {
	config := manager.GetConfig()
	starfield := effects.NewStarfield(
		config.ScreenSize.Width,
		config.ScreenSize.Height,
		100, // star count
		2.0, // speed
	)

	// Try to load title logo
	var titleLogo *ebiten.Image
	if resourceMgr != nil {
		if logo, ok := resourceMgr.GetSprite(context.Background(), "title_logo"); ok {
			titleLogo = logo
		}
	}

	return &TitleScreenScene{
		manager:      manager,
		startTime:    time.Now(),
		font:         font,
		starfield:    starfield,
		resourceMgr:  resourceMgr,
		scoreManager: scoreManager,
		titleLogo:    titleLogo,
	}
}

func (s *TitleScreenScene) Update() error {
	// Update starfield animation
	deltaTime := 1.0 / 60.0 // Assume 60 FPS
	s.starfield.Update(deltaTime)

	// Log input event for debugging
	event := s.manager.GetInputHandler().GetLastEvent()
	s.manager.GetLogger().Debug("TitleScreen input event", "event", event)

	// Transition on any key or mouse event
	if event == common.InputEventAny {
		s.manager.SwitchScene(scenes.SceneMenu)
	}
	return nil
}

func (s *TitleScreenScene) Draw(screen *ebiten.Image) {
	// Draw starfield background
	screen.Fill(color.Black)
	s.starfield.Draw(screen)

	config := s.manager.GetConfig()
	centerX := float64(config.ScreenSize.Width) / 2
	centerY := float64(config.ScreenSize.Height) / 2

	// Draw title logo if available, otherwise draw text
	if s.titleLogo != nil {
		op := &ebiten.DrawImageOptions{}
		logoWidth := float64(s.titleLogo.Bounds().Dx())
		logoHeight := float64(s.titleLogo.Bounds().Dy())
		op.GeoM.Translate(centerX-logoWidth/2, centerY-logoHeight/2-80)
		// Don't set ColorScale - let Ebiten use defaults which preserve source alpha
		screen.DrawImage(s.titleLogo, op)
	} else {
		// Fallback to text
		scenes.DrawCenteredTextWithOptions(screen, scenes.TextDrawOptions{
			Text:  "GIMBAL",
			X:     centerX,
			Y:     centerY - 50,
			Alpha: 1.0,
			Font:  s.font,
		})
	}

	// Draw "PRESS START" with blinking effect
	elapsed := time.Since(s.startTime).Seconds()
	blinkAlpha := (math.Sin(elapsed*4.0) + 1.0) / 2.0 // Smooth blink between 0 and 1
	if blinkAlpha > 0.3 {                             // Only show when above threshold
		scenes.DrawCenteredTextWithOptions(screen, scenes.TextDrawOptions{
			Text:  "PRESS START",
			X:     centerX,
			Y:     centerY + 100,
			Alpha: blinkAlpha,
			Font:  s.font,
		})
	}
}

func (s *TitleScreenScene) Enter() {
	s.manager.GetLogger().Debug("Entering title screen scene")
	s.startTime = time.Now()
	s.starfield.Reset()

	// Start main theme music
	s.startMusic("game_music_main")
}

func (s *TitleScreenScene) Exit() {
	s.manager.GetLogger().Debug("Exiting title screen scene")
	s.stopMusic("game_music_main")
}

// startMusic starts playing a music track
func (s *TitleScreenScene) startMusic(trackName string) {
	if s.resourceMgr == nil {
		return
	}
	audioPlayer := s.resourceMgr.GetAudioPlayer()
	if audioPlayer == nil {
		return
	}
	musicRes, ok := s.resourceMgr.GetAudio(context.Background(), trackName)
	if !ok {
		return
	}
	if err := audioPlayer.PlayMusic(trackName, musicRes, 0.7); err != nil {
		s.manager.GetLogger().Warn("Failed to play music", "track", trackName, "error", err)
		return
	}
	s.musicPlaying = true
}

// stopMusic stops playing a music track
func (s *TitleScreenScene) stopMusic(trackName string) {
	if s.resourceMgr == nil || !s.musicPlaying {
		return
	}
	audioPlayer := s.resourceMgr.GetAudioPlayer()
	if audioPlayer != nil {
		audioPlayer.StopMusic(trackName)
	}
	s.musicPlaying = false
}

func (s *TitleScreenScene) GetType() scenes.SceneType {
	return scenes.SceneTitleScreen
}
