#!/bin/bash

# scripts/fix-error-codes.sh
# Fix inconsistent error code naming

set -e

echo "🔧 Fixing Error Code Naming Inconsistencies"
echo "==========================================="

# Fix the error code names in resource files
echo "Updating resource manager files..."

# Replace ErrorCodeAssetLoadFailed with AssetLoadFailed
find internal/ecs/managers/resource -name "*.go" -exec sed -i.bak 's/errors\.ErrorCodeAssetLoadFailed/errors.AssetLoadFailed/g' {} \;

# Replace ErrorCodeAssetInvalid with AssetInvalid  
find internal/ecs/managers/resource -name "*.go" -exec sed -i.bak 's/errors\.ErrorCodeAssetInvalid/errors.AssetInvalid/g' {} \;

# Fix the rendering error
find internal/ecs/core -name "*.go" -exec sed -i.bak 's/errors\.ErrorCodeRenderingFailed/errors.RenderFailed/g' {} \;

echo "✅ Fixed error code naming in resource managers"
echo "✅ Fixed error code naming in core rendering"

# Clean up backup files
find internal/ -name "*.go.bak" -delete

echo ""
echo "🎯 Updated files:"
echo "- internal/ecs/managers/resource/*.go"
echo "- internal/ecs/core/*.go"

echo ""
echo "🚀 Now try: task lint:all"
