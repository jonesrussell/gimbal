#!/bin/bash

# scripts/analyze-architecture-new.sh
# Modular Go project architecture analysis

set -e

# Source all library modules
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/lib/common.sh"
source "$SCRIPT_DIR/lib/project-structure.sh"
source "$SCRIPT_DIR/lib/code-metrics.sh"
source "$SCRIPT_DIR/lib/static-analysis.sh"
source "$SCRIPT_DIR/lib/architecture-recommendations.sh"

# Usage function
usage() {
    echo "Usage: $0 [options]"
    echo ""
    echo "Analyze Go project architecture and identify improvement opportunities"
    echo ""
    echo "Options:"
    echo "  -s, --structure   Show detailed project structure analysis"
    echo "  -m, --metrics     Show detailed code metrics"
    echo "  -a, --analysis    Show static analysis results"
    echo "  -i, --issues      Show potential issues and anti-patterns"
    echo "  -r, --recommend   Show architecture recommendations"
    echo "  -o, --output      Output to file instead of stdout"
    echo "  -h, --help        Show this help"
    echo ""
    echo "Examples:"
    echo "  $0                    # Full analysis"
    echo "  $0 --structure --metrics"
    echo "  $0 --issues --recommend --output analysis.txt"
}

# Default options
SHOW_STRUCTURE=false
SHOW_METRICS=false
SHOW_ANALYSIS=false
SHOW_ISSUES=false
SHOW_RECOMMEND=false
OUTPUT_FILE=""

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -s|--structure)
            SHOW_STRUCTURE=true
            shift
            ;;
        -m|--metrics)
            SHOW_METRICS=true
            shift
            ;;
        -a|--analysis)
            SHOW_ANALYSIS=true
            shift
            ;;
        -i|--issues)
            SHOW_ISSUES=true
            shift
            ;;
        -r|--recommend)
            SHOW_RECOMMEND=true
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
            echo "Unexpected argument $1"
            usage
            exit 1
            ;;
    esac
done

# Setup output redirection
if [ -n "$OUTPUT_FILE" ]; then
    exec > "$OUTPUT_FILE"
fi

# Change to project root (assuming script is in scripts/ directory)
cd "$(dirname "$0")/.."

# Main analysis function
analyze_architecture() {
    print_header "🔍 Go Project Architecture Analysis"
    
    # Always show basic structure and dependencies
    analyze_project_structure
    analyze_dependencies
    analyze_entry_points
    analyze_build_tools
    
    # Optional detailed analyses
    if [ "$SHOW_METRICS" = true ] || [ "$SHOW_STRUCTURE" = false ] && [ "$SHOW_ANALYSIS" = false ] && [ "$SHOW_ISSUES" = false ] && [ "$SHOW_RECOMMEND" = false ]; then
        analyze_code_metrics
        analyze_patterns
    fi
    
    if [ "$SHOW_ANALYSIS" = true ] || [ "$SHOW_STRUCTURE" = false ] && [ "$SHOW_METRICS" = false ] && [ "$SHOW_ISSUES" = false ] && [ "$SHOW_RECOMMEND" = false ]; then
        analyze_static_analysis
    fi
    
    if [ "$SHOW_ISSUES" = true ] || [ "$SHOW_STRUCTURE" = false ] && [ "$SHOW_METRICS" = false ] && [ "$SHOW_ANALYSIS" = false ] && [ "$SHOW_RECOMMEND" = false ]; then
        analyze_potential_issues
    fi
    
    if [ "$SHOW_RECOMMEND" = true ] || [ "$SHOW_STRUCTURE" = false ] && [ "$SHOW_METRICS" = false ] && [ "$SHOW_ANALYSIS" = false ] && [ "$SHOW_ISSUES" = false ]; then
        generate_architecture_recommendations
    fi
    
    # Always show summary and recommendations if no specific options
    if [ "$SHOW_STRUCTURE" = false ] && [ "$SHOW_METRICS" = false ] && [ "$SHOW_ANALYSIS" = false ] && [ "$SHOW_ISSUES" = false ] && [ "$SHOW_RECOMMEND" = false ]; then
        # Full analysis mode
        analyze_code_metrics
        analyze_patterns
        analyze_static_analysis
        analyze_potential_issues
        generate_architecture_recommendations
    fi
    
    generate_analysis_summary
}

# Run the analysis
analyze_architecture 