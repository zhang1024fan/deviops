# N9E (夜莺) 监控系统 Docker 部署模板

## 简介

夜莺 (Nightingale) 是一款开源的企业级监控系统，由滴滴开源。

**特点**:
- 🎯 企业级监控解决方案
- 📊 支持 Prometheus 生态
- 🔔 强大的告警能力
- 📈 内置丰富的监控仪表板
- 🌐 多租户支持
- 🔌 插件化架构

## 架构组件

| 组件 | 说明 | 端口 |
|------|------|------|
| **Nightingale** | 核心服务 | 17000(HTTP), 20090(RPC) |
| **MySQL** | 存储元数据 | 3306 |
| **Redis** | 缓存 | 6379 |
| **VictoriaMetrics** | 时序数据库 | 8428 |
| **Categraf** | 数据采集器 | - |

## 环境变量配置

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| N9E_VERSION | 夜莺版本 | latest |
| MYSQL_ROOT_PASSWORD | MySQL root 密码 | 1234 |
| MYSQL_PORT | MySQL 端口 | 3306 |
| REDIS_PORT | Redis 端口 | 6379 |
| VM_PORT | VictoriaMetrics 端口 | 8428 |
| VM_RETENTION | 数据保留时间 | 30d |
| N9E_HTTP_PORT | 夜莺 HTTP 端口 | 17000 |
| N9E_RPC_PORT | 夜莺 RPC 端口 | 20090 |
| DATA_DIR | 数据目录 | /data/n9e |
| CONFIG_DIR | 配置目录 | ./config |

## 快速部署

### 方法1: 使用部署脚本

```bash
cd templates/

# 基础部署
./deploy.sh n9e n9e-stack deploy

# 自定义密码
MYSQL_ROOT_PASSWORD=MySecurePass123 \
./deploy.sh n9e n9e-stack deploy
```

### 方法2: 手动部署

```bash
cd templates/n9e

# 创建数据目录
mkdir -p /data/n9e/{mysql,redis,victoriametrics}

# 准备配置文件（从现有的 n9e 目录复制）
cp -r /path/to/existing/n9e/etc-nightingale config/
cp -r /path/to/existing/n9e/etc-categraf config/

# 启动服务
docker-compose -f versions/n9e-stack-docker-compose.yml up -d
```

## 访问服务

### 夜莺 Web 界面
- URL: http://your-ip:17000
- 默认账号: `root`
- 默认密码: `root.2020`

### VictoriaMetrics
- URL: http://your-ip:8428

## 初始化配置

### 1. 首次登录

访问 http://localhost:17000

1. 使用默认账号登录: root / root.2020
2. **立即修改默认密码**
3. 配置组织和团队

### 2. 配置数据源

1. 左侧菜单 -> 配置中心 -> 数据源管理
2. 添加 VictoriaMetrics:
   - 名称: VictoriaMetrics
   - 类型: VictoriaMetrics
   - URL: http://victoriametrics:8428

### 3. 配置采集器

Categraf 会自动采集以下指标：
- CPU
- 内存
- 磁盘
- 网络
- 进程

### 4. 导入仪表板

夜莺内置了丰富的监控仪表板，直接使用即可。

## 监控配置

### 添加监控对象

1. 左侧菜单 -> 监控对象
2. 点击 "添加监控对象"
3. 填写信息（IP、标签等）

### 创建告警规则

1. 左侧菜单 -> 告警管理 -> 告警规则
2. 点击 "创建告警规则"
3. 配置规则:

```promql
# CPU 使用率超过 80%
100 - (avg by(instance) (rate(node_cpu_seconds_total{mode="idle"}[5m])) * 100) > 80

# 内存使用率超过 90%
(1 - (node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes)) * 100 > 90

# 磁盘使用率超过 85%
(1 - (node_filesystem_avail_bytes{fstype!~"tmpfs|overlay"} / node_filesystem_size_bytes{fstype!~"tmpfs|overlay"})) * 100 > 85
```

### 配置告警通知

1. 左侧菜单 -> 告警管理 -> 通知设置
2. 配置通知渠道:
   - 企业微信
   - 钉钉
   - 飞书
   - 邮件
   - Webhook

## Categraf 配置

### 配置文件位置

`config/etc-categraf/config.toml`

### 添加监控目标

#### 监控 MySQL

编辑 `config/etc-categraf/input.mysql/mysql.toml`:

```toml
[[instances]]
address = "root:password@tcp(mysql:3306)/"
```

#### 监控 Redis

编辑 `config/etc-categraf/input.redis/redis.toml`:

```toml
[[instances]]
address = "redis:6379"
password = ""
```

#### 监控 HTTP 接口

编辑 `config/etc-categraf/input.http_response/http_response.toml`:

```toml
[[instances]]
targets = [
    "http://localhost:8080",
    "https://api.example.com"
]
```

## 常用操作

### 查看日志
```bash
docker logs -f nightingale
docker logs -f categraf
```

### 重启服务
```bash
./deploy.sh n9e n9e-stack restart
```

### 备份配置
```bash
tar -czf n9e-config-backup.tar.gz config/
```

### 备份数据
```bash
# 备份 MySQL
docker exec n9e-mysql mysqldump -uroot -p1234 n9e > n9e-db-backup.sql

# 备份 VictoriaMetrics
tar -czf vm-data-backup.tar.gz /data/n9e/victoriametrics/
```

## 性能优化

### VictoriaMetrics 优化

```bash
# 增加数据保留时间
VM_RETENTION=90d docker-compose up -d victoriametrics

# 调整内存限制
docker update --memory=4g victoriametrics
```

### MySQL 优化

编辑 `config/my.cnf`:

```ini
[mysqld]
innodb_buffer_pool_size=1G
max_connections=500
```

## 高可用部署

### 集群模式

夜莺支持集群部署：

1. 部署多个 Nightingale 实例
2. 共享 MySQL 和 Redis
3. 使用 Nginx 负载均衡

### 监控高可用

1. VictoriaMetrics 集群
2. MySQL 主从复制
3. Redis 哨兵模式

## 集成 Grafana

### 配置数据源

在 Grafana 中添加 Prometheus 数据源：
- URL: http://victoriametrics:8428

### 导入仪表板

Grafana 官方仪表板兼容 N9E。

## 常见问题

### Q1: 服务无法启动

检查依赖服务是否正常：
```bash
docker ps
docker logs nightingale
```

### Q2: 数据未采集

检查 Categraf 状态：
```bash
docker logs categraf
```

### Q3: 告警未发送

检查通知渠道配置和网络连通性。

### Q4: Web 界面无法访问

检查防火墙和端口映射：
```bash
netstat -tulpn | grep 17000
```

## 监控指标

### 系统指标
- CPU 使用率
- 内存使用率
- 磁盘 IO
- 网络流量
- 进程状态

### 应用指标
- HTTP 请求量
- 响应时间
- 错误率
- 数据库连接数

### 业务指标
- 自定义指标采集
- Pushgateway 支持

## 安全建议

1. **修改默认密码**: 首次登录后立即修改
2. **启用 HTTPS**: 配置 SSL 证书
3. **网络隔离**: 使用防火墙限制访问
4. **定期备份**: 自动化备份配置和数据
5. **权限管理**: 合理分配用户权限

## 注意事项

1. **首次启动**: 需要 3-5 分钟初始化数据库
2. **资源要求**: 建议至少 4GB 内存
3. **数据保留**: 定期清理旧数据
4. **配置文件**: 需要从现有 N9E 部署复制配置文件
5. **网络模式**: Categraf 需要访问宿主机指标
6. **时区设置**: 已统一设置为 Asia/Shanghai

## 参考资源

- 官方文档: https://n9e.github.io/
- GitHub: https://github.com/ccfos/nightingale
- 社区论坛: https://answer.flashcat.cloud/
