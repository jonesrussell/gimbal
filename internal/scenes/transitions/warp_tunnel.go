package transitions

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// WarpTunnelTransition implements a warp tunnel transition using animation frames
type WarpTunnelTransition struct {
	duration     float64
	elapsed      float64
	frames       []*ebiten.Image
	frameCount   int
	currentFrame int
	complete     bool
	screenWidth  int
	screenHeight int
}

// NewWarpTunnelTransition creates a new warp tunnel transition
func NewWarpTunnelTransition(duration float64, frames []*ebiten.Image, screenWidth, screenHeight int) *WarpTunnelTransition {
	return &WarpTunnelTransition{
		duration:     duration,
		frames:       frames,
		frameCount:   len(frames),
		screenWidth:  screenWidth,
		screenHeight: screenHeight,
	}
}

// Update updates the warp tunnel transition
func (w *WarpTunnelTransition) Update(deltaTime float64) bool {
	if w.complete {
		return true
	}

	w.elapsed += deltaTime
	progress := w.GetProgress()

	// Calculate current frame based on progress
	if w.frameCount > 0 {
		frameIndex := int(progress * float64(w.frameCount))
		if frameIndex >= w.frameCount {
			frameIndex = w.frameCount - 1
		}
		w.currentFrame = frameIndex
	}

	if w.elapsed >= w.duration {
		w.elapsed = w.duration
		w.complete = true
	}

	return w.complete
}

// Draw draws the warp tunnel transition
func (w *WarpTunnelTransition) Draw(screen *ebiten.Image, from, to *ebiten.Image) {
	// Draw the "to" scene in the background (scaled down to simulate tunnel effect)
	if to != nil {
		op := &ebiten.DrawImageOptions{}
		scale := 0.3 + (w.GetProgress() * 0.7) // Scale from 0.3 to 1.0
		op.GeoM.Scale(scale, scale)
		centerX := float64(w.screenWidth) / 2
		centerY := float64(w.screenHeight) / 2
		op.GeoM.Translate(centerX-float64(to.Bounds().Dx())*scale/2, centerY-float64(to.Bounds().Dy())*scale/2)
		screen.DrawImage(to, op)
	}

	// Draw warp tunnel frame overlay
	if w.frameCount > 0 && w.currentFrame < len(w.frames) && w.frames[w.currentFrame] != nil {
		op := &ebiten.DrawImageOptions{}
		// Fade out the tunnel effect as we progress
		alpha := 1.0 - w.GetProgress()
		op.ColorScale.SetA(float32(alpha))
		screen.DrawImage(w.frames[w.currentFrame], op)
	}

	// Draw radial tunnel effect (fallback if no frames)
	if w.frameCount == 0 {
		w.drawRadialTunnel(screen, to)
	}
}

// drawRadialTunnel draws a procedural radial tunnel effect
func (w *WarpTunnelTransition) drawRadialTunnel(screen *ebiten.Image, to *ebiten.Image) {
	progress := w.GetProgress()
	centerX := float64(w.screenWidth) / 2
	centerY := float64(w.screenHeight) / 2

	// Create tunnel effect using concentric circles
	tunnelImg := ebiten.NewImage(w.screenWidth, w.screenHeight)

	// Draw radial lines
	for i := 0; i < 16; i++ {
		angle := float64(i) * (2 * math.Pi / 16)
		radius := float64(w.screenWidth) * (1.0 - progress)

		// Draw line from center to edge
		// This is simplified - in practice, you'd use a more sophisticated drawing method
		// For now, we'll use the frame-based approach when frames are available
	}

	op := &ebiten.DrawImageOptions{}
	op.ColorScale.SetA(float32(1.0 - progress))
	screen.DrawImage(tunnelImg, op)
}

// Reset resets the warp tunnel transition
func (w *WarpTunnelTransition) Reset() {
	w.elapsed = 0
	w.complete = false
	w.currentFrame = 0
}

// GetProgress returns the transition progress (0.0 to 1.0)
func (w *WarpTunnelTransition) GetProgress() float64 {
	if w.duration == 0 {
		return 1.0
	}
	progress := w.elapsed / w.duration
	if progress > 1.0 {
		return 1.0
	}
	return progress
}
