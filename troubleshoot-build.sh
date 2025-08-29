#!/bin/bash

# Build Troubleshooting Script for Air Quality Monitor

echo "=== Build Troubleshooting for Air Quality Monitor ==="
echo

echo "🔍 Checking Go installation..."
if command -v go &> /dev/null; then
    echo "✅ Go is installed: $(go version)"
else
    echo "❌ Go is not installed"
    exit 1
fi

echo
echo "🔍 Checking Go module..."
if [ -f "go.mod" ]; then
    echo "✅ go.mod exists"
    echo "   Module: $(grep '^module' go.mod | cut -d' ' -f2)"
    echo "   Go version: $(grep '^go' go.mod | cut -d' ' -f2)"
else
    echo "❌ go.mod not found"
    exit 1
fi

echo
echo "🔍 Checking Go dependencies..."
if [ -f "go.sum" ]; then
    echo "✅ go.sum exists"
    echo "   Dependencies:"
    go mod graph | head -5
else
    echo "❌ go.sum not found"
fi

echo
echo "🔍 Checking source files..."
GO_FILES=$(find . -name "*.go" -not -path "./.git/*" | wc -l)
echo "✅ Found $GO_FILES Go source files"

echo
echo "🔍 Checking for syntax errors..."
if go build -o /dev/null . 2>&1; then
    echo "✅ Go build syntax check passed"
else
    echo "❌ Go build syntax check failed"
    echo "   Errors:"
    go build -o /dev/null . 2>&1 | head -10
fi

echo
echo "🔍 Running tests..."
if go test -v ./... 2>&1; then
    echo "✅ All tests passed"
else
    echo "❌ Some tests failed"
    echo "   Test output:"
    go test -v ./... 2>&1 | tail -10
fi

echo
echo "🔍 Building application..."
if go build -o air-quality-monitor . 2>&1; then
    echo "✅ Application built successfully"
    echo "   Binary size: $(du -h air-quality-monitor | cut -f1)"
else
    echo "❌ Application build failed"
    echo "   Build errors:"
    go build -o air-quality-monitor . 2>&1 | head -10
fi

echo
echo "🔍 Checking Docker..."
if command -v docker &> /dev/null; then
    echo "✅ Docker is installed: $(docker --version)"
    
    # Check if we can run Docker commands
    if docker info &> /dev/null; then
        echo "✅ Docker daemon is accessible"
        
        echo
        echo "🔍 Testing Docker build..."
        if docker build -t air-quality-monitor:test . 2>&1; then
            echo "✅ Docker build successful"
            echo "   Image created: air-quality-monitor:test"
        else
            echo "❌ Docker build failed"
            echo "   Build errors:"
            docker build -t air-quality-monitor:test . 2>&1 | tail -10
        fi
    else
        echo "⚠️  Docker daemon not accessible (permission issue)"
        echo "   This is normal on some systems and won't affect GitHub Actions"
    fi
else
    echo "❌ Docker is not installed"
fi

echo
echo "🔍 Checking file permissions..."
echo "   Current user: $(whoami)"
echo "   Current directory: $(pwd)"
echo "   Directory permissions: $(ls -ld .)"

echo
echo "🔍 Checking for large files..."
echo "   Large files (>10MB):"
find . -type f -size +10M -not -path "./.git/*" 2>/dev/null | while read file; do
    size=$(du -h "$file" | cut -f1)
    echo "     $file ($size)"
done

echo
echo "🔍 Checking .dockerignore..."
if [ -f ".dockerignore" ]; then
    echo "✅ .dockerignore exists"
    echo "   Excluded patterns:"
    cat .dockerignore | grep -v '^#' | grep -v '^$' | head -10
else
    echo "❌ .dockerignore not found"
fi

echo
echo "🔍 Checking for potential build issues..."
echo "   Files that might cause issues:"

# Check for database files
if [ -f "air_quality.db" ]; then
    size=$(du -h "air_quality.db" | cut -f1)
    echo "     ⚠️  air_quality.db ($size) - should be in .dockerignore"
fi

# Check for log files
if [ -f "server.log" ]; then
    echo "     ⚠️  server.log - should be in .dockerignore"
fi

# Check for binary files
if [ -f "air-quality-monitor" ]; then
    size=$(du -h "air-quality-monitor" | cut -f1)
    echo "     ⚠️  air-quality-monitor ($size) - should be in .dockerignore"
fi

echo
echo "🎯 Recommendations:"
echo "=================="

if [ -f "air_quality.db" ]; then
    echo "1. Remove air_quality.db from git tracking:"
    echo "   git rm --cached air_quality.db"
    echo "   echo '*.db' >> .gitignore"
fi

if [ -f "server.log" ]; then
    echo "2. Remove server.log from git tracking:"
    echo "   git rm --cached server.log"
    echo "   echo '*.log' >> .gitignore"
fi

if [ -f "air-quality-monitor" ]; then
    echo "3. Remove binary from git tracking:"
    echo "   git rm --cached air-quality-monitor"
    echo "   echo 'air-quality-monitor' >> .gitignore"
fi

echo
echo "4. Commit the .dockerignore file:"
echo "   git add .dockerignore"
echo "   git commit -m 'Add .dockerignore for Docker builds'"

echo
echo "5. Push changes to trigger GitHub Actions:"
echo "   git push origin main"

echo
echo "🎉 Troubleshooting completed!"
