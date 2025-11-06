# ELK Stack Docker 部署模板

## 简介

完整的日志分析栈，包含：
- **Elasticsearch**: 搜索和分析引擎
- **Kibana**: 日志可视化和分析界面
- **Filebeat**: 轻量级日志采集器（可选）

## 环境变量配置

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| ELK_VERSION | ELK 版本 | 8.12.0 |
| ES_CONTAINER_NAME | Elasticsearch 容器名 | elasticsearch |
| ES_NODE_NAME | ES 节点名 | elasticsearch |
| ES_CLUSTER_NAME | ES 集群名 | elk-cluster |
| ES_HTTP_PORT | ES HTTP 端口 | 9200 |
| ES_TCP_PORT | ES TCP 端口 | 9300 |
| ES_JAVA_OPTS | ES JVM 参数 | -Xms512m -Xmx512m |
| ES_SECURITY_ENABLED | 是否启用安全认证 | false |
| KIBANA_CONTAINER_NAME | Kibana 容器名 | kibana |
| KIBANA_PORT | Kibana 端口 | 5601 |
| FILEBEAT_CONTAINER_NAME | Filebeat 容器名 | filebeat |
| DATA_DIR | 数据目录 | /data/elk |
| CONFIG_DIR | 配置目录 | ./config |
| LOG_DIR | 日志目录 | ./logs |

## 快速部署

### 方法1: 使用部署脚本

```bash
cd templates/
./deploy.sh elk elk-stack deploy
```

### 方法2: 手动部署

```bash
cd templates/elk

# 创建数据目录
mkdir -p /data/elk/{elasticsearch,filebeat}

# 启动服务
docker-compose -f versions/elk-stack-docker-compose.yml up -d
```

## 访问服务

- **Elasticsearch**: http://your-ip:9200
- **Kibana**: http://your-ip:5601

## 配置示例

### Filebeat 配置

创建 `config/filebeat.yml`:

```yaml
filebeat.inputs:
  - type: log
    enabled: true
    paths:
      - /var/log/app/*.log
      - /var/log/*.log

  - type: container
    enabled: true
    paths:
      - '/var/lib/docker/containers/*/*.log'

output.elasticsearch:
  hosts: ["elasticsearch:9200"]

setup.kibana:
  host: "kibana:5601"

setup.dashboards.enabled: true
setup.template.enabled: true
```

## 常用操作

### 查看 Elasticsearch 健康状态
```bash
curl http://localhost:9200/_cluster/health?pretty
```

### 查看索引
```bash
curl http://localhost:9200/_cat/indices?v
```

### 创建索引模板
```bash
curl -X PUT "localhost:9200/_index_template/logs-template" -H 'Content-Type: application/json' -d'
{
  "index_patterns": ["logs-*"],
  "template": {
    "settings": {
      "number_of_shards": 1,
      "number_of_replicas": 0
    }
  }
}'
```

### 查询日志
```bash
curl -X GET "localhost:9200/filebeat-*/_search?pretty" -H 'Content-Type: application/json' -d'
{
  "query": {
    "match_all": {}
  },
  "size": 10
}'
```

## Kibana 使用

### 创建 Index Pattern

1. 访问 Kibana: http://localhost:5601
2. 左侧菜单 -> Management -> Stack Management
3. Kibana -> Index Patterns
4. Create index pattern
5. 输入: `filebeat-*`
6. 选择时间字段: `@timestamp`

### 查看日志

1. 左侧菜单 -> Analytics -> Discover
2. 选择 index pattern: `filebeat-*`
3. 设置时间范围
4. 添加过滤条件

## 性能优化

### Elasticsearch 内存配置

```bash
# 4GB 内存服务器
ES_JAVA_OPTS="-Xms2g -Xmx2g" docker-compose up -d

# 8GB 内存服务器
ES_JAVA_OPTS="-Xms4g -Xmx4g" docker-compose up -d
```

### 索引生命周期管理 (ILM)

```bash
curl -X PUT "localhost:9200/_ilm/policy/logs-policy" -H 'Content-Type: application/json' -d'
{
  "policy": {
    "phases": {
      "hot": {
        "actions": {
          "rollover": {
            "max_size": "50GB",
            "max_age": "7d"
          }
        }
      },
      "delete": {
        "min_age": "30d",
        "actions": {
          "delete": {}
        }
      }
    }
  }
}'
```

### 清理旧数据

```bash
# 删除 7 天前的索引
curl -X DELETE "localhost:9200/filebeat-$(date -d '7 days ago' +%Y.%m.%d)"

# 批量删除
curator_cli delete_indices --filter_list '[{"filtertype":"age","source":"creation_date","direction":"older","unit":"days","unit_count":7}]'
```

## 监控

### 查看 Elasticsearch 统计
```bash
curl http://localhost:9200/_stats?pretty
```

### 查看节点信息
```bash
curl http://localhost:9200/_nodes/stats?pretty
```

## 常见问题

### Q1: Elasticsearch 启动失败

检查系统配置:
```bash
# 增加虚拟内存
sudo sysctl -w vm.max_map_count=262144

# 永久生效
echo "vm.max_map_count=262144" | sudo tee -a /etc/sysctl.conf
```

### Q2: 磁盘空间不足

```bash
# 查看磁盘使用
curl http://localhost:9200/_cat/allocation?v

# 清理旧索引
curl -X DELETE "localhost:9200/filebeat-2024.01.*"
```

### Q3: Kibana 无法连接 Elasticsearch

检查网络和 Elasticsearch 健康状态。

## 安全配置

### 启用安全认证

```yaml
environment:
  - xpack.security.enabled=true
  - ELASTIC_PASSWORD=your_password
```

### 创建用户
```bash
docker exec -it elasticsearch bin/elasticsearch-users useradd myuser -p mypassword -r superuser
```

## 注意事项

1. **首次启动**: Elasticsearch 需要 2-3 分钟初始化
2. **内存要求**: 至少 4GB 内存
3. **磁盘空间**: 预留足够空间存储日志数据
4. **数据备份**: 定期备份 Elasticsearch 数据目录
5. **索引管理**: 配置 ILM 自动清理旧数据
6. **安全加固**: 生产环境启用 X-Pack Security
