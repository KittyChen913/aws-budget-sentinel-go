#!/bin/bash
# Bash build script for Go Lambda function

echo "Building Go Lambda function for Linux..."

export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=0

# 切換到專案根目錄
cd "$(dirname "$0")/.."

# 編譯
go build -ldflags="-s -w" -o bootstrap ./cmd/lambda

if [ $? -eq 0 ]; then
    echo "Build successful! Output: bootstrap"
else
    echo "Build failed!"
    exit 1
fi
