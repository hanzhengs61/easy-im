.PHONY: help dev build test lint gen tidy docker-up docker-down

# 项目基础变量
PROJECT     = easy-im
MODULE      = github.com/yourname/easy-im
GO          = go
GOCTL       = goctl
SERVICES    = gateway user auth message group push ws

# 颜色输出
GREEN  = \033[0;32m
YELLOW = \033[0;33m
RESET  = \033[0m

help: ## 显示帮助
	@echo "$(GREEN)easy-im 开发命令$(RESET)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  $(YELLOW)%-20s$(RESET) %s\n", $$1, $$2}'

dev: docker-up ## 启动本地开发环境（依赖服务）

tidy: ## 整理 go.mod
	$(GO) mod tidy

lint: ## 代码静态检查
	golangci-lint run ./...

test: ## 运行单元测试
	$(GO) test -v -race -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html

build: ## 编译所有服务
	@for svc in $(SERVICES); do \
		echo "$(GREEN)Building $$svc...$(RESET)"; \
		$(GO) build -o bin/$$svc ./cmd/$$svc/; \
	done

gen: ## 生成 go-zero 代码（需指定 svc=user）
	$(GOCTL) api go -api api/$(svc).api -dir internal/$(svc) -style goZero

gen-model: ## 生成数据库 model（需指定 svc=user table=users）
	$(GOCTL) model mysql datasource -url "root:123456@tcp(127.0.0.1:3306)/easy_im" \
		-table $(table) -dir internal/$(svc)/model

docker-up: ## 启动依赖服务（MySQL Redis Kafka etcd）
	docker compose -f deploy/docker-compose.yml up -d

docker-down: ## 停止依赖服务
	docker compose -f deploy/deploy/docker-compose.yml down

docker-build: ## 构建所有服务镜像
	@for svc in $(SERVICES); do \
		docker build -f build/docker/$$svc/Dockerfile -t easy-im-$$svc:latest .; \
	done

clean: ## 清理编译产物
	rm -rf bin/ coverage.out coverage.html

goctl: ## 生成api
	goctl api go -api api/user.api -dir internal/user
gomodel: ## 生成model
	goctl model mysql ddl -src deploy/sql/user.sql -dir internal/user/model --style go_zero
.DEFAULT_GOAL := help