# JumpServer Docker 部署模板

## 简介

JumpServer 是开源堡垒机，提供：
- SSH/RDP 协议跳转
- 录像审计
- 命令记录
- 资产管理
- 权限控制

## 环境变量配置

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| SECRET_KEY | 加密密钥 | 随机生成 |
| BOOTSTRAP_TOKEN | 初始Token | 随机生成 |
| HTTP_PORT | Web端口 | 80 |
| SSH_PORT | SSH端口 | 2222 |
| DB_PASSWORD | 数据库密码 | jumpserver123 |
| DATA_DIR | 数据目录 | /data/jumpserver |

## 快速部署

```bash
# 生成随机密钥
if [ "$SECRET_KEY" = "" ]; then SECRET_KEY=`cat /dev/urandom | tr -dc A-Za-z0-9 | head -c 50`; echo "SECRET_KEY=$SECRET_KEY" >> .env; fi
if [ "$BOOTSTRAP_TOKEN" = "" ]; then BOOTSTRAP_TOKEN=`cat /dev/urandom | tr -dc A-Za-z0-9 | head -c 16`; echo "BOOTSTRAP_TOKEN=$BOOTSTRAP_TOKEN" >> .env; fi

docker-compose -f jumpserver-latest-docker-compose.yml up -d
```

## 访问

- Web: `http://your-ip`
- 默认账号: `admin`
- 默认密码: `admin`

首次登录会要求修改密码。

## 注意事项

1. 首次启动需要 2-3 分钟
2. 必须修改 SECRET_KEY 和 BOOTSTRAP_TOKEN
3. 生产环境必须修改默认密码
4. 定期备份数据目录
