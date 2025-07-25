package scenes

import (
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type TitleScreenScene struct {
	manager   *SceneManager
	startTime time.Time
	font      text.Face
}

func NewTitleScreenScene(manager *SceneManager, font text.Face) *TitleScreenScene {
	return &TitleScreenScene{
		manager:   manager,
		startTime: time.Now(),
		font:      font,
	}
}

func (s *TitleScreenScene) Update() error {
	// Check for any key press to continue (handled by input system)
	return nil
}

func (s *TitleScreenScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)
	drawCenteredTextWithOptions(screen, TextDrawOptions{
		Text:  "GIMBAL",
		X:     float64(s.manager.config.ScreenSize.Width) / 2,
		Y:     float64(s.manager.config.ScreenSize.Height)/2 - 50,
		Alpha: 1.0,
		Font:  s.font,
	})
	drawCenteredTextWithOptions(screen, TextDrawOptions{
		Text:  "Exoplanetary Gyruss-Inspired Shooter",
		X:     float64(s.manager.config.ScreenSize.Width) / 2,
		Y:     float64(s.manager.config.ScreenSize.Height) / 2,
		Alpha: 1.0,
		Font:  s.font,
	})
	elapsed := time.Since(s.startTime).Seconds()
	blink := (elapsed * 2) < 1.0 // Blink every 0.5 seconds
	if blink {
		drawCenteredTextWithOptions(screen, TextDrawOptions{
			Text:  "Press any key to continue",
			X:     float64(s.manager.config.ScreenSize.Width) / 2,
			Y:     float64(s.manager.config.ScreenSize.Height)/2 + 100,
			Alpha: 1.0,
			Font:  s.font,
		})
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
