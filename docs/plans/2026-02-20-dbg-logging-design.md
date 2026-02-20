# Game-Loop Debug Logging Layer (dbg) — Design

**Status:** Approved  
**Date:** 2026-02-20  
**Scope:** Purpose-built, at-a-glance debug logging for the real-time game loop; replace Zap/fmt in the loop with a minimal `dbg` package. Application-level output (startup, shutdown, config, errors) uses `log.Printf` or nothing.

---

## 1. Goal and rationale

Game-loop logging today is a mix of Zap `Debug` and ad-hoc `fmt.Printf`/`fmt.Println`. At 60 FPS this produces unstructured, high-noise output that is hard to scan. The goal is a **single, human-readable, category-based debug layer** used only inside the game loop: events, state transitions, system order, spawns, and world identity — with frame markers, a global on/off, and optional one-frame trace and in-game viewer.

**Principle:** Inside the loop it’s dbg only. No Zap, no fmt spam, no JSON.

---

## 2. Package and API

**Location:** `internal/dbg`

**Categories (v1):** `Event`, `State`, `System`, `World`, `Spawn`. Add `Input` only if input breadcrumbs are needed later.

**Core API:**

- `dbg.Log(cat Category, msg string, args ...any)` — formats with `fmt.Sprintf(msg, args...)`, prints `[CAT] message`. No timestamps, no JSON, no structured fields in v1.
- `dbg.Enable()` / `dbg.Disable()` — global toggle implemented with `sync/atomic` so it is safe from the game thread and any debug UI.

**No per-category toggles in v1** — keeps surface area minimal; category filtering can be added later (e.g. with the in-game viewer).

**Dependencies:** Stdlib only (`fmt`, `sync/atomic`). No Zap or other log package in this package.

---

## 3. Integration points and frame marker

**Frame marker:** In the main gameplay update path, once every N frames (e.g. 30), call:

```go
dbg.Log(dbg.System, "=== FRAME %d ===", g.frame)
```

**Breadcrumbs at key points:**

- **Event emission** — e.g. `dbg.Log(dbg.Event, "EmitBossDefeated (world=%p eventSystem=%p)", world, eventSystem)`
- **Event reception** — e.g. `dbg.Log(dbg.Event, "StageStateMachine.onBossDefeated fired (state=%v)", ssm.state)`
- **State transitions** — e.g. `dbg.Log(dbg.State, "%v → %v", old, new)` (with optional prefix like `StageStateMachine:` if needed)
- **System update boundaries** — e.g. `dbg.Log(dbg.System, "updateGameplaySystems start")` at start of gameplay update
- **Overlay / debug view** — e.g. `dbg.Log(dbg.State, "Overlay sees state=%v bossExists=%v", st, boss != nil)` when drawing overlay
- **Spawn events** — e.g. `dbg.Log(dbg.Spawn, "boss spawned")` where useful
- **World / event-system identity** — e.g. one line at init or first use when the world or event system is set/used

All calls are one-line; no structured fields. More breadcrumbs can be added incrementally.

---

## 4. One-frame trace mode

**Mechanism:** Package-level “trace next frame” flag (e.g. `traceNextFrame` or atomic), set by `dbg.Trace()` (or `dbg.TraceNextFrame()`).

**Behavior:**

- At **start** of the game’s `Update()` (or start of the gameplay-update slice): if trace flag is set, call `dbg.Enable()`.
- Run the rest of the update as usual.
- At **end** of that same update: if trace was set, call `dbg.Disable()` and clear the flag so the next frame is silent again.

**Trigger:** Debug key or small debug console command that calls `dbg.Trace()` so the user gets one frame of full trace on demand.

---

## 5. Removing Zap from the loop

**Rule:** All debug output inside the game loop and from systems invoked every frame uses **only** `dbg`. No Zap, no `fmt.Printf`/`fmt.Println` in the loop.

**Concrete steps:**

- Replace or remove every `g.logger.Debug` (and any `logger.Debug`) that runs inside the loop or from systems called from the loop (Gyruss, stage state machine, collision, health, weapon, movement, behavior, path, power-up, spawn, etc.) with either a `dbg.Log(dbg.<Category>, ...)` call at the same semantic point or nothing if the log is low-value.
- Remove ad-hoc `fmt.Printf`/`fmt.Println` in init_systems, game_systems, gyruss_system, stage_state_machine and replace with one-line `dbg.Log(...)` where a breadcrumb is desired (e.g. event system identity as a single `dbg.Log(dbg.World, ...)` at init or first use).

**Outside the loop:** Use `log.Printf` for startup, shutdown, config, and non-frame errors, or no logging where appropriate. Zap is removed from the project (Approach A).

---

## 6. Optional in-game log viewer (phase 2)

**Scope:** Ring buffer of the last N lines (e.g. 200) of dbg output; one key (e.g. `L`) toggles a scrollable overlay; optional filter by category.

**Implementation note:** `dbg.Log` optionally writes to a sink (e.g. callback or interface registered by the game) so the same call both prints and appends to the ring buffer. No file I/O or second format. Exact keybind and layout are left to the implementation plan.

---

## 7. Summary

- **internal/dbg:** Categories, `Log`, `Enable`/`Disable`, atomic, stdlib-only.
- **Loop:** Frame marker every N frames; breadcrumbs at events, state, systems, overlay, spawn, world.
- **Trace:** One-frame trace via `dbg.Trace()` and enable/disable around a single update.
- **Zap:** Removed from loop; removed from project; app-level output uses `log.Printf` or nothing.
- **Optional:** In-game viewer (ring buffer + key + category filter) as phase 2.
