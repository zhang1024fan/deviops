package service

import (
	"bytes"
	"errors"
	"fmt"
	"dodevops-api/api/cmdb/dao"
	configModel "dodevops-api/api/configcenter/model"
	"dodevops-api/common"
	"dodevops-api/common/util"
	"dodevops-api/common/util/websocket"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/ssh"
)

type CmdbHostSSHServiceInterface interface {
	ConnectTerminal(c *gin.Context, hostID uint) (*websocket.WebSSH, error)
	ExecuteCommand(c *gin.Context, hostID uint, command string) (*CommandResponse, error)
	UploadFile(c *gin.Context, hostID uint, filePath string, destPath string) error
}

type CmdbHostSSHServiceImpl struct{}

func (s *CmdbHostSSHServiceImpl) ConnectTerminal(c *gin.Context, hostID uint) (*websocket.WebSSH, error) {
	host, err := dao.NewCmdbHostSSHDao().GetHostSSHInfo(hostID)
	if err != nil {
		log.Printf("获取主机ID=%d信息失败: %v", hostID, err)
		return nil, fmt.Errorf("获取主机信息失败: %v", err)
	}

	if host.SSHIP == "" || host.SSHPort == 0 || host.SSHName == "" {
		log.Printf("主机ID=%d SSH信息不完整", hostID)
		return nil, errors.New("主机SSH信息不完整")
	}

	var authMethod ssh.AuthMethod
	if host.SSHKeyID != 0 {
		credential, err := dao.NewCmdbHostSSHDao().GetSSHCredentials(host.SSHKeyID)
		if err != nil {
			log.Printf("获取SSH凭据失败: %v", err)
			return nil, fmt.Errorf("获取SSH凭据失败: %v", err)
		}

		// 查询认证类型
		var auth configModel.EcsAuth
		if err := common.GetDB().Table("config_ecsauth").Where("id = ?", host.SSHKeyID).First(&auth).Error; err != nil {
			log.Printf("获取认证类型失败: %v", err)
			return nil, fmt.Errorf("获取认证类型失败: %v", err)
		}

		// 根据认证类型创建认证方法
		switch auth.Type {
		case 1: // 密码认证
			authMethod = ssh.Password(credential)
		case 2: // 密钥认证
			authMethod, err = util.NewSSHUtil().PublicKeyAuth(credential)
			if err != nil {
				log.Printf("创建密钥认证失败: %v", err)
				return nil, fmt.Errorf("创建密钥认证失败: %v", err)
			}
		case 3: // 公钥免认证 - 自动查找本地私钥
			// 先尝试用户主目录
			if userKeyAuth, err := util.NewSSHUtil().UserKeyAuth(); err == nil {
				authMethod = userKeyAuth
			} else {
				// 再尝试系统默认路径
				if defaultKeyAuth, err := util.NewSSHUtil().DefaultKeyAuth(); err == nil {
					authMethod = defaultKeyAuth
				} else {
					log.Printf("type=3 公钥免认证失败: 本地未找到私钥文件")
					return nil, errors.New("type=3 公钥免认证失败: 本地未找到私钥文件")
				}
			}
		default:
			return nil, errors.New("不支持的认证方式")
		}
	} else {
		return nil, errors.New("主机未配置SSH凭据")
	}

	webSSH, err := websocket.NewWebSSHWithAuth(host.SSHIP, strconv.Itoa(host.SSHPort), host.SSHName, authMethod)
	if err != nil {
		log.Printf("创建WebSSH连接失败: %v", err)
		return nil, fmt.Errorf("创建WebSSH连接失败: %v", err)
	}

	return webSSH, nil
}

type CommandResponse struct {
	Command string `json:"command"`
	HostID  string `json:"hostId"` 
	Output  string `json:"output"`
}

func (s *CmdbHostSSHServiceImpl) ExecuteCommand(c *gin.Context, hostID uint, command string) (*CommandResponse, error) {
	// 创建独立的SSH连接，不使用WebSSH
	host, err := dao.NewCmdbHostSSHDao().GetHostSSHInfo(hostID)
	if err != nil {
		return nil, fmt.Errorf("获取主机信息失败: %v", err)
	}

	if host.SSHIP == "" || host.SSHPort == 0 || host.SSHName == "" {
		return nil, errors.New("主机SSH信息不完整")
	}

	var authMethod ssh.AuthMethod
	if host.SSHKeyID != 0 {
		credential, err := dao.NewCmdbHostSSHDao().GetSSHCredentials(host.SSHKeyID)
		if err != nil {
			return nil, fmt.Errorf("获取SSH凭据失败: %v", err)
		}

		var auth configModel.EcsAuth
		if err := common.GetDB().Table("config_ecsauth").Where("id = ?", host.SSHKeyID).First(&auth).Error; err != nil {
			return nil, fmt.Errorf("获取认证类型失败: %v", err)
		}

		switch auth.Type {
		case 1: // 密码认证
			authMethod = ssh.Password(credential)
		case 2: // 密钥认证
			authMethod, err = util.NewSSHUtil().PublicKeyAuth(credential)
			if err != nil {
				return nil, fmt.Errorf("创建密钥认证失败: %v", err)
			}
		case 3: // 公钥免认证 - 自动查找本地私钥
			// 先尝试用户主目录
			if userKeyAuth, err := util.NewSSHUtil().UserKeyAuth(); err == nil {
				authMethod = userKeyAuth
			} else {
				// 再尝试系统默认路径
				if defaultKeyAuth, err := util.NewSSHUtil().DefaultKeyAuth(); err == nil {
					authMethod = defaultKeyAuth
				} else {
					return nil, errors.New("type=3 公钥免认证失败: 本地未找到私钥文件")
				}
			}
		default:
			return nil, errors.New("不支持的认证方式")
		}
	} else {
		return nil, errors.New("主机未配置SSH凭据")
	}

	config := &ssh.ClientConfig{
		User:            host.SSHName,
		Auth:            []ssh.AuthMethod{authMethod},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         15 * time.Second,
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", host.SSHIP, host.SSHPort), config)
	if err != nil {
		return nil, fmt.Errorf("创建SSH连接失败: %v", err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return nil, fmt.Errorf("创建SSH会话失败: %v", err)
	}
	defer session.Close()

	// 执行命令并获取输出
	output, err := session.CombinedOutput(command)
	if err != nil {
		return nil, fmt.Errorf("命令执行失败: %v", err)
	}

	// 清理输出
	ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*[mK]`)
	cleanOutput := ansiRegex.ReplaceAllString(string(output), "")

	promptRegex := regexp.MustCompile(`.*@.*:[~#]\s*`)
	parts := promptRegex.Split(cleanOutput, -1)

	var finalOutput string
	if len(parts) > 2 {
		finalOutput = strings.TrimSpace(parts[1])
	} else {
		finalOutput = strings.TrimSpace(cleanOutput)
	}

	return &CommandResponse{
		Command: command,
		HostID:  strconv.Itoa(int(hostID)),
		Output:  finalOutput,
	}, nil
}

func (s *CmdbHostSSHServiceImpl) UploadFile(c *gin.Context, hostID uint, filePath string, destPath string) error {
	// 验证文件路径是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Printf("文件不存在: %s", filePath)
		return fmt.Errorf("文件不存在: %s", filePath)
	}

	host, err := dao.NewCmdbHostSSHDao().GetHostSSHInfo(hostID)
	if err != nil {
		log.Printf("获取主机信息失败: %v", err)
		return fmt.Errorf("获取主机信息失败: %v", err)
	}

	if host.SSHIP == "" || host.SSHPort == 0 || host.SSHName == "" {
		return errors.New("主机SSH信息不完整")
	}

	var password string
	if host.SSHKeyID != 0 {
		password, err = dao.NewCmdbHostSSHDao().GetSSHCredentials(host.SSHKeyID)
		if err != nil {
			log.Printf("获取SSH凭据失败: %v", err)
			return fmt.Errorf("获取SSH凭据失败: %v", err)
		}
	} else {
		return errors.New("主机未配置SSH凭据")
	}

	// 使用绝对路径
	absFilePath, err := filepath.Abs(filePath)
	if err != nil {
		log.Printf("获取绝对路径失败: %v", err)
		return fmt.Errorf("获取文件绝对路径失败: %v", err)
	}

	// 验证文件存在
	if _, err := os.Stat(absFilePath); os.IsNotExist(err) {
		log.Printf("文件不存在(绝对路径): %s", absFilePath)
		return fmt.Errorf("文件不存在: %s", absFilePath)
	}

	// 生成唯一临时目录路径
	tmpDir := fmt.Sprintf("/tmp/ssh_uploads_%d_%d", time.Now().UnixNano(), os.Getpid())
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		log.Printf("创建临时目录失败: %v", err)
		return fmt.Errorf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(tmpDir) // 完成后自动清理

	// 实现文件复制
	tmpFilePath := filepath.Join(tmpDir, filepath.Base(filePath))
	srcFile, err := os.Open(absFilePath)
	if err != nil {
		log.Printf("打开源文件失败: %v", err)
		return fmt.Errorf("打开源文件失败: %v", err)
	}
	defer srcFile.Close()

	dstFile, err := os.Create(tmpFilePath)
	if err != nil {
		log.Printf("创建目标文件失败: %v", err)
		return fmt.Errorf("创建目标文件失败: %v", err)
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		log.Printf("复制文件内容失败: %v", err)
		return fmt.Errorf("复制文件内容失败: %v", err)
	}
	absFilePath = tmpFilePath // 使用临时文件路径

	// 构建SCP命令并添加详细调试
	log.Printf("文件上传调试信息:")
	log.Printf(" - 本地文件路径: %s", absFilePath)
	log.Printf(" - 文件存在: %t", func() bool {
		_, err := os.Stat(absFilePath)
		return err == nil
	}())
	log.Printf(" - 文件权限: %v", func() string {
		fi, err := os.Stat(absFilePath)
		if err != nil {
			return err.Error()
		}
		return fi.Mode().String()
	}())
	log.Printf(" - 当前工作目录: %s", func() string {
		wd, _ := os.Getwd()
		return wd
	}())
	log.Printf(" - 文件所在目录: %s", filepath.Dir(absFilePath))
	log.Printf(" - 文件名: %s", filepath.Base(absFilePath))

	config := &ssh.ClientConfig{
		User: host.SSHName,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         30 * time.Second, // 增加SSH连接超时时间
	}

	addr := fmt.Sprintf("%s:%d", host.SSHIP, host.SSHPort)

	// 先测试基本SSH连接
	testClient, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		log.Printf("SSH连接测试失败: %v", err)
		return fmt.Errorf("无法连接到SSH服务器(%s): %v (请检查网络连接和服务器状态)", addr, err)
	}
	defer testClient.Close()

	// 测试执行简单命令
	testSession, err := testClient.NewSession()
	if err != nil {
		log.Printf("创建测试会话失败: %v", err)
		return fmt.Errorf("SSH会话创建失败: %v", err)
	}
	defer testSession.Close()

	if _, err := testSession.CombinedOutput("echo connection-test"); err != nil {
		log.Printf("SSH命令测试失败: %v", err)
		return fmt.Errorf("SSH命令执行失败: %v (请检查远程服务器配置)", err)
	}

	// 创建正式连接
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		log.Printf("创建SSH连接失败: %v", err)
		return fmt.Errorf("创建SSH连接失败: %v", err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		log.Printf("创建SSH会话失败: %v", err)
		return fmt.Errorf("创建SSH会话失败: %v", err)
	}
	defer session.Close()

	// 获取当前工作目录
	wd, _ := os.Getwd()
	log.Printf("当前工作目录: %s", wd)

	// 检查上传锁请求头
	if lockHeader := c.GetHeader("X-Upload-Lock"); lockHeader != "" {
		log.Printf("检测到上传锁请求头: %s", lockHeader)
	}

	// 直接上传文件，依赖前端并发控制
	// 设置60秒超时(与前端一致)
	log.Printf("开始文件上传流程(超时60秒)...")

	// 使用Go的SCP实现上传文件
	file, err := os.Open(absFilePath)
	if err != nil {
		log.Printf("打开本地文件失败: %v", err)
		return fmt.Errorf("打开本地文件失败: %v", err)
	}
	defer file.Close()

	// 创建远程目录 (使用单独的session)
	log.Printf("开始创建远程目录...")
	mkdirSession, err := client.NewSession()
	if err != nil {
		log.Printf("创建mkdir会话失败: %v", err)
		return fmt.Errorf("创建mkdir会话失败: %v", err)
	}
	defer mkdirSession.Close()

	mkdirCmd := fmt.Sprintf("mkdir -p %q", filepath.Dir(destPath))
	if output, err := mkdirSession.CombinedOutput(mkdirCmd); err != nil {
		log.Printf("创建远程目录失败: %v", err)
		log.Printf("命令输出: %s", string(output))
		return fmt.Errorf("创建远程目录失败: %v", err)
	}
	log.Printf("远程目录创建成功")

	// 尝试使用SCP协议上传文件
	scpSession, err := client.NewSession()
	if err != nil {
		log.Printf("创建SCP会话失败: %v", err)
		return fmt.Errorf("创建SCP会话失败: %v", err)
	}
	defer scpSession.Close()

	// 准备SCP命令
	go func() {
		w, _ := scpSession.StdinPipe()
		defer w.Close()

		file, err := os.Open(absFilePath)
		if err != nil {
			log.Printf("打开本地文件失败: %v", err)
			return
		}
		defer file.Close()

		fileInfo, _ := file.Stat()
		fmt.Fprintf(w, "C%04o %d %s\n", fileInfo.Mode().Perm(), fileInfo.Size(), filepath.Base(destPath))
		io.Copy(w, file)
		fmt.Fprint(w, "\x00")
	}()

	// 确保目标路径包含文件名
	finalDestPath := filepath.Join(destPath, filepath.Base(absFilePath))

	// 执行SCP命令
	if err := scpSession.Run(fmt.Sprintf("scp -t %s", finalDestPath)); err != nil {
		log.Printf("SCP上传失败: %v", err)

		// 回退到rsync协议
		log.Printf("尝试使用rsync协议上传...")
		rsyncCmd := fmt.Sprintf("timeout 60 rsync -az --progress --timeout=30 %s %s@%s:%s",
			absFilePath,
			host.SSHName,
			host.SSHIP,
			destPath)
		log.Printf("执行rsync命令: %s", rsyncCmd)

		rsyncSession, err := client.NewSession()
		if err != nil {
			log.Printf("创建rsync会话失败: %v", err)
			return fmt.Errorf("创建rsync会话失败: %v", err)
		}
		defer rsyncSession.Close()

		var stderr bytes.Buffer
		rsyncSession.Stderr = &stderr

		if err := rsyncSession.Run(rsyncCmd); err != nil {
			log.Printf("文件上传失败(SCP和rsync): %v\n详细错误: %s", err, stderr.String())
			return fmt.Errorf("文件上传失败(SCP和rsync): %v\n详细错误: %s", err, stderr.String())
		}
	}

	log.Printf("文件上传成功: %s -> %s@%s:%s", absFilePath, host.SSHName, host.SSHIP, destPath)

	// 验证远程文件
	checkCmd := fmt.Sprintf("ls -la %q", destPath)
	if output, err := session.CombinedOutput(checkCmd); err != nil {
		log.Printf("远程文件验证失败: %v", err)
	} else {
		log.Printf("远程文件验证成功: %s", string(output))
	}

	return nil
}

func GetCmdbHostSSHService() CmdbHostSSHServiceInterface {
	return &CmdbHostSSHServiceImpl{}
}
