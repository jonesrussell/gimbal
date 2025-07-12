package viewport

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

// CachedLayout represents a cached layout calculation
type CachedLayout struct {
	ViewportHash string
	Elements     []RenderCommand
	Valid        bool
}

// RenderCommand represents a rendering operation
type RenderCommand struct {
	ElementType string
	Position    Position2D
	Size        Position2D
	Layer       int
	Content     interface{}
	Priority    int // Higher priority renders first
}

// ResponsiveRenderer implements 2025 efficient viewport-based rendering
type ResponsiveRenderer struct {
	lastViewport  *AdvancedViewportManager
	layoutCache   map[string]CachedLayout
	renderQueue   []RenderCommand
	batchedRender bool

	// Performance metrics
	frameCount    int
	lastFrameTime float64
	cacheHits     int
	cacheMisses   int
}

// NewResponsiveRenderer creates a new performance-optimized renderer
func NewResponsiveRenderer() *ResponsiveRenderer {
	return &ResponsiveRenderer{
		layoutCache:   make(map[string]CachedLayout),
		renderQueue:   make([]RenderCommand, 0),
		batchedRender: true,
	}
}

// ShouldRelayout implements 2025 performance optimization
func (rr *ResponsiveRenderer) ShouldRelayout(current *AdvancedViewportManager) bool {
	if rr.lastViewport == nil {
		return true
	}

	// Use the viewport's built-in relayout detection
	return current.NeedsRelayout()
}

// AddRenderCommand adds a command to the render queue
func (rr *ResponsiveRenderer) AddRenderCommand(cmd RenderCommand) {
	rr.renderQueue = append(rr.renderQueue, cmd)
}

// ClearRenderQueue clears the render queue
func (rr *ResponsiveRenderer) ClearRenderQueue() {
	rr.renderQueue = rr.renderQueue[:0]
}

// SortRenderQueue sorts the render queue by priority and layer
func (rr *ResponsiveRenderer) SortRenderQueue() {
	// Sort by priority (highest first), then by layer (lowest first)
	for i := 0; i < len(rr.renderQueue)-1; i++ {
		for j := i + 1; j < len(rr.renderQueue); j++ {
			if rr.renderQueue[i].Priority < rr.renderQueue[j].Priority ||
				(rr.renderQueue[i].Priority == rr.renderQueue[j].Priority &&
					rr.renderQueue[i].Layer > rr.renderQueue[j].Layer) {
				rr.renderQueue[i], rr.renderQueue[j] = rr.renderQueue[j], rr.renderQueue[i]
			}
		}
	}
}

// RenderFrame renders a complete frame with 2025 optimizations
func (rr *ResponsiveRenderer) RenderFrame(screen *ebiten.Image, viewport *AdvancedViewportManager, hud *GameHUD) {
	rr.frameCount++

	// Check if we need to relayout
	if rr.ShouldRelayout(viewport) {
		rr.performRelayout(viewport)
		rr.lastViewport = viewport
	}

	// Sort render queue for proper layering
	rr.SortRenderQueue()

	// Render all commands
	for _, cmd := range rr.renderQueue {
		rr.renderCommand(screen, cmd, viewport, hud)
	}

	// Clear queue for next frame
	rr.ClearRenderQueue()
}

// performRelayout performs a complete layout recalculation
func (rr *ResponsiveRenderer) performRelayout(viewport *AdvancedViewportManager) {
	// Generate viewport hash for caching
	viewportHash := rr.generateViewportHash(viewport)

	// Check cache first
	if cached, exists := rr.layoutCache[viewportHash]; exists && cached.Valid {
		rr.cacheHits++
		// Use cached layout
		rr.renderQueue = append(rr.renderQueue, cached.Elements...)
		return
	}

	rr.cacheMisses++

	// Perform new layout calculation
	rr.calculateLayout(viewport)

	// Cache the result
	rr.layoutCache[viewportHash] = CachedLayout{
		ViewportHash: viewportHash,
		Elements:     make([]RenderCommand, len(rr.renderQueue)),
		Valid:        true,
	}
	copy(rr.layoutCache[viewportHash].Elements, rr.renderQueue)

	// Limit cache size to prevent memory leaks
	rr.cleanupCache()
}

// generateViewportHash generates a hash for viewport state
func (rr *ResponsiveRenderer) generateViewportHash(viewport *AdvancedViewportManager) string {
	width, height := viewport.GetCurrentDimensions()
	deviceClass := viewport.GetDeviceClass()
	orientation := viewport.GetOrientation()

	// Simple hash based on key viewport properties
	return fmt.Sprintf("%dx%d-%s-%s", width, height, deviceClass, orientation)
}

// calculateLayout calculates the layout for the current viewport
func (rr *ResponsiveRenderer) calculateLayout(viewport *AdvancedViewportManager) {
	// Clear existing queue
	rr.ClearRenderQueue()

	// Add UI elements based on device class
	switch viewport.GetDeviceClass() {
	case string(DeviceClassMobile):
		rr.calculateMobileLayout(viewport)
	case string(DeviceClassTablet):
		rr.calculateTabletLayout(viewport)
	case string(DeviceClassUltrawide):
		rr.calculateUltrawideLayout(viewport)
	default:
		rr.calculateDesktopLayout(viewport)
	}
}

// calculateMobileLayout calculates layout for mobile devices
func (rr *ResponsiveRenderer) calculateMobileLayout(viewport *AdvancedViewportManager) {
	width, height := viewport.GetCurrentDimensions()

	// Health bar at top
	rr.AddRenderCommand(RenderCommand{
		ElementType: "health_bar",
		Position:    Position2D{X: 20, Y: 20},
		Size:        Position2D{X: 200, Y: 20},
		Layer:       1,
		Priority:    10,
		Content:     0.8, // Example health value
	})

	// Score display at top right
	rr.AddRenderCommand(RenderCommand{
		ElementType: "score_display",
		Position:    Position2D{X: float64(width) - 220, Y: 20},
		Size:        Position2D{X: 200, Y: 30},
		Layer:       1,
		Priority:    10,
		Content:     "Score: 1250",
	})

	// Ammo counter at bottom
	rr.AddRenderCommand(RenderCommand{
		ElementType: "ammo_counter",
		Position:    Position2D{X: 20, Y: float64(height) - 40},
		Size:        Position2D{X: 150, Y: 20},
		Layer:       1,
		Priority:    10,
		Content:     5, // Example ammo count
	})
}

// calculateTabletLayout calculates layout for tablet devices
func (rr *ResponsiveRenderer) calculateTabletLayout(viewport *AdvancedViewportManager) {
	width, height := viewport.GetCurrentDimensions()

	// Larger UI elements for tablet
	rr.AddRenderCommand(RenderCommand{
		ElementType: "health_bar",
		Position:    Position2D{X: 30, Y: 30},
		Size:        Position2D{X: 300, Y: 25},
		Layer:       1,
		Priority:    10,
		Content:     0.8,
	})

	rr.AddRenderCommand(RenderCommand{
		ElementType: "score_display",
		Position:    Position2D{X: float64(width) - 330, Y: 30},
		Size:        Position2D{X: 300, Y: 40},
		Layer:       1,
		Priority:    10,
		Content:     "Score: 1250",
	})

	rr.AddRenderCommand(RenderCommand{
		ElementType: "ammo_counter",
		Position:    Position2D{X: 30, Y: float64(height) - 50},
		Size:        Position2D{X: 200, Y: 25},
		Layer:       1,
		Priority:    10,
		Content:     5,
	})
}

// calculateDesktopLayout calculates layout for desktop devices
func (rr *ResponsiveRenderer) calculateDesktopLayout(viewport *AdvancedViewportManager) {
	width, height := viewport.GetCurrentDimensions()

	// Standard desktop layout
	rr.AddRenderCommand(RenderCommand{
		ElementType: "health_bar",
		Position:    Position2D{X: 40, Y: 40},
		Size:        Position2D{X: 400, Y: 30},
		Layer:       1,
		Priority:    10,
		Content:     0.8,
	})

	rr.AddRenderCommand(RenderCommand{
		ElementType: "score_display",
		Position:    Position2D{X: float64(width) - 440, Y: 40},
		Size:        Position2D{X: 400, Y: 50},
		Layer:       1,
		Priority:    10,
		Content:     "Score: 1250",
	})

	rr.AddRenderCommand(RenderCommand{
		ElementType: "ammo_counter",
		Position:    Position2D{X: 40, Y: float64(height) - 60},
		Size:        Position2D{X: 250, Y: 30},
		Layer:       1,
		Priority:    10,
		Content:     5,
	})
}

// calculateUltrawideLayout calculates layout for ultrawide displays
func (rr *ResponsiveRenderer) calculateUltrawideLayout(viewport *AdvancedViewportManager) {
	width, height := viewport.GetCurrentDimensions()

	// Ultrawide: spread UI elements across the wider screen
	rr.AddRenderCommand(RenderCommand{
		ElementType: "health_bar",
		Position:    Position2D{X: 60, Y: 50},
		Size:        Position2D{X: 500, Y: 35},
		Layer:       1,
		Priority:    10,
		Content:     0.8,
	})

	rr.AddRenderCommand(RenderCommand{
		ElementType: "score_display",
		Position:    Position2D{X: float64(width) - 560, Y: 50},
		Size:        Position2D{X: 500, Y: 60},
		Layer:       1,
		Priority:    10,
		Content:     "Score: 1250",
	})

	rr.AddRenderCommand(RenderCommand{
		ElementType: "ammo_counter",
		Position:    Position2D{X: 60, Y: float64(height) - 70},
		Size:        Position2D{X: 300, Y: 35},
		Layer:       1,
		Priority:    10,
		Content:     5,
	})

	// Additional status indicators for ultrawide
	rr.AddRenderCommand(RenderCommand{
		ElementType: "status_indicator",
		Position:    Position2D{X: float64(width) - 100, Y: 50},
		Size:        Position2D{X: 30, Y: 30},
		Layer:       1,
		Priority:    10,
		Content:     "ready",
	})
}

// renderCommand renders a single command
func (rr *ResponsiveRenderer) renderCommand(screen *ebiten.Image, cmd RenderCommand, viewport *AdvancedViewportManager, hud *GameHUD) {
	// Create UI element for HUD rendering
	element := UIElement{
		Position:    cmd.Position,
		Size:        cmd.Size,
		Layer:       cmd.Layer,
		ElementType: cmd.ElementType,
		Content:     cmd.Content,
	}

	// Render with HUD effects
	hud.RenderWithDepth(screen, element, viewport)
}

// cleanupCache removes old cache entries to prevent memory leaks
func (rr *ResponsiveRenderer) cleanupCache() {
	const maxCacheSize = 50

	if len(rr.layoutCache) > maxCacheSize {
		// Simple cleanup: remove oldest entries
		// In a production system, you might want LRU or more sophisticated caching
		for key := range rr.layoutCache {
			delete(rr.layoutCache, key)
			if len(rr.layoutCache) <= maxCacheSize {
				break
			}
		}
	}
}

// GetPerformanceMetrics returns performance metrics
func (rr *ResponsiveRenderer) GetPerformanceMetrics() map[string]interface{} {
	return map[string]interface{}{
		"frame_count":  rr.frameCount,
		"cache_hits":   rr.cacheHits,
		"cache_misses": rr.cacheMisses,
		"cache_ratio":  float64(rr.cacheHits) / float64(rr.cacheHits+rr.cacheMisses),
		"cache_size":   len(rr.layoutCache),
	}
}

// SetBatchedRender enables or disables batched rendering
func (rr *ResponsiveRenderer) SetBatchedRender(enabled bool) {
	rr.batchedRender = enabled
}
