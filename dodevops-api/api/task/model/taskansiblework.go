package model

import "time"

// TaskAnsibleWork Ansible子任务表
type TaskAnsibleWork struct {
    ID           uint      `gorm:"primaryKey;comment:'主键ID'"`
    TaskID       uint      `gorm:"not null;index:idx_task_id;index:idx_task_work_composite;comment:'父任务ID'"`
    EntryFileName string   `gorm:"size:255;not null;comment:'入口文件名'"`
    EntryFilePath string   `gorm:"size:255;not null;comment:'入口文件路径'"`
    LogPath      string    `gorm:"size:255;comment:'日志路径'"`
    Status       int       `json:"status" gorm:"not null;default:1;index:idx_task_work_composite;comment:'子任务状态:1-等待中,2-运行中,3-成功,4-异常'"`
    StartTime    *time.Time `gorm:"comment:'开始时间'"`
    EndTime      *time.Time `gorm:"comment:'结束时间'"`
    Duration     int       `gorm:"comment:'执行耗时(秒)'"`
    ExitCode     int       `gorm:"comment:'退出代码'"`
    ErrorMsg     string    `gorm:"type:text;comment:'错误信息'"`
    Task         TaskAnsible `gorm:"foreignKey:TaskID;comment:'父任务'"`
}

// TaskAnsibleStatus 任务状态结构体
type TaskAnsibleStatus struct {
    Status    int       `json:"status"`
    StartTime time.Time `json:"start_time"`
    EndTime   time.Time `json:"end_time"`
    Log       string    `json:"log"`
}

func (TaskAnsibleWork) TableName() string {
	return "task_ansiblework"
}
