# DevOps运维管理系统 Docker 一键部署

## 部署说明

本目录提供了 DevOps 运维管理系统的 Docker Compose 一键部署方案，所有配置文件和数据都已持久化到本地。

## 目录结构

```
docker/
├── docker-compose.yml          # Docker Compose 编排文件
├── .env                        # 环境变量配置
├── api/                        # 后端配置
│   ├── config.yaml            # API 配置文件
│   ├── logs/                  # 日志目录(持久化)
│   └── upload/                # 上传文件目录(持久化)
├── web/                       # 前端配置
│   └── devops.conf           # Nginx 配置文件
├── mysql/                     # MySQL 配置
│   ├── devops.sql            # 初始化 SQL 脚本
│   ├── conf.d/               # MySQL 配置文件
│   │   └── my.cnf
│   └── data/                 # 数据目录(持久化)
├── redis/                     # Redis 配置
│   ├── redis.conf            # Redis 配置文件
│   └── data/                 # 数据目录(持久化)
├── prometheus/                # Prometheus 配置
│   ├── prometheus.yml        # Prometheus 配置文件
│   └── data/                 # 数据目录(持久化)
└── pushgateway/               # Pushgateway 配置
    └── data/                 # 数据目录(持久化)
```

## 服务列表

| 服务名 | 容器名 | 端口映射 | 说明 |
|--------|--------|----------|------|
| mysql | devops-mysql | 3306:3306 | MySQL 8.0.33 数据库 |
| redis | devops-redis | 6379:6379 | Redis 7.0 缓存 |
| pushgateway | devops-pushgateway | 9091:9091 | Pushgateway 指标推送 |
| prometheus | devops-prometheus | 9090:9090 | Prometheus 监控 |
| devops-api | devops-api | 8000:8000 | DevOps API 后端服务 |
| devops-web | devops-web | 8080:80 | DevOps Web 前端服务 |

## 部署前准备
```bash
docker pull crpi-aj3vgoxp9kzh2jx4.cn-hangzhou.personal.cr.aliyuncs.com/zhangfan_k8s/deviops-api:v1.0         
docker pull crpi-aj3vgoxp9kzh2jx4.cn-hangzhou.personal.cr.aliyuncs.com/zhangfan_k8s/deviops-web:v1.0         
docker pull crpi-aj3vgoxp9kzh2jx4.cn-hangzhou.personal.cr.aliyuncs.com/zhangfan_k8s/pushgateway:v1.9.0       
docker pull crpi-aj3vgoxp9kzh2jx4.cn-hangzhou.personal.cr.aliyuncs.com/zhangfan_k8s/redis:7.0-alpine  
docker pull crpi-aj3vgoxp9kzh2jx4.cn-hangzhou.personal.cr.aliyuncs.com/zhangfan_k8s/prometheus:v2.47.0      
docker pull crpi-aj3vgoxp9kzh2jx4.cn-hangzhou.personal.cr.aliyuncs.com/zhangfan_k8s/mysql:8.0.33
```
### 1. 修改环境变量

编辑 `.env` 文件,修改以下配置:

```bash
# 修改外网访问地址(必改)
IMAGE_HOST=http://your-domain.com  # 或 http://192.168.1.100:8080

# 如果需要,可以修改数据库密码
MYSQL_ROOT_PASSWORD=your-password
REDIS_PASSWORD=your-password
```

### 2. 修改 API 配置

编辑 `api/config.yaml`,确保 `imageHost` 与 `.env` 中的 `IMAGE_HOST` 一致:

```yaml
imageSettings:
  uploadDir: ./upload/
  imageHost: http://your-domain.com  # 修改为实际地址
```

## 快速启动

### 1. 启动所有服务

```bash
cd /Users/apple/Desktop/git/deviops/docker
docker-compose up -d
```

### 2. 查看服务状态

```bash
docker-compose ps
```

### 3. 查看服务日志

```bash
# 查看所有服务日志
docker-compose logs -f

# 查看指定服务日志
docker-compose logs -f devops-api
docker-compose logs -f devops-web
docker-compose logs -f mysql
```

### 4. 访问系统

- **Web 前端**: http://localhost:8080
- **默认账号**: admin / admin@2025
- **Prometheus**: http://localhost:9090
- **Pushgateway**: http://localhost:9091

## 服务管理

### 停止服务

```bash
docker-compose stop
```

### 重启服务

```bash
docker-compose restart
```

### 停止并删除容器

```bash
docker-compose down
```

### 停止并删除容器及数据卷(危险操作)

```bash
docker-compose down -v
```

### 重新构建并启动

```bash
docker-compose up -d --build
```

### 单独重启某个服务

```bash
docker-compose restart devops-api
```

## 数据持久化

以下目录数据已持久化到本地,停止容器后数据不会丢失:

- **MySQL数据**: `./mysql/data/`
- **Redis数据**: `./redis/data/`
- **Prometheus数据**: `./prometheus/data/`
- **Pushgateway数据**: `./pushgateway/data/`
- **API日志**: `./api/logs/`
- **上传文件**: `./api/upload/`

## 健康检查

所有服务都配置了健康检查,启动时会自动等待依赖服务就绪:

- MySQL: 检查数据库连接
- Redis: 检查 Redis 连接
- Prometheus: 检查 HTTP 健康端点
- Pushgateway: 检查 HTTP 健康端点
- API: 检查 HTTP 健康端点(启动后30秒开始检查)
- Web: 检查 Nginx HTTP 服务

## 网络配置

所有服务运行在独立的 Docker 网络 `devops-network` 中,使用子网 `172.20.0.0/16`。

服务之间通过容器名通信:
- API 连接 MySQL: `mysql:3306`
- API 连接 Redis: `redis:6379`
- Prometheus 采集 Pushgateway: `pushgateway:9091`
- Web 代理 API: `devops-api:8000`

## 故障排查

### 1. 容器启动失败

```bash
# 查看容器日志
docker-compose logs <service-name>

# 查看容器状态
docker-compose ps
```

### 2. 数据库连接失败

```bash
# 进入 MySQL 容器检查
docker-compose exec mysql mysql -uroot -pdevops@2025

# 检查数据库是否初始化
docker-compose exec mysql mysql -uroot -pdevops@2025 -e "SHOW DATABASES;"
```

### 3. API 无法连接数据库

```bash
# 检查网络连通性
docker-compose exec devops-api ping mysql

# 检查 MySQL 是否健康
docker-compose exec mysql mysqladmin ping -h localhost -uroot -pdevops@2025
```

### 4. 重置数据库

```bash
# 停止服务
docker-compose down

# 删除 MySQL 数据
rm -rf ./mysql/data/*

# 重新启动(会自动初始化数据库)
docker-compose up -d
```

## 端口冲突

如果本机端口被占用,可以修改 `docker-compose.yml` 中的端口映射:

```yaml
ports:
  - "8080:80"  # 修改为 "18080:80"
```

## 性能优化

### MySQL 内存优化

编辑 `mysql/conf.d/my.cnf`,根据服务器内存调整:

```ini
innodb_buffer_pool_size=512M  # 建议设置为物理内存的 50-80%
```

### Redis 内存优化

编辑 `redis/redis.conf`:

```
maxmemory 1gb  # 根据实际情况调整
```

## 备份与恢复

### 备份数据库

```bash
docker-compose exec mysql mysqldump -uroot -pdevops@2025 devops > backup_$(date +%Y%m%d).sql
```

### 恢复数据库

```bash
docker-compose exec -T mysql mysql -uroot -pdevops@2025 devops < backup.sql
```

## 升级说明

### 升级镜像版本

1. 修改 `docker-compose.yml` 中的镜像版本
2. 拉取新镜像: `docker-compose pull`
3. 重启服务: `docker-compose up -d`

## 安全建议

1. **修改默认密码**: 生产环境务必修改 MySQL 和 Redis 密码
2. **配置防火墙**: 限制外网访问数据库端口(3306, 6379)
3. **HTTPS配置**: 生产环境建议配置 HTTPS 证书
4. **定期备份**: 建议定期备份数据库和上传文件

## 生产环境部署建议

1. 使用外部 MySQL 和 Redis(提高可用性)
2. 配置 Nginx 反向代理和 SSL 证书
3. 配置日志轮转,避免日志文件过大
4. 设置资源限制,避免单个容器占用过多资源
5. 配置监控告警,及时发现问题

## 技术支持

- 项目地址: https://github.com/zhang1024fan/deviops.git
- 微信技术交流: zf5391621
- QQ技术交流: 545118130
- 建议邮箱: zfwh1024@163.com
