package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"golang.org/x/crypto/ssh"
)

// WebSSH 简化的WebSSH实现
type WebSSH struct {
	conn      *websocket.Conn
	client    *ssh.Client
	session   *ssh.Session
	stdinPipe io.WriteCloser
	cancel    context.CancelFunc
}

// GetStdinPipe 获取输入管道
func (w *WebSSH) GetStdinPipe() io.WriteCloser {
	if w.stdinPipe == nil && w.session != nil {
		var err error
		w.stdinPipe, err = w.session.StdinPipe()
		if err != nil {
			log.Printf("初始化SSH输入管道失败: %v", err)
			return nil
		}
	}
	return w.stdinPipe
}

// GetSession 获取SSH会话
func (w *WebSSH) GetSession() *ssh.Session {
	return w.session
}

// NewWebSSH 创建新的WebSSH连接(密码认证)
func NewWebSSH(host, port, user, password string) (*WebSSH, error) {
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
			ssh.KeyboardInteractive(func(user, instruction string, questions []string, echos []bool) (answers []string, err error) {
				if len(questions) == 0 {
					return []string{}, nil
				}
				return []string{password}, nil
			}),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         15 * time.Second,
		Config: ssh.Config{
			Ciphers: []string{
				"aes128-ctr",
				"aes192-ctr", 
				"aes256-ctr",
				"aes128-gcm@openssh.com",
				"arcfour256",
				"arcfour128",
				"arcfour",
				"aes128-cbc",
				"3des-cbc",
			},
		},
	}

	return newWebSSHWithConfig(host, port, config)
}

// NewWebSSHWithAuth 创建新的WebSSH连接(自定义认证方法)
func NewWebSSHWithAuth(host, port, user string, authMethod ssh.AuthMethod) (*WebSSH, error) {
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{authMethod},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         15 * time.Second,
		Config: ssh.Config{
			Ciphers: []string{
				"aes128-ctr",
				"aes192-ctr", 
				"aes256-ctr",
				"aes128-gcm@openssh.com",
				"arcfour256",
				"arcfour128",
				"arcfour",
				"aes128-cbc",
				"3des-cbc",
			},
		},
	}

	webSSH, err := newWebSSHWithConfig(host, port, config)
	if err != nil {
		return nil, err
	}

	// 立即初始化session和stdinPipe
	session, err := webSSH.client.NewSession()
	if err != nil {
		return nil, fmt.Errorf("创建SSH会话失败: %v", err)
	}
	webSSH.session = session

	stdinPipe, err := session.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("创建输入管道失败: %v", err)
	}
	webSSH.stdinPipe = stdinPipe

	return webSSH, nil
}

// newWebSSHWithConfig 使用配置创建WebSSH连接
func newWebSSHWithConfig(host, port string, config *ssh.ClientConfig) (*WebSSH, error) {
	addr := net.JoinHostPort(host, port)
	log.Printf("尝试连接SSH服务器: %s", addr)
	log.Printf("使用的认证方法: %v", config.Auth)
	
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		log.Printf("SSH连接失败详情: %v", err)
		return nil, fmt.Errorf("failed to dial SSH server: %v", err)
	}
	
	log.Println("SSH连接成功建立")

	return &WebSSH{
		client: client,
	}, nil
}

// Connect 建立SSH会话并处理WebSocket连接
func (w *WebSSH) Connect(wsConn *websocket.Conn) error {
	session, err := w.client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %v", err)
	}
	w.session = session
	w.conn = wsConn

	// 获取输入管道
	w.stdinPipe, err = session.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdin pipe: %v", err)
	}

	// 设置输出
	stdout, err := session.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdout pipe: %v", err)
	}

	// 启动shell
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
		ssh.ONLCR:         1,      // 启用换行转换
		ssh.OCRNL:         0,      // 禁止将回车转换为换行
		ssh.INLCR:         0,      // 禁止将换行转换为回车
		ssh.IGNCR:         0,      // 不忽略回车
		ssh.ICRNL:         1,      // 启用回车转换
		ssh.OPOST:         1,      // 启用输出处理
		ssh.ONLRET:        0,      // 允许输出时执行回车
		ssh.ONOCR:         0,      // 允许在列0时输出回车
	}
	width := 160 // 默认宽度
	height := 40
	if err := session.RequestPty("xterm-256color", height, width, modes); err != nil {
		return fmt.Errorf("request pty failed: %v", err)
	}
	if err := session.Shell(); err != nil {
		return fmt.Errorf("start shell failed: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	w.cancel = cancel

	// 处理WebSocket输入
	go w.handleInput(ctx)

	// 处理SSH输出
	go w.handleOutput(ctx, stdout)

	return nil
}

func (w *WebSSH) handleInput(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			_, data, err := w.conn.ReadMessage()
			if err != nil {
				log.Printf("WebSocket read error: %v", err)
				return
			}

			// 处理窗口调整
			var msg map[string]interface{}
			if err := json.Unmarshal(data, &msg); err == nil {
				if msg["type"] == "resize" {
					cols, _ := strconv.Atoi(fmt.Sprintf("%v", msg["cols"]))
					rows, _ := strconv.Atoi(fmt.Sprintf("%v", msg["rows"]))
					if cols > 0 && rows > 0 {
						w.session.WindowChange(rows, cols)
					}
					continue
				}
			}

			// 发送输入到SSH
			if _, err := w.stdinPipe.Write(data); err != nil {
				log.Printf("SSH write error: %v", err)
				return
			}
		}
	}
}

func (w *WebSSH) handleOutput(ctx context.Context, stdout io.Reader) {
	buf := make([]byte, 1024)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			n, err := stdout.Read(buf)
			if err != nil {
				if err != io.EOF {
					log.Printf("SSH read error: %v", err)
				}
				return
			}
			if n > 0 {
				if err := w.conn.WriteMessage(websocket.TextMessage, buf[:n]); err != nil {
					log.Printf("WebSocket write error: %v", err)
					return
				}
			}
		}
	}
}

// Close 关闭连接
func (w *WebSSH) Close() error {
	if w.cancel != nil {
		w.cancel()
	}
	if w.session != nil {
		w.session.Close()
	}
	if w.client != nil {
		return w.client.Close()
	}
	return nil
}
