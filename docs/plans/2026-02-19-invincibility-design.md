# Invincibility: Developer God Mode + Player Power-Up — Design

**Status:** Approved  
**Date:** 2026-02-19  
**Scope:** Two distinct systems — dev invincibility (tool) and player invincibility (power-up).

---

## 1. Context and goal

Gimbal already has (a) debug-only CLI invincibility (`-invincible` when `DEBUG=true`) and (b) post-hit i-frames on the player. This design adds a **pause-menu God mode toggle** for developers and a **player invincibility power-up** as a third power-up type alongside double_shot and extra_life.

**Goals**

- **Developer invincibility:** Session-only tool for iteration; activation via CLI and/or pause-menu toggle when Debug is on; absolute immunity, not saved.
- **Player invincibility:** Gameplay reward via existing power-up pipeline; temporary full invulnerability using the same `Health` invincibility state as i-frames; no menu toggle for players.

---

## 2. Priority order and damage flow

**Single decision point:** In `HealthSystem.DamageEntity` (and any other place that applies damage to the player), apply this order:

1. **Developer invincibility** — If active (CLI or pause toggle when Debug), return immediately; no damage, no i-frames.
2. **Player invincibility (power-up or i-frames)** — If `Health.IsInvincible` is true (set by power-up or by previous hit), return immediately.
3. **Apply damage** — Subtract health, set `IsInvincible = true` and `InvincibilityTime = InvincibilityDuration`, handle death/respawn as today.

Power-up and i-frames both use the same `Health` fields; no separate “power-up invincibility” flag. The power-up simply sets `IsInvincible` and `InvincibilityTime` (e.g. 5s) on collect; existing `HealthSystem.Update` countdown and renderer flash behavior stay unchanged.

---

## 3. Developer invincibility (God mode)

**Activation**

- **CLI (unchanged):** `-invincible` when `DEBUG=true`; ignored when Debug is false.
- **Pause menu (new):** When `config.Debug` is true, show a fourth pause option: “God mode: ON” or “God mode: OFF” that toggles dev invincibility for the session.

**Implementation**

- **Config:** Keep `config.Invincible` as the single source of truth for “dev invincible on”. Add `SetDevInvincible(bool)` on `GameConfig` that sets `Invincible = value` only when `Debug` is true (no-op when Debug is false).
- **Pause scene:** If `manager.GetConfig().Debug`, add a menu option that displays current state (“God mode: ON” / “God mode: OFF”) and on select calls `config.SetDevInvincible(!config.Invincible)`. SceneManager and PausedScene already hold the same `*config.GameConfig` reference as the game.
- **Health system:** No logic change; continue to check `config.Debug && config.Invincible` first in `DamageEntity` and return nil when true (player only).

**Persistence:** Session-only; not saved. Resetting or restarting the game clears it (CLI can set it again).

---

## 4. Player invincibility power-up

**Integration:** Third power-up type in the existing pipeline: same spawn/collect flow as double_shot and extra_life; uses `PowerUpType` and `PowerUpTypeData`; collected effect applied in `PowerUpSystem.collectPowerUp`.

**Core types**

- **`core.PowerUpType`:** Add `PowerUpInvincibility` (after `PowerUpExtraLife`).
- **`PowerUpTypeData`:** Already has `Duration`; use it for invincibility duration (e.g. 5s default).

**Collection behavior**

- In `collectPowerUp`, when `data.Type == core.PowerUpInvincibility`: find player entry (same as extra_life), get `Health`, set `IsInvincible = true` and `InvincibilityTime = max(health.InvincibilityTime, data.Duration)`. If `data.Duration == 0`, use a default (e.g. 5s). No separate state in PowerUpSystem for invincibility; the Health component is the source of truth and is already ticked down by HealthSystem.Update.

**Spawning**

- **Sprite:** Add a distinct sprite in `createSprites()` (e.g. white or cyan) for `PowerUpInvincibility`.
- **SpawnPowerUp:** When `powerUpType == core.PowerUpInvincibility`, set `PowerUpTypeData.Duration` to desired duration (e.g. 5s).
- **TrySpawnPowerUp:** Extend weighted random: e.g. 60% double_shot, 25% extra_life, 15% invincibility (tunable constants).

**Visual/audio:** Reuse existing invincibility flashing (render system already keys off `Health.IsInvincible`). Optional: log or future sound on collect (“Invincibility!”); no requirement for this design.

---

## 5. Edge cases and testing

- **Dev + power-up:** Dev invincibility is checked first; if on, power-up state is irrelevant for damage. No conflict.
- **Power-up then hit:** While power-up is active, `IsInvincible` is true so damage is no-op; timer is not reset by hits. When timer expires, HealthSystem clears `IsInvincible` as today.
- **Stacking:** Collecting a second invincibility power-up: use `max(InvincibilityTime, data.Duration)` so a new pickup extends or grants invincibility without shortening.
- **Testing:** Unit tests for HealthSystem damage order (dev invincible → no damage; power-up/i-frames → no damage; otherwise damage). Power-up: test collect sets Health invincibility and duration; optional integration test that spawns invincibility power-up and collects it.

---

## 6. Files to touch (summary)

- **Config:** `internal/config/config.go` — add `SetDevInvincible(bool)`.
- **Health:** `internal/ecs/systems/health/system.go` — keep current order; ensure first check is dev invincible (already is).
- **Core:** `internal/ecs/core/components.go` — add `PowerUpInvincibility` to `PowerUpType` const.
- **Power-up:** `internal/ecs/systems/powerup/powerup_system.go` — handle new type in `createSprites`, `collectPowerUp`, `SpawnPowerUp`, `TrySpawnPowerUp`.
- **Pause:** `internal/scenes/pause/paused.go` — when `config.Debug`, add “God mode” option and call `config.SetDevInvincible`.

No change to stage JSON or entity config for this design; power-up drop weights are in code. Config-driven weights can be a later follow-up.
