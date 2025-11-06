// 操作日志相关模型
// author xiaoRui

package model

import "dodevops-api/common/util"

// 操作日志
type SysOperationLog struct {
	ID          uint       `gorm:"column:id;comment:'主键';primaryKey;NOT NULL" json:"id"`                         // ID
	AdminId     uint       `gorm:"column:admin_id;comment:'管理员id';NOT NULL" json:"adminId"`                      // 管理员id
	Username    string     `gorm:"column:username;type:varchar(64);comment:'管理员账号';NOT NULL" json:"username"`    // 管理员账号
	Method      string     `gorm:"column:method;type:varchar(64);comment:'请求方式';NOT NULL" json:"method"`          // 请求方式
	Ip          string     `gorm:"column:ip;type:varchar(64);comment:'IP'" json:"ip"`                             // IP
	Url         string     `gorm:"column:url;type:varchar(500);comment:'URL'" json:"url"`                         // URL
	Description string     `gorm:"column:description;type:varchar(255);comment:'操作描述'" json:"description"`        // 操作描述
	CreateTime  util.HTime `gorm:"column:create_time;comment:'创建时间';NOT NULL" json:"createTime"`                 // 创建时间
}

func (SysOperationLog) TableName() string {
	return "sys_operation_log"
}

// Id参数
type SysOperationLogIdDto struct {
	Id uint `json:"id"` // ID
}

// 批量删除id参数
type BatchDeleteSysOperationLogDto struct {
	Ids []uint //id列表
}
