package intro

import (
	"fmt"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	"github.com/jonesrussell/gimbal/internal/common"
	scenes "github.com/jonesrussell/gimbal/internal/scenes"
)

type TitleScreenScene struct {
	manager   *scenes.SceneManager
	startTime time.Time
	font      text.Face
}

func NewTitleScreenScene(manager *scenes.SceneManager, font text.Face) *TitleScreenScene {
	return &TitleScreenScene{
		manager:   manager,
		startTime: time.Now(),
		font:      font,
	}
}

func (s *TitleScreenScene) Update() error {
	// Log input event for debugging
	event := s.manager.GetInputHandler().GetLastEvent()
	s.manager.GetLogger().Debug("TitleScreen input event", "event", event)

	// Transition on any key or mouse event
	if event == common.InputEventAny {
		s.manager.SwitchScene(scenes.SceneMenu) // Or scenes.ScenePlaying if you want to go straight to gameplay
	}
	return nil
}

func (s *TitleScreenScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)
	scenes.DrawCenteredTextWithOptions(screen, scenes.TextDrawOptions{
		Text:  "GIMBAL",
		X:     float64(s.manager.GetConfig().ScreenSize.Width) / 2,
		Y:     float64(s.manager.GetConfig().ScreenSize.Height)/2 - 50,
		Alpha: 1.0,
		Font:  s.font,
	})
	scenes.DrawCenteredTextWithOptions(screen, scenes.TextDrawOptions{
		Text:  "Exoplanetary Gyruss-Inspired Shooter",
		X:     float64(s.manager.GetConfig().ScreenSize.Width) / 2,
		Y:     float64(s.manager.GetConfig().ScreenSize.Height) / 2,
		Alpha: 1.0,
		Font:  s.font,
	})
	elapsed := time.Since(s.startTime).Seconds()
	blink := (elapsed * 2) < 1.0 // Blink every 0.5 seconds
	if blink {
		scenes.DrawCenteredTextWithOptions(screen, scenes.TextDrawOptions{
			Text:  "Press any key to continue",
			X:     float64(s.manager.GetConfig().ScreenSize.Width) / 2,
			Y:     float64(s.manager.GetConfig().ScreenSize.Height)/2 + 100,
			Alpha: 1.0,
			Font:  s.font,
		})
	}
	// Draw debug info at the bottom
	debugText := fmt.Sprintf(
		"Resolution: %dx%d | TPS: %.1f",
		s.manager.GetConfig().ScreenSize.Width,
		s.manager.GetConfig().ScreenSize.Height,
		ebiten.ActualTPS(),
	)
	scenes.DrawCenteredTextWithOptions(screen, scenes.TextDrawOptions{
		Text:  debugText,
		X:     float64(s.manager.GetConfig().ScreenSize.Width) / 2,
		Y:     float64(s.manager.GetConfig().ScreenSize.Height) - 30,
		Alpha: 0.5,
		Font:  s.font,
	})
}

func (s *TitleScreenScene) Enter() {
	s.manager.GetLogger().Debug("Entering title screen scene")
	s.startTime = time.Now()
	// No music on title screen - music starts in menu
}

func (s *TitleScreenScene) Exit() {
	s.manager.GetLogger().Debug("Exiting title screen scene")
	// No music to stop - title screen doesn't play music
}

func (s *TitleScreenScene) GetType() scenes.SceneType {
	return scenes.SceneTitleScreen
}
