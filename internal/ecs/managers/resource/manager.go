package resources

import (
	"context"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	"github.com/jonesrussell/gimbal/internal/common"
)

// ResourceType represents different types of resources
type ResourceType int

const (
	ResourceSprite ResourceType = iota
	ResourceSound
	ResourceFont
	ResourceData
)

// Resource represents a loaded game resource
type Resource struct {
	Type ResourceType
	Name string
	Data interface{}
}

// ResourceManager manages all game resources
type ResourceManager struct {
	resources   map[string]*Resource
	mutex       sync.RWMutex
	logger      common.Logger
	defaultFont text.Face

	scaledCache map[string]*ebiten.Image // Cache for scaled sprites
	audioPlayer *AudioPlayer             // Audio player for background music
}

// NewResourceManager creates a new resource management system with the provided logger
func NewResourceManager(ctx context.Context, logger common.Logger) *ResourceManager {
	rm := &ResourceManager{
		resources:   make(map[string]*Resource),
		logger:      logger,
		scaledCache: make(map[string]*ebiten.Image),
	}

	// Initialize audio player (44100 Hz sample rate)
	// Audio is optional - if initialization fails (e.g., no audio device in container),
	// the game will continue without audio
	audioPlayer, audioErr := NewAudioPlayer(44100, logger)
	if audioErr != nil {
		logger.Warn("Failed to create audio player, audio will be disabled", "error", audioErr)
		rm.audioPlayer = nil
	} else if audioPlayer == nil {
		logger.Debug("Audio player not available (no audio device), continuing without audio")
		rm.audioPlayer = nil
	} else {
		rm.audioPlayer = audioPlayer
	}

	if fontErr := rm.loadDefaultFont(ctx); fontErr != nil {
		logger.Error("failed to load default font", "error", fontErr)
	}
	return rm
}

// GetResourceCount returns the number of loaded resources
func (rm *ResourceManager) GetResourceCount() int {
	rm.mutex.RLock()
	defer rm.mutex.RUnlock()
	return len(rm.resources)
}

// GetResourceInfo returns information about loaded resources
func (rm *ResourceManager) GetResourceInfo() map[string]interface{} {
	rm.mutex.RLock()
	defer rm.mutex.RUnlock()

	info := make(map[string]interface{})
	for name, resource := range rm.resources {
		info[name] = map[string]interface{}{
			"type": resource.Type,
		}
	}
	return info
}

// GetAudioPlayer returns the audio player instance
func (rm *ResourceManager) GetAudioPlayer() *AudioPlayer {
	return rm.audioPlayer
}

// Cleanup releases all resources
func (rm *ResourceManager) Cleanup(ctx context.Context) error {
	// Check for cancellation
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	// Cleanup audio player
	if rm.audioPlayer != nil {
		rm.audioPlayer.Cleanup()
	}

	rm.logger.Info("Cleaning up resources", "count", len(rm.resources))
	rm.resources = make(map[string]*Resource)
	return nil
}

// Predefined sprite names for easy access
const (
	SpritePlayer     = "player"
	SpriteStar       = "star"
	SpriteButton     = "button"
	SpriteBackground = "background"
)
