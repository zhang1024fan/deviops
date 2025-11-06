package model

import (
	"time"
)

// TaskWork 任务作业模型
type TaskWork struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	TaskID    uint      `json:"task_id" gorm:"index;comment:关联的任务ID"`
	TemplateID uint     `json:"template_id" gorm:"index;comment:任务模板ID"`
	HostID    uint      `json:"host_id" gorm:"comment:执行主机ID"`
	Type      int       `json:"type" gorm:"comment:任务类型 1=普通任务,2=定时任务"`
	Status    int       `json:"status" gorm:"comment:任务状态 1=等待中,2=运行中,3=成功,4=异常"`
	Log       string    `json:"log" gorm:"type:text;comment:任务日志"`
	LogPath   string    `json:"log_path" gorm:"type:text;comment:日志文件路径"`
	StartTime *time.Time `json:"start_time" gorm:"comment:任务开始时间"`
	EndTime   *time.Time `json:"end_time" gorm:"comment:任务结束时间"`
	Duration  int       `json:"duration" gorm:"comment:执行耗时(秒)"`
	CreatedAt time.Time `json:"created_at" gorm:"comment:创建时间"`
	ScheduledTime *time.Time `json:"scheduled_time" gorm:"comment:定时任务执行时间"`
}



func (TaskWork) TableName() string {
	return "task_work"
}
// TaskHost 任务主机模型
type TaskHost struct {
	ID       uint   `gorm:"primarykey" json:"id"`
	TaskID   uint   `json:"task_id" gorm:"index;comment:关联的任务ID"`
	IP       string `json:"ip" gorm:"comment:主机IP"`
	Port     int    `json:"port" gorm:"comment:SSH端口"`
	User     string `json:"user" gorm:"comment:SSH用户名"`
	Password string `json:"password" gorm:"comment:SSH密码"`
	Key      string `json:"key" gorm:"type:text;comment:SSH密钥"`
}

// TaskJobStatus 任务作业状态
type TaskWorkStatus struct {
	Status    int    `json:"status"`
	Log       string `json:"log"`
	StartTime string `json:"start_time,omitempty"`
	EndTime   string `json:"end_time,omitempty"`
	Duration  int    `json:"duration,omitempty"`
}
