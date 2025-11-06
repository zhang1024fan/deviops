package controller

import (
	"dodevops-api/api/system/model"
	"dodevops-api/api/system/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

var sysPost model.SysPost

// @Tags System系统管理
// @Summary 新增岗位接口
// @Produce json
// @Description 新增岗位接口
// @Param data body model.SysPost true "data"
// @Success 200 {object} result.Result
// @router /api/v1/post/add [post]
// @Security ApiKeyAuth
func CreateSysPost(c *gin.Context) {
	_ = c.BindJSON(&sysPost)
	service.SysPostService().CreateSysPost(c, sysPost)
}

// @Tags System系统管理
// 分页查询岗位列表
// @Summary 分页查询岗位列表
// @Produce json
// @Description 分页查询岗位列表
// @Param pageNum query int false "分页数"
// @Param pageSize query int false "每页数"
// @Param postName query string false "岗位名称"
// @Param postStatus query string false "状态：1->启用,2->禁用"
// @Param beginTime query string false "开始时间"
// @Param endTime query string false "结束时间"
// @Success 200 {object} result.Result
// @router /api/v1/post/list [get]
// @Security ApiKeyAuth
func GetSysPostList(c *gin.Context) {
	PageNum, _ := strconv.Atoi(c.Query("pageNum"))
	PageSize, _ := strconv.Atoi(c.Query("pageSize"))
	PostName := c.Query("postName")
	PostStatus := c.Query("postStatus")
	BeginTime := c.Query("beginTime")
	EndTime := c.Query("endTime")
	service.SysPostService().GetSysPostList(c, PageNum, PageSize, PostName, PostStatus, BeginTime, EndTime)
}

// @Tags System系统管理
// 根据id查询岗位
// @Summary 根据id查询岗位
// @Produce json
// @Description 根据id查询岗位
// @Param id query int true "ID"
// @Success 200 {object} result.Result
// @router /api/v1/post/info [get]
// @Security ApiKeyAuth
func GetSysPostById(c *gin.Context) {
	Id, _ := strconv.Atoi(c.Query("id"))
	service.SysPostService().GetSysPostById(c, Id)
}

// @Tags System系统管理
// 修改岗位
// @Summary 修改岗位接口
// @Produce json
// @Description 修改岗位接口
// @Param data body model.SysPost true "data"
// @Success 200 {object} result.Result
// @router /api/v1/post/update [put]
// @Security ApiKeyAuth
func UpdateSysPost(c *gin.Context) {
	_ = c.BindJSON(&sysPost)
	service.SysPostService().UpdateSysPost(c, sysPost)
}

// @Tags System系统管理
// 根据id删除岗位
// @Summary 根据id删除岗位接口
// @Produce json
// @Description 根据id删除岗位接口
// @Param data body model.SysPostIdDto true "data"
// @Success 200 {object} result.Result
// @router /api/v1/post/delete [delete]
// @Security ApiKeyAuth
func DeleteSysPostById(c *gin.Context) {
	var dto model.SysPostIdDto
	_ = c.BindJSON(&dto)
	service.SysPostService().DeleteSysPostById(c, dto)
}

// @Tags System系统管理
// 批量删除岗位
// @Summary 批量删除岗位接口
// @Produce json
// @Description 批量删除岗位接口
// @Param data body model.DelSysPostDto true "data"
// @Success 200 {object} result.Result
// @router /api/v1/post/batch/delete [delete]
// @Security ApiKeyAuth
func BatchDeleteSysPost(c *gin.Context) {
	var dto model.DelSysPostDto
	_ = c.BindJSON(&dto)
	service.SysPostService().BatchDeleteSysPost(c, dto)
}

// @Tags System系统管理
// 岗位状态修改
// @Summary 岗位状态修改接口
// @Produce json
// @Description 岗位状态修改接口
// @Param data body model.UpdateSysPostStatusDto true "data"
// @Success 200 {object} result.Result
// @router /api/v1/post/updateStatus [put]
// @Security ApiKeyAuth
func UpdateSysPostStatus(c *gin.Context) {
	var dto model.UpdateSysPostStatusDto
	_ = c.BindJSON(&dto)
	service.SysPostService().UpdateSysPostStatus(c, dto)
}

// @Tags System系统管理
// 岗位下拉列表
// @Summary 岗位下拉列表
// @Produce json
// @Description 岗位下拉列表
// @Success 200 {object} result.Result
// @router /api/v1/post/vo/list [get]
// @Security ApiKeyAuth
func QuerySysPostVoList(c *gin.Context) {
	service.SysPostService().QuerySysPostVoList(c)
}
