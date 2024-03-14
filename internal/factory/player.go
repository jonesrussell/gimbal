package factory

import (
	"github.com/jonesrussell/gimbal/internal/archetypes"
	"github.com/jonesrussell/gimbal/internal/assets"
	"github.com/jonesrussell/gimbal/internal/components"
	dresolv "github.com/jonesrussell/gimbal/internal/resolv"
	"github.com/solarlune/resolv"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func CreatePlayer(ecs *ecs.ECS) *donburi.Entry {
	player := archetypes.Player.Spawn(ecs)

	obj := resolv.NewObject(640/2, 480/2, 16, 16)
	dresolv.SetObject(player, obj)
	components.Player.SetValue(player, components.PlayerData{
		Sprite: assets.LoadPlayerSprite(),
	})

	obj.SetShape(resolv.NewRectangle(0, 0, 16, 16))

	return player
}
