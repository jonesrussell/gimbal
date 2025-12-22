package resources

import (
	"bytes"
	"context"
	"io"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"

	"github.com/jonesrussell/gimbal/assets"
	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/errors"
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
		rm.logger.Debug("[AUDIO_CACHE] Audio reused from cache", "name", name)
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
	rm.logger.Debug("[AUDIO_LOAD] Loading audio from embed", "name", name, "path", path)

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

	rm.logger.Debug("[AUDIO_LOAD] Audio loaded successfully", "name", name, "length", length)
	return audioRes, nil
}

// loadAudioFile loads audio file from embedded assets
func (rm *ResourceManager) loadAudioFile(path string) ([]byte, error) {
	audioData, err := assets.Assets.ReadFile(path)
	if err != nil {
		rm.logger.Error("[AUDIO_ERROR] Failed to read audio file", "path", path, "error", err)
		return nil, errors.NewGameErrorWithCause(errors.AssetLoadFailed, "failed to read audio file", err)
	}

	rm.logger.Debug("[AUDIO_LOAD] Audio file read successfully", "path", path, "size", len(audioData))
	return audioData, nil
}

// decodeVorbisData decodes OGG/Vorbis data to get length information
func (rm *ResourceManager) decodeVorbisData(name, path string, audioData []byte) (int64, error) {
	// Create a reader from the audio data
	reader := bytes.NewReader(audioData)

	// Decode OGG/Vorbis stream to get length
	stream, err := vorbis.DecodeWithoutResampling(reader)
	if err != nil {
		rm.logger.Error("[AUDIO_ERROR] Failed to decode Vorbis audio", "name", name, "path", path, "error", err)
		return 0, errors.NewGameErrorWithCause(errors.AssetInvalid, "failed to decode audio", err)
	}

	length := stream.Length()
	rm.logger.Debug("[AUDIO_DECODE] Audio decoded successfully", "name", name, "length", length)

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

// LoadAllAudio loads all required audio resources for the game
func (rm *ResourceManager) LoadAllAudio(ctx context.Context) error {
	// Check for cancellation at the start
	if err := common.CheckContextCancellation(ctx); err != nil {
		return err
	}

	audioConfigs := []struct {
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
	}

	for _, cfg := range audioConfigs {
		if _, err := rm.LoadAudio(ctx, cfg.name, cfg.path); err != nil {
			rm.logger.Warn("Failed to load audio, continuing without it", "name", cfg.name, "error", err)
			// Don't fail completely if audio fails to load
			continue
		}
	}

	rm.logger.Info("[AUDIO_LOAD] All audio resources loaded successfully")
	return nil
}

// AudioPlayer manages audio playback
type AudioPlayer struct {
	audioContext *audio.Context
	players      map[string]*audio.Player
	audioData    map[string][]byte // Store raw audio data for looping
	logger       common.Logger
}

// NewAudioPlayer creates a new audio player
// Returns nil, nil if audio initialization fails (e.g., no audio device available)
// This allows the game to run without audio in environments like containers
func NewAudioPlayer(sampleRate int, logger common.Logger) (*AudioPlayer, error) {
	// Try to create audio context - this may fail in containers without audio devices
	// We use a recover to catch any panics that might occur during audio initialization
	var audioContext *audio.Context
	func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Debug("Audio context creation panicked (no audio device available)", "panic", r)
				audioContext = nil
			}
		}()
		audioContext = audio.NewContext(sampleRate)
	}()

	if audioContext == nil {
		// Return nil without error - audio is optional
		logger.Debug("Audio context not available (no audio device), continuing without audio")
		return nil, nil
	}

	return &AudioPlayer{
		audioContext: audioContext,
		players:      make(map[string]*audio.Player),
		audioData:    make(map[string][]byte),
		logger:       logger,
	}, nil
}

// PlayMusic plays a background music track by name
func (ap *AudioPlayer) PlayMusic(name string, audioRes *AudioResource, volume float64) error {
	// Check if audio player is initialized
	if ap == nil || ap.audioContext == nil {
		if ap != nil && ap.logger != nil {
			ap.logger.Debug("Audio player not available, skipping music playback", "name", name)
		}
		return nil // Not an error - audio is optional
	}

	ap.logger.Debug("PlayMusic: Stopping any existing music", "name", name)
	// Stop any currently playing music
	ap.StopMusic(name)

	ap.logger.Debug("PlayMusic: Storing audio data", "name", name, "data_size", len(audioRes.Data))
	// Store the audio data if not already stored
	if _, exists := ap.audioData[name]; !exists {
		ap.audioData[name] = audioRes.Data
	}

	// Get the stored audio data
	audioData := ap.audioData[name]
	ap.logger.Debug("PlayMusic: Decoding audio first", "name", name, "data_size", len(audioData))

	// Decode the audio first (this might take a moment for large files)
	decodedOnce, err := vorbis.DecodeWithoutResampling(bytes.NewReader(audioData))
	if err != nil {
		ap.logger.Error("Failed to decode audio", "name", name, "error", err)
		return errors.NewGameErrorWithCause(errors.AssetInvalid, "failed to decode audio", err)
	}

	// Get the length of the decoded stream
	streamLength := decodedOnce.Length()
	ap.logger.Debug("PlayMusic: Audio decoded", "name", name, "stream_length", streamLength)

	// Create an infinite loop from the decoded stream
	// We need to read the decoded data and create a loop from it
	ap.logger.Debug("PlayMusic: Reading decoded audio data", "name", name)
	// Read all data from decoded stream into a buffer
	decodedData, err := io.ReadAll(decodedOnce)
	if err != nil {
		ap.logger.Error("Failed to read decoded audio data", "name", name, "error", err)
		return errors.NewGameErrorWithCause(errors.AssetLoadFailed, "failed to read decoded audio", err)
	}

	// Create infinite loop from decoded data
	ap.logger.Debug("PlayMusic: Creating infinite loop from decoded stream",
		"name", name, "decoded_size", len(decodedData))
	loopStream := audio.NewInfiniteLoop(bytes.NewReader(decodedData), int64(len(decodedData)))

	ap.logger.Debug("PlayMusic: Creating audio player", "name", name)
	// Create a new player from the looping stream
	player, err := ap.audioContext.NewPlayer(loopStream)
	if err != nil {
		ap.logger.Error("Failed to create audio player", "name", name, "error", err)
		return errors.NewGameErrorWithCause(errors.SystemInitFailed, "failed to create audio player", err)
	}

	ap.logger.Debug("PlayMusic: Setting volume and playing", "name", name, "volume", volume)
	// Set volume (0.0 to 1.0)
	player.SetVolume(volume)

	// Play the music (it will loop automatically)
	player.Play()

	ap.players[name] = player
	ap.logger.Debug("Music started playing", "name", name, "volume", volume)

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
		if ap.logger != nil {
			ap.logger.Debug("Music stopped", "name", name)
		}
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
	if ap.logger != nil {
		ap.logger.Debug("Audio player cleaned up")
	}
}
