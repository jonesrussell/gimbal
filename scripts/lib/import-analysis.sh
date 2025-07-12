#!/bin/bash

# scripts/lib/import-analysis.sh
# Import analysis functionality

LIB_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$LIB_DIR/common.sh"

# Analyze imports in detail
analyze_imports() {
    local go_files="$1"
    local module_name
    module_name=$(get_module_info)
    
    print_subheader "📥 Import Analysis"
    
    echo "All imports in package:"
    echo "$go_files" | xargs grep "^import\|^\s*\"" | \
        grep -v "^import (" | \
        sed 's/.*"\(.*\)".*/\1/' | \
        sort | uniq -c | sort -nr | head -20
    
    echo ""
    echo "Standard library imports:"
    echo "$go_files" | xargs grep "^\s*\"[^/]*\"" | \
        sed 's/.*"\(.*\)".*/\1/' | \
        sort | uniq -c | sort -nr
    
    echo ""
    echo "External dependencies:"
    echo "$go_files" | xargs grep "^\s*\".*\..*/" | \
        sed 's/.*"\(.*\)".*/\1/' | \
        sort | uniq -c | sort -nr
    
    echo ""
    echo "Internal imports:"
    echo "$go_files" | xargs grep "^\s*\"$module_name/" | \
        sed 's/.*"\(.*\)".*/\1/' | \
        sort | uniq -c | sort -nr
} 