// 角色 相关模型
// author xiaoRui

package model

import "dodevops-api/common/util"

// 角色模型
type SysRole struct {
	ID          uint       `gorm:"column:id;comment:'主键';primaryKey;NOT NULL" json:"id"`                        // ID
	RoleName    string     `gorm:"column:role_name;varchar(64);comment:'角色名称';NOT NULL" json:"roleName"`        // 角色名称
	RoleKey     string     `gorm:"column:role_key;varchar(64);comment:'权限字符串';NOT NULL" json:"roleKey"`         // 权限字符串
	Status      int        `gorm:"column:status;default:1;comment:'帐号启用状态：1->启用,2->禁用';NOT NULL" json:"status"` // 帐号启用状态：1->启用,2->禁用
	Description string     `gorm:"column:description;varchar(500);comment:'描述'" json:"description"`             // 描述
	CreateTime  util.HTime `gorm:"column:create_time;comment:'创建时间';NOT NULL" json:"createTime"`                // 创建时间
}

func (SysRole) TableName() string {
	return "sys_role"
}

// 新增参数
type AddSysRoleDto struct {
	RoleName    string // 角色名称
	RoleKey     string // 角色key
	Status      int    // 状态：1->启用,2->禁用
	Description string // 描述
}

// UpdateSysRoleDto 修改参数
type UpdateSysRoleDto struct {
	Id          uint   // Id
	RoleName    string // 角色名称
	RoleKey     string // 角色key
	Status      int    // 状态：1->启用,2->禁用
	Description string // 描述
}

// Id参数
type SysRoleIdDto struct {
	Id uint `json:"id"` // ID
}

// UpdateSysRoleStatusDto 设置状态参数
type UpdateSysRoleStatusDto struct {
	Id     uint // ID
	Status int  // 状态：1->启用,2->禁用
}

// 角色下拉列表
type SysRoleVo struct {
	Id       int    `json:"id"`       // ID
	RoleName string `json:"roleName"` // 角色名称
}

// IdVo 当前角色的菜单权限id
type IdVo struct {
	Id int `json:"id"` // ID
}

// RoleMenu 角色id,菜单id视图
type RoleMenu struct {
	Id      uint   `json:"id" binding:"required"`      // ID
	MenuIds []uint `json:"menuIds" binding:"required"` // 菜单id列表
}
