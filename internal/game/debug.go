package game

import (
	"image"
	"log/slog"
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
		slog.Info(format, v...)
	}
}

func (d *Debugger) DebugPrint() {
	d.logIfDebugEnabled("message", "Debug mode is enabled.")
}

func (d *Debugger) DebugPrintOrientation(viewAngle float64) {
	d.logIfDebugEnabled("Player", "viewAngle", viewAngle)
}

func (d *Debugger) DebugPrintDirection(direction float64) {
	d.logIfDebugEnabled("Player", "direction", direction)
}

func (d *Debugger) DebugPrintAngle(angle float64) {
	d.logIfDebugEnabled("Player", "angle", angle)
}

func (d *Debugger) DebugPrintPosition(position image.Point) {
	d.logIfDebugEnabled("Player", "X", float64(position.X), "Y", float64(position.Y))
}

func (d *Debugger) DebugPlayer(player *Player) {
	d.DebugPrintOrientation(player.viewAngle)
	d.DebugPrintDirection(player.direction)
	d.DebugPrintAngle(player.angle)
	pos := image.Point{X: int(player.Object.Position.X), Y: int(player.Object.Position.Y)}
	d.DebugPrintPosition(pos)
}
