# Linting Issues Tracker

## Package Organization Issues
- [ ] Found packages `core` and `game` in core directory
  - Need to resolve package organization in `core/constants.go` and `core/stars.go`

## Import Issues
- [ ] Cannot import `github.com/jonesrussell/gimbal/player`
  - Related to package organization issues above

## Constants Issues
- [ ] `radius` redeclared in `player/constants.go`
  - Need to remove duplicate declaration (lines 11 and 26)
- [ ] `AngleStep` undefined in `game/player.go`
  - Need to define or import constant

## Function Signature Issues
- [ ] `NewGimlarGame` called with incorrect arguments in tests
  - Update calls in `game/game_test.go` (lines 17, 31, 48)
  - Function expects `(*zap.Logger, *Config)`, receiving only `(float64)`

## Test Issues
- [x] `MockHandler` missing `HandleInput` method implementation
- [ ] `calculateCoordinates` called with wrong number of arguments
  - Fix in `player/player_calculations_test.go` line 93
- [ ] Type mismatch in test comparisons
  - Fix in `player/player_calculations_test.go`:
    - Line 94: `float64` vs `int` comparison
    - Line 97: `float64` vs `int` comparison

## Progress
- Total Issues: 8
- Resolved: 1
- Remaining: 7 