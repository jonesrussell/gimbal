package scenes

import (
	"image/color"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jonesrussell/gimbal/internal/factory"
	"github.com/jonesrussell/gimbal/internal/layers"
	dresolv "github.com/jonesrussell/gimbal/internal/resolv"
	"github.com/jonesrussell/gimbal/internal/systems"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

type GalaxyScene struct {
	ecs  *ecs.ECS
	once sync.Once
}

func (ps *GalaxyScene) Update() {
	ps.once.Do(ps.configure)
	ps.ecs.Update()
}

func (ps *GalaxyScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{20, 20, 40, 255})
	ps.ecs.Draw(screen)
}

func (ps *GalaxyScene) configure() {
	ecs := ecs.NewECS(donburi.NewWorld())

	ecs.AddSystem(systems.UpdatePlayer)
	ecs.AddSystem(systems.UpdateObjects)

	ecs.AddRenderer(layers.Default, systems.DrawPlayer)

	ps.ecs = ecs

	// Define the world's Space. Here, a Space is essentially a grid (the game's width and height, or 640x360), made up of 16x16 cells. Each cell can have 0 or more Objects within it,
	// and collisions can be found by checking the Space to see if the Cells at specific positions contain (or would contain) Objects. This is a broad, simplified approach to collision
	// detection.
	space := factory.CreateSpace(ps.ecs)

	dresolv.Add(space,
		// Create the Player. NewPlayer adds it to the world's Space.
		factory.CreatePlayer(ps.ecs),
	)

}
