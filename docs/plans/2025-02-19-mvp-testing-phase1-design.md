# MVP Testing Phase 1 — Design

**Status:** Approved  
**Date:** 2025-02-19  
**Scope:** CI + correctness baseline (Phase 1 only)

---

## 1. Context and goal

Gimbal is a Gyruss-style arcade game (Go, Ebiten, Donburi). This design establishes a minimal, industry-aligned testing foundation: **CI plus correctness baseline** first, with integration and performance work in later phases.

**Goal:** Every push/PR is validated by CI; wave → boss → stage complete is covered by unit and integration tests; all stage JSON is validated so bad data cannot ship.

---

## 2. Scope and architecture (Phase 1)

**In scope**

- **CI:** One workflow that runs on push/PR: `go build`, lint, `go test ./...` (including integration via `-tags=integration`). All steps blocking.
- **Unit tests:** Fill gaps in wave manager (wave start, spawn sequence, wave complete, edge cases), spawner (spawn index, orbit angle, group sequencing), and collision outcome (enemy death, entity removal). No new tests for already-covered math/utils.
- **Asset/config validation:** A test that loads all `assets/stages/*.json`, validates structure (waves, groups, enemy types, patterns, boss config), and fails on invalid or missing data.
- **One minimal integration test:** Run a slice of the game loop without `ebiten.RunGame`; assert wave 1 → wave 2 → boss spawn → boss defeated → stage complete. Use existing ECS/wave/spawn logic; simulate time and destroy enemies to drive progression.

**Out of scope for Phase 1**

- Scene transitions and scene-switching tests (Phase 2).
- Golden-path simulated player and performance baseline (Phase 3).
- New gameplay features.

**Architecture**

- CI: single config file (e.g. `.github/workflows/ci.yml`); four steps in order: build, lint, test (with integration tag).
- Tests remain in-tree (`*_test.go`). Integration test under `internal/ecs/systems/gyruss/` or `internal/game/`, using same ECS/managers as production, minimal mocks.
- Asset validation: Go test using existing stage loader; part of `go test ./...`.

---

## 3. CI and asset validation (concrete)

**CI**

- **Location:** `.github/workflows/ci.yml` (or equivalent).
- **Steps (order):** checkout → setup Go (match `go.mod`) → `go build ./...` → `task lint` → `go test ./... -race -tags=integration -timeout=60s` (or `task test:all` with integration enabled). Fail job on any failure.
- **Triggers:** Push and PRs to default branch.
- **Caching:** Optional (e.g. Go module cache); not required for Phase 1.

**Asset validation**

- **Location:** Test file next to stage loading (e.g. `internal/ecs/managers/stage_loader_test.go` or dedicated `*_test.go` in same package).
- **Behavior:** Load every stage from `assets/stages/` via existing StageLoader; assert: loadable JSON, non-empty waves, valid spawn groups (enemy type, count), known pattern names, boss config when present. Test fails on first validation error; CI fails accordingly.
- **No separate CI step:** Covered by `go test ./...`.

---

## 4. Unit-test coverage and integration harness

**Unit tests to add**

- **Wave manager:** Start/sequence (first wave after level-start delay), spawn flow (ShouldSpawnEnemy / MarkEnemySpawned), wave completion (all spawned + no enemies → next wave or boss), edge cases (empty group, inter-wave delay, OnClear behavior).
- **Spawner:** Spawn index progression, orbit angle distribution for a multi-enemy group.
- **Collision/game:** At least one test for “enemy destroyed” path (entity removed, no leak).

**Integration test**

- **Goal:** Prove stage 1 can complete (wave 1 → wave 2 → boss → stage complete) without a display.
- **Placement:** `internal/ecs/systems/gyruss/integration_test.go` (or `internal/game/integration_test.go`) with `//go:build integration`.
- **Setup:** Real world, StageLoader, GyrussSystem (same as existing Gyruss tests). Real stage 1 from assets.
- **Execution:** Load stage 1; loop calling `gs.Update(ctx, dt)`; to complete waves, let spawns happen then destroy non-boss enemies via `gs.DestroyEnemy`; when boss spawns, destroy boss; assert `IsStageComplete()`. Time cap (e.g. 10–15 s) to avoid hangs.
- **No Ebiten:** Do not call `ebiten.RunGame`; only ECS + GyrussSystem.Update. Mock or stub any dependency that would require a display.

---

## 5. Error handling and CI failure behavior

- **Test failure:** Any failing test exits non-zero; CI job fails; no retries.
- **Asset validation failure:** Same (it’s a normal test).
- **Integration:** Build tag `//go:build integration`; default `go test ./...` does not run it; CI runs with `-tags=integration`. Optionally skip in `-short` mode.
- **Build/lint:** Fail fast; order is build → lint → test so a broken build doesn’t run tests.
- **Strict:** All checks blocking; no flake retries in Phase 1.

---

## 6. Phase 1 acceptance criteria and out-of-scope

**Done when**

- CI runs on push/PR and blocks on build, lint, or test failure.
- Unit tests cover wave manager (start, spawn, completion, edge case), spawner (index/angle), and at least one collision-outcome test.
- Asset validation test loads and validates all stage JSONs.
- One integration test proves stage 1 completes (waves + boss) without display, with bounded runtime.
- Correctness baseline: wave 1 → wave 2 → boss → stage complete is validated; no softlock in that path.

**Explicitly out of scope**

- Scene transition tests; ECS cleanup tests; golden-path simulated player; performance baseline; Phases 2 and 3.

---

## 7. Next step

Invoke the **writing-plans** skill to produce the detailed implementation plan (bite-sized tasks, file paths, commands, commits) for Phase 1.
