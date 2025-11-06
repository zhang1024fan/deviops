# ============================================
# Makefile for dodevops-api
# ============================================

.PHONY: help build run test clean docker-build docker-run docker-stop docker-logs

# 默认目标
.DEFAULT_GOAL := help

# 项目信息
APP_NAME := dodevops-api
IMAGE_NAME := dodevops-api
IMAGE_TAG := latest
CONTAINER_NAME := dodevops-api

# 帮助信息
help: ## 显示帮助信息
	@echo "dodevops-api Makefile 命令:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

# ==================== 本地开发 ====================

build: ## 编译应用程序
	@echo "正在编译 $(APP_NAME)..."
	@go build -o $(APP_NAME) .
	@echo "编译完成: ./$(APP_NAME)"

run: ## 运行应用程序
	@echo "启动 $(APP_NAME)..."
	@go run main.go

test: ## 运行测试
	@echo "运行测试..."
	@go test -v ./...

tidy: ## 整理依赖
	@echo "整理 Go 模块..."
	@go mod tidy
	@go mod download

clean: ## 清理构建产物
	@echo "清理构建产物..."
	@rm -f $(APP_NAME)
	@rm -rf log/*
	@rm -rf tmp/*
	@echo "清理完成"

# ==================== Docker 相关 ====================

docker-build: ## 构建 Docker 镜像
	@echo "构建 Docker 镜像: $(IMAGE_NAME):$(IMAGE_TAG)..."
	@docker build -t $(IMAGE_NAME):$(IMAGE_TAG) .
	@echo "镜像构建完成"

docker-build-no-cache: ## 不使用缓存构建 Docker 镜像
	@echo "无缓存构建 Docker 镜像: $(IMAGE_NAME):$(IMAGE_TAG)..."
	@docker build --no-cache -t $(IMAGE_NAME):$(IMAGE_TAG) .
	@echo "镜像构建完成"

docker-run: ## 运行 Docker 容器
	@echo "启动 Docker 容器: $(CONTAINER_NAME)..."
	@docker run -d \
		--name $(CONTAINER_NAME) \
		-p 8000:8000 \
		-v $(PWD)/config.yaml:/app/config.yaml:ro \
		-v $(PWD)/log:/app/log \
		-v $(PWD)/upload:/app/upload \
		--restart unless-stopped \
		$(IMAGE_NAME):$(IMAGE_TAG)
	@echo "容器已启动"

docker-stop: ## 停止并删除 Docker 容器
	@echo "停止 Docker 容器: $(CONTAINER_NAME)..."
	@docker stop $(CONTAINER_NAME) 2>/dev/null || true
	@docker rm $(CONTAINER_NAME) 2>/dev/null || true
	@echo "容器已停止并删除"

docker-restart: docker-stop docker-run ## 重启 Docker 容器

docker-logs: ## 查看 Docker 容器日志
	@docker logs -f $(CONTAINER_NAME)

docker-exec: ## 进入 Docker 容器
	@docker exec -it $(CONTAINER_NAME) /bin/sh

docker-ps: ## 查看容器状态
	@docker ps -a | grep $(CONTAINER_NAME)

docker-images: ## 查看镜像
	@docker images | grep $(IMAGE_NAME)

# ==================== Docker Compose 相关 ====================

up: ## 启动服务（docker-compose）
	@echo "启动服务..."
	@docker-compose up -d
	@echo "服务已启动"

down: ## 停止服务（docker-compose）
	@echo "停止服务..."
	@docker-compose down
	@echo "服务已停止"

logs: ## 查看服务日志（docker-compose）
	@docker-compose logs -f

restart: down up ## 重启服务（docker-compose）

ps: ## 查看服务状态（docker-compose）
	@docker-compose ps

# ==================== 生产部署 ====================

deploy: docker-build up ## 部署到生产环境（构建 + 启动）

deploy-prod: ## 生产环境部署（使用标签）
	@echo "生产环境部署..."
	@docker build -t $(IMAGE_NAME):$(IMAGE_TAG) -t $(IMAGE_NAME):prod .
	@docker-compose -f docker-compose.yml up -d
	@echo "生产环境部署完成"

# ==================== 清理 ====================

prune: ## 清理 Docker 无用资源
	@echo "清理 Docker 资源..."
	@docker system prune -f
	@docker volume prune -f
	@echo "清理完成"

clean-all: clean docker-stop prune ## 清理所有资源（本地 + Docker）

# ==================== 其他工具 ====================

swagger: ## 生成 Swagger 文档
	@echo "生成 Swagger 文档..."
	@swag init
	@echo "Swagger 文档生成完成"

fmt: ## 格式化代码
	@echo "格式化代码..."
	@go fmt ./...
	@echo "代码格式化完成"

lint: ## 代码检查
	@echo "运行代码检查..."
	@golangci-lint run
	@echo "代码检查完成"
