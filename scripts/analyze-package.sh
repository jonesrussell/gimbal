#!/bin/bash

# scripts/analyze-package.sh
# Deep analysis of a specific Go package

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Usage function
usage() {
    echo "Usage: $0 <package_path> [options]"
    echo ""
    echo "Analyze a Go package in detail"
    echo ""
    echo "Arguments:"
    echo "  package_path    Path to package (e.g., internal/ecs/systems/collision)"
    echo ""
    echo "Options:"
    echo "  -f, --full      Show full file contents (default: summary only)"
    echo "  -s, --structs   Show detailed struct analysis"
    echo "  -i, --imports   Show detailed import analysis"
    echo "  -m, --methods   Show method signatures and complexity"
    echo "  -d, --deps      Show dependency graph"
    echo "  -o, --output    Output to file instead of stdout"
    echo "  -h, --help      Show this help"
    echo ""
    echo "Examples:"
    echo "  $0 internal/ecs/systems/collision"
    echo "  $0 internal/input -f -s -m"
    echo "  $0 internal/game --full --methods --output game_analysis.txt"
}

# Default options
PACKAGE_PATH=""
SHOW_FULL=false
SHOW_STRUCTS=false
SHOW_IMPORTS=false
SHOW_METHODS=false
SHOW_DEPS=false
OUTPUT_FILE=""

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -f|--full)
            SHOW_FULL=true
            shift
            ;;
        -s|--structs)
            SHOW_STRUCTS=true
            shift
            ;;
        -i|--imports)
            SHOW_IMPORTS=true
            shift
            ;;
        -m|--methods)
            SHOW_METHODS=true
            shift
            ;;
        -d|--deps)
            SHOW_DEPS=true
            shift
            ;;
        -o|--output)
            OUTPUT_FILE="$2"
            shift 2
            ;;
        -h|--help)
            usage
            exit 0
            ;;
        -*|--*)
            echo "Unknown option $1"
            usage
            exit 1
            ;;
        *)
            if [ -z "$PACKAGE_PATH" ]; then
                PACKAGE_PATH="$1"
            else
                echo "Multiple package paths not supported"
                usage
                exit 1
            fi
            shift
            ;;
    esac
done

# Check if package path is provided
if [ -z "$PACKAGE_PATH" ]; then
    echo "Error: Package path is required"
    usage
    exit 1
fi

# Check if package exists
if [ ! -d "$PACKAGE_PATH" ]; then
    echo "Error: Package directory '$PACKAGE_PATH' does not exist"
    exit 1
fi

# Setup output redirection
if [ -n "$OUTPUT_FILE" ]; then
    exec > "$OUTPUT_FILE"
fi

# Helper functions
print_header() {
    echo -e "${BLUE}$1${NC}"
    printf '%*s\n' ${#1} | tr ' ' '='
}

print_subheader() {
    echo -e "\n${CYAN}$1${NC}"
    printf '%*s\n' ${#1} | tr ' ' '-'
}

print_info() {
    echo -e "${GREEN}ℹ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

# Get go module info
get_module_info() {
    if [ -f "go.mod" ]; then
        MODULE_NAME=$(head -1 go.mod | awk '{print $2}')
        echo "$MODULE_NAME"
    else
        echo "unknown"
    fi
}

# Main analysis function
analyze_package() {
    local pkg_path="$1"
    local module_name
    module_name=$(get_module_info)
    
    print_header "🔍 Go Package Analysis: $pkg_path"
    
    # Basic package info
    print_subheader "📦 Package Information"
    echo "Package Path: $pkg_path"
    echo "Import Path: $module_name/$pkg_path"
    echo "Analysis Date: $(date)"
    
    # Find Go files
    local go_files
    go_files=$(find "$pkg_path" -name "*.go" -type f | sort)
    local file_count
    file_count=$(echo "$go_files" | wc -l)
    
    if [ -z "$go_files" ] || [ "$file_count" -eq 0 ]; then
        print_error "No Go files found in $pkg_path"
        return 1
    fi
    
    echo "Go Files: $file_count"
    echo ""
    
    # List files with line counts
    print_subheader "📄 Files Overview"
    echo "$go_files" | while read -r file; do
        if [ -f "$file" ]; then
            local lines
            lines=$(wc -l < "$file")
            filename=$(basename "$file")
            printf "%s%*s %6d lines\n" "$filename" $((40 - ${#filename})) "" "$lines"
        fi
    done
    
    # Package declaration analysis
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
    
    # Import analysis
    if [ "$SHOW_IMPORTS" = true ]; then
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
    fi
    
    # Struct analysis
    if [ "$SHOW_STRUCTS" = true ]; then
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
        
        echo ""
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
        done
    fi
    
    # Method analysis
    if [ "$SHOW_METHODS" = true ]; then
        print_subheader "🔧 Method Analysis"
        
        echo "Function signatures:"
        echo "$go_files" | xargs grep -n "^func " | while read -r line; do
            local file_line
            file_line=$(echo "$line" | cut -d: -f1,2 | sed 's|.*/||')
            local func_sig
            func_sig=$(echo "$line" | cut -d: -f3- | sed 's/^func //')
            local func_name
            func_name=$(echo "$func_sig" | awk '{print $1}' | sed 's/(.*//')
            
            echo "  $func_name ($file_line)"
            
            # Estimate complexity by counting branches
            local file
            file=$(echo "$line" | cut -d: -f1)
            local start_line
            start_line=$(echo "$line" | cut -d: -f2)
            
            # Find function end (rough estimate)
            local complexity
            complexity=$(sed -n "${start_line},/^}/p" "$file" | grep -c "if\|for\|switch\|case\|range" || echo "0")
            if [ "$complexity" -gt 7 ]; then
                print_warning "    High complexity: $complexity"
            elif [ "$complexity" -gt 4 ]; then
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
            receiver=$(echo "$method_sig" | sed 's/^(\([^)]*\)).*/\1/')
            local method_name
            method_name=$(echo "$method_sig" | sed 's/^([^)]*) \([^(]*\).*/\1/')
            
            echo "  $method_name ($receiver) ($file_line)"
        done
    fi
    
    # Dependency analysis
    if [ "$SHOW_DEPS" = true ]; then
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
    fi
    
    # Code quality metrics
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
    
    if [ "$total_funcs" -gt 0 ]; then
        local avg_lines_per_func
        avg_lines_per_func=$((total_lines / total_funcs))
        echo "Average Lines per Function: $avg_lines_per_func"
    fi
    
    if [ "$total_structs" -gt 0 ] && [ "$total_interfaces" -gt 0 ]; then
        local interface_ratio
        interface_ratio=$(echo "scale=2; $total_interfaces / $total_structs" | bc -l 2>/dev/null || echo "0")
        echo "Interface/Struct Ratio: $interface_ratio"
    fi
    
    # Show file contents if requested
    if [ "$SHOW_FULL" = true ]; then
        print_subheader "📝 Full File Contents"
        
        echo "$go_files" | while read -r file; do
            echo ""
            print_info "File: $file"
            echo "$(printf '─%.0s' $(seq 1 50))"
            cat "$file"
            echo ""
        done
    else
        print_subheader "📋 File Summaries"
        
        echo "$go_files" | while read -r file; do
            echo ""
            print_info "File: $file"
            echo "Package: $(grep "^package " "$file" | awk '{print $2}' | head -1)"
            echo "Imports: $(grep -c "^\s*\".*\"" "$file" || echo "0")"
            echo "Functions: $(grep -c "^func " "$file" || echo "0")"
            echo "Types: $(grep -c "^type " "$file" || echo "0")"
            echo "First 10 lines:"
            head -10 "$file" | sed 's/^/  /'
            if [ $(wc -l < "$file") -gt 20 ]; then
                echo "  ..."
                echo "Last 5 lines:"
                tail -5 "$file" | sed 's/^/  /'
            fi
        done
    fi
    
    # Recommendations
    print_subheader "💡 Recommendations"
    
    if [ "$total_lines" -gt 500 ]; then
        print_warning "Package is large ($total_lines lines) - consider splitting"
    fi
    
    if [ "$total_funcs" -gt 0 ] && [ "$exported_funcs" -gt 0 ]; then
        local export_ratio
        export_ratio=$(echo "scale=2; $exported_funcs / $total_funcs" | bc -l 2>/dev/null || echo "0")
        if (( $(echo "$export_ratio > 0.7" | bc -l 2>/dev/null || echo "0") )); then
            print_warning "High ratio of exported functions ($export_ratio) - consider internal functions"
        fi
    fi
    
    if [ "$total_interfaces" -eq 0 ] && [ "$total_structs" -gt 0 ]; then
        print_warning "No interfaces found - consider adding abstractions"
    fi
    
    # Check for missing documentation
    local undocumented_exports
    undocumented_exports=$(echo "$go_files" | xargs grep -B1 "^func [A-Z]\|^type [A-Z].*struct\|^type [A-Z].*interface" | grep -v "^--$" | grep -v "^//" | grep "^func [A-Z]\|^type [A-Z]" | wc -l || echo "0")
    if [ "$undocumented_exports" -gt 0 ]; then
        print_warning "$undocumented_exports exported symbols lack documentation"
    fi
    
    print_info "Analysis complete!"
}

# Run the analysis
analyze_package "$PACKAGE_PATH"
