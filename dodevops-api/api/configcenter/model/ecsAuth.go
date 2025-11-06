// ECS认证凭证模型
// xiaoRui
package model

import (
	"dodevops-api/common/util"
)

// ECS认证凭证模型
type EcsAuth struct {
	ID         uint       `gorm:"column:id;comment:'主键';primaryKey;NOT NULL" json:"id"`
	Name       string     `gorm:"column:name;varchar(64);comment:'凭证名称';NOT NULL" json:"name"`
	Type       int        `gorm:"column:type;comment:'认证类型:1->密码,2->私钥,3->公钥(免认证)';NOT NULL" json:"type"`
	Username   string     `gorm:"column:username;varchar(64);comment:'用户名'" json:"username"`
	Password   string     `gorm:"column:password;varchar(256);comment:'密码(type=1时使用)'" json:"password"`
	PublicKey  string     `gorm:"column:public_key;type:text;comment:'私钥内容(type=2时使用，字段名历史原因)'" json:"publicKey"` // 实际存储私钥
	Port       int        `gorm:"column:port;comment:'端口号';default:22" json:"port"`
	CreateTime util.HTime `gorm:"column:create_time;comment:'创建时间';NOT NULL" json:"createTime"`
	Remark     string     `gorm:"column:remark;varchar(500);comment:'备注'" json:"remark"`
}

func (EcsAuth) TableName() string {
	return "config_ecsauth"
}

// 创建ECS密码认证DTO
type CreateEcsPasswordAuthDto struct {
	Name      string `validate:"required"` // 凭证名称
	Type      int    `validate:"required"` // 认证类型:1->密码
	Username  string `validate:"required"` // 用户名
	Password  string `validate:"required"` // 密码
	PublicKey string // 公钥
	Port      int    `validate:"required"` // 端口号
	Remark    string // 备注
}

// 创建ECS密钥认证DTO
type CreateEcsKeyAuthDto struct {
	Name      string `validate:"required"` // 凭证名称
	Type      int    `validate:"required"` // 认证类型:2->密钥
	PublicKey string `validate:"required"` // 私钥内容(字段名历史原因，实际传入私钥)
	Username  string `validate:"required"` // 用户名
	Port      int    `validate:"required"` // 端口号
	Remark    string // 备注
}

// 创建ECS公钥认证DTO (type=3, 免认证)
type CreateEcsPublicKeyAuthDto struct {
	Name     string `validate:"required"` // 凭证名称
	Type     int    `validate:"required"` // 认证类型:3->公钥免认证
	Username string `validate:"required"` // 用户名
	Port     int    `validate:"required"` // 端口号
	Remark   string // 备注
}

// ID参数
type EcsAuthIdDto struct {
	Id uint `json:"id"` // ID
}

// 更新ECS认证DTO
type UpdateEcsAuthDto struct {
	EcsAuthIdDto
	CreateEcsPasswordAuthDto
}

// 认证凭证列表VO
type EcsAuthVo struct {
	ID         uint       `json:"id"`
	Name       string     `json:"name"`
	Type       int        `json:"type"`
	Username   string     `json:"username"`
	Password   string     `json:"password"`
	PublicKey  string     `json:"publicKey"`
	Port       int        `json:"port"`
	CreateTime util.HTime `json:"createTime"`
	Remark     string     `json:"remark"`
}
