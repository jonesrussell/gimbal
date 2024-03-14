package systems

import (
	"github.com/jonesrussell/gimbal/internal/components"
	dresolv "github.com/jonesrussell/gimbal/internal/resolv"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func UpdateObjects(ecs *ecs.ECS) {
	components.Object.Each(ecs.World, func(e *donburi.Entry) {
		obj := dresolv.GetObject(e)
		obj.Update()
	})
}
