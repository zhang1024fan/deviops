#!/bin/bash

################################################################################
# DevOps 镜像构建和推送脚本
# 用途：构建前后端Docker镜像并推送到私有镜像仓库
################################################################################

set -e  # 遇到错误立即退出

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 配置变量
API_IMAGE_NAME="devops-api"
WEB_IMAGE_NAME="devops-web"
VERSION=${1:-"latest"}

# 镜像仓库配置（根据需要修改）
REGISTRY="your-registry.com"              # 私有仓库地址
NAMESPACE="devops"                        # 命名空间
REGISTRY_USER="your-username"             # 仓库用户名（可选）

# 完整镜像名称
API_IMAGE="${REGISTRY}/${NAMESPACE}/${API_IMAGE_NAME}"
WEB_IMAGE="${REGISTRY}/${NAMESPACE}/${WEB_IMAGE_NAME}"

# 函数：打印信息
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

# 函数：显示banner
show_banner() {
    cat << "EOF"
========================================
   DevOps 镜像构建和推送工具
========================================
EOF
}

# 函数：检查命令是否存在
check_command() {
    if ! command -v $1 &> /dev/null; then
        print_error "$1 命令不存在，请先安装"
        exit 1
    fi
}

# 函数：检查必要文件
check_files() {
    print_info "检查必要文件..."

    # 检查后端二进制文件
    if [ ! -f "../api/devops" ]; then
        print_error "后端二进制文件不存在: ../api/devops"
        print_error "请先编译 devops 二进制文件并放到 api/ 目录下"
        exit 1
    fi

    # 检查前端构建文件
    if [ ! -d "../web/dist" ]; then
        print_error "前端构建文件不存在: ../web/dist"
        print_error "请先构建前端项目"
        exit 1
    fi

    # 检查Dockerfile
    if [ ! -f "../api/Dockerfile" ]; then
        print_error "后端Dockerfile不存在: ../api/Dockerfile"
        exit 1
    fi

    if [ ! -f "../web/Dockerfile" ]; then
        print_error "前端Dockerfile不存在: ../web/Dockerfile"
        exit 1
    fi

    print_success "文件检查通过"
}

# 函数：登录镜像仓库
login_registry() {
    print_info "登录镜像仓库: ${REGISTRY}"

    if [ -n "${REGISTRY_USER}" ]; then
        docker login ${REGISTRY} -u ${REGISTRY_USER}
    else
        docker login ${REGISTRY}
    fi

    if [ $? -eq 0 ]; then
        print_success "登录成功"
    else
        print_error "登录失败"
        exit 1
    fi
}

# 函数：构建后端镜像
build_api_image() {
    print_info "构建后端镜像: ${API_IMAGE}:${VERSION}"

    cd ../api
    docker build -t ${API_IMAGE}:${VERSION} -f Dockerfile .

    if [ $? -eq 0 ]; then
        # 同时打上latest标签
        docker tag ${API_IMAGE}:${VERSION} ${API_IMAGE}:latest
        print_success "后端镜像构建成功"
    else
        print_error "后端镜像构建失败"
        exit 1
    fi

    cd ../docker
}

# 函数：构建前端镜像
build_web_image() {
    print_info "构建前端镜像: ${WEB_IMAGE}:${VERSION}"

    cd ../web
    docker build -t ${WEB_IMAGE}:${VERSION} -f Dockerfile .

    if [ $? -eq 0 ]; then
        # 同时打上latest标签
        docker tag ${WEB_IMAGE}:${VERSION} ${WEB_IMAGE}:latest
        print_success "前端镜像构建成功"
    else
        print_error "前端镜像构建失败"
        exit 1
    fi

    cd ../docker
}

# 函数：推送镜像
push_images() {
    print_info "推送镜像到仓库..."

    # 推送后端镜像
    print_info "推送后端镜像: ${API_IMAGE}:${VERSION}"
    docker push ${API_IMAGE}:${VERSION}

    print_info "推送后端镜像: ${API_IMAGE}:latest"
    docker push ${API_IMAGE}:latest

    # 推送前端镜像
    print_info "推送前端镜像: ${WEB_IMAGE}:${VERSION}"
    docker push ${WEB_IMAGE}:${VERSION}

    print_info "推送前端镜像: ${WEB_IMAGE}:latest"
    docker push ${WEB_IMAGE}:latest

    print_success "所有镜像推送完成"
}

# 函数：显示镜像信息
show_image_info() {
    echo ""
    print_info "构建的镜像列表:"
    docker images | grep -E "${API_IMAGE_NAME}|${WEB_IMAGE_NAME}" | head -10

    echo ""
    print_info "镜像拉取命令:"
    echo "  docker pull ${API_IMAGE}:${VERSION}"
    echo "  docker pull ${WEB_IMAGE}:${VERSION}"

    echo ""
    print_info "使用docker-compose部署:"
    echo "  export DEVOPS_API_IMAGE=${API_IMAGE}:${VERSION}"
    echo "  export DEVOPS_WEB_IMAGE=${WEB_IMAGE}:${VERSION}"
    echo "  docker-compose up -d"
}

# 函数：显示使用说明
show_usage() {
    cat << EOF
使用方法:
    ./build-and-push.sh [VERSION] [OPTIONS]

参数说明:
    VERSION   - 镜像版本号 (默认: latest)

选项:
    --build-only    仅构建镜像，不推送
    --push-only     仅推送镜像，不构建
    --api-only      仅操作后端镜像
    --web-only      仅操作前端镜像
    -h, --help      显示帮助信息

示例:
    ./build-and-push.sh v1.0.0              # 构建v1.0.0版本并推送
    ./build-and-push.sh latest --build-only # 仅构建latest版本
    ./build-and-push.sh v2.0.0 --api-only   # 仅构建和推送后端镜像

构建前准备:
    1. 编译后端二进制文件，放到 api/ 目录
    2. 构建前端项目，确保 web/dist/ 目录存在
    3. 修改脚本中的镜像仓库配置

镜像仓库配置:
    REGISTRY="${REGISTRY}"
    NAMESPACE="${NAMESPACE}"

EOF
}

# 主程序
main() {
    # 显示banner
    show_banner
    echo ""

    # 检查参数
    if [ "$1" == "-h" ] || [ "$1" == "--help" ]; then
        show_usage
        exit 0
    fi

    # 检查Docker是否安装
    check_command docker

    # 解析参数
    BUILD_ONLY=false
    PUSH_ONLY=false
    API_ONLY=false
    WEB_ONLY=false

    for arg in "$@"; do
        case $arg in
            --build-only)
                BUILD_ONLY=true
                ;;
            --push-only)
                PUSH_ONLY=true
                ;;
            --api-only)
                API_ONLY=true
                ;;
            --web-only)
                WEB_ONLY=true
                ;;
        esac
    done

    # 显示配置信息
    print_info "镜像仓库: ${REGISTRY}"
    print_info "命名空间: ${NAMESPACE}"
    print_info "版本号: ${VERSION}"
    echo ""

    # 检查文件
    check_files
    echo ""

    # 构建镜像
    if [ "$PUSH_ONLY" = false ]; then
        if [ "$WEB_ONLY" = false ]; then
            build_api_image
        fi

        if [ "$API_ONLY" = false ]; then
            build_web_image
        fi
        echo ""
    fi

    # 推送镜像
    if [ "$BUILD_ONLY" = false ]; then
        # 登录仓库
        login_registry
        echo ""

        # 推送镜像
        push_images
    fi

    # 显示镜像信息
    show_image_info

    echo ""
    print_success "所有操作完成!"
}

# 执行主程序
main "$@"
