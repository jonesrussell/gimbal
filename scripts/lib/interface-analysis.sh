#!/bin/bash

# scripts/lib/interface-analysis.sh
# Interface analysis functionality

source "$(dirname "$0")/common.sh"

# Analyze interfaces in detail
analyze_interfaces() {
    local go_files="$1"
    
    print_subheader "🔌 Interface Analysis"
    
    echo "Interface definitions:"
    echo "$go_files" | xargs grep -n "^type.*interface" | while read -r line; do
        local file_line
        file_line=$(echo "$line" | cut -d: -f1,2)
        local interface_def
        interface_def=$(echo "$line" | cut -d: -f3-)
        local interface_name
        interface_name=$(echo "$interface_def" | awk '{print $2}')
        
        echo "  $interface_name ($file_line)"
        
        # Count methods (rough estimate)
        local file
        file=$(echo "$line" | cut -d: -f1)
        local line_num
        line_num=$(echo "$line" | cut -d: -f2)
        local method_count
        method_count=$(sed -n "${line_num},/^}/p" "$file" | grep -c "^\s*[A-Z].*(" || echo "0")
        echo "    Methods: ~$method_count"
        
        # Show method signatures for interfaces
        echo "    Method signatures:"
        sed -n "${line_num},/^}/p" "$file" | grep "^\s*[A-Z].*(" | while read -r method_line; do
            local clean_method
            clean_method=$(echo "$method_line" | sed 's/^\s*//' | sed 's/\s*$//')
            echo "      $clean_method"
        done
    done
    
    echo ""
    echo "Interface usage analysis:"
    echo "Types implementing interfaces:"
    echo "$go_files" | xargs grep -n "func (.*) " | grep -E "\([^)]*\) [A-Z]" | while read -r line; do
        local file_line
        file_line=$(echo "$line" | cut -d: -f1,2)
        local method_def
        method_def=$(echo "$line" | cut -d: -f3-)
        local receiver
        receiver=$(extract_receiver "$method_def")
        local method_name
        method_name=$(echo "$method_def" | sed 's/^func ([^)]*) \([^(]*\).*/\1/')
        
        echo "  $method_name ($receiver) ($file_line)"
    done
} 