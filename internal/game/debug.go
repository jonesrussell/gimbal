package game

import (
	"image"
	"log"
)

type Debugger struct {
	debug bool
}

func NewDebugger() *Debugger {
	return &Debugger{
		debug: false,
	}
}

func (d *Debugger) IsDebugMode() bool {
	return d.debug
}

func (d *Debugger) DebugMode(debug bool) {
	d.debug = debug
}

func (d *Debugger) logIfDebugEnabled(format string, v ...interface{}) {
	if d.debug {
		log.Printf(format, v...)
	}
}

func (d *Debugger) DebugPrint() {
	d.logIfDebugEnabled("Debug mode is enabled.")
}

func (d *Debugger) DebugPrintOrientation(viewAngle float64) {
	d.logIfDebugEnabled("Player viewAngle: %f", viewAngle)
}

func (d *Debugger) DebugPrintDirection(direction float64) {
	d.logIfDebugEnabled("Player direction: %f", direction)
}

func (d *Debugger) DebugPrintAngle(angle float64) {
	d.logIfDebugEnabled("Player angle: %f", angle)
}

func (d *Debugger) DebugPrintPosition(position image.Point) {
	d.logIfDebugEnabled("Player position: (%.2f, %.2f)", float64(position.X), float64(position.Y))
}

func (d *Debugger) DebugPlayer(player *Player) {
	d.DebugPrintOrientation(player.viewAngle)
	d.DebugPrintDirection(player.direction)
	d.DebugPrintAngle(player.angle)
	pos := image.Point{X: int(player.Object.Position.X), Y: int(player.Object.Position.Y)}
	d.DebugPrintPosition(pos)
}
