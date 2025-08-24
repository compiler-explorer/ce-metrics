#!/bin/bash
set -e

echo "Killing any running ce-node-exporter processes..."
pkill -f ce-node-exporter || true

# Try to find go in PATH first, otherwise use the specific version
GO_BIN=$(which go 2>/dev/null || echo "/opt/compiler-explorer/golang-1.24.2/go/bin/go")

echo "Using Go binary: $GO_BIN"
echo "Building CE Node Exporter..."
$GO_BIN mod download
$GO_BIN build -o ce-node-exporter main.go
echo "Build complete!"