#!/bin/bash

# NodeImage Backup Tool 跨平台构建脚本
# 支持 macOS, Windows, Linux

set -e

# 版本信息
VERSION="1.0.0"
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# 构建标志
LDFLAGS="-X main.Version=$VERSION -X main.BuildTime=$BUILD_TIME -X main.GitCommit=$GIT_COMMIT -s -w"

echo "🚀 开始构建 NodeImage Backup Tool (nib) v$VERSION"
echo "构建时间: $BUILD_TIME"
echo "Git提交: $GIT_COMMIT"
echo ""

# 创建构建目录
mkdir -p build
rm -rf build/*

# 构建 macOS (Intel)
echo "📦 构建 macOS (Intel)..."
GOOS=darwin GOARCH=amd64 go build -ldflags "$LDFLAGS" -o build/nib-darwin-amd64 main.go

# 构建 macOS (Apple Silicon)
echo "📦 构建 macOS (Apple Silicon)..."
GOOS=darwin GOARCH=arm64 go build -ldflags "$LDFLAGS" -o build/nib-darwin-arm64 main.go

# 构建 Linux (Intel)
echo "📦 构建 Linux (Intel)..."
GOOS=linux GOARCH=amd64 go build -ldflags "$LDFLAGS" -o build/nib-linux-amd64 main.go

# 构建 Linux (ARM64)
echo "📦 构建 Linux (ARM64)..."
GOOS=linux GOARCH=arm64 go build -ldflags "$LDFLAGS" -o build/nib-linux-arm64 main.go

# 构建 Windows (Intel)
echo "📦 构建 Windows (Intel)..."
GOOS=windows GOARCH=amd64 go build -ldflags "$LDFLAGS" -o build/nib-windows-amd64.exe main.go

# 构建 Windows (ARM64)
echo "📦 构建 Windows (ARM64)..."
GOOS=windows GOARCH=arm64 go build -ldflags "$LDFLAGS" -o build/nib-windows-arm64.exe main.go

# 创建压缩包
echo ""
echo "🗜️  创建压缩包..."

# macOS
cd build
tar -czf nib-darwin-amd64.tar.gz nib-darwin-amd64
tar -czf nib-darwin-arm64.tar.gz nib-darwin-arm64

# Linux
tar -czf nib-linux-amd64.tar.gz nib-linux-amd64
tar -czf nib-linux-arm64.tar.gz nib-linux-arm64

# Windows
zip nib-windows-amd64.zip nib-windows-amd64.exe
zip nib-windows-arm64.zip nib-windows-arm64.exe

cd ..

echo ""
echo "✅ 构建完成!"
echo ""
echo "📁 构建文件位置:"
ls -la build/
echo ""
echo "📋 文件大小:"
du -h build/*

echo ""
echo "🎯 使用说明:"
echo "  macOS Intel:   ./build/nib-darwin-amd64"
echo "  macOS ARM64:   ./build/nib-darwin-arm64"
echo "  Linux Intel:   ./build/nib-linux-amd64"
echo "  Linux ARM64:   ./build/nib-linux-arm64"
echo "  Windows Intel: ./build/nib-windows-amd64.exe"
echo "  Windows ARM64: ./build/nib-windows-arm64.exe"
echo ""
echo "�� 快速开始:"
echo "  1. 复制配置文件: cp nib.yaml.example nib.yaml"
echo "  2. 编辑配置文件: nano nib.yaml (填入你的API Token)"
echo "  3. 执行同步: ./nib"
echo "  4. 查看列表: ./nib list" 