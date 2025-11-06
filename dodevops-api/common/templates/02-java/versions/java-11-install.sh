#!/bin/bash
# Java 11 LTS 安装脚本
# 从Docker镜像提取的二进制文件已就绪，现在配置环境变量
# 参数1: 安装路径 (如: /usr/local/java11)
# 参数2: 环境变量文件名 (如: java11)

set -e

INSTALL_PATH="${1:-/usr/local/java11}"
PROFILE_NAME="${2:-java11}"

echo "===== Java 11 LTS 安装配置 ====="
echo "安装路径: $INSTALL_PATH"
echo "环境变量文件: /etc/profile.d/$PROFILE_NAME.sh"

# 配置环境变量
PROFILE_FILE="/etc/profile.d/$PROFILE_NAME.sh"

echo "配置 Java 环境变量..."
cat > "$PROFILE_FILE" << EOF
# Java Environment
export JAVA_HOME=$INSTALL_PATH
export PATH=\$PATH:\$JAVA_HOME/bin
export CLASSPATH=.:\$JAVA_HOME/lib
EOF

# 使环境变量立即生效
source "$PROFILE_FILE"

echo "环境变量已配置: $PROFILE_FILE"
echo "JAVA_HOME=$INSTALL_PATH"
echo "PATH已添加: \$JAVA_HOME/bin"

# 验证安装
if command -v java &> /dev/null; then
    echo "===== 安装成功 ====="
    java -version
else
    echo "===== 安装失败：找不到 java 命令 ====="
    exit 1
fi

echo ""
echo "提示: 请执行 'source /etc/profile.d/$PROFILE_NAME.sh' 或重新登录使环境变量生效"
