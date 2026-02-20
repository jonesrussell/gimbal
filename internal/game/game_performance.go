package game

import (
	"github.com/hajimehoshi/ebiten/v2"
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
		g.logger.Debug("Debug overlay toggled", "enabled", g.showDebugInfo)
	} else if !ebiten.IsKeyPressed(ebiten.KeyF3) {
		g.debugKeyPressed = false
	}
}

// endPerformanceMonitoring ends performance monitoring for the frame
func (g *ECSGame) endPerformanceMonitoring() {
	if g.perfMonitor != nil {
		g.perfMonitor.EndFrame()
	}
}
