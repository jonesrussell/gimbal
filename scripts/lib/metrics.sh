#!/bin/bash

# scripts/lib/metrics.sh
# Code quality metrics and recommendations

LIB_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$LIB_DIR/common.sh"

# Calculate and display code quality metrics
analyze_metrics() {
    local go_files="$1"
    
    print_subheader "📊 Code Quality Metrics"
    
    local total_lines
    total_lines=$(echo "$go_files" | xargs wc -l | tail -1 | awk '{print $1}')
    local total_funcs
    total_funcs=$(echo "$go_files" | xargs grep -c "^func " | awk -F: '{sum+=$2} END {print sum}' || echo "0")
    local total_structs
    total_structs=$(echo "$go_files" | xargs grep -c "^type.*struct" | awk -F: '{sum+=$2} END {print sum}' || echo "0")
    local total_interfaces
    total_interfaces=$(echo "$go_files" | xargs grep -c "^type.*interface" | awk -F: '{sum+=$2} END {print sum}' || echo "0")
    local exported_funcs
    exported_funcs=$(echo "$go_files" | xargs grep -c "^func [A-Z]" | awk -F: '{sum+=$2} END {print sum}' || echo "0")
    
    echo "Total Lines: $total_lines"
    echo "Total Functions: $total_funcs"
    echo "Total Structs: $total_structs"
    echo "Total Interfaces: $total_interfaces"
    echo "Exported Functions: $exported_funcs"
    
    if safe_compare "$total_funcs" -gt 0; then
        local avg_lines_per_func
        avg_lines_per_func=$((total_lines / total_funcs))
        echo "Average Lines per Function: $avg_lines_per_func"
    fi
    
    if safe_compare "$total_structs" -gt 0 && safe_compare "$total_interfaces" -gt 0; then
        local interface_ratio
        interface_ratio=$(echo "scale=2; $total_interfaces / $total_structs" | bc -l 2>/dev/null || echo "0")
        echo "Interface/Struct Ratio: $interface_ratio"
    fi
}

# Generate recommendations based on metrics
generate_recommendations() {
    local go_files="$1"
    
    print_subheader "💡 Recommendations"
    
    local total_lines
    total_lines=$(echo "$go_files" | xargs wc -l | tail -1 | awk '{print $1}')
    local total_funcs
    total_funcs=$(echo "$go_files" | xargs grep -c "^func " | awk -F: '{sum+=$2} END {print sum}' || echo "0")
    local exported_funcs
    exported_funcs=$(echo "$go_files" | xargs grep -c "^func [A-Z]" | awk -F: '{sum+=$2} END {print sum}' || echo "0")
    local total_structs
    total_structs=$(echo "$go_files" | xargs grep -c "^type.*struct" | awk -F: '{sum+=$2} END {print sum}' || echo "0")
    local total_interfaces
    total_interfaces=$(echo "$go_files" | xargs grep -c "^type.*interface" | awk -F: '{sum+=$2} END {print sum}' || echo "0")
    
    if safe_compare "$total_lines" -gt 500; then
        print_warning "Package is large ($total_lines lines) - consider splitting"
    fi
    
    if safe_compare "$total_funcs" -gt 0 && safe_compare "$exported_funcs" -gt 0; then
        local export_ratio
        export_ratio=$(echo "scale=2; $exported_funcs / $total_funcs" | bc -l 2>/dev/null || echo "0")
        if (( $(echo "$export_ratio > 0.7" | bc -l 2>/dev/null || echo "0") )); then
            print_warning "High ratio of exported functions ($export_ratio) - consider internal functions"
        fi
    fi
    
    if safe_compare "$total_interfaces" -eq 0 && safe_compare "$total_structs" -gt 0; then
        print_warning "No interfaces found - consider adding abstractions"
    fi
    
    # Check for missing documentation
    local undocumented_exports
    undocumented_exports=$(echo "$go_files" | xargs grep -B1 "^func [A-Z]\|^type [A-Z].*struct\|^type [A-Z].*interface" | grep -v "^--$" | grep -v "^//" | grep "^func [A-Z]\|^type [A-Z]" | wc -l || echo "0")
    if safe_compare "$undocumented_exports" -gt 0; then
        print_warning "$undocumented_exports exported symbols lack documentation"
    fi
    
    print_info "Analysis complete!"
} 