# Architecture validation: Boss and wave lifecycles

This document answers explicit questions about the current design and recommends adjustments before implementing the boss-state and wave-timing fixes.

---

## 1. Boss lifecycle architecture

### 1.1 Does the system currently model the boss lifecycle as explicit states, or is it inferred from booleans and entity presence?

**Answer: Inferred from booleans and entity presence.** There is no explicit state enum.

- **Wave manager:** `bossTriggered` (bool) — set when the final wave’s `on_clear` is "boss".
- **Gyruss system:** `bossSpawned` (bool), `bossTimer` (float64). Boss is spawned after `SpawnDelay` elapses; then `bossSpawned = true`.
- **“Defeated”:** In `IsBossDefeated()` the system only checks world for a boss entity; it returns `!bossExists`. So “defeated” is inferred as “no boss entity in world,” which is also true before the boss has ever spawned.

Effective lifecycle is therefore: **NotTriggered** (implicit) → **Triggered** (`bossTriggered`) → **Spawning** (timer) → **Spawned** (`bossSpawned` + entity exists) → **Defeated** (no entity). The “Spawning” and “Defeated” vs “not yet spawned” distinction is not explicit and causes the debug “Boss: Defeated” bug.

### 1.2 Should boss lifecycle be represented as a single authoritative enum/state machine instead of multiple booleans?

**Recommendation: Yes, for clarity and to avoid bugs.**

- Introduce a single **BossLifecycleState** (e.g. `NotTriggered | Triggered | Spawning | Active | Defeated`) owned by one place (e.g. GyrussSystem or a dedicated component/manager).
- Transitions: Triggered when final wave clears; Spawning when `bossTriggered` and timer &lt; SpawnDelay; Active when boss entity exists after spawn; Defeated when `bossSpawned` and entity removed.
- Replace `bossTriggered`, `bossSpawned`, and “no entity = defeated” with this state so UI and other systems query one source of truth.

**Short-term (minimal fix):** Keep booleans but fix semantics: `IsBossDefeated()` = `bossSpawned && !bossExists`. No new state machine yet.

**Medium-term:** Add an explicit enum and transition logic; deprecate the boolean-based API.

### 1.3 Are there any systems that implicitly assume “no boss entity = defeated” and would break once semantics change?

**Answer: Yes — but the planned fix is compatible.**

- **`internal/game/debug_wave.go`:** Uses `IsBossDefeated()` to choose “Boss: Defeated” vs “Boss: Spawning soon…”. It currently shows Defeated whenever there is no boss entity (including pre-spawn). Changing `IsBossDefeated()` to `bossSpawned && !bossExists` fixes this; no other change in debug_wave needed.
- **`internal/game/debug_boss.go`:** When there is no boss entry it uses `WasBossSpawned()` for status: Defeated vs Spawning soon. So it already distinguishes “never spawned” vs “was spawned”; it does not rely on “no entity = defeated” in a way that breaks.
- **`internal/game/game_level.go`:** Uses `IsStageComplete()` = `bossSpawned && IsBossDefeated()`. With `IsBossDefeated()` = `bossSpawned && !bossExists`, stage complete remains “boss was spawned and is now gone” — correct.
- **`internal/scenes/gameplay/playing.go`:** Does **not** use GyrussSystem for boss presence. It uses `isBossActive()` which **queries the world** for a boss entity. So music/intro react to entity presence only; they do not depend on `IsBossDefeated()` and will not break.

So: the only place that “assumes no boss entity = defeated” in a harmful way is the debug overlay branch that uses `IsBossDefeated()`. Tightening `IsBossDefeated()` to require `bossSpawned` is safe and fixes the bug.

### 1.4 Should boss telegraphs, spawn delays, multi-phase bosses, or cutscenes be first-class lifecycle states?

**Recommendation: Not yet; design so they can be added later.**

- **Spawn delay:** Already first-class (timer in GyrussSystem + config `SpawnDelay`). An explicit state machine would have a `Spawning` state that lasts until delay elapses.
- **Telegraphs / cutscenes:** Could be extra states (e.g. `Telegraphing`, `IntroCutscene`) before `Active`. Not required for the current fix.
- **Multi-phase bosses:** Would be phases *within* Active (or sub-states). The single “boss entity exists” check can stay; phase could live on the boss entity or a separate state. A unified StageStateMachine (see section 5) could own “phase” when you add it.

For now: fix semantics with minimal change; document that an explicit BossLifecycleState enum is the desired direction for future telegraphs/phases.

---

## 2. Wave lifecycle architecture

### 2.1 Should wave progression be strictly world-driven (enemy count, spawn completion), or should timers ever be allowed to advance waves?

**Recommendation: Strictly world-driven for normal progression.**

- **Current:** Wave completes when (a) `allSpawned && activeEnemies == 0` or (b) `wave.Timing.Timeout` elapses. So the next wave (or boss) can start while the previous wave’s enemies are still on screen if the timeout fires.
- **Desired:** Next wave starts only when the current wave is clear: `allSpawned && activeEnemies == 0`. No timeout-based advance.
- **Optional:** A separate, much larger “safety” timeout (e.g. 120s) could force-advance to avoid softlock (e.g. one enemy stuck and never removed). Not required for the first pass.

So: remove the timeout-based completion branch in `checkWaveCompletion()`. Keep `timing.timeout` in JSON for possible future use (e.g. safety or UI hints); it will no longer drive completion.

### 2.2 Does the current architecture clearly separate wave spawning, wave completion detection, stage progression, and boss triggering?

**Answer: Partially.**

- **Wave spawning:** Handled in GyrussWaveManager (ShouldSpawnEnemy, MarkEnemySpawned) and GyrussSystem (handleSpawning). Clear.
- **Wave completion detection:** Inside GyrussWaveManager.Update() → checkWaveCompletion(). Completion is “inferred” each frame (timeout or clear), not an explicit transition. Mixed with spawn timing (waveTimer, groupSpawnTimer, interWaveTimer).
- **Stage progression:** In game/game_level.go: each frame `checkLevelCompletion()` calls `IsStageComplete()` and then handleLevelComplete() (next stage, scene transition). So “stage complete” is separate from wave logic, but it depends on GyrussSystem’s boss state.
- **Boss triggering:** Wave manager sets `bossTriggered` in completeWave() when wave’s on_clear is "boss". GyrussSystem handles actual spawn (handleBossSpawning) and owns `bossSpawned`. So “trigger” is in wave manager; “spawn” and “defeated” are in Gyruss system.

Separation is reasonable but could be clearer: wave completion is buried inside Update and could be an explicit “wave completed” transition or event.

### 2.3 Should wave completion be an explicit state transition rather than inferred inside checkWaveCompletion()?

**Recommendation: Yes, as a medium-term improvement.**

- **Now:** checkWaveCompletion() runs every frame and calls completeWave() when conditions hold. No explicit “wave state” (e.g. Spawning | Clearing | Completed).
- **Better:** Treat “wave completed” as a discrete transition: when `allSpawned && activeEnemies == 0` (and optionally no timeout), set a wave state to Completed (or emit WaveCompleted), then advance to next wave or boss in one place. That makes it obvious when and why progression happens and simplifies adding events or UI.

For the immediate fix, removing the timeout path is enough; explicit wave state can be a follow-up refactor.

### 2.4 Are there any systems that rely on timeout-based wave completion (e.g. music, UI, pacing)?

**Answer: No.**

- Music and boss intro are driven by boss **entity** presence (PlayingScene.isBossActive()) and level/stage changes, not by wave completion or timeout.
- Debug overlay shows wave index and boss status from wave manager and Gyruss system; it does not depend on timeout to show “wave complete.”
- Level completion uses IsStageComplete() (boss spawned and defeated), not wave timeout.

So removing timeout-based wave completion does not break music, UI, or pacing. The integration test advances by destroying enemies (screen clear), not by timeout, so it still passes.

---

## 3. ECS responsibility boundaries

### 3.1 Is boss state owned by the GyrussSystem, or should it be a dedicated BossLifecycleSystem?

**Current:** GyrussSystem owns `bossSpawned` and `bossTimer`; it performs handleBossSpawning and exposes IsBossActive(), IsBossDefeated(), WasBossSpawned(), IsStageComplete(). Wave manager owns `bossTriggered`.

**Recommendation:** Keep boss state in GyrussSystem for now.

- Boss spawning is tightly coupled to wave manager (triggered when final wave clears) and spawner (SpawnBoss). A separate BossLifecycleSystem would need the same dependencies and would duplicate coordination that GyrussSystem already does.
- A dedicated system becomes useful when you add multiple boss phases, telegraphs, or cutscenes; then a small BossLifecycleSystem that only advances state and emits events could sit alongside GyrussSystem.
- For the planned fix, keeping ownership in GyrussSystem and fixing IsBossDefeated() is sufficient.

### 3.2 Should wave management be split into WaveSpawnerSystem, WaveCompletionSystem, StageProgressionSystem?

**Recommendation: Not for the current scope.**

- GyrussWaveManager already encapsulates wave index, spawn timing, and completion; GyrussSystem drives it and handles boss spawn. Splitting into three systems would spread state (wave index, timers, “completed” flags) and add coordination overhead.
- Clearer separation can be achieved with **events** (WaveCompleted, BossTriggered, etc.) and a single “stage progression” owner that reacts to those events, without necessarily splitting into three systems. If the codebase grows (e.g. mid-wave reinforcements, optional waves), then a WaveCompletionSystem that only evaluates “is current wave clear?” and emits WaveCompleted could make sense.

So: keep one wave manager and GyrussSystem; improve semantics and optionally add events; consider splitting only if requirements grow.

### 3.3 Are there any cross-system couplings that should subscribe to events instead of polling?

**Answer: Yes.**

- **Debug overlay (game/debug_wave.go, debug_boss.go):** Polls GyrussSystem and wave manager every frame to choose text. Could subscribe to BossSpawned / BossDefeated / WaveStarted / WaveCompleted and update a cached string, or keep polling with corrected semantics (simplest for now).
- **PlayingScene (boss music, boss intro):** Polls the world every frame via isBossActive() and compares to previous frame to detect “boss just appeared” / “boss just died.” This is effectively a state-change detector. It would be cleaner to subscribe to BossSpawned and BossDefeated so the scene reacts to events instead of diffing entity presence.
- **game_level.checkLevelCompletion():** Polls IsStageComplete() every frame. Could subscribe to StageCompleted (or BossDefeated when it implies stage complete) and then transition scene once.

Recommendation: Introduce progression events (see section 4) and migrate these consumers to subscribe where it simplifies logic (especially PlayingScene and level completion). Debug overlay can stay polling until events are in place.

---

## 4. Event model

### 4.1 Should the game emit explicit events: WaveStarted, WaveCompleted, BossSpawned, BossDefeated, StageCompleted?

**Recommendation: Yes.**

- The codebase already has an EventSystem (internal/ecs/events/) with EnemyDestroyed, LevelChanged, GameOver, etc. No WaveStarted, WaveCompleted, BossSpawned, BossDefeated, or StageCompleted today.
- Adding these would:
  - Give a single, ordered record of progression (useful for debugging and analytics).
  - Let PlayingScene, HUD, and level completion react to transitions instead of polling.
  - Make it easier to add features (e.g. “on WaveCompleted play sound,” “on BossSpawned show telegraph”) without touching GyrussSystem internals.

Proposed events:

- **WaveStarted** (waveIndex, waveID or description) — when a wave begins (after level start delay or inter-wave delay).
- **WaveCompleted** (waveIndex, waveID, onClear: "next_wave" | "boss") — when checkWaveCompletion() calls completeWave().
- **BossTriggered** — when final wave clears with on_clear "boss" (optional; could be folded into WaveCompleted).
- **BossSpawned** — when handleBossSpawning actually spawns the boss entity.
- **BossDefeated** — when the boss entity is removed (e.g. in collision when DestroyEnemy is called for the boss).
- **StageCompleted** — when IsStageComplete() becomes true (boss spawned and defeated).

Emit from: GyrussSystem and/or wave manager (and game_level for StageCompleted if desired). Collision system already emits EnemyDestroyed; it could detect “was boss” and also emit BossDefeated, or GyrussSystem could detect “boss was present, now not” and emit (latter keeps boss semantics in one place).

### 4.2 Would event-driven progression reduce the need for polling and boolean juggling?

**Yes.**

- PlayingScene would subscribe to BossSpawned / BossDefeated for music and intro instead of isBossActive() and bossWasActive.
- Level completion could subscribe to StageCompleted and call handleLevelComplete() once.
- Debug overlay could either subscribe and cache last event state or keep polling a single “stage progression state” that is updated from events.

Boss state could remain as booleans/state inside GyrussSystem for gameplay logic, but transitions would be published as events so UI and scene logic don’t need to infer from entity presence or multiple flags.

---

## 5. Future-proofing

### 5.1 If we later add multi-phase bosses, mid-wave reinforcements, optional waves, or difficulty modifiers, will the current lifecycle model scale?

**Assessment:**

- **Multi-phase bosses:** Current model has one “boss entity exists” and one “defeated” notion. Multi-phase would need “phase” (e.g. 1, 2, 3) and possibly “phase transition” (telegraph, invuln). The current single boolean/entity check would need to become a small state machine or phase counter. An explicit BossLifecycleState enum (with phases as sub-states or extra enum values) would scale.
- **Mid-wave reinforcements:** Would require “partial” wave completion or additional spawn triggers during a wave. Current wave completion is “all groups spawned + no enemies.” Reinforcements could be new groups that spawn on a timer or event; wave manager would need a way to add groups or trigger “spawn group X at T.” The current single “spawn sequence” per wave could be extended (e.g. delayed groups, or WaveCompleted used to spawn a “reinforcement” wave that doesn’t advance the wave index). Doable but not first-class today.
- **Optional waves:** Would need “skip wave” or conditional waves (e.g. if player has X, skip wave 2). Wave manager currently has a linear sequence. Would require conditional logic in completeWave() or in stage config (e.g. “next_wave_id” or conditions). Current model doesn’t support this; adding it would be a config + state change.
- **Difficulty modifiers:** Already present in stage config (DifficultySettings). They could affect spawn rate, health, speed. No change to lifecycle needed; only to how spawner/behavior use config.

So: boss lifecycle should be made explicit (state enum) for multi-phase; wave model can be extended for reinforcements and optional waves with config and a bit of state.

### 5.2 Should we formalize a unified “StageStateMachine” that owns boss/wave transitions?

**Recommendation: Yes, as a medium-term target.**

- A **StageStateMachine** (or “StageProgression” component/manager) could own:
  - **Wave state:** Which wave is active; whether we’re in “level start delay,” “spawning,” “inter-wave delay,” or “wave complete.”
  - **Boss state:** NotTriggered → Triggered → Spawning → Active → Defeated.
  - **Stage state:** InProgress | BossPhase | StageComplete.

- Transitions would be driven by:
  - World (enemy count, boss entity presence).
  - Timers (level start, inter-wave, boss spawn delay).
  - One place would evaluate “wave clear” and “boss defeated” and advance state, and optionally emit events.

- Benefits: single place to reason about progression, easier to add phases/optional waves/reinforcements, and UI/scenes can subscribe to one state or to events emitted from the state machine.

For the immediate boss and wave fixes, you don’t need to implement this. Fix IsBossDefeated() and remove timeout-based wave completion; then document “StageStateMachine” as the desired direction and introduce events (and optionally an explicit boss state enum) as a next step.

---

## Summary of recommendations

| Area | Immediate (for current fix) | Follow-up |
|------|-----------------------------|-----------|
| **Boss lifecycle** | Fix `IsBossDefeated()` = `bossSpawned && !bossExists`. No new state machine. | Introduce BossLifecycleState enum; optionally BossLifecycleSystem if phases/telegraphs are added. |
| **Wave completion** | Remove timeout-based completion in `checkWaveCompletion()`; progression only when `allSpawned && activeEnemies == 0`. | Explicit wave state (e.g. Completed) and/or WaveCompleted event. |
| **Ownership** | Keep boss and wave state in GyrussSystem and GyrussWaveManager. | Consider StageStateMachine that owns wave + boss transitions. |
| **Events** | Not required for the bug fix. | Add WaveStarted, WaveCompleted, BossSpawned, BossDefeated, StageCompleted; migrate PlayingScene and level completion to subscribe. |
| **Debug / UI** | Rely on corrected semantics; keep polling. | Switch to event-driven or single “progression state” query. |

Implementing the two concrete fixes (boss semantics + wave completion rule) is consistent with this architecture and does not block later refactors to state machines or events.
