package gameover

import (
	"fmt"
	"image/color"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	"github.com/jonesrussell/gimbal/internal/ecs/managers"
	resources "github.com/jonesrussell/gimbal/internal/ecs/managers/resource"
	"github.com/jonesrussell/gimbal/internal/scenes"
)

const (
	continueCountdownDuration = 10.0 // seconds
	continueTickInterval      = 1.0  // seconds
)

type GameOverScene struct {
	manager          *scenes.SceneManager
	font             text.Face
	resourceMgr      *resources.ResourceManager
	scoreManager     *managers.ScoreManager
	startTime        time.Time
	fadeAlpha        float64
	countdown        float64
	lastTickTime     time.Time
	continueAccepted bool
}

func NewGameOverScene(
	manager *scenes.SceneManager,
	font text.Face,
	scoreManager *managers.ScoreManager,
	resourceMgr *resources.ResourceManager,
) *GameOverScene {
	return &GameOverScene{
		manager:      manager,
		font:         font,
		resourceMgr:  resourceMgr,
		scoreManager: scoreManager,
		startTime:    time.Now(),
		fadeAlpha:    0.0,
		countdown:    continueCountdownDuration,
		lastTickTime: time.Now(),
	}
}

func (s *GameOverScene) Update() error {
	elapsed := time.Since(s.startTime).Seconds()

	// Fade in effect
	if elapsed < 0.5 {
		s.fadeAlpha = elapsed / 0.5
	} else {
		s.fadeAlpha = 1.0
	}

	// Update countdown
	if s.countdown > 0 {
		s.countdown = continueCountdownDuration - elapsed
		if s.countdown < 0 {
			s.countdown = 0
		}

		// Play tick sound every second
		timeSinceLastTick := time.Since(s.lastTickTime).Seconds()
		if timeSinceLastTick >= continueTickInterval && s.countdown > 0 {
			s.playContinueTick()
			s.lastTickTime = time.Now()
		}
	}

	// Handle continue input
	inputHandler := s.manager.GetInputHandler()
	if inputHandler.IsShootPressed() || inputHandler.IsPausePressed() {
		if s.countdown > 0 {
			// Continue game
			s.continueAccepted = true
			s.manager.SwitchScene(scenes.ScenePlaying)
			return nil
		}
	}

	// Auto-return to menu when countdown reaches 0
	if s.countdown <= 0 && !s.continueAccepted {
		s.manager.SwitchScene(scenes.SceneMenu)
		return nil
	}

	return nil
}

func (s *GameOverScene) Draw(screen *ebiten.Image) {
	// Fade in from black
	fadeOverlay := ebiten.NewImage(screen.Bounds().Dx(), screen.Bounds().Dy())
	overlayAlpha := uint8((1.0 - s.fadeAlpha) * 255)
	fadeOverlay.Fill(color.RGBA{0, 0, 0, overlayAlpha})
	screen.Fill(color.Black)

	// Draw game over text
	config := s.manager.GetConfig()
	centerX := float64(config.ScreenSize.Width) / 2
	centerY := float64(config.ScreenSize.Height) / 2

	scenes.DrawCenteredTextWithOptions(screen, scenes.TextDrawOptions{
		Text:  "GAME OVER",
		X:     centerX,
		Y:     centerY - 80,
		Alpha: s.fadeAlpha,
		Font:  s.font,
	})

	// Draw continue prompt
	if s.countdown > 0 {
		continueText := "CONTINUE?"
		scenes.DrawCenteredTextWithOptions(screen, scenes.TextDrawOptions{
			Text:  continueText,
			X:     centerX,
			Y:     centerY,
			Alpha: s.fadeAlpha,
			Font:  s.font,
		})

		// Draw countdown
		countdownText := fmt.Sprintf("%.0f", math.Ceil(s.countdown))
		// Pulse effect for countdown
		pulseAlpha := 0.7 + 0.3*math.Sin(time.Since(s.startTime).Seconds()*4.0)
		scenes.DrawCenteredTextWithOptions(screen, scenes.TextDrawOptions{
			Text:  countdownText,
			X:     centerX,
			Y:     centerY + 60,
			Alpha: s.fadeAlpha * pulseAlpha,
			Font:  s.font,
		})

		// Draw instruction
		instructionText := "Press SPACE or ESC to continue"
		scenes.DrawCenteredTextWithOptions(screen, scenes.TextDrawOptions{
			Text:  instructionText,
			X:     centerX,
			Y:     centerY + 120,
			Alpha: s.fadeAlpha * 0.6,
			Font:  s.font,
		})
	} else {
		// Draw return to menu text
		scenes.DrawCenteredTextWithOptions(screen, scenes.TextDrawOptions{
			Text:  "Returning to menu...",
			X:     centerX,
			Y:     centerY + 60,
			Alpha: s.fadeAlpha * 0.8,
			Font:  s.font,
		})
	}

	// Draw fade overlay
	screen.DrawImage(fadeOverlay, nil)
}

func (s *GameOverScene) Enter() {
	s.manager.GetLogger().Debug("Entering game over scene")
	s.startTime = time.Now()
	s.fadeAlpha = 0.0
	s.countdown = continueCountdownDuration
	s.lastTickTime = time.Now()
	s.continueAccepted = false

	// Play game over sound
	s.playGameOverSound()
}

func (s *GameOverScene) Exit() {
	s.manager.GetLogger().Debug("Exiting game over scene")
}

func (s *GameOverScene) playGameOverSound() {
	if s.resourceMgr == nil {
		return
	}

	audioPlayer := s.resourceMgr.GetAudioPlayer()
	if audioPlayer == nil {
		return
	}

	// Play game over sound effect
	s.manager.GetLogger().Debug("Game over sound should play")
}

func (s *GameOverScene) playContinueTick() {
	if s.resourceMgr == nil {
		return
	}

	audioPlayer := s.resourceMgr.GetAudioPlayer()
	if audioPlayer == nil {
		return
	}

	// Play continue tick sound effect
	s.manager.GetLogger().Debug("Continue tick sound should play")
}

func (s *GameOverScene) GetType() scenes.SceneType {
	return scenes.SceneGameOver
}
