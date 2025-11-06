package dao

import (
	"dodevops-api/api/system/model"
	"dodevops-api/common/util"
	. "dodevops-api/pkg/db"
	"time"
)

// 查询部门列表
func GetSysDeptList(DeptName string, DeptStatus string) (sysDept []model.SysDept) {
	curDb := Db.Table("sys_dept")
	if DeptName != "" {
		curDb = curDb.Where("dept_name = ?", DeptName)
	}
	if DeptStatus != "" {
		curDb = curDb.Where("dept_status = ?", DeptStatus)
	}
	curDb.Find(&sysDept)
	return sysDept
}

// 根据部门名称查询
func GetSysDeptByName(deptName string) (sysDept model.SysDept) {
	Db.Where("dept_name = ?", deptName).First(&sysDept)
	return sysDept
}

// 新增部门 之前先查询部门是否存在。
func CreateSysDept(sysDept model.SysDept) bool {
	sysDeptByName := GetSysDeptByName(sysDept.DeptName)
	if sysDeptByName.ID > 0 {
		return false
	}
	if sysDept.DeptType == 1 {
		sysDept := model.SysDept{
			DeptStatus: sysDept.DeptStatus,
			ParentId:   0,
			DeptType:   sysDept.DeptType,
			DeptName:   sysDept.DeptName,
			CreateTime: util.HTime{Time: time.Now()},
		}
		Db.Create(&sysDept)
		return true
	} else {
		sysDept := model.SysDept{
			DeptStatus: sysDept.DeptStatus,
			ParentId:   sysDept.ParentId,
			DeptType:   sysDept.DeptType,
			DeptName:   sysDept.DeptName,
			CreateTime: util.HTime{Time: time.Now()},
		}
		Db.Create(&sysDept)
		return true
	}
}

// 根据id查询部门
func GetSysDeptById(Id int) (sysDept model.SysDept) {
	Db.First(&sysDept, Id)
	return sysDept
}

// 修改部门
func UpdateSysDept(dept model.SysDept) (sysDept model.SysDept) {
	Db.First(&sysDept, dept.ID)
	sysDept.ParentId = dept.ParentId
	sysDept.DeptType = dept.DeptType
	sysDept.DeptName = dept.DeptName
	sysDept.DeptStatus = dept.DeptStatus
	Db.Save(&sysDept)
	return sysDept
}

// 查询部门是否有人
func GetSysAdminDept(id int) (sysAdmin model.SysAdmin) {
	Db.Where("dept_id = ?", id).First(&sysAdmin)
	return sysAdmin
}

// 删除部门前先确认部门是否有人
func DeleteSysDeptById(dto model.SysDeptIdDto) bool {
	sysAdmin := GetSysAdminDept(dto.Id)
	if sysAdmin.ID > 0 {
		return false
	}
	Db.Where("parent_id = ?", dto.Id).Delete(&model.SysDept{})
	Db.Delete(&model.SysDept{}, dto.Id)
	return true
}

// 部门下拉列表
func QuerySysDeptVoList() (sysDeptVo []model.SysDeptVo) {
	Db.Table("sys_dept").Select("id, dept_name AS label, parent_id").Scan(&sysDeptVo)
	return sysDeptVo
}

// 获取某部门下的所有用户
func GetUsersByDeptId(deptId int) []model.SysAdminVo {
	var users []model.SysAdminVo
	Db.Table("sys_admin").
		Select(`
			sys_admin.*,
			sys_admin_role.role_id AS role_id,
			sys_dept.dept_name AS dept_name,
			sys_post.post_name AS post_name,
			sys_role.role_name AS role_name
		`).
		Joins("LEFT JOIN sys_admin_role ON sys_admin.id = sys_admin_role.admin_id").
		Joins("LEFT JOIN sys_dept ON sys_admin.dept_id = sys_dept.id").
		Joins("LEFT JOIN sys_post ON sys_admin.post_id = sys_post.id").
		Joins("LEFT JOIN sys_role ON sys_admin_role.role_id = sys_role.id").
		Where("sys_admin.dept_id = ?", deptId).
		Scan(&users)
	return users
}
