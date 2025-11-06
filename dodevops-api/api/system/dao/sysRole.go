// 角色 数据层
// author xiaoRui

package dao

import (
	"dodevops-api/api/system/model"
	"dodevops-api/common/util"
	. "dodevops-api/pkg/db"
	"gorm.io/gorm"
	"time"
)

// 根据名称查询
func GetSysRoleByName(roleName string) (sysRole model.SysRole) {
	Db.Where("role_name = ?", roleName).First(&sysRole)
	return sysRole
}

// 根据角色Key查询
func GetSysRoleByKey(roleKey string) (sysRole model.SysRole) {
	Db.Where("role_key = ?", roleKey).First(&sysRole)
	return sysRole
}

// 新增角色
func CreateSysRole(dto model.AddSysRoleDto) bool {
	sysRoleByName := GetSysRoleByName(dto.RoleName)
	if sysRoleByName.ID > 0 {
		return false
	}
	sysRoleByKey := GetSysRoleByKey(dto.RoleKey)
	if sysRoleByKey.ID > 0 {
		return false
	}
	addSysRole := model.SysRole{
		RoleName:    dto.RoleName,
		RoleKey:     dto.RoleKey,
		Description: dto.Description,
		Status:      dto.Status,
		CreateTime:  util.HTime{Time: time.Now()},
	}
	tx := Db.Create(&addSysRole)
	if tx.RowsAffected > 0 {
		return true
	}
	return false
}

// 根据id获取详情
func GetSysRoleById(Id int) (sysRole model.SysRole) {
	Db.First(&sysRole, Id)
	return sysRole
}

// 修改角色
func UpdateSysRole(dto model.UpdateSysRoleDto) (sysRole model.SysRole) {
	Db.First(&sysRole, dto.Id)
	sysRole.RoleName = dto.RoleName
	sysRole.RoleKey = dto.RoleKey
	sysRole.Status = dto.Status
	if dto.Description != "" {
		sysRole.Description = dto.Description
	}
	Db.Save(&sysRole)
	return sysRole
}

// 根据id删除角色
func DeleteSysRoleById(dto model.SysRoleIdDto) {
	Db.Table("sys_role").Delete(&model.SysRole{}, dto.Id)
	Db.Table("sys_role_menu").Where("role_id = ?", dto.Id).Delete(&model.SysRoleMenu{})
}

// 角色状态启用/停用
func UpdateSysRoleStatus(dto model.UpdateSysRoleStatusDto) bool {
	var sysRole model.SysRole
	Db.First(&sysRole, dto.Id)
	sysRole.Status = dto.Status
	tx := Db.Save(&sysRole)
	if tx.RowsAffected > 0 {
		return true
	}
	return false
}

// 分页查询角色列表
func GetSysRoleList(PageNum, PageSize int, RoleName, status, BeginTime, EndTime string) (sysRole []*model.SysRole, count int64) {
	curDb := Db.Table("sys_role")
	if RoleName != "" {
		curDb = curDb.Where("role_name = ?", RoleName)
	}
	if BeginTime != "" && EndTime != "" {
		curDb = curDb.Where("create_time BETWEEN ? AND ?", BeginTime, EndTime)
	}
	if status != "" {
		curDb = curDb.Where("status = ?", status)
	}
	curDb.Count(&count)
	curDb.Limit(PageSize).Offset((PageNum - 1) * PageSize).Order("create_time DESC").Find(&sysRole)
	return sysRole, count
}

// 角色下拉列表
func QuerySysRoleVoList() (sysRoleVo []model.SysRoleVo) {
	Db.Table("sys_role").Select("id, role_name").Scan(&sysRoleVo)
	return sysRoleVo
}

// 根据角色的id查询菜单权限数据列表
func QueryRoleMenuIdList(Id int) (idVo []model.IdVo) {
	Db.Table("sys_menu").
		Select("id").
		Where("id IN (SELECT menu_id FROM sys_role_menu WHERE role_id = ?)", Id).
		Scan(&idVo)
	return idVo
}

// 分配权限
func AssignPermissions(menu model.RoleMenu) (err error) {
	err = Db.Transaction(func(tx *gorm.DB) error {
		// 删除角色权限
		if err := tx.Table("sys_role_menu").Where("role_id = ?", menu.Id).Delete(&model.SysRoleMenu{}).Error; err != nil {
			return err
		}
		
		// 批量插入新权限
		if len(menu.MenuIds) > 0 {
			var entities []model.SysRoleMenu
			for _, menuId := range menu.MenuIds {
				entities = append(entities, model.SysRoleMenu{
					RoleId: menu.Id,
					MenuId: menuId,
				})
			}
			if err := tx.CreateInBatches(entities, 100).Error; err != nil {
				return err
			}
		}
		return nil
	})
	return err
}
