#!/usr/bin/env bash

# AutoOps 8 服务一键部署脚本
# 用法: ./devops-start.sh <version> <server_host> <web_port> [api_port] [mysql_port] [redis_port] [victoriametrics_port] [victorialogs_port] [kafka_port]
# 示例: ./devops-start.sh v5.0 192.168.1.100 8088

set -Eeuo pipefail

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

REGISTRY="crpi-aj3vgoxp9kzh2jx4.cn-hangzhou.personal.cr.aliyuncs.com/zhangfan_k8s"
API_REPOSITORY="${AUTOOPS_API_REPOSITORY:-${REGISTRY}/deviops-api}"
WEB_REPOSITORY="${AUTOOPS_WEB_REPOSITORY:-${REGISTRY}/deviops-web}"
KAFKA_TOPIC="aiops-logs"
KAFKA_PARTITIONS=6
WAIT_TIMEOUT="${AUTOOPS_WAIT_TIMEOUT:-300}"
PODMAN_DOCKER_NOTICE="Emulate Docker CLI using podman. Create /etc/containers/nodocker to quiet msg."

docker_cli() {
    docker "$@" 2> >(
        while IFS= read -r line || [ -n "$line" ]; do
            [ "$line" = "$PODMAN_DOCKER_NOTICE" ] || printf '%s\n' "$line" >&2
        done
    )
}

usage() {
    echo "用法: $0 <version> <server_host> <web_port> [api_port] [mysql_port] [redis_port] [victoriametrics_port] [victorialogs_port] [kafka_port]"
    echo ""
    echo "参数说明:"
    echo "  version               API/Web 镜像版本，例如 v5.0"
    echo "  server_host           Agent 可访问的服务器 IP 或域名"
    echo "  web_port              Web 访问端口"
    echo "  api_port              API 访问端口，默认 8000"
    echo "  mysql_port            MySQL 映射端口，默认 3307"
    echo "  redis_port            Redis 映射端口，默认 6379"
    echo "  victoriametrics_port  VictoriaMetrics 端口，默认 8428"
    echo "  victorialogs_port     VictoriaLogs 端口，默认 9428"
    echo "  kafka_port            Kafka 外部端口，默认 9092"
    echo ""
    echo "示例:"
    echo "  $0 v5.0 192.168.1.100 8088"
    echo "  $0 v5.0 autoops.example.com 80 8000 3307 6379 8428 9428 9092"
}

fail() {
    echo -e "${RED}错误: $*${NC}" >&2
    exit 1
}

on_error() {
    local exit_code=$?
    echo -e "${RED}部署失败，出错行: ${BASH_LINENO[0]}，退出码: ${exit_code}${NC}" >&2
    exit "$exit_code"
}
trap on_error ERR

if [ "$#" -lt 3 ] || [ "$#" -gt 9 ]; then
    usage
    exit 1
fi

VERSION="$1"
SERVER_HOST="$2"
WEB_PORT="$3"
API_PORT="${4:-8000}"
MYSQL_PORT="${5:-3307}"
REDIS_PORT="${6:-6379}"
VICTORIAMETRICS_PORT="${7:-8428}"
VICTORIALOGS_PORT="${8:-9428}"
KAFKA_PORT="${9:-9092}"

[[ "$VERSION" =~ ^[A-Za-z0-9][A-Za-z0-9._-]*$ ]] || fail "镜像版本格式无效: $VERSION"
[[ "$SERVER_HOST" =~ ^[A-Za-z0-9][A-Za-z0-9.-]*$ ]] || fail "服务器 IP 或域名格式无效: $SERVER_HOST"
[[ "$WAIT_TIMEOUT" =~ ^[0-9]+$ ]] || fail "AUTOOPS_WAIT_TIMEOUT 必须是正整数"

PORTS=("$WEB_PORT" "$API_PORT" "$MYSQL_PORT" "$REDIS_PORT" "$VICTORIAMETRICS_PORT" "$VICTORIALOGS_PORT" "$KAFKA_PORT")
for port in "${PORTS[@]}"; do
    [[ "$port" =~ ^[0-9]+$ ]] || fail "端口号无效: $port"
    (( port >= 1 && port <= 65535 )) || fail "端口号超出范围: $port"
done
for ((i = 0; i < ${#PORTS[@]}; i++)); do
    for ((j = i + 1; j < ${#PORTS[@]}; j++)); do
        [ "${PORTS[$i]}" != "${PORTS[$j]}" ] || fail "端口不能重复: ${PORTS[$i]}"
    done
done

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

COMPOSE_FILE="$SCRIPT_DIR/docker-compose.yml"
ENV_FILE="$SCRIPT_DIR/.env"
CONFIG_TEMPLATE="$SCRIPT_DIR/api/config.template.yaml"
RUNTIME_CONFIG="$SCRIPT_DIR/api/config.runtime.yaml"
VECTOR_CONFIG="$SCRIPT_DIR/vector/vector.yaml"

[ -f "$COMPOSE_FILE" ] || fail "找不到 $COMPOSE_FILE"
[ -f "$CONFIG_TEMPLATE" ] || fail "找不到 $CONFIG_TEMPLATE"
[ -f "$VECTOR_CONFIG" ] || fail "找不到 $VECTOR_CONFIG"
touch "$ENV_FILE"

command -v docker >/dev/null 2>&1 || fail "未安装 Docker"
docker_cli info >/dev/null 2>&1 || fail "Docker/Podman 未运行或当前用户无访问权限"

if docker_cli compose version >/dev/null 2>&1; then
    COMPOSE=(docker_cli compose)
    COMPOSE_DISPLAY="docker compose"
elif command -v docker-compose >/dev/null 2>&1; then
    COMPOSE=(docker-compose)
    COMPOSE_DISPLAY="docker-compose"
else
    fail "找不到 docker compose 或 docker-compose"
fi

get_env_value() {
    local key="$1"
    awk -v key="$key" '
        index($0, key "=") == 1 {
            print substr($0, length(key) + 2)
            exit
        }
    ' "$ENV_FILE"
}

set_env_value() {
    local key="$1"
    local value="$2"
    local temp_file
    temp_file="$(mktemp "${SCRIPT_DIR}/.env.XXXXXX")"
    ENV_KEY="$key" ENV_VALUE="$value" awk '
        BEGIN {
            key = ENVIRON["ENV_KEY"]
            value = ENVIRON["ENV_VALUE"]
            found = 0
        }
        index($0, key "=") == 1 {
            if (!found) {
                print key "=" value
                found = 1
            }
            next
        }
        { print }
        END {
            if (!found) {
                print key "=" value
            }
        }
    ' "$ENV_FILE" > "$temp_file"
    mv "$temp_file" "$ENV_FILE"
}

yaml_escape() {
    local value="$1"
    value="${value//\\/\\\\}"
    value="${value//\"/\\\"}"
    printf '%s' "$value"
}

MYSQL_ROOT_PASSWORD_VALUE="${MYSQL_ROOT_PASSWORD:-$(get_env_value MYSQL_ROOT_PASSWORD)}"
MYSQL_ROOT_PASSWORD_VALUE="${MYSQL_ROOT_PASSWORD_VALUE:-devops@2025}"
MYSQL_DATABASE_VALUE="${MYSQL_DATABASE:-$(get_env_value MYSQL_DATABASE)}"
MYSQL_DATABASE_VALUE="${MYSQL_DATABASE_VALUE:-autoops}"
REDIS_PASSWORD_VALUE="${REDIS_PASSWORD:-$(get_env_value REDIS_PASSWORD)}"
REDIS_PASSWORD_VALUE="${REDIS_PASSWORD_VALUE:-devops@2025}"
MYSQL_BOOTSTRAP_CHECK_STORED="$(get_env_value AUTOOPS_MYSQL_BOOTSTRAP_CHECK)"

if [ -n "${AUTOOPS_MYSQL_BOOTSTRAP_CHECK:-}" ]; then
    MYSQL_BOOTSTRAP_CHECK_VALUE="$AUTOOPS_MYSQL_BOOTSTRAP_CHECK"
elif [ ! -d "$SCRIPT_DIR/mysql/data/mysql" ]; then
    MYSQL_BOOTSTRAP_CHECK_VALUE="required"
elif [ "$MYSQL_BOOTSTRAP_CHECK_STORED" = "required" ]; then
    # Keep strict verification after an interrupted first initialization.
    MYSQL_BOOTSTRAP_CHECK_VALUE="required"
else
    # Existing installations predate the completion marker and must remain usable.
    MYSQL_BOOTSTRAP_CHECK_VALUE="legacy"
fi

[[ "$MYSQL_DATABASE_VALUE" =~ ^[A-Za-z0-9_]+$ ]] || fail "MYSQL_DATABASE 只能包含字母、数字和下划线"
[[ "$MYSQL_ROOT_PASSWORD_VALUE" != *$'\n'* ]] || fail "MYSQL_ROOT_PASSWORD 不能包含换行"
[[ "$REDIS_PASSWORD_VALUE" != *$'\n'* ]] || fail "REDIS_PASSWORD 不能包含换行"
[[ "$MYSQL_BOOTSTRAP_CHECK_VALUE" = "required" || "$MYSQL_BOOTSTRAP_CHECK_VALUE" = "legacy" ]] || \
    fail "AUTOOPS_MYSQL_BOOTSTRAP_CHECK 只能是 required 或 legacy"

AUTOOPS_API_IMAGE="${API_REPOSITORY}:${VERSION}"
AUTOOPS_WEB_IMAGE="${WEB_REPOSITORY}:${VERSION}"

set_env_value AUTOOPS_API_IMAGE "$AUTOOPS_API_IMAGE"
set_env_value AUTOOPS_WEB_IMAGE "$AUTOOPS_WEB_IMAGE"
set_env_value API_CONFIG_FILE "./api/config.runtime.yaml"
set_env_value SERVER_HOST "$SERVER_HOST"
set_env_value WEB_PORT "$WEB_PORT"
set_env_value API_PORT "$API_PORT"
set_env_value MYSQL_PORT "$MYSQL_PORT"
set_env_value REDIS_PORT "$REDIS_PORT"
set_env_value VICTORIAMETRICS_PORT "$VICTORIAMETRICS_PORT"
set_env_value VICTORIALOGS_PORT "$VICTORIALOGS_PORT"
set_env_value KAFKA_PORT "$KAFKA_PORT"
set_env_value KAFKA_ADVERTISED_HOST "$SERVER_HOST"
set_env_value MYSQL_ROOT_PASSWORD "$MYSQL_ROOT_PASSWORD_VALUE"
set_env_value MYSQL_DATABASE "$MYSQL_DATABASE_VALUE"
set_env_value REDIS_PASSWORD "$REDIS_PASSWORD_VALUE"
set_env_value AUTOOPS_MYSQL_BOOTSTRAP_CHECK "$MYSQL_BOOTSTRAP_CHECK_VALUE"

render_api_config() {
    local temp_file
    temp_file="$(mktemp "${SCRIPT_DIR}/api/.config.runtime.XXXXXX")"

    CFG_SERVER_HOST="$SERVER_HOST" \
    CFG_WEB_PORT="$WEB_PORT" \
    CFG_API_PORT="$API_PORT" \
    CFG_KAFKA_PORT="$KAFKA_PORT" \
    CFG_VICTORIAMETRICS_PORT="$VICTORIAMETRICS_PORT" \
    CFG_MYSQL_DATABASE="$MYSQL_DATABASE_VALUE" \
    CFG_MYSQL_ROOT_PASSWORD="$(yaml_escape "$MYSQL_ROOT_PASSWORD_VALUE")" \
    CFG_REDIS_PASSWORD="$(yaml_escape "$REDIS_PASSWORD_VALUE")" \
    awk '
        function replace_all(text, needle, replacement, position, output) {
            output = ""
            while ((position = index(text, needle)) > 0) {
                output = output substr(text, 1, position - 1) replacement
                text = substr(text, position + length(needle))
            }
            return output text
        }
        {
            line = $0
            line = replace_all(line, "__SERVER_HOST__", ENVIRON["CFG_SERVER_HOST"])
            line = replace_all(line, "__WEB_PORT__", ENVIRON["CFG_WEB_PORT"])
            line = replace_all(line, "__API_PORT__", ENVIRON["CFG_API_PORT"])
            line = replace_all(line, "__KAFKA_PORT__", ENVIRON["CFG_KAFKA_PORT"])
            line = replace_all(line, "__VICTORIAMETRICS_PORT__", ENVIRON["CFG_VICTORIAMETRICS_PORT"])
            line = replace_all(line, "__MYSQL_DATABASE__", ENVIRON["CFG_MYSQL_DATABASE"])
            line = replace_all(line, "__MYSQL_ROOT_PASSWORD__", ENVIRON["CFG_MYSQL_ROOT_PASSWORD"])
            line = replace_all(line, "__REDIS_PASSWORD__", ENVIRON["CFG_REDIS_PASSWORD"])
            print line
        }
    ' "$CONFIG_TEMPLATE" > "$temp_file"

    chmod 600 "$temp_file"
    mv "$temp_file" "$RUNTIME_CONFIG"
}

prepare_persistent_paths() {
    mkdir -p api/logs api/upload api/ssh_keys mysql/data redis/data
    chmod 700 api/ssh_keys
    find api/ssh_keys -type f -name '*.pub' -exec chmod 644 {} +
    find api/ssh_keys -type f ! -name '*.pub' -exec chmod 600 {} +
}

migrate_legacy_autoops_data() {
    local legacy_data_dir="$SCRIPT_DIR/api/data"

    if [ ! -f "$legacy_data_dir/instance.id" ] && [ ! -f "$legacy_data_dir/license.key" ]; then
        return 0
    fi

    echo -e "${YELLOW}检测到旧 AutoOps 数据，迁移到 autoops-data 数据卷...${NC}"
    "${COMPOSE[@]}" run --rm --no-deps \
        --entrypoint /bin/sh \
        -v "$legacy_data_dir:/legacy-autoops-data:ro" \
        autoops-api -ec '
            mkdir -p /var/lib/autoops
            for name in instance.id license.key; do
                if [ -f "/legacy-autoops-data/$name" ] && [ ! -e "/var/lib/autoops/$name" ]; then
                    cp "/legacy-autoops-data/$name" "/var/lib/autoops/$name"
                    chmod 600 "/var/lib/autoops/$name"
                fi
            done
        '
    echo -e "${GREEN}✓ 旧机器码和许可证迁移完成${NC}"
}

wait_for_container() {
    local container="$1"
    local deadline=$((SECONDS + WAIT_TIMEOUT))
    local status health

    while ((SECONDS < deadline)); do
        if ! docker_cli inspect "$container" >/dev/null 2>&1; then
            sleep 2
            continue
        fi

        status="$(docker_cli inspect --format '{{.State.Status}}' "$container")"
        health="$(docker_cli inspect --format '{{if .State.Health}}{{.State.Health.Status}}{{else}}none{{end}}' "$container")"

        if [ "$status" = "running" ] && { [ "$health" = "none" ] || [ "$health" = "healthy" ]; }; then
            echo -e "${GREEN}✓ $container 已就绪${NC}"
            return 0
        fi
        if [ "$status" = "exited" ] || [ "$status" = "dead" ] || [ "$health" = "unhealthy" ]; then
            echo -e "${RED}✗ $container 启动失败 (status=$status, health=$health)${NC}"
            if [ "$health" != "none" ]; then
                echo -e "${YELLOW}最近的健康检查结果:${NC}"
                docker_cli inspect --format \
                    '{{range .State.Health.Log}}{{.End}} exit={{.ExitCode}}: {{.Output}}{{println}}{{end}}' \
                    "$container" 2>&1 || true
            fi
            docker_cli logs --tail 80 "$container" 2>&1 || true
            return 1
        fi
        sleep 2
    done

    echo -e "${RED}✗ 等待 $container 超时${NC}"
    health="$(docker_cli inspect --format '{{if .State.Health}}{{.State.Health.Status}}{{else}}none{{end}}' "$container" 2>/dev/null || echo none)"
    if [ "$health" != "none" ]; then
        echo -e "${YELLOW}最近的健康检查结果:${NC}"
        docker_cli inspect --format \
            '{{range .State.Health.Log}}{{.End}} exit={{.ExitCode}}: {{.Output}}{{println}}{{end}}' \
            "$container" 2>&1 || true
    fi
    docker_cli logs --tail 80 "$container" 2>&1 || true
    return 1
}

wait_for_kafka() {
    local deadline=$((SECONDS + WAIT_TIMEOUT))
    while ((SECONDS < deadline)); do
        if docker_cli exec autoops-kafka /opt/kafka/bin/kafka-topics.sh \
            --bootstrap-server autoops-kafka:29092 --list >/dev/null 2>&1; then
            echo -e "${GREEN}✓ autoops-kafka 已就绪${NC}"
            return 0
        fi
        if ! docker_cli inspect --format '{{.State.Running}}' autoops-kafka 2>/dev/null | grep -q true; then
            docker_cli logs --tail 80 autoops-kafka 2>&1 || true
            return 1
        fi
        sleep 3
    done

    echo -e "${RED}✗ 等待 autoops-kafka 超时${NC}"
    docker_cli logs --tail 80 autoops-kafka 2>&1 || true
    return 1
}

ensure_kafka_topic() {
    local description partition_count

    docker_cli exec autoops-kafka /opt/kafka/bin/kafka-topics.sh \
        --bootstrap-server autoops-kafka:29092 \
        --create --if-not-exists \
        --topic "$KAFKA_TOPIC" \
        --partitions "$KAFKA_PARTITIONS" \
        --replication-factor 1

    description="$(docker_cli exec autoops-kafka /opt/kafka/bin/kafka-topics.sh \
        --bootstrap-server autoops-kafka:29092 \
        --describe --topic "$KAFKA_TOPIC")"
    partition_count="$(printf '%s\n' "$description" | awk -F'PartitionCount: ' 'NF > 1 {split($2, fields, " "); print fields[1]; exit}')"

    if [[ "$partition_count" =~ ^[0-9]+$ ]] && ((partition_count < KAFKA_PARTITIONS)); then
        docker_cli exec autoops-kafka /opt/kafka/bin/kafka-topics.sh \
            --bootstrap-server autoops-kafka:29092 \
            --alter --topic "$KAFKA_TOPIC" \
            --partitions "$KAFKA_PARTITIONS"
        partition_count="$KAFKA_PARTITIONS"
    fi

    [ "$partition_count" = "$KAFKA_PARTITIONS" ] || fail "Kafka Topic $KAFKA_TOPIC 分区数异常: ${partition_count:-unknown}"
    echo -e "${GREEN}✓ Kafka Topic $KAFKA_TOPIC 已就绪 (${KAFKA_PARTITIONS} partitions)${NC}"
}

render_api_config
prepare_persistent_paths

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}AutoOps 8 服务一键部署${NC}"
echo -e "${GREEN}========================================${NC}"
echo "API 镜像:       $AUTOOPS_API_IMAGE"
echo "Web 镜像:       $AUTOOPS_WEB_IMAGE"
echo "服务器地址:     $SERVER_HOST"
echo "Web/API:        $WEB_PORT / $API_PORT"
echo "MySQL/Redis:    $MYSQL_PORT / $REDIS_PORT"
echo "MySQL 初始化校验: $MYSQL_BOOTSTRAP_CHECK_VALUE"
echo "VM/VL/Kafka:    $VICTORIAMETRICS_PORT / $VICTORIALOGS_PORT / $KAFKA_PORT"
echo ""

echo -e "${YELLOW}[1/6] 校验 Compose 和运行时配置...${NC}"
"${COMPOSE[@]}" config >/dev/null
echo -e "${GREEN}✓ 配置校验通过${NC}"

echo -e "${YELLOW}[2/6] 拉取 8 个服务镜像...${NC}"
if ! "${COMPOSE[@]}" pull; then
    fail "镜像拉取失败，请确认已执行 docker login $REGISTRY 且版本 $VERSION 已推送"
fi

echo -e "${YELLOW}[3/6] 停止旧栈并启动基础组件...${NC}"
"${COMPOSE[@]}" down --remove-orphans >/dev/null 2>&1 || true
migrate_legacy_autoops_data
"${COMPOSE[@]}" up -d --remove-orphans \
    autoops-mysql \
    autoops-redis \
    autoops-victoriametrics \
    autoops-victorialogs \
    autoops-kafka

wait_for_container autoops-mysql
wait_for_container autoops-redis
wait_for_container autoops-victoriametrics
wait_for_container autoops-victorialogs
wait_for_kafka

echo -e "${YELLOW}[4/6] 初始化 Kafka Topic...${NC}"
ensure_kafka_topic

echo -e "${YELLOW}[5/6] 启动 Vector、API 和 Web...${NC}"
"${COMPOSE[@]}" up -d --remove-orphans \
    autoops-vector-aggregator \
    autoops-api \
    autoops-web

echo -e "${YELLOW}[6/6] 检查全部服务...${NC}"
SERVICES=(
    autoops-mysql
    autoops-redis
    autoops-victoriametrics
    autoops-victorialogs
    autoops-kafka
    autoops-vector-aggregator
    autoops-api
    autoops-web
)
for service in "${SERVICES[@]}"; do
    wait_for_container "$service"
done

echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}✓ AutoOps 8 个服务部署完成${NC}"
echo -e "${GREEN}========================================${NC}"
echo "Web:             http://$SERVER_HOST:$WEB_PORT"
echo "API:             http://$SERVER_HOST:$API_PORT"
echo "Swagger:         http://$SERVER_HOST:$API_PORT/swagger/index.html"
echo "VictoriaMetrics: http://$SERVER_HOST:$VICTORIAMETRICS_PORT"
echo "VictoriaLogs:    http://$SERVER_HOST:$VICTORIALOGS_PORT/select/vmui/"
echo "Kafka:           $SERVER_HOST:$KAFKA_PORT"
echo "MySQL:           $SERVER_HOST:$MYSQL_PORT"
echo "Redis:           $SERVER_HOST:$REDIS_PORT"
echo ""
echo "查看状态: $COMPOSE_DISPLAY ps"
echo "查看日志: $COMPOSE_DISPLAY logs -f <service>"
