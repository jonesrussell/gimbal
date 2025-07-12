#!/bin/bash

# Go Project Architecture Analysis Script
# Run from project root: ./scripts/analyze-architecture.sh

set -e

echo "🔍 Go Project Architecture Analysis"
echo "====================================="
echo ""

# Change to project root (assuming script is in scripts/ directory)
cd "$(dirname "$0")/.."

echo "📁 PROJECT STRUCTURE"
echo "---------------------"
echo "Project root: $(pwd)"
echo ""

# Check if tree is available, fallback to find
if command -v tree &> /dev/null; then
    echo "Directory structure:"
    tree -I 'vendor|node_modules|.git|*.exe|*.so|*.dylib' -L 3
else
    echo "Directory structure (using find):"
    find . -type d -name ".git" -prune -o -type d -name "vendor" -prune -o -type d -name "node_modules" -prune -o -type d -print | head -20
fi

echo ""
echo "Go files structure:"
find . -type f -name "*.go" | grep -v vendor | head -20

echo ""
echo "📦 DEPENDENCIES & MODULES"
echo "-------------------------"
if [ -f "go.mod" ]; then
    echo "Go module info:"
    head -10 go.mod
    echo ""
    echo "Dependency graph (top 20):"
    go mod graph | head -20
else
    echo "❌ No go.mod found - not a Go module project"
fi

echo ""
echo "📊 CODE METRICS"
echo "---------------"
echo "Go files count:"
find . -name "*.go" -not -path "./vendor/*" | wc -l

echo ""
echo "Total lines of Go code:"
find . -name "*.go" -not -path "./vendor/*" -exec wc -l {} + | tail -1

echo ""
echo "Lines per file (largest files first):"
find . -name "*.go" -not -path "./vendor/*" -exec wc -l {} + | sort -nr | head -10

echo ""
echo "Function/method counts:"
echo "Total functions: $(grep -r "^func " . --include="*.go" --exclude-dir=vendor | wc -l)"
echo "Methods (receiver functions): $(grep -r "^func (" . --include="*.go" --exclude-dir=vendor | wc -l)"

echo ""
echo "🎯 MAIN ENTRY POINTS"
echo "--------------------"
echo "Main files found:"
find . -name "main.go" -not -path "./vendor/*"

echo ""
echo "Main functions:"
grep -r "func main()" . --include="*.go" --exclude-dir=vendor

echo ""
echo "🛠️  BUILD TOOLS & TASKS"
echo "----------------------"
if command -v task &> /dev/null; then
    echo "Available tasks:"
    task --list
    echo ""
else
    echo "❌ Task not found - install from https://taskfile.dev"
fi

echo ""
echo "🔧 STATIC ANALYSIS"
echo "------------------"

# Check for common static analysis tools
echo "Checking for static analysis tools..."

if command -v gocyclo &> /dev/null; then
    echo "Cyclomatic complexity (top 10 most complex functions):"
    gocyclo -top 10 .
    echo ""
else
    echo "📝 Install gocyclo for complexity analysis: go install github.com/fzipp/gocyclo/cmd/gocyclo@latest"
fi

if command -v golint &> /dev/null; then
    echo "Lint issues (first 20):"
    golint ./... | head -20
    echo ""
else
    echo "📝 Install golint: go install golang.org/x/lint/golint@latest"
fi

if command -v staticcheck &> /dev/null; then
    echo "Static analysis issues (first 20):"
    staticcheck ./... | head -20
    echo ""
else
    echo "📝 Install staticcheck: go install honnef.co/go/tools/cmd/staticcheck@latest"
fi

# Basic pattern analysis
echo ""
echo "🔍 PATTERN ANALYSIS"
echo "-------------------"
echo "Struct definitions:"
grep -r "^type.*struct" . --include="*.go" --exclude-dir=vendor | wc -l

echo ""
echo "Interface definitions:"
grep -r "^type.*interface" . --include="*.go" --exclude-dir=vendor | wc -l

echo ""
echo "Package declarations:"
grep -r "^package " . --include="*.go" --exclude-dir=vendor | cut -d: -f2 | sort | uniq -c | sort -nr

echo ""
echo "🚨 POTENTIAL ISSUES"
echo "-------------------"
echo "Files with >500 lines (potential refactoring candidates):"
find . -name "*.go" -not -path "./vendor/*" -exec wc -l {} + | awk '$1 > 500 {print $2 " (" $1 " lines)"}' | sort

echo ""
echo "Long function names (>30 chars, might indicate SRP violations):"
grep -r "^func " . --include="*.go" --exclude-dir=vendor | grep -E "func [^(]{30,}\(" | head -10

echo ""
echo "Global variables (potential state management issues):"
grep -r "^var " . --include="*.go" --exclude-dir=vendor | head -10

echo ""
echo "Init functions (initialization complexity):"
grep -r "func init()" . --include="*.go" --exclude-dir=vendor

echo ""
echo "✅ ANALYSIS COMPLETE"
echo "===================="
echo "Review the output above to identify:"
echo "- Large files that need refactoring"
echo "- Complex functions that violate SRP"
echo "- Missing interfaces (low interface count vs struct count)"
echo "- Package organization issues"
echo "- Global state management problems"
echo ""
echo "Next steps:"
echo "1. Install missing static analysis tools"
echo "2. Review largest/most complex files first"
echo "3. Look for common anti-patterns"
echo "4. Check for proper separation of concerns"
