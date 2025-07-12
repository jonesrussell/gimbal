# UI System Documentation

## Architecture Overview

The UI system is built on a component-based architecture using the `UIElement` interface. It provides:

- **Component-based system** using `UIElement` interface
- **Layout system** for automatic positioning and arrangement
- **Theme system** for consistent styling across components
- **Debug mode** for development and troubleshooting

## Core Interface

All UI components implement the `UIElement` interface:

```go
type UIElement interface {
    Draw(renderer *UIRenderer, pos Position)
    GetSize() (width, height float64)
}
```

## Creating New Components

### 1. Implement UIElement Interface

```go
type MyComponent struct {
    // Component data
}

func (mc *MyComponent) Draw(renderer *UIRenderer, pos Position) {
    // Drawing logic
}

func (mc *MyComponent) GetSize() (width, height float64) {
    // Size calculation
    return width, height
}
```

### 2. Add Styling Through Theme System

```go
// Use DefaultTheme for consistent styling
textStyle := TextStyle{
    Font:  DefaultTheme.Fonts.HUD,
    Color: DefaultTheme.Colors.Text,
    Size:  16,
    Alignment: AlignCenter,
}
```

### 3. Use Layouts for Automatic Positioning

```go
// Horizontal arrangement
layout := NewHorizontalLayout(8.0)
layout.Add(textElement)
layout.Add(iconElement)
layout.Draw(renderer, pos)
```

## Available Components

### TextElement
Styled text rendering with alignment support.

```go
text := NewTextWithStyle("Hello World", TextStyle{
    Font:      font,
    Color:     color.White,
    Size:      16,
    Alignment: AlignCenter,
})
```

### IconElement
Scaled sprite rendering with consistent sizing.

```go
icon := NewIcon(sprite, MediumIcon) // 32px size
```

### LivesDisplay
Composite component showing hearts + text.

```go
lives := NewLivesDisplay(currentLives, maxLives)
```

### ScoreDisplay
Formatted score text with consistent styling.

```go
score := NewScoreDisplay(scoreValue)
```

## Layout System

### HorizontalLayout
Arranges elements left-to-right with consistent spacing.

```go
layout := NewHorizontalLayout(spacing)
layout.Add(element1)
layout.Add(element2)
layout.Draw(renderer, pos)
```

### VerticalLayout
Arranges elements top-to-bottom with consistent spacing.

```go
layout := NewVerticalLayout(spacing)
layout.Add(element1)
layout.Add(element2)
layout.Draw(renderer, pos)
```

## Positioning Helpers

### Basic Positioning
- `TopLeft(x, y)` - Position at top-left coordinates
- `TopRightRelative(screenWidth, x, y)` - Position from top-right

### Responsive Positioning
- `CenterHorizontal(screenWidth, elementWidth)` - Center horizontally
- `CenterVertical(screenHeight, elementHeight)` - Center vertically
- `ResponsivePosition(screenWidth, screenHeight, anchor)` - Anchor-based positioning

### Anchor Options
- `"center"` - Screen center
- `"top-center"` - Top center with margin
- `"bottom-center"` - Bottom center with margin

## Theme System

The theme system provides consistent styling across all components:

```go
// Colors
DefaultTheme.Colors.Text      // Primary text color
DefaultTheme.Colors.TextLight // Secondary text color
DefaultTheme.Colors.Heart     // Heart icon color
DefaultTheme.Colors.Warning   // Warning text color
DefaultTheme.Colors.Debug     // Debug overlay color

// Fonts
DefaultTheme.Fonts.UI    // General UI text
DefaultTheme.Fonts.HUD   // HUD elements
DefaultTheme.Fonts.Title // Title text

// Sizes
DefaultTheme.Sizes.SmallIcon  // 24px
DefaultTheme.Sizes.MediumIcon // 32px
DefaultTheme.Sizes.LargeIcon  // 48px
DefaultTheme.Sizes.Padding    // 8px
DefaultTheme.Sizes.Margin     // 16px
```

## Usage Examples

### Basic Text Display
```go
text := NewTextWithStyle("Score: 1000", TextStyle{
    Font:  DefaultTheme.Fonts.HUD,
    Color: DefaultTheme.Colors.Text,
    Size:  16,
})
text.Draw(renderer, TopLeft(10, 10))
```

### Centered Text
```go
text := NewTextWithStyle("Game Over", TextStyle{
    Font:      DefaultTheme.Fonts.Title,
    Color:     DefaultTheme.Colors.Text,
    Size:      24,
    Alignment: AlignCenter,
})
text.Draw(renderer, ResponsivePosition(screenWidth, screenHeight, "center"))
```

### Horizontal Layout with Icons
```go
layout := NewHorizontalLayout(8.0)
layout.Add(NewTextWithStyle("Lives:", textStyle))
for i := 0; i < lives; i++ {
    layout.Add(NewIcon(heartSprite, SmallIcon))
}
layout.Draw(renderer, TopLeft(20, 20))
```

### Vertical Layout for Menu
```go
layout := NewVerticalLayout(16.0)
layout.Add(NewTextWithStyle("Start Game", menuStyle))
layout.Add(NewTextWithStyle("Options", menuStyle))
layout.Add(NewTextWithStyle("Quit", menuStyle))
layout.Draw(renderer, ResponsivePosition(screenWidth, screenHeight, "center"))
```

## Debug Mode

Enable debug mode to see component boundaries and positioning:

```go
renderer.SetDebug(true)
```

Debug mode shows:
- Component bounding boxes
- Corner markers for precise positioning
- Visual feedback for layout calculations

## Best Practices

1. **Use Theme System**: Always use `DefaultTheme` for consistent styling
2. **Layout Components**: Use layouts for automatic positioning and spacing
3. **Responsive Design**: Use responsive positioning helpers for different screen sizes
4. **Text Alignment**: Use alignment options for centered or right-aligned text
5. **Debug Mode**: Enable debug mode during development for visual feedback

## Performance Considerations

- Text measurement is cached internally
- Layout calculations are optimized for common use cases
- Debug rendering is only active when debug mode is enabled
- Component reuse is encouraged for similar elements 