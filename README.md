# 项目总览：easy-im

# 系统架构图
<img width="784" height="871" alt="image" src="https://github.com/user-attachments/assets/ec19cf74-883c-4236-a3c9-ea34109e8442" />


## 技术栈选型
```
1. 微服务框架：go-zero
2. 通信协议：HTTP REST + gRPC + WebSocket
3. 数据库：MySQL 8 + Redis 7 + MongoDB
4. 消息队列：Kafka
5. 服务注册：etcd
6. 缓存框架：Redis
7. 可观测性：Prometheus + Grafana + Jaeger
8. 容器化：Docker + Kubernetes
9. CI/CD：GitHub Actions + Makefile
10. 代码规范：golangci-lint + gofmt
```
## 项目结构
```
easy-im/
├── .github/
│   └── workflows/
│       └── ci.yml                  # GitHub Actions 流水线
├── api/                            # Proto / API 定义文件（OpenAPI / .proto）
│   ├── user.proto
│   └── message.proto
├── build/                          # Docker、K8s 构建配置
│   ├── docker/
│   │   ├── user/Dockerfile
│   │   └── message/Dockerfile
│   └── k8s/
│       ├── user-deployment.yaml
│       └── message-deployment.yaml
├── cmd/                            # 各微服务入口（每个服务一个目录）
│   ├── user/
│   │   └── main.go
│   ├── auth/
│   │   └── main.go
│   ├── message/
│   │   └── main.go
│   ├── group/
│   │   └── main.go
│   └── gateway/
│       └── main.go
├── deploy/                         # docker-compose 本地开发环境
│   └── docker-compose.yml
├── internal/                       # 私有应用代码（Go 编译器强制不可外部导入）
│   ├── user/                       # User 服务业务代码
│   │   ├── handler/                # HTTP 处理器（go-zero 生成）
│   │   ├── logic/                  # 核心业务逻辑
│   │   ├── model/                  # DB 数据模型（goctl model 生成）
│   │   ├── svc/                    # 服务上下文（依赖注入）
│   │   └── types/                  # 请求/响应结构体
│   ├── message/
│   ├── group/
│   └── gateway/
├── pkg/                            # 可复用公共库（多人协作共享）
│   ├── errorx/                     # 统一错误码定义
│   │   └── error.go
│   ├── middleware/                 # HTTP 中间件
│   │   ├── logger.go               # 请求日志
│   │   ├── recovery.go             # panic 恢复
│   │   └── auth.go                 # JWT 校验
│   ├── response/                   # 统一响应体封装
│   │   └── response.go
│   ├── jwt/                        # JWT 工具
│   │   └── jwt.go
│   ├── cache/                      # Redis 封装
│   │   └── redis.go
│   └── logger/                     # 日志封装（zap）
│       └── logger.go
├── scripts/                        # 构建、迁移、工具脚本
│   ├── gen.sh                      # goctl 代码生成脚本
│   └── migrate.sh                  # DB migration 脚本
├── test/                           # 集成测试
├── docs/                           # 项目文档、API 文档
├── go.work                         # Go workspace（Mono-Repo 多模块管理）
├── Makefile                        # 统一构建入口
├── .golangci.yml                   # Lint 配置
└── README.md
```