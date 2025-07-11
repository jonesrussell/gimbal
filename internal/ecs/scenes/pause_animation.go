package ecs

import (
	"math"
	"time"
)

// AnimationState tracks current animation values
type AnimationState struct {
	Time              float64
	FadeIn            float64
	SelectionChanged  bool
	LastSelectionTime time.Time
}

// AnimationCalculator computes animation values
type AnimationCalculator struct {
	config AnimationConfig
	state  *AnimationState
}

// NewAnimationCalculator creates a new animation calculator
func NewAnimationCalculator(config AnimationConfig, state *AnimationState) *AnimationCalculator {
	return &AnimationCalculator{
		config: config,
		state:  state,
	}
}

// UpdateAnimationState updates the animation state with delta time
func (ac *AnimationCalculator) UpdateAnimationState(dt float64) {
	ac.state.Time += dt
	ac.state.FadeIn = math.Min(1.0, ac.state.FadeIn+dt/ac.config.FadeInDuration)

	if ac.state.SelectionChanged && time.Since(ac.state.LastSelectionTime).Seconds() > 0.1 {
		ac.state.SelectionChanged = false
	}
}

// ResetAnimationState resets all animation values
func (ac *AnimationCalculator) ResetAnimationState() {
	ac.state.Time = 0
	ac.state.FadeIn = 0
	ac.state.SelectionChanged = false
}

// MarkSelectionChanged marks that selection has changed
func (ac *AnimationCalculator) MarkSelectionChanged() {
	ac.state.SelectionChanged = true
	ac.state.LastSelectionTime = time.Now()
}

// AnimationCalculator methods (SRP - only calculates animations)
func (ac *AnimationCalculator) GetFadeInAlpha() float64 {
	return ac.state.FadeIn
}

func (ac *AnimationCalculator) GetPulseValue(speed float64) float64 {
	return 0.8 + 0.2*math.Sin(ac.state.Time*speed)
}

func (ac *AnimationCalculator) GetChevronPosition(baseX float64) float64 {
	return baseX + ac.config.ChevronAmplitude*math.Sin(ac.state.Time*ac.config.ChevronSpeed)
}

func (ac *AnimationCalculator) GetScaleValue() float64 {
	return 1.0 + ac.config.ScaleAmplitude*math.Sin(ac.state.Time*ac.config.PulseSpeed)
}

func (ac *AnimationCalculator) GetTitlePulse() float64 {
	return 0.9 + 0.1*math.Sin(ac.state.Time*ac.config.TitlePulseSpeed)
}

func (ac *AnimationCalculator) GetHintPulse() float64 {
	return 0.8 + 0.2*math.Sin(ac.state.Time*ac.config.HintPulseSpeed)
}
