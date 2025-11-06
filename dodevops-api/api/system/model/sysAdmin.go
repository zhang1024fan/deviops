// 用户相关结构体
// xiaoRui
package model

import (
	"dodevops-api/common/util"
)

// 用户模型对象
type SysAdmin struct {
	ID         uint       `gorm:"column:id;comment:'主键';primaryKey;NOT NULL" json:"id"`                        //ID
	PostId     int        `gorm:"column:post_id;comment:'岗位id'" json:"postId"`                                 // 岗位id
	DeptId     int        `gorm:"column:dept_id;comment:'部门id'" json:"deptId"`                                 // 部门id
	Username   string     `gorm:"column:username;varchar(64);comment:'用户账号';NOT NULL" json:"username"`         // 用户账号
	Password   string     `gorm:"column:password;varchar(64);comment:'密码';NOT NULL" json:"password"`           // 密码
	Nickname   string     `gorm:"column:nickname;varchar(64);comment:'昵称'" json:"nickname"`                    // 昵称
	Status     int        `gorm:"column:status;default:1;comment:'帐号启用状态:1->启用,2->禁用';NOT NULL" json:"status"` // 帐号启用状态：1->启用,2->禁用
	Icon       string     `gorm:"column:icon;varchar(500);comment:'头像'" json:"icon"`                           //  头像
	Email      string     `gorm:"column:email;varchar(64);comment:'邮箱'" json:"email"`                          // 邮箱
	Phone      string     `gorm:"column:phone;varchar(64);comment:'电话'" json:"phone"`                          // 电话
	Note       string     `gorm:"column:note;varchar(500);comment:'备注'" json:"note"`                           // 备注
	CreateTime util.HTime `gorm:"column:create_time;comment:'创建时间';NOT NULL" json:"createTime"`                // 创建时间
}

func (SysAdmin) TableName() string {
	return "sys_admin"
}

// 鉴权用户结构体
type JwtAdmin struct {
	ID       uint   `json:"id"`      //ID
	Username string `json:"username"` //用户名
	Nickname string `json:"nickname"` //昵称
	Icon     string `json:"icon"`     //头像
	Email    string `json:"email"`    //邮箱
	Phone    string `json:"phone"`    //电话
	Note     string `json:"note"`     //备注
}

// 登录对象
type LoginDto struct {
	Username string `json:"username" validate:"required"`          //用户名
	Password string `json:"password" validate:"required"`          //密码
	Image    string `json:"image" validate:"required,min=4,max=6"` //验证码
	IdKey    string `json:"idKey" validate:"required"`             //uuid
}

// AddSysAdminDto 新增参数
type AddSysAdminDto struct {
	PostId   int    `validate:"required"` // 岗位id
	RoleId   uint   `validate:"required"` // 角色id
	DeptId   int    `validate:"required"` // 部门id
	Username string `validate:"required"` // 用户名
	Password string `validate:"required"` // 密码
	Nickname string `validate:"required"` // 昵称
	Phone    string `validate:"required"` // 手机号
	Email    string `validate:"required"` // 邮箱
	Note     string // 备注
	Status   int    `validate:"required"` // 状态：1->启用,2->禁用
}

// 详情视图
type SysAdminInfo struct {
	ID       uint   `json:"id"`       // ID
	Username string `json:"username"` // 用户名
	Nickname string `json:"nickname"` // 昵称
	Status   int    `json:"status"`   // 状态：1->启用,2->禁用
	PostId   int    `json:"postId"`   // 岗位id
	DeptId   int    `json:"deptId"`   // 部门id
	RoleId   uint   `json:"roleId" `  // 角色id
	Email    string `json:"email"`    // 邮箱
	Phone    string `json:"phone"`    // 手机号
	Note     string `json:"note"`     // 备注
}

// 修改参数
type UpdateSysAdminDto struct {
	Id       uint   // ID
	PostId   int    // 岗位id
	DeptId   int    // 部门id
	RoleId   uint   // 角色id
	Username string // 用户名
	Nickname string // 昵称
	Phone    string // 手机号
	Email    string // 邮箱
	Note     string // 备注
	Status   int    // 状态：1->启用,2->禁用
}

// Id参数
type SysAdminIdDto struct {
	Id uint `json:"id"` // ID
}

// 设置状态参数
type UpdateSysAdminStatusDto struct {
	Id     uint // ID
	Status int  // 状态：1->启用,2->禁用
}

// 重置密码参数
type ResetSysAdminPasswordDto struct {
	Id       uint   // ID
	Password string //密码
}

// 用户列表的vo视图
type SysAdminVo struct {
	ID         uint       `json:"id"`         // ID
	Username   string     `json:"username"`   // 用户名
	Nickname   string     `json:"nickname"`   // 昵称
	Status     int        `json:"status"`     // 状态：1->启用,2->禁用
	PostId     int        `json:"postId"`     // 岗位id
	DeptId     int        `json:"deptId"`     // 部门id
	RoleId     uint       `json:"roleId" `    // 角色id
	PostName   string     `json:"postName"`   // 岗位名称
	DeptName   string     `json:"deptName"`   // 部门名称
	RoleName   string     `json:"roleName"`   // 角色名称
	Icon       string     `json:"icon"`       // 头像
	Email      string     `json:"email"`      // 邮箱
	Phone      string     `json:"phone"`      // 电话
	Note       string     `json:"note"`       // 备注
	CreateTime util.HTime `json:"createTime"` // 创建时间
}

// 修改个人信息参数
type UpdatePersonalDto struct {
	Id       uint   //ID
	Icon     string // 头像
	Username string `validate:"required"` //用户名
	Nickname string `validate:"required"` // 昵称
	Phone    string `validate:"required"` // 电话
	Email    string `validate:"required"` // 邮箱
	Note     string `validate:"required"` // 备注
}

// 修改个人密码
type UpdatePersonalPasswordDto struct {
	Id            uint   //ID
	Password      string `validate:"required"` // 密码
	NewPassword   string `validate:"required"` // 新密码
	ResetPassword string `validate:"required"` // 重复密码
}
