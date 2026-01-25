// Package presenter provides UI presentation logic with event-driven updates.
package presenter

import (
	"sync"

	"github.com/yohamta/donburi"

	"github.com/jonesrussell/gimbal/internal/ecs/events"
	"github.com/jonesrussell/gimbal/internal/ui/state"
)

// HUDPresenter manages HUD state with event-driven updates.
// It subscribes to game events and only updates the UI when state changes.
type HUDPresenter struct {
	mu    sync.RWMutex
	data  state.HUDData
	dirty bool
}

// NewHUDPresenter creates a new HUD presenter with initial values.
func NewHUDPresenter(initialScore, initialLives, initialLevel int, initialHealth float64) *HUDPresenter {
	return &HUDPresenter{
		data: state.HUDData{
			Score:  initialScore,
			Lives:  initialLives,
			Level:  initialLevel,
			Health: initialHealth,
		},
		dirty: true, // Start dirty to ensure initial render
	}
}

// Subscribe registers event handlers with the event system.
func (p *HUDPresenter) Subscribe(eventSystem *events.EventSystem) {
	eventSystem.SubscribeToScoreChanged(p.onScoreChanged)
	eventSystem.SubscribeToPlayerDamaged(p.onPlayerDamaged)
	eventSystem.SubscribeToLifeAdded(p.onLifeAdded)
	eventSystem.SubscribeToLevelChanged(p.onLevelChanged)
}

// IsDirty returns true if HUD state has changed since last flush.
func (p *HUDPresenter) IsDirty() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.dirty
}

// GetData returns the current HUD data and clears the dirty flag.
func (p *HUDPresenter) GetData() state.HUDData {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.dirty = false
	return p.data
}

// PeekData returns the current HUD data without clearing the dirty flag.
func (p *HUDPresenter) PeekData() state.HUDData {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.data
}

// SetHealth updates the health value directly (for continuous updates like health regen).
func (p *HUDPresenter) SetHealth(current, maximum int) {
	p.mu.Lock()
	defer p.mu.Unlock()

	newHealth := 1.0
	if maximum > 0 {
		newHealth = float64(current) / float64(maximum)
	}

	if p.data.Health != newHealth || p.data.Lives != current {
		p.data.Health = newHealth
		p.data.Lives = current
		p.dirty = true
	}
}

// Event handlers

func (p *HUDPresenter) onScoreChanged(_ donburi.World, e events.ScoreChangedEvent) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.data.Score != e.NewScore {
		p.data.Score = e.NewScore
		p.dirty = true
	}
}

func (p *HUDPresenter) onPlayerDamaged(_ donburi.World, e events.PlayerDamagedEvent) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.data.Lives != e.RemainingLives {
		p.data.Lives = e.RemainingLives
		p.dirty = true
	}
}

func (p *HUDPresenter) onLifeAdded(_ donburi.World, e events.LifeAddedEvent) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.data.Lives != e.NewLives {
		p.data.Lives = e.NewLives
		p.dirty = true
	}
}

func (p *HUDPresenter) onLevelChanged(_ donburi.World, e events.LevelChangedEvent) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.data.Level != e.NewLevel {
		p.data.Level = e.NewLevel
		p.dirty = true
	}
}

// MarkDirty forces a HUD update on the next frame.
func (p *HUDPresenter) MarkDirty() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.dirty = true
}
