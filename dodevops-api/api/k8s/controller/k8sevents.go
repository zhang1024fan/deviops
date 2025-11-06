package controller

import (
	"net/http"
	"strconv"

	"dodevops-api/api/k8s/service"
	"dodevops-api/common/result"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// K8sEventsController K8s事件控制器
type K8sEventsController struct {
	eventsService *service.K8sEventsServiceImpl
}

func NewK8sEventsController(db *gorm.DB) *K8sEventsController {
	return &K8sEventsController{
		eventsService: service.NewK8sEventsServiceImpl(db),
	}
}

// GetEvents 获取指定命名空间的事件列表
// @Summary 获取K8s事件列表
// @Description 获取指定命名空间的事件列表，支持按资源类型和名称过滤
// @Tags K8s事件管理
// @Accept json
// @Produce json
// @Param id path uint true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param kind query string false "资源类型" Enums(Pod,Deployment,StatefulSet,DaemonSet,Service)
// @Param name query string false "资源名称"
// @Param limit query int false "限制返回数量" default(100)
// @Success 200 {object} result.Result{data=model.EventListResponse} "获取成功"
// @Failure 400 {object} result.Result "参数错误"
// @Failure 404 {object} result.Result "集群不存在"
// @Failure 500 {object} result.Result "服务器错误"
// @Router /api/v1/k8s/cluster/{id}/namespaces/{namespaceName}/events [get]
// @Security ApiKeyAuth
func (ctrl *K8sEventsController) GetEvents(c *gin.Context) {
	clusterIdStr := c.Param("id")
	clusterId, err := strconv.ParseUint(clusterIdStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的集群ID")
		return
	}

	namespaceName := c.Param("namespaceName")
	if namespaceName == "" {
		result.Failed(c, http.StatusBadRequest, "命名空间名称不能为空")
		return
	}

	// 获取查询参数
	kind := c.Query("kind")       // 资源类型过滤
	name := c.Query("name")       // 资源名称过滤
	limitStr := c.Query("limit")  // 限制返回数量
	
	limit := 100 // 默认返回100条
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	ctrl.eventsService.GetEvents(c, uint(clusterId), namespaceName, kind, name, limit)
}

// GetClusterEvents 获取整个集群的事件列表
// @Summary 获取集群事件列表
// @Description 获取整个集群的事件列表
// @Tags K8s事件管理
// @Accept json
// @Produce json
// @Param id path uint true "集群ID"
// @Param limit query int false "限制返回数量" default(100)
// @Success 200 {object} result.Result{data=model.EventListResponse} "获取成功"
// @Failure 400 {object} result.Result "参数错误"
// @Failure 404 {object} result.Result "集群不存在"
// @Failure 500 {object} result.Result "服务器错误"
// @Router /api/v1/k8s/cluster/{id}/events [get]
// @Security ApiKeyAuth
func (ctrl *K8sEventsController) GetClusterEvents(c *gin.Context) {
	clusterIdStr := c.Param("id")
	clusterId, err := strconv.ParseUint(clusterIdStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的集群ID")
		return
	}

	// 获取查询参数
	limitStr := c.Query("limit")
	limit := 100 // 默认返回100条
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	ctrl.eventsService.GetClusterEvents(c, uint(clusterId), limit)
}