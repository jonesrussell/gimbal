# Graphics and Polish Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** One pass to remove the scanning grid, drop unused assets from loading, wire level 2 music, clean up the debug overlay (no grid by default), optimize image assets with ImageMagick/optipng/pngcrush, update docs, and apply general HUD/intro polish.

**Architecture:** Code changes in stage intro, resource manager (sprite/audio configs), debug renderer, and gameplay music selection. Asset optimization via a repeatable script. No new systems; removal and configuration only.

**Tech Stack:** Go, Ebiten, existing ECS; ImageMagick, optipng, pngcrush for assets. Design reference: `docs/plans/2026-02-20-graphics-polish-design.md`.

---

## Task 1: Remove scanning grid from stage intro

**Files:**
- Modify: `internal/scenes/stageintro/stage_intro.go`
- Modify: `internal/ecs/managers/resource/sprite_creation.go`

**Step 1:** Remove grid field and draw/load usage from stage intro.

In `internal/scenes/stageintro/stage_intro.go`:
- Delete the `gridOverlay *ebiten.Image` field from the struct (line ~34).
- Delete the entire draw block that uses `s.gridOverlay` (lines ~100–109: the "Draw scanning grid overlay" block).
- In `Enter()`, delete the block that loads `"scanning_grid"` (lines ~186–189: the "Load scanning grid" block).

**Step 2:** Remove scanning_grid from cutscene sprite config.

In `internal/ecs/managers/resource/sprite_creation.go`, in `loadCutsceneSprites`, remove the slice entry for `scanning_grid` (Name: `"scanning_grid"`, Path: `"cutscenes/scanning_grid.png"`). Keep the warp tunnel frame configs. If the slice becomes only warp frames, keep the loop; no need to remove the function.

**Step 3:** Delete the grid asset file if present.

If the file exists, delete it so it is no longer embedded:
- Delete: `assets/cutscenes/scanning_grid.png` (if it exists).

Embed is `//go:embed cutscenes/*.png` in `assets/assets.go`; removing the file is enough.

**Step 4:** Run tests and smoke check.

Run: `task test:all` (or `go test ./...`). Start the game and trigger the stage intro (e.g. complete stage 1 or jump to stage intro scene). Confirm no grid is drawn and no errors in log.

**Step 5:** Commit.

```bash
git add internal/scenes/stageintro/stage_intro.go internal/ecs/managers/resource/sprite_creation.go
# If file was deleted: git add assets/cutscenes/scanning_grid.png
git commit -m "chore: remove scanning grid from stage intro and asset"
```

---

## Task 2: Remove scanline_overlay from loading

**Files:**
- Modify: `internal/ecs/managers/resource/sprite_creation.go`

**Step 1:** Remove scanline_overlay from UI sprite config.

In `internal/ecs/managers/resource/sprite_creation.go`, in `loadUISprites`, remove the slice entry for `scanline_overlay` (Name: `"scanline_overlay"`, Path: `"ui/scanline_overlay.png"`). Do not remove `warning_overlay`.

**Step 2:** Verify no code references scanline_overlay.

Run: `grep -r scanline_overlay --include='*.go' .` (or grep in repo). Expected: no references. If any remain (e.g. in a scene), remove those references.

**Step 3:** Run tests.

Run: `task test:all`. Optional: start game and confirm no errors on scenes that use UI sprites.

**Step 4:** Commit.

```bash
git add internal/ecs/managers/resource/sprite_creation.go
git commit -m "chore: remove scanline_overlay from UI sprite loading"
```

---

## Task 3: Update ASSETS_REQUIRED for player_spritesheet and scanline

**Files:**
- Modify: `docs/ASSETS_REQUIRED.md`

**Step 1:** Mark player_spritesheet as not used / removed.

In `docs/ASSETS_REQUIRED.md`, update the Player Spritesheet row: set Status to "Removed" or "Not used" and Notes to "Removed from loading; not used in code." Optionally remove the row if you prefer to drop it from the list entirely.

**Step 2:** Remove or update scanline_overlay in docs.

If ASSETS_REQUIRED or any other doc lists `scanline_overlay`, remove that row or mark as "Removed from loading."

**Step 3:** Commit.

```bash
git add docs/ASSETS_REQUIRED.md
git commit -m "docs: update ASSETS_REQUIRED for removed sprites"
```

---

## Task 4: Wire level 2 music to stage 2

**Files:**
- Modify: `internal/ecs/managers/resource/audio.go`
- Modify: `internal/scenes/gameplay/playing.go`

**Step 1:** Add level 2 music to audio config.

In `internal/ecs/managers/resource/audio.go`, in the slice returned by `getAudioConfigs()` (or the inline config slice), add an entry for level 2:
- name: `"game_music_level_2"`
- path: `"sounds/game_music_level_2.ogg"`

Place it after `game_music_level_1`. Existing `LoadAllAudio` continues to load all configs and skips on error (no crash if file missing).

**Step 2:** Return level 2 music for stage 2 in gameplay.

In `internal/scenes/gameplay/playing.go`, add a constant for level 2 music (e.g. `musicTrackLevel2 = "game_music_level_2"`). In `getLevelMusicName()`, add a branch: when `levelManager.GetLevel() == 2`, return `musicTrackLevel2`. Keep level 1 returning `musicTrackLevel1` and fallback to `musicTrackMain` for other levels.

**Step 3:** Run tests and quick manual check.

Run: `task test:all`. Start game, reach stage 2 (or force level 2 in dev); confirm level 2 music plays when available, and no crash when file is missing.

**Step 4:** Commit.

```bash
git add internal/ecs/managers/resource/audio.go internal/scenes/gameplay/playing.go
git commit -m "feat: wire level 2 music to stage 2 playback"
```

---

## Task 5: Debug overlay — no grid in default view

**Files:**
- Modify: `internal/ecs/debug/debug_renderer.go`

**Step 1:** Draw grid only in Detailed mode.

In `internal/ecs/debug/debug_renderer.go`, in `Render()`, the call `dr.drawGrid(screen)` currently runs whenever debug is enabled. Change so the grid is drawn only when `dr.level == DebugDetailed`. Example: wrap `dr.drawGrid(screen)` in `if dr.level == DebugDetailed { dr.drawGrid(screen) }`. Performance metrics (and entity/collision debug) remain as-is: Basic shows performance only, Detailed adds grid and entity/collision info.

**Step 2:** Confirm default is off.

In `NewDebugRenderer`, `enabled` is already `false`. No change needed unless your codebase sets it elsewhere; verify with grep for `showDebugInfo` / `enabled` that default is off.

**Step 3:** Document toggle key.

In `CLAUDE.md` or `docs/CODING_STANDARDS.md`, add a short line: e.g. "Debug overlay: off by default; toggle with [key, e.g. F3]. Basic = FPS/entity count; Detailed = adds grid and entity/collision debug."

**Step 4:** Run tests and quick manual check.

Run: `task test:all`. Start game, toggle debug: Basic view has no grid; cycle to Detailed to see grid.

**Step 5:** Commit.

```bash
git add internal/ecs/debug/debug_renderer.go CLAUDE.md
git commit -m "chore: show debug grid only in Detailed mode; document toggle"
```

---

## Task 6: Asset optimization script and one-off run

**Files:**
- Create: `scripts/optimize-assets.sh` (or `docs/asset-optimization.md` if you prefer doc-only)
- Modify: image files under `assets/sprites/`, `assets/ui/`, `assets/cutscenes/` (after Task 1, no scanning_grid)

**Step 1:** Add optimization script or doc.

Create a script that:
1. Strips metadata: e.g. `mogrify -strip assets/sprites/*.png assets/ui/*.png assets/cutscenes/*.png` (adjust paths to your layout; if assets live in a single tree, run over that).
2. Resize only if needed: document or script a check (e.g. compare to ASSETS_REQUIRED/sprite_creation sizes); only resize files larger than max needed (e.g. 2x display size). If all assets are already at or below target sizes, skip resize.
3. PNG compression: run `optipng -o2` and/or `pngcrush` on the same PNGs. Example: `for f in assets/sprites/*.png; do optipng -o2 "$f"; done` and similar for ui/cutscenes.

Use ImageMagick for strip (and optional resize); use optipng/pngcrush for lossless compression. Document in the script or in `docs/asset-optimization.md` the exact commands and the rule "resize only if oversized."

**Step 2:** Run script on a copy or branch.

Run the script (or commands) against the repo's asset directories. Prefer running on a copy first to verify no visual regressions.

**Step 3:** Replace originals and run game.

Replace original assets with optimized versions. Start game and do a quick visual check: sprites and UI look correct.

**Step 4:** Commit.

```bash
git add scripts/optimize-assets.sh assets/
# or: git add docs/asset-optimization.md assets/
git commit -m "chore: add asset optimization script and run on sprites/ui/cutscenes"
```

---

## Task 7: General polish (HUD and stage intro)

**Files:**
- Modify: `internal/ui/responsive/responsive_ui.go` and/or HUD builder/layout code
- Modify: `internal/scenes/stageintro/stage_intro.go` (optional small timing/fade)

**Step 1:** HUD alignment and spacing.

In the responsive UI (or wherever HUD layout is defined), ensure one clear anchor (e.g. top-left or top-right), consistent padding from screen edge, and spacing between score, lives, and level so nothing is cramped or overlapping. Adjust only layout constants or container padding; no new widgets.

**Step 2:** Stage intro balance (optional).

If the stage intro feels empty without the grid, consider a small tweak: e.g. slightly longer fade-in for the planet or a brief delay before route text. Keep changes minimal (e.g. one constant).

**Step 3:** Run game and verify.

Start game; check HUD readability and stage intro flow. Run `task test:all`.

**Step 4:** Commit.

```bash
git add internal/ui/responsive/responsive_ui.go internal/scenes/stageintro/stage_intro.go
git commit -m "polish: HUD alignment and optional stage intro timing"
```

---

## Task 8: Final smoke test and doc update

**Files:**
- Modify: `docs/ASSETS_REQUIRED.md` (if any remaining updates)
- Read: `docs/plans/2026-02-20-graphics-polish-design.md`

**Step 1:** Smoke checklist.

- Start game; play to stage intro: no grid, no errors.
- Play to stage 2: level 2 music plays (if file present).
- Toggle debug: off by default; Basic = no grid; Detailed = grid + entity debug.
- Check HUD: score, lives, level readable and aligned.
- Complete a stage and see transition; no full-screen artifacts.

**Step 2:** Run full test suite.

Run: `task test:all` (or `go test ./...` with race if desired). All tests pass.

**Step 3:** Update ASSETS_REQUIRED if needed.

Ensure ASSETS_REQUIRED reflects: scanning_grid removed, scanline_overlay removed, player_spritesheet not used, level 2 music used for stage 2.

**Step 4:** Commit any final doc tweaks.

```bash
git add docs/ASSETS_REQUIRED.md
git commit -m "docs: final ASSETS_REQUIRED updates for graphics polish pass"
```

---

## Execution options

Plan complete and saved to `docs/plans/2026-02-20-graphics-polish.md`.

**Two execution options:**

1. **Subagent-driven (this session)** — I dispatch a fresh subagent per task, review between tasks, fast iteration.

2. **Parallel session (separate)** — Open a new session with executing-plans and run through the plan with checkpoints.

Which approach do you want?
