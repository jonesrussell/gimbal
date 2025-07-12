#!/bin/bash

# scripts/lib/method-analysis.sh
# Method analysis functionality

LIB_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$LIB_DIR/common.sh"

# Analyze methods and functions
analyze_methods() {
    local go_files="$1"
    
    print_subheader "🔧 Method Analysis"
    
    echo "Function signatures:"
    echo "$go_files" | xargs grep -n "^func " | while read -r line; do
        local file_line
        file_line=$(echo "$line" | cut -d: -f1,2 | sed 's|.*/||')
        local func_sig
        func_sig=$(echo "$line" | cut -d: -f3- | sed 's/^func //')
        local func_name
        func_name=$(extract_function_name "$func_sig")
        
        echo "  $func_name ($file_line)"
        
        # Estimate complexity by counting branches
        local file
        file=$(echo "$line" | cut -d: -f1)
        local start_line
        start_line=$(echo "$line" | cut -d: -f2)
        
        # Find function end (rough estimate)
        local complexity
        complexity=$(sed -n "${start_line},/^}/p" "$file" | grep -c "if\|for\|switch\|case\|range" 2>/dev/null || echo "0")
        # Ensure complexity is a number
        if ! [[ "$complexity" =~ ^[0-9]+$ ]]; then
            complexity=0
        fi
        if safe_compare "$complexity" -gt 7; then
            print_warning "    High complexity: $complexity"
        elif safe_compare "$complexity" -gt 4; then
            echo "    Moderate complexity: $complexity"
        else
            echo "    Low complexity: $complexity"
        fi
    done
    
    echo ""
    echo "Method signatures (with receivers):"
    echo "$go_files" | xargs grep -n "^func (.*)" | while read -r line; do
        local file_line
        file_line=$(echo "$line" | cut -d: -f1,2 | sed 's|.*/||')
        local method_sig
        method_sig=$(echo "$line" | cut -d: -f3- | sed 's/^func //')
        local receiver
        receiver=$(extract_receiver "$method_sig")
        local method_name
        method_name=$(echo "$method_sig" | sed 's/^([^)]*) \([^(]*\).*/\1/')
        
        echo "  $method_name ($receiver) ($file_line)"
    done
} 