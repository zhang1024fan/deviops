package model

import "time"

// TaskAnsibleView Ansible任务视图表
type TaskAnsibleView struct {
	ID        uint      `gorm:"primaryKey;comment:'主键ID'"`
	Name      string    `gorm:"size:100;not null;uniqueIndex;comment:'视图名称'"`
	CreatedAt time.Time `gorm:"autoCreateTime;comment:'创建时间'"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;comment:'更新时间'"`
}

func (TaskAnsibleView) TableName() string {
	return "task_ansible_view"
}