#!/bin/bash
# Java 17 LTS 安装脚本
# 参数1: 安装路径 (如: /usr/local/java17)
# 参数2: 环境变量文件名 (如: java17)

set -e

INSTALL_PATH="${1:-/usr/local/java17}"
PROFILE_NAME="${2:-java17}"

echo "===== Java 17 LTS 安装配置 ====="
echo "安装路径: $INSTALL_PATH"
echo "环境变量文件: /etc/profile.d/$PROFILE_NAME.sh"

PROFILE_FILE="/etc/profile.d/$PROFILE_NAME.sh"

cat > "$PROFILE_FILE" << ENVEOF
# Java Environment
export JAVA_HOME=$INSTALL_PATH
export PATH=\$PATH:\$JAVA_HOME/bin
export CLASSPATH=.:\$JAVA_HOME/lib
ENVEOF

source "$PROFILE_FILE"

echo "环境变量已配置: $PROFILE_FILE"
if command -v java &> /dev/null; then
    echo "===== 安装成功 ====="
    java -version
else
    echo "===== 安装失败 ====="
    exit 1
fi
