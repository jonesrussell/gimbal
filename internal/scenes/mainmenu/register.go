package mainmenu

import (
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	"github.com/jonesrussell/gimbal/internal/ecs/managers"
	resources "github.com/jonesrussell/gimbal/internal/ecs/managers/resource"
	scenes "github.com/jonesrussell/gimbal/internal/scenes"
)

// Register registers the main menu scene with the scene registry.
// This should be called explicitly during application initialization.
func Register() {
	scenes.RegisterScene(scenes.SceneMenu, createMenuScene)
}

func init() {
	// Auto-register for backward compatibility
	Register()
}

func createMenuScene(
	manager *scenes.SceneManager,
	font text.Face,
	scoreManager *managers.ScoreManager,
	resourceMgr *resources.ResourceManager,
) scenes.Scene {
	return NewMenuScene(manager, font)
}
