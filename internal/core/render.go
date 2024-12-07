package core

import (
	"fmt"
	"image"
	"sort"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"go.uber.org/zap"
)

// Implementation-specific types
type (
	sprite struct {
		image *ebiten.Image
		pos   image.Point
		z     int
	}

	rendererImpl struct {
		sprites []sprite
		cache   assetCache
		mu      sync.RWMutex
		logger  *zap.Logger
		config  *RenderConfig
	}
)

// Implementation of NewRenderer constructor type from types.go
var _ NewRenderer = NewRendererImpl

func NewRendererImpl(logger *zap.Logger, opts ...RenderOption) (Renderer, error) {
	if logger == nil {
		return nil, fmt.Errorf("logger is required")
	}

	config := &RenderConfig{
		VSync:        true,
		DoubleBuffer: true,
		MaxSprites:   1000,
		Debug:        false,
	}

	for _, opt := range opts {
		opt(config)
	}

	return &rendererImpl{
		logger:  logger,
		sprites: make([]sprite, 0, config.MaxSprites),
		config:  config,
	}, nil
}

func (r *rendererImpl) Draw(screen *ebiten.Image) error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if screen == nil {
		return fmt.Errorf("screen is nil")
	}

	// Sort sprites by z-index
	sort.SliceStable(r.sprites, func(i, j int) bool {
		return r.sprites[i].z < r.sprites[j].z
	})

	// Draw all sprites
	for _, s := range r.sprites {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(s.pos.X), float64(s.pos.Y))
		screen.DrawImage(s.image, op)
	}

	if r.config.Debug {
		r.drawDebugInfo(screen)
	}

	return nil
}

func (r *rendererImpl) AddSprite(id string, spriteImg *ebiten.Image, pos image.Point, z int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if spriteImg == nil {
		return fmt.Errorf("sprite image is nil for id: %s", id)
	}

	if len(r.sprites) >= r.config.MaxSprites {
		return fmt.Errorf("max sprites limit reached (%d)", r.config.MaxSprites)
	}

	r.sprites = append(r.sprites, struct {
		image *ebiten.Image
		pos   image.Point
		z     int
	}{
		image: spriteImg,
		pos:   pos,
		z:     z,
	})

	r.logger.Debug("added sprite",
		zap.String("id", id),
		zap.Int("x", pos.X),
		zap.Int("y", pos.Y),
		zap.Int("z", z))

	return nil
}

func (r *rendererImpl) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.sprites = r.sprites[:0]
	r.logger.Debug("cleared sprite buffer")
}

func (r *rendererImpl) drawDebugInfo(screen *ebiten.Image) {
	// TODO: Implement debug rendering
	// - FPS counter
	// - Sprite count
	// - Memory usage
	// - Other debug metrics
}

// Renderer options
func WithVSync(enable bool) RenderOption {
	return func(c *RenderConfig) {
		c.VSync = enable
	}
}

func WithDoubleBuffer(enable bool) RenderOption {
	return func(c *RenderConfig) {
		c.DoubleBuffer = enable
	}
}

func WithMaxSprites(maxCount int) RenderOption {
	return func(c *RenderConfig) {
		c.MaxSprites = maxCount
	}
}

func WithRenderDebug(debug bool) RenderOption {
	return func(c *RenderConfig) {
		c.Debug = debug
	}
}
