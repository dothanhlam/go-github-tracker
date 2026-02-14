#!/bin/bash
set -e

echo "Building Lambda deployment package for ARM64..."

# Build the Go binary for ARM64 Linux
echo "Compiling Go binary..."
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build \
  -tags lambda.norpc \
  -o bootstrap \
  cmd/collector/main.go

# Create deployment package
echo "Creating deployment package..."
zip -j lambda-deployment.zip bootstrap

# Clean up
rm bootstrap

echo "✅ Lambda deployment package created: lambda-deployment.zip"
ls -lh lambda-deployment.zip
