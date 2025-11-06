package controller

import (
	"dodevops-api/api/dashboard/service"
	"dodevops-api/common/result"
	"dodevops-api/pkg/log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DashboardController struct {
	service service.IDashboardService
}

func NewDashboardController() *DashboardController {
	return &DashboardController{
		service: service.NewDashboardService(),
	}
}

// GetDashboardStats 获取看板统计数据
// @Summary 获取看板统计数据
// @Description 获取系统看板的各项统计数据，包括主机、K8s集群、发布、任务、服务和数据库统计
// @Tags 看板管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Success 200 {object} result.Result{data=model.DashboardStats}
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /dashboard/stats [get]
func (ctrl *DashboardController) GetDashboardStats(c *gin.Context) {
	logger := log.Log()

	// 获取看板统计数据
	stats, err := ctrl.service.GetDashboardStats()
	if err != nil {
		logger.Errorf("获取看板统计数据失败: %v", err)
		result.Failed(c, http.StatusInternalServerError, "获取看板统计数据失败")
		return
	}

	result.Success(c, stats)
}

// GetBusinessDistributionStats 获取业务分布统计数据
// @Summary 获取业务分布统计数据
// @Description 获取各业务线的服务数量分布，包括总服务数量和各业务线的服务占比
// @Tags 看板管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Success 200 {object} result.Result{data=model.BusinessDistributionStats}
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /dashboard/business-distribution [get]
func (ctrl *DashboardController) GetBusinessDistributionStats(c *gin.Context) {
	logger := log.Log()

	// 获取业务分布统计数据
	stats, err := ctrl.service.GetBusinessDistributionStats()
	if err != nil {
		logger.Errorf("获取业务分布统计数据失败: %v", err)
		result.Failed(c, http.StatusInternalServerError, "获取业务分布统计数据失败")
		return
	}

	result.Success(c, stats)
}

// GetAssetStats 获取资产统计数据
// @Summary 获取资产统计数据
// @Description 获取系统资产统计数据，包括主机(按云平台)、数据库(按类型)、K8s集群等资产分布
// @Tags 看板管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Success 200 {object} result.Result{data=model.AssetStats}
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /dashboard/assets [get]
func (ctrl *DashboardController) GetAssetStats(c *gin.Context) {
	logger := log.Log()

	// 获取资产统计数据
	stats, err := ctrl.service.GetAssetStats()
	if err != nil {
		logger.Errorf("获取资产统计数据失败: %v", err)
		result.Failed(c, http.StatusInternalServerError, "获取资产统计数据失败")
		return
	}

	result.Success(c, stats)
}