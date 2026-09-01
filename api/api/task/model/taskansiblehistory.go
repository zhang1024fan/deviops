package model

import "time"

// TaskAnsibleHistory 任务执行历史记录主表
type TaskAnsibleHistory struct {
	ID            uint   `gorm:"primaryKey;comment:'主键ID'"`
	TaskID        uint   `gorm:"not null;index:idx_history_task_id;comment:'关联的任务ID'"`
	UniqId        string `gorm:"size:50;not null;comment:'任务唯一标识(每次执行生成)'"`
	Status        int    `json:"status" gorm:"not null;default:1;comment:'执行状态:1-等待中,2-运行中,3-成功,4-异常'"`
	ErrorMsg      string `gorm:"type:text;comment:'错误信息'"`
	TotalDuration int    `gorm:"not null;default:0;comment:'任务执行总耗时(秒)'"`
	Trigger       int    `gorm:"not null;default:1;comment:'触发方式:1-手动,2-定时,3-API'"`
	OperatorID    uint   `gorm:"comment:'操作人ID'"`
	OperatorName  string `gorm:"size:50;comment:'操作人姓名'"`
	StartedAt     *time.Time
	FinishedAt    *time.Time
	CreatedAt     time.Time

	TaskAnsible   *TaskAnsible             `gorm:"foreignKey:TaskID"`
	WorkHistories []TaskAnsibleworkHistory `gorm:"foreignKey:HistoryID"`
}

func (TaskAnsibleHistory) TableName() string {
	return "task_ansible_history"
}

// TaskAnsibleworkHistory 任务执行历史记录详情表（对应每个host的执行结果）
type TaskAnsibleworkHistory struct {
	ID        uint   `gorm:"primaryKey;comment:'主键ID'"`
	HistoryID uint   `gorm:"not null;index:idx_work_history_id;comment:'关联的历史记录ID'"`
	TaskID    uint   `gorm:"not null;comment:'关联的任务ID'"` // 为了方便查询保留
	WorkID    uint   `gorm:"comment:'关联的WorkID(如果有)'"`
	HostName  string `gorm:"size:255;not null;comment:'主机名/IP'"`
	Status    int    `gorm:"not null;default:1;comment:'状态:1-等待,2-执行中,3-成功,4-失败,5-跳过'"`
	LogPath   string `gorm:"size:255;comment:'日志文件路径'"`
	Duration  int    `gorm:"not null;default:0;comment:'耗时(秒)'"`
	CreatedAt time.Time
}

func (TaskAnsibleworkHistory) TableName() string {
	return "task_ansiblework_history"
}