# GitLab Docker 部署模板

## 支持版本

- GitLab CE (社区版，免费)
- GitLab EE (企业版，需要许可证)

## 系统要求

**最低配置**:
- CPU: 4 核心
- 内存: 4GB
- 磁盘: 50GB

**推荐配置**:
- CPU: 8 核心
- 内存: 8GB
- 磁盘: 100GB SSD

## 环境变量配置

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| CONTAINER_NAME | 容器名称 | gitlab-ce/gitlab-ee |
| GITLAB_HOSTNAME | GitLab域名 | gitlab.example.com |
| GITLAB_HTTP_PORT | HTTP端口 | 80 |
| GITLAB_HTTPS_PORT | HTTPS端口 | 443 |
| GITLAB_SSH_PORT | SSH端口 | 2222 |
| DATA_DIR | 数据目录 | /data/gitlab |

## 快速部署

### 1. 配置环境变量

创建 `.env` 文件：
```bash
CONTAINER_NAME=gitlab-ce
GITLAB_HOSTNAME=gitlab.yourdomain.com
GITLAB_HTTP_PORT=80
GITLAB_SSH_PORT=2222
DATA_DIR=/data/gitlab
```

### 2. 启动 GitLab

```bash
docker-compose -f gitlab-ce-docker-compose.yml up -d
```

### 3. 等待初始化

首次启动需要 **5-10 分钟**，查看日志：
```bash
docker logs -f gitlab-ce
```

### 4. 获取初始密码

```bash
docker exec -it gitlab-ce grep 'Password:' /etc/gitlab/initial_root_password
```

初始用户名: `root`

### 5. 访问 GitLab

浏览器访问: `http://your-ip`

## 数据持久化

数据存储在以下目录：
- `${DATA_DIR}/config` - GitLab 配置
- `${DATA_DIR}/logs` - 日志文件
- `${DATA_DIR}/data` - Git 仓库、数据库、附件

## 常用操作

### 查看日志
```bash
docker logs -f gitlab-ce
```

### 重新配置
```bash
docker exec -it gitlab-ce gitlab-ctl reconfigure
```

### 重启服务
```bash
docker exec -it gitlab-ce gitlab-ctl restart
```

### 查看状态
```bash
docker exec -it gitlab-ce gitlab-ctl status
```

### 进入 Rails 控制台
```bash
docker exec -it gitlab-ce gitlab-rails console
```

### 备份 GitLab
```bash
# 创建备份
docker exec -it gitlab-ce gitlab-backup create

# 备份文件位置
ls /data/gitlab/data/backups/
```

### 恢复 GitLab
```bash
# 停止服务
docker exec -it gitlab-ce gitlab-ctl stop puma
docker exec -it gitlab-ce gitlab-ctl stop sidekiq

# 恢复备份 (假设备份文件为 1234567890_2023_12_01_16.0.0_gitlab_backup.tar)
docker exec -it gitlab-ce gitlab-backup restore BACKUP=1234567890_2023_12_01_16.0.0

# 重启服务
docker exec -it gitlab-ce gitlab-ctl restart
docker exec -it gitlab-ce gitlab-rake gitlab:check SANITIZE=true
```

## 邮件配置

编辑 `docker-compose.yml` 中的 GITLAB_OMNIBUS_CONFIG：

```yaml
GITLAB_OMNIBUS_CONFIG: |
  # SMTP 配置
  gitlab_rails['smtp_enable'] = true
  gitlab_rails['smtp_address'] = "smtp.gmail.com"
  gitlab_rails['smtp_port'] = 587
  gitlab_rails['smtp_user_name'] = "your-email@gmail.com"
  gitlab_rails['smtp_password'] = "your-password"
  gitlab_rails['smtp_domain'] = "gmail.com"
  gitlab_rails['smtp_authentication'] = "login"
  gitlab_rails['smtp_enable_starttls_auto'] = true
  gitlab_rails['smtp_tls'] = false
  gitlab_rails['gitlab_email_from'] = 'your-email@gmail.com'
  gitlab_rails['gitlab_email_reply_to'] = 'noreply@example.com'
```

重新配置：
```bash
docker exec -it gitlab-ce gitlab-ctl reconfigure
```

测试邮件：
```bash
docker exec -it gitlab-ce gitlab-rails console
Notify.test_email('user@example.com', 'Test Subject', 'Test Body').deliver_now
```

## HTTPS 配置

### 方法1：使用 Let's Encrypt

```yaml
GITLAB_OMNIBUS_CONFIG: |
  external_url 'https://gitlab.example.com'
  letsencrypt['enable'] = true
  letsencrypt['contact_emails'] = ['admin@example.com']
```

### 方法2：使用自签名证书

```bash
# 生成证书
mkdir -p /data/gitlab/config/ssl
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout /data/gitlab/config/ssl/gitlab.key \
  -out /data/gitlab/config/ssl/gitlab.crt
```

```yaml
GITLAB_OMNIBUS_CONFIG: |
  external_url 'https://gitlab.example.com'
  nginx['ssl_certificate'] = "/etc/gitlab/ssl/gitlab.crt"
  nginx['ssl_certificate_key'] = "/etc/gitlab/ssl/gitlab.key"
```

## GitLab Runner 集成

```bash
# 获取 Registration Token
# GitLab UI: Admin Area -> CI/CD -> Runners

# 启动 Runner
docker run -d --name gitlab-runner \
  --restart unless-stopped \
  -v /data/gitlab-runner/config:/etc/gitlab-runner \
  -v /var/run/docker.sock:/var/run/docker.sock \
  gitlab/gitlab-runner:latest

# 注册 Runner
docker exec -it gitlab-runner gitlab-runner register \
  --non-interactive \
  --url "http://your-gitlab-url" \
  --registration-token "YOUR_TOKEN" \
  --executor "docker" \
  --docker-image "alpine:latest" \
  --description "docker-runner" \
  --tag-list "docker,aws" \
  --run-untagged="true" \
  --locked="false"
```

## 性能优化

### 内存优化 (4GB 内存服务器)

```yaml
GITLAB_OMNIBUS_CONFIG: |
  puma['worker_processes'] = 2
  sidekiq['max_concurrency'] = 10
  postgresql['shared_buffers'] = "128MB"
  postgresql['max_worker_processes'] = 4
  prometheus_monitoring['enable'] = false
```

### 清理旧数据

```bash
# 清理旧日志
docker exec -it gitlab-ce gitlab-ctl cleanup-log --days 7

# 清理旧备份
docker exec -it gitlab-ce find /var/opt/gitlab/backups -type f -mtime +7 -delete
```

## 监控

GitLab 内置 Prometheus + Grafana 监控：

访问: `http://your-gitlab/-/grafana`

## 常见问题

### 1. 502 Bad Gateway

GitLab 还在初始化，等待 5-10 分钟。

### 2. 内存不足

至少需要 4GB 内存，调整配置减少进程数。

### 3. Git Push 失败

检查 SSH 端口映射和防火墙配置。

### 4. 升级 GitLab

```bash
# 备份数据
docker exec -it gitlab-ce gitlab-backup create

# 停止容器
docker-compose down

# 拉取新版本
docker pull gitlab/gitlab-ce:latest

# 启动新版本
docker-compose up -d
```

## 注意事项

1. **首次启动**: 需要 5-10 分钟初始化
2. **资源要求**: 至少 4GB 内存
3. **定期备份**: 每天自动备份或手动备份
4. **版本升级**: 遵循官方升级路径
5. **SSH 端口**: 默认映射到 2222，避免与系统 SSH 冲突
6. **磁盘空间**: 预留足够空间存储 Git 仓库
