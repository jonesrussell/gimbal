package pause

import (
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	"github.com/jonesrussell/gimbal/internal/scenes"
	"github.com/jonesrussell/gimbal/internal/scenes/menu"
)

// PauseOption represents the available pause menu options
type PauseOption int

const (
	PauseOptionResume PauseOption = iota
	PauseOptionReturnToMenu
	PauseOptionQuit
)

const (
	// Animation constants
	frameRate      = 60.0
	fadeInDuration = 0.3
	selectionDelay = 0.1

	// Layout constants
	titleScale     = 1.5
	titleY         = 80
	menuSpacing    = 50
	overlayAlpha   = 128
	paddingX       = 30
	paddingY       = 10
	hitAreaPadding = 50
	hintTextY      = 60

	// Color constants
	dimmedAlpha   = 0.7
	hintBaseAlpha = 0.6
)

// PausedScene manages the pause menu state and rendering
type PausedScene struct {
	manager           *scenes.SceneManager
	menu              *menu.MenuSystem
	overlayImage      *ebiten.Image
	animationTime     float64
	fadeIn            float64
	lastSelectionTime time.Time
	selectionChanged  bool
	escWasPressed     bool
	canUnpause        bool
	font              text.Face
}

// NewPausedScene creates a new pause scene instance
func NewPausedScene(manager *scenes.SceneManager, font text.Face) *PausedScene {
	scene := &PausedScene{
		manager:           manager,
		animationTime:     0,
		selectionChanged:  false,
		lastSelectionTime: time.Now(),
		fadeIn:            0,
		font:              font,
	}

	options := []menu.MenuOption{
		{Text: "Resume", Action: func() {
			// Resume callback is handled by the game when switching to playing scene
			manager.SwitchScene(scenes.ScenePlaying)
		}},
		{Text: "Return to Menu", Action: func() { manager.SwitchScene(scenes.SceneMenu) }},
		{Text: "Quit", Action: func() { os.Exit(0) }},
	}

	config := menu.PausedMenuConfig()
	config.MenuY = float64(manager.GetConfig().ScreenSize.Height) / 2

	scene.menu = menu.NewMenuSystem(options, &config, manager.GetConfig().ScreenSize.Width,
		manager.GetConfig().ScreenSize.Height, font)

	// Create overlay image once (TODO: handle resizing if needed)
	scene.overlayImage = ebiten.NewImage(manager.GetConfig().ScreenSize.Width, manager.GetConfig().ScreenSize.Height)
	return scene
}

// Update handles input and animations for the pause scene
func (s *PausedScene) Update() error {
	dt := 1.0 / frameRate
	s.animationTime += dt
	s.updateFadeIn(dt)
	s.updateSelectionAnimation()
	s.handleInput()
	s.menu.Update()
	return nil
}

// Draw renders the pause menu
func (s *PausedScene) Draw(screen *ebiten.Image) {
	s.drawOverlay(screen)
	s.drawTitle(screen)
	s.menu.Draw(screen, s.fadeIn)
	s.drawHintText(screen)
}

// Enter is called when the scene becomes active
func (s *PausedScene) Enter() {
	s.manager.GetLogger().Debug("Entering paused scene")
	s.fadeIn = 0
	s.animationTime = 0
	s.selectionChanged = false

	// Check ESC state when entering
	s.escWasPressed = ebiten.IsKeyPressed(ebiten.KeyEscape)
	s.canUnpause = false
}

// Exit is called when the scene becomes inactive
func (s *PausedScene) Exit() {
	s.manager.GetLogger().Debug("Exiting paused scene")
}

// GetType returns the scene type identifier
func (s *PausedScene) GetType() scenes.SceneType {
	return scenes.ScenePaused
}
