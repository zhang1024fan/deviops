# DevOps运维管理系统 - Docker部署文档

## 概述

本文档提供基于Docker的部署方案。项目采用前后端分离架构：
- **后端镜像**：封装已编译的devops二进制文件
- **前端镜像**：封装nginx + dist静态文件
- **配置文件**：通过volume挂载config.yaml和devops.conf
- **其他组件**：MySQL、Redis、Prometheus、Pushgateway

## 架构说明

```
devops/
├── api/
│   ├── devops              # 已编译的后端二进制文件（需要提供）
│   ├── Dockerfile          # 后端镜像构建文件
│   └── config.yaml         # 后端配置文件（部署时挂载）
├── web/
│   ├── dist/               # 前端构建文件（需要提供）
│   ├── Dockerfile          # 前端镜像构建文件
│   └── devops.conf         # Nginx配置文件（部署时挂载）
└── docker/
    ├── docker-compose.yml  # 部署编排文件
    └── build-and-push.sh   # 镜像构建推送脚本
```

## 镜像说明

### 后端镜像 (devops-api)
- **基础镜像**：alpine:latest
- **包含内容**：devops二进制文件
- **运行时需要**：挂载config.yaml配置文件
- **镜像大小**：约20MB

### 前端镜像 (devops-web)
- **基础镜像**：nginx:1.24-alpine
- **包含内容**：dist静态文件
- **运行时需要**：挂载devops.conf配置文件
- **镜像大小**：约30MB

## 一、构建和推送镜像

### 1. 准备文件

```bash
# 确保以下文件存在
api/devops           # 后端二进制文件（已编译）
web/dist/            # 前端构建文件
api/config.yaml      # 后端配置文件
web/devops.conf      # Nginx配置文件
```

### 2. 配置镜像仓库

编辑 `docker/build-and-push.sh`，修改镜像仓库配置：

```bash
# 镜像仓库配置
REGISTRY="your-registry.com"              # 私有仓库地址
NAMESPACE="devops"                        # 命名空间
REGISTRY_USER="your-username"             # 仓库用户名（可选）
```

### 3. 构建镜像

```bash
cd docker

# 查看帮助
./build-and-push.sh --help

# 构建并推送v1.0.0版本
./build-and-push.sh v1.0.0

# 仅构建不推送
./build-and-push.sh v1.0.0 --build-only

# 仅构建后端镜像
./build-and-push.sh v1.0.0 --api-only

# 仅构建前端镜像
./build-and-push.sh v1.0.0 --web-only
```

### 4. 手动构建（可选）

```bash
# 构建后端镜像
cd api
docker build -t your-registry.com/devops/devops-api:v1.0.0 .

# 构建前端镜像
cd web
docker build -t your-registry.com/devops/devops-web:v1.0.0 .

# 推送镜像
docker push your-registry.com/devops/devops-api:v1.0.0
docker push your-registry.com/devops/devops-web:v1.0.0
```

## 二、部署系统

### 前置要求

- Docker Engine 20.10+
- Docker Compose 2.0+
- 至少 4GB 可用内存
- 至少 20GB 可用磁盘空间

### 1. 克隆项目

```bash
git clone https://github.com/zhang1024fan/deviops.git
cd deviops
```

### 2. 修改配置文件

#### 修改后端配置 `api/config.yaml`

```yaml
# 数据库配置（连接docker容器）
db:
  dialects: mysql
  host: mysql              # 改为容器服务名
  port: 3306
  db: devops
  username: root
  password: devops@2025
  charset: utf8
  maxIdle: 50
  maxOpen: 150

# redis配置（连接docker容器）
redis:
  address: redis:6379      # 改为容器服务名
  password: "zhangfan@123"

# 监控配置（连接docker容器）
monitor:
  prometheus:
    url: "http://prometheus:9090"
  pushgateway:
    url: "http://pushgateway:9091"
  agent:
    heartbeat_server_url: "http://devops-api:8000/api/v1/monitor/agent/heartbeat"
    heartbeat_token: "agent-heartbeat-token-2024"
```

#### 修改前端配置 `web/devops.conf`

```nginx
# API 代理配置
location /api/ {
    proxy_pass http://devops-api:8000;  # 改为容器服务名
    proxy_http_version 1.1;
    proxy_set_header Upgrade $http_upgrade;
    proxy_set_header Connection $connection_upgrade;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
    proxy_read_timeout 86400;
}
```

#### 修改Prometheus配置 `api/monitor/prometheus.yml`

```yaml
scrape_configs:
  - job_name: 'pushgateway'
    honor_labels: true
    static_configs:
      - targets: ['pushgateway:9091']  # 改为容器服务名
```

### 3. 配置docker-compose.yml

编辑 `docker/docker-compose.yml`，设置镜像地址：

```yaml
services:
  devops-api:
    image: your-registry.com/devops/devops-api:v1.0.0
    # ...

  nginx:
    image: your-registry.com/devops/devops-web:v1.0.0
    # ...
```

或者使用环境变量：

```bash
export DEVOPS_API_IMAGE=your-registry.com/devops/devops-api:v1.0.0
export DEVOPS_WEB_IMAGE=your-registry.com/devops/devops-web:v1.0.0
```

### 4. 启动服务

```bash
cd docker

# 拉取镜像
docker-compose pull

# 启动所有服务
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f
```

### 5. 访问系统

- **前端地址**: http://localhost 或 http://localhost:8080
- **默认账号**: admin
- **默认密码**: admin@2025

其他服务：
- API接口: http://localhost:8000
- Prometheus: http://localhost:9090
- Pushgateway: http://localhost:9091

## 三、常用管理命令

### 查看服务状态
```bash
docker-compose ps
```

### 查看日志
```bash
# 查看所有服务日志
docker-compose logs -f

# 查看特定服务
docker-compose logs -f devops-api
docker-compose logs -f nginx
```

### 停止服务
```bash
docker-compose stop
```

### 重启服务
```bash
docker-compose restart

# 重启特定服务
docker-compose restart devops-api
```

### 更新镜像
```bash
# 拉取新镜像
docker-compose pull

# 重启服务应用新镜像
docker-compose up -d
```

### 删除服务
```bash
# 停止并删除容器（保留数据）
docker-compose down

# 删除所有数据（包括数据卷）
docker-compose down -v
```

### 进入容器调试
```bash
# 进入API容器
docker-compose exec devops-api sh

# 进入Nginx容器
docker-compose exec nginx sh

# 进入MySQL容器
docker-compose exec mysql bash
```

### 数据备份
```bash
# 备份MySQL数据
docker-compose exec mysql mysqldump -uroot -pdevops@2025 devops > backup_$(date +%Y%m%d).sql

# 备份上传文件
docker cp devops-api:/app/upload ./upload_backup
```

## 四、数据持久化

Docker Compose自动创建以下数据卷：

| 数据卷 | 说明 | 映射路径 |
|--------|------|----------|
| mysql-data | MySQL数据库文件 | /var/lib/mysql |
| redis-data | Redis持久化文件 | /data |
| prometheus-data | Prometheus时序数据 | /prometheus |
| devops-logs | API服务日志 | /app/log |
| devops-upload | 用户上传文件 | /app/upload |

查看数据卷：
```bash
docker volume ls | grep devops
```

## 五、端口映射

| 服务 | 容器端口 | 主机端口 | 说明 |
|------|---------|---------|------|
| MySQL | 3306 | 3306 | 数据库 |
| Redis | 6379 | 6379 | 缓存 |
| DevOps API | 8000 | 8000 | 后端API |
| Nginx | 80/8080 | 80/8080 | 前端服务 |
| Prometheus | 9090 | 9090 | 监控 |
| Pushgateway | 9091 | 9091 | 指标推送 |

## 六、故障排查

### 1. 镜像拉取失败
```bash
# 检查镜像仓库配置
docker login your-registry.com

# 手动拉取镜像测试
docker pull your-registry.com/devops/devops-api:v1.0.0
```

### 2. 容器启动失败
```bash
# 查看容器日志
docker-compose logs devops-api
docker-compose logs nginx

# 检查容器状态
docker-compose ps
```

### 3. 配置文件未生效
```bash
# 检查配置文件挂载
docker-compose exec devops-api cat /app/config.yaml
docker-compose exec nginx cat /etc/nginx/conf.d/devops.conf

# 重启服务使配置生效
docker-compose restart devops-api
docker-compose restart nginx
```

### 4. 数据库连接失败
```bash
# 检查MySQL是否就绪
docker-compose exec mysql mysqladmin ping -h localhost -uroot -pdevops@2025

# 检查网络连通性
docker-compose exec devops-api ping mysql
```

### 5. 前端访问404
```bash
# 检查Nginx配置
docker-compose exec nginx nginx -t

# 检查静态文件
docker-compose exec nginx ls -la /var/www/html/dist
```

## 七、升级说明

### 升级后端
```bash
# 1. 构建新版本镜像
cd docker
./build-and-push.sh v1.1.0 --api-only

# 2. 更新docker-compose.yml中的镜像版本
# 或使用环境变量
export DEVOPS_API_IMAGE=your-registry.com/devops/devops-api:v1.1.0

# 3. 拉取并重启
docker-compose pull devops-api
docker-compose up -d devops-api
```

### 升级前端
```bash
# 1. 构建新版本镜像
cd docker
./build-and-push.sh v1.1.0 --web-only

# 2. 更新并重启
export DEVOPS_WEB_IMAGE=your-registry.com/devops/devops-web:v1.1.0
docker-compose pull nginx
docker-compose up -d nginx
```

### 更新配置文件
```bash
# 修改配置文件后重启服务
vi api/config.yaml
docker-compose restart devops-api

vi web/devops.conf
docker-compose restart nginx
```

## 八、生产环境建议

### 1. 安全配置
- 修改所有默认密码（MySQL、Redis、管理员）
- 配置SSL证书启用HTTPS
- 不要暴露MySQL、Redis端口到外网
- 使用防火墙限制访问

### 2. 性能优化
```yaml
# MySQL优化
mysql:
  command:
    - --max_connections=1000
    - --innodb_buffer_pool_size=2G

# Redis优化
redis:
  command: redis-server --requirepass xxx --maxmemory 1gb
```

### 3. 资源限制
```yaml
services:
  devops-api:
    deploy:
      resources:
        limits:
          cpus: '2'
          memory: 2G
```

### 4. 定期备份
- 设置定时任务备份MySQL数据库
- 备份上传文件目录
- 备份配置文件

### 5. 监控告警
- 配置Prometheus告警规则
- 监控容器资源使用情况
- 设置日志轮转

## 九、快速参考

### 完整部署流程
```bash
# 1. 准备文件
# 确保 api/devops 和 web/dist/ 存在

# 2. 修改镜像仓库配置
vi docker/build-and-push.sh

# 3. 构建并推送镜像
cd docker
./build-and-push.sh v1.0.0

# 4. 修改配置文件
vi api/config.yaml
vi web/devops.conf

# 5. 部署
cd docker
export DEVOPS_API_IMAGE=your-registry.com/devops/devops-api:v1.0.0
export DEVOPS_WEB_IMAGE=your-registry.com/devops/devops-web:v1.0.0
docker-compose up -d

# 6. 验证
docker-compose ps
docker-compose logs -f
```

## 技术支持

- **作者微信**: zf5391621
- **QQ群**: 925037324
- **邮箱**: zfwh1024@163.com
- **项目地址**: https://github.com/zhang1024fan/deviops.git
