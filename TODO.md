# TODO: Gimbal Refactor & Architecture Improvements

## High Priority (ASAP)

- [x] **Consolidate Game Configs**
    - Merge duplicate GameConfig structs (internal/common and internal/game) into a single source of truth.
    - Update all usages to reference the unified config.

- [x] **Refactor Build & Run Scripts**
    - Update Taskfile.yml and Makefile to use the new main.go location (root directory).
    - Remove references to the old cmd/gimbal path.

- [x] **Abstract Entities**
    - Refactor player and stars to use a common Entity interface.
    - Prepare code for ECS migration (ensure all game objects implement basic interfaces: Update, Draw, GetPosition, etc.).

- [x] **Introduce ECS Skeleton**
    - Add basic ECS scaffolding: EntityManager, Component registry, System runner.
    - Migrate player and stars to ECS-style management.

- [x] **Remove Dead Code & Duplicates**
    - Clean up any unused files, duplicate types, or legacy code from the old structure.

## Medium Priority

- [x] **ECS Migration Complete**
    - Successfully migrated from old entity system to Donburi ECS.
    - Removed all old entity code and dependencies.
    - Game now runs with ECS architecture.

- [x] **Event System**
    - Implement a simple event/message bus for decoupled communication (e.g., player death, score updates).

- [x] **Scene/State Management**
    - Abstract game states (menu, playing, game over) for better flow control and separation of concerns.

- [ ] **Resource Management**
    - Centralize asset loading (sprites, sounds) to avoid duplication and leaks.
    - Move assets to shared location accessible by ECS.

- [ ] **Testing**
    - Increase test coverage, especially for systems and components.
    - Use dependency injection for easier testing.

## Long-Term

- [ ] **Full ECS Migration**
    - Move all game logic (enemies, bullets, power-ups, etc.) to ECS.

- [ ] **Performance Profiling & Optimization**
    - Profile the game and optimize bottlenecks as new features are added.

- [ ] **Documentation**
    - Update README and add architecture docs as the codebase evolves. 