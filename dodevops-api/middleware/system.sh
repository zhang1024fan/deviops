#!/bin/bash

# 获取私有 IP（默认网卡）
privateIp=$(hostname -I | awk '{print $1}' || echo "unknown")
# 获取公网 IP（如果存在）
publicIp=$(curl -s ipinfo.io/ip 2>/dev/null || echo "")
# 获取操作系统版本
os=$(cat /etc/os-release 2>/dev/null | grep PRETTY_NAME | cut -d= -f2 | tr -d '"' | sed 's/ LTS//;s/ //g' || echo "unknown")
# 获取 CPU 核心数
cpu=$(nproc 2>/dev/null || echo "unknown")
# 获取内存大小（单位为 MB，转换为 G）
memory=$(free -m 2>/dev/null | awk '/^Mem:/{printf "%.0fG\n",$2/1024}' || echo "unknown")
# 获取磁盘总容量（单位为 G）
disk=$(df -h 2>/dev/null | awk '/\/$/ {print $2}' || echo "unknown")
# 输出 JSON 格式（确保只输出一次）
json_output=$(cat <<EOF
{"privateIp":"$privateIp","publicIp":"$publicIp","os":"$os","cpu":"$cpu","memory":"$memory","disk":"$disk"}
EOF
)

# 验证并输出JSON
if jq -e . >/dev/null 2>&1 <<<"$json_output"; then
    echo "SSH配置获取成功"
    echo "$json_output"
else
    echo "SSH配置获取失败"
    echo '{"error":"invalid json generated"}'
fi
