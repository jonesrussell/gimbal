package scenes

import (
	"image/color"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	"github.com/jonesrussell/gimbal/internal/scenes/menu"
)

// MenuScene manages the main menu state and rendering
type MenuScene struct {
	manager *SceneManager
	menu    *menu.MenuSystem
	font    text.Face
}

// NewMenuScene creates a new menu scene instance
func NewMenuScene(manager *SceneManager, font text.Face) *MenuScene {
	options := []menu.MenuOption{
		{Text: "Start Game", Action: func() { manager.SwitchScene(ScenePlaying) }},
		{Text: "Options", Action: func() { manager.SwitchScene(SceneOptions) }},
		{Text: "Credits", Action: func() { manager.SwitchScene(SceneCredits) }},
		{Text: "Quit", Action: func() { manager.logger.Debug("Quitting game"); os.Exit(0) }},
	}
	config := menu.DefaultMenuConfig()
	config.MenuY = float64(manager.config.ScreenSize.Height) / 2
	return &MenuScene{
		manager: manager,
		menu: menu.NewMenuSystem(options, &config, manager.config.ScreenSize.Width,
			manager.config.ScreenSize.Height, font),
		font: font,
	}
}

// Update handles input and animations for the menu scene
func (s *MenuScene) Update() error {
	s.menu.Update()
	return nil
}

// Draw renders the main menu
func (s *MenuScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)
	s.drawTitle(screen)
	s.menu.Draw(screen, 1.0)
}

// drawTitle renders the game title
func (s *MenuScene) drawTitle(screen *ebiten.Image) {
	drawCenteredTextWithOptions(screen, TextDrawOptions{
		Text:  "GIMBAL",
		X:     float64(s.manager.config.ScreenSize.Width) / 2,
		Y:     titleY,
		Alpha: titleScale,
		Font:  s.font,
	})
}

// Enter is called when the scene becomes active
func (s *MenuScene) Enter() {
	s.manager.logger.Debug("Entering menu scene")
	s.menu.Reset()
}

// Exit is called when the scene becomes inactive
func (s *MenuScene) Exit() {
	s.manager.logger.Debug("Exiting menu scene")
}

// GetType returns the scene type identifier
func (s *MenuScene) GetType() SceneType {
	return SceneMenu
}
