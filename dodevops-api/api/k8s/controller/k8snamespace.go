package controller

import (
	"net/http"
	"strconv"

	"dodevops-api/api/k8s/service"
	"dodevops-api/api/k8s/model"
	"dodevops-api/common/result"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type K8sNamespaceController struct {
	service service.IK8sNamespaceService
}

func NewK8sNamespaceController(db *gorm.DB) *K8sNamespaceController {
	return &K8sNamespaceController{
		service: service.NewK8sNamespaceService(db),
	}
}

func NewK8sNamespaceControllerWithCache(db *gorm.DB, redisClient *redis.Client) *K8sNamespaceController {
	return &K8sNamespaceController{
		service: service.NewK8sNamespaceServiceWithCache(db, redisClient),
	}
}

// GetNamespaces 获取命名空间列表
// @Summary 获取K8s命名空间列表
// @Description 获取指定集群的所有命名空间信息
// @Tags K8s命名空间管理
// @Accept json
// @Produce json
// @Param id path uint true "集群ID"
// @Success 200 {object} result.Result{data=model.NamespaceListResponse} "获取成功"
// @Failure 400 {object} result.Result "参数错误"
// @Failure 404 {object} result.Result "集群不存在"
// @Failure 500 {object} result.Result "服务器错误"
// @Router /api/v1/k8s/cluster/{id}/namespaces [get]
// @Security ApiKeyAuth
func (ctrl *K8sNamespaceController) GetNamespaces(c *gin.Context) {
	clusterIdStr := c.Param("id")
	clusterId, err := strconv.ParseUint(clusterIdStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的集群ID")
		return
	}

	ctrl.service.GetNamespaces(c, uint(clusterId))
}

// GetNamespace 获取单个命名空间详情
// @Summary 获取K8s命名空间详情
// @Description 获取指定命名空间的详细信息
// @Tags K8s命名空间管理
// @Accept json
// @Produce json
// @Param id path uint true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Success 200 {object} result.Result{data=model.K8sNamespace} "获取成功"
// @Failure 400 {object} result.Result "参数错误"
// @Failure 404 {object} result.Result "命名空间不存在"
// @Failure 500 {object} result.Result "服务器错误"
// @Router /api/v1/k8s/cluster/{id}/namespaces/{namespaceName} [get]
// @Security ApiKeyAuth
func (ctrl *K8sNamespaceController) GetNamespace(c *gin.Context) {
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

	ctrl.service.GetNamespace(c, uint(clusterId), namespaceName)
}

// CreateNamespace 创建命名空间
// @Summary 创建K8s命名空间
// @Description 在指定集群中创建新的命名空间
// @Tags K8s命名空间管理
// @Accept json
// @Produce json
// @Param id path uint true "集群ID"
// @Param namespace body model.CreateNamespaceRequest true "命名空间信息"
// @Success 200 {object} result.Result{data=model.K8sNamespace} "创建成功"
// @Failure 400 {object} result.Result "参数错误"
// @Failure 404 {object} result.Result "集群不存在"
// @Failure 500 {object} result.Result "服务器错误"
// @Router /api/v1/k8s/cluster/{id}/namespaces [post]
// @Security ApiKeyAuth
func (ctrl *K8sNamespaceController) CreateNamespace(c *gin.Context) {
	clusterIdStr := c.Param("id")
	clusterId, err := strconv.ParseUint(clusterIdStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的集群ID")
		return
	}

	var req model.CreateNamespaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Failed(c, http.StatusBadRequest, "参数验证失败: "+err.Error())
		return
	}

	ctrl.service.CreateNamespace(c, uint(clusterId), &req)
}

// UpdateNamespace 更新命名空间
// @Summary 更新K8s命名空间
// @Description 更新指定命名空间的标签和注释
// @Tags K8s命名空间管理
// @Accept json
// @Produce json
// @Param id path uint true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param namespace body model.UpdateNamespaceRequest true "更新信息"
// @Success 200 {object} result.Result{data=model.K8sNamespace} "更新成功"
// @Failure 400 {object} result.Result "参数错误"
// @Failure 404 {object} result.Result "命名空间不存在"
// @Failure 500 {object} result.Result "服务器错误"
// @Router /api/v1/k8s/cluster/{id}/namespaces/{namespaceName} [put]
// @Security ApiKeyAuth
func (ctrl *K8sNamespaceController) UpdateNamespace(c *gin.Context) {
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

	var req model.UpdateNamespaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Failed(c, http.StatusBadRequest, "参数验证失败: "+err.Error())
		return
	}

	ctrl.service.UpdateNamespace(c, uint(clusterId), namespaceName, &req)
}

// DeleteNamespace 删除命名空间
// @Summary 删除K8s命名空间
// @Description 删除指定的命名空间（会同时删除其中的所有资源）
// @Tags K8s命名空间管理
// @Accept json
// @Produce json
// @Param id path uint true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Success 200 {object} result.Result "删除成功"
// @Failure 400 {object} result.Result "参数错误"
// @Failure 404 {object} result.Result "命名空间不存在"
// @Failure 500 {object} result.Result "服务器错误"
// @Router /api/v1/k8s/cluster/{id}/namespaces/{namespaceName} [delete]
// @Security ApiKeyAuth
func (ctrl *K8sNamespaceController) DeleteNamespace(c *gin.Context) {
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

	// 检查是否为系统命名空间
	systemNamespaces := []string{"default", "kube-system", "kube-public", "kube-node-lease"}
	for _, sysNs := range systemNamespaces {
		if namespaceName == sysNs {
			result.Failed(c, http.StatusBadRequest, "不能删除系统命名空间: "+namespaceName)
			return
		}
	}

	ctrl.service.DeleteNamespace(c, uint(clusterId), namespaceName)
}

// GetResourceQuotas 获取命名空间的ResourceQuota列表
// @Summary 获取资源配额列表
// @Description 获取指定命名空间的所有资源配额
// @Tags K8s命名空间管理
// @Accept json
// @Produce json
// @Param id path uint true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Success 200 {object} result.Result{data=model.ResourceQuotaListResponse} "获取成功"
// @Failure 400 {object} result.Result "参数错误"
// @Failure 404 {object} result.Result "命名空间不存在"
// @Failure 500 {object} result.Result "服务器错误"
// @Router /api/v1/k8s/cluster/{id}/namespaces/{namespaceName}/resourcequotas [get]
// @Security ApiKeyAuth
func (ctrl *K8sNamespaceController) GetResourceQuotas(c *gin.Context) {
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

	ctrl.service.GetResourceQuotas(c, uint(clusterId), namespaceName)
}

// CreateResourceQuota 创建ResourceQuota
// @Summary 创建资源配额
// @Description 为指定命名空间创建资源配额
// @Tags K8s命名空间管理
// @Accept json
// @Produce json
// @Param id path uint true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param quota body model.CreateResourceQuotaRequest true "资源配额信息"
// @Success 200 {object} result.Result{data=model.ResourceQuotaDetail} "创建成功"
// @Failure 400 {object} result.Result "参数错误"
// @Failure 404 {object} result.Result "命名空间不存在"
// @Failure 500 {object} result.Result "服务器错误"
// @Router /api/v1/k8s/cluster/{id}/namespaces/{namespaceName}/resourcequotas [post]
// @Security ApiKeyAuth
func (ctrl *K8sNamespaceController) CreateResourceQuota(c *gin.Context) {
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

	var req model.CreateResourceQuotaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Failed(c, http.StatusBadRequest, "参数验证失败: "+err.Error())
		return
	}

	ctrl.service.CreateResourceQuota(c, uint(clusterId), namespaceName, &req)
}

// UpdateResourceQuota 更新ResourceQuota
// @Summary 更新资源配额
// @Description 更新指定的资源配额
// @Tags K8s命名空间管理
// @Accept json
// @Produce json
// @Param id path uint true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param quotaName path string true "ResourceQuota名称"
// @Param quota body model.UpdateResourceQuotaRequest true "更新信息"
// @Success 200 {object} result.Result{data=model.ResourceQuotaDetail} "更新成功"
// @Failure 400 {object} result.Result "参数错误"
// @Failure 404 {object} result.Result "资源配额不存在"
// @Failure 500 {object} result.Result "服务器错误"
// @Router /api/v1/k8s/cluster/{id}/namespaces/{namespaceName}/resourcequotas/{quotaName} [put]
// @Security ApiKeyAuth
func (ctrl *K8sNamespaceController) UpdateResourceQuota(c *gin.Context) {
	clusterIdStr := c.Param("id")
	clusterId, err := strconv.ParseUint(clusterIdStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的集群ID")
		return
	}

	namespaceName := c.Param("namespaceName")
	quotaName := c.Param("quotaName")
	if namespaceName == "" || quotaName == "" {
		result.Failed(c, http.StatusBadRequest, "命名空间名称和配额名称不能为空")
		return
	}

	var req model.UpdateResourceQuotaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Failed(c, http.StatusBadRequest, "参数验证失败: "+err.Error())
		return
	}

	ctrl.service.UpdateResourceQuota(c, uint(clusterId), namespaceName, quotaName, &req)
}

// DeleteResourceQuota 删除ResourceQuota
// @Summary 删除资源配额
// @Description 删除指定的资源配额
// @Tags K8s命名空间管理
// @Accept json
// @Produce json
// @Param id path uint true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param quotaName path string true "ResourceQuota名称"
// @Success 200 {object} result.Result "删除成功"
// @Failure 400 {object} result.Result "参数错误"
// @Failure 404 {object} result.Result "资源配额不存在"
// @Failure 500 {object} result.Result "服务器错误"
// @Router /api/v1/k8s/cluster/{id}/namespaces/{namespaceName}/resourcequotas/{quotaName} [delete]
// @Security ApiKeyAuth
func (ctrl *K8sNamespaceController) DeleteResourceQuota(c *gin.Context) {
	clusterIdStr := c.Param("id")
	clusterId, err := strconv.ParseUint(clusterIdStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的集群ID")
		return
	}

	namespaceName := c.Param("namespaceName")
	quotaName := c.Param("quotaName")
	if namespaceName == "" || quotaName == "" {
		result.Failed(c, http.StatusBadRequest, "命名空间名称和配额名称不能为空")
		return
	}

	ctrl.service.DeleteResourceQuota(c, uint(clusterId), namespaceName, quotaName)
}

// GetLimitRanges 获取命名空间的LimitRange列表
// @Summary 获取默认资源限制列表
// @Description 获取指定命名空间的所有默认资源限制
// @Tags K8s命名空间管理
// @Accept json
// @Produce json
// @Param id path uint true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Success 200 {object} result.Result{data=model.LimitRangeListResponse} "获取成功"
// @Failure 400 {object} result.Result "参数错误"
// @Failure 404 {object} result.Result "命名空间不存在"
// @Failure 500 {object} result.Result "服务器错误"
// @Router /api/v1/k8s/cluster/{id}/namespaces/{namespaceName}/limitranges [get]
// @Security ApiKeyAuth
func (ctrl *K8sNamespaceController) GetLimitRanges(c *gin.Context) {
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

	ctrl.service.GetLimitRanges(c, uint(clusterId), namespaceName)
}

// CreateLimitRange 创建LimitRange
// @Summary 创建默认资源限制
// @Description 为指定命名空间创建默认资源限制
// @Tags K8s命名空间管理
// @Accept json
// @Produce json
// @Param id path uint true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param limitrange body model.CreateLimitRangeRequest true "默认资源限制信息"
// @Success 200 {object} result.Result{data=model.LimitRangeDetail} "创建成功"
// @Failure 400 {object} result.Result "参数错误"
// @Failure 404 {object} result.Result "命名空间不存在"
// @Failure 500 {object} result.Result "服务器错误"
// @Router /api/v1/k8s/cluster/{id}/namespaces/{namespaceName}/limitranges [post]
// @Security ApiKeyAuth
func (ctrl *K8sNamespaceController) CreateLimitRange(c *gin.Context) {
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

	var req model.CreateLimitRangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Failed(c, http.StatusBadRequest, "参数验证失败: "+err.Error())
		return
	}

	ctrl.service.CreateLimitRange(c, uint(clusterId), namespaceName, &req)
}

// UpdateLimitRange 更新LimitRange
// @Summary 更新默认资源限制
// @Description 更新指定的默认资源限制
// @Tags K8s命名空间管理
// @Accept json
// @Produce json
// @Param id path uint true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param limitRangeName path string true "LimitRange名称"
// @Param limitrange body model.UpdateLimitRangeRequest true "更新信息"
// @Success 200 {object} result.Result{data=model.LimitRangeDetail} "更新成功"
// @Failure 400 {object} result.Result "参数错误"
// @Failure 404 {object} result.Result "默认资源限制不存在"
// @Failure 500 {object} result.Result "服务器错误"
// @Router /api/v1/k8s/cluster/{id}/namespaces/{namespaceName}/limitranges/{limitRangeName} [put]
// @Security ApiKeyAuth
func (ctrl *K8sNamespaceController) UpdateLimitRange(c *gin.Context) {
	clusterIdStr := c.Param("id")
	clusterId, err := strconv.ParseUint(clusterIdStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的集群ID")
		return
	}

	namespaceName := c.Param("namespaceName")
	limitRangeName := c.Param("limitRangeName")
	if namespaceName == "" || limitRangeName == "" {
		result.Failed(c, http.StatusBadRequest, "命名空间名称和限制范围名称不能为空")
		return
	}

	var req model.UpdateLimitRangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Failed(c, http.StatusBadRequest, "参数验证失败: "+err.Error())
		return
	}

	ctrl.service.UpdateLimitRange(c, uint(clusterId), namespaceName, limitRangeName, &req)
}

// DeleteLimitRange 删除LimitRange
// @Summary 删除默认资源限制
// @Description 删除指定的默认资源限制
// @Tags K8s命名空间管理
// @Accept json
// @Produce json
// @Param id path uint true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param limitRangeName path string true "LimitRange名称"
// @Success 200 {object} result.Result "删除成功"
// @Failure 400 {object} result.Result "参数错误"
// @Failure 404 {object} result.Result "默认资源限制不存在"
// @Failure 500 {object} result.Result "服务器错误"
// @Router /api/v1/k8s/cluster/{id}/namespaces/{namespaceName}/limitranges/{limitRangeName} [delete]
// @Security ApiKeyAuth
func (ctrl *K8sNamespaceController) DeleteLimitRange(c *gin.Context) {
	clusterIdStr := c.Param("id")
	clusterId, err := strconv.ParseUint(clusterIdStr, 10, 32)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的集群ID")
		return
	}

	namespaceName := c.Param("namespaceName")
	limitRangeName := c.Param("limitRangeName")
	if namespaceName == "" || limitRangeName == "" {
		result.Failed(c, http.StatusBadRequest, "命名空间名称和限制范围名称不能为空")
		return
	}

	ctrl.service.DeleteLimitRange(c, uint(clusterId), namespaceName, limitRangeName)
}