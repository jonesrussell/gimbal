#!/bin/bash

# scripts/lib/package-info.sh
# Basic package information analysis

LIB_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$LIB_DIR/common.sh"

# Analyze basic package information
analyze_package_info() {
    local pkg_path="$1"
    local module_name
    module_name=$(get_module_info)
    
    print_subheader "📦 Package Information"
    echo "Package Path: $pkg_path"
    echo "Import Path: $module_name/$pkg_path"
    echo "Analysis Date: $(date)"
}

# Analyze package declaration and consistency
analyze_package_declaration() {
    local go_files="$1"
    
    print_subheader "📋 Package Declaration"
    local package_name
    package_name=$(head -10 $(echo "$go_files" | head -1) | grep "^package " | awk '{print $2}' | head -1)
    echo "Package Name: $package_name"
    
    # Check for consistent package names
    local inconsistent_packages
    inconsistent_packages=$(echo "$go_files" | xargs grep "^package " | awk '{print $2}' | sort | uniq -c | sort -nr)
    if [ $(echo "$inconsistent_packages" | wc -l) -gt 1 ]; then
        print_warning "Inconsistent package names found:"
        echo "$inconsistent_packages"
    else
        print_info "Package naming is consistent"
    fi
}

# Analyze files overview
analyze_files_overview() {
    local go_files="$1"
    local file_count
    file_count=$(echo "$go_files" | wc -l)
    
    if [ -z "$go_files" ] || [ "$file_count" -eq 0 ]; then
        print_error "No Go files found"
        return 1
    fi
    
    echo "Go Files: $file_count"
    echo ""
    
    print_subheader "📄 Files Overview"
    echo "$go_files" | while read -r file; do
        if [ -f "$file" ]; then
            local lines
            lines=$(count_lines "$file")
            filename=$(basename "$file")
            printf "%s%*s %6d lines\n" "$filename" $((40 - ${#filename})) "" "$lines"
        fi
    done
} 