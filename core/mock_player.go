package game

import "github.com/hajimehoshi/ebiten/v2"

type MockImage struct {
	drawn bool
}

func (mi *MockImage) DrawImage(image *ebiten.Image, options *ebiten.DrawImageOptions) {
	mi.drawn = true
}

// Implement other necessary methods...
