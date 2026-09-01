package model

import "time"

// TaskAnsible Ansible任务主表
type TaskAnsible struct {
	ID                 uint              `gorm:"primaryKey;comment:'主键ID'"`
	Name               string            `gorm:"size:100;not null;uniqueIndex;comment:'任务名称'"`
	Description        string            `gorm:"type:text;comment:'任务描述'"`
	Type               int               `gorm:"not null;default:1;comment:'任务类型:1-手动,2-Git,3-K8s'"`
	GitRepo            string            `gorm:"size:255;comment:'Git仓库地址'"`
	HostGroups         string            `gorm:"type:text;not null;comment:'主机分组JSON'"`
	AllHostIDs         string            `gorm:"type:text;not null;comment:'所有主机ID JSON数组'"`
	GlobalVars         string            `gorm:"type:text;comment:'全局变量JSON'"`
	ExtraVars          string            `gorm:"type:text;comment:'额外参数YAML/JSON'"`
	CliArgs            string            `gorm:"type:text;comment:'cli命令行参数'"`
	Status             int               `json:"status" gorm:"not null;default:1;index:idx_task_status;comment:'任务状态:1-等待中,2-运行中,3-成功,4-异常'"`
	ErrorMsg           string            `gorm:"type:text;comment:'错误信息'"`
	TaskCount          int               `gorm:"not null;default:0;comment:'任务数量(Type=1时为上传文件数,Type=2时为解析的playbook数,Type=3时固定为1)'"`
	TotalDuration      int               `gorm:"not null;default:0;comment:'任务执行总耗时(秒,所有子任务耗时总和)'"`
	UseConfig          int               `gorm:"not null;default:0;comment:'是否使用配置管理中的参数 0-不使用,1-使用'"`
	InventoryConfigID  *uint             `gorm:"comment:'选用的inventory配置ID'"`
	GlobalVarsConfigID *uint             `gorm:"comment:'选用的global_vars配置ID'"`
	ExtraVarsConfigID  *uint             `gorm:"comment:'选用的extra_vars配置ID'"`
	CliArgsConfigID    *uint             `gorm:"comment:'选用的cli_args配置ID'"`
	MaxHistoryKeep     int               `gorm:"default:3;comment:'最大保留历史记录数'"`
	CreatedAt          time.Time         `gorm:"not null;comment:'创建时间'"`
	UpdatedAt          time.Time         `gorm:"not null;comment:'更新时间'"`
	Works              []TaskAnsibleWork `gorm:"foreignKey:TaskID;comment:'子任务列表'"`
	CronExpr    	   string           `gorm:"size:64;comment:'定时表达式'"`
	IsRecurring        int              `gorm:"not null;default:0;comment:'是否周期性任务:0-否,1-是'"`
	ViewID             *uint            `gorm:"comment:'视图ID'"`
	View               *TaskAnsibleView `gorm:"foreignKey:ViewID"`
	InventoryConfig  *ConfigAnsible `gorm:"foreignKey:InventoryConfigID"`
	GlobalVarsConfig *ConfigAnsible `gorm:"foreignKey:GlobalVarsConfigID"`
	ExtraVarsConfig  *ConfigAnsible `gorm:"foreignKey:ExtraVarsConfigID"`
	CliArgsConfig    *ConfigAnsible `gorm:"foreignKey:CliArgsConfigID"`
}

func (TaskAnsible) TableName() string {
	return "task_ansible"
}