#!/bin/zsh

# Navigate to the gimbal project root
cd ~/dev/gimbal

echo "=== GIMBAL INTERNAL/COMMON ARCHITECTURE ANALYSIS ==="
echo "Generated: $(date)"
echo "Branch: $(git branch --show-current)"
echo "Commit: $(git rev-parse --short HEAD)"
echo ""

# Function to display file with header
show_file() {
    local file="$1"
    if [[ -f "$file" ]]; then
        echo "=== $file ==="
        echo "Lines: $(wc -l < "$file")"
        echo "Size: $(wc -c < "$file") bytes"
        echo ""
        cat "$file"
        echo ""
        echo "=== END $file ==="
        echo ""
    else
        echo "=== $file === (FILE NOT FOUND)"
        echo ""
    fi
}

# Show all files in internal/common/
for file in internal/common/*.go; do
    show_file "$file"
done

echo "=== IMPORT ANALYSIS ==="
echo "Files that import from internal/common:"
grep -r "internal/common" --include="*.go" . | grep -v "internal/common/" | cut -d: -f1 | sort | uniq
echo ""

echo "=== ANALYSIS COMPLETE ==="
