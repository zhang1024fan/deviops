# Prometheus 监控栈 Docker 部署模板

## 简介

完整的监控栈，包含：
- **Prometheus**: 时序数据库和监控系统
- **Grafana**: 数据可视化面板
- **Node Exporter**: 系统指标采集器

## 环境变量配置

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| CONTAINER_NAME | Prometheus 容器名 | prometheus |
| PROMETHEUS_VERSION | Prometheus 版本 | v2.47.0 |
| PROMETHEUS_PORT | Prometheus 端口 | 9090 |
| GRAFANA_CONTAINER_NAME | Grafana 容器名 | grafana |
| GRAFANA_VERSION | Grafana 版本 | 10.1.0 |
| GRAFANA_PORT | Grafana 端口 | 3000 |
| GRAFANA_ADMIN_USER | Grafana 管理员用户 | admin |
| GRAFANA_ADMIN_PASSWORD | Grafana 管理员密码 | admin123 |
| NODE_EXPORTER_PORT | Node Exporter 端口 | 9100 |
| RETENTION_TIME | 数据保留时间 | 15d |
| DATA_DIR | 数据目录 | /data/prometheus |
| CONFIG_DIR | 配置目录 | ./config |

## 快速部署

### 方法1: 使用部署脚本

```bash
# 使用统一部署脚本
cd templates/
./deploy.sh prometheus prometheus-stack deploy

# 指定 Grafana 密码
GRAFANA_ADMIN_PASSWORD=MySecurePass123 ./deploy.sh prometheus prometheus-stack deploy
```

### 方法2: 手动部署

```bash
cd templates/prometheus

# 复制配置文件
cp -r config-example config

# 启动服务
docker-compose -f versions/prometheus-stack-docker-compose.yml up -d
```

## 访问服务

- **Prometheus**: http://your-ip:9090
- **Grafana**: http://your-ip:3000 (admin/admin123)
- **Node Exporter**: http://your-ip:9100/metrics

## 配置文件

### Prometheus 配置

配置文件位于: `config/prometheus.yml`

```yaml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

scrape_configs:
  - job_name: 'prometheus'
    static_configs:
      - targets: ['localhost:9090']

  - job_name: 'node-exporter'
    static_configs:
      - targets: ['node-exporter:9100']
```

### Grafana 数据源自动配置

创建 `config/grafana/provisioning/datasources/prometheus.yml`:

```yaml
apiVersion: 1
datasources:
  - name: Prometheus
    type: prometheus
    access: proxy
    url: http://prometheus:9090
    isDefault: true
```

## 常用操作

### 查看日志
```bash
docker logs -f prometheus
docker logs -f grafana
```

### 重载 Prometheus 配置
```bash
curl -X POST http://localhost:9090/-/reload
```

### 导入 Grafana Dashboard

1. 访问 Grafana: http://localhost:3000
2. 登录 (admin/admin123)
3. 点击 "+" -> "Import"
4. 输入 Dashboard ID:
   - **1860**: Node Exporter Full
   - **3662**: Prometheus 2.0 Overview
   - **7362**: Docker and System Monitoring

### 添加监控目标

编辑 `config/prometheus.yml`:

```yaml
scrape_configs:
  - job_name: 'my-app'
    static_configs:
      - targets: ['192.168.1.100:8080']
```

重载配置:
```bash
curl -X POST http://localhost:9090/-/reload
```

## 监控 MySQL

```yaml
# 在 Prometheus 配置中添加
scrape_configs:
  - job_name: 'mysql'
    static_configs:
      - targets: ['mysql-exporter:9104']
```

启动 MySQL Exporter:
```bash
docker run -d --name mysql-exporter \
  --network monitoring \
  -e DATA_SOURCE_NAME="user:password@(mysql:3306)/" \
  -p 9104:9104 \
  prom/mysqld-exporter
```

## 监控 Redis

```bash
docker run -d --name redis-exporter \
  --network monitoring \
  -e REDIS_ADDR=redis:6379 \
  -p 9121:9121 \
  oliver006/redis_exporter
```

## 性能优化

### 数据保留策略

```bash
# 修改数据保留时间（默认15天）
RETENTION_TIME=30d docker-compose up -d
```

### 内存限制

编辑 docker-compose.yml 添加资源限制:

```yaml
services:
  prometheus:
    deploy:
      resources:
        limits:
          memory: 2G
        reservations:
          memory: 1G
```

## 告警配置

创建 `config/alert_rules.yml`:

```yaml
groups:
  - name: node_alerts
    interval: 30s
    rules:
      - alert: HighCPUUsage
        expr: 100 - (avg by(instance) (irate(node_cpu_seconds_total{mode="idle"}[5m])) * 100) > 80
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High CPU usage detected"
          description: "CPU usage is above 80% for 5 minutes"

      - alert: HighMemoryUsage
        expr: (node_memory_MemTotal_bytes - node_memory_MemAvailable_bytes) / node_memory_MemTotal_bytes * 100 > 90
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "High memory usage detected"
          description: "Memory usage is above 90%"
```

在 `config/prometheus.yml` 中引用:

```yaml
rule_files:
  - '/etc/prometheus/alert_rules.yml'
```

## 常见问题

### Q1: Grafana 无法连接 Prometheus

检查网络配置，确保两个容器在同一网络中。

### Q2: Node Exporter 数据不显示

检查 Prometheus 配置中的 targets 是否正确。

### Q3: 数据丢失

检查数据目录权限和磁盘空间。

## 注意事项

1. **首次启动**: Grafana 需要 1-2 分钟初始化
2. **默认密码**: 生产环境必须修改 Grafana 密码
3. **数据备份**: 定期备份 `/data/prometheus` 目录
4. **资源消耗**: Prometheus 和 Grafana 需要至少 2GB 内存
5. **网络配置**: 确保防火墙开放相应端口
