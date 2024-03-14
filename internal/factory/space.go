package factory

import (
	"github.com/jonesrussell/gimbal/internal/archetypes"
	"github.com/jonesrussell/gimbal/internal/components"
	"github.com/jonesrussell/gimbal/internal/config"
	"github.com/solarlune/resolv"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func CreateSpace(ecs *ecs.ECS) *donburi.Entry {
	space := archetypes.Space.Spawn(ecs)

	cfg := config.C
	spaceData := resolv.NewSpace(cfg.Width, cfg.Height, 16, 16)
	components.Space.Set(space, spaceData)

	return space
}
