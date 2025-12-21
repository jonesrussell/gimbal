package core

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
)

// GetPlayerEntries returns all entities tagged as players
func GetPlayerEntries(world donburi.World) []*donburi.Entry {
	var entries []*donburi.Entry
	query.NewQuery(
		filter.And(
			filter.Contains(PlayerTag),
			filter.Contains(Position),
			filter.Contains(Size),
		),
	).Each(world, func(entry *donburi.Entry) {
		entries = append(entries, entry)
	})
	return entries
}

// GetEnemyEntries returns all entities tagged as enemies
func GetEnemyEntries(world donburi.World) []*donburi.Entry {
	var entries []*donburi.Entry
	query.NewQuery(
		filter.And(
			filter.Contains(EnemyTag),
			filter.Contains(Position),
			filter.Contains(Size),
		),
	).Each(world, func(entry *donburi.Entry) {
		entries = append(entries, entry)
	})
	return entries
}

// GetProjectileEntries returns all player projectile entities
func GetProjectileEntries(world donburi.World) []*donburi.Entry {
	var entries []*donburi.Entry
	query.NewQuery(
		filter.And(
			filter.Contains(ProjectileTag),
			filter.Contains(Position),
			filter.Contains(Size),
		),
	).Each(world, func(entry *donburi.Entry) {
		entries = append(entries, entry)
	})
	return entries
}

// GetEnemyProjectileEntries returns all enemy projectile entities
func GetEnemyProjectileEntries(world donburi.World) []*donburi.Entry {
	var entries []*donburi.Entry
	query.NewQuery(
		filter.And(
			filter.Contains(EnemyProjectileTag),
			filter.Contains(Position),
			filter.Contains(Size),
		),
	).Each(world, func(entry *donburi.Entry) {
		entries = append(entries, entry)
	})
	return entries
}

// GetPlayerWithHealthEntries returns all player entities with health component
func GetPlayerWithHealthEntries(world donburi.World) []*donburi.Entry {
	var entries []*donburi.Entry
	query.NewQuery(
		filter.And(
			filter.Contains(PlayerTag),
			filter.Contains(Health),
		),
	).Each(world, func(entry *donburi.Entry) {
		entries = append(entries, entry)
	})
	return entries
}

// QueryEntriesWithFilters performs a query with the given layout filters
// This is a generic helper for custom queries not covered by the specific helpers above
func QueryEntriesWithFilters(world donburi.World, filters ...filter.LayoutFilter) []*donburi.Entry {
	var entries []*donburi.Entry
	query.NewQuery(filter.And(filters...)).Each(world, func(entry *donburi.Entry) {
		entries = append(entries, entry)
	})
	return entries
}
