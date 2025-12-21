package game

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/jonesrussell/gimbal/internal/ecs/core"
	enemysys "github.com/jonesrussell/gimbal/internal/ecs/systems/enemy"
)

// drawBossDebugInfo draws boss debug information
func (g *ECSGame) drawBossDebugInfo(screen *ebiten.Image, x, screenHeight, lineHeight float64) {
	bossEntry := g.findBossEntity()
	if bossEntry == nil {
		g.drawBossStatus(screen, x, screenHeight, lineHeight)
		return
	}
	g.drawBossDetails(screen, bossEntry, x, screenHeight, lineHeight)
}

// findBossEntity finds the boss entity in the world
func (g *ECSGame) findBossEntity() *donburi.Entry {
	var bossEntry *donburi.Entry
	query.NewQuery(
		filter.And(
			filter.Contains(core.EnemyTag),
			filter.Contains(core.EnemyTypeID),
		),
	).Each(g.world, func(entry *donburi.Entry) {
		typeID := core.EnemyTypeID.Get(entry)
		if enemysys.EnemyType(*typeID) == enemysys.EnemyTypeBoss {
			bossEntry = entry
		}
	})
	return bossEntry
}

// drawBossStatus draws boss spawn/defeat status
func (g *ECSGame) drawBossStatus(screen *ebiten.Image, x, screenHeight, lineHeight float64) {
	if g.enemySystem.WasBossSpawned() {
		g.drawDebugText(screen, "Boss: Defeated", x, screenHeight-lineHeight)
	} else {
		g.drawDebugText(screen, "Boss: Spawning soon...", x, screenHeight-lineHeight)
	}
}

// drawBossDetails draws detailed boss information
func (g *ECSGame) drawBossDetails(screen *ebiten.Image, bossEntry *donburi.Entry, x, screenHeight, lineHeight float64) {
	pos := core.Position.Get(bossEntry)
	health := core.Health.Get(bossEntry)
	orbital := core.Orbital.Get(bossEntry)
	size := core.Size.Get(bossEntry)

	// Calculate number of lines for boss info
	numLines := 6 // Boss, Health, Position, Orbital Angle, Size, Status
	startY := screenHeight - float64(numLines)*lineHeight - 20

	// Draw boss information from bottom up
	y := startY
	g.drawDebugText(screen, "BOSS", x, y)
	y += lineHeight

	if health != nil {
		healthPercent := float64(health.Current) / float64(health.Maximum) * 100
		healthText := fmt.Sprintf("Health: %d/%d (%.0f%%)",
			health.Current, health.Maximum, healthPercent)
		g.drawDebugText(screen, healthText, x, y)
	} else {
		g.drawDebugText(screen, "Health: Unknown", x, y)
	}
	y += lineHeight

	if pos != nil {
		g.drawDebugText(screen, fmt.Sprintf("Position: (%.0f, %.0f)", pos.X, pos.Y), x, y)
	} else {
		g.drawDebugText(screen, "Position: Unknown", x, y)
	}
	y += lineHeight

	if orbital != nil {
		g.drawDebugText(screen, fmt.Sprintf("Orbital Angle: %.1f°", float64(orbital.OrbitalAngle)), x, y)
	} else {
		g.drawDebugText(screen, "Orbital: Unknown", x, y)
	}
	y += lineHeight

	if size != nil {
		g.drawDebugText(screen, fmt.Sprintf("Size: %dx%d", size.Width, size.Height), x, y)
	} else {
		g.drawDebugText(screen, "Size: Unknown", x, y)
	}
	y += lineHeight

	status := "Active"
	if health != nil && health.Current <= 0 {
		status = "Defeated"
	}
	g.drawDebugText(screen, fmt.Sprintf("Status: %s", status), x, y)
}
