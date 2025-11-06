// 菜单 服务层
// author xiaoRui

package service

import (
	"dodevops-api/api/system/dao"
	"dodevops-api/api/system/model"
	"dodevops-api/common/result"
	"github.com/gin-gonic/gin"
)

type ISysMenuService interface {
	CreateSysMenu(c *gin.Context, menu model.SysMenu)                  // 创建菜单
	QuerySysMenuVoList(c *gin.Context)                                 // 查询菜单
	GetSysMenu(c *gin.Context, Id int)                                 // 获取菜单
	UpdateSysMenu(c *gin.Context, menu model.SysMenu)                  // 修改菜单
	DeleteSysMenu(c *gin.Context, dto model.SysMenuIdDto)              // 删除菜单
	GetSysMenuList(c *gin.Context, MenuName string, MenuStatus string) // 查询菜单列表
}

type SysMenuServiceImpl struct{}

// 新增菜单
func (s SysMenuServiceImpl) CreateSysMenu(c *gin.Context, sysMenu model.SysMenu) {
	bool := dao.CreateSysMenu(sysMenu)
	if !bool {
		result.Failed(c, int(result.ApiCode.MENUISEXIST), result.ApiCode.GetMessage(result.ApiCode.MENUISEXIST))
		return
	}
	result.Success(c, true)
}

// 查询新增选项列表
func (s SysMenuServiceImpl) QuerySysMenuVoList(c *gin.Context) {
	result.Success(c, dao.QuerySysMenuVoList())

}

// 根据id查询菜单详情
func (s SysMenuServiceImpl) GetSysMenu(c *gin.Context, Id int) {
	result.Success(c, dao.GetSysMenu(Id))
}

// 修改菜单
func (s SysMenuServiceImpl) UpdateSysMenu(c *gin.Context, menu model.SysMenu) {
	result.Success(c, dao.UpdateSysMenu(menu))
}

// 删除菜单
func (s SysMenuServiceImpl) DeleteSysMenu(c *gin.Context, dto model.SysMenuIdDto) {
	bool := dao.DeleteSysMenu(dto)
	if !bool {
		result.Failed(c, int(result.ApiCode.DELSYSMENUFAILED), result.ApiCode.GetMessage(result.ApiCode.DELSYSMENUFAILED))
		return
	}
	result.Success(c, true)
}

// 查询菜单列表
func (s SysMenuServiceImpl) GetSysMenuList(c *gin.Context, MenuName string, MenuStatus string) {
	result.Success(c, dao.GetSysMenuList(MenuName, MenuStatus))
}

var sysMenuService = SysMenuServiceImpl{}

func SysMenuService() ISysMenuService {
	return &sysMenuService
}
