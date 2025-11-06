package controller

import (
	"net/http"
	"strconv"

	"dodevops-api/api/k8s/service"
	"dodevops-api/api/k8s/model"
	"dodevops-api/common/result"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type KubeClusterController struct {
	service service.IKubeClusterService
}

func NewKubeClusterController(db *gorm.DB) *KubeClusterController {
	return &KubeClusterController{
		service: service.NewKubeClusterService(db),
	}
}

// CreateCluster 创建K8s集群
// @Summary 创建K8s集群
// @Description 创建K8s集群，可选择是否自动部署
// @Tags K8s集群管理
// @Accept json
// @Produce json
// @Param cluster body model.CreateKubeClusterRequest true "集群创建参数"
// @Success 200 {object} result.Result{data=model.KubeCluster} "创建成功"
// @Failure 400 {object} result.Result "参数错误"
// @Failure 500 {object} result.Result "服务器错误"
// @Router /api/v1/k8s/cluster [post]
// @Security ApiKeyAuth
func (ctrl *KubeClusterController) CreateCluster(c *gin.Context) {
	var req model.CreateKubeClusterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Failed(c, http.StatusBadRequest, "参数验证失败: "+err.Error())
		return
	}

	ctrl.service.CreateCluster(c, &req)
}

// GetCluster 获取集群详情
// @Summary 获取K8s集群详情
// @Description 根据集群ID获取集群详细信息
// @Tags K8s集群管理
// @Accept json
// @Produce json
// @Param id path uint true "集群ID"
// @Success 200 {object} result.Result{data=model.KubeCluster} "获取成功"
// @Failure 400 {object} result.Result "参数错误"
// @Failure 404 {object} result.Result "集群不存在"
// @Failure 500 {object} result.Result "服务器错误"
// @Router /api/v1/k8s/cluster/{id} [get]
// @Security ApiKeyAuth
func (ctrl *KubeClusterController) GetCluster(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的集群ID")
		return
	}

	ctrl.service.GetCluster(c, uint(id))
}

// GetClusterList 获取集群列表
// @Summary 获取K8s集群列表
// @Description 分页获取K8s集群列表
// @Tags K8s集群管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Success 200 {object} result.Result{data=model.KubeClusterListResponse} "获取成功"
// @Failure 500 {object} result.Result "服务器错误"
// @Router /api/v1/k8s/cluster [get]
// @Security ApiKeyAuth
func (ctrl *KubeClusterController) GetClusterList(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	sizeStr := c.DefaultQuery("size", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	size, err := strconv.Atoi(sizeStr)
	if err != nil || size < 1 || size > 100 {
		size = 10
	}

	ctrl.service.GetClusterList(c, page, size)
}

// DeleteCluster 删除集群
// @Summary 删除K8s集群
// @Description 删除指定的K8s集群（只能删除已停止的集群）
// @Tags K8s集群管理
// @Accept json
// @Produce json
// @Param id path uint true "集群ID"
// @Success 200 {object} result.Result "删除成功"
// @Failure 400 {object} result.Result "参数错误或集群状态不允许删除"
// @Failure 404 {object} result.Result "集群不存在"
// @Failure 500 {object} result.Result "服务器错误"
// @Router /api/v1/k8s/cluster/{id} [delete]
// @Security ApiKeyAuth
func (ctrl *KubeClusterController) DeleteCluster(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的集群ID")
		return
	}

	ctrl.service.DeleteCluster(c, uint(id))
}

// GetClusterStatus 获取集群状态
// @Summary 获取K8s集群状态
// @Description 获取集群运行状态和节点信息
// @Tags K8s集群管理
// @Accept json
// @Produce json
// @Param id path uint true "集群ID"
// @Success 200 {object} result.Result{data=map[string]interface{}} "获取成功"
// @Failure 400 {object} result.Result "参数错误"
// @Failure 404 {object} result.Result "集群不存在"
// @Failure 500 {object} result.Result "服务器错误"
// @Router /api/v1/k8s/cluster/{id}/status [get]
// @Security ApiKeyAuth
func (ctrl *KubeClusterController) GetClusterStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的集群ID")
		return
	}

	ctrl.service.GetClusterStatus(c, uint(id))
}

// UpdateCluster 更新集群信息
// @Summary 更新K8s集群信息
// @Description 更新集群的基本信息（名称、描述、kubeconfig等）
// @Tags K8s集群管理
// @Accept json
// @Produce json
// @Param id path uint true "集群ID"
// @Param cluster body model.UpdateKubeClusterRequest true "集群更新参数"
// @Success 200 {object} result.Result{data=model.KubeCluster} "更新成功"
// @Failure 400 {object} result.Result "参数错误"
// @Failure 404 {object} result.Result "集群不存在"
// @Failure 500 {object} result.Result "服务器错误"
// @Router /api/v1/k8s/cluster/{id} [put]
// @Security ApiKeyAuth
func (ctrl *KubeClusterController) UpdateCluster(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的集群ID")
		return
	}

	var req model.UpdateKubeClusterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Failed(c, http.StatusBadRequest, "参数验证失败: "+err.Error())
		return
	}

	ctrl.service.UpdateCluster(c, uint(id), &req)
}

// SyncCluster 同步集群信息
// @Summary 同步K8s集群信息
// @Description 通过K8s API同步集群版本、节点数量、集群状态等信息并更新数据库
// @Tags K8s集群管理
// @Accept json
// @Produce json
// @Param id path uint true "集群ID"
// @Success 200 {object} result.Result{data=model.KubeCluster} "同步成功"
// @Failure 400 {object} result.Result "参数错误"
// @Failure 404 {object} result.Result "集群不存在"
// @Failure 500 {object} result.Result "服务器错误"
// @Router /api/v1/k8s/cluster/{id}/sync [post]
// @Security ApiKeyAuth
func (ctrl *KubeClusterController) SyncCluster(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的集群ID")
		return
	}

	ctrl.service.SyncCluster(c, uint(id))
}

// GetClusterDetail 获取集群详细信息
// @Summary 获取K8s集群详细信息
// @Description 获取集群的完整详细信息，包括节点、工作负载、组件、网络配置、监控信息等
// @Tags K8s集群管理
// @Accept json
// @Produce json
// @Param id path uint true "集群ID"
// @Success 200 {object} result.Result{data=model.ClusterDetailResponse} "获取成功"
// @Failure 400 {object} result.Result "参数错误"
// @Failure 404 {object} result.Result "集群不存在"
// @Failure 500 {object} result.Result "服务器错误"
// @Router /api/v1/k8s/cluster/{id}/detail [get]
// @Security ApiKeyAuth
func (ctrl *KubeClusterController) GetClusterDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的集群ID")
		return
	}

	ctrl.service.GetClusterDetail(c, uint(id))
}