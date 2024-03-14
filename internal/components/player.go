package components

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
	"github.com/yohamta/donburi"
)

type PlayerData struct {
	Sprite    *ebiten.Image
	SubSprite *ebiten.Image
	Path      []resolv.Vector
	ViewAngle float64
	Direction float64
	Angle     float64
	AngleStep float64
}

var Player = donburi.NewComponentType[PlayerData]()
