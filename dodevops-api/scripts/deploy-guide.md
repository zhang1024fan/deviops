# 生产环境部署指南

## 📋 部署流程总览

```
本地开发 → 构建镜像 → 推送到私有仓库 → 服务器拉取 → 启动服务
```

---

## 🏗️ 第一步：构建并推送镜像

### 方式一：使用自动化脚本（推荐）

```bash
# 1. 设置环境变量
export REGISTRY_URL="registry.example.com"
export IMAGE_TAG="v1.0.0"
export REGISTRY_USERNAME="your-username"
export REGISTRY_PASSWORD="your-password"

# 2. 运行构建脚本
chmod +x scripts/build-and-push.sh
./scripts/build-and-push.sh
```

### 方式二：手动构建推送

```bash
# 1. 构建镜像
docker build -t registry.example.com/dodevops-api:v1.0.0 .

# 2. 登录私有仓库
docker login registry.example.com

# 3. 推送镜像
docker push registry.example.com/dodevops-api:v1.0.0

# 4. 推送 latest 标签（可选）
docker tag registry.example.com/dodevops-api:v1.0.0 \
           registry.example.com/dodevops-api:latest
docker push registry.example.com/dodevops-api:latest
```

---

## 🚀 第二步：服务器部署

### 1. 准备部署目录

```bash
# 在服务器上创建部署目录
mkdir -p /opt/dodevops-api
cd /opt/dodevops-api

# 创建必要的子目录
mkdir -p log upload
```

### 2. 准备配置文件

```bash
# 复制配置文件到服务器
scp config.yaml user@server:/opt/dodevops-api/

# 或在服务器上直接编辑
vim /opt/dodevops-api/config.yaml
```

**重要**：修改 `config.yaml` 中的配置：
```yaml
server:
  address: 0.0.0.0:8000  # 监听所有网卡

db:
  host: your-db-host
  password: your-db-password

redis:
  address: your-redis-host:6379
  password: your-redis-password
```

### 3. 创建 docker-compose.yml

```bash
# 在服务器上创建 docker-compose.yml
cat > /opt/dodevops-api/docker-compose.yml << 'EOF'
version: '3.8'

services:
  dodevops-api:
    image: registry.example.com/dodevops-api:latest
    container_name: dodevops-api
    restart: unless-stopped

    ports:
      - "8000:8000"

    environment:
      - TZ=Asia/Shanghai
      - GIN_MODE=release

    volumes:
      - ./config.yaml:/app/config.yaml:ro
      - ./log:/app/log
      - ./upload:/app/upload
      - go-mod-cache:/go/pkg/mod
      - go-build-cache:/go-cache

    networks:
      - dodevops-network

    deploy:
      resources:
        limits:
          cpus: '2'
          memory: 2G

    healthcheck:
      test: ["CMD", "curl", "-f", "http://127.0.0.1:8000/api/v1/health"]
      interval: 30s
      timeout: 3s
      retries: 3

    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"

networks:
  dodevops-network:
    driver: bridge

volumes:
  go-mod-cache:
  go-build-cache:
EOF
```

### 4. 登录私有仓库并拉取镜像

```bash
# 登录私有仓库
docker login registry.example.com

# 拉取镜像
docker pull registry.example.com/dodevops-api:latest
```

### 5. 启动服务

```bash
# 启动服务
docker-compose up -d

# 查看日志
docker-compose logs -f

# 检查服务状态
docker-compose ps
```

---

## ✅ 第三步：验证部署

```bash
# 1. 检查容器状态
docker ps | grep dodevops-api

# 2. 查看健康状态（应显示 healthy）
docker inspect dodevops-api | grep Health -A 10

# 3. 测试 API
curl http://localhost:8000/api/v1/health

# 4. 查看日志
tail -f log/sys.log

# 5. 测试 agent 部署（确保 Go 编译功能正常）
# 通过 API 部署 agent 到某台主机
```

---

## 🔄 更新部署

### 方式一：拉取最新镜像

```bash
# 1. 拉取最新镜像
docker-compose pull

# 2. 重启服务
docker-compose up -d

# 3. 清理旧镜像
docker image prune -f
```

### 方式二：指定版本更新

```bash
# 1. 修改 docker-compose.yml 中的镜像标签
vim docker-compose.yml
# image: registry.example.com/dodevops-api:v1.0.1

# 2. 拉取并重启
docker-compose pull
docker-compose up -d
```

---

## 📂 目录结构

部署后的服务器目录结构：

```
/opt/dodevops-api/
├── docker-compose.yml      # Docker Compose 配置
├── config.yaml             # 应用配置（挂载到容器）
├── log/                    # 日志目录（持久化）
│   └── sys.log
└── upload/                 # 上传文件目录（持久化）
```

---

## 🛠️ 常用运维命令

```bash
# 查看服务状态
docker-compose ps

# 查看实时日志
docker-compose logs -f

# 重启服务
docker-compose restart

# 停止服务
docker-compose stop

# 停止并删除容器
docker-compose down

# 进入容器调试
docker exec -it dodevops-api /bin/sh

# 查看资源使用
docker stats dodevops-api

# 更新并重启
docker-compose pull && docker-compose up -d
```

---

## 🔒 安全建议

1. **配置文件安全**
   - 不要将 `config.yaml` 提交到 Git
   - 使用只读挂载（`:ro`）
   - 定期轮换密码

2. **网络安全**
   - 使用防火墙限制端口访问
   - 配置 HTTPS（使用 Nginx 反向代理）
   - 内网部署使用私有网络

3. **日志管理**
   - 定期清理日志文件
   - 配置日志轮转
   - 敏感信息脱敏

4. **镜像安全**
   - 定期更新基础镜像
   - 扫描镜像漏洞
   - 使用私有仓库

---

## ❓ 常见问题

### Q1: 容器无法启动？

```bash
# 查看详细日志
docker-compose logs dodevops-api

# 检查配置文件路径
ls -la config.yaml

# 检查端口占用
netstat -tuln | grep 8000
```

### Q2: 无法拉取镜像？

```bash
# 检查网络
ping registry.example.com

# 重新登录
docker login registry.example.com

# 手动拉取测试
docker pull registry.example.com/dodevops-api:latest
```

### Q3: Agent 编译失败？

```bash
# 进入容器检查 Go 环境
docker exec -it dodevops-api /bin/sh
go version
go env

# 检查 Go 缓存权限
ls -la /go /go-cache
```

---

## 📞 支持

如有问题，请联系运维团队或提交 Issue。
