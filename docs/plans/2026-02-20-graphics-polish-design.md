# Graphics and Polish Design

**Date:** 2026-02-20  
**Status:** Approved  
**Scope:** One pass covering grid removal, lives display, asset optimization, unused-asset cleanup, debug overlay, and general visual polish.

---

## 1. Scope and goals

**In scope**
- Remove the scanning grid from the stage intro and delete the asset; remove its sprite registration and any fallback.
- Decide lives display: keep heart icons (Gyruss-style lives) or switch to number-only; implement that choice in the HUD.
- Optimize all image assets with ImageMagick (and optipng/pngcrush): consistent sizing, format (PNG), compression; document or script so the process is repeatable.
- Audit assets: remove or repurpose unused ones; for scanline_overlay, player_spritesheet, level-2 music—implement where they add value, otherwise remove from loading or document as future.
- Debug overlay: clarify what’s shown (e.g. FPS, entity count, stage state), default off or behind a key, tidy layout/readability.
- General polish: HUD layout and styling, stage intro look without the grid, transition visuals, and quick wins for consistency (e.g. colours, alignment).

**Out of scope**
- New gameplay features, new stages, or new enemy types.
- Full asset pipeline (e.g. Makefile/task that runs on every asset add); we only document/script the one-off optimization.
- Sound effects or new music; only decisions about existing music (e.g. level-2 track) and asset list updates.
- Web/mobile-specific optimizations (e.g. separate atlases or formats per platform).

**Success criteria**
- No scanning grid in the game or repo; stage intro still reads clearly.
- Lives are shown in one agreed form (hearts or number) and match game behaviour.
- Image assets are optimized and the process is repeatable (script or doc).
- Unused assets are either used or removed from loading; ASSETS_REQUIRED (or equivalent) updated.
- Debug overlay is off by default (or behind a key), readable, and only shows agreed info.
- Game looks more consistent and polished (HUD, intro, transitions) without changing behaviour.

---

## 2. Scanning grid removal

**Code changes**
- **Stage intro** (`internal/scenes/stageintro/stage_intro.go`): Remove the `gridOverlay` field, the draw block that uses it (lines ~100–109), and in `Enter()` remove the load of `"scanning_grid"` (lines ~186–189). Stage intro continues to show: fade-in, "STAGE N", planet sprite, route text, message.
- **Sprite registration** (`internal/ecs/managers/resource/sprite_creation.go`): In `loadCutsceneSprites`, remove the `scanning_grid` entry (path `cutscenes/scanning_grid.png`). If that leaves only warp frames, keep the function and load only warp tunnel frames.
- No other references use `scanning_grid`.

**Asset deletion**
- Delete `assets/cutscenes/scanning_grid.png` (or the actual path under the project’s asset layout). If the repo uses embedded assets, remove it from the embed list so the file is no longer shipped.

**Fallback**
- Removing the sprite config removes the fallback; no replacement needed since the grid is no longer drawn.

**Testing**
- Run stage intro: no grid overlay, no errors; stage number, planet, and route text still display correctly.

---

## 3. Lives display

**Context:** Original Gyruss used lives (one hit = one life), not a health bar. The game already models this; the HUD shows heart icons plus lives value.

**Decision:** Keep heart icons for lives. No HUD behaviour change for lives in this pass. Optional: small layout/readability tweaks for the lives area in the general polish pass.

---

## 4. Asset optimization (ImageMagick)

**Goal:** All image assets are optimized for size and consistency without changing how they look in-game. Process is documented or scripted for repeatability.

**What gets optimized**
- Every image referenced by the game: `sprites/`, `ui/`, `cutscenes/` (warp frames only after grid removal).

**How**
- **Strip metadata:** `convert -strip` or `mogrify -strip`.
- **Resize only if oversized:** Compare each asset to its use in code (sprite_creation fallback sizes, ASSETS_REQUIRED). Resize to max needed (e.g. 2x for hidpi) only when file is larger; otherwise leave dimensions.
- **PNG compression:** Lossless: ImageMagick `-quality` for PNG (e.g. 90–95) and/or run `optipng` and `pngcrush` (both installed). No lossy re-encode unless explicitly chosen for a specific asset.
- **Format:** Keep PNG for sprites/UI.
- **Backups:** Run on copies or in a branch until satisfied; then replace in repo.

**Deliverables**
- Script or doc: shell script (e.g. `scripts/optimize-assets.sh` or under `task`) that runs ImageMagick and optional optipng/pngcrush on `assets/sprites`, `assets/ui`, `assets/cutscenes`; or a short doc in `docs/` with exact commands and max-dimension rules.
- One-off run applied to current assets; commit optimized files.

**Dependencies:** ImageMagick, optipng, pngcrush.

---

## 5. Unused and missing assets

**Scanning grid:** Remove from loading and delete file (Section 2); update ASSETS_REQUIRED if listed.

**scanline_overlay:** Remove from loading and from docs; idea didn’t pan out. Optionally delete the file or leave it out of embed.

**warning_overlay:** Keep; used by boss intro.

**player_spritesheet:** Remove from loading; do not implement animation in this pass. Document as removed/not used or drop from ASSETS_REQUIRED.

**Level 2 music:** Wire `game_music_level_2.ogg` to stage 2 playback in this pass. Use same pattern as existing level music; if file missing or fails to load, fall back and log; don’t crash.

**Missing assets (from ASSETS_REQUIRED):** No new art or audio creation. Keep “missing” list as-is; optionally note that this pass did not add new assets.

**Deliverables**
- Code: load only assets that are used; scanning_grid and scanline_overlay removed; player_spritesheet removed from loading.
- Level 2 music wired to stage 2.
- ASSETS_REQUIRED (and any asset list) updated to match.

---

## 6. Debug overlay cleanup

**Goal:** Debug overlay off by default, shows only agreed info, readable layout. No new debug features.

**Proposed behaviour**
- **Default:** Off.
- **Toggle:** One key (e.g. F3 or existing) to show/hide; document in CLAUDE.md or dev doc.
- **Contents:** FPS, entity count, current scene/stage if available; optionally one line per slow system (e.g. over 5ms). Do not show the 50px debug grid in default view.
- **Layout:** One compact block (e.g. top-left or top-right), small font, high contrast.
- **Code:** Single place responsible for what to draw; remove or gate duplicate overlay drawing.

**Deliverables**
- Overlay off by default; one key toggles; shows FPS, entity count, scene/stage; no 50px grid by default; compact, readable layout.
- Toggle key and “no grid by default” documented.

---

## 7. General polish

**HUD:** Consistent anchor (e.g. top-left or top-right), padding, spacing; readable text and contrast; no new widgets.

**Stage intro:** Centered, balanced layout without grid; optional small timing/fade tweak so the intro doesn’t feel empty.

**Transitions:** Full-screen overlays match window size; consistent opacity if touched.

**Global:** No broad palette change; fix contrast only where an element is hard to read. Keep Press Start 2P; no font change.

**Deliverables:** HUD and intro polished; transitions correct; no new assets or UI behaviour.

---

## 8. Error handling and testing

**Error handling**
- No code requests removed sprites (scanning_grid, scanline_overlay); remove lookups and fallbacks.
- Stage intro: if planet sprite missing, keep current behaviour (stage number + route text).
- Level 2 music: fall back and log if missing or load fails; don’t crash.
- Asset lists and embeds stay in sync after optimization and deletions.

**Testing**
- Manual: stage intro (no grid, no errors); stage 2 with level 2 music; debug overlay (off by default, key toggles, no grid); HUD layout; one full transition.
- Automated: existing tests stay green; add tests only where low-cost (e.g. stage 2 music or debug default off).
- Regression: run full test suite and smoke run before considering the pass done.

**Deliverables**
- No requests for removed sprites; level 2 music has safe fallback.
- Manual checklist and existing tests green; optional small tests; smoke pass in implementation plan.
