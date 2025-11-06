// 角色 控制层
// author xiaoRui

package controller

import (
	"dodevops-api/api/system/model"
	"dodevops-api/api/system/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Tags System系统管理
// 新增角色
// @Summary 新增角色接口
// @Produce json
// @Description 新增角色接口
// @Param data body model.AddSysRoleDto true "data"
// @Success 200 {object} result.Result
// @router /api/v1/role/add [post]
// @Security ApiKeyAuth
func CreateSysRole(c *gin.Context) {
	var dto model.AddSysRoleDto
	_ = c.BindJSON(&dto)
	service.SysRoleService().CreateSysRole(c, dto)
}

// @Tags System系统管理
// 根据id查询角色
// @Summary 根据id查询角色接口
// @Produce json
// @Description 根据id查询角色接口
// @Param id query int true "Id"
// @Success 200 {object} result.Result
// @router /api/v1/role/info [get]
// @Security ApiKeyAuth
func GetSysRoleById(c *gin.Context) {
	Id, _ := strconv.Atoi(c.Query("id"))
	service.SysRoleService().GetSysRoleById(c, Id)
}

// @Tags System系统管理
// 修改角色
// @Summary 修改角色
// @Produce json
// @Description 修改角色
// @Param data body model.UpdateSysRoleDto true "data"
// @Success 200 {object} result.Result
// @router /api/v1/role/update [put]
// @Security ApiKeyAuth
func UpdateSysRole(c *gin.Context) {
	var dto model.UpdateSysRoleDto
	_ = c.BindJSON(&dto)
	service.SysRoleService().UpdateSysRole(c, dto)
}

// @Tags System系统管理
// 根据id删除角色
// @Summary 根据id删除角色接口
// @Produce json
// @Description 根据id删除角色接口
// @Param data body model.SysRoleIdDto true "data"
// @Success 200 {object} result.Result
// @router /api/v1/role/delete [delete]
// @Security ApiKeyAuth
func DeleteSysRoleById(c *gin.Context) {
	var dto model.SysRoleIdDto
	_ = c.BindJSON(&dto)
	service.SysRoleService().DeleteSysRoleById(c, dto)
}

// @Tags System系统管理
// 角色状态启用/停用
// @Summary 角色状态启用/停用接口
// @Produce json
// @Description 角色状态启用/停用接口
// @Param data body model.UpdateSysRoleStatusDto true "data"
// @Success 200 {object} result.Result
// @router /api/v1/role/updateStatus [put]
// @Security ApiKeyAuth
func UpdateSysRoleStatus(c *gin.Context) {
	var dto model.UpdateSysRoleStatusDto
	_ = c.BindJSON(&dto)
	service.SysRoleService().UpdateSysRoleStatus(c, dto)
}

// @Tags System系统管理
// 分页查询角色列表
// @Summary 分页查询角色列表接口
// @Produce json
// @Description 分页查询角色列表接口
// @Param pageNum query int false "分页数"
// @Param pageSize query int false "每页数"
// @Param roleName query string false "角色名称"
// @Param status query string false "帐号启用状态：1->启用,2->禁用"
// @Param beginTime query string false "开始时间"
// @Param endTime query string false "结束时间"
// @Success 200 {object} result.Result
// @router /api/v1/role/list [get]
// @Security ApiKeyAuth
func GetSysRoleList(c *gin.Context) {
	PageNum, _ := strconv.Atoi(c.Query("pageNum"))
	PageSize, _ := strconv.Atoi(c.Query("pageSize"))
	RoleName := c.Query("roleName")
	Status := c.Query("status")
	BeginTime := c.Query("beginTime")
	EndTime := c.Query("endTime")
	service.SysRoleService().GetSysRoleList(c, PageNum, PageSize, RoleName, Status, BeginTime, EndTime)
}

// @Tags System系统管理
// 角色下拉列表
// @Summary 角色下拉列表
// @Produce json
// @Description 角色下拉列表
// @Success 200 {object} result.Result
// @router /api/v1/role/vo/list [get]
// @Security ApiKeyAuth
func QuerySysRoleVoList(c *gin.Context) {
	service.SysRoleService().QuerySysRoleVoList(c)
}

// @Tags System系统管理
// 根据角色id查询菜单数据
// @Summary 根据角色id查询菜单数据接口
// @Produce json
// @Description 根据角色id查询菜单数据接口
// @Param id query int true "Id"
// @Success 200 {object} result.Result
// @router /api/v1/role/vo/idList [get]
// @Security ApiKeyAuth
func QueryRoleMenuIdList(c *gin.Context) {
	Id, _ := strconv.Atoi(c.Query("id"))
	service.SysRoleService().QueryRoleMenuIdList(c, Id)
}

// @Tags System系统管理
// AssignPermissions 分配权限
// @Summary 分配权限接口
// @Produce json
// @Description 分配权限接口
// @Param data body model.RoleMenu true "data"
// @Success 200 {object} result.Result
// @router /api/v1/role/assignPermissions [put]
// @Security ApiKeyAuth
func AssignPermissions(c *gin.Context) {
	var RoleMenu model.RoleMenu
	_ = c.BindJSON(&RoleMenu)
	service.SysRoleService().AssignPermissions(c, RoleMenu)
}
