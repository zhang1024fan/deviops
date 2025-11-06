#!/bin/bash
# Node.js 20 LTS 安装脚本
# 参数1: 安装路径 (如: /usr/local/node20)
# 参数2: 环境变量文件名 (如: nodejs20)

set -e

INSTALL_PATH="${1:-/usr/local/node20}"
PROFILE_NAME="${2:-nodejs20}"

echo "===== Node.js 20 LTS 安装配置 ====="
echo "安装路径: $INSTALL_PATH"
echo "环境变量文件: /etc/profile.d/$PROFILE_NAME.sh"

PROFILE_FILE="/etc/profile.d/$PROFILE_NAME.sh"

cat > "$PROFILE_FILE" << ENVEOF
# Node.js Environment
export NODE_HOME=$INSTALL_PATH
export PATH=\$PATH:\$NODE_HOME/bin
ENVEOF

source "$PROFILE_FILE"

echo "环境变量已配置: $PROFILE_FILE"
if command -v node &> /dev/null; then
    echo "===== 安装成功 ====="
    node -v
    npm -v
else
    echo "===== 安装失败 ====="
    exit 1
fi
