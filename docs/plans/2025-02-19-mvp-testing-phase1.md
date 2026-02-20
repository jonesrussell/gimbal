# MVP Testing Phase 1 — Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Establish CI and a correctness baseline so every push/PR is validated and wave→boss→stage-complete is covered by unit and integration tests.

**Architecture:** Add a single CI workflow (build, lint, test with integration tag). Add asset validation test, fill unit-test gaps in wave manager and spawner, add one collision-outcome test, then add one headless integration test that runs stage 1 to completion. All tests in-tree; integration test behind `//go:build integration`.

**Tech Stack:** Go, Ebiten, Donburi, GitHub Actions, existing Taskfile (lint, tests). Reference design: `docs/plans/2025-02-19-mvp-testing-phase1-design.md`.

---

## Task 1: Add CI workflow

**Files:**
- Create: `.github/workflows/ci.yml`

**Step 1: Add workflow file**

Create `.github/workflows/ci.yml` that:
- Triggers on `push` and `pull_request` to default branch (no tag filter).
- Single job: checkout → setup Go (`go-version-file: go.mod`, cache `go.sum`) → build → lint → test.
- Build: `go build ./...`
- Lint: run `task lint` (ensure Task is available: `curl -sL https://taskfile.dev/install.sh | sh` and `~/go/bin` or install task; or run `go fmt ./...`, `go vet ./...`, `golangci-lint run ./...` if you prefer not to depend on Task in CI).
- Test: `go test ./... -v -race -tags=integration -timeout=60s`

**Step 2: Verify workflow syntax**

- Commit and push to a branch or run `act` locally if available; otherwise rely on PR check.
- Expected: workflow runs and all steps pass (tests may fail until later tasks are done).

**Step 3: Commit**

```bash
git add .github/workflows/ci.yml
git commit -m "ci: add workflow for build, lint, and test (with integration)"
```

---

## Task 2: Asset validation test

**Files:**
- Create or modify: `internal/ecs/managers/stage_loader_test.go` (or `internal/ecs/managers/stage_validation_test.go`)

**Step 1: Write the failing test**

Add a test that:
- Uses `StageLoader` with the real embedded assets FS (from `github.com/jonesrussell/gimbal/assets`).
- For each stage number from 1 to `loader.GetTotalStages()`, calls `LoadStage(n)`.
- Asserts: config non-nil, `len(config.Waves) > 0`, each wave has `len(wave.SpawnSequence) > 0`, each group has `EnemyType != ""` and `Count > 0`, and if `config.Boss.Enabled` then `config.Boss.BossType != ""` and `config.Boss.Health > 0`.
- Use a table-driven loop over stage numbers; fail with clear message (stage number and what failed).

**Step 2: Run test to verify it fails or passes**

- If LoadStage returns nil for missing files (or uses default): adjust test to load only stages that exist (e.g. read embed or use GetTotalStages). Ensure at least one stage file is required (e.g. stages 1–6 must load).
- Run: `go test ./internal/ecs/managers/... -v -run TestStageValidation` (or the test name you chose).
- Expected: PASS if current stages are valid; if not, fix data or test until it reflects “all stages must be valid.”

**Step 3: Commit**

```bash
git add internal/ecs/managers/stage_loader_test.go
git commit -m "test: add asset validation for all stage JSONs"
```

---

## Task 3: Wave manager — start next wave and spawn flow

**Files:**
- Modify: `internal/ecs/systems/enemy/gyruss_wave_manager_test.go`
- Reference: `internal/ecs/systems/enemy/gyruss_wave_manager.go` (startNextWave, ShouldSpawnEnemy, MarkEnemySpawned, Update)

**Step 1: Write failing test for wave start after level-start delay**

- In `gyruss_wave_manager_test.go`, add test: load stage with one wave and one group (e.g. count 2). Call `wm.Update(dt)` in a loop with small dt until `!wm.IsWaitingForLevelStart()` (or a max iteration cap). Assert `wm.GetCurrentWaveIndex() == 0` and that spawning has started (e.g. after more updates, `ShouldSpawnEnemy` returns true at least once).
- Run: `go test ./internal/ecs/systems/enemy/... -v -run TestWaveManager_LevelStartThenSpawning`
- Expected: FAIL if behavior doesn’t match (e.g. wrong timing); adjust test or code until intent is clear.

**Step 2: Write failing test for spawn flow and wave completion**

- Add test: load stage with one wave, one group (count 2). Advance past level start. Call `ShouldSpawnEnemy` and `MarkEnemySpawned` until the group is exhausted (or use Update and count spawns). Then remove all enemy entities from the world (query EnemyTag, world.Remove). Call `wm.Update(dt)` repeatedly until wave completes (e.g. boss triggered or currentWaveIndex advances). Assert expected state (e.g. `wm.IsBossTriggered()` if OnClear is "boss").
- Run: `go test ./internal/ecs/systems/enemy/... -v -run TestWaveManager_SpawnFlowAndWaveComplete`
- Expected: FAIL until logic is correct; fix implementation or test as needed.

**Step 3: Run all wave manager tests**

- Run: `go test ./internal/ecs/systems/enemy/... -v -run TestGyrussWaveManager`
- Expected: All PASS.

**Step 4: Commit**

```bash
git add internal/ecs/systems/enemy/gyruss_wave_manager_test.go
git commit -m "test: wave manager level start, spawn flow, and wave completion"
```

---

## Task 4: Spawner — spawn index and orbit angle

**Files:**
- Modify: `internal/ecs/systems/enemy/gyruss_spawner_test.go`
- Reference: spawner code that uses `GetCurrentGroupSpawnIndex` and orbit angle calculation

**Step 1: Write failing test for spawn index progression**

- Add test: set up world, spawner, and a wave with one group (count 3). Drive spawning (either by calling spawner with the same group config and incrementing spawn index, or by using wave manager + spawner together). Assert that the spawn index (or resulting entity positions/angles) differ for each of the 3 spawns (e.g. angles not all equal).
- Run: `go test ./internal/ecs/systems/enemy/... -v -run TestGyrussSpawner_SpawnIndexOrAngle`
- Expected: FAIL or PASS; if FAIL, implement or fix so that spawn index influences placement.

**Step 2: Run spawner tests**

- Run: `go test ./internal/ecs/systems/enemy/... -v -run TestGyrussSpawner`
- Expected: All PASS.

**Step 3: Commit**

```bash
git add internal/ecs/systems/enemy/gyruss_spawner_test.go
git commit -m "test: spawner spawn index / orbit angle distribution"
```

---

## Task 5: Collision outcome — enemy destroyed

**Files:**
- Locate where “enemy hit” leads to entity removal (e.g. collision system or game loop calling `DestroyEnemy`).
- Create or modify: test in that package (e.g. `internal/ecs/systems/collision/..._test.go` or `internal/ecs/systems/gyruss/gyruss_system_test.go`).

**Step 1: Write failing test**

- Add test: create world and GyrussSystem (or minimal set of systems). Spawn one enemy (or add one enemy entity with EnemyTag). Call the code path that destroys the enemy (e.g. `gs.DestroyEnemy(entity)`). Assert: entity is removed (world.Valid(entity) is false or entry no longer exists). Optionally assert no panic and score/event if applicable.
- Run: `go test ./internal/ecs/systems/gyruss/... -v -run TestDestroyEnemy` (or the chosen package).
- Expected: PASS (DestroyEnemy already exists); if not, implement or fix.

**Step 2: Commit**

```bash
git add <modified test file>
git commit -m "test: collision outcome — enemy destroyed and entity removed"
```

---

## Task 6: Integration test — stage 1 wave1 → wave2 → boss → complete

**Files:**
- Create: `internal/ecs/systems/gyruss/integration_test.go` (with `//go:build integration` at top)

**Step 1: Add build tag and test skeleton**

- First line: `//go:build integration`
- Package gyruss. In TestStage1Completes (or similar): create test world, config, logger, resource manager, GyrussSystem (reuse createTestGyrussSystem pattern from gyruss_system_test.go). Load stage 1 with `gs.LoadStage(1)`. Add `t.Skip("skipping until loop implemented")` and run test.
- Run: `go test -tags=integration ./internal/ecs/systems/gyruss/... -v -run TestStage1Completes`
- Expected: SKIP.

**Step 2: Implement run loop and assertions**

- Remove skip. Loop: call `gs.Update(ctx, dt)` with dt = 1/60 and a max duration (e.g. 15 seconds). After each Update, query for EnemyTag entities (excluding boss if desired); once a wave’s worth have spawned, destroy them via `gs.DestroyEnemy(entity)`. When boss appears (query for EnemyTypeBoss), destroy boss. Assert `gs.IsStageComplete()` before timeout. Use context.WithTimeout for the loop.
- Run: `go test -tags=integration ./internal/ecs/systems/gyruss/... -v -run TestStage1Completes -timeout=20s`
- Expected: PASS. If timeout or wrong state, fix timing or wave-completion logic (e.g. ensure level-start delay is satisfied and enemies are actually removed so wave completion can trigger).

**Step 3: Ensure CI runs integration tests**

- In `.github/workflows/ci.yml`, test step must include `-tags=integration`. If using `task test:all`, add a task or use explicit `go test ./... -race -tags=integration -timeout=60s`.
- Run full test suite locally: `go test ./... -race -tags=integration -timeout=60s`
- Expected: All PASS including integration test.

**Step 4: Commit**

```bash
git add internal/ecs/systems/gyruss/integration_test.go
git add .github/workflows/ci.yml
git commit -m "test: add headless integration test for stage 1 completion"
```

---

## Task 7: Wire CI to run integration and fix any failures

**Files:**
- Modify: `.github/workflows/ci.yml` if not already done
- Optionally modify: `Taskfile.test.yml` to add a `ci` task that runs `go test ./... -race -tags=integration -timeout=60s`

**Step 1: Confirm CI runs integration**

- Ensure workflow uses `go test ./... -race -tags=integration -timeout=60s` (or equivalent).
- Push and open a PR or push to branch; confirm CI job runs and test step includes integration.
- Expected: Job runs; fix any lint or path issues.

**Step 2: Fix any CI-only failures**

- If tests pass locally but fail in CI (e.g. missing assets, different working directory), fix paths or embed usage so that `go test` from repo root finds assets.
- Re-run CI until green.

**Step 3: Commit**

```bash
git add .github/workflows/ci.yml Taskfile.test.yml
git commit -m "ci: ensure integration tests run and fix CI test paths"
```

---

## Execution handoff

Plan complete and saved to `docs/plans/2025-02-19-mvp-testing-phase1.md`.

**Two execution options:**

1. **Subagent-driven (this session)** — I dispatch a fresh subagent per task, review between tasks, fast iteration.
2. **Parallel session (separate)** — You open a new session with executing-plans and run the plan task-by-task with checkpoints.

Which approach do you want?
