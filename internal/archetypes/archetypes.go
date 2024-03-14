package archetypes

import (
	"github.com/jonesrussell/gimbal/internal/components"
	"github.com/jonesrussell/gimbal/internal/layers"
	"github.com/jonesrussell/gimbal/internal/tags"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

var (
	Player = newArchetype(
		tags.Player,
		components.Player,
		components.Object,
	)
	Space = newArchetype(
		components.Space,
	)
)

type archetype struct {
	components []donburi.IComponentType
}

func newArchetype(cs ...donburi.IComponentType) *archetype {
	return &archetype{
		components: cs,
	}
}

func (a *archetype) Spawn(ecs *ecs.ECS, cs ...donburi.IComponentType) *donburi.Entry {
	e := ecs.World.Entry(ecs.Create(
		layers.Default,
		append(a.components, cs...)...,
	))
	return e
}
