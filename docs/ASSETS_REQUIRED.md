# Game Assets Required

This document lists all assets required for the Gimbal game, including what exists and what is missing.

## Sprites

| Asset Name | Path | Status | Size | Notes |
|------------|------|--------|------|-------|
| Player | `sprites/player.png` | ✅ Exists | 32x32 | Main player sprite |
| Player Spritesheet | `sprites/player_spritesheet.png` | ✅ Exists | - | Not currently used in code |
| Heart | `sprites/heart.png` | ✅ Exists | 16x16 | Health indicator |
| Enemy (Basic) | `sprites/enemy.png` | ✅ Exists | 32x32 | Basic enemy type |
| Enemy Heavy | `sprites/enemy_heavy.png` | ✅ Exists | 32x32 | Heavy enemy type |
| Enemy Boss (Generic) | `sprites/enemy_boss.png` | ✅ Exists | 64x64 | Generic boss sprite (fallback) |
| Enemy Ammo | `sprites/enemy_ammo.png` | ✅ Exists | 6x6 | Basic enemy projectiles |
| Enemy Heavy Ammo | `sprites/enemy_heavy_ammo.png` | ✅ Exists | 6x6 | Heavy enemy projectiles |
| Star | `sprites/star.png` | ✅ Exists | - | UI star sprite |
| **Satellite** | `sprites/satellite.png` | ❌ **Missing** | 32x32 | Satellite enemy type (used in stages 1-6) |
| **Player Ammo** | `sprites/player_ammo.png` | ❌ **Missing** | 4x4 | Player projectiles (currently generated programmatically) |
| **Power-up: Double Shot** | `sprites/powerup_double_shot.png` | ❌ **Missing** | 16x16 | Double shot power-up (currently generated programmatically) |
| **Power-up: Extra Life** | `sprites/powerup_extra_life.png` | ❌ **Missing** | 16x16 | Extra life power-up (currently generated programmatically) |
| **Boss: Earth** | `sprites/earth_boss.png` | ❌ **Missing** | 64x64 | Stage 1 boss |
| **Boss: Mars** | `sprites/mars_boss.png` | ❌ **Missing** | 72x72 | Stage 2 boss |
| **Boss: Jupiter** | `sprites/jupiter_boss.png` | ❌ **Missing** | 80x80 | Stage 3 boss |
| **Boss: Saturn** | `sprites/saturn_boss.png` | ❌ **Missing** | 88x88 | Stage 4 boss |
| **Boss: Uranus** | `sprites/uranus_boss.png` | ❌ **Missing** | 96x96 | Stage 5 boss |
| **Boss: Final** | `sprites/final_boss.png` | ❌ **Missing** | 112x112 | Stage 6 final boss |

## Audio

| Asset Name | Path | Status | Notes |
|------------|------|--------|-------|
| Main Menu Music | `sounds/game_music_main.ogg` | ✅ Exists | Main menu background music |
| Level 1 Music | `sounds/game_music_level_1.ogg` | ✅ Exists | Stage 1 gameplay music |
| Level 2 Music | `sounds/game_music_level_2.ogg` | ✅ Exists | Not currently loaded in code |
| Boss Music | `sounds/game_music_boss.ogg` | ✅ Exists | Boss fight music |
| **Level 3 Music** | `sounds/game_music_level_3.ogg` | ❌ **Missing** | For stage 3+ |
| **Level 4 Music** | `sounds/game_music_level_4.ogg` | ❌ **Missing** | For stage 4+ |
| **Level 5 Music** | `sounds/game_music_level_5.ogg` | ❌ **Missing** | For stage 5+ |
| **Level 6 Music** | `sounds/game_music_level_6.ogg` | ❌ **Missing** | For stage 6 |
| **Sound Effects** | `sounds/` | ❌ **Missing** | No sound effects currently implemented |

## Fonts

| Asset Name | Path | Status | Notes |
|------------|------|--------|-------|
| Press Start 2P | `fonts/PressStart2P.ttf` | ✅ Exists | Main game font |

## Stage Configuration Files

| Asset Name | Path | Status | Notes |
|------------|------|--------|-------|
| Stage 1 | `stages/stage_01.json` | ✅ Exists | Earth to Mars |
| Stage 2 | `stages/stage_02.json` | ✅ Exists | Mars to Jupiter |
| Stage 3 | `stages/stage_03.json` | ✅ Exists | Jupiter to Saturn |
| Stage 4 | `stages/stage_04.json` | ✅ Exists | Saturn to Uranus |
| Stage 5 | `stages/stage_05.json` | ✅ Exists | Uranus to Neptune |
| Stage 6 | `stages/stage_06.json` | ✅ Exists | Final stage |

## Entity Configuration Files

| Asset Name | Path | Status | Notes |
|------------|------|--------|-------|
| Player Config | `entities/player.json` | ✅ Exists | Player entity configuration |
| Enemies Config | `entities/enemies.json` | ✅ Exists | Enemy types configuration |

## Summary

### ✅ Existing Assets
- **Sprites**: 9/19 (47%)
- **Audio**: 4/9 (44%)
- **Fonts**: 1/1 (100%)
- **Config Files**: 8/8 (100%)

### ❌ Missing Assets
- **Sprites**: 10 missing
  - Satellite enemy sprite
  - Player ammo sprite
  - 2 power-up sprites
  - 6 boss-specific sprites
- **Audio**: 5 missing
  - 4 level music tracks (levels 3-6)
  - Sound effects (not yet implemented)

### Notes
- Player projectiles and power-ups are currently generated programmatically with colored rectangles, but sprites would improve visual quality
- Boss sprites currently fall back to generic `enemy_boss.png`, but stage configs reference specific boss types
- Satellite enemy type is used in multiple stages but has no sprite asset
- Level music tracks beyond level 2 are not loaded in the code, but stages 3-6 exist
- Sound effects system is not yet implemented
