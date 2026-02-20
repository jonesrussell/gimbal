// Package mainmenu provides the main menu scene for the game.
// It handles menu navigation, option selection, and scene transitions.
package mainmenu

import (
	"context"
	"image/color"
	"log"

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
		{Text: "Start Game", Action: func() {
			// Start with stage 1 intro - set info after switching
			manager.SwitchScene(scenes.SceneStageIntro)
			// Set stage info on the new scene
			if stageIntroScene, ok := manager.GetCurrentScene().(interface {
				SetStageInfo(int, string, string, string)
			}); ok {
				stageIntroScene.SetStageInfo(1, "Earth", "Mars", "ENEMY ACTIVITY DETECTED")
			}
		}},
		{Text: "Options", Action: func() { manager.SwitchScene(scenes.SceneOptions) }},
		{Text: "Credits", Action: func() { manager.SwitchScene(scenes.SceneCredits) }},
		{Text: "Quit", Action: func() { manager.RequestQuit() }},
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
	s.menu.Reset()
	s.startMenuMusic()
}

// Exit is called when the scene becomes inactive
func (s *MenuScene) Exit() {
	s.stopMenuMusic()
}

// startMenuMusic starts playing the menu background music
func (s *MenuScene) startMenuMusic() {
	resourceMgr := s.manager.GetResourceManager()
	if resourceMgr == nil {
		return
	}

	audioPlayer := resourceMgr.GetAudioPlayer()
	if audioPlayer == nil {
		return
	}

	musicRes, ok := resourceMgr.GetAudio(context.Background(), "game_music_main")
	if !ok {
		return
	}

	if err := audioPlayer.PlayMusic("game_music_main", musicRes, 0.7); err != nil {
		log.Printf("[WARN] Failed to play menu music: error=%v", err)
	}
}

// stopMenuMusic stops the menu background music
func (s *MenuScene) stopMenuMusic() {
	resourceMgr := s.manager.GetResourceManager()
	if resourceMgr == nil {
		return
	}

	audioPlayer := resourceMgr.GetAudioPlayer()
	if audioPlayer == nil {
		return
	}

	audioPlayer.StopMusic("game_music_main")
}

// GetType returns the scene type identifier
func (s *MenuScene) GetType() scenes.SceneType {
	return scenes.SceneMenu
}
