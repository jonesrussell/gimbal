# EbitenUI Responsive Design System

This package implements a professional responsive design system using EbitenUI for the Gimbal game, providing industry-standard UI components with automatic scaling and positioning.

## 🚀 Features

### **Professional UI Framework**
- **EbitenUI Integration**: Built on the industry-standard EbitenUI library
- **Responsive Layouts**: Automatic scaling and positioning across all devices
- **Widget System**: Professional buttons, text, containers, and progress bars
- **Nine-Slice Support**: Crisp UI graphics at any resolution
- **Event Handling**: Built-in input handling for all devices

### **Responsive Design**
- **Device Classification**: Automatic detection of mobile, tablet, desktop, and ultrawide
- **Adaptive Layouts**: UI elements automatically adjust to screen size
- **Touch Optimization**: Proper touch targets for mobile devices
- **Orientation Support**: Portrait and landscape awareness

### **Gaming UI Components**
- **Lives Display**: Heart icons with responsive positioning
- **Score Display**: Professional score panel with background
- **Health Bar**: Dynamic progress bar with color coding
- **Ammo Counter**: Visual ammo indicators with icons

## 📱 Device Support

### **Responsive Breakpoints**
```go
// Mobile: < 768px width
// Tablet: 768px - 1024px width  
// Desktop: 1024px - 1920px width
// Ultrawide: > 1920px width
```

### **Layout Adaptations**
- **Mobile**: Compact UI with larger touch targets
- **Tablet**: Balanced layout with medium-sized elements
- **Desktop**: Standard layout with optimal spacing
- **Ultrawide**: Spread UI elements across wider screens

## 🎮 UI Components

### **ResponsiveUI Structure**
```go
type ResponsiveUI struct {
    ui *ebitenui.UI
    
    // UI containers
    hudContainer     *widget.Container
    livesContainer   *widget.Container
    scoreContainer   *widget.Container
    healthContainer  *widget.Container
    ammoContainer    *widget.Container
    
    // Widgets
    livesText        *widget.Text
    scoreText        *widget.Text
    healthBar        *widget.ProgressBar
    ammoIcons        []*widget.Graphic
}
```

### **Layout System**
- **Anchor Layout**: Professional positioning system
- **Row Layout**: Horizontal arrangement for lives and ammo
- **Responsive Padding**: Automatic spacing based on device class
- **Nine-Slice Backgrounds**: Scalable UI panels

## 🔧 Usage

### **Basic Setup**
```go
// Create responsive UI system
font := resourceManager.GetDefaultFont()
heartSprite, _ := resourceManager.GetSprite("heart")
ammoSprite, _ := resourceManager.GetSprite("ammo")

responsiveUI := ui.NewResponsiveUI(font, heartSprite, ammoSprite)
```

### **Game Integration**
```go
// In game Update loop
func (g *Game) Update() error {
    // Update EbitenUI system
    g.responsiveUI.Update()
    
    // Update game logic
    // ...
    
    return nil
}

// In game Draw method
func (g *Game) Draw(screen *ebiten.Image) {
    // Draw game content
    // ...
    
    // Draw responsive UI
    g.responsiveUI.Draw(screen)
}
```

### **Updating UI State**
```go
// Update lives display
responsiveUI.UpdateLives(3)

// Update score display
responsiveUI.UpdateScore(1250)

// Update health bar (0.0 to 1.0)
responsiveUI.UpdateHealth(0.75)

// Update ammo counter
responsiveUI.UpdateAmmo(8)
```

### **Responsive Layout Updates**
```go
// Update layout for current screen size
width, height := ebiten.WindowSize()
responsiveUI.UpdateResponsiveLayout(width, height)
```

## 🎯 Layout System

### **Anchor Layout Positioning**
```go
// Top-left positioning
widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
    HorizontalPosition: widget.AnchorLayoutPositionStart,
    VerticalPosition:   widget.AnchorLayoutPositionStart,
    Padding: widget.Insets{
        Left: 20, Top: 20,
    },
})

// Top-right positioning
widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
    HorizontalPosition: widget.AnchorLayoutPositionEnd,
    VerticalPosition:   widget.AnchorLayoutPositionStart,
    Padding: widget.Insets{
        Right: 20, Top: 20,
    },
})
```

### **Row Layout for Horizontal Elements**
```go
widget.ContainerOpts.Layout(widget.NewRowLayout(
    widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
    widget.RowLayoutOpts.Spacing(10),
))
```

## 🎨 UI Components

### **Lives Display**
- **Location**: Top-left corner
- **Layout**: Horizontal row with text and heart icons
- **Responsive**: Automatic spacing and sizing
- **Dynamic**: Updates based on current lives

### **Score Display**
- **Location**: Top-right corner
- **Background**: Semi-transparent panel
- **Text**: Centered score display
- **Responsive**: Adapts to screen width

### **Health Bar**
- **Location**: Bottom-left corner
- **Type**: Progress bar with background
- **Colors**: Green (high), Yellow (medium), Red (low)
- **Dynamic**: Updates based on health percentage

### **Ammo Counter**
- **Location**: Bottom-right corner
- **Layout**: Horizontal row with text and ammo icons
- **Limit**: Maximum 10 visible icons
- **Responsive**: Automatic spacing

## 📱 Responsive Features

### **Device Classification**
```go
func (ui *ResponsiveUI) UpdateResponsiveLayout(width, height int) {
    // Determine device class
    if width < 768 {
        ui.deviceClass = "mobile"
    } else if width < 1024 {
        ui.deviceClass = "tablet"
    } else if width > 1920 {
        ui.deviceClass = "ultrawide"
    } else {
        ui.deviceClass = "desktop"
    }
}
```

### **Automatic Scaling**
- **Font Scaling**: Text automatically scales with screen size
- **Icon Scaling**: Sprites scale appropriately for device
- **Padding Adaptation**: Spacing adjusts for touch vs mouse
- **Layout Optimization**: Elements reposition for optimal viewing

## 🔧 Integration with Game Systems

### **Health System Integration**
```go
// Update health bar from health system
healthSystem := game.GetHealthSystem()
if healthSystem != nil {
    current, max := healthSystem.GetPlayerHealth()
    healthPercent := float64(current) / float64(max)
    responsiveUI.UpdateHealth(healthPercent)
}
```

### **Score System Integration**
```go
// Update score from score manager
scoreManager := game.GetScoreManager()
if scoreManager != nil {
    score := scoreManager.GetScore()
    responsiveUI.UpdateScore(score)
}
```

### **Resource Management**
```go
// Get sprites from resource manager
heartSprite, _ := resourceManager.GetSprite("heart")
ammoSprite, _ := resourceManager.GetSprite("ammo")

// Create responsive UI with resources
responsiveUI := ui.NewResponsiveUI(font, heartSprite, ammoSprite)
```

## 🎯 Performance Benefits

### **EbitenUI Advantages**
- **Optimized Rendering**: Hardware-accelerated UI rendering
- **Efficient Layouts**: Professional layout algorithms
- **Memory Management**: Automatic cleanup and optimization
- **Event Handling**: Built-in input processing

### **Responsive Optimization**
- **Layout Caching**: Efficient layout calculations
- **Conditional Updates**: Only update when values change
- **Batch Rendering**: Optimized drawing operations
- **Memory Efficiency**: Minimal memory footprint

## 🚀 Migration from Custom UI

### **Benefits of EbitenUI**
- ✅ **Professional Framework**: Industry-standard UI library
- ✅ **Reduced Complexity**: No custom UI implementation needed
- ✅ **Better Maintainability**: Well-documented, actively maintained
- ✅ **Touch Support**: Built-in mobile input handling
- ✅ **Scalable Architecture**: Easy to add new UI elements

### **Migration Steps**
1. **Remove Custom UI**: Delete `internal/ecs/ui/` package
2. **Add EbitenUI**: `go get github.com/ebitenui/ebitenui`
3. **Update Game**: Replace custom UI with EbitenUI system
4. **Update Scenes**: Remove UI rendering from individual scenes
5. **Test Responsive**: Verify layout across different screen sizes

## 🔮 Future Enhancements

### **Planned Features**
- **Menu System**: Pause menu, settings, main menu
- **Dialog System**: Confirmation dialogs, tooltips
- **Animation System**: Smooth transitions and effects
- **Theme System**: Customizable UI themes
- **Accessibility**: Screen reader support, high contrast

### **Advanced Features**
- **Localization**: Multi-language support
- **Custom Widgets**: Game-specific UI components
- **Animation Framework**: Smooth UI animations
- **Theme Engine**: Dynamic theme switching
- **Performance Monitoring**: UI performance metrics

## 📚 Best Practices

### **1. Responsive Design**
- Always test on multiple screen sizes
- Use relative positioning when possible
- Consider touch targets on mobile devices
- Implement proper safe areas

### **2. Performance**
- Update UI only when values change
- Use efficient layout algorithms
- Monitor memory usage
- Optimize sprite loading

### **3. User Experience**
- Maintain consistent spacing
- Use appropriate colors and contrast
- Provide clear visual feedback
- Ensure accessibility compliance

### **4. Code Organization**
- Separate UI logic from game logic
- Use clear naming conventions
- Document complex layouts
- Maintain clean architecture

## 🐛 Troubleshooting

### **Common Issues**

**UI not updating**
- Check if Update() is called in game loop
- Verify state values are changing
- Ensure proper initialization

**Layout issues**
- Test on different screen sizes
- Check anchor layout positioning
- Verify padding and spacing values

**Performance problems**
- Monitor update frequency
- Check for unnecessary redraws
- Optimize sprite usage

**Touch input issues**
- Verify touch target sizes
- Check event handling setup
- Test on actual mobile devices

## 📄 License

This responsive design system uses EbitenUI and follows the same licensing terms as the Gimbal game project.

---

**Built with EbitenUI for professional game development** 