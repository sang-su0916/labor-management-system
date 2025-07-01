#!/bin/bash
set -e

echo "Simple Go build script"
echo "Working directory: $(pwd)"
echo "Go version: $(go version)"

# Clean and create bin directory
rm -rf bin/
mkdir -p bin

# Try simple build first
echo "Attempting simple build..."
if go build -o bin/main ./cmd/server/main.go; then
    echo "Simple build successful"
    ls -la bin/main
    file bin/main
    chmod +x bin/main
    exit 0
fi

# Try with CGO disabled
echo "Attempting CGO disabled build..."
if CGO_ENABLED=0 go build -o bin/main ./cmd/server/main.go; then
    echo "CGO disabled build successful"
    ls -la bin/main
    file bin/main
    chmod +x bin/main
    exit 0
fi

echo "All build attempts failed"
exit 1