package behavior

import (
	"math"
	"time"

	"github.com/yohamta/donburi"

	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/config"
	"github.com/jonesrussell/gimbal/internal/dbg"
	"github.com/jonesrussell/gimbal/internal/ecs/core"
)

// HoveringState handles hovering near the center before orbiting
type HoveringState struct {
	config      *config.GameConfig
	logger      common.Logger
	hoverRadius float64       // Radius to hover at
	hoverTime   time.Duration // How long to hover before moving to orbit
}

// NewHoveringState creates a new hovering state handler
func NewHoveringState(cfg *config.GameConfig, logger common.Logger) *HoveringState {
	return &HoveringState{
		config:      cfg,
		logger:      logger,
		hoverRadius: 80.0,            // Hover 80 pixels from center
		hoverTime:   2 * time.Second, // Hover for 2 seconds
	}
}

// StateType returns the state type
func (hs *HoveringState) StateType() core.BehaviorStateType {
	return core.StateHovering
}

// Enter is called when transitioning into this state
func (hs *HoveringState) Enter(entry *donburi.Entry, data *core.BehaviorStateData) {
	dbg.Log(dbg.State, "Entering hover state")
}

// Update is called every frame while in this state
func (hs *HoveringState) Update(entry *donburi.Entry, data *core.BehaviorStateData, deltaTime float64) {
	// Gentle circular motion while hovering
	hs.updateHoverPosition(entry, data, deltaTime)
}

// updateHoverPosition creates a gentle hovering motion near center
func (hs *HoveringState) updateHoverPosition(entry *donburi.Entry, data *core.BehaviorStateData, deltaTime float64) {
	if !entry.HasComponent(core.Position) {
		return
	}

	pos := core.Position.Get(entry)
	centerX := float64(hs.config.ScreenSize.Width) / 2
	centerY := float64(hs.config.ScreenSize.Height) / 2

	// Calculate hover angle (slowly rotating)
	hoverSpeed := 30.0 // degrees per second
	hoverAngle := float64(data.StateTime.Seconds()) * hoverSpeed * math.Pi / 180

	// Calculate target hover position
	targetX := centerX + hs.hoverRadius*math.Cos(hoverAngle)
	targetY := centerY + hs.hoverRadius*math.Sin(hoverAngle)

	// Smoothly move toward hover position
	moveSpeed := 100.0 // pixels per second
	dx := targetX - pos.X
	dy := targetY - pos.Y
	distance := math.Sqrt(dx*dx + dy*dy)

	if distance > 1 {
		dx /= distance
		dy /= distance
		moveDistance := math.Min(moveSpeed*deltaTime, distance)
		pos.X += dx * moveDistance
		pos.Y += dy * moveDistance
	}
}

// Exit is called when transitioning out of this state
func (hs *HoveringState) Exit(entry *donburi.Entry, data *core.BehaviorStateData) {
	dbg.Log(dbg.State, "Exiting hover state")
}

// NextState determines the next state
func (hs *HoveringState) NextState(entry *donburi.Entry, data *core.BehaviorStateData) core.BehaviorStateType {
	// Transition to orbiting after hover time
	if data.StateTime >= hs.hoverTime {
		return core.StateOrbiting
	}

	return core.StateHovering
}
