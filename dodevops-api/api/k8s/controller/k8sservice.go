package controller

import (
	"strconv"

	"dodevops-api/api/k8s/model"
	"dodevops-api/api/k8s/service"
	"dodevops-api/common/result"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// K8sServiceController K8s Service控制器
type K8sServiceController struct {
	service service.IK8sServiceService
}

func NewK8sServiceController(db *gorm.DB) *K8sServiceController {
	return &K8sServiceController{
		service: service.NewK8sServiceService(db),
	}
}

// GetServices 获取Service列表
// @Summary 获取Service列表
// @Description 获取指定命名空间的Service列表
// @Tags K8s Service管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Success 200 {object} result.Result{data=model.ServiceListResponse}
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/services [get]
func (ctrl *K8sServiceController) GetServices(c *gin.Context) {
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

	ctrl.service.GetServices(c, uint(clusterId), namespaceName)
}

// GetServiceDetail 获取Service详情
// @Summary 获取Service详情
// @Description 获取指定Service的详细信息
// @Tags K8s Service管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param serviceName path string true "Service名称"
// @Success 200 {object} result.Result{data=model.ServiceDetail}
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 404 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/services/{serviceName} [get]
func (ctrl *K8sServiceController) GetServiceDetail(c *gin.Context) {
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

	serviceName := c.Param("serviceName")
	if serviceName == "" {
		result.Failed(c, 400, "Service名称不能为空")
		return
	}

	ctrl.service.GetServiceDetail(c, uint(clusterId), namespaceName, serviceName)
}

// CreateService 创建Service
// @Summary 创建Service
// @Description 在指定命名空间中创建新的Service
// @Tags K8s Service管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param service body model.CreateServiceRequest true "Service配置"
// @Success 200 {object} result.Result{data=model.K8sService}
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/services [post]
func (ctrl *K8sServiceController) CreateService(c *gin.Context) {
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

	var req model.CreateServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Failed(c, 400, "参数验证失败: "+err.Error())
		return
	}

	ctrl.service.CreateService(c, uint(clusterId), namespaceName, &req)
}

// UpdateService 更新Service
// @Summary 更新Service
// @Description 更新指定的Service配置
// @Tags K8s Service管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param serviceName path string true "Service名称"
// @Param service body model.UpdateServiceRequest true "更新配置"
// @Success 200 {object} result.Result{data=model.K8sService}
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 404 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/services/{serviceName} [put]
func (ctrl *K8sServiceController) UpdateService(c *gin.Context) {
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

	serviceName := c.Param("serviceName")
	if serviceName == "" {
		result.Failed(c, 400, "Service名称不能为空")
		return
	}

	var req model.UpdateServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Failed(c, 400, "参数验证失败: "+err.Error())
		return
	}

	ctrl.service.UpdateService(c, uint(clusterId), namespaceName, serviceName, &req)
}

// DeleteService 删除Service
// @Summary 删除Service
// @Description 删除指定的Service
// @Tags K8s Service管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param serviceName path string true "Service名称"
// @Success 200 {object} result.Result{data=string}
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 404 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/services/{serviceName} [delete]
func (ctrl *K8sServiceController) DeleteService(c *gin.Context) {
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

	serviceName := c.Param("serviceName")
	if serviceName == "" {
		result.Failed(c, 400, "Service名称不能为空")
		return
	}

	ctrl.service.DeleteService(c, uint(clusterId), namespaceName, serviceName)
}

// GetServiceYaml 获取Service的YAML配置
// @Summary 获取Service的YAML配置
// @Description 获取指定Service的完整YAML配置
// @Tags K8s Service管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param serviceName path string true "Service名称"
// @Success 200 {object} result.Result{data=map[string]interface{}}
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 404 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/services/{serviceName}/yaml [get]
func (ctrl *K8sServiceController) GetServiceYaml(c *gin.Context) {
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

	serviceName := c.Param("serviceName")
	if serviceName == "" {
		result.Failed(c, 400, "Service名称不能为空")
		return
	}

	ctrl.service.GetServiceYaml(c, uint(clusterId), namespaceName, serviceName)
}

// UpdateServiceYaml 通过YAML更新Service
// @Summary 通过YAML更新Service
// @Description 通过提供的YAML内容更新Service配置
// @Tags K8s Service管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param serviceName path string true "Service名称"
// @Param yaml body map[string]interface{} true "YAML内容"
// @Success 200 {object} result.Result{data=string}
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 404 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/services/{serviceName}/yaml [put]
func (ctrl *K8sServiceController) UpdateServiceYaml(c *gin.Context) {
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

	serviceName := c.Param("serviceName")
	if serviceName == "" {
		result.Failed(c, 400, "Service名称不能为空")
		return
	}

	var yamlData map[string]interface{}
	if err := c.ShouldBindJSON(&yamlData); err != nil {
		result.Failed(c, 400, "YAML内容格式错误: "+err.Error())
		return
	}

	ctrl.service.UpdateServiceYaml(c, uint(clusterId), namespaceName, serviceName, yamlData)
}

// GetServiceEvents 获取Service事件
// @Summary 获取Service事件
// @Description 获取指定Service的相关事件列表
// @Tags K8s Service管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param serviceName path string true "Service名称"
// @Success 200 {object} result.Result{data=[]model.K8sEvent}
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 404 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/services/{serviceName}/events [get]
func (ctrl *K8sServiceController) GetServiceEvents(c *gin.Context) {
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

	serviceName := c.Param("serviceName")
	if serviceName == "" {
		result.Failed(c, 400, "Service名称不能为空")
		return
	}

	ctrl.service.GetServiceEvents(c, uint(clusterId), namespaceName, serviceName)
}