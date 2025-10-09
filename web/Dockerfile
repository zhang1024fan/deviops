# DevOps 前端镜像
# 将dist静态文件打包到nginx镜像中

FROM nginx:1.24-alpine

LABEL maintainer="zfwh1024@163.com"
LABEL description="DevOps运维管理系统 前端服务"
LABEL version="1.0"

# 删除nginx默认配置
RUN rm -rf /etc/nginx/conf.d/default.conf

# 复制前端构建文件到nginx目录
COPY dist /var/www/html/dist

# 设置文件权限
RUN chmod -R 755 /var/www/html/dist

# 创建上传文件目录（用于volume挂载）
RUN mkdir -p /home/devops/upload

# 设置时区
ENV TZ=Asia/Shanghai

# 暴露端口
EXPOSE 80 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:80/ || exit 1

# 启动nginx（配置文件通过volume挂载）
CMD ["nginx", "-g", "daemon off;"]
