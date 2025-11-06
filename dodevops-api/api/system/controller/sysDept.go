package controller

import (
	"dodevops-api/api/system/model"
	"dodevops-api/api/system/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

var sysDept model.SysDept

// @Tags System系统管理
// 查询部门列表
// @Summary 查询部门列表接口
// @Produce json
// @Description 查询部门列表接口
// @Param deptName query string false "部门名称"
// @Param deptStatus query string false "部门状态"
// @Success 200 {object} result.Result
// @router /api/v1/dept/list [get]
// @Security ApiKeyAuth
func GetSysDeptList(c *gin.Context) {
	DeptName := c.Query("deptName")
	DeptStatus := c.Query("deptStatus")
	service.SysDeptService().GetSysDeptList(c, DeptName, DeptStatus)
}

// @Tags System系统管理
// 新增部门
// @Summary 新增部门接口
// @Produce json
// @Description 新增部门接口
// @Param data body model.SysDept true "data"
// @Success 200 {object} result.Result
// @router /api/v1/dept/add [post]
// @Security ApiKeyAuth
func CreateSysDept(c *gin.Context) {
	_ = c.BindJSON(&sysDept)
	service.SysDeptService().CreateSysDept(c, sysDept)
}

// @Tags System系统管理
// 根据id查询部门
// @Summary 根据id查询部门接口
// @Produce json
// @Description 根据id查询部门接口
// @Param id query int true "ID"
// @Success 200 {object} result.Result
// @router /api/v1/dept/info [get]
// @Security ApiKeyAuth
func GetSysDeptById(c *gin.Context) {
	Id, _ := strconv.Atoi(c.Query("id"))
	service.SysDeptService().GetSysDeptById(c, Id)
}

// @Tags System系统管理
// 修改部门
// @Summary 修改部门接口
// @Produce json
// @Description 修改部门接口
// @Param data body model.SysDept true "data"
// @Success 200 {object} result.Result
// @router /api/v1/dept/update [put]
// @Security ApiKeyAuth
func UpdateSysDept(c *gin.Context) {
	_ = c.BindJSON(&sysDept)
	service.SysDeptService().UpdateSysDept(c, sysDept)
}

// @Tags System系统管理
// 根据id删除部门
// @Summary 根据id删除部门接口
// @Produce json
// @Description 根据id删除部门接口
// @Param data body model.SysDeptIdDto true "data"
// @Success 200 {object} result.Result
// @router /api/v1/dept/delete [delete]
// @Security ApiKeyAuth
func DeleteSysDeptById(c *gin.Context) {
	var dto model.SysDeptIdDto
	_ = c.BindJSON(&dto)
	service.SysDeptService().DeleteSysDeptById(c, dto)
}

// @Tags System系统管理
// 部门下拉列表
// @Summary 部门下拉列表接口
// @Produce json
// @Description 部门下拉列表接口
// @Success 200 {object} result.Result
// @router /api/v1/dept/vo/list [get]
// @Security ApiKeyAuth
func QuerySysDeptVoList(c *gin.Context) {
	service.SysDeptService().QuerySysDeptVoList(c)
}

// @Tags System系统管理
// 获取某部门下的所有用户
// @Summary 获取某部门下的所有用户接口
// @Produce json
// @Description 获取某部门下的所有用户
// @Param deptId query int true "部门ID"
// @Success 200 {object} result.Result
// @router /api/v1/dept/users [get]
// @Security ApiKeyAuth
func GetDeptUsers(c *gin.Context) {
	deptId, _ := strconv.Atoi(c.Query("deptId"))
	service.SysDeptService().GetDeptUsers(c, deptId)
}
