# Redis Docker 部署模板

## 支持版本

- Redis 6.2 (稳定版)
- Redis 7.0
- Redis 7.2 (最新版)

## 环境变量配置

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| CONTAINER_NAME | 容器名称 | redis62/redis70/redis72 |
| REDIS_PASSWORD | Redis密码 | redis123456 |
| REDIS_PORT | 映射端口 | 6379 |
| REDIS_MAXMEMORY | 最大内存 | 2gb |
| DATA_DIR | 数据目录 | /data/redis* |

## 快速部署

### 方法1：使用环境变量

```bash
REDIS_PASSWORD=MySecurePass123 docker-compose -f redis-7.2-docker-compose.yml up -d
```

### 方法2：使用 .env 文件

创建 `.env`:
```bash
CONTAINER_NAME=redis72
REDIS_PASSWORD=MyPass123
REDIS_PORT=6379
REDIS_MAXMEMORY=4gb
DATA_DIR=/data/redis72
```

启动：
```bash
docker-compose -f redis-7.2-docker-compose.yml up -d
```

## 数据持久化

- RDB: 定期快照
- AOF: 每秒同步 (appendfsync everysec)
- 数据目录: `${DATA_DIR}/data`

## 常用操作

### 连接 Redis
```bash
docker exec -it redis72 redis-cli -a redis123456
```

### 查看信息
```bash
docker exec redis72 redis-cli -a redis123456 INFO
```

### 查看日志
```bash
docker logs -f redis72
```

### 备份数据
```bash
docker exec redis72 redis-cli -a redis123456 BGSAVE
cp /data/redis72/data/dump.rdb /backup/
```

### 监控实时命令
```bash
docker exec -it redis72 redis-cli -a redis123456 MONITOR
```

## 性能优化

### 内存策略

- `allkeys-lru`: 所有key LRU淘汰 (默认)
- `volatile-lru`: 设置过期时间的key LRU淘汰
- `allkeys-lfu`: 所有key LFU淘汰 (Redis 4.0+)
- `noeviction`: 不淘汰，写满报错

### 持久化策略

生产环境建议：
```bash
# RDB + AOF 双重持久化
appendonly yes
appendfsync everysec
save 900 1
save 300 10
save 60 10000
```

高性能场景：
```bash
# 仅 AOF
appendonly yes
appendfsync everysec
save ""
```

## 集群部署

如需 Redis Cluster 或 Sentinel 高可用，请参考 `redis-cluster` 模板。

## 注意事项

1. **密码安全**: 生产环境必须设置强密码
2. **内存限制**: 根据实际情况调整 `maxmemory`
3. **持久化**: RDB+AOF 双重保障
4. **监控**: 建议配置 redis_exporter 监控
5. **备份**: 定期备份 RDB 和 AOF 文件
