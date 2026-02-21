#!/usr/bin/env bash
# Optimize image assets: strip metadata (ImageMagick), then lossless PNG compression (optipng, pngcrush).
# Run from repository root. Resize only if files exceed max display size (see docs/ASSETS_REQUIRED.md).
set -e
ROOT="${1:-.}"
cd "$ROOT"
echo "Stripping metadata (mogrify -strip)..."
for dir in assets/sprites assets/ui assets/cutscenes assets/planets assets/bosses assets/ending; do
  [ -d "$dir" ] || continue
  mogrify -strip "$dir"/*.png 2>/dev/null || true
done
echo "Running optipng -o2..."
for dir in assets/sprites assets/ui assets/cutscenes assets/planets assets/bosses assets/ending; do
  [ -d "$dir" ] || continue
  for f in "$dir"/*.png; do
    [ -f "$f" ] || continue
    optipng -o2 "$f" 2>/dev/null || true
  done
done
echo "Running pngcrush..."
for dir in assets/sprites assets/ui assets/cutscenes assets/planets assets/bosses assets/ending; do
  [ -d "$dir" ] || continue
  for f in "$dir"/*.png; do
    [ -f "$f" ] || continue
    tmp=$(mktemp)
    pngcrush -q "$f" "$tmp" 2>/dev/null && mv "$tmp" "$f" || rm -f "$tmp"
  done
done
echo "Done."
