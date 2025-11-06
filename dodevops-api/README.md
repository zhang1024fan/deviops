# DevOps运维管理系统 API

![Go Version](https://img.shields.io/badge/Go-1.24-blue.svg)
![Gin](https://img.shields.io/badge/Gin-v1.11.0-green.svg)
![License](https://img.shields.io/badge/License-MIT-yellow.svg)

## 项目简介

DevOps运维管理系统是一个基于 Go + Gin 框架开发的企业级运维自动化平台，提供主机管理、配置中心、任务调度、K8s集群管理、监控告警等功能模块。

## 技术栈

### 后端框架
- **Gin** - Web框架
- **GORM** - ORM框架
- **JWT** - 身份认证
- **Swagger** - API文档
- **Logrus** - 日志管理
- **Robfig/Cron** - 定时任务调度

### 数据存储
- **MySQL** - 关系型数据库
- **Redis** - 缓存/消息队列

### 监控与追踪
- **Prometheus** - 监控指标采集
- **Pushgateway** - 指标推送

### 容器编排
- **Kubernetes** - K8s集群管理
- **client-go** - K8s客户端

### 云平台SDK
- **Tencent Cloud SDK** - 腾讯云资源管理
- **Baidu Cloud SDK** - 百度云资源管理

## 项目结构

```
dodevops-api/
├── api/                    # 业务逻辑层
│   ├── app/               # 应用管理模块（Jenkins、应用发布）
│   ├── cmdb/              # 配置管理数据库（主机、SQL、云资源）
│   ├── configcenter/      # 配置中心（密钥管理、账号认证、同步调度）
│   ├── dashboard/         # 数据看板
│   ├── k8s/               # K8s集群管理
│   ├── monitor/           # 监控告警（Prometheus、Agent心跳）
│   ├── system/            # 系统管理（用户、角色、菜单、部门）
│   └── task/              # 任务中心（Ansible、任务模板、定时任务）
├── common/                # 公共组件
│   ├── agent/            # Agent通信
│   ├── cache/            # 缓存层
│   ├── config/           # 配置管理
│   ├── constant/         # 常量定义
│   ├── result/           # 响应封装
│   ├── util/             # 工具函数（SSH、百度云服务等）
│   └── valid/            # 数据验证
├── docs/                  # Swagger文档
├── middleware/            # 中间件（JWT、日志、CORS、API描述）
├── pkg/                   # 基础包
│   ├── db/               # 数据库连接
│   ├── jwt/              # JWT工具
│   ├── log/              # 日志工具
│   └── redis/            # Redis连接
├── router/                # 路由注册
├── scheduler/             # 调度器管理
├── scripts/               # 运维脚本
├── task/                  # 任务执行目录
├── config.yaml            # 配置文件
└── main.go               # 程序入口

```

## 核心功能模块

### 1. CMDB（配置管理数据库）
- 主机资产管理
- 主机分组管理
- SSH连接管理
- SQL数据库管理
- 云资源管理（腾讯云、百度云）
- SQL操作日志审计

### 2. 配置中心
- 密钥管理（SSH密钥、API密钥）
- ECS云主机认证
- 账号权限管理
- 配置同步调度

### 3. 任务中心
- Ansible任务编排
- 任务模板管理
- 定时任务调度
- 任务执行监控
- WebSocket实时日志
- 任务队列系统

### 4. K8s集群管理
- 多集群管理
- Namespace管理
- Workload管理（Deployment、StatefulSet、DaemonSet）
- Service & Ingress管理
- ConfigMap & Secret管理
- 存储管理（PV、PVC、StorageClass）
- 节点管理
- 事件查看
- WebShell终端

### 5. 应用管理
- Jenkins集成
- 应用发布管理
- 构建历史

### 6. 监控告警
- Prometheus指标采集
- Pushgateway指标推送
- Agent心跳监控
- 系统资源监控

### 7. 系统管理
- 用户管理
- 角色权限
- 菜单管理
- 部门管理
- 岗位管理
- 操作日志
- 登录日志
- 验证码

### 8. 数据看板
- 运维数据统计
- 可视化展示

## 快速开始

### 环境要求
- Go 1.24+
- MySQL 5.7+
- Redis 5.0+

### 安装步骤

1. **克隆项目**
```bash
git clone <repository-url>
cd dodevops-api
```

2. **安装依赖**
```bash
go mod tidy
```

3. **配置文件**

编辑 `config.yaml` 文件：

```yaml
# 服务配置
server:
  address: 127.0.0.1:8000
  model: debug                # debug/release
  enableSwagger: true         # 是否启用Swagger文档

# 数据库配置
db:
  dialects: mysql
  host: localhost
  port: 3306
  db: gin-api
  username: root
  password: your_password
  charset: utf8
  maxIdle: 50
  maxOpen: 150

# Redis配置
redis:
  address: localhost:6379
  password: ""

# 监控配置
monitor:
  prometheus:
    url: "http://localhost:9090"
  pushgateway:
    url: "http://localhost:9091"
  agent:
    heartbeat_server_url: "http://127.0.0.1:8000/api/v1/monitor/agent/heartbeat"
    heartbeat_token: "agent-heartbeat-token-2024"

# 图片上传配置
imageSettings:
  uploadDir: ./upload/
  imageHost: https://www.deviops.cn

# 日志配置
log:
  path: ./log
  name: sys
  model: console              # console/file
```

4. **运行项目**

开发模式：
```bash
go run main.go
```

使用指定配置文件：
```bash
go run main.go -c /path/to/config.yaml
```

生产环境编译：
```bash
go build -o devops main.go
./devops
```

5. **访问服务**

- API服务：http://127.0.0.1:8000
- Swagger文档：http://127.0.0.1:8000/swagger/index.html

## API文档

项目集成了 Swagger 文档，启动服务后访问：

```
http://127.0.0.1:8000/swagger/index.html
```

### Swagger开关控制

通过配置文件控制 Swagger 文档的启用/禁用：

```yaml
server:
  enableSwagger: true  # true=启用, false=禁用
```

**建议：**
- 开发环境：`enableSwagger: true`
- 生产环境：`enableSwagger: false`

更新Swagger文档：
```bash
swag init
```

## 认证机制

API 使用 JWT (JSON Web Token) 进行身份认证。

### 获取Token

```bash
POST /api/v1/login
Content-Type: application/json

{
  "username": "admin",
  "password": "password",
  "code": "1234"
}
```

### 使用Token

在请求头中添加：
```
Authorization: Bearer <your-token>
```

## WebSocket支持

### 任务日志实时推送
```
ws://127.0.0.1:8000/ws/task/ansible/{id}/log/{work_id}
```

### K8s WebShell终端
```
ws://127.0.0.1:8000/api/v1/k8s/terminal
```

## 调度任务系统

项目内置强大的任务调度系统，支持：

- **Cron定时任务** - 基于Cron表达式的定时执行
- **任务队列** - 异步任务队列处理
- **并发控制** - 任务并发执行控制
- **失败重试** - 任务失败自动重试
- **任务监控** - 实时任务执行状态监控

## 主要依赖

| 依赖 | 版本 | 说明 |
|------|------|------|
| gin-gonic/gin | v1.11.0 | Web框架 |
| gorm.io/gorm | v1.31.0 | ORM框架 |
| go-redis/redis | v8.11.5 | Redis客户端 |
| robfig/cron | v3.0.1 | 定时任务 |
| k8s.io/client-go | v0.34.1 | K8s客户端 |
| prometheus/client_golang | v1.23.2 | Prometheus客户端 |
| swaggo/gin-swagger | v1.6.1 | Swagger文档 |
| gorilla/websocket | v1.5.4 | WebSocket |
| sirupsen/logrus | v1.9.3 | 日志库 |

## 开发规范

### 项目分层
```
Controller -> Service -> DAO -> Model
```

### 命名规范
- 文件名：小驼峰 `camelCase.go`
- 包名：小写 `package name`
- 结构体：大驼峰 `StructName`
- 接口：大驼峰 `InterfaceName`
- 函数：大驼峰 `FunctionName`（导出）/ 小驼峰 `functionName`（私有）

### Git提交规范
```
feat: 新功能
fix: 修复bug
docs: 文档更新
style: 代码格式
refactor: 重构
test: 测试
chore: 构建/工具
```

## 部署说明

### Docker部署
```bash
# 构建镜像
docker build -t dodevops-api:latest .

# 运行容器
docker run -d \
  -p 8000:8000 \
  -v $(pwd)/config.yaml:/app/config.yaml \
  -v $(pwd)/logs:/app/logs \
  --name dodevops-api \
  dodevops-api:latest
```

### 生产环境建议
1. 使用 `release` 模式
2. 禁用 Swagger 文档（`enableSwagger: false`）
3. 配置日志输出到文件（`log.model: file`）
4. 使用反向代理（Nginx/Caddy）
5. 启用HTTPS
6. 配置数据库连接池参数
7. 设置Redis密码

## 性能优化

- **数据库连接池**：根据并发量调整 `maxIdle` 和 `maxOpen`
- **Redis缓存**：热点数据缓存
- **日志分级**：生产环境使用 `info` 或 `warn` 级别
- **Gzip压缩**：API响应启用压缩
- **静态资源CDN**：图片等静态资源使用CDN

## 监控告警

### Prometheus指标
系统自动暴露以下指标：
- HTTP请求数量/延迟
- 数据库连接池状态
- Redis连接状态
- 任务执行状态
- Agent心跳状态

### 告警规则
在Prometheus中配置告警规则，通过AlertManager推送告警。

## 常见问题

### 1. 数据库连接失败
检查 `config.yaml` 中的数据库配置是否正确，确保MySQL服务已启动。

### 2. Redis连接失败
检查Redis服务状态，确认地址和密码配置正确。

### 3. Swagger文档404
确认 `enableSwagger: true` 并且已执行 `swag init` 生成文档。

### 4. JWT认证失败
检查Token是否过期，请求头格式是否正确。

### 5. K8s集群连接失败
确认kubeconfig配置正确，集群网络可达。

## 路线图

- [ ] 完善权限控制（RBAC）
- [ ] 支持多云平台（AWS、阿里云）
- [ ] 增加CI/CD流水线管理
- [ ] 增加审计日志分析
- [ ] 优化任务调度性能
- [ ] 支持分布式部署
- [ ] 增加API限流
- [ ] 支持OAuth2认证

## 贡献指南

欢迎提交 Issue 和 Pull Request！

1. Fork 本项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 提交 Pull Request

## 许可证

本项目采用 MIT 许可证。

## 联系方式

- 作者：xiaoRui
- 邮箱：your-email@example.com
- 项目地址：https://github.com/your-username/dodevops-api

---

**注意：** 生产环境部署前请务必修改默认密码和密钥配置！
