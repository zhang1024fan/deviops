package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// UserIDs 用户ID数组(关联sys_admin表)
type UserIDs []uint

// Value 实现driver.Valuer接口
func (u UserIDs) Value() (driver.Value, error) {
	if len(u) == 0 {
		return nil, nil
	}
	return json.Marshal(u)
}

// Scan 实现sql.Scanner接口
func (u *UserIDs) Scan(value interface{}) error {
	if value == nil {
		*u = UserIDs{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, u)
}

// ResourceIDs 资源ID数组
type ResourceIDs []uint

// Value 实现driver.Valuer接口
func (r ResourceIDs) Value() (driver.Value, error) {
	if len(r) == 0 {
		return nil, nil
	}
	return json.Marshal(r)
}

// Scan 实现sql.Scanner接口
func (r *ResourceIDs) Scan(value interface{}) error {
	if value == nil {
		*r = ResourceIDs{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, r)
}

// DomainsJSON 域名JSON类型
type DomainsJSON []string

// Value 实现driver.Valuer接口
func (d DomainsJSON) Value() (driver.Value, error) {
	if len(d) == 0 {
		return nil, nil
	}
	return json.Marshal(d)
}

// Scan 实现sql.Scanner接口
func (d *DomainsJSON) Scan(value interface{}) error {
	if value == nil {
		*d = DomainsJSON{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, d)
}

// OtherResources 其他资源配置
type OtherResources struct {
	RabbitMQ  []string `json:"rabbitmq,omitempty"`
	Zookeeper []string `json:"zookeeper,omitempty"`
	Kafka     []string `json:"kafka,omitempty"`
	Redis     []string `json:"redis,omitempty"`
	Other     []string `json:"other,omitempty"`
}

// Value 实现driver.Valuer接口
func (o OtherResources) Value() (driver.Value, error) {
	return json.Marshal(o)
}

// Scan 实现sql.Scanner接口
func (o *OtherResources) Scan(value interface{}) error {
	if value == nil {
		*o = OtherResources{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, o)
}

// Application 应用管理主表
type Application struct {
	ID          uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string         `gorm:"type:varchar(255);not null" json:"name"`                    // 应用名称
	Code        string         `gorm:"type:varchar(255);not null;uniqueIndex" json:"code"`       // 应用编码

	// 基本信息
	BusinessGroupID uint   `gorm:"not null" json:"business_group_id"`          // 业务组ID(关联cmdb_group)
	BusinessDeptID  uint   `gorm:"not null" json:"business_dept_id"`           // 业务部门ID(关联sys_dept)
	Description     string `gorm:"type:text" json:"description"`               // 应用介绍
	RepoURL         string `gorm:"type:varchar(500)" json:"repo_url"`          // 仓库地址

	// 负责人信息 (多个用户ID，关联sys_admin表)
	DevOwners   UserIDs `gorm:"type:json" json:"dev_owners"`   // 研发负责人
	TestOwners  UserIDs `gorm:"type:json" json:"test_owners"`  // 测试负责人
	OpsOwners   UserIDs `gorm:"type:json" json:"ops_owners"`   // 运维负责人

	// 技术信息
	ProgrammingLang string `gorm:"type:varchar(100)" json:"programming_lang"`  // 编程语言
	StartCommand    string `gorm:"type:text" json:"start_command"`             // 启动命令
	StopCommand     string `gorm:"type:text" json:"stop_command"`              // 停止命令
	HealthAPI       string `gorm:"type:varchar(500)" json:"health_api"`        // 健康检查接口

	// 关联资源 (存储资源ID)
	Domains     DomainsJSON    `gorm:"type:json" json:"domains"`     // 关联域名
	Hosts       ResourceIDs    `gorm:"type:json" json:"hosts"`       // 关联主机(cmdb_host表ID)
	Databases   ResourceIDs    `gorm:"type:json" json:"databases"`   // 关联数据库(cmdb_sql表ID)
	OtherRes    OtherResources `gorm:"type:json" json:"other_res"`   // 关联其他资源

	Status      int            `gorm:"type:tinyint;default:1" json:"status"`     // 状态
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`

	// 关联的Jenkins环境配置（级联删除）
	JenkinsEnvs []JenkinsEnv `gorm:"foreignKey:AppID;constraint:OnDelete:CASCADE" json:"jenkins_envs,omitempty"`
}


// JenkinsEnv Jenkins环境配置表
type JenkinsEnv struct {
	ID               uint               `gorm:"primaryKey;autoIncrement" json:"id"`
	AppID            uint               `gorm:"not null;index" json:"app_id"`                     // 应用ID
	EnvName          string             `gorm:"type:varchar(50);not null" json:"env_name"`        // 环境名称
	JenkinsServerID  *uint              `gorm:"default:null" json:"jenkins_server_id"`            // Jenkins服务器ID(关联account_auth)
	JobName          string             `gorm:"type:varchar(255);default:''" json:"job_name"`     // Jenkins任务名称

	CreatedAt        time.Time          `json:"created_at"`
	UpdatedAt        time.Time          `json:"updated_at"`

	// 关联
	Application   Application `gorm:"foreignKey:AppID" json:"application,omitempty"`
}

// TableName 设置表名
func (Application) TableName() string {
	return "app_application"
}

// TableName 设置表名
func (JenkinsEnv) TableName() string {
	return "app_jenkins_env"
}

// 请求响应结构体

// CreateApplicationRequest 创建应用请求
type CreateApplicationRequest struct {
	Name        string `json:"name" binding:"required"`        // 应用名称
	Code        string `json:"code"`                           // 应用编码(可选，不提供则根据名称自动生成)
	Description string `json:"description"`                    // 应用介绍
	RepoURL     string `json:"repo_url"`                       // 仓库地址

	BusinessGroupID uint `json:"business_group_id" binding:"required"` // 业务组ID
	BusinessDeptID  uint `json:"business_dept_id" binding:"required"`  // 业务部门ID

	// 负责人信息
	DevOwners  []uint `json:"dev_owners"`  // 研发负责人ID数组
	TestOwners []uint `json:"test_owners"` // 测试负责人ID数组
	OpsOwners  []uint `json:"ops_owners"`  // 运维负责人ID数组

	// 技术信息
	ProgrammingLang string `json:"programming_lang"` // 编程语言
	StartCommand    string `json:"start_command"`    // 启动命令
	StopCommand     string `json:"stop_command"`     // 停止命令
	HealthAPI       string `json:"health_api"`       // 健康检查接口

	// 关联资源
	Domains   []string       `json:"domains"`   // 关联域名
	Hosts     []uint         `json:"hosts"`     // 关联主机ID
	Databases []uint         `json:"databases"` // 关联数据库ID
	OtherRes  OtherResources `json:"other_res"` // 其他资源

	// Jenkins环境配置(可选，如果不提供则创建默认的3套环境：prod, test, dev)
	JenkinsEnvs []CreateJenkinsEnvRequest `json:"jenkins_envs,omitempty"` // Jenkins环境配置
}

// DeploymentRequest 部署请求
type DeploymentRequest struct {
	AppID       uint                   `json:"app_id" binding:"required"`       // 应用ID
	EnvName     string                 `json:"env_name" binding:"required"`     // 环境名称
	Branch      string                 `json:"branch"`                          // 分支名称(可选，默认master)
	Version     string                 `json:"version"`                         // 版本号(可选)
}

// DeploymentResponse 部署响应
type DeploymentResponse struct {
	AppID         uint   `json:"app_id"`         // 应用ID
	EnvName       string `json:"env_name"`       // 环境名称
	JenkinsJobURL string `json:"jenkins_job_url"` // Jenkins任务URL
	BuildNumber   int    `json:"build_number"`   // 构建编号
	Status        string `json:"status"`         // 部署状态
	Message       string `json:"message"`        // 状态消息
}

// GetDeploymentStatusRequest 获取部署状态请求
type GetDeploymentStatusRequest struct {
	AppID       uint   `form:"app_id" binding:"required"`       // 应用ID
	EnvName     string `form:"env_name" binding:"required"`     // 环境名称
	BuildNumber *int   `form:"build_number"`                    // 构建编号(可选)
}

// DeploymentHistoryRequest 部署历史请求
type DeploymentHistoryRequest struct {
	AppID    uint   `form:"app_id" binding:"required"` // 应用ID
	EnvName  string `form:"env_name"`                  // 环境名称(可选)
	Page     int    `form:"page" binding:"min=1"`      // 页码
	PageSize int    `form:"pageSize" binding:"min=1"`  // 每页数量
}

// DeploymentHistoryResponse 部署历史响应
type DeploymentHistoryResponse struct {
	Total int                 `json:"total"` // 总数
	List  []DeploymentHistory `json:"list"`  // 部署历史列表
}

// DeploymentHistory 部署历史记录
type DeploymentHistory struct {
	ID            uint                   `json:"id"`             // ID
	AppID         uint                   `json:"app_id"`         // 应用ID
	AppName       string                 `json:"app_name"`       // 应用名称
	EnvName       string                 `json:"env_name"`       // 环境名称
	Branch        string                 `json:"branch"`         // 分支
	Version       string                 `json:"version"`        // 版本
	BuildNumber   int                    `json:"build_number"`   // 构建编号
	Status        string                 `json:"status"`         // 状态：building/success/failed
	StartTime     time.Time              `json:"start_time"`     // 开始时间
	EndTime       *time.Time             `json:"end_time"`       // 结束时间
	Duration      int                    `json:"duration"`       // 耗时(秒)
	Operator      string                 `json:"operator"`       // 操作人
	JenkinsJobURL string                 `json:"jenkins_job_url"` // Jenkins任务URL
	LogURL        string                 `json:"log_url"`        // 日志URL
}

// UpdateApplicationRequest 更新应用请求
type UpdateApplicationRequest struct {
	Name        *string `json:"name"`        // 应用名称
	Description *string `json:"description"` // 应用介绍
	RepoURL     *string `json:"repo_url"`    // 仓库地址
	Status      *int    `json:"status"`      // 状态

	BusinessGroupID *uint `json:"business_group_id"` // 业务组ID
	BusinessDeptID  *uint `json:"business_dept_id"`  // 业务部门ID

	// 负责人信息
	DevOwners  *[]uint `json:"dev_owners"`  // 研发负责人ID数组
	TestOwners *[]uint `json:"test_owners"` // 测试负责人ID数组
	OpsOwners  *[]uint `json:"ops_owners"`  // 运维负责人ID数组

	// 技术信息
	ProgrammingLang *string `json:"programming_lang"` // 编程语言
	StartCommand    *string `json:"start_command"`    // 启动命令
	StopCommand     *string `json:"stop_command"`     // 停止命令
	HealthAPI       *string `json:"health_api"`       // 健康检查接口

	// 关联资源
	Domains   *[]string       `json:"domains"`   // 关联域名
	Hosts     *[]uint         `json:"hosts"`     // 关联主机ID
	Databases *[]uint         `json:"databases"` // 关联数据库ID
	OtherRes  *OtherResources `json:"other_res"` // 其他资源

	// Jenkins环境配置(可选，如果提供则完全替换现有配置)
	JenkinsEnvs *[]UpdateJenkinsEnvRequest `json:"jenkins_envs,omitempty"` // Jenkins环境配置
}

// ApplicationListRequest 应用列表请求
type ApplicationListRequest struct {
	Page            int    `form:"page" binding:"min=1"`     // 页码
	PageSize        int    `form:"pageSize" binding:"min=1"` // 每页数量
	Name            string `form:"name"`                     // 应用名称(模糊查询)
	Code            string `form:"code"`                     // 应用编码(模糊查询)
	BusinessGroupID *uint  `form:"business_group_id"`        // 业务组ID
	BusinessDeptID  *uint  `form:"business_dept_id"`         // 业务部门ID
	ProgrammingLang string `form:"programming_lang"`         // 编程语言
	Status          *int   `form:"status"`                   // 状态
}

// ApplicationListResponse 应用列表响应
type ApplicationListResponse struct {
	Total int           `json:"total"` // 总数
	List  []Application `json:"list"`  // 列表
}

// CreateJenkinsEnvRequest 创建Jenkins环境配置请求
type CreateJenkinsEnvRequest struct {
	AppID           uint   `json:"app_id"`                              // 应用ID(由控制器自动设置)
	EnvName         string `json:"env_name" binding:"required"`         // 环境名称
	JenkinsServerID *uint  `json:"jenkins_server_id"`                   // Jenkins服务器ID(关联account_auth表)
	JobName         string `json:"job_name"`                            // Jenkins任务名称
}

// UpdateJenkinsEnvRequest 更新Jenkins环境配置请求
type UpdateJenkinsEnvRequest struct {
	ID              *uint   `json:"id,omitempty"`     // 环境配置ID(可选，不提供则创建新的)
	EnvName         *string `json:"env_name"`         // 环境名称
	JenkinsServerID *uint   `json:"jenkins_server_id"` // Jenkins服务器ID(关联account_auth表)
	JobName         *string `json:"job_name"`         // Jenkins任务名称
}

// JenkinsEnvListRequest Jenkins环境配置列表请求
type JenkinsEnvListRequest struct {
	AppID   uint   `form:"app_id"`   // 应用ID
	EnvName string `form:"env_name"` // 环境名称
}

// JenkinsEnvListResponse Jenkins环境配置列表响应
type JenkinsEnvListResponse struct {
	Total int          `json:"total"` // 总数
	List  []JenkinsEnv `json:"list"`  // 列表
}

// JenkinsServerOption Jenkins服务器选项
type JenkinsServerOption struct {
	ID   uint   `json:"id"`   // 服务器ID
	Name string `json:"name"` // 服务器名称(别名)
}

// ValidateJenkinsJobRequest 验证Jenkins任务存在性请求
type ValidateJenkinsJobRequest struct {
	JenkinsServerID uint   `json:"jenkins_server_id" binding:"required"` // Jenkins服务器ID
	JobName         string `json:"job_name" binding:"required"`          // 任务名称
}

// ValidateJenkinsJobResponse 验证Jenkins任务存在性响应
type ValidateJenkinsJobResponse struct {
	Exists    bool   `json:"exists"`     // 任务是否存在
	JobName   string `json:"job_name"`   // 任务名称
	JobURL    string `json:"job_url"`    // 任务URL(如果存在)
	Message   string `json:"message"`    // 验证消息
	ServerID  uint   `json:"server_id"`  // 服务器ID
}

// 快速发布相关模型

// QuickDeployment 快速发布主表
type QuickDeployment struct {
	ID               uint                     `gorm:"primaryKey;autoIncrement" json:"id"`
	Title            string                   `gorm:"type:varchar(255);not null" json:"title"`                 // 发布标题
	BusinessGroupID  uint                     `gorm:"not null" json:"business_group_id"`                       // 业务组ID
	BusinessDeptID   uint                     `gorm:"not null" json:"business_dept_id"`                        // 业务部门ID
	Description      string                   `gorm:"type:text" json:"description"`                            // 发布描述
	Status           int                      `gorm:"type:tinyint;default:1" json:"status"`                    // 发布状态: 1=待发布 2=发布中 3=发布成功 4=发布失败 5=已取消
	TaskCount        int                      `gorm:"not null;default:0" json:"task_count"`                    // 任务数量，记录用户提交的发布任务数量
	ExecutionMode    int                      `gorm:"type:tinyint;default:1" json:"execution_mode"`            // 执行模式: 1=并行 2=串行
	CreatorID        uint                     `gorm:"not null" json:"creator_id"`                              // 创建人ID
	CreatorName      string                   `gorm:"type:varchar(100)" json:"creator_name"`                   // 创建人姓名
	StartTime        *time.Time               `json:"start_time"`                                               // 开始发布时间
	EndTime          *time.Time               `json:"end_time"`                                                 // 结束发布时间
	Duration         int                      `json:"duration"`                                                 // 发布耗时(秒)
	CreatedAt        time.Time                `json:"created_at"`
	UpdatedAt        time.Time                `json:"updated_at"`

	// 关联的发布任务
	Tasks []QuickDeploymentTask `gorm:"foreignKey:DeploymentID;constraint:OnDelete:CASCADE" json:"tasks,omitempty"`
}

// QuickDeploymentTask 快速发布任务表
type QuickDeploymentTask struct {
	ID            uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	DeploymentID  uint       `gorm:"not null;index" json:"deployment_id"`                         // 发布ID
	AppID         uint       `gorm:"not null" json:"app_id"`                                       // 应用ID
	AppName       string     `gorm:"type:varchar(255)" json:"app_name"`                           // 应用名称
	AppCode       string     `gorm:"type:varchar(255)" json:"app_code"`                           // 应用编码
	Environment   string     `gorm:"type:varchar(50)" json:"environment"`                         // 环境名称
	JenkinsEnvID  uint       `gorm:"not null" json:"jenkins_env_id"`                              // Jenkins环境配置ID
	JenkinsJobURL string     `gorm:"type:varchar(500)" json:"jenkins_job_url"`                    // Jenkins任务URL
	BuildNumber   int        `json:"build_number"`                                                 // 构建编号
	Status        int        `gorm:"type:tinyint;default:1" json:"status"`                        // 任务状态: 1=未部署 2=部署中 3=成功 4=异常
	ExecuteOrder  int        `gorm:"not null" json:"execute_order"`                               // 执行顺序
	StartTime     *time.Time `json:"start_time"`                                                   // 任务开始时间
	EndTime       *time.Time `json:"end_time"`                                                     // 任务结束时间
	Duration      int        `json:"duration"`                                                     // 任务耗时(秒)
	ErrorMessage  string     `gorm:"type:text" json:"error_message"`                              // 错误信息
	LogURL        string     `gorm:"type:varchar(500)" json:"log_url"`                            // 日志URL
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`

	// 关联
	Application Application `gorm:"foreignKey:AppID" json:"application,omitempty"`
	JenkinsEnv  JenkinsEnv  `gorm:"foreignKey:JenkinsEnvID" json:"jenkins_env,omitempty"`
}

// 表名设置
func (QuickDeployment) TableName() string {
	return "quick_deployments"
}

func (QuickDeploymentTask) TableName() string {
	return "quick_deployment_tasks"
}

// 快速发布请求响应模型

// CreateQuickDeploymentRequest 创建快速发布请求
type CreateQuickDeploymentRequest struct {
	Title           string                         `json:"title" binding:"required"`            // 发布标题
	BusinessGroupID uint                           `json:"business_group_id" binding:"required"` // 业务组ID
	BusinessDeptID  uint                           `json:"business_dept_id" binding:"required"`  // 业务部门ID
	Description     string                         `json:"description"`                          // 发布描述
	Applications    []QuickDeploymentAppRequest    `json:"applications" binding:"required,dive"` // 应用列表
}

// QuickDeploymentAppRequest 快速发布应用请求
type QuickDeploymentAppRequest struct {
	AppID       uint   `json:"app_id" binding:"required"`       // 应用ID（按数组顺序执行）
	Environment string `json:"environment" binding:"required"`  // 应用发布环境
}

// ExecuteQuickDeploymentRequest 执行快速发布请求
type ExecuteQuickDeploymentRequest struct {
	DeploymentID  uint `json:"deployment_id" binding:"required"`                       // 发布ID
	ExecutionMode *int `json:"execution_mode" binding:"omitempty,oneof=1 2" default:"1"` // 执行模式: 1=并行(默认) 2=串行
}

// QuickDeploymentListRequest 快速发布列表请求
type QuickDeploymentListRequest struct {
	Page            int    `form:"page" binding:"min=1"`     // 页码
	PageSize        int    `form:"pageSize" binding:"min=1"` // 每页数量
	BusinessGroupID *uint  `form:"business_group_id"`        // 业务组ID
	BusinessDeptID  *uint  `form:"business_dept_id"`         // 业务部门ID
	Environment     string `form:"environment"`              // 环境名称
	Status          *int   `form:"status"`                   // 状态
	CreatorID       *uint  `form:"creator_id"`               // 创建人ID
}

// QuickDeploymentListResponse 快速发布列表响应
type QuickDeploymentListResponse struct {
	Total int               `json:"total"` // 总数
	List  []QuickDeployment `json:"list"`  // 列表
}

// QuickDeploymentStatusInfo 快速发布状态信息
type QuickDeploymentStatusInfo struct {
	Status      int    `json:"status"`       // 状态码
	StatusText  string `json:"status_text"`  // 状态文本
	StatusColor string `json:"status_color"` // 状态颜色
}

// GetApplicationsForDeploymentRequest 获取可发布应用列表请求
type GetApplicationsForDeploymentRequest struct {
	BusinessGroupID uint   `form:"business_group_id" binding:"required"` // 业务组ID
	BusinessDeptID  uint   `form:"business_dept_id" binding:"required"`  // 业务部门ID
	Environment     string `form:"environment" binding:"required"`       // 环境名称
}

// ApplicationForDeployment 可发布应用信息
type ApplicationForDeployment struct {
	ID          uint   `json:"id"`           // 应用ID
	Name        string `json:"name"`         // 应用名称
	Code        string `json:"code"`         // 应用编码
	Environment string `json:"environment"`  // 环境名称
	JenkinsEnvID uint  `json:"jenkins_env_id"` // Jenkins环境配置ID
	JobName     string `json:"job_name"`     // Jenkins任务名称
	CanDeploy   bool   `json:"can_deploy"`   // 是否可以部署
	Reason      string `json:"reason"`       // 不能部署的原因
}

// 业务线服务树相关模型

// GetServiceTreeRequest 获取服务树请求
type GetServiceTreeRequest struct {
	BusinessGroupIDs []uint `form:"business_group_ids[]"` // 业务组ID列表，为空则查询所有
	Status           *int   `form:"status"`               // 应用状态筛选，为空则查询所有状态
	Environment      string `form:"environment"`          // 环境名称筛选，为空则不筛选环境配置
}

// GetAppEnvironmentRequest 获取应用环境配置请求
type GetAppEnvironmentRequest struct {
	AppID       uint   `form:"app_id" binding:"required"`       // 应用ID
	Environment string `form:"environment" binding:"required"`  // 环境名称
}

// AppEnvironmentResponse 应用环境配置响应
type AppEnvironmentResponse struct {
	AppID               uint   `json:"app_id"`               // 应用ID
	AppName             string `json:"app_name"`             // 应用名称
	AppCode             string `json:"app_code"`             // 应用编码
	Environment         string `json:"environment"`          // 环境名称
	IsConfigured        bool   `json:"is_configured"`        // 是否已配置
	JenkinsServerID     *uint  `json:"jenkins_server_id"`    // Jenkins服务器ID
	JenkinsServerName   string `json:"jenkins_server_name"`  // Jenkins服务器名称
	JobName             string `json:"job_name"`             // Jenkins任务名称
	JenkinsJobURL       string `json:"jenkins_job_url"`      // Jenkins任务URL
	Status              int    `json:"status"`               // 应用状态
	StatusText          string `json:"status_text"`          // 应用状态文本
	BusinessGroupID     uint   `json:"business_group_id"`    // 业务组ID
	BusinessDeptID      uint   `json:"business_dept_id"`     // 业务部门ID
	ProgrammingLang     string `json:"programming_lang"`     // 编程语言
}

// BusinessLineServiceTree 业务线服务树
type BusinessLineServiceTree struct {
	BusinessGroupID   uint                 `json:"business_group_id"`   // 业务组ID
	BusinessGroupName string               `json:"business_group_name"` // 业务组名称
	ServiceCount      int                  `json:"service_count"`       // 服务数量
	Services          []ServiceTreeNode    `json:"services"`            // 服务列表
}

// ServiceTreeNode 服务树节点
type ServiceTreeNode struct {
	ID                uint                    `json:"id"`                  // 应用ID
	Name              string                  `json:"name"`                // 应用名称
	Code              string                  `json:"code"`                // 应用编码
	Status            int                     `json:"status"`              // 应用状态
	StatusText        string                  `json:"status_text"`         // 状态文本
	ProgrammingLang   string                  `json:"programming_lang"`    // 编程语言
	BusinessDeptID    uint                    `json:"business_dept_id"`    // 业务部门ID
	BusinessDeptName  string                  `json:"business_dept_name"`  // 业务部门名称
	CreatedAt         string                  `json:"created_at"`          // 创建时间
	JenkinsEnvs       []ServiceJenkinsEnv     `json:"jenkins_envs"`        // Jenkins环境配置
}

// ServiceJenkinsEnv 服务Jenkins环境配置
type ServiceJenkinsEnv struct {
	ID              uint   `json:"id"`               // 环境配置ID
	EnvName         string `json:"env_name"`         // 环境名称
	JenkinsServerID *uint  `json:"jenkins_server_id"` // Jenkins服务器ID
	JobName         string `json:"job_name"`         // Jenkins任务名称
	IsConfigured    bool   `json:"is_configured"`    // 是否已配置完整
}

// TaskStatusResponse 任务状态响应
type TaskStatusResponse struct {
	TaskID       uint       `json:"task_id"`        // 任务ID
	Status       int        `json:"status"`         // 任务状态: 1=未部署 2=部署中 3=成功 4=异常
	StatusText   string     `json:"status_text"`    // 状态文本
	AppName      string     `json:"app_name"`       // 应用名称
	AppCode      string     `json:"app_code"`       // 应用编码
	Environment  string     `json:"environment"`    // 环境名称
	BuildNumber  int        `json:"build_number"`   // 构建编号
	StartTime    *time.Time `json:"start_time"`     // 开始时间
	EndTime      *time.Time `json:"end_time"`       // 结束时间
	Duration     int        `json:"duration"`       // 耗时(秒)
	ErrorMessage string     `json:"error_message"`  // 错误信息
	LogURL       string     `json:"log_url"`        // 日志URL
	Progress     int        `json:"progress"`       // 进度百分比(0-100)
}