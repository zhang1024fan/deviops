package model

import (
	"dodevops-api/common/util"
)

// Agent Agent管理模型 (只支持Linux主机)
type Agent struct {
	ID              uint       `gorm:"column:id;comment:'主键';primaryKey;NOT NULL" json:"id"`
	HostID          uint       `gorm:"column:host_id;comment:'主机ID';NOT NULL;index" json:"hostId"`
	HostName        string     `gorm:"column:host_name;varchar(128);comment:'主机名称'" json:"hostName"`
	Version         string     `gorm:"column:version;varchar(32);comment:'Agent版本';default:'1.0.0'" json:"version"`
	Status          int        `gorm:"column:status;comment:'状态:1->部署中,2->部署失败,3->运行中,4->已停止'" json:"status"`
	InstallPath     string     `gorm:"column:install_path;varchar(256);comment:'安装路径'" json:"installPath"`
	Port            int        `gorm:"column:port;comment:'监听端口';default:9100" json:"port"`
	PID             int        `gorm:"column:pid;comment:'进程ID'" json:"pid"`
	LastHeartbeat   util.HTime `gorm:"column:last_heartbeat;comment:'最后心跳时间'" json:"lastHeartbeat"`
	UpdateTime      util.HTime `gorm:"column:update_time;comment:'更新时间'" json:"updateTime"`
	CreateTime      util.HTime `gorm:"column:create_time;comment:'创建时间';NOT NULL" json:"createTime"`
	ErrorMsg        string     `gorm:"column:error_msg;type:text;comment:'错误信息'" json:"errorMsg"`
	InstallProgress int        `gorm:"column:install_progress;comment:'安装进度(0-100)';default:0" json:"installProgress"` // 安装进度
}

func (Agent) TableName() string {
	return "monitor_agent"
}

// AgentStatus Agent状态常量
const (
	AgentStatusDeploying    = 1 // 部署中
	AgentStatusDeployFailed = 2 // 部署失败
	AgentStatusRunning      = 3 // 运行中
	AgentStatusStartError   = 4 // 启动异常
)

// AgentInstallProgress Agent安装进度常量
const (
	InstallProgressStart       = 0   // 开始
	InstallProgressCompiling   = 10  // 开始编译
	InstallProgressCompiled    = 30  // 编译完成
	InstallProgressTransfer    = 50  // 传输中
	InstallProgressTransferred = 70  // 传输完成
	InstallProgressConfigured  = 90  // 配置完成
	InstallProgressStarted     = 100 // 启动成功
)

// GetStatusText 获取状态文本
func (a *Agent) GetStatusText() string {
	statusMap := map[int]string{
		AgentStatusDeploying:    "部署中",
		AgentStatusDeployFailed: "部署失败",
		AgentStatusRunning:      "运行中",
		AgentStatusStartError:   "启动异常",
	}
	if text, ok := statusMap[a.Status]; ok {
		return text
	}
	return "未知状态"
}

// GetInstallProgressText 获取安装进度文本描述
func (a *Agent) GetInstallProgressText() string {
	progressMap := map[int]string{
		InstallProgressStart:       "开始部署",
		InstallProgressCompiling:   "开始编译",
		InstallProgressCompiled:    "编译完成",
		InstallProgressTransfer:    "传输中",
		InstallProgressTransferred: "传输完成",
		InstallProgressConfigured:  "配置完成",
		InstallProgressStarted:     "启动成功",
	}
	if text, ok := progressMap[a.InstallProgress]; ok {
		return text
	}
	return "未知进度"
}

// IsHealthy 判断Agent是否健康
func (a *Agent) IsHealthy() bool {
	return a.Status == AgentStatusRunning
}

// CreateAgentDto 创建Agent DTO (Linux主机专用)
type CreateAgentDto struct {
	HostID  uint   `validate:"required" json:"hostId"` // 主机ID
	Version string `json:"version"`                    // 版本
}

// BatchDeployAgentDto 批量部署Agent DTO
type BatchDeployAgentDto struct {
	HostIDs []uint `validate:"required" json:"hostIds"` // 主机ID列表
	Version string `json:"version"`                     // 版本
}

// UpdateAgentDto 更新Agent DTO
type UpdateAgentDto struct {
	ID      uint   `json:"id"`      // Agent ID
	Status  int    `json:"status"`  // 状态
	Version string `json:"version"` // 版本
}

// AgentListDto Agent列表查询DTO (Linux主机专用)
type AgentListDto struct {
	HostID   uint `json:"hostId"`   // 主机ID
	Status   int  `json:"status"`   // 状态
	Page     int  `json:"page"`     // 页码
	PageSize int  `json:"pageSize"` // 页大小
}

// AgentVO Agent视图对象 (Linux主机专用)
type AgentVO struct {
	ID                  uint       `json:"id"`
	HostID              uint       `json:"hostId"`
	HostName            string     `json:"hostName"`
	SSHIP               string     `json:"sshIp"`               // SSH连接IP
	Version             string     `json:"version"`
	Status              int        `json:"status"`
	StatusText          string     `json:"statusText"`
	InstallPath         string     `json:"installPath"`
	Port                int        `json:"port"`
	PID                 int        `json:"pid"`
	LastHeartbeat       util.HTime `json:"lastHeartbeat"`
	UpdateTime          util.HTime `json:"updateTime"`
	IsHealthy           bool       `json:"isHealthy"`
	CreateTime          util.HTime `json:"createTime"`
	ErrorMsg            string     `json:"errorMsg"`
	InstallProgress     int        `json:"installProgress"`     // 安装进度
	InstallProgressText string     `json:"installProgressText"` // 安装进度文本描述
}

// ToVO 转换为VO (Linux主机专用)
func (a *Agent) ToVO() *AgentVO {
	return &AgentVO{
		ID:                  a.ID,
		HostID:              a.HostID,
		HostName:            a.HostName,
		Version:             a.Version,
		Status:              a.Status,
		StatusText:          a.GetStatusText(),
		InstallPath:         a.InstallPath,
		Port:                a.Port,
		PID:                 a.PID,
		LastHeartbeat:       a.LastHeartbeat,
		UpdateTime:          a.UpdateTime,
		IsHealthy:           a.IsHealthy(),
		CreateTime:          a.CreateTime,
		ErrorMsg:            a.ErrorMsg,
		InstallProgress:     a.InstallProgress,
		InstallProgressText: a.GetInstallProgressText(),
	}
}

// AgentHeartbeatDto Agent心跳DTO (Linux主机专用)
type AgentHeartbeatDto struct {
	PID      int    `json:"pid"`      // Agent进程ID
	IP       string `json:"ip"`       // Agent所在主机的IP地址 (内网IP)
	Hostname string `json:"hostname"` // Agent所在主机的hostname
	Port     int    `json:"port"`     // Agent监听端口
	Token    string `json:"token"`    // 认证token
}
