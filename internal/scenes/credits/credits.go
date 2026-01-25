package credits

import (
	"context"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/ecs/managers"
	resources "github.com/jonesrussell/gimbal/internal/ecs/managers/resource"
	"github.com/jonesrussell/gimbal/internal/scenes"
	"github.com/jonesrussell/gimbal/internal/scenes/effects"
)

const (
	scrollSpeed   = 30.0 // pixels per second
	lineSpacing   = 40.0
	creditsStartY = 600.0  // Start below screen
	creditsEndY   = -100.0 // End above screen
)

var creditsText = []string{
	"GIMBAL",
	"",
	"Created by",
	"Gimbal Studios",
	"",
	"Programming",
	"[Your Name]",
	"",
	"Music",
	"[Composer Name]",
	"",
	"Art",
	"[Artist Name]",
	"",
	"Special Thanks",
	"To all playtesters",
	"",
	"2025",
	"",
	"Thank you for playing!",
}

type CreditsScene struct {
	manager      *scenes.SceneManager
	font         text.Face
	resourceMgr  *resources.ResourceManager
	scoreManager *managers.ScoreManager
	startTime    time.Time
	starfield    *effects.Starfield
	scrollY      float64
	musicPlaying bool
}

func NewCreditsScene(
	manager *scenes.SceneManager,
	font text.Face,
	scoreManager *managers.ScoreManager,
	resourceMgr *resources.ResourceManager,
) *CreditsScene {
	config := manager.GetConfig()
	starfield := effects.NewStarfield(
		config.ScreenSize.Width,
		config.ScreenSize.Height,
		100, // star count
		1.0, // speed
	)

	return &CreditsScene{
		manager:      manager,
		font:         font,
		resourceMgr:  resourceMgr,
		scoreManager: scoreManager,
		startTime:    time.Now(),
		starfield:    starfield,
		scrollY:      creditsStartY,
	}
}

func (s *CreditsScene) Update() error {
	deltaTime := 1.0 / 60.0

	// Update starfield
	s.starfield.Update(deltaTime)

	// Update scroll position
	s.scrollY -= scrollSpeed * deltaTime

	// Check if credits have scrolled past the top
	totalHeight := float64(len(creditsText)) * lineSpacing
	if s.scrollY+totalHeight < creditsEndY {
		// Loop back to start or return to title
		s.manager.SwitchScene(scenes.SceneTitleScreen)
		return nil
	}

	// Handle input to skip credits
	event := s.manager.GetInputHandler().GetLastEvent()
	if event != common.InputEventNone {
		s.manager.SwitchScene(scenes.SceneTitleScreen)
		return nil
	}

	return nil
}

func (s *CreditsScene) Draw(screen *ebiten.Image) {
	config := s.manager.GetConfig()
	centerX := float64(config.ScreenSize.Width) / 2

	// Draw starfield background
	screen.Fill(color.Black)
	s.starfield.Draw(screen)

	// Draw scrolling credits text
	currentY := s.scrollY
	for _, line := range creditsText {
		if line == "" {
			currentY += lineSpacing
			continue
		}

		// Only draw if on screen
		if currentY > -50 && currentY < float64(config.ScreenSize.Height)+50 {
			alpha := 1.0
			// Fade in/out at edges
			if currentY < 50 {
				alpha = currentY / 50.0
			} else if currentY > float64(config.ScreenSize.Height)-50 {
				alpha = (float64(config.ScreenSize.Height) - currentY) / 50.0
			}

			// Special styling for title
			if line == "GIMBAL" {
				scenes.DrawCenteredTextWithOptions(screen, scenes.TextDrawOptions{
					Text:  line,
					X:     centerX,
					Y:     currentY,
					Alpha: alpha,
					Font:  s.font,
				})
			} else {
				scenes.DrawCenteredTextWithOptions(screen, scenes.TextDrawOptions{
					Text:  line,
					X:     centerX,
					Y:     currentY,
					Alpha: alpha * 0.9,
					Font:  s.font,
				})
			}
		}

		currentY += lineSpacing
	}
}

func (s *CreditsScene) Enter() {
	s.manager.GetLogger().Debug("Entering credits scene")
	s.startTime = time.Now()
	s.scrollY = creditsStartY
	s.starfield.Reset()

	// Start main theme music
	if s.resourceMgr != nil {
		audioPlayer := s.resourceMgr.GetAudioPlayer()
		if audioPlayer != nil {
			if musicRes, ok := s.resourceMgr.GetAudio(context.Background(), "game_music_main"); ok {
				if err := audioPlayer.PlayMusic("game_music_main", musicRes, 0.7); err != nil {
					s.manager.GetLogger().Warn("Failed to play credits music", "error", err)
				} else {
					s.musicPlaying = true
				}
			}
		}
	}
}

func (s *CreditsScene) Exit() {
	s.manager.GetLogger().Debug("Exiting credits scene")
	// Stop music if playing
	if s.resourceMgr != nil && s.musicPlaying {
		audioPlayer := s.resourceMgr.GetAudioPlayer()
		if audioPlayer != nil {
			audioPlayer.StopMusic("game_music_main")
		}
		s.musicPlaying = false
	}
}

func (s *CreditsScene) GetType() scenes.SceneType {
	return scenes.SceneCredits
}
