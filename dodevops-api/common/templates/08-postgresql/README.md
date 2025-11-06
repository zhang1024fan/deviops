# PostgreSQL Docker 部署模板

## 支持版本

- PostgreSQL 13
- PostgreSQL 14
- PostgreSQL 16 (最新版)

## 环境变量配置

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| CONTAINER_NAME | 容器名称 | postgresql13/14/16 |
| POSTGRES_USER | 超级用户名 | postgres |
| POSTGRES_PASSWORD | 超级用户密码 | postgres123456 |
| POSTGRES_DB | 默认数据库 | postgres |
| POSTGRES_PORT | 映射端口 | 5432 |
| DATA_DIR | 数据目录 | /data/postgresql* |

## 快速部署

### 使用环境变量

```bash
POSTGRES_PASSWORD=MySecurePass123 docker-compose -f postgresql-16-docker-compose.yml up -d
```

### 使用 .env 文件

```bash
CONTAINER_NAME=postgresql16
POSTGRES_USER=admin
POSTGRES_PASSWORD=Admin123!
POSTGRES_DB=myapp
POSTGRES_PORT=5432
DATA_DIR=/data/postgresql16
```

## 常用操作

### 连接数据库
```bash
docker exec -it postgresql16 psql -U postgres
```

### 执行 SQL 文件
```bash
docker exec -i postgresql16 psql -U postgres -d mydb < init.sql
```

### 备份数据库
```bash
# 单个数据库
docker exec postgresql16 pg_dump -U postgres mydb > backup.sql

# 所有数据库
docker exec postgresql16 pg_dumpall -U postgres > backup_all.sql
```

### 恢复数据库
```bash
docker exec -i postgresql16 psql -U postgres mydb < backup.sql
```

### 查看日志
```bash
docker logs -f postgresql16
```

### 创建新用户和数据库
```bash
docker exec -it postgresql16 psql -U postgres -c "CREATE USER appuser WITH PASSWORD 'apppass123';"
docker exec -it postgresql16 psql -U postgres -c "CREATE DATABASE appdb OWNER appuser;"
docker exec -it postgresql16 psql -U postgres -c "GRANT ALL PRIVILEGES ON DATABASE appdb TO appuser;"
```

## 性能优化

### 根据机器配置调整参数

**2GB 内存服务器**:
```yaml
shared_buffers=512MB
effective_cache_size=1536MB
maintenance_work_mem=128MB
work_mem=2MB
```

**8GB 内存服务器**:
```yaml
shared_buffers=2GB
effective_cache_size=6GB
maintenance_work_mem=512MB
work_mem=8MB
```

### 启用慢查询日志

在 docker-compose 添加：
```yaml
- "-c"
- "log_min_duration_statement=1000"  # 记录超过1秒的查询
- "-c"
- "log_line_prefix=%t [%p]: [%l-1] user=%u,db=%d,app=%a,client=%h"
```

## 高可用部署

如需主从复制或 Patroni 高可用，请参考 `postgresql-ha` 模板。

## 监控

推荐使用 postgres_exporter 配合 Prometheus + Grafana 监控。

```bash
docker run -d --name postgres-exporter \
  --network postgresql_network \
  -e DATA_SOURCE_NAME="postgresql://postgres:password@postgresql16:5432/postgres?sslmode=disable" \
  -p 9187:9187 \
  prometheuscommunity/postgres-exporter
```

## 注意事项

1. **首次启动**: 初始化需要几秒钟
2. **密码安全**: 生产环境必须使用强密码
3. **数据备份**: 定期备份 PGDATA 目录
4. **版本升级**: 大版本升级需要 pg_upgrade 或逻辑导出导入
5. **连接池**: 生产环境建议使用 PgBouncer
6. **时区设置**: 已默认设置为 Asia/Shanghai
