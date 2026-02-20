# Game-Loop Debug Logging (dbg) Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Add a purpose-built `internal/dbg` package for at-a-glance, category-based debug output in the game loop; add frame marker and breadcrumbs; add one-frame trace mode; remove Zap and ad-hoc fmt from the loop; use `log.Printf` for app-level output.

**Architecture:** Single `internal/dbg` package with categories (Event, State, System, World, Spawn), atomic Enable/Disable, and `Log(cat, msg, args...)`. Game loop and all systems invoked every frame call only `dbg.Log` for debug output. Trace mode: flag set by `dbg.Trace()`, enable at start of Update and disable at end for one frame.

**Tech Stack:** Go stdlib (fmt, sync/atomic, log). No Zap in loop; Zap removed from project; app-level uses log.

**Design reference:** `docs/plans/2026-02-20-dbg-logging-design.md`

---

## Task 1: Create internal/dbg package (core API)

**Files:**
- Create: `internal/dbg/log.go`
- Create: `internal/dbg/log_test.go`

**Step 1: Write failing tests**

In `internal/dbg/log_test.go`:

```go
package dbg

import (
	"bytes"
	"os"
	"testing"
)

func TestLog_outputFormat(t *testing.T) {
	Enable()
	defer Enable() // leave enabled for other tests
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	defer func() { os.Stdout = old }()

	Log(Event, "test %s", "message")
	w.Close()
	var buf bytes.Buffer
	buf.ReadFrom(r)
	out := buf.String()
	if out != "[EVENT] test message\n" {
		t.Errorf("Log output = %q, want [EVENT] test message\n", out)
	}
}

func TestLog_respectsDisable(t *testing.T) {
	Disable()
	defer Enable()
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	defer func() { os.Stdout = old }()

	Log(Event, "should not appear")
	w.Close()
	var buf bytes.Buffer
	buf.ReadFrom(r)
	if buf.Len() != 0 {
		t.Errorf("Log should produce no output when disabled, got %q", buf.String())
	}
}
```

**Step 2: Run tests (expect fail)**

Run: `go test ./internal/dbg/... -v`  
Expected: FAIL (package dbg or Log not defined)

**Step 3: Implement dbg package**

Create `internal/dbg/log.go`:

```go
package dbg

import (
	"fmt"
	"sync/atomic"
)

type Category string

const (
	Event  Category = "EVENT"
	State  Category = "STATE"
	System Category = "SYSTEM"
	World  Category = "WORLD"
	Spawn  Category = "SPAWN"
)

var enabled uint32 = 1

func Enable()  { atomic.StoreUint32(&enabled, 1) }
func Disable() { atomic.StoreUint32(&enabled, 0) }

func Log(cat Category, msg string, args ...any) {
	if atomic.LoadUint32(&enabled) == 0 {
		return
	}
	fmt.Printf("[%s] %s\n", cat, fmt.Sprintf(msg, args...))
}
```

**Step 4: Run tests (expect pass)**

Run: `go test ./internal/dbg/... -v`  
Expected: PASS

**Step 5: Commit**

```bash
git add internal/dbg/log.go internal/dbg/log_test.go
git commit -m "feat: add internal/dbg package with categories and atomic enable/disable"
```

---

## Task 2: Frame marker in gameplay update

**Files:**
- Modify: `internal/game/game_systems.go` (top of updateGameplaySystems, after isPlayingScene check)

**Step 1: Add frame marker**

In `updateGameplaySystems`, after the `if !isPlayingScene { return nil }` block, add:

```go
import "github.com/jonesrussell/gimbal/internal/dbg"

// After: if !isPlayingScene { return nil }
if g.frameCount%30 == 0 {
	dbg.Log(dbg.System, "=== FRAME %d ===", g.frameCount)
}
```

Remove the existing `fmt.Printf("Gameplay using EventSystem %p\n", g.eventSystem)` line (replaced by a World breadcrumb in a later task if desired, or drop).

**Step 2: Verify**

Run: `go build ./...`  
Run game with DEBUG, enter play; in terminal you should see `[SYSTEM] === FRAME 30 ===` etc. every 30 frames.

**Step 3: Commit**

```bash
git add internal/game/game_systems.go
git commit -m "game: add dbg frame marker every 30 frames in gameplay update"
```

---

## Task 3: Breadcrumbs — init and gameplay (World, System)

**Files:**
- Modify: `internal/game/init_systems.go` (replace fmt.Printf and add dbg.Log where useful)
- Modify: `internal/game/game_systems.go` (add dbg.Log at start of updateGameplaySystems)

**Step 1: init_systems.go**

- Add import: `"github.com/jonesrussell/gimbal/internal/dbg"`
- Remove: `fmt.Printf("EventSystem created at %p\n", g.eventSystem)` and `fmt.Printf("StageStateMachine using EventSystem %p\n", g.eventSystem)`.
- Add once after EventSystem creation: `dbg.Log(dbg.World, "EventSystem created at %p", g.eventSystem)`.
- Add once after StageStateMachine creation: `dbg.Log(dbg.World, "StageStateMachine using EventSystem %p", g.eventSystem)`.
- Leave or remove other `g.logger.Debug` in this file for now (migrated in Task 8).

**Step 2: game_systems.go**

- At the very start of the gameplay block (after frame marker), add: `dbg.Log(dbg.System, "updateGameplaySystems start")`.
- Ensure `internal/dbg` is imported.

**Step 3: Build and run**

Run: `go build ./...`  
Expected: build OK. Run game and confirm one-off World lines at startup and System line each gameplay frame (or every 30 if you only log when frame marker fires; design says one line at start of update — so every frame is fine for "start", or only when frameCount%30==0 to reduce noise; design says "at start of gameplay update", so add it unconditionally; frame marker is every 30, so "start" every frame may be noisy — design doc says "at start of gameplay update" so we add it; user can disable dbg if too noisy).

**Step 4: Commit**

```bash
git add internal/game/init_systems.go internal/game/game_systems.go
git commit -m "game: add dbg World and System breadcrumbs in init and gameplay"
```

---

## Task 4: Breadcrumbs — GyrussSystem (Event, Spawn)

**Files:**
- Modify: `internal/ecs/systems/gyruss/gyruss_system.go`

**Step 1: Add import**

`"github.com/jonesrussell/gimbal/internal/dbg"`

**Step 2: EmitBossDefeated**

- Remove: `fmt.Println("GyrussSystem: EmitBossDefeated called")`.
- Add: `dbg.Log(dbg.Event, "EmitBossDefeated (world=%p eventSystem=%p)", world, gs.eventSystem)` (or equivalent; pass world and eventSystem into the function that emits, or use gs fields if available). Check gyruss_system.go for the exact signature and add the line where the event is emitted.

**Step 3: Spawn (optional)**

- Where boss is spawned, add: `dbg.Log(dbg.Spawn, "boss spawned")`. Where stage/wave is loaded, you can add `dbg.Log(dbg.Spawn, "stage loaded")` or similar if useful. Prefer one line at boss spawn and one at EmitBossDefeated.

**Step 4: Build**

Run: `go build ./...`

**Step 5: Commit**

```bash
git add internal/ecs/systems/gyruss/gyruss_system.go
git commit -m "gyruss: add dbg Event and Spawn breadcrumbs"
```

---

## Task 5: Breadcrumbs — StageStateMachine (Event, State)

**Files:**
- Modify: `internal/ecs/systems/stage/stage_state_machine.go`

**Step 1: Add import**

`"github.com/jonesrussell/gimbal/internal/dbg"`

**Step 2: onBossDefeated**

- Remove: `fmt.Println("StageStateMachine: onBossDefeated fired")` and the duplicate `ssm.logger.Debug(...)` for the same event.
- Add: `dbg.Log(dbg.Event, "StageStateMachine.onBossDefeated fired (state=%v)", ssm.state)`.

**Step 3: State transitions**

- Where state is updated (e.g. `ssm.setState(newState)` or assignment to `ssm.state`), add before or after: `dbg.Log(dbg.State, "StageStateMachine: %v → %v", oldState, newState)`. Ensure every transition path logs once (old → new).

**Step 4: Build**

Run: `go build ./...`

**Step 5: Commit**

```bash
git add internal/ecs/systems/stage/stage_state_machine.go
git commit -m "stage: add dbg Event and State breadcrumbs in StageStateMachine"
```

---

## Task 6: Breadcrumb — Overlay (drawWaveDebugInfo)

**Files:**
- Modify: `internal/game/debug_wave.go`

**Step 1: Add overlay breadcrumb**

- Add import: `"github.com/jonesrussell/gimbal/internal/dbg"`.
- In `drawWaveDebugInfo`, you need boss existence. Check how to get "boss exists" (e.g. query world for boss entity or ask stage state). Add one line: `dbg.Log(dbg.State, "Overlay sees state=%v bossExists=%v", st, bossExists)` where `bossExists` is a bool (e.g. from a helper that checks for a boss entity in the world). If no simple helper exists, use `dbg.Log(dbg.State, "Overlay sees state=%v", st)` for now.
- Remove or replace the `g.logger.Debug("UI reading stage state", ...)` call with the dbg line so we don’t duplicate.

**Step 2: Build**

Run: `go build ./...`

**Step 3: Commit**

```bash
git add internal/game/debug_wave.go
git commit -m "game: add dbg Overlay breadcrumb in drawWaveDebugInfo"
```

---

## Task 7: One-frame trace mode (dbg.Trace)

**Files:**
- Modify: `internal/dbg/log.go` (add Trace, traceNextFrame)
- Modify: `internal/game/game.go` (Update: at start if trace set then Enable; at end if trace set then Disable and clear)

**Step 1: dbg package**

In `internal/dbg/log.go`:

- Add: `var traceNextFrame uint32` (atomic).
- Add: `func Trace() { atomic.StoreUint32(&traceNextFrame, 1) }`.
- Add helper: `func traceRequested() bool { return atomic.LoadUint32(&traceNextFrame) == 1 }` and `func clearTrace() { atomic.StoreUint32(&traceNextFrame, 0) }`.

No change to `Log`; trace only toggles global Enable/Disable for one frame. Export: `func TraceRequested() bool { return atomic.LoadUint32(&traceNextFrame) == 1 }` and `func ClearTrace() { atomic.StoreUint32(&traceNextFrame, 0) }`.

**Step 2: Game Update**

In `internal/game/game.go`:

- Add import: `"github.com/jonesrussell/gimbal/internal/dbg"`.
- At the very start of `Update()`: if `dbg.TraceRequested()`, call `dbg.Enable()`.
- At the very end of `Update()` (before `return nil`): if `dbg.TraceRequested()`, call `dbg.Disable()` and `dbg.ClearTrace()`.

**Step 3: Trigger (debug key)**

In `internal/game/game_performance.go` (or game_input.go), in the debug key handling (e.g. F3 for overlay), add another key (e.g. F4) that calls `dbg.Trace()`. So when the user presses F4, the next frame is fully traced then dbg is disabled again.

**Step 4: Build and test**

Run: `go build ./...`. Run game, enter play, press F4; next frame should produce a burst of dbg lines, then silence.

**Step 5: Commit**

```bash
git add internal/dbg/log.go internal/game/game.go internal/game/game_performance.go
git commit -m "feat: add dbg one-frame trace mode and F4 trigger"
```

---

## Task 8: Replace Zap/fmt with dbg in game loop (game package)

**Files:**
- Modify: `internal/game/init_systems.go` (remove remaining g.logger.Debug or replace with dbg where it’s loop-relevant; init is one-off so logger.Debug here could become log.Printf or be removed)
- Modify: `internal/game/game_systems.go` (g.logger in updateSystemWithTiming: keep Error/Warn for failures and slow system — those are not “debug trace”, they are operational; design says “inside the loop it’s dbg only” for debug output; so replace only Debug calls, keep Error/Warn as-is until we remove Zap in Task 9 — then those become log.Printf)
- Modify: `internal/game/debug_wave.go` (already done in Task 6)
- Modify: `internal/game/game_performance.go` (g.logger.Debug for overlay toggled → dbg.Log(dbg.System, "Debug overlay toggled") or remove)
- Modify: `internal/game/game_level.go`, `game_init.go`, `init_entities.go`, `game_events.go`, `game_state.go`, `game_input.go`, `game.go` (Cleanup) — replace or remove every g.logger.Debug in these; for non-loop code (init, cleanup, one-off events) we can use log.Printf or remove; for loop code use dbg.

**Step 1: init_systems.go**

Remove all remaining `g.logger.Debug` calls (or replace with dbg if they run during init and you want them in the same stream; init runs once so dbg.Log is fine). Leave g.logger.Warn for “failed to load stage” etc. for now (migrated in Task 9).

**Step 2: game_systems.go**

No g.logger.Debug in this file currently; g.logger is used for Error/Warn in updateSystemWithTiming. Leave those for Task 9.

**Step 3: game_performance.go**

Replace `g.logger.Debug("Debug overlay toggled", ...)` with `dbg.Log(dbg.System, "Debug overlay toggled (enabled=%v)", g.showDebugInfo)` or remove. Add dbg import.

**Step 4: Other game/*.go**

- game_level.go: g.logger.Debug → dbg.Log(dbg.State, ...) or remove; g.logger.Error → leave for Task 9.
- game_init.go: logger.Debug → remove or log.Printf (app-level).
- init_entities.go: g.logger.Debug → remove or dbg (one-off at init, dbg is fine).
- game_events.go: g.logger.Debug in event handlers → dbg.Log(dbg.Event, ...) or remove.
- game_state.go: gsm.logger.Debug → dbg.Log(dbg.State, ...).
- game_input.go: g.logger.Debug → dbg.Log(dbg.System, ...) or remove.
- game.go Cleanup: g.logger.Debug/Error → leave for Task 9 (app-level).

**Step 5: Build**

Run: `go build ./...`

**Step 6: Commit**

```bash
git add internal/game/*.go
git commit -m "game: replace logger.Debug with dbg in game package loop paths"
```

---

## Task 9: Replace Zap with dbg or log in ECS systems (loop only)

**Files:**
- Modify all files under `internal/ecs/systems/` and `internal/ecs/managers/` that use logger.Debug in code paths that run every frame or from the game loop. Replace with dbg.Log(cat, msg, args) or remove. Do not change Error/Warn yet (Task 10).
- Key files: gyruss (already done in Task 4), stage (Task 5), collision, health, movement, weapon, behavior, path, powerup, enemy (gyruss_spawner, gyruss_wave_manager), resource (sprite, audio — these are init/cache, not loop; can leave or switch to log in Task 10).

**Step 1: List and replace**

For each system that runs in the loop (health, movement, collision, weapon, behavior/*, path, powerup, enemy), replace every logger.Debug with dbg.Log(appropriate category) or remove. Add `"github.com/jonesrussell/gimbal/internal/dbg"` where needed. Keep logger for Error/Warn for now.

**Step 2: Build and test**

Run: `go build ./...` and `go test ./...` (or at least `go test ./internal/ecs/... ./internal/game/...`).

**Step 3: Commit**

```bash
git add internal/ecs/systems/...
git commit -m "ecs: replace logger.Debug with dbg in loop systems"
```

---

## Task 10: Remove Zap from project; use log for app-level

**Files:**
- Modify: `internal/app/container.go`, `main.go`, and any remaining code that uses the injected logger (scenes/manager.go, etc.). Replace logger.Info/Warn/Error with `log.Printf` with a prefix or level in the message (e.g. `log.Printf("[INFO] ...")`).
- Remove: logger from `common.Logger` interface and from all constructors that take a logger; remove Zap dependency from go.mod and imports.
- Modify: `internal/game/game.go` and game constructors to not take or store `logger`; remove from ECSGame struct.
- Modify: all systems that currently receive a logger to no longer take it (only use dbg in loop); for non-loop errors they can use log.Printf or return errors to the caller.

**Step 1: App and main**

- In container.go, replace c.logger.Info/Warn/Error/Debug with log.Printf. Remove logger creation and injection.
- In main.go, replace logger.Warn/Info with log.Printf. Remove logger init.

**Step 2: Game and scenes**

- Remove logger field from ECSGame; remove from NewECSGame and all places that pass g.logger. Replace any remaining g.logger in game package with log.Printf (for errors) or dbg (already done for debug).
- In scenes/manager.go, remove sceneMgr.logger and use log.Printf for errors/info.

**Step 3: ECS systems and managers**

- Remove logger parameter from all system constructors and structs. Use dbg in the loop; use log.Printf for rare errors (e.g. “failed to emit event”) or return errors.
- Update resource manager, stage loader, and any other manager that used logger to use log.Printf for errors.

**Step 4: common.Logger**

- Remove or simplify: if nothing implements Logger anymore, remove the interface and the dependency from common package.

**Step 5: go.mod**

- Run: `go mod tidy` and remove the zap (and any zap-dependent logger package) dependency.

**Step 6: Build and test**

Run: `go build ./...` and `go test ./...`. Fix any remaining references to logger.

**Step 7: Commit**

```bash
git add -A
git commit -m "chore: remove Zap; use log for app-level and dbg in loop"
```

---

## Optional (Phase 2): In-game log viewer

**Scope:** Ring buffer of last 200 dbg lines; key `L` toggles scrollable overlay; optional category filter. Implement a sink in dbg (e.g. optional callback or buffer) that the game registers; dbg.Log writes to it when set. Omit from initial implementation plan; add as a follow-up plan or task list when needed.

---

## Execution

Plan complete and saved to `docs/plans/2026-02-20-dbg-logging-implementation.md`.

**Two execution options:**

1. **Subagent-driven (this session)** — I dispatch a fresh subagent per task, review between tasks, fast iteration.
2. **Parallel session (separate)** — Open a new session with executing-plans and run through the plan task-by-task with checkpoints.

Which approach do you want?
