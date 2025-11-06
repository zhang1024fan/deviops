package controller

import (
	"strconv"

	"dodevops-api/api/k8s/model"
	"dodevops-api/api/k8s/service"
	"dodevops-api/common/result"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// K8sWorkloadController K8s工作负载控制器
type K8sWorkloadController struct {
	service service.IK8sWorkloadService
}

func NewK8sWorkloadController(db *gorm.DB) *K8sWorkloadController {
	return &K8sWorkloadController{
		service: service.NewK8sWorkloadService(db),
	}
}

// GetWorkloads 获取工作负载列表
// @Summary 获取工作负载列表
// @Description 获取指定命名空间的工作负载列表
// @Tags K8s工作负载管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param type query string false "工作负载类型" Enums(deployment,statefulset,daemonset,job,cronjob,all)
// @Success 200 {object} result.Result{data=model.WorkloadListResponse}
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/workloads [get]
func (ctrl *K8sWorkloadController) GetWorkloads(c *gin.Context) {
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

	workloadType := c.DefaultQuery("type", "all")

	ctrl.service.GetWorkloads(c, uint(clusterId), namespaceName, workloadType)
}

// GetWorkloadDetail 获取工作负载详情
// @Summary 获取工作负载详情
// @Description 获取指定工作负载的详细信息
// @Tags K8s工作负载管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param type path string true "工作负载类型" Enums(deployment,statefulset,daemonset,job,cronjob)
// @Param workloadName path string true "工作负载名称"
// @Success 200 {object} result.Result{data=model.K8sWorkloadDetail}
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 404 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/workloads/{type}/{workloadName} [get]
func (ctrl *K8sWorkloadController) GetWorkloadDetail(c *gin.Context) {
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

	workloadType := c.Param("type")
	if workloadType == "" {
		result.Failed(c, 400, "工作负载类型不能为空")
		return
	}

	workloadName := c.Param("workloadName")
	if workloadName == "" {
		result.Failed(c, 400, "工作负载名称不能为空")
		return
	}

	ctrl.service.GetWorkloadDetail(c, uint(clusterId), namespaceName, workloadType, workloadName)
}

// ===================== Deployment 管理 =====================

// CreateDeployment 创建Deployment
// @Summary 创建Deployment
// @Description 在指定命名空间中创建新的Deployment
// @Tags Deployment管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param deployment body model.CreateDeploymentRequest true "Deployment配置"
// @Success 200 {object} result.Result{data=model.K8sWorkload}
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/deployments [post]
func (ctrl *K8sWorkloadController) CreateDeployment(c *gin.Context) {
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

	var req model.CreateDeploymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Failed(c, 400, "参数验证失败: "+err.Error())
		return
	}

	ctrl.service.CreateDeployment(c, uint(clusterId), namespaceName, &req)
}

// UpdateDeployment 更新Deployment
// @Summary 更新Deployment
// @Description 更新指定的Deployment配置
// @Tags Deployment管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param deploymentName path string true "Deployment名称"
// @Param deployment body model.UpdateWorkloadRequest true "更新配置"
// @Success 200 {object} result.Result{data=model.K8sWorkload}
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 404 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/deployments/{deploymentName} [put]
func (ctrl *K8sWorkloadController) UpdateDeployment(c *gin.Context) {
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

	deploymentName := c.Param("deploymentName")
	if deploymentName == "" {
		result.Failed(c, 400, "Deployment名称不能为空")
		return
	}

	var req model.UpdateWorkloadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Failed(c, 400, "参数验证失败: "+err.Error())
		return
	}

	ctrl.service.UpdateDeployment(c, uint(clusterId), namespaceName, deploymentName, &req)
}

// DeleteDeployment 删除Deployment
// @Summary 删除Deployment
// @Description 删除指定的Deployment
// @Tags Deployment管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param deploymentName path string true "Deployment名称"
// @Success 200 {object} result.Result{data=string}
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 404 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/deployments/{deploymentName} [delete]
func (ctrl *K8sWorkloadController) DeleteDeployment(c *gin.Context) {
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

	deploymentName := c.Param("deploymentName")
	if deploymentName == "" {
		result.Failed(c, 400, "Deployment名称不能为空")
		return
	}

	ctrl.service.DeleteDeployment(c, uint(clusterId), namespaceName, deploymentName)
}

// ScaleDeployment 伸缩Deployment
// @Summary 伸缩Deployment
// @Description 调整Deployment的副本数
// @Tags Deployment管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param deploymentName path string true "Deployment名称"
// @Param scale body model.ScaleWorkloadRequest true "伸缩配置"
// @Success 200 {object} result.Result{data=string}
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 404 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/deployments/{deploymentName}/scale [post]
func (ctrl *K8sWorkloadController) ScaleDeployment(c *gin.Context) {
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

	deploymentName := c.Param("deploymentName")
	if deploymentName == "" {
		result.Failed(c, 400, "Deployment名称不能为空")
		return
	}

	var req model.ScaleWorkloadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Failed(c, 400, "参数验证失败: "+err.Error())
		return
	}

	ctrl.service.ScaleDeployment(c, uint(clusterId), namespaceName, deploymentName, &req)
}

// RestartDeployment 重启Deployment
// @Summary 重启Deployment
// @Description 通过更新Pod模板来重启Deployment
// @Tags Deployment管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param deploymentName path string true "Deployment名称"
// @Success 200 {object} result.Result{data=string}
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 404 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/deployments/{deploymentName}/restart [post]
func (ctrl *K8sWorkloadController) RestartDeployment(c *gin.Context) {
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

	deploymentName := c.Param("deploymentName")
	if deploymentName == "" {
		result.Failed(c, 400, "Deployment名称不能为空")
		return
	}

	ctrl.service.RestartDeployment(c, uint(clusterId), namespaceName, deploymentName)
}

// ===================== Pod 管理 =====================

// GetPods 获取Pod列表
// @Summary 获取Pod列表
// @Description 获取指定命名空间的Pod列表
// @Tags Pod管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Success 200 {object} result.Result{data=map[string]interface{}}
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/pods [get]
func (ctrl *K8sWorkloadController) GetPods(c *gin.Context) {
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

	ctrl.service.GetPods(c, uint(clusterId), namespaceName)
}

// GetPodDetail 获取Pod详情
// @Summary 获取Pod详情
// @Description 获取指定Pod的详细信息
// @Tags Pod管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param podName path string true "Pod名称"
// @Success 200 {object} result.Result{data=model.K8sPodDetail}
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 404 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/pods/{podName} [get]
func (ctrl *K8sWorkloadController) GetPodDetail(c *gin.Context) {
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

	podName := c.Param("podName")
	if podName == "" {
		result.Failed(c, 400, "Pod名称不能为空")
		return
	}

	ctrl.service.GetPodDetail(c, uint(clusterId), namespaceName, podName)
}

// DeletePod 删除Pod
// @Summary 删除Pod
// @Description 删除指定的Pod
// @Tags Pod管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param podName path string true "Pod名称"
// @Success 200 {object} result.Result{data=string}
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 404 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/pods/{podName} [delete]
func (ctrl *K8sWorkloadController) DeletePod(c *gin.Context) {
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

	podName := c.Param("podName")
	if podName == "" {
		result.Failed(c, 400, "Pod名称不能为空")
		return
	}

	ctrl.service.DeletePod(c, uint(clusterId), namespaceName, podName)
}

// GetPodLogs 获取Pod日志
// @Summary 获取Pod日志
// @Description 获取指定Pod容器的日志
// @Tags Pod管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param podName path string true "Pod名称"
// @Param container query string false "容器名称"
// @Success 200 {object} result.Result{data=map[string]interface{}}
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 404 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/pods/{podName}/logs [get]
func (ctrl *K8sWorkloadController) GetPodLogs(c *gin.Context) {
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

	podName := c.Param("podName")
	if podName == "" {
		result.Failed(c, 400, "Pod名称不能为空")
		return
	}

	containerName := c.Query("container")

	ctrl.service.GetPodLogs(c, uint(clusterId), namespaceName, podName, containerName)
}

// GetPodEvents 获取Pod事件
// @Summary 获取Pod事件
// @Description 获取指定Pod的相关事件列表
// @Tags Pod管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param podName path string true "Pod名称"
// @Success 200 {object} result.Result{data=map[string]interface{}}
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 404 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/pods/{podName}/events [get]
func (ctrl *K8sWorkloadController) GetPodEvents(c *gin.Context) {
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

	podName := c.Param("podName")
	if podName == "" {
		result.Failed(c, 400, "Pod名称不能为空")
		return
	}

	ctrl.service.GetPodEvents(c, uint(clusterId), namespaceName, podName)
}

// GetPodYaml 获取Pod的YAML配置
// @Summary 获取Pod的YAML配置
// @Description 获取指定Pod的完整YAML配置
// @Tags Pod管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param podName path string true "Pod名称"
// @Success 200 {object} result.Result{data=map[string]interface{}}
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 404 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/pods/{podName}/yaml [get]
func (ctrl *K8sWorkloadController) GetPodYaml(c *gin.Context) {
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

	podName := c.Param("podName")
	if podName == "" {
		result.Failed(c, 400, "Pod名称不能为空")
		return
	}

	ctrl.service.GetPodYaml(c, uint(clusterId), namespaceName, podName)
}

// UpdatePodYaml 更新Pod的YAML配置
// @Summary 更新Pod的YAML配置
// @Description 通过YAML内容更新指定的Pod配置，支持校验模式和DryRun模式
// @Tags Pod管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param podName path string true "Pod名称"
// @Param request body model.UpdatePodYAMLRequest true "更新请求"
// @Success 200 {object} result.Result{data=model.UpdatePodYAMLResponse}
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 404 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/pods/{podName}/yaml [put]
func (ctrl *K8sWorkloadController) UpdatePodYaml(c *gin.Context) {
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

	podName := c.Param("podName")
	if podName == "" {
		result.Failed(c, 400, "Pod名称不能为空")
		return
	}

	var req model.UpdatePodYAMLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Failed(c, 400, "请求参数错误: "+err.Error())
		return
	}

	ctrl.service.UpdatePodYaml(c, uint(clusterId), namespaceName, podName, &req)
}

// ===================== 通用工作负载YAML管理API =====================

// GetWorkloadYaml 获取工作负载的YAML配置
// @Summary 获取工作负载的YAML配置
// @Description 获取指定工作负载的完整YAML配置，支持deployment,statefulset,daemonset,job,cronjob
// @Tags 工作负载YAML管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param workloadType path string true "工作负载类型" Enums(deployment,statefulset,daemonset,job,cronjob)
// @Param workloadName path string true "工作负载名称"
// @Success 200 {object} result.Result{data=model.GetWorkloadYAMLResponse}
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 404 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/workloads/{workloadType}/{workloadName}/yaml [get]
func (ctrl *K8sWorkloadController) GetWorkloadYaml(c *gin.Context) {
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

	workloadType := c.Param("workloadType")
	if workloadType == "" {
		result.Failed(c, 400, "工作负载类型不能为空")
		return
	}

	workloadName := c.Param("workloadName")
	if workloadName == "" {
		result.Failed(c, 400, "工作负载名称不能为空")
		return
	}

	ctrl.service.GetWorkloadYaml(c, uint(clusterId), namespaceName, workloadType, workloadName)
}

// UpdateWorkloadYaml 更新工作负载的YAML配置
// @Summary 更新工作负载的YAML配置
// @Description 通过YAML内容更新指定的工作负载配置，支持deployment,statefulset,daemonset,job,cronjob
// @Tags 工作负载YAML管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param request body model.UpdateWorkloadYAMLRequest true "更新请求"
// @Success 200 {object} result.Result{data=model.UpdateWorkloadYAMLResponse}
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 404 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/workloads/yaml [put]
func (ctrl *K8sWorkloadController) UpdateWorkloadYaml(c *gin.Context) {
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

	var req model.UpdateWorkloadYAMLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Failed(c, 400, "请求参数错误: "+err.Error())
		return
	}

	ctrl.service.UpdateWorkloadYaml(c, uint(clusterId), namespaceName, &req)
}

// ===================== 监控API =====================

// GetPodMetrics 获取Pod监控指标
// @Summary 获取Pod监控指标
// @Description 获取指定Pod的CPU、内存等监控指标和使用率
// @Tags K8s监控
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param podName path string true "Pod名称"
// @Success 200 {object} result.Result{data=model.PodMetricsInfo}
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 404 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/pods/{podName}/metrics [get]
func (ctrl *K8sWorkloadController) GetPodMetrics(c *gin.Context) {
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

	podName := c.Param("podName")
	if podName == "" {
		result.Failed(c, 400, "Pod名称不能为空")
		return
	}

	ctrl.service.GetPodMetrics(c, uint(clusterId), namespaceName, podName)
}

// GetNodeMetrics 获取节点监控指标
// @Summary 获取节点监控指标
// @Description 获取指定节点的CPU、内存等监控指标和使用率
// @Tags K8s监控
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param nodeName path string true "节点名称"
// @Success 200 {object} result.Result{data=model.NodeMetricsInfo}
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 404 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/nodes/{nodeName}/metrics [get]
func (ctrl *K8sWorkloadController) GetNodeMetrics(c *gin.Context) {
	clusterIdStr := c.Param("id")
	clusterId, err := strconv.Atoi(clusterIdStr)
	if err != nil || clusterId <= 0 {
		result.Failed(c, 400, "无效的集群ID")
		return
	}

	nodeName := c.Param("nodeName")
	if nodeName == "" {
		result.Failed(c, 400, "节点名称不能为空")
		return
	}

	ctrl.service.GetNodeMetrics(c, uint(clusterId), nodeName)
}

// GetNamespaceMetrics 获取命名空间监控指标
// @Summary 获取命名空间监控指标
// @Description 获取指定命名空间下所有Pod的CPU、内存等监控指标汇总
// @Tags K8s监控
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Success 200 {object} result.Result{data=model.NamespaceMetricsInfo}
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 404 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/metrics [get]
func (ctrl *K8sWorkloadController) GetNamespaceMetrics(c *gin.Context) {
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

	ctrl.service.GetNamespaceMetrics(c, uint(clusterId), namespaceName)
}

// ===================== YAML管理API =====================

// CreatePodFromYAML 通过YAML创建Pod
// @Summary 通过YAML创建Pod
// @Description 通过提供的YAML内容创建Pod，支持校验模式和DryRun模式
// @Tags Pod管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param request body model.CreatePodFromYAMLRequest true "创建请求"
// @Success 200 {object} result.Result{data=model.CreatePodFromYAMLResponse}
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/pods/yaml [post]
func (ctrl *K8sWorkloadController) CreatePodFromYAML(c *gin.Context) {
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

	var req model.CreatePodFromYAMLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Failed(c, 400, "请求参数错误: "+err.Error())
		return
	}

	ctrl.service.CreatePodFromYAML(c, uint(clusterId), namespaceName, &req)
}

// ValidateYAML 校验YAML格式
// @Summary 校验YAML格式
// @Description 校验提供的YAML内容是否符合Kubernetes资源规范
// @Tags YAML校验
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param request body model.ValidateYAMLRequest true "校验请求"
// @Success 200 {object} result.Result{data=model.ValidateYAMLResponse}
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/yaml/validate [post]
func (ctrl *K8sWorkloadController) ValidateYAML(c *gin.Context) {
	clusterIdStr := c.Param("id")
	clusterId, err := strconv.Atoi(clusterIdStr)
	if err != nil || clusterId <= 0 {
		result.Failed(c, 400, "无效的集群ID")
		return
	}

	var req model.ValidateYAMLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Failed(c, 400, "请求参数错误: "+err.Error())
		return
	}

	ctrl.service.ValidateYAML(c, uint(clusterId), &req)
}

// GetWorkloadPods 获取工作负载下的Pod列表
// @Summary 获取工作负载下的Pod列表
// @Description 获取指定工作负载下的所有Pod信息
// @Tags Pod管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param type path string true "工作负载类型" Enums(deployment,statefulset,daemonset,job,cronjob)
// @Param workloadName path string true "工作负载名称"
// @Success 200 {object} result.Result{data=[]model.K8sPodInfo}
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/workloads/{type}/{workloadName}/pods [get]
func (ctrl *K8sWorkloadController) GetWorkloadPods(c *gin.Context) {
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

	workloadType := c.Param("type")
	if workloadType == "" {
		result.Failed(c, 400, "工作负载类型不能为空")
		return
	}

	workloadName := c.Param("workloadName")
	if workloadName == "" {
		result.Failed(c, 400, "工作负载名称不能为空")
		return
	}

	ctrl.service.GetWorkloadPods(c, uint(clusterId), namespaceName, workloadType, workloadName)
}

// ===================== Deployment 版本回滚管理 =====================

// GetDeploymentHistory 获取Deployment版本历史
// @Summary 获取Deployment版本历史
// @Description 获取指定Deployment的所有版本历史信息
// @Tags Deployment版本管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param deploymentName path string true "Deployment名称"
// @Success 200 {object} result.Result{data=model.DeploymentRolloutHistoryResponse}
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 404 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/deployments/{deploymentName}/history [get]
func (ctrl *K8sWorkloadController) GetDeploymentHistory(c *gin.Context) {
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

	deploymentName := c.Param("deploymentName")
	if deploymentName == "" {
		result.Failed(c, 400, "Deployment名称不能为空")
		return
	}

	ctrl.service.GetDeploymentHistory(c, uint(clusterId), namespaceName, deploymentName)
}

// GetDeploymentRevision 获取Deployment指定版本详情
// @Summary 获取Deployment指定版本详情
// @Description 获取指定Deployment特定版本的详细配置信息
// @Tags Deployment版本管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param deploymentName path string true "Deployment名称"
// @Param revision path int true "版本号"
// @Success 200 {object} result.Result{data=model.DeploymentRevisionDetail}
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 404 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/deployments/{deploymentName}/revisions/{revision} [get]
func (ctrl *K8sWorkloadController) GetDeploymentRevision(c *gin.Context) {
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

	deploymentName := c.Param("deploymentName")
	if deploymentName == "" {
		result.Failed(c, 400, "Deployment名称不能为空")
		return
	}

	revisionStr := c.Param("revision")
	revision, err := strconv.ParseInt(revisionStr, 10, 64)
	if err != nil || revision <= 0 {
		result.Failed(c, 400, "无效的版本号")
		return
	}

	ctrl.service.GetDeploymentRevision(c, uint(clusterId), namespaceName, deploymentName, revision)
}

// RollbackDeployment 回滚Deployment到指定版本
// @Summary 回滚Deployment到指定版本
// @Description 将Deployment回滚到指定的历史版本
// @Tags Deployment版本管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param deploymentName path string true "Deployment名称"
// @Param request body model.RollbackDeploymentRequest true "回滚请求"
// @Success 200 {object} result.Result{data=model.RollbackDeploymentResponse}
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 404 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/deployments/{deploymentName}/rollback [post]
func (ctrl *K8sWorkloadController) RollbackDeployment(c *gin.Context) {
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

	deploymentName := c.Param("deploymentName")
	if deploymentName == "" {
		result.Failed(c, 400, "Deployment名称不能为空")
		return
	}

	var req model.RollbackDeploymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Failed(c, 400, "参数验证失败: "+err.Error())
		return
	}

	ctrl.service.RollbackDeployment(c, uint(clusterId), namespaceName, deploymentName, &req)
}

// PauseDeployment 暂停Deployment滚动更新
// @Summary 暂停Deployment滚动更新
// @Description 暂停正在进行的Deployment滚动更新过程
// @Tags Deployment版本管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param deploymentName path string true "Deployment名称"
// @Success 200 {object} result.Result{data=model.PauseDeploymentResponse}
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 404 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/deployments/{deploymentName}/pause [post]
func (ctrl *K8sWorkloadController) PauseDeployment(c *gin.Context) {
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

	deploymentName := c.Param("deploymentName")
	if deploymentName == "" {
		result.Failed(c, 400, "Deployment名称不能为空")
		return
	}

	ctrl.service.PauseDeployment(c, uint(clusterId), namespaceName, deploymentName)
}

// ResumeDeployment 恢复Deployment滚动更新
// @Summary 恢复Deployment滚动更新
// @Description 恢复被暂停的Deployment滚动更新过程
// @Tags Deployment版本管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param deploymentName path string true "Deployment名称"
// @Success 200 {object} result.Result{data=model.ResumeDeploymentResponse}
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 404 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/deployments/{deploymentName}/resume [post]
func (ctrl *K8sWorkloadController) ResumeDeployment(c *gin.Context) {
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

	deploymentName := c.Param("deploymentName")
	if deploymentName == "" {
		result.Failed(c, 400, "Deployment名称不能为空")
		return
	}

	ctrl.service.ResumeDeployment(c, uint(clusterId), namespaceName, deploymentName)
}

// GetDeploymentRolloutStatus 获取Deployment滚动发布状态
// @Summary 获取Deployment滚动发布状态
// @Description 获取指定Deployment的当前滚动发布状态和进度信息
// @Tags Deployment版本管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户Token"
// @Param id path int true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param deploymentName path string true "Deployment名称"
// @Success 200 {object} result.Result{data=model.DeploymentRolloutStatusResponse}
// @Failure 400 {object} result.Result
// @Failure 401 {object} result.Result
// @Failure 404 {object} result.Result
// @Failure 500 {object} result.Result
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/deployments/{deploymentName}/rollout-status [get]
func (ctrl *K8sWorkloadController) GetDeploymentRolloutStatus(c *gin.Context) {
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

	deploymentName := c.Param("deploymentName")
	if deploymentName == "" {
		result.Failed(c, 400, "Deployment名称不能为空")
		return
	}

	ctrl.service.GetDeploymentRolloutStatus(c, uint(clusterId), namespaceName, deploymentName)
}