package intro

import (
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	"github.com/jonesrussell/gimbal/internal/ecs/managers"
	resources "github.com/jonesrussell/gimbal/internal/ecs/managers/resource"
	scenes "github.com/jonesrussell/gimbal/internal/scenes"
)

// Register registers the intro scenes with the scene registry.
// This should be called explicitly during application initialization.
func Register() {
	scenes.RegisterScene(scenes.SceneStudioIntro, createStudioIntroScene)
	scenes.RegisterScene(scenes.SceneTitleScreen, createTitleScreenScene)
}

func init() {
	// Auto-register for backward compatibility
	Register()
}

func createStudioIntroScene(
	manager *scenes.SceneManager,
	font text.Face,
	scoreManager *managers.ScoreManager,
	resourceMgr *resources.ResourceManager,
) scenes.Scene {
	return NewStudioIntroScene(manager, resourceMgr)
}

func createTitleScreenScene(
	manager *scenes.SceneManager,
	font text.Face,
	scoreManager *managers.ScoreManager,
	resourceMgr *resources.ResourceManager,
) scenes.Scene {
	return NewTitleScreenScene(manager, font, scoreManager, resourceMgr)
}
