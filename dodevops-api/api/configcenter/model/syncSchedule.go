package model

import (
	"dodevops-api/common/util"
)

// SyncSchedule 定时同步配置表
type SyncSchedule struct {
	ID          uint        `gorm:"primaryKey" json:"id"`
	Name        string      `gorm:"type:varchar(100);not null" json:"name"`        // 配置名称
	CronExpr    string      `gorm:"type:varchar(100);not null" json:"cronExpr"`    // cron表达式
	KeyTypes    string      `gorm:"type:text;not null" json:"keyTypes"`            // 要同步的云厂商类型（JSON数组格式：[1,2,3]）
	Status      int         `gorm:"not null;default:1" json:"status"`              // 状态：1=启用，0=禁用
	LastRunTime *util.HTime `json:"lastRunTime"`                                   // 上次执行时间
	NextRunTime *util.HTime `json:"nextRunTime"`                                   // 下次执行时间
	SyncLog     string      `gorm:"type:text" json:"syncLog"`                      // 最近一次同步日志
	Remark      string      `gorm:"type:text" json:"remark"`                       // 备注信息
	CreatedAt   util.HTime  `json:"createdAt"`                                     // 创建时间
	UpdatedAt   util.HTime  `json:"updatedAt"`                                     // 更新时间
}

// TableName 表名
func (SyncSchedule) TableName() string {
	return "config_sync_schedule"
}

// CreateSyncScheduleDto 创建定时同步配置DTO
type CreateSyncScheduleDto struct {
	Name     string `json:"name" binding:"required"`     // 配置名称
	CronExpr string `json:"cronExpr" binding:"required"` // cron表达式
	KeyTypes string `json:"keyTypes" binding:"required"` // 要同步的云厂商类型（JSON数组格式：[1,2,3]）
	Status   int    `json:"status"`                      // 状态：1=启用，0=禁用
	Remark   string `json:"remark"`                      // 备注信息
}

// UpdateSyncScheduleDto 更新定时同步配置DTO
type UpdateSyncScheduleDto struct {
	ID       uint   `json:"id" binding:"required"`       // 配置ID
	Name     string `json:"name" binding:"required"`     // 配置名称
	CronExpr string `json:"cronExpr" binding:"required"` // cron表达式
	KeyTypes string `json:"keyTypes" binding:"required"` // 要同步的云厂商类型（JSON数组格式：[1,2,3]）
	Status   int    `json:"status"`                      // 状态：1=启用，0=禁用
	Remark   string `json:"remark"`                      // 备注信息
}

// SyncScheduleListDto 定时同步配置列表DTO（用于展示，包含下次执行时间）
type SyncScheduleListDto struct {
	ID          uint        `json:"id"`
	Name        string      `json:"name"`
	CronExpr    string      `json:"cronExpr"`
	KeyTypes    string      `json:"keyTypes"`
	Status      int         `json:"status"`
	LastRunTime *util.HTime `json:"lastRunTime"`
	NextRunTime *util.HTime `json:"nextRunTime"`
	SyncLog     string      `json:"syncLog"`
	Remark      string      `json:"remark"`
	CreatedAt   util.HTime  `json:"createdAt"`
	UpdatedAt   util.HTime  `json:"updatedAt"`
}