#!/bin/bash

# scripts/analyze-package-new.sh
# Modular Go package analysis using library components

set -e

# Source all library modules
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/lib/common.sh"
source "$SCRIPT_DIR/lib/package-info.sh"
source "$SCRIPT_DIR/lib/import-analysis.sh"
source "$SCRIPT_DIR/lib/struct-analysis.sh"
source "$SCRIPT_DIR/lib/interface-analysis.sh"
source "$SCRIPT_DIR/lib/method-analysis.sh"
source "$SCRIPT_DIR/lib/dependency-analysis.sh"
source "$SCRIPT_DIR/lib/metrics.sh"

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
    echo "  -I, --interfaces Show detailed interface analysis"
    echo "  -m, --methods   Show method signatures and complexity"
    echo "  -d, --deps      Show dependency graph"
    echo "  -o, --output    Output to file instead of stdout"
    echo "  -h, --help      Show this help"
    echo ""
    echo "Examples:"
    echo "  $0 internal/ecs/systems/collision"
    echo "  $0 internal/input -f -s -m"
    echo "  $0 internal/game --full --methods --output game_analysis.txt"
    echo "  $0 internal/ecs --interfaces --methods"
}

# Default options
PACKAGE_PATH=""
SHOW_FULL=false
SHOW_STRUCTS=false
SHOW_IMPORTS=false
SHOW_INTERFACES=false
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
        -I|--interfaces)
            SHOW_INTERFACES=true
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

# Validate package path
if ! validate_package_path "$PACKAGE_PATH"; then
    usage
    exit 1
fi

# Setup output redirection
if [ -n "$OUTPUT_FILE" ]; then
    exec > "$OUTPUT_FILE"
fi

# Main analysis function
analyze_package() {
    local pkg_path="$1"
    local go_files
    go_files=$(find_go_files "$pkg_path")
    
    print_header "🔍 Go Package Analysis: $pkg_path"
    
    # Basic package information
    analyze_package_info "$pkg_path"
    analyze_files_overview "$go_files"
    analyze_package_declaration "$go_files"
    
    # Optional detailed analyses
    if [ "$SHOW_IMPORTS" = true ]; then
        analyze_imports "$go_files"
    fi
    
    if [ "$SHOW_STRUCTS" = true ]; then
        analyze_structs "$go_files"
    fi
    
    if [ "$SHOW_INTERFACES" = true ]; then
        analyze_interfaces "$go_files"
    fi
    
    if [ "$SHOW_METHODS" = true ]; then
        analyze_methods "$go_files"
    fi
    
    if [ "$SHOW_DEPS" = true ]; then
        analyze_dependencies "$go_files" "$pkg_path"
    fi
    
    # Always show metrics and recommendations
    analyze_metrics "$go_files"
    generate_recommendations "$go_files"
    
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
}

# Run the analysis
analyze_package "$PACKAGE_PATH" 