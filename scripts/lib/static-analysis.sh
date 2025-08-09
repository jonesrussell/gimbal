#!/bin/bash

# scripts/lib/static-analysis.sh
# Static analysis functionality

LIB_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$LIB_DIR/common.sh"

# Analyze with static analysis tools
analyze_static_analysis() {
    print_subheader "🔧 STATIC ANALYSIS"
    
    echo "Checking for static analysis tools..."
    
    # Cyclomatic complexity analysis
    if command -v gocyclo &> /dev/null; then
        echo "Cyclomatic complexity (top 10 most complex functions):"
        gocyclo -top 10 .
        echo ""
    else
        print_info "Install gocyclo for complexity analysis: go install github.com/fzipp/gocyclo/cmd/gocyclo@latest"
    fi
    
    # Linting analysis
    if command -v golint &> /dev/null; then
        echo "Lint issues (first 20):"
        golint ./... | head -20
        echo ""
    else
        print_info "Install golint: go install golang.org/x/lint/golint@latest"
    fi
    
    # Static analysis
    if command -v staticcheck &> /dev/null; then
        echo "Static analysis issues (first 20):"
        staticcheck ./... | head -20
        echo ""
    else
        print_info "Install staticcheck: go install honnef.co/go/tools/cmd/staticcheck@latest"
    fi
    
    # Go vet
    echo "Go vet analysis:"
    go vet ./... 2>&1 | head -10
    echo ""
}

# Analyze potential issues
analyze_potential_issues() {
    print_subheader "🚨 POTENTIAL ISSUES"
    
    # Large files
    echo "Files with >500 lines (potential refactoring candidates):"
    find . -name "*.go" -not -path "./vendor/*" -exec wc -l {} + | awk '$1 > 500 {print $2 " (" $1 " lines)"}' | sort
    
    echo ""
    echo "Files with >1000 lines (high priority refactoring candidates):"
    find . -name "*.go" -not -path "./vendor/*" -exec wc -l {} + | awk '$1 > 1000 {print $2 " (" $1 " lines)"}' | sort
    
    echo ""
    echo "Long function names (>30 chars, might indicate SRP violations):"
    grep -r "^func " . --include="*.go" --exclude-dir=vendor | grep -E "func [^(]{30,}\(" | head -10
    
    echo ""
    echo "Global variables (potential state management issues):"
    grep -r "^var " . --include="*.go" --exclude-dir=vendor | head -10
    
    echo ""
    echo "Init functions (initialization complexity):"
    local init_count
    init_count=$(grep -r "func init()" . --include="*.go" --exclude-dir=vendor | wc -l)
    echo "Total init functions: $init_count"
    if safe_compare "$init_count" -gt 5; then
        print_warning "Many init functions detected - consider dependency injection"
    fi
    grep -r "func init()" . --include="*.go" --exclude-dir=vendor
    
    echo ""
    echo "Panic usage (error handling issues):"
    local panic_count
    panic_count=$(grep -r "panic(" . --include="*.go" --exclude-dir=vendor | wc -l)
    echo "Total panic calls: $panic_count"
    if safe_compare "$panic_count" -gt 0; then
        print_warning "Panic calls detected - consider proper error handling"
        grep -r "panic(" . --include="*.go" --exclude-dir=vendor | head -5
    fi
} 