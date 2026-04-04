# 项目总览：easy-im

## 系统架构图
<img width="784" height="871" alt="image" src="https://github.com/user-attachments/assets/ec19cf74-883c-4236-a3c9-ea34109e8442" />


## [技术栈选型](https://github.com/hanzhengs61/easy-im/README_选型.md)
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
## WebSocket 长连接服务
```
客户端──WS握手──▶Gateway(HTTP升级)──▶ConnManager(连接池)
                                              │
                          ┌───────────────────┤
                          ▼                   ▼
                       读消息循环            心跳检测
                          │
                          ▼
                       消息路由器 ──▶ 找到目标连接 ──▶ 推送
                          │
                          ▼
                        Kafka ──▶ Message服务持久化
```
## 项目结构
```
easy-im/
├── .github/
│   └── workflows/
│       └── ci.yml                  # GitHub Actions 流水线
├── api/                            # Proto / API 定义文件（OpenAPI / .proto）
├── build/                          # Docker、K8s 构建配置
│   ├── docker/
├── cmd/                            # 各微服务入口（每个服务一个目录）
├── deploy/                         # docker-compose 本地开发环境
│   └── docker-compose.yml
├── internal/                       # 私有应用代码（Go 编译器强制不可外部导入）
├── pkg/                            # 可复用公共库（多人协作共享）
│   ├── errorx/                     # 统一错误码定义
│   ├── jwt/                        # JWT 工具
│   ├── logger/                     # 日志封装（zap）
│   ├── middleware/                 # HTTP 中间件
│   ├── protocol/                   # 协议（ws）
│   └── response/                   # 统一响应体封装
├── test/                           # 集成测试
├── go.work                         # Go workspace（Mono-Repo 多模块管理）
├── Makefile                        # 统一构建入口
├── .golangci.yml                   # Lint 配置
└── README.md
```