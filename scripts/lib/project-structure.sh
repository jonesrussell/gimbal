#!/bin/bash

# scripts/lib/project-structure.sh
# Project structure analysis functionality

LIB_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$LIB_DIR/common.sh"

# Analyze project structure
analyze_project_structure() {
    print_subheader "📁 PROJECT STRUCTURE"
    echo "Project root: $(pwd)"
    echo ""
    
    # Check if tree is available, fallback to find
    if command -v tree &> /dev/null; then
        echo "Directory structure:"
        tree -I 'vendor|node_modules|.git|*.exe|*.so|*.dylib' -L 3
    else
        echo "Directory structure (using find):"
        find . -type d -name ".git" -prune -o -type d -name "vendor" -prune -o -type d -name "node_modules" -prune -o -type d -print | head -20
    fi
    
    echo ""
    echo "Go files structure:"
    find . -type f -name "*.go" | grep -v vendor | head -20
}

# Analyze dependencies and modules
analyze_dependencies() {
    print_subheader "📦 DEPENDENCIES & MODULES"
    if [ -f "go.mod" ]; then
        echo "Go module info:"
        head -10 go.mod
        echo ""
        echo "Dependency graph (top 20):"
        go mod graph | head -20
    else
        print_error "No go.mod found - not a Go module project"
    fi
}

# Analyze main entry points
analyze_entry_points() {
    print_subheader "🎯 MAIN ENTRY POINTS"
    echo "Main files found:"
    find . -name "main.go" -not -path "./vendor/*"
    
    echo ""
    echo "Main functions:"
    grep -r "func main()" . --include="*.go" --exclude-dir=vendor
}

# Analyze build tools and tasks
analyze_build_tools() {
    print_subheader "🛠️ BUILD TOOLS & TASKS"
    if command -v task &> /dev/null; then
        echo "Available tasks:"
        task --list
        echo ""
    else
        print_warning "Task not found - install from https://taskfile.dev"
    fi
} 