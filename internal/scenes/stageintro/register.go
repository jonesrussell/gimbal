package stageintro

import (
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	"github.com/jonesrussell/gimbal/internal/ecs/managers"
	resources "github.com/jonesrussell/gimbal/internal/ecs/managers/resource"
	"github.com/jonesrussell/gimbal/internal/scenes"
)

// Register registers the stage intro scene with the scene registry
func Register() {
	scenes.RegisterScene(scenes.SceneStageIntro, createStageIntroScene)
}

func init() {
	// Auto-register for backward compatibility
	Register()
}

func createStageIntroScene(
	manager *scenes.SceneManager,
	font text.Face,
	scoreManager *managers.ScoreManager,
	resourceMgr *resources.ResourceManager,
) scenes.Scene {
	return NewStageIntroScene(manager, font, scoreManager, resourceMgr)
}
