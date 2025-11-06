# Docker Compose 运维工具箱

一键部署常用开发和运维服务的 Docker Compose 模板库。

## 📦 支持的服务

### 数据库服务

| 服务 | 支持版本 | 说明 |
|------|---------|------|
| **MySQL** | 5.7, 8.0, 8.4 | 关系型数据库 |
| **PostgreSQL** | 13, 14, 16 | 关系型数据库 |
| **Redis** | 6.2, 7.0, 7.2 | 内存数据库/缓存 |

### CI/CD 工具

| 服务 | 支持版本 | 说明 |
|------|---------|------|
| **Jenkins** | LTS, Latest | 持续集成/持续部署 |
| **GitLab** | CE, EE | 代码托管平台 |

### 监控与日志

| 服务 | 支持版本 | 说明 |
|------|---------|------|
| **Prometheus** | - | 监控系统 (已有) |
| **Grafana** | Latest | 可视化面板 |
| **Elasticsearch** | 8.x | 搜索引擎 |
| **Loki** | Latest | 日志聚合系统 |
| **Fluentd** | Latest | 日志收集器 |
| **ELK Stack** | - | 日志分析平台 (已有) |
| **夜莺 (N9E)** | - | 监控系统 (已有) |

### 运维工具

| 服务 | 支持版本 | 说明 |
|------|---------|------|
| **JumpServer** | Latest | 开源堡垒机 |

### 开发环境

| 服务 | 支持版本 | 说明 |
|------|---------|------|
| **Node.js** | 18, 20 | JavaScript 运行时 |
| **Java** | 8, 11, 17 | Java 运行环境 |
| **Golang** | 1.20, 1.21 | Go 运行环境 |

## 🚀 快速开始

### 1. 目录结构

```
templates/
├── deploy.sh                    # 统一部署脚本
├── README.md                    # 总体说明文档
├── mysql/
│   ├── versions/
│   │   ├── mysql-5.7-docker-compose.yml
│   │   ├── mysql-8.0-docker-compose.yml
│   │   └── mysql-8.4-docker-compose.yml
│   ├── config/
│   │   └── my.cnf
│   └── README.md
├── redis/
│   ├── versions/
│   │   ├── redis-6.2-docker-compose.yml
│   │   ├── redis-7.0-docker-compose.yml
│   │   └── redis-7.2-docker-compose.yml
│   ├── config/
│   │   └── redis.conf
│   └── README.md
└── ...
```

### 2. 使用部署脚本

```bash
# 基本语法
./deploy.sh <service> <version> <action> [options]

# 示例: 部署 MySQL 8.0
./deploy.sh mysql mysql-8.0 deploy

# 示例: 部署 Redis 7.2 并指定端口
./deploy.sh redis redis-7.2 deploy --port 6380

# 示例: 查看服务状态
./deploy.sh mysql mysql-8.0 status

# 示例: 查看日志
./deploy.sh mysql mysql-8.0 logs

# 示例: 停止服务
./deploy.sh mysql mysql-8.0 stop

# 示例: 删除服务
./deploy.sh mysql mysql-8.0 remove
```

### 3. 多实例部署

使用不同的项目名称部署多个实例：

```bash
# 部署实例1
./deploy.sh mysql mysql-8.0 deploy -p app-db-1 --port 3306

# 部署实例2
./deploy.sh mysql mysql-8.0 deploy -p app-db-2 --port 3307

# 部署实例3
./deploy.sh mysql mysql-8.0 deploy -p app-db-3 --port 3308
```

### 4. 使用环境变量文件

创建 `.env` 文件：

```bash
# mysql-app1.env
CONTAINER_NAME=mysql-app1
MYSQL_ROOT_PASSWORD=SecurePass123
MYSQL_DATABASE=app1_db
MYSQL_PORT=3306
DATA_DIR=/data/mysql-app1
```

部署：

```bash
./deploy.sh mysql mysql-8.0 deploy -e mysql-app1.env
```

## 🔧 后端 API 集成

### Shell 脚本调用示例

```bash
#!/bin/bash
# 后端 API 调用部署脚本示例

SERVICE="mysql"
VERSION="mysql-8.0"
PROJECT_NAME="app-db-production"
DATA_DIR="/data/mysql-prod"
MYSQL_PASSWORD="SecurePass123"

# 执行部署
/path/to/templates/deploy.sh \
    $SERVICE \
    $VERSION \
    deploy \
    -p "$PROJECT_NAME" \
    -d "$DATA_DIR" \
    --password "$MYSQL_PASSWORD"

# 检查部署结果
if [ $? -eq 0 ]; then
    echo "部署成功"
else
    echo "部署失败"
    exit 1
fi
```

### Python 调用示例

```python
import subprocess
import json

def deploy_service(service, version, project_name, **kwargs):
    """
    部署 Docker Compose 服务

    Args:
        service: 服务名称 (mysql, redis, etc.)
        version: 版本标识 (mysql-8.0, redis-7.2, etc.)
        project_name: 项目名称
        **kwargs: 其他参数 (port, password, data_dir, etc.)

    Returns:
        dict: {"success": bool, "message": str, "output": str}
    """

    cmd = [
        "/path/to/templates/deploy.sh",
        service,
        version,
        "deploy",
        "-p", project_name
    ]

    # 添加可选参数
    if 'data_dir' in kwargs:
        cmd.extend(["-d", kwargs['data_dir']])

    if 'port' in kwargs:
        cmd.extend(["--port", str(kwargs['port'])])

    if 'password' in kwargs:
        cmd.extend(["--password", kwargs['password']])

    try:
        result = subprocess.run(
            cmd,
            capture_output=True,
            text=True,
            timeout=300  # 5分钟超时
        )

        return {
            "success": result.returncode == 0,
            "message": "部署成功" if result.returncode == 0 else "部署失败",
            "output": result.stdout,
            "error": result.stderr
        }

    except subprocess.TimeoutExpired:
        return {
            "success": False,
            "message": "部署超时",
            "error": "操作超过5分钟"
        }

    except Exception as e:
        return {
            "success": False,
            "message": "部署异常",
            "error": str(e)
        }

# 使用示例
if __name__ == "__main__":
    result = deploy_service(
        service="mysql",
        version="mysql-8.0",
        project_name="app-db-1",
        data_dir="/data/mysql-app1",
        port=3306,
        password="MySecurePass123"
    )

    print(json.dumps(result, indent=2, ensure_ascii=False))
```

### Node.js 调用示例

```javascript
const { exec } = require('child_process');
const util = require('util');
const execPromise = util.promisify(exec);

async function deployService(service, version, projectName, options = {}) {
  const cmd = [
    '/path/to/templates/deploy.sh',
    service,
    version,
    'deploy',
    '-p', projectName
  ];

  if (options.dataDir) {
    cmd.push('-d', options.dataDir);
  }

  if (options.port) {
    cmd.push('--port', options.port);
  }

  if (options.password) {
    cmd.push('--password', options.password);
  }

  try {
    const { stdout, stderr } = await execPromise(cmd.join(' '), {
      timeout: 300000 // 5分钟
    });

    return {
      success: true,
      message: '部署成功',
      output: stdout
    };
  } catch (error) {
    return {
      success: false,
      message: '部署失败',
      error: error.message,
      stderr: error.stderr
    };
  }
}

// 使用示例
(async () => {
  const result = await deployService('mysql', 'mysql-8.0', 'app-db-1', {
    dataDir: '/data/mysql-app1',
    port: 3306,
    password: 'MySecurePass123'
  });

  console.log(JSON.stringify(result, null, 2));
})();
```

## 📋 环境变量说明

每个服务都支持通过环境变量自定义配置，常用环境变量包括：

### 通用变量

- `CONTAINER_NAME`: 容器名称
- `DATA_DIR`: 数据目录
- `TZ`: 时区 (默认: Asia/Shanghai)

### 数据库特定变量

**MySQL**:
- `MYSQL_ROOT_PASSWORD`: root 密码
- `MYSQL_DATABASE`: 初始数据库
- `MYSQL_USER`: 普通用户
- `MYSQL_PASSWORD`: 普通用户密码
- `MYSQL_PORT`: 端口

**PostgreSQL**:
- `POSTGRES_USER`: 超级用户
- `POSTGRES_PASSWORD`: 超级用户密码
- `POSTGRES_DB`: 初始数据库
- `POSTGRES_PORT`: 端口

**Redis**:
- `REDIS_PASSWORD`: Redis 密码
- `REDIS_PORT`: 端口
- `REDIS_MAXMEMORY`: 最大内存

详细配置请查看各服务目录下的 README.md。

## 🛠 常用操作

### 查看运行中的容器

```bash
docker ps
```

### 查看所有容器（包括停止的）

```bash
docker ps -a
```

### 查看容器日志

```bash
./deploy.sh <service> <version> logs -p <project-name>
```

### 进入容器

```bash
./deploy.sh <service> <version> exec -p <project-name>
```

### 备份数据

```bash
# 直接备份数据目录
tar -czf backup-$(date +%Y%m%d).tar.gz /data/<service>

# MySQL 备份
docker exec <container> mysqldump -u root -p<password> --all-databases > backup.sql

# PostgreSQL 备份
docker exec <container> pg_dumpall -U postgres > backup.sql

# Redis 备份
docker exec <container> redis-cli -a <password> BGSAVE
cp /data/redis/data/dump.rdb backup/
```

## 📊 监控与维护

### 健康检查

所有服务都配置了健康检查 (healthcheck)，可以通过以下命令查看：

```bash
docker inspect <container-name> | grep -A 10 Health
```

### 资源监控

```bash
# 查看容器资源使用情况
docker stats

# 查看单个容器
docker stats <container-name>
```

### 清理资源

```bash
# 清理停止的容器
docker container prune

# 清理未使用的镜像
docker image prune

# 清理未使用的卷
docker volume prune

# 清理所有未使用的资源
docker system prune -a
```

## 🔒 安全建议

1. **修改默认密码**: 所有服务的默认密码必须修改
2. **使用强密码**: 密码长度至少 12 位，包含大小写字母、数字、特殊字符
3. **限制网络访问**: 使用防火墙限制端口访问
4. **定期备份**: 每天自动备份重要数据
5. **及时更新**: 定期更新镜像版本，修复安全漏洞
6. **最小权限**: 避免使用 root 用户运行容器
7. **SSL/TLS**: 生产环境启用 HTTPS

## 🆘 故障排查

### 容器无法启动

```bash
# 查看容器日志
docker logs <container-name>

# 查看最近 100 行日志
docker logs --tail 100 <container-name>

# 实时查看日志
docker logs -f <container-name>
```

### 端口冲突

```bash
# 查看端口占用
netstat -tulpn | grep <port>
lsof -i :<port>

# 修改端口映射
./deploy.sh <service> <version> deploy --port <new-port>
```

### 磁盘空间不足

```bash
# 查看磁盘使用
df -h

# 清理 Docker 资源
docker system prune -a

# 清理日志
docker exec <container> sh -c "truncate -s 0 /var/log/*.log"
```

### 性能问题

```bash
# 查看资源使用
docker stats

# 限制容器资源
docker update --cpus=2 --memory=2g <container-name>
```

## 📚 更多文档

每个服务的详细文档请查看对应目录：

- [MySQL 文档](mysql/README.md)
- [Redis 文档](redis/README.md)
- [PostgreSQL 文档](postgresql/README.md)
- [Jenkins 文档](jenkins/README.md)
- [GitLab 文档](gitlab/README.md)
- [Grafana 文档](grafana/README.md)
- [JumpServer 文档](jumpserver/README.md)

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📄 许可证

MIT License
