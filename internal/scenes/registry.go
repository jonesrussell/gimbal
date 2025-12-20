package scenes

import (
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	"github.com/jonesrussell/gimbal/internal/ecs/managers"
	resources "github.com/jonesrussell/gimbal/internal/ecs/managers/resource"
)

// SceneFactory is a function type that creates a scene
type SceneFactory func(
	manager *SceneManager,
	font text.Face,
	scoreManager *managers.ScoreManager,
	resourceMgr *resources.ResourceManager,
) Scene

// sceneRegistry holds the registered scene factories
var sceneRegistry = make(map[SceneType]SceneFactory)

// RegisterScene registers a scene factory for a given scene type
func RegisterScene(sceneType SceneType, factory SceneFactory) {
	sceneRegistry[sceneType] = factory
}

// GetSceneFactory returns the factory for a given scene type
func GetSceneFactory(sceneType SceneType) (SceneFactory, bool) {
	factory, exists := sceneRegistry[sceneType]
	return factory, exists
}

// CreateScene creates a scene using the registered factory
func CreateScene(
	sceneType SceneType,
	manager *SceneManager,
	font text.Face,
	scoreManager *managers.ScoreManager,
	resourceMgr *resources.ResourceManager,
) (Scene, bool) {
	factory, exists := GetSceneFactory(sceneType)
	if !exists {
		return nil, false
	}
	return factory(manager, font, scoreManager, resourceMgr), true
}
