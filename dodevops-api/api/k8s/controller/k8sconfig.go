package controller

import (
	"dodevops-api/api/k8s/model"
	"dodevops-api/api/k8s/service"
	"dodevops-api/common/result"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type K8sConfigController struct {
	DB *gorm.DB
}

func NewK8sConfigController(db *gorm.DB) *K8sConfigController {
	return &K8sConfigController{DB: db}
}

// ===================== ConfigMap 管理接口 =====================

// GetConfigMaps 获取ConfigMap列表
// @Summary 获取ConfigMap列表
// @Description 获取指定集群和命名空间下的ConfigMap列表
// @Tags K8s配置管理
// @Accept json
// @Produce json
// @Param id path int true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页大小" default(10)
// @Success 200 {object} result.Result{data=model.ConfigMapListResponse}
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/configmaps [get]
func (ctrl *K8sConfigController) GetConfigMaps(c *gin.Context) {
	clusterIdStr := c.Param("id")
	clusterId, err := strconv.Atoi(clusterIdStr)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的集群ID")
		return
	}

	namespaceName := c.Param("namespaceName")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	configService := service.NewK8sConfigService(ctrl.DB)
	response, err := configService.GetConfigMaps(clusterId, namespaceName, page, pageSize)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, err.Error())
		return
	}

	result.Success(c, response)
}

// GetConfigMapDetail 获取ConfigMap详情
// @Summary 获取ConfigMap详情
// @Description 获取指定ConfigMap的详细信息
// @Tags K8s配置管理
// @Accept json
// @Produce json
// @Param id path int true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param configMapName path string true "ConfigMap名称"
// @Success 200 {object} result.Result{data=model.ConfigMapDetail}
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/configmaps/{configMapName} [get]
func (ctrl *K8sConfigController) GetConfigMapDetail(c *gin.Context) {
	clusterIdStr := c.Param("id")
	clusterId, err := strconv.Atoi(clusterIdStr)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的集群ID")
		return
	}

	namespaceName := c.Param("namespaceName")
	configMapName := c.Param("configMapName")

	configService := service.NewK8sConfigService(ctrl.DB)
	response, err := configService.GetConfigMapDetail(clusterId, namespaceName, configMapName)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, err.Error())
		return
	}

	result.Success(c, response)
}

// CreateConfigMap 创建ConfigMap
// @Summary 创建ConfigMap
// @Description 创建新的ConfigMap
// @Tags K8s配置管理
// @Accept json
// @Produce json
// @Param id path int true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param request body model.CreateConfigMapRequest true "创建ConfigMap请求"
// @Success 200 {object} result.Result{data=model.K8sConfigMap}
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/configmaps [post]
func (ctrl *K8sConfigController) CreateConfigMap(c *gin.Context) {
	clusterIdStr := c.Param("id")
	clusterId, err := strconv.Atoi(clusterIdStr)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的集群ID")
		return
	}

	namespaceName := c.Param("namespaceName")

	var req model.CreateConfigMapRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Failed(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	configService := service.NewK8sConfigService(ctrl.DB)
	response, err := configService.CreateConfigMap(clusterId, namespaceName, &req)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, err.Error())
		return
	}

	result.Success(c, response)
}

// UpdateConfigMap 更新ConfigMap
// @Summary 更新ConfigMap
// @Description 更新指定的ConfigMap
// @Tags K8s配置管理
// @Accept json
// @Produce json
// @Param id path int true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param configMapName path string true "ConfigMap名称"
// @Param request body model.UpdateConfigMapRequest true "更新ConfigMap请求"
// @Success 200 {object} result.Result{data=model.K8sConfigMap}
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/configmaps/{configMapName} [put]
func (ctrl *K8sConfigController) UpdateConfigMap(c *gin.Context) {
	clusterIdStr := c.Param("id")
	clusterId, err := strconv.Atoi(clusterIdStr)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的集群ID")
		return
	}

	namespaceName := c.Param("namespaceName")
	configMapName := c.Param("configMapName")

	var req model.UpdateConfigMapRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Failed(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	configService := service.NewK8sConfigService(ctrl.DB)
	response, err := configService.UpdateConfigMap(clusterId, namespaceName, configMapName, &req)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, err.Error())
		return
	}

	result.Success(c, response)
}

// DeleteConfigMap 删除ConfigMap
// @Summary 删除ConfigMap
// @Description 删除指定的ConfigMap
// @Tags K8s配置管理
// @Accept json
// @Produce json
// @Param id path int true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param configMapName path string true "ConfigMap名称"
// @Success 200 {object} result.Result
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/configmaps/{configMapName} [delete]
func (ctrl *K8sConfigController) DeleteConfigMap(c *gin.Context) {
	clusterIdStr := c.Param("id")
	clusterId, err := strconv.Atoi(clusterIdStr)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的集群ID")
		return
	}

	namespaceName := c.Param("namespaceName")
	configMapName := c.Param("configMapName")

	configService := service.NewK8sConfigService(ctrl.DB)
	err = configService.DeleteConfigMap(clusterId, namespaceName, configMapName)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, err.Error())
		return
	}

	result.Success(c, "ConfigMap删除成功")
}

// GetConfigMapYaml 获取ConfigMap YAML
// @Summary 获取ConfigMap YAML
// @Description 获取指定ConfigMap的YAML配置
// @Tags K8s配置管理
// @Accept json
// @Produce json
// @Param id path int true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param configMapName path string true "ConfigMap名称"
// @Success 200 {object} result.Result{data=string}
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/configmaps/{configMapName}/yaml [get]
func (ctrl *K8sConfigController) GetConfigMapYaml(c *gin.Context) {
	clusterIdStr := c.Param("id")
	clusterId, err := strconv.Atoi(clusterIdStr)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的集群ID")
		return
	}

	namespaceName := c.Param("namespaceName")
	configMapName := c.Param("configMapName")

	configService := service.NewK8sConfigService(ctrl.DB)
	yamlContent, err := configService.GetConfigMapYaml(clusterId, namespaceName, configMapName)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, err.Error())
		return
	}

	result.Success(c, yamlContent)
}

// UpdateConfigMapYaml 更新ConfigMap YAML
// @Summary 更新ConfigMap YAML
// @Description 通过YAML更新ConfigMap配置
// @Tags K8s配置管理
// @Accept json
// @Produce json
// @Param id path int true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param configMapName path string true "ConfigMap名称"
// @Param request body map[string]interface{} true "YAML内容"
// @Success 200 {object} result.Result{data=model.K8sConfigMap}
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/configmaps/{configMapName}/yaml [put]
func (ctrl *K8sConfigController) UpdateConfigMapYaml(c *gin.Context) {
	clusterIdStr := c.Param("id")
	clusterId, err := strconv.Atoi(clusterIdStr)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的集群ID")
		return
	}

	namespaceName := c.Param("namespaceName")
	configMapName := c.Param("configMapName")

	var yamlData map[string]interface{}
	if err := c.ShouldBindJSON(&yamlData); err != nil {
		result.Failed(c, http.StatusBadRequest, "YAML内容格式错误: "+err.Error())
		return
	}

	configService := service.NewK8sConfigService(ctrl.DB)
	response, err := configService.UpdateConfigMapYaml(clusterId, namespaceName, configMapName, yamlData)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, err.Error())
		return
	}

	result.Success(c, response)
}

// ===================== Secret 管理接口 =====================

// GetSecrets 获取Secret列表
// @Summary 获取Secret列表
// @Description 获取指定集群和命名空间下的Secret列表
// @Tags K8s配置管理
// @Accept json
// @Produce json
// @Param id path int true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页大小" default(10)
// @Success 200 {object} result.Result{data=model.SecretListResponse}
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/secrets [get]
func (ctrl *K8sConfigController) GetSecrets(c *gin.Context) {
	clusterIdStr := c.Param("id")
	clusterId, err := strconv.Atoi(clusterIdStr)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的集群ID")
		return
	}

	namespaceName := c.Param("namespaceName")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	configService := service.NewK8sConfigService(ctrl.DB)
	response, err := configService.GetSecrets(clusterId, namespaceName, page, pageSize)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, err.Error())
		return
	}

	result.Success(c, response)
}

// GetSecretDetail 获取Secret详情
// @Summary 获取Secret详情
// @Description 获取指定Secret的详细信息
// @Tags K8s配置管理
// @Accept json
// @Produce json
// @Param id path int true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param secretName path string true "Secret名称"
// @Success 200 {object} result.Result{data=model.SecretDetail}
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/secrets/{secretName} [get]
func (ctrl *K8sConfigController) GetSecretDetail(c *gin.Context) {
	clusterIdStr := c.Param("id")
	clusterId, err := strconv.Atoi(clusterIdStr)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的集群ID")
		return
	}

	namespaceName := c.Param("namespaceName")
	secretName := c.Param("secretName")

	configService := service.NewK8sConfigService(ctrl.DB)
	response, err := configService.GetSecretDetail(clusterId, namespaceName, secretName)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, err.Error())
		return
	}

	result.Success(c, response)
}

// CreateSecret 创建Secret
// @Summary 创建Secret
// @Description 创建新的Secret
// @Tags K8s配置管理
// @Accept json
// @Produce json
// @Param id path int true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param request body model.CreateSecretRequest true "创建Secret请求"
// @Success 200 {object} result.Result{data=model.K8sSecret}
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/secrets [post]
func (ctrl *K8sConfigController) CreateSecret(c *gin.Context) {
	clusterIdStr := c.Param("id")
	clusterId, err := strconv.Atoi(clusterIdStr)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的集群ID")
		return
	}

	namespaceName := c.Param("namespaceName")

	var req model.CreateSecretRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Failed(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	configService := service.NewK8sConfigService(ctrl.DB)
	response, err := configService.CreateSecret(clusterId, namespaceName, &req)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, err.Error())
		return
	}

	result.Success(c, response)
}

// UpdateSecret 更新Secret
// @Summary 更新Secret
// @Description 更新指定的Secret
// @Tags K8s配置管理
// @Accept json
// @Produce json
// @Param id path int true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param secretName path string true "Secret名称"
// @Param request body model.UpdateSecretRequest true "更新Secret请求"
// @Success 200 {object} result.Result{data=model.K8sSecret}
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/secrets/{secretName} [put]
func (ctrl *K8sConfigController) UpdateSecret(c *gin.Context) {
	clusterIdStr := c.Param("id")
	clusterId, err := strconv.Atoi(clusterIdStr)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的集群ID")
		return
	}

	namespaceName := c.Param("namespaceName")
	secretName := c.Param("secretName")

	var req model.UpdateSecretRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Failed(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	configService := service.NewK8sConfigService(ctrl.DB)
	response, err := configService.UpdateSecret(clusterId, namespaceName, secretName, &req)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, err.Error())
		return
	}

	result.Success(c, response)
}

// DeleteSecret 删除Secret
// @Summary 删除Secret
// @Description 删除指定的Secret
// @Tags K8s配置管理
// @Accept json
// @Produce json
// @Param id path int true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param secretName path string true "Secret名称"
// @Success 200 {object} result.Result
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/secrets/{secretName} [delete]
func (ctrl *K8sConfigController) DeleteSecret(c *gin.Context) {
	clusterIdStr := c.Param("id")
	clusterId, err := strconv.Atoi(clusterIdStr)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的集群ID")
		return
	}

	namespaceName := c.Param("namespaceName")
	secretName := c.Param("secretName")

	configService := service.NewK8sConfigService(ctrl.DB)
	err = configService.DeleteSecret(clusterId, namespaceName, secretName)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, err.Error())
		return
	}

	result.Success(c, "Secret删除成功")
}

// GetSecretYaml 获取Secret YAML
// @Summary 获取Secret YAML
// @Description 获取指定Secret的YAML配置
// @Tags K8s配置管理
// @Accept json
// @Produce json
// @Param id path int true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param secretName path string true "Secret名称"
// @Success 200 {object} result.Result{data=string}
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/secrets/{secretName}/yaml [get]
func (ctrl *K8sConfigController) GetSecretYaml(c *gin.Context) {
	clusterIdStr := c.Param("id")
	clusterId, err := strconv.Atoi(clusterIdStr)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的集群ID")
		return
	}

	namespaceName := c.Param("namespaceName")
	secretName := c.Param("secretName")

	configService := service.NewK8sConfigService(ctrl.DB)
	yamlContent, err := configService.GetSecretYaml(clusterId, namespaceName, secretName)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, err.Error())
		return
	}

	result.Success(c, yamlContent)
}

// UpdateSecretYaml 更新Secret YAML
// @Summary 更新Secret YAML
// @Description 通过YAML更新Secret配置
// @Tags K8s配置管理
// @Accept json
// @Produce json
// @Param id path int true "集群ID"
// @Param namespaceName path string true "命名空间名称"
// @Param secretName path string true "Secret名称"
// @Param request body map[string]interface{} true "YAML内容"
// @Success 200 {object} result.Result{data=model.K8sSecret}
// @Router /k8s/cluster/{id}/namespaces/{namespaceName}/secrets/{secretName}/yaml [put]
func (ctrl *K8sConfigController) UpdateSecretYaml(c *gin.Context) {
	clusterIdStr := c.Param("id")
	clusterId, err := strconv.Atoi(clusterIdStr)
	if err != nil {
		result.Failed(c, http.StatusBadRequest, "无效的集群ID")
		return
	}

	namespaceName := c.Param("namespaceName")
	secretName := c.Param("secretName")

	var yamlData map[string]interface{}
	if err := c.ShouldBindJSON(&yamlData); err != nil {
		result.Failed(c, http.StatusBadRequest, "YAML内容格式错误: "+err.Error())
		return
	}

	configService := service.NewK8sConfigService(ctrl.DB)
	response, err := configService.UpdateSecretYaml(clusterId, namespaceName, secretName, yamlData)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, err.Error())
		return
	}

	result.Success(c, response)
}