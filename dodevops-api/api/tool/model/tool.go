// 快捷导航工具模型
// author xiaoRui
package model

import (
	"dodevops-api/common/util"
	"time"
)

// Tool 快捷导航工具模型
type Tool struct {
	ID         uint       `gorm:"column:id;comment:'主键';primaryKey;autoIncrement;NOT NULL" json:"id"`
	Title      string     `gorm:"column:title;varchar(100);comment:'导航标题';NOT NULL" json:"title"`
	Icon       string     `gorm:"column:icon;varchar(500);comment:'导航图标'" json:"icon"`
	Link       string     `gorm:"column:link;varchar(500);comment:'链接地址';NOT NULL" json:"link"`
	Sort       int        `gorm:"column:sort;default:0;comment:'排序'" json:"sort"`
	CreateTime util.HTime `gorm:"column:create_time;comment:'创建时间';NOT NULL" json:"createTime"`
	UpdateTime time.Time  `gorm:"column:update_time;comment:'更新时间'" json:"updateTime"`
}

func (Tool) TableName() string {
	return "tool_link"
}

// AddToolDto 新增导航工具参数
type AddToolDto struct {
	Title string `json:"title" validate:"required,min=1,max=100"` // 导航标题
	Icon  string `json:"icon"`                                     // 导航图标
	Link  string `json:"link" validate:"required,url,max=500"`    // 链接地址
	Sort  int    `json:"sort"`                                     // 排序
}

// UpdateToolDto 更新导航工具参数
type UpdateToolDto struct {
	ID    uint   `json:"id" validate:"required"`                  // ID
	Title string `json:"title" validate:"required,min=1,max=100"` // 导航标题
	Icon  string `json:"icon"`                                     // 导航图标
	Link  string `json:"link" validate:"required,url,max=500"`    // 链接地址
	Sort  int    `json:"sort"`                                     // 排序
}

// ToolQueryDto 查询参数
type ToolQueryDto struct {
	Title    string `form:"title"`    // 标题（模糊查询）
	PageNum  int    `form:"pageNum"`  // 页码
	PageSize int    `form:"pageSize"` // 每页数量
}

// ========== 服务部署相关模型 ==========

// ServiceDeploy 服务部署记录
type ServiceDeploy struct {
	ID            uint       `gorm:"column:id;comment:'主键';primaryKey;autoIncrement;NOT NULL" json:"id"`
	ServiceName   string     `gorm:"column:service_name;varchar(64);comment:'服务名称';NOT NULL" json:"serviceName"`
	ServiceID     string     `gorm:"column:service_id;varchar(64);comment:'服务ID';NOT NULL" json:"serviceId"`
	Version       string     `gorm:"column:version;varchar(64);comment:'服务版本';NOT NULL" json:"version"`
	HostID        uint       `gorm:"column:host_id;comment:'主机ID';NOT NULL" json:"hostId"`
	HostIP        string     `gorm:"column:host_ip;varchar(64);comment:'主机IP';NOT NULL" json:"hostIp"`
	InstallDir    string     `gorm:"column:install_dir;varchar(255);comment:'安装目录';NOT NULL" json:"installDir"`
	ContainerName string     `gorm:"column:container_name;varchar(128);comment:'容器名称'" json:"containerName"`
	Ports         string     `gorm:"column:ports;varchar(255);comment:'端口映射(JSON)'" json:"ports"`
	EnvVars       string     `gorm:"column:env_vars;text;comment:'环境变量(JSON)'" json:"envVars"`
	Status        int        `gorm:"column:status;default:0;comment:'状态:0->部署中,1->运行中,2->已停止,3->部署失败'" json:"status"`
	DeployLog     string     `gorm:"column:deploy_log;text;comment:'部署日志'" json:"deployLog"`
	CreateTime    util.HTime `gorm:"column:create_time;comment:'创建时间';NOT NULL" json:"createTime"`
	UpdateTime    time.Time  `gorm:"column:update_time;comment:'更新时间'" json:"updateTime"`
}

func (ServiceDeploy) TableName() string {
	return "tool_service_deploy"
}

// CreateDeployDto 创建部署任务DTO
type CreateDeployDto struct {
	ServiceID   string                 `json:"serviceId" validate:"required"`   // 服务ID (如: mysql)
	Version     string                 `json:"version" validate:"required"`     // 版本 (如: 5.7)
	HostID      uint                   `json:"hostId" validate:"required"`      // 主机ID
	InstallDir  string                 `json:"installDir" validate:"required"`  // 安装目录
	EnvVars     map[string]interface{} `json:"envVars"`                         // 环境变量
	AutoStart   bool                   `json:"autoStart"`                       // 是否自动启动
}

// DeployQueryDto 部署记录查询DTO
type DeployQueryDto struct {
	ServiceName string `form:"serviceName"` // 服务名称
	HostID      uint   `form:"hostId"`      // 主机ID
	Status      *int   `form:"status"`      // 状态: 不传=全部, 0=部署中, 1=运行中, 2=已停止, 3=部署失败
	PageNum     int    `form:"pageNum"`     // 页码
	PageSize    int    `form:"pageSize"`    // 每页数量
}

// ServiceDeployVo 部署记录VO
type ServiceDeployVo struct {
	ID            uint                   `json:"id"`
	ServiceName   string                 `json:"serviceName"`
	ServiceID     string                 `json:"serviceId"`
	Version       string                 `json:"version"`
	HostID        uint                   `json:"hostId"`
	HostIP        string                 `json:"hostIp"`
	HostName      string                 `json:"hostName"`
	InstallDir    string                 `json:"installDir"`
	ContainerName string                 `json:"containerName"`
	Ports         string                 `json:"ports"`
	EnvVars       map[string]interface{} `json:"envVars"`
	Status        int                    `json:"status"`
	StatusText    string                 `json:"statusText"`
	DeployLog     string                 `json:"deployLog"`
	CreateTime    util.HTime             `json:"createTime"`
	UpdateTime    time.Time              `json:"updateTime"`
}

// ServiceInfo 服务信息（来自services.json）
type ServiceInfo struct {
	ID          string              `json:"id"`
	Name        string              `json:"name"`
	Category    string              `json:"category"`
	Description string              `json:"description"`
	Icon        string              `json:"icon"`
	Versions    []ServiceVersion    `json:"versions"`
	DefaultPort int                 `json:"default_port"`
	EnvVars     []ServiceEnvVar     `json:"env_vars"`
	MinMemory   string              `json:"min_memory,omitempty"`
	MinCPU      string              `json:"min_cpu,omitempty"`
}

// ServiceVersion 服务版本信息
type ServiceVersion struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	File        string            `json:"file"`
	Stable      bool              `json:"stable"`
	Recommended bool              `json:"recommended"`
	DeployType  string            `json:"deploy_type,omitempty"`  // 部署类型: container(容器) 或 binary(二进制)
	ExtractPaths map[string]string `json:"extract_paths,omitempty"` // 需要从镜像提取的路径映射
}

// ServiceEnvVar 服务环境变量
type ServiceEnvVar struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Required    bool   `json:"required"`
	Default     string `json:"default"`
}

// ServiceCategory 服务分类
type ServiceCategory struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Icon string `json:"icon"`
}

// ServicesConfig 服务配置（services.json整体结构）
type ServicesConfig struct {
	Services   []ServiceInfo     `json:"services"`
	Categories []ServiceCategory `json:"categories"`
}
