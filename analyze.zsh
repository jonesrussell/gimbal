#!/bin/zsh

# Gimbal Context Integration Analysis Script
# This script gathers all necessary information for context integration planning

echo "🚀 GIMBAL CONTEXT INTEGRATION ANALYSIS"
echo "======================================"

# Create analysis output directory
mkdir -p analysis_output
cd analysis_output

echo "\n📁 RESOURCE MANAGER ANALYSIS"
echo "----------------------------"

# 1. Current ResourceManager interface and implementation
echo "=== Resource Manager Files ==="
find ../internal/ecs/managers/resource -name "*.go" -exec echo "📄 {}" \; -exec cat {} \; -exec echo "\n" \;

echo "\n=== Resource Manager Interface Definitions ==="
# Look for ResourceManager interface definitions
grep -r "ResourceManager" ../internal --include="*.go" -n

echo "\n📊 LOADSPRITE USAGE ANALYSIS"
echo "----------------------------"

# 2. Find all LoadSprite calls
echo "=== All LoadSprite() Calls ==="
grep -r "LoadSprite" ../internal --include="*.go" -n -A 2 -B 2

echo "\n=== All LoadSprite Method Definitions ==="
grep -r "func.*LoadSprite" ../internal --include="*.go" -n -A 5

echo "\n🔍 RELATED LOADING METHODS"
echo "-------------------------"

# 3. Find other resource loading methods that might need context
echo "=== Other Load Methods ==="
grep -r "func.*Load" ../internal/ecs/managers/resource --include="*.go" -n

echo "\n=== Font Loading ==="
grep -r "LoadFont\|Font.*Load" ../internal --include="*.go" -n

echo "\n=== Audio Loading ==="
grep -r "LoadAudio\|Audio.*Load" ../internal --include="*.go" -n

echo "\n🏗️ INTERFACE DEFINITIONS"
echo "------------------------"

# 4. Find interface definitions that might need updating
echo "=== Common Interfaces ==="
cat ../internal/common/interfaces.go

echo "\n=== ECS Contracts ==="
cat ../internal/ecs/contracts/contracts.go

echo "\n📞 CALL SITE ANALYSIS"
echo "--------------------"

# 5. Find all places where resource manager is used
echo "=== Resource Manager Usage ==="
grep -r "resourceManager\|ResourceManager" ../internal --include="*.go" -n -A 2 -B 2

echo "\n=== Container Usage (Dependency Injection) ==="
grep -r "container\." ../internal --include="*.go" -n | grep -i resource

echo "\n🎯 INITIALIZATION PATTERNS"
echo "-------------------------"

# 6. Find initialization patterns
echo "=== Game Initialization ==="
cat ../internal/game/game_init.go

echo "\n=== Container Setup ==="
cat ../internal/app/container.go

echo "\n📦 CONTEXT USAGE PATTERNS"
echo "------------------------"

# 7. Find existing context usage
echo "=== Current Context Usage ==="
grep -r "context\." ../internal --include="*.go" -n -A 1 -B 1

echo "\n=== Context Import Statements ==="
grep -r "\"context\"" ../internal --include="*.go" -n

echo "\n🔧 METHOD SIGNATURES"
echo "-------------------"

# 8. Get all method signatures that might need context
echo "=== All Manager Method Signatures ==="
grep -r "func.*Manager" ../internal/ecs/managers --include="*.go" -n -A 1

echo "\n=== System Constructor Signatures ==="
grep -r "func New.*System" ../internal/ecs/systems --include="*.go" -n -A 3

echo "\n✅ ANALYSIS COMPLETE"
echo "==================="
echo "All analysis files saved to: $(pwd)"
echo "Review the output above to understand:"
echo "  • Current ResourceManager interface"
echo "  • All LoadSprite() call sites"
echo "  • Related loading methods"
echo "  • Interface definitions needing updates"
echo "  • Initialization patterns"
echo "  • Existing context usage patterns"

# Create a summary file
cat > context_integration_summary.md << 'EOF'
# Context Integration Analysis Summary

## Files to Update

### 1. Interface Definitions
- [ ] internal/common/interfaces.go (if ResourceManager is defined here)
- [ ] internal/ecs/contracts/contracts.go (if ResourceManager is defined here)
- [ ] internal/ecs/managers/resource/manager.go (main implementation)

### 2. Resource Loading Methods
- [ ] LoadSprite() method signature
- [ ] LoadFont() method signature (if exists)
- [ ] LoadAudio() method signature (if exists)
- [ ] Any other Load*() methods

### 3. Call Sites (to be identified from grep results)
- [ ] Game initialization
- [ ] Scene setup
- [ ] System constructors
- [ ] Any other resource loading calls

### 4. Import Statements
- [ ] Add "context" import to files that need it

## Implementation Strategy

1. **Update Interface First**: Modify the ResourceManager interface
2. **Update Implementation**: Add context parameter to methods
3. **Update Call Sites**: Systematically update all callers
4. **Add Context Checks**: Add cancellation checks where appropriate
5. **Test**: Verify clean build after each step

## Context Usage Patterns

- Use `context.Background()` for initialization
- Use `context.WithTimeout()` for long-running operations
- Use `context.WithCancel()` for cancellable operations
- Check `ctx.Done()` in long-running loops
EOF

echo "\n📋 Summary file created: context_integration_summary.md"
echo "🎯 Ready for context integration implementation!"