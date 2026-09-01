package model

import "time"

// ConfigAnsible Ansible配置中心
type ConfigAnsible struct {
	ID        uint      `gorm:"primaryKey;autoIncrement;comment:'主键ID'"`
	Name      string    `gorm:"size:100;not null;uniqueIndex:uk_config_ansible_name;comment:'配置名称'"`
	Type      int       `gorm:"not null;index:idx_config_ansible_type;comment:'1-inventory 2-global_vars 3-extra_vars 4-cli_args'"`
	Content   string    `gorm:"type:longtext;not null;comment:'内容：inventory为文本，vars/args为JSON'"`
	Remark    string    `gorm:"size:500;comment:'备注'"`
	CreatedBy string    `gorm:"size:50;comment:'创建人'"`
	UpdatedBy string    `gorm:"size:50;comment:'更新人'"`
	CreatedAt time.Time `gorm:"not null;comment:'创建时间'"`
	UpdatedAt time.Time `gorm:"not null;comment:'更新时间'"`
}

func (ConfigAnsible) TableName() string {
	return "config_ansible"
}