# MySQL Docker 部署模板

## 支持版本

- MySQL 5.7
- MySQL 8.0
- MySQL 8.4 (最新版)

## 环境变量配置

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| CONTAINER_NAME | 容器名称 | mysql57/mysql80/mysql84 |
| MYSQL_ROOT_PASSWORD | root密码 | root123456 |
| MYSQL_DATABASE | 初始数据库 | 空 |
| MYSQL_USER | 普通用户名 | 空 |
| MYSQL_PASSWORD | 普通用户密码 | 空 |
| MYSQL_PORT | 映射端口 | 3306 |
| DATA_DIR | 数据目录 | /data/mysql* |

## 快速部署

### 方法1：使用环境变量文件

创建 `.env` 文件：
```bash
CONTAINER_NAME=mysql80
MYSQL_ROOT_PASSWORD=MySecurePass123
MYSQL_DATABASE=myapp
MYSQL_USER=appuser
MYSQL_PASSWORD=AppPass123
MYSQL_PORT=3306
DATA_DIR=/data/mysql80
```

启动：
```bash
docker-compose -f mysql-8.0-docker-compose.yml up -d
```

### 方法2：命令行传参

```bash
MYSQL_ROOT_PASSWORD=MyPass123 MYSQL_PORT=3307 docker-compose -f mysql-8.0-docker-compose.yml up -d
```

## 数据持久化

数据存储在以下目录：
- `${DATA_DIR}/data` - 数据库文件
- `${DATA_DIR}/conf` - 自定义配置
- `${DATA_DIR}/logs` - 日志文件

## 自定义配置

将配置文件放入 `${DATA_DIR}/conf/` 目录，示例：
```bash
mkdir -p /data/mysql80/conf
cp config/my.cnf /data/mysql80/conf/
```

## 常用操作

### 查看日志
```bash
docker logs -f mysql80
```

### 进入容器
```bash
docker exec -it mysql80 mysql -uroot -p
```

### 备份数据库
```bash
docker exec mysql80 mysqldump -uroot -proot123456 --all-databases > backup.sql
```

### 恢复数据库
```bash
docker exec -i mysql80 mysql -uroot -proot123456 < backup.sql
```

## 注意事项

1. **首次启动**: 初始化需要几分钟，等待 healthcheck 通过
2. **密码安全**: 生产环境必须修改默认密码
3. **数据备份**: 定期备份 `/data` 目录
4. **版本升级**: 大版本升级需要导出导入数据
5. **性能优化**: 根据实际情况调整 `my.cnf` 配置
