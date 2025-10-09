# DevOps API 镜像
# 使用已编译的二进制文件直接构建镜像

FROM alpine:latest

LABEL maintainer="zfwh1024@163.com"
LABEL description="DevOps运维管理系统 API服务"
LABEL version="1.0"

# 安装运行时依赖
RUN apk --no-cache add \
    ca-certificates \
    tzdata \
    wget \
    && rm -rf /var/cache/apk/*

# 创建工作目录
WORKDIR /app

# 复制已编译的二进制文件（需要先编译好devops二进制文件）
COPY devops /app/devops

# 创建日志和上传目录
RUN mkdir -p /app/log /app/upload

# 添加执行权限
RUN chmod +x /app/devops

# 设置时区
ENV TZ=Asia/Shanghai

# 暴露端口
EXPOSE 8000

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=40s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8000/api/v1/health || exit 1

# 启动命令（配置文件通过volume挂载）
CMD ["./devops", "-c", "/app/config.yaml"]
