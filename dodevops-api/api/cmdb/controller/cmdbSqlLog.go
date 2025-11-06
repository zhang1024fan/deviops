package controller

import (
	"dodevops-api/api/cmdb/dao"
	"dodevops-api/api/cmdb/model"
	"dodevops-api/api/cmdb/service"
	"dodevops-api/common"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Tags CMDB数据库
// 分页获取SQL操作日志列表
// @Summary 分页获取SQL操作日志列表接口
// @Produce json
// @Description 分页获取SQL操作日志列表接口
// @Param pageSize query int false "每页数"
// @Param pageNum query int false "分页数"
// @Param execUser query string false "执行用户"
// @Param beginTime query string false "开始时间"
// @Param endTime query string false "结束时间"
// @Success 200 {object} result.Result
// @router /api/v1/cmdb/sqlLog/list [get]
// @Security ApiKeyAuth
func GetCmdbSqlLogList(c *gin.Context) {
	ExecUser := c.Query("execUser")
	BeginTime := c.Query("beginTime")
	EndTime := c.Query("endTime")
	
	// 设置默认分页参数
	PageSize, err := strconv.Atoi(c.Query("pageSize"))
	if err != nil || PageSize <= 0 {
		PageSize = 10 // 默认每页10条
	}
	
	PageNum, err := strconv.Atoi(c.Query("pageNum")) 
	if err != nil || PageNum <= 0 {
		PageNum = 1 // 默认第一页
	}

	service.NewCmdbSQLRecordService(dao.NewCmdbSQLRecordDao(common.GetDB())).GetCmdbSqlLogList(c, ExecUser, BeginTime, EndTime, PageSize, PageNum)
}

// @Tags CMDB数据库
// 根据id删除SQL操作日志
// @Summary 根据id删除SQL操作日志
// @Produce json
// @Description 根据id删除SQL操作日志
// @Param data body model.CmdbSqlLogIdDto true "data"
// @Success 200 {object} result.Result
// @router /api/v1/cmdb/sqlLog/delete [delete]
// @Security ApiKeyAuth
func DeleteCmdbSqlLogById(c *gin.Context) {
	var dto model.CmdbSqlLogIdDto
	_ = c.BindJSON(&dto)
	service.NewCmdbSQLRecordService(dao.NewCmdbSQLRecordDao(common.GetDB())).DeleteCmdbSqlLogById(c, dto)
}

// @Tags CMDB数据库
// 清空SQL操作日志
// @Summary 清空SQL操作日志接口
// @Produce json
// @Description 清空SQL操作日志接口
// @Success 200 {object} result.Result
// @router /api/v1/cmdb/sqlLog/clean [delete]
// @Security ApiKeyAuth
func CleanCmdbSqlLog(c *gin.Context) {
	service.NewCmdbSQLRecordService(dao.NewCmdbSQLRecordDao(common.GetDB())).CleanCmdbSqlLog(c)
}
