package core

import "github.com/hajimehoshi/ebiten/v2"

const (
	screenWidth  = 800
	screenHeight = 600
)

var (
	Debug = false
)

type GimlarGame interface {
	Update() error
	Draw(screen *ebiten.Image)
}
