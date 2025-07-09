package ecs

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"

	"github.com/jonesrussell/gimbal/internal/common"
)

// MovementSystemWrapper wraps the MovementSystem to implement System interface
type MovementSystemWrapper struct{}

func (msw *MovementSystemWrapper) Update(world donburi.World, args ...interface{}) error {
	MovementSystem(world)
	return nil
}

func (msw *MovementSystemWrapper) Name() string {
	return "MovementSystem"
}

// OrbitalMovementSystemWrapper wraps the OrbitalMovementSystem to implement System interface
type OrbitalMovementSystemWrapper struct{}

func (omsw *OrbitalMovementSystemWrapper) Update(world donburi.World, args ...interface{}) error {
	OrbitalMovementSystem(world)
	return nil
}

func (omsw *OrbitalMovementSystemWrapper) Name() string {
	return "OrbitalMovementSystem"
}

// RenderSystemWrapper wraps the RenderSystem to implement System interface
type RenderSystemWrapper struct {
	screen *ebiten.Image
}

func NewRenderSystemWrapper(screen *ebiten.Image) *RenderSystemWrapper {
	return &RenderSystemWrapper{
		screen: screen,
	}
}

func (rsw *RenderSystemWrapper) Update(world donburi.World, args ...interface{}) error {
	if rsw.screen == nil {
		return common.NewGameError(common.ErrorCodeRenderingFailed, "screen is nil")
	}
	RenderSystem(world, rsw.screen)
	return nil
}

func (rsw *RenderSystemWrapper) Name() string {
	return "RenderSystem"
}

// StarMovementSystemWrapper wraps the StarMovementSystem to implement System interface
type StarMovementSystemWrapper struct {
	ecsInstance *ecs.ECS
	config      *common.GameConfig
}

func NewStarMovementSystemWrapper(ecsInstance *ecs.ECS, config *common.GameConfig) *StarMovementSystemWrapper {
	return &StarMovementSystemWrapper{
		ecsInstance: ecsInstance,
		config:      config,
	}
}

func (smsw *StarMovementSystemWrapper) Update(world donburi.World, args ...interface{}) error {
	if smsw.ecsInstance == nil {
		return common.NewGameError(common.ErrorCodeSystemFailed, "ecs instance is nil")
	}
	if smsw.config == nil {
		return common.NewGameError(common.ErrorCodeConfigMissing, "config is nil")
	}
	StarMovementSystem(smsw.ecsInstance, smsw.config)
	return nil
}

func (smsw *StarMovementSystemWrapper) Name() string {
	return "StarMovementSystem"
}

// PlayerInputSystemWrapper wraps the PlayerInputSystem to implement System interface
type PlayerInputSystemWrapper struct {
	inputAngle common.Angle
}

func NewPlayerInputSystemWrapper(inputAngle common.Angle) *PlayerInputSystemWrapper {
	return &PlayerInputSystemWrapper{
		inputAngle: inputAngle,
	}
}

func (pisw *PlayerInputSystemWrapper) Update(world donburi.World, args ...interface{}) error {
	PlayerInputSystem(world, pisw.inputAngle)
	return nil
}

func (pisw *PlayerInputSystemWrapper) Name() string {
	return "PlayerInputSystem"
}
