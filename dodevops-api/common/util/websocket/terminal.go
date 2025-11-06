/*

 */

package websocket

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"time"
	"unicode/utf8"
	"github.com/dnsjia/luban/common"
	"github.com/dnsjia/luban/pkg/utils"
	"golang.org/x/crypto/ssh"
)

type utf8Writer struct {
	w io.Writer
}

func (u *utf8Writer) Write(p []byte) (n int, err error) {
	// 验证是否为有效UTF-8
	if !utf8.Valid(p) {
		p = []byte(string(p)) // 强制转换为有效UTF-8
	}
	return u.w.Write(p)
}

type Terminal struct {
	Client       *ssh.Client
	TERM         string
	session      *ssh.Session
	stdinPipe    io.WriteCloser
	config       Config
	closeHandler func() error
	closed       bool
}

type Config struct {
	UserName      string
	IpAddress     string //IP地址
	Port          string
	Password      string // 密码连接
	PrivateKey    string // 私钥连接
	KeyPassphrase string // 私钥密码
	Width         int    // pty width
	Height        int    // pty height
}

func (t *Terminal) SetCloseHandler(h func() error) {
	t.closeHandler = h
}

func (t *Terminal) SetWinSize(h int, w int) {
	common.LOG.Info(fmt.Sprintf("Setting terminal size to: %dx%d", w, h))
	if err := t.session.WindowChange(h, w); err != nil {
		common.LOG.Info(fmt.Sprintf("ssh pty change windows size failed:\t %v", err))
	} else {
		common.LOG.Info("Terminal size changed successfully")
	}
}

// IsClosed 终端是否已关闭
func (t *Terminal) IsClosed() bool {
	return t.closed
}

func (t *Terminal) Close() (err error) {
	if t.IsClosed() {
		return nil
	}
	defer func() {
		if t.closeHandler != nil {
			err = t.closeHandler()
		}
		t.closed = true
	}()

	if t.stdinPipe != nil {
		_ = t.stdinPipe.Close()
	}

	if err = t.session.Close(); err != nil {
		return
	}

	if err = t.Client.Close(); err != nil {
		return
	}

	return
}

// GetSession 获取SSH会话
func (t *Terminal) GetSession() *ssh.Session {
	return t.session
}
func getTerm() (term string) {
	if term = os.Getenv("TERM"); term == "" {
		term = "xterm-256color"  // 使用标准终端类型
	}
	return
}
func (t *Terminal) Connect(stdin io.Reader, stdout io.Writer, stderr io.Writer) error {
	var err error
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // 禁用回显
		ssh.TTY_OP_ISPEED: 14400,  // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400,  // output speed = 14.4kbaud
		ssh.ONLCR:         1,      // 启用换行转换
		ssh.OCRNL:         0,      // 禁止将回车转换为换行
		ssh.INLCR:         0,      // 禁止将换行转换为回车
		ssh.IGNCR:         0,      // 不忽略回车
		ssh.ICRNL:         1,      // 启用回车转换
		ssh.OPOST:         1,      // 启用输出处理
		ssh.ONLRET:        0,      // 允许输出时执行回车
		ssh.ONOCR:         0,      // 允许在列0时输出回车
	}

	// 1. 获取stdin管道
	t.stdinPipe, err = t.session.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdin pipe: %v", err)
	}

	// 2. 设置标准输出和错误输出，确保UTF-8编码
	t.session.Stdout = &utf8Writer{w: stdout}
	t.session.Stderr = &utf8Writer{w: stderr}

	// 3. 请求PTY (确保最小80x24终端尺寸)
	width := t.config.Width
	height := t.config.Height
	if width < 80 {
		width = 80
	}
	if height < 24 {
		height = 24
	}
	if err = t.session.RequestPty(t.TERM, height, width, modes); err != nil {
		return fmt.Errorf("request PTY failed: %v", err)
	}

	// 4. 启动shell
	if err = t.session.Shell(); err != nil {
		return fmt.Errorf("start shell failed: %v", err)
	}

	// 5. 监控会话状态
	go func() {
		err := t.session.Wait()
		if err != nil {
			common.LOG.Error(fmt.Sprintf("SSH session wait error: %v", err))
		}
		_ = t.Close()
	}()

		// 6. 添加连接健康检查
		go func() {
			ticker := time.NewTicker(30 * time.Second)
			defer ticker.Stop()
			for range ticker.C {
				if _, err := t.session.SendRequest("keepalive@golang.org", true, nil); err != nil {
					common.LOG.Error(fmt.Sprintf("SSH keepalive failed: %v", err))
					_ = t.Close()
					return
				}
			}
		}()

	return nil
}

func NewTerminal(config Config) (*Terminal, error) {
	var authMethods []ssh.AuthMethod

	// 验证配置
	if config.UserName == "" || config.IpAddress == "" || config.Port == "" {
		return nil, fmt.Errorf("invalid SSH configuration")
	}

	sshConfig := &ssh.ClientConfig{
		User:            config.UserName,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		BannerCallback:  ssh.BannerDisplayStderr(),
		Timeout:         15 * time.Second,
	}

	// 设置认证方法
	if config.PrivateKey != "" {
		pk, err := getPrivateKey(config.PrivateKey, config.KeyPassphrase)
		if err != nil {
			return nil, fmt.Errorf("failed to parse private key: %v", err)
		}
		authMethods = append(authMethods, pk)
	} else if config.Password != "" {
		authMethods = append(authMethods, ssh.Password(config.Password))
	} else {
		return nil, fmt.Errorf("no authentication method provided")
	}

	sshConfig.Auth = authMethods

	// 建立SSH连接
	addr := net.JoinHostPort(config.IpAddress, config.Port)
	client, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to dial SSH server: %v", err)
	}

	// 创建SSH会话
	session, err := client.NewSession()
	if err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to create SSH session: %v", err)
	}

	// 初始化终端
	term := &Terminal{
		TERM:    getTerm(),
		Client:  client,
		config:  config,
		session: session,
	}

	// 设置默认窗口大小
	if term.config.Width <= 0 {
		term.config.Width = 120  // 增加默认宽度
	}
	if term.config.Height <= 0 {
		term.config.Height = 40  // 增加默认高度
	}

	return term, nil
}

func getPrivateKey(privateKeyPath string, privateKeyPassphrase string) (ssh.AuthMethod, error) {
	if !utils.FileExist(privateKeyPath) {
		privateKeyPath = filepath.Join(os.Getenv("HOME"), ".ssh/id_rsa")
	}
	key, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("unable to parse private key: %v", err)
	}
	var signer ssh.Signer
	if privateKeyPassphrase != "" {
		signer, err = ssh.ParsePrivateKeyWithPassphrase(key, []byte(privateKeyPassphrase))
	} else {
		signer, err = ssh.ParsePrivateKey(key)
	}
	if err != nil {
		return nil, fmt.Errorf("parse private key failed: %v", err)
	}
	return ssh.PublicKeys(signer), nil
}
