// 菜单 控制层
// author xiaoRui

package controller

import (
	"dodevops-api/api/system/model"
	"dodevops-api/api/system/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

var sysMenu model.SysMenu

// @Tags System系统管理
// 新增菜单
// @Summary 新增菜单接口
// @Produce json
// @Description 新增菜单接口
// @Param data body model.SysMenu true "data"
// @Success 200 {object} result.Result
// @router /api/v1/menu/add [post]
// @Security ApiKeyAuth
func CreateSysMenu(c *gin.Context) {
	_ = c.BindJSON(&sysMenu)
	service.SysMenuService().CreateSysMenu(c, sysMenu)
}

// @Tags System系统管理
// 查询新增选项列表
// @Summary 查询新增选项列表接口
// @Produce json
// @Description 查询新增选项列表接口
// @Success 200 {object} result.Result
// @router /api/v1/menu/vo/list [get]
// @Security ApiKeyAuth
func QuerySysMenuVoList(c *gin.Context) {
	service.SysMenuService().QuerySysMenuVoList(c)
}

// @Tags System系统管理
// 根据id查询菜单
// @Summary 根据id查询菜单
// @Produce json
// @Description 根据id查询菜单
// @Param id query int true "id"
// @Success 200 {object} result.Result
// @router /api/v1/menu/info [get]
// @Security ApiKeyAuth
func GetSysMenu(c *gin.Context) {
	Id, _ := strconv.Atoi(c.Query("id"))
	service.SysMenuService().GetSysMenu(c, Id)
}

// @Tags System系统管理
// 修改菜单
// @Summary 修改菜单接口
// @Produce json
// @Description 修改菜单接口
// @Param data body model.SysMenu true "data"
// @Success 200 {object} result.Result
// @router /api/v1/menu/update [put]
// @Security ApiKeyAuth
func UpdateSysMenu(c *gin.Context) {
	_ = c.BindJSON(&sysMenu)
	service.SysMenuService().UpdateSysMenu(c, sysMenu)
}

// @Tags System系统管理
// 根据id删除菜单
// @Summary 根据id删除菜单接口
// @Produce json
// @Description 根据id删除菜单接口
// @Param data body model.SysMenuIdDto true "data"
// @Success 200 {object} result.Result
// @router /api/v1/menu/delete [delete]
// @Security ApiKeyAuth
func DeleteSysMenu(c *gin.Context) {
	var dto model.SysMenuIdDto
	_ = c.BindJSON(&dto)
	service.SysMenuService().DeleteSysMenu(c, dto)
}

// @Tags System系统管理
// 查询菜单列表
// @Summary 查询菜单列表
// @Produce json
// @Description 查询菜单列表
// @Param menuName query string false "菜单名称"
// @Param menuStatus query string false "菜单状态"
// @Success 200 {object} result.Result
// @router /api/v1/menu/list [get]
// @Security ApiKeyAuth
func GetSysMenuList(c *gin.Context) {
	MenuName := c.Query("menuName")
	MenuStatus := c.Query("menuStatus")
	service.SysMenuService().GetSysMenuList(c, MenuName, MenuStatus)
}
