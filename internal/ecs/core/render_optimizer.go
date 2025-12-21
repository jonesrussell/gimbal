package core

import (
	"fmt"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/config"
)

// RenderBatch represents a batch of entities to render together
type RenderBatch struct {
	Sprite   *ebiten.Image
	Entities []*donburi.Entry
	ZOrder   int // Lower numbers render first (background)
}

// RenderOptimizer optimizes rendering performance through batching and culling
type RenderOptimizer struct {
	config     *config.GameConfig
	screenSize config.Size
	batches    map[string]*RenderBatch // Key: sprite hash
	cullMargin float64
}

// NewRenderOptimizer creates a new render optimization system with the provided configuration
func NewRenderOptimizer(cfg *config.GameConfig) *RenderOptimizer {
	return &RenderOptimizer{
		config:     cfg,
		screenSize: cfg.ScreenSize,
		batches:    make(map[string]*RenderBatch),
		cullMargin: 100.0, // Cull objects 100 pixels outside screen
	}
}

// OptimizedRenderSystem renders entities with batching and culling
func (ro *RenderOptimizer) OptimizedRenderSystem(world donburi.World, screen *ebiten.Image) {
	// Clear previous batches
	ro.batches = make(map[string]*RenderBatch)

	// Collect and batch entities
	ro.collectEntities(world)

	// Render batches in Z-order
	ro.renderBatches(screen)
}

// collectEntities collects entities and organizes them into batches
func (ro *RenderOptimizer) collectEntities(world donburi.World) {
	query.NewQuery(
		filter.And(
			filter.Contains(Position),
			filter.Contains(Sprite),
		),
	).Each(world, func(entry *donburi.Entry) {
		pos := Position.Get(entry)
		sprite := Sprite.Get(entry)

		if pos == nil || sprite == nil {
			return
		}

		// Frustum culling - skip off-screen entities
		if ro.isOffScreen(*pos, *sprite) {
			return
		}

		// Create batch key based on sprite and rendering properties
		batchKey := ro.createBatchKey(entry, *sprite)

		// Get or create batch
		batch, exists := ro.batches[batchKey]
		if !exists {
			batch = &RenderBatch{
				Sprite:   *sprite,
				Entities: make([]*donburi.Entry, 0),
				ZOrder:   ro.getZOrder(entry),
			}
			ro.batches[batchKey] = batch
		}

		batch.Entities = append(batch.Entities, entry)
	})
}

// isOffScreen checks if an entity is outside the visible area
func (ro *RenderOptimizer) isOffScreen(pos common.Point, sprite *ebiten.Image) bool {
	bounds := sprite.Bounds()
	width := float64(bounds.Dx())
	height := float64(bounds.Dy())

	// Check if entity is outside screen bounds with margin
	return pos.X+width < -ro.cullMargin ||
		pos.X > float64(ro.screenSize.Width)+ro.cullMargin ||
		pos.Y+height < -ro.cullMargin ||
		pos.Y > float64(ro.screenSize.Height)+ro.cullMargin
}

// createBatchKey creates a unique key for batching based on sprite and rendering properties
func (ro *RenderOptimizer) createBatchKey(entry *donburi.Entry, sprite *ebiten.Image) string {
	// Base key includes sprite pointer to ensure different sprites get different batches
	// This is critical - entities with the same sprite bounds but different sprites must be in different batches
	key := fmt.Sprintf("%p_%s", sprite, sprite.Bounds().String())

	// Add scale information
	if entry.HasComponent(Scale) {
		scale := Scale.Get(entry)
		key += "_scale_" + string(rune(int(*scale*100)))
	}

	// Add rotation information
	if entry.HasComponent(Orbital) {
		key += "_rotated"
	} else if entry.HasComponent(Angle) {
		key += "_rotated"
	}

	// Add invincibility state
	if entry.HasComponent(Health) {
		health := Health.Get(entry)
		if health.IsInvincible {
			key += "_invincible"
		}
	}

	return key
}

// getZOrder determines the rendering order for an entity
func (ro *RenderOptimizer) getZOrder(entry *donburi.Entry) int {
	// Background elements (stars, etc.)
	if entry.HasComponent(StarTag) {
		return 0
	}

	// Player and enemies
	if entry.HasComponent(PlayerTag) {
		return 10
	}
	if entry.HasComponent(EnemyTag) {
		return 10
	}

	// Projectiles
	if entry.HasComponent(ProjectileTag) {
		return 20
	}

	// UI elements
	return 30
}

// renderBatches renders all batches in Z-order
func (ro *RenderOptimizer) renderBatches(screen *ebiten.Image) {
	// Sort batches by Z-order
	batchList := make([]*RenderBatch, 0, len(ro.batches))
	for _, batch := range ro.batches {
		batchList = append(batchList, batch)
	}

	// Simple insertion sort for small number of batches
	for i := 1; i < len(batchList); i++ {
		key := batchList[i]
		j := i - 1
		for j >= 0 && batchList[j].ZOrder > key.ZOrder {
			batchList[j+1] = batchList[j]
			j--
		}
		batchList[j+1] = key
	}

	// Render each batch
	for _, batch := range batchList {
		ro.renderBatch(screen, batch)
	}
}

// renderBatch renders a single batch of entities
func (ro *RenderOptimizer) renderBatch(screen *ebiten.Image, batch *RenderBatch) {
	for _, entry := range batch.Entities {
		pos := Position.Get(entry)
		if pos == nil {
			continue
		}

		op := &ebiten.DrawImageOptions{}
		ApplySpriteTransform(entry, batch.Sprite, op)

		// Apply invincibility flashing
		if entry.HasComponent(Health) {
			health := Health.Get(entry)
			if health.IsInvincible {
				flashRate := 0.2 * float64(time.Second) // Convert to time.Duration
				remainingTime := health.InvincibilityDuration - health.InvincibilityTime
				flashPhase := int(remainingTime / time.Duration(flashRate))
				if flashPhase%2 == 0 {
					op.ColorScale.SetR(1)
					op.ColorScale.SetG(1)
					op.ColorScale.SetB(1)
					op.ColorScale.SetA(0.5)
				}
			}
		}

		// Apply position translation
		op.GeoM.Translate(pos.X, pos.Y)
		screen.DrawImage(batch.Sprite, op)
	}
}
