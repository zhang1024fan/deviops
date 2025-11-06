// 操作日志 控制层
// author xiaoRui

package controller

import (
	"dodevops-api/api/system/model"
	"dodevops-api/api/system/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Tags System系统管理
// 分页获取操作日志列表
// @Summary 分页获取操作日志列表接口
// @Produce json
// @Description 分页获取操作日志列表接口
// @Param pageSize query int false "每页数"
// @Param pageNum query int false "分页数"
// @Param username query string false "用户名"
// @Param beginTime query string false "开始时间"
// @Param endTime query string false "结束时间"
// @Success 200 {object} result.Result
// @router /api/v1/sysOperationLog/list [get]
// @Security ApiKeyAuth
func GetSysOperationLogList(c *gin.Context) {
	Username := c.Query("username")
	BeginTime := c.Query("beginTime")
	EndTime := c.Query("endTime")
	PageSize, _ := strconv.Atoi(c.Query("pageSize"))
	PageNum, _ := strconv.Atoi(c.Query("pageNum"))
	service.SysOperationLogService().GetSysOperationLogList(c, Username, BeginTime, EndTime, PageSize, PageNum)
}

// @Tags System系统管理
// 根据id删除操作日志
// @Summary 根据id删除操作日志
// @Produce json
// @Description 根据id删除操作日志
// @Param data body model.SysOperationLogIdDto true "data"
// @Success 200 {object} result.Result
// @router /api/v1/sysOperationLog/delete [delete]
// @Security ApiKeyAuth
func DeleteSysOperationLogById(c *gin.Context) {
	var dto model.SysOperationLogIdDto
	_ = c.BindJSON(&dto)
	service.SysOperationLogService().DeleteSysOperationLogById(c, dto)
}

// @Tags System系统管理
// 批量删除操作日志
// @Summary 批量删除操作日志接口
// @Produce json
// @Description 批量删除操作日志接口
// @Param data body model.BatchDeleteSysOperationLogDto true "data"
// @Success 200 {object} result.Result
// @router /api/v1/sysOperationLog/batch/delete [delete]
// @Security ApiKeyAuth
func BatchDeleteSysOperationLog(c *gin.Context) {
	var dto model.BatchDeleteSysOperationLogDto
	_ = c.BindJSON(&dto)
	service.SysOperationLogService().BatchDeleteSysOperationLog(c, dto)
}

// @Tags System系统管理
// 清空操作日志
// @Summary 清空操作日志接口
// @Produce json
// @Description 清空操作日志接口
// @Success 200 {object} result.Result
// @router /api/v1/sysOperationLog/clean [delete]
// @Security ApiKeyAuth
func CleanSysOperationLog(c *gin.Context) {
	service.SysOperationLogService().CleanSysOperationLog(c)
}
