package model

import "time"

// 任务状态常量
const (
	TaskStatusPending   = 1 // 等待中
	TaskStatusRunning   = 2 // 运行中
	TaskStatusSuccess   = 3 // 成功
	TaskStatusFailed    = 4 // 异常
	TaskStatusPaused    = 5 // 已暂停
)

// 任务类型常量
const (
	TaskTypeImmediate = 1 // 普通任务（立即执行）
	TaskTypeScheduled = 2 // 定时任务
	TaskTypeAnsible   = 3 // Ansible任务
)

type TaskIDRequest struct {
	ID uint `json:"id" binding:"required"`
}

// Task 任务模型
type Task struct {
	ID           uint      `gorm:"primarykey" json:"id"`
	Name         string    `json:"name" gorm:"size:255;comment:任务标题"`
	Type         int       `json:"type" gorm:"comment:任务类型 1=普通任务,2=定时任务,3=ansible任务"`
	Shell        string    `json:"shell" gorm:"type:text;comment:任务内容(任务模板ID,多个用逗号分隔)"`
	HostIDs      string    `json:"host_ids" gorm:"type:text;comment:主机ID(多个用逗号分隔)"`
	CronExpr     string    `json:"cron_expr" gorm:"size:255;comment:定时表达式(* * * * *)"`
	Tasklog      string    `json:"tasklog" gorm:"type:text;comment:任务执行日志"`
	Status       int       `json:"status" gorm:"comment:任务状态 1=等待中,2=运行中,3=成功,4=异常,5=已暂停"`
	Duration     int       `json:"duration" gorm:"comment:执行耗时(秒)"`
	TaskCount    int       `json:"task_count" gorm:"comment:任务数量"`
	ExecuteCount int       `json:"execute_count" gorm:"default:0;comment:执行次数"`
	NextRunTime  *time.Time `json:"next_run_time" gorm:"comment:下次执行时间" time_format:"2006-01-02 15:04:05"`
	Remark       string    `json:"remark" gorm:"type:text;comment:任务备注"`
	StartTime    *time.Time `json:"start_time" gorm:"comment:任务开始时间" time_format:"2006-01-02 15:04:05"`
	EndTime      *time.Time `json:"end_time" gorm:"comment:任务结束时间" time_format:"2006-01-02 15:04:05"`
	CreatedAt    time.Time  `json:"created_at" gorm:"comment:任务创建时间;autoCreateTime" time_format:"2006-01-02 15:04:05"`
}

func (Task) TableName() string {
	return "task_job"
}

// GetStatusName 获取状态名称
func GetStatusName(status int) string {
	switch status {
	case TaskStatusPending:
		return "等待中"
	case TaskStatusRunning:
		return "运行中"
	case TaskStatusSuccess:
		return "成功"
	case TaskStatusFailed:
		return "异常"
	case TaskStatusPaused:
		return "已暂停"
	default:
		return "未知状态"
	}
}

// GetTypeName 获取类型名称
func GetTypeName(taskType int) string {
	switch taskType {
	case TaskTypeImmediate:
		return "立即执行"
	case TaskTypeScheduled:
		return "定时任务"
	case TaskTypeAnsible:
		return "Ansible任务"
	default:
		return "未知类型"
	}
}

// CanPause 检查任务是否可以暂停
func (t *Task) CanPause() bool {
	// 只有定时任务且处于运行中状态时才可以暂停
	return t.Type == TaskTypeScheduled && t.Status == TaskStatusRunning
}

// CanResume 检查任务是否可以恢复
func (t *Task) CanResume() bool {
	// 只有定时任务且处于暂停状态时才可以恢复
	return t.Type == TaskTypeScheduled && t.Status == TaskStatusPaused
}

// CanStop 检查任务是否可以停止
func (t *Task) CanStop() bool {
	// 运行中或等待中的任务可以停止
	return t.Status == TaskStatusRunning || t.Status == TaskStatusPending
}

// IsScheduledTask 检查是否为定时任务
func (t *Task) IsScheduledTask() bool {
	return t.Type == TaskTypeScheduled
}

// IsActiveScheduledTask 检查是否为活跃的定时任务（运行中且为定时任务）
func (t *Task) IsActiveScheduledTask() bool {
	return t.Type == TaskTypeScheduled && t.Status == TaskStatusRunning
}

// TaskWithDetails 任务详细信息DTO（包含关联信息）
type TaskWithDetails struct {
	ID           uint      `json:"id"`
	Name         string    `json:"name"`
	Type         int       `json:"type"`
	Shell        string    `json:"shell"`
	HostIDs      string    `json:"host_ids"`
	CronExpr     string    `json:"cron_expr"`
	Tasklog      string    `json:"tasklog"`
	Status       int       `json:"status"`
	Duration     int       `json:"duration"`
	TaskCount    int       `json:"task_count"`
	ExecuteCount int       `json:"execute_count"`
	NextRunTime  *time.Time `json:"next_run_time"`
	Remark       string    `json:"remark"`
	StartTime    *time.Time `json:"start_time"`
	EndTime      *time.Time `json:"end_time"`
	CreatedAt    time.Time  `json:"created_at"`
	// 关联信息
	Templates []TaskTemplateInfo `json:"templates"`  // 关联的模板信息
	Hosts     []HostInfo         `json:"hosts"`      // 关联的主机信息
}

// TaskTemplateInfo 任务模板信息
type TaskTemplateInfo struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Type    int    `json:"type"`
	Content string `json:"content"`
}

// HostInfo 主机信息
type HostInfo struct {
	ID        uint   `json:"id"`
	HostName  string `json:"hostName"`
	PrivateIP string `json:"privateIp"`
	PublicIP  string `json:"publicIp"`
	SSHIP     string `json:"sshIp"`
	Status    int    `json:"status"`
}
