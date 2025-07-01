#!/bin/bash
echo "🔨 Building Labor Management System..."
echo "Current directory: $(pwd)"
echo "Directory contents:"
ls -la

# Build the Go application
echo "📦 Building Go application..."
go build -o bin/main cmd/server/main.go

# Check if web directory exists
if [ -d "web" ]; then
    echo "✅ Web directory found"
    echo "Web directory contents:"
    ls -la web/
else
    echo "❌ Web directory not found!"
fi

echo "✅ Build complete!"