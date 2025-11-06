#!/bin/bash
# ============================================
# 构建并推送 Docker 镜像到私有仓库
# ============================================

set -e

# ==================== 配置区 ====================
# 私有仓库地址（请替换为实际地址）
REGISTRY_URL="${REGISTRY_URL:-registry.example.com}"
IMAGE_NAME="dodevops-api"
IMAGE_TAG="${IMAGE_TAG:-latest}"

# 完整镜像名称
FULL_IMAGE_NAME="${REGISTRY_URL}/${IMAGE_NAME}:${IMAGE_TAG}"

# ==================== 颜色输出 ====================
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

function echo_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

function echo_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

function echo_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# ==================== 主流程 ====================
echo_info "开始构建 Docker 镜像..."
echo_info "镜像名称: ${FULL_IMAGE_NAME}"
echo ""

# 1. 检查 Docker 是否运行
if ! docker info > /dev/null 2>&1; then
    echo_error "Docker 未运行，请先启动 Docker"
    exit 1
fi

# 2. 构建镜像
echo_info "步骤 1/3: 构建镜像..."
docker build \
    -t "${FULL_IMAGE_NAME}" \
    -t "${REGISTRY_URL}/${IMAGE_NAME}:latest" \
    --build-arg BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
    --build-arg VERSION="${IMAGE_TAG}" \
    .

if [ $? -eq 0 ]; then
    echo_info "✅ 镜像构建成功"
else
    echo_error "❌ 镜像构建失败"
    exit 1
fi

# 3. 登录私有仓库（如果需要）
echo_info "步骤 2/3: 登录私有仓库..."
echo_warn "请输入私有仓库凭证（如果已登录可跳过）"

# 如果设置了环境变量，使用环境变量登录
if [ -n "${REGISTRY_USERNAME}" ] && [ -n "${REGISTRY_PASSWORD}" ]; then
    echo "${REGISTRY_PASSWORD}" | docker login "${REGISTRY_URL}" -u "${REGISTRY_USERNAME}" --password-stdin
else
    # 交互式登录
    docker login "${REGISTRY_URL}" || {
        echo_error "❌ 登录失败"
        exit 1
    }
fi

# 4. 推送镜像
echo_info "步骤 3/3: 推送镜像到私有仓库..."
docker push "${FULL_IMAGE_NAME}"

if [ $? -eq 0 ]; then
    echo_info "✅ 镜像推送成功"
else
    echo_error "❌ 镜像推送失败"
    exit 1
fi

# 可选：同时推送 latest 标签
if [ "${IMAGE_TAG}" != "latest" ]; then
    echo_info "推送 latest 标签..."
    docker push "${REGISTRY_URL}/${IMAGE_NAME}:latest"
fi

# 5. 清理本地镜像（可选）
echo_info "清理本地构建缓存..."
docker image prune -f

# ==================== 完成 ====================
echo ""
echo_info "=========================================="
echo_info "🎉 构建和推送完成！"
echo_info "=========================================="
echo_info "镜像地址: ${FULL_IMAGE_NAME}"
echo_info ""
echo_info "使用方法:"
echo_info "  1. 在服务器上创建 docker-compose.yml"
echo_info "  2. 修改镜像地址为: ${FULL_IMAGE_NAME}"
echo_info "  3. 运行: docker-compose up -d"
echo_info "=========================================="
