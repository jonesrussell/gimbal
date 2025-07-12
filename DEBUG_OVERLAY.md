# Debug Overlay System

The game now includes a comprehensive debug overlay system to help with development and troubleshooting.

## Features

### 1. Debug Controls
- **F3**: Toggle debug mode on/off (Basic level)
- **F4**: Cycle through debug levels (Basic → Detailed → Off)
- **Status**: Shows current debug level in the top-left corner

### 2. Debug Levels

#### Basic Level (F3)
- **Grid Overlay**: 50x50 pixel grid with subtle white lines
- **Performance Metrics**: Condensed single-line display
- **Entity Dots**: Colored dots for different entity types

#### Detailed Level (F4)
- **Everything from Basic level**
- **Mouse-Hover Entity Info**: Detailed info only for entities near mouse cursor
- **Entity Bounding Boxes**: Colored outlines for different entity types
- **Collision Debug**: Visual collision detection areas
- **Sprite Boundaries**: Sprite rendering boundaries

### 3. Performance Metrics
Displays in the top-left corner with dark background for readability:
- **Format**: `FPS:60 TPS:60 E:25 M:1024K [F3]`
- **FPS**: Current frames per second
- **TPS**: Current ticks per second  
- **E**: Total number of entities in the world
- **M**: Current memory usage in KB

### 4. Entity Debug Visualization
- **Green dots**: Player entities
- **Red dots**: Enemy entities
- **Yellow dots**: Projectile entities
- **Blue dots**: Star entities
- **White dots**: Other entities

### 5. Mouse-Hover Entity Info
When in Detailed mode, moving the mouse near entities (within 50 pixels) shows:
- **Position coordinates**: (x, y) position
- **Size information**: Width x Height dimensions
- **Sprite boundaries**: Visual sprite rendering area
- **Collision boxes**: Entity collision detection areas

### 6. Visual Hierarchy
- **Text backgrounds**: Semi-transparent black backgrounds behind all text
- **Entity colors**: Different colors for different entity types
- **Subtle effects**: Low alpha values to avoid overwhelming gameplay
- **Clean layout**: Organized information display

## Usage

1. **Start the game**
2. **Press F3** to enable basic debug mode
3. **Press F4** to cycle to detailed debug mode
4. **Move mouse** near entities to see detailed information
5. **Press F3 or F4** to disable debug mode

## Debug Information

### Entity Information (Detailed Mode Only)
When debug mode is in detailed level and mouse is near entities:
- Entity positions with coordinates
- Entity sizes (width x height)
- Sprite boundaries and centers
- Collision detection areas

### Performance Monitoring
- Monitor FPS/TPS for performance issues
- Track entity count for memory management
- Watch memory usage for potential leaks

### Visual Debugging
- Grid helps with positioning issues
- Collision boxes show detection areas
- Sprite boundaries reveal rendering issues
- Entity colors help identify different types

## Implementation Details

The debug system is implemented in:
- `internal/ecs/debug/debug_renderer.go` - Main debug renderer with level management
- `internal/ecs/scenes/manager.go` - Integration with scene manager
- `internal/input/handler.go` - F3/F4 key handling

### Key Improvements
- **Selective Display**: Entity info only shows for nearby entities
- **Multiple Levels**: Basic and detailed debug modes
- **Better Readability**: Text backgrounds and improved contrast
- **Entity Type Colors**: Visual distinction between entity types
- **Condensed Performance**: Single-line performance display
- **Mouse Integration**: Hover-based entity information display

The debug overlay is drawn on top of all game content and can be toggled at any time during gameplay.
