package model

import ()

// Jenkins账号类型常量 (对应AccountAuth.Type字段)
const (
	JenkinsAccountType = 4 // Jenkins账号类型
)

// JenkinsJob Jenkins任务信息
type JenkinsJob struct {
	Name         string            `json:"name"`         // 任务名称
	DisplayName  string            `json:"displayName"`  // 显示名称
	Description  string            `json:"description"`  // 任务描述
	URL          string            `json:"url"`          // 任务URL
	Buildable    bool              `json:"buildable"`    // 是否可构建
	Color        string            `json:"color"`        // 状态颜色(blue/red/yellow等)
	Class        string            `json:"_class"`       // 任务类型
	LastBuild    *JenkinsBuild     `json:"lastBuild"`    // 最后一次构建
	LastStableBuild *JenkinsBuild  `json:"lastStableBuild"` // 最后一次稳定构建
	LastSuccessfulBuild *JenkinsBuild `json:"lastSuccessfulBuild"` // 最后一次成功构建
	LastFailedBuild *JenkinsBuild  `json:"lastFailedBuild"` // 最后一次失败构建
	Properties   []JobProperty     `json:"property"`     // 任务属性
	Actions      []JobAction       `json:"actions"`      // 任务操作
}

// JenkinsBuild Jenkins构建信息
type JenkinsBuild struct {
	Number      int       `json:"number"`      // 构建编号
	URL         string    `json:"url"`         // 构建URL
	DisplayName string    `json:"displayName"` // 显示名称
	FullDisplayName string `json:"fullDisplayName"` // 完整显示名称
	Description string    `json:"description"` // 构建描述
	Result      string    `json:"result"`      // 构建结果 SUCCESS/FAILURE/UNSTABLE/ABORTED
	Building    bool      `json:"building"`    // 是否正在构建
	Duration    int64     `json:"duration"`    // 构建时长(毫秒)
	EstimatedDuration int64 `json:"estimatedDuration"` // 预计时长(毫秒)
	Timestamp   int64     `json:"timestamp"`   // 开始时间戳
	KeepLog     bool      `json:"keepLog"`     // 是否保留日志
	QueueId     int       `json:"queueId"`     // 队列ID
	Executor    *BuildExecutor `json:"executor"` // 执行器信息
	Actions     []BuildAction  `json:"actions"`  // 构建操作
	ChangeSet   *ChangeSet     `json:"changeSet"` // 变更集
	Culprits    []User         `json:"culprits"`  // 责任人
}

// BuildExecutor 构建执行器信息
type BuildExecutor struct {
	Number int    `json:"number"` // 执行器编号
	Node   string `json:"node"`   // 节点名称
}

// BuildAction 构建操作
type BuildAction struct {
	Class string `json:"_class"` // 操作类型
}

// ChangeSet 变更集
type ChangeSet struct {
	Items []ChangeSetItem `json:"items"` // 变更项目
	Kind  string          `json:"kind"`  // 变更类型
}

// ChangeSetItem 变更项目
type ChangeSetItem struct {
	AffectedPaths []string `json:"affectedPaths"` // 影响路径
	Author        User     `json:"author"`        // 作者
	Comment       string   `json:"comment"`       // 提交注释
	Date          string   `json:"date"`          // 提交日期
	Id            string   `json:"id"`            // 提交ID
	Msg           string   `json:"msg"`           // 提交消息
	Timestamp     int64    `json:"timestamp"`     // 时间戳
}

// User 用户信息
type User struct {
	AbsoluteUrl string `json:"absoluteUrl"` // 绝对URL
	FullName    string `json:"fullName"`    // 全名
}

// JobProperty 任务属性
type JobProperty struct {
	Class string `json:"_class"` // 属性类型
}

// JobAction 任务操作
type JobAction struct {
	Class string `json:"_class"` // 操作类型
}

// BuildStatus 构建状态枚举
const (
	BuildStatusSuccess   = "SUCCESS"   // 成功
	BuildStatusFailure   = "FAILURE"   // 失败
	BuildStatusUnstable  = "UNSTABLE"  // 不稳定
	BuildStatusAborted   = "ABORTED"   // 已终止
	BuildStatusInProgress = "IN_PROGRESS" // 进行中
)

// JobColor 任务状态颜色枚举
const (
	JobColorBlue       = "blue"        // 成功
	JobColorRed        = "red"         // 失败
	JobColorYellow     = "yellow"      // 不稳定
	JobColorGrey       = "grey"        // 未构建
	JobColorDisabled   = "disabled"    // 已禁用
	JobColorAborted    = "aborted"     // 已终止
	JobColorNotBuilt   = "notbuilt"    // 未构建
	JobColorBlueAnime  = "blue_anime"  // 构建中(成功)
	JobColorRedAnime   = "red_anime"   // 构建中(失败)
	JobColorYellowAnime = "yellow_anime" // 构建中(不稳定)
	JobColorGreyAnime  = "grey_anime"  // 构建中(未知)
)

// ================== 请求和响应结构体 ==================

// JenkinsServerInfo Jenkins服务器信息 (基于AccountAuth)
type JenkinsServerInfo struct {
	ID          uint   `json:"id"`          // 账号ID
	Alias       string `json:"alias"`       // 别名(服务器名称)
	Host        string `json:"host"`        // Jenkins服务器地址
	Port        int    `json:"port"`        // 端口
	Username    string `json:"username"`    // 用户名
	Description string `json:"description"` // 描述(备注)
	CreatedAt   string `json:"createdAt"`   // 创建时间
	UpdatedAt   string `json:"updatedAt"`   // 更新时间
}

// JenkinsServerListResponse Jenkins服务器列表响应
type JenkinsServerListResponse struct {
	List  []JenkinsServerInfo `json:"list"`
	Total int64               `json:"total"`
}

// JenkinsJobListResponse Jenkins任务列表响应
type JenkinsJobListResponse struct {
	Jobs    []JenkinsJob `json:"jobs"`
	Total   int          `json:"total"`
	Server  string       `json:"server"`  // 服务器名称
}

// JenkinsJobDetailResponse Jenkins任务详情响应
type JenkinsJobDetailResponse struct {
	Job    JenkinsJob      `json:"job"`
	Builds []JenkinsBuild  `json:"builds"`   // 构建历史
	Server string          `json:"server"`   // 服务器名称
}

// JenkinsBuildDetailResponse Jenkins构建详情响应
type JenkinsBuildDetailResponse struct {
	Build  JenkinsBuild `json:"build"`
	Log    string       `json:"log"`    // 构建日志
	Server string       `json:"server"` // 服务器名称
}

// StartJobRequest 启动任务请求
type StartJobRequest struct {
	Parameters map[string]string `json:"parameters"` // 构建参数
	Reason     string            `json:"reason"`     // 构建原因
}

// StartJobResponse 启动任务响应
type StartJobResponse struct {
	Success     bool   `json:"success"`     // 是否启动成功
	Message     string `json:"message"`     // 响应消息
	QueueId     int    `json:"queueId"`     // 队列ID
	JobName     string `json:"jobName"`     // 任务名称
	Server      string `json:"server"`      // 服务器名称
	BuildNumber int    `json:"buildNumber"` // 构建编号(如果已知)
}

// StopBuildRequest 停止构建请求
type StopBuildRequest struct {
	Reason string `json:"reason"` // 停止原因
}

// StopBuildResponse 停止构建响应
type StopBuildResponse struct {
	Success bool   `json:"success"` // 是否停止成功
	Message string `json:"message"` // 响应消息
	JobName string `json:"jobName"` // 任务名称
	BuildNumber int `json:"buildNumber"` // 构建编号
	Server  string `json:"server"`  // 服务器名称
}

// GetBuildLogRequest 获取构建日志请求
type GetBuildLogRequest struct {
	Start int `json:"start"` // 开始位置
	Html  bool `json:"html"`  // 是否返回HTML格式
}

// GetBuildLogResponse 获取构建日志响应
type GetBuildLogResponse struct {
	Log        string `json:"log"`        // 日志内容
	HasMore    bool   `json:"hasMore"`    // 是否有更多日志
	TextSize   int    `json:"textSize"`   // 文本大小
	MoreData   bool   `json:"moreData"`   // 是否有更多数据
	JobName    string `json:"jobName"`    // 任务名称
	BuildNumber int   `json:"buildNumber"` // 构建编号
	Server     string `json:"server"`     // 服务器名称
}

// JenkinsSystemInfo Jenkins系统信息
type JenkinsSystemInfo struct {
	Version          string            `json:"version"`          // Jenkins版本
	Mode             string            `json:"mode"`             // 运行模式
	NodeDescription  string            `json:"nodeDescription"`  // 节点描述
	NodeName         string            `json:"nodeName"`         // 节点名称
	NumExecutors     int               `json:"numExecutors"`     // 执行器数量
	UseCrumbs        bool              `json:"useCrumbs"`        // 是否使用CSRF保护
	UseSecurity      bool              `json:"useSecurity"`      // 是否使用安全
	Views            []JenkinsView     `json:"views"`            // 视图列表
	PrimaryView      *JenkinsView      `json:"primaryView"`      // 主视图
	UnlabeledLoad    map[string]int    `json:"unlabeledLoad"`    // 未标记负载
	AssignedLabels   []JenkinsLabel    `json:"assignedLabels"`   // 分配的标签
	OverallLoad      map[string]int    `json:"overallLoad"`      // 总体负载
	Computers        []JenkinsComputer `json:"computers"`        // 计算机列表
}

// JenkinsView Jenkins视图
type JenkinsView struct {
	Name        string       `json:"name"`        // 视图名称
	URL         string       `json:"url"`         // 视图URL
	Description string       `json:"description"` // 视图描述
	Jobs        []JenkinsJob `json:"jobs"`        // 视图中的任务
}

// JenkinsLabel Jenkins标签
type JenkinsLabel struct {
	Name string `json:"name"` // 标签名称
}

// JenkinsComputer Jenkins计算机(节点)
type JenkinsComputer struct {
	DisplayName     string                  `json:"displayName"`     // 显示名称
	Executors       []JenkinsExecutor       `json:"executors"`       // 执行器列表
	Icon            string                  `json:"icon"`            // 图标
	IconClassName   string                  `json:"iconClassName"`   // 图标类名
	Idle            bool                    `json:"idle"`            // 是否空闲
	JnlpAgent       bool                    `json:"jnlpAgent"`       // 是否JNLP代理
	LaunchSupported bool                    `json:"launchSupported"` // 是否支持启动
	LoadStatistics  JenkinsLoadStatistics   `json:"loadStatistics"`  // 负载统计
	ManualLaunchAllowed bool                `json:"manualLaunchAllowed"` // 是否允许手动启动
	MonitorData     map[string]interface{}  `json:"monitorData"`     // 监控数据
	NumExecutors    int                     `json:"numExecutors"`    // 执行器数量
	Offline         bool                    `json:"offline"`         // 是否离线
	OfflineCause    interface{}             `json:"offlineCause"`    // 离线原因
	OneOffExecutors []JenkinsExecutor       `json:"oneOffExecutors"` // 一次性执行器
	TemporarilyOffline bool                 `json:"temporarilyOffline"` // 是否临时离线
}

// JenkinsExecutor Jenkins执行器
type JenkinsExecutor struct {
	CurrentExecutable interface{} `json:"currentExecutable"` // 当前执行的任务
	CurrentWorkUnit   interface{} `json:"currentWorkUnit"`   // 当前工作单元
	Idle              bool        `json:"idle"`              // 是否空闲
	LikelyStuck       bool        `json:"likelyStuck"`       // 是否可能卡住
	Number            int         `json:"number"`            // 执行器编号
	Progress          int         `json:"progress"`          // 进度
}

// JenkinsLoadStatistics Jenkins负载统计
type JenkinsLoadStatistics struct {
	BusyExecutors   int `json:"busyExecutors"`   // 忙碌执行器数
	IdleExecutors   int `json:"idleExecutors"`   // 空闲执行器数
	TotalExecutors  int `json:"totalExecutors"`  // 总执行器数
	QueueLength     int `json:"queueLength"`     // 队列长度
}

// JenkinsQueue Jenkins队列信息
type JenkinsQueue struct {
	Items []JenkinsQueueItem `json:"items"` // 队列项目
}

// JenkinsQueueItem Jenkins队列项目
type JenkinsQueueItem struct {
	Actions            []interface{} `json:"actions"`            // 操作
	Blocked            bool          `json:"blocked"`            // 是否阻塞
	Buildable          bool          `json:"buildable"`          // 是否可构建
	Id                 int           `json:"id"`                 // 队列项目ID
	InQueueSince       int64         `json:"inQueueSince"`       // 入队时间
	Params             string        `json:"params"`             // 参数
	Stuck              bool          `json:"stuck"`              // 是否卡住
	Task               JenkinsTask   `json:"task"`               // 任务信息
	URL                string        `json:"url"`                // URL
	Why                string        `json:"why"`                // 等待原因
	BuildableStartMilliseconds int64 `json:"buildableStartMilliseconds"` // 可构建开始时间
}

// JenkinsTask Jenkins任务
type JenkinsTask struct {
	Name  string `json:"name"`  // 任务名称
	URL   string `json:"url"`   // 任务URL
	Color string `json:"color"` // 任务颜色
}

// TestJenkinsConnectionRequest 测试Jenkins连接请求
type TestJenkinsConnectionRequest struct {
	URL      string `json:"url" binding:"required"`      // Jenkins服务器地址
	Username string `json:"username" binding:"required"` // 用户名
	Password string `json:"password" binding:"required"` // 密码或API Token
}

// TestJenkinsConnectionResponse 测试Jenkins连接响应
type TestJenkinsConnectionResponse struct {
	Success     bool              `json:"success"`     // 是否连接成功
	Message     string            `json:"message"`     // 响应消息
	SystemInfo  *JenkinsSystemInfo `json:"systemInfo"`  // 系统信息
	Error       string            `json:"error"`       // 错误信息
}

// BuildParameter 构建参数
type BuildParameter struct {
	Name         string      `json:"name"`         // 参数名称
	Type         string      `json:"type"`         // 参数类型
	DefaultValue interface{} `json:"defaultValue"` // 默认值
	Description  string      `json:"description"`  // 参数描述
	Choices      []string    `json:"choices"`      // 选择项(对于choice参数)
}

// JobConfigResponse 任务配置响应
type JobConfigResponse struct {
	JobName    string           `json:"jobName"`    // 任务名称
	Config     string           `json:"config"`     // 配置XML
	Parameters []BuildParameter `json:"parameters"` // 构建参数
	Server     string           `json:"server"`     // 服务器名称
}

// UpdateJobConfigRequest 更新任务配置请求
type UpdateJobConfigRequest struct {
	Config string `json:"config" binding:"required"` // 配置XML
}

// CreateJobRequest 创建任务请求
type CreateJobRequest struct {
	Name   string `json:"name" binding:"required"`   // 任务名称
	Config string `json:"config" binding:"required"` // 配置XML
}

// DeleteJobRequest 删除任务请求
type DeleteJobRequest struct {
	Confirm bool `json:"confirm"` // 确认删除
}

// JobStatistics 任务统计信息
type JobStatistics struct {
	TotalJobs       int `json:"totalJobs"`       // 总任务数
	ActiveJobs      int `json:"activeJobs"`      // 活跃任务数
	DisabledJobs    int `json:"disabledJobs"`    // 禁用任务数
	SuccessfulBuilds int `json:"successfulBuilds"` // 成功构建数
	FailedBuilds    int `json:"failedBuilds"`    // 失败构建数
	RunningBuilds   int `json:"runningBuilds"`   // 正在运行的构建数
}

// SystemStatistics 系统统计信息
type SystemStatistics struct {
	ServerCount     int           `json:"serverCount"`     // 服务器数量
	JobStatistics   JobStatistics `json:"jobStatistics"`   // 任务统计
	QueueLength     int           `json:"queueLength"`     // 队列长度
	BusyExecutors   int           `json:"busyExecutors"`   // 忙碌执行器数
	TotalExecutors  int           `json:"totalExecutors"`  // 总执行器数
}