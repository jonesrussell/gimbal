#!/bin/bash

# scripts/standardize-error-codes.sh
# Standardize all error code naming to match the constants in errors.go

set -e

echo "🔧 Standardizing Error Code Naming"
echo "================================="

# Function to replace error codes in files
replace_errors() {
    local pattern="$1"
    local replacement="$2" 
    local description="$3"
    
    echo "Replacing $pattern -> $replacement ($description)"
    
    # Find all Go files and replace the pattern
    find internal -name "*.go" -type f -exec sed -i.bak "s/errors\.${pattern}/errors.${replacement}/g" {} \;
}

# Asset errors
replace_errors "ErrorCodeAssetLoadFailed" "AssetLoadFailed" "asset load failures"
replace_errors "ErrorCodeAssetInvalid" "AssetInvalid" "invalid assets"
replace_errors "ErrorCodeAssetNotFound" "AssetNotFound" "missing assets"

# System errors
replace_errors "ErrorCodeSystemFailed" "SystemInitFailed" "system failures"
replace_errors "ErrorCodeSystemInitFailed" "SystemInitFailed" "system init failures" 
replace_errors "ErrorCodeSystemUpdateFailed" "SystemUpdateFailed" "system update failures"
replace_errors "ErrorCodeSystemCleanupFailed" "SystemCleanupFailed" "system cleanup failures"

# Config errors
replace_errors "ErrorCodeConfigMissing" "ConfigMissing" "missing config"
replace_errors "ErrorCodeConfigInvalid" "ConfigInvalid" "invalid config"
replace_errors "ErrorCodeConfigValidation" "ConfigValidation" "config validation"

# Entity errors
replace_errors "ErrorCodeEntityCreationFailed" "EntityInvalid" "entity creation failures"
replace_errors "ErrorCodeEntityNotFound" "EntityNotFound" "missing entities"
replace_errors "ErrorCodeEntityInvalid" "EntityInvalid" "invalid entities"

# Rendering errors
replace_errors "ErrorCodeRenderingFailed" "RenderFailed" "rendering failures"
replace_errors "ErrorCodeRenderFailed" "RenderFailed" "render failures"
replace_errors "ErrorCodeRenderTimeout" "RenderTimeout" "render timeouts"
replace_errors "ErrorCodeRenderUnsupported" "RenderUnsupported" "unsupported rendering"

# Scene errors
replace_errors "ErrorCodeSceneNotFound" "SceneNotFound" "missing scenes"
replace_errors "ErrorCodeSceneTransition" "SceneTransition" "scene transitions"
replace_errors "ErrorCodeSceneLoadFailed" "SceneLoadFailed" "scene load failures"

# Sprite/Resource errors - map to appropriate existing constants
replace_errors "ErrorCodeSpriteNotFound" "AssetNotFound" "missing sprites"
replace_errors "ErrorCodeResourceNotFound" "ResourceNotFound" "missing resources"
replace_errors "ErrorCodeResourceLoadFailed" "ResourceLoadFailed" "resource load failures"
replace_errors "ErrorCodeResourceExhausted" "ResourceExhausted" "exhausted resources"
replace_errors "ErrorCodeResourceLocked" "ResourceLocked" "locked resources"

# Input errors
replace_errors "ErrorCodeInputInvalid" "InputInvalid" "invalid input"
replace_errors "ErrorCodeInputUnsupported" "InputUnsupported" "unsupported input"
replace_errors "ErrorCodeInputTimeout" "InputTimeout" "input timeouts"

# Validation errors
replace_errors "ErrorCodeValidationFailed" "ValidationFailed" "validation failures"
replace_errors "ErrorCodeValidationTimeout" "ValidationTimeout" "validation timeouts"

echo ""
echo "🧹 Cleaning up backup files..."
find internal -name "*.go.bak" -delete

echo ""
echo "✅ Standardization complete!"
echo ""
echo "📝 Summary of changes:"
echo "- All error codes now use the standard naming (e.g., errors.AssetLoadFailed)"
echo "- Removed 'ErrorCode' prefix for consistency"
echo "- Mapped logical equivalents (e.g., SpriteNotFound -> AssetNotFound)"
echo ""
echo "🚀 Test the build: task lint:all"
