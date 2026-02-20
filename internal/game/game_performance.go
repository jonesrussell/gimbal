package game

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/jonesrussell/gimbal/internal/dbg"
)

// updatePerformanceMonitoring handles performance monitoring for the frame
func (g *ECSGame) updatePerformanceMonitoring() {
	if g.perfMonitor != nil {
		g.perfMonitor.StartFrame()
	}
}

// updateDebugLogging handles periodic debug logging (no per-tick log to avoid spam)
func (g *ECSGame) updateDebugLogging() {
	g.frameCount++
}

// updateDebugInput handles debug key input
func (g *ECSGame) updateDebugInput() {
	if ebiten.IsKeyPressed(ebiten.KeyF3) && !g.debugKeyPressed {
		g.showDebugInfo = !g.showDebugInfo
		if g.renderDebugger != nil {
			g.renderDebugger.Toggle()
		}
		g.debugKeyPressed = true
		dbg.Log(dbg.System, "Debug overlay toggled (enabled=%v)", g.showDebugInfo)
	} else if !ebiten.IsKeyPressed(ebiten.KeyF3) {
		g.debugKeyPressed = false
	}
	if ebiten.IsKeyPressed(ebiten.KeyF4) && !g.traceKeyPressed {
		dbg.Trace()
		g.traceKeyPressed = true
	} else if !ebiten.IsKeyPressed(ebiten.KeyF4) {
		g.traceKeyPressed = false
	}
}

// endPerformanceMonitoring ends performance monitoring for the frame
func (g *ECSGame) endPerformanceMonitoring() {
	if g.perfMonitor != nil {
		g.perfMonitor.EndFrame()
	}
}
