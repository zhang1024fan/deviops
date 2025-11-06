# Jenkins Docker 部署模板

## 支持版本

- Jenkins LTS (长期支持版，推荐生产环境)
- Jenkins Latest (最新版)

## 环境变量配置

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| CONTAINER_NAME | 容器名称 | jenkins-lts/jenkins-latest |
| JENKINS_HTTP_PORT | Web端口 | 8080 |
| JENKINS_AGENT_PORT | Agent端口 | 50000 |
| DATA_DIR | 数据目录 | /data/jenkins |

## 快速部署

### 启动 Jenkins

```bash
docker-compose -f jenkins-lts-docker-compose.yml up -d
```

### 获取初始密码

```bash
# 方法1：查看日志
docker logs jenkins-lts

# 方法2：读取文件
docker exec jenkins-lts cat /var/jenkins_home/secrets/initialAdminPassword
```

### 访问 Jenkins

浏览器访问: `http://your-ip:8080`

输入初始密码，安装推荐插件。

## 数据持久化

所有数据存储在 `${DATA_DIR}/home` 目录：
- 构建历史
- 插件
- 配置文件
- 凭据

## 常用操作

### 查看日志
```bash
docker logs -f jenkins-lts
```

### 重启 Jenkins
```bash
docker restart jenkins-lts
```

### 备份 Jenkins
```bash
tar -czf jenkins-backup-$(date +%Y%m%d).tar.gz /data/jenkins/home/
```

### 恢复 Jenkins
```bash
docker-compose down
tar -xzf jenkins-backup-20231201.tar.gz -C /
docker-compose up -d
```

### 更新插件
```bash
# 进入容器
docker exec -it jenkins-lts bash

# 批量更新插件
java -jar /var/jenkins_home/war/WEB-INF/jenkins-cli.jar -s http://localhost:8080/ install-plugin <plugin-name>
```

## 推荐插件

### 基础插件
- Git
- Pipeline
- Docker Pipeline
- Blue Ocean (现代化UI)
- Configuration as Code (JCasC)

### 集成插件
- GitLab Plugin
- GitHub Plugin
- DingTalk (钉钉通知)
- Email Extension

### 部署插件
- Kubernetes
- SSH
- Ansible
- Publish Over SSH

## Docker in Docker

本模板已配置 Docker Socket 挂载，支持在 Jenkins Pipeline 中使用 Docker 命令：

```groovy
pipeline {
    agent any
    stages {
        stage('Build') {
            steps {
                sh 'docker build -t myapp:latest .'
            }
        }
        stage('Push') {
            steps {
                sh 'docker push myapp:latest'
            }
        }
    }
}
```

## 配置为代码 (JCasC)

创建 `jenkins.yaml` 并挂载：

```yaml
volumes:
  - ${DATA_DIR}/home:/var/jenkins_home
  - ./jenkins.yaml:/var/jenkins_home/jenkins.yaml
```

## 性能优化

### Java 内存配置

根据服务器配置调整：

**4GB 内存服务器**:
```yaml
JAVA_OPTS=-Xmx2048m -Xms512m
```

**8GB 内存服务器**:
```yaml
JAVA_OPTS=-Xmx4096m -Xms1024m
```

### 构建优化

1. 使用 Pipeline 而非 Freestyle
2. 合理使用缓存 (Maven/npm)
3. 并行构建
4. 定期清理旧构建

## 安全加固

1. **启用 HTTPS**: 使用 Nginx 反向代理
2. **启用认证**: 配置 LDAP/OAuth
3. **权限管理**: 使用 Role-based Authorization
4. **凭据管理**: 使用 Credentials Plugin 存储敏感信息
5. **定期备份**: 每天自动备份配置

## 高可用部署

如需 Jenkins 集群，推荐使用：
- Jenkins 主从架构
- Kubernetes + Jenkins Operator
- CloudBees Jenkins Enterprise

## 常见问题

### 1. 容器启动后无法访问

等待 1-2 分钟，Jenkins 初始化需要时间。查看日志：
```bash
docker logs -f jenkins-lts
```

### 2. 插件安装失败

更换插件源为国内镜像：
```bash
# 进入容器修改
sed -i 's#https://updates.jenkins.io/update-center.json#https://mirrors.tuna.tsinghua.edu.cn/jenkins/updates/update-center.json#g' \
  /var/jenkins_home/hudson.model.UpdateCenter.xml
```

### 3. 磁盘空间不足

定期清理构建历史：
```bash
# Jenkins UI: 管理Jenkins -> 磁盘使用
# 或使用脚本清理旧构建
```

## 注意事项

1. **首次启动**: 需要 2-3 分钟初始化
2. **初始密码**: 保存好初始管理员密码
3. **数据备份**: 定期备份 `/var/jenkins_home`
4. **资源限制**: 至少 2GB 内存
5. **防火墙**: 开放 8080 和 50000 端口
