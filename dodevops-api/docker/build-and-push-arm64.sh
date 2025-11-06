#!/bin/bash

# 设置变量
IMAGE_NAME="crpi-aj3vgoxp9kzh2jx4.cn-hangzhou.personal.cr.aliyuncs.com/zhangfan_k8s/deviops-api"
VERSION="v1.0-arm64"
FULL_IMAGE_NAME="${IMAGE_NAME}:${VERSION}"

echo "=========================================="
echo "开始构建和推送 ARM64 镜像: ${FULL_IMAGE_NAME}"
echo "=========================================="

# 1. 编译项目 (ARM64)
echo ""
echo "[1/4] 编译项目 (ARM64)..."
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
cd "${SCRIPT_DIR}/.."
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o devops .
if [ $? -ne 0 ]; then
    echo "编译失败，退出脚本"
    exit 1
fi
echo "编译成功"

# 2. 构建镜像 (ARM64)
echo ""
echo "[2/4] 构建 Docker 镜像 (ARM64)..."
docker build --platform linux/arm64 -f docker/Dockerfile -t ${FULL_IMAGE_NAME} .
if [ $? -ne 0 ]; then
    echo "镜像构建失败，退出脚本"
    rm -f devops
    exit 1
fi
echo "镜像构建成功"

# 3. 推送镜像
echo ""
echo "[3/4] 推送 Docker 镜像..."
docker push ${FULL_IMAGE_NAME}
if [ $? -ne 0 ]; then
    echo "镜像推送失败，退出脚本"
    rm -f devops
    docker rmi ${FULL_IMAGE_NAME} 2>/dev/null
    exit 1
fi
echo "镜像推送成功"

# 4. 清理本地文件和镜像
echo ""
echo "[4/4] 清理本地文件和镜像..."
rm -f devops
echo "已删除本地devops文件"

docker rmi ${FULL_IMAGE_NAME}
if [ $? -eq 0 ]; then
    echo "已删除本地Docker镜像: ${FULL_IMAGE_NAME}"
else
    echo "删除本地Docker镜像失败（可能不影响使用）"
fi

echo ""
echo "=========================================="
echo "所有操作已完成！"
echo "=========================================="
