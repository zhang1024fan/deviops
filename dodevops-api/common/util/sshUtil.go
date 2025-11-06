package util

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
)

type SSHUtil struct{}

// NewSSHUtil 创建SSHUtil实例
func NewSSHUtil() *SSHUtil {
	return &SSHUtil{}
}

// ExecuteRemoteCommand 执行远程命令并返回输出
func (s *SSHUtil) ExecuteRemoteCommand(auth *SSHConfig, command string) (string, error) {
	config, err := s.getSSHConfig(auth)
	if err != nil {
		return "", fmt.Errorf("failed to create SSH config: %v", err)
	}

	// 建立SSH连接
	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", auth.IP, auth.Port), config)
	if err != nil {
		fmt.Printf("SSH连接失败: %v\n", err)
		return "", fmt.Errorf("failed to dial: %v", err)
	}
	defer conn.Close()

	// 创建会话
	session, err := conn.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create session: %v", err)
	}
	defer session.Close()

	// 设置输出缓冲区
	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf

	// 执行命令
	err = session.Run(command)
	if err != nil {
		return "", fmt.Errorf("failed to run command: %v", err)
	}

	return stdoutBuf.String(), nil
}

// getSSHConfig 根据认证类型创建SSH配置
func (s *SSHUtil) getSSHConfig(auth *SSHConfig) (*ssh.ClientConfig, error) {
	var authMethods []ssh.AuthMethod

	switch auth.Type {
	case 1: // 密码认证
		authMethods = []ssh.AuthMethod{ssh.Password(auth.Password)}
	case 2: // 私钥认证 - 注意：PublicKey字段实际存储的是私钥内容
		authMethod, err := s.PrivateKeyAuth(auth.PublicKey) // 这里的PublicKey实际是私钥
		if err != nil {
			return nil, err
		}
		authMethods = []ssh.AuthMethod{authMethod}
	case 3: // 公钥免认证 - 公钥已提前拷贝到服务器，自动查找本地私钥
		// 按顺序尝试找到可用的私钥
		// 1. 先尝试用户主目录
		if userKeyAuth, err := s.UserKeyAuth(); err == nil {
			authMethods = append(authMethods, userKeyAuth)
		} else {
			// 2. 再尝试系统默认路径
			if defaultKeyAuth, err := s.DefaultKeyAuth(); err == nil {
				authMethods = append(authMethods, defaultKeyAuth)
			}
		}

		// 如果没有找到任何私钥，返回清晰的错误信息
		if len(authMethods) == 0 {
			return nil, errors.New("type=3 公钥免认证失败: 本地未找到私钥文件，请确保以下路径之一存在私钥:\n" +
				"用户目录: ~/.ssh/id_rsa, ~/.ssh/id_ed25519, ~/.ssh/id_ecdsa\n" +
				"系统目录: /root/.ssh/id_rsa, /root/.ssh/id_ed25519, /root/.ssh/id_ecdsa")
		}
	default:
		return nil, errors.New("unsupported authentication type")
	}

	return &ssh.ClientConfig{
		User: auth.Username,
		Auth: authMethods,
		HostKeyCallback: ssh.HostKeyCallback(func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil // 忽略主机密钥检查
		}),
		Timeout: 30 * time.Second,
	}, nil
}

// PrivateKeyAuth 使用私钥创建认证方法
func (s *SSHUtil) PrivateKeyAuth(privateKey string) (ssh.AuthMethod, error) {
	// 尝试解析私钥
	signer, err := ssh.ParsePrivateKey([]byte(privateKey))
	if err != nil {
		// 尝试添加密码为空的情况
		signer, err = ssh.ParsePrivateKeyWithPassphrase([]byte(privateKey), []byte(""))
		if err != nil {
			return nil, fmt.Errorf("failed to parse private key (with/without passphrase): %v", err)
		}
	}

	return ssh.PublicKeys(signer), nil
}

// PublicKeyAuth 保持向后兼容，实际调用PrivateKeyAuth
func (s *SSHUtil) PublicKeyAuth(privateKey string) (ssh.AuthMethod, error) {
	return s.PrivateKeyAuth(privateKey)
}

// DefaultKeyAuth 尝试使用默认路径的私钥文件
func (s *SSHUtil) DefaultKeyAuth() (ssh.AuthMethod, error) {
	// 常见的默认私钥路径
	defaultKeyPaths := []string{
		"/root/.ssh/id_rsa",
		"/root/.ssh/id_ed25519",
		"/root/.ssh/id_ecdsa",
	}

	for _, keyPath := range defaultKeyPaths {
		if keyData, err := os.ReadFile(keyPath); err == nil {
			// 尝试解析私钥
			if signer, err := ssh.ParsePrivateKey(keyData); err == nil {
				return ssh.PublicKeys(signer), nil
			}
			// 尝试无密码解析
			if signer, err := ssh.ParsePrivateKeyWithPassphrase(keyData, []byte("")); err == nil {
				return ssh.PublicKeys(signer), nil
			}
		}
	}

	return nil, errors.New("no default SSH key found")
}

// UserKeyAuth 尝试使用用户主目录的私钥文件
func (s *SSHUtil) UserKeyAuth() (ssh.AuthMethod, error) {
	// 获取当前用户主目录
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %v", err)
	}

	// 用户主目录下的私钥路径
	userKeyPaths := []string{
		homeDir + "/.ssh/id_rsa",
		homeDir + "/.ssh/id_ed25519",
		homeDir + "/.ssh/id_ecdsa",
	}

	for _, keyPath := range userKeyPaths {
		if keyData, err := os.ReadFile(keyPath); err == nil {
			// 尝试解析私钥
			if signer, err := ssh.ParsePrivateKey(keyData); err == nil {
				return ssh.PublicKeys(signer), nil
			}
			// 尝试无密码解析
			if signer, err := ssh.ParsePrivateKeyWithPassphrase(keyData, []byte("")); err == nil {
				return ssh.PublicKeys(signer), nil
			}
		}
	}

	return nil, errors.New("no user SSH key found")
}

// ExecuteScript 执行远程脚本
func (s *SSHUtil) ExecuteScript(auth *SSHConfig, script string) (string, error) {
	return s.ExecuteRemoteCommand(auth, script)
}

// TerminalLogin 建立SSH终端连接
func (s *SSHUtil) TerminalLogin(auth *SSHConfig) (*ssh.Client, error) {
	config, err := s.getSSHConfig(auth)
	if err != nil {
		return nil, fmt.Errorf("failed to create SSH config: %v", err)
	}

	return ssh.Dial("tcp", fmt.Sprintf("%s:%d", auth.IP, auth.Port), config)
}

// UploadFile 上传文件到远程主机
func (s *SSHUtil) UploadFile(auth *SSHConfig, localPath, remotePath string) error {
	client, err := s.TerminalLogin(auth)
	if err != nil {
		return err
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	// 使用scp命令上传文件
	cmd := fmt.Sprintf("scp %s %s@%s:%s", localPath, auth.Username, auth.IP, remotePath)
	_, err = s.ExecuteRemoteCommand(auth, cmd)
	return err
}

// GetSystemInfo 获取远程主机系统信息
func (s *SSHUtil) GetSystemInfo(auth *SSHConfig) (map[string]string, error) {
	// 1. 测试基本连接
	if _, err := s.ExecuteRemoteCommand(auth, "echo 'Testing SSH connection...'"); err != nil {
		return nil, fmt.Errorf("SSH connection test failed: %v", err)
	}

	// 2. 创建并执行脚本
	output, err := s.ExecuteRemoteCommand(auth, getScriptContent())
	if err != nil {
		return nil, fmt.Errorf("failed to execute script: %v", err)
	}

	// 3. 解析结果
	var info map[string]string
	if err := json.Unmarshal([]byte(output), &info); err != nil {
		return nil, fmt.Errorf("invalid script output: %v (output: %s)", err, output)
	}

	// 4. 验证必要字段
	requiredFields := []string{"privateIp", "os", "cpu", "memory", "disk", "name"}
	var missingFields []string
	for _, field := range requiredFields {
		if _, ok := info[field]; !ok {
			missingFields = append(missingFields, field)
		}
	}
	if len(missingFields) > 0 {
		return info, fmt.Errorf("missing required fields: %v (partial info: %+v)", missingFields, info)
	}

	return info, nil
}

// getScriptContent 返回可直接执行的脚本命令
func getScriptContent() string {
	return `#!/bin/bash

# 获取系统信息
privateIp=$(hostname -I | awk '{print $1}' || echo "unknown")
publicIp=$(curl -s ipinfo.io/ip 2>/dev/null || echo "")
os=$(cat /etc/os-release 2>/dev/null | grep PRETTY_NAME | cut -d= -f2 | tr -d '"' | sed 's/ LTS//;s/ //g' || echo "unknown")

# 获取CPU核心数(去掉单位)
cpu=$(nproc 2>/dev/null || echo "unknown")

# 获取内存大小(去掉单位)
memory=$(free -m 2>/dev/null | awk '/^Mem:/{printf "%.0f\n",$2/1024}' || echo "unknown")

# 获取磁盘总容量(去掉单位)
disk=$(df -h 2>/dev/null | awk '/\/$/ {print $2}' | sed 's/G//' || echo "unknown")

# 确保获取主机名称
name=$(hostname 2>/dev/null)
if [ -z "$name" ]; then
    name=$(uname -n 2>/dev/null || echo "unknown")
fi

# 输出JSON结果
echo '{"privateIp":"'$privateIp'","publicIp":"'$publicIp'","os":"'$os'","cpu":"'$cpu'","memory":"'$memory'","disk":"'$disk'","name":"'$name'"}'
`
}

// SSHExec 执行SSH命令
func SSHExec(ip string, port int, user string, password string, command string) (string, error) {
	sshUtil := NewSSHUtil()
	auth := &SSHConfig{
		IP:       ip,
		Port:     port,
		Username: user,
		Password: password,
		Type:     1, // 密码认证
	}
	return sshUtil.ExecuteRemoteCommand(auth, command)
}

// GetSSHClientByConfig 通过 SSHConfig 获取SSH客户端连接
func GetSSHClientByConfig(auth *SSHConfig) (*ssh.Client, error) {
	sshUtil := NewSSHUtil()
	config, err := sshUtil.getSSHConfig(auth)
	if err != nil {
		return nil, fmt.Errorf("创建SSH配置失败: %v", err)
	}

	// 建立SSH连接
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", auth.IP, auth.Port), config)
	if err != nil {
		return nil, fmt.Errorf("SSH连接失败: %v", err)
	}

	return client, nil
}

// ExecuteSSHCommand 在SSH客户端上执行命令（不返回输出）
func ExecuteSSHCommand(client *ssh.Client, command string) error {
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("创建SSH会话失败: %v", err)
	}
	defer session.Close()

	return session.Run(command)
}

// ExecuteSSHCommandWithOutput 在SSH客户端上执行命令并返回输出
func ExecuteSSHCommandWithOutput(client *ssh.Client, command string) (string, error) {
	session, err := client.NewSession()
	if err != nil {
		return "", fmt.Errorf("创建SSH会话失败: %v", err)
	}
	defer session.Close()

	var stdoutBuf, stderrBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	session.Stderr = &stderrBuf

	err = session.Run(command)
	output := stdoutBuf.String()
	if err != nil {
		// 即使出错也返回输出，方便调试
		return output + "\nSTDERR:\n" + stderrBuf.String(), err
	}

	return output, nil
}

