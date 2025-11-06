// 服务部署 Controller层
// author xiaoRui
package controller

import (
	"dodevops-api/api/tool/model"
	"dodevops-api/api/tool/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetServicesList 获取可部署服务列表
// @Tags Tool运维工具箱
// @Summary 获取可部署服务列表
// @Description 获取所有可部署的服务列表及分类
// @Produce json
// @Success 200 {object} result.Result
// @Router /api/v1/tool/services [get]
// @Security ApiKeyAuth
func GetServicesList(c *gin.Context) {
	service.ServiceDeployService().GetServicesList(c)
}

// GetServiceDetail 获取服务详情
// @Tags Tool运维工具箱
// @Summary 获取服务详情
// @Description 根据服务ID获取服务的详细信息
// @Produce json
// @Param serviceId path string true "服务ID"
// @Success 200 {object} result.Result
// @Router /api/v1/tool/services/{serviceId} [get]
// @Security ApiKeyAuth
func GetServiceDetail(c *gin.Context) {
	serviceID := c.Param("serviceId")
	service.ServiceDeployService().GetServiceDetail(c, serviceID)
}

// CreateDeploy 创建部署任务
// @Tags Tool运维工具箱
// @Summary 创建部署任务
// @Description 创建一个服务部署任务，异步执行部署
// @Accept json
// @Produce json
// @Param data body model.CreateDeployDto true "部署参数"
// @Success 200 {object} result.Result
// @Router /api/v1/tool/deploy [post]
// @Security ApiKeyAuth
func CreateDeploy(c *gin.Context) {
	var dto model.CreateDeployDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	service.ServiceDeployService().CreateDeploy(c, dto)
}

// GetDeployList 获取部署历史列表
// @Tags Tool运维工具箱
// @Summary 获取部署历史列表
// @Description 获取服务部署历史记录列表（分页）
// @Produce json
// @Param serviceName query string false "服务名称（模糊查询）"
// @Param hostId query int false "主机ID"
// @Param status query int false "状态: 0=部署中, 1=运行中, 2=已停止, 3=部署失败"
// @Param pageNum query int false "页码"
// @Param pageSize query int false "每页数量"
// @Success 200 {object} result.Result
// @Router /api/v1/tool/deploy/list [get]
// @Security ApiKeyAuth
func GetDeployList(c *gin.Context) {
	var dto model.DeployQueryDto
	if err := c.ShouldBindQuery(&dto); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	service.ServiceDeployService().GetDeployList(c, dto)
}

// GetDeployStatus 获取部署状态
// @Tags Tool运维工具箱
// @Summary 获取部署状态
// @Description 根据部署ID获取部署状态和日志
// @Produce json
// @Param id path int true "部署ID"
// @Success 200 {object} result.Result
// @Router /api/v1/tool/deploy/{id}/status [get]
// @Security ApiKeyAuth
func GetDeployStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"code": 400, "message": "ID格式错误"})
		return
	}
	service.ServiceDeployService().GetDeployStatus(c, uint(id))
}

// DeleteDeploy 删除部署记录（卸载服务）
// @Tags Tool运维工具箱
// @Summary 卸载服务
// @Description 停止并删除已部署的服务
// @Produce json
// @Param id path int true "部署ID"
// @Success 200 {object} result.Result
// @Router /api/v1/tool/deploy/{id} [delete]
// @Security ApiKeyAuth
func DeleteDeploy(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"code": 400, "message": "ID格式错误"})
		return
	}
	service.ServiceDeployService().DeleteDeploy(c, uint(id))
}
