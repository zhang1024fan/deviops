package controller

import (
	"strconv"

	"dodevops-api/api/k8s/model"
	"dodevops-api/api/k8s/service"
	"dodevops-api/common/result"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// K8sIngressController K8s Ingress控制器
type K8sIngressController struct {
	service service.IK8sIngressService
}

func NewK8sIngressController(db *gorm.DB) *K8sIngressController {
	return &K8sIngressController{
		service: service.NewK8sIngressService(db),
	}
}

// GetIngresses 获取Ingress列表
// @Summary 获取Ingress列表
// @Description 获取指定命名空间的Ingress列表
// @Tags K8s Ingress管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Success 200 {object} result.Result{data=model.IngressListResponse}
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/ingresses [get]
func (ctrl *K8sIngressController) GetIngresses(c *gin.Context) {
	clusterIdStr := c.Param("id")
	clusterId, err := strconv.Atoi(clusterIdStr)
	if err != nil || clusterId <= 0 {
		result.Failed(c, 400, "无效的集群ID")
		return
	}

	namespaceName := c.Param("namespaceName")
	if namespaceName == "" {
		result.Failed(c, 400, "命名空间名称不能为空")
		return
	}

	ctrl.service.GetIngresses(c, uint(clusterId), namespaceName)
}

// GetIngressDetail 获取Ingress详情
// @Summary 获取Ingress详情
// @Description 获取指定Ingress的详细信息
// @Tags K8s Ingress管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param ingressName path string true "Ingress名称"
// @Success 200 {object} result.Result{data=model.IngressDetail}
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 404 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/ingresses/{ingressName} [get]
func (ctrl *K8sIngressController) GetIngressDetail(c *gin.Context) {
	clusterIdStr := c.Param("id")
	clusterId, err := strconv.Atoi(clusterIdStr)
	if err != nil || clusterId <= 0 {
		result.Failed(c, 400, "无效的集群ID")
		return
	}

	namespaceName := c.Param("namespaceName")
	if namespaceName == "" {
		result.Failed(c, 400, "命名空间名称不能为空")
		return
	}

	ingressName := c.Param("ingressName")
	if ingressName == "" {
		result.Failed(c, 400, "Ingress名称不能为空")
		return
	}

	ctrl.service.GetIngressDetail(c, uint(clusterId), namespaceName, ingressName)
}

// CreateIngress 创建Ingress
// @Summary 创建Ingress
// @Description 在指定命名空间中创建新的Ingress
// @Tags K8s Ingress管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param ingress body model.CreateIngressRequest true "Ingress配置"
// @Success 200 {object} result.Result{data=model.K8sIngress}
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/ingresses [post]
func (ctrl *K8sIngressController) CreateIngress(c *gin.Context) {
	clusterIdStr := c.Param("id")
	clusterId, err := strconv.Atoi(clusterIdStr)
	if err != nil || clusterId <= 0 {
		result.Failed(c, 400, "无效的集群ID")
		return
	}

	namespaceName := c.Param("namespaceName")
	if namespaceName == "" {
		result.Failed(c, 400, "命名空间名称不能为空")
		return
	}

	var req model.CreateIngressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Failed(c, 400, "参数验证失败: "+err.Error())
		return
	}

	ctrl.service.CreateIngress(c, uint(clusterId), namespaceName, &req)
}

// UpdateIngress 更新Ingress
// @Summary 更新Ingress
// @Description 更新指定的Ingress配置
// @Tags K8s Ingress管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param ingressName path string true "Ingress名称"
// @Param ingress body model.UpdateIngressRequest true "更新配置"
// @Success 200 {object} result.Result{data=model.K8sIngress}
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 404 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/ingresses/{ingressName} [put]
func (ctrl *K8sIngressController) UpdateIngress(c *gin.Context) {
	clusterIdStr := c.Param("id")
	clusterId, err := strconv.Atoi(clusterIdStr)
	if err != nil || clusterId <= 0 {
		result.Failed(c, 400, "无效的集群ID")
		return
	}

	namespaceName := c.Param("namespaceName")
	if namespaceName == "" {
		result.Failed(c, 400, "命名空间名称不能为空")
		return
	}

	ingressName := c.Param("ingressName")
	if ingressName == "" {
		result.Failed(c, 400, "Ingress名称不能为空")
		return
	}

	var req model.UpdateIngressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Failed(c, 400, "参数验证失败: "+err.Error())
		return
	}

	ctrl.service.UpdateIngress(c, uint(clusterId), namespaceName, ingressName, &req)
}

// DeleteIngress 删除Ingress
// @Summary 删除Ingress
// @Description 删除指定的Ingress
// @Tags K8s Ingress管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param ingressName path string true "Ingress名称"
// @Success 200 {object} result.Result{data=string}
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 404 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/ingresses/{ingressName} [delete]
func (ctrl *K8sIngressController) DeleteIngress(c *gin.Context) {
	clusterIdStr := c.Param("id")
	clusterId, err := strconv.Atoi(clusterIdStr)
	if err != nil || clusterId <= 0 {
		result.Failed(c, 400, "无效的集群ID")
		return
	}

	namespaceName := c.Param("namespaceName")
	if namespaceName == "" {
		result.Failed(c, 400, "命名空间名称不能为空")
		return
	}

	ingressName := c.Param("ingressName")
	if ingressName == "" {
		result.Failed(c, 400, "Ingress名称不能为空")
		return
	}

	ctrl.service.DeleteIngress(c, uint(clusterId), namespaceName, ingressName)
}

// GetIngressYaml 获取Ingress的YAML配置
// @Summary 获取Ingress的YAML配置
// @Description 获取指定Ingress的完整YAML配置
// @Tags K8s Ingress管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param ingressName path string true "Ingress名称"
// @Success 200 {object} result.Result{data=map[string]interface{}}
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 404 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/ingresses/{ingressName}/yaml [get]
func (ctrl *K8sIngressController) GetIngressYaml(c *gin.Context) {
	clusterIdStr := c.Param("id")
	clusterId, err := strconv.Atoi(clusterIdStr)
	if err != nil || clusterId <= 0 {
		result.Failed(c, 400, "无效的集群ID")
		return
	}

	namespaceName := c.Param("namespaceName")
	if namespaceName == "" {
		result.Failed(c, 400, "命名空间名称不能为空")
		return
	}

	ingressName := c.Param("ingressName")
	if ingressName == "" {
		result.Failed(c, 400, "Ingress名称不能为空")
		return
	}

	ctrl.service.GetIngressYaml(c, uint(clusterId), namespaceName, ingressName)
}

// UpdateIngressYaml 通过YAML更新Ingress
// @Summary 通过YAML更新Ingress
// @Description 通过提供的YAML内容更新Ingress配置
// @Tags K8s Ingress管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param ingressName path string true "Ingress名称"
// @Param yaml body map[string]interface{} true "YAML内容"
// @Success 200 {object} result.Result{data=string}
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 404 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/ingresses/{ingressName}/yaml [put]
func (ctrl *K8sIngressController) UpdateIngressYaml(c *gin.Context) {
	clusterIdStr := c.Param("id")
	clusterId, err := strconv.Atoi(clusterIdStr)
	if err != nil || clusterId <= 0 {
		result.Failed(c, 400, "无效的集群ID")
		return
	}

	namespaceName := c.Param("namespaceName")
	if namespaceName == "" {
		result.Failed(c, 400, "命名空间名称不能为空")
		return
	}

	ingressName := c.Param("ingressName")
	if ingressName == "" {
		result.Failed(c, 400, "Ingress名称不能为空")
		return
	}

	var yamlData map[string]interface{}
	if err := c.ShouldBindJSON(&yamlData); err != nil {
		result.Failed(c, 400, "YAML内容格式错误: "+err.Error())
		return
	}

	ctrl.service.UpdateIngressYaml(c, uint(clusterId), namespaceName, ingressName, yamlData)
}

// GetIngressEvents 获取Ingress事件
// @Summary 获取Ingress事件
// @Description 获取指定Ingress的相关事件列表
// @Tags K8s Ingress管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param ingressName path string true "Ingress名称"
// @Success 200 {object} result.Result{data=[]model.K8sEvent}
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 404 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/ingresses/{ingressName}/events [get]
func (ctrl *K8sIngressController) GetIngressEvents(c *gin.Context) {
	clusterIdStr := c.Param("id")
	clusterId, err := strconv.Atoi(clusterIdStr)
	if err != nil || clusterId <= 0 {
		result.Failed(c, 400, "无效的集群ID")
		return
	}

	namespaceName := c.Param("namespaceName")
	if namespaceName == "" {
		result.Failed(c, 400, "命名空间名称不能为空")
		return
	}

	ingressName := c.Param("ingressName")
	if ingressName == "" {
		result.Failed(c, 400, "Ingress名称不能为空")
		return
	}

	ctrl.service.GetIngressEvents(c, uint(clusterId), namespaceName, ingressName)
}

// GetIngressMonitoring 获取Ingress监控信息
// @Summary 获取Ingress监控信息
// @Description 获取指定Ingress的监控指标和状态信息
// @Tags K8s Ingress管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param ingressName path string true "Ingress名称"
// @Success 200 {object} result.Result{data=map[string]interface{}}
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 404 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/ingresses/{ingressName}/monitoring [get]
func (ctrl *K8sIngressController) GetIngressMonitoring(c *gin.Context) {
	clusterIdStr := c.Param("id")
	clusterId, err := strconv.Atoi(clusterIdStr)
	if err != nil || clusterId <= 0 {
		result.Failed(c, 400, "无效的集群ID")
		return
	}

	namespaceName := c.Param("namespaceName")
	if namespaceName == "" {
		result.Failed(c, 400, "命名空间名称不能为空")
		return
	}

	ingressName := c.Param("ingressName")
	if ingressName == "" {
		result.Failed(c, 400, "Ingress名称不能为空")
		return
	}

	ctrl.service.GetIngressMonitoring(c, uint(clusterId), namespaceName, ingressName)
}