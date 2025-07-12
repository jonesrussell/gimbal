#!/bin/bash

# scripts/lib/code-metrics.sh
# Code metrics analysis functionality

LIB_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$LIB_DIR/common.sh"

# Analyze code metrics
analyze_code_metrics() {
    print_subheader "📊 CODE METRICS"
    
    local go_files_count
    go_files_count=$(find . -name "*.go" -not -path "./vendor/*" | wc -l)
    echo "Go files count: $go_files_count"
    
    echo ""
    echo "Total lines of Go code:"
    find . -name "*.go" -not -path "./vendor/*" -exec wc -l {} + | tail -1
    
    echo ""
    echo "Lines per file (largest files first):"
    find . -name "*.go" -not -path "./vendor/*" -exec wc -l {} + | sort -nr | head -10
    
    echo ""
    echo "Function/method counts:"
    local total_funcs
    total_funcs=$(grep -r "^func " . --include="*.go" --exclude-dir=vendor | wc -l)
    local method_funcs
    method_funcs=$(grep -r "^func (" . --include="*.go" --exclude-dir=vendor | wc -l)
    echo "Total functions: $total_funcs"
    echo "Methods (receiver functions): $method_funcs"
    
    # Calculate ratios
    if safe_compare "$total_funcs" -gt 0; then
        local method_ratio
        method_ratio=$(echo "scale=2; $method_funcs / $total_funcs" | bc -l 2>/dev/null || echo "0")
        echo "Method/Function ratio: $method_ratio"
    fi
}

# Analyze pattern statistics
analyze_patterns() {
    print_subheader "🔍 PATTERN ANALYSIS"
    
    local struct_count
    struct_count=$(grep -r "^type.*struct" . --include="*.go" --exclude-dir=vendor | wc -l)
    echo "Struct definitions: $struct_count"
    
    echo ""
    local interface_count
    interface_count=$(grep -r "^type.*interface" . --include="*.go" --exclude-dir=vendor | wc -l)
    echo "Interface definitions: $interface_count"
    
    # Calculate interface/struct ratio
    if safe_compare "$struct_count" -gt 0 && safe_compare "$interface_count" -gt 0; then
        local interface_ratio
        interface_ratio=$(echo "scale=2; $interface_count / $struct_count" | bc -l 2>/dev/null || echo "0")
        echo "Interface/Struct ratio: $interface_ratio"
        
        if safe_compare "$interface_ratio" -lt 0.3; then
            print_warning "Low interface/struct ratio - consider adding more abstractions"
        fi
    fi
    
    echo ""
    echo "Package declarations:"
    grep -r "^package " . --include="*.go" --exclude-dir=vendor | cut -d: -f2 | sort | uniq -c | sort -nr
} 