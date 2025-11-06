package model

import (
	"dodevops-api/common/util"
)

// SQLRequest SQL执行请求参数
type SQLRequest struct {
	InstanceID string `json:"instanceId" binding:"required"` // 实例ID
	Database   string `json:"database" binding:"required"`   // 数据库名称
	SQL        string `json:"sql" binding:"required"`         // SQL语句
	ExecUser   string `json:"execUser" binding:"required"`    // 执行用户
}
// SQL执行结果
type CmdbSQLRecord struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	InstanceID    string    `gorm:"size:64;not null" json:"instanceId"`      // 实例ID
	Database      string    `gorm:"size:128;not null" json:"database"`     // 数据库名
	OperationType string    `gorm:"size:32;not null" json:"operationType"`  // 操作类型(SELECT/INSERT/UPDATE/DELETE)
	SQLContent    string    `gorm:"type:text;not null" json:"sqlContent"`    // SQL内容
	ExecUser      string    `gorm:"size:64;not null" json:"execUser"`       // 执行用户
	IP            string    `gorm:"size:64;not null" json:"ip"`           // 用户IP
	ScannedRows   int64     `gorm:"default:0" json:"scannedRows"`           // 扫描行数
	AffectedRows  int64     `gorm:"default:0" json:"affectedRows"`          // 影响行数
	ExecutionTime int64     `gorm:"default:0" json:"executionTime"`         // 执行耗时(毫秒)
	ReturnedRows  int64     `gorm:"default:0" json:"returnedRows"`         // 返回行数
	Result        string    `gorm:"size:32;not null" json:"result"`        // 执行结果(SUCCESS/FAILED)
	QueryTime     util.HTime `gorm:"index;not null" json:"queryTime"`       // 查询时间
}

// Id参数
type CmdbSqlLogIdDto struct {
	Id uint `json:"id"` // ID
}

// 批量删除id参数
type BatchDeleteCmdbSqlLogDto struct {
	Ids []uint `json:"ids"` //id列表
}

func (CmdbSQLRecord) TableName() string {
	return "cmdb_sql_log"
}
