# NodeImage Backup Tool Makefile

.PHONY: help build clean test run list sync cross-build install

# 默认目标
help:
	@echo "NodeImage Backup Tool 构建工具"
	@echo ""
	@echo "可用命令:"
	@echo "  build       - 构建当前平台版本"
	@echo "  cross-build - 构建所有平台版本"
	@echo "  clean       - 清理构建产物"
	@echo "  test        - 运行测试"
	@echo "  run         - 本地运行程序"
	@echo "  install     - 安装到系统路径"
	@echo "  list        - 查看远程图片列表"
	@echo "  sync        - 同步远程图片到本地"

# 构建当前平台版本
build:
	@echo "🔨 构建当前平台版本..."
	go build -ldflags "-s -w" -o nib main.go
	@echo "✅ 构建完成: ./nib"

# 清理构建产物
clean:
	@echo "🧹 清理构建产物..."
	rm -f nib nib.exe
	rm -rf build/
	@echo "✅ 清理完成"

# 运行测试
test:
	@echo "🧪 运行测试..."
	./test.sh

# 本地运行程序
run:
	@echo "🚀 本地运行程序..."
	go run main.go

# 安装到系统路径
install: build
	@echo "📦 安装到系统路径..."
	@if [ "$(shell uname)" = "Darwin" ]; then \
		sudo cp nib /usr/local/bin/; \
		echo "✅ 已安装到 /usr/local/bin/nib"; \
	elif [ "$(shell uname)" = "Linux" ]; then \
		sudo cp nib /usr/local/bin/; \
		echo "✅ 已安装到 /usr/local/bin/nib"; \
	else \
		echo "❌ 不支持的平台"; \
	fi

# 跨平台构建
cross-build:
	@echo "🌍 跨平台构建..."
	./build.sh

# 查看远程图片列表 (需要提供token)
list:
	@echo "📋 查看远程图片列表..."
	@echo "请提供API Token:"
	@read -p "Token: " token; \
	./nib list -t $$token

# 同步远程图片 (需要提供token)
sync:
	@echo "🔄 同步远程图片..."
	@echo "请提供API Token:"
	@read -p "Token: " token; \
	echo "请提供本地目录 (默认: 程序目录/images):"; \
	read -p "目录: " dir; \
	if [ -z "$$dir" ]; then dir=""; fi; \
	./nib -t $$token -d $$dir

# 开发模式 (监听文件变化并重新构建)
dev:
	@echo "👨‍💻 开发模式..."
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "❌ 请先安装 air: go install github.com/cosmtrek/air@latest"; \
	fi

# 检查依赖
deps:
	@echo "📦 检查依赖..."
	go mod tidy
	go mod verify

# 格式化代码
fmt:
	@echo "🎨 格式化代码..."
	go fmt ./...

# 代码检查
lint:
	@echo "🔍 代码检查..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "❌ 请先安装 golangci-lint"; \
	fi

# 显示版本信息
version:
	@echo "📋 版本信息:"
	@echo "  Go版本: $(shell go version)"
	@echo "  构建时间: $(shell date)"
	@echo "  Git提交: $(shell git rev-parse --short HEAD 2>/dev/null || echo 'unknown')" 