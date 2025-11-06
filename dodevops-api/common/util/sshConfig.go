package util

import "time"

// SSHConfig SSH连接配置
type SSHConfig struct {
	IP        string        // 主机IP
	Port      int           // SSH端口
	Type      int           // 认证类型:1->密码,2->私钥,3->公钥(免认证)
	Username  string        // 用户名
	Password  string        // 密码(type=1时使用)
	PublicKey string        // 私钥内容(type=2时使用) - 注意：字段名为PublicKey但实际存储私钥
	Timeout   time.Duration // 超时时间
}
