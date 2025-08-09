package mainmenu

import (
	"image/color"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	"github.com/jonesrussell/gimbal/internal/scenes"
	"github.com/jonesrussell/gimbal/internal/scenes/menu"
)

const (
	titleY     = 80
	titleScale = 1.5
)

// MenuScene manages the main menu state and rendering
type MenuScene struct {
	manager *scenes.SceneManager
	menu    *menu.MenuSystem
	font    text.Face
}

// NewMenuScene creates a new menu scene instance
func NewMenuScene(manager *scenes.SceneManager, font text.Face) *MenuScene {
	options := []menu.MenuOption{
		{Text: "Start Game", Action: func() { manager.SwitchScene(scenes.ScenePlaying) }},
		{Text: "Options", Action: func() { manager.SwitchScene(scenes.SceneOptions) }},
		{Text: "Credits", Action: func() { manager.SwitchScene(scenes.SceneCredits) }},
		{Text: "Quit", Action: func() { manager.GetLogger().Debug("Quitting game"); os.Exit(0) }},
	}
	config := menu.DefaultMenuConfig()
	config.MenuY = float64(manager.GetConfig().ScreenSize.Height) / 2
	return &MenuScene{
		manager: manager,
		menu: menu.NewMenuSystem(options, &config, manager.GetConfig().ScreenSize.Width,
			manager.GetConfig().ScreenSize.Height, font),
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
	scenes.DrawCenteredTextWithOptions(screen, scenes.TextDrawOptions{
		Text:  "GIMBAL",
		X:     float64(s.manager.GetConfig().ScreenSize.Width) / 2,
		Y:     titleY,
		Alpha: titleScale,
		Font:  s.font,
	})
}

// Enter is called when the scene becomes active
func (s *MenuScene) Enter() {
	s.manager.GetLogger().Debug("Entering menu scene")
	s.menu.Reset()
}

// Exit is called when the scene becomes inactive
func (s *MenuScene) Exit() {
	s.manager.GetLogger().Debug("Exiting menu scene")
}

// GetType returns the scene type identifier
func (s *MenuScene) GetType() scenes.SceneType {
	return scenes.SceneMenu
}
