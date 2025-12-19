.PHONY: help build run test clean docker-build docker-up docker-down migrate lint

# 默认目标
.DEFAULT_GOAL := help

# 变量定义
APP_NAME := shoppee
BUILD_DIR := ./bin
MAIN_FILE := ./cmd/api/main.go
DOCKER_IMAGE := shoppee:latest

# 帮助信息
help: ## 显示帮助信息
	@echo "Shoppee 电商系统 - Makefile命令"
	@echo ""
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ 开发

build: ## 编译应用
	@echo "编译应用..."
	@mkdir -p $(BUILD_DIR)
	@go build -ldflags="-w -s" -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_FILE)
	@echo "编译完成: $(BUILD_DIR)/$(APP_NAME)"

run: ## 运行应用
	@echo "运行应用..."
	@go run $(MAIN_FILE)

dev: ## 开发模式（热重载需要安装air）
	@echo "开发模式..."
	@air

test: ## 运行测试
	@echo "运行测试..."
	@go test -v -cover ./...

test-coverage: ## 生成测试覆盖率报告
	@echo "生成覆盖率报告..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "覆盖率报告已生成: coverage.html"

lint: ## 代码检查
	@echo "运行代码检查..."
	@golangci-lint run ./...

fmt: ## 格式化代码
	@echo "格式化代码..."
	@go fmt ./...
	@goimports -w .

clean: ## 清理编译产物
	@echo "清理编译产物..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@echo "清理完成"

##@ 依赖管理

deps: ## 下载依赖
	@echo "下载依赖..."
	@go mod download

tidy: ## 整理依赖
	@echo "整理依赖..."
	@go mod tidy

vendor: ## 创建vendor目录
	@echo "创建vendor..."
	@go mod vendor

##@ Docker

docker-build: ## 构建Docker镜像
	@echo "构建Docker镜像..."
	@docker build -t $(DOCKER_IMAGE) .
	@echo "镜像构建完成: $(DOCKER_IMAGE)"

docker-up: ## 启动Docker Compose服务
	@echo "启动Docker服务..."
	@docker-compose up -d
	@echo "服务已启动"
	@docker-compose ps

docker-down: ## 停止Docker Compose服务
	@echo "停止Docker服务..."
	@docker-compose down
	@echo "服务已停止"

docker-logs: ## 查看Docker日志
	@docker-compose logs -f app

docker-restart: ## 重启Docker服务
	@echo "重启Docker服务..."
	@docker-compose restart app

##@ 数据库

migrate: ## 运行数据库迁移
	@echo "运行数据库迁移..."
	@go run $(MAIN_FILE) migrate

db-reset: ## 重置数据库（危险操作）
	@echo "重置数据库..."
	@docker-compose exec postgres psql -U postgres -c "DROP DATABASE IF EXISTS shoppee;"
	@docker-compose exec postgres psql -U postgres -c "CREATE DATABASE shoppee;"
	@echo "数据库已重置"

##@ 部署

build-linux: ## 交叉编译Linux版本
	@echo "编译Linux版本..."
	@mkdir -p $(BUILD_DIR)
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o $(BUILD_DIR)/$(APP_NAME)-linux-amd64 $(MAIN_FILE)
	@echo "编译完成: $(BUILD_DIR)/$(APP_NAME)-linux-amd64"

build-windows: ## 交叉编译Windows版本
	@echo "编译Windows版本..."
	@mkdir -p $(BUILD_DIR)
	@CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-w -s" -o $(BUILD_DIR)/$(APP_NAME)-windows-amd64.exe $(MAIN_FILE)
	@echo "编译完成: $(BUILD_DIR)/$(APP_NAME)-windows-amd64.exe"

build-mac: ## 交叉编译macOS版本
	@echo "编译macOS版本..."
	@mkdir -p $(BUILD_DIR)
	@CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="-w -s" -o $(BUILD_DIR)/$(APP_NAME)-darwin-amd64 $(MAIN_FILE)
	@echo "编译完成: $(BUILD_DIR)/$(APP_NAME)-darwin-amd64"

build-all: build-linux build-windows build-mac ## 编译所有平台版本
	@echo "所有平台编译完成"
