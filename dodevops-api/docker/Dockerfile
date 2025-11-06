FROM crpi-aj3vgoxp9kzh2jx4.cn-hangzhou.personal.cr.aliyuncs.com/zhangfan_k8s/golang:1.25-alpine

WORKDIR /app

# 设置 Alpine 国内镜像源
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

# 安装运行时依赖
RUN apk --no-cache add \
    ca-certificates \
    openssh-client \
    tzdata \
    curl \
    bash \
    git \
    gcc \
    musl-dev && \
    rm -rf /var/cache/apk/*

# 设置时区
ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# 配置 Go 环境变量
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \
    GOPATH=/go \
    GOCACHE=/go-cache

# 创建必要的目录
RUN mkdir -p /app/log /app/upload/xlsl /app/docs /go/pkg/mod /go-cache && \
    chmod -R 777 /go /go-cache /app/upload

# 复制编译好的二进制文件
COPY devops /app/devops
RUN chmod +x /app/devops

# 复制配置文件和静态资源
COPY config.yaml /app/config.yaml
COPY docs/ /app/docs/
COPY common/ /app/common/

# 复制模板文件（确保在本地构建时该文件存在）
COPY upload/xlsl/host.xlsx /app/upload/xlsl/host.xlsx

# 暴露端口
EXPOSE 8000

# 启动命令
CMD ["/app/devops"]
