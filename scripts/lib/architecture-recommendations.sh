#!/bin/bash

# scripts/lib/architecture-recommendations.sh
# Architecture recommendations and summary

LIB_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$LIB_DIR/common.sh"

# Generate architecture recommendations
generate_architecture_recommendations() {
    print_subheader "💡 ARCHITECTURE RECOMMENDATIONS"
    
    # Analyze overall project health
    local total_files
    total_files=$(find . -name "*.go" -not -path "./vendor/*" | wc -l)
    local large_files
    large_files=$(find . -name "*.go" -not -path "./vendor/*" -exec wc -l {} + | awk '$1 > 500 {count++} END {print count+0}')
    local very_large_files
    very_large_files=$(find . -name "*.go" -not -path "./vendor/*" -exec wc -l {} + | awk '$1 > 1000 {count++} END {print count+0}')
    local init_functions
    init_functions=$(grep -r "func init()" . --include="*.go" --exclude-dir=vendor | wc -l)
    local panic_calls
    panic_calls=$(grep -r "panic(" . --include="*.go" --exclude-dir=vendor | wc -l)
    
    echo "📊 Project Health Summary:"
    echo "  Total Go files: $total_files"
    echo "  Large files (>500 lines): $large_files"
    echo "  Very large files (>1000 lines): $very_large_files"
    echo "  Init functions: $init_functions"
    echo "  Panic calls: $panic_calls"
    echo ""
    
    # Generate specific recommendations
    if safe_compare "$large_files" -gt 5; then
        print_warning "Many large files detected - prioritize refactoring"
        echo "  → Break down files >500 lines into smaller, focused modules"
        echo "  → Apply single responsibility principle"
    fi
    
    if safe_compare "$very_large_files" -gt 0; then
        print_error "Very large files detected - high priority refactoring needed"
        echo "  → Files >1000 lines violate maintainability best practices"
        echo "  → Consider splitting into multiple packages"
    fi
    
    if safe_compare "$init_functions" -gt 5; then
        print_warning "Many init functions detected"
        echo "  → Consider dependency injection patterns"
        echo "  → Reduce global state initialization"
    fi
    
    if safe_compare "$panic_calls" -gt 0; then
        print_warning "Panic calls detected"
        echo "  → Replace panics with proper error handling"
        echo "  → Use error wrapping and context"
    fi
    
    echo ""
    echo "🎯 Priority Actions:"
    echo "1. Review largest files first (highest impact)"
    echo "2. Address panic calls (reliability)"
    echo "3. Reduce init functions (testability)"
    echo "4. Add missing interfaces (abstraction)"
    echo "5. Improve error handling patterns"
    echo ""
}

# Generate final summary
generate_analysis_summary() {
    print_subheader "✅ ANALYSIS COMPLETE"
    echo "===================="
    echo "Review the output above to identify:"
    echo "- Large files that need refactoring"
    echo "- Complex functions that violate SRP"
    echo "- Missing interfaces (low interface count vs struct count)"
    echo "- Package organization issues"
    echo "- Global state management problems"
    echo ""
    echo "Next steps:"
    echo "1. Install missing static analysis tools"
    echo "2. Review largest/most complex files first"
    echo "3. Look for common anti-patterns"
    echo "4. Check for proper separation of concerns"
    echo "5. Consider architectural patterns (Clean Architecture, Hexagonal)"
    echo ""
    print_info "Use './scripts/analyze-package.sh <package>' for detailed package analysis"
} 