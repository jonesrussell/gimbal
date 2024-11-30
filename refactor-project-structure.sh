#!/bin/bash

# Parse command line arguments
DRY_RUN=false
while [[ "$#" -gt 0 ]]; do
    case $1 in
        --dry-run|-n) DRY_RUN=true ;;
        *) echo "Unknown parameter: $1"; exit 1 ;;
    esac
    shift
done

# Function to execute or echo commands based on dry run flag
execute() {
    if [ "$DRY_RUN" = true ]; then
        echo "Would execute: $*"
    else
        eval "$@"
    fi
}

# Ensure we are in the project root directory
execute 'cd "$(git rev-parse --show-toplevel)" || exit'

# Move files to the new structure
execute 'git mv internal/types.go core/types.go'
execute 'git mv internal/game.go core/game.go'
execute 'git mv internal/input.go core/input.go'
execute 'git mv internal/render.go core/render.go'
execute 'git mv internal/assets assets'
execute 'git mv internal/systems systems'

# Create new directories for config and cmd (if needed)
execute 'mkdir -p config cmd'

# Move relevant files to config
execute 'git mv internal/config/* config/'

# Cleanup any now-empty directories
execute 'rmdir internal/config'
execute 'rmdir internal'

# Commit the changes
if [ "$DRY_RUN" = false ]; then
    execute 'git commit -m "Reorganize project structure to modern Go standards"'
else
    echo "Would commit changes with message: Reorganize project structure to modern Go standards"
fi
