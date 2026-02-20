# Invincibility (God Mode + Power-Up) Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Add developer God mode (CLI + pause-menu toggle when Debug) and player invincibility power-up as a third power-up type, with damage priority: dev → power-up/i-frames → apply damage.

**Architecture:** Dev invincibility remains `config.Invincible` with new `SetDevInvincible(bool)` for pause toggle; HealthSystem keeps current check order. Player invincibility is a new `PowerUpInvincibility` type that sets `Health.IsInvincible` and `InvincibilityTime` on collect; same spawn/collect pipeline as existing power-ups.

**Tech Stack:** Go, Ebiten, Donburi, existing config/health/powerup systems.

**Design reference:** `docs/plans/2026-02-19-invincibility-design.md`

---

## Task 1: Config — SetDevInvincible

**Files:**
- Modify: `internal/config/config.go` (after `WithInvincible`, before `DefaultConfig`)
- Create or modify: `internal/config/config_test.go` (add test for SetDevInvincible; create file if it does not exist)

**Step 1: Add test for SetDevInvincible**

In `internal/config/config_test.go` add (or create file if missing):

```go
func TestSetDevInvincible_OnlyWhenDebug(t *testing.T) {
	// When Debug is false, SetDevInvincible does not set Invincible
	cfg := NewConfig(WithDebug(false))
	cfg.SetDevInvincible(true)
	if cfg.Invincible {
		t.Error("Invincible should stay false when Debug is false")
	}
	// When Debug is true, SetDevInvincible sets Invincible
	cfg2 := NewConfig(WithDebug(true))
	cfg2.SetDevInvincible(true)
	if !cfg2.Invincible {
		t.Error("Invincible should be true after SetDevInvincible(true) when Debug is true")
	}
	cfg2.SetDevInvincible(false)
	if cfg2.Invincible {
		t.Error("Invincible should be false after SetDevInvincible(false)")
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/config/... -run TestSetDevInvincible -v`  
Expected: FAIL (SetDevInvincible not defined or Invincible not updated)

**Step 3: Implement SetDevInvincible**

In `internal/config/config.go`, after the `WithInvincible` function (around line 96), add:

```go
// SetDevInvincible sets Invincible for runtime toggle (e.g. pause menu).
// Only takes effect when Debug is true; no-op when Debug is false.
func (c *GameConfig) SetDevInvincible(invincible bool) {
	if c.Debug {
		c.Invincible = invincible
	}
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./internal/config/... -run TestSetDevInvincible -v`  
Expected: PASS

**Step 5: Commit**

```bash
git add internal/config/config.go internal/config/config_test.go
git commit -m "config: add SetDevInvincible for runtime God mode toggle"
```

---

## Task 2: Pause menu — God mode option when Debug

**Files:**
- Modify: `internal/scenes/pause/paused.go` (options built in Enter; NewPausedScene needs to build initial options or defer to Enter)
- Test: Manual or existing scene test; optional unit test for option count when Debug

**Context:** Pause menu currently builds a fixed list of 3 options in `NewPausedScene`. To show "God mode: ON/OFF" when Debug, build the options list in `Enter()` so the label reflects current `config.Invincible` and the toggle action is correct. That requires the pause scene to create/update the menu when entering.

**Step 1: Build options in Enter when Debug**

- In `PausedScene`, store a reference to the options-building logic so Enter can refresh the menu.
- Option A: In `NewPausedScene`, create the menu with a base list; add a method `RefreshPauseOptions()` that rebuilds options (Resume, Return to Menu, [if config.Debug] God mode toggle, Quit) and replaces the menu’s options if the menu package supports it.
- Option B: In `Enter()`, build the full options slice (Resume, Return to Menu; if `s.manager.GetConfig().Debug` append God mode option with text "God mode: ON" or "God mode: OFF" from `s.manager.GetConfig().Invincible`, action calling `s.manager.GetConfig().SetDevInvincible(!s.manager.GetConfig().Invincible)`; append Quit). Then create a new `menu.NewMenuSystem(options, ...)` and assign to `s.menu`.

Implement Option B: in `Enter()`, build `options` slice with conditional God mode, then `s.menu = menu.NewMenuSystem(options, &config, ...)` so the menu is fresh each time with correct label and action.

**Step 2: Keep NewPausedScene working**

- NewPausedScene must still create a valid PausedScene. Build the same options list in NewPausedScene (with optional God mode when `manager.GetConfig().Debug`) so the initial `s.menu` is valid. Then in Enter(), rebuild and replace `s.menu` so every time the player opens pause they see the current God mode state.

**Step 3: Implement**

In `internal/scenes/pause/paused.go`:

- Add a helper that returns the pause menu options for the current config, e.g. `func (s *PausedScene) buildPauseOptions() []menu.MenuOption`. It returns Resume, Return to Menu; if `s.manager.GetConfig().Debug` then append `{Text: "God mode: ON" or "God mode: OFF", Action: toggle}`; append Quit.
- In `NewPausedScene`, call that helper (you’ll need to pass manager and get config; if the helper is on *PausedScene, build a minimal list in NewPausedScene for the initial menu, then in Enter() call `s.buildPauseOptions()` and create menu from that). Simplest: in NewPausedScene, build options with a function that takes `*scenes.SceneManager` and returns `[]menu.MenuOption`, and use it for both NewPausedScene and Enter. So in NewPausedScene set `scene.menu = menu.NewMenuSystem(buildPauseOptions(manager), ...)`. In Enter do `s.menu = menu.NewMenuSystem(buildPauseOptions(s.manager), &config, ...)` so the menu is refreshed with current God mode state. Extract `buildPauseOptions(manager *scenes.SceneManager) []menu.MenuOption` as a package-level or PausedScene method.
- Implement `buildPauseOptions`: base options Resume, Return to Menu; if `manager.GetConfig().Debug` append God mode (text from `manager.GetConfig().Invincible`, action `manager.GetConfig().SetDevInvincible(!manager.GetConfig().Invincible)`); append Quit.

**Step 4: Verify**

Run game with `DEBUG=true`, enter play, pause, confirm fourth option "God mode: OFF" (or ON if already on). Toggle and confirm it flips. Unpause and verify player does not take damage when God mode is ON.

**Step 5: Commit**

```bash
git add internal/scenes/pause/paused.go
git commit -m "pause: add God mode toggle when Debug is enabled"
```

---

## Task 3: Core — PowerUpInvincibility type

**Files:**
- Modify: `internal/ecs/core/components.go` (PowerUpType const block, ~line 187)

**Step 1: Add constant**

In `internal/ecs/core/components.go`, in the `PowerUpType` const block, add:

```go
PowerUpInvincibility  // Temporary invincibility (duration from PowerUpTypeData)
```

So the block becomes:

```go
const (
	PowerUpDoubleShot PowerUpType = iota
	PowerUpExtraLife
	PowerUpInvincibility
)
```

**Step 2: Run build**

Run: `go build ./...`  
Expected: success (no other code references the new value yet).

**Step 3: Commit**

```bash
git add internal/ecs/core/components.go
git commit -m "ecs: add PowerUpInvincibility type"
```

---

## Task 4: Power-up system — sprite and collect

**Files:**
- Modify: `internal/ecs/systems/powerup/powerup_system.go`

**Step 1: Add sprite in createSprites**

In `createSprites()`, add after extra life sprite:

```go
invincibleSprite := ebiten.NewImage(PowerUpSize, PowerUpSize)
invincibleSprite.Fill(color.RGBA{R: 255, G: 255, B: 255, A: 255}) // white
ps.sprites[core.PowerUpInvincibility] = invincibleSprite
```

**Step 2: Handle collect in collectPowerUp**

In `collectPowerUp`, add case in the switch:

```go
case core.PowerUpInvincibility:
	duration := data.Duration
	if duration <= 0 {
		duration = 5 * time.Second
	}
	if playerEntry.HasComponent(core.Health) {
		health := core.Health.Get(playerEntry)
		health.IsInvincible = true
		if duration > health.InvincibilityTime {
			health.InvincibilityTime = duration
		}
		core.Health.SetValue(playerEntry, *health)
	}
	ps.logger.Info("Invincibility power-up collected", "duration", duration)
```

**Step 3: Run build and tests**

Run: `go build ./...` and `go test ./internal/ecs/systems/powerup/... -v`  
Expected: build OK; tests pass (or add a small test that collects invincibility and checks Health if desired).

**Step 4: Commit**

```bash
git add internal/ecs/systems/powerup/powerup_system.go
git commit -m "powerup: add invincibility sprite and collection effect"
```

---

## Task 5: Power-up system — SpawnPowerUp and TrySpawnPowerUp

**Files:**
- Modify: `internal/ecs/systems/powerup/powerup_system.go`

**Step 1: Default duration in SpawnPowerUp**

In `SpawnPowerUp`, where `duration` is set (e.g. for double shot), set duration for invincibility:

```go
duration := time.Duration(0)
switch powerUpType {
case core.PowerUpDoubleShot:
	duration = 10 * time.Second
case core.PowerUpInvincibility:
	duration = 5 * time.Second
}
```

Use the same `core.PowerUpData.SetValue(entry, core.PowerUpTypeData{..., Duration: duration, ...})` so the payload has the correct duration.

**Step 2: Add to TrySpawnPowerUp weighted random**

In `TrySpawnPowerUp`, change the weighted random to three-way. Example (tunable): 60% double, 25% extra life, 15% invincibility:

```go
r := rand.Float64()
if r < 0.60 {
	ps.SpawnPowerUp(position, core.PowerUpDoubleShot)
} else if r < 0.85 {
	ps.SpawnPowerUp(position, core.PowerUpExtraLife)
} else {
	ps.SpawnPowerUp(position, core.PowerUpInvincibility)
}
```

**Step 3: Run build and tests**

Run: `go build ./...` and `go test ./internal/ecs/systems/powerup/...`  
Expected: pass.

**Step 4: Commit**

```bash
git add internal/ecs/systems/powerup/powerup_system.go
git commit -m "powerup: spawn invincibility with duration and drop weight"
```

---

## Task 6: Health system — confirm damage order

**Files:**
- Modify: `internal/ecs/systems/health/system.go` (review only; no change if order already correct)

**Step 1: Verify order in DamageEntity**

Open `DamageEntity`. Confirm the first check is `if hs.config.Debug && hs.config.Invincible` for player (return nil). Second is `if health.IsInvincible` (return nil). Then apply damage and set i-frames. No code change if already in this order.

**Step 2: Add unit test for damage order (optional but recommended)**

In `internal/ecs/systems/health/system_test.go` (or create), add test that when config.Debug and config.Invincible are true, DamagePlayer does not reduce health. And when Health.IsInvincible is true, DamageEntity returns without applying damage. This locks in the priority.

**Step 3: Commit**

```bash
git add internal/ecs/systems/health/system.go internal/ecs/systems/health/system_test.go
git commit -m "health: verify damage priority (dev invincible then i-frames)"
```

---

## Task 7: Docs and final check

**Files:**
- Modify: `CLAUDE.md` or game-mechanics rule if invincibility is documented there
- Design doc already in `docs/plans/2026-02-19-invincibility-design.md`

**Step 1: Update CLAUDE.md (optional)**

If CLAUDE.md mentions invincibility or power-ups, add a line that player invincibility is a power-up type and dev invincibility is CLI/pause when Debug.

**Step 2: Full test and run**

Run: `task lint`, `go test ./...`, then run game with `DEBUG=true -invincible`, run with DEBUG and toggle God mode in pause, run and collect an invincibility power-up and verify no damage for duration.

**Step 3: Commit**

```bash
git add CLAUDE.md  # if changed
git commit -m "docs: note invincibility modes in CLAUDE.md"
```

---

## Execution summary

- **Task 1:** Config setter + test  
- **Task 2:** Pause menu God mode option when Debug  
- **Task 3:** New PowerUpType constant  
- **Task 4:** Power-up sprite and collect effect  
- **Task 5:** Spawn and drop weight for invincibility  
- **Task 6:** Health damage order verification (+ optional test)  
- **Task 7:** Docs and full verification  

Plan complete and saved to `docs/plans/2026-02-19-invincibility-implementation.md`.
