package debug

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/jonesrussell/gimbal/internal/common"
)

// PerformanceMonitor tracks detailed performance metrics
type PerformanceMonitor struct {
	logger common.Logger

	// Frame timing
	frameTimes     []time.Duration
	frameTimeIndex int
	maxFrameTimes  int

	// System timing
	systemTimes map[string][]time.Duration
	lastUpdate  time.Time

	// Performance thresholds
	slowFrameThreshold  time.Duration
	slowSystemThreshold time.Duration
}

// NewPerformanceMonitor creates a new performance monitor
func NewPerformanceMonitor(logger common.Logger) *PerformanceMonitor {
	return &PerformanceMonitor{
		logger:              logger,
		frameTimes:          make([]time.Duration, 60), // Track 60 frames
		maxFrameTimes:       60,
		systemTimes:         make(map[string][]time.Duration),
		slowFrameThreshold:  16 * time.Millisecond, // 60 FPS = 16.67ms per frame
		slowSystemThreshold: 5 * time.Millisecond,
	}
}

// StartFrame marks the beginning of a frame
func (pm *PerformanceMonitor) StartFrame() {
	pm.lastUpdate = time.Now()
}

// EndFrame marks the end of a frame and records timing
func (pm *PerformanceMonitor) EndFrame() {
	frameTime := time.Since(pm.lastUpdate)
	pm.frameTimes[pm.frameTimeIndex] = frameTime
	pm.frameTimeIndex = (pm.frameTimeIndex + 1) % pm.maxFrameTimes

	// Log slow frames
	if frameTime > pm.slowFrameThreshold {
		pm.logger.Warn("Slow frame detected",
			"frame_time", frameTime,
			"threshold", pm.slowFrameThreshold,
			"fps", ebiten.ActualFPS())
	}
}

// RecordSystemTime records timing for a specific system
func (pm *PerformanceMonitor) RecordSystemTime(systemName string, duration time.Duration) {
	if pm.systemTimes[systemName] == nil {
		pm.systemTimes[systemName] = make([]time.Duration, 10)
	}

	// Store last 10 measurements
	times := pm.systemTimes[systemName]
	for i := len(times) - 1; i > 0; i-- {
		times[i] = times[i-1]
	}
	times[0] = duration

	// Log slow systems
	if duration > pm.slowSystemThreshold {
		pm.logger.Warn("Slow system detected",
			"system", systemName,
			"duration", duration,
			"threshold", pm.slowSystemThreshold)
	}
}

// GetAverageFrameTime returns the average frame time over the last 60 frames
func (pm *PerformanceMonitor) GetAverageFrameTime() time.Duration {
	var total time.Duration
	count := 0
	for _, t := range pm.frameTimes {
		if t > 0 {
			total += t
			count++
		}
	}
	if count == 0 {
		return 0
	}
	return total / time.Duration(count)
}

// GetAverageSystemTime returns the average time for a specific system
func (pm *PerformanceMonitor) GetAverageSystemTime(systemName string) time.Duration {
	times := pm.systemTimes[systemName]
	if len(times) == 0 {
		return 0
	}

	var total time.Duration
	count := 0
	for _, t := range times {
		if t > 0 {
			total += t
			count++
		}
	}
	if count == 0 {
		return 0
	}
	return total / time.Duration(count)
}

// GetPerformanceReport returns a comprehensive performance report
func (pm *PerformanceMonitor) GetPerformanceReport() map[string]interface{} {
	avgFrameTime := pm.GetAverageFrameTime()

	report := map[string]interface{}{
		"current_fps":        ebiten.ActualFPS(),
		"current_tps":        ebiten.ActualTPS(),
		"avg_frame_time":     avgFrameTime,
		"frame_time_ms":      avgFrameTime.Milliseconds(),
		"slow_frames":        pm.countSlowFrames(),
		"system_performance": make(map[string]interface{}),
	}

	// Add system performance data
	for systemName := range pm.systemTimes {
		avgTime := pm.GetAverageSystemTime(systemName)
		if systemPerf, ok := report["system_performance"].(map[string]interface{}); ok {
			systemPerf[systemName] = map[string]interface{}{
				"avg_time_ms": avgTime.Milliseconds(),
				"is_slow":     avgTime > pm.slowSystemThreshold,
			}
		}
	}

	return report
}

// countSlowFrames counts how many frames were slow in the last 60 frames
func (pm *PerformanceMonitor) countSlowFrames() int {
	count := 0
	for _, t := range pm.frameTimes {
		if t > pm.slowFrameThreshold {
			count++
		}
	}
	return count
}
