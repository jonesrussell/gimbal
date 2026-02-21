package resources

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"

	"github.com/jonesrussell/gimbal/assets"
	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/errors"
)

const (
	disableAudioEnvValue = "1"
	disableAudioEnvTrue  = "true"
)

// AudioResource represents a loaded audio resource
type AudioResource struct {
	Data   []byte
	Length int64
}

// LoadAudio loads an audio resource from embedded assets
func (rm *ResourceManager) LoadAudio(ctx context.Context, name, path string) (*AudioResource, error) {
	// Check context cancellation
	if err := common.CheckContextCancellation(ctx); err != nil {
		return nil, err
	}

	// Check cache first
	if cached := rm.getCachedAudio(name); cached != nil {
		return cached, nil
	}

	// Load and decode audio
	return rm.loadAndCacheAudio(ctx, name, path)
}

// getCachedAudio retrieves audio from cache if exists
func (rm *ResourceManager) getCachedAudio(name string) *AudioResource {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	if resource, exists := rm.resources[name]; exists {
		if audioRes, ok := resource.Data.(*AudioResource); ok {
			return audioRes
		}
	}
	return nil
}

// loadAndCacheAudio loads audio from assets and caches it
func (rm *ResourceManager) loadAndCacheAudio(ctx context.Context, name, path string) (*AudioResource, error) {
	// Load from embedded assets
	audioData, err := rm.loadAudioFile(path)
	if err != nil {
		return nil, err
	}

	// Get audio length by decoding
	length, err := rm.decodeVorbisData(name, path, audioData)
	if err != nil {
		return nil, err
	}

	// Create audio resource (store raw data for later decoding)
	audioRes := &AudioResource{
		Data:   audioData,
		Length: length,
	}

	// Cache the audio resource
	rm.cacheAudio(name, audioRes)

	return audioRes, nil
}

// loadAudioFile loads audio file from embedded assets
func (rm *ResourceManager) loadAudioFile(path string) ([]byte, error) {
	audioData, err := assets.Assets.ReadFile(path)
	if err != nil {
		return nil, errors.NewGameErrorWithCause(errors.AssetLoadFailed, "failed to read audio file", err)
	}

	return audioData, nil
}

// decodeVorbisData decodes OGG/Vorbis data to get length information
func (rm *ResourceManager) decodeVorbisData(name, path string, audioData []byte) (int64, error) {
	// Create a reader from the audio data
	reader := bytes.NewReader(audioData)

	// Decode OGG/Vorbis stream to get length
	stream, err := vorbis.DecodeWithoutResampling(reader)
	if err != nil {
		return 0, errors.NewGameErrorWithCause(errors.AssetInvalid, "failed to decode audio", err)
	}

	length := stream.Length()
	return length, nil
}

// cacheAudio stores audio resource in the resource cache
func (rm *ResourceManager) cacheAudio(name string, audioRes *AudioResource) {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	if rm.resources == nil {
		rm.resources = make(map[string]*Resource)
	}

	rm.resources[name] = &Resource{
		Name: name,
		Type: ResourceSound,
		Data: audioRes,
	}
}

// GetAudio retrieves a loaded audio resource
func (rm *ResourceManager) GetAudio(ctx context.Context, name string) (*AudioResource, bool) {
	// Check for cancellation
	if err := common.CheckContextCancellation(ctx); err != nil {
		return nil, false
	}

	rm.mutex.RLock()
	defer rm.mutex.RUnlock()

	if resource, exists := rm.resources[name]; exists {
		if audioRes, ok := resource.Data.(*AudioResource); ok {
			return audioRes, true
		}
	}
	return nil, false
}

// getAudioConfigs returns the list of audio files to load
func getAudioConfigs() []struct {
	name string
	path string
} {
	return []struct {
		name string
		path string
	}{
		{
			name: "game_music_main",
			path: "sounds/game_music_main.ogg",
		},
		{
			name: "game_music_level_1",
			path: "sounds/game_music_level_1.ogg",
		},
		{
			name: "game_music_boss",
			path: "sounds/game_music_boss.ogg",
		},
		// UI audio cues
		{
			name: "ui_blip",
			path: "sounds/ui_blip.ogg",
		},
		{
			name: "stage_intro_whoosh",
			path: "sounds/stage_intro_whoosh.ogg",
		},
		{
			name: "boss_warning",
			path: "sounds/boss_warning.ogg",
		},
		{
			name: "warp_transition",
			path: "sounds/warp_transition.ogg",
		},
		{
			name: "game_over",
			path: "sounds/game_over.ogg",
		},
		{
			name: "continue_tick",
			path: "sounds/continue_tick.ogg",
		},
		{
			name: "victory_fanfare",
			path: "sounds/victory_fanfare.ogg",
		},
	}
}

// LoadAllAudio loads all required audio resources for the game
func (rm *ResourceManager) LoadAllAudio(ctx context.Context) error {
	// Skip audio loading if audio is disabled
	if isAudioDisabled() {
		return nil
	}

	// Check for cancellation at the start
	if err := common.CheckContextCancellation(ctx); err != nil {
		return err
	}

	audioConfigs := getAudioConfigs()
	for _, cfg := range audioConfigs {
		if _, err := rm.LoadAudio(ctx, cfg.name, cfg.path); err != nil {
			// Don't fail completely if audio fails to load
			continue
		}
	}

	return nil
}

// AudioPlayer manages audio playback
type AudioPlayer struct {
	audioContext *audio.Context
	players      map[string]*audio.Player
	audioData    map[string][]byte // Store raw audio data for looping
}

// isAudioDisabled checks if audio is disabled via environment variable
func isAudioDisabled() bool {
	val := os.Getenv("DISABLE_AUDIO")
	return val == disableAudioEnvValue || val == disableAudioEnvTrue
}

// tryCreateAudioContext attempts to create an audio context with panic recovery
func tryCreateAudioContext(sampleRate int) (*audio.Context, error) {
	var audioContext *audio.Context
	var initErr error

	func() {
		defer func() {
			if r := recover(); r != nil {
				audioContext = nil
				if err, ok := r.(error); ok {
					initErr = err
				} else {
					initErr = fmt.Errorf("audio context panic: %v", r)
				}
			}
		}()

		audioContext = audio.NewContext(sampleRate)
	}()

	return audioContext, initErr
}

// handleAudioInitFailure sets DISABLE_AUDIO so Ebiten does not retry audio init
func handleAudioInitFailure(initErr error) {
	if os.Getenv("DISABLE_AUDIO") == "" {
		_ = os.Setenv("DISABLE_AUDIO", disableAudioEnvValue)
	}
}

// NewAudioPlayer creates a new audio player
// Returns nil, nil if audio initialization fails (e.g., no audio device available)
func NewAudioPlayer(sampleRate int) (*AudioPlayer, error) {
	if isAudioDisabled() {
		return nil, nil
	}

	audioContext, initErr := tryCreateAudioContext(sampleRate)

	if audioContext == nil {
		handleAudioInitFailure(initErr)
		return nil, nil
	}

	return &AudioPlayer{
		audioContext: audioContext,
		players:      make(map[string]*audio.Player),
		audioData:    make(map[string][]byte),
	}, nil
}

// PlayMusic plays a background music track by name
func (ap *AudioPlayer) PlayMusic(name string, audioRes *AudioResource, volume float64) error {
	// Check if audio player is initialized
	if ap == nil || ap.audioContext == nil {
		return nil // Not an error - audio is optional
	}

	// Stop any currently playing music
	ap.StopMusic(name)

	// Store the audio data if not already stored
	if _, exists := ap.audioData[name]; !exists {
		ap.audioData[name] = audioRes.Data
	}

	audioData := ap.audioData[name]

	// Decode the audio first (this might take a moment for large files)
	decodedOnce, err := vorbis.DecodeWithoutResampling(bytes.NewReader(audioData))
	if err != nil {
		return errors.NewGameErrorWithCause(errors.AssetInvalid, "failed to decode audio", err)
	}

	streamLength := decodedOnce.Length()
	_ = streamLength

	// Read all data from decoded stream into a buffer
	decodedData, err := io.ReadAll(decodedOnce)
	if err != nil {
		return errors.NewGameErrorWithCause(errors.AssetLoadFailed, "failed to read decoded audio", err)
	}

	loopStream := audio.NewInfiniteLoop(bytes.NewReader(decodedData), int64(len(decodedData)))

	player, err := ap.audioContext.NewPlayer(loopStream)
	if err != nil {
		return err
	}

	player.SetVolume(volume)
	player.Play()

	ap.players[name] = player

	return nil
}

// StopMusic stops a music track by name
func (ap *AudioPlayer) StopMusic(name string) {
	if ap == nil || ap.players == nil {
		return
	}
	if player, exists := ap.players[name]; exists {
		player.Close()
		delete(ap.players, name)
	}
}

// StopAllMusic stops all playing music
func (ap *AudioPlayer) StopAllMusic() {
	if ap == nil || ap.players == nil {
		return
	}
	for name := range ap.players {
		ap.StopMusic(name)
	}
}

// IsPlaying checks if a music track is currently playing
func (ap *AudioPlayer) IsPlaying(name string) bool {
	if ap == nil || ap.players == nil {
		return false
	}
	if player, exists := ap.players[name]; exists {
		return player.IsPlaying()
	}
	return false
}

// SetVolume sets the volume for a music track
func (ap *AudioPlayer) SetVolume(name string, volume float64) {
	if ap == nil || ap.players == nil {
		return
	}
	if player, exists := ap.players[name]; exists {
		player.SetVolume(volume)
	}
}

// Cleanup releases all audio resources
func (ap *AudioPlayer) Cleanup() {
	if ap == nil {
		return
	}
	ap.StopAllMusic()
}
