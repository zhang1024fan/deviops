package service

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"dodevops-api/api/cmdb/dao"
	"dodevops-api/api/cmdb/model"
	agentDao "dodevops-api/api/monitor/dao"
	agentModel "dodevops-api/api/monitor/model"
	"dodevops-api/common"
	"dodevops-api/common/agent"
	"dodevops-api/common/config"
	"dodevops-api/common/util"
	"dodevops-api/pkg/redis"

	"dodevops-api/common/result"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/ssh"
	"gorm.io/gorm"
)

// EcsAuth SSH认证信息
type EcsAuth struct {
	ID         uint       `gorm:"column:id;comment:'主键';primaryKey;NOT NULL" json:"id"`
	Name       string     `gorm:"column:name;varchar(64);comment:'凭证名称';NOT NULL" json:"name"`
	Type       int        `gorm:"column:type;comment:'认证类型:1->密码,2->密钥';NOT NULL" json:"type"`
	Username   string     `gorm:"column:username;varchar(64);comment:'用户名(type=1时使用)'" json:"username"`
	Password   string     `gorm:"column:password;varchar(256);comment:'密码(type=1时使用)'" json:"password"`
	PublicKey  string     `gorm:"column:public_key;type:text;comment:'公钥(type=2时使用)'" json:"publicKey"`
	Port       int        `gorm:"column:port;comment:'端口号';default:22" json:"port"`
	CreateTime util.HTime `gorm:"column:create_time;comment:'创建时间';NOT NULL" json:"createTime"`
	Remark     string     `gorm:"column:remark;varchar(500);comment:'备注'" json:"remark"`
}

// TableName 指定表名
func (EcsAuth) TableName() string {
	return "config_ecsauth"
}

// EcsAuthDao SSH密钥认证DAO接口
type EcsAuthDao interface {
	GetEcsAuthById(id uint) (EcsAuth, error)
}

// EcsAuthDaoImpl SSH密钥认证DAO实现
type EcsAuthDaoImpl struct{}

// NewEcsAuthDao 创建EcsAuthDao实例
func NewEcsAuthDao() EcsAuthDao {
	return &EcsAuthDaoImpl{}
}

// GetEcsAuthById 根据ID获取SSH密钥
func (d *EcsAuthDaoImpl) GetEcsAuthById(id uint) (EcsAuth, error) {
	// 从数据库查询SSH密钥
	var ecsAuth EcsAuth
	err := common.GetDB().Where("id = ?", id).First(&ecsAuth).Error
	if err != nil {
		return EcsAuth{}, fmt.Errorf("查询SSH密钥失败: %v", err)
	}
	return ecsAuth, nil
}

// AgentServiceInterface agent服务接口
type AgentServiceInterface interface {
	DeployAgent(c *gin.Context, hostID uint)                                              // 部署agent到指定主机
	BatchDeployAgent(c *gin.Context, dto *agentModel.BatchDeployAgentDto)                 // 批量部署agent
	UninstallAgent(c *gin.Context, hostID uint)                                           // 卸载指定主机的agent
	BatchUninstallAgent(c *gin.Context, dto *agentModel.BatchDeployAgentDto)              // 批量卸载agent
	DeleteAgent(c *gin.Context, hostID uint)                                              // 删除agent数据(用于离线服务器，通过hostID)
	DeleteAgentByID(c *gin.Context, agentID uint)                                         // 删除agent数据(通过agentID)
	GetAgentStatus(c *gin.Context, hostID uint)                                           // 获取agent状态
	RestartAgent(c *gin.Context, hostID uint)                                             // 重启agent
	GetAgentList(c *gin.Context, dto *agentModel.AgentListDto)                            // 获取agent列表
	UpdateHeartbeat(c *gin.Context, hostID uint, heartbeat *agentModel.AgentHeartbeatDto) // 更新心跳(旧版本，保持兼容)
	UpdateHeartbeatByIP(c *gin.Context, heartbeat *agentModel.AgentHeartbeatDto)          // 通过IP更新心跳
	GetAgentStatistics(c *gin.Context)                                                    // 获取统计信息
	CheckOfflineAgents()                                                                  // 检查离线agent
}

// AgentServiceImpl agent服务实现
type AgentServiceImpl struct {
	hostDao    dao.CmdbHostDao
	agentDao   agentDao.AgentDao
	ecsAuthDao EcsAuthDao // 添加ECS认证DAO
}

// NewAgentService 创建agent服务实例
func NewAgentService() AgentServiceInterface {
	return &AgentServiceImpl{
		hostDao:    dao.NewCmdbHostDao(),
		agentDao:   agentDao.NewAgentDao(),
		ecsAuthDao: NewEcsAuthDao(), // 初始化ECS认证DAO
	}
}

// 编译目标平台类型
type BuildTarget struct {
	GOOS   string
	GOARCH string
	Suffix string
}

// 支持的编译目标平台
var buildTargets = map[string]BuildTarget{
	"linux":   {GOOS: "linux", GOARCH: "amd64", Suffix: ""},
	"windows": {GOOS: "windows", GOARCH: "amd64", Suffix: ".exe"},
	"darwin":  {GOOS: "darwin", GOARCH: "amd64", Suffix: ""},
}

// DeployAgent 部署agent到指定主机
func (s *AgentServiceImpl) DeployAgent(c *gin.Context, hostID uint) {
	// 参数校验
	if hostID == 0 {
		result.FailedWithCode(c, int(result.ApiCode.ValidationParameterError), "主机ID不能为空")
		return
	}

	// 获取主机信息
	host, err := s.hostDao.GetCmdbHostById(hostID)
	if err != nil {
		result.FailedWithCode(c, int(result.ApiCode.FAILED), "获取主机信息失败")
		return
	}

	if host.Name == "" {
		result.FailedWithCode(c, int(result.ApiCode.FAILED), "主机名称为空")
		return
	}

	// 创建或更新Agent记录（Linux主机专用）
	agentRecord := &agentModel.Agent{
		HostID:      hostID,
		HostName:    host.Name,
		Version:     "1.0.0",
		Status:      agentModel.AgentStatusDeploying,
		InstallPath: "/opt/agent", // Linux的标准安装路径
		Port:        9100,
		CreateTime:  util.HTime{Time: time.Now()},
	}

	// 先创建或更新数据库记录
	if err := s.createOrUpdateAgent(agentRecord); err != nil {
		result.FailedWithCode(c, int(result.ApiCode.FAILED), "创建Agent记录失败")
		return
	}

	// 异步部署
	go func() {
		deployErr := s.deployAgentAsync(host.ID, host)
		if deployErr != nil {
			log.Printf("Deploy agent to host %d failed: %v", hostID, deployErr)
			// 更新状态为部署失败
			s.agentDao.UpdateByHostID(hostID, map[string]interface{}{
				"status":    agentModel.AgentStatusDeployFailed,
				"error_msg": deployErr.Error(),
			})
		}
	}()

	result.Success(c, map[string]interface{}{
		"message": "Agent部署任务已启动，正在后台执行",
		"hostId":  hostID,
		"status":  "deploying",
	})
}

// BatchDeployAgent 批量部署agent到多台主机
// 优化策略：编译一次，并行拷贝和启动到多台主机
func (s *AgentServiceImpl) BatchDeployAgent(c *gin.Context, dto *agentModel.BatchDeployAgentDto) {
	// 参数校验
	if len(dto.HostIDs) == 0 {
		result.FailedWithCode(c, int(result.ApiCode.ValidationParameterError), "主机ID列表不能为空")
		return
	}

	// 设置默认版本
	version := dto.Version
	if version == "" {
		version = "1.0.0"
	}

	var successHosts []uint
	var failedHosts []map[string]interface{}

	// 收集所有主机信息并创建Agent记录
	var validHosts []model.CmdbHost
	for _, hostID := range dto.HostIDs {
		// 获取主机信息
		host, err := s.hostDao.GetCmdbHostById(hostID)
		if err != nil {
			failedHosts = append(failedHosts, map[string]interface{}{
				"hostId": hostID,
				"error":  "获取主机信息失败",
			})
			continue
		}

		if host.Name == "" {
			failedHosts = append(failedHosts, map[string]interface{}{
				"hostId": hostID,
				"error":  "主机名称为空",
			})
			continue
		}

		// 创建或更新Agent记录
		agentRecord := &agentModel.Agent{
			HostID:      hostID,
			HostName:    host.Name,
			Version:     version,
			Status:      agentModel.AgentStatusDeploying,
			InstallPath: "/opt/agent",
			Port:        9100,
			CreateTime:  util.HTime{Time: time.Now()},
		}

		// 创建或更新数据库记录
		if err := s.createOrUpdateAgent(agentRecord); err != nil {
			failedHosts = append(failedHosts, map[string]interface{}{
				"hostId": hostID,
				"error":  "创建Agent记录失败",
			})
			continue
		}

		successHosts = append(successHosts, hostID)
		validHosts = append(validHosts, host)
	}

	// 异步批量部署
	go func() {
		batchErr := s.batchDeployAgentAsync(validHosts)
		if batchErr != nil {
			log.Printf("Batch deploy agents failed: %v", batchErr)
		}
	}()

	result.Success(c, map[string]interface{}{
		"message":      fmt.Sprintf("批量部署任务已启动，成功启动 %d 台，失败 %d 台", len(successHosts), len(failedHosts)),
		"successCount": len(successHosts),
		"failedCount":  len(failedHosts),
		"successHosts": successHosts,
		"failedHosts":  failedHosts,
		"status":       "deploying",
	})
}

// GetAgentList 获取agent列表
func (s *AgentServiceImpl) GetAgentList(c *gin.Context, dto *agentModel.AgentListDto) {
	agents, total, err := s.agentDao.GetList(dto)
	if err != nil {
		result.FailedWithCode(c, int(result.ApiCode.FAILED), "获取Agent列表失败")
		return
	}

	// 转换为VO并关联SSH IP
	var agentVOs []*agentModel.AgentVO
	for _, agent := range agents {
		agentVO := agent.ToVO()
		// 获取主机信息以获得SSH IP
		if host, err := s.hostDao.GetCmdbHostById(agent.HostID); err == nil {
			agentVO.SSHIP = host.SSHIP
		}
		agentVOs = append(agentVOs, agentVO)
	}

	result.Success(c, map[string]interface{}{
		"list":     agentVOs,
		"total":    total,
		"page":     dto.Page,
		"pageSize": dto.PageSize,
	})
}

// GetAgentStatus 获取agent状态
func (s *AgentServiceImpl) GetAgentStatus(c *gin.Context, hostID uint) {
	if hostID == 0 {
		result.FailedWithCode(c, int(result.ApiCode.ValidationParameterError), "主机ID不能为空")
		return
	}

	// 从数据库获取Agent信息
	agent, err := s.agentDao.GetByHostID(hostID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result.FailedWithCode(c, int(result.ApiCode.FAILED), "Agent未部署")
			return
		}
		result.FailedWithCode(c, int(result.ApiCode.FAILED), "获取Agent状态失败")
		return
	}

	// 如果Agent状态为运行中，启动异步健康检查
	if agent.Status == agentModel.AgentStatusRunning {
		go s.performHealthCheck(hostID)
	}

	// 转换为VO并关联SSH IP
	agentVO := agent.ToVO()
	if host, err := s.hostDao.GetCmdbHostById(agent.HostID); err == nil {
		agentVO.SSHIP = host.SSHIP
	}

	result.Success(c, agentVO)
}

// UpdateHeartbeat 更新心跳(旧版本，保持兼容性)
func (s *AgentServiceImpl) UpdateHeartbeat(c *gin.Context, hostID uint, heartbeat *agentModel.AgentHeartbeatDto) {
	if hostID == 0 {
		result.FailedWithCode(c, int(result.ApiCode.ValidationParameterError), "主机ID不能为空")
		return
	}

	if err := s.agentDao.UpdateHeartbeat(hostID, heartbeat); err != nil {
		result.FailedWithCode(c, int(result.ApiCode.FAILED), "更新心跳失败")
		return
	}

	result.Success(c, map[string]interface{}{
		"message": "心跳更新成功",
		"hostId":  hostID,
	})
}

// UpdateHeartbeatByIP 通过IP更新心跳
func (s *AgentServiceImpl) UpdateHeartbeatByIP(c *gin.Context, heartbeat *agentModel.AgentHeartbeatDto) {
	var host model.CmdbHost
	var err error

	// 优先使用hostname匹配（因为agent上报的IP可能是内网IP，与cmdb中的ssh_ip不匹配）
	if heartbeat.Hostname != "" {
		host, err = s.hostDao.GetCmdbHostByName(heartbeat.Hostname)
		if err != nil {
			log.Printf("通过hostname查找主机失败: %s, 尝试使用IP匹配", heartbeat.Hostname)
		}
	}

	// 如果hostname匹配失败或为空，尝试使用IP匹配
	if host.ID == 0 && heartbeat.IP != "" {
		host, err = s.hostDao.GetCmdbHostByIP(heartbeat.IP)
		if err != nil {
			result.FailedWithCode(c, int(result.ApiCode.FAILED), fmt.Sprintf("未找到对应主机: hostname=%s, ip=%s", heartbeat.Hostname, heartbeat.IP))
			return
		}
	}

	if host.ID == 0 {
		result.FailedWithCode(c, int(result.ApiCode.FAILED), "无法匹配到主机记录")
		return
	}

	// 更新心跳信息
	if err := s.agentDao.UpdateHeartbeatByIP(host.ID, heartbeat); err != nil {
		result.FailedWithCode(c, int(result.ApiCode.FAILED), "更新心跳失败")
		return
	}

	result.Success(c, map[string]interface{}{
		"message":  "心跳更新成功",
		"hostId":   host.ID,
		"hostName": host.Name,
		"hostIP":   heartbeat.IP,
	})
}

// GetAgentStatistics 获取统计信息
func (s *AgentServiceImpl) GetAgentStatistics(c *gin.Context) {
	stats, err := s.agentDao.GetAgentStatistics()
	if err != nil {
		result.FailedWithCode(c, int(result.ApiCode.FAILED), "获取统计信息失败")
		return
	}

	result.Success(c, stats)
}

// CheckOfflineAgents 检查离线agent
func (s *AgentServiceImpl) CheckOfflineAgents() {
	// 检查5分钟内没有心跳的Agent
	offlineAgents, err := s.agentDao.GetOfflineAgents(5 * time.Minute)
	if err != nil {
		log.Printf("检查离线Agent失败: %v", err)
		return
	}

	for _, agent := range offlineAgents {
		log.Printf("Agent %d (主机: %s) 离线", agent.ID, agent.HostName)
		// 更新状态为启动异常
		s.agentDao.UpdateStatus(agent.ID, agentModel.AgentStatusStartError, "心跳超时")
	}
}

// cleanupPushgatewayMetrics 清理Pushgateway中的监控数据
func (s *AgentServiceImpl) cleanupPushgatewayMetrics(jobName string) {
	if jobName == "" {
		log.Printf("job名称为空，跳过Pushgateway清理")
		return
	}

	// 从配置中获取Pushgateway URL
	pushgatewayURL := "http://8.130.14.34:9091"
	if config.Config != nil && config.Config.Monitor.Pushgateway.URL != "" {
		pushgatewayURL = config.Config.Monitor.Pushgateway.URL
	}

	// 构建DELETE请求的URL
	// Pushgateway API: DELETE /metrics/job/<job_name>
	deleteURL := fmt.Sprintf("%s/metrics/job/%s", pushgatewayURL, jobName)

	log.Printf("开始清理Pushgateway监控数据: %s", deleteURL)

	// 创建HTTP DELETE请求
	req, err := http.NewRequest("DELETE", deleteURL, nil)
	if err != nil {
		log.Printf("创建DELETE请求失败: %v", err)
		return
	}

	// 发送请求
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("清理Pushgateway数据失败: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 202 || resp.StatusCode == 200 {
		log.Printf("成功清理Pushgateway数据: job=%s", jobName)
	} else {
		log.Printf("清理Pushgateway数据返回状态码: %d, job=%s", resp.StatusCode, jobName)
	}
}

// DeleteAgent 删除agent数据(通过hostID)
func (s *AgentServiceImpl) DeleteAgent(c *gin.Context, hostID uint) {
	// 参数校验
	if hostID == 0 {
		result.FailedWithCode(c, int(result.ApiCode.ValidationParameterError), "主机ID不能为空")
		return
	}

	// 检查Agent是否存在
	agent, err := s.agentDao.GetByHostID(hostID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result.FailedWithCode(c, int(result.ApiCode.FAILED), "Agent记录不存在")
			return
		}
		result.FailedWithCode(c, int(result.ApiCode.FAILED), "获取Agent信息失败")
		return
	}

	// 获取主机信息以便清理Pushgateway
	host, err := s.hostDao.GetCmdbHostById(hostID)
	if err != nil {
		log.Printf("获取主机信息失败: %v, 跳过Pushgateway清理", err)
	} else {
		// 清理Pushgateway中的监控数据
		s.cleanupPushgatewayMetrics(host.Name)
	}

	// 直接删除数据库中的Agent记录
	if err := s.agentDao.DeleteByHostID(hostID); err != nil {
		result.FailedWithCode(c, int(result.ApiCode.FAILED), "删除Agent记录失败")
		return
	}

	result.Success(c, map[string]interface{}{
		"message":  "Agent数据删除成功",
		"hostId":   hostID,
		"hostName": agent.HostName,
	})
}

// DeleteAgentByID 删除agent数据(通过agentID)
func (s *AgentServiceImpl) DeleteAgentByID(c *gin.Context, agentID uint) {
	// 参数校验
	if agentID == 0 {
		result.FailedWithCode(c, int(result.ApiCode.ValidationParameterError), "Agent ID不能为空")
		return
	}

	// 检查Agent是否存在
	agent, err := s.agentDao.GetByID(agentID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result.FailedWithCode(c, int(result.ApiCode.FAILED), "Agent记录不存在")
			return
		}
		result.FailedWithCode(c, int(result.ApiCode.FAILED), "获取Agent信息失败")
		return
	}

	// 直接删除数据库中的Agent记录
	if err := s.agentDao.Delete(agentID); err != nil {
		result.FailedWithCode(c, int(result.ApiCode.FAILED), "删除Agent记录失败")
		return
	}

	result.Success(c, map[string]interface{}{
		"message": "Agent数据删除成功",
		"agentId": agentID,
		"hostId":  agent.HostID,
		"hostName": agent.HostName,
	})
}

// getSSHKeyByID 根据ID获取SSH密钥
func (s *AgentServiceImpl) getSSHKeyByID(id uint) (*EcsAuth, error) {
	ecsAuth, err := s.ecsAuthDao.GetEcsAuthById(id)
	if err != nil {
		return nil, fmt.Errorf("获取SSH密钥失败: %v", err)
	}
	return &ecsAuth, nil
}

// batchDeployAgentAsync 批量部署agent - 优化版本：编译一次，并行拷贝启动
func (s *AgentServiceImpl) batchDeployAgentAsync(hosts []model.CmdbHost) error {
	if len(hosts) == 0 {
		return fmt.Errorf("主机列表为空")
	}

	log.Printf("开始批量部署Agent到 %d 台主机", len(hosts))

	// 1. 统一编译一次（使用第一台主机的OS信息作为目标平台）
	firstHost := hosts[0]
	targetOS := "linux"
	if strings.Contains(strings.ToLower(firstHost.OS), "windows") {
		targetOS = "windows"
	} else if strings.Contains(strings.ToLower(firstHost.OS), "darwin") {
		targetOS = "darwin"
	}

	buildTarget, exists := buildTargets[targetOS]
	if !exists {
		return fmt.Errorf("不支持的操作系统: %s", firstHost.OS)
	}

	// 更新所有主机进度：开始编译
	for _, host := range hosts {
		s.agentDao.UpdateInstallProgress(host.ID, agentModel.InstallProgressCompiling)
	}

	log.Printf("开始为批量部署编译Agent，目标OS: %s", targetOS)

	// 创建临时目录
	tempDir := fmt.Sprintf("/tmp/agent-batch-deploy-%d", time.Now().Unix())
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return fmt.Errorf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// 2. 生成并编译agent（只编译一次）
	heartbeatURL := "https://www.deviops.cn/api/v1/monitor/agent/heartbeat"
	heartbeatToken := "agent-heartbeat-token-2024"
	pushgatewayURL := "http://8.130.14.34:9091"

	if config.Config != nil && config.Config.Monitor.Agent.HeartbeatServerURL != "" {
		heartbeatURL = config.Config.Monitor.Agent.HeartbeatServerURL
	}
	if config.Config != nil && config.Config.Monitor.Agent.HeartbeatToken != "" {
		heartbeatToken = config.Config.Monitor.Agent.HeartbeatToken
	}
	if config.Config != nil && config.Config.Monitor.Pushgateway.URL != "" {
		pushgatewayURL = config.Config.Monitor.Pushgateway.URL
	}

	// 生成通用的main.go（不包含任何主机特定信息）
	mainGoContent := agent.CreateAgentMainFile(heartbeatURL, heartbeatToken, pushgatewayURL)
	if len(mainGoContent) == 0 {
		return fmt.Errorf("生成的main.go内容为空")
	}

	mainGoPath := filepath.Join(tempDir, "main.go")
	log.Printf("准备写入main.go到: %s, 内容长度: %d", mainGoPath, len(mainGoContent))
	if err := os.WriteFile(mainGoPath, []byte(mainGoContent), 0644); err != nil {
		return fmt.Errorf("创建main.go文件失败: %v", err)
	}
	log.Printf("main.go文件创建成功")

	// 创建go.mod
	goModContent := `module agent

go 1.24

require (
	github.com/prometheus/client_golang v1.17.0
	github.com/shirou/gopsutil/v3 v3.23.10
)`
	goModPath := filepath.Join(tempDir, "go.mod")
	log.Printf("准备写入go.mod到: %s", goModPath)
	if err := os.WriteFile(goModPath, []byte(goModContent), 0644); err != nil {
		return fmt.Errorf("创建go.mod文件失败: %v", err)
	}
	log.Printf("go.mod文件创建成功")

	// 设置Go环境
	goTempDir := filepath.Join(tempDir, ".gocache")
	goCacheDir := filepath.Join(goTempDir, "build")
	goModCacheDir := filepath.Join(goTempDir, "mod")

	os.MkdirAll(goCacheDir, 0755)
	os.MkdirAll(goModCacheDir, 0755)

	baseEnvVars := []string{
		"GO111MODULE=on",
		"CGO_ENABLED=0",
		"GOFLAGS=-buildvcs=false",
		fmt.Sprintf("GOMODCACHE=%s", goModCacheDir),
		fmt.Sprintf("GOCACHE=%s", goCacheDir),
		"GOPROXY=https://goproxy.cn,direct",
		"GOSUMDB=off",
	}

	// 下载依赖
	log.Printf("开始下载依赖...")
	downloadCmd := exec.Command("go", "mod", "download")
	downloadCmd.Dir = tempDir
	downloadCmd.Env = append(os.Environ(), baseEnvVars...)
	if output, err := downloadCmd.CombinedOutput(); err != nil {
		log.Printf("go mod download 警告: %v, 输出: %s", err, string(output))
	} else {
		log.Printf("依赖下载完成")
	}

	// go mod tidy
	log.Printf("开始 go mod tidy...")
	tidyCmd := exec.Command("go", "mod", "tidy")
	tidyCmd.Dir = tempDir
	tidyCmd.Env = append(os.Environ(), baseEnvVars...)
	if output, err := tidyCmd.CombinedOutput(); err != nil {
		log.Printf("go mod tidy 警告: %v, 输出: %s", err, string(output))
	} else {
		log.Printf("go mod tidy 完成")
	}

	// 编译二进制文件
	binaryName := fmt.Sprintf("dodevops-agent%s", buildTarget.Suffix)
	binaryPath := filepath.Join(tempDir, binaryName)

	buildFlags := []string{
		"build",
		"-a",
		"-ldflags", "-s -w -extldflags '-static'",
		"-installsuffix", "netgo",
		"-tags", "netgo,osusergo",
		"-o", binaryPath,
		mainGoPath,
	}

	buildCmd := exec.Command("go", buildFlags...)
	buildCmd.Dir = tempDir
	buildEnvVars := append(baseEnvVars,
		fmt.Sprintf("GOOS=%s", buildTarget.GOOS),
		fmt.Sprintf("GOARCH=%s", buildTarget.GOARCH),
	)
	buildCmd.Env = append(os.Environ(), buildEnvVars...)

	log.Printf("批量部署编译命令: %s", strings.Join(append([]string{"go"}, buildFlags...), " "))
	log.Printf("开始编译，目标平台: GOOS=%s, GOARCH=%s", buildTarget.GOOS, buildTarget.GOARCH)

	// 使用context设置5分钟超时
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	buildCmdWithTimeout := exec.CommandContext(ctx, "go", buildFlags...)
	buildCmdWithTimeout.Dir = tempDir
	buildCmdWithTimeout.Env = append(os.Environ(), buildEnvVars...)

	output, err := buildCmdWithTimeout.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		return fmt.Errorf("编译超时（5分钟），可能是网络问题或依赖下载失败")
	}
	if err != nil {
		return fmt.Errorf("编译失败: %v, 输出: %s", err, string(output))
	}

	log.Printf("批量部署编译成功！二进制文件: %s", binaryPath)

	// 更新所有主机进度：编译完成
	for _, host := range hosts {
		s.agentDao.UpdateInstallProgress(host.ID, agentModel.InstallProgressCompiled)
	}

	// 3. 并行拷贝和启动到所有主机
	type deployResult struct {
		hostID uint
		err    error
	}

	resultChan := make(chan deployResult, len(hosts))

	for _, host := range hosts {
		go func(h model.CmdbHost) {
			err := s.deployToSingleHost(h, binaryPath, binaryName, targetOS)
			resultChan <- deployResult{hostID: h.ID, err: err}
		}(host)
	}

	// 收集结果
	successCount := 0
	failedCount := 0
	for i := 0; i < len(hosts); i++ {
		result := <-resultChan
		if result.err != nil {
			failedCount++
			log.Printf("部署到主机 %d 失败: %v", result.hostID, result.err)
			s.agentDao.UpdateByHostID(result.hostID, map[string]interface{}{
				"status":    agentModel.AgentStatusDeployFailed,
				"error_msg": result.err.Error(),
			})
		} else {
			successCount++
			log.Printf("部署到主机 %d 成功", result.hostID)
		}
	}

	log.Printf("批量部署完成: 成功 %d 台, 失败 %d 台", successCount, failedCount)
	return nil
}

// deployToSingleHost 部署已编译的agent到单个主机（拷贝+启动）
func (s *AgentServiceImpl) deployToSingleHost(host model.CmdbHost, binaryPath, binaryName, targetOS string) error {
	log.Printf("开始部署到主机 %d (%s)", host.ID, host.HostName)

	// 更新进度：传输中
	s.agentDao.UpdateInstallProgress(host.ID, agentModel.InstallProgressTransfer)

	// 获取SSH配置
	var sshConfig util.SSHConfig
	if host.SSHKeyID > 0 {
		sshKey, err := s.getSSHKeyByID(host.SSHKeyID)
		if err != nil {
			return fmt.Errorf("获取SSH认证信息失败: %v", err)
		}

		sshConfig = util.SSHConfig{
			IP:       host.SSHIP,
			Port:     host.SSHPort,
			Username: host.SSHName,
			Type:     sshKey.Type,
			Timeout:  30 * time.Second,
		}

		switch sshKey.Type {
		case 1:
			sshConfig.Password = sshKey.Password
		case 2:
			sshConfig.PublicKey = sshKey.PublicKey
		case 3:
			// 公钥免认证
		default:
			return fmt.Errorf("不支持的SSH认证类型: %d", sshKey.Type)
		}
	} else {
		return fmt.Errorf("主机未配置SSH认证信息")
	}

	// 创建远程目录
	remotePath := "/opt/agent"
	if targetOS == "windows" {
		remotePath = "C:\\dodevops-agent"
	}

	createDirCmd := fmt.Sprintf("mkdir -p %s", remotePath)
	if targetOS == "windows" {
		createDirCmd = fmt.Sprintf("mkdir %s", remotePath)
	}

	if err := ExecuteSSHCommand(sshConfig, createDirCmd); err != nil {
		return fmt.Errorf("创建远程目录失败: %v", err)
	}

	// 拷贝二进制文件
	remoteBinaryPath := filepath.Join(remotePath, binaryName)
	if err := CopyFileViaSSH(sshConfig, binaryPath, remoteBinaryPath); err != nil {
		return fmt.Errorf("拷贝文件失败: %v", err)
	}

	log.Printf("主机 %d: 文件拷贝成功", host.ID)

	// 更新进度：传输完成
	s.agentDao.UpdateInstallProgress(host.ID, agentModel.InstallProgressTransferred)

	// 设置权限和systemd服务
	if targetOS != "windows" {
		chmodCmd := fmt.Sprintf("chmod +x %s", remoteBinaryPath)
		if err := ExecuteSSHCommand(sshConfig, chmodCmd); err != nil {
			return fmt.Errorf("设置执行权限失败: %v", err)
		}

		// 创建systemd服务
		serviceContent := `[Unit]
Description=DevOps Monitoring Agent
After=network.target

[Service]
Type=simple
ExecStart=/opt/agent/dodevops-agent
Restart=always
RestartSec=5
User=root
WorkingDirectory=/opt/agent
StandardOutput=append:/var/log/dodevops-agent.log
StandardError=append:/var/log/dodevops-agent.log

[Install]
WantedBy=multi-user.target`

		createServiceCmd := fmt.Sprintf("sudo bash -c 'cat > /etc/systemd/system/agent.service << \"EOF\"\n%s\nEOF'", serviceContent)
		if err := ExecuteSSHCommand(sshConfig, createServiceCmd); err != nil {
			log.Printf("主机 %d: 创建systemd服务失败，使用备用方法: %v", host.ID, err)
			// 备用方法
			tempServiceFile := "/tmp/dodevops-agent.service"
			createTempCmd := fmt.Sprintf("cat > %s << 'EOF'\n%s\nEOF", tempServiceFile, serviceContent)
			if err := ExecuteSSHCommand(sshConfig, createTempCmd); err == nil {
				moveCmd := fmt.Sprintf("sudo mv %s /etc/systemd/system/agent.service", tempServiceFile)
				ExecuteSSHCommand(sshConfig, moveCmd)
			}
		}

		ExecuteSSHCommand(sshConfig, "sudo systemctl daemon-reload")
		ExecuteSSHCommand(sshConfig, "sudo systemctl enable agent.service")
	}

	// 更新进度：配置完成
	s.agentDao.UpdateInstallProgress(host.ID, agentModel.InstallProgressConfigured)

	// 停止已存在的进程
	stopCommands := []string{
		"sudo systemctl stop agent.service 2>/dev/null || true",
		"sudo pkill -f dodevops-agent 2>/dev/null || true",
	}
	for _, stopCmd := range stopCommands {
		ExecuteSSHCommand(sshConfig, stopCmd)
	}
	time.Sleep(2 * time.Second)

	// 启动服务
	startCmd := s.getStartAgentCommand(targetOS, remoteBinaryPath)
	if err := ExecuteSSHCommand(sshConfig, startCmd); err != nil {
		log.Printf("主机 %d: systemd启动失败，尝试备用方式: %v", host.ID, err)
		// 备用启动
		directRunCmd := fmt.Sprintf("cd %s && nohup %s > /var/log/dodevops-agent.log 2>&1 </dev/null &", remotePath, remoteBinaryPath)
		if err := ExecuteSSHCommand(sshConfig, directRunCmd); err != nil {
			return fmt.Errorf("启动服务失败: %v", err)
		}
	}

	// 验证启动
	time.Sleep(5 * time.Second)
	checkProcessCmd := "ps aux | grep dodevops-agent | grep -v grep"
	if err := ExecuteSSHCommand(sshConfig, checkProcessCmd); err != nil {
		return fmt.Errorf("服务启动验证失败: 进程未运行")
	}

	// 更新状态为运行中
	s.agentDao.UpdateByHostID(host.ID, map[string]interface{}{
		"status":           agentModel.AgentStatusRunning,
		"install_progress": agentModel.InstallProgressStarted,
		"error_msg":        "",
		"update_time":      time.Now(),
	})

	log.Printf("主机 %d: Agent部署成功", host.ID)
	return nil
}

// deployAgentAsync 异步部署agent
func (s *AgentServiceImpl) deployAgentAsync(hostID uint, host model.CmdbHost) error {
	// 检测目标操作系统
	targetOS := "linux" // 默认linux
	if strings.Contains(strings.ToLower(host.OS), "windows") {
		targetOS = "windows"
	} else if strings.Contains(strings.ToLower(host.OS), "darwin") {
		targetOS = "darwin"
	}

	buildTarget, exists := buildTargets[targetOS]
	if !exists {
		return fmt.Errorf("不支持的操作系统: %s", host.OS)
	}

	// 更新进度：开始编译 (10%)
	s.agentDao.UpdateInstallProgress(hostID, agentModel.InstallProgressCompiling)

	log.Printf("开始为主机 %d (%s) 部署Agent，目标OS: %s", hostID, host.HostName, targetOS)

	// 1. 创建临时目录
	tempDir := fmt.Sprintf("/tmp/agent-deploy-%d-%d", hostID, time.Now().Unix())
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return fmt.Errorf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(tempDir) // 清理临时目录

	// 2. 生成main.go文件
	log.Printf("开始生成Agent main.go文件...")

	// 从配置文件读取心跳URL和Token（在部署时读取并写死到agent中）
	heartbeatURL := "https://www.deviops.cn/api/v1/monitor/agent/heartbeat"
	heartbeatToken := "agent-heartbeat-token-2024"
	pushgatewayURL := "http://8.130.14.34:9091"

	if config.Config != nil && config.Config.Monitor.Agent.HeartbeatServerURL != "" {
		heartbeatURL = config.Config.Monitor.Agent.HeartbeatServerURL
		log.Printf("从配置文件读取心跳URL: %s", heartbeatURL)
	}
	if config.Config != nil && config.Config.Monitor.Agent.HeartbeatToken != "" {
		heartbeatToken = config.Config.Monitor.Agent.HeartbeatToken
	}
	if config.Config != nil && config.Config.Monitor.Pushgateway.URL != "" {
		pushgatewayURL = config.Config.Monitor.Pushgateway.URL
		log.Printf("从配置文件读取Pushgateway URL: %s", pushgatewayURL)
	}

	mainGoContent := agent.CreateAgentMainFile(heartbeatURL, heartbeatToken, pushgatewayURL)
	log.Printf("Agent main.go内容长度: %d", len(mainGoContent))

	if len(mainGoContent) == 0 {
		return fmt.Errorf("生成的main.go内容为空")
	}

	// 验证生成的代码
	if err := agent.ValidateAgentCode(mainGoContent); err != nil {
		log.Printf("Agent代码验证失败: %v", err)
		return fmt.Errorf("生成的Agent代码有语法错误: %v", err)
	}
	log.Printf("Agent代码验证通过")

	mainGoPath := filepath.Join(tempDir, "main.go")
	log.Printf("写入main.go文件到: %s", mainGoPath)

	if err := os.WriteFile(mainGoPath, []byte(mainGoContent), 0644); err != nil {
		return fmt.Errorf("创建main.go文件失败: %v", err)
	}

	log.Printf("成功创建main.go文件，大小: %d bytes", len(mainGoContent))

	// 3. 创建go.mod文件 - 使用独立模块，不依赖本地源码
	goModContent := `module agent

go 1.24

require (
	github.com/prometheus/client_golang v1.17.0
	github.com/shirou/gopsutil/v3 v3.23.10
)`

	goModPath := filepath.Join(tempDir, "go.mod")
	if err := os.WriteFile(goModPath, []byte(goModContent), 0644); err != nil {
		return fmt.Errorf("创建go.mod文件失败: %v", err)
	}

	// 3.1 设置Go环境变量 - 增强生产环境兼容性
	goTempDir := filepath.Join(tempDir, ".gocache")
	goCacheDir := filepath.Join(goTempDir, "build")
	goTmpDir := filepath.Join(goTempDir, "tmp")
	goModCacheDir := filepath.Join(goTempDir, "mod")

	// 创建必要的Go缓存目录
	if err := os.MkdirAll(goCacheDir, 0755); err != nil {
		return fmt.Errorf("创建Go构建缓存目录失败: %v", err)
	}
	if err := os.MkdirAll(goTmpDir, 0755); err != nil {
		return fmt.Errorf("创建Go临时目录失败: %v", err)
	}
	if err := os.MkdirAll(goModCacheDir, 0755); err != nil {
		return fmt.Errorf("创建Go模块缓存目录失败: %v", err)
	}

	log.Printf("Go环境目录设置: GOMODCACHE=%s, GOCACHE=%s", goModCacheDir, goCacheDir)

	// 基础环境变量（所有命令都会使用）
	baseEnvVars := []string{
		"GO111MODULE=on",
		"CGO_ENABLED=0",
		"GOFLAGS=-buildvcs=false", // 禁用VCS信息，避免生产环境权限问题
		fmt.Sprintf("GOMODCACHE=%s", goModCacheDir), // 设置模块缓存目录
		fmt.Sprintf("GOCACHE=%s", goCacheDir), // 设置构建缓存
		fmt.Sprintf("GOTMPDIR=%s", goTmpDir), // 设置临时目录
		fmt.Sprintf("GOPATH=%s", goTempDir), // 设置GOPATH作为备用
	}

	// 下载专用环境变量
	downloadEnvVars := append(baseEnvVars,
		"GOPROXY=https://goproxy.cn,https://proxy.golang.org,direct",
		"GOSUMDB=sum.golang.google.cn",
	)

	// 构建专用环境变量
	goEnvVars := append(baseEnvVars,
		"GOPROXY=https://goproxy.cn,direct",
		"GOSUMDB=off", // 生产环境兼容性
	)

	// 3.2 运行 go mod download 先下载依赖
	log.Printf("开始下载Go模块依赖...")
	downloadCmd := exec.Command("go", "mod", "download")
	downloadCmd.Dir = tempDir
	downloadCmd.Env = append(os.Environ(), downloadEnvVars...)

	// 打印调试信息
	log.Printf("执行命令: go mod download")
	log.Printf("工作目录: %s", tempDir)
	for _, env := range downloadEnvVars {
		if strings.HasPrefix(env, "GO") {
			log.Printf("环境变量: %s", env)
		}
	}

	if output, err := downloadCmd.CombinedOutput(); err != nil {
		log.Printf("go mod download 失败: %v, 输出: %s", err, string(output))
		// 尝试使用直连模式并设置完整环境变量
		log.Printf("尝试直连模式重新下载...")
		retryDownloadCmd := exec.Command("go", "mod", "download")
		retryDownloadCmd.Dir = tempDir
		retryDownloadEnv := append(baseEnvVars, "GOPROXY=direct", "GOSUMDB=off")
		retryDownloadCmd.Env = append(os.Environ(), retryDownloadEnv...)

		if output, err := retryDownloadCmd.CombinedOutput(); err != nil {
			log.Printf("直连模式下载依赖失败: %v, 输出: %s", err, string(output))
			// 继续尝试，不中断流程
		} else {
			log.Printf("直连模式下载依赖成功")
		}
	} else {
		log.Printf("go mod download 成功")
	}

	// 3.3 运行 go mod tidy 来确保依赖正确
	log.Printf("开始运行 go mod tidy...")
	tidyCmd := exec.Command("go", "mod", "tidy")
	tidyCmd.Dir = tempDir
	tidyCmd.Env = append(os.Environ(), goEnvVars...)

	if output, err := tidyCmd.CombinedOutput(); err != nil {
		log.Printf("go mod tidy 失败: %v, 输出: %s", err, string(output))

		// 尝试清理模块缓存并重试
		log.Printf("清理模块缓存并重试...")
		cleanCmd := exec.Command("go", "clean", "-modcache")
		cleanCmd.Dir = tempDir
		cleanCmd.Env = append(os.Environ(), goEnvVars...)
		cleanCmd.Run()

		// 重新创建缓存目录
		os.MkdirAll(goModCacheDir, 0755)
		os.MkdirAll(goCacheDir, 0755)

		// 重试with更严格的环境变量设置
		retryTidyCmd := exec.Command("go", "mod", "tidy")
		retryTidyCmd.Dir = tempDir
		retryEnv := append(baseEnvVars,
			"GOPROXY=https://goproxy.cn,direct",
			"GOSUMDB=off",
		)
		retryTidyCmd.Env = append(os.Environ(), retryEnv...)

		if output, err := retryTidyCmd.CombinedOutput(); err != nil {
			return fmt.Errorf("运行 go mod tidy 失败: %v, 输出: %s", err, string(output))
		} else {
			log.Printf("重试 go mod tidy 成功")
		}
	} else {
		log.Printf("go mod tidy 成功")
	}

	// 4. 编译二进制文件
	// 更新进度：编译完成 (30%)
	s.agentDao.UpdateInstallProgress(hostID, agentModel.InstallProgressCompiled)

	binaryName := fmt.Sprintf("dodevops-agent%s", buildTarget.Suffix)
	binaryPath := filepath.Join(tempDir, binaryName)

	// 生产环境独立编译 - 确保完全静态链接
	log.Printf("开始生产环境独立编译...")
	buildFlags := []string{
		"build",
		"-a", // 强制重新构建所有包
		"-ldflags", "-s -w -extldflags '-static'", // 静态链接
		"-installsuffix", "netgo", // 使用纯Go网络栈
		"-tags", "netgo,osusergo", // 禁用CGO依赖
		"-o", binaryPath,
		mainGoPath,
	}

	cmd := exec.Command("go", buildFlags...)
	cmd.Dir = tempDir

	// 构建环境变量（包含目标平台设置）
	buildEnvVars := append(goEnvVars,
		fmt.Sprintf("GOOS=%s", buildTarget.GOOS),
		fmt.Sprintf("GOARCH=%s", buildTarget.GOARCH),
		"GODEBUG=netdns=go", // 使用Go DNS解析器
	)
	cmd.Env = append(os.Environ(), buildEnvVars...)

	log.Printf("开始编译二进制文件...")
	log.Printf("编译命令: %s", strings.Join(append([]string{"go"}, buildFlags...), " "))
	log.Printf("工作目录: %s", tempDir)
	log.Printf("目标平台: GOOS=%s GOARCH=%s", buildTarget.GOOS, buildTarget.GOARCH)
	log.Printf("Go模块缓存: GOMODCACHE=%s", goModCacheDir)
	log.Printf("Go构建缓存: GOCACHE=%s", goCacheDir)

	// 打印关键环境变量用于调试
	for _, env := range buildEnvVars {
		if strings.HasPrefix(env, "GO") {
			log.Printf("构建环境变量: %s", env)
		}
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("编译失败，详细信息:")
		log.Printf("  错误: %v", err)
		log.Printf("  输出: %s", string(output))
		log.Printf("  二进制路径: %s", binaryPath)
		log.Printf("  工作目录内容:")
		if files, err := os.ReadDir(tempDir); err == nil {
			for _, file := range files {
				log.Printf("    %s", file.Name())
			}
		}
		return fmt.Errorf("编译失败: %v, 输出: %s", err, string(output))
	}

	log.Printf("编译成功！")
	if len(output) > 0 {
		log.Printf("编译输出: %s", string(output))
	}

	// 本地验证编译结果
	localVerifyCommands := []string{
		fmt.Sprintf("ls -la %s", binaryPath),
		fmt.Sprintf("file %s", binaryPath),
		fmt.Sprintf("du -h %s", binaryPath),
	}

	// 如果是Linux二进制文件且在Linux环境下，可以检查依赖
	if buildTarget.GOOS == "linux" {
		localVerifyCommands = append(localVerifyCommands, fmt.Sprintf("ldd %s 2>/dev/null || echo 'Static binary (no dynamic dependencies)'", binaryPath))
	}

	log.Printf("本地验证编译结果:")
	for _, verifyCmd := range localVerifyCommands {
		if verifyOutput, err := exec.Command("bash", "-c", verifyCmd).CombinedOutput(); err == nil {
			log.Printf("  %s -> %s", verifyCmd, strings.TrimSpace(string(verifyOutput)))
		} else {
			log.Printf("  %s -> 验证失败: %v", verifyCmd, err)
		}
	}

	// 5. 通过SSH拷贝文件到目标主机
	// 更新进度：开始传输 (50%)
	s.agentDao.UpdateInstallProgress(hostID, agentModel.InstallProgressTransfer)

	remotePath := "/opt/agent"
	if targetOS == "windows" {
		remotePath = "C:\\dodevops-agent"
	}

	// 获取SSH认证信息
	var sshConfig util.SSHConfig
	if host.SSHKeyID > 0 {
		sshKey, err := s.getSSHKeyByID(host.SSHKeyID)
		if err != nil {
			return fmt.Errorf("获取SSH认证信息失败: %v", err)
		}

		// 创建SSH连接配置
		sshConfig = util.SSHConfig{
			IP:       host.SSHIP,
			Port:     host.SSHPort,
			Username: host.SSHName,
			Type:     sshKey.Type, // 使用数据库中的认证类型
			Timeout:  30 * time.Second,
		}

		// 根据认证类型设置相应的认证信息
		switch sshKey.Type {
		case 1: // 密码认证
			sshConfig.Password = sshKey.Password
		case 2: // 私钥认证
			sshConfig.PublicKey = sshKey.PublicKey
		case 3: // 公钥免认证
			// 不需要设置额外信息，自动查找本地私钥
		default:
			return fmt.Errorf("不支持的SSH认证类型: %d", sshKey.Type)
		}
	} else {
		return fmt.Errorf("主机未配置SSH认证信息")
	}


	// 创建远程目录并进行环境预检查
	log.Printf("创建远程目录并检查环境兼容性...")

	// 生产环境兼容性检查 - 不依赖Go环境
	envCheckCmds := []string{
		"whoami",
		"id",
		"uname -a",
		"cat /etc/os-release | head -5 || cat /etc/redhat-release || echo 'Unknown OS'",
		"df -h /opt || df -h /",
		"free -h",
		"which systemctl && echo 'systemctl available' || echo 'systemctl not found, will use alternative startup'",
		"ps --version 2>/dev/null && echo 'ps available' || echo 'ps command check'",
		"netstat --version 2>/dev/null || ss --version 2>/dev/null || echo 'network tools check'",
		"ls -la /lib/x86_64-linux-gnu/ 2>/dev/null | head -3 || echo 'lib check'",
		"ldd --version 2>/dev/null || echo 'ldd not available'",
	}

	log.Printf("环境检查结果:")
	for _, checkCmd := range envCheckCmds {
		if output, err := ExecuteSSHCommandWithOutput(sshConfig, checkCmd); err == nil {
			log.Printf("  %s: %s", checkCmd, strings.TrimSpace(output))
		}
	}

	createDirCmd := fmt.Sprintf("mkdir -p %s", remotePath)
	if targetOS == "windows" {
		createDirCmd = fmt.Sprintf("mkdir %s", remotePath)
	}

	if err := ExecuteSSHCommand(sshConfig, createDirCmd); err != nil {
		return fmt.Errorf("创建远程目录失败: %v", err)
	}

	// 确保完整的生产环境目录结构
	if targetOS != "windows" {
		log.Printf("设置生产环境目录结构...")

		// 创建必要的目录结构
		dirSetupCmds := []string{
			// 创建主目录
			fmt.Sprintf("sudo mkdir -p %s", remotePath),
			// 创建日志目录
			"sudo mkdir -p /var/log",
			// 创建临时目录(确保agent有写权限)
			"sudo mkdir -p /tmp/agent",
			// 设置目录权限
			fmt.Sprintf("sudo chown -R root:root %s", remotePath),
			fmt.Sprintf("sudo chmod 755 %s", remotePath),
			"sudo chmod 777 /tmp/agent",
			// 创建systemd目录(如果不存在)
			"sudo mkdir -p /etc/systemd/system",
			// 确保日志目录权限
			"sudo chmod 755 /var/log",
		}

		for _, cmd := range dirSetupCmds {
			if err := ExecuteSSHCommand(sshConfig, cmd); err != nil {
				log.Printf("目录设置命令失败 [%s]: %v", cmd, err)
			}
		}

		// 检查依赖库是否存在(生产环境兼容性)
		libCheckCmds := []string{
			"ls -la /lib/x86_64-linux-gnu/libc.so* 2>/dev/null | head -2 || echo 'glibc check'",
			"ls -la /lib64/ld-linux-x86-64.so* 2>/dev/null || echo 'loader check'",
			"cat /etc/ld.so.conf 2>/dev/null | head -3 || echo 'ld config check'",
		}

		log.Printf("检查系统依赖库:")
		for _, cmd := range libCheckCmds {
			if output, err := ExecuteSSHCommandWithOutput(sshConfig, cmd); err == nil && output != "" {
				log.Printf("  %s", strings.TrimSpace(output))
			}
		}
	}

	// 拷贝二进制文件
	remoteBinaryPath := filepath.Join(remotePath, binaryName)
	if err := CopyFileViaSSH(sshConfig, binaryPath, remoteBinaryPath); err != nil {
		return fmt.Errorf("拷贝文件失败: %v", err)
	}

	log.Printf("验证生产环境二进制文件完整性...")

	// 验证二进制文件是否正确传输
	verifyCommands := []string{
		// 检查文件是否存在
		fmt.Sprintf("ls -la %s", remoteBinaryPath),
		// 检查文件大小
		fmt.Sprintf("du -h %s", remoteBinaryPath),
		// 检查文件类型
		fmt.Sprintf("file %s", remoteBinaryPath),
		// 检查是否可执行
		fmt.Sprintf("test -x %s && echo 'executable' || echo 'not executable'", remoteBinaryPath),
	}

	for _, cmd := range verifyCommands {
		if output, err := ExecuteSSHCommandWithOutput(sshConfig, cmd); err == nil {
			log.Printf("文件验证: %s -> %s", cmd, strings.TrimSpace(output))
		} else {
			log.Printf("文件验证失败: %s -> %v", cmd, err)
		}
	}

	// 检查动态链接依赖(生产环境关键)
	if targetOS != "windows" {
		lddCmd := fmt.Sprintf("ldd %s 2>/dev/null | head -10 || echo 'Static binary or ldd not available'", remoteBinaryPath)
		if output, err := ExecuteSSHCommandWithOutput(sshConfig, lddCmd); err == nil && output != "" {
			log.Printf("二进制依赖检查: %s", strings.TrimSpace(output))
		}

		// 尝试快速测试运行(超时保护)
		testCmd := fmt.Sprintf("timeout 3 %s --help 2>/dev/null || timeout 3 %s -h 2>/dev/null || echo 'Binary help test completed'", remoteBinaryPath, remoteBinaryPath)
		if output, err := ExecuteSSHCommandWithOutput(sshConfig, testCmd); err == nil {
			log.Printf("二进制测试运行: %s", strings.TrimSpace(output))
		}
	}

	// 更新进度：传输完成 (70%)
	s.agentDao.UpdateInstallProgress(hostID, agentModel.InstallProgressTransferred)

	// 6. 设置文件权限并创建systemd服务
	if targetOS != "windows" {
		// Linux/Unix系统设置执行权限
		chmodCmd := fmt.Sprintf("chmod +x %s", remoteBinaryPath)
		if err := ExecuteSSHCommand(sshConfig, chmodCmd); err != nil {
			return fmt.Errorf("设置执行权限失败: %v", err)
		}
		// 创建systemd服务文件 - 增强生产环境兼容性
		serviceContent := `[Unit]
Description=DevOps Monitoring Agent
Documentation=https://dodevops.com
After=network.target network-online.target
Wants=network-online.target
StartLimitIntervalSec=60
StartLimitBurst=3

[Service]
Type=simple
ExecStart=/opt/agent/dodevops-agent
ExecReload=/bin/kill -HUP $MAINPID
Restart=always
RestartSec=5
User=root
Group=root
WorkingDirectory=/opt/agent
StandardOutput=append:/var/log/dodevops-agent.log
StandardError=append:/var/log/dodevops-agent.log
Environment=HOME=/root
Environment=PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
KillMode=mixed
KillSignal=SIGTERM
TimeoutStartSec=60
TimeoutStopSec=30

[Install]
WantedBy=multi-user.target`

		// 方法1: 使用cat heredoc写入服务文件
		createServiceCmd := fmt.Sprintf("sudo bash -c 'cat > /etc/systemd/system/agent.service << \"EOF\"\n%s\nEOF'", serviceContent)
		if err := ExecuteSSHCommand(sshConfig, createServiceCmd); err != nil {
			log.Printf("方法1创建systemd服务文件失败: %v，尝试方法2", err)

			// 方法2: 使用tee命令
			createServiceCmd = fmt.Sprintf("echo '%s' | sudo tee /etc/systemd/system/agent.service > /dev/null", strings.ReplaceAll(serviceContent, "'", "'\"'\"'"))
			if err := ExecuteSSHCommand(sshConfig, createServiceCmd); err != nil {
				log.Printf("方法2创建systemd服务文件失败: %v，尝试方法3", err)

				// 方法3: 先创建临时文件，再移动
				tempServiceFile := "/tmp/dodevops-agent.service"
				createTempCmd := fmt.Sprintf("cat > %s << 'EOF'\n%s\nEOF", tempServiceFile, serviceContent)
				if err := ExecuteSSHCommand(sshConfig, createTempCmd); err == nil {
					moveCmd := fmt.Sprintf("sudo mv %s /etc/systemd/system/agent.service && sudo chown root:root /etc/systemd/system/agent.service", tempServiceFile)
					if err := ExecuteSSHCommand(sshConfig, moveCmd); err != nil {
						return fmt.Errorf("移动systemd服务文件失败: %v", err)
					}
				} else {
					return fmt.Errorf("创建systemd服务文件失败: %v", err)
				}
			}
		}

		// 验证服务文件是否创建成功
		checkServiceCmd := "sudo test -f /etc/systemd/system/agent.service && echo 'Service file created successfully'"
		if err := ExecuteSSHCommand(sshConfig, checkServiceCmd); err != nil {
			return fmt.Errorf("验证systemd服务文件失败: %v", err)
		}

		// 重新加载systemd配置
		if err := ExecuteSSHCommand(sshConfig, "sudo systemctl daemon-reload"); err != nil {
			return fmt.Errorf("重新加载systemd配置失败: %v", err)
		}

		// 启用服务开机自启
		if err := ExecuteSSHCommand(sshConfig, "sudo systemctl enable agent.service"); err != nil {
			return fmt.Errorf("启用agent服务失败: %v", err)
		}
	}

	// 更新进度：配置完成 (90%)
	s.agentDao.UpdateInstallProgress(hostID, agentModel.InstallProgressConfigured)

	// 7. 停止已存在的agent进程（如果有）
	log.Printf("停止现有agent进程...")
	stopCommands := []string{
		"sudo systemctl stop agent.service 2>/dev/null || true",
		"sudo pkill -f dodevops-agent 2>/dev/null || true",
		"killall dodevops-agent 2>/dev/null || true",
		"ps aux | grep dodevops-agent | grep -v grep | awk '{print $2}' | xargs sudo kill -9 2>/dev/null || true",
	}
	for _, stopCmd := range stopCommands {
		ExecuteSSHCommand(sshConfig, stopCmd) // 忽略错误，进程可能不存在
	}

	// 等待进程完全停止
	time.Sleep(2 * time.Second)

	// 8. 预检查：验证二进制文件可执行
	testRunCmd := fmt.Sprintf("timeout 5 %s --help 2>/dev/null || echo 'Binary test completed'", remoteBinaryPath)
	ExecuteSSHCommand(sshConfig, testRunCmd) // 忽略错误，只是测试

	// 9. 启动agent服务 - 增强启动策略
	log.Printf("启动agent服务...")
	startCmd := s.getStartAgentCommand(targetOS, remoteBinaryPath)
	startSuccess := false

	// 尝试systemd启动
	if err := ExecuteSSHCommand(sshConfig, startCmd); err != nil {
		log.Printf("systemd启动失败: %v，尝试重置服务状态", err)

		// 重置服务状态并重试
		resetCommands := []string{
			"sudo systemctl reset-failed agent.service 2>/dev/null || true",
			"sudo systemctl daemon-reload",
		}
		for _, resetCmd := range resetCommands {
			ExecuteSSHCommand(sshConfig, resetCmd)
		}

		// 再次尝试启动
		if err := ExecuteSSHCommand(sshConfig, startCmd); err != nil {
			log.Printf("systemd二次启动失败: %v，尝试直接运行二进制文件", err)

			// 备用方案1：使用screen后台运行
			screenCmd := fmt.Sprintf("screen -dmS dodevops-agent %s", remoteBinaryPath)
			if err := ExecuteSSHCommand(sshConfig, screenCmd); err != nil {
				log.Printf("screen启动失败: %v，尝试nohup方式", err)

				// 备用方案2：直接后台运行
				directRunCmd := fmt.Sprintf("cd %s && nohup %s > /var/log/dodevops-agent.log 2>&1 </dev/null &", remotePath, remoteBinaryPath)
				if err := ExecuteSSHCommand(sshConfig, directRunCmd); err != nil {
					// 最后尝试：简单后台启动
					simpleCmd := fmt.Sprintf("%s &", remoteBinaryPath)
					if err := ExecuteSSHCommand(sshConfig, simpleCmd); err != nil {
						return fmt.Errorf("所有启动方式都失败: %v", err)
					}
					log.Printf("使用简单后台启动成功")
					startSuccess = true
				} else {
					log.Printf("使用nohup启动成功")
					startSuccess = true
				}
			} else {
				log.Printf("使用screen启动成功")
				startSuccess = true
			}
		} else {
			log.Printf("systemd二次启动成功")
			startSuccess = true
		}
	} else {
		log.Printf("systemd启动成功")
		startSuccess = true
	}

	if !startSuccess {
		return fmt.Errorf("无法启动agent服务")
	}

	// 10. 综合验证服务是否启动成功 - 增强验证策略
	log.Printf("开始验证agent服务启动状态...")

	// 分阶段验证，增加容错性
	maxRetries := 3
	verifySuccess := false

	for i := 0; i < maxRetries && !verifySuccess; i++ {
		log.Printf("验证阶段 %d/%d...", i+1, maxRetries)

		// 等待服务启动
		waitTime := time.Duration(5+i*2) * time.Second
		log.Printf("等待 %v 后验证...", waitTime)
		time.Sleep(waitTime)

		// 步骤1: 检查进程是否运行
		processRunning := false
		checkProcessCmds := []string{
			"ps aux | grep dodevops-agent | grep -v grep",
			"pgrep -f dodevops-agent",
			"pidof dodevops-agent",
		}

		for _, processCmd := range checkProcessCmds {
			if err := ExecuteSSHCommand(sshConfig, processCmd); err == nil {
				processRunning = true
				log.Printf("进程检查通过: %s", processCmd)
				break
			}
		}

		if !processRunning {
			log.Printf("第%d次验证: agent进程未运行", i+1)
			if i == maxRetries-1 {
				// 最后一次失败，收集详细错误日志
				logCommands := []string{
					"tail -50 /var/log/dodevops-agent.log 2>/dev/null",
					"journalctl -u agent.service -n 50 --no-pager 2>/dev/null",
					"tail -20 /tmp/agent.log 2>/dev/null",
					"systemctl status agent.service --no-pager 2>/dev/null",
				}

				allLogs := make([]string, 0)
				for _, logCmd := range logCommands {
					if output, err := ExecuteSSHCommandWithOutput(sshConfig, logCmd); err == nil && output != "" {
						allLogs = append(allLogs, fmt.Sprintf("=== %s ===\n%s", logCmd, output))
					}
				}

				errorLog := "无可用日志"
				if len(allLogs) > 0 {
					errorLog = strings.Join(allLogs, "\n\n")
				}

				return fmt.Errorf("agent进程未运行，详细日志:\n%s", errorLog)
			}
			continue
		}

		// 步骤2: 检查端口是否监听 (可选检查，允许失败)
		portListening := false
		checkPortCmds := []string{
			"netstat -tunlp | grep :9100",
			"ss -tunlp | grep :9100",
			"lsof -i :9100",
		}

		for _, portCmd := range checkPortCmds {
			if err := ExecuteSSHCommand(sshConfig, portCmd); err == nil {
				portListening = true
				log.Printf("端口检查通过: %s", portCmd)
				break
			}
		}

		if !portListening {
			log.Printf("警告: 端口9100暂未监听，可能需要更多时间初始化")
			// 不作为失败条件，继续验证其他方面
		}

		// 步骤3: 尝试HTTP健康检查 (可选)
		healthCheck := false
		checkHealthCmd := "curl -s http://localhost:9100/metrics | head -1 2>/dev/null || wget -qO- http://localhost:9100/metrics | head -1 2>/dev/null"
		if err := ExecuteSSHCommand(sshConfig, checkHealthCmd); err == nil {
			healthCheck = true
			log.Printf("HTTP健康检查通过")
		} else {
			log.Printf("HTTP健康检查失败，可能服务还在初始化")
		}

		// 综合判断
		if processRunning {
			if portListening || healthCheck {
				verifySuccess = true
				log.Printf("验证成功: 进程运行中且服务正常")
			} else if i == maxRetries-1 {
				// 最后一次验证：进程运行但端口/服务异常，记录警告但不失败
				log.Printf("警告: 进程运行但端口检查失败，可能需要更长时间初始化")
				verifySuccess = true
			}
		}
	}

	if !verifySuccess {
		return fmt.Errorf("服务验证失败: 经过%d次验证仍无法确认服务正常启动", maxRetries)
	}

	// 更新进度：启动成功 (100%) 并同时更新状态为运行中
	updateErr := s.agentDao.UpdateByHostID(hostID, map[string]interface{}{
		"status": agentModel.AgentStatusRunning,
		"install_progress": agentModel.InstallProgressStarted,
		"error_msg": "",
		"update_time": time.Now(),
	})
	if updateErr != nil {
		log.Printf("Failed to update agent status to running for host %d: %v", hostID, updateErr)
		return fmt.Errorf("更新Agent状态失败: %v", updateErr)
	}

	log.Printf("Agent successfully deployed to host %d, status updated to running", hostID)
	return nil
}

// getStopAgentCommand 获取停止agent的命令
func (s *AgentServiceImpl) getStopAgentCommand(targetOS string) string {
	switch targetOS {
	case "windows":
		return "taskkill /F /IM dodevops-agent.exe"
	default:
		return "sudo systemctl stop agent.service"
	}
}

// getStartAgentCommand 获取启动agent的命令
func (s *AgentServiceImpl) getStartAgentCommand(targetOS, binaryPath string) string {
	switch targetOS {
	case "windows":
		return fmt.Sprintf("start /B %s", binaryPath)
	default:
		return "sudo systemctl start agent.service"
	}
}

// getCheckAgentCommand 获取检查agent状态的命令
func (s *AgentServiceImpl) getCheckAgentCommand(targetOS string) string {
	switch targetOS {
	case "windows":
		return "tasklist | findstr dodevops-agent.exe"
	default:
		return "systemctl status agent.service | grep Active | awk '{print $3}' | awk -F'(' '{print $2}' | awk -F')' '{print $1}'"
	}
}

// UninstallAgent 卸载指定主机的agent
func (s *AgentServiceImpl) UninstallAgent(c *gin.Context, hostID uint) {
	// 参数校验
	if hostID == 0 {
		result.FailedWithCode(c, int(result.ApiCode.ValidationParameterError), "主机ID不能为空")
		return
	}

	// 获取主机信息
	host, err := s.hostDao.GetCmdbHostById(hostID)
	if err != nil {
		result.FailedWithCode(c, int(result.ApiCode.FAILED), "获取主机信息失败")
		return
	}

	// 异步卸载
	go func() {
		uninstallErr := s.uninstallAgentAsync(host.ID, host)
		if uninstallErr != nil {
			log.Printf("Uninstall agent from host %d failed: %v", hostID, uninstallErr)
			redis.RedisDb.Set(c, fmt.Sprintf("agent:uninstall:error:%d", hostID), uninstallErr.Error(), 30*time.Minute)
		}
	}()

	result.Success(c, map[string]interface{}{
		"message": "Agent卸载任务已启动，正在后台执行",
		"hostId":  hostID,
		"status":  "uninstalling",
	})
}

// BatchUninstallAgent 批量卸载agent
func (s *AgentServiceImpl) BatchUninstallAgent(c *gin.Context, dto *agentModel.BatchDeployAgentDto) {
	// 参数校验
	if len(dto.HostIDs) == 0 {
		result.FailedWithCode(c, int(result.ApiCode.ValidationParameterError), "主机ID列表不能为空")
		return
	}

	var successHosts []uint
	var failedHosts []map[string]interface{}

	// 并发卸载每台主机
	for _, hostID := range dto.HostIDs {
		// 获取主机信息
		host, err := s.hostDao.GetCmdbHostById(hostID)
		if err != nil {
			failedHosts = append(failedHosts, map[string]interface{}{
				"hostId": hostID,
				"error":  "获取主机信息失败",
			})
			continue
		}

		successHosts = append(successHosts, hostID)

		// 异步卸载
		go func(hID uint, h model.CmdbHost) {
			uninstallErr := s.uninstallAgentAsync(hID, h)
			if uninstallErr != nil {
				log.Printf("Uninstall agent from host %d failed: %v", hID, uninstallErr)
			}
		}(hostID, host)
	}

	result.Success(c, map[string]interface{}{
		"message":      fmt.Sprintf("批量卸载任务已启动，成功启动 %d 台，失败 %d 台", len(successHosts), len(failedHosts)),
		"successCount": len(successHosts),
		"failedCount":  len(failedHosts),
		"successHosts": successHosts,
		"failedHosts":  failedHosts,
		"status":       "uninstalling",
	})
}

// uninstallAgentAsync 异步卸载agent
func (s *AgentServiceImpl) uninstallAgentAsync(hostID uint, host model.CmdbHost) error {
	log.Printf("开始卸载主机 %d 的 agent...", hostID)

	// 检测目标操作系统
	targetOS := "linux"
	if strings.Contains(strings.ToLower(host.OS), "windows") {
		targetOS = "windows"
	} else if strings.Contains(strings.ToLower(host.OS), "darwin") {
		targetOS = "darwin"
	}

	// 获取SSH认证信息
	var sshConfig util.SSHConfig
	if host.SSHKeyID > 0 {
		sshKey, err := s.getSSHKeyByID(host.SSHKeyID)
		if err != nil {
			return fmt.Errorf("获取SSH认证信息失败: %v", err)
		}

		// 创建SSH连接配置
		sshConfig = util.SSHConfig{
			IP:       host.SSHIP,
			Port:     host.SSHPort,
			Username: host.SSHName,
			Type:     sshKey.Type, // 使用数据库中的认证类型
			Timeout:  30 * time.Second,
		}

		// 根据认证类型设置相应的认证信息
		switch sshKey.Type {
		case 1: // 密码认证
			sshConfig.Password = sshKey.Password
		case 2: // 私钥认证
			sshConfig.PublicKey = sshKey.PublicKey
		case 3: // 公钥免认证
			// 不需要设置额外信息，自动查找本地私钥
		default:
			return fmt.Errorf("不支持的SSH认证类型: %d", sshKey.Type)
		}
	} else {
		return fmt.Errorf("主机未配置SSH认证信息")
	}

	log.Printf("开始卸载主机 %d 的 agent...", hostID)

	// 1. 停止agent进程
	stopCmd := s.getStopAgentCommand(targetOS)
	log.Printf("执行停止命令: %s", stopCmd)
	ExecuteSSHCommand(sshConfig, stopCmd)

	// 2. 禁用并删除systemd服务(仅Linux)
	if targetOS != "windows" {
		// 禁用服务
		ExecuteSSHCommand(sshConfig, "sudo systemctl disable agent.service")
		// 删除服务文件
		ExecuteSSHCommand(sshConfig, "sudo rm -rf /etc/systemd/system/agent.service")
		// 重新加载systemd配置
		ExecuteSSHCommand(sshConfig, "sudo systemctl daemon-reload")
	}

	// 3. 删除agent文件和目录
	remotePath := "/opt/agent"
	if targetOS == "windows" {
		remotePath = "C:\\dodevops-agent"
	}

	var removeCmd string
	switch targetOS {
	case "windows":
		removeCmd = fmt.Sprintf("rmdir /S /Q %s", remotePath)
	default:
		removeCmd = fmt.Sprintf("rm -rf %s", remotePath)
	}

	log.Printf("执行删除命令: %s", removeCmd)
	if err := ExecuteSSHCommand(sshConfig, removeCmd); err != nil {
		return fmt.Errorf("删除agent文件失败: %v", err)
	}

	// 4. 清理Pushgateway中的监控数据
	s.cleanupPushgatewayMetrics(host.Name)

	// 5. 删除数据库中的Agent记录
	if err := s.agentDao.DeleteByHostID(hostID); err != nil {
		log.Printf("删除数据库Agent记录失败: %v", err)
		// 虽然删除数据库记录失败，但远程文件已删除，仍可认为卸载成功
		// 只记录错误日志，不返回错误
	}

	log.Printf("主机 %d 的 agent 卸载成功，已删除远程文件、Pushgateway数据和数据库记录", hostID)
	return nil
}

// RestartAgent 重启agent
func (s *AgentServiceImpl) RestartAgent(c *gin.Context, hostID uint) {
	// 参数校验
	if hostID == 0 {
		result.FailedWithCode(c, int(result.ApiCode.ValidationParameterError), "主机ID不能为空")
		return
	}

	// 获取Agent信息
	agent, err := s.agentDao.GetByHostID(hostID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Agent未部署，执行部署
			s.DeployAgent(c, hostID)
			return
		}
		result.FailedWithCode(c, int(result.ApiCode.FAILED), "获取Agent信息失败")
		return
	}

	// 检查Agent状态
	if agent.Status != agentModel.AgentStatusRunning {
		// 状态不是运行中，执行部署
		s.DeployAgent(c, hostID)
		return
	}

	// 获取主机信息
	host, err := s.hostDao.GetCmdbHostById(hostID)
	if err != nil {
		result.FailedWithCode(c, int(result.ApiCode.FAILED), "获取主机信息失败")
		return
	}

	// 异步重启
	go func() {
		// 1. 停止agent进程
		targetOS := s.detectPlatform(host.OS)
		stopCmd := s.getStopAgentCommand(targetOS)
		var sshConfig util.SSHConfig
		if host.SSHKeyID > 0 {
			sshKey, err := s.getSSHKeyByID(host.SSHKeyID)
			if err == nil {
				// 创建SSH连接配置
				sshConfig = util.SSHConfig{
					IP:       host.SSHIP,
					Port:     host.SSHPort,
					Username: host.SSHName,
					Type:     sshKey.Type, // 使用数据库中的认证类型
					Timeout:  30 * time.Second,
				}

				// 根据认证类型设置相应的认证信息
				switch sshKey.Type {
				case 1: // 密码认证
					sshConfig.Password = sshKey.Password
				case 2: // 私钥认证
					sshConfig.PublicKey = sshKey.PublicKey
				case 3: // 公钥免认证
					// 不需要设置额外信息，自动查找本地私钥
				}

			}
		}
		ExecuteSSHCommand(sshConfig, stopCmd)

		// 2. 启动agent进程
		startCmd := s.getStartAgentCommand(targetOS, agent.InstallPath)
		if err := ExecuteSSHCommand(sshConfig, startCmd); err != nil {
			s.agentDao.UpdateByHostID(hostID, map[string]interface{}{
				"status":    agentModel.AgentStatusStartError,
				"error_msg": fmt.Sprintf("重启失败: %v", err),
			})
			return
		}

		// 3. 验证启动状态
		time.Sleep(3 * time.Second)
		checkCmd := s.getCheckAgentCommand(targetOS)
		if err := ExecuteSSHCommand(sshConfig, checkCmd); err != nil {
			s.agentDao.UpdateByHostID(hostID, map[string]interface{}{
				"status":    agentModel.AgentStatusStartError,
				"error_msg": "重启后验证失败",
			})
			return
		}

		// 更新状态为运行中
		s.agentDao.UpdateByHostID(hostID, map[string]interface{}{
			"status":      agentModel.AgentStatusRunning,
			"update_time": time.Now(),
		})
	}()

	result.Success(c, map[string]interface{}{
		"message": "Agent重启任务已启动，正在后台执行",
		"hostId":  hostID,
		"status":  "restarting",
	})
}

// getCurrentWorkingDir 获取当前工作目录，确保返回包含go.mod的项目根目录
func getCurrentWorkingDir() string {
	dir, err := os.Getwd()
	if err != nil {
		// 如果获取失败，尝试找到项目根目录
		return findProjectRoot()
	}

	// 检查当前目录是否包含go.mod文件
	if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
		return dir
	}

	// 如果当前目录没有go.mod，尝试找到项目根目录
	return findProjectRoot()
}

// findProjectRoot 查找包含go.mod的项目根目录
func findProjectRoot() string {
	// 获取当前执行文件的路径
	execPath, err := os.Executable()
	if err != nil {
		return "/tmp" // 如果都失败了，返回临时目录
	}

	// 从执行文件所在目录向上查找go.mod
	dir := filepath.Dir(execPath)
	for i := 0; i < 10; i++ { // 最多向上查找10级目录
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir { // 到达根目录
			break
		}
		dir = parent
	}

	// 如果都找不到，使用一个临时的方案：不依赖本地go.mod
	return "/tmp"
}

// ExecuteSSHCommand 执行SSH命令
func ExecuteSSHCommand(config util.SSHConfig, command string) error {
	log.Printf("执行SSH命令: %s@%s:%d -> %s", config.Username, config.IP, config.Port, command)

	sshUtil := util.NewSSHUtil()
	output, err := sshUtil.ExecuteRemoteCommand(&config, command)

	if err != nil {
		log.Printf("SSH命令执行失败: %v, 输出: %s", err, output)
		return fmt.Errorf("SSH命令执行失败: %v (输出: %s)", err, output)
	}

	if output != "" {
		log.Printf("SSH命令输出: %s", strings.TrimSpace(output))
	}

	return nil
}

// ExecuteSSHCommandWithOutput 执行SSH命令并返回输出
func ExecuteSSHCommandWithOutput(config util.SSHConfig, command string) (string, error) {
	log.Printf("执行SSH命令(带输出): %s@%s:%d -> %s", config.Username, config.IP, config.Port, command)

	sshUtil := util.NewSSHUtil()
	output, err := sshUtil.ExecuteRemoteCommand(&config, command)

	if err != nil {
		log.Printf("SSH命令执行失败: %v, 输出: %s", err, output)
		return output, fmt.Errorf("SSH命令执行失败: %v", err)
	}

	if output != "" {
		log.Printf("SSH命令输出: %s", strings.TrimSpace(output))
	}

	return output, nil
}

// CopyFileViaSSH 通过SSH拷贝文件(使用SCP协议)
func CopyFileViaSSH(config util.SSHConfig, localPath, remotePath string) error {
	log.Printf("开始拷贝文件: %s -> %s@%s:%s", localPath, config.Username, config.IP, remotePath)

	// 验证文件路径是否存在并获取文件信息
	fileInfo, err := os.Stat(localPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("文件不存在: %s", localPath)
	}
	if err != nil {
		return fmt.Errorf("获取文件信息失败: %v", err)
	}

	log.Printf("源文件信息: %s (大小: %d bytes)", localPath, fileInfo.Size())

	// 创建SSH客户端配置
	sshConfig := &ssh.ClientConfig{
		User:            config.Username,
		Auth:            []ssh.AuthMethod{},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         30 * time.Second,
	}

	// 设置认证方式
	switch config.Type {
	case 1: // 密码认证
		if config.Password != "" {
			sshConfig.Auth = append(sshConfig.Auth, ssh.Password(config.Password))
		} else {
			return fmt.Errorf("密码认证但密码为空")
		}
	case 2: // 私钥认证
		if config.PublicKey != "" {
			signer, err := ssh.ParsePrivateKey([]byte(config.PublicKey))
			if err != nil {
				return fmt.Errorf("解析SSH密钥失败: %v", err)
			}
			sshConfig.Auth = append(sshConfig.Auth, ssh.PublicKeys(signer))
		} else {
			return fmt.Errorf("私钥认证但私钥为空")
		}
	case 3: // 公钥免认证 - 自动查找本地私钥
		sshUtil := util.NewSSHUtil()
		// 先尝试用户主目录
		if userKeyAuth, err := sshUtil.UserKeyAuth(); err == nil {
			sshConfig.Auth = append(sshConfig.Auth, userKeyAuth)
		} else {
			// 再尝试系统默认路径
			if defaultKeyAuth, err := sshUtil.DefaultKeyAuth(); err == nil {
				sshConfig.Auth = append(sshConfig.Auth, defaultKeyAuth)
			} else {
				return fmt.Errorf("type=3 公钥免认证失败: 本地未找到私钥文件")
			}
		}
	default:
		return fmt.Errorf("不支持的SSH认证类型: %d", config.Type)
	}

	// 创建SSH连接
	addr := fmt.Sprintf("%s:%d", config.IP, config.Port)
	client, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		return fmt.Errorf("创建SSH连接失败: %v", err)
	}
	defer client.Close()

	// 创建SCP会话
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("创建SCP会话失败: %v", err)
	}
	defer session.Close()

	// 打开本地文件
	file, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("打开本地文件失败: %v", err)
	}
	defer file.Close()

	// 获取文件信息
	fileStat, err := file.Stat()
	if err != nil {
		return fmt.Errorf("获取文件信息失败: %v", err)
	}

	// 准备SCP传输
	go func() {
		w, _ := session.StdinPipe()
		defer w.Close()

		fmt.Fprintf(w, "C%04o %d %s\n", fileStat.Mode().Perm(), fileStat.Size(), filepath.Base(remotePath))
		io.Copy(w, file)
		fmt.Fprint(w, "\x00")
	}()

	// 执行SCP命令
	if err := session.Run(fmt.Sprintf("scp -t %s", remotePath)); err != nil {
		// 回退到rsync
		log.Printf("SCP上传失败，尝试使用rsync: %v", err)
		rsyncCmd := fmt.Sprintf("rsync -avz -e 'ssh -p %d' %s %s@%s:%s",
			config.Port,
			localPath,
			config.Username,
			config.IP,
			remotePath)

		cmd := exec.Command("sh", "-c", rsyncCmd)
		if config.Type == 1 && config.Password != "" {
			cmd.Env = append(os.Environ(), fmt.Sprintf("SSH_PASSWORD=%s", config.Password))
		}
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("rsync传输失败: %v, 输出: %s", err, string(output))
		}
	}

	return nil
}

// updateDeployError 更新部署错误信息
func (s *AgentServiceImpl) updateDeployError(hostID uint, errorMsg string) {
	s.agentDao.UpdateByHostID(hostID, map[string]interface{}{
		"status":    agentModel.AgentStatusDeployFailed,
		"error_msg": errorMsg,
	})
}

// 新增辅助方法
func (s *AgentServiceImpl) detectPlatform(osInfo string) string {
	osInfo = strings.ToLower(osInfo)
	if strings.Contains(osInfo, "windows") {
		return "windows"
	} else if strings.Contains(osInfo, "darwin") || strings.Contains(osInfo, "macos") {
		return "darwin"
	}
	return "linux"
}

func (s *AgentServiceImpl) getInstallPath(osInfo string) string {
	if s.detectPlatform(osInfo) == "windows" {
		return "C:\\dodevops-agent"
	}
	return "/opt/agent"
}

func (s *AgentServiceImpl) createOrUpdateAgent(agent *agentModel.Agent) error {
	existing, err := s.agentDao.GetByHostID(agent.HostID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	if existing != nil {
		// 更新现有记录
		agent.ID = existing.ID
		agent.CreateTime = existing.CreateTime
		return s.agentDao.Update(agent)
	}

	// 创建新记录
	return s.agentDao.Create(agent)
}

func (s *AgentServiceImpl) performHealthCheck(hostID uint) {
	// 获取主机信息
	host, err := s.hostDao.GetCmdbHostById(hostID)
	if err != nil {
		return
	}

	// 获取SSH认证信息
	var sshConfig util.SSHConfig
	if host.SSHKeyID > 0 {
		sshKey, err := s.getSSHKeyByID(host.SSHKeyID)
		if err != nil {
			log.Printf("获取SSH密钥失败: %v", err)
			return
		}

		// 设置SSH连接基本信息
		sshConfig.IP = host.SSHIP
		sshConfig.Port = int(sshKey.Port)
		sshConfig.Username = sshKey.Username
		sshConfig.Type = sshKey.Type

		// 根据认证类型设置相应的认证信息
		switch sshKey.Type {
		case 1: // 密码认证
			sshConfig.Password = sshKey.Password
		case 2: // 私钥认证
			sshConfig.PublicKey = sshKey.PublicKey
		case 3: // 公钥免认证
			// 不需要设置额外信息，自动查找本地私钥
		default:
			log.Printf("不支持的SSH认证类型: %d", sshKey.Type)
			return
		}
	}

	// 检查agent进程是否运行
	targetOS := s.detectPlatform(host.OS)
	checkCmd := s.getCheckAgentCommand(targetOS)

	output, err := ExecuteSSHCommandWithOutput(sshConfig, checkCmd)
	if err != nil || (targetOS != "windows" && strings.TrimSpace(output) != "running") {
		// 进程不存在或状态不是运行中，更新状态为启动异常
		s.agentDao.UpdateByHostID(hostID, map[string]interface{}{
			"status": agentModel.AgentStatusStartError,
		})
	}
}
