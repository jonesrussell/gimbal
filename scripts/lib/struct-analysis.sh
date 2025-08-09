#!/bin/bash

# scripts/lib/struct-analysis.sh
# Struct analysis functionality

LIB_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$LIB_DIR/common.sh"

# Analyze structs in detail
analyze_structs() {
    local go_files="$1"
    
    print_subheader "🏗️ Struct Analysis"
    
    echo "Struct definitions:"
    echo "$go_files" | xargs grep -n "^type.*struct" | while read -r line; do
        local file_line
        file_line=$(echo "$line" | cut -d: -f1,2)
        local struct_def
        struct_def=$(echo "$line" | cut -d: -f3-)
        local struct_name
        struct_name=$(echo "$struct_def" | awk '{print $2}')
        
        echo "  $struct_name ($file_line)"
        
        # Count fields (rough estimate)
        local file
        file=$(echo "$line" | cut -d: -f1)
        local line_num
        line_num=$(echo "$line" | cut -d: -f2)
        local field_count
        field_count=$(sed -n "${line_num},/^}/p" "$file" | grep -c "^\s*[A-Z]" || echo "0")
        echo "    Fields: ~$field_count"
    done
} 