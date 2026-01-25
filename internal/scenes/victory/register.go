package victory

import (
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	"github.com/jonesrussell/gimbal/internal/ecs/managers"
	resources "github.com/jonesrussell/gimbal/internal/ecs/managers/resource"
	"github.com/jonesrussell/gimbal/internal/scenes"
)

// Register registers the victory scene with the scene registry
func Register() {
	scenes.RegisterScene(scenes.SceneVictory, createVictoryScene)
}

func init() {
	// Auto-register for backward compatibility
	Register()
}

func createVictoryScene(
	manager *scenes.SceneManager,
	font text.Face,
	scoreManager *managers.ScoreManager,
	resourceMgr *resources.ResourceManager,
) scenes.Scene {
	return NewVictoryScene(manager, font, scoreManager, resourceMgr)
}
