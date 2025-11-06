// 用户与角色关系相关模型
// author xiaoRui

package model

type SysAdminRole struct {
	RoleId  uint `gorm:"column:role_id;comment:'角色id';NOT NULL" json:"roleId"`  // 角色id
	AdminId uint `gorm:"column:admin_id;comment:'用户id';NOT NULL" json:"menuId"` // 用户id
}

func (SysAdminRole) TableName() string {
	return "sys_admin_role"
}
