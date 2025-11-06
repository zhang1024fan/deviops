#!/bin/bash
# Golang 1.21 安装脚本
# 参数1: 安装路径 (如: /usr/local/go1.21)
# 参数2: 环境变量文件名 (如: golang1.21)

set -e

INSTALL_PATH="${1:-/usr/local/go1.21}"
PROFILE_NAME="${2:-golang1.21}"

echo "===== Golang 1.21 安装配置 ====="
echo "安装路径: $INSTALL_PATH"
echo "环境变量文件: /etc/profile.d/$PROFILE_NAME.sh"

PROFILE_FILE="/etc/profile.d/$PROFILE_NAME.sh"

cat > "$PROFILE_FILE" << ENVEOF
# Golang Environment
export GOROOT=$INSTALL_PATH
export PATH=\$PATH:\$GOROOT/bin
export GOPATH=\$HOME/go
export PATH=\$PATH:\$GOPATH/bin
ENVEOF

source "$PROFILE_FILE"

echo "环境变量已配置: $PROFILE_FILE"
if command -v go &> /dev/null; then
    echo "===== 安装成功 ====="
    go version
else
    echo "===== 安装失败 ====="
    exit 1
fi
