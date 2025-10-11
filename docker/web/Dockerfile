FROM crpi-aj3vgoxp9kzh2jx4.cn-hangzhou.personal.cr.aliyuncs.com/zhangfan_k8s/nginx:latest

# 删除默认的 nginx 配置和静态文件
RUN rm -rf /usr/share/nginx/html/*

# 复制编译好的静态文件到 nginx 目录
COPY dist /usr/share/nginx/html

# 暴露80端口
EXPOSE 80

# 启动 nginx
CMD ["nginx", "-g", "daemon off;"]
