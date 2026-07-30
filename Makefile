APP_NAME    := hypersku-cli
VERSION     := 0.1.0
COMMIT      := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE        := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS     := -X github.com/hypersku/hypersku-cli/internal/version.Version=$(VERSION) \
               -X github.com/hypersku/hypersku-cli/internal/version.Commit=$(COMMIT) \
               -X github.com/hypersku/hypersku-cli/internal/version.Date=$(DATE)

GO          := go
GOFLAGS     := -ldflags "$(LDFLAGS)"
OUTPUT_DIR  := build

.PHONY: all build clean test lint run help

all: clean build

build: ## 编译项目
	@mkdir -p $(OUTPUT_DIR)
	$(GO) build $(GOFLAGS) -o $(OUTPUT_DIR)/$(APP_NAME).exe .

build-linux: ## 编译 Linux 版本
	@mkdir -p $(OUTPUT_DIR)
	SET GOOS=linux&& SET GOARCH=amd64&& $(GO) build $(GOFLAGS) -o $(OUTPUT_DIR)/$(APP_NAME)-linux-amd64 .

build-macos: ## 编译 macOS 版本
	@mkdir -p $(OUTPUT_DIR)
	SET GOOS=darwin&& SET GOARCH=amd64&& $(GO) build $(GOFLAGS) -o $(OUTPUT_DIR)/$(APP_NAME)-darwin-amd64 .

clean: ## 清理构建产物
	@rm -rf $(OUTPUT_DIR)
	@echo "✓ 清理完成"

test: ## 运行测试
	$(GO) test -v ./...

lint: ## 代码检查
	$(GO) vet ./...

run: ## 直接运行
	$(GO) run . --help

tidy: ## 整理依赖
	$(GO) mod tidy

help: ## 显示帮助
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
