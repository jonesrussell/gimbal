#!/bin/bash

# scripts/lib/dependency-analysis.sh
# Dependency analysis functionality

source "$(dirname "$0")/common.sh"

# Analyze dependencies
analyze_dependencies() {
    local go_files="$1"
    local pkg_path="$2"
    local module_name
    module_name=$(get_module_info)
    
    print_subheader "🔗 Dependency Analysis"
    
    echo "Dependencies on other internal packages:"
    echo "$go_files" | xargs grep "\"$module_name/" | \
        sed "s|.*\"$module_name/\([^\"]*\)\".*|\1|" | \
        grep -v "^$pkg_path" | \
        sort | uniq -c | sort -nr
    
    echo ""
    echo "Files that import this package:"
    if command -v rg >/dev/null 2>&1; then
        rg -l "\"$module_name/$pkg_path\"" --type go . 2>/dev/null || echo "None found"
    else
        grep -r "\"$module_name/$pkg_path\"" --include="*.go" . | cut -d: -f1 | sort | uniq || echo "None found"
    fi
} 