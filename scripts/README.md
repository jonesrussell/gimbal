# Scripts Directory

This directory contains analysis and utility scripts for the Gimbal project.

## Architecture

The scripts follow a modular architecture with shared libraries:

```
scripts/
├── analyze-package.sh          # Original monolithic script (483 lines)
├── analyze-package-new.sh      # New modular orchestrator (189 lines)
├── lib/                        # Shared library modules
│   ├── common.sh              # Shared utilities (colors, helpers)
│   ├── package-info.sh        # Basic package information
│   ├── file-analysis.sh       # File overview and summaries
│   ├── import-analysis.sh     # Import analysis
│   ├── struct-analysis.sh     # Struct analysis
│   ├── interface-analysis.sh  # Interface analysis
│   ├── method-analysis.sh     # Method analysis
│   ├── dependency-analysis.sh # Dependency analysis
│   └── metrics.sh             # Code quality metrics
└── tools/                      # Standalone analysis tools
    ├── complexity-checker.sh  # Cyclomatic complexity analysis
    └── dependency-graph.sh    # Dependency visualization
```

## Benefits of Modular Architecture

### 1. **Separation of Concerns**
- Each module handles one specific aspect of analysis
- Easy to understand and maintain
- Clear boundaries between functionality

### 2. **Reusability**
- Library modules can be used by other scripts
- Common utilities shared across all modules
- No code duplication

### 3. **Maintainability**
- Smaller, focused files (30-104 lines vs 483 lines)
- Easier to debug and test individual components
- Changes isolated to specific modules

### 4. **Extensibility**
- Easy to add new analysis types
- Simple to modify existing analysis logic
- Plug-and-play architecture

### 5. **Testing**
- Individual modules can be tested in isolation
- Easier to write unit tests for specific functionality
- Better error isolation

## Usage

### Basic Analysis
```bash
# Use the new modular script
./scripts/analyze-package-new.sh internal/ecs/systems/collision

# Use the original monolithic script
./scripts/analyze-package.sh internal/ecs/systems/collision
```

### Detailed Analysis
```bash
# Interface and method analysis
./scripts/analyze-package-new.sh internal/ecs --interfaces --methods

# Full analysis with all options
./scripts/analyze-package-new.sh internal/game -f -s -i -I -m -d
```

## Module Descriptions

### `common.sh`
- **Purpose**: Shared utilities and helper functions
- **Functions**: Colors, printing, file operations, validation
- **Size**: 104 lines
- **Dependencies**: None

### `package-info.sh`
- **Purpose**: Basic package information and file overview
- **Functions**: Package declaration, file listing, consistency checks
- **Size**: 63 lines
- **Dependencies**: `common.sh`

### `import-analysis.sh`
- **Purpose**: Detailed import analysis
- **Functions**: Import categorization, dependency counting
- **Size**: 39 lines
- **Dependencies**: `common.sh`

### `struct-analysis.sh`
- **Purpose**: Struct definition analysis
- **Functions**: Struct discovery, field counting
- **Size**: 34 lines
- **Dependencies**: `common.sh`

### `interface-analysis.sh`
- **Purpose**: Interface analysis and usage
- **Functions**: Interface discovery, method signatures, implementation analysis
- **Size**: 58 lines
- **Dependencies**: `common.sh`

### `method-analysis.sh`
- **Purpose**: Function and method analysis
- **Functions**: Complexity estimation, signature extraction
- **Size**: 61 lines
- **Dependencies**: `common.sh`

### `dependency-analysis.sh`
- **Purpose**: Package dependency analysis
- **Functions**: Internal/external dependency tracking
- **Size**: 30 lines
- **Dependencies**: `common.sh`

### `metrics.sh`
- **Purpose**: Code quality metrics and recommendations
- **Functions**: Statistical analysis, best practice recommendations
- **Size**: 85 lines
- **Dependencies**: `common.sh`

## Migration Path

The original `analyze-package.sh` script is preserved for backward compatibility. The new modular approach provides:

1. **Same functionality** - All features preserved
2. **Better organization** - Clear separation of concerns
3. **Easier maintenance** - Smaller, focused modules
4. **Future extensibility** - Easy to add new analysis types

## Best Practices

### Adding New Analysis Types
1. Create a new module in `lib/`
2. Source `common.sh` for shared utilities
3. Implement focused analysis functions
4. Add to the main orchestrator script
5. Update this README

### Module Design Principles
- **Single Responsibility**: Each module does one thing well
- **Dependency Injection**: Pass data as parameters, don't rely on globals
- **Error Handling**: Use the shared error handling from `common.sh`
- **Documentation**: Include clear comments and function descriptions

### Testing
- Test individual modules in isolation
- Verify integration with the main orchestrator
- Ensure backward compatibility with existing usage patterns 