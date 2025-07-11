package scenes

import (
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type TitleScreenScene struct {
	manager   *SceneManager
	startTime time.Time
}

func NewTitleScreenScene(manager *SceneManager) *TitleScreenScene {
	return &TitleScreenScene{
		manager:   manager,
		startTime: time.Now(),
	}
}

func (s *TitleScreenScene) Update() error {
	// Check for any key press to continue (handled by input system)
	return nil
}

func (s *TitleScreenScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)
	drawCenteredText(screen, "GIMBAL",
		float64(s.manager.config.ScreenSize.Width)/2,
		float64(s.manager.config.ScreenSize.Height)/2-50,
		1.0)
	drawCenteredText(screen, "Exoplanetary Gyruss-Inspired Shooter",
		float64(s.manager.config.ScreenSize.Width)/2,
		float64(s.manager.config.ScreenSize.Height)/2,
		1.0)
	elapsed := time.Since(s.startTime).Seconds()
	blink := (elapsed * 2) < 1.0 // Blink every 0.5 seconds
	if blink {
		drawCenteredText(screen, "Press any key to continue",
			float64(s.manager.config.ScreenSize.Width)/2,
			float64(s.manager.config.ScreenSize.Height)/2+100,
			1.0)
	}
}

func (s *TitleScreenScene) Enter() {
	s.manager.logger.Debug("Entering title screen scene")
	s.startTime = time.Now()
}

func (s *TitleScreenScene) Exit() {
	s.manager.logger.Debug("Exiting title screen scene")
}

func (s *TitleScreenScene) GetType() SceneType {
	return SceneTitleScreen
}
