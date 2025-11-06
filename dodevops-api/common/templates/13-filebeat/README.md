# Filebeat Docker 部署模板

## 简介

Filebeat 是轻量级的日志采集器，用于转发和集中日志数据。

**特点**:
- 轻量级，占用资源少
- 支持多种日志格式
- 内置 Docker 日志采集
- 自动重试机制
- 与 Elasticsearch 深度集成

## 环境变量配置

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| FILEBEAT_VERSION | Filebeat 版本 | 8.12.0 |
| PROJECT_NAME | 项目名称（用于多实例） | default |
| ELASTICSEARCH_HOSTS | Elasticsearch 地址 | http://127.0.0.1:9200 |
| KIBANA_HOST | Kibana 地址 | http://127.0.0.1:5601 |
| ELASTICSEARCH_USERNAME | ES 用户名 | 空 |
| ELASTICSEARCH_PASSWORD | ES 密码 | 空 |
| NETWORK_MODE | 网络模式 | host |
| DATA_DIR | 数据目录 | /data/filebeat-default |
| CONFIG_DIR | 配置目录 | ./config |
| LOG_DIR | 应用日志目录 | ./logs |

## 快速部署

### 方法1: 使用部署脚本

```bash
cd templates/

# 部署到本地 Elasticsearch
ELASTICSEARCH_HOSTS=http://localhost:9200 \
./deploy.sh filebeat filebeat-standalone deploy

# 部署到远程 Elasticsearch
ELASTICSEARCH_HOSTS=http://192.168.1.100:9200 \
./deploy.sh filebeat filebeat-standalone deploy
```

### 方法2: 多实例部署

```bash
# 实例1 - 收集应用A的日志
COMPOSE_PROJECT_NAME=app-a \
ELASTICSEARCH_HOSTS=http://localhost:9200 \
LOG_DIR=/var/log/app-a \
./deploy.sh filebeat filebeat-standalone deploy -p app-a

# 实例2 - 收集应用B的日志
COMPOSE_PROJECT_NAME=app-b \
ELASTICSEARCH_HOSTS=http://localhost:9200 \
LOG_DIR=/var/log/app-b \
./deploy.sh filebeat filebeat-standalone deploy -p app-b
```

## 配置说明

### 基础配置

配置文件: `config/filebeat.yml`

主要配置项：
- `filebeat.inputs`: 定义日志输入源
- `output.elasticsearch`: Elasticsearch 输出配置
- `setup.kibana`: Kibana 集成配置
- `processors`: 数据处理器

### 采集应用日志

```yaml
filebeat.inputs:
  - type: log
    enabled: true
    paths:
      - /var/log/app/*.log
    fields:
      app: myapp
      env: production
```

### 采集 Docker 日志

```yaml
filebeat.inputs:
  - type: container
    enabled: true
    paths:
      - '/var/lib/docker/containers/*/*.log'
```

### 采集 JSON 日志

```yaml
filebeat.inputs:
  - type: log
    enabled: true
    paths:
      - /var/log/app/*.json
    json.keys_under_root: true
    json.add_error_key: true
```

### 多行日志合并

```yaml
filebeat.inputs:
  - type: log
    enabled: true
    paths:
      - /var/log/app/*.log
    multiline.pattern: '^\d{4}-\d{2}-\d{2}'
    multiline.negate: true
    multiline.match: after
```

## 常用操作

### 查看日志
```bash
docker logs -f filebeat-default
```

### 测试配置
```bash
docker exec filebeat-default filebeat test config
```

### 测试输出
```bash
docker exec filebeat-default filebeat test output
```

### 重启服务
```bash
./deploy.sh filebeat filebeat-standalone restart -p default
```

### 查看采集状态
```bash
# 在 Kibana 中查看
# Management -> Stack Monitoring -> Beats
```

## Kibana 集成

### 自动加载仪表板

Filebeat 会自动加载内置的仪表板到 Kibana：

```bash
# 手动加载仪表板
docker exec filebeat-default filebeat setup --dashboards
```

### 查看仪表板

1. 访问 Kibana
2. 左侧菜单 -> Analytics -> Dashboard
3. 搜索 "Filebeat"

## 性能优化

### 批量发送

```yaml
output.elasticsearch:
  hosts: ["localhost:9200"]
  bulk_max_size: 2048
  worker: 2
```

### 压缩传输

```yaml
output.elasticsearch:
  hosts: ["localhost:9200"]
  compression_level: 3
```

### 内存限制

```yaml
# docker-compose.yml
services:
  filebeat:
    deploy:
      resources:
        limits:
          memory: 256M
```

## 常见问题

### Q1: Filebeat 无法连接 Elasticsearch

检查 Elasticsearch 地址和网络连通性：
```bash
curl http://localhost:9200
```

### Q2: 日志未采集

检查文件路径和权限：
```bash
docker exec filebeat-default ls -la /var/log/app/
```

### Q3: 数据重复采集

检查 registry 文件：
```bash
docker exec filebeat-default cat /usr/share/filebeat/data/registry/filebeat/log.json
```

## 高级用法

### 过滤日志

```yaml
processors:
  - drop_event:
      when:
        regexp:
          message: "DEBUG"
```

### 添加字段

```yaml
processors:
  - add_fields:
      target: ''
      fields:
        env: production
        region: us-east-1
```

### 解析 JSON

```yaml
processors:
  - decode_json_fields:
      fields: ["message"]
      target: ""
      overwrite_keys: true
```

## 监控告警

### 配置告警

在 Elasticsearch 中创建 Watcher:

```json
{
  "trigger": {
    "schedule": {
      "interval": "5m"
    }
  },
  "input": {
    "search": {
      "request": {
        "indices": ["filebeat-*"],
        "body": {
          "query": {
            "match": {
              "log.level": "ERROR"
            }
          }
        }
      }
    }
  },
  "condition": {
    "compare": {
      "ctx.payload.hits.total": {
        "gt": 10
      }
    }
  },
  "actions": {
    "notify": {
      "email": {
        "to": "admin@example.com",
        "subject": "High error rate detected"
      }
    }
  }
}
```

## 注意事项

1. **文件权限**: Filebeat 需要读取日志文件的权限
2. **磁盘空间**: Registry 文件会记录采集进度
3. **网络连接**: 确保能连接到 Elasticsearch
4. **日志轮转**: 支持日志文件轮转，自动检测新文件
5. **多实例**: 使用不同的项目名避免冲突
6. **资源消耗**: 监控 CPU 和内存使用情况
