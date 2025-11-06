package model

import (
	"dodevops-api/common/util"
)

// CmdbSQL 定义CMDB中的数据库模型
type CmdbSQL struct {
	ID          uint      `gorm:"primaryKey" json:"id"`              // 主键ID
	Name        string    `gorm:"size:100;not null" json:"name"`     // 数据库名称
	Type        int       `gorm:"type:integer;not null" json:"type"` // 数据库类型(1=MySQL 2=PostgreSQL 3=Redis 4=MongoDB 5=Elasticsearch)
	AccountID   uint      `gorm:"not null" json:"accountId"`         // 所属账号ID
	GroupID     uint      `gorm:"not null" json:"groupId"`           // 所属业务组ID
	Tags        string    `gorm:"size:255" json:"tags"`              // 标签(多个标签用逗号分隔)
	Description string    `gorm:"size:500" json:"description"`       // 描述/备注
	CreatedAt   util.HTime `json:"createdAt"`                         // 创建时间
	UpdatedAt   util.HTime `json:"updatedAt"`                         // 更新时间
}

func (CmdbSQL) TableName() string {
	return "cmdb_sql"
}
