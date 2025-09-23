# AI Shell - Makefile
# 提供便捷的构建、测试、安装命令

.PHONY: help build install clean test fmt vet lint dev release run

# 默认目标
.DEFAULT_GOAL := help

# 项目配置
BINARY_NAME := aishell
BUILD_DIR := build
CMD_DIR := cmd/aishell
PKG_DIRS := ./pkg/...

# 构建信息
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "unknown")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GO_VERSION := $(shell go version | awk '{print $$3}')

# ldflags
LDFLAGS := -X main.Version=$(VERSION) \
           -X main.Commit=$(COMMIT) \
           -X main.BuildTime=$(BUILD_TIME) \
           -X main.GoVersion=$(GO_VERSION)

# 构建标志
BUILD_FLAGS := -ldflags="$(LDFLAGS)"
RELEASE_FLAGS := -ldflags="$(LDFLAGS) -s -w"

## help: 显示此帮助信息
help:
	@echo "AI Shell - Makefile 帮助"
	@echo ""
	@echo "可用命令:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'
	@echo ""
	@echo "示例:"
	@echo "  make build          # 标准构建"
	@echo "  make dev            # 开发模式构建"
	@echo "  make release        # 发布模式构建"
	@echo "  make install        # 构建并安装"
	@echo "  make test           # 运行测试"
	@echo "  make clean          # 清理构建文件"

## build: 标准模式构建
build: clean-build
	@echo "🔨 构建 $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./$(CMD_DIR)
	@ln -sf build/$(BINARY_NAME) $(BINARY_NAME)
	@echo "✅ 构建完成: $(BUILD_DIR)/$(BINARY_NAME)"

## dev: 开发模式构建（保留调试信息）
dev: clean-build
	@echo "🔨 开发模式构建 $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build $(BUILD_FLAGS) -race -o $(BUILD_DIR)/$(BINARY_NAME) ./$(CMD_DIR)
	@ln -sf build/$(BINARY_NAME) $(BINARY_NAME)
	@echo "✅ 开发构建完成: $(BUILD_DIR)/$(BINARY_NAME)"

## release: 发布模式构建（优化体积）
release: clean-build test
	@echo "🔨 发布模式构建 $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build $(RELEASE_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./$(CMD_DIR)
	@ln -sf build/$(BINARY_NAME) $(BINARY_NAME)
	@echo "✅ 发布构建完成: $(BUILD_DIR)/$(BINARY_NAME)"

## install: 构建并安装到系统
install: release
	@echo "📦 安装 $(BINARY_NAME)..."
	@./scripts/install.sh

## install-user: 构建并安装到用户目录
install-user: release
	@echo "📦 安装 $(BINARY_NAME) 到用户目录..."
	@./scripts/install.sh --user

## test: 运行所有测试
test:
	@echo "🧪 运行测试..."
	@go test -v $(PKG_DIRS)

## test-race: 运行测试（启用竞态检测）
test-race:
	@echo "🧪 运行测试（竞态检测）..."
	@go test -race -v $(PKG_DIRS)

## test-cover: 运行测试并生成覆盖率报告
test-cover:
	@echo "🧪 生成测试覆盖率报告..."
	@go test -coverprofile=coverage.out $(PKG_DIRS)
	@go tool cover -html=coverage.out -o coverage.html
	@echo "📊 覆盖率报告: coverage.html"

## bench: 运行基准测试
bench:
	@echo "⚡ 运行基准测试..."
	@go test -bench=. -benchmem $(PKG_DIRS)

## fmt: 格式化代码
fmt:
	@echo "🎨 格式化代码..."
	@go fmt ./...

## vet: 静态检查
vet:
	@echo "🔍 静态检查..."
	@go vet ./...

## lint: 代码检查（需要 golangci-lint）
lint:
	@echo "🔧 代码检查..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "⚠️  golangci-lint 未安装，跳过检查"; \
		echo "安装方法: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

## tidy: 整理依赖
tidy:
	@echo "📦 整理依赖..."
	@go mod tidy
	@go mod verify

## deps: 下载依赖
deps:
	@echo "📦 下载依赖..."
	@go mod download

## run: 构建并运行
run: build
	@echo "🚀 运行 $(BINARY_NAME)..."
	@./$(BUILD_DIR)/$(BINARY_NAME)

## run-dev: 开发模式运行
run-dev: dev
	@echo "🚀 开发模式运行 $(BINARY_NAME)..."
	@AISHELL_DEBUG=true ./$(BUILD_DIR)/$(BINARY_NAME)

## clean: 清理构建文件
clean:
	@echo "🧹 清理构建文件..."
	@rm -rf $(BUILD_DIR)
	@rm -f $(BINARY_NAME)
	@rm -f coverage.out coverage.html
	@echo "✅ 清理完成"

## clean-build: 清理构建目录
clean-build:
	@rm -rf $(BUILD_DIR)

## docker-build: 构建Docker镜像
docker-build:
	@echo "🐳 构建Docker镜像..."
	@docker build -t $(BINARY_NAME):$(VERSION) .
	@docker tag $(BINARY_NAME):$(VERSION) $(BINARY_NAME):latest

## info: 显示构建信息
info:
	@echo "📋 构建信息:"
	@echo "  版本: $(VERSION)"
	@echo "  提交: $(COMMIT)"
	@echo "  时间: $(BUILD_TIME)"
	@echo "  Go版本: $(GO_VERSION)"
	@echo "  二进制名: $(BINARY_NAME)"
	@echo "  构建目录: $(BUILD_DIR)"

## size: 显示二进制文件大小
size:
	@if [ -f "$(BUILD_DIR)/$(BINARY_NAME)" ]; then \
		echo "📏 文件大小:"; \
		ls -lh $(BUILD_DIR)/$(BINARY_NAME); \
		echo ""; \
		du -h $(BUILD_DIR)/$(BINARY_NAME); \
	else \
		echo "❌ 二进制文件不存在，请先构建"; \
	fi

## check: 运行所有检查
check: fmt vet lint test

## ci: CI/CD 流水线
ci: deps check test-cover release

## all: 完整构建流程
all: clean deps check test build
