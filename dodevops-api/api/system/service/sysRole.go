// 角色 服务层
// author xiaoRui

package service

import (
	"dodevops-api/api/system/dao"
	"dodevops-api/api/system/model"
	"dodevops-api/common/result"
	"github.com/gin-gonic/gin"
)

type ISysRoleService interface {
	CreateSysRole(c *gin.Context, dto model.AddSysRoleDto)                                             // 新建角色
	GetSysRoleById(c *gin.Context, Id int)                                                             // 查询角色
	UpdateSysRole(c *gin.Context, dto model.UpdateSysRoleDto)                                          // 修改角色
	DeleteSysRoleById(c *gin.Context, dto model.SysRoleIdDto)                                          // 删除角色
	UpdateSysRoleStatus(c *gin.Context, dto model.UpdateSysRoleStatusDto)                              // 修改角色状态
	GetSysRoleList(c *gin.Context, PageNum, PageSize int, RoleName, Status, BeginTime, EndTime string) // 查询角色列表
	QuerySysRoleVoList(c *gin.Context)                                                                 // 查询角色列表
	QueryRoleMenuIdList(c *gin.Context, Id int)                                                        // 查询角色菜单id列表
	AssignPermissions(c *gin.Context, menu model.RoleMenu)                                             // 分配权限
}

type SysRoleServiceImpl struct{}

// 新建角色
func (s SysRoleServiceImpl) CreateSysRole(c *gin.Context, dto model.AddSysRoleDto) {
	bool := dao.CreateSysRole(dto)
	if !bool {
		result.Failed(c, int(result.ApiCode.ROLENAMEALREADYEXISTS), result.ApiCode.GetMessage(result.ApiCode.ROLENAMEALREADYEXISTS))
		return
	}
	result.Success(c, true)
}

// 根据id查询角色
func (s SysRoleServiceImpl) GetSysRoleById(c *gin.Context, Id int) {
	sysRole := dao.GetSysRoleById(Id)
	result.Success(c, sysRole)
}

// 修改角色
func (s SysRoleServiceImpl) UpdateSysRole(c *gin.Context, dto model.UpdateSysRoleDto) {
	sysRole := dao.UpdateSysRole(dto)
	result.Success(c, sysRole)
}

// 根据id删除角色
func (s SysRoleServiceImpl) DeleteSysRoleById(c *gin.Context, dto model.SysRoleIdDto) {
	dao.DeleteSysRoleById(dto)
	result.Success(c, true)
}

// 角色状态启用/停用
func (s SysRoleServiceImpl) UpdateSysRoleStatus(c *gin.Context, dto model.UpdateSysRoleStatusDto) {
	bool := dao.UpdateSysRoleStatus(dto)
	if !bool {
		return
	}
	result.Success(c, true)
}

// 分页查询角色列表
func (s SysRoleServiceImpl) GetSysRoleList(c *gin.Context, PageNum, PageSize int, RoleName, Status, BeginTime, EndTime string) {
	if PageSize < 1 {
		PageSize = 10
	}
	if PageNum < 1 {
		PageNum = 1
	}
	sysRole, count := dao.GetSysRoleList(PageNum, PageSize, RoleName, Status, BeginTime, EndTime)
	result.Success(c, map[string]interface{}{"total": count, "pageSize": PageSize, "pageNum": PageNum, "list": sysRole})
}

// 角色下拉列表
func (s SysRoleServiceImpl) QuerySysRoleVoList(c *gin.Context) {
	result.Success(c, dao.QuerySysRoleVoList())
}

// 根据角色id查询菜单数据
func (s *SysRoleServiceImpl) QueryRoleMenuIdList(c *gin.Context, Id int) {
	roleMenuIdList := dao.QueryRoleMenuIdList(Id)
	var idList = make([]int, 0)
	for _, id := range roleMenuIdList {
		idList = append(idList, id.Id)
	}
	result.Success(c, idList)
}

// 分配权限
func (s *SysRoleServiceImpl) AssignPermissions(c *gin.Context, menu model.RoleMenu) {
	// 立即返回成功响应
	result.Success(c, true)
	
	// 异步处理权限分配
	go func() {
		dao.AssignPermissions(menu)
	}()
}

var sysRoleService = SysRoleServiceImpl{}

func SysRoleService() ISysRoleService {
	return &sysRoleService
}
