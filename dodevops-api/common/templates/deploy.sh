#!/bin/bash

###############################################################################
# Docker Compose 运维工具箱 - 统一部署脚本
# 用途: 后端 API 调用此脚本来部署各种服务
# 使用: ./deploy.sh <service> <version> <action> [options]
###############################################################################

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 日志函数
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 显示使用说明
usage() {
    cat << EOF
使用方法:
    $0 <service> <version> <action> [options]

参数说明:
    service     服务名称 (mysql, redis, postgresql, jenkins, gitlab, etc.)
    version     版本标识 (如: mysql-8.0, redis-7.2, nodejs-20)
    action      操作类型 (deploy, stop, restart, remove, status, logs)

可选参数:
    -e, --env-file      环境变量文件路径
    -p, --project       项目名称 (用于多实例部署)
    -d, --data-dir      数据目录
    --port              端口映射
    --password          密码

示例:
    # 部署 MySQL 8.0
    $0 mysql mysql-8.0 deploy -e /path/to/.env

    # 部署 Redis 7.2 并指定端口
    $0 redis redis-7.2 deploy --port 6380

    # 停止服务
    $0 mysql mysql-8.0 stop

    # 查看日志
    $0 mysql mysql-8.0 logs

    # 多实例部署
    $0 mysql mysql-8.0 deploy -p app-db-1 -e app1.env
    $0 mysql mysql-8.0 deploy -p app-db-2 -e app2.env

支持的服务:
    数据库: mysql, postgresql, redis
    CI/CD: jenkins, gitlab
    监控: prometheus, grafana, elasticsearch, loki
    开发环境: nodejs, java, golang
    其他: jumpserver, fluentd

EOF
    exit 1
}

# 检查参数
if [ $# -lt 3 ]; then
    usage
fi

SERVICE=$1
VERSION=$2
ACTION=$3
shift 3

# 默认配置
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
COMPOSE_FILE="${SCRIPT_DIR}/${SERVICE}/versions/${VERSION}-docker-compose.yml"
ENV_FILE=""
PROJECT_NAME="${SERVICE}"
DATA_DIR="/data/${SERVICE}"
EXTRA_ARGS=""

# 解析可选参数
while [[ $# -gt 0 ]]; do
    case $1 in
        -e|--env-file)
            ENV_FILE="$2"
            shift 2
            ;;
        -p|--project)
            PROJECT_NAME="$2"
            shift 2
            ;;
        -d|--data-dir)
            DATA_DIR="$2"
            EXTRA_ARGS="$EXTRA_ARGS DATA_DIR=$DATA_DIR"
            shift 2
            ;;
        --port)
            EXTRA_ARGS="$EXTRA_ARGS PORT=$2"
            shift 2
            ;;
        --password)
            EXTRA_ARGS="$EXTRA_ARGS PASSWORD=$2"
            shift 2
            ;;
        *)
            log_error "未知参数: $1"
            usage
            ;;
    esac
done

# 检查 docker-compose 文件是否存在
if [ ! -f "$COMPOSE_FILE" ]; then
    log_error "Compose 文件不存在: $COMPOSE_FILE"
    log_info "可用版本:"
    ls -1 "${SCRIPT_DIR}/${SERVICE}/versions/" 2>/dev/null || log_error "服务 ${SERVICE} 不存在"
    exit 1
fi

# 构建 docker-compose 命令
COMPOSE_CMD="docker-compose -f ${COMPOSE_FILE} -p ${PROJECT_NAME}"

# 添加环境变量文件
if [ -n "$ENV_FILE" ] && [ -f "$ENV_FILE" ]; then
    COMPOSE_CMD="$COMPOSE_CMD --env-file $ENV_FILE"
fi

# 执行操作
log_info "服务: $SERVICE"
log_info "版本: $VERSION"
log_info "操作: $ACTION"
log_info "项目名称: $PROJECT_NAME"
log_info "Compose文件: $COMPOSE_FILE"

case $ACTION in
    deploy|up)
        log_info "开始部署服务..."

        # 创建数据目录
        mkdir -p "$DATA_DIR"
        log_info "数据目录: $DATA_DIR"

        # 启动服务
        eval "$EXTRA_ARGS $COMPOSE_CMD up -d"

        log_info "服务部署成功!"
        log_info "查看状态: $0 $SERVICE $VERSION status -p $PROJECT_NAME"
        log_info "查看日志: $0 $SERVICE $VERSION logs -p $PROJECT_NAME"
        ;;

    stop)
        log_info "停止服务..."
        eval "$COMPOSE_CMD stop"
        log_info "服务已停止"
        ;;

    start)
        log_info "启动服务..."
        eval "$COMPOSE_CMD start"
        log_info "服务已启动"
        ;;

    restart)
        log_info "重启服务..."
        eval "$COMPOSE_CMD restart"
        log_info "服务已重启"
        ;;

    remove|down)
        log_warn "即将删除服务容器 (数据目录将保留)"
        read -p "确认删除? (yes/no): " confirm
        if [ "$confirm" = "yes" ]; then
            eval "$COMPOSE_CMD down"
            log_info "服务已删除"
        else
            log_info "取消删除"
        fi
        ;;

    status|ps)
        eval "$COMPOSE_CMD ps"
        ;;

    logs)
        eval "$COMPOSE_CMD logs -f --tail=100"
        ;;

    exec)
        CONTAINER_NAME="${PROJECT_NAME}-${SERVICE}"
        log_info "进入容器: $CONTAINER_NAME"
        docker exec -it "$CONTAINER_NAME" /bin/sh || docker exec -it "$CONTAINER_NAME" /bin/bash
        ;;

    *)
        log_error "未知操作: $ACTION"
        usage
        ;;
esac

exit 0
