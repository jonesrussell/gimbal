#!/bin/bash

# scripts/lib/common.sh
# Shared utilities for analysis scripts

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Print functions
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

# Find Go files in a package
find_go_files() {
    local pkg_path="$1"
    find "$pkg_path" -name "*.go" -type f | sort
}

# Count lines in a file safely
count_lines() {
    local file="$1"
    if [ -f "$file" ]; then
        wc -l < "$file"
    else
        echo "0"
    fi
}

# Safe integer comparison
safe_compare() {
    local a="$1"
    local op="$2"
    local b="$3"
    
    if [[ "$a" =~ ^[0-9]+$ ]] && [[ "$b" =~ ^[0-9]+$ ]]; then
        [ "$a" "$op" "$b" ]
    else
        false
    fi
}

# Extract function name from signature
extract_function_name() {
    local sig="$1"
    echo "$sig" | sed 's/^func //' | awk '{print $1}' | sed 's/(.*$//'
}

# Extract receiver from method signature
extract_receiver() {
    local sig="$1"
    echo "$sig" | sed 's/^func (\([^)]*\)).*/\1/'
}

# Validate package path
validate_package_path() {
    local pkg_path="$1"
    if [ -z "$pkg_path" ]; then
        print_error "Package path is required"
        return 1
    fi
    
    if [ ! -d "$pkg_path" ]; then
        print_error "Package directory '$pkg_path' does not exist"
        return 1
    fi
    
    return 0
} 