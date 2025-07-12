package intro

import (
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	"github.com/jonesrussell/gimbal/internal/ecs/managers"
	resources "github.com/jonesrussell/gimbal/internal/ecs/managers/resource"
	scenes "github.com/jonesrussell/gimbal/internal/scenes"
)

func init() {
	// Register intro scenes with the scene registry
	scenes.RegisterScene(scenes.SceneStudioIntro, createStudioIntroScene)
	scenes.RegisterScene(scenes.SceneTitleScreen, createTitleScreenScene)
}

func createStudioIntroScene(
	manager *scenes.SceneManager,
	font text.Face,
	scoreManager *managers.ScoreManager,
	resourceMgr *resources.ResourceManager,
) scenes.Scene {
	return NewStudioIntroScene(manager, font)
}

func createTitleScreenScene(
	manager *scenes.SceneManager,
	font text.Face,
	scoreManager *managers.ScoreManager,
	resourceMgr *resources.ResourceManager,
) scenes.Scene {
	return NewTitleScreenScene(manager, font)
}
