package stagetransition

import (
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	"github.com/jonesrussell/gimbal/internal/ecs/managers"
	resources "github.com/jonesrussell/gimbal/internal/ecs/managers/resource"
	"github.com/jonesrussell/gimbal/internal/scenes"
)

// Register registers the stage transition scene with the scene registry
func Register() {
	scenes.RegisterScene(scenes.SceneStageTransition, createStageTransitionScene)
}

func init() {
	// Auto-register for backward compatibility
	Register()
}

func createStageTransitionScene(
	manager *scenes.SceneManager,
	font text.Face,
	scoreManager *managers.ScoreManager,
	resourceMgr *resources.ResourceManager,
) scenes.Scene {
	return NewStageTransitionScene(manager, font, scoreManager, resourceMgr)
}
