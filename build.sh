#!/bin/bash
set -e

echo "Killing any running ce-node-exporter processes..."
pkill -f ce-node-exporter || true

echo "Building CE Node Exporter..."
/opt/compiler-explorer/golang-1.24.2/go/bin/go mod download
/opt/compiler-explorer/golang-1.24.2/go/bin/go build -o ce-node-exporter main.go
echo "Build complete!"