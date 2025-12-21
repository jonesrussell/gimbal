package debug

import (
	"fmt"
	"runtime"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

// drawPerformanceMetrics draws condensed performance information
func (dr *DebugRenderer) drawPerformanceMetrics(screen *ebiten.Image, world donburi.World) {
	if dr.font == nil {
		return
	}

	// Get memory stats
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// Get entity count
	entityCount := world.Len()

	// Calculate frame time
	fps := ebiten.ActualFPS()
	frameTime := 0.0
	if fps > 0 {
		frameTime = 1000.0 / fps // Convert to milliseconds
	}

	// Format enhanced debug text
	debugText := fmt.Sprintf("FPS:%.0f(%.1fms) TPS:%.0f E:%d M:%dK [F3]",
		fps,
		frameTime,
		ebiten.ActualTPS(),
		entityCount,
		m.Alloc/1024)

	// Draw text with background for better readability
	dr.drawTextWithBackground(screen, debugText, 10, 30)

	// Draw additional performance warnings
	if fps < 50 {
		warningText := "PERFORMANCE WARNING: Low FPS detected"
		dr.drawTextWithBackground(screen, warningText, 10, 50)
	}

	if m.Alloc > 100*1024*1024 { // 100MB
		warningText := "MEMORY WARNING: High memory usage"
		dr.drawTextWithBackground(screen, warningText, 10, 70)
	}
}
